'use strict';

viewModel.DatabrowserScada = new Object();
var dbs = viewModel.DatabrowserScada;

dbs.InitScadaGrid = function() {
    dbr.oemvis(true);

    // var turbine = [];
    // if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
    //     turbine = turbineval;
    // } else {
    //     turbine = $("#turbineList").data("kendoMultiSelect").value();
    // }
    var misc = {
        "tipe": "scadaoem",
        "needtotalturbine": true,
        "period": fa.period,
    }
    var param = {"misc": misc};

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

    $('#scadaGrid').html("");
    $('#scadaGrid').kendoGrid({
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
                    dbr.oemvis(true);
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function(res) {
                    app.loading(false);
                    dbr.oemvis(false);
                    app.isFine(res);
                    
                    return res.data.Data
                },
                total: function(res) {
                    if (!app.isFine(res)) {
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
        pageable: {
            pageSize: 10,
            input:true, 
        },
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
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Active Power",
            field: "AI_intern_ActivPower",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern I1",
            field: "AI_intern_I1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern I2",
            field: "AI_intern_I2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern I3",
            field: "AI_intern_I3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Nacelle Drill",
            field: "AI_intern_NacelleDrill",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Nacelle Pos",
            field: "AI_intern_NacellePos",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Akku V1",
            field: "AI_intern_PitchAkku_V1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Akku V2",
            field: "AI_intern_PitchAkku_V2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Akku V3",
            field: "AI_intern_PitchAkku_V3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Angle1",
            field: "AI_intern_PitchAngle1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Angle2",
            field: "AI_intern_PitchAngle2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Angle3",
            field: "AI_intern_PitchAngle3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Conv Current1",
            field: "AI_intern_PitchConv_Current1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Conv Current2",
            field: "AI_intern_PitchConv_Current2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Conv Current3",
            field: "AI_intern_PitchConv_Current3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Angle Sp Diff1",
            field: "AI_intern_PitchAngleSP_Diff1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Angle Sp Diff2",
            field: "AI_intern_PitchAngleSP_Diff2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchangle Sp Diff3",
            field: "AI_intern_PitchAngleSP_Diff3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Reactive Power",
            field: "AI_intern_ReactivPower",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Rpm Diff",
            field: "AI_intern_RpmDiff",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern U1",
            field: "AI_intern_U1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern U2",
            field: "AI_intern_U2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern U3",
            field: "AI_intern_U3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Wind Direction",
            field: "AI_intern_WindDirection",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Wind Speed",
            field: "AI_intern_WindSpeed",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Wind Speed Dif",
            field: "AI_Intern_WindSpeedDif",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Speed Rot Fr",
            field: "AI_speed_RotFR",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Wind Speed1",
            field: "AI_WindSpeed1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Wind Speed2",
            field: "AI_WindSpeed2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Wind Vane1",
            field: "AI_WindVane1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Wind Vane2",
            field: "AI_WindVane2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Current Asym",
            field: "AI_internCurrentAsym",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Gear Box Ims Nde",
            field: "Temp_GearBox_IMS_NDE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Wind Vane Diff",
            field: "AI_intern_WindVaneDiff",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "C Intern Speed Generator",
            field: "C_intern_SpeedGenerator",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "C Intern Speed Rotor",
            field: "C_intern_SpeedRotor",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Speed Rpm Diff Fr1 Rot Cnt",
            field: "AI_intern_Speed_RPMDiff_FR1_RotCNT",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Frequency Grid",
            field: "AI_intern_Frequency_Grid",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Gear Box Hss Nde",
            field: "Temp_GearBox_HSS_NDE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Dr Tr Vib Value",
            field: "AI_DrTrVibValue",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Last Error Conv1",
            field: "AI_intern_InLastErrorConv1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Last Error Conv2",
            field: "AI_intern_InLastErrorConv2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Last Error Conv3",
            field: "AI_intern_InLastErrorConv3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Conv1",
            field: "AI_intern_TempConv1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Conv2",
            field: "AI_intern_TempConv2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Conv3",
            field: "AI_intern_TempConv3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch Speed1",
            field: "AI_intern_PitchSpeed1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Brake 1",
            field: "Temp_YawBrake_1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Brake 2",
            field: "Temp_YawBrake_2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp G1L1",
            field: "Temp_G1L1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp G1L2",
            field: "Temp_G1L2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp G1L3",
            field: "Temp_G1L3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Brake 3",
            field: "Temp_YawBrake_3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Hydr System Pressure",
            field: "AI_HydrSystemPressure",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Bottom Control Section Low",
            field: "Temp_BottomControlSection_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Gearbox Hss De",
            field: "Temp_GearBox_HSS_DE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Gear Oil Sump",
            field: "Temp_GearOilSump",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Generator Bearing De",
            field: "Temp_GeneratorBearing_DE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Generator Bearing Nde",
            field: "Temp_GeneratorBearing_NDE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Main Bearing",
            field: "Temp_MainBearing",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Gear Box Ims De",
            field: "Temp_GearBox_IMS_DE",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Nacelle",
            field: "Temp_Nacelle",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Outdoor",
            field: "Temp_Outdoor",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Tower Vib Value Axial",
            field: "AI_TowerVibValueAxial",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Diff Gen Speed Sp To Act",
            field: "AI_intern_DiffGenSpeedSPToAct",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Brake 4",
            field: "Temp_YawBrake_4",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Speed Generator Proximity",
            field: "AI_intern_SpeedGenerator_Proximity",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Speed Diff Encoder Proximity",
            field: "AI_intern_SpeedDiff_Encoder_Proximity",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Gear Oil Pressure",
            field: "AI_GearOilPressure",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Cabinet Top Box Low",
            field: "Temp_CabinetTopBox_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Cabinet Top Box",
            field: "Temp_CabinetTopBox",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Bottom Control Section",
            field: "Temp_BottomControlSection",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Bottom Power Section",
            field: "Temp_BottomPowerSection",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Bottom Power Section Low",
            field: "Temp_BottomPowerSection_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch1 Status High",
            field: "AI_intern_Pitch1_Status_High",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch2 Status High",
            field: "AI_intern_Pitch2_Status_High",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch3 Status High",
            field: "AI_intern_Pitch3_Status_High",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Position1 Ch2",
            field: "AI_intern_InPosition1_ch2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Position2 Ch2",
            field: "AI_intern_InPosition2_ch2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern In Position3 Ch2",
            field: "AI_intern_InPosition3_ch2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Brake Blade1",
            field: "AI_intern_Temp_Brake_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Brake Blade2",
            field: "AI_intern_Temp_Brake_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Brake Blade3",
            field: "AI_intern_Temp_Brake_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch Motor Blade1",
            field: "AI_intern_Temp_PitchMotor_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch Motor Blade2",
            field: "AI_intern_Temp_PitchMotor_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch Motor Blade3",
            field: "AI_intern_Temp_PitchMotor_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Additional1",
            field: "AI_intern_Temp_Hub_Additional1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Additional2",
            field: "AI_intern_Temp_Hub_Additional2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Additional3",
            field: "AI_intern_Temp_Hub_Additional3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch1 Status Low",
            field: "AI_intern_Pitch1_Status_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch2 Status Low",
            field: "AI_intern_Pitch2_Status_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitch3 Status Low",
            field: "AI_intern_Pitch3_Status_Low",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Voltage Blade1 Center",
            field: "AI_intern_Battery_VoltageBlade1_center",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Voltage Blade2 Center",
            field: "AI_intern_Battery_VoltageBlade2_center",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Voltage Blade3 Center",
            field: "AI_intern_Battery_VoltageBlade3_center",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Charging Cur Blade1",
            field: "AI_intern_Battery_ChargingCur_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Charging Cur Blade2",
            field: "AI_intern_Battery_ChargingCur_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Charging Cur Blade3",
            field: "AI_intern_Battery_ChargingCur_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Discharging Cur Blade1",
            field: "AI_intern_Battery_DischargingCur_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Discharging Cur Blade2",
            field: "AI_intern_Battery_DischargingCur_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Battery Discharging Cur Blade3",
            field: "AI_intern_Battery_DischargingCur_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Voltage Blade1",
            field: "AI_intern_PitchMotor_BrakeVoltage_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Voltage Blade2",
            field: "AI_intern_PitchMotor_BrakeVoltage_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Voltage Blade3",
            field: "AI_intern_PitchMotor_BrakeVoltage_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Current Blade1",
            field: "AI_intern_PitchMotor_BrakeCurrent_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Current Blade2",
            field: "AI_intern_PitchMotor_BrakeCurrent_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Pitchmotor Brake Current Blade3",
            field: "AI_intern_PitchMotor_BrakeCurrent_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Box Blade1",
            field: "AI_intern_Temp_HubBox_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Box Blade2",
            field: "AI_intern_Temp_HubBox_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Hub Box Blade3",
            field: "AI_intern_Temp_HubBox_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch1 Heat Sink",
            field: "AI_intern_Temp_Pitch1_HeatSink",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch2 Heat Sink",
            field: "AI_intern_Temp_Pitch2_HeatSink",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Pitch3 Heat Sink",
            field: "AI_intern_Temp_Pitch3_HeatSink",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Error Stack Blade1",
            field: "AI_intern_ErrorStackBlade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Error Stack Blade2",
            field: "AI_intern_ErrorStackBlade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Error Stack Blade3",
            field: "AI_intern_ErrorStackBlade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Battery Box Blade1",
            field: "AI_intern_Temp_BatteryBox_Blade1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Battery Box Blade2",
            field: "AI_intern_Temp_BatteryBox_Blade2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Temp Battery Box Blade3",
            field: "AI_intern_Temp_BatteryBox_Blade3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dc Linkvoltage1",
            field: "AI_intern_DC_LinkVoltage1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dc Linkvoltage2",
            field: "AI_intern_DC_LinkVoltage2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dc Linkvoltage3",
            field: "AI_intern_DC_LinkVoltage3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Motor1",
            field: "Temp_Yaw_Motor1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Motor2",
            field: "Temp_Yaw_Motor2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Motor3",
            field: "Temp_Yaw_Motor3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Temp Yaw Motor4",
            field: "Temp_Yaw_Motor4",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ao Dfig Power Setpiont",
            field: "AO_DFIG_Power_Setpiont",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ao Dfig Q Setpoint",
            field: "AO_DFIG_Q_Setpoint",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Dfig Torque Actual",
            field: "AI_DFIG_Torque_actual",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Dfig Speed Generator Encoder",
            field: "AI_DFIG_SpeedGenerator_Encoder",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Dc Link Voltage Actual",
            field: "AI_intern_DFIG_DC_Link_Voltage_actual",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Msc Current",
            field: "AI_intern_DFIG_MSC_current",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Main Voltage",
            field: "AI_intern_DFIG_Main_voltage",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Main Current",
            field: "AI_intern_DFIG_Main_current",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Active Power Actual",
            field: "AI_intern_DFIG_active_power_actual",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Reactive Power Actual",
            field: "AI_intern_DFIG_reactive_power_actual",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Active Power Actual Lsc",
            field: "AI_intern_DFIG_active_power_actual_LSC",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Lsc Current",
            field: "AI_intern_DFIG_LSC_current",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Dfig Data Log Number",
            field: "AI_intern_DFIG_Data_log_number",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Damper Osc Magnitude",
            field: "AI_intern_Damper_OscMagnitude",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Damper Passband Full Load",
            field: "AI_intern_Damper_PassbandFullLoad",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Yaw Brake Temp Rise1",
            field: "AI_YawBrake_TempRise1",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Yaw Brake Temp Rise2",
            field: "AI_YawBrake_TempRise2",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Yaw Brake Temp Rise3",
            field: "AI_YawBrake_TempRise3",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Yaw Brake Temp Rise4",
            field: "AI_YawBrake_TempRise4",
            width: 90,
            attributes: {
                class: "align-center"
            },
            format: "{0:n2}",
            filterable: false
        }, {
            title: "Ai Intern Nacelle Drill At North Pos Sensor",
            field: "AI_intern_NacelleDrill_at_NorthPosSensor",
            width: 90,
            attributes: {
                class: "align-center"
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
    $('#scadaGrid').data('kendoGrid').refresh();
}