'use strict';

viewModel.DatabrowserEvent = new Object();
var dbe = viewModel.DatabrowserEvent;

dbe.InitEventGrid = function() {
    dbr.eventrawvis(true);
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));
    dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));

    var turbine = [];
    if ($("#turbineMulti").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineMulti").data("kendoMultiSelect").value();
    }

    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Turbine: turbine,
    };


    $('#EventGrid').html("");
    $('#EventGrid').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/geteventlist",
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
                    // app.loading(false);
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
                class: "align-right"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Alarm Id",
            field: "AlarmId",
            width: 120,
            attributes: {
                class: "align-right"
            },
            format: "{0:n2}",
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
}