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
monitoring.last_minute = ko.observable();
monitoring.last_date = ko.observable();
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

monitoring.showDetail = function(project, turbine){
    var param = {
        turbine: turbine,
        project: project
    };

    $("#modalDetail").on("shown.bs.modal", function () { 
        monitoring.windSpeedChart();
        monitoring.productionChart ();
        monitoring.dataAvailChart();
        setTimeout(function () {
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
        }, 300);
        monitoring.turbine = [turbine];
        monitoring.project = project;
        wr.GetData();
    }).modal('show');
}

monitoring.windSpeedChart = function(){
    var dataSource = [{
        "timestamp": "2016-11-18 00:00:00",
        "value": 1.125
      },
      {
        "timestamp": "2016-11-18 00:10:00",
        "value": 1.220833
      },
      {
        "timestamp": "2016-11-18 00:20:00",
        "value": 1.258333
      },
      {
        "timestamp": "2016-11-18 00:30:00",
        "value": 1.416667
      },
      {
        "timestamp": "2016-11-18 00:40:00",
        "value": 1.716667
      },
      {
        "timestamp": "2016-11-18 00:50:00",
        "value": 1.891667
      },
      {
        "timestamp": "2016-11-18 01:00:00",
        "value": 1.970833
      },
      {
        "timestamp": "2016-11-18 01:10:00",
        "value": 1.966667
      },
      {
        "timestamp": "2016-11-18 01:20:00",
        "value": 2.008333
      },
      {
        "timestamp": "2016-11-18 01:30:00",
        "value": 2.054167
      },
      {
        "timestamp": "2016-11-18 01:40:00",
        "value": 2.05
      },
      {
        "timestamp": "2016-11-18 01:50:00",
        "value": 1.9625
      },
      {
        "timestamp": "2016-11-18 02:00:00",
        "value": 1.833333
      },
      {
        "timestamp": "2016-11-18 02:10:00",
        "value": 1.633333
      },
      {
        "timestamp": "2016-11-18 02:20:00",
        "value": 1.420833
      },
      {
        "timestamp": "2016-11-18 02:30:00",
        "value": 1.275
      },
      {
        "timestamp": "2016-11-18 02:40:00",
        "value": 1.191667
      },
      {
        "timestamp": "2016-11-18 02:50:00",
        "value": 1.091667
      },
      {
        "timestamp": "2016-11-18 03:00:00",
        "value": 1.1
      },
      {
        "timestamp": "2016-11-18 03:10:00",
        "value": 1.079167
      },
      {
        "timestamp": "2016-11-18 03:20:00",
        "value": 0.958333
      }
    ]
    $("#wsChart").kendoChart({
      title: {
        text: "Wind Speed",
        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      theme: "flat",
      chartArea:{
        height: 175,
        margin:0, 
        padding: 0
      },
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      series: [{
        type: "area",
        field: "value",
        aggregate: "avg", 
        color: "#337ab7",
      }],
      valueAxis: {
        title: {
            text: "m/s",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
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
monitoring.productionChart = function(){
    var dataSource = [{
        "timestamp": "2016-11-18 00:00:00",
        "value": 1.125
      },
      {
        "timestamp": "2016-11-18 00:10:00",
        "value": 1.220833
      },
      {
        "timestamp": "2016-11-18 00:20:00",
        "value": 1.258333
      },
      {
        "timestamp": "2016-11-18 00:30:00",
        "value": 1.416667
      },
      {
        "timestamp": "2016-11-18 00:40:00",
        "value": 1.716667
      },
      {
        "timestamp": "2016-11-18 00:50:00",
        "value": 1.891667
      },
      {
        "timestamp": "2016-11-18 01:00:00",
        "value": 1.970833
      },
      {
        "timestamp": "2016-11-18 01:10:00",
        "value": 1.966667
      },
      {
        "timestamp": "2016-11-18 01:20:00",
        "value": 2.008333
      },
      {
        "timestamp": "2016-11-18 01:30:00",
        "value": 2.054167
      },
      {
        "timestamp": "2016-11-18 01:40:00",
        "value": 2.05
      },
      {
        "timestamp": "2016-11-18 01:50:00",
        "value": 1.9625
      },
      {
        "timestamp": "2016-11-18 02:00:00",
        "value": 1.833333
      },
      {
        "timestamp": "2016-11-18 02:10:00",
        "value": 1.633333
      },
      {
        "timestamp": "2016-11-18 02:20:00",
        "value": 1.420833
      },
      {
        "timestamp": "2016-11-18 02:30:00",
        "value": 1.275
      },
      {
        "timestamp": "2016-11-18 02:40:00",
        "value": 1.191667
      },
      {
        "timestamp": "2016-11-18 02:50:00",
        "value": 1.091667
      },
      {
        "timestamp": "2016-11-18 03:00:00",
        "value": 1.1
      },
      {
        "timestamp": "2016-11-18 03:10:00",
        "value": 1.079167
      },
      {
        "timestamp": "2016-11-18 03:20:00",
        "value": 0.958333
      }
    ]
    $("#prodChart").kendoChart({
      title: {
        text: "Production",
        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      theme: "flat",
      chartArea:{
        height: 175,
        margin:0, 
        padding: 0
      },
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      series: [{
        type: "area",
        field: "value",
        aggregate: "avg", 
        color: "#ea5b19",
      }],
      valueAxis: {
        title: {
            text: "MWh",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
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
monitoring.dataAvailChart = function(){
    var dataSource = [{
        "timestamp": "2016-11-18 00:00:00",
        "value": 1.125
      },
      {
        "timestamp": "2016-11-18 00:10:00",
        "value": 1.220833
      },
      {
        "timestamp": "2016-11-18 00:20:00",
        "value": 1.258333
      },
      {
        "timestamp": "2016-11-18 00:30:00",
        "value": 1.416667
      },
      {
        "timestamp": "2016-11-18 00:40:00",
        "value": 1.716667
      },
      {
        "timestamp": "2016-11-18 00:50:00",
        "value": 1.891667
      },
      {
        "timestamp": "2016-11-18 01:00:00",
        "value": 1.970833
      },
      {
        "timestamp": "2016-11-18 01:10:00",
        "value": 1.966667
      },
      {
        "timestamp": "2016-11-18 01:20:00",
        "value": 2.008333
      },
      {
        "timestamp": "2016-11-18 01:30:00",
        "value": 2.054167
      },
      {
        "timestamp": "2016-11-18 01:40:00",
        "value": 2.05
      },
      {
        "timestamp": "2016-11-18 01:50:00",
        "value": 1.9625
      },
      {
        "timestamp": "2016-11-18 02:00:00",
        "value": 1.833333
      },
      {
        "timestamp": "2016-11-18 02:10:00",
        "value": 1.633333
      },
      {
        "timestamp": "2016-11-18 02:20:00",
        "value": 1.420833
      },
      {
        "timestamp": "2016-11-18 02:30:00",
        "value": 1.275
      },
      {
        "timestamp": "2016-11-18 02:40:00",
        "value": 1.191667
      },
      {
        "timestamp": "2016-11-18 02:50:00",
        "value": 1.091667
      },
      {
        "timestamp": "2016-11-18 03:00:00",
        "value": 1.1
      },
      {
        "timestamp": "2016-11-18 03:10:00",
        "value": 1.079167
      },
      {
        "timestamp": "2016-11-18 03:20:00",
        "value": 0.958333
      }
    ]
    $("#dataAvailChart").kendoChart({
      title: {
        text: "Data Available",
        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      theme: "flat",
      chartArea:{
        height: 175,
        margin:0, 
        padding: 0
      },
      seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
      series: [{
        type: "area",
        field: "value",
        aggregate: "avg", 
        color: "#ea5b19",
      }],
      valueAxis: {
        title: {
            text: "MWh",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
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
function secondsToHms(d) {
    d = Number(d);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var res = (h > 0 ? (h < 10 ? "0" + h : h) : "00") + ":" + (m > 0 ? (m < 10 ? "0" + m : m) : "00") + ":" + (s > 0 ? s : "00")

    return res;
}

// ============================ WINDROSE ====================================


viewModel.WRFlexiDetail = new Object();
var wr = viewModel.WRFlexiDetail;

wr.dataWindrose = ko.observableArray([]);
wr.dataWindroseGrid = ko.observableArray([]);
wr.dataWindroseEachTurbine = ko.observableArray([]);
wr.sectorDerajat = ko.observable(0);

wr.sectionsBreakdownList = ko.observableArray([
    { "text": 36, "value": 36 },
    { "text": 24, "value": 24 },
    { "text": 12, "value": 12 },
]);
var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
var listOfChart = [];
var listOfButton = {};
var listOfCategory = [
    { "category": "0 to 4m/s", "color": colorFieldsWR[0] },
    { "category": "4 to 8m/s", "color": colorFieldsWR[1] },
    { "category": "8 to 12m/s", "color": colorFieldsWR[2] },
    { "category": "12 to 16m/s", "color": colorFieldsWR[3] },
    { "category": "16 to 20m/s", "color": colorFieldsWR[4] },
    { "category": "20m/s and above", "color": colorFieldsWR[5] },
];

wr.ExportWindRose = function () {
    var chart = $("#wr-chart").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "WindRose.pdf",
        });
    });
}
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
        // app.loading(true)
        setTimeout(function () {
            if ($(idChart).data("kendoChart") != null) {
                $(idChart).data("kendoChart").refresh();
            }
        }, 200);
    });
}

wr.RefreshChart = function(source) {
  setTimeout(function(){
      $.each(listOfChart, function(idx, elem){
          $(elem).data("kendoChart").refresh();
      });
  }, 300);
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
        $(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
        $(idChart).data("kendoChart").refresh();
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