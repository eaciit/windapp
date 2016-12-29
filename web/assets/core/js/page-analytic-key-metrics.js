'use strict';

viewModel.KeyMetrics = new Object();
var km = viewModel.KeyMetrics;

km.breakdown = ko.observableArray([
    { "value": "$dateinfo.dateid", "text": "Date" },
    { "value": "$dateinfo.monthid", "text": "Month" },
    { "value": "$dateinfo.year", "text": "Year" },
    { "value": "$turbine", "text": "Turbine" },
    { "value": "$projectname", "text": "Project" }
]);

km.ExportKeyMetrics = function () {
	var chart = $("#km-chart").getKendoChart();
	chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
		kendo.saveAs({
			dataURI: data,
			fileName: "KeyMetrics.pdf",
		});
	});
}

km.createChart = function (dataSource) {
	var key1Val = $("#key1").data("kendoDropDownList").value();
	var key2Val = $("#key2").data("kendoDropDownList").value();
	var breakdownVal = $("#breakdownlist").data("kendoDropDownList").value();

	var filters = [
		{ field: "dateinfo.dateid", operator: "gte", value: fa.dateStart },
		{ field: "dateinfo.dateid", operator: "lte", value: fa.dateEnd },
		{ field: "turbine", operator: "in", value: fa.turbine },
	];
	/*var listOfMonths = [];
	var monthCount = fa.dateStart.getMonth();
	for (monthCount = fa.dateStart.getMonth(); monthCount <= fa.dateEnd.getMonth(); monthCount++) {
    	listOfMonths.push(monthCount);
	}*/
	if (fa.project != "") {
		filters.push({ field: "projectname", operator: "eq", value: fa.project })
	}

	var filter = { filters: filters }
	var misc = {
		key1: key1Val,
		key2: key2Val,
		breakdown: breakdownVal,
		duration: ((fa.dateEnd - fa.dateStart) / 86400000) + 1,
		totalturbine: fa.turbine.length,
		period: fa.period,
		/*monthlist: listOfMonths,
		startdate: fa.dateStart.getDate(),
		enddate: fa.dateEnd.getDate(),
		year: fa.dateEnd.getYear(),*/
	};
	var param = { filter: filter, misc: misc };

	var request = toolkit.ajaxPost(viewModel.appName + "analytickeymetrics/getkeymetrics", param, function (res) {
		if (!app.isFine(res)) {
			app.loading(false);
			return;
		}

		var series = res.data.Series;
		var categories = res.data.Categories;
		var minKey1 = res.data.MinKey1;
		var maxKey1 = res.data.MaxKey1;
		var minKey2 = res.data.MinKey2;
		var maxKey2 = res.data.MaxKey2;
		var catTitle = res.data.CatTitle;
		var rotation = 0;
		if (breakdownVal == "$turbine") {
			rotation = -330;
		}

		$("#km-chart").kendoChart({
			theme: "flat",
			title: {
				text: ""
			},
			legend: {
				position: "top",
				visible: true,
			},
			chartArea: {
				height : 370,
				// width : 900
			},
			series: series,
			seriesColor: colorField,
			valueAxes: [{
				name: "Key1",
				title: {
					text: $("#key1").data("kendoDropDownList").text() + " (" + $("#key1").data("kendoDropDownList").dataSource.data()[$("#key1").data("kendoDropDownList").select()].unit + ")",
					visible: true,
					font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
				},
				labels: {
					step: 2
				},
				line: {
					visible: false
				},
				axisCrossingValue: -10,
				majorGridLines: {
					visible: true,
					color: "#eee",
					width: 0.8,
				},
	            /*min: 0,
	            max: maxBar*2-(maxBar/4),*/
				min: minKey1,
				max: maxKey1,
			},
			{
				name: "Key2",
				title: {
					text: $("#key2").data("kendoDropDownList").text() + " (" + $("#key2").data("kendoDropDownList").dataSource.data()[$("#key2").data("kendoDropDownList").select()].unit + ")",
					visible: true,
					font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
				},
				visible: true,
	            /*min: -(maxLine*4),
	            max: maxLine*2,*/
				min: minKey2,
				max: maxKey2,
			}],
			categoryAxis: {
				categories: categories,
				majorGridLines: {
					visible: false
				},
				title: {
					text: catTitle,
					font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
				},
				labels: {
					rotation: rotation
				},
				majorTickType: "none",
				axisCrossingValues: [0, 1000],
			},
			tooltip: {
				visible: true,
				format: "{0:n1}",
				background: "rgb(255,255,255, 0.9)",
				shared: true,
				sharedTemplate: kendo.template($("#template").html()),
				color: "#58666e",
				font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
				border: {
					color: "#eee",
					width: "2px",
				},
			}
		});
	});

	$.when(request).done(function(){
		setTimeout(function(){
			app.loading(false);
		},300);
	});

}

vm.currentMenu('Compare Metrics');
vm.currentTitle('Compare Metrics');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Compare Metrics', href: viewModel.appName + 'page/analytickeymetrics' }]);

km.setBreakDown = function () {
	km.breakdownList = [];

	setTimeout(function () {
		/*$.each(fa.GetBreakDown(), function (i, val) {
			km.breakdownList.push(val);
		});*/

		$.each(km.breakdown(), function (i, valx) {
			$.each(fa.GetBreakDown(), function (i, valy) {
				if (valx.text == valy.text) {
					km.breakdownList.push(valx);
				}
			});
		});

		$("#breakdownlist").data("kendoDropDownList").dataSource.data(km.breakdownList);
		$("#breakdownlist").data("kendoDropDownList").dataSource.query();
		if ($("#breakdownlist").data("kendoDropDownList").value() == "") {
			$("#breakdownlist").data("kendoDropDownList").select(0);
		}
	}, 1000);
}

km.getData = function () {
	app.loading(true);
	var request = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
		if (!app.isFine(res)) {
            return;
        }
        var minDatetemp = new Date(res.data.ScadaData[0]);
        var maxDatetemp = new Date(res.data.ScadaData[1]);
        $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
    });
	fa.getProjectInfo();
	fa.LoadData();

	setTimeout(function () {
		km.setBreakDown();
		km.createChart();
	}, 1000);
}

$(document).ready(function () {
	$('#btnRefresh').on('click', function () {
		setTimeout(function () {
			km.getData();
		}, 200);
	});

	// smart filter :)

	$('#periodList').kendoDropDownList({
		data: fa.periodList,
		dataValueField: 'value',
		dataTextField: 'text',
		suggest: true,
		change: function () { fa.showHidePeriod(km.setBreakDown()) }
	});

	setTimeout(function () {
		$('#projectList').kendoDropDownList({
			data: fa.projectList,
			dataValueField: 'value',
			dataTextField: 'text',
			suggest: true,
			change: function () { km.setBreakDown() }
		});

		$("#dateStart").change(function () { fa.DateChange(km.setBreakDown()) });
		$("#dateEnd").change(function () { fa.DateChange(km.setBreakDown()) });

		km.getData();
	}, 1500);
});

$(document).bind("kendo:skinChange", km.createChart);