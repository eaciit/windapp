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
var height = $(".content").width() * 0.125;

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
    app.loading(true);
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
            app.loading(false);
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

        app.loading(false);
    });
};

pg.GridLoss = function () {
    app.loading(true);

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
            height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
            pageable: true,
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
    
    $.when(requestGridLoss).done(function(){
        setTimeout(function(){
            app.loading(false);
        },500)
    });
};

pg.DTDuration = function (dataSource) {
    $("#chartDTDuration").kendoChart({
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
            visible: true,
        },
        chartArea: {
            height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 100) / 2,
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [{
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
        }],
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
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
            field: "_id",
            title: {
                text: "Turbine",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
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
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n2') # Hours",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        seriesClick: function (e) {
            pg.toDetailDTTop(e, "Hours");
        }
    });
}

pg.DTFrequency = function (dataSource) {
    $("#chartDTFrequency").kendoChart({
        dataSource: {
            data: dataSource,
            // group: [{field: "_id.id4"}],
            sort: { field: "Total", dir: 'desc' }
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
            height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 100) / 2,
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [{
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
        }],
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
                text: "Turbine",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            dir: "desc",
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
        seriesClick: function (e) {
            pg.toDetailDTTop(e, "Times");
        }
    });
}

pg.TopTurbineLoss = function (dataSource) {
    $("#chartTopTurbineLoss").kendoChart({
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
            visible: true,
        },
        chartArea: {
            height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 100) / 2,
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [{
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
        }],
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
                text: "Turbine",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
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
}

pg.TLossCat = function (id, byTotalLostenergy, dataSource, measurement) {
    var gapVal = 1
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
            height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
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
            template: (byTotalLostenergy == true) ? "<b>#: category # :</b> #: kendo.toString(value/1000, 'n1')# " + measurement : "<b>#: category # :</b> #: kendo.toString(value, 'n1')# " + measurement,
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
    app.loading(true);

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
}

pg.ChartWindAvail = function () {
    fa.getProjectInfo();
    fa.LoadData();
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwindavailability/getdata", param, function (res) {
        var data = res.data;

        $("#windAvailabilityChart").html("");
        $("#windAvailabilityChart").kendoChart({
            dataSource: {
                data: data,
                sort: { field: "WindSpeed", dir: 'asc' }
            },
            theme: "Flat",
            chartArea: {
                height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
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

        app.loading(false);
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
            height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
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
            format: "{0:n1}",
            // sharedTemplate: kendo.template($("#templateDTLEbyType").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        /*seriesClick: function (e) {
            avail.toDetailDTLostEnergy(e, false, "chartbytype");
        }*/
    });
}

pg.DTTopDetail = function (turbine, type) {
    // app.loading(true);
    // var project = fa.project;
    // var date = new Date(Date.UTC(2016,5,30,23,59,59,0));
    // var param = {ProjectName : project, Date: date, Type: type, Turbine: turbine};

    // var template = (type == 'Hours' ? "#: category # : #:  kendo.toString(value, 'n1') #" : "#: category # : #:  kendo.toString(value, 'n0') #")
    // toolkit.ajaxPost(viewModel.appName +"dashboard/getdowntimetopdetail",param, function (res) {
    //     if (!toolkit.isFine(res)) {
    //         return;
    //     }

    //     var dataSource = res.data

    //     $("#chartDTTopDetail").kendoChart({
    //         dataSource: {
    //             data: dataSource,
    //         },
    //         theme: "flat",
    //         title: {
    //             text: ""
    //         },
    //         legend: {
    //             position: "top",
    //             visible : false,
    //         },
    //         chartArea: {
    //         height : 160
    //         },
    //         seriesDefaults: {
    //             area: {
    //                 line: {
    //                     style: "smooth"
    //                 }
    //             }
    //         },
    //         series: [{
    //             // name : "Lost Energy",
    //             field : "result",
    //             // opacity : 0.7,
    //         }],
    //         seriesColors: colorField,
    //         valueAxis: {
    //         //majorUnit: 100,
    //             labels:{
    //                 step : 2
    //             },
    //             line: {
    //                 visible: false
    //             },
    //             axisCrossingValue: -10,
    //             majorGridLines: {
    //                 visible: true,
    //                 color: "#eee",
    //                 width: 0.8,
    //             }
    //         },
    //         categoryAxis: {
    //             field: "_id.id2",
    //             majorGridLines: {
    //                 visible: false
    //             },
    //             labels:{
    //                 template: '#=  value.substring(0,3) #'
    //             },
    //             majorTickType: "none"
    //         },
    //         tooltip: {
    //             visible: true,
    //             template : template,
    //             background: "rgb(255,255,255, 0.9)",
    //             color : "#58666e",
    //             font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
    //             border : {
    //                 color : "#eee",
    //                 width : "2px",
    //             },

    //         }
    //     });

    //     app.loading(false);
    //     $("#chartDTTopDetail").data("kendoChart").refresh();
    // });

    // $('#gridDTTopDetail').html("");
    // $('#gridDTTopDetail').kendoGrid({
    //     dataSource: {
    //         serverPaging: true,
    //         serverSorting: true,
    //         transport: {
    //             read: {
    //             url: viewModel.appName + "dashboard/getdowntimetopdetailtable",
    //             type: "POST",
    //             data: param,
    //             dataType: "json",
    //             contentType: "application/json; charset=utf-8"
    //             },
    //             parameterMap: function(options) {
    //             return JSON.stringify(options);
    //             }
    //         },
    //         pageSize: 10,
    //         schema: {
    //             data: function(res){
    //                 if (!app.isFine(res)) {
    //                     return;
    //                 }
    //                 return res.data.data
    //             },
    //             total: function(res){
    //                 if (!app.isFine(res)) {
    //                     return;
    //                 }
    //                 return res.data.total;
    //             }
    //         },
    //         sort: [
    //             { field: 'StartDate', dir: 'asc' }
    //         ],
    //     },
    //     groupable: false,
    //     sortable: true,
    //     filterable: false,
    //     pageable: true,
    //     //resizable: true,
    //     columns: [
    //         { title: "Date", field: "StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80},
    //         { title: "Start Time", field: "StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 75, attributes:{style:"text-align:center;"} },
    //         { title: "End Date", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80},
    //         { title: "End Time", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 70, attributes:{style:"text-align:center;"}},
    //         { title: "Alert Description", field: "AlertDescription", width: 200 },
    //         { title: "External Stop", field: "ExternalStop", width: 80 , sortable: false, template:'# if (ExternalStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "Grid Down", field: "GridDown", width: 80 , sortable: false, template:'# if (GridDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "Internal Grid", field: "InternalGrid", width: 80 , sortable: false, template:'# if (InternalGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "Machine Down", field: "MachineDown", width: 80 , sortable: false, template:'# if (MachineDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "AEbOK", field: "AEbOK", width: 80 , sortable: false, template:'# if (AEbOK == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "Unknown", field: "Unknown", width: 80 , sortable: false, template:'# if (Unknown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},
    //         { title: "Weather Stop", field: "WeatherStop", width: 80 , sortable: false, template:'# if (WeatherStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style:"text-align: center" }, attributes:{style:"text-align:center;"}},            
    //     ]
    // });

}

pg.toDetailDTTop = function (e, type) {
    // vm.isDashboard(false);

    // if (type == "Times") {
    //     pg.detailDTTopTxt("("+e.category+") - Frequency");
    // }else{
    //     pg.detailDTTopTxt("("+e.category+") - "+type);
    // }
    // pg.isDetailDTTop(true);

    // // get the data and push into the chart    
    // pg.DTTopDetail(e.category, type);
}

pg.backToDownTime = function () {
    // vm.isDashboard(true);

    // pg.isSummary(false);
    // pg.isProduction(false);
    // pg.isAvailability(false);
    // pg.isDowntime(true);

    // pg.isDetailDTLostEnergy(false);
    // pg.detailDTLostEnergyTxt("Lost Energy for Last 12 months");

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
    fa.getProjectInfo();
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
            var minDatetemp = new Date(res.ScadaData[0]);
            var maxDatetemp = new Date(res.ScadaData[1]);
            $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            $('#availabledatestartscada2').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            $('#availabledateendscada2').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            $('#availabledatestartalarm').html(kendo.toString(moment.utc(res.Alarm[0]).format('DD-MMMM-YYYY')));
            $('#availabledateendalarm').html(kendo.toString(moment.utc(res.Alarm[1]).format('DD-MMMM-YYYY')));

            $('#availabledatestartscada3').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            $('#availabledateendscada3').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

            $('#availabledatestartalarm2').html(kendo.toString(moment.utc(res.Alarm[0]).format('DD-MMMM-YYYY')));
            $('#availabledateendalarm2').html(kendo.toString(moment.utc(res.Alarm[1]).format('DD-MMMM-YYYY')));
        })

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettop10", param, function (res) {
            // console.log(res);
            if (!toolkit.isFine(res)) {
                return;
            }

            pg.DTDuration(res.data.duration);
            pg.DTFrequency(res.data.frequency);
            pg.TopTurbineLoss(res.data.loss);
            // pg.TLossCat('chartLCByLTE', true, res.data.catloss, 'MWh');
            pg.TLossCat('chartLCByTEL', true, res.data.catloss, 'MWh');
            pg.TLossCat('chartLCByDuration', true, res.data.catlossduration, 'Hours');
            pg.TLossCat('chartLCByFreq', true, res.data.catlossfreq, 'Times');
        });
        var paramdown = {
            Period: fa.period,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Turbine: fa.turbine,
            Project: fa.project
        };
        toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeloss", paramdown, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }
            pg.DTLEbyType(res.data);
        });
    }, 100);
    // console.log(app.dateAvail);
    var tes = app.dateAvail;
    for(var key in tes){
        // The key is key
        // The value is obj[key]
        console.log(key);
    }
}

