'use strict';

viewModel.DatabrowserScadaHFD = new Object();
var dbsh = viewModel.DatabrowserScadaHFD;


dbsh.selectedColumn = ko.observableArray([]);
dbsh.unselectedColumn = ko.observableArray([]);
dbsh.ColumnList = ko.observableArray([]);
dbsh.ColList = ko.observableArray([]);

dbsh.defaultSelectedColumn = ko.observableArray([
 {
    "_id": "Fast_ActivePowerOutPWCSell_kW",
    "label": "Fast ActivePowerOutPWCSell kW",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_ActivePower_kW",
    "label": "Fast ActivePower kW",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_DrTrVibValue",
    "label": "Fast DrTrVibValue",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_Frequency_Hz",
    "label": "Fast Frequency Hz",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_GenSpeed_RPM",
    "label": "Fast GenSpeed RPM",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_PitchAngle",
    "label": "Fast PitchAngle",
    "source": "ScadaDataHFD"
  }, {
    "_id": "Fast_PitchAngle1",
    "label": "Fast PitchAngle1",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_PitchAngle2",
    "label": "Fast PitchAngle2",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_PitchAngle3",
    "label": "Fast PitchAngle3",
    "source": "ScadaDataHFD"
  },{
    "_id": "Fast_PitchSpeed1",
    "label": "Fast PitchSpeed1",
    "source": "ScadaDataHFD"
  },  {
    "_id": "Fast_RotorSpeed_RPM",
    "label": "Fast RotorSpeed RPM",
    "source": "ScadaDataHFD"
  },  {
    "_id": "Fast_WindSpeed_ms",
    "label": "Fast WindSpeed ms",
    "source": "ScadaDataHFD"
  },  {
    "_id": "Fast_YawAngle",
    "label": "Fast YawAngle",
    "source": "ScadaDataHFD"
  },  {
    "_id": "Fast_YawService",
    "label": "Fast YawService",
    "source": "ScadaDataHFD"
  },  {
    "_id": "Slow_NacelleDrill",
    "label": "Slow NacelleDrill",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_NacellePos",
    "label": "Slow NacellePos",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempG1L1",
    "label": "Slow TempG1L1",
    "source": "ScadaDataHFD"
  }, {
    "_id": "Slow_TempG1L2",
    "label": "Slow TempG1L2",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempG1L3",
    "label": "Slow TempG1L3",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGearBoxHSSDE",
    "label": "Slow TempGearBoxHSSDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGearBoxHSSNDE",
    "label": "Slow TempGearBoxHSSNDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGearBoxIMSDE",
    "label": "Slow TempGearBoxIMSDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGearBoxIMSNDE",
    "label": "Slow TempGearBoxIMSNDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGearBoxOilSump",
    "label": "Slow TempGearBoxOilSump",
    "source": "ScadaDataHFD"
  }, {
    "_id": "Slow_TempGeneratorBearingDE",
    "label": "Slow TempGeneratorBearingDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempGeneratorBearingNDE",
    "label": "Slow TempGeneratorBearingNDE",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempHubBearing",
    "label": "Slow TempHubBearing",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempNacelle",
    "label": "Slow TempNacelle",
    "source": "ScadaDataHFD"
  },{
    "_id": "Slow_TempOutdoor",
    "label": "Slow TempOutdoor",
    "source": "ScadaDataHFD"
  }, {
    "_id": "Slow_WindDirection",
    "label": "Slow WindDirection",
    "source": "ScadaDataHFD"
  }]);

dbsh.InitScadaHFDGrid= function() {
    dbr.hfdvis(true);
    // var turbine = [];
    // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    //     turbine = turbineval;
    // } else {
    //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // }

    var misc = {
        "tipe": "scadahfd",
        "needtotalturbine": true,
        "period": fa.period,
    }
    var param = {
        "Custom": {
            "ColumnList": (dbsh.selectedColumn() == "" ? dbsh.defaultSelectedColumn() : dbsh.selectedColumn())
        },
        "misc": misc
    };

    var filters = [{
        field: "TimeStamp",
        operator: "gte",
        value: fa.dateStart
    }, {
        field: "TimeStamp",
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

    var columns = [
            { title: "Time Stamp", field: "TimeStamp", template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #", width: 130, locked: true, filterable: false },
            { title: "Turbine", field: "Turbine", attributes: { class: "align-center" }, width: 90, locked: true, filterable: false },
    ];

    var gColumns = dbsh.selectedColumn();
    if (dbsh.selectedColumn().length == 0) {
        gColumns = dbsh.defaultSelectedColumn();
    }

    $.each(gColumns, function(i, val) {
        var col = {
            field: val._id,
            title: val.label,
            width: 120,
            headerAttributes: {
                style: "text-align:center"
            },
            attributes: {
                style: "text-align:center"
            },
            template: "#=kendo.toString("+val._id+", 'n2')#"
        };

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
                    $('#totalactivepowerhfd').html(kendo.toString(res.data.TotalActivePower / 1000, 'n2') + ' MWh');
                    $('#totalprodhfd').html(kendo.toString(res.data.TotalEnergy / 1000, 'n2') + ' MWh');
                    $('#avgwindspeedhfd').html(kendo.toString(res.data.AvgWindSpeed, 'n1') + ' m/s');
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
    $('#scadahfdGrid').data("kendoGrid").showColumn("Turbine");
    $('#scadahfdGrid').data("kendoGrid").showColumn("TimeStamp");
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
            source: val.source
        });
    });

    $.each($('#columnListHFD').data('kendoGrid').dataSource.data(), function(i, val) {
        dbsh.unselectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source
        });
    });

    $.each(dbsh.ColumnList(), function(idx, data) {
        columnList.push(data.id);
    })

    dbsh.InitScadaHFDGrid();
    $('#modalShowHideHFD').modal("hide");
}