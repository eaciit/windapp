'use strict';

viewModel.Downtime = new Object();
var dt = viewModel.Downtime;
var SeriesDowntime = [{
    field: "AEBOK",
    name: "AEBOK"
}, {
    field: "ExternalStop",
    name: "External Stop"
}, {
    field: "GridDown",
    name: "Grid Down"
}, {
    field: "InternalGrid",
    name: "InternalGrid"
}, {
    field: "MachineDown",
    name: "Machine Down"
}, {
    field: "WeatherStop",
    name: "Weather Stop"
}, {
    field: "Unknown",
    name: "Unknown"
}]

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
                pg.GenChartDownAlarmComponent(res.data.duration,'chartDTDuration',SeriesDowntime,true,"Turbine", "Hours",false,-330,HDowntime,wAll,"N1");
                pg.GenChartDownAlarmComponent(res.data.frequency,'chartDTFrequency',SeriesDowntime,true,"Turbine", "Times",false,-330,HDowntime,wAll,"N0");
                pg.GenChartDownAlarmComponent(res.data.loss,'chartTopTurbineLoss',SeriesDowntime,true,"Turbine","MWh",false,-330,HDowntime,wAll,"N1");

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