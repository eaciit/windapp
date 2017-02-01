'use strict';


viewModel.PCComparison = new Object();
var pc = viewModel.PCComparison;


vm.currentMenu('Comparison');
vm.currentTitle('Comparison');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Comparison',
    href: viewModel.appName + 'page/analyticpccomparison'
}]);


pc.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Last 24 hours" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);

pc.turbineList = ko.observableArray([]);
pc.projectList = ko.observableArray([]);
pc.dateStart = ko.observable();
pc.dateEnd = ko.observable();
pc.turbine = ko.observableArray([]);
pc.project = ko.observable();

var lastPeriod = "";
var turbineval = [];


pc.InitFirst = function () {
    $.when(
        app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                pc.turbineList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                var dataturbine = [];
                if (res.data.length > 0) {
                    var allturbine = {}
                    $.each(res.data, function (key, val) {
                        turbineval.push(val);
                    });
                    // allturbine.value = "All Turbine";
                    // allturbine.text = "All Turbines";
                    // datavalue.push(allturbine);
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                        dataturbine.push(val);
                    });
                }
                pc.turbineList(datavalue);
                pc.turbine(dataturbine);
            }
        }),
        app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                pc.projectList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                if (res.data.length > 0) {
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                    });
                }
                pc.projectList(datavalue);
            }
        })

    ).then(function () {
        // $('#turbineList1').data('kendoDropDownList').value(["All Turbine"])
        // $('#turbineList2').data('kendoDropDownList').value(["All Turbine"])
        // override to set the value
        // $("#projectList1").data("kendoDropDownList").value("Tejuva");
        // $("#projectList2").data("kendoDropDownList").value("Tejuva");

        pc.project = $("#projectList").data("kendoDropDownList").value();
    });
}

pc.populateTurbine = function (data) {
    if (data.length == 0) {
        data = [];;
        pc.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        var dataturbine = [];
        if (data.length > 0) {
            var allturbine = {}
            $.each(data, function (key, val) {
                turbineval.push(val);
            });
            // allturbine.value = "All Turbine";
            // allturbine.text = "All Turbines";
            // datavalue.push(allturbine);
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
                dataturbine.push(val);
            });
        }
        pc.turbineList(datavalue);
        pc.turbine(dataturbine);
    }

    setTimeout(function () {
        // $('#turbineList1').data('kendoDropDownList').value(["All Turbine"])
        // $('#turbineList2').data('kendoDropDownList').value(["All Turbine"])
    }, 50);
};

pc.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        pc.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        pc.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            // $("#projectList").data("kendoDropDownList").select(1);
            // pc.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};


pc.showHidePeriod = function (callback) {
    var period = $('#periodList').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    // var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), 1, 0, 0, 0, 0));
    // var endMonthDate = new Date(app.toUTC(maxDateData));
    // var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    // var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    // var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    // var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
    var startMonthDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), 1, 0, 0, 0, 0);
    var endMonthDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0);
    var startYearDate = new Date(maxDateData.getFullYear(), 0, 1, 0, 0, 0, 0);
    var endYearDate = new Date(maxDateData.getFullYear(), 0, 1, 0, 0, 0, 0);
    var last24hours = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0);
    var lastweek = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0);
    if (period == "custom") {
        $(".show_hide").show();
        $('#dateStart').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#dateEnd').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });

        $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startMonthDate);
            $('#dateEnd').data('kendoDatePicker').value(endMonthDate);

            $(".show_hide").show();
        } else if (period == "annual") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",

            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startYearDate);
            $('#dateEnd').data('kendoDatePicker').value(endYearDate);

            $(".show_hide").show();
        } else {
            if (period == 'last24hours') {
                $('#dateStart').data('kendoDatePicker').value(last24hours);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            } else if (period == 'last7days') {
                $('#dateStart').data('kendoDatePicker').value(lastweek);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            }
            $(".show_hide").hide();
        }
        lastPeriod = period;
    }

    setTimeout(function () {
        callback;
    }, 50);
}

