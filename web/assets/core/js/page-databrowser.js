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
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
        fa.LoadData();
        dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));
        dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));

        if (dateStart - dateEnd > 25200000) {
            toolkit.showError("Invalid Date Range Selection");
            return;
        } else {
            var idParent = $(".panel-body").find(".nav-tabs").find(".active")[0].id
            switch (idParent) {
                case 'mainTab':
                    var idTab = $("#Main").find(".nav-tabs").find(".active")[0].id
                    Data.IdCheck(idTab);
                    break;
                case 'exceptionTab':
                    var idTab = $("#Exception").find(".nav-tabs").find(".active")[0].id
                    Data.IdCheck(idTab);
                    break;
            }
        }
    },
    IdCheck: function(param) {
        switch(param) {
            case 'scadaTab':
                if(!dbr.isScadaLoaded()) {
                    app.loading(true);
                    dbr.isScadaLoaded(true);
                    dbs.InitScadaGrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/getscadadataoemavaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        //Scada Data
                        if (res.data.ScadaDataOEM.length == 0) {
                            res.data.ScadaDataOEM = [];
                        } else {
                            if (res.data.ScadaDataOEM.length > 0) {
                                var minDatetemp = new Date(res.data.ScadaDataOEM[0]);
                                var maxDatetemp = new Date(res.data.ScadaDataOEM[1]);
                                availDateList.availabledatestartscadadataoem = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendscadadataoem = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartscadadataoem);
                                $('#availabledateend').html(availDateList.availabledateendscadadataoem);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartscadadataoem);
                    $('#availabledateend').html(availDateList.availabledateendscadadataoem);
                }
                break;
            case 'scadahfdTab':
                if(!dbr.isScadaHFDLoaded()) {
                    app.loading(true);
                    dbr.isScadaHFDLoaded(true);
                    dbsh.InitScadaHFDGrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/getscadadatahfdavaildate", {}, function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        //Scada Data HFD
                        if (res.data.ScadaDataHFD.length == 0) {
                            res.data.ScadaDataHFD = [];
                        } else {
                            if (res.data.ScadaDataHFD.length > 0) {
                                var minDatetemp = new Date(res.data.ScadaDataHFD[0]);
                                var maxDatetemp = new Date(res.data.ScadaDataHFD[1]);
                                availDateList.availabledatestartscadadatahfd = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendscadadatahfd = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartscadadatahfd);
                                $('#availabledateend').html(availDateList.availabledateendscadadatahfd);
                            }
                        }         
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartscadadatahfd);
                    $('#availabledateend').html(availDateList.availabledateendscadadatahfd);
                }
                break;
            case 'downtimeeventTab':
                if(!dbr.isDowntimeEventLoaded()) {
                    app.loading(true);
                    dbr.isDowntimeEventLoaded(true);
                    dbd.InitDEgrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/getdowntimeeventvaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        //Scada Data
                        if (res.data.DowntimeEvent.length == 0) {
                            res.data.DowntimeEvent = [];
                        } else {
                            if (res.data.DowntimeEvent.length > 0) {
                                var minDatetemp = new Date(res.data.DowntimeEvent[0]);
                                var maxDatetemp = new Date(res.data.DowntimeEvent[1]);
                                availDateList.availabledatestartDE = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendDE = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartDE);
                                $('#availabledateend').html(availDateList.availabledateendDE);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartDE);
                    $('#availabledateend').html(availDateList.availabledateendDE);
                }
                break;
            case 'customTab':
                if(!dbr.isCustomLoaded()) {
                    app.loading(true);
                    dbr.isCustomLoaded(true);
                    dbc.InitCustomGrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/getcustomavaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        if (res.data.CustomDate.length == 0) {
                            res.data.CustomDate = [];
                        } else {
                            if (res.data.CustomDate.length > 0) {
                                var arrDate = res.data.CustomDate.sort();
                                var minDatetemp = new Date(arrDate[0]);
                                var maxDatetemp = new Date(arrDate[3]);
                                availDateList.availabledatestartCustom = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendCustom = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartCustom);
                                $('#availabledateend').html(availDateList.availabledateendCustom);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartCustom);
                    $('#availabledateend').html(availDateList.availabledateendCustom);
                }
                break;
            case 'eventTab':
                if(!dbr.isEventLoaded()) {
                    app.loading(true);
                    dbr.isEventLoaded(true);
                    dbe.InitEventGrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/geteventavaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        //EventDate Data
                        if (res.data.EventDate.length == 0) {
                            res.data.EventDate = [];
                        } else {
                            if (res.data.EventDate.length > 0) {
                                var minDatetemp = new Date(res.data.EventDate[0]);
                                var maxDatetemp = new Date(res.data.EventDate[1]);
                                availDateList.availabledatestarteventraw = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendeventraw = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestarteventraw);
                                $('#availabledateend').html(availDateList.availabledateendeventraw);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestarteventraw);
                    $('#availabledateend').html(availDateList.availabledateendeventraw);
                }
                break;
            case 'metTab':
                if(!dbr.isMetLoaded()) {
                    app.loading(true);
                    dbr.isMetLoaded(true);
                    dbm.InitMet();
                    app.ajaxPost(viewModel.appName + "/databrowser/getmetavaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        //MET Tower Data
                        if (res.data.MET.length == 0) {
                            res.data.MET = [];
                        } else {
                            if (res.data.MET.length > 0) {
                                var minDatetemp = new Date(res.data.MET[0]);
                                var maxDatetemp = new Date(res.data.MET[1]);
                                availDateList.availabledatestartmet = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendmet = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartmet);
                                $('#availabledateend').html(availDateList.availabledateendmet);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartmet);
                    $('#availabledateend').html(availDateList.availabledateendmet);
                }
                break;
            case 'jmrTab':
                if(!dbr.isJMRLoaded()) {
                    app.loading(true);
                    dbr.isJMRLoaded(true);
                    dbj.InitGridJMR();
                    app.ajaxPost(viewModel.appName + "/databrowser/getjmravaildate", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        //JMR Data
                        if (res.data.JMR.length == 0) {
                            res.data.JMR = [];
                        } else {
                            if (res.data.JMR.length > 0) {
                                var minDatetemp = new Date(res.data.JMR[0]);
                                var maxDatetemp = new Date(res.data.JMR[1]);
                                availDateList.availabledatestartjmr = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendjmr = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartjmr);
                                $('#availabledateend').html(availDateList.availabledateendjmr);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartjmr);
                    $('#availabledateend').html(availDateList.availabledateendjmr);
                }
                break;
            case 'scadaExceptionTimeDurationTab':
                if(!dbr.isScadaExceptionTimeDurationLoaded()) {
                    app.loading(true);
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
                }
                break;
            case 'scadaAnomaliesTab':
                if(!dbr.isScadaAnomaliesLoaded()) {
                    app.loading(true);
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
                }
                break;
            case 'alarmOverlappingTab':
                if(!dbr.isAlarmOverlappingLoaded()) {
                    app.loading(true);
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
                }
                break;
            case 'alarmAnomaliesTab':
                if(!dbr.isAlarmAnomaliesLoaded()) {
                    app.loading(true);
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
                }
                break;
            case 'downtimeeventhfdTab':
                if(!dbr.isDowntimeeventhfdLoaded()) {
                    app.loading(true);
                    dbr.isDowntimeeventhfdLoaded(true);
                    dbdhfd.InitDEHFDgrid();
                    app.ajaxPost(viewModel.appName + "/databrowser/getdowntimeeventvaildatehfd", {}, function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        //Scada Data
                        if (res.data.DowntimeEvent.length == 0) {
                            res.data.DowntimeEvent = [];
                        } else {
                            if (res.data.DowntimeEvent.length > 0) {
                                var minDatetemp = new Date(res.data.DowntimeEvent[0]);
                                var maxDatetemp = new Date(res.data.DowntimeEvent[1]);
                                availDateList.availabledatestartDEHFD = kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY'));
                                availDateList.availabledateendDEHFD = kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY'));
                                $('#availabledatestart').html(availDateList.availabledatestartDEHFD);
                                $('#availabledateend').html(availDateList.availabledateendDEHFD);
                            }
                        }
                    });
                } else {
                    $('#availabledatestart').html(availDateList.availabledatestartDEHFD);
                    $('#availabledateend').html(availDateList.availabledateendDEHFD);
                }
                break;
        }
    },
    RefreshGrid: function(param) {
        setTimeout(function() {
            switch (param) {
                case 'Main':
                    var idTab = $("#Main").find(".nav-tabs").find(".active")[0].id
                    Data.IdCheck(idTab);
                    break;
                case 'Exception':
                    var idTab = $("#Exception").find(".nav-tabs").find(".active")[0].id
                    Data.IdCheck(idTab);
                    break;
                default:
                    Data.IdCheck(param);
            }
        }, 5);
    },

    InitDefault: function() {
        var maxDateData = new Date(app.getUTCDate(app.currentDateData));
        var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
        var lastEndDate = new Date(app.toUTC(maxDateData));

        $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
        $('#dateStart').data('kendoDatePicker').value(lastStartDate);

        setTimeout(function() {
            Data.LoadData();
        }, 500);
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

    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
    }

    var Filter = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Turbine: turbine,
    };

    app.ajaxPost(viewModel.appName + "databrowser/" + functionName, Filter, function(res) {
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
        $("#"+idGrid).getKendoGrid().saveAsExcel()
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

        Data.LoadData();
    });

    setTimeout(function() {
        Data.InitDefault();
        // dbc.getColumnCustom();
        dbsh.getColumnListHFD();
    }, 1000);
});