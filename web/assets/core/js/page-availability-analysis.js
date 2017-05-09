'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;

pg.typeChart = ko.observable();
pg.breakDownVal = ko.observable();
pg.dataSource = ko.observableArray();
var height = $(".content").width() * 0.125;

pg.ChartAvailability = function () {
    app.loading(true);

    pg.breakDownVal = $("#breakdownlist").data("kendoDropDownList").value();

    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        breakDown: pg.breakDownVal,
    };

    toolkit.ajaxPost(viewModel.appName + "analyticavailability/getdata", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        pg.dataSource(res.data);
        pg.createChartAvailability(pg.dataSource());
        pg.createChartProduction(pg.dataSource());
        app.loading(false);

        if ($("#availabilityChart").data("kendoChart") != null) {
            $("#availabilityChart").data("kendoChart").refresh();
        }
        if ($("#productionChart").data("kendoChart") != null) {
            $("#productionChart").data("kendoChart").refresh();
        }

    });
};

pg.createChartAvailability = function (dataSource) {
    var series = dataSource.SeriesAvail;
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#944dff";

    $("#availabilityChart").replaceWith('<div id="availabilityChart"></div>');
    $("#availabilityChart").height(height);
    $("#availabilityChart").kendoChart({
        theme: "Flat",
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        series: series,
        seriesColors: colorField,
        valueAxes: [{
            line: {
                visible: false
            },
            max: 100,
            min: 0,
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availpercentage",
            title: { text: "Availability (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }, {
            visible: false,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availline",
            title: { text: "Production (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            categories: categories,
            title: {
                text: $("#breakdownlist").data("kendoDropDownList").value(),
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValues: [0, 1],
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            shared: true,
            sharedTemplate: kendo.template($("#template").html()),
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
}
pg.createChartProduction = function (dataSource) {
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#ff880e";

    $("#productionChart").replaceWith('<div id="productionChart"></div>');
    $("#productionChart").height(height);
    $("#productionChart").kendoChart({
        height: "150px",
        theme: "Flat",
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            background: "transparent",
        },
        series: seriesProd,
        seriesColors: colorField,
        valueAxes: [{
            visible: true,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availline",
            title: { text: "Production (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            categories: categories,
            title: {
                text: $("#breakdownlist").data("kendoDropDownList").value(),
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: "#= series.name # at #= category # : #= kendo.toString(value, 'n2')# MWh",
            shared: false,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });
}
pg.SetBreakDown = function () {
    pg.columnsBreakdownList = [];
    pg.rowsBreakdownList = [];

    setTimeout(function () {
        var project = $('#projectList').data("kendoDropDownList").value();
        fa.populateTurbine(project);

        if (true) {
            pg.rowsBreakdownList.push({ "value": "Project", "text": "Project" });
        }

        $.each(fa.GetBreakDown(), function (i, val) {
            if (val.value == "Turbine") {
                pg.rowsBreakdownList.push(val);
            } else if (val.value != "Project") {
                pg.columnsBreakdownList.push(val);
            }
        });

        $("#breakdownlist").data("kendoDropDownList").dataSource.data(fa.GetBreakDown());
        $("#breakdownlist").data("kendoDropDownList").dataSource.query();
        if ($("#breakdownlist").data("kendoDropDownList").value() == "") {
            $("#breakdownlist").data("kendoDropDownList").select(0);
        }
    }, 1000);
}

pg.loadData = function () {
    app.loading(true);
    fa.getProjectInfo();
    setTimeout(function () {
        fa.LoadData();
        pg.SetBreakDown();
        pg.ChartAvailability();
    }, 1000);
}

vm.currentMenu('Availability');
vm.currentTitle('Availability');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Availability', href: viewModel.appName + 'page/analyticcomparison' }]);

$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        pg.loadData();
    });
    $('#breakdownlist').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.loadData() },
    });

    // smart filter :)

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { fa.showHidePeriod(pg.SetBreakDown()) }
    });

    setTimeout(function () {
        $('#projectList').kendoDropDownList({
            data: fa.projectList,
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { pg.SetBreakDown() }
        });

        $("#dateStart").change(function () { fa.DateChange(pg.SetBreakDown()) });
        $("#dateEnd").change(function () { fa.DateChange(pg.SetBreakDown()) });

        pg.loadData();
    }, 1500);
});