'use strict';

viewModel.DatabrowserAlarmAnomaly = new Object();
var dbaa = viewModel.DatabrowserAlarmAnomaly;

dbaa.InitGridAlarmAnomalies = function() {
    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
    }

    var filters = [{
        field: "startdate",
        operator: "gte",
        value: fa.dateStart
    }, {
        field: "startdate",
        operator: "lte",
        value: fa.dateEnd
    }, {
        field: "turbine",
        operator: "in",
        value: turbine
    }, ];
    
    var filter = {
        filters: filters
    }
    var param = {
        filter: filter
    };

    $('#dataGridAlarmAnomalies').html("");
    $('#dataGridAlarmAnomalies').kendoGrid({
        selectable: "multiple",
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getalarmscadaanomalylist",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function(res) {
                    app.isFine(res);
                    app.loading(false);
                    return res.data.Data
                },
                total: function(res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbinealarmAnomalies').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataalarmAnomalies').html(kendo.toString(res.data.Total, 'n0'));
                    return res.data.Total;
                }
            },
            sort: [{
                field: 'StartDate',
                dir: 'asc'
            }, {
                field: 'Turbine',
                dir: 'asc'
            }],
        },
        // toolbar: [
        //     "excel"
        // ],
        excel: {
            fileName: "Alarm Anomalies.xlsx",
            filterable: true,
            allPages: true
        },
        groupable: false,
        sortable: true,
        filterable: false,
        pageable: {
            pageSize: 10,
            input:true, 
        },
        resizable: true,
        columns: [{
                title: "Date",
                field: "StartDate",
                template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                width: 80
            }, {
                title: "Turbine",
                field: "Turbine",
                width: 90,
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Start Time",
                field: "StartDate",
                template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 75,
                attributes: {
                    style: "text-align:center;"
                }
            },
            /*{ title: "Farm", field: "Farm", width: 100 },*/
            {
                title: "End Date",
                field: "EndDate",
                template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                width: 80
            }, {
                title: "End Time",
                field: "EndDate",
                template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 70,
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Alert Description",
                field: "AlertDescription",
                width: 200
            }, {
                title: "External Stop",
                field: "ExternalStop",
                width: 80,
                sortable: false,
                template: '# if (ExternalStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Grid Down",
                field: "GridDown",
                width: 80,
                sortable: false,
                template: '# if (GridDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Internal Grid",
                field: "InternalGrid",
                width: 80,
                sortable: false,
                template: '# if (InternalGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Machine Down",
                field: "MachineDown",
                width: 80,
                sortable: false,
                template: '# if (MachineDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "AEbOK",
                field: "AEbOK",
                width: 80,
                sortable: false,
                template: '# if (AEbOK == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Unknown",
                field: "Unknown",
                width: 80,
                sortable: false,
                template: '# if (Unknown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Weather Stop",
                field: "WeatherStop",
                width: 80,
                sortable: false,
                template: '# if (WeatherStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Alarm Ok Time",
                field: "IsAlarmOk",
                width: 80,
                sortable: false,
                template: '# if (IsAlarmOk == true ) { # <span class="glyphicon glyphicon-ok"></span> # } else {# <span class="glyphicon glyphicon-remove"></span> #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            },
        ]
    });
    $('#dataGridAlarmAnomalies').data("kendoGrid").refresh();
}