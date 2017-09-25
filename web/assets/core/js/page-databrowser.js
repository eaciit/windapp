'use strict';

viewModel.Databrowser = new Object();
var dbr = viewModel.Databrowser;

dbr.turbineList = ko.observableArray([]);
dbr.modelList = ko.observableArray([{
    "value": 1,
    "text": "Regen"
}, {
    "value": 2,
    "text": "Suzlon"
}, ]);
dbr.projectList = ko.observableArray([{
    "value": 1,
    "text": "WindFarm-01"
}, {
    "value": 2,
    "text": "WindFarm-02"
}, ]);

dbr.jmrvis = ko.observable(true);
dbr.mettowervis = ko.observable(true);
dbr.oemvis = ko.observable(true);
dbr.hfdvis = ko.observable(true);
dbr.downeventvis = ko.observable(true);
dbr.downeventhfdvis = ko.observable(true);
dbr.customvis = ko.observable(true);
dbr.eventrawvis = ko.observable(true);

dbr.isScadaLoaded = ko.observable(false);
dbr.isScadaHFDLoaded = ko.observable(false);
dbr.isDowntimeEventLoaded = ko.observable(false);
dbr.isCustomLoaded = ko.observable(false);
dbr.isEventLoaded = ko.observable(false);
dbr.isMetLoaded = ko.observable(false);
dbr.isJMRLoaded = ko.observable(false);
dbr.isScadaExceptionTimeDurationLoaded = ko.observable(false);
dbr.isScadaAnomaliesLoaded = ko.observable(false);
dbr.isAlarmOverlappingLoaded = ko.observable(false);
dbr.isAlarmAnomaliesLoaded = ko.observable(false);
dbr.isDowntimeeventhfdLoaded = ko.observable(false);

dbr.gridColumnsScada = ko.observableArray([]);
dbr.gridColumnsScadaException = ko.observableArray([]);
dbr.gridColumnsScadaAnomaly = ko.observableArray([]);
dbr.filterJMR = ko.observableArray([]);
var turbineval = [];
var availDateList = {};
var availDateAll;
dbr.LastFilter;
dbr.LastSort;
dbr.selectedColumn = ko.observableArray([]);
dbr.unselectedColumn = ko.observableArray([]);
dbr.ColumnList = ko.observableArray([]);
dbr.ColList = ko.observableArray([]);
dbr.defaultSelectedColumn = ko.observableArray();
dbr.columnMustHaveHFD = [{
    _id: "timestamp",
    label: "TimeStamp",
    source: "ScadaDataHFD",
}, {
    _id: "turbine",
    label: "Turbine",
    source: "ScadaDataHFD",
}];

dbr.columnMustHaveOEM = [{
    _id: "timestamp",
    label: "TimeStamp",
    source: "ScadaDataOEM",
}, {
    _id: "turbine",
    label: "Turbine",
    source: "ScadaDataOEM",
}];

dbr.ShowHideColumnScada = function(gridID, field, id, index) {
    if ($('#' + id).is(":checked")) {
        $('#' + gridID).data("kendoGrid").showColumn(index);
    } else {
        $('#' + gridID).data("kendoGrid").hideColumn(index);
    }
}

