'use strict';

viewModel.FilterAnalytic = new Object();
var fa = viewModel.FilterAnalytic;

fa.turbineList = ko.observableArray([]);
fa.projectList = ko.observableArray([]);

fa.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Last 24 hours" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
fa.periodType = ko.observable();
fa.dateStart = ko.observable();
fa.dateEnd = ko.observable();
fa.turbine = ko.observableArray([]);
fa.project = ko.observable();
fa.period = ko.observable();
fa.infoPeriodRange = ko.observable();
fa.infoPeriodIcon = ko.observable(false);

var lastPeriod = "";
var turbineval = [];

fa.InitFirst = function () {
    $.when(
        app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                fa.turbineList([{ value: "", text: "" }]);
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
                fa.turbineList(datavalue);
            }
        }),
        app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                fa.projectList([{ value: "", text: "" }]);
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
                fa.projectList(datavalue);
            }
        })

    ).then(function () {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"])
        // override to set the value
        $("#projectList").data("kendoDropDownList").value("Tejuva");
        fa.project = $("#projectList").data("kendoDropDownList").value();
    });
}

/*fa.populateTurbine = function () {
    app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.length == 0) {
            res.data = [];;
            fa.turbineList([{ value: "", text: "" }]);
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
            fa.turbineList(datavalue);
        }

        setTimeout(function () {
            $('#turbineList').data('kendoMultiSelect').value(["All Turbine"])
        }, 300);
    });
};

fa.populateProject = function () {
    app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.length == 0) {
            res.data = [];;
            fa.projectList([{ value: "", text: "" }]);
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
            fa.projectList(datavalue);

            // override to set the value
            setTimeout(function () {
                $("#projectList").data("kendoDropDownList").value("Tejuva");
                fa.project = $("#projectList").data("kendoDropDownList").value();
            }, 300);
        }
    });
};*/

fa.populateTurbine = function (data) {
    if (data.length == 0) {
        data = [];;
        fa.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            var allturbine = {}
            $.each(data, function (key, val) {
                turbineval.push(val);
            });
            allturbine.value = "All Turbine";
            allturbine.text = "All Turbines";
            datavalue.push(allturbine);
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        fa.turbineList(datavalue);
    }

    setTimeout(function () {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"]);
    }, 300);
};

fa.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        fa.projectList([{ value: "", text: "" }]);
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
        fa.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").value("Tejuva");
            fa.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

fa.getProjectInfo = function () {
    var project = $("#projectList").data("kendoDropDownList").value();
    var turbines = $('#turbineList').data('kendoMultiSelect').value();

    if (turbines[0] == "All Turbine") {
        turbines = [];
    }

    var param = {
        Project: project,
        Turbines: turbines
    }

    app.ajaxPost(viewModel.appName + "helper/getprojectinfo", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        $("#project-info").html($("#projectList").data("kendoDropDownList").value());
        $("#total-turbine-info").html('<i class="fa fa-flash tooltipster tooltipstered" aria-hidden="true" title="Total Turbine"></i>&nbsp;' + res.data.TotalTurbine);
        $("#total-capacity-info").html('<i class="fa fa-tachometer tooltipster tooltipstered" aria-hidden="true" title="Total Capacity"></i>&nbsp;' + res.data.TotalCapacity + "MW");
    });
};

