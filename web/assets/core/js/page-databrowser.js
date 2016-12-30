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
dbr.customvis = ko.observable(true);
dbr.eventrawvis = ko.observable(true);

dbr.gridColumnsScada = ko.observableArray([]);
dbr.gridColumnsScadaException = ko.observableArray([]);
dbr.gridColumnsScadaAnomaly = ko.observableArray([]);
dbr.filterJMR = ko.observableArray([]);
var turbineval = [];

dbr.selectedColumn = ko.observableArray([]);
dbr.unselectedColumn = ko.observableArray([]);
dbr.ColumnList = ko.observableArray([]);
dbr.ColList = ko.observableArray([]);
dbr.defaultSelectedColumn = ko.observableArray([{
    "_id": "timestamp",
    "label": "Time Stamp",
    "source": "ScadaDataOEM"
}, {
    "_id": "turbine",
    "label": "Turbine",
    "source": "ScadaDataOEM"
}, {
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


dbr.populateTurbine = function() {
    app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0) {
            res.data = [];
            dbr.turbineList([{
                value: "",
                text: ""
            }]);
        } else {
            var datavalue = [];
            if (res.data.length > 0) {
                var allturbine = {}
                $.each(res.data, function(key, val) {
                    turbineval.push(val);
                });
                allturbine.value = "All Turbine";
                allturbine.text = "All Turbines";
                datavalue.push(allturbine);
                $.each(res.data, function(key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    datavalue.push(data);
                });
            }
            dbr.turbineList(datavalue);
        }
        setTimeout(function() {
            $('#turbineMulti').data('kendoMultiSelect').value(["All Turbine"])
        }, 300);
    });
};

