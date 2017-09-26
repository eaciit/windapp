'use strict';

viewModel.DatabrowserScadaHFD = new Object();
var dbsh = viewModel.DatabrowserScadaHFD;


dbsh.selectedColumn = ko.observableArray([]);
dbsh.unselectedColumn = ko.observableArray([]);
dbsh.ColumnList = ko.observableArray([]);
dbsh.AllProjectColumnList = ko.observableArray([]);
dbsh.ColList = ko.observableArray([]);
dbsh.defaultSelectedColumn = ko.observableArray();

dbsh.InitScadaHFDGrid= function() {
    dbr.hfdvis(true);
    // var turbine = [];
    // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    //     turbine = turbineval;
    // } else {
    //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // }

    var misc = {
        "tipe": "ScadaHFD",
        "period": fa.period,
    }

    var param = {
        "Custom": {
            "ColumnList": dbr.columnMustHaveHFD.concat((dbsh.selectedColumn() == "" ? dbsh.defaultSelectedColumn() : dbsh.selectedColumn()))
        },
        "misc": misc
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
        value: fa.turbine()
    }, ];

    if(fa.project != "") {
        filters.push({
            field: "projectname",
            operator: "eq",
            value: fa.project
        })
    }

    // var columns = [
    //         { title: "Time Stamp", field: "TimeStamp", template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #", width: 130, locked: true, filterable: false },
    //         { title: "Turbine", field: "Turbine", attributes: { class: "align-center" }, width: 90, locked: true, filterable: false },
    // ];

    var columns = [];
    var gColumns = dbr.columnMustHaveHFD.concat(dbsh.selectedColumn());
    if (dbsh.selectedColumn().length == 0) {
        gColumns = dbr.columnMustHaveHFD.concat(dbsh.defaultSelectedColumn());
    }
    
    var widthVal = 90;
    var lowerLabel = "";
    $.each(gColumns, function(i, val) {
        lowerLabel = val.label.toLowerCase();
        if(lowerLabel.indexOf("direction") >= 0 ||
            lowerLabel.indexOf("react") >= 0) {
            widthVal = 100;
        } else if(lowerLabel.indexOf("generator") >= 0 ||
            lowerLabel.indexOf("frequency") >= 0 ||
            lowerLabel.indexOf("ambient") >= 0 ||
            lowerLabel.indexOf("pressure") >= 0 ||
            lowerLabel.indexOf("gearbox") >= 0) {
            widthVal = 110;
        } else {
            widthVal = 90;
        }
        var col = {
            field: val._id,
            title: val.label,
            type: (val._id == "turbine") || (val._id == "statedescription") ? "string" : "number",
            width: widthVal,
            headerAttributes: {
                style: "text-align:center"
            },
            attributes: {
                style: "text-align:center"
            },
            locked: (val._id == "turbine" ? true : false),
        };

        if (val._id == "timestamp") {
            col = {
                field: val._id,
                title: val.label,
                type: "date",
                width: 130,
                template: "#= kendo.toString(moment.utc(timestamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                value: true,
                locked: true,
            }
        }
        columns.push(col);
    });

    $('#scadahfdGrid').html("");
    $('#scadahfdGrid').kendoGrid({
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
                    app.loading(false);
                    dbr.hfdvis(false);
                    dbr.LastFilter = res.data.LastFilter;
                    dbr.LastSort = res.data.LastSort;
                    return res.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalturbinehfd').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                    $('#totaldatahfd').html(kendo.toString(res.data.Total, 'n0'));
                    $('#totalactivepowerhfd').html(kendo.toString(res.data.TotalActivePower / 1000, 'n0') + ' MWh');
                    $('#totalprodhfd').html(kendo.toString(res.data.TotalEnergy / 1000, 'n0') + ' MWh');
                    $('#avgwindspeedhfd').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
                    return res.data.Total;
                },
            },
            sort: [{
                field: 'timestamp',
                dir: 'asc'
            }, 
            // {
            //     field: 'turbine',
            //     dir: 'asc'
            // }
            ],
        },
        excel: {
            fileName: "Scada HFD.xlsx",
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
        dataBound: function(){
            setTimeout(function(){
                $("#scadahfdGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","73px");
                $("#scadahfdGrid >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","73px");
            },200);
        }
    });

    var grid = $('#scadahfdGrid').data('kendoGrid');
    var columns = grid.columns;

    $.each(columns, function(i, val) {
        $('#scadahfdGrid').data("kendoGrid").hideColumn(val.field);
    });

    $('#scadahfdGrid').data("kendoGrid").showColumn("timestamp");
    $('#scadahfdGrid').data("kendoGrid").showColumn("turbine");
    if (dbsh.selectedColumn() == "") {
        $.each(dbsh.defaultSelectedColumn(), function(idx, data) {
            $('#scadahfdGrid').data("kendoGrid").showColumn(data._id);
        });
    } else {
        $.each(dbsh.selectedColumn(), function(idx, data) {
            $('#scadahfdGrid').data("kendoGrid").showColumn(data._id);
        });
    }
    $('.k-grid-showHideColumnHFD').on("click", function() {
        Data.InitColumnListHFD();
        $("#modalShowHideHFD").modal();
        return false;
    });
    $('#scadahfdGrid').data("kendoGrid").showColumn("turbine");
    $('#scadahfdGrid').data("kendoGrid").showColumn("timestamp");
    $("#scadahfdGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","73px");
    
    $('#scadahfdGrid').data("kendoGrid").refresh();
}

dbsh.selectRowHFD = function() {
    var grid1 = $('#columnListHFD').data('kendoGrid');
    var grid2 = $('#selectedListHFD').data('kendoGrid');
    dbr.gridMoveTo(grid2, grid1, true);
}

dbsh.unselectRowHFD = function() {
    var grid1 = $('#columnListHFD').data('kendoGrid');
    var grid2 = $('#selectedListHFD').data('kendoGrid');
    dbr.gridMoveTo(grid1, grid2, true);
}
dbsh.getColumnListHFD = function(){
    var a = dbsh.defaultSelectedColumn();
    var b = dbsh.ColumnList();

    var onlyInA = a.filter(function(current){
        return b.filter(function(current_b){
            return current_b._id == current._id && current_b.label == current.label
        }).length == 0
    });

    var onlyInB = b.filter(function(current){
        return a.filter(function(current_a){
            return current_a._id == current._id && current_a.label == current.label
        }).length == 0
    });

    var result = onlyInA.concat(onlyInB);

    dbsh.ColumnList(result);
    // console.log(dbsh.ColumnList());
}
dbsh.showColumnHFD = function() {
    app.loading(true);
    dbsh.selectedColumn([]);
    dbsh.unselectedColumn([]);
    var grid = $('#selectedListHFD').data('kendoGrid');
    var dataSources = grid.dataSource.data();
    var selectedList = [];
    var columnList = [];

    $.each(dataSources, function(i, val) {
        selectedList.push(val.id);
        dbsh.selectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source,
            order: val.order,
            projectname: val.projectname,
        });
    });
    dbsh.selectedColumn().sort(function(a, b){
        return b.order < a.order ? 1
        : b.order > a.order ? -1
        : 0;
    });

    $.each($('#columnListHFD').data('kendoGrid').dataSource.data(), function(i, val) {
        dbsh.unselectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source,
            order: val.order,
            projectname: val.projectname,
        });
    });

    $.each(dbsh.ColumnList(), function(idx, data) {
        columnList.push(data.id);
    })

    dbsh.InitScadaHFDGrid();
    $('#modalShowHideHFD').modal("hide");
}