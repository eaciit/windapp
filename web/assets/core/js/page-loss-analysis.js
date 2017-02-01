'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.type = ko.observable();
pg.detailDTTopTxt = ko.observable();
pg.isDetailDTTop = ko.observable(false);
pg.periodDesc = ko.observable();

pg.breakDown = ko.observableArray([]);
pg.typeChart = ko.observable();
pg.dataSource = ko.observableArray();

pg.availabledatestartscada = ko.observable();
pg.availabledateendscada = ko.observable();
pg.availabledatestartscada2 = ko.observable();
pg.availabledateendscada2 = ko.observable();

pg.availabledatestartalarm = ko.observable();
pg.availabledateendalarm = ko.observable();

pg.availabledatestartscada3 = ko.observable();
pg.availabledateendscada3 = ko.observable();

pg.availabledatestartalarm2 = ko.observable();
pg.availabledateendalarm2 = ko.observable();

pg.availabledatestartwarning = ko.observable();
pg.availabledateendwarning = ko.observable();

pg.dtCompponentAlarm = ko.observable();
pg.labelAlarm = ko.observable("Downtime ");

var height = $(".content").width() * 0.125;

var SeriesDowntime = [{
    field: "AEBOK",
    name: "AEBOK"
}, {
    field: "ExternalStop",
    name: "External Stop"
}, {
    field: "GridDown",
    name: "Grid Down"
}, {
    field: "InternalGrid",
    name: "InternalGrid"
}, {
    field: "MachineDown",
    name: "Machine Down"
}, {
    field: "WeatherStop",
    name: "Weather Stop"
}, {
    field: "Unknown",
    name: "Unknown"
}]

var SeriesAlarm =  [{
                type: "pie",
                field: "result",
                categoryField: "_id",
            }]

pg.breakDownList = ko.observableArray([
    { "value": "dateinfo.dateid", "text": "Date" },
    { "value": "dateinfo.monthdesc", "text": "Month" },
    { "value": "dateinfo.year", "text": "Year" },
    { "value": "projectname", "text": "Project" },
    { "value": "turbine", "text": "Turbine" },
]);

pg.dummyData = ko.observableArray([
    { "name": "TJ001", "prod": 90, "machdown": 50, "griddown": 40, "pcdev": 50, "elecloss": 55, "others": 30 },
    { "name": "TJ002", "prod": 90, "machdown": 50, "griddown": 40, "pcdev": 50, "elecloss": 55, "others": 30 }
]);

