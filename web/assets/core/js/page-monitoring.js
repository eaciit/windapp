'use strict';

viewModel.monitoring = {};
var monitoring = viewModel.monitoring;


vm.currentMenu('Monitoring');
vm.currentTitle('Monitoring');
vm.breadcrumb([{
    title: 'Monitoring',
    href: viewModel.appName + 'page/monitoring'
}, {
    title: 'Monitoring',
    href: '#'
}]);

monitoring.turbineList = ko.observableArray([]);
monitoring.projectList = ko.observableArray([]);
monitoring.turbine = ko.observableArray([]);
monitoring.project = ko.observable();
monitoring.data = ko.observableArray([]);
monitoring.event = ko.observableArray([]);
monitoring.detailEvent = ko.observableArray([]);
monitoring.last_minute = ko.observable();
monitoring.last_date = ko.observable();
monitoring.selectedProject = ko.observable();
monitoring.selectedTurbine = ko.observable();
monitoring.selectedMonitoring = ko.observable();
monitoring.selectedMonitoring({
    pitchangle: 0,
    production: 0,
    project:"",
    rotorspeedrpm:0,
    status:"",
    statuscode:"",
    statusdesc:"",
    timestamp:"",
    timestampstr:"",
    totalProduction:0,
    turbine:"",
    winddirection:0,
    windspeed:0,
});
var turbineval = [];


monitoring.populateTurbine = function (data) {
    if (data.length == 0) {
        data = [];
        monitoring.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            var allturbine = {}
            $.each(data, function (key, val) {
                turbineval.push(val);
            });
            allturbine.value = "All Turbine";
            allturbine.text = "All Turbines";
            datavalue.push(allturbine);
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        monitoring.turbineList(datavalue);
    }

    setTimeout(function () {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"]);
    }, 300);
};

monitoring.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        monitoring.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        monitoring.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(1);
            monitoring.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

monitoring.checkTurbine = function () {
    var arr = $('#turbineList').data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbine");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#turbineList').data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    } else if (arr.length == 0) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    }
}

monitoring.getData = function(){
    // app.loading(true);

    var turbine = $("#turbineList").data("kendoMultiSelect").value()
    var project = $("#projectList").data("kendoDropDownList").value()
    var param = {
        turbine: (turbine == "All Turbine" ? []: turbine),
        project: project
    };

    var request = toolkit.ajaxPost(viewModel.appName + "monitoring/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        monitoring.last_minute(res.data.Data.timestamp.minute);
        monitoring.last_date(res.data.Data.timestamp.date);

        monitoring.data([]);
        $.each(res.data.Data.data, function (index, item) {   
            monitoring.data.push(item);                    
        });
    });

    var requestEvent = toolkit.ajaxPost(viewModel.appName + "monitoring/getevent", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
       monitoring.event(res.data.Data)
    });


    $.when(request, requestEvent).done(function(){
        $(".red").addClass("flash");
        setTimeout(function(){
            // app.loading(false);
            app.prepareTooltipster();            
        },500);

        setTimeout(function() {
            $(".red").removeClass("flash");
        }, 2500);
    });
}

var interval = null;

