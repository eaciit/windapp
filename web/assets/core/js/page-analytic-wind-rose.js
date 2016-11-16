'use strict';

viewModel.WindRose = new Object();
var wr = viewModel.WindRose;
 

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
	fa.LoadData();
	var filters = [
		{ field: "dateinfo.dateid", operator: "gte", value: fa.dateStart },
		{ field: "dateinfo.dateid", operator: "lte", value: fa.dateEnd },
		{ field: "turbineid", operator: "in", value: fa.turbine },
		{ field: "projectid", operator: "eq", value: fa.project },
	];
	var filter = {filters : filters}
	var param = {filter : filter};

	$('#gridContribution').html("");
	$('#gridContribution').kendoGrid({
	  theme:"flat",
      dataSource: {
      	transport: {
      		read: {
                url: viewModel.appName + "analyticwindrose/getwscategory",
            	type: "POST",
            	data: param,
            	dataType: "json",
            	contentType: "application/json; charset=utf-8"
            },
            parameterMap: function(options) {
            	return JSON.stringify(options);
          	}
        },
        schema: {
        	data: function(res){
          		return res.data.WSCategory
          	}
        },
      },
      columns: [
        { title: "Category", field: "wscategorydesc", width: 70, headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:left;"} },
        { title: "Hours", field: "hours", width: 50, format: "{0:n2}", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:right;"} },
		{ title: "Times", field: "frequency", width: 50, format: "{0:n0}", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:right;"} },
      ]
    });
}

wr.createChart = function() {
	var filters = [
		{ field: "dateinfo.dateid", operator: "gte", value: fa.dateStart },
		{ field: "dateinfo.dateid", operator: "lte", value: fa.dateEnd },
		{ field: "turbineid", operator: "in", value: fa.turbine },
		{ field: "projectid", operator: "eq", value: fa.project },
	];
	var filter = {filters : filters}

	var param = {filter : filter};
	$("#wr-chart").kendoChart({
		theme : "nova",
	    title: {
	        // text: "Contribution each Speed Category by Direction",
	        text: "Speed Distribution",
	        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	    },
	    legend: {
	        position: "bottom",
	        labels: {
	            template: "#= (series.data[0] || {}).wscategorydesc #"
	        },
	        // position: "custom",
     		/*offsetX: 700,
     		offsetY: 150*/
	    },
	    dataSource: {
	        transport: {
	            read: {
	                url: viewModel.appName + "analyticwindrose/getwsdata",
	            	type: "POST",
	            	data: param,
	            	dataType: "json",
	            	contentType: "application/json; charset=utf-8"
	            },
	            parameterMap: function(options) {
	            	return JSON.stringify(options);
	          	}
	        },
	        group: {
	            field: "wscategoryno",
	            dir: "asc"
	        },
	        sort: {
	            field: "directionno",
	            dir: "asc"
	        },
	        schema: {
	        	data: function(res){
	          		return res.data.WindRose
	          	}
	        },
	    },
		seriesColors: colorField,
	    series: [{
	        type: "radarColumn",
	        stack: true,
	        field: "contribute"
	    }],
	    categoryAxis: {
	        field: "directiondesc"
	    },
	    valueAxis: {
	        visible: false
	    },
		tooltip: {
	        visible: true,
	        template: "#= category # (#= dataItem.wscategorydesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.hours, 'n2') # hours",
	        background: "rgb(255,255,255, 0.9)",
	        color : "#58666e",
	        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	        border : {
	            color : "#eee",
	            width : "2px",
	        },

	    }
	});
}

wr.createChartFreq = function() {
	var filters = [
		{ field: "dateinfo.dateid", operator: "gte", value: fa.dateStart },
		{ field: "dateinfo.dateid", operator: "lte", value: fa.dateEnd },
		{ field: "turbineid", operator: "in", value: fa.turbine },
		{ field: "projectid", operator: "eq", value: fa.project },
	];
	var filter = {filters : filters}

	var param = {filter : filter};
	$("#wr-chart-freq").kendoChart({
		theme : "nova",
	    title: {
	        // text: "Frequency each Speed Category by Direction",
	        text: "Frequency Distribution",
	        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	    },
	    legend: {
	        position: "bottom",
	        labels: {
	            template: "#= (series.data[0] || {}).wscategorydesc #"
	        },
	    },
	    dataSource: {
	        transport: {
	            read: {
	                url: viewModel.appName + "analyticwindrose/getwsdata",
	            	type: "POST",
	            	data: param,
	            	dataType: "json",
	            	contentType: "application/json; charset=utf-8"
	            },
	            parameterMap: function(options) {
	            	return JSON.stringify(options);
	          	}
	        },
	        group: {
	            field: "wscategoryno",
	            dir: "asc"
	        },
	        sort: {
	            field: "directionno",
	            dir: "asc"
	        },
	        schema: {
	        	data: function(res){
	          		return res.data.WindRose
	          	}
	        },
	    },
	    seriesColors: ["#ff9ea5", "#db7bd8", "#7b9adb", "#b5e61d"],
	    // seriesColors: ["#ff4350", "#ff9ea5", "#ffbf46", "#80deea", "#00acc1"],
		// seriesColors: colorField,
	    series: [{
	        type: "radarColumn",
	        stack: true,
	        field: "contribute"
	    }],
	    categoryAxis: {
	        field: "directiondesc"
	    },
	    valueAxis: {
	        visible: false
	    },
		tooltip: {
	        visible: true,
	        template: "#= category # (#= dataItem.wscategorydesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.frequency, 'n0') # times",
	        background: "rgb(255,255,255, 0.9)",
	        color : "#58666e",
	        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	        border : {
	            color : "#eee",
	            width : "2px",
	        },

	    }
	});
}

vm.currentMenu('Wind Rose');
vm.currentTitle('Wind Rose');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind Rose', href: viewModel.appName + 'page/analyticwindrose' }]);

$(document).ready(function(){
	$('#btnRefresh').on('click', function(){
		app.loading(true);
		setTimeout(function(){ 
			wr.GetData();
			wr.createChart();
			wr.createChartFreq();
			app.loading(false);
		},200);
	});
	
	$('#btnRefresh').click();
});

$(document).bind("kendo:skinChange", wr.createChart);
$(document).bind("kendo:skinChange", wr.createChartFreq);