pg.ChartLoss = function () {
    var breakDownVal = $("#breakdownlist").data("kendoDropDownList").value();

    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        breakdown: breakDownVal
    };

    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarychart", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var cSeries = res.data.Series
        var cCategories = res.data.Categories;

        $('#lossChart').html("");

        $("#lossChart").kendoChart({
            theme: "Flat",
            legend: {
                position: "top",
                visible: true,
            },
            series: cSeries,
            seriesColors: colorField,
            valueAxes: [{
                line: {
                    visible: false
                },
                max: 100,
                majorUnit: 10,
                labels: {
                    format: "{0}",
                },
                majorGridLines: {
                    visible: false,
                    color: "#eee",
                    width: 0.8,
                },
                name: "lossLine",
                title: { text: "Percentage (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }, {
                line: {
                    visible: false
                },
                // max: 100,
                // majorUnit: 2,
                labels: {
                    format: "{0:n2}",
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "lossBar",
                title: { text: "Capacity (MW)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }],
            categoryAxis: {
                categories: cCategories,
                title: {
                    text: $("#breakdownlist").data("kendoDropDownList").text(),
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                axisCrossingValues: [0, 1000],
                justified: true,
                majorGridLines: {
                    visible: false
                },
            },
            tooltip: {
                visible: true,
                shared: false,
                template: '#if(series.name == "Energy Lost Due to Machine Down (%)"){ ##= "Lost Due to MD" ## }else if(series.name == "Energy Lost Due to Grid Down (%)"){ ##= "Lost Due to GD" ## }else if(series.name == "Others (MWh)"){ ##= "Others" ## }else if(series.name == "Electrical Losses (MWh)"){ ##= "Electrical Losses" ## }else if(series.name == "PC Deviation (MW)"){ ##= "PC Deviation" ## }else{##= series.name ##}# : # if (series.name  == "Electrical Losses (MWh)" || series.name  == "Others (MWh)" ) {##= value # MWh#} else if(series.name == "PC Deviation (MW)"){# #= value # MW #} else {# #= value # % #}#',
                // sharedTemplate:kendo.template($("#templateChart").html()),
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                }
            }
        });
        setTimeout(function () {
            if ($("#lossChart").data("kendoChart") != null) {
                $("#lossChart").data("kendoChart").refresh();
            }
        }, 10);

    });
};

pg.GridLoss = function () {
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
    };

    var requestGridLoss = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarylist", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var gData = res.data.Data

        $('#lossGrid').html("");
        $('#lossGrid').kendoGrid({
            dataSource: {
                data: gData,
                pageSize: 10,
                aggregate: [
                    { field: "Production", aggregate: "sum" },
                    { field: "LossEnergy", aggregate: "sum" },
                    { field: "MachineDownHours", aggregate: "sum" },
                    { field: "GridDownHours", aggregate: "sum" },
                    { field: "EnergyyMD", aggregate: "sum" },
                    { field: "EnergyyGD", aggregate: "sum" },
                    { field: "ElectricLoss", aggregate: "sum" },
                    { field: "PCDeviation", aggregate: "sum" },
                    { field: "Others", aggregate: "sum" },
                ]
            },
            groupable: false,
            sortable: true,
            filterable: false,
            // height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
            height: 399,
            pageable: {
                pageSize: 10,
                input: true, 
            },
            columns: [
                { title: pg.type,field: "Id",width: 100,attributes: {style: "text-align:center;"},headerAttributes: {style: "text-align:center;"},footerTemplate: "<center>Total (All Turbines)</center>"}, 
                { title: "Production (MWh)", headerAttributes: { tyle: "text-align:center;"}, field: "Production",width: 100,attributes: { class: "align-center" },format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" }, 
                { title: "Lost Energy (MWh)",headerAttributes: {style: "text-align:center;"},field: "LossEnergy", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"},
                {
                    title: "Downtime : Duration (Hrs)",
                    headerAttributes: {
                        style: 'font-weight: bold; text-align: center;'
                    },
                    columns: [
                        {
                            title: "Machine",
                            headerAttributes: { style: "text-align:center;" },
                            field: "MachineDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                        },
                        {
                            title: "Grid",
                            headerAttributes: { style: "text-align:center;" },
                            field: "GridDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                        },
                    ]
                }, {
                    title: "Downtime : Energy Loss (MWh)",
                    headerAttributes: {
                        style: 'font-weight: bold; text-align: center;'
                    },
                    columns: [
                        {
                            title: "Machine",
                            headerAttributes: { style: "text-align:center;" },
                            field: "EnergyyMD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                        },
                        {
                            title: "Grid",
                            headerAttributes: { style: "text-align:center;" },
                            field: "EnergyyGD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                        },
                    ]
                }, {
                    title: "Electrical Losses (MWh)",
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    field: "ElectricLoss", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                }, {
                    title: "Power Curve Deviation (MW)", //Sepertinya ini MW
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    field: "PCDeviation", width: 120, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                }, {
                    title: "Others (MWh)", //Sepertinya ini KWh
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    field: "Others", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                }]
        })
    });
};

pg.DTDuration = function (dataSource,id,Series,legend,name,vislabel,rotate,heightParam) {
    // $("#" + id).height("300px")
    // var Height = $('#CompAlarm').width() * heightParam

    $("#" + id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: [
                { "field": "Total", "dir": "desc" }
            ],
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: legend,
        },
        chartArea: {
            height: heightParam, 
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            labels: {
                        visible: vislabel,
                        background: "transparent",
                        template: "#= category #: \n #= kendo.format('{0:N1}', value)# Hours",
                    }
        },
        series: Series,
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: "Hours",
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
            }
        },
        categoryAxis: {
            // visible: legend,
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: rotate
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n1') # Hours",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        // seriesClick: function (e) {
        //     pg.toDetailDTTop(e, "Hours");
        // }
    });

    setTimeout(function () {
        if ($("#" + id).data("kendoChart") != null) {
            $("#" + id).data("kendoChart").refresh();
        }
    }, 100);
}

pg.DTFrequency = function (dataSource,id,Series,legend,name,vislabel,rotate,heightParam) {
    // $("#" + id).height("300px")
    // var Height = $('#CompAlarm').width() * height
    $("#" + id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "Total", dir: 'desc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: legend,
        },
        chartArea: {
            height: heightParam,
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            labels: {
                        visible: vislabel,
                        background: "transparent",
                        template: "#= category #: \n #= kendo.format('{0:N1}', value)#",
                    }
        },
        series: Series,
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: "Times",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            name: "result",
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
            }
        },
        categoryAxis: {
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            dir: "desc",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: rotate
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            // format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            // template: "#: category #: #: kendo.toString(value, 'n1') # Hours",
            template: "#: category #: #: value # ",
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        // seriesClick: function (e) {
        //     pg.toDetailDTTop(e, "Times");
        // }
    });

    setTimeout(function () {
        if ($("#" + id).data("kendoChart") != null) {
            $("#" + id).data("kendoChart").refresh();
        }
    }, 100);
}

