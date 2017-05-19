'use strict';

viewModel.DatabrowserEvent = new Object();
var dbe = viewModel.DatabrowserEvent;

dbe.InitEventGrid = function() {
    dbr.eventrawvis(true);

    // var turbine = [];
    // // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    // //     turbine = turbineval;
    // // } else {
    // //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // // }

    var filters = [{
        field: "timestamp",
        operator: "gte",
        value: fa.dateStart
    }, {
        field: "timestamp",
        operator: "lte",
        value: fa.dateEnd
    }, {
        field: "turbine",
        operator: "in",
        value: fa.turbine
    }, ];

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }


    $('#EventGrid').html("");
    $('#EventGrid').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            filter: filters,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/geteventlist",
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
                data: function(res) {
                    app.loading(false);
                    dbr.eventrawvis(false);
                    return res.data.Data
                },
                total: function(res) {
                    if (!app.isFine(res)) {
                        return;
                    }

                    $('#totalturbineEvent').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataEvent').html(kendo.toString(res.data.Total, 'n0'));

                    return res.data.Total;
                }
            },
            sort: [{
                field: 'TimeStamp',
                dir: 'asc'
            }, {
                field: 'Turbine',
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
            title: "Time Stamp",
            field: "TimeStamp",
            template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
            width: 130,
            filterable: false
        }, {
            title: "Project Name",
            field: "ProjectName",
            attributes: {
                class: "align-center"
            },
            width: 90,
            filterable: false
        }, {
            title: "Turbine",
            field: "Turbine",
            attributes: {
                class: "align-center"
            },
            width: 90,
            filterable: false
        }, {
            title: "Event Type",
            field: "EventType",
            attributes: {
                class: "align-center"
            },
            width: 100,
            filterable: false
        }, {
            title: "Alarm Description",
            field: "AlarmDescription",
            attributes: {
                class: "align-center"
            },
            width: 150,
            filterable: false
        }, {
            title: "Turbine Status",
            field: "TurbineStatus",
            attributes: {
                class: "align-center"
            },
            width: 120,
            filterable: false
        }, {
            title: "Brake Type",
            field: "BrakeType",
            attributes: {
                class: "align-center"
            },
            width: 150,
            filterable: false
        }, {
            title: "Brake Program",
            field: "BrakeProgram",
            width: 120,
            attributes: {
                class: "align-center"
            },
            format: "{0}",
            filterable: false
        }, {
            title: "Alarm Id",
            field: "AlarmId",
            width: 120,
            attributes: {
                class: "align-center"
            },
            format: "{0}",
            filterable: false
        }, {
            title: "Alarm Toggle",
            field: "AlarmToggle",
            width: 120,
            sortable: false,
            template: '# if (AlarmToggle == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
            headerAttributes: {
                style: "text-align: center"
            },
            attributes: {
                style: "text-align:center;"
            }
        }, ]
    });
    $('#EventGrid').data("kendoGrid").refresh();
}