package helper

import (
	"log"
	"regexp"
	"strconv"

	tk "github.com/eaciit/toolkit"
)

func LogProcess(processName string, totalData float64, totalDone float64) {
	donePCT := tk.ToFloat64(totalDone/totalData*100, 2, tk.RoundingAuto)
	strPCT := ""

	if int64(donePCT)%10 == 0 {
		for count := 1; count <= 10; count++ {
			if donePCT > (float64(count) * 10) {
				strPCT = strPCT + "#"
			} else {
				strPCT = strPCT + "_"
			}
		}

		if totalDone >= totalData {
			log.Printf(">> %v  [%v] %v/%v (%v PCT) \n", processName, "##########", formatCommas(int(totalDone)), formatCommas(int(totalData)), donePCT)
			log.Printf("==================================c DONE: %v \n", processName)
		} else {
			log.Printf(">> %v  [%v] %v/%v (%v PCT) \n", processName, strPCT, formatCommas(int(totalDone)), formatCommas(int(totalData)), donePCT)
		}
	}
}

func formatCommas(num int) string {
	numString := strconv.Itoa(num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for {
		formatted := re.ReplaceAllString(numString, "$1,$2")
		if formatted == numString {
			return formatted
		}
		numString = formatted
	}
}
