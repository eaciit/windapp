'use strict';

viewModel.DatabrowserJMR = new Object();
var dbj = viewModel.DatabrowserJMR;

dbj.InitGridJMR = function() {
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
}