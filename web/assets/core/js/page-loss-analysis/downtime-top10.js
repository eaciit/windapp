'use strict';

viewModel.Downtime = new Object();
var dt = viewModel.Downtime;
var SeriesDowntime = [
// {
//     field: "AEBOK",
//     name: "AEBOK"
// }, 
// {
//     field: "ExternalStop",
//     name: "External Stop"
// }, 
{
    field: "GridDown",
    name: "Grid Down"
}, 
// {
//     field: "InternalGrid",
//     name: "InternalGrid"
// }, 
{
    field: "MachineDown",
    name: "Machine Down"
}, 
// {
//     field: "WeatherStop",
//     name: "Weather Stop"
// }, 
{
    field: "Unknown",
    name: "Unknown"
}]

dt.GenChartDownAlarmComponent = function (tab, dataSource,id,Series,legend,name,axisLabel, vislabel,rotate,heightParam,wParam,format) {
    // var colours = ["#ff880e", "#21c4af", "#f44336","#feb64e","#a2df53", "#69d2e7","#4589b0","#ed5784"];
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
            visible: legend,
            labels: {              
                // template: "#: kendo.toString(replaceString(text))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        chartArea: {
            height: heightParam, 
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
        series: Series,
        seriesColors: (tab == "component" ? colours : colorField),
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

dt.Downtime = function(){
    var valid = fa.LoadData();
    if (valid) {
        
        if(pg.isFirstDowntime() === true){
            app.loading(true);

            // pg.setAvailableDate(true);

            setTimeout(function(){
                var dateStart = $('#dateStart').data('kendoDatePicker').value();
                var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 
                var param = {
                    period: fa.period,
                    dateStart: dateStart,
                    dateEnd: dateEnd,
                    turbine: fa.turbine(),
                    project: fa.project,
                }

                var reqDuration = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getdowntimetabduration", param, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    setTimeout(function(){
                        var HDowntime = $('#filter-analytic').width() * 0.2
                        var wAll = $('#filter-analytic').width() * 0.275

                        /*Downtime Tab*/
                        dt.GenChartDownAlarmComponent("downtime",res.data.duration,'chartDTDuration',SeriesDowntime,true,"Turbine", "Hours",false,-330,HDowntime,wAll,"N1");

                        pg.isFirstDowntime(false);
                    },300);
                });
                var reqFreq = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getdowntimetabfreq", param, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    setTimeout(function(){
                        var HDowntime = $('#filter-analytic').width() * 0.2
                        var wAll = $('#filter-analytic').width() * 0.275

                        /*Downtime Tab*/
                        dt.GenChartDownAlarmComponent("downtime",res.data.frequency,'chartDTFrequency',SeriesDowntime,true,"Turbine", "Times",false,-330,HDowntime,wAll,"N0");

                        pg.isFirstDowntime(false);
                    },300);
                });
                var reqLoss = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getdowntimetabloss", param, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    setTimeout(function(){
                        var HDowntime = $('#filter-analytic').width() * 0.2
                        var wAll = $('#filter-analytic').width() * 0.275

                        /*Downtime Tab*/
                        dt.GenChartDownAlarmComponent("downtime",res.data.loss,'chartTopTurbineLoss',SeriesDowntime,true,"Turbine","MWh",false,-330,HDowntime,wAll,"N1");

                        pg.isFirstDowntime(false);
                    },300);
                });
                $.when(reqDuration, reqFreq, reqLoss).done(function(){
                    setTimeout(function(){
                        app.loading(false);
                    }, 100);
                });
            },1000);
        }else{
            pg.setAvailableDate(false);
            setTimeout(function(){
                $("#chartDTDuration").data("kendoChart").refresh();
                $("#chartDTFrequency").data("kendoChart").refresh();
                $("#chartTopTurbineLoss").data("kendoChart").refresh();
                app.loading(false);
            },300);
        }
    }
}