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
        turbine: [turbine],
        project: project
    };

    $("#modalDetail").on("shown.bs.modal", function () { 
        toolkit.ajaxPost(viewModel.appName + "monitoring/getdetailchart", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            monitoring.chartWindSpeed(res.data.Data.ws);
            monitoring.chartProduction(res.data.Data.prod);
            monitoring.dataAvailChart();
        });
        toolkit.ajaxPost(viewModel.appName + "monitoring/getevent", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
           monitoring.detailEvent(res.data.Data)
        });

    }).modal('show');
}

monitoring.chartWindSpeed = function(dataSource){
    $("#chartWindSpeed").kendoStockChart({
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
      chartArea:{
        height : 250,
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
        color: "#337ab7",
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
          color: "#337ab7",
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
        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      chartArea:{
        height : 250,
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
        color: "#ea5b19",
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
          color: "#ea5b19",
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