var Data = {
    LoadData: function() {
        setTimeout(function(){
            $(".panel-body").find(".nav-tabs").find("li.active").find('a').trigger( "click" );
        }, 100)
    },
    InitDefault: function() {
        // dbr.getAvailDate();

        Data.LoadData();
    },
    InitColumnList: function() {
        $("#columnList").kendoGrid({
            theme: "flat",
            dataSource: {
                data: (dbr.selectedColumn() == "" ? dbr.ColumnList() : dbr.unselectedColumn()),
                sort: [{
                    field: 'label',
                    dir: 'asc'
                }],
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [{
                field: "label",
                title: "Columns List",
                headerAttributes: {
                    style: "text-align: center"
                }
            }, ],
            change: function(arg) {
                var selected = $.map(this.select(), function(item) {
                    return $(item).find('td').first().text();
                });
                var grid1 = $('#columnList').data('kendoGrid');
                var grid2 = $('#selectedList').data('kendoGrid');
                dbr.gridMoveTo(grid1, grid2, false);
            },
        });

        setTimeout(function() {
            $('#columnList').data('kendoGrid').refresh();
            $('#selectedList').data('kendoGrid').refresh();
        }, 300);

        $("#selectedList").kendoGrid({
            theme: "flat",
            dataSource: {
                data: dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn(),
                sort: [{
                    field: 'label',
                    dir: 'asc'
                }],
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [{
                field: "label",
                title: "Selected Columns",
                headerAttributes: {
                    style: "text-align: center"
                }
            }, ],
            change: function(arg) {
                var selected = $.map(this.select(), function(item) {
                    return $(item).find('td').first().text();
                });
                var grid1 = $('#columnList').data('kendoGrid');
                var grid2 = $('#selectedList').data('kendoGrid');
                dbr.gridMoveTo(grid2, grid1, false);
            },
        });
    },
    InitColumnListHFD: function() {
        $("#columnListHFD").kendoGrid({
            theme: "flat",
            dataSource: {
                data: (dbsh.selectedColumn() == "" ? dbsh.ColumnList() : dbsh.unselectedColumn()),
                sort: [{
                    field: 'label',
                    dir: 'asc'
                }],
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [{
                field: "label",
                title: "Columns List",
                headerAttributes: {
                    style: "text-align: center"
                }
            }, ],
            change: function(arg) {
                var selected = $.map(this.select(), function(item) {
                    return $(item).find('td').first().text();
                });
                var grid1 = $('#columnListHFD').data('kendoGrid');
                var grid2 = $('#selectedListHFD').data('kendoGrid');

                var dataSource = grid2.dataSource;
                var recordsOnCurrentView = dataSource.view().length;
                
                if(recordsOnCurrentView == 30){
                    app.showError("Max. 30 Columns")
                }else{
                    dbr.gridMoveTo(grid1, grid2, false);
                }
            },
        });

        setTimeout(function() {
            $('#columnListHFD').data('kendoGrid').refresh();
            $('#selectedListHFD').data('kendoGrid').refresh();
        }, 300);

        $("#selectedListHFD").kendoGrid({
            theme: "flat",
            dataSource: {
                data: dbsh.selectedColumn() == "" ? dbsh.defaultSelectedColumn() : dbsh.selectedColumn(),
                // sort: [{
                //     field: 'label',
                //     dir: 'asc'
                // }],
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [{
                field: "label",
                title: "Selected Columns",
                headerAttributes: {
                    style: "text-align: center"
                }
            }, ],
            change: function(arg) {
                var selected = $.map(this.select(), function(item) {
                    return $(item).find('td').first().text();
                });
                var grid1 = $('#columnListHFD').data('kendoGrid');
                var grid2 = $('#selectedListHFD').data('kendoGrid');
                dbr.gridMoveTo(grid2, grid1, false);
            },
        });
    }
};

dbr.setAvailableDate = function(isFirst) {

    setTimeout(function(){
        var tabType = $(".panel-body").find(".nav-tabs").find("li.active").attr('id');
        var tipeTab = "MET";
        switch (tabType) {
            case "scadahfdTab" :
                tipeTab = "ScadaDataHFD"
                break;
            case "downtimeeventTab" :
                tipeTab = "EventDown"
                break;
            case "customTab" :
                tipeTab = "MET"
                break;
            case "eventTab" :
                tipeTab = "EventRaw"
                break;
            case "metTab" :
                tipeTab = "MET"
                break;
            case "jmrTab" :
                tipeTab = "JMR"
                break;
            case "downtimeeventhfdTab" :
                tipeTab = "EventDownHFD"
                break;
            default:
                tipeTab = "ScadaDataOEM"
                break;
        }

        var namaproject = $('#projectList').data("kendoDropDownList").value();

        if(namaproject == "") {
            namaproject = "Tejuva";
        }

        var startDate = kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][0]).format('DD-MMM-YYYY'));
        var endDate = kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][1]).format('DD-MMM-YYYY'));

        $('#availabledatestart').html(startDate);
        $('#availabledateend').html(endDate);

        var maxDateData = new Date(availDateAll[namaproject][tipeTab][1]);

        if(moment(maxDateData).get('year') !== 1){
            if(isFirst === true){
                var startDatepicker = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

                $('#dateStart').data('kendoDatePicker').value(startDatepicker);
                $('#dateEnd').data('kendoDatePicker').value(endDate);
            }
        }

    }, 500);
}

