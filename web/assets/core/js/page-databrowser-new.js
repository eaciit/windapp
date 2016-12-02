'use strict';

viewModel.Databrowser = new Object();
var dbr = viewModel.Databrowser;

dbr.turbineList = ko.observableArray([]);
dbr.modelList = ko.observableArray([
    { "value": 1, "text": "Regen" },
    { "value": 2, "text": "Suzlon" },
]);
dbr.projectList = ko.observableArray([
    { "value": 1, "text": "WindFarm-01" },
    { "value": 2, "text": "WindFarm-02" },
]);

dbr.gridColumnsScada = ko.observableArray([]);
dbr.gridColumnsScadaException = ko.observableArray([]);
dbr.gridColumnsScadaAnomaly = ko.observableArray([]);
dbr.filterJMR = ko.observableArray([]);
dbr.selectedColumn = ko.observableArray([]);
dbr.unselectedColumn = ko.observableArray([]);
dbr.ColumnList = ko.observableArray([]);
dbr.ColList = ko.observableArray([]);
dbr.defaultSelectedColumn = ko.observableArray([
	 { "_id": "timestamp", "label": "Time Stamp", "source": "ScadaDataOEM"},
	 { "_id": "turbine", "label": "Turbine", "source": "ScadaDataOEM"},
	 { "_id": "ai_intern_r_pidangleout", "label": "Ai Intern R Pid Angle Out", "source": "ScadaDataOEM"},
	 { "_id": "ai_intern_activpower", "label": "Ai Intern Active Power", "source": "ScadaDataOEM"},
	 { "_id": "ai_intern_i1", "label": "Ai Intern I1", "source": "ScadaDataOEM"},
	 { "_id": "ai_intern_i2", "label": "Ai Intern I2", "source": "ScadaDataOEM"},
]);

var turbineval = [];

dbr.populateTurbine = function () {
    app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0) {
            res.data = [];
            dbr.turbineList([{ value: "", text: "" }]);
        } else {
            var datavalue = [];
            if (res.data.length > 0) {
                var allturbine = {}
                $.each(res.data, function (key, val) {
                    turbineval.push(val);
                });
                allturbine.value = "All Turbine";
                allturbine.text = "All Turbines";
                datavalue.push(allturbine);
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    datavalue.push(data);
                });
            }
            dbr.turbineList(datavalue);
        }
        setTimeout(function () {
            $('#turbineMulti').data('kendoMultiSelect').value(["All Turbine"])
        }, 300);
    });
};

