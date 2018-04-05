'use strict';

viewModel.EventAnalysis = new Object();
var ea = viewModel.EventAnalysis;

var dataByGroup;
var dataByError;
var dataByTurbine;
var breakDownEa = "detailgroup";
var additionalFilter = {};
var realDesc = {};
var categoryLvl1;
ea.firstLoad = ko.observable(true);
ea.labelEventDetail1 = ko.observable();
ea.labelEventDetail2 = ko.observable();


ea.data = {level0: {},level1 : {}, level2:{}}

vm.currentMenu('Event Analysis ');
vm.currentTitle('Event Analysis ');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Event Analysis', href: viewModel.appName + 'page/eventanalysis' }]);

var availDateAll;

ea.getDataAvailableInfo =  function(isFirstLoad){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateAll = res.data;

        ea.setAvailableDate(true);
    })
}

ea.setAvailableDate = function(isFirstLoad) {
    setTimeout(function(){
        var tipeTab = "Alarm";

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
ea.checkType = function(level){
    var value = $('input[name=convert'+level+']:checked').val();
    if(value == ("tohours"+level)){
        return "hours";
    }else{
        return "percentage";
    }
}

ea.autoGenerateLevel1 = function(params = {}, type) {
    app.loading(true);
    ea.RefreshData(params, type);
}

ea.autoGenerateLevel2 = function(params = {}, type) {
    app.loading(true);
    ea.RefreshData(params, type)
}

ea.RefreshData = function(params = {}, type){  
    var valid = fa.LoadData();
    if (valid) {        
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));

        
        breakDownEa = params.breakDownEa == undefined ? breakDownEa : params.breakDownEa;

        additionalFilter = params.additionalfilter == undefined ? additionalFilter : params.additionalfilter;

        var param = {
            period: fa.period,
            dateStart: dateStart,
            dateEnd: dateEnd,
            turbine: fa.turbine(),
            project: fa.project,
            breakdown: breakDownEa, // parent bdgroup , lvl1 => alarmdesc, lvl2 => turbine
            additionalfilter:  additionalFilter, // lvl 1 => { detailgroup : 'Machine'} , lvl 2 => {bdgroup : 'Machine', alarmdesc : 'Sembaranglah' }
            realdesc: realDesc
        }

        var reqData = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/geteventanalysistab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            var results;
            var labelAxis;

            ea.data[params.level]  = {
                hours : res.data.data,
                percentage : res.data.datapercentage
            }
            if(type == "percentage"){
                labelAxis = "%"
                results = res.data.datapercentage;
            }else{
                labelAxis = "Hours";
                results = res.data.data;
            }
            dataByGroup = _.sortBy(results, '_id');
            realDesc = res.data.realdesc;

            var chartId;

            if(params.level == "level0"){
                chartId = "chartEventAnalysis";
            }else if(params.level == "level1"){
                chartId = "chartEventAnalysisLevel1";
                ea.labelEventDetail1(additionalFilter.detailgroup);
            }else{
                chartId = "chartEventAnalysisLevel2";
                ea.labelEventDetail2(additionalFilter.alarmdesc);
            }
            
            ea.GenEventAnalysisChart(dataByGroup, chartId, "", labelAxis, false, "N2", params.level);
        });
    }
}

ea.refreshChart = function(type, axisLabel){
    var dataByGroup1 = _.sortBy(ea.data["level0"][type], '_id');
    var dataByGroup2 = _.sortBy(ea.data["level1"][type], '_id');
    var dataByGroup3 = _.sortBy(ea.data["level2"][type], '_id');
    ea.GenEventAnalysisChart(dataByGroup1, "chartEventAnalysis", "", axisLabel, false, "N2", "level0")
    ea.GenEventAnalysisChart(dataByGroup2, "chartEventAnalysisLevel1", "", axisLabel, false, "N2", "level1")
    ea.GenEventAnalysisChart(dataByGroup3, "chartEventAnalysisLevel2", "", axisLabel, false, "N2", "level2")
}

ea.GenEventAnalysisChart = function (dataSource,id,name,axisLabel, vislabel,format, level) {
    var CONTAINER_SIZE = 370;
    var LEGEND_SIZE = 100;
    var LEGEND_OFFSET = CONTAINER_SIZE - LEGEND_SIZE;

    // var colours = ["#ff880e","#21c4af","#b71c1c","#F0638B","#a2df53","#1c9ec4","#880d4e","#4a148c","#053872","#b1b2ac","#ffcf49","#605c5c","#b1b2ac","#ffcf49","#605c5c"];
    var colours = ["#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"];

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
            position: "custom",
            orientation: "horizontal",
            offsetY: LEGEND_OFFSET,
            labels: {              
                // template: "#: kendo.toString(replaceString(text))#",
                align: "center",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        // plotArea: plotArea,
        chartArea: {
            width: 330,
            // height: 800,
            padding: 0,
            margin: {
                top: -110,
            },
            background:"transparent"
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
        series: [{
            type: "pie",
            field: "result",
            categoryField: "_id"
        }],
        seriesColors: colours,
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
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none",
        },
        seriesClick : function (e){
            if(level == "level0"){
                ea.firstLoad(true);
                var param = {
                    level : "level1", 
                    additionalfilter:{ detailgroup : e.category },
                    breakDownEa:  "alarmdesc"
                };

                categoryLvl1 = e.category;

                ea.autoGenerateLevel2(param, ea.checkType(0));
            } else if(level == "level1"){
                var param = {
                    level : "level2", 
                    additionalfilter:{ detailgroup : categoryLvl1 , alarmdesc : e.category},
                    breakDownEa: "turbine"
                };
                app.loading(true);
                ea.RefreshData(param, ea.checkType(0));
            }
            
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, '" + format + "') # " + axisLabel,
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
        /* check for auto generating */
        if(ea.firstLoad() && level == "level0" ){
            setTimeout(function(){
                var category1 = $("#chartEventAnalysis").data("kendoChart").options.categoryAxis.categories[0];
                var paramLevel1 = {
                    level : "level1", 
                    additionalfilter:{ detailgroup : category1 },
                    breakDownEa:  "alarmdesc"
                };
                ea.RefreshData(paramLevel1, ea.checkType(0));
            }, 100);
        } else if (ea.firstLoad() && level == "level1" ){
            setTimeout(function(){
                var category2 = $("#chartEventAnalysisLevel1").data("kendoChart").options.categoryAxis.categories[0]
                var paramLevel2 = {
                    level : "level2", 
                    additionalfilter:{ detailgroup : categoryLvl1 , alarmdesc : category2},
                    breakDownEa: "turbine"
                };
                ea.RefreshData(paramLevel2, ea.checkType(0));
            },100);
        } else if (level == "level2") {
            setTimeout(function(){
                ea.firstLoad(false);
                app.loading(false);
            }, 500);
        }
    }, 100);
}

$(function() {

    $("input[name=convert0][value=tohours0]").prop('checked', true);
    $("input[name=convert1][value=tohours1]").prop('checked', true);
    $("input[name=convert2][value=tohours2]").prop('checked', true);

    $(".btnhours0").addClass("active");
    $(".btnhours1").addClass("active");
    $(".btnhours2").addClass("active");

    $('#btnRefresh').on('click', function () {
        ea.autoGenerateLevel1({level : "level0"}, ea.checkType(0));
        fa.checkTurbine();
    });

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { fa.showHidePeriod() }
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            ea.getDataAvailableInfo(true);
        }
    });

    setTimeout(function(){
        ea.getDataAvailableInfo(true);
        ea.autoGenerateLevel1({level : "level0"}, ea.checkType(0));
    }, 300)
});