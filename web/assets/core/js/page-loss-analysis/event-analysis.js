'use strict';

viewModel.EventAnalysis = new Object();
var ea = viewModel.EventAnalysis;

var dataByGroup;
var dataByError;
var dataByTurbine;
var breakDownEa = "bdgroup";
var additionalFilter = {};
var realDesc = {};
ea.labelEvent = ko.observable("Event Analysis by Detail Group");

ea.RefreshData = function(){    
    var valid = fa.LoadData();
    if (valid) {
        app.loading(true);
        pg.setAvailableDate(false);
        if(pg.isFirstEventAnalysis() === true){
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));

            var param = {
                period: fa.period,
                dateStart: dateStart,
                dateEnd: dateEnd,
                turbine: fa.turbine(),
                project: fa.project,
                breakdown: breakDownEa,
                additionalfilter: additionalFilter,
                realdesc: realDesc
            }

            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/geteventanalysistab", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                setTimeout(function(){
                    dataByGroup = _.sortBy(res.data.data, '_id');
                    realDesc = res.data.realdesc;
                    ea.labelEvent("Event Analysis by Detail Group");

                    ea.GenEventAnalysisChart(dataByGroup, 'chartEventAnalysis', "", "Hours", false, "N1");

                    app.loading(false);
                    pg.isFirstEventAnalysis(false);
                },300);
            }); 
        }else{
            setTimeout(function(){
                $("#chartEventAnalysis").data("kendoChart").refresh();
                app.loading(false);
            },200); 
        }
    }
}


ea.GenEventAnalysisChart = function (dataSource,id,name,axisLabel, vislabel,format) {
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
            position: "top",
            align: "center",
            visible: true,
            labels: {              
                // template: "#: kendo.toString(replaceString(text))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        chartArea: {
            // height: heightParam, 
            // width: wParam, 
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
                // template: "#: kendo.toString(replaceString(value))#",
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
            // template: "#: kendo.toString(replaceString(category)) #: #: kendo.toString(value, '" + format + "') # " + axisLabel,
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