dbr.checkTurbine = function() {
    var arr = $('#turbineMulti').data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbine");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#turbineMulti').data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $("#turbineMulti").data("kendoMultiSelect").value(["All Turbine"]);
    } else if (arr.length == 0) {
        $("#turbineMulti").data("kendoMultiSelect").value(["All Turbine"]);
    }
}

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

        dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));
        dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));

        if ($("#turbineMulti").data("kendoMultiSelect").value() == "") {
            $('#turbineMulti').data('kendoMultiSelect').value(["All Turbine"])
        }

        if ((new Date(dateStart).getTime() > new Date(dateEnd).getTime())) {
            toolkit.showError("Invalid Date Range Selection");
            return;
        } else {
            app.loading(true);

            // MAIN
            dbs.InitScadaGrid();
            dbsh.InitScadaHFDGrid();
            dbd.InitDEgrid();
            dbc.InitCustomGrid();
            dbe.InitEventGrid();
            dbm.InitMet();
            dbj.InitGridJMR();

            // Exception
            dbt.InitGridExceptionTimeDuration();
            dbsa.InitGridAnomalies();
            dbao.InitGridAlarmOverlapping();
            dbaa.InitGridAlarmAnomalies();
        }
    },
    LoadAvailDate: function() {
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
                    $('#availabledatestartCustom').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendCustom').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartscadadataoem').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendscadadataoem').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });

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
                    $('#availabledatestartscadadatahfd').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendscadadatahfd').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }         
        });


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
                    $('#availabledatestartDE').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendDE').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
        app.ajaxPost(viewModel.appName + "/databrowser/getscadaavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //Scada Data
            if (res.data.ScadaData.length == 0) {
                res.data.ScadaData = [];
            } else {
                if (res.data.ScadaData.length > 0) {
                    var minDatetemp = new Date(res.data.ScadaData[0]);
                    var maxDatetemp = new Date(res.data.ScadaData[1]);
                    $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
        app.ajaxPost(viewModel.appName + "/databrowser/getalarmavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            //Alarm Data
            if (res.data.Alarm.length == 0) {
                res.data.Alarm = [];
            } else {
                if (res.data.Alarm.length > 0) {
                    var minDatetemp = new Date(res.data.Alarm[0]);
                    var maxDatetemp = new Date(res.data.Alarm[1]);
                    $('#availabledatestartalarm').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendalarm').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartjmr').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendjmr').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartmet').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendmet').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartduration').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendduration').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartscadaanomaly').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendscadaanomaly').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartalarmoverlapping').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendalarmoverlapping').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
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
                    $('#availabledatestartalarmscadaanomaly').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#availabledateendalarmscadaanomaly').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });

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
                    $('#eventdatestart').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                    $('#eventdateend').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
                }
            }
        });
    },
    RefreshGrid: function() {
        setTimeout(function() {
            // MAIN 

            $('#scadaGrid').data('kendoGrid').refresh();
            $('#scadahfdGrid').data('kendoGrid').refresh();
            $('#DEgrid').data('kendoGrid').refresh();
            $('#customGrid').data('kendoGrid').refresh();
            $('#EventGrid').data('kendoGrid').refresh();
            $('#dataGridJMR').data('kendoGrid').refresh();
            $('#dataGridMet').data('kendoGrid').refresh();

            // EXCEPTION
            $('#dataGridExceptionTimeDuration').data('kendoGrid').refresh();
            $('#dataGridAnomalies').data('kendoGrid').refresh();
            $('#dataGridAlarmOverlapping').data('kendoGrid').refresh();
            $('#dataGridAlarmAnomalies').data('kendoGrid').refresh();


        }, 5);
    },

    // INIT GRID MAIN

    // INIT GRID EXCEPTION
    InitOverlapDetail: function(e) {
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
        var param = {};

        $("<div/>").appendTo(e.detailCell).kendoGrid({
            selectable: "multiple",
            dataSource: {
                serverPaging: false,
                serverSorting: false,
                serverFiltering: true,
                filter: [{
                    field: "_id",
                    operator: "eq",
                    value: e.data.ID
                }, {
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
                }],
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getalarmoverlappingdetails",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function(options) {
                        return JSON.stringify(options);
                    }
                },
                schema: {
                    data: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
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
            scrollable: true,
            sortable: false,
            pageable: {
                pageSize: 10,
                input:true, 
            },
            //resizable: true,
            columns: [{
                    title: "Date",
                    field: "StartDate",
                    template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    width: 80
                }, {
                    title: "Turbine",
                    field: "Turbine",
                    width: 90,
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Start Time",
                    field: "StartDate",
                    template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                    width: 65,
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                },
                /*{ title: "Farm", field: "Farm", width: 100 },*/
                {
                    title: "End Date",
                    field: "EndDate",
                    template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    width: 80,
                    sortable: false
                }, {
                    title: "End Time",
                    field: "EndDate",
                    template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #",
                    width: 65,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    },
                    sortable: false
                }, {
                    title: "Alert Description",
                    field: "AlertDescription",
                    width: 200,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    sortable: false
                },
                // { title: "External Stop", field: "ExternalStop", width: 90 , sortable: false, template:"<img src='../res/img/green-dot.png'>", attributes:{style:"text-align:center;"}},
                {
                    title: "External Stop",
                    field: "ExternalStop",
                    width: 80,
                    sortable: false,
                    template: '# if (ExternalStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Grid Down",
                    field: "GridDown",
                    width: 80,
                    sortable: false,
                    template: '# if (GridDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Internal Grid",
                    field: "InternalGrid",
                    width: 80,
                    sortable: false,
                    template: '# if (InternalGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Machine Down",
                    field: "MachineDown",
                    width: 80,
                    sortable: false,
                    template: '# if (MachineDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "AEbOK",
                    field: "AEbOK",
                    width: 80,
                    sortable: false,
                    template: '# if (AEbOK == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Unknown",
                    field: "Unknown",
                    width: 80,
                    sortable: false,
                    template: '# if (Unknown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "WeatherStop",
                    field: "WeatherStop",
                    width: 80,
                    sortable: false,
                    template: '# if (WeatherStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                },
            ]
        });
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
    if ($("#turbineMulti").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineMulti").data("kendoMultiSelect").value();
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
    app.loading(true);
    dbr.populateTurbine();

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
    $('#btnRefresh').on('click', function() {
        Data.LoadData();
    });

    setTimeout(function() {
        Data.InitDefault();
        dbc.InitCustomGrid();
    }, 1000);
    Data.LoadAvailDate();
});