pc.showHidePeriod2 = function (callback) {
    var period = $('#periodList2').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    // var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), 1, 0, 0, 0, 0));
    // var endMonthDate = new Date(app.toUTC(maxDateData));
    // var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    // var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    // var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    // var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
    var startMonthDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), 1, 0, 0, 0, 0);
    var endMonthDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0);
    var startYearDate = new Date(maxDateData.getFullYear(), 0, 1, 0, 0, 0, 0);
    var endYearDate = new Date(maxDateData.getFullYear(), 0, 1, 0, 0, 0, 0);
    var last24hours = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0);
    var lastweek = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0);
    if (period == "custom") {
        $(".show_hide2").show();
        $('#dateStart2').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#dateEnd2').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });

        $('#dateStart2').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart2').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });
            $('#dateEnd2').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });

            $('#dateStart2').data('kendoDatePicker').value(startMonthDate);
            $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);

            $(".show_hide2").show();
        } else if (period == "annual") {
            $('#dateStart2').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",

            });
            $('#dateEnd2').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",
            });

            $('#dateStart2').data('kendoDatePicker').value(startYearDate);
            $('#dateEnd2').data('kendoDatePicker').value(endYearDate);

            $(".show_hide2").show();
        } else {
            if (period == 'last24hours') {
                $('#dateStart2').data('kendoDatePicker').value(last24hours);
                $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
            } else if (period == 'last7days') {
                $('#dateStart2').data('kendoDatePicker').value(lastweek);
                $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
            }
            $(".show_hide2").hide();
        }
        lastPeriod = period;
    }

    setTimeout(function () {
        callback;
    }, 50);
}


