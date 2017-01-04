'use strict';

viewModel.DatabrowserJMR = new Object();
var dbj = viewModel.DatabrowserJMR;

dbj.InitGridJMR = function() {
    dbr.jmrvis(true);
    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
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
                    app.isFine(res);
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
        pageable: {
            pageSize: 10,
            input:true, 
        },
        detailInit: dbj.InitJMRDetail,
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
}

dbj.InitJMRDetail = function(e) {
    var turbine = [];
    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        turbine = turbineval;
    } else {
        turbine = $("#turbineList").data("kendoMultiSelect").value();
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
        pageable: {
            pageSize: 10,
            input:true, 
        },
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
}