monitoring.showDetail = function(project, turbine){
    var param = {
        turbine: [turbine],
        project: project
    };

    monitoring.selectedProject(project);
    monitoring.selectedTurbine(turbine);    
    $("#modalDetail").on("shown.bs.modal", function () { 
        var param = {
            turbine: [monitoring.selectedTurbine()],
            project: monitoring.selectedProject()
        };
        var getDetail = toolkit.ajaxPost(viewModel.appName + "monitoring/getdetailchart", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            monitoring.chartWindSpeed(res.data.Data.ws);
            monitoring.chartProduction(res.data.Data.prod);
            monitoring.dataAvailChart(res.data.Data.avail);
            monitoring.dataChartLine(res.data.Data.line);
            monitoring.selectedMonitoring(res.data.Data.monitoring);
        });
        var getEvent = toolkit.ajaxPost(viewModel.appName + "monitoring/getevent", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
           monitoring.detailEvent(res.data.Data)
        });

         /*WINDROSE INITIAL*/
        $("#legend-list").html("");
        $.each(listOfCategory, function (idx, val) {
            var idName = "btn" + idx;
            listOfButton[idName] = true;
            $("#legend-list").append(
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="wr.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
                '<span class="span-legend">' + val.category + '</span>'
            );
        });
        $("#nosection").data("kendoDropDownList").value(12);
        monitoring.turbine = [turbine];
        monitoring.project = project;
        wr.GetData();
        monitoring.changeRotation();

        interval = setInterval(function(){
            getDetail 
            getEvent
            wr.GetData();
         },1000*120);

    }).modal('show');

    $('#modalDetail').on('hidden.bs.modal', function (e) {
        clearInterval(interval);
        $('#modalDetail').off();
    });
}



