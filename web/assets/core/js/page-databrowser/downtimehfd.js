'use strict';

viewModel.DatabrowserDowntimeHFD = new Object();
var dbdhfd = viewModel.DatabrowserDowntimeHFD;

dbdhfd.InitDEHFDgrid = function() {
    dbr.downeventhfdvis(true);

    // var turbine = [];
    // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    //     turbine = turbineval;
    // } else {
    //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // }

    var misc = {
        "tipe": "eventdownhfd",
        "needtotalturbine": true,
        "period": fa.period,
    }

    var param = {"misc": misc}
    
    var filters = [{
        field: "isdeleted",
        operator: "eq",
        value: false
    },{
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
        value: fa.turbine()
    }, ];

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }

    $('#DEHFDgrid').html("");
    $('#DEHFDgrid').kendoGrid({
        dataSource: {
            serverSorting: true,
            serverFiltering: true,
            filter: filters,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getdatabrowserlist",
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
                data: function(ress) {
                    app.loading(false);
                    dbr.downeventhfdvis(false);
                    app.isFine(ress);
                    dbr.LastFilter = ress.data.LastFilter;
                    dbr.LastSort = ress.data.LastSort;
                    return ress.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbineDEHFD').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataDEHFD').html(kendo.toString(res.data.Total, 'n0'));

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
        filterable: {
            extra: false,
            operators: {
                string: {
                    contains: "Contains",
                    eq: "Is equal to"
                },
            }
        },
        columns: [{
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 70,
                filterable: false
            }, {
                title: "Time Start",
                field: "TimeStart",
                template: "#= kendo.toString(moment.utc(TimeStart).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 100,
                filterable: false
            }, {
                title: "Time End",
                field: "TimeEnd",
                template: "#= (moment.utc(TimeEnd).format('DD-MM-YYYY') == '01-01-0001'?'Not yet finished' : kendo.toString(moment.utc(TimeEnd).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss')) #",
                width: 100,
                filterable: false
            }, {
                title: "Duration (hh:mm:ss)",
                field: "Duration",
                template: '#= kendo.toString(secondsToHmsDatabrowser(Duration)) #',
                width: 90,
                attributes: {
                    class: "align-center"
                },
                filterable: false
            }, {
                title: "BreakDown Group",
                field: "BDGroup",
                width: 90,
                sortable: false,
                filterable: false
            }, 
            // {
            //     title: "Turbine State",
            //     field: "TurbineState",
            //     width: 70,
            //     attributes: {
            //         class: "align-center"
            //     },
            //     filterable: false
            // }, {
            //     title: "Alarm Code",
            //     field: "AlarmCode",
            //     width: 70,
            //      attributes: {
            //         class: "align-center"
            //     },
            //     filterable: false
            // }, 
            {
                title: "Alarm Description",
                field: "AlarmDesc",
                width: 190,
            },{
                title: "Reduce Availability",
                field: "ReduceAvailability",
                width: 70,
                template: '# if (ReduceAvailability == true ) { # Yes # } else {# No #}#',
                attributes: {
                    class: "align-center"
                },
                filterable: false
            },
        ]
    });
    $('#DEHFDgrid').data("kendoGrid").refresh();
}