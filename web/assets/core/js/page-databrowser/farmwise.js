'use strict';

viewModel.DatabrowserFarmWise = new Object();
var dbfs = viewModel.DatabrowserFarmWise;


dbfs.selectedColumn = ko.observableArray([]);
dbfs.unselectedColumn = ko.observableArray([]);
dbfs.ColumnList = ko.observableArray([]);
dbfs.AllProjectColumnList = ko.observableArray([]);
dbfs.ColList = ko.observableArray([]);
dbfs.defaultSelectedColumn = ko.observableArray();


dbfs.getDate = function(date){
    var newvalue = new Date(date);
    var dates = moment(newvalue).utc().format("YYYY-MM-DD HH:mm:ss");

    return dates;
}

dbfs.InitFarmWiseGrid= function() {

    dbr.farmwisevis(true);



    var misc = {
        "tipe": "ScadaHFD",
        "period": fa.period,
    }

    var param = {
        "Custom": {
            "ColumnList": dbr.columnMustHaveFarmWise.concat((dbfs.selectedColumn() == "" ? dbfs.defaultSelectedColumn() : dbfs.selectedColumn()))
        },
        "misc": misc
    };

    var filters = [{
        field: "timestamp",
        operator: "gte",
        value: dbfs.getDate(fa.dateStart)
    }, {
        field: "timestamp",
        operator: "lte",
        value: dbfs.getDate(fa.dateEnd)
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

    var columns = [];
    var gColumns = dbr.columnMustHaveFarmWise.concat(dbfs.selectedColumn());
    if (dbfs.selectedColumn().length == 0) {
        gColumns = dbr.columnMustHaveFarmWise.concat(dbfs.defaultSelectedColumn());
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
            type: (val._id == "projectname") || (val._id == "statedescription") ? "string" : "number",
            width: widthVal,
            headerAttributes: {
                style: "text-align:center"
            },
            attributes: {
                style: "text-align:center"
            },
            locked: (val._id == "projectname" ? true : false),
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
            if (fa.dateEnd - fa.dateStart < 86400000) {
                col["filterable"] = {ui: function(element){
                    element.kendoTimePicker({
                        interval : 10,
                        format : "HH:mm",
                    });
                    element.data("kendoTimePicker").options.max = fa.dateEnd;
                    element.data("kendoTimePicker").options.min = fa.dateStart;
                }};
            }else{
                col["filterable"] = {ui: function(element){
                    element.kendoDatePicker()
                    element.data("kendoDatePicker").max(fa.dateEnd)
                    element.data("kendoDatePicker").min(fa.dateStart)
                }};
            }
        }
        columns.push(col);
    });

    console.log(columns);

    $('#farmWiseGrid').html("");
    $('#farmWiseGrid').kendoGrid({
        dataSource: {
            filter: filters,
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            transport: {
                read: {
                    url: viewModel.appName + "databrowser/getcustomfarmwise",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function(options) {
                    $.each(options.filter.filters, function(key, val){
                        if(val.logic !== undefined){
                            var index = 0;
                            $.each(val.filters, function(i, res){
                                if(res.field == "timestamp"){
                                    if (fa.dateEnd - fa.dateStart < 86400000) {
                                        var waktu = $("form").find("[data-role='timepicker'][data-bind='value:filters["+index+"].value']").val();
                                        if(index > 0){
                                            waktu = $("form").find("[data-role='timepicker'][data-bind='value: filters["+index+"].value']").val();
                                        }
                                        var splitWaktu = waktu.split(":");
                                        res.value = new Date(Date.UTC(moment(fa.dateEnd).get('year'), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), parseInt(splitWaktu[0]), parseInt(splitWaktu[1]), 0, 0));
                                    }else{
                                        if(index == 0){
                                            var tanggal = $("form").find("[data-role='datepicker'][data-bind='value:filters["+index+"].value']").val();
                                            var splitTanggal = tanggal.split("/");
                                            res.value = new Date(Date.UTC(parseInt(splitTanggal[2]), parseInt(splitTanggal[0]-1), parseInt(splitTanggal[1]), 0, 0, 0, 0));
                                        } else {
                                            tanggal = $("form").find("[data-role='datepicker'][data-bind='value: filters["+index+"].value']").val();
                                            splitTanggal = tanggal.split("/");
                                            res.value = new Date(Date.UTC(parseInt(splitTanggal[2]), parseInt(splitTanggal[0]-1), parseInt(splitTanggal[1]), 23, 59, 59, 999));
                                        }
                                    }
                                }
                                index++;
                            })
                        }
                    });
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function(res) {
                    app.loading(false);
                    dbr.farmwisevis(false);
                    dbr.LastFilter = res.data.LastFilter;
                    dbr.LastSort = res.data.LastSort;
                    return res.data.Data
                },
                total: function(res) {

                    if (!app.isFine(res)) {
                        return;
                    }
                    $('#totalfarm').html(kendo.toString(res.data.TotalProject, 'n0'));
                    $('#totaldatafarm').html(kendo.toString(res.data.Total, 'n0'));
                    $('#totalactivepowerfarm').html(kendo.toString(res.data.TotalActivePower / 1000, 'n2') + ' MW');
                    $('#totalprodfarm').html(kendo.toString(res.data.TotalEnergy / 1000, 'n2') + ' MWh');
                    $('#avgwindspeedfarm').html(kendo.toString(res.data.AvgWindSpeed, 'n2') + ' m/s');
                    return res.data.Total;
                },
            },
            sort: [{
                field: 'timestamp',
                dir: 'asc'
            }, 
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
                $("#farmWiseGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","73px");
                $("#farmWiseGrid >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","73px");
            },200);
        }
    });

    var grid = $('#farmWiseGrid').data('kendoGrid');
    var columns = grid.columns;

    $.each(columns, function(i, val) {
        $('#farmWiseGrid').data("kendoGrid").hideColumn(val.field);
    });

    $('#farmWiseGrid').data("kendoGrid").showColumn("timestamp");
    $('#farmWiseGrid').data("kendoGrid").showColumn("projectname");
    if (dbfs.selectedColumn() == "") {
        $.each(dbfs.defaultSelectedColumn(), function(idx, data) {
            $('#farmWiseGrid').data("kendoGrid").showColumn(data._id);
        });
    } else {
        $.each(dbfs.selectedColumn(), function(idx, data) {
            $('#farmWiseGrid').data("kendoGrid").showColumn(data._id);
        });
    }
    $('.k-grid-showHideColumnFarmWise').on("click", function() {
        Data.InitColumnListHFD();
        $("#modalShowHideHFD").modal();
        return false;
    });
    $('#farmWiseGrid').data("kendoGrid").showColumn("projectname");
    $('#farmWiseGrid').data("kendoGrid").showColumn("timestamp");
    $("#farmWiseGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","73px");
    
    $('#farmWiseGrid').data("kendoGrid").refresh();
}

