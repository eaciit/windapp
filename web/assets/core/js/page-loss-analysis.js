'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.isDetailDTTop = ko.observable(false);

pg.availabledatestartscada = ko.observable();
pg.availabledateendscada = ko.observable();
pg.availabledatestartscada2 = ko.observable();
pg.availabledateendscada2 = ko.observable();

pg.availabledatestartalarm = ko.observable();
pg.availabledateendalarm = ko.observable();

pg.availabledatestartscada3 = ko.observable();
pg.availabledateendscada3 = ko.observable();

pg.availabledatestartalarm2 = ko.observable();
pg.availabledateendalarm2 = ko.observable();

pg.availabledatestartwarning = ko.observable();
pg.availabledateendwarning = ko.observable();

pg.labelAlarm = ko.observable("Downtime ");
var availDateListLoss = {};

var SeriesAlarm =  [{
    type: "pie",
    field: "result",
    categoryField: "_id",
}]

pg.isFirstStaticView = ko.observable(true);
pg.isFirstDowntime = ko.observable(true);
pg.isFirstAvailability = ko.observable(true);
pg.isFirstLostEnergy = ko.observable(true);
pg.isFirstReliability = ko.observable(true);
pg.isFirstWindSpeed = ko.observable(true);
pg.isFirstWarning = ko.observable(true);
pg.isFirstComponentAlarm = ko.observable(true);
// pg.isFirstEventAnalysis = ko.observable(true);
pg.isFirstMTBF = ko.observable(true);

var availDateAll;

pg.getDataAvailableInfo =  function(isFirstLoad){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateAll = res.data;

        pg.setAvailableDate(true);
    })
}

pg.setAvailableDate = function(isFirstLoad) {
    setTimeout(function(){
        var tabType = $(".panel-body").find(".nav-tabs").find("li.active").attr('id');
        var tipeTab = "ScadaData";
        switch (tabType) {
            case "staticViewTab" :
                tipeTab = "ScadaData"
                break;
            case "Top10DowntimeTab" :
                tipeTab = "Alarm"
                break;
            case "availabilityTab" :
                tipeTab = "ScadaData"
                break;
            case "lostenergyTab" :
                tipeTab = "Alarm"
                break;
            case "reliabilitykpiTab" :
                tipeTab = "ScadaData"
                break;
            case "windspeedavailTab" :
                tipeTab = "ScadaData"
                break;
            case "warningTab" :
                tipeTab = "Warning"
                break;
            case "CompAlarmTab" :
                tipeTab = "Alarm"
                break;
            case "mtbfTab" :
                tipeTab = "ScadaDataOEM"
                break;
            default:
                tipeTab = "ScadaData"
                break;
        }

        var namaproject = $("#projectList").data("kendoDropDownList").value();
        if(namaproject == "") {
            namaproject = "Tejuva";
        }

        var startDate = kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][0]).format('DD-MMM-YYYY'));
        var endDate = kendo.toString(moment.utc(availDateAll[namaproject][tipeTab][1]).format('DD-MMM-YYYY'));

        var maxDateData = new Date(availDateAll[namaproject][tipeTab][1]);

        var startDatepicker = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

        $('#availabledatestart').html(startDate);
        $('#availabledateend').html(endDate);

        if(isFirstLoad === true){
            $('#dateStart').data('kendoDatePicker').value(startDatepicker);
            $('#dateEnd').data('kendoDatePicker').value(endDate);  
        }

    }, 500);
}
pg.backToDownTime = function () {
    pg.isDetailDTTop(false);
    pg.detailDTTopTxt("");
}

pg.LoadData = function(){
    fa.LoadData();
    $.when(fa.LoadData()).done(function(){
        setTimeout(function(){
            fa.getDataAvailability();
        }, 300);
    });
}