pg.TopTurbineLoss = function (dataSource,id,Series,legend,name,vislabel,rotate,heightParam) {
    // $("#" + id).height("300px")
    // var Height = $('#CompAlarm').width() * height
    $("#" + id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: [
                { "field": "Total", "dir": "desc" }
            ],
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: legend,
        },
        chartArea: {
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 100) / 2,
            height: heightParam, 
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            labels: {
                        visible: vislabel,
                        background: "transparent",
                        template: "#= category #: \n #= kendo.format('{0:N1}', value/1000)# MWh",
                    }
            // opacity : 0.7
        },
        series: Series,
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
            title: {
                text: "MWh",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                step: 2,
                template: "#: kendo.toString(value/1000, 'n1') #",
            },
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            }
        },
        categoryAxis: {
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: rotate
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value/1000, 'n1') # MWh",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        // seriesClick: function (e) {
        //     avail.toDetailDTTop(e, "MWh");
        // }
    });

    setTimeout(function () {
        if ($("#" + id).data("kendoChart") != null) {
            $("#" + id).data("kendoChart").refresh();
        }
    }, 100);
}

pg.TLossCat = function (id, byTotalLostenergy, dataSource, measurement) {
    var gapVal = 1
    var templateLossCat = ''
    switch (dataSource.length) {
        case 1:
            gapVal = 5;
            break;
        case 2:
            gapVal = 3;
            break;
        case 3:
            gapVal = 1;
            break;
        case 4:
            gapVal = 1;
            break;
        case 5:
            gapVal = 1;
            break;
    } 

    if(measurement == "MWh") {
       templateLossCat = "<b>#: category # :</b> #: kendo.toString(value/1000, 'n1')# " + measurement
    } else if(measurement == "Hours") {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n1')# " + measurement
    } else {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n0')# "
    }

    $('#' + id).html("");
    $('#' + id).kendoChart({
        dataSource: {
            data: dataSource,
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true
        },
        chartArea: {
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
            height: 163,
        },
        seriesDefaults: {
            type: "column",
            gap: gapVal,
        },
        series: [{
            type: "column",
            field: "result",
        }],
        seriesColor: colorField,
        valueAxis: {
            labels: {
                step: 2,
                template: (byTotalLostenergy == true) ? "#= value / 1000 #" : "#= value#"
            },
            title: {
                text: measurement,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
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

        },
        categoryAxis: {
            field: "_id.id2",
            title: {
                text: "Loss Categories",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            template: templateLossCat,
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

pg.ChartAvailability = function () {

    pg.breakDownVal = $("#breakdownlistavail").data("kendoDropDownList").value();

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
            return;
        }
        pg.dataSource(res.data);
        pg.createChartAvailability(pg.dataSource());
        pg.createChartProduction(pg.dataSource());

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
    $("#availabilityChart").kendoChart({
        theme: "Flat",
        legend: {
            position: "top",
            visible: true,
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

pg.ChartWindAvail = function () {
    // fa.getProjectInfo();
    fa.LoadData();
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwindavailability/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var data = res.data;

        $("#windAvailabilityChart").html("");
        $("#windAvailabilityChart").kendoChart({
            dataSource: {
                data: data,
                sort: { field: "WindSpeed", dir: 'asc' }
            },
            theme: "Flat",
            chartArea: {
                // height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
                height: 400,
            },
            legend: {
                position: "top",
                visible: true,
            },
            series: [{
                type: "column",
                field: "TotalAvail",
                axis: "windPercentage",
                name: "Total Availability [%]",
                opacity: 0.6
            }, {
                type: "line",
                style: "smooth",
                field: "Time",
                axis: "windPercentage",
                name: "Cumulative % of Time",
                markers: {
                    visible: false,
                },
                width: 3,
            }, {
                type: "line",
                style: "smooth",
                field: "Energy",
                axis: "cumProd",
                name: "Cumulative % of Energy Delivered",
                markers: {
                    visible: false,
                },
                width: 3,
            }],
            seriesColors: colorFields2,
            valueAxes: [{
                line: {
                    visible: false
                },
                max: 100,
                majorUnit: 20,
                labels: {
                    format: "{0}%",
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "windPercentage",
                title: { text: "Availability (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }, {
                line: {
                    visible: false
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                max: 100,
                labels: {
                    format: "{0}%",
                },
                name: "cumProd",
                title: { text: "Cumulative Production (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }],
            categoryAxis: {
                field: "WindSpeed",
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                axisCrossingValues: [0, 1000],
                justified: true,
                majorGridLines: {
                    visible: false
                },
            },
            tooltip: {
                visible: true,
                shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                // template: "#= series.name # : #= kendo.toString(value, 'n2')# at #= category #",
                template: "#= kendo.toString(value, 'n2')#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
        });

        $("#windAvailabilityChart").data("kendoChart").refresh();
    });
};

pg.DTLEbyType = function (dataSource) {
    $("#chartDTLEbyType").kendoChart({
        dataSource: {
            data: dataSource,
            group: [{ field: "_id.id3" }],
            sort: { field: "_id.id1", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
        },
        chartArea: {
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
            height: 163,
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "powerlost",
            // opacity : 0.7,
            stacked: true,
            axis: "PowerLost"
        },
        {
            name: function () {
                return "Duration";
            },
            type: "line",
            field: "duration",
            axis: "Duration",
            markers: {
                visible: false
            }
        },
        {
            name: function () {
                return "Frequency";
            },
            type: "line",
            field: "frequency",
            axis: "Frequency",
            markers: {
                visible: false
            }
        }],
        seriesColor: colorField,
        valueAxis: [{
            name: "PowerLost",
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
        },
        {
            name: "Duration",
            title: { visible: false },
            visible: false,
        },
        {
            name: "Frequency",
            title: { visible: false },
            visible: false,
        }],
        categoryAxis: {
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            sharedTemplate: kendo.template($("#templateDTLE").html()),
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });
}

pg.backToDownTime = function () {

    pg.isDetailDTTop(false);
    pg.detailDTTopTxt("");
}

pg.getPeriodDesc = function () {
    var duration = ((fa.dateEnd - fa.dateStart) / 86400000) + 1;
    var breakDownVal = $("#breakdownlist").data("kendoDropDownList").value();
    if (breakDownVal == "dateinfo.dateid") {
        pg.periodDesc = fa.dateStart + " to " + fa.dateEnd;
    } else {
        pg.periodDesc = "";
    }
}

pg.loadData = function () {

    // fa.getProjectInfo();
    setTimeout(function () {
        pg.SetBreakDown();

        if (fa.project == "") {
            pg.type = "Project Name";
        } else {
            pg.type = "Turbine";
        }

        pg.ChartAvailability();
        pg.ChartWindAvail();
        pg.ChartLoss();
        pg.GridLoss();
        warn.loadData();

        // moment(Date.UTC(date1.getFullYear(), date1.getMonth(), date1.getDate(), 0, 0, 0)).toISOString()

        // Top 10 Downtime
        // var param = {ProjectName : fa.project, Date: maxdate};
        var param = {
            // dateStart: fa.dateStart,
            period: fa.period,
            dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
            dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
            turbine: fa.turbine,
            project: fa.project,
        }

        // console.log(param);
        pg.getPeriodDesc();
        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            var minDatetemp = new Date(res.data.ScadaData[0]);
            var maxDatetemp = new Date(res.data.ScadaData[1]);

            pg.availabledatestartscada(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            pg.availabledateendscada(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            pg.availabledatestartscada2(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            pg.availabledateendscada2(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            pg.availabledatestartalarm(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
            pg.availabledateendalarm(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

            pg.availabledatestartscada3(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            pg.availabledateendscada3(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            pg.availabledatestartalarm2(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
            pg.availabledateendalarm2(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

            pg.availabledatestartwarning(kendo.toString(moment.utc(res.data.Warning[0]).format('DD-MMMM-YYYY')));
            pg.availabledateendwarning(kendo.toString(moment.utc(res.data.Warning[1]).format('DD-MMMM-YYYY')));

            $('#availabledatestart').html(pg.availabledatestartscada());
            $('#availabledateend').html(pg.availabledateendscada());
        })

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettop10", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            pg.dtCompponentAlarm(res.data)

            var HDowntime = $('#filter-analytic').width() * 0.2
            var HAlarm = $('#filter-analytic').width() * 0.235
            
            // ===== Downtime =====
            pg.DTDuration(res.data.duration,'chartDTDuration',SeriesDowntime,true,"Turbine",false,-330,HDowntime);
            pg.DTFrequency(res.data.frequency,'chartDTFrequency',SeriesDowntime,true,"Turbine",false,-330,HDowntime);
            pg.TopTurbineLoss(res.data.loss,'chartTopTurbineLoss',SeriesDowntime,true,"Turbine",false,-330,HDowntime);

            // ===== Alarm =====
            pg.DTDuration(res.data.componentduration,'chartCADuration',SeriesAlarm,true, "",false,-90,HAlarm);
            pg.DTFrequency(res.data.componentfrequency,'chartCAFrequency',SeriesAlarm,true, "",false,-90,HAlarm);
            pg.TopTurbineLoss(res.data.componentloss,'chartCATurbineLoss',SeriesAlarm,true, "",false,-90,HAlarm);

            pg.TLossCat('chartLCByTEL', true, res.data.catloss, 'MWh');
            pg.TLossCat('chartLCByDuration', false, res.data.catlossduration, 'Hours');
            pg.TLossCat('chartLCByFreq', false, res.data.catlossfreq, 'Times');

            app.loading(false);
        });
        var paramdown = {
            Period: fa.period,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Turbine: fa.turbine,
            Project: fa.project
        };
        toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeloss", paramdown, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            pg.DTLEbyType(res.data);
        });
    }, 100);
}

pg.refreshGrid = function (param) {
    app.loading(true);
    if(param != "undefined"){
        if(param == "staticView"){
            $('#availabledatestart').html(pg.availabledatestartscada());
            $('#availabledateend').html(pg.availabledateendscada());
        }
        if(param == 'top10downtime'){
            $('#availabledatestart').html(pg.availabledatestartalarm2());
            $('#availabledateend').html(pg.availabledateendalarm2());
        }
        if(param == 'availability'){
            $('#availabledatestart').html(pg.availabledatestartscada3());
            $('#availabledateend').html(pg.availabledateendscada3());
        }
        if(param == 'lostenergy'){
            $('#availabledatestart').html(pg.availabledatestartalarm());
            $('#availabledateend').html(pg.availabledateendalarm());
        }
        if(param == 'reliabilitykpi'){
            $('#availabledatestart').html(pg.availabledatestartscada2());
            $('#availabledateend').html(pg.availabledateendscada2());
        }
        if(param == 'windspeedavail'){
            $('#availabledatestart').html(pg.availabledatestartscada2());
            $('#availabledateend').html(pg.availabledateendscada2());
        }
        if(param == 'warning'){
            $('#availabledatestart').html(pg.availabledatestartwarning());
            $('#availabledateend').html(pg.availabledateendwarning());
        }
    }
    setTimeout(function () {
        if ($("#gridLoss").data("kendoGrid") != null) {
            $("#gridLoss").data("kendoGrid").refresh();
        }
        if ($("#lossChart").data("kendoChart") != null) {
            $("#lossChart").data("kendoChart").refresh();
        }

        if ($("#chartDTDuration").data("kendoChart") != null) {
            $("#chartDTDuration").data("kendoChart").refresh();
        }

        if ($("#chartDTFrequency").data("kendoChart") != null) {
            $("#chartDTFrequency").data("kendoChart").refresh();
        }

        if ($("#chartTopTurbineLoss").data("kendoChart") != null) {
            $("#chartTopTurbineLoss").data("kendoChart").refresh();
        }

        if ($("#availabilityChart").data("kendoChart") != null) {
            $("#availabilityChart").data("kendoChart").refresh();
        }

        if ($("#chartLCByTEL").data("kendoChart") != null) {
            $("#chartLCByTEL").data("kendoChart").refresh();
        }

        if ($("#chartLCByDuration").data("kendoChart") != null) {
            $("#chartLCByDuration").data("kendoChart").refresh();
        }

        if ($("#chartLCByFreq").data("kendoChart") != null) {
            $("#chartLCByFreq").data("kendoChart").refresh();
        }

        if ($("#chartDTLEbyType").data("kendoChart") != null) {
            $("#chartDTLEbyType").data("kendoChart").refresh();
        }

        if ($("#windAvailabilityChart").data("kendoChart") != null) {
            $("#windAvailabilityChart").data("kendoChart").refresh();
        }
        
       if ($("#productionChart").data("kendoChart") != null) {
            $("#productionChart").data("kendoChart").refresh();
        }

        if ($("#warningGrid").data("kendoGrid") != null) {
            $("#warningGrid").data("kendoGrid").refresh();
        }

        if ($("#chartCADuration").data("kendoChart") != null) {
            $("#chartCADuration").data("kendoChart").refresh();
        }

        if ($("#chartCAFrequency").data("kendoChart") != null) {
            $("#chartCAFrequency").data("kendoChart").refresh();
        }

        if ($("#chartCATurbineLoss").data("kendoChart") != null) {
            $("#chartCATurbineLoss").data("kendoChart").refresh();
        }

        app.loading(false);
    },1500);
}

pg.SetBreakDown = function () {
    fa.disableRefreshButton(true);
    pg.breakDown = [];

    setTimeout(function () {
        $.each(fa.GetBreakDown(), function (i, val) {
            if (val.value == "Turbine" || val.value == "Project") {
                return false;
            } else {
               pg.breakDown.push(val);
            }
        });

        $("#breakdownlistavail").data("kendoDropDownList").dataSource.data(pg.breakDown);
        $("#breakdownlistavail").data("kendoDropDownList").dataSource.query();
        $("#breakdownlistavail").data("kendoDropDownList").select(0);


        fa.disableRefreshButton(false);
    }, 500);
}

viewModel.Warning = new Object();
var warn = viewModel.Warning;
warn.dataSource = ko.observableArray();

warn.generateGrid = function () {
    var config = {
        dataSource: {
            data: warn.dataSource(),
            pageSize: 10
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Warning Description", field: "desc", attributes: { class: "align-left row-custom" }, width: 200, locked: true, filterable: false },
        ],
        dataBound: function(){
            setTimeout(function(){
                $("#warningGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","20px");
                $("#warningGrid >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","20px");
                // app.loading(false);
            },200);
        },
    };

    $.each(warn.dataSource()[0].turbines, function (i, val) {
        var column = {
            title: val.turbine,
            field: "turbines["+i+"].count",
            attributes: { class: "align-center" },
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            width: 100
        }

        config.columns.push(column);
    });

    $('#warningGrid').html("");
    $('#warningGrid').kendoGrid(config);
    $('#warningGrid').data('kendoGrid').refresh();

    // setTimeout(function() {
    //     app.loading(false);
    // }, 500);
}

warn.loadData = function() {
    var param = {
        period: fa.period,
        Turbine: fa.turbine,
        DateStart: fa.dateStart,
        DateEnd: fa.dateEnd,
        Project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getwarning", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.Data.length == 0) {
            return;
        }
        warn.dataSource(res.data.Data);
        warn.generateGrid();
    });
}


vm.currentMenu('Losses and Efficiency');
vm.currentTitle('Losses and Efficiency');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticloss' }]);


$(document).ready(function () {
    fa.LoadData();
    $('#btnRefresh').on('click', function () {
         app.loading(true);
        fa.LoadData();
        pg.loadData();
    });

    $('#breakdownlist').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.loadData() },
    });

    $('#breakdownlistavail').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.loadData() },
    });

    setTimeout(function () {
        fa.LoadData();
        pg.loadData();
    }, 1000);

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

    }, 1500);


    $("input[name=IsAlarm]").on("change", function() {

            var HAlarm = $('#filter-analytic').width() * 0.235
        
            var data = pg.dtCompponentAlarm()
            if(this.id == "alarm"){   
                SeriesAlarm =  [{
                    field: "result",
                    name: "Downtime"
                }]             
                // ===== Alarm =====
                pg.DTDuration(data.alarmduration,'chartCADuration',SeriesAlarm,false, "",false,-90,HAlarm);
                pg.DTFrequency(data.alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "",false,-90,HAlarm);
                pg.TopTurbineLoss(data.alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "",false,-90,HAlarm);

                pg.labelAlarm(" Top 10 Downtime")
            }else{     
                 SeriesAlarm = [{
                    type: "pie",
                    field: "result",
                    categoryField: "_id",
                }]           
                // ===== Component =====
                pg.DTDuration(data.componentduration,'chartCADuration',SeriesAlarm,true, "",false,-90,HAlarm);
                pg.DTFrequency(data.componentfrequency,'chartCAFrequency',SeriesAlarm,true, "",false,-90,HAlarm);
                pg.TopTurbineLoss(data.componentloss,'chartCATurbineLoss',SeriesAlarm,true, "",false,-90,HAlarm);

                pg.labelAlarm(" Downtime")
            }
    });

});