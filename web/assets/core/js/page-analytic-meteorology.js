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
    { "category": "0 to 4m/s", "color": colorFieldsWR[0] , "id" : "wr1" },
    { "category": "4 to 8m/s", "color": colorFieldsWR[1], "id" : "wr2"},
    { "category": "8 to 12m/s", "color": colorFieldsWR[2], "id" : "wr3"},
    { "category": "12 to 16m/s", "color": colorFieldsWR[3],"id" : "wr4"},
    { "category": "16 to 20m/s", "color": colorFieldsWR[4],"id" : "wr5"},
    { "category": "20m/s and above", "color": colorFieldsWR[5],"id" : "wr6"},
];

pm.valueCategory = ko.observableArray([
    { "value": "powerGeneration", "text": "Power Generation (MW)" },
    { "value": "machine", "text": "Machine Availability" },
    { "value": "scada", "text": "Scada Availability" },
    { "value": "grid", "text": "Grid Availability" },
]);

// var color = ["#87c5da","#cc2a35", "#d66b76", "#5d1b62", "#f1c175","#95204c","#8f4bc5","#7d287d","#00818e","#c8c8c8","#546698","#66c99a","#f3d752","#20adb8","#333d6b","#d077b1","#aab664","#01a278","#c1d41a","#807063","#ff5975","#01a3d4","#ca9d08","#026e51","#4c653f","#007ca7"];
// var color = ["#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"]
var color = ["#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"]
pm.isMet = ko.observable(false);
pm.isFirstAverage = ko.observable(true);
pm.isFirstWindRose = ko.observable(true);
pm.isFirstWindDis = ko.observable(true);
pm.isFirstNacelleDis = ko.observable(true);
pm.isFirstTurbulence = ko.observable(true);
pm.isFirstTemperature = ko.observable(true);
pm.isFirstTurbine = ko.observable(true);
pm.isFirstEnergy = ko.observable(true);
pm.isFirstTwelve = ko.observable(true);
pm.isFirstWindRoseComparison = ko.observable(true);
pm.projectName = ko.observable();
pm.dateEnd = ko.observable();
pm.dateStart = ko.observable();


pm.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data;

        var namaproject = $("#projectList").data("kendoDropDownList").value();
        if(namaproject == "") {
            namaproject = "Tejuva";
        }

        availDateList.availabledatestartscada = kendo.toString(moment.utc(data[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY'));
        availDateList.availabledateendscada = kendo.toString(moment.utc(data[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY'));

        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

        availDateList.startScadaHFD = kendo.toString(moment.utc(data[namaproject]["ScadaDataHFD"][0]).format('DD-MMM-YYYY'));
        availDateList.endScadaHFD = kendo.toString(moment.utc(data[namaproject]["ScadaDataHFD"][1]).format('DD-MMM-YYYY'));

        availDateList.startScadaOEM = kendo.toString(moment.utc(data[namaproject]["ScadaDataOEM"][0]).format('DD-MMM-YYYY'));
        availDateList.endScadaOEM = kendo.toString(moment.utc(data[namaproject]["ScadaDataOEM"][1]).format('DD-MMM-YYYY'));

        availDateList.availabledatestartmet = kendo.toString(moment.utc(data[namaproject]["MET"][0]).format('DD-MMM-YYYY'));
        availDateList.availabledateendmet = kendo.toString(moment.utc(data[namaproject]["MET"][1]).format('DD-MMM-YYYY'));


        var startDate = kendo.toString(moment.utc(data[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY'));
        var endDate = kendo.toString(moment.utc(data[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY'));

        var maxDateData = new Date(data[namaproject]["ScadaData"][1]);

        var startDatepicker = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));


        $('#dateStart').data('kendoDatePicker').value(startDatepicker);
        $('#dateEnd').data('kendoDatePicker').value(endDate);  

    })
}
pm.loadData = function () {
    setTimeout(function () {
        if (fa.project == "") {
            pm.type = "Project Name";
        } else {
            pm.type = "Turbine";
        }

    }, 100);
}

pm.resetStatus= function(){
    pm.isFirstAverage(true);
    pm.isFirstWindRose(true);
    pm.isFirstWindDis(true);
    pm.isFirstNacelleDis(true);
    pm.isFirstTurbulence(true);
    pm.isFirstTemperature(true);
    pm.isFirstTurbine(true);
    pm.isFirstEnergy(true);
    pm.isFirstTwelve(true);
    pm.isFirstWindRoseComparison(true);

    // reset for turbulence intensity
    pm.ShowScatter(false);
    $('#wCbScatter').find('input[type=checkbox]').removeAttr('checked');
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

pm.addslashes = function ( str ) {
    var slash = "\\‍";

    var string = "["+slash+"\""+str+slash+"\"]"
    // console.log(string);
    return string;
}
pm.sortObject = function(o) {
    var sorted = {},
    key, a = [];

    for (key in o) {
        if (o.hasOwnProperty(key)) {
            a.push(key);
        }
    }

    a.sort();

    for (key = 0; key < a.length; key++) {
        sorted[a[key]] = o[a[key]];
    }
    return sorted;
}
$(function(){
    pm.getAvailDate();

    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        pm.resetStatus();
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value(); 
        var project = $('#projectList').data("kendoDropDownList").value();

        pm.dateEnd(dateEnd);
        pm.dateStart(dateStart);
        pm.projectName(project);

        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $("input[name=isMet]").on("change", function() {
        if($("#gridDineural").is(':checked')) {
            tb.generateGridTable(this.id);
        } else {
            tb.generateGraph();
        }
        
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

    $("input[name=convertTable]").on("change", function() {
        if($("#gridDineural").is(':checked')) {
            tb.generateGridTable();
        } else {
            tb.generateGraph();
        }
    });

    setTimeout(function () {
        pm.loadData();
        aw.AverageWindSpeed();
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value(); 
        var project = $('#projectList').data("kendoDropDownList").value();

        pm.dateEnd(dateEnd);
        pm.dateStart(dateStart);
        pm.projectName(project);

        $('#projectList').kendoDropDownList({
            change: function () {  
                pm.resetStatus();
                pm.getAvailDate();
                var project = $('#projectList').data("kendoDropDownList").value();
                fa.populateTurbine(project);
                pm.projectName(project);


            }
        });

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
    }, 500);
});