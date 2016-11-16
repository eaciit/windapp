package datareader

import (
	"bufio"
	. "eaciit/wfdemo/library/models"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/eaciit/orm"
	. "github.com/eaciit/sshclient"
	tk "github.com/eaciit/toolkit"
)

type DataReader struct {
	FileLocation string
	PathProcess  string
	PathRoot     string
}

var (
	DataTranspose    tk.M
	mutex            = &sync.Mutex{}
	mutexData        = &sync.Mutex{}
	idx              = 0
	FileCount        = 0
	DraftDir         = "Draft"
	ProcessDir       = "Process"
	SuccessDir       = "Success"
	ReaderConfigFile = "conf/reader.conf"
)

func NewDataReader(fileLocation string, pathProcess string, pathRoot string) *DataReader {
	dr := new(DataReader)
	dr.FileLocation = fileLocation
	dr.PathProcess = pathProcess
	dr.PathRoot = pathRoot

	return dr
}

func fileExists(fileLocation string) bool {
	if _, err := os.Stat(fileLocation); err == nil {
		return true
	}

	return false
}

func (c *DataReader) readerConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(c.PathRoot + ReaderConfigFile)
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}

func (c *DataReader) writeConfig(lastFileName string, lastIndex int) {
	input, err := ioutil.ReadFile(c.PathRoot + ReaderConfigFile)
	if err != nil {
		tk.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "LastFileName") {
			lines[i] = "LastFileName=" + lastFileName
		}
		if strings.Contains(line, "LastIndex") {
			lines[i] = tk.Sprintf("LastIndex=%v", lastIndex)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(c.PathRoot+ReaderConfigFile, []byte(output), 0644)
	if err != nil {
		tk.Println(err)
	}
}

func (c *DataReader) Start() {
	fileToProcess := c.copyFile(c.FileLocation, c.PathProcess+"\\"+DraftDir)
	if fileToProcess != "" {
		if fileExists(fileToProcess) {
			file, _ := os.Stat(fileToProcess)
			start := time.Now()

			DataTranspose = tk.M{}
			FileCount++
			c.readFile(file.Name())

			duration := time.Now().Sub(start).Seconds()
			tk.Println(tk.Sprintf("Loading file %v data about %v sec(s)", file.Name(), duration))
		}
	}
}

func (c *DataReader) readFile(fileName string) {
	var wg sync.WaitGroup

	conf := c.readerConfig()
	lastFileName := conf["LastFileName"]
	lastIndex := tk.ToInt(conf["LastIndex"], "0")

	if lastFileName != "" && lastFileName != fileName {
		c.writeConfig("", 0)

		conf = c.readerConfig()
		lastFileName = conf["LastFileName"]
		lastIndex = tk.ToInt(conf["LastIndex"], "0")
	}

	tk.Println("Start processing file: " + fileName)

	draftFile := c.PathProcess + "\\" + DraftDir + "\\" + fileName
	processFile := c.PathProcess + "\\" + ProcessDir + "\\" + fileName
	successFile := c.PathProcess + "\\" + SuccessDir + "\\" + fileName

	err := os.Rename(draftFile, processFile)
	if err != nil {
		tk.Println("Error Move Draft File : ", err.Error())
	}

	f, _ := os.Open(processFile)
	r, err := csv.NewReader(bufio.NewReader(f)).ReadAll()
	if err != nil {
		tk.Println("Error Read File : ", err.Error())
	}
	countData := len(r)

	countPerProcess := 5000
	counter := 0
	startIndex := counter * countPerProcess
	if lastIndex > 0 {
		startIndex = lastIndex
	}
	endIndex := (counter+1)*countPerProcess - 1
	isFinish := false

	if startIndex >= countData {
		isFinish = true
		c.writeConfig("", 0)
	}

	for !isFinish {
		startIndex = counter * countPerProcess
		endIndex = (counter+1)*countPerProcess - 1

		if endIndex > countData {
			endIndex = countData
		}

		data := r[startIndex:endIndex]

		wg.Add(1)
		go func(data [][]string) {
			for _, d := range data {
				mutex.Lock()

				contents := d
				parseContent(contents)

				mutex.Unlock()
			}
			wg.Done()
		}(data)

		counter++

		if endIndex >= countData {
			isFinish = true
		}
	}

	f.Close()

	wg.Wait()

	err = os.Rename(processFile, successFile)
	if err != nil {
		tk.Println("Error Move Process File : ", err.Error())
	}

	if len(DataTranspose) > 0 {
		tk.Println("Start create log file...")
		fName := c.createLog()
		tk.Println("Finish created log file")
		fZip := c.createZip(fName)
		tk.Println("Start sending file: " + fZip)
		c.sendFile(fZip)
		c.writeConfig(fileName, endIndex)
	}
}

func (c *DataReader) sendFile(filename string) {
	locationTarget := "/home/developer/wfdemo-watch/draft"

	ssh := new(SshSetting)

	ssh.SSHAuthType = SSHAuthType_Certificate
	ssh.SSHHost = "go.eaciit.com:22"
	ssh.SSHUser = "developer"
	ssh.SSHKeyLocation = c.PathRoot + "\\conf\\key\\developer.pem"

	_, err := ssh.Connect()
	if err != nil {
		tk.Println("Error connecting to server: " + err.Error())
	} else {
		tk.Println("Connected to server!")
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		tk.Println("Error opening file: " + err.Error())
		os.Exit(1)
	}

	fileStat, err := file.Stat()
	if err != nil {
		tk.Println("Error opening file: " + err.Error())
		os.Exit(1)
	}

	err = ssh.SshCopyByFile(file, fileStat.Size(), fileStat.Mode().Perm(), filepath.Base(fileStat.Name()), locationTarget)
	if err != nil {
		tk.Println("Error: ", err.Error())
	} else {
		tk.Println("Sending file successfully")
	}
}

func (c *DataReader) copyFile(src string, pathTarget string) string {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		tk.Println("Error read source file: " + err.Error())
		os.Exit(1)
	}

	srcFileStat, _ := srcFile.Stat()
	destFile, err := os.Create(pathTarget + "\\" + srcFileStat.Name())
	defer destFile.Close()
	if err != nil {
		tk.Println("Error read target file: " + err.Error())
		os.Exit(1)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		tk.Println("Error copy file: " + err.Error())
		os.Exit(1)
	}

	err = destFile.Sync()
	if err != nil {
		tk.Println("Error sync file: " + err.Error())
		os.Exit(1)
	}

	return destFile.Name()
}