dbfs.selectRowHFD = function() {
    var grid1 = $('#columnListFarmWise').data('kendoGrid');
    var grid2 = $('#selectedListFarmWise').data('kendoGrid');
    dbr.gridMoveTo(grid2, grid1, true);
}

dbfs.unselectRowHFD = function() {
    var grid1 = $('#columnListFarmWise').data('kendoGrid');
    var grid2 = $('#selectedListFarmWise').data('kendoGrid');
    dbr.gridMoveTo(grid1, grid2, true);
}
dbfs.getColumnListHFD = function(){
    var a = dbfs.defaultSelectedColumn();
    var b = dbfs.ColumnList();

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

    dbfs.ColumnList(result);
    // console.log(dbfs.ColumnList());
}
dbfs.showColumnHFD = function() {
    app.loading(true);
    dbfs.selectedColumn([]);
    dbfs.unselectedColumn([]);
    var grid = $('#selectedListFarmWise').data('kendoGrid');
    var dataSources = grid.dataSource.data();
    var selectedList = [];
    var columnList = [];

    $.each(dataSources, function(i, val) {
        selectedList.push(val.id);
        dbfs.selectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source,
            order: val.order,
            projectname: val.projectname,
        });
    });
    dbfs.selectedColumn().sort(function(a, b){
        return b.order < a.order ? 1
        : b.order > a.order ? -1
        : 0;
    });

    $.each($('#columnListFarmWise').data('kendoGrid').dataSource.data(), function(i, val) {
        dbfs.unselectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source,
            order: val.order,
            projectname: val.projectname,
        });
    });

    $.each(dbfs.ColumnList(), function(idx, data) {
        columnList.push(data.id);
    })

    dbfs.InitFarmWiseGrid();
    $('#modalShowHideFarmWise').modal("hide");
}