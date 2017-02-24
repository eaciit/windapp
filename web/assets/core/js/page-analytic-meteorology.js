'use strict';

viewModel.MeteorologyAnalysis = new Object();
var pm = viewModel.MeteorologyAnalysis;

vm.currentMenu('Meteorology');
vm.currentTitle('Meteorology');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Meteorology', href: viewModel.appName + 'page/analyticmeteorology' }]);

var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));
var availDateList = {};

pm.type = ko.observable();
pm.detailDTTopTxt = ko.observable();
pm.isDetailDTTop = ko.observable(false);
pm.periodDesc = ko.observable();
pm.breakDown = ko.observableArray([]);
pm.breakDownList = ko.observableArray([
    { "value": "dateinfo.dateid", "text": "Date" },
    { "value": "dateinfo.monthdesc", "text": "Month" },
    { "value": "dateinfo.year", "text": "Year" },
    { "value": "projectname", "text": "Project" },
    { "value": "turbine", "text": "Turbine" },
]);
pm.sectionsBreakdownList = ko.observableArray([
    { "text": 36, "value": 36 },
    { "text": 24, "value": 24 },
    { "text": 12, "value": 12 },
]);
// var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
var colorFieldsWR = ["#224373","#186baa","#25FEDF", "#f5a265","#eb5b19", "#9E0000","#e4cc37"];
var listOfCategory = [
    { "category": "0 to 4m/s", "color": colorFieldsWR[0] },
    { "category": "4 to 8m/s", "color": colorFieldsWR[1] },
    { "category": "8 to 12m/s", "color": colorFieldsWR[2] },
    { "category": "12 to 16m/s", "color": colorFieldsWR[3] },
    { "category": "16 to 20m/s", "color": colorFieldsWR[4] },
    { "category": "20m/s and above", "color": colorFieldsWR[5] },
];

pm.valueCategory = ko.observableArray([
    { "value": "powerGeneration", "text": "Power Generation (MW)" },
    { "value": "machine", "text": "Machine Availability" },
    { "value": "scada", "text": "Scada Availability" },
    { "value": "grid", "text": "Grid Availability" },
]);

var color = ["#87c5da","#cc2a35", "#d66b76", "#5d1b62", "#f1c175","#95204c","#8f4bc5","#7d287d","#00818e","#c8c8c8","#546698","#66c99a","#f3d752","#20adb8","#333d6b","#d077b1","#aab664","#01a278","#c1d41a","#807063","#ff5975","#01a3d4","#ca9d08","#026e51","#4c653f","#007ca7"];

pm.isMet = ko.observable(false);
pm.isFirstAverage = ko.observable(true);
pm.isFirstWindRose = ko.observable(true);
pm.isFirstWindDis = ko.observable(true);
pm.isFirstTurbulence = ko.observable(true);
pm.isFirstTemperature = ko.observable(true);
pm.isFirstTurbine = ko.observable(true);
pm.isFirstTwelve = ko.observable(true);
pm.isFirstWindRoseComparison = ko.observable(true);
pm.isFirstMTBF = ko.observable(true);


pm.loadData = function () {
    setTimeout(function () {
        if (fa.project == "") {
            pm.type = "Project Name";
        } else {
            pm.type = "Turbine";
        }

    }, 100);
    
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateList.availabledatestartscada = kendo.toString(moment.utc(res.data.ScadaData[0]).format('DD-MMMM-YYYY'));
        availDateList.availabledateendscada = kendo.toString(moment.utc(res.data.ScadaData[1]).format('DD-MMMM-YYYY'));

        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

        availDateList.startScadaHFD = kendo.toString(moment.utc(res.data.ScadaDataHFD[0]).format('DD-MMMM-YYYY'));
        availDateList.endScadaHFD = kendo.toString(moment.utc(res.data.ScadaDataHFD[1]).format('DD-MMMM-YYYY'));

        availDateList.startScadaOEM = kendo.toString(moment.utc(res.data.ScadaDataOEM[0]).format('DD-MMMM-YYYY'));
        availDateList.endScadaOEM = kendo.toString(moment.utc(res.data.ScadaDataOEM[1]).format('DD-MMMM-YYYY'));

        availDateList.availabledatestartmet = kendo.toString(moment.utc(res.data.MET[0]).format('DD-MMMM-YYYY'));
        availDateList.availabledateendmet = kendo.toString(moment.utc(res.data.MET[1]).format('DD-MMMM-YYYY'));
    })
}

pm.resetStatus= function(){
    pm.isFirstAverage(true);
    pm.isFirstWindRose(true);
    pm.isFirstWindDis(true);
    pm.isFirstTurbulence(true);
    pm.isFirstMTBF(true);
    pm.isFirstTemperature(true);
    pm.isFirstTurbine(true);
    pm.isFirstTwelve(true);
    pm.isFirstWindRoseComparison(true);
}
pm.showFilter = function(){
    $("#periodList").closest(".k-widget").show();
    $("#dateStart").closest(".k-widget").show();
    $("#dateEnd").closest(".k-widget").show();
    $(".control-label:contains('Period')").show();
    $(".control-label:contains('to')").show();
}
pm.hideFilter = function(){
    $("#periodList").closest(".k-widget").hide();
    $("#dateStart").closest(".k-widget").hide();
    $("#dateEnd").closest(".k-widget").hide();
    $(".control-label:contains('Period')").hide();
    $(".control-label:contains('to')").hide();
}
$(function(){
    pm.loadData();
    aw.AverageWindSpeed();

    $('#btnRefresh').on('click', function () {
        pm.resetStatus();
        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $("input[name=isMet]").on("change", function() {
        tb.generateGridTable(this.id);
        if($("#met").is(':checked')) {
            pm.isMet(true);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
        } else {
             pm.isMet(false);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        }
    });

    setTimeout(function () {
        $("#legend-list").html("");
        $.each(listOfCategory, function (idx, val) {
            var idName = "btn" + idx;
            listOfButton[idName] = true;
            $("#legend-list").append(
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="wr.showHideLegendWR(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
                '<span class="span-legend">' + val.category + '</span>'
            );
        });
        $("#nosection").data("kendoDropDownList").value(12);
    }, 300);
});