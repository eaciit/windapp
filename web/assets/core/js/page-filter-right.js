'use strict';

viewModel.FilterRight = new Object();
var fr = viewModel.FilterRight;

var keys = [
    { "value": "Actual Production", "text": "Production", "unit": "MWh" },
    { "value": "Actual PLF", "text": "PLF", "unit": "%" },
    { "value": "Total Availability", "text": "Total Availability", "unit": "%" },
    { "value": "Grid Availability", "text": "Grid Availability", "unit": "%" },
    { "value": "Machine Availability", "text": "Machine Availability", "unit": "%" },
    { "value": "Data Availability", "text": "Data Availability", "unit": "%" },
    // { "value": "MTTR / MTTF", "text": "MTTR / MTTF" },
    { "value": "P50 Generation", "text": "P50 Generation", "unit": "MWh" },
    { "value": "P50 PLF", "text": "P50 PLF", "unit": "%" },
    { "value": "P75 Generation", "text": "P75 Generation", "unit": "MWh" },
    { "value": "P75 PLF", "text": "P75 PLF", "unit": "%h" },
    { "value": "P90 Generation", "text": "P90 Generation", "unit": "MWh" },
    { "value": "P90 PLF", "text": "P90 PLF", "unit": "%" },
];
fr.key1 = ko.observableArray([]);
fr.key2 = ko.observableArray([]);
fr.key1(keys);
fr.key2(keys);
var isFirst = true;

fr.InitDefaultValue = function () {
    $("#key1").data("kendoDropDownList").value("Actual Production");
    $("#key2").data("kendoDropDownList").value("Actual PLF");
}

fr.breakdown = ko.observableArray([
    { "value": "$dateinfo.dateid", "text": "Date" },
    { "value": "$dateinfo.monthid", "text": "Month" },
    { "value": "$dateinfo.year", "text": "Year" },
    { "value": "$turbine", "text": "Turbine" },
    { "value": "$projectname", "text": "Project" }
]);

fr.getBreakDown = function () {
    fa.periodType = $("#periodList").data("kendoDropDownList").value();
    fa.project = $("#projectList").data("kendoDropDownList").value();

    var result = [];
    if (fa.periodType == "today" || fa.periodType == "lastday" || fa.periodType == "lastweek" || fa.periodType == "lastmonth") {
        result.push({ "value": "$dateinfo.dateid", "text": "Date" });
        result.push({ "value": "$dateinfo.monthid", "text": "Month" });
        result.push({ "value": "$dateinfo.year", "text": "Year" });
    } else if (fa.periodType == "lastyear") {
        result.push({ "value": "$dateinfo.monthid", "text": "Month" });
        result.push({ "value": "$dateinfo.year", "text": "Year" });
    } else if (fa.periodType == "custom") {
        if ((fa.dateEnd - fa.dateStart) / 86400000 <= 31) {
            result.push({ "value": "$dateinfo.dateid", "text": "Date" });
        }
        result.push({ "value": "$dateinfo.monthid", "text": "Month" });
        result.push({ "value": "$dateinfo.year", "text": "Year" });
    }

    if (fa.project !== "") {
        result.push({ "value": "$turbine", "text": "Turbine" });
    } else {
        result.push({ "value": "$projectname", "text": "Project" });
    }

    return result;
}

fr.checkKey1 = function () {
    var key1 = $("#key1").data("kendoDropDownList").value();
    var key2 = $("#key2").data("kendoDropDownList").value();
    fr.key2([]);
    $.each(keys, function (i) {
        if (keys[i].value == key1) {
            return true;
        }
        fr.key2.push(keys[i]);
    });
    $("#key2").data("kendoDropDownList").value(key2);
    if (isFirst == false) {
        km.getData();
    }
}

fr.checkKey2 = function () {
    var key2 = $("#key2").data("kendoDropDownList").value();
    var key1 = $("#key1").data("kendoDropDownList").value();
    fr.key1([]);
    $.each(keys, function (i) {
        if (keys[i].value == key2) {
            return true;
        }
        fr.key1.push(keys[i]);
    });
    $("#key1").data("kendoDropDownList").value(key1);
    if (isFirst == false) {
        km.getData();
    }
}

$(document).ready(function () {
    fr.InitDefaultValue();
    fr.checkKey1();
    fr.checkKey2();
    isFirst = false;
});