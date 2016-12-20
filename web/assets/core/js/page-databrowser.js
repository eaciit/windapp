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
            this.InitScadaGrid();
            this.InitScadaHFDGrid();
            this.InitDEgrid();
            this.InitCustomGrid();
            this.InitEventGrid();
            this.InitMet();
            this.InitGridJMR();

            // Exception
            this.InitGridExceptionTimeDuration();
            this.InitGridAnomalies();
            this.InitGridAlarmOverlapping();
            this.InitGridAlarmAnomalies();
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
            $('#dataGridExceptionTimeDuration').data('kendoGrid').refresh();
            $('#dataGridAnomalies').data('kendoGrid').refresh();
            $('#dataGridAlarmOverlapping').data('kendoGrid').refresh();
            $('#dataGridAlarmAnomalies').data('kendoGrid').refresh();

            $('#dataGridJMR').data('kendoGrid').refresh();
            $('#dataGridMet').data('kendoGrid').refresh();

            $('#scadaGrid').data('kendoGrid').refresh();
            $('#customGrid').data('kendoGrid').refresh();
            $('#DEgrid').data('kendoGrid').refresh();
            $('#EventGrid').data('kendoGrid').refresh();
        }, 5);
    },

    // INIT GRID MAIN
        InitScadaHFDGrid: function () {
        dbr.hfdvis(true);
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

        $('#scadahfdGrid').html("");
        $('#scadahfdGrid').kendoGrid({
            dataSource: {
                filter: [
                    { field: "timestamp", operator: "gte", value: dateStart },
                    { field: "timestamp", operator: "lte", value: dateEnd },
                    { field: "turbine", operator: "in", value: turbine }
                ],
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getscadahfdlist",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function (options) {
                        dbr.hfdvis(true);
                        return JSON.stringify(options);
                    }
                },
                pageSize: 10,
                schema: {
                    data: function (res) {
                        app.loading(false);
                        dbr.hfdvis(false);
                        if (!app.isFine(res)) {
                            return;
                        }
                        console.log(res);
                        return res.data.Data
                    },
                    total: function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        $('#totalturbinehfd').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldatahfd').html(kendo.toString(res.data.Total, 'n0'));
                        // $('#totalactivepowerhfd').html(kendo.toString(res.data.TotalActivePower / 1000, 'n0') + ' MWh');
                        // $('#totalprodhfd').html(kendo.toString(res.data.TotalEnergy / 1000, 'n0') + ' MWh');
                        // $('#avgwindspeedhfd').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
                        return res.data.Total;
                    }
                },
                sort: [
                    { field: 'TimeStamp', dir: 'asc' },
                    { field: 'Turbine', dir: 'asc' }
                ],
            },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            pageable: true,
            columns: [
                { title: "Time Stamp", field: "TimeStamp", template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #", width: 130, locked: true, filterable: false },
                { title: "Turbine", field: "Turbine", attributes: { class: "align-center" }, width: 90, locked: true, filterable: false },
                { title: "Fast ActivePower kW", field: "Fast_ActivePower_kW", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePower kW StdDev", field: "Fast_ActivePower_kW_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePower kW Min", field: "Fast_ActivePower_kW_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePower kW Max", field: "Fast_ActivePower_kW_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePower kW Count", field: "Fast_ActivePower_kW_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast WindSpeed ms", field: "Fast_WindSpeed_ms", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast WindSpeed ms StdDev", field: "Fast_WindSpeed_ms_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast WindSpeed ms Min", field: "Fast_WindSpeed_ms_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast WindSpeed ms Max", field: "Fast_WindSpeed_ms_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast WindSpeed ms Count", field: "Fast_WindSpeed_ms_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacellePos", field: "Slow_NacellePos", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacellePos StdDev", field: "Slow_NacellePos_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacellePos Min", field: "Slow_NacellePos_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacellePos Max", field: "Slow_NacellePos_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacellePos Count", field: "Slow_NacellePos_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow WindDirection", field: "Slow_WindDirection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow WindDirection StdDev", field: "Slow_WindDirection_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow WindDirection Min", field: "Slow_WindDirection_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow WindDirection Max", field: "Slow_WindDirection_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow WindDirection Count", field: "Slow_WindDirection_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL3", field: "Fast_CurrentL3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL3 StdDev", field: "Fast_CurrentL3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL3 Min", field: "Fast_CurrentL3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL3 Max", field: "Fast_CurrentL3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL3 Count", field: "Fast_CurrentL3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL1", field: "Fast_CurrentL1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL1 StdDev", field: "Fast_CurrentL1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL1 Min", field: "Fast_CurrentL1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL1 Max", field: "Fast_CurrentL1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL1 Count", field: "Fast_CurrentL1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerSetpoint kW", field: "Fast_ActivePowerSetpoint_kW", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerSetpoint kW StdDev", field: "Fast_ActivePowerSetpoint_kW_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerSetpoint kW Min", field: "Fast_ActivePowerSetpoint_kW_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerSetpoint kW Max", field: "Fast_ActivePowerSetpoint_kW_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerSetpoint kW Count", field: "Fast_ActivePowerSetpoint_kW_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL2", field: "Fast_CurrentL2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL2 StdDev", field: "Fast_CurrentL2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL2 Min", field: "Fast_CurrentL2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL2 Max", field: "Fast_CurrentL2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast CurrentL2 Count", field: "Fast_CurrentL2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast DrTrVibValue", field: "Fast_DrTrVibValue", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast DrTrVibValue StdDev", field: "Fast_DrTrVibValue_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast DrTrVibValue Min", field: "Fast_DrTrVibValue_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast DrTrVibValue Max", field: "Fast_DrTrVibValue_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast GenSpeed RPM", field: "Fast_GenSpeed_RPM", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast GenSpeed RPM StdDev", field: "Fast_GenSpeed_RPM_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast GenSpeed RPM Min", field: "Fast_GenSpeed_RPM_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast GenSpeed RPM Max", field: "Fast_GenSpeed_RPM_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast GenSpeed RPM Count", field: "Fast_GenSpeed_RPM_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV1", field: "Fast_PitchAccuV1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV1 StdDev", field: "Fast_PitchAccuV1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV1 Min", field: "Fast_PitchAccuV1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV1 Max", field: "Fast_PitchAccuV1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV1 Count", field: "Fast_PitchAccuV1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle", field: "Fast_PitchAngle", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle StdDev", field: "Fast_PitchAngle_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle Min", field: "Fast_PitchAngle_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle Max", field: "Fast_PitchAngle_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle Count", field: "Fast_PitchAngle_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle3", field: "Fast_PitchAngle3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle3 StdDev", field: "Fast_PitchAngle3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle3 Min", field: "Fast_PitchAngle3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle3 Max", field: "Fast_PitchAngle3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle3 Count", field: "Fast_PitchAngle3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle2", field: "Fast_PitchAngle2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle2 StdDev", field: "Fast_PitchAngle2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle2 Min", field: "Fast_PitchAngle2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle2 Max", field: "Fast_PitchAngle2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle2 Count", field: "Fast_PitchAngle2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent1", field: "Fast_PitchConvCurrent1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent1 StdDev", field: "Fast_PitchConvCurrent1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent1 Min", field: "Fast_PitchConvCurrent1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent1 Max", field: "Fast_PitchConvCurrent1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent1 Count", field: "Fast_PitchConvCurrent1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent3", field: "Fast_PitchConvCurrent3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent3 StdDev", field: "Fast_PitchConvCurrent3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent3 Min", field: "Fast_PitchConvCurrent3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent3 Max", field: "Fast_PitchConvCurrent3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent3 Count", field: "Fast_PitchConvCurrent3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent2", field: "Fast_PitchConvCurrent2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent2 StdDev", field: "Fast_PitchConvCurrent2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent2 Min", field: "Fast_PitchConvCurrent2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent2 Max", field: "Fast_PitchConvCurrent2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchConvCurrent2 Count", field: "Fast_PitchConvCurrent2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PowerFactor", field: "Fast_PowerFactor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PowerFactor StdDev", field: "Fast_PowerFactor_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PowerFactor Min", field: "Fast_PowerFactor_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PowerFactor Max", field: "Fast_PowerFactor_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PowerFactor Count", field: "Fast_PowerFactor_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePowerSetpointPPC kVA", field: "Fast_ReactivePowerSetpointPPC_kVA", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePowerSetpointPPC kVA StdDev", field: "Fast_ReactivePowerSetpointPPC_kVA_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePowerSetpointPPC kVA Min", field: "Fast_ReactivePowerSetpointPPC_kVA_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePowerSetpointPPC kVA Max", field: "Fast_ReactivePowerSetpointPPC_kVA_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePowerSetpointPPC kVA Count", field: "Fast_ReactivePowerSetpointPPC_kVA_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePower kVAr", field: "Fast_ReactivePower_kVAr", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePower kVAr StdDev", field: "Fast_ReactivePower_kVAr_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePower kVAr Min", field: "Fast_ReactivePower_kVAr_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePower kVAr Max", field: "Fast_ReactivePower_kVAr_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ReactivePower kVAr Count", field: "Fast_ReactivePower_kVAr_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast RotorSpeed RPM", field: "Fast_RotorSpeed_RPM", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast RotorSpeed RPM StdDev", field: "Fast_RotorSpeed_RPM_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast RotorSpeed RPM Min", field: "Fast_RotorSpeed_RPM_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast RotorSpeed RPM Max", field: "Fast_RotorSpeed_RPM_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast RotorSpeed RPM Count", field: "Fast_RotorSpeed_RPM_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL1", field: "Fast_VoltageL1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL1 StdDev", field: "Fast_VoltageL1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL1 Min", field: "Fast_VoltageL1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL1 Max", field: "Fast_VoltageL1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL1 Count", field: "Fast_VoltageL1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL2", field: "Fast_VoltageL2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL2 StdDev", field: "Fast_VoltageL2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL2 Min", field: "Fast_VoltageL2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL2 Max", field: "Fast_VoltageL2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL2 Count", field: "Fast_VoltageL2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitiveReactPwr kVAr", field: "Slow_CapableCapacitiveReactPwr_kVAr", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitiveReactPwr kVAr StdDev", field: "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitiveReactPwr kVAr Min", field: "Slow_CapableCapacitiveReactPwr_kVAr_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitiveReactPwr kVAr Max", field: "Slow_CapableCapacitiveReactPwr_kVAr_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitiveReactPwr kVAr Count", field: "Slow_CapableCapacitiveReactPwr_kVAr_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductiveReactPwr kVAr", field: "Slow_CapableInductiveReactPwr_kVAr", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductiveReactPwr kVAr StdDev", field: "Slow_CapableInductiveReactPwr_kVAr_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductiveReactPwr kVAr Min", field: "Slow_CapableInductiveReactPwr_kVAr_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductiveReactPwr kVAr Max", field: "Slow_CapableInductiveReactPwr_kVAr_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductiveReactPwr kVAr Count", field: "Slow_CapableInductiveReactPwr_kVAr_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow DateTime Sec", field: "Slow_DateTime_Sec", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow DateTime Sec StdDev", field: "Slow_DateTime_Sec_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow DateTime Sec Min", field: "Slow_DateTime_Sec_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow DateTime Sec Max", field: "Slow_DateTime_Sec_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow DateTime Sec Count", field: "Slow_DateTime_Sec_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle1", field: "Fast_PitchAngle1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle1 StdDev", field: "Fast_PitchAngle1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle1 Min", field: "Fast_PitchAngle1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle1 Max", field: "Fast_PitchAngle1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAngle1 Count", field: "Fast_PitchAngle1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL3", field: "Fast_VoltageL3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL3 StdDev", field: "Fast_VoltageL3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL3 Min", field: "Fast_VoltageL3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL3 Max", field: "Fast_VoltageL3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast VoltageL3 Count", field: "Fast_VoltageL3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitivePwrFactor", field: "Slow_CapableCapacitivePwrFactor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitivePwrFactor StdDev", field: "Slow_CapableCapacitivePwrFactor_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitivePwrFactor Min", field: "Slow_CapableCapacitivePwrFactor_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitivePwrFactor Max", field: "Slow_CapableCapacitivePwrFactor_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableCapacitivePwrFactor Count", field: "Slow_CapableCapacitivePwrFactor_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Production kWh", field: "Fast_Total_Production_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Production kWh StdDev", field: "Fast_Total_Production_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Production kWh Min", field: "Fast_Total_Production_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Production kWh Max", field: "Fast_Total_Production_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Production kWh Count", field: "Fast_Total_Production_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Day kWh", field: "Fast_Total_Prod_Day_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Day kWh StdDev", field: "Fast_Total_Prod_Day_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Day kWh Min", field: "Fast_Total_Prod_Day_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Day kWh Max", field: "Fast_Total_Prod_Day_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Day kWh Count", field: "Fast_Total_Prod_Day_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Month kWh", field: "Fast_Total_Prod_Month_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Month kWh StdDev", field: "Fast_Total_Prod_Month_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Month kWh Min", field: "Fast_Total_Prod_Month_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Month kWh Max", field: "Fast_Total_Prod_Month_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Prod Month kWh Count", field: "Fast_Total_Prod_Month_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerOutPWCSell kW", field: "Fast_ActivePowerOutPWCSell_kW", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerOutPWCSell kW StdDev", field: "Fast_ActivePowerOutPWCSell_kW_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerOutPWCSell kW Min", field: "Fast_ActivePowerOutPWCSell_kW_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerOutPWCSell kW Max", field: "Fast_ActivePowerOutPWCSell_kW_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerOutPWCSell kW Count", field: "Fast_ActivePowerOutPWCSell_kW_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Frequency Hz", field: "Fast_Frequency_Hz", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Frequency Hz StdDev", field: "Fast_Frequency_Hz_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Frequency Hz Min", field: "Fast_Frequency_Hz_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Frequency Hz Max", field: "Fast_Frequency_Hz_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Frequency Hz Count", field: "Fast_Frequency_Hz_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L2", field: "Slow_TempG1L2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L2 StdDev", field: "Slow_TempG1L2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L2 Min", field: "Slow_TempG1L2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L2 Max", field: "Slow_TempG1L2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L2 Count", field: "Slow_TempG1L2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L3", field: "Slow_TempG1L3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L3 StdDev", field: "Slow_TempG1L3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L3 Min", field: "Slow_TempG1L3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L3 Max", field: "Slow_TempG1L3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L3 Count", field: "Slow_TempG1L3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSDE", field: "Slow_TempGearBoxHSSDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSDE StdDev", field: "Slow_TempGearBoxHSSDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSDE Min", field: "Slow_TempGearBoxHSSDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSDE Max", field: "Slow_TempGearBoxHSSDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSDE Count", field: "Slow_TempGearBoxHSSDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSNDE", field: "Slow_TempGearBoxIMSNDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSNDE StdDev", field: "Slow_TempGearBoxIMSNDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSNDE Min", field: "Slow_TempGearBoxIMSNDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSNDE Max", field: "Slow_TempGearBoxIMSNDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSNDE Count", field: "Slow_TempGearBoxIMSNDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempOutdoor", field: "Slow_TempOutdoor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempOutdoor StdDev", field: "Slow_TempOutdoor_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempOutdoor Min", field: "Slow_TempOutdoor_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempOutdoor Max", field: "Slow_TempOutdoor_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempOutdoor Count", field: "Slow_TempOutdoor_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV3", field: "Fast_PitchAccuV3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV3 StdDev", field: "Fast_PitchAccuV3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV3 Min", field: "Fast_PitchAccuV3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV3 Max", field: "Fast_PitchAccuV3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV3 Count", field: "Fast_PitchAccuV3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineActiveHours", field: "Slow_TotalTurbineActiveHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineActiveHours StdDev", field: "Slow_TotalTurbineActiveHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineActiveHours Min", field: "Slow_TotalTurbineActiveHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineActiveHours Max", field: "Slow_TotalTurbineActiveHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineActiveHours Count", field: "Slow_TotalTurbineActiveHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineOKHours", field: "Slow_TotalTurbineOKHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineOKHours StdDev", field: "Slow_TotalTurbineOKHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineOKHours Min", field: "Slow_TotalTurbineOKHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineOKHours Max", field: "Slow_TotalTurbineOKHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineOKHours Count", field: "Slow_TotalTurbineOKHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineTimeAllHours", field: "Slow_TotalTurbineTimeAllHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineTimeAllHours StdDev", field: "Slow_TotalTurbineTimeAllHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineTimeAllHours Min", field: "Slow_TotalTurbineTimeAllHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineTimeAllHours Max", field: "Slow_TotalTurbineTimeAllHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalTurbineTimeAllHours Count", field: "Slow_TotalTurbineTimeAllHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L1", field: "Slow_TempG1L1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L1 StdDev", field: "Slow_TempG1L1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L1 Min", field: "Slow_TempG1L1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L1 Max", field: "Slow_TempG1L1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempG1L1 Count", field: "Slow_TempG1L1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxOilSump", field: "Slow_TempGearBoxOilSump", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxOilSump StdDev", field: "Slow_TempGearBoxOilSump_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxOilSump Min", field: "Slow_TempGearBoxOilSump_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxOilSump Max", field: "Slow_TempGearBoxOilSump_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxOilSump Count", field: "Slow_TempGearBoxOilSump_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV2", field: "Fast_PitchAccuV2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV2 StdDev", field: "Fast_PitchAccuV2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV2 Min", field: "Fast_PitchAccuV2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV2 Max", field: "Fast_PitchAccuV2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchAccuV2 Count", field: "Fast_PitchAccuV2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalGridOkHours", field: "Slow_TotalGridOkHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalGridOkHours StdDev", field: "Slow_TotalGridOkHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalGridOkHours Min", field: "Slow_TotalGridOkHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalGridOkHours Max", field: "Slow_TotalGridOkHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalGridOkHours Count", field: "Slow_TotalGridOkHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOut kWh", field: "Slow_TotalActPowerOut_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOut kWh StdDev", field: "Slow_TotalActPowerOut_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOut kWh Min", field: "Slow_TotalActPowerOut_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOut kWh Max", field: "Slow_TotalActPowerOut_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOut kWh Count", field: "Slow_TotalActPowerOut_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawService", field: "Fast_YawService", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawService StdDev", field: "Fast_YawService_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawService Min", field: "Fast_YawService_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawService Max", field: "Fast_YawService_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawService Count", field: "Fast_YawService_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawAngle", field: "Fast_YawAngle", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawAngle StdDev", field: "Fast_YawAngle_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawAngle Min", field: "Fast_YawAngle_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawAngle Max", field: "Fast_YawAngle_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast YawAngle Count", field: "Fast_YawAngle_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductivePwrFactor", field: "Slow_CapableInductivePwrFactor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductivePwrFactor StdDev", field: "Slow_CapableInductivePwrFactor_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductivePwrFactor Min", field: "Slow_CapableInductivePwrFactor_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductivePwrFactor Max", field: "Slow_CapableInductivePwrFactor_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CapableInductivePwrFactor Count", field: "Slow_CapableInductivePwrFactor_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSNDE", field: "Slow_TempGearBoxHSSNDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSNDE StdDev", field: "Slow_TempGearBoxHSSNDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSNDE Min", field: "Slow_TempGearBoxHSSNDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSNDE Max", field: "Slow_TempGearBoxHSSNDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxHSSNDE Count", field: "Slow_TempGearBoxHSSNDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempHubBearing", field: "Slow_TempHubBearing", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempHubBearing StdDev", field: "Slow_TempHubBearing_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempHubBearing Min", field: "Slow_TempHubBearing_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempHubBearing Max", field: "Slow_TempHubBearing_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempHubBearing Count", field: "Slow_TempHubBearing_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG1ActiveHours", field: "Slow_TotalG1ActiveHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG1ActiveHours StdDev", field: "Slow_TotalG1ActiveHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG1ActiveHours Min", field: "Slow_TotalG1ActiveHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG1ActiveHours Max", field: "Slow_TotalG1ActiveHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG1ActiveHours Count", field: "Slow_TotalG1ActiveHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG1 kWh", field: "Slow_TotalActPowerOutG1_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG1 kWh StdDev", field: "Slow_TotalActPowerOutG1_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG1 kWh Min", field: "Slow_TotalActPowerOutG1_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG1 kWh Max", field: "Slow_TotalActPowerOutG1_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG1 kWh Count", field: "Slow_TotalActPowerOutG1_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG1 kVArh", field: "Slow_TotalReactPowerInG1_kVArh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG1 kVArh StdDev", field: "Slow_TotalReactPowerInG1_kVArh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG1 kVArh Min", field: "Slow_TotalReactPowerInG1_kVArh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG1 kVArh Max", field: "Slow_TotalReactPowerInG1_kVArh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG1 kVArh Count", field: "Slow_TotalReactPowerInG1_kVArh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacelleDrill", field: "Slow_NacelleDrill", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacelleDrill StdDev", field: "Slow_NacelleDrill_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacelleDrill Min", field: "Slow_NacelleDrill_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacelleDrill Max", field: "Slow_NacelleDrill_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow NacelleDrill Count", field: "Slow_NacelleDrill_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSDE", field: "Slow_TempGearBoxIMSDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSDE StdDev", field: "Slow_TempGearBoxIMSDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSDE Min", field: "Slow_TempGearBoxIMSDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSDE Max", field: "Slow_TempGearBoxIMSDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGearBoxIMSDE Count", field: "Slow_TempGearBoxIMSDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Operating hrs", field: "Fast_Total_Operating_hrs", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Operating hrs StdDev", field: "Fast_Total_Operating_hrs_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Operating hrs Min", field: "Fast_Total_Operating_hrs_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Operating hrs Max", field: "Fast_Total_Operating_hrs_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Operating hrs Count", field: "Fast_Total_Operating_hrs_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempNacelle", field: "Slow_TempNacelle", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempNacelle StdDev", field: "Slow_TempNacelle_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempNacelle Min", field: "Slow_TempNacelle_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempNacelle Max", field: "Slow_TempNacelle_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempNacelle Count", field: "Slow_TempNacelle_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Grid OK hrs", field: "Fast_Total_Grid_OK_hrs", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Grid OK hrs StdDev", field: "Fast_Total_Grid_OK_hrs_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Grid OK hrs Min", field: "Fast_Total_Grid_OK_hrs_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Grid OK hrs Max", field: "Fast_Total_Grid_OK_hrs_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Grid OK hrs Count", field: "Fast_Total_Grid_OK_hrs_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total WTG OK hrs", field: "Fast_Total_WTG_OK_hrs", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total WTG OK hrs StdDev", field: "Fast_Total_WTG_OK_hrs_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total WTG OK hrs Min", field: "Fast_Total_WTG_OK_hrs_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total WTG OK hrs Max", field: "Fast_Total_WTG_OK_hrs_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total WTG OK hrs Count", field: "Fast_Total_WTG_OK_hrs_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempCabinetTopBox", field: "Slow_TempCabinetTopBox", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempCabinetTopBox StdDev", field: "Slow_TempCabinetTopBox_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempCabinetTopBox Min", field: "Slow_TempCabinetTopBox_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempCabinetTopBox Max", field: "Slow_TempCabinetTopBox_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempCabinetTopBox Count", field: "Slow_TempCabinetTopBox_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingNDE", field: "Slow_TempGeneratorBearingNDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingNDE StdDev", field: "Slow_TempGeneratorBearingNDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingNDE Min", field: "Slow_TempGeneratorBearingNDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingNDE Max", field: "Slow_TempGeneratorBearingNDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingNDE Count", field: "Slow_TempGeneratorBearingNDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Access hrs", field: "Fast_Total_Access_hrs", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Access hrs StdDev", field: "Fast_Total_Access_hrs_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Access hrs Min", field: "Fast_Total_Access_hrs_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Access hrs Max", field: "Fast_Total_Access_hrs_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast Total Access hrs Count", field: "Fast_Total_Access_hrs_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomPowerSection", field: "Slow_TempBottomPowerSection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomPowerSection StdDev", field: "Slow_TempBottomPowerSection_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomPowerSection Min", field: "Slow_TempBottomPowerSection_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomPowerSection Max", field: "Slow_TempBottomPowerSection_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomPowerSection Count", field: "Slow_TempBottomPowerSection_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingDE", field: "Slow_TempGeneratorBearingDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingDE StdDev", field: "Slow_TempGeneratorBearingDE_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingDE Min", field: "Slow_TempGeneratorBearingDE_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingDE Max", field: "Slow_TempGeneratorBearingDE_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempGeneratorBearingDE Count", field: "Slow_TempGeneratorBearingDE_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerIn kVArh", field: "Slow_TotalReactPowerIn_kVArh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerIn kVArh StdDev", field: "Slow_TotalReactPowerIn_kVArh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerIn kVArh Min", field: "Slow_TotalReactPowerIn_kVArh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerIn kVArh Max", field: "Slow_TotalReactPowerIn_kVArh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerIn kVArh Count", field: "Slow_TotalReactPowerIn_kVArh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomControlSection", field: "Slow_TempBottomControlSection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomControlSection StdDev", field: "Slow_TempBottomControlSection_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomControlSection Min", field: "Slow_TempBottomControlSection_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomControlSection Max", field: "Slow_TempBottomControlSection_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomControlSection Count", field: "Slow_TempBottomControlSection_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv1", field: "Slow_TempConv1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv1 StdDev", field: "Slow_TempConv1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv1 Min", field: "Slow_TempConv1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv1 Max", field: "Slow_TempConv1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv1 Count", field: "Slow_TempConv1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerRated kW", field: "Fast_ActivePowerRated_kW", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerRated kW StdDev", field: "Fast_ActivePowerRated_kW_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerRated kW Min", field: "Fast_ActivePowerRated_kW_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerRated kW Max", field: "Fast_ActivePowerRated_kW_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast ActivePowerRated kW Count", field: "Fast_ActivePowerRated_kW_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast NodeIP", field: "Fast_NodeIP", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast NodeIP StdDev", field: "Fast_NodeIP_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast NodeIP Min", field: "Fast_NodeIP_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast NodeIP Max", field: "Fast_NodeIP_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast NodeIP Count", field: "Fast_NodeIP_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchSpeed1", field: "Fast_PitchSpeed1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchSpeed1 StdDev", field: "Fast_PitchSpeed1_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchSpeed1 Min", field: "Fast_PitchSpeed1_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchSpeed1 Max", field: "Fast_PitchSpeed1_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Fast PitchSpeed1 Count", field: "Fast_PitchSpeed1_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSize", field: "Slow_CFCardSize", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSize StdDev", field: "Slow_CFCardSize_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSize Min", field: "Slow_CFCardSize_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSize Max", field: "Slow_CFCardSize_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSize Count", field: "Slow_CFCardSize_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CPU Number", field: "Slow_CPU_Number", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CPU Number StdDev", field: "Slow_CPU_Number_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CPU Number Min", field: "Slow_CPU_Number_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CPU Number Max", field: "Slow_CPU_Number_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CPU Number Count", field: "Slow_CPU_Number_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSpaceLeft", field: "Slow_CFCardSpaceLeft", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSpaceLeft StdDev", field: "Slow_CFCardSpaceLeft_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSpaceLeft Min", field: "Slow_CFCardSpaceLeft_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSpaceLeft Max", field: "Slow_CFCardSpaceLeft_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow CFCardSpaceLeft Count", field: "Slow_CFCardSpaceLeft_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomCapSection", field: "Slow_TempBottomCapSection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomCapSection StdDev", field: "Slow_TempBottomCapSection_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomCapSection Min", field: "Slow_TempBottomCapSection_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomCapSection Max", field: "Slow_TempBottomCapSection_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempBottomCapSection Count", field: "Slow_TempBottomCapSection_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow RatedPower", field: "Slow_RatedPower", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow RatedPower StdDev", field: "Slow_RatedPower_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow RatedPower Min", field: "Slow_RatedPower_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow RatedPower Max", field: "Slow_RatedPower_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow RatedPower Count", field: "Slow_RatedPower_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv3", field: "Slow_TempConv3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv3 StdDev", field: "Slow_TempConv3_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv3 Min", field: "Slow_TempConv3_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv3 Max", field: "Slow_TempConv3_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv3 Count", field: "Slow_TempConv3_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv2", field: "Slow_TempConv2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv2 StdDev", field: "Slow_TempConv2_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv2 Min", field: "Slow_TempConv2_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv2 Max", field: "Slow_TempConv2_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TempConv2 Count", field: "Slow_TempConv2_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerIn kWh", field: "Slow_TotalActPowerIn_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerIn kWh StdDev", field: "Slow_TotalActPowerIn_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerIn kWh Min", field: "Slow_TotalActPowerIn_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerIn kWh Max", field: "Slow_TotalActPowerIn_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerIn kWh Count", field: "Slow_TotalActPowerIn_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG1 kWh", field: "Slow_TotalActPowerInG1_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG1 kWh StdDev", field: "Slow_TotalActPowerInG1_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG1 kWh Min", field: "Slow_TotalActPowerInG1_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG1 kWh Max", field: "Slow_TotalActPowerInG1_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG1 kWh Count", field: "Slow_TotalActPowerInG1_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG2 kWh", field: "Slow_TotalActPowerInG2_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG2 kWh StdDev", field: "Slow_TotalActPowerInG2_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG2 kWh Min", field: "Slow_TotalActPowerInG2_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG2 kWh Max", field: "Slow_TotalActPowerInG2_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerInG2 kWh Count", field: "Slow_TotalActPowerInG2_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG2 kWh", field: "Slow_TotalActPowerOutG2_kWh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG2 kWh StdDev", field: "Slow_TotalActPowerOutG2_kWh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG2 kWh Min", field: "Slow_TotalActPowerOutG2_kWh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG2 kWh Max", field: "Slow_TotalActPowerOutG2_kWh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalActPowerOutG2 kWh Count", field: "Slow_TotalActPowerOutG2_kWh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG2ActiveHours", field: "Slow_TotalG2ActiveHours", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG2ActiveHours StdDev", field: "Slow_TotalG2ActiveHours_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG2ActiveHours Min", field: "Slow_TotalG2ActiveHours_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG2ActiveHours Max", field: "Slow_TotalG2ActiveHours_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalG2ActiveHours Count", field: "Slow_TotalG2ActiveHours_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG2 kVArh", field: "Slow_TotalReactPowerInG2_kVArh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG2 kVArh StdDev", field: "Slow_TotalReactPowerInG2_kVArh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG2 kVArh Min", field: "Slow_TotalReactPowerInG2_kVArh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG2 kVArh Max", field: "Slow_TotalReactPowerInG2_kVArh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerInG2 kVArh Count", field: "Slow_TotalReactPowerInG2_kVArh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerOut kVArh", field: "Slow_TotalReactPowerOut_kVArh", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerOut kVArh StdDev", field: "Slow_TotalReactPowerOut_kVArh_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerOut kVArh Min", field: "Slow_TotalReactPowerOut_kVArh_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerOut kVArh Max", field: "Slow_TotalReactPowerOut_kVArh_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow TotalReactPowerOut kVArh Count", field: "Slow_TotalReactPowerOut_kVArh_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow UTCoffset int", field: "Slow_UTCoffset_int", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow UTCoffset int StdDev", field: "Slow_UTCoffset_int_StdDev", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow UTCoffset int Min", field: "Slow_UTCoffset_int_Min", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow UTCoffset int Max", field: "Slow_UTCoffset_int_Max", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Slow UTCoffset int Count", field: "Slow_UTCoffset_int_Count", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
            ]
        });

        var grid = $('#scadahfdGrid').data('kendoGrid');
        var columns = grid.columns;
        dbr.gridColumnsScada([]);
        $.each(columns, function (i, val) {
            $('#scadahfdGrid').data("kendoGrid").showColumn(val.field);
            var result = {
                field: val.field,
                title: val.title,
                value: true
            }

            dbr.gridColumnsScada.push(result);

                
        });
    },
    InitScadaGrid: function() {
        dbr.oemvis(true);
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

        $('#scadaGrid').html("");
        $('#scadaGrid').kendoGrid({
            dataSource: {
                filter: [{
                    field: "timestamp",
                    operator: "gte",
                    value: dateStart
                }, {
                    field: "timestamp",
                    operator: "lte",
                    value: dateEnd
                }, {
                    field: "turbine",
                    operator: "in",
                    value: turbine
                }],
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getscadaoemlist",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function(options) {
                        dbr.oemvis(true);
                        return JSON.stringify(options);
                    }
                },
                pageSize: 10,
                schema: {
                    data: function(res) {
                        app.loading(false);
                        dbr.oemvis(false);
                        if (!app.isFine(res)) {   
                            return;
                        }
                        return res.data.Data
                    },
                    total: function(res) {
                        console.log(res);
                        if (!app.isFine(res)) {

                        console.log(res);
                            return;
                        }

                        $('#totalturbine').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldata').html(kendo.toString(res.data.Total, 'n0'));
                        $('#totalactivepower').html(kendo.toString(res.data.TotalActivePower / 1000, 'n0') + ' MWh');
                        $('#totalprodoem').html(kendo.toString(res.data.TotalEnergy / 1000, 'n0') + ' MWh');
                        $('#avgwindspeedoem').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
                        return res.data.Total;
                    }
                },
                sort: [{
                    field: 'TimeStamp',
                    dir: 'asc'
                }, {
                    field: 'Turbine',
                    dir: 'asc'
                }],
            },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            pageable: true,
            columns: [{
                title: "Time Stamp",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 130,
                locked: true,
                filterable: false
            }, {
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                locked: true,
                filterable: false
            }, {
                title: "Ai Intern R Pid Angle Out",
                field: "AI_intern_R_PidAngleOut",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Active Power",
                field: "AI_intern_ActivPower",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern I1",
                field: "AI_intern_I1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern I2",
                field: "AI_intern_I2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern I3",
                field: "AI_intern_I3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Nacelle Drill",
                field: "AI_intern_NacelleDrill",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Nacelle Pos",
                field: "AI_intern_NacellePos",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Akku V1",
                field: "AI_intern_PitchAkku_V1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Akku V2",
                field: "AI_intern_PitchAkku_V2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Akku V3",
                field: "AI_intern_PitchAkku_V3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Angle1",
                field: "AI_intern_PitchAngle1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Angle2",
                field: "AI_intern_PitchAngle2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Angle3",
                field: "AI_intern_PitchAngle3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Conv Current1",
                field: "AI_intern_PitchConv_Current1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Conv Current2",
                field: "AI_intern_PitchConv_Current2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Conv Current3",
                field: "AI_intern_PitchConv_Current3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Angle Sp Diff1",
                field: "AI_intern_PitchAngleSP_Diff1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Angle Sp Diff2",
                field: "AI_intern_PitchAngleSP_Diff2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchangle Sp Diff3",
                field: "AI_intern_PitchAngleSP_Diff3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Reactive Power",
                field: "AI_intern_ReactivPower",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Rpm Diff",
                field: "AI_intern_RpmDiff",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern U1",
                field: "AI_intern_U1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern U2",
                field: "AI_intern_U2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern U3",
                field: "AI_intern_U3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Wind Direction",
                field: "AI_intern_WindDirection",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Wind Speed",
                field: "AI_intern_WindSpeed",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Wind Speed Dif",
                field: "AI_Intern_WindSpeedDif",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Speed Rot Fr",
                field: "AI_speed_RotFR",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Wind Speed1",
                field: "AI_WindSpeed1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Wind Speed2",
                field: "AI_WindSpeed2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Wind Vane1",
                field: "AI_WindVane1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Wind Vane2",
                field: "AI_WindVane2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Current Asym",
                field: "AI_internCurrentAsym",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Gear Box Ims Nde",
                field: "Temp_GearBox_IMS_NDE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Wind Vane Diff",
                field: "AI_intern_WindVaneDiff",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "C Intern Speed Generator",
                field: "C_intern_SpeedGenerator",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "C Intern Speed Rotor",
                field: "C_intern_SpeedRotor",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Speed Rpm Diff Fr1 Rot Cnt",
                field: "AI_intern_Speed_RPMDiff_FR1_RotCNT",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Frequency Grid",
                field: "AI_intern_Frequency_Grid",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Gear Box Hss Nde",
                field: "Temp_GearBox_HSS_NDE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Dr Tr Vib Value",
                field: "AI_DrTrVibValue",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Last Error Conv1",
                field: "AI_intern_InLastErrorConv1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Last Error Conv2",
                field: "AI_intern_InLastErrorConv2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Last Error Conv3",
                field: "AI_intern_InLastErrorConv3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Conv1",
                field: "AI_intern_TempConv1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Conv2",
                field: "AI_intern_TempConv2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Conv3",
                field: "AI_intern_TempConv3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch Speed1",
                field: "AI_intern_PitchSpeed1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Brake 1",
                field: "Temp_YawBrake_1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Brake 2",
                field: "Temp_YawBrake_2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp G1L1",
                field: "Temp_G1L1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp G1L2",
                field: "Temp_G1L2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp G1L3",
                field: "Temp_G1L3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Brake 3",
                field: "Temp_YawBrake_3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Hydr System Pressure",
                field: "AI_HydrSystemPressure",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Bottom Control Section Low",
                field: "Temp_BottomControlSection_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Gearbox Hss De",
                field: "Temp_GearBox_HSS_DE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Gear Oil Sump",
                field: "Temp_GearOilSump",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Generator Bearing De",
                field: "Temp_GeneratorBearing_DE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Generator Bearing Nde",
                field: "Temp_GeneratorBearing_NDE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Main Bearing",
                field: "Temp_MainBearing",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Gear Box Ims De",
                field: "Temp_GearBox_IMS_DE",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Nacelle",
                field: "Temp_Nacelle",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Outdoor",
                field: "Temp_Outdoor",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Tower Vib Value Axial",
                field: "AI_TowerVibValueAxial",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Diff Gen Speed Sp To Act",
                field: "AI_intern_DiffGenSpeedSPToAct",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Brake 4",
                field: "Temp_YawBrake_4",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Speed Generator Proximity",
                field: "AI_intern_SpeedGenerator_Proximity",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Speed Diff Encoder Proximity",
                field: "AI_intern_SpeedDiff_Encoder_Proximity",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Gear Oil Pressure",
                field: "AI_GearOilPressure",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Cabinet Top Box Low",
                field: "Temp_CabinetTopBox_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Cabinet Top Box",
                field: "Temp_CabinetTopBox",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Bottom Control Section",
                field: "Temp_BottomControlSection",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Bottom Power Section",
                field: "Temp_BottomPowerSection",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Bottom Power Section Low",
                field: "Temp_BottomPowerSection_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch1 Status High",
                field: "AI_intern_Pitch1_Status_High",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch2 Status High",
                field: "AI_intern_Pitch2_Status_High",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch3 Status High",
                field: "AI_intern_Pitch3_Status_High",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Position1 Ch2",
                field: "AI_intern_InPosition1_ch2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Position2 Ch2",
                field: "AI_intern_InPosition2_ch2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern In Position3 Ch2",
                field: "AI_intern_InPosition3_ch2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Brake Blade1",
                field: "AI_intern_Temp_Brake_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Brake Blade2",
                field: "AI_intern_Temp_Brake_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Brake Blade3",
                field: "AI_intern_Temp_Brake_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch Motor Blade1",
                field: "AI_intern_Temp_PitchMotor_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch Motor Blade2",
                field: "AI_intern_Temp_PitchMotor_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch Motor Blade3",
                field: "AI_intern_Temp_PitchMotor_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Additional1",
                field: "AI_intern_Temp_Hub_Additional1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Additional2",
                field: "AI_intern_Temp_Hub_Additional2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Additional3",
                field: "AI_intern_Temp_Hub_Additional3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch1 Status Low",
                field: "AI_intern_Pitch1_Status_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch2 Status Low",
                field: "AI_intern_Pitch2_Status_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitch3 Status Low",
                field: "AI_intern_Pitch3_Status_Low",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Voltage Blade1 Center",
                field: "AI_intern_Battery_VoltageBlade1_center",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Voltage Blade2 Center",
                field: "AI_intern_Battery_VoltageBlade2_center",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Voltage Blade3 Center",
                field: "AI_intern_Battery_VoltageBlade3_center",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Charging Cur Blade1",
                field: "AI_intern_Battery_ChargingCur_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Charging Cur Blade2",
                field: "AI_intern_Battery_ChargingCur_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Charging Cur Blade3",
                field: "AI_intern_Battery_ChargingCur_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Discharging Cur Blade1",
                field: "AI_intern_Battery_DischargingCur_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Discharging Cur Blade2",
                field: "AI_intern_Battery_DischargingCur_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Battery Discharging Cur Blade3",
                field: "AI_intern_Battery_DischargingCur_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Voltage Blade1",
                field: "AI_intern_PitchMotor_BrakeVoltage_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Voltage Blade2",
                field: "AI_intern_PitchMotor_BrakeVoltage_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Voltage Blade3",
                field: "AI_intern_PitchMotor_BrakeVoltage_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Current Blade1",
                field: "AI_intern_PitchMotor_BrakeCurrent_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Current Blade2",
                field: "AI_intern_PitchMotor_BrakeCurrent_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Pitchmotor Brake Current Blade3",
                field: "AI_intern_PitchMotor_BrakeCurrent_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Box Blade1",
                field: "AI_intern_Temp_HubBox_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Box Blade2",
                field: "AI_intern_Temp_HubBox_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Hub Box Blade3",
                field: "AI_intern_Temp_HubBox_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch1 Heat Sink",
                field: "AI_intern_Temp_Pitch1_HeatSink",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch2 Heat Sink",
                field: "AI_intern_Temp_Pitch2_HeatSink",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Pitch3 Heat Sink",
                field: "AI_intern_Temp_Pitch3_HeatSink",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Error Stack Blade1",
                field: "AI_intern_ErrorStackBlade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Error Stack Blade2",
                field: "AI_intern_ErrorStackBlade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Error Stack Blade3",
                field: "AI_intern_ErrorStackBlade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Battery Box Blade1",
                field: "AI_intern_Temp_BatteryBox_Blade1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Battery Box Blade2",
                field: "AI_intern_Temp_BatteryBox_Blade2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Temp Battery Box Blade3",
                field: "AI_intern_Temp_BatteryBox_Blade3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dc Linkvoltage1",
                field: "AI_intern_DC_LinkVoltage1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dc Linkvoltage2",
                field: "AI_intern_DC_LinkVoltage2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dc Linkvoltage3",
                field: "AI_intern_DC_LinkVoltage3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Motor1",
                field: "Temp_Yaw_Motor1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Motor2",
                field: "Temp_Yaw_Motor2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Motor3",
                field: "Temp_Yaw_Motor3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Temp Yaw Motor4",
                field: "Temp_Yaw_Motor4",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ao Dfig Power Setpiont",
                field: "AO_DFIG_Power_Setpiont",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ao Dfig Q Setpoint",
                field: "AO_DFIG_Q_Setpoint",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Dfig Torque Actual",
                field: "AI_DFIG_Torque_actual",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Dfig Speed Generator Encoder",
                field: "AI_DFIG_SpeedGenerator_Encoder",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Dc Link Voltage Actual",
                field: "AI_intern_DFIG_DC_Link_Voltage_actual",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Msc Current",
                field: "AI_intern_DFIG_MSC_current",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Main Voltage",
                field: "AI_intern_DFIG_Main_voltage",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Main Current",
                field: "AI_intern_DFIG_Main_current",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Active Power Actual",
                field: "AI_intern_DFIG_active_power_actual",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Reactive Power Actual",
                field: "AI_intern_DFIG_reactive_power_actual",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Active Power Actual Lsc",
                field: "AI_intern_DFIG_active_power_actual_LSC",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Lsc Current",
                field: "AI_intern_DFIG_LSC_current",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Dfig Data Log Number",
                field: "AI_intern_DFIG_Data_log_number",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Damper Osc Magnitude",
                field: "AI_intern_Damper_OscMagnitude",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Damper Passband Full Load",
                field: "AI_intern_Damper_PassbandFullLoad",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Yaw Brake Temp Rise1",
                field: "AI_YawBrake_TempRise1",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Yaw Brake Temp Rise2",
                field: "AI_YawBrake_TempRise2",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Yaw Brake Temp Rise3",
                field: "AI_YawBrake_TempRise3",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Yaw Brake Temp Rise4",
                field: "AI_YawBrake_TempRise4",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ai Intern Nacelle Drill At North Pos Sensor",
                field: "AI_intern_NacelleDrill_at_NorthPosSensor",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, ]
        });

        var grid = $('#scadaGrid').data('kendoGrid');
        var columns = grid.columns;
        dbr.gridColumnsScada([]);
        $.each(columns, function(i, val) {
            $('#scadaGrid').data("kendoGrid").showColumn(val.field);
            var result = {
                field: val.field,
                title: val.title,
                value: true
            }

            dbr.gridColumnsScada.push(result);


        });
    },
    InitDEgrid: function() {
        dbr.downeventvis(true);
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
        var param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Turbine: turbine,
        };

        $('#DEgrid').html("");
        $('#DEgrid').kendoGrid({
            dataSource: {
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getdowntimeeventlist",
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
                        // app.loading(false);
                        dbr.downeventvis(false);
                        if (!app.isFine(ress)) {
                            return;
                        }
                        return ress.data.Data
                    },
                    total: function(res) {

                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totalturbineDE').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataDE').html(kendo.toString(res.data.Total, 'n0'));

                        return res.data.Total;
                    }
                },
                sort: [{
                    field: 'timestart',
                    dir: 'asc'
                }, {
                    field: 'turbine',
                    dir: 'asc'
                }],
            },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            pageable: true,
            columns: [{
                    title: "Time Start",
                    field: "TimeStart",
                    template: "#= kendo.toString(moment.utc(TimeStart).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                    width: 100,
                    filterable: false
                },

                {
                    title: "Time End",
                    field: "TimeEnd",
                    template: "#= kendo.toString(moment.utc(TimeEnd).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                    width: 100
                }, {
                    title: "Grid Down",
                    field: "DownGrid",
                    width: 80,
                    sortable: false,
                    template: '# if (DownGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Machine Down",
                    field: "DownMachine",
                    width: 80,
                    sortable: false,
                    template: '# if (DownMachine == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Environment Down",
                    field: "DownEnvironment",
                    width: 80,
                    sortable: false,
                    template: '# if (DownEnvironment == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Turbine",
                    field: "Turbine",
                    attributes: {
                        class: "align-center"
                    },
                    width: 90,
                    filterable: false
                }, {
                    title: "Alarm Description",
                    field: "AlarmDescription",
                    width: 100,
                    filterable: false
                }, {
                    title: "Duration (hh:mm:ss)",
                    field: "Duration",
                    template: '#= kendo.toString(secondsToHms(Duration)) #',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },


            ]
        });
    },
    InitCustomGrid: function() {
        dbr.customvis(true);
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

        var param = {
            "Custom": {
                "ColumnList": (dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn())
            }
        };

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
                filter: [{
                    field: "timestamp",
                    operator: "gte",
                    value: dateStart
                }, {
                    field: "timestamp",
                    operator: "lte",
                    value: dateEnd
                }, {
                    field: "turbine",
                    operator: "in",
                    value: turbine
                }],
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
                        if (!app.isFine(res)) {
                            return;
                        }

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
            pageable: true,
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
    },
    InitEventGrid: function() {
        dbr.eventrawvis(true);
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

        var param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Turbine: turbine,
        };


        $('#EventGrid').html("");
        $('#EventGrid').kendoGrid({
            dataSource: {
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/geteventlist",
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
                        // app.loading(false);
                        dbr.eventrawvis(false);
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }

                        $('#totalturbineEvent').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataEvent').html(kendo.toString(res.data.Total, 'n0'));

                        return res.data.Total;
                    }
                },
                sort: [{
                    field: 'TimeStamp',
                    dir: 'asc'
                }, {
                    field: 'Turbine',
                    dir: 'asc'
                }],
            },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            pageable: true,
            columns: [{
                title: "Time Stamp",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                width: 130,
                filterable: false
            }, {
                title: "Project Name",
                field: "ProjectName",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                filterable: false
            }, {
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                filterable: false
            }, {
                title: "Event Type",
                field: "EventType",
                attributes: {
                    class: "align-center"
                },
                width: 100,
                filterable: false
            }, {
                title: "Alarm Description",
                field: "AlarmDescription",
                attributes: {
                    class: "align-center"
                },
                width: 150,
                filterable: false
            }, {
                title: "Turbine Status",
                field: "TurbineStatus",
                attributes: {
                    class: "align-center"
                },
                width: 120,
                filterable: false
            }, {
                title: "Brake Type",
                field: "BrakeType",
                attributes: {
                    class: "align-center"
                },
                width: 150,
                filterable: false
            }, {
                title: "Brake Program",
                field: "BrakeProgram",
                width: 120,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Alarm Id",
                field: "AlarmId",
                width: 120,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Alarm Toggle",
                field: "AlarmToggle",
                width: 120,
                sortable: false,
                template: '# if (AlarmToggle == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#',
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align:center;"
                }
            }, ]
        });
    },
    InitMet: function() {
        dbr.mettowervis(true);
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();

        dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));
        dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));

        var param = {};

        $('#dataGridMet').html("");
        $('#dataGridMet').kendoGrid({
            dataSource: {
                filter: [{
                    field: "timestamp",
                    operator: "gte",
                    value: dateStart
                }, {
                    field: "timestamp",
                    operator: "lte",
                    value: dateEnd
                }],
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getmetlist",
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
                    model: {
                        fields: {
                            VHubWS90mAvg: {
                                type: "number"
                            },
                            VHubWS90mMax: {
                                type: "number"
                            },
                            VHubWS90mMin: {
                                type: "number"
                            },
                            VHubWS90mStdDev: {
                                type: "number"
                            },
                            VHubWS90mCount: {
                                type: "number"
                            },

                            VRefWS88mAvg: {
                                type: "number"
                            },
                            VRefWS88mMax: {
                                type: "number"
                            },
                            VRefWS88mMin: {
                                type: "number"
                            },
                            VRefWS88mStdDev: {
                                type: "number"
                            },
                            VRefWS88mCount: {
                                type: "number"
                            },

                            VTipWS42mAvg: {
                                type: "number"
                            },
                            VTipWS42mMax: {
                                type: "number"
                            },
                            VTipWS42mMin: {
                                type: "number"
                            },
                            VTipWS42mStdDev: {
                                type: "number"
                            },
                            VTipWS42mCount: {
                                type: "number"
                            },

                            DHubWD88mAvg: {
                                type: "number"
                            },
                            DHubWD88mMax: {
                                type: "number"
                            },
                            DHubWD88mMin: {
                                type: "number"
                            },
                            DHubWD88mStdDev: {
                                type: "number"
                            },
                            DHubWD88mCount: {
                                type: "number"
                            },

                            DRefWD86mAvg: {
                                type: "number"
                            },
                            DRefWD86mMax: {
                                type: "number"
                            },
                            DRefWD86mMin: {
                                type: "number"
                            },
                            DRefWD86mStdDev: {
                                type: "number"
                            },
                            DRefWD86mCount: {
                                type: "number"
                            },

                            THubHHubHumid855mAvg: {
                                type: "number"
                            },
                            THubHHubHumid855mMax: {
                                type: "number"
                            },
                            THubHHubHumid855mMin: {
                                type: "number"
                            },
                            THubHHubHumid855mStdDev: {
                                type: "number"
                            },
                            THubHHubHumid855mCount: {
                                type: "number"
                            },

                            TRefHRefHumid855mAvg: {
                                type: "number"
                            },
                            TRefHRefHumid855mMax: {
                                type: "number"
                            },
                            TRefHRefHumid855mMin: {
                                type: "number"
                            },
                            TRefHRefHumid855mStdDev: {
                                type: "number"
                            },
                            TRefHRefHumid855mCount: {
                                type: "number"
                            },

                            THubHHubTemp855mAvg: {
                                type: "number"
                            },
                            THubHHubTemp855mMax: {
                                type: "number"
                            },
                            THubHHubTemp855mMin: {
                                type: "number"
                            },
                            THubHHubTemp855mStdDev: {
                                type: "number"
                            },
                            THubHHubTemp855mCount: {
                                type: "number"
                            },

                            TRefHRefTemp855mAvg: {
                                type: "number"
                            },
                            TRefHRefTemp855mMax: {
                                type: "number"
                            },
                            TRefHRefTemp855mMin: {
                                type: "number"
                            },
                            TRefHRefTemp855mStdDev: {
                                type: "number"
                            },
                            TRefHRefTemp855mCount: {
                                type: "number"
                            },

                            BaroAirPress855mAvg: {
                                type: "number"
                            },
                            BaroAirPress855mMax: {
                                type: "number"
                            },
                            BaroAirPress855mMin: {
                                type: "number"
                            },
                            BaroAirPress855mStdDev: {
                                type: "number"
                            },
                            BaroAirPress855mCount: {
                                type: "number"
                            },

                            YawAngleVoltageAvg: {
                                type: "number"
                            },
                            YawAngleVoltageMax: {
                                type: "number"
                            },
                            YawAngleVoltageMin: {
                                type: "number"
                            },
                            YawAngleVoltageStdDev: {
                                type: "number"
                            },
                            YawAngleVoltageCount: {
                                type: "number"
                            },
                            OtherSensorVoltageAI1Avg: {
                                type: "number"
                            },
                            OtherSensorVoltageAI1Max: {
                                type: "number"
                            },
                            OtherSensorVoltageAI1Min: {
                                type: "number"
                            },
                            OtherSensorVoltageAI1StdDev: {
                                type: "number"
                            },
                            OtherSensorVoltageAI1Count: {
                                type: "number"
                            },
                            OtherSensorVoltageAI2Avg: {
                                type: "number"
                            },
                            OtherSensorVoltageAI2Max: {
                                type: "number"
                            },
                            OtherSensorVoltageAI2Min: {
                                type: "number"
                            },
                            OtherSensorVoltageAI2StdDev: {
                                type: "number"
                            },
                            OtherSensorVoltageAI2Count: {
                                type: "number"
                            },
                            OtherSensorVoltageAI3Avg: {
                                type: "number"
                            },
                            OtherSensorVoltageAI3Max: {
                                type: "number"
                            },
                            OtherSensorVoltageAI3Min: {
                                type: "number"
                            },
                            OtherSensorVoltageAI3StdDev: {
                                type: "number"
                            },
                            OtherSensorVoltageAI3Count: {
                                type: "number"
                            },
                            OtherSensorVoltageAI4Avg: {
                                type: "number"
                            },
                            OtherSensorVoltageAI4Max: {
                                type: "number"
                            },
                            OtherSensorVoltageAI4Min: {
                                type: "number"
                            },
                            OtherSensorVoltageAI4StdDev: {
                                type: "number"
                            },
                            OtherSensorVoltageAI4Count: {
                                type: "number"
                            },
                            GenRPMCurrentAvg: {
                                type: "number"
                            },
                            GenRPMCurrentMax: {
                                type: "number"
                            },
                            GenRPMCurrentMin: {
                                type: "number"
                            },
                            GenRPMCurrentStdDev: {
                                type: "number"
                            },
                            GenRPMCurrentCount: {
                                type: "number"
                            },
                            WS_SCSCurrentAvg: {
                                type: "number"
                            },
                            WS_SCSCurrentMax: {
                                type: "number"
                            },
                            WS_SCSCurrentMin: {
                                type: "number"
                            },
                            WS_SCSCurrentStdDev: {
                                type: "number"
                            },
                            WS_SCSCurrentCount: {
                                type: "number"
                            },
                            RainStatusCount: {
                                type: "number"
                            },
                            RainStatusSum: {
                                type: "number"
                            },
                            OtherSensor2StatusIO1Avg: {
                                type: "number"
                            },
                            OtherSensor2StatusIO1Max: {
                                type: "number"
                            },
                            OtherSensor2StatusIO1Min: {
                                type: "number"
                            },
                            OtherSensor2StatusIO1StdDev: {
                                type: "number"
                            },
                            OtherSensor2StatusIO1Count: {
                                type: "number"
                            },
                            OtherSensor2StatusIO2Avg: {
                                type: "number"
                            },
                            OtherSensor2StatusIO2Max: {
                                type: "number"
                            },
                            OtherSensor2StatusIO2Min: {
                                type: "number"
                            },
                            OtherSensor2StatusIO2StdDev: {
                                type: "number"
                            },
                            OtherSensor2StatusIO2Count: {
                                type: "number"
                            },
                            OtherSensor2StatusIO3Avg: {
                                type: "number"
                            },
                            OtherSensor2StatusIO3Max: {
                                type: "number"
                            },
                            OtherSensor2StatusIO3Min: {
                                type: "number"
                            },
                            OtherSensor2StatusIO3StdDev: {
                                type: "number"
                            },
                            OtherSensor2StatusIO3Count: {
                                type: "number"
                            },
                            OtherSensor2StatusIO4Avg: {
                                type: "number"
                            },
                            OtherSensor2StatusIO4Max: {
                                type: "number"
                            },
                            OtherSensor2StatusIO4Min: {
                                type: "number"
                            },
                            OtherSensor2StatusIO4StdDev: {
                                type: "number"
                            },
                            OtherSensor2StatusIO4Count: {
                                type: "number"
                            },
                            OtherSensor2StatusIO5Avg: {
                                type: "number"
                            },
                            OtherSensor2StatusIO5Max: {
                                type: "number"
                            },
                            OtherSensor2StatusIO5Min: {
                                type: "number"
                            },
                            OtherSensor2StatusIO5StdDev: {
                                type: "number"
                            },
                            OtherSensor2StatusIO5Count: {
                                type: "number"
                            },
                            A1Avg: {
                                type: "number"
                            },
                            A1Max: {
                                type: "number"
                            },
                            A1Min: {
                                type: "number"
                            },
                            A1StdDev: {
                                type: "number"
                            },
                            A1Count: {
                                type: "number"
                            },
                            A2Avg: {
                                type: "number"
                            },
                            A2Max: {
                                type: "number"
                            },
                            A2Min: {
                                type: "number"
                            },
                            A2StdDev: {
                                type: "number"
                            },
                            A2Count: {
                                type: "number"
                            },
                            A3Avg: {
                                type: "number"
                            },
                            A3Max: {
                                type: "number"
                            },
                            A3Min: {
                                type: "number"
                            },
                            A3StdDev: {
                                type: "number"
                            },
                            A3Count: {
                                type: "number"
                            },
                            A4Avg: {
                                type: "number"
                            },
                            A4Max: {
                                type: "number"
                            },
                            A4Min: {
                                type: "number"
                            },
                            A4StdDev: {
                                type: "number"
                            },
                            A4Count: {
                                type: "number"
                            },
                            A5Avg: {
                                type: "number"
                            },
                            A5Max: {
                                type: "number"
                            },
                            A5Min: {
                                type: "number"
                            },
                            A5StdDev: {
                                type: "number"
                            },
                            A5Count: {
                                type: "number"
                            },
                            A6Avg: {
                                type: "number"
                            },
                            A6Max: {
                                type: "number"
                            },
                            A6Min: {
                                type: "number"
                            },
                            A6StdDev: {
                                type: "number"
                            },
                            A6Count: {
                                type: "number"
                            },
                            A7Avg: {
                                type: "number"
                            },
                            A7Max: {
                                type: "number"
                            },
                            A7Min: {
                                type: "number"
                            },
                            A7StdDev: {
                                type: "number"
                            },
                            A7Count: {
                                type: "number"
                            },
                            A8Avg: {
                                type: "number"
                            },
                            A8Max: {
                                type: "number"
                            },
                            A8Min: {
                                type: "number"
                            },
                            A8StdDev: {
                                type: "number"
                            },
                            A8Count: {
                                type: "number"
                            },
                            A9Avg: {
                                type: "number"
                            },
                            A9Max: {
                                type: "number"
                            },
                            A9Min: {
                                type: "number"
                            },
                            A9StdDev: {
                                type: "number"
                            },
                            A9Count: {
                                type: "number"
                            },
                            A10Avg: {
                                type: "number"
                            },
                            A10Max: {
                                type: "number"
                            },
                            A10Min: {
                                type: "number"
                            },
                            A10StdDev: {
                                type: "number"
                            },
                            A10Count: {
                                type: "number"
                            },
                            AC1Avg: {
                                type: "number"
                            },
                            AC1Max: {
                                type: "number"
                            },
                            AC1Min: {
                                type: "number"
                            },
                            AC1StdDev: {
                                type: "number"
                            },
                            AC1Count: {
                                type: "number"
                            },
                            AC2Avg: {
                                type: "number"
                            },
                            AC2Max: {
                                type: "number"
                            },
                            AC2Min: {
                                type: "number"
                            },
                            AC2StdDev: {
                                type: "number"
                            },
                            AC2Count: {
                                type: "number"
                            },
                            C1Avg: {
                                type: "number"
                            },
                            C1Max: {
                                type: "number"
                            },
                            C1Min: {
                                type: "number"
                            },
                            C1StdDev: {
                                type: "number"
                            },
                            C1Count: {
                                type: "number"
                            },
                            C2Avg: {
                                type: "number"
                            },
                            C2Max: {
                                type: "number"
                            },
                            C2Min: {
                                type: "number"
                            },
                            C2StdDev: {
                                type: "number"
                            },
                            C2Count: {
                                type: "number"
                            },
                            C3Avg: {
                                type: "number"
                            },
                            C3Max: {
                                type: "number"
                            },
                            C3Min: {
                                type: "number"
                            },
                            C3StdDev: {
                                type: "number"
                            },
                            C3Count: {
                                type: "number"
                            },
                            D1Avg: {
                                type: "number"
                            },
                            D1Max: {
                                type: "number"
                            },
                            D1Min: {
                                type: "number"
                            },
                            D1StdDev: {
                                type: "number"
                            },
                            M1_1Avg: {
                                type: "number"
                            },
                            M1_1Max: {
                                type: "number"
                            },
                            M1_1Min: {
                                type: "number"
                            },
                            M1_1StdDev: {
                                type: "number"
                            },
                            M1_1Count: {
                                type: "number"
                            },
                            M1_2Avg: {
                                type: "number"
                            },
                            M1_2Max: {
                                type: "number"
                            },
                            M1_2Min: {
                                type: "number"
                            },
                            M1_2StdDev: {
                                type: "number"
                            },
                            M1_2Count: {
                                type: "number"
                            },
                            M1_3Avg: {
                                type: "number"
                            },
                            M1_3Max: {
                                type: "number"
                            },
                            M1_3Min: {
                                type: "number"
                            },
                            M1_3StdDev: {
                                type: "number"
                            },
                            M1_3Count: {
                                type: "number"
                            },
                            M1_4Avg: {
                                type: "number"
                            },
                            M1_4Max: {
                                type: "number"
                            },
                            M1_4Min: {
                                type: "number"
                            },
                            M1_4StdDev: {
                                type: "number"
                            },
                            M1_4Count: {
                                type: "number"
                            },
                            M1_5Avg: {
                                type: "number"
                            },
                            M1_5Max: {
                                type: "number"
                            },
                            M1_5Min: {
                                type: "number"
                            },
                            M1_5StdDev: {
                                type: "number"
                            },
                            M1_5Count: {
                                type: "number"
                            },
                            M2_1Avg: {
                                type: "number"
                            },
                            M2_1Max: {
                                type: "number"
                            },
                            M2_1Min: {
                                type: "number"
                            },
                            M2_1StdDev: {
                                type: "number"
                            },
                            M2_1Count: {
                                type: "number"
                            },
                            M2_2Avg: {
                                type: "number"
                            },
                            M2_2Max: {
                                type: "number"
                            },
                            M2_2Min: {
                                type: "number"
                            },
                            M2_2StdDev: {
                                type: "number"
                            },
                            M2_2Count: {
                                type: "number"
                            },
                            M2_3Avg: {
                                type: "number"
                            },
                            M2_3Max: {
                                type: "number"
                            },
                            M2_3Min: {
                                type: "number"
                            },
                            M2_3StdDev: {
                                type: "number"
                            },
                            M2_3Count: {
                                type: "number"
                            },
                            M2_4Avg: {
                                type: "number"
                            },
                            M2_4Max: {
                                type: "number"
                            },
                            M2_4Min: {
                                type: "number"
                            },
                            M2_4StdDev: {
                                type: "number"
                            },
                            M2_4Count: {
                                type: "number"
                            },
                            M2_5Avg: {
                                type: "number"
                            },
                            M2_5Max: {
                                type: "number"
                            },
                            M2_5Min: {
                                type: "number"
                            },
                            M2_5StdDev: {
                                type: "number"
                            },
                            M2_5Count: {
                                type: "number"
                            },
                            M2_6Avg: {
                                type: "number"
                            },
                            M2_6Max: {
                                type: "number"
                            },
                            M2_6Min: {
                                type: "number"
                            },
                            M2_6StdDev: {
                                type: "number"
                            },
                            M2_6Count: {
                                type: "number"
                            },
                            M2_7Avg: {
                                type: "number"
                            },
                            M2_7Max: {
                                type: "number"
                            },
                            M2_7Min: {
                                type: "number"
                            },
                            M2_7StdDev: {
                                type: "number"
                            },
                            M2_7Count: {
                                type: "number"
                            },
                            M2_8Avg: {
                                type: "number"
                            },
                            M2_8Max: {
                                type: "number"
                            },
                            M2_8Min: {
                                type: "number"
                            },
                            M2_8StdDev: {
                                type: "number"
                            },
                            M2_8Count: {
                                type: "number"
                            },
                            VAvg: {
                                type: "number"
                            },
                            VMax: {
                                type: "number"
                            },
                            VMin: {
                                type: "number"
                            },
                            IAvg: {
                                type: "number"
                            },
                            IMax: {
                                type: "number"
                            },
                            IMin: {
                                type: "number"
                            },
                            T: {
                                type: "number"
                            },
                            addr: {
                                type: "number"
                            },
                        }
                    },
                    data: function(res) {
                        if (!app.isFine(res)) {
                            dbr.mettowervis(false);
                            return;
                        }

                        dbr.mettowervis(false);
                        return res.data.Data
                    },
                    total: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totaldatamet').html(kendo.toString(res.data.Total, 'n0'));
                        return res.data.Total;
                    }
                },
                sort: [{
                    field: 'TimeStamp',
                    dir: 'asc'
                }],
            },
            // toolbar: [
            //     "excel"
            // ],
            // excel: {
            //     fileName: "Permanent Met Tower Data.xlsx",
            //     filterable: true,
            //      allPages: true
            // },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            filterable: {
                extra: false,
                operators: {
                    string: {
                        eq: "Is equal to",
                        neq: "Is not equal to",
                        gt: "Is greater than",
                        gte: "Is greater than or equal to",
                        lt: "Is less than",
                        lte: "Is less than or equal to"
                    }
                }
            },
            filterMenuInit: function(e) {
                e.container.data("kendoPopup").bind("open", function() {
                    // console.log(e.container);
                    if (e.container.is(".k-grid-filter")) {
                        e.container.find("form").removeClass("k-state-border-up");
                        e.container.find("form").addClass("k-state-border-down");
                        // console.log(e.container[0]);
                    } else {
                        // console.log("test");
                    }


                });
            },
            pageable: true,
            //resizable: true,
            columns: [{
                    title: "Date",
                    field: "TimeStamp",
                    template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
                    width: 150,
                    locked: true,
                    filterable: false
                },

                {
                    title: "V Hub</br>WS 90m Avg",
                    field: "VHubWS90mAvg",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Hub</br>WS 90m Max",
                    field: "VHubWS90mMax",
                    format: "{0:n2}",
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Hub</br>WS 90m Min",
                    field: "VHubWS90mMin",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Hub</br>WS 90m Std Dev",
                    field: "VHubWS90mStdDev",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Hub</br>WS 90m Count",
                    field: "VHubWS90mCount",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "V Ref</br>WS 88m Avg",
                    field: "VRefWS88mAvg",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Ref</br>WS 88m Max",
                    field: "VRefWS88mMax",
                    format: "{0:n2}",
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Ref</br>WS 88m Min",
                    field: "VRefWS88mMin",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Ref</br>WS 88m Std Dev",
                    field: "VRefWS88mStdDev",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Ref</br>WS 88m Count",
                    field: "VRefWS88mCount",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "V Tip</br>WS 42m Avg",
                    field: "VTipWS42mAvg",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Tip</br>WS 42m Max",
                    field: "VTipWS42mMax",
                    format: "{0:n2}",
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Tip</br>WS 42m Min",
                    field: "VTipWS42mMin",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Tip</br>WS 42m Std Dev",
                    field: "VTipWS42mStdDev",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Tip</br>WS 42m Count",
                    field: "VTipWS42mCount",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "D Hub</br>WD 88m Avg",
                    field: "DHubWD88mAvg",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Hub</br>WD 88m Max",
                    field: "DHubWD88mMax",
                    format: "{0:n2}",
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Hub</br>WD 88m Min",
                    field: "DHubWD88mMin",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Hub</br>WD 88m Std Dev",
                    field: "DHubWD88mStdDev",
                    format: "{0:n2}",
                    width: 110,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Hub</br>WD 88m Count",
                    field: "DHubWD88mCount",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "D Ref</br>WD 86m Avg",
                    field: "DRefWD86mAvg",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Ref</br>WD 86m Max",
                    field: "DRefWD86mMax",
                    format: "{0:n2}",
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Ref</br>WD 86m Min",
                    field: "DRefWD86mMin",
                    format: "{0:n2}",
                    width: 80,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Ref</br>WD 86m Std Dev",
                    field: "DRefWD86mStdDev",
                    format: "{0:n2}",
                    width: 110,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D Ref</br>WD 86m Count",
                    field: "DRefWD86mCount",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "T Hub & H Hub</br>Humid 85m Avg",
                    format: "{0:n2}",
                    field: "THubHHubHumid855mAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Humid 85m Max",
                    format: "{0:n2}",
                    field: "THubHHubHumid855mMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Humid 85m Min",
                    format: "{0:n2}",
                    field: "THubHHubHumid855mMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Humid 85m Std Dev",
                    format: "{0:n2}",
                    field: "THubHHubHumid855mStdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Humid 85m Count",
                    format: "{0:n2}",
                    field: "THubHHubHumid855mCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "T Ref & H Ref</br>Humid 85.5m Avg",
                    format: "{0:n2}",
                    field: "TRefHRefHumid855mAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Humid 85.5m Max",
                    format: "{0:n2}",
                    field: "TRefHRefHumid855mMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Humid 85.5m Min",
                    format: "{0:n2}",
                    field: "TRefHRefHumid855mMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Humid 85.5m Std Dev",
                    format: "{0:n2}",
                    field: "TRefHRefHumid855mStdDev",
                    width: 130,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Humid 85.5m Count",
                    format: "{0:n2}",
                    field: "TRefHRefHumid855mCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "T Hub & H Hub</br>Temp 85.5m Avg",
                    format: "{0:n2}",
                    field: "THubHHubTemp855mAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Temp 85.5m Max",
                    format: "{0:n2}",
                    field: "THubHHubTemp855mMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Temp 85.5m Min",
                    format: "{0:n2}",
                    field: "THubHHubTemp855mMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Temp 85.5m Std Dev",
                    format: "{0:n2}",
                    field: "THubHHubTemp855mStdDev",
                    width: 130,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Hub & H Hub</br>Temp 85.5m Count",
                    format: "{0:n2}",
                    field: "THubHHubTemp855mCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "T Ref & H Ref</br>Temp 85.5 Avg",
                    format: "{0:n2}",
                    field: "TRefHRefTemp855mAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Temp 85.5 Max",
                    format: "{0:n2}",
                    field: "TRefHRefTemp855mMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Temp 85.5 Min",
                    format: "{0:n2}",
                    field: "TRefHRefTemp855mMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Temp 85.5 Std Dev",
                    format: "{0:n2}",
                    field: "TRefHRefTemp855mStdDev",
                    width: 130,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T Ref & H Ref</br>Temp 85.5 Count",
                    format: "{0:n2}",
                    field: "TRefHRefTemp855mCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "Baro Air Pressure</br>85.5m Avg",
                    format: "{0:n2}",
                    field: "BaroAirPress855mAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Baro Air Pressure</br>85.5m Max",
                    format: "{0:n2}",
                    field: "BaroAirPress855mMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Baro Air Pressure</br>85.5m Min",
                    format: "{0:n2}",
                    field: "BaroAirPress855mMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Baro Air Pressure</br>85.5m Std Dev",
                    format: "{0:n2}",
                    field: "BaroAirPress855mStdDev",
                    width: 130,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Baro Air Pressure</br>85.5m Count",
                    format: "{0:n2}",
                    field: "BaroAirPress855mCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

                {
                    title: "Yaw angle Voltage Avg",
                    format: "{0:n2}",
                    field: "YawAngleVoltageAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Yaw angle Voltage Max",
                    format: "{0:n2}",
                    field: "YawAngleVoltageMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Yaw angle Voltage Min",
                    format: "{0:n2}",
                    field: "YawAngleVoltageMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Yaw angle Voltage StdDev",
                    format: "{0:n2}",
                    field: "YawAngleVoltageStdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Yaw angle Voltage Count",
                    format: "{0:n2}",
                    field: "YawAngleVoltageCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI1 Avg",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI1 Max",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI1 Min",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI1 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI1 Count",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI2 Avg",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI2 Max",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI2 Min",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI2 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI2 Count",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI3 Avg",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI3 Max",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI3 Min",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI3 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI3 Count",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI4 Avg",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI4Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI4 Max",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI4Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI4 Min",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI4Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI4 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI4StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor Voltage AI4 Count",
                    format: "{0:n2}",
                    field: "OtherSensorVoltageAI4Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Gen RPM Current Avg",
                    format: "{0:n2}",
                    field: "GenRPMCurrentAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Gen RPM Current Max",
                    format: "{0:n2}",
                    field: "GenRPMCurrentMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Gen RPM Current Min",
                    format: "{0:n2}",
                    field: "GenRPMCurrentMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Gen RPM Current StdDev",
                    format: "{0:n2}",
                    field: "GenRPMCurrentStdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Gen RPM Current Count",
                    format: "{0:n2}",
                    field: "GenRPMCurrentCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "ws SCS Current Avg",
                    format: "{0:n2}",
                    field: "WS_SCSCurrentAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "ws SCS Current Max",
                    format: "{0:n2}",
                    field: "WS_SCSCurrentMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "ws SCS Current Min",
                    format: "{0:n2}",
                    field: "WS_SCSCurrentMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "ws SCS Current StdDev",
                    format: "{0:n2}",
                    field: "WS_SCSCurrentStdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "ws SCS Current Count",
                    format: "{0:n2}",
                    field: "WS_SCSCurrentCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Rain Status Count",
                    format: "{0:n2}",
                    field: "RainStatusCount",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Rain Status Sum",
                    format: "{0:n2}",
                    field: "RainStatusSum",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO1 Avg",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO1 Max",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO1 Min",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO1 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO1 Count",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO2 Avg",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO2 Max",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO2 Min",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO2 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO2 Count",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO3 Avg",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO3 Max",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO3 Min",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO3 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO3 Count",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO4 Avg",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO4Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO4 Max",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO4Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO4 Min",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO4Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO4 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO4StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO4 Count",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO4Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO5 Avg",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO5Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO5 Max",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO5Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO5 Min",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO5Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO5 StdDev",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO5StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Other Sensor 2 Status IO5 Count",
                    format: "{0:n2}",
                    field: "OtherSensor2StatusIO5Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A1 Avg",
                    format: "{0:n2}",
                    field: "A1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A1 Max",
                    format: "{0:n2}",
                    field: "A1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A1 Min",
                    format: "{0:n2}",
                    field: "A1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A1 StdDev",
                    format: "{0:n2}",
                    field: "A1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A1 Count",
                    format: "{0:n2}",
                    field: "A1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A2 Avg",
                    format: "{0:n2}",
                    field: "A2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A2 Max",
                    format: "{0:n2}",
                    field: "A2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A2 Min",
                    format: "{0:n2}",
                    field: "A2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A2 StdDev",
                    format: "{0:n2}",
                    field: "A2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A2 Count",
                    format: "{0:n2}",
                    field: "A2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A3 Avg",
                    format: "{0:n2}",
                    field: "A3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A3 Max",
                    format: "{0:n2}",
                    field: "A3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A3 Min",
                    format: "{0:n2}",
                    field: "A3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A3 StdDev",
                    format: "{0:n2}",
                    field: "A3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A3 Count",
                    format: "{0:n2}",
                    field: "A3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A4 Avg",
                    format: "{0:n2}",
                    field: "A4Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A4 Max",
                    format: "{0:n2}",
                    field: "A4Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A4 Min",
                    format: "{0:n2}",
                    field: "A4Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A4 StdDev",
                    format: "{0:n2}",
                    field: "A4StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A4 Count",
                    format: "{0:n2}",
                    field: "A4Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A5 Avg",
                    format: "{0:n2}",
                    field: "A5Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A5 Max",
                    format: "{0:n2}",
                    field: "A5Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A5 Min",
                    format: "{0:n2}",
                    field: "A5Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A5 StdDev",
                    format: "{0:n2}",
                    field: "A5StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A5 Count",
                    format: "{0:n2}",
                    field: "A5Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A6 Avg",
                    format: "{0:n2}",
                    field: "A6Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A6 Max",
                    format: "{0:n2}",
                    field: "A6Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A6 Min",
                    format: "{0:n2}",
                    field: "A6Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A6 StdDev",
                    format: "{0:n2}",
                    field: "A6StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A6 Count",
                    format: "{0:n2}",
                    field: "A6Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A7 Avg",
                    format: "{0:n2}",
                    field: "A7Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A7 Max",
                    format: "{0:n2}",
                    field: "A7Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A7 Min",
                    format: "{0:n2}",
                    field: "A7Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A7 StdDev",
                    format: "{0:n2}",
                    field: "A7StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A7 Count",
                    format: "{0:n2}",
                    field: "A7Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A8 Avg",
                    format: "{0:n2}",
                    field: "A8Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A8 Max",
                    format: "{0:n2}",
                    field: "A8Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A8 Min",
                    format: "{0:n2}",
                    field: "A8Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A8 StdDev",
                    format: "{0:n2}",
                    field: "A8StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A8 Count",
                    format: "{0:n2}",
                    field: "A8Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A9 Avg",
                    format: "{0:n2}",
                    field: "A9Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A9 Max",
                    format: "{0:n2}",
                    field: "A9Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A9 Min",
                    format: "{0:n2}",
                    field: "A9Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A9 StdDev",
                    format: "{0:n2}",
                    field: "A9StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A9 Count",
                    format: "{0:n2}",
                    field: "A9Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A10 Avg",
                    format: "{0:n2}",
                    field: "A10Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A10 Max",
                    format: "{0:n2}",
                    field: "A10Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A10 Min",
                    format: "{0:n2}",
                    field: "A10Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A10 StdDev",
                    format: "{0:n2}",
                    field: "A10StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "A10 Count",
                    format: "{0:n2}",
                    field: "A10Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC1 Avg",
                    format: "{0:n2}",
                    field: "AC1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC1 Max",
                    format: "{0:n2}",
                    field: "AC1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC1 Min",
                    format: "{0:n2}",
                    field: "AC1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC1 StdDev",
                    format: "{0:n2}",
                    field: "AC1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC1 Count",
                    format: "{0:n2}",
                    field: "AC1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC2 Avg",
                    format: "{0:n2}",
                    field: "AC2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC2 Max",
                    format: "{0:n2}",
                    field: "AC2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC2 Min",
                    format: "{0:n2}",
                    field: "AC2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC2 StdDev",
                    format: "{0:n2}",
                    field: "AC2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "AC2 Count",
                    format: "{0:n2}",
                    field: "AC2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C1 Avg",
                    format: "{0:n2}",
                    field: "C1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C1 Max",
                    format: "{0:n2}",
                    field: "C1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C1 Min",
                    format: "{0:n2}",
                    field: "C1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C1 StdDev",
                    format: "{0:n2}",
                    field: "C1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C1 Count",
                    format: "{0:n2}",
                    field: "C1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C2 Avg",
                    format: "{0:n2}",
                    field: "C2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C2 Max",
                    format: "{0:n2}",
                    field: "C2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C2 Min",
                    format: "{0:n2}",
                    field: "C2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C2 StdDev",
                    format: "{0:n2}",
                    field: "C2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C2 Count",
                    format: "{0:n2}",
                    field: "C2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C3 Avg",
                    format: "{0:n2}",
                    field: "C3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C3 Max",
                    format: "{0:n2}",
                    field: "C3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C3 Min",
                    format: "{0:n2}",
                    field: "C3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C3 StdDev",
                    format: "{0:n2}",
                    field: "C3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "C3 Count",
                    format: "{0:n2}",
                    field: "C3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D1 Avg",
                    format: "{0:n2}",
                    field: "D1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D1 Max",
                    format: "{0:n2}",
                    field: "D1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D1 Min",
                    format: "{0:n2}",
                    field: "D1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "D1 StdDev",
                    format: "{0:n2}",
                    field: "D1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 1 Avg",
                    format: "{0:n2}",
                    field: "M1_1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 1 Max",
                    format: "{0:n2}",
                    field: "M1_1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 1 Min",
                    format: "{0:n2}",
                    field: "M1_1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 1 StdDev",
                    format: "{0:n2}",
                    field: "M1_1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 1 Count",
                    format: "{0:n2}",
                    field: "M1_1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 2 Avg",
                    format: "{0:n2}",
                    field: "M1_2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 2 Max",
                    format: "{0:n2}",
                    field: "M1_2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 2 Min",
                    format: "{0:n2}",
                    field: "M1_2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 2 StdDev",
                    format: "{0:n2}",
                    field: "M1_2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 2 Count",
                    format: "{0:n2}",
                    field: "M1_2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 3 Avg",
                    format: "{0:n2}",
                    field: "M1_3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 3 Max",
                    format: "{0:n2}",
                    field: "M1_3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 3 Min",
                    format: "{0:n2}",
                    field: "M1_3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 3 StdDev",
                    format: "{0:n2}",
                    field: "M1_3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 3 Count",
                    format: "{0:n2}",
                    field: "M1_3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 4 Avg",
                    format: "{0:n2}",
                    field: "M1_4Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 4 Max",
                    format: "{0:n2}",
                    field: "M1_4Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 4 Min",
                    format: "{0:n2}",
                    field: "M1_4Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 4 StdDev",
                    format: "{0:n2}",
                    field: "M1_4StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 4 Count",
                    format: "{0:n2}",
                    field: "M1_4Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 5 Avg",
                    format: "{0:n2}",
                    field: "M1_5Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 5 Max",
                    format: "{0:n2}",
                    field: "M1_5Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 5 Min",
                    format: "{0:n2}",
                    field: "M1_5Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 5 StdDev",
                    format: "{0:n2}",
                    field: "M1_5StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M1 5 Count",
                    format: "{0:n2}",
                    field: "M1_5Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 1 Avg",
                    format: "{0:n2}",
                    field: "M2_1Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 1 Max",
                    format: "{0:n2}",
                    field: "M2_1Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 1 Min",
                    format: "{0:n2}",
                    field: "M2_1Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 1 StdDev",
                    format: "{0:n2}",
                    field: "M2_1StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 1 Count",
                    format: "{0:n2}",
                    field: "M2_1Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 2 Avg",
                    format: "{0:n2}",
                    field: "M2_2Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 2 Max",
                    format: "{0:n2}",
                    field: "M2_2Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 2 Min",
                    format: "{0:n2}",
                    field: "M2_2Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 2 StdDev",
                    format: "{0:n2}",
                    field: "M2_2StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 2 Count",
                    format: "{0:n2}",
                    field: "M2_2Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 3 Avg",
                    format: "{0:n2}",
                    field: "M2_3Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 3 Max",
                    format: "{0:n2}",
                    field: "M2_3Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 3 Min",
                    format: "{0:n2}",
                    field: "M2_3Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 3 StdDev",
                    format: "{0:n2}",
                    field: "M2_3StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 3 Count",
                    format: "{0:n2}",
                    field: "M2_3Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 4 Avg",
                    format: "{0:n2}",
                    field: "M2_4Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 4 Max",
                    format: "{0:n2}",
                    field: "M2_4Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 4 Min",
                    format: "{0:n2}",
                    field: "M2_4Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 4 StdDev",
                    format: "{0:n2}",
                    field: "M2_4StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 4 Count",
                    format: "{0:n2}",
                    field: "M2_4Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 5 Avg",
                    format: "{0:n2}",
                    field: "M2_5Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 5 Max",
                    format: "{0:n2}",
                    field: "M2_5Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 5 Min",
                    format: "{0:n2}",
                    field: "M2_5Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 5 StdDev",
                    format: "{0:n2}",
                    field: "M2_5StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 5 Count",
                    format: "{0:n2}",
                    field: "M2_5Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 6 Avg",
                    format: "{0:n2}",
                    field: "M2_6Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 6 Max",
                    format: "{0:n2}",
                    field: "M2_6Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 6 Min",
                    format: "{0:n2}",
                    field: "M2_6Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 6 StdDev",
                    format: "{0:n2}",
                    field: "M2_6StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 6 Count",
                    format: "{0:n2}",
                    field: "M2_6Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 7 Avg",
                    format: "{0:n2}",
                    field: "M2_7Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 7 Max",
                    format: "{0:n2}",
                    field: "M2_7Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 7 Min",
                    format: "{0:n2}",
                    field: "M2_7Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 7 StdDev",
                    format: "{0:n2}",
                    field: "M2_7StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 7 Count",
                    format: "{0:n2}",
                    field: "M2_7Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 8 Avg",
                    format: "{0:n2}",
                    field: "M2_8Avg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 8 Max",
                    format: "{0:n2}",
                    field: "M2_8Max",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 8 Min",
                    format: "{0:n2}",
                    field: "M2_8Min",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 8 StdDev",
                    format: "{0:n2}",
                    field: "M2_8StdDev",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "M2 8 Count",
                    format: "{0:n2}",
                    field: "M2_8Count",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Avg",
                    format: "{0:n2}",
                    field: "VAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Max",
                    format: "{0:n2}",
                    field: "VMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "V Min",
                    format: "{0:n2}",
                    field: "VMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "I Avg",
                    format: "{0:n2}",
                    field: "IAvg",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "I Max",
                    format: "{0:n2}",
                    field: "IMax",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "I Min",
                    format: "{0:n2}",
                    field: "IMin",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "T",
                    format: "{0:n2}",
                    field: "T",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "addr",
                    format: "{0:n2}",
                    field: "addr",
                    width: 120,
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                },

            ]
        });
    },
    InitGridJMR: function() {
        dbr.jmrvis(true);
        var turbine = [];
        if ($("#turbineMulti").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
            turbine = turbineval;
        } else {
            turbine = $("#turbineMulti").data("kendoMultiSelect").value();
        }

        var dateStart = kendo.toString($('#dateStart').data('kendoDatePicker').value(), "yyyyMM");
        var dateEnd = kendo.toString($('#dateEnd').data('kendoDatePicker').value(), "yyyyMM");

        var monthId = [];

        if (dateStart != dateEnd) {
            var dateStartInt = parseInt(dateStart);
            var dateEndInt = parseInt(dateEnd);
            var dsYear = parseInt(dateStart.substring(0, 4));
            var dsMonth = parseInt(dateStart.substring(4, 6));
            var deYear = parseInt(dateEnd.substring(0, 4));
            var deMonth = parseInt(dateEnd.substring(4, 6));
            var exit = false;

            monthId.push(dateStartInt);

            do {
                if (dateStartInt < dateEndInt) {
                    if (dsMonth < 12) {
                        dsMonth++;
                    } else {
                        dsYear++;
                        dsMonth = 1;
                    }

                    if (dsMonth > 9) {
                        dateStartInt = parseInt(dsYear + "" + dsMonth)
                    } else {
                        dateStartInt = parseInt(dsYear + "0" + dsMonth)
                    }

                    monthId.push(dateStartInt);
                } else {
                    exit = true;
                }
            } while (exit == false);
        } else {
            monthId.push(parseInt(dateStart));
        }

        var filters = [{
            field: "dateinfo.monthid",
            operator: "in",
            value: monthId
        }, {
            field: "sections.turbine",
            operator: "in",
            value: turbine
        }, ];

        dbr.filterJMR(filters);

        var filter = {
            filters: filters
        }
        var param = {
            filter: filter
        };

        $('#dataGridJMR').html("");
        $('#dataGridJMR').kendoGrid({
            dataSource: {
                serverPaging: true,
                serverSorting: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getjmrlist",
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
                        if (!app.isFine(res)) {
                            dbr.jmrvis(false);
                            return;
                        }
                        dbr.jmrvis(false);
                        return res.data.Data
                    },
                    total: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totaldatajmr').html(kendo.toString(res.data.Total, 'n0'));
                        return res.data.Total;
                    }
                },
                sort: [{
                    field: 'DateInfo.DateId',
                    dir: 'asc'
                }, ],
            },
            groupable: false,
            sortable: true,
            filterable: false,
            pageable: true,
            detailInit: Data.InitJMRDetail,
            columns: [{
                title: "Month",
                field: "DateInfo.DateId",
                attributes: {
                    style: "text-align: center"
                },
                template: "#= kendo.toString(moment.utc(DateInfo.DateId).format('MMMM YYYY'), 'dd-MMM-yyyy') #"
            }, {
                title: "Description",
                field: "Description"
            }, ]
        });
    },

    // INIT GRID EXCEPTION
    InitGridExceptionTimeDuration: function() {
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

        var param = {
            DateStart: dateStart,
            DateEnd: dateEnd,
            Turbine: turbine,
        };


        $('#dataGridExceptionTimeDuration').html("");
        $('#dataGridExceptionTimeDuration').kendoGrid({
            dataSource: {
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getscadalist",
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
                    model: {
                        fields: {
                            AlarmOkTime: {
                                type: "number"
                            },
                            OkTime: {
                                type: "number"
                            },
                            Power: {
                                type: "number"
                            },
                            PowerLost: {
                                type: "number"
                            },
                        }
                    },
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
                        $('#totalpowerExceptionTimeDuration').html(kendo.toString(res.data.TotalPower / 1000, 'n2') + ' MW');
                        $('#totalpowerlostExceptionTimeDuration').html(kendo.toString(res.data.TotalPowerLost / 1000, 'n2') + ' MW');
                        $('#totalturbineExceptionTimeDuration').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataExceptionTimeDuration').html(kendo.toString(res.data.Total, 'n0'));

                        // $('#totprodExceptionTimeDuration').html(kendo.toString(res.data.totalProduction / 1000, 'n2') + ' MWh');
                        // $('#avgwindspeedExceptionTimeDuration').html(kendo.toString(res.data.avgWindSpeed, 'n2') + ' m/s');

                        $('#totprodExceptionTimeDuration').html(kendo.toString(res.data.TotalProduction / 1000, 'n2') + ' MWh');
                        $('#avgwindspeedExceptionTimeDuration').html(kendo.toString(res.data.AvgWindSpeed, 'n2') + ' m/s');
                        return res.data.Total;
                    }
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
            //     "excel"
            // ],
            excel: {
                fileName: "Scada Exception Time Duration.xlsx",
                filterable: true,
                allPages: true
            },
            selectable: "multiple",
            groupable: false,
            sortable: true,
            filterable: {
                extra: false,
                operators: {
                    string: {
                        eq: "Is equal to",
                        neq: "Is not equal to",
                        gt: "Is greater than",
                        gte: "Is greater than or equal to",
                        lt: "Is less than",
                        lte: "Is less than or equal to"
                    }
                }
            },
            pageable: true,
            //resizable: true,

            columns: [{
                title: "Date",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                width: 80,
                locked: true,
                filterable: false
            }, {
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                locked: true,
                filterable: false
            }, {
                title: "Start Time",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 65,
                locked: true,
                attributes: {
                    style: "text-align:center;"
                },
                filterable: false
            }, {
                title: "Grid Frequency",
                field: "GridFrequency",
                template: '# if (GridFrequency < -99998 ) { # - # } else {#' + '#: GridFrequency #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                filterable: false
            }, {
                title: "Reactive Power",
                field: "ReactivePower",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                filterable: false
            }, {
                title: "Alarm",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "Alarm Ext Stop Time",
                    field: "AlarmExtStopTime",
                    width: 90,
                    template: '# if (AlarmExtStopTime < -99998 ) { # - # } else {#' + '#: AlarmExtStopTime #' + '#}#',
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Alarm Grid Down Time",
                    field: "AlarmGridDownTime",
                    template: '# if (AlarmGridDownTime < -99998 ) { # - # } else {#' + '#: AlarmGridDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Inter Line Down",
                    field: "AlarmInterLineDown",
                    template: '# if (AlarmInterLineDown < -99998 ) { # - # } else {#' + '#: AlarmInterLineDown #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Mach Down Time",
                    field: "AlarmMachDownTime",
                    template: '# if (AlarmMachDownTime < -99998 ) { # - # } else {#' + '#: AlarmMachDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm OK Time",
                    field: "AlarmOkTime",
                    template: '# if (AlarmOkTime < -99998 ) { # - # } else {#' + '#: AlarmOkTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: true,
                    headerAttributes: {
                        class: 'gridAlarmOkTime'
                    }
                }, {
                    title: "Alarm Unknown Time",
                    field: "AlarmUnknownTime",
                    template: '# if (AlarmUnknownTime < -99998 ) { # - # } else {#' + '#: AlarmUnknownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Weather Stop",
                    field: "AlarmWeatherStop",
                    template: '# if (AlarmWeatherStop < -99998 ) { # - # } else {#' + '#: AlarmWeatherStop #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }]
            }, {
                title: "Time",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "Ext Stop Time",
                    field: "ExternalStopTime",
                    template: '# if (ExternalStopTime < -99998 ) { # - # } else {#' + '#: ExternalStopTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Grid Down Time",
                    field: "GridDownTime",
                    template: '# if (GridDownTime < -99998 ) { # - # } else {#' + '#: GridDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Grid OK Secs",
                    field: "GridOkSecs",
                    template: '# if (GridOkSecs < -99998 ) { # - # } else {#' + '#: GridOkSecs #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Internal Line Down",
                    field: "InternalLineDown",
                    template: '# if (InternalLineDown < -99998 ) { # - # } else {#' + '#: InternalLineDown #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Machine Down Time",
                    field: "MachineDownTime",
                    template: '# if (MachineDownTime < -99998 ) { # - # } else {#' + '#: MachineDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "OK Secs",
                    field: "OkSecs",
                    template: '# if (OkSecs < -99998 ) { # - # } else {#' + '#: OkSecs #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "OK Time",
                    field: "OkTime",
                    template: '# if (OkTime < -99998 ) { # - # } else {#' + '#: OkTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: true
                }, {
                    title: "Unknown Time",
                    field: "UnknownTime",
                    template: '# if (UnknownTime < -99998 ) { # - # } else {#' + '#: UnknownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }]
            }, {
                title: "Total Time",
                field: "TotalTime",
                template: '# if (TotalTime < -99998 ) { # - # } else {#' + '#: TotalTime #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Generator RPM",
                field: "GeneratorRPM",
                template: '# if (GeneratorRPM < -99998 ) { # - # } else {#' + '#: GeneratorRPM #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacelle Yaw Position Untwist",
                field: "NacelleYawPositionUntwist",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacelle Temperature",
                field: "NacelleTemperature",
                template: '# if (NacelleTemperature < -99998 ) { # - # } else {#' + '#: NacelleTemperature #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Adj Wind Speed",
                field: "AdjWindSpeed",
                template: '# if (AdjWindSpeed < -99998 ) { # - # } else {#' + '#: AdjWindSpeed #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ambient Temperature",
                field: "AmbientTemperature",
                template: '# if (AmbientTemperature < -99998 ) { # - # } else {#' + '#: AmbientTemperature #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Avg Blade Angle",
                field: "AvgBladeAngle",
                template: '# if (AvgBladeAngle < -99998 ) { # - # } else {#' + '#: AvgBladeAngle #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Avg Wind Speed",
                field: "AvgWindSpeed",
                template: '# if (AvgWindSpeed < -99998 ) { # - # } else {#' + '#: AvgWindSpeed #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Units Generated",
                field: "UnitsGenerated",
                template: '# if (UnitsGenerated < -99998 ) { # - # } else {#' + '#: UnitsGenerated #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Estimated Power",
                field: "EstimatedPower",
                template: '# if (EstimatedPower < -99998 ) { # - # } else {#' + '#: EstimatedPower #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacel Direction",
                field: "NacelDirection",
                template: '# if (NacelDirection < -99998 ) { # - # } else {#' + '#: NacelDirection #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Power",
                field: "Power",
                template: '# if (Power < -99998 ) { # - # } else {#' + '#: Power #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: true
            }, {
                title: "Power Lost",
                field: "PowerLost",
                template: '# if (PowerLost < -99998 ) { # - # } else {#' + '#: PowerLost #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: true,
                headerAttributes: {
                    class: 'gridPowerLost'
                }
            }, {
                title: "Rotor RPM",
                field: "RotorRPM",
                template: '# if (RotorRPM < -99998 ) { # - # } else {#' + '#: RotorRPM #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Wind Direction",
                field: "WindDirection",
                template: '# if (WindDirection < -99998 ) { # - # } else {#' + '#: WindDirection #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }]
        });
        var grid = $('#dataGridExceptionTimeDuration').data('kendoGrid');
        var columns = grid.columns;
        dbr.gridColumnsScadaException([]);
        $.each(columns, function(i, val) {
            $('#dataGridExceptionTimeDuration').data("kendoGrid").showColumn(val.field);
            var result = {
                field: val.field,
                title: val.title,
                value: true
            }
            dbr.gridColumnsScadaException.push(result);
        });
    },
    InitGridAnomalies: function() {
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

        $('#dataGridAnomalies').html("");
        $('#dataGridAnomalies').kendoGrid({
            selectable: "multiple",
            dataSource: {
                filter: [{
                    field: "timestamp",
                    operator: "gte",
                    value: dateStart
                }, {
                    field: "timestamp",
                    operator: "lte",
                    value: dateEnd
                }, {
                    field: "isvalidtimeduration",
                    operator: "eq",
                    value: true
                }, {
                    field: "turbine",
                    operator: "in",
                    value: turbine
                }],
                serverPaging: true,
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getscadaanomalylist",
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
                    model: {
                        fields: {
                            AlarmOkTime: {
                                type: "number"
                            },
                            OkTime: {
                                type: "number"
                            },
                            Power: {
                                type: "number"
                            },
                            PowerLost: {
                                type: "number"
                            },
                        }
                    },
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
                        $('#totalpowerAnomalies').html(kendo.toString(res.data.TotalPower / 1000, 'n2') + ' MW');
                        $('#totalpowerlostAnomalies').html(kendo.toString(res.data.TotalPowerLost / 1000, 'n2') + ' MW');
                        $('#totalturbineAnomalies').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataAnomalies').html(kendo.toString(res.data.Total, 'n0'));

                        $('#totprodAnomalies').html(kendo.toString(res.data.TotalProduction / 1000, 'n2') + ' MWh');
                        $('#avgwindspeedAnomalies').html(kendo.toString(res.data.AvgWindSpeed, 'n2') + ' m/s');
                        return res.data.Total;
                    }
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
                fileName: "Scada Anomaly.xlsx",
                filterable: true,
                allPages: true
            },
            groupable: false,
            sortable: true,
            filterable: {
                extra: false,
                operators: {
                    string: {
                        eq: "Is equal to",
                        neq: "Is not equal to",
                        gt: "Is greater than",
                        gte: "Is greater than or equal to",
                        lt: "Is less than",
                        lte: "Is less than or equal to"
                    }
                }
            },
            pageable: true,
            resizable: true,

            columns: [{
                title: "Date",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                width: 80,
                locked: true,
                filterable: false
            }, {
                title: "Turbine",
                field: "Turbine",
                attributes: {
                    class: "align-center"
                },
                width: 90,
                locked: true,
                filterable: false
            }, {
                title: "Start Time",
                field: "TimeStamp",
                template: "#= kendo.toString(moment.utc(TimeStamp).format('HH:mm:ss'), 'HH:mm:ss') #",
                width: 65,
                locked: true,
                attributes: {
                    style: "text-align:center;"
                },
                filterable: false
            }, {
                title: "Grid Frequency",
                field: "GridFrequency",
                template: '# if (GridFrequency < -99998 ) { # - # } else {#' + '#: GridFrequency #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                filterable: false
            }, {
                title: "Reactive Power",
                field: "ReactivePower",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                filterable: false
            }, {
                title: "Alarm",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "Alarm Ext Stop Time",
                    field: "AlarmExtStopTime",
                    width: 90,
                    template: '# if (AlarmExtStopTime < -99998 ) { # - # } else {#' + '#: AlarmExtStopTime #' + '#}#',
                    attributes: {
                        class: "align-right"
                    },
                    filterable: false
                }, {
                    title: "Alarm Grid Down Time",
                    field: "AlarmGridDownTime",
                    template: '# if (AlarmGridDownTime < -99998 ) { # - # } else {#' + '#: AlarmGridDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Inter Line Down",
                    field: "AlarmInterLineDown",
                    template: '# if (AlarmInterLineDown < -99998 ) { # - # } else {#' + '#: AlarmInterLineDown #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Mach Down Time",
                    field: "AlarmMachDownTime",
                    template: '# if (AlarmMachDownTime < -99998 ) { # - # } else {#' + '#: AlarmMachDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm OK Time",
                    field: "AlarmOkTime",
                    template: '# if (AlarmOkTime < -99998 ) { # - # } else {#' + '#: AlarmOkTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: true,
                    headerAttributes: {
                        class: 'gridAlarmOkTime'
                    }
                }, {
                    title: "Alarm Unknown Time",
                    field: "AlarmUnknownTime",
                    template: '# if (AlarmUnknownTime < -99998 ) { # - # } else {#' + '#: AlarmUnknownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Alarm Weather Stop",
                    field: "AlarmWeatherStop",
                    template: '# if (AlarmWeatherStop < -99998 ) { # - # } else {#' + '#: AlarmWeatherStop #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }]
            }, {
                title: "Time",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "Ext Stop Time",
                    field: "ExternalStopTime",
                    template: '# if (ExternalStopTime < -99998 ) { # - # } else {#' + '#: ExternalStopTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Grid Down Time",
                    field: "GridDownTime",
                    template: '# if (GridDownTime < -99998 ) { # - # } else {#' + '#: GridDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Grid OK Secs",
                    field: "GridOkSecs",
                    template: '# if (GridOkSecs < -99998 ) { # - # } else {#' + '#: GridOkSecs #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Internal Line Down",
                    field: "InternalLineDown",
                    template: '# if (InternalLineDown < -99998 ) { # - # } else {#' + '#: InternalLineDown #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "Machine Down Time",
                    field: "MachineDownTime",
                    template: '# if (MachineDownTime < -99998 ) { # - # } else {#' + '#: MachineDownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "OK Secs",
                    field: "OkSecs",
                    template: '# if (OkSecs < -99998 ) { # - # } else {#' + '#: OkSecs #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }, {
                    title: "OK Time",
                    field: "OkTime",
                    template: '# if (OkTime < -99998 ) { # - # } else {#' + '#: OkTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: true
                }, {
                    title: "Unknown Time",
                    field: "UnknownTime",
                    template: '# if (UnknownTime < -99998 ) { # - # } else {#' + '#: UnknownTime #' + '#}#',
                    width: 90,
                    attributes: {
                        class: "align-right"
                    },
                    format: "{0:n2}",
                    filterable: false
                }]
            }, {
                title: "Total Time",
                field: "TotalTime",
                template: '# if (TotalTime < -99998 ) { # - # } else {#' + '#: TotalTime #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Generator RPM",
                field: "GeneratorRPM",
                template: '# if (GeneratorRPM < -99998 ) { # - # } else {#' + '#: GeneratorRPM #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacelle Yaw Position Untwist",
                field: "NacelleYawPositionUntwist",
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacelle Temperature",
                field: "NacelleTemperature",
                template: '# if (NacelleTemperature < -99998 ) { # - # } else {#' + '#: NacelleTemperature #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Adj Wind Speed",
                field: "AdjWindSpeed",
                template: '# if (AdjWindSpeed < -99998 ) { # - # } else {#' + '#: AdjWindSpeed #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Ambient Temperature",
                field: "AmbientTemperature",
                template: '# if (AmbientTemperature < -99998 ) { # - # } else {#' + '#: AmbientTemperature #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Avg Blade Angle",
                field: "AvgBladeAngle",
                template: '# if (AvgBladeAngle < -99998 ) { # - # } else {#' + '#: AvgBladeAngle #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Avg Wind Speed",
                field: "AvgWindSpeed",
                template: '# if (AvgWindSpeed < -99998 ) { # - # } else {#' + '#: AvgWindSpeed #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Units Generated",
                field: "UnitsGenerated",
                template: '# if (UnitsGenerated < -99998 ) { # - # } else {#' + '#: UnitsGenerated #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Estimated Power",
                field: "EstimatedPower",
                template: '# if (EstimatedPower < -99998 ) { # - # } else {#' + '#: EstimatedPower #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Nacel Direction",
                field: "NacelDirection",
                template: '# if (NacelDirection < -99998 ) { # - # } else {#' + '#: NacelDirection #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Power",
                field: "Power",
                template: '# if (Power < -99998 ) { # - # } else {#' + '#: Power #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: true
            }, {
                title: "Power Lost",
                field: "PowerLost",
                template: '# if (PowerLost < -99998 ) { # - # } else {#' + '#: PowerLost #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: true,
                headerAttributes: {
                    class: 'gridPowerLost'
                }
            }, {
                title: "Rotor RPM",
                field: "RotorRPM",
                template: '# if (RotorRPM < -99998 ) { # - # } else {#' + '#: RotorRPM #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }, {
                title: "Wind Direction",
                field: "WindDirection",
                template: '# if (WindDirection < -99998 ) { # - # } else {#' + '#: WindDirection #' + '#}#',
                width: 90,
                attributes: {
                    class: "align-right"
                },
                format: "{0:n2}",
                filterable: false
            }]
        });
        var grid = $('#dataGridAnomalies').data('kendoGrid');
        var columns = grid.columns;
        dbr.gridColumnsScadaAnomaly([]);
        $.each(columns, function(i, val) {
            $('#dataGridAnomalies').data("kendoGrid").showColumn(val.field);
            var result = {
                field: val.field,
                title: val.title,
                value: true
            }
            dbr.gridColumnsScadaAnomaly.push(result);
        });
    },
    InitGridAlarmAnomalies: function() {
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

        $('#dataGridAlarmAnomalies').html("");
        $('#dataGridAlarmAnomalies').kendoGrid({
            selectable: "multiple",
            dataSource: {
                serverPaging: true,
                serverSorting: true,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getalarmscadaanomalylist",
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
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totalturbinealarmAnomalies').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataalarmAnomalies').html(kendo.toString(res.data.Total, 'n0'));
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
            // toolbar: [
            //     "excel"
            // ],
            excel: {
                fileName: "Alarm Anomalies.xlsx",
                filterable: true,
                allPages: true
            },
            groupable: false,
            sortable: true,
            filterable: false,
            pageable: true,
            resizable: true,
            columns: [{
                    title: "Date",
                    field: "StartDate",
                    template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",
                    width: 80
                }, {
                    title: "Turbine",
                    field: "Turbine",
                    width: 90,
                    attributes: {
                        style: "text-align:center;"
                    }
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
                    width: 70,
                    attributes: {
                        style: "text-align:center;"
                    }
                }, {
                    title: "Alert Description",
                    field: "AlertDescription",
                    width: 200
                }, {
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
                    title: "Weather Stop",
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
                }, {
                    title: "Alarm Ok Time",
                    field: "IsAlarmOk",
                    width: 80,
                    sortable: false,
                    template: '# if (IsAlarmOk == true ) { # <span class="glyphicon glyphicon-ok"></span> # } else {# <span class="glyphicon glyphicon-remove"></span> #}#',
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
            pageable: false,
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
    InitGridAlarmOverlapping: function() {
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
                        if (!app.isFine(res)) {
                            return;
                        }
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
            pageable: true,
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
    },
    InitJMRDetail: function(e) {
        var turbine = [];
        if ($("#turbineMulti").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
            turbine = turbineval;
        } else {
            turbine = $("#turbineMulti").data("kendoMultiSelect").value();
        }

        var filters = [{
            field: "dateinfo.monthid",
            operator: "in",
            value: [e.data.DateInfo.MonthId]
        }, {
            field: "sections.turbine",
            operator: "in",
            value: turbine
        }, ];

        var param = {};

        $("<div/>").appendTo(e.detailCell).kendoGrid({
            selectable: "multiple",
            dataSource: {
                serverPaging: false,
                serverSorting: false,
                serverFiltering: true,
                filter: filters,
                transport: {
                    read: {
                        url: viewModel.appName + "databrowser/getjmrdetails",
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
                    model: {
                        fields: {
                            ContrGen: {
                                type: "number"
                            },

                            BoEExport: {
                                type: "number"
                            },
                            BoEImport: {
                                type: "number"
                            },
                            BoENet: {
                                type: "number"
                            },

                            BoLExport: {
                                type: "number"
                            },
                            BoLImport: {
                                type: "number"
                            },
                            BoLNet: {
                                type: "number"
                            },

                            BoE2Export: {
                                type: "number"
                            },
                            BoE2Import: {
                                type: "number"
                            },
                            BoE2Net: {
                                type: "number"
                            },
                        }
                    },
                    data: function(res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data
                    }
                },
            },
            scrollable: true,
            sortable: false,
            pageable: false,
            columns: [{
                title: "Description",
                field: "Description",
                width: 130,
                headerAttributes: {
                    style: "text-align: center"
                },
                sortable: false
            }, {
                title: "Turbine",
                field: "Turbine",
                width: 70,
                headerAttributes: {
                    style: "text-align: center"
                },
                attributes: {
                    style: "text-align: center"
                },
                sortable: false
            }, {
                title: "Company",
                field: "Company",
                width: 150,
                headerAttributes: {
                    style: "text-align: center"
                },
                sortable: false
            }, {
                title: "Controller Gen.",
                field: "ContrGen",
                format: "{0:n2}",
                width: 100,
                attributes: {
                    style: "text-align: center"
                },
                sortable: false,
                headerAttributes: {
                    style: "text-align: center"
                }
            }, {
                title: "Break of Energy",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "KWh Export",
                    field: "BoEExport",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    }
                }, {
                    title: "KWh Import",
                    field: "BoEImport",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    }
                }, {
                    title: "KWh Net",
                    field: "BoENet",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    }
                }, ]
            }, {
                title: "Break of Losses",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "KWh Export",
                    field: "BoLExport",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoLExport==0){#  #}else {# #: kendo.toString(BoLExport, 'n2') # #}#"
                }, {
                    title: "KWh Import",
                    field: "BoLImport",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoLImport==0){#  #}else {# #: kendo.toString(BoLImport, 'n2') # #}#"
                }, {
                    title: "KWh Net",
                    field: "BoLNet",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoLNet==0){#  #}else {# #: kendo.toString(BoLNet, 'n2') # #}#"
                }, ]
            }, {
                title: "Break of Energy",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [{
                    title: "KWh Export",
                    field: "BoE2Export",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoE2Export==0){#  #}else {# #: kendo.toString(BoE2Export, 'n2') # #}#"
                }, {
                    title: "KWh Import",
                    field: "BoE2Import",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoE2Import==0){#  #}else {# #: kendo.toString(BoE2Import, 'n2') # #}#"
                }, {
                    title: "KWh Net",
                    field: "BoE2Net",
                    format: "{0:n2}",
                    width: 100,
                    attributes: {
                        style: "text-align: center"
                    },
                    sortable: false,
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    template: "#if(BoE2Net==0){#  #}else {# #: kendo.toString(BoE2Net, 'n2') # #}#"
                }, ]
            }, ]
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

    Data.InitCustomGrid();

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
        Data.InitCustomGrid();
    }, 1000);
    Data.LoadAvailDate();
});