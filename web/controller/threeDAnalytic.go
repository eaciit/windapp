package controller

type ThreeDAnalyticController struct {
	App
}

func CreateThreeDAnalyticController() *ThreeDAnalyticController {
	var controller = new(ThreeDAnalyticController)
	return controller
}
