package main

import (
	"flag"
	tk "github.com/eaciit/toolkit"
	"strings"

	c3to10 "eaciit/wfdemo-git-dev/processapp/threeextractor/convert3secto10min/lib"
)

var (
	intstartdate = int(20160801)
	intenddate   = int(20160831)
	intyear      = int(2016)

	strarraydate = ""
	strfilename  = ""
	config       = map[string]string{}
)

func main() {

	flag.IntVar(&intstartdate, "sdate", 20160821, "Start date for processing data")
	flag.IntVar(&intenddate, "edate", 20160831, "End date for processing data")
	flag.StringVar(&strarraydate, "adate", "", "Interval date in string")
	flag.StringVar(&strfilename, "file", "", "Filename will be first priority to execute")
	flag.IntVar(&intyear, "year", 0, "Full year for processing data")
	flag.Parse()

	param := tk.M{}
	if strfilename != "" {
		param.
			Set("selector", "file").
			Set("file", strfilename)
	} else if strarraydate != "" {
		arridate := strings.Split(strarraydate, "|")
		param.
			Set("selector", "adate").
			Set("adate", arridate)
	} else if intyear > 0 {
		param.
			Set("selector", "year").
			Set("year", intyear)
	} else {
		param.
			Set("selector", "date").
			Set("sdate", tk.ToString(intstartdate)).
			Set("edate", tk.ToString(intenddate))
	}

	err := c3to10.Generate(param)
	if err != nil {
		tk.Println(err)
	} else {
		tk.Println(">> DONE <<")
	}

}

//>> Param >>
// selector = date | file | adate | year
// >> selector = date
// sdate / edate exp. "20160821" / "20160831"
// exp tk.M{}.Set("selector", "date").Set("sdate", "20160821").Set("edate","20160831")
// >> selector = file
// file exp. "DataFile20160821-01.csv"
// exp tk.M{}.Set("selector", "file").Set("file", "DataFile20160821-01.csv")
// >> selector = adate
// adate exp []string{"20160821", "20160831"}
// exp tk.M{}.Set("selector", "adate").Set("adate", []string{"20160821", "20160831"})
// >> selector = year
// year exp 2016
// exp tk.M{}.Set("selector", "year").Set("year", 2016)
//>>>>>>>>>>>