pc.InitDefaultValue = function () {
    // $("#projectList1").data("kendoDropDownList").value("Tejuva (24 | 50.4 MWh)")
    // $("#projectList2").data("kendoDropDownList").value("Tejuva (24 | 50.4 MWh)")
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    $("#periodList2").data("kendoDropDownList").value("custom");
    $("#periodList2").data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    // var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-30, 0, 0, 0, 0));
    // var lastEndDate = new Date(app.toUTC(maxDateData));
    var lastStartDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate()-30, 0, 0, 0);
    var lastEndDate = new Date(maxDateData.getFullYear(), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0);

    // var dateEnd2 = new Date(Date.UTC(moment(lastStartDate).get('year'), lastStartDate.getMonth(), lastStartDate.getDate()-30, 0, 0, 0, 0));
    // var dateStart2 =new Date(Date.UTC(moment(dateEnd2).get('year'), dateEnd2.getMonth(), dateEnd2.getDate()-30, 0, 0, 0, 0));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd2').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart2').data('kendoDatePicker').value(lastStartDate);
}
pc.initChart = function() {
        var p1DateStart = $('#dateStart').data('kendoDatePicker').value();
            // p1DateStart = new Date(Date.UTC(p1DateStart.getFullYear(), p1DateStart.getMonth(), p1DateStart.getDate(), 0, 0, 0));
            p1DateStart = new Date(p1DateStart.getFullYear(), p1DateStart.getMonth(), p1DateStart.getDate(), 0, 0, 0);

        var p1DateEnd  = $('#dateEnd').data('kendoDatePicker').value();
            // p1DateEnd = new Date(Date.UTC(p1DateEnd.getFullYear(), p1DateEnd.getMonth(), p1DateEnd.getDate(), 0, 0, 0));
            p1DateEnd = new Date(p1DateEnd.getFullYear(), p1DateEnd.getMonth(), p1DateEnd.getDate(), 0, 0, 0);

        var p2DateStart = $('#dateStart2').data('kendoDatePicker').value();
            // p2DateStart = new Date(Date.UTC(p2DateStart.getFullYear(), p2DateStart.getMonth(), p2DateStart.getDate(), 0, 0, 0));
            p2DateStart = new Date(p2DateStart.getFullYear(), p2DateStart.getMonth(), p2DateStart.getDate(), 0, 0, 0);

        var p2DateEnd  = $('#dateEnd2').data('kendoDatePicker').value();
            // p2DateEnd = new Date(Date.UTC(p2DateEnd.getFullYear(), p2DateEnd.getMonth(), p2DateEnd.getDate(), 0, 0, 0));
            p2DateEnd = new Date(p2DateEnd.getFullYear(), p2DateEnd.getMonth(), p2DateEnd.getDate(), 0, 0, 0);

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function(res) {
            if (!app.isFine(res)) {
                return;
            }
            var minDatetemp = new Date(res.data.ScadaData[0]);
            var maxDatetemp = new Date(res.data.ScadaData[1]);
            $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
        })

        var link = "analyticpowercurve/getlistpowercurvecomparison"

        app.loading(true);
        var param = {
            PC1Period       : $('#periodList').data('kendoDropDownList').value(),
            PC1Project      :  $("#projectList1").data("kendoDropDownList").value(),
            PC1Turbine      :  $("#turbineList1").data('kendoDropDownList').value(),// == "All Turbine" || $("#turbineList1").data('kendoDropDownList').value() == undefined ? pc.turbine() : $("#turbineList1").data('kendoDropDownList').value(),
            PC1DateStart    : p1DateStart,
            PC1DateEnd      : p1DateEnd,

            PC2Period       : $('#periodList2').data('kendoDropDownList').value(),
            PC2Project      :  $("#projectList1").data("kendoDropDownList").value(),
            PC2Turbine      : $("#turbineList2").data('kendoDropDownList').value(),// == "All Turbine" || $("#turbineList2").data('kendoDropDownList').value() == undefined  ? pc.turbine() : $("#turbineList2").data('kendoDropDownList').value(),
            PC2DateStart    : p2DateStart,
            PC2DateEnd      : p2DateEnd

        };

        toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            
            var dataTurbine = res.data.Data;

            
            $('#chartPCcomparison').html("");
            $("#chartPCcomparison").kendoChart({
                pdf: {
                  fileName: "DetailPowerCurve.pdf",
                },
                theme: "flat",
                title: {
                    text: "Power Curves",
                    visible: false,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                legend: {
                    position: "bottom",
                    visible: true,
                },
                chartArea: {
                    height: 375,
                },
                seriesDefaults: {
                    type: "scatterLine",
                    style: "smooth",
                    dashType: "longDash",
                    markers: {
                        visible: false,
                        size: 4,
                    },
                },
                seriesColors: colorField,
                series: dataTurbine,
                categoryAxis: {
                    labels: {
                        step: 1
                    }
                },
                valueAxis: [{
                    labels: {
                        format: "N0",
                    }
                }],
                xAxis: {
                    majorUnit: 1,
                    title: {
                        text: "Wind Speed (m/s)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555",
                        visible: true,
                    },
                    crosshair: {
                        visible: true,
                        tooltip: {
                            visible: true,
                            format: "N1",
                            background: "rgb(255,255,255, 0.9)",
                            color: "#58666e",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            border: {
                                color: "#eee",
                                width: "2px",
                            },
                        }
                    },
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    max: 25
                },
                yAxis: {
                    title: {
                        text: "Generation (KW)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555"
                    },
                    labels: {
                        format: "N0",
                    },
                    axisCrossingValue: -5,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    crosshair: {
                        visible: true,
                        tooltip: {
                            visible: true,
                            format: "N1",
                            background: "rgb(255,255,255, 0.9)",
                            color: "#58666e",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            border: {
                                color: "#eee",
                                width: "2px",
                            },
                        }
                    },
                },
                tooltip: {
                    visible: true,
                    template: "#= series.name #",
                    shared: true,
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                },
                pannable: false,
                zoomable: false
            });
            app.loading(false);
            $("#chartPCcomparison").data("kendoChart").refresh();

            
        });


}



$(document).ready(function () {
    $('#btnRefresh').on('click', function() {
        app.loading(true);
        setTimeout(function() {
            pc.initChart()
        }, 300);
    });

    app.loading(true);
    pc.InitDefaultValue();
    pc.initChart();
});