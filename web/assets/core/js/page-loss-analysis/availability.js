'use strict';

viewModel.Availability = new Object();
var av = viewModel.Availability;

av.breakDown = ko.observableArray([]);
av.typeChart = ko.observable();
av.breakDownVal = ko.observable();
av.dataSource = ko.observableArray();
var height = $(".content").width() * 0.125;

av.Availability = function(){
    var valid = fa.LoadData();
    if (valid) {
        if(pg.isFirstAvailability() === true){
            app.loading(true);
            av.breakDownVal = $("#breakdownlistavail").data("kendoDropDownList").value();
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine,
                project: fa.project,
                breakDown: av.breakDownVal,
            };
            toolkit.ajaxPost(viewModel.appName + "analyticavailability/getdata", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                setTimeout(function(){
                    av.dataSource(res.data);
                    av.createChartAvailability(av.dataSource());
                    av.createChartProduction(av.dataSource());
                    pg.isFirstAvailability(false);
                    app.loading(false);
                },200);
            });
            $('#availabledatestart').html(pg.availabledatestartscada3());
            $('#availabledateend').html(pg.availabledateendscada3());
        }else{
            setTimeout(function(){
                $('#availabledatestart').html(pg.availabledatestartscada3());
                $('#availabledateend').html(pg.availabledateendscada3());
                $("#availabilityChart").data("kendoChart").refresh();
                $("#productionChart").data("kendoChart").refresh();
                // app.loading(false);
            },200);
        }
    }
}

av.SetBreakDown = function () {
    fa.disableRefreshButton(true);
    av.breakDown = [];

    setTimeout(function () {
        var project = $('#projectList').data("kendoDropDownList").value();
        fa.populateTurbine(project);
        
        $.each(fa.GetBreakDown(), function (i, val) {
            if (val.value == "Turbine" || val.value == "Project") {
                return false;
            } else {
               av.breakDown.push(val);
            }
        });

        $("#breakdownlistavail").data("kendoDropDownList").dataSource.data(av.breakDown);
        $("#breakdownlistavail").data("kendoDropDownList").dataSource.query();
        $("#breakdownlistavail").data("kendoDropDownList").select(0);


        fa.disableRefreshButton(false);
    }, 500);
}

av.createChartAvailability = function (dataSource) {
    var series = dataSource.SeriesAvail;
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#eb5b19";

    $("#availabilityChart").replaceWith('<div id="availabilityChart"></div>');
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
        chartArea :{
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209 + 100) ) / 2,
            height: 200,
            padding: 0,
            background: "transparent",
        },
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
                text: $("#breakdownlistavail").data("kendoDropDownList").value(),
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
    $("#availabilityChart").css("top",-60);
    $("#availabilityChart").css("margin-bottom",-25);
}
av.createChartProduction = function (dataSource) {
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#ff9933";

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
        chartArea :{
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209 + 100) ) / 2,
            height: 200, 
            margin : 0,
            padding: 0,
            background: "transparent",
        },
        series: seriesProd,
        seriesColors: colorField,
        valueAxes: [{
            max: max,
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
                text: $("#breakdownlistavail").data("kendoDropDownList").value(),
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
     $("#productionChart").css("top",-60);
     $("#productionChart").css("margin-bottom",-60);
}