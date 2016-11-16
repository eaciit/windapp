'use strict';

viewModel.WindRose = new Object();
var wr = viewModel.WindRose;

wr.dataWindrose = ko.observableArray([]);
wr.dataWindroseGrid = ko.observableArray([]);
wr.dataWindroseEachTurbine = ko.observableArray([]);

wr.sectionsBreakdownList = ko.observableArray([
    {"text" : 8, "value" : 8},
    {"text" : 10, "value" :10},
    {"text" : 12, "value" : 12},
    {"text" : 15, "value" : 15},
    {"text" : 18, "value" : 18},
    {"text" : 20, "value" : 20},
    {"text" : 24, "value" : 24},
    {"text" : 36, "value" : 36},
]); 

wr.ExportWindRose = function(){
    var chart = $("#wr-chart").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function(data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "WindRose.pdf",
        });
    });
}

wr.ExportWindRoseFreq = function(){
    var chart = $("#wr-chart-freq").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function(data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "WindRose.pdf",
        });
    });
}


wr.GetData = function() {
	app.loading(true);
	fa.LoadData();
	var breakDownVal = $("#nosection").data("kendoDropDownList").value();
	var param = {
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        breakDown: breakDownVal,
    };
	toolkit.ajaxPost(viewModel.appName + "analyticwindroseflexi/getflexidataeachturbine", param, function (res) {
		if (!app.isFine(res)) {
			app.loading(false);
			return;
		}
		console.log(res.data)
        var metData = res.data.WindRose;
        wr.dataWindroseEachTurbine(metData.reverse());

        wr.initChart();
    
        app.loading(false)

	})
}


wr.initChart = function() {
	 app.loading(true)
	$.each(wr.dataWindroseEachTurbine(), function(i, val){
		var name = val.Name
		if(name =="All"){
			name = "Met Tower"
		}

   		var idChart = "#chart-" + val.Name
		$(idChart).kendoChart({
			theme : "nova",
		    title: {
		        text: name,
		        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		    },
		    legend: {
		        position: "bottom",
		        labels: {
		            template: "#= (series.data[0] || {}).WsCategoryDesc #"
		        },
		    },
		    dataSource: {
		    	data : val.Data,
		        group: {
		            field: "WsCategoryNo",
		            dir: "asc"
		        },
		        sort: {
		            field: "DirectionNo",
		            dir: "asc"
		        }
		    },
		    seriesColors: colorField,
		    series: [{
		        type: "radarColumn",
		        stack: true,
		        field: "Contribution"
		    }],
		    categoryAxis: {
		        field: "DirectionDesc"
		    },
		    valueAxis: {
		        visible: false
		    },
			tooltip: {
		        visible: true,
		        template: "#= category # (#= dataItem.WsCategoryDesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
		        background: "rgb(255,255,255, 0.9)",
		        color : "#58666e",
		        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		        border : {
		            color : "#eee",
		            width : "2px",
		        },
		    }
		});
		app.loading(true)
		setTimeout(function() {
            if ($(idChart).data("kendoChart") != null) {
                $(idChart).data("kendoChart").refresh();
            }
        }, 10);
        
    });
}

wr.GetDataOld = function() {
	app.loading(true);
	fa.LoadData();
	var breakDownVal = $("#nosection").data("kendoDropDownList").value();
	var param = {
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        breakDown: breakDownVal,
    };
	toolkit.ajaxPost(viewModel.appName + "analyticwindroseflexi/getflexidata", param, function (res) {
		if (!app.isFine(res)) {
			app.loading(false);
			return;
		}
		wr.dataWindrose(res.data.WindRose);
		wr.dataWindroseGrid(res.data.GridWindrose);

		wr.initGridContribution()
		wr.createChartSpeed()
		wr.createChartFreq()

		app.loading(false);

	})

	// wr.dataWindrose(res.data.WindRose);
}