fa.showHidePeriod = function (callback) {
    var period = $('#periodList').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.toUTC(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
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
        // if (lastPeriod == "monthly") {
        //     $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        //     $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
        // } else if (lastPeriod == "annual") {
        //     $('#dateStart').data('kendoDatePicker').value(startYearDate);
        //     $('#dateEnd').data('kendoDatePicker').value(endYearDate);
        // }
        $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMMM yyyy",
            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMMM yyyy",
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

fa.LoadData = function () {
    if ($("#turbineList").data("kendoMultiSelect").value() == "") {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"])
    }

    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    if (fa.dateStart > fa.dateEnd) {
        toolkit.showError("Invalid Date Range Selection");
        return;
    } else {
        fa.InitFilter();
    }

    var period = $('#periodList').data('kendoDropDownList').value();

    fa.checkCompleteDate();
}

fa.checkTurbine = function () {
    var arr = $('#turbineList').data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbine");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#turbineList').data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    } else if (arr.length == 0) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    }
}

fa.InitFilter = function () {
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    fa.project = $("#projectList").data("kendoDropDownList").value();
    fa.period = $("#periodList").data("kendoDropDownList").value();
    fa.isDownTime = $("#isDownTime").is(":checked");

    if ($("#turbineList").data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
        fa.turbine = [];
    } else {
        fa.turbine = $("#turbineList").data("kendoMultiSelect").value();
    }

    fa.periodType = $("#periodList").data("kendoDropDownList").value();

    fa.GetBreakDown();
}

fa.InitDefaultValue = function () {
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-7, 0, 0, 0, 0));
    var lastEndDate = new Date(app.toUTC(maxDateData));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
}

fa.GetBreakDown = function () {
    fa.periodType = $("#periodList").data("kendoDropDownList").value();
    fa.project = $("#projectList").data("kendoDropDownList").value();
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    fa.dateStart = new Date(Date.UTC(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), fa.dateStart.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), 0, 0, 0));

    var result = [];

    if (fa.periodType == "last24hours" || fa.periodType == "last7days") {
        result.push({ "value": "Date", "text": "Date" });
    } else if (fa.periodType == "monthly" || fa.periodType == "annual") {
        result.push({ "value": "Month", "text": "Month" });
        result.push({ "value": "Year", "text": "Year" });
    } else if (fa.periodType == "custom") {
        if ((fa.dateEnd - fa.dateStart) / 86400000 + 1 <= 30) {
            result.push({ "value": "Date", "text": "Date" });
        }
        result.push({ "value": "Month", "text": "Month" });
        result.push({ "value": "Year", "text": "Year" });
    }

    if (fa.project == "") {
        result.push({ "value": "Project", "text": "Project" });
    } else {
        result.push({ "value": "Turbine", "text": "Turbine" });
    }

    return result;
}

fa.DateChange = function () {
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    fa.dateStart = new Date(Date.UTC(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), fa.dateStart.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), 0, 0, 0));
}

fa.checkCompleteDate = function () {
    var period = $('#periodList').data('kendoDropDownList').value();

    var monthNames = moment.months();

    var currentDateData = moment(app.currentDateData).format('YYYY-MM-DD');
    var today = moment().format('YYYY-MM-DD');
    var thisMonth = moment().get('month');
    var firstDayMonth = moment(new Date(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), 1)).format("YYYY-MM-DD");
    var lastDayMonth = moment(new Date(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth() + 1, 0)).format("YYYY-MM-DD");
    var firstDayYear = moment().startOf('year').format('YYYY-MM-DD');
    var endDayYear = moment().endOf('year').format('YYYY-MM-DD');

    var dateStart = moment(fa.dateStart).format('YYYY-MM-DD');
    var dateEnd = moment(fa.dateEnd).format('YYYY-MM-DD');

    if (period === 'custom') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range on start date and end date");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodRange("* Incomplete period data ange on start date");
            fa.infoPeriodIconmozilla(true);
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range on end date");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'annual') {
        if ((moment(fa.dateEnd).get('year') == moment(app.currentDateData).get('year')) && (currentDateData < today)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range in end year");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'monthly') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in start month and start month");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in start month");
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in end month");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else {
        fa.infoPeriodRange("");
        fa.infoPeriodIcon(false);
    }


}

$(document).ready(function () {
    app.loading(true);
    fa.showHidePeriod();
    // fa.populateTurbine();
    // fa.populateProject();
    fa.InitDefaultValue();
});