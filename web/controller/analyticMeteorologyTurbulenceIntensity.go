package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/helper"
	// . "eaciit/wfdemo-git/library/models"
	// "eaciit/wfdemo-git/web/helper"
	// "sort"
	// "strings"
	// "time"

	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
)

func (m *AnalyticMeteorologyController) GetTurbulenceIntensity(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return nil
}