dbr.checkTurbine = function () {
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

dbr.ShowHideColumnScada = function (gridID, field, id, index) {
    if ($('#' + id).is(":checked")) {
        $('#' + gridID).data("kendoGrid").showColumn(index);
    } else {
        $('#' + gridID).data("kendoGrid").hideColumn(index);
    }
}

var Data = {
    LoadData: function () {
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();

        if ($("#turbineMulti").data("kendoMultiSelect").value() == "") {
            $('#turbineMulti').data('kendoMultiSelect').value(["All Turbine"])
        }

        if ((new Date(dateStart).getTime() > new Date(dateEnd).getTime())) {
            toolkit.showError("Invalid Date Range Selection");
            return;
        } else {
            app.loading(true);
            this.InitScadaGrid();
            this.InitDEgrid();
            this.InitCustomGrid();

        }

        this.LoadAvailDate();
        this.LoadAvailDateDE();
        this.LoadAvailDateCustom();
    },
    LoadAvailDate: function () {
        app.ajaxPost(viewModel.appName + "/databrowsernew/getscadadataoemavaildate", {}, function (res) {
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
    },
    LoadAvailDateDE: function () {
        app.ajaxPost(viewModel.appName + "/databrowsernew/getdowntimeeventvaildate", {}, function (res) {
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
    },
    LoadAvailDateCustom: function () {
        app.ajaxPost(viewModel.appName + "/databrowsernew/getcustomavaildate", {}, function (res) {
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
    },
    RefreshGrid: function () {
        setTimeout(function () {
            $('#scadaGrid').data('kendoGrid').refresh();
            $('#customGrid').data('kendoGrid').refresh();
            // $('#customGrid').data("kendoGrid").hideColumn(0);
            $('#DEgrid').data('kendoGrid').refresh();
        }, 200);
    },
    InitScadaGrid: function () {
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
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
                        url: viewModel.appName + "databrowsernew/getscadalist",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function (options) {
                        return JSON.stringify(options);
                    }
                },
                pageSize: 10,
                schema: {
                    data: function (res) {
                        app.loading(false);
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totalturbine').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldata').html(kendo.toString(res.data.Total, 'n0'));
                        $('#totalprodoem').html(kendo.toString(res.data.TotalProduction / 1000, 'n0') + ' MWh');
                        $('#avgwindspeedoem').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
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
                { title: "Ai Intern R Pid Angle Out", field: "AI_intern_R_PidAngleOut", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Active Power", field: "AI_intern_ActivPower", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern I1", field: "AI_intern_I1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern I2", field: "AI_intern_I2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern I3", field: "AI_intern_I3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Nacelle Drill", field: "AI_intern_NacelleDrill", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Nacelle Pos", field: "AI_intern_NacellePos", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Akku V1", field: "AI_intern_PitchAkku_V1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Akku V2", field: "AI_intern_PitchAkku_V2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Akku V3", field: "AI_intern_PitchAkku_V3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Angle1", field: "AI_intern_PitchAngle1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Angle2", field: "AI_intern_PitchAngle2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Angle3", field: "AI_intern_PitchAngle3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Conv Current1", field: "AI_intern_PitchConv_Current1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Conv Current2", field: "AI_intern_PitchConv_Current2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Conv Current3", field: "AI_intern_PitchConv_Current3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Angle Sp Diff1", field: "AI_intern_PitchAngleSP_Diff1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Angle Sp Diff2", field: "AI_intern_PitchAngleSP_Diff2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchangle Sp Diff3", field: "AI_intern_PitchAngleSP_Diff3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Reactive Power", field: "AI_intern_ReactivPower", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Rpm Diff", field: "AI_intern_RpmDiff", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern U1", field: "AI_intern_U1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern U2", field: "AI_intern_U2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern U3", field: "AI_intern_U3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Wind Direction", field: "AI_intern_WindDirection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Wind Speed", field: "AI_intern_WindSpeed", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Wind Speed Dif", field: "AI_Intern_WindSpeedDif", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Speed Rot Fr", field: "AI_speed_RotFR", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Wind Speed1", field: "AI_WindSpeed1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Wind Speed2", field: "AI_WindSpeed2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Wind Vane1", field: "AI_WindVane1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Wind Vane2", field: "AI_WindVane2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Current Asym", field: "AI_internCurrentAsym", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Gear Box Ims Nde", field: "Temp_GearBox_IMS_NDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Wind Vane Diff", field: "AI_intern_WindVaneDiff", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "C Intern Speed Generator", field: "C_intern_SpeedGenerator", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "C Intern Speed Rotor", field: "C_intern_SpeedRotor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Speed Rpm Diff Fr1 Rot Cnt", field: "AI_intern_Speed_RPMDiff_FR1_RotCNT", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Frequency Grid", field: "AI_intern_Frequency_Grid", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Gear Box Hss Nde", field: "Temp_GearBox_HSS_NDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Dr Tr Vib Value", field: "AI_DrTrVibValue", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Last Error Conv1", field: "AI_intern_InLastErrorConv1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Last Error Conv2", field: "AI_intern_InLastErrorConv2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Last Error Conv3", field: "AI_intern_InLastErrorConv3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Conv1", field: "AI_intern_TempConv1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Conv2", field: "AI_intern_TempConv2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Conv3", field: "AI_intern_TempConv3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch Speed1", field: "AI_intern_PitchSpeed1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Brake 1", field: "Temp_YawBrake_1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Brake 2", field: "Temp_YawBrake_2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp G1L1", field: "Temp_G1L1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp G1L2", field: "Temp_G1L2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp G1L3", field: "Temp_G1L3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Brake 3", field: "Temp_YawBrake_3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Hydr System Pressure", field: "AI_HydrSystemPressure", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Bottom Control Section Low", field: "Temp_BottomControlSection_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Gearbox Hss De", field: "Temp_GearBox_HSS_DE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Gear Oil Sump", field: "Temp_GearOilSump", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Generator Bearing De", field: "Temp_GeneratorBearing_DE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Generator Bearing Nde", field: "Temp_GeneratorBearing_NDE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Main Bearing", field: "Temp_MainBearing", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Gear Box Ims De", field: "Temp_GearBox_IMS_DE", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Nacelle", field: "Temp_Nacelle", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Outdoor", field: "Temp_Outdoor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Tower Vib Value Axial", field: "AI_TowerVibValueAxial", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Diff Gen Speed Sp To Act", field: "AI_intern_DiffGenSpeedSPToAct", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Brake 4", field: "Temp_YawBrake_4", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Speed Generator Proximity", field: "AI_intern_SpeedGenerator_Proximity", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Speed Diff Encoder Proximity", field: "AI_intern_SpeedDiff_Encoder_Proximity", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Gear Oil Pressure", field: "AI_GearOilPressure", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Cabinet Top Box Low", field: "Temp_CabinetTopBox_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Cabinet Top Box", field: "Temp_CabinetTopBox", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Bottom Control Section", field: "Temp_BottomControlSection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Bottom Power Section", field: "Temp_BottomPowerSection", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Bottom Power Section Low", field: "Temp_BottomPowerSection_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch1 Status High", field: "AI_intern_Pitch1_Status_High", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch2 Status High", field: "AI_intern_Pitch2_Status_High", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch3 Status High", field: "AI_intern_Pitch3_Status_High", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Position1 Ch2", field: "AI_intern_InPosition1_ch2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Position2 Ch2", field: "AI_intern_InPosition2_ch2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern In Position3 Ch2", field: "AI_intern_InPosition3_ch2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Brake Blade1", field: "AI_intern_Temp_Brake_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Brake Blade2", field: "AI_intern_Temp_Brake_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Brake Blade3", field: "AI_intern_Temp_Brake_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch Motor Blade1", field: "AI_intern_Temp_PitchMotor_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch Motor Blade2", field: "AI_intern_Temp_PitchMotor_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch Motor Blade3", field: "AI_intern_Temp_PitchMotor_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Additional1", field: "AI_intern_Temp_Hub_Additional1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Additional2", field: "AI_intern_Temp_Hub_Additional2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Additional3", field: "AI_intern_Temp_Hub_Additional3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch1 Status Low", field: "AI_intern_Pitch1_Status_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch2 Status Low", field: "AI_intern_Pitch2_Status_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitch3 Status Low", field: "AI_intern_Pitch3_Status_Low", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Voltage Blade1 Center", field: "AI_intern_Battery_VoltageBlade1_center", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Voltage Blade2 Center", field: "AI_intern_Battery_VoltageBlade2_center", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Voltage Blade3 Center", field: "AI_intern_Battery_VoltageBlade3_center", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Charging Cur Blade1", field: "AI_intern_Battery_ChargingCur_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Charging Cur Blade2", field: "AI_intern_Battery_ChargingCur_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Charging Cur Blade3", field: "AI_intern_Battery_ChargingCur_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Discharging Cur Blade1", field: "AI_intern_Battery_DischargingCur_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Discharging Cur Blade2", field: "AI_intern_Battery_DischargingCur_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Battery Discharging Cur Blade3", field: "AI_intern_Battery_DischargingCur_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Voltage Blade1", field: "AI_intern_PitchMotor_BrakeVoltage_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Voltage Blade2", field: "AI_intern_PitchMotor_BrakeVoltage_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Voltage Blade3", field: "AI_intern_PitchMotor_BrakeVoltage_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Current Blade1", field: "AI_intern_PitchMotor_BrakeCurrent_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Current Blade2", field: "AI_intern_PitchMotor_BrakeCurrent_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Pitchmotor Brake Current Blade3", field: "AI_intern_PitchMotor_BrakeCurrent_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Box Blade1", field: "AI_intern_Temp_HubBox_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Box Blade2", field: "AI_intern_Temp_HubBox_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Hub Box Blade3", field: "AI_intern_Temp_HubBox_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch1 Heat Sink", field: "AI_intern_Temp_Pitch1_HeatSink", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch2 Heat Sink", field: "AI_intern_Temp_Pitch2_HeatSink", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Pitch3 Heat Sink", field: "AI_intern_Temp_Pitch3_HeatSink", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Error Stack Blade1", field: "AI_intern_ErrorStackBlade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Error Stack Blade2", field: "AI_intern_ErrorStackBlade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Error Stack Blade3", field: "AI_intern_ErrorStackBlade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Battery Box Blade1", field: "AI_intern_Temp_BatteryBox_Blade1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Battery Box Blade2", field: "AI_intern_Temp_BatteryBox_Blade2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Temp Battery Box Blade3", field: "AI_intern_Temp_BatteryBox_Blade3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dc Linkvoltage1", field: "AI_intern_DC_LinkVoltage1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dc Linkvoltage2", field: "AI_intern_DC_LinkVoltage2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dc Linkvoltage3", field: "AI_intern_DC_LinkVoltage3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Motor1", field: "Temp_Yaw_Motor1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Motor2", field: "Temp_Yaw_Motor2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Motor3", field: "Temp_Yaw_Motor3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Temp Yaw Motor4", field: "Temp_Yaw_Motor4", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ao Dfig Power Setpiont", field: "AO_DFIG_Power_Setpiont", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ao Dfig Q Setpoint", field: "AO_DFIG_Q_Setpoint", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Dfig Torque Actual", field: "AI_DFIG_Torque_actual", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Dfig Speed Generator Encoder", field: "AI_DFIG_SpeedGenerator_Encoder", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Dc Link Voltage Actual", field: "AI_intern_DFIG_DC_Link_Voltage_actual", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Msc Current", field: "AI_intern_DFIG_MSC_current", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Main Voltage", field: "AI_intern_DFIG_Main_voltage", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Main Current", field: "AI_intern_DFIG_Main_current", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Active Power Actual", field: "AI_intern_DFIG_active_power_actual", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Reactive Power Actual", field: "AI_intern_DFIG_reactive_power_actual", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Active Power Actual Lsc", field: "AI_intern_DFIG_active_power_actual_LSC", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Lsc Current", field: "AI_intern_DFIG_LSC_current", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Dfig Data Log Number", field: "AI_intern_DFIG_Data_log_number", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Damper Osc Magnitude", field: "AI_intern_Damper_OscMagnitude", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Damper Passband Full Load", field: "AI_intern_Damper_PassbandFullLoad", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Yaw Brake Temp Rise1", field: "AI_YawBrake_TempRise1", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Yaw Brake Temp Rise2", field: "AI_YawBrake_TempRise2", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Yaw Brake Temp Rise3", field: "AI_YawBrake_TempRise3", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Yaw Brake Temp Rise4", field: "AI_YawBrake_TempRise4", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
                { title: "Ai Intern Nacelle Drill At North Pos Sensor", field: "AI_intern_NacelleDrill_at_NorthPosSensor", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },
            ]
        });

        var grid = $('#scadaGrid').data('kendoGrid');
        var columns = grid.columns;
        dbr.gridColumnsScada([]);
        $.each(columns, function (i, val) {
            $('#scadaGrid').data("kendoGrid").showColumn(val.field);
            var result = {
                field: val.field,
                title: val.title,
                value: true
            }

            dbr.gridColumnsScada.push(result);

                
        });
	},
	InitCustomGrid: function(){
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
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
        if(dbr.selectedColumn().length == 0){
            gColumns = dbr.defaultSelectedColumn();
        }

        $.each(gColumns, function(i, val){
         var col = {
             field : val._id, 
             title : val.label,
             type: val._id == "turbine" ? "string" : "number",
             width: 120,
             headerAttributes: {style:"text-align:center"}
         };

        if(val._id == "timestamp"){
            col = {
                field: val._id,
                title: val.label,
                type: "date",
                width:140,
                template: "#= kendo.toString(moment.utc(timestamp).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #", 
                value: true
            }
        }
         columns.push(col);
        });

        $('#customGrid').html("");
        $('#customGrid').kendoGrid({
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
                        url: viewModel.appName + "databrowsernew/getcustomlist",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function (options) {
                        return JSON.stringify(options);
                    }
                },
                pageSize: 10,
                schema: {
                    data: function (res) {
                        app.loading(false);
                        if (!app.isFine(res)) {
                            return;
                        }

                        return res.data.Data
                    },
                    total: function (res) {

                        if (!app.isFine(res)) {
                            return;
                        }
                        $('#totalturbineCustom').html(kendo.toString(res.data.TotalTurbine, 'n0'));
                        $('#totaldataCustom').html(kendo.toString(res.data.Total, 'n0'));
                        $('#totalprodCustom').html(kendo.toString(res.data.TotalProduction / 1000, 'n0') + ' MWh');
                        $('#avgwindspeedCustom').html(kendo.toString(res.data.AvgWindSpeed, 'n0') + ' m/s');
                        return res.data.Total;
                    },
                },
                sort: [
                    { field: 'TimeStamp', dir: 'asc' },
                    { field: 'Turbine', dir: 'asc' }
                ],
            },
	         toolbar: [
                "excel",
                {
                    text: "Show Hide Columns",
                    name: "showHideColumn",
                    imageClass: "fa fa-eye-slash ",
                }
            ],
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
            columns : columns,
        });
        
        var grid = $('#customGrid').data('kendoGrid');
        var columns = grid.columns;

        $.each(columns, function (i, val) {
            $('#customGrid').data("kendoGrid").hideColumn(val.field);
        });

        if(dbr.selectedColumn() == ""){
        	$.each(dbr.defaultSelectedColumn(), function (idx, data) {
        		$('#customGrid').data("kendoGrid").showColumn(data._id);
	        });
        }else{
        	$.each(dbr.selectedColumn(), function (idx, data) {
	            $('#customGrid').data("kendoGrid").showColumn(data._id);
	        });
        }
        $('.k-grid-showHideColumn').on("click", function(){
            Data.InitColumnList();
            $("#modalShowHide").modal();
            return false;
        });
    },
    InitDEgrid: function() {
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
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
                    url: viewModel.appName + "databrowsernew/getdowntimeeventlist",
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
            sort: [
            {
                field: 'timestart',
                dir: 'asc'
            }, 
            {
                field: 'turbine',
                dir: 'asc'
            }],
        },
        selectable: "multiple",
        groupable: false,
        sortable: true,
        pageable: true,
        columns: [
        {
            title: "Time Start",
            field: "TimeStart",
            template: "#= kendo.toString(moment.utc(TimeStart).format('DD-MMM-YYYY HH:mm:ss'), 'dd-MMM-yyyy HH:mm:ss') #",
            width: 130,
            filterable: false
        }, 
        {
            title: "Turbine",
            field: "Turbine",
            attributes: {
                class: "align-center"
            },
            width: 90,
            filterable: false
        },
        { title: "Alarm Description", field: "AlarmDescription", width: 100, filterable: false },
        { title: "Duration (Second)", field: "Duration", width: 90, attributes: { class: "align-right" }, format: "{0:n2}", filterable: false },

        ]
    });
	},
    InitDefault: function () {
    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-7, 0, 0, 0, 0));
    var lastEndDate = new Date(app.toUTC(maxDateData));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);

        setTimeout(function () {
            Data.LoadData();
        }, 500);
    },
    InitColumnList : function(){
    	$("#columnList").kendoGrid({
            theme: "flat",
            dataSource: {
                data: (dbr.selectedColumn() == "" ? dbr.ColumnList() : dbr.unselectedColumn()),
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [
                { field: "label", title: "Columns List", headerAttributes: {style: "text-align: center" }},
            ],
            change: function (arg) {
		        var selected = $.map(this.select(), function (item) {
		            return $(item).find('td').first().text();
		        });
		        var grid1 = $('#columnList').data('kendoGrid');
				var grid2 = $('#selectedList').data('kendoGrid');
			    dbr.gridMoveTo(grid1, grid2, false);
		    },
        });

     	setTimeout(function(){
	    	$('#columnList').data('kendoGrid').refresh();
	    	$('#selectedList').data('kendoGrid').refresh();
    	},300);

        $("#selectedList").kendoGrid({
            theme: "flat",
            dataSource: {
                data: dbr.selectedColumn() == "" ? dbr.defaultSelectedColumn() : dbr.selectedColumn(),
            },
            height: 300,
            scrollable: true,
            sortable: true,
            selectable: "multiple",
            columns: [
                { field: "label", title: "Selected Columns", headerAttributes: {style: "text-align: center" }},
            ],
            change: function (arg) {
		        var selected = $.map(this.select(), function (item) {
		            return $(item).find('td').first().text();
		        });
		        var grid1 = $('#columnList').data('kendoGrid');
				var grid2 = $('#selectedList').data('kendoGrid');
			    dbr.gridMoveTo(grid2, grid1, false);
		    },
        });
    }
};

