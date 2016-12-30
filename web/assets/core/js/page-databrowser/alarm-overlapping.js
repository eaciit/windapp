'use strict';

viewModel.DatabrowserAlarmOverlap = new Object();
var dbao = viewModel.DatabrowserAlarmOverlap;

dbao.InitGridAlarmOverlapping = function() {
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

    var filters = [{
        field: "startdate",
        operator: "gte",
        value: dateStart
    }, {
        field: "startdate",
        operator: "lte",
        value: dateEnd
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
        detailInit: Data.InitOverlapDetail,
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