'use strict';

viewModel.DatabrowserScadaAnomaly = new Object();
var dbsa = viewModel.DatabrowserScadaAnomaly;

dbsa.InitGridAnomalies = function() {
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
                    app.isFine(res);
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
        pageable: {
            pageSize: 10,
            input:true, 
        },
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
}