wr.initGridContribution = function() {
	app.loading(true);
	$('#gridContribution').html("");
	$('#gridContribution').kendoGrid({
	  theme:"flat",
      dataSource: {
      	data : wr.dataWindroseGrid(),
        sort: [
            { field: 'WsCategoryNo', dir: 'asc' },
        ],
      },
      sortable: true,
      columns: [
        { title: "Category", field: "WsCategoryDesc", width: 70, headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:left;"} },
        { title: "Hours", field: "Hours", width: 50, format: "{0:n2}", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:right;"} },
		{ title: "Times", field: "Frequency", width: 50, format: "{0:n0}", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:right;"} },
      ]
    });
    app.loading(false);
	setTimeout(function() {
        if ($("#gridContribution").data("kendoGrid") != null) {
            $("#gridContribution").data("kendoGrid").refresh();
        }
    }, 10);
}

wr.createChartSpeed = function() {
	app.loading(true);
	$("#wr-chart").kendoChart({
		theme : "nova",
	    title: {
	        text: "Speed Distribution",
	        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	    },
	    legend: {
	        position: "bottom",
	        labels: {
	            template: "#= (series.data[0] || {}).WsCategoryDesc #"
	        },
	    },
	    dataSource: {
	        data : wr.dataWindrose(),
	        group: {
	            field: "WsCategoryNo",
	            dir: "asc"
	        },
	        sort: {
	            field: "DirectionNo",
	            dir: "asc"
	        },
	    },
		seriesColors: colorField,
	    series: [{
	        type: "radarColumn",
	        stack: true,
	        field: "Contribution"
	    }],
	    categoryAxis: {
	        field: "DirectionDesc",
	        labels : {
	            rotation: 0
	        }
	    },
	    valueAxis: {
	        visible: false
	    },
		tooltip: {
	        visible: true,
	        template: "#= category # (#= dataItem.WsCategoryDesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # hours",
	        background: "rgb(255,255,255, 0.9)",
	        color : "#58666e",
	        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	        border : {
	            color : "#eee",
	            width : "2px",
	        },

	    }
	});
	app.loading(false);

	setTimeout(function() {
        if ($("#wr-chart").data("kendoChart") != null) {
            $("#wr-chart").data("kendoChart").refresh();
        }
    }, 10);
}

wr.createChartFreq = function() {
	app.loading(true);
	$("#wr-chart-freq").kendoChart({
		theme : "nova",
	    title: {
	        text: "Frequency Distribution",
	        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	    },
	    legend: {
	        position: "bottom",
	        labels: {
	           template: "#= (series.data[0] || {}).WsCategoryDesc #"
	        },
	    },
	    dataSource: {
	        data : wr.dataWindrose(),
	        group: {
	            field: "WsCategoryNo",
	            dir: "asc"
	        },
	        sort: {
	            field: "DirectionNo",
	            dir: "asc"
	        },
	    },
	    seriesColors: ["#ff9ea5", "#db7bd8", "#7b9adb", "#b5e61d"],
	    series: [{
	        type: "radarColumn",
	        stack: true,
	        field: "Contribution"
	    }],
	    categoryAxis: {
	        field: "DirectionDesc",
	        labels : {
	            rotation: 0
	        }
	    },
	    valueAxis: {
	        visible: false
	    },
		tooltip: {
	        visible: true,
	        template: "#= category # (#= dataItem.WsCategoryDesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.Frequency, 'n0') # times",
	        background: "rgb(255,255,255, 0.9)",
	        color : "#58666e",
	        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	        border : {
	            color : "#eee",
	            width : "2px",
	        },

	    }
	});
	app.loading(false);
	setTimeout(function() {
        if ($("#wr-chart-freq").data("kendoChart") != null) {
            $("#wr-chart-freq").data("kendoChart").refresh();
        }
    }, 10);
}

vm.currentMenu('Wind Rose');
vm.currentTitle('Wind Rose');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind Rose', href: viewModel.appName + 'page/analyticwindroseflexi' }]);

$(document).ready(function(){
	$('#btnRefresh').on('click', function(){
		app.loading(true);
		setTimeout(function(){ 
			wr.GetData();
		},200);
	});
	
	$('#btnRefresh').click();
});

$(document).bind("kendo:skinChange", wr.GetData);