func (c *DataReader) createLog() string {
	scada := new(ScadaThreeSecs)
	ref := reflect.ValueOf(scada).Elem()
	typeOf := ref.Type()

	content := ""
	delim := ""
	for i := 0; i < ref.NumField(); i++ {
		if typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" && typeOf.Field(i).Name != "ModelBase" {
			content += delim + typeOf.Field(i).Name
			delim = ","
		}
	}

	f, _ := os.Create(c.PathProcess + "\\Results\\result_" + tk.ToString(FileCount) + ".csv")
	defer f.Close()

	f.WriteString(content + "\n")
	for _, value := range DataTranspose {
		content = ""
		delim = ""
		for i := 0; i < ref.NumField(); i++ {
			if typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" && typeOf.Field(i).Name != "ModelBase" {
				field := typeOf.Field(i).Name
				fieldType := ref.Field(i).Type()

				if field == "ID" {
					field = "Id"
				}

				var valData interface{}
				if value.(tk.M).Has(field) {
					if fieldType.String() == "time.Time" {
						valTime := value.(tk.M).Get(field).(time.Time)
						valData = valTime.Format("2006-01-02 15:04:05")
					} else {
						valData = value.(tk.M).Get(field)
					}
				} else {
					valData = ""
				}
				content += tk.Sprintf("%v%v", delim, valData)
				delim = ","
			}
		}
		f.WriteString(content + "\n")
	}

	return f.Name()
}

func (c *DataReader) createZip(fileName string) string {
	filetarget := c.PathProcess + "\\Results\\result_" + tk.ToString(FileCount) + ".zip"
	err := tk.ZipCompress(fileName, filetarget)
	if err != nil {
		tk.Println("Error compressing file: ", err.Error())
	}

	return filetarget
}

func parseContent(contents []string) {
	time1, _ := time.Parse("02-Jan-2006 15:04:05", contents[0])
	time2, _ := time.Parse("02-Jan-2006 15:04:05", contents[1])
	date1, _ := time.Parse("2006-01-02", time1.Format("2006-01-02"))
	date2, _ := time.Parse("2006-01-02", time2.Format("2006-01-02"))

	thour := time1.Hour()
	tminute := time1.Minute()
	tsecond := time1.Second()
	tminutevalue := float64(tminute) + tk.Div(float64(tsecond), 60.0)
	tminutecategory := tk.ToInt(tk.RoundingUp64(tk.Div(tminutevalue, 10), 0)*10, "0")
	if tminutecategory == 60 {
		tminutecategory = 0
		thour = thour + 1
	}
	newTimeStamp := date1.Add(time.Duration(thour) * time.Hour).Add(time.Duration(tminutecategory) * time.Minute)
	timestampconverted := newTimeStamp.UTC()

	infos := strings.Split(contents[2], ".")
	value := tk.ToFloat64(contents[3], 6, tk.RoundingAuto)

	project := "Tejuva"
	turbine := infos[2]
	column := infos[3] + "_" + infos[4]

	id := time1.Format("20060102_150405") + "_" + time2.Format("20060102_150405") + "_" + project + "_" + turbine

	if DataTranspose.Get(id) == nil {
		DataTranspose.Set(id, tk.M{}.Set("Id", id).Set("ProjectName", project).Set("Turbine", turbine).Set("TimeStamp1", time1).Set("TimeStamp2", time2).Set("DateId1", date1).Set("DateId2", date2).Set("THour", thour).Set("TMinute", tminute).Set("TSecond", tsecond).Set("TMinuteValue", tminutevalue).Set("TMinuteCategory", tminutecategory).Set("TimeStampConverted", timestampconverted).Set(column, value))
	} else {
		newData := DataTranspose.Get(id).(tk.M)
		DataTranspose.Set(id, newData.Set(column, value))
	}
}
