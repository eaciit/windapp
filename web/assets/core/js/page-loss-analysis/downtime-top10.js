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

dt.GenChartDownAlarmComponent = function (dataSource,id,Series,legend,name,axisLabel, vislabel,rotate,heightParam,wParam,format) {

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

dt.Downtime = function(){
    app.loading(true);
    fa.LoadData();
    if(pg.isFirstDowntime() === true){
        var param = {
            period: fa.period,
            dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
            dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
            turbine: fa.turbine,
            project: fa.project,
        }

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getdowntimetab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                var HDowntime = $('#filter-analytic').width() * 0.2
                var wAll = $('#filter-analytic').width() * 0.275

                /*Downtime Tab*/
                dt.GenChartDownAlarmComponent(res.data.duration,'chartDTDuration',SeriesDowntime,true,"Turbine", "Hours",false,-330,HDowntime,wAll,"N1");
                dt.GenChartDownAlarmComponent(res.data.frequency,'chartDTFrequency',SeriesDowntime,true,"Turbine", "Times",false,-330,HDowntime,wAll,"N0");
                dt.GenChartDownAlarmComponent(res.data.loss,'chartTopTurbineLoss',SeriesDowntime,true,"Turbine","MWh",false,-330,HDowntime,wAll,"N1");

                pg.isFirstDowntime(false);
                app.loading(false);
            },300);
           
        });
        $('#availabledatestart').html(pg.availabledatestartalarm2());
        $('#availabledateend').html(pg.availabledateendalarm2());
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartalarm2());
            $('#availabledateend').html(pg.availabledateendalarm2());
            $("#chartDTDuration").data("kendoChart").refresh();
            $("#chartDTFrequency").data("kendoChart").refresh();
            $("#chartTopTurbineLoss").data("kendoChart").refresh();
            app.loading(false);
        },300);
    }
}