dbr.Scada = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaLoaded()) {
        dbr.isScadaLoaded(true);
        dbs.InitScadaGrid();
    } else {
        app.loading(false);
    }
}
dbr.getAvailDate = function(){
    app.ajaxPost(viewModel.appName + "/analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        availDateAll = res.data;
        dbr.setAvailableDate(true);
    });
}

dbr.ScadaHFD = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isScadaHFDLoaded());
    if(!dbr.isScadaHFDLoaded()) {
        dbr.isScadaHFDLoaded(true);
        dbsh.InitScadaHFDGrid();
    } else {
        app.loading(false);
    }
}

dbr.Downtime = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isDowntimeEventLoaded());
    if(!dbr.isDowntimeEventLoaded()) {
        dbr.isDowntimeEventLoaded(true);
        dbd.InitDEgrid();
    } else {
        app.loading(false);
    }
}

dbr.Custom = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isCustomLoaded());
    if(!dbr.isCustomLoaded()) {
        dbr.isCustomLoaded(true);
        dbc.InitCustomGrid();
    } else {
        app.loading(false);
    }
}

dbr.Event = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isEventLoaded());
    if(!dbr.isEventLoaded()) {
        dbr.isEventLoaded(true);
        dbe.InitEventGrid();
    } else {
        app.loading(false);
    }
}

dbr.Met = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(br.isMetLoaded());
    if(!dbr.isMetLoaded()) {
        dbr.isMetLoaded(true);
        dbm.InitMet();
    } else {
        app.loading(false);
    }
}

dbr.JMR = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isJMRLoaded());
    if(!dbr.isJMRLoaded()) {
        dbr.isJMRLoaded(true);
        dbj.InitGridJMR();
    } else {
        app.loading(false);
    }
}

dbr.DowntimeHFD = function(id) {
    fa.LoadData();
    app.loading(true);
    dbr.setAvailableDate(dbr.isDowntimeeventhfdLoaded());
    if(!dbr.isDowntimeeventhfdLoaded()) {
        dbr.isDowntimeeventhfdLoaded(true);
        dbdhfd.InitDEHFDgrid();
    } else {
        app.loading(false);
    }
}

dbr.ScadaException = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaExceptionTimeDurationLoaded()) {
        dbr.isScadaExceptionTimeDurationLoaded(true);
        dbt.InitGridExceptionTimeDuration();
    } else {
        app.loading(false);
    }
}

dbr.ScadaAnomalies = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaAnomaliesLoaded()) {
        dbr.isScadaAnomaliesLoaded(true);
        dbsa.InitGridAnomalies();
    } else {
        app.loading(false);
    }
}

dbr.AlarmOverlapping = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isAlarmOverlappingLoaded()) {
        dbr.isAlarmOverlappingLoaded(true);
        dbao.InitGridAlarmOverlapping();
    } else {
        app.loading(false);
    }
}

dbr.AlarmAnomalies = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isAlarmAnomaliesLoaded()) {
        dbr.isAlarmAnomaliesLoaded(true);
        dbaa.InitGridAlarmAnomalies();
    } else {
        app.loading(false);
    }
}