dbr.selectRow = function(){
	var grid1 = $('#columnList').data('kendoGrid');
	var grid2 = $('#selectedList').data('kendoGrid');
    dbr.gridMoveTo(grid2, grid1, true);
}

dbr.unselectRow = function(){
	var grid1 = $('#columnList').data('kendoGrid');
	var grid2 = $('#selectedList').data('kendoGrid');
    dbr.gridMoveTo(grid1, grid2,true);
}

dbr.gridMoveTo = function (from, to ,all) {
	if(all == true){
		from.select(from.tbody.find(">tr"));
	}
	var selected = from.select();

    if (selected.length > 0) {
        var items = [];
        $.each(selected, function (idx, elem) {
            items.push(from.dataItem(elem));
        });
        var fromDS = from.dataSource;
        var toDS = to.dataSource;
        $.each(items, function (idx, elem) {
            toDS.add({ _id: elem._id, label: elem.label, source: elem.source});
            fromDS.remove(elem);
        });
        toDS.sync();
        fromDS.sync();
    }
}

dbr.showColumn = function(){
	dbr.selectedColumn([]);
	dbr.unselectedColumn([]);
	var grid = $('#selectedList').data('kendoGrid');
	var dataSources = grid.dataSource.data();
	var selectedList = [];
	var columnList = [];

	$.each(dataSources, function(i, val){
		selectedList.push(val.id);
		dbr.selectedColumn.push({
			_id: val._id,
            label : val.label,
            source: val.source
		});
	});

	$.each($('#columnList').data('kendoGrid').dataSource.data(), function(i, val){
		dbr.unselectedColumn.push({
			_id: val._id,
            label : val.label,
            source: val.source
		});
	});

	$.each(dbr.ColumnList(), function(idx, data){
		columnList.push(data.id);
	})

	// $.grep(columnList, function(el) {
 //        if ($.inArray(el, selectedList) == -1){
 //        	$('#customGrid').data("kendoGrid").hideColumn(el);
 //        }else{
 //        	$('#customGrid').data("kendoGrid").showColumn(el);
 //        }
	// });
	Data.InitCustomGrid();

	$('#modalShowHide').modal("hide");
}


vm.currentMenu('New Data');
vm.currentTitle('New Data');
vm.breadcrumb([{ title: 'Databrowser', href: '#' }, { title: 'New Databrowser', href: viewModel.appName + 'page/databrowsernew' }]);

$(document).ready(function () {
    app.loading(true);
    dbr.populateTurbine();
     $('.k-grid-showHideColumn').on("click", function(){
        $("#modalShowHide").modal();

        $("#myModal").on('shown.bs.modal', function () {
            Data.InitColumnList();
        });
        return false;
     });
    $('#btnRefresh').on('click', function () {
        Data.LoadData();
    });

    setTimeout(function () {
        Data.InitDefault();
        Data.InitCustomGrid();
    }, 1000);


});