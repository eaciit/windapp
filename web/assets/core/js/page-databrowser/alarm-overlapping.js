'use strict';

viewModel.DatabrowserAlarmOverlap = new Object();
var dbao = viewModel.DatabrowserAlarmOverlap;

dbao.InitGridAlarmOverlapping = function() {
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

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }

    var filter = {
        filters: filters
    }
    var param = {
        filter: filter
    };

    $('#dataGridAlarmOverlapping').html("");
    $('#dataGridAlarmOverlapping').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getalarmoverlappinglist",
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
                    return res.data.Data
                },
                total: function(res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbinealarmo').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataalarmo').html(kendo.toString(res.data.Total, 'n0'));
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
        groupable: false,
        sortable: true,
        filterable: false,
        pageable: {
            pageSize: 10,
            input:true, 
        },
        detailInit: dbao.InitOverlapDetail,
        columns: [{
                title: "Date",
                field: "StartDate",
                template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                width: 80
            }, {
                title: "Turbine",
                field: "Turbine",
                width: 90
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
                width: 75,
                attributes: {
                    style: "text-align:center;"
                }
            },
        ]
    });
}

dbao.InitOverlapDetail = function(e) {
    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
    }
    var param = {};

    var filters = [{
        field: "_id",
        operator: "eq",
        value: e.data.ID
    },{
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

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }

    $("<div/>").appendTo(e.detailCell).kendoGrid({
        selectable: "multiple",
        dataSource: {
            serverPaging: false,
            serverSorting: false,
            serverFiltering: true,
            filter: filters,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getalarmoverlappingdetails",
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
                data: function(res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.Data
                },
                total: function(res) {
                    if (!app.isFine(res)) {
                        return;
                    }
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
        scrollable: true,
        sortable: false,
        pageable: {
            pageSize: 10,
            input:true, 
        },
        //resizable: true,
        columns: [{
                title: "Date",
                field: "StartDate",
                template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                headerAttributes: {
                    style: "text-align: center"
                },
                width: 80
            }, {
                title: "Turbine",
                field: "Turbine",
                width: 90,
                sortable: false,
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Start Time",
                field: "StartDate",
                template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 65,
                sortable: false,
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            },
            /*{ title: "Farm", field: "Farm", width: 100 },*/
            {
                title: "End Date",
                field: "EndDate",
                template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                headerAttributes: {
                    style: "text-align: center"
                },
                width: 80,
                sortable: false
            }, {
                title: "End Time",
                field: "EndDate",
                template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 65,
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                },
                sortable: false
            }, {
                title: "Alert Description",
                field: "AlertDescription",
                width: 200,
                headerAttributes: {
                    style: "text-align: center"
                },
                sortable: false
            },
            // { title: "External Stop", field: "ExternalStop", width: 90 , sortable: false, template:"<img src='../res/img/green-dot.png'>", attributes:{style:"text-align:center;"}},
            {
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
                title: "WeatherStop",
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
            },
        ]
    });
}