dbr.ResetFlagLoaded = function() {
    dbr.isScadaLoaded(false);
    dbr.isScadaHFDLoaded(false);
    dbr.isDowntimeEventLoaded(false);
    dbr.isCustomLoaded(false);
    dbr.isEventLoaded(false);
    dbr.isMetLoaded(false);
    dbr.isJMRLoaded(false);
    dbr.isScadaExceptionTimeDurationLoaded(false);
    dbr.isScadaAnomaliesLoaded(false);
    dbr.isAlarmOverlappingLoaded(false);
    dbr.isAlarmAnomaliesLoaded(false);
    dbr.isDowntimeeventhfdLoaded(false);
}

dbr.selectRow = function() {
    var grid1 = $('#columnList').data('kendoGrid');
    var grid2 = $('#selectedList').data('kendoGrid');
    dbr.gridMoveTo(grid2, grid1, true);
}

dbr.unselectRow = function() {
    var grid1 = $('#columnList').data('kendoGrid');
    var grid2 = $('#selectedList').data('kendoGrid');
    dbr.gridMoveTo(grid1, grid2, true);
}

dbr.gridMoveTo = function(from, to, all) {
    if (all == true) {
        from.select(from.tbody.find(">tr"));
    }
    var selected = from.select();

    if (selected.length > 0) {
        var items = [];
        $.each(selected, function(idx, elem) {
            items.push(from.dataItem(elem));
        });
        var fromDS = from.dataSource;
        var toDS = to.dataSource;
        $.each(items, function(idx, elem) {
            toDS.add({
                _id: elem._id,
                label: elem.label,
                source: elem.source,
                order: elem.order,
                projectname: elem.projectname,
            });
            fromDS.remove(elem);
        });
        toDS.sync();
        fromDS.sync();
    }
}

dbr.showColumn = function() {
    dbr.selectedColumn([]);
    dbr.unselectedColumn([]);
    var grid = $('#selectedList').data('kendoGrid');
    var dataSources = grid.dataSource.data();
    var selectedList = [];
    var columnList = [];

    $.each(dataSources, function(i, val) {
        selectedList.push(val.id);
        dbr.selectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source
        });
    });

    $.each($('#columnList').data('kendoGrid').dataSource.data(), function(i, val) {
        dbr.unselectedColumn.push({
            _id: val._id,
            label: val.label,
            source: val.source
        });
    });

    $.each(dbr.ColumnList(), function(idx, data) {
        columnList.push(data.id);
    })

    dbc.InitCustomGrid();

    $('#modalShowHide').modal("hide");
}

function secondsToHms(d) {
    d = Number(d);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var res = (h > 0 ? (h < 10 ? "0" + h : h) : "00") + ":" + (m > 0 ? (m < 10 ? "0" + m : m) : "00") + ":" + (s > 0 ? s : "00")

    return res;
}

function DataBrowserExporttoExcel(functionName) {
    app.loading(true);
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));
    dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));

    var misc = {
        tipe: functionName,
        "period": fa.period,
    }
    var columnList = dbr.columnMustHaveOEM.concat(dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn());
    if (functionName == "ScadaHFDCustom") {
        columnList = dbr.columnMustHaveHFD.concat(dbsh.selectedColumn() == "" ? dbsh.defaultSelectedColumn() : dbsh.selectedColumn());
    }

    var param = {
        Project: fa.project,
        "Custom": {
            "ColumnList": columnList,
        },
        "misc": misc,
        filter: dbr.LastFilter,
        sort: dbr.LastSort,
    };

    var urlName = viewModel.appName + "databrowser/genexceldata";
    if(functionName.toLowerCase().indexOf("custom") >= 0) {
        urlName = viewModel.appName + "databrowser/genexcelcustom10minutes";
    }

    app.ajaxPost(urlName, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        window.location = viewModel.appName + "/".concat(res.data);
        app.loading(false);
    });
}