monitoring.chartWindSpeed = function(dataSource){
    $("#chartWindSpeed").kendoStockChart({
      title: {
        text: "Wind Speed",
        font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      chartArea:{
        height : 220,
      },
      theme: "flat",
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      dateField: "timestamp",
      series: [{
        type: "area",
        field: "value",
        aggregate: "avg", 
        color: "#ea5b19",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        pane: {
            height: 50,
        },
        series: [{
          type: "area",
          field: "value",
          aggregate: "avg",
          color: "#ea5b19",
        }]
      },
      valueAxis: {
        title: {
            text: "m/s",
            visible: true,
            font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
      categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # m/s",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
} 
monitoring.chartProduction = function(dataSource){
    $("#chartProduction").kendoStockChart({
      title: {
        text: "Production",
        font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      chartArea:{
        height : 220,
      },
      theme: "flat",
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      dateField: "timestamp",
      series: [{
        type: "area",
        field: "value",
        aggregate: "sum", 
        color: "#ee7a44",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        pane:{
            height: 50,
        },
        series: [{
          type: "area",
          field: "value",
          aggregate: "sum",
          color: "#ee7a44",
        }]
      },
      valueAxis: {
        title: {
            text: "MWh",
            visible: true,
            font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
      categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # MWh",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
} 
monitoring.dataAvailChart = function(dataSource){
    $("#dataAvailChart").kendoStockChart({
      title: {
        text: "Availability",
        font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      chartArea:{
        height : 220,
      },
      theme: "flat",
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      dateField: "timestamp",
      series: [{
        type: "area",
        field: "value",
        // aggregate: "sum", 
        color: "#f4ac8a",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        pane:{
            height: 50,
        },
        series: [{
          type: "area",
          field: "value",
          // aggregate: "sum",
          color: "#f4ac8a",
        }]
      },
      valueAxis: {
        title: {
            text: "%",
            visible: true,
            font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0,
        max: 100,
      },
      categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # %",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
}
function secondsToHms(d) {
    d = Number(d);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var res = (h > 0 ? (h < 10 ? "0" + h : h) : "00") + ":" + (m > 0 ? (m < 10 ? "0" + m : m) : "00") + ":" + (s > 0 ? s : "00")

    return res;
}

// Minimum/maximum number of visible items
var MIN_SIZE = 40;
var MAX_SIZE = 100;

// Optional sort expression
// var SORT = { field: "val", dir: "asc" };
var SORT = {};

// Minimum distance in px to start dragging
var DRAG_THR = 50;

// State variables
var viewStart = 0;
var viewSize = MIN_SIZE;
var newStart;

// Drag handler
function onDrag(e) {
    var chart = e.sender;
    var ds = chart.dataSource;
    var delta = Math.round(e.originalEvent.x.initialDelta / DRAG_THR);

    if (delta != 0) {
    newStart = Math.max(0, viewStart - delta);
    newStart = Math.min(data.length - viewSize, newStart);
    ds.query({
        skip: newStart,
        page: 0,
        pageSize: viewSize,
        sort: SORT
    });
    }
}

function onDragEnd() {
    viewStart = newStart;
}

monitoring.createLineChartZoom = function(e) {
    var chart = e.sender;
    var ds = chart.dataSource;
    viewSize = Math.min(Math.max(viewSize + e.delta, MIN_SIZE), MAX_SIZE);
    ds.query({
        skip: viewStart,
        page: 0,
        pageSize: viewSize,
        sort: SORT
    });

    // Prevent document scrolling
    e.originalEvent.preventDefault();
}

monitoring.dataChartLine = function (data) {
    $("#chartline").html("");
    $("#chartline").kendoChart({
        zoomable: true,
        /*zoom: monitoring.createLineChartZoom,
        transitions: false,
        drag: onDrag,
        dragEnd: onDragEnd,*/
        dataSource: {
            data: data,
            sort: { field: "Timestamp", dir: 'asc' }
        },
        theme: "Flat",
        chartArea: {
            height: 414,
        },
        legend: {
            position: "top",
            visible: true,
        },
        series: [{
            type: "line",
            style: "smooth",
            field: "ws",
            axis: "ws",
            name: "Wind Speed(m/s)",
            markers: {
                visible: false,
            },
            width: 3,
        }, {
            type: "line",
            style: "smooth",
            field: "production",
            axis: "prod",
            name: "Production(KWh)",
            markers: {
                visible: false,
            },
            width: 3,
        }, {
            type: "area",
            // style: "smooth",
            field: "avail",
            axis: "percentage",
            name: "Availability(%)",
            markers: {
                visible: false,
            },
            width: 3,
        }],
        seriesColors: colorFields2,
        valueAxes: [{
            line: {
                visible: false
            },
            max: 100,
            // majorUnit: 20,
            labels: {
                format: "{0}%",
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "percentage",
            title: { text: "Availability (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }, {
            line: {
                visible: false
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            max: 25,
            // majorUnit: 1,
            labels: {
                format: "{0}(m/s)",
            },
            name: "ws",
            title: { text: "Average Wind Speed (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }, {
            line: {
                visible: false
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            // max: 25,
            // majorUnit: 20,
            labels: {
                format: "{0}(KWh)",
            },
            name: "prod",
            title: { text: "Production (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            field: "timestamp",
            title: {
                text: "Time",
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            axisCrossingValues: [0, 1000],
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            shared: true,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            // template: "#= series.name # : #= kendo.toString(value, 'n2')# at #= category #",
            template: "#= kendo.toString(value, 'n2')#",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });

    $("#chartline").data("kendoChart").refresh();
};

// ============================ WINDROSE ====================================


viewModel.WRFlexiDetail = new Object();
var wr = viewModel.WRFlexiDetail;

wr.dataWindroseEachTurbine = ko.observableArray([]);
wr.sectorDerajat = ko.observable(0);

wr.sectionsBreakdownList = ko.observableArray([
    { "text": 36, "value": 36 },
    { "text": 24, "value": 24 },
    { "text": 12, "value": 12 },
]);
var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
// var colorFieldsWR = ["#2d6a9f", "#337ab7", "#4c91cd", "#74a9d8", "#9cc2e3", "#c3daee", "#ebf3f9"];
var listOfChart = [];
var listOfButton = {};
var listOfCategory = [
    { "category": "0 to 4 m/s", "color": colorFieldsWR[0] },
    { "category": "4 to 8 m/s", "color": colorFieldsWR[1] },
    { "category": "8 to 12 m/s", "color": colorFieldsWR[2] },
    { "category": "12 to 16 m/s", "color": colorFieldsWR[3] },
    { "category": "16 to 20 m/s", "color": colorFieldsWR[4] },
    { "category": "20 m/s and above", "color": colorFieldsWR[5] },
];

var maxValue = 0;

wr.GetData = function () {
  if(monitoring.turbine.length > 0) {
    app.loading(true);
    setTimeout(function () {
        var breakDownVal = $("#nosection").data("kendoDropDownList").value();
        var secDer = 360 / breakDownVal;
        wr.sectorDerajat(secDer);
        var param = {
            turbine: monitoring.turbine,
            project: monitoring.project,
            breakDown: breakDownVal,
            isMonitoring: true,
        };
        toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getflexidataeachturbine", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            if (res.data.WindRose != null) {
                var metData = res.data.WindRose;
                maxValue = res.data.MaxValue;
                wr.dataWindroseEachTurbine(metData);
                wr.initChart();
            }

            app.loading(false);

        })
    }, 300);
  }
}

wr.initChart = function () {
    // app.loading(true)
    listOfChart = [];
    var breakDownVal = $("#nosection").data("kendoDropDownList").value();
    var stepNum = 1
    var gapNum = 1
    if (breakDownVal == 36) {
        stepNum = 3
        gapNum = 0
    } else if (breakDownVal == 24) {
        stepNum = 2
        gapNum = 0
    } else if (breakDownVal == 12) {
        stepNum = 1
        gapNum = 0
    }

    $.each(wr.dataWindroseEachTurbine(), function (i, val) {
        var name = val.Name
        var idChart = "#chart-" + val.Name
        listOfChart.push(idChart);
        $(idChart).kendoChart({
            chartArea: {
                height: 350
            },
            theme: "nova",
            title: {
                text: name,
                visible: false,
                font: '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            legend: {
                position: "bottom",
                labels: {
                    template: "#= (series.data[0] || {}).WsCategoryDesc #"
                },
                visible: false,
            },
            dataSource: {
                data: val.Data,
                group: {
                    field: "WsCategoryNo",
                    dir: "asc"
                },
                sort: {
                    field: "DirectionNo",
                    dir: "asc"
                }
            },
            seriesColors: colorFieldsWR,
            series: [{
                type: "radarColumn",
                stack: true,
                field: "Contribution",
                gap: gapNum,
                border: {
                    width: 1,
                    color: "#7f7f7f",
                    opacity: 0.5
                },
            }],
            categoryAxis: {
                field: "DirectionDesc",
                visible: true,
                majorGridLines: {
                    visible: true,
                    step: stepNum
                },
                labels: {
                    font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    visible: true,
                    step: stepNum
                }
            },
            valueAxis: {
                labels: {
                    template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                    font: '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                majorUnit: 10,
                max: maxValue,
                min: 0
            },
            tooltip: {
                visible: true,
                template: "#= category # (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });
    });
}

wr.showHideLegend = function (index) {
    var idName = "btn" + index;
    listOfButton[idName] = !listOfButton[idName];
    if (listOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChart, function (idx, idChart) {
        if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
            $(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
            $(idChart).data("kendoChart").refresh();
        }
    });
}

monitoring.changeRotation = function(){
    $.each( $('.rotation'), function( key, value ) {
        var deg = $(value).attr("rotationval")
        $(value).attr("style", $(value).attr("style")+"-ms-transform: rotate("+deg+"deg);-webkit-transform: rotate("+deg+"deg);transform: rotate("+deg+"deg);");
    });
}

$(function () {

    $("#restore-screen").hide();

    $("#max-screen").click(function(){
        $("html").addClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 80);
        $(".multicol").height($(window).innerHeight() - 80 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 80-50);
        $("#max-screen").hide();
        $("#restore-screen").show();  
    });

    $("#restore-screen").click(function(){
        $("html").removeClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150-50);
        $("#max-screen").show();  
        $("#restore-screen").hide();  
    });

    $('#btnRefresh').on('click', function() {
        monitoring.getData();
    });

    setInterval(function(){monitoring.getData()},1000*120);

    setTimeout(function() {
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150 - 50);
        monitoring.getData();
    }, 500);
});