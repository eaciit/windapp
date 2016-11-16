'use strict';

viewModel.WRFlexiDetail = new Object();
var wd = viewModel.WRFlexiDetail;

wd.dataWindrose = ko.observableArray([]);
wd.dataWindroseGrid = ko.observableArray([]);
wd.dataWindroseEachTurbine = ko.observableArray([]);
wd.sectorDerajat = ko.observable(0);

wd.sectionsBreakdownList = ko.observableArray([
	{ "text": 36, "value": 36 },
	{ "text": 24, "value": 24 },
	{ "text": 12, "value": 12 },
]);
var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
var listOfChart = [];
var listOfButton = {};
var listOfCategory = [
	{ "category": "0 to 4m/s", "color": colorFieldsWR[0] },
	{ "category": "4 to 8m/s", "color": colorFieldsWR[1] },
	{ "category": "8 to 12m/s", "color": colorFieldsWR[2] },
	{ "category": "12 to 16m/s", "color": colorFieldsWR[3] },
	{ "category": "16 to 20m/s", "color": colorFieldsWR[4] },
	{ "category": "20m/s and above", "color": colorFieldsWR[5] },
];

wd.ExportWindRose = function () {
	var chart = $("#wr-chart").getKendoChart();
	chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
		kendo.saveAs({
			dataURI: data,
			fileName: "WindRose.pdf",
		});
	});
}
var maxValue = 0;

wd.GetData = function () {
	app.loading(true);
	fa.LoadData();

	setTimeout(function () {
		fa.getProjectInfo();
		var breakDownVal = $("#nosection").data("kendoDropDownList").value();
		var secDer = 360 / breakDownVal;
		wd.sectorDerajat(secDer);
		var param = {
			period: fa.period,
			dateStart: fa.dateStart,
			dateEnd: fa.dateEnd,
			turbine: fa.turbine,
			project: fa.project,
			breakDown: breakDownVal,
		};
		toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getflexidataeachturbine", param, function (res) {
			if (!app.isFine(res)) {
				app.loading(false);
				return;
			}
			if (res.data.WindRose != null) {
				var metData = res.data.WindRose;
				maxValue = res.data.MaxValue;
				wd.dataWindroseEachTurbine(metData);
				wd.initChart();
			}

			app.loading(false)

		})
	}, 300);
}


wd.initChart = function () {
	app.loading(true)
	listOfChart = [];
	var breakDownVal = $("#nosection").data("kendoDropDownList").value();
	var stepNum = 1
	var gapNum = 1
	if (breakDownVal == 36) {
		stepNum = 3
		gapNum = 0
	} else if (breakDownVal == 24) {
		stepNum = 2
		gapNum = 0
	} else if (breakDownVal == 12) {
		stepNum = 1
		gapNum = 0
	}

	$.each(wd.dataWindroseEachTurbine(), function (i, val) {
		var name = val.Name
		if (name == "MetTower") {
			name = "Met Tower"
		}

		var idChart = "#chart-" + val.Name
		listOfChart.push(idChart);
		var pWidth = $('body').width() * ($(idChart).closest('div.windrose-item').width() - 2) / 100;
		$(idChart).kendoChart({
			theme: "nova",
			chartArea: {
				width: pWidth,
				height: pWidth
			},

			title: {
				text: name,
				font: '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
			},
			legend: {
				position: "bottom",
				labels: {
					template: "#= (series.data[0] || {}).WsCategoryDesc #"
				},
				visible: false,
			},
			dataSource: {
				data: val.Data,
				group: {
					field: "WsCategoryNo",
					dir: "asc"
				},
				sort: {
					field: "DirectionNo",
					dir: "asc"
				}
			},
			seriesColors: colorFieldsWR,
			series: [{
				type: "radarColumn",
				stack: true,
				field: "Contribution",
				gap: gapNum,
				border: {
					width: 1,
					color: "#7f7f7f",
					opacity: 0.5
				},
			}],
			categoryAxis: {
				field: "DirectionDesc",
				visible: true,
				majorGridLines: {
					visible: true,
					step: stepNum
				},
				labels: {
					font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
					visible: true,
					step: stepNum
				}
			},
			valueAxis: {
				labels: {
					template: kendo.template("#= kendo.toString(value, 'n0') #%"),
					font: '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
				},
				majorUnit: 10,
				max: maxValue,
				min: 0
			},
			tooltip: {
				visible: true,
				template: "#= category # (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
				background: "rgb(255,255,255, 0.9)",
				color: "#58666e",
				font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
				border: {
					color: "#eee",
					width: "2px",
				},
			}
		});
		app.loading(true)
		setTimeout(function () {
			if ($(idChart).data("kendoChart") != null) {
				$(idChart).data("kendoChart").refresh();
			}
		}, 200);
	});
}