dbr.exportToExcel = function(idGrid){
    app.loading(true);
    setTimeout(function(){
        $("#"+idGrid).getKendoGrid().saveAsExcel();
        app.loading(false);
    },500);
}

vm.currentMenu('Data Browser');
vm.currentTitle('Data Browser');
vm.isShowDataAvailability(false);
vm.breadcrumb([{
    title: 'Data Browser',
    href: viewModel.appName + 'page/databrowser'
}]);

dbr.ChangeColumnList = function() {
    dbsh.ColumnList([]);
    $.each(dbsh.AllProjectColumnList(), function(idx, val) {
        if(val.projectname === $("#projectList").data("kendoDropDownList").value() ||
            val.source === "MetTower") {
            dbsh.ColumnList.push(val);
        }
    });
    dbsh.defaultSelectedColumn(dbsh.ColumnList().slice(0, 28));
    var tempSelected = dbsh.selectedColumn();
    dbsh.selectedColumn([]);
    $.each(dbsh.AllProjectColumnList(), function(idxAll, valAll) {
        $.each(tempSelected, function(idxSelect, valSelect) {
            if(valAll.projectname === $("#projectList").data("kendoDropDownList").value() && valAll.source === valSelect.source &&
                valAll._id === valSelect._id) {
                dbsh.selectedColumn.push(valAll);
            }
        });
    });

    var tempUnselected = dbsh.unselectedColumn();
    dbsh.unselectedColumn([]);
    $.each(dbsh.AllProjectColumnList(), function(idxAll, valAll) {
        $.each(tempUnselected, function(idxSelect, valSelect) {
            if(valAll.projectname === $("#projectList").data("kendoDropDownList").value() && valAll.source === valSelect.source &&
                valAll._id === valSelect._id) {
                dbsh.unselectedColumn.push(valAll);
            }
        });
    });
    dbsh.getColumnListHFD();
}

$(document).ready(function() {
    // app.loading(true);
    dbr.getAvailDate();
    $("#scadaExceptionbtn").click(function(){
        $("#scadaException").slideToggle("slow");
        $("#scadaException").css("display","inline-table");
    });
    $("#scadaAnomalybtn").click(function(){
        $("#scadaAnomaly").slideToggle("slow");
        $("#scadaAnomaly").css("display","inline-table");
    });

    $('.k-grid-showHideColumn').on("click", function() {
        $("#modalShowHide").modal();

        $("#myModal").on('shown.bs.modal', function() {
            Data.InitColumnList();
        });
        return false;
    });

    $('.k-grid-showHideColumnHFD').on("click", function() {
        $("#modalShowHideHFD").modal();

        $("#modalShowHideHFD").on('shown.bs.modal', function() {
            Data.InitColumnListHFD();
        });
        return false;
    });
    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        dbr.ResetFlagLoaded();
        Data.LoadData();
    });

    setTimeout(function() {
        dbsh.ColumnList([]);
        $.each(dbsh.AllProjectColumnList(), function(idx, val) {
            if(val.projectname === $("#projectList").data("kendoDropDownList").value() || 
                val.source === "MetTower") {
                dbsh.ColumnList.push(val);
            }
        });
        dbr.defaultSelectedColumn(dbr.ColumnList().slice(0, 28));
        dbsh.defaultSelectedColumn(dbsh.ColumnList().slice(0, 28));
        fa.checkTurbine();
        Data.InitDefault();
        dbc.getColumnCustom();
        dbsh.getColumnListHFD();
        $("#projectList").on("change", function(event) { 
             dbr.ChangeColumnList();
        });
    }, 1000);

    $('#projectList').kendoDropDownList({
		change: function () {  
            dbr.ResetFlagLoaded();
            dbr.getAvailDate();
			var project = $('#projectList').data("kendoDropDownList").value();
			fa.populateTurbine(project);
		}
	});
});
