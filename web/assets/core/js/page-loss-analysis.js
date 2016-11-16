'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.type = ko.observable();
pg.detailDTTopTxt = ko.observable();
pg.isDetailDTTop = ko.observable(false);
pg.periodDesc = ko.observable();

pg.breakDown = ko.observableArray([]);

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
        // console.log(res)
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
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarylist", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        // console.log(res)
        var gData = res.data.Data

        $('#lossGrid').html("");
        $('#lossGrid').kendoGrid({
            dataSource: {
                data: gData,
                pageSize: 10
            },
            groupable: false,
            sortable: true,
            filterable: false,
            pageable: true,
            columns: [
                {
                    title: pg.type,
                    field: "Id",
                    width: 100,
                    attributes: {
                        style: "text-align:center;"
                    },
                    headerAttributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Production",
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    columns: [
                        { title: "(Hours)", field: "Production", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                    ]
                }, {
                    title: "Lost Energy",
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    columns: [
                        { title: "(MWh)", field: "LossEnergy", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                    ]
                    // field: "LossEnergy",
                    // width: 100,
                    // attributes: {
                    //     class: "align-right"
                    // },
                    // format: "{0:n2}",
                    // headerAttributes: {
                    //     style: "text-align:center;"
                    // }
                },
                {
                    title: "Downtime : Duration",
                    headerAttributes: {
                        style: 'font-weight: bold; text-align: center;'
                    },
                    columns: [
                        {
                            title: "Machine",
                            headerAttributes: { style: "text-align:center;" },
                            columns: [
                                { title: "(Hours)", field: "MachineDownHours", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                            ]
                            // field: "MachineDownHours", width: 100, attributes: { class: "align-right" },format: "{0:n2}",headerAttributes: { style: "text-align:center;" } 
                        },
                        {
                            title: "Grid",
                            headerAttributes: { style: "text-align:center;" },
                            columns: [
                                { title: "(Hours)", field: "GridDownHours", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                            ]
                            // field: "GridDownHours", width: 100, attributes: { class: "align-right" }, format: "{0:n2}",headerAttributes: { style: "text-align:center;" }
                        },
                    ]
                }, {
                    title: "Downtime : Energy Loss",
                    headerAttributes: {
                        style: 'font-weight: bold; text-align: center;'
                    },
                    columns: [
                        {
                            title: "Machine",
                            headerAttributes: { style: "text-align:center;" },
                            columns: [
                                { title: "(MWh)", field: "EnergyyMD", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                            ]
                            // field: "EnergyyMD", width: 100, attributes: { class: "align-right" },format: "{0:n2}",headerAttributes: { style: "text-align:center;" } 
                        },
                        {
                            title: "Grid",
                            headerAttributes: { style: "text-align:center;" },
                            columns: [
                                { title: "(MWh)", field: "EnergyyMD", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                            ]
                            // field: "EnergyyGD", width: 100, attributes: { class: "align-right" }, format: "{0:n2}",headerAttributes: { style: "text-align:center;" }
                        },
                    ]
                }, {
                    title: "Electrical Losses",
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    columns: [
                        { title: "(MWh)", field: "ElectricLoss", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                    ]
                    // field: "ElectricLoss",
                    // width: 100,
                    // attributes: {
                    //     class: "align-right"
                    // },
                    // format: "{0:n2}",
                    // headerAttributes: {
                    //     style: "text-align:center;"
                    // }
                }, {
                    title: "Power Curve Deviation", //Sepertinya ini MW
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    columns: [
                        { title: "(MWh)", field: "PCDeviation", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                    ]
                    // field: "PCDeviation",
                    // width: 100,
                    // attributes: {
                    //     class: "align-right"
                    // },
                    // format: "{0:n2}",
                    // headerAttributes: {
                    //     style: "text-align:center;"
                    // }
                }, {
                    title: "Others", //Sepertinya ini KWh
                    headerAttributes: {
                        style: "text-align:center;"
                    },
                    columns: [
                        { title: "(MWh)", field: "Others", width: 100, attributes: { class: "align-right" }, format: "{0:n2}" },
                    ]
                    // field: "Others",
                    // width: 100,
                    // attributes: {
                    //     class: "align-right"
                    // },
                    // format: "{0:n2}",
                    // headerAttributes: {
                    //     style: "text-align:center;"
                    // }
                }]
        })
        app.loading(false);
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
            height: 220,
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
            height: 220,
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
            height: 220,
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
            height: 220
        },
        seriesDefaults: {
            type: "column",
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

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettop10", param, function (res) {
            // console.log(res);
            if (!toolkit.isFine(res)) {
                return;
            }

            pg.DTDuration(res.data.duration);
            pg.DTFrequency(res.data.frequency);
            pg.TopTurbineLoss(res.data.loss);
            pg.TLossCat('chartLCByLTE', true, res.data.catloss, 'MWh');
        });

        // toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettop10", param, function (res) {
        //     if (!toolkit.isFine(res)) {
        //         return;
        //     }

        // });


    }, 100);
}

pg.refreshGrid = function () {
    setTimeout(function () {
        if ($("#gridLoss").data("kendoGrid") != null) {
            $("#gridLoss").data("kendoGrid").refresh();
        }
        if ($("#lossChart").data("kendoChart") != null) {
            $("#lossChart").data("kendoChart").refresh();
        }

        $("#chartDTDuration").data("kendoChart").refresh();
        $("#chartDTFrequency").data("kendoChart").refresh();
        $("#chartTopTurbineLoss").data("kendoChart").refresh();
        $("#chartLCByLTE").data("kendoChart").refresh();

    }, 10);
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
    }, 1000);
}

vm.currentMenu('Losses and Efficiency');
vm.currentTitle('Losses and Efficiency');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticwindavailability' }]);


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