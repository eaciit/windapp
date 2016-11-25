package model
import(
	"strconv"
	"strings"
)
type ProcessedLog struct{
	Filename 	string
	StartTime 	string
	EndTime		string
	RowIndex	int
}
func (p *ProcessedLog) ToString()string{
	output := ""
	output+=p.Filename+";"
	output+=p.StartTime+";"
	output+=p.EndTime+";"
	output+=strconv.Itoa(p.RowIndex)+";"
	return output
}
func FromString(line string)*ProcessedLog{
	newPF := new(ProcessedLog)
	sval := strings.Split(string(line), ";")
	newPF.Filename = sval[0]
	newPF.StartTime = sval[1]
	newPF.EndTime = sval[2]
	newPF.RowIndex,_ = strconv.Atoi(sval[3])
	return newPF
}