pg.Reliability = function(){
    if(pg.isFirstReliability() === true){
        $('#availabledatestart').html(pg.availabledatestartscada2());
        $('#availabledateend').html(pg.availabledateendscada2());
    }else{
        $('#availabledatestart').html(pg.availabledatestartscada2());
        $('#availabledateend').html(pg.availabledateendscada2());
    }
}

pg.showFilter = function(){
    $("#periodList").closest(".k-widget").show();
    $("#dateStart").closest(".k-widget").show();
    $("#dateEnd").closest(".k-widget").show();
    $(".control-label:contains('Period')").show();
    $(".control-label:contains('to')").show();
}

pg.resetStatus = function(){
    pg.isFirstStaticView(true);
    pg.isFirstDowntime(true);
    pg.isFirstAvailability(true);
    pg.isFirstLostEnergy(true);
    pg.isFirstReliability(true);
    pg.isFirstWindSpeed(true);
    pg.isFirstWarning(true);
    pg.isFirstComponentAlarm(true);
    // pg.isFirstEventAnalysis(true);
    pg.isFirstMTBF(true);
}
vm.currentMenu('Losses and Efficiencies ');
vm.currentTitle('Losses and Efficiencies ');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticloss' }]);

function replaceString(value) {
    return value.replace(/_/gi, "  ");
}

$(function(){

    pg.getDataAvailableInfo(true);

    $('#btnRefresh').on('click', function () {
        app.loading(true);
        pg.LoadData();
        fa.checkTurbine();
        pg.resetStatus();
        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $('#breakdownlistavail').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.isFirstAvailability(true); av.Availability(); },
    });

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { fa.showHidePeriod(av.SetBreakDown()) }
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            // app.loading(true);
            av.SetBreakDown();
            pg.resetStatus();
            pg.getDataAvailableInfo(true);

            // setTimeout(function(){
            //     fa.checkTurbine();
            //     pg.resetStatus();
            //     $('.nav').find('li.active').find('a').trigger( "click" );
            // },1000)
        }
    });

    $("#dateStart").change(function () { fa.DateChange(av.SetBreakDown()) });
    $("#dateEnd").change(function () { fa.DateChange(av.SetBreakDown()) });

    $("input[name=IsAlarm]").on("change", function() {
        var HAlarm = $('#filter-analytic').width() * 0.235
        var wAll = $('#filter-analytic').width() * 0.275
    
        var data = ca.dtCompponentAlarm()
        if(this.id == "alarm"){   
            SeriesAlarm =  [{
                field: "result",
                name: "Downtime"
            }]             
            // ===== Alarm =====
            dt.GenChartDownAlarmComponent("alarm",data.alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll,"N1");
            dt.GenChartDownAlarmComponent("alarm",data.alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll,"N0");
            dt.GenChartDownAlarmComponent("alarm",data.alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll,"N1");

            pg.labelAlarm(" Top 10 Downtime")
        }else{     
             SeriesAlarm = [{
                type: "pie",
                field: "result",
                categoryField: "_id",
            }]           
            // ===== Component =====
            var componentduration = _.sortBy(data.componentduration, '_id');
            var componentfrequency = _.sortBy(data.componentfrequency, '_id');
            var componentloss = _.sortBy(data.componentloss, '_id');
            dt.GenChartDownAlarmComponent("component",componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll,"N1");
            dt.GenChartDownAlarmComponent("component",componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll,"N0");
            dt.GenChartDownAlarmComponent("component",componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll,"N1");

            pg.labelAlarm(" Downtime")
        }
    });

    $("input[name=convertStatic]").on("change", function() {
        if($("#gridStatic").is(':checked')) {
            sv.isGrid(true);
        } else {
            sv.isGrid(false);
        }

        var view= this.id;
        
        setTimeout(function(){
            sv.refreshView(view);
        },300);
    });


    setTimeout(function(){
        pg.LoadData();
        sv.StaticView();
    },500);
    /*$(window).resize(function() {
        $("#chartCADuration").data("kendoChart").refresh();
    });*/

})