pg.refreshGrid = function () {
    app.loading(true);
    var refresh = setTimeout(function () {
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

        app.loading(false);
    }, 1000);
}

pg.SetBreakDown = function () {
    pg.breakDown = [];

    setTimeout(function () {
        $.each(pg.breakDownList(), function (i, valx) {
            $.each(fa.GetBreakDown(), function (i, valy) {
                if (valx.text == valy.text) {
                    pg.breakDown.push(valx);
                }
            });
        });

        $("#breakdownlist").data("kendoDropDownList").dataSource.data(pg.breakDown);
        $("#breakdownlist").data("kendoDropDownList").dataSource.query();
        if ($("#breakdownlist").data("kendoDropDownList").value() == "") {
            $("#breakdownlist").data("kendoDropDownList").select(0);
        }

        $("#breakdownlistavail").data("kendoDropDownList").dataSource.data(fa.GetBreakDown());
        $("#breakdownlistavail").data("kendoDropDownList").dataSource.query();
        if ($("#breakdownlistavail").data("kendoDropDownList").value() == "") {
            $("#breakdownlistavail").data("kendoDropDownList").select(0);
        }
    }, 1000);
}

vm.currentMenu('Losses and Efficiency');
vm.currentTitle('Losses and Efficiency');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticloss' }]);


$(document).ready(function () {
    app.loading(true);
    fa.LoadData();
    $('#btnRefresh').on('click', function () {
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

});