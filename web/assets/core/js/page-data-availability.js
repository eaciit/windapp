viewModel.DataAvailability = new Object();
var page = viewModel.DataAvailability;

vm.currentMenu('Data Availability');
vm.currentTitle('Data Availability');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'KPI Table', href: viewModel.appName + 'page/dataavailability' }]);

page.dataDummy = ko.observableArray([
		{
			Category : "Scada Data", 
			Turbine	 : [],
			Data 	 : [
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "20%"
				},
				{
					"tooltip" : "11 jan 2017 until 1 feb",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "20%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "15%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "15%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "5%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},

			]
		},
		{
			Category : "Event Down", 
			Turbine	 : [],
			Data 	 : [
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},
				{
					"tooltip" : "11 jan 2017 until 1 feb",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "20%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "2%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "5%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "1%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "10%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},

			]
		}	
	])
page.hideFilter = function(){
    $("#periodList").closest(".k-widget").hide();
    $("#dateStart").closest(".k-widget").hide();
    $("#dateEnd").closest(".k-widget").hide();
    $(".control-label:contains('Period')").hide();
    $(".control-label:contains('to')").hide();
}

$(function () {
	page.hideFilter();
    setTimeout(function() {
        app.loading(false);    
    }, 200);
});
