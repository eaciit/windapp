'use strict';

viewModel.DatabrowserDowntime = new Object();
var dbd = viewModel.DatabrowserDowntime;

dbd.InitDEgrid = function() {
    dbr.downeventvis(true);

    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
    }
    
    var filters = [{
        field: "timestart",
        operator: "gte",
        value: fa.dateStart
    }, {
        field: "timestart",
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

    $('#DEgrid').html("");
    $('#DEgrid').kendoGrid({
        dataSource: {
            serverSorting: true,
            serverFiltering: true,
            filter: filters,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getdowntimeeventlist",
                    type: "POST",
                    data: {},
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function(ress) {
                    // app.loading(false);
                    dbr.downeventvis(false);
                    app.isFine(ress);
                    return ress.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbineDE').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataDE').html(kendo.toString(res.data.Total, 'n0'));

                    return res.data.Total;
                }
            },
            sort: [{
                field: 'timestart',
                dir: 'asc'
            }, {
                field: 'turbine',
                dir: 'asc'
            }],
        },
        selectable: "multiple",
        groupable: false,
        sortable: true,
        pageable: {
            pageSize: 10,
            input:true, 
        },
        columns: [{
                title: "Time Start",
                field: "TimeStart",
                template: "#= kendo.toString(moment.utc(TimeStart).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 100,
                filterable: false
            },

            {
                title: "Time End",
                field: "TimeEnd",
                template: "#= kendo.toString(moment.utc(TimeEnd).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 100
            }, {
                title: "Grid Down",
                field: "DownGrid",
                width: 80,
                sortable: false,
                template: '# if (DownGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Machine Down",
                field: "DownMachine",
                width: 80,
                sortable: false,
                template: '# if (DownMachine == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Environment Down",
                field: "DownEnvironment",
                width: 80,
                sortable: false,
                template: '# if (DownEnvironment == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, {
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                filterable: false
            }, {
                title: "Alarm Description",
                field: "AlarmDescription",
                width: 100,
                filterable: false
            }, {
                title: "Duration (hh:mm:ss)",
                field: "Duration",
                template: '#= kendo.toString(secondsToHms(Duration)) #',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                filterable: false
            },


        ]
    });
}