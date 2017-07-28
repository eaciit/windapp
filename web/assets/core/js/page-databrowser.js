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

dbr.selectedColumn = ko.observableArray([]);
dbr.unselectedColumn = ko.observableArray([]);
dbr.ColumnList = ko.observableArray([]);
dbr.ColList = ko.observableArray([]);
dbr.defaultSelectedColumn = ko.observableArray([{
    "_id": "timestamp",
    "label": "Time Stamp",
    "source": "ScadaDataOEM"
  },{
    "_id": "turbine",
    "label": "Turbine",
    "source": "ScadaDataOEM"
  },{
    "_id": "ai_intern_r_pidangleout",
    "label": "Ai Intern R Pid Angle Out",
    "source": "ScadaDataOEM"
}, {
    "_id": "ai_intern_activpower",
    "label": "Ai Intern Active Power",
    "source": "ScadaDataOEM"
}, {
    "_id": "ai_intern_i1",
    "label": "Ai Intern I1",
    "source": "ScadaDataOEM"
}, {
    "_id": "ai_intern_i2",
    "label": "Ai Intern I2",
    "source": "ScadaDataOEM"
}, ]);

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
            // var idParent = $(".panel-body").find(".nav-tabs").find(".active")[0].id;
            $(".panel-body").find(".nav-tabs").find("li.active").find('a').trigger( "click" );
            // switch (idParent) {
            //     case 'mainTab':
            //         $("#Main").find(".nav-tabs").find("li.active").find('a').trigger( "click" );
            //         break;
            //     case 'exceptionTab':
            //         $("#Exception").find(".nav-tabs").find("li.active").find('a').trigger( "click" );
            //         break;
            // }
        }, 100)
    },
    InitDefault: function() {
        var maxDateData = new Date(app.getUTCDate(app.currentDateData));
        var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
        var lastEndDate = new Date(app.getDateMax(maxDateData));

        $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
        $('#dateStart').data('kendoDatePicker').value(lastStartDate);

        
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
                var grid1 = $('#columnListHFD').data('kendoGrid');
                var grid2 = $('#selectedListHFD').data('kendoGrid');
                dbr.gridMoveTo(grid2, grid1, false);
            },
        });
    }
};

dbr.setAvailableDate = function() {
    setTimeout(function(){
        var tabType = $(".panel-body").find(".nav-tabs").find("li.active").attr('id');
        var tipeTab = "ScadaDataOEM";
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
        var namaproject = fa.project;
        if(namaproject == "") {
            namaproject = "Tejuva";
        }
        $('#availabledatestart').html(kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][0]).format('DD-MMMM-YYYY')));
        $('#availabledateend').html(kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][1]).format('DD-MMMM-YYYY')));
    }, 300);
}

dbr.Scada = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaLoaded()) {
        dbr.isScadaLoaded(true);
        dbs.InitScadaGrid();
        app.ajaxPost(viewModel.appName + "/analyticlossanalysis/getavaildateall", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }
            availDateAll = res.data;
            dbr.setAvailableDate();
        });
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.ScadaHFD = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaHFDLoaded()) {
        dbr.isScadaHFDLoaded(true);
        dbsh.InitScadaHFDGrid();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.Downtime = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isDowntimeEventLoaded()) {
        dbr.isDowntimeEventLoaded(true);
        dbd.InitDEgrid();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.Custom = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isCustomLoaded()) {
        dbr.isCustomLoaded(true);
        dbc.InitCustomGrid();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.Event = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isEventLoaded()) {
        dbr.isEventLoaded(true);
        dbe.InitEventGrid();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.Met = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isMetLoaded()) {
        dbr.isMetLoaded(true);
        dbm.InitMet();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.JMR = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isJMRLoaded()) {
        dbr.isJMRLoaded(true);
        dbj.InitGridJMR();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.DowntimeHFD = function(id) {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isDowntimeeventhfdLoaded()) {
        dbr.isDowntimeeventhfdLoaded(true);
        dbdhfd.InitDEHFDgrid();
        dbr.setAvailableDate();
    } else {
        dbr.setAvailableDate();
        app.loading(false);
    }
}

dbr.ScadaException = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaExceptionTimeDurationLoaded()) {
        dbr.isScadaExceptionTimeDurationLoaded(true);
        dbt.InitGridExceptionTimeDuration();
        app.ajaxPost(viewModel.appName + "/databrowser/getdurationavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //Duration Data
            if (res.data.Duration.length == 0) {
                res.data.Duration = [];
            } else {
                if (res.data.Duration.length > 0) {
                    var minDatetemp = new Date(res.data.Duration[0]);
                    var maxDatetemp = new Date(res.data.Duration[1]);
                    availDateList.availabledatestartduration = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                    availDateList.availabledateendduration = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                    $('#availabledatestart').html(availDateList.availabledatestartduration);
                    $('#availabledateend').html(availDateList.availabledateendduration);
                }
            }
        });
    } else {
        $('#availabledatestart').html(availDateList.availabledatestartduration);
        $('#availabledateend').html(availDateList.availabledateendduration);
        app.loading(false);
    }
}