wd.showHideLegend = function (index) {
	var idName = "btn" + index;
	listOfButton[idName] = !listOfButton[idName];
	if (listOfButton[idName] == false) {
		$("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
	} else {
		$("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
	}
	$.each(listOfChart, function (idx, idChart) {
		$(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
		$(idChart).data("kendoChart").refresh();
	});
}

wd.setPeriod = function(){		
	var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), 6, 23, 0, 0, 0, 0));
    var lastEndDate = new Date(Date.UTC(moment(maxDateData).get('year'), 6+1, 0, 0, 0, 0, 0));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);

	$("#periodList").change(function () {
		var period = $('#periodList').data('kendoDropDownList').value();
		
		if(period == "monthly")  {
			setTimeout(function(){
				$('#dateStart').data('kendoDatePicker').value(lastStartDate);
	    		$('#dateEnd').data('kendoDatePicker').value(lastEndDate);
			}, 60);
		}else if(period == "custom") {
			setTimeout(function(){
				$('#dateStart').data('kendoDatePicker').value(lastStartDate);
	    		$('#dateEnd').data('kendoDatePicker').value(lastEndDate);
			}, 60);
		}
	});
}

wd.checkPeriod = function(){
	var period = $('#periodList').data('kendoDropDownList').value();
    var monthNames = moment.months();

    var currentDateData = moment(new Date(2016, 6 + 1, 0)).format("YYYY-MM-DD");
    var today = moment().format('YYYY-MM-DD');
    var thisMonth = moment().get('month');
    var firstDayMonth = moment(new Date(2016, 6, 1)).format("YYYY-MM-DD");
    var lastDayMonth = moment(new Date(2016, 6 + 1, 0)).format("YYYY-MM-DD");
    var firstDayYear = moment().startOf('year').format('YYYY-MM-DD');
    var endDayYear = moment().endOf('year').format('YYYY-MM-DD');

    var dateStart = moment(fa.dateStart).format('YYYY-MM-DD');
    var dateEnd = moment(fa.dateEnd).format('YYYY-MM-DD');

    if (period === 'custom') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period range on start date and end date");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodRange("* Incomplete period range on start date");
            fa.infoPeriodIconmozilla(true);
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period range on end date");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'annual') {
        if ((moment(fa.dateEnd).get('year') == moment(new Date(2016, 6 + 1, 0)).get('year')) && (currentDateData < today)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period range in end year");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'monthly') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period range in start month and start month");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period range in start month");
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period range in end month");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else {
        fa.infoPeriodRange("");
        fa.infoPeriodIcon(false);
    }
}

vm.currentMenu('Wind Rose');
vm.currentTitle('Wind Rose');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind Rose', href: viewModel.appName + 'page/analyticwindroseflexi' }]);

$(document).ready(function () {
	$('#btnRefresh').on('click', function () {
		app.loading(true);
		setTimeout(function () {
			wd.GetData();
			wd.checkPeriod();
		}, 200);
	});

	// $('#btnRefresh').click();
	app.loading(true);
	setTimeout(function () {
		$("#legend-list").html("");
		$.each(listOfCategory, function (idx, val) {
			var idName = "btn" + idx;
			listOfButton[idName] = true;
			$("#legend-list").append(
				'<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="wd.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
				'<span class="span-legend">' + val.category + '</span>'
			);
		});
		$("#nosection").data("kendoDropDownList").value(12);
		// wd.GetData();
		wd.setPeriod();
		$( "#btnRefresh" ).trigger( "click" );
	}, 300);

	

});

$(document).bind("kendo:skinChange", wd.GetData);

