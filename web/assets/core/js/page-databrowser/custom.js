'use strict';

viewModel.DatabrowserCustom = new Object();
var dbc = viewModel.DatabrowserCustom;

dbc.InitCustomGrid = function() {
    dbr.customvis(true);
   
    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
    }

    var param = {
        "Custom": {
            "ColumnList": (dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn())
        }
    };

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
        value: turbine
    }, ];

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }

    var columns = [];
    var gColumns = dbr.selectedColumn();
    if (dbr.selectedColumn().length == 0) {
        gColumns = dbr.defaultSelectedColumn();
    }

    $.each(gColumns, function(i, val) {
        var col = {
            field: val._id,
            title: val.label,
            type: val._id == "turbine" ? "string" : "number",
            width: 120,
            headerAttributes: {
                style: "text-align:center"
            }
        };

        if (val._id == "timestamp") {
            col = {
                field: val._id,
                title: val.label,
                type: "date",
                width: 140,
                template: "#= kendo.toString(moment.utc(timestamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                value: true
            }
        }
        columns.push(col);
    });

    $('#customGrid').html("");
    $('#customGrid').kendoGrid({
        dataSource: {
            filter: filters,
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getcustomlist",
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
                    dbr.customvis(false);
                    // app.loading(false);
                    return res.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbineCustom').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldataCustom').html(kendo.toString(res.data.Total, 'n0'));
                    $('#totalactivepowerCustom').html(kendo.toString(res.data.TotalActivePower / 1000, 'n0') + ' MWh');
                    $('#totalprodCustom').html(kendo.toString(res.data.TotalEnergy / 1000, 'n0') + ' MWh');
                    $('#avgwindspeedCustom').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
                    return res.data.Total;
                },
            },
            sort: [{
                field: 'TimeStamp',
                dir: 'asc'
            }, {
                field: 'Turbine',
                dir: 'asc'
            }],
        },
        // toolbar: [
        //     "excel", {
        //         text: "Show Hide Columns",
        //         name: "showHideColumn",
        //         imageClass: "fa fa-eye-slash ",
        //     }
        // ],
        excel: {
            fileName: "Custom 10 Minutes Data.xlsx",
            filterable: true,
            allPages: true
        },
        selectable: "multiple",
        reorderable: true,
        groupable: false,
        sortable: true,
        pageable: {
            pageSize: 10,
            input:true, 
        },
        filterable: true,
        scrollable: true,
        columns: columns,
    });

    var grid = $('#customGrid').data('kendoGrid');
    var columns = grid.columns;

    $.each(columns, function(i, val) {
        $('#customGrid').data("kendoGrid").hideColumn(val.field);
    });

    if (dbr.selectedColumn() == "") {
        $.each(dbr.defaultSelectedColumn(), function(idx, data) {
            $('#customGrid').data("kendoGrid").showColumn(data._id);
        });
    } else {
        $.each(dbr.selectedColumn(), function(idx, data) {
            $('#customGrid').data("kendoGrid").showColumn(data._id);
        });
    }
    $('.k-grid-showHideColumn').on("click", function() {
        Data.InitColumnList();
        $("#modalShowHide").modal();
        return false;
    });
}