dbr.ScadaAnomalies = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isScadaAnomaliesLoaded()) {
        dbr.isScadaAnomaliesLoaded(true);
        dbsa.InitGridAnomalies();
        app.ajaxPost(viewModel.appName + "/databrowser/getscadaanomalyavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //Scada Anomaly Data
            if (res.data.ScadaAnomaly.length == 0) {
                res.data.ScadaAnomaly = [];
            } else {
                if (res.data.ScadaAnomaly.length > 0) {
                    var minDatetemp = new Date(res.data.ScadaAnomaly[0]);
                    var maxDatetemp = new Date(res.data.ScadaAnomaly[1]);
                    availDateList.availabledatestartscadaanomaly = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                    availDateList.availabledateendscadaanomaly = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                    $('#availabledatestart').html(availDateList.availabledatestartscadaanomaly);
                    $('#availabledateend').html(availDateList.availabledateendscadaanomaly);
                }
            }
        });
    } else {
        $('#availabledatestart').html(availDateList.availabledatestartscadaanomaly);
        $('#availabledateend').html(availDateList.availabledateendscadaanomaly);
        app.loading(false);
    }
}

dbr.AlarmOverlapping = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isAlarmOverlappingLoaded()) {
        dbr.isAlarmOverlappingLoaded(true);
        dbao.InitGridAlarmOverlapping();
        app.ajaxPost(viewModel.appName + "/databrowser/getalarmoverlappingavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //AlarmOverlapping Data
            if (res.data.AlarmOverlapping.length == 0) {
                res.data.AlarmOverlapping = [];
            } else {
                if (res.data.AlarmOverlapping.length > 0) {
                    var minDatetemp = new Date(res.data.AlarmOverlapping[0]);
                    var maxDatetemp = new Date(res.data.AlarmOverlapping[1]);

                    availDateList.availabledatestartalarmoverlapping = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                    availDateList.availabledateendalarmoverlapping = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                    $('#availabledatestart').html(availDateList.availabledatestartalarmoverlapping);
                    $('#availabledateend').html(availDateList.availabledateendalarmoverlapping);
                }
            }
        });
    } else {
        $('#availabledatestart').html(availDateList.availabledatestartalarmoverlapping);
        $('#availabledateend').html(availDateList.availabledateendalarmoverlapping);
        app.loading(false);
    }
}

dbr.AlarmAnomalies = function() {
    fa.LoadData();
    app.loading(true);
    if(!dbr.isAlarmAnomaliesLoaded()) {
        dbr.isAlarmAnomaliesLoaded(true);
        dbaa.InitGridAlarmAnomalies();
        app.ajaxPost(viewModel.appName + "/databrowser/getalarmscadaanomalyavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //AlarmScadaAnomaly Data
            if (res.data.AlarmScadaAnomaly.length == 0) {
                res.data.AlarmScadaAnomaly = [];
            } else {
                if (res.data.AlarmScadaAnomaly.length > 0) {
                    var minDatetemp = new Date(res.data.AlarmScadaAnomaly[0]);
                    var maxDatetemp = new Date(res.data.AlarmScadaAnomaly[1]);
                    availDateList.availabledatestartalarmscadaanomaly = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                    availDateList.availabledateendalarmscadaanomaly = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                    $('#availabledatestart').html(availDateList.availabledatestartalarmscadaanomaly);
                    $('#availabledateend').html(availDateList.availabledateendalarmscadaanomaly);
                }
            }
        });
    } else {
        $('#availabledatestart').html(availDateList.availabledatestartalarmscadaanomaly);
        $('#availabledateend').html(availDateList.availabledateendalarmscadaanomaly);
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
                source: elem.source
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

    // var turbine = [];
    // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    //     turbine = turbineval;
    // } else {
    //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // }

    var Filter = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Turbine: fa.turbine(),
        Project: fa.project,
        Misc: {tipe: functionName}
    };
    var urlName = viewModel.appName + "databrowser/genexceldata";
    if(functionName === "genexcelcustom10minutes") {
        urlName = viewModel.appName + "databrowser/" + functionName;
        Filter["Custom"] = {"ColumnList": (dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn())}
    }

    app.ajaxPost(urlName, Filter, function(res) {
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
vm.breadcrumb([{
    title: 'Data Browser',
    href: viewModel.appName + 'page/databrowser'
}]);

$(document).ready(function() {
    // app.loading(true);

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
        fa.checkTurbine();
        Data.InitDefault();
        dbc.getColumnCustom();
        dbsh.getColumnListHFD();
    }, 1000);

    $('#projectList').kendoDropDownList({
		change: function () {  
			var project = $('#projectList').data("kendoDropDownList").value();
			fa.populateTurbine(project);
		}
	});
});
