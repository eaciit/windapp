'use strict';

viewModel.Forecasting = new Object();
var pg = viewModel.Forecasting;

vm.currentMenu('Forecasting & Scheduling');
vm.currentTitle('Forecasting & Scheduling');
vm.breadcrumb([{ title: 'Forecasting & Scheduling', href: viewModel.appName + 'page/forecasting' }]);

pg.DataSource = ko.observableArray([]);
pg.CurrentTab = ko.observable('grid');

pg.getData = function() {
    app.loading(true);
    var url = viewModel.appName + 'forecast/getlist';
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 
    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
    };
    var getdata = toolkit.ajaxPostDeffered(url, param, function(res) {});
    $.when(getdata).done(function(d){
        pg.DataSource(d.data);
        if(pg.CurrentTab()=='grid')
            pg.genereateGrid();
        else
            pg.genereateChart();
        app.loading(false);
    });
}

pg.genereateGrid = function(){
    app.loading(true);
    setTimeout(function(){ 
        $("#gridForecasting").html('');
        $("#gridForecasting").kendoGrid({
            dataSource: {
                data: pg.DataSource(),
                pageSize: 20
            },
            height: $('body').height() - 260,
            // scrollable: true,
            sortable: true,
            filterable: false,
            pageable: {
                input: true,
                numeric: false
            },
            columns: [
                { field: "Date", title: "Date"},
                { field: "TimeBlock", title: "Time Block" },
                { field: "AvaCap", title: "Ava Cap", template : "#: (AvaCap==null?'-':kendo.toString(AvaCap, 'n0')) #"},
                { field: "Forecast", template : "#: (Forecast==null?'-':kendo.toString(Forecast, 'n2')) #"},
                { title: "Sch Fcast <br>(SLDC)", field: "SchFcast", template : "#: (SchFcast==null?'-':kendo.toString(SchFcast, 'n2')) #" },
                { title: "Exp Prod <br> (Pwr Curv)", field: "ExpProd", template : "#: (ExpProd==null?'-':kendo.toString(ExpProd, 'n2')) #" },
                { field: "Actual", title: "Actual Prod", template : "#: (Actual==null?'-':kendo.toString(Actual, 'n2')) #" },
                { field: "FcastWs", title: "Fcast ws <br> (m/s)", template : "#: (FcastWs==null?'-':kendo.toString(FcastWs, 'n2')) #" },
                { field: "ActualWs", title: "Actual ws <br> (m/s)", template : "#: (ActualWs==null?'-':kendo.toString(ActualWs, 'n2')) #" },
                { field: "DevFcast", title: "% Dev B/W <br> Fcast & Act", template : "#: kendo.toString(DevFcast, 'p2') #" },
                { field: "DevSchAct", title: "% Dev B/W <br> Sch & Act", template : "#: kendo.toString(DevSchAct, 'p2') #" },
                { field: "DSMPenalty", title: "DSM Penalty"},
            ]
        });
        $("#gridForecasting").data("kendoGrid").refresh();
        app.loading(false);
    }, 300);
}

pg.genereateChart = function(){
    app.loading(true);
    var date1 = $('#dateStart').data('kendoDatePicker').value();
    var date2 = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));
    var timeDiff = Math.abs(date2.getTime() - date1.getTime());
    var diffDays = Math.ceil(timeDiff / (1000 * 3600 * 24)); 
    var mindays = 2;
    setTimeout(function(){
        $("#chartForecasting").html("");
        $("#chartForecasting").kendoChart({
            dataSource: {
                data: pg.DataSource(),
            },
            title: {
                text: "Forecasting and Scheduling",
                 font: '18px bold Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            chartArea : {
                height : $('body').height() - 260,
            },
            legend: {
                visible: true,
                position : "top",
            },
            seriesDefaults: {
                type: "line",
                labels: {
                    visible: false,
                    background: "transparent"
                },
                style: "smooth",
            },
            series: [{
                field: "Forecast",
                name: "Forecast",
                markers : {
                    visible : false
                },
                color: "#3d8dbd",
                dashType: "longDash",
                axis: "forecast",
            },{
                field: "SchFcast",
                name: "Sch Fcast (SLDC)",
                markers : {
                    visible : false
                },
                color: "#e91e63",
                dashType: "longDash",
                axis: "forecast",
            },{
                field: "ExpProd",
                name: "Exp Prod (Pwr Curv)",
                markers : {
                    visible : false
                },
                color: "#8bc34a",
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "Actual",
                name: "Actual Prod",
                markers : {
                    visible : false
                },
                color: "#9c27b0",
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "FcastWs",
                name: "Fcast ws (m/s)",
                markers : {
                    visible : false
                },
                color: "#00bcd4",
                dashType: "longDash",
                axis: "forecast",
            },{
                field: "ActualWs",
                name: "Actual ws (m/s)",
                markers : {
                    visible : false
                },
                color: "#ff9800",
                dashType: "solid",
                axis: "dynamic",
            }],
            valueAxes: [{
                line: {
                    visible: false
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n0') #",
                    visible: true,
                },
                name: "dynamic",
            },{
                line: {
                    visible: false
                }, 
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n0') #",
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "forecast",
            }],
            categoryAxis: {
                field: (diffDays>mindays?'Date':'TimeBlock'),
                axisCrossingValues: [0, 1000],
                majorGridLines: {
                    visible: false
                },
                majorTickType: "none",
                labels: {
                  rotation: (diffDays>mindays?45:'auto'),
                  step: (diffDays>mindays?96:4),
                }
            },
            tooltip: {
                visible: true,
                template: "${series.name} on #= moment(dataItem.TimeStamp).format('DD-MM-YYYY HH:mm') # = <b>#= kendo.toString(value, 'n2') #</b>"
            },
        });
        $("#chartForecasting").data("kendoChart").refresh();
        app.loading(false);
    },200);
}

pg.initLoad = function() {
    window.setTimeout(function(){
        fa.LoadData();
        di.getAvailDate();
        app.loading(false);

        pg.refresh();
    }, 200);
}

pg.refresh = function() {
    fa.checkTurbine();

    pg.getData();
}

$(function(){
    $('#projectList').kendoDropDownList({
        change: function () {  
            di.getAvailDate();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.project = project;
            fa.populateTurbine(project);
        }
    });
    $('#btnRefresh').on('click', function () {
        pg.refresh();
    });

    pg.initLoad();
})
