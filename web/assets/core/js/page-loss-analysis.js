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


pg.getDataAvailableInfo =  function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var minDatetemp = new Date(res.data.ScadaData[0]);
        var maxDatetemp = new Date(res.data.ScadaData[1]);

        pg.availabledatestartscada(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartscada2(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada2(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartalarm(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendalarm(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

        pg.availabledatestartscada3(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada3(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartalarm2(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendalarm2(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

        pg.availabledatestartwarning(kendo.toString(moment.utc(res.data.Warning[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendwarning(kendo.toString(moment.utc(res.data.Warning[1]).format('DD-MMMM-YYYY')));

        $('#availabledatestart').html(pg.availabledatestartscada());
        $('#availabledateend').html(pg.availabledateendscada());
    })
}
pg.backToDownTime = function () {
    pg.isDetailDTTop(false);
    pg.detailDTTopTxt("");
}

pg.LoadData = function(){
    fa.LoadData();
    if (fa.project == "") {
        sv.type = "Project Name";
    } else {
        sv.type = "Turbine";
    }
    pg.getDataAvailableInfo();
}
pg.GenChartDownAlarmComponent = function (dataSource,id,Series,legend,name,axisLabel, vislabel,rotate,heightParam,wParam,format) {

    $("#" + id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: [
                { "field": "Total", "dir": "desc" }
            ],
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: legend,
            labels: {              
                template: "#: kendo.toString(replaceString(text))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        chartArea: {
            height: heightParam, 
            width: wParam, 
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            labels: {
                        visible: vislabel,
                        background: "transparent",
                        template: "#= category #: \n #= kendo.format('{0:" + format + "}', value)# " + axisLabel,
                    }
        },
        series: Series,
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: axisLabel,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                step: 2,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            }
        },
        categoryAxis: {
            // visible: legend,
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: rotate,                
                template: "#: kendo.toString(replaceString(value))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none",
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: kendo.toString(replaceString(category)) #: #: kendo.toString(value, '" + format + "') # " + axisLabel,
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });

    setTimeout(function () {
        if ($("#" + id).data("kendoChart") != null) {
            $("#" + id).data("kendoChart").refresh();
        }
    }, 100);
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

pg.resetStatus = function(){
    pg.isFirstStaticView(true);
    pg.isFirstDowntime(true);
    pg.isFirstAvailability(true);
    pg.isFirstLostEnergy(true);
    pg.isFirstReliability(true);
    pg.isFirstWindSpeed(true);
    pg.isFirstWarning(true);
    pg.isFirstComponentAlarm(true);
}
vm.currentMenu('Losses and Efficiency');
vm.currentTitle('Losses and Efficiency');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticloss' }]);

function replaceString(value) {
    return value.replace(/_/gi, "  ");
}

$(function(){
    setTimeout(function(){
        pg.LoadData();
        sv.StaticView();
    },200);

    $('#btnRefresh').on('click', function () {
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
        change: function () { av.SetBreakDown() }
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
            pg.GenChartDownAlarmComponent(data.alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll,"N1");
            pg.GenChartDownAlarmComponent(data.alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll,"N0");
            pg.GenChartDownAlarmComponent(data.alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll,"N1");

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
            pg.GenChartDownAlarmComponent(componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll,"N1");
            pg.GenChartDownAlarmComponent(componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll,"N0");
            pg.GenChartDownAlarmComponent(componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll,"N1");

            pg.labelAlarm(" Downtime")
        }
    });

    /*$(window).resize(function() {
        $("#chartCADuration").data("kendoChart").refresh();
    });*/

})
