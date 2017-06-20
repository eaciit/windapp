package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type App struct {
	Server *knot.Server
	Reff   toolkit.M
}

var (
	LayoutFile    = "layout.html"
	AppBasePath   = func(dir string, err error) string { return dir }(os.Getwd())
	DATA_PATH     = filepath.Join(AppBasePath, "data")
	CONFIG_PATH   = filepath.Join(AppBasePath, "config")
	MenuList      = []string{}
	ServerAddress = ""
)

func init() {
	fmt.Println("Base Path ===> ", AppBasePath)

	if DATA_PATH != "" {
		fmt.Println("DATA_PATH ===> ", DATA_PATH)
		fmt.Println("CONFIG_PATH ===> ", CONFIG_PATH)
	}
}

type Payload struct {
	Period              string
	DateStart           time.Time
	DateEnd             time.Time
	Project             string
	Turbine             string
	Skip                int
	Take                int
	Sort                []Sorting
	Filter              toolkit.M
	IsValidTimeDuration bool
}

type PayloadAnalytic struct {
	Period     string
	Project    string
	Turbine    []interface{}
	DateStart  time.Time
	DateEnd    time.Time
	IsClean    bool
	IsAverage  bool
	IsDownTime bool
	Color      []interface{}
	BreakDown  string
}

type PayloadAnalyticTLP struct {
	Period          string
	Project         string
	Turbine         []interface{}
	DateStart       time.Time
	DateEnd         time.Time
	ColName         string
	DeviationStatus bool
	Deviation       float64
}

type PayloadAnalyticPC struct {
	Period       string
	DateStart    time.Time
	DateEnd      time.Time
	Turbine      []interface{}
	Project      string
	IsClean      bool
	IsDeviation  bool
	DeviationVal string
	IsAverage    bool
	Color        []interface{}
	ColorDeg     []interface{}
	IsDownTime   bool
	BreakDown    string
	ViewSession  string
}

type PayloadKPI struct {
	Project         string
	Turbine         []interface{}
	Period          string
	DateStart       time.Time
	DateEnd         time.Time
	ColumnBreakDown string
	RowBreakDown    string
	KeyA            string
	KeyB            string
	KeyC            string
}

type PayloadPCComparison struct {
	PC1Period    string
	PC1Project   string
	PC1Turbine   string //[]interface{}
	PC1DateStart time.Time
	PC1DateEnd   time.Time

	PC2Period    string
	PC2Project   string
	PC2Turbine   string //[]interface{}
	PC2DateStart time.Time
	PC2DateEnd   time.Time
}

type PayloadTimeSeries struct {
	Period    string
	Project   string
	Turbine   string
	DateStart time.Time
	DateEnd   time.Time
	TagList   []string
	DataType  string
	PageType  string
	IsHour    bool
}
