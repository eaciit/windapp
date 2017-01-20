'use strict';


viewModel.TLPlots = new Object();
var tlp = viewModel.TLPlots;


vm.currentMenu('TrendLinePlots');
vm.currentTitle('TrendLinePlots');
vm.breadcrumb([ {
    title: 'Analysis Tool Box',
    href: '#'
}, {
    title: 'Trend Line Plots',
    href: viewModel.appName + 'page/analytictrendlineplots'
}]);

tlp.compTemp = ko.observableArray([
    { "value": 2, "text": "Ambient Temp" },
    { "value": 5, "text": "Temp_GearBox_IMS_NDE" },
    { "value": 5, "text": "Temp_GearBox_HSS_NDE" },
    { "value": 5, "text": "Temp_G1L1" },
    { "value": 5, "text": "Temp_G1L2" },
    { "value": 5, "text": "Temp_G1L3" },
    { "value": 5, "text": "Temp_GearBox_HSS_DE" },
    { "value": 5, "text": "Temp_GearOilSump" },
    { "value": 5, "text": "Temp_GeneratorBearing_DE" },
    { "value": 5, "text": "Temp_GeneratorBearing_NDE" },
    { "value": 5, "text": "Temp_MainBearing" },
    { "value": 5, "text": "Temp_GearBox_IMS_DE" },
    { "value": 5, "text": "Converter-1,2 temp" },
    { "value": 5, "text": "Nacelle Temp" },
]);

tlp.compTempVal = ko.observable("2");




tlp.initChart = function() {
    $("#charttlp").kendoChart({
        legend: {
            position: "bottom"
        },
        chartArea: {
            background: ""
        },
        seriesDefaults: {
            type: "line",
            style: "smooth"
        },
        series: [{
            name: "Turbine 1",
            color: colorField[1],
            data: [3907, 7943, 7848, 9284, 9263, 9801, 3890, 8238, 9552, 6855]
        }, {
            name: "Turbine 2",
            color: colorField[2],
            data: [1988, 2733, 3994, 3464, 4001, 3939, 1333, 2245, 4339, 2727]
        }, {
            name: "Turbine 3",
            color: colorField[3],
            data: [4743, 7295, 7175, 6376, 8153, 8535, 5247, 7832, 8832, 9832]
        }, {
            name: "Avg Value",
            color: "black",
            data: [1253, 2362, 3519, 4799, 3252, 4343, 5843, 6877, 7416, 8590]
        }],
        valueAxis: {
            labels: {
                format: "N1"
            },
            line: {
                visible: false
            },
        },
        categoryAxis: {
            categories: ["01-jan-2016", "01-jan-2016", "02-jan-2016", "03-jan-2016", "04-jan-2016", "05-jan-2016", "06-jan-2016", "07-jan-2016", "08-jan-2016", "09-jan-2016"],
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: 45
            }
        },
        tooltip: {
            visible: true,
            template: "#= series.name #: #= value #"
        }
    });
}

setTimeout(function() {
    app.loading(false);
    tlp.initChart()
}, 300);