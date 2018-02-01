'use strict';

viewModel.Forecasting = new Object();
var pg = viewModel.Forecasting;


pg.DataDummy = ko.observableArray([
     {Date : "15/01/18", TimeBlock: "00:00 - 00.15" , AvaCap : 60, Forecast : 11.08, SchFcast : 11.08, ExpProd : 15, Actual : 8, FcastWs : 9.00 , ActualWs : 8, DevFcast : -0.7, DevSchAct: -0.7, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "00:15 - 00.30" , AvaCap : 60, Forecast : 10.67, SchFcast : 10.67, ExpProd : 10, Actual : 9, FcastWs : 8.67 , ActualWs : 9, DevFcast : 11.1, DevSchAct: 11.1, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "00:30 - 00.45" , AvaCap : 60, Forecast : 10.17, SchFcast : 10.17, ExpProd : 12, Actual : 9, FcastWs : 8.26 , ActualWs : 9, DevFcast : -1.75, DevSchAct: -1.75, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "00:45 - 01.00" , AvaCap : 60, Forecast : 9.68, SchFcast : 9.68, ExpProd : 10, Actual : 6, FcastWs : 7.86 , ActualWs : 6, DevFcast : -21.0, DevSchAct: -21.0, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "01:00 - 01.15" , AvaCap : 60, Forecast : 9.18, SchFcast : 9.18, ExpProd : 10, Actual : 6, FcastWs : 7.46 , ActualWs : 6, DevFcast : -53.0, DevSchAct: -53.0, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "01:15 - 01.30" , AvaCap : 60, Forecast : 8.68, SchFcast : 8.68, ExpProd : 9, Actual : 5, FcastWs :7.05 , ActualWs : 5, DevFcast : -73.6, DevSchAct: -73.6, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "01:30 - 01.45" , AvaCap : 60, Forecast : 8.21, SchFcast : 8.21, ExpProd : null, Actual : null, FcastWs : 6.67 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "01:45 - 02.00" , AvaCap : 60, Forecast : 7.74, SchFcast : 7.74, ExpProd : null, Actual : null, FcastWs : 6.29 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "02:00 - 02.15" , AvaCap : 60, Forecast : 7.27, SchFcast : 7.27, ExpProd : null, Actual : null, FcastWs : 5.91 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "02:15 - 02.30" , AvaCap : 60, Forecast : 6.79, SchFcast : 6.79, ExpProd : null, Actual : null, FcastWs : 5.52 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "02:30 - 02.45" , AvaCap : 60, Forecast : 6.48, SchFcast : 6.48, ExpProd : null, Actual : null, FcastWs : 5.26 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "02:45 - 03.00" , AvaCap : 60, Forecast : 6.17, SchFcast : 6.17, ExpProd : null, Actual : null, FcastWs : 5.01 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "03:00 - 03.15" , AvaCap : 60, Forecast : 5.86, SchFcast : 5.86, ExpProd : null, Actual : null, FcastWs : 4.76 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "03:15 - 03.30" , AvaCap : 60, Forecast : 5.55, SchFcast : 5.55, ExpProd : null, Actual : null, FcastWs : 4.51 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "03:30 - 03.45" , AvaCap : 60, Forecast : 5.35, SchFcast : 5.35, ExpProd : null, Actual : null, FcastWs : 4.35 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "03:45 - 04.00" , AvaCap : 60, Forecast : 5.14, SchFcast : 5.14, ExpProd : null, Actual : null, FcastWs : 4.18 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "04:00 - 04.15" , AvaCap : 60, Forecast : 4.94, SchFcast : 4.94, ExpProd : null, Actual : null, FcastWs : 4.01 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "04:15 - 04.30" , AvaCap : 60, Forecast : 4.74, SchFcast : 4.74, ExpProd : null, Actual : null, FcastWs : 3.85 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},
     {Date : "15/01/18", TimeBlock: "04:30 - 04.45" , AvaCap : 60, Forecast : 4.64, SchFcast : 4.64, ExpProd : null, Actual : null, FcastWs : 3.77 , ActualWs : null, DevFcast : null, DevSchAct: null, DSMPenalty: ""},

    ])

vm.currentMenu('Forecasting & Scheduling');
vm.currentTitle('Forecasting & Scheduling');
vm.breadcrumb([{ title: 'Forecasting & Scheduling', href: viewModel.appName + 'page/forecasting' }]);



pg.genereateGrid = function(){
    $("#gridForecasting").kendoGrid({
        dataSource: {
            data: pg.DataDummy(),
            pageSize: 20
        },
        height: 550,
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
            { field: "AvaCap", title: "Ava Cap"},
            { field: "Forecast"},
            { title: "Sch Fcast <br>(SLDC)", field: "SchFcast"},
            { title: "Exp Prod <br> (Pwr Curv)", field: "ExpProd"},
            { field: "Actual", title: "Actual Prod"},
            { field: "FcastWs", title: "Fcast ws <br> (m/s)"},
            { field: "ActualWs", title: "Actual ws <br> (m/s)"},
            { field: "DevFcast", title: "% Dev B/W <br> Fcast & Act", template : "#: kendo.toString(DevFcast, 'n0') # %" },
            { field: "DevSchAct", title: "% Dev B/W <br> Sch & Act", template : "#: kendo.toString(DevSchAct, 'n0') # %" },
            { field: "DSMPenalty", title: "DSM Penalty"},
        ]
    });
}

pg.genereateChart = function(){

    setTimeout(function(){
        $("#chartForecasting").html("");
        $("#chartForecasting").kendoChart({
            dataSource: {
                data: pg.DataDummy(),
            },
            title: {
                text: "Forecasting and Scheduling",
                 font: '18px bold Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            chartArea : {
                height : 500,
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
                    visible: true,
                },
                name: "dynamic",
            },{
                line: {
                    visible: false
                }, 
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n2') #",
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "forecast",
            }],
            categoryAxis: {
                field: "TimeBlock",
                axisCrossingValues: [0, 1000],
                majorGridLines: {
                    visible: false
                },
                majorTickType: "none",
                labels: {
                  rotation: "auto"    
                }
            }
        });

        $("#chartForecasting").data("kendoChart").refresh();
    },200);
}
$(function(){
    setTimeout(function(){
        app.loading(false);
        fa.LoadData();
        di.getAvailDate();
        pg.genereateGrid();
    },200);



    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();

    });

    $('#projectList').kendoDropDownList({
        change: function () {  
            di.getAvailDate();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
        }
    });
})
