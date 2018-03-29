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
ea.firsLoad = ko.observable(true);
ea.labelEvent = ko.observable("Event Analysis by Detail Group");


ea.LoadData = function(){
    if(pg.isFirstEventAnalysis() === true){
        ea.RefreshData();
    }else{
        setTimeout(function(){
            $("#chartEventAnalysis").data("kendoChart").refresh();
            app.loading(false);
        },200); 
    }
}

ea.RefreshData = function(params = {}){  
    var valid = fa.LoadData();
    if (valid) {
        console.log(params);
        app.loading(true);
        pg.setAvailableDate(false);
        
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

        console.log(param);


        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/geteventanalysistab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                dataByGroup = _.sortBy(res.data.data, '_id');
                realDesc = res.data.realdesc;

                var chartId;

                if(params.level == undefined){
                    chartId = "chartEventAnalysis";
                    ea.labelEvent("Event Analysis by Detail Group");
                }else if(params.level == 1){
                    chartId = "chartEventAnalysisLevel1";
                     ea.labelEvent("Event Analysis by Detail Group > "+ additionalFilter.detailgroup);
                }else{
                    chartId = "chartEventAnalysisLevel2";
                     ea.labelEvent("Event Analysis by Detail Group > "+ additionalFilter.detailgroup + " > " + additionalFilter.alarmdesc );
                }
                

                if(ea.firsLoad() == true){
                    $.when(ea.GenEventAnalysisChart(dataByGroup, chartId, "", "Hours", false, "N1", params.level)).done(function(){

                    }); 
                }

                app.loading(false);
                pg.isFirstEventAnalysis(false);
            },300);
        }); 
        
    }
}


ea.GenEventAnalysisChart = function (dataSource,id,name,axisLabel, vislabel,format, level) {

    var CONTAINER_SIZE = 300;
    var LEGEND_SIZE = 50;
    var LEGEND_OFFSET = CONTAINER_SIZE - LEGEND_SIZE;

    var legend = {
        position: "custom",
        orientation: "horizontal",
        offsetY: LEGEND_OFFSET
    };

    var plotArea = {
        height: LEGEND_OFFSET
    };

    var colours = ["#ff880e","#21c4af","#b71c1c","#F0638B","#a2df53","#1c9ec4","#880d4e","#4a148c","#053872","#b1b2ac","#ffcf49","#605c5c","#b1b2ac","#ffcf49","#605c5c"];
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
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        plotArea: plotArea,
        chartArea: {
            width: 300,
            height: 300,
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
            if(level == undefined){
                var param = {
                    level : 1, 
                    additionalfilter:{ detailgroup : e.category },
                    breakDownEa:  "alarmdesc"
                };

                categoryLvl1 = e.category;
                ea.RefreshData(param);

            }else if(level == 1){
                var param = {
                    level : 2, 
                    additionalfilter:{ detailgroup : categoryLvl1 , alarmdesc : e.category},
                    breakDownEa: "turbine"
                };

                ea.RefreshData(param);
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
    }, 100);
}