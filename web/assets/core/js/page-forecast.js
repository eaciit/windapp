'use strict';

viewModel.Forecasting = new Object();
var pg = viewModel.Forecasting;
var heightSub = 200;

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
                pageSize: 15
            },
            height: 520, //$('body').height() - heightSub + 30,
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
                { field: "AvaCap", title: "Avg. Capacity<br />(MW)", template : "#: (AvaCap==null?'-':kendo.toString(AvaCap, 'n0')) #"},
                { field: "Forecast", title: "Forecast<br />(MW)", template : "#: (Forecast==null?'-':kendo.toString(Forecast, 'n2')) #"},
                { title: "Sch Fcast /<br>SLDC (MW)", field: "SchFcast", template : "#: (SchFcast==null?'-':kendo.toString(SchFcast, 'n2')) #" },
                { field: "Actual", title: "Actual Prod<br/>(MW)", template : "#: (Actual==null?'-':kendo.toString(Actual, 'n2')) #" },
                { title: "Exp Prod /<br>Pwr Curv (MW)", width: 120, field: "ExpProd", template : "#: (ExpProd==null?'-':kendo.toString(ExpProd, 'n2')) #" },
                { field: "FcastWs", title: "Fcast ws<br>(m/s)", template : "#: (FcastWs==null?'-':kendo.toString(FcastWs, 'n2')) #" },
                { field: "ActualWs", title: "Actual ws<br>(m/s)", template : "#: (ActualWs==null?'-':kendo.toString(ActualWs, 'n2')) #" },
                { field: "DevFcast", title: "% Dev B/W<br>Fcast & Act", template : "#: kendo.toString(DevFcast, 'p2') #" },
                { field: "DevSchAct", title: "% Dev B/W<br>Sch & Act", template : "#: kendo.toString(DevSchAct, 'p2') #" },
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
                margin: {
                    top: 0,
                }
            },
            chartArea : {
                height : 500, // $('body').height() - heightSub + 30,
                background: "transparent",
            },
            legend: {
                visible: true,
                position : "top",
                offsetY: -10,
                item: {
                    visual: function (e) {
                        var color = e.options.markers.background;
                        var dashType = e.series.dashType;
                        var labelColor = e.options.labels.color;
                        var rect = new kendo.geometry.Rect([0, 0], [200, 50]);
                        var layout = new kendo.drawing.Layout(rect, {
                          spacing: 5,
                          alignItems: "center"
                        });
              
                        var svgPath = "M0 5.5 L 28 5.5 28 8 0 8Z";
                        if(dashType!="solid") {
                            svgPath = "M0 5.5 L 8 5.5 8 8 0 8Z M11 5.5 L 19 5.5 19 8 11 8Z M22 5.5 L 30 5.5 30 8 22 8Z";
                        }
                        var path = kendo.drawing.Path.parse(svgPath, {
                            stroke: {
                                color: color,
                                width: 0
                            },
                            fill: {
                                color: color,
                            },
                            cursor: "pointer",
                        });
              
                        var label = new kendo.drawing.Text(e.series.name, [0, 0], {
                          fill: {
                            color: '#232323',
                          },
                          cursor: "pointer",
                          font: "12px 'Source Sans Pro',Lato,'Open Sans','Helvetica Neue',Arial,sans-serif!important",
                          opacity: 0.8,
                        });
              
                        layout.append(path, label);
                        layout.reflow()
              
                        return layout;
                    }
                }
            },
            seriesDefaults: {
                type: "line",
                labels: {
                    visible: false,
                    background: "transparent"
                },
                style: "smooth",
            },
            axisDefaults: {
                crosshair: {
                    visible: true,
                    opacity: 0.175,
                    width: 0.7,
                },
                majorTicks: {
                    visible: false,
                    step: 3,
                    width: 1,
                    size: 2,
                },
            },
            series: [{
                field: "Forecast",
                name: "Forecast",
                markers : {
                    visible : false
                },
                color: "#9c27b0",
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
                field: "Actual",
                name: "Actual Prod",
                markers : {
                    visible : false
                },
                color: "#3d8dbd",
                dashType: "solid",
                axis: "dynamic",
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
                    template: "#= kendo.toString(value, 'n1') #",
                    visible: true,
                },
                name: "dynamic",
                axisCrossingValue: [-10],
            },{
                line: {
                    visible: false
                }, 
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n1') #",
                    margin: {
                        left: 10,
                    },
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
                axisCrossingValues: [0, 10000],
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
                template: "${series.name} on #= moment.utc(dataItem.TimeStamp).format('DD-MM-YYYY HH:mm') # = <b>#= kendo.toString(value, 'n2') #</b>"
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
