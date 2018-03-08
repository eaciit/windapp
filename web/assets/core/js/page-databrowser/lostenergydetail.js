'use strict';

viewModel.DatabrowserLostEnergyDetail = new Object();
var dbled = viewModel.DatabrowserLostEnergyDetail;

dbled.InitLEDgrid = function() {
    dbr.lostenergydetail(true);

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
        field: "startdate",
        operator: "$gte",
        value: fa.dateStart
    }, {
        field: "startdate",
        operator: "$lte",
        value: fa.dateEnd
    }, {
        field: "turbine",
        operator: "$in",
        value: fa.turbine()
    }, ];

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "$eq",
            value: fa.project
        })
    }

    $('#LEDgrid').html("");
    $('#LEDgrid').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            filter: filters,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getlostenergydetail",
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
                    dbr.lostenergydetail(false);
                    app.isFine(ress);
                    dbr.LastFilter = ress.data.LastFilter;
                    dbr.LastSort = ress.data.LastSort;
                    return ress.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbineLED').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataLED').html(kendo.toString(res.data.Total, 'n0'));

                    return res.data.Total;
                }
            },
            sort: [{
                field: 'startdate',
                dir: 'asc'
            }, {
                field: 'detail.startdate',
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
                title: "Turbine",
                field: "turbinename",
                attributes: {
                    class: "align-center"
                },
                width: 70,
                filterable: false,
                sortable: false,
            }, {
                title: "Time Start",
                field: "detail.startdate",
                template: "#= kendo.toString(moment.utc(detail.startdate).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 100,
                filterable: false,
                sortable: false,
            }, {
                title: "Time End",
                field: "detail.enddate",
                template: "#= kendo.toString(moment.utc(detail.enddate).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 100,
                filterable: false,
                sortable: false,
            }, {
                title: "Duration (hh:mm:ss)",
                field: "detail.duration",
                template: '#= kendo.toString(secondsToHmsDatabrowser(detail.duration)) #',
                width: 90,
                attributes: {
                    class: "align-center"
                },
                filterable: false,
                sortable: false,
            }, {
                title: "Power Lost",
                field: "detail.powerlost",
                width: 90,
                sortable: false,
                format: "{0:n2}",
                attributes: {
                    class: "align-right"
                },
            }, 
            {
                title: "Alarm Description",
                field: "alertdescription",
                width: 190,
                filterable: false,
                sortable: false,
            },{
                title: "Reduce Availability",
                field: "reduceavailability",
                width: 70,
                template: '# if (reduceavailability == true ) { # Yes # } else {# No #}#',
                attributes: {
                    class: "align-center"
                },
                filterable: false,
                sortable: false,
            },
        ]
    });
    $('#LEDgrid').data("kendoGrid").refresh();
}