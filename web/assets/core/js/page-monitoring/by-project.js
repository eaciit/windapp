'use strict';


viewModel.ByProject = new Object();
var bp = viewModel.ByProject;


vm.currentMenu('By Project');
vm.currentTitle('By Project');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'By Project', href: viewModel.appName + 'page/monitoringbyproject' }]);

bp.Turbines1 = ko.observableArray([]);
bp.Turbines2 = ko.observableArray([]);
bp.selectedTurbine = ko.observable("");
bp.selectedProject = ko.observable();
bp.selectedMonitoring = ko.observable();
bp.selectedMonitoring({
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
bp.detailEvent = ko.observableArray([]);
bp.projectList = ko.observableArray([]);
bp.project = ko.observable();
bp.turbine = ko.observableArray([]);

bp.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        bp.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val.Value;
                data.text = val.Name;
                datavalue.push(data);
            });
        }
        bp.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(0);
            bp.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

bp.PowerChartOpt = function(data) {
    return {
        type: "bullet",
        data: data,
        valueAxis: {
            min: 0,
            max: 2200,
            plotBands: [{
                from: 0, to: 800, color: "#787878", opacity: 0.15
            }, {
                from: 800, to: 1600, color: "#787878", opacity: 0.20
            }, {
                from: 1600, to: 2200, color: "#787878", opacity: 0.3
            }]
        },
        tooltip: {
            visible: true,
            template: 'Curr. Generation : #= value.current # KW<br />' +
                'Max. Capacity : #= value.target # KW'
        }
    };
};

bp.CheckWeather = function() {
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    var param = { "q": "Jaisalmer,in", "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };
    $.ajax({
        type: "GET",
        url: surl,
        data: param,
        dataType: "jsonp",
        success:function(data){
            $('#project_img_weather').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
            $('#project_weather').text(data.weather[0].description);
            $('#project_temperature').text(data.main.temp);
        },
        error:function(){
            // do nothing
        }  
    });
};

bp.GetData = function(data) {
    var COOKIES = {};
    var cookieStr = document.cookie;
    var project = "";

    if(cookieStr.indexOf("project=") >= 0) {
        document.cookie = "project=;expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        cookieStr.split(/; /).forEach(function(keyValuePair) {
            var cookieName = keyValuePair.replace(/=.*$/, "");
            var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
            COOKIES[cookieName] = cookieValue;
        });
        project = COOKIES["project"];
        $('#projectList').data('kendoDropDownList').value(project);
    } else {
        project = $('#projectList').data('kendoDropDownList').value();
    }

    var param = {
        Project: project
    };
    var getDetail = toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
        bp.Turbines1(res.data.ListOfTurbine["Feeder 5"]);
        bp.Turbines2(res.data.ListOfTurbine["Feeder 8"]);
        bp.PlotData(res.data);
    });
    bp.CheckWeather();

    app.loading(false);
};

bp.PlotData = function(data) {
    // var feeder5 = data.Data["Feeder 5"];
    // var allData = feeder5.concat(data.Data["Feeder 8"]);
    var allData = data.Detail
    var lastUpdate = moment.utc(data.TimeStamp);
    // var lastUpdate = moment.utc(data.TimeStamp).add(5.5, 'hour');
    // var lastUpdate = moment.utc().add(5.5, 'hour');
    // var energy = data.Data.ActivePower * 1000 * (3/3600);
    // var plf = energy / (2100 * (3/3600) * 24 * 1000) * 100;
    $('#project_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));
    // $('#turbine_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));
    //$('#project_generation').text(data.ActivePower.toFixed(2));
    //$('#project_wind_speed').text((data.WindSpeed>-999999?data.WindSpeed.toFixed(2):'N/A'));
    //$('#project_production').text(energy.toFixed(2));
    //$('#project_plf').text(plf.toFixed(2));
    $('#project_turbine_active').text(24);
    $('#project_turbine_down').text(0);

    var totalPower = 0;
    var totalWs = 0;
    var countWs = 0;

    $.each(allData, function(idx, val){
        var turbine = val.Turbine;

        if(val.ActivePower > -999999) { 
            if(val.ActivePower.toFixed(2)!=$('#power_'+ turbine).text()) {
                $('#power_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
             
            window.setTimeout(function(){ 
                if(val.ActivePower < 0) {
                    $('#power_'+ turbine).addClass('redvalue'); 
                } else {
                    $('#power_'+ turbine).removeClass('redvalue');  
                }
                $('#power_'+ turbine).text(val.ActivePower.toFixed(2));
                totalPower += parseFloat($('#power_'+ turbine).text());
                $('#power_'+ turbine).css('background-color', 'transparent');
            }, 750);
        }
        if(val.WindSpeed > -999999) {
            if(val.WindSpeed.toFixed(2)!=$('#wind_'+ turbine).text()) {
                $('#wind_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)'); 
            }
            
            window.setTimeout(function(){ 
                if(val.WindSpeed < 3.5) {
                    $('#wind_'+ turbine).addClass('orangevalue');   
                } else {
                    $('#wind_'+ turbine).removeClass('orangevalue');    
                }
                $('#wind_'+ turbine).text(val.WindSpeed.toFixed(2));
                totalWs += parseFloat($('#wind_'+ turbine).text());
                countWs++;
                $('#wind_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.NacellePosition > -999999) {
            if(val.NacellePosition.toFixed(2)!=$('#dir_'+ turbine).text()) {
                $('#dir_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
            }
            
            window.setTimeout(function(){ 
                $('#dir_'+ turbine).text(val.NacellePosition.toFixed(2));
                $('#dir_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.RotorRPM > -999999) {
            if(val.RotorRPM.toFixed(2)!=$('#rotor_'+ turbine).text()) {
                $('#rotor_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
            
            window.setTimeout(function(){ 
                $('#rotor_'+ turbine).text(val.RotorRPM.toFixed(2));
                $('#rotor_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.PitchAngle > -999999) {
            if(val.PitchAngle.toFixed(2)!=$('#pitch_'+ turbine).text()) {
                $('#pitch_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
            
            window.setTimeout(function(){ 
                $('#pitch_'+ turbine).text(val.PitchAngle.toFixed(2));
                $('#pitch_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.Temperature > -999999) {
            if(val.Temperature.toFixed(2)!=$('#temperature_'+ turbine).text()) {
                $('#temperature_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
            }
            
            window.setTimeout(function(){ 
                $('#temperature_'+ turbine).removeClass('orangevalue');
                $('#temperature_'+ turbine).removeClass('redvalue');
                if(val.Temperature > 38) {
                    $('#temperature_'+ turbine).addClass('redvalue');   
                } else if(val.Temperature >= 30) {
                    $('#temperature_'+ turbine).addClass('orangevalue');    
                }
                $('#temperature_'+ turbine).text(val.Temperature.toFixed(2));
                $('#temperature_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }

        /* TURBINE STATUS PART */
        if(val.AlarmDesc!="") {
            $('#alarmdesc_'+ turbine).text(val.AlarmCode);
            $('#alarmdesc_'+ turbine).attr('data-original-title', val.AlarmDesc);
        } else {
            $('#alarmdesc_'+ turbine).text("-");
            $('#alarmdesc_'+ turbine).attr('data-original-title', "This turbine already UP");
        }

        var iconStatus = "fa fa-circle fa-project-info fa-green";
        if(val.Status==0) {
            iconStatus = "fa fa-circle fa-project-info fa-red"; // faa-flash animated
        } else if(val.Status === 1 && val.IsWarning === true) {
            iconStatus = "fa fa-circle fa-project-info fa-orange";
        }
        if(val.DataComing==0) {
            iconStatus = "fa fa-circle fa-project-info fa-grey";
        }
        $('#status_'+ turbine).attr('class', iconStatus);

        /* END OF TURBINE STATUS PART */
    });

   
    /*$.each(data.TurbineStatus, function(idx, val){
        
    });*/

    $('#project_turbine_down').text(data.TurbineDown);
    $('#project_turbine_active').text(data.TurbineActive);

    window.setTimeout(function(){ 
        $('#project_generation').text(totalPower.toFixed(2));
        $('#project_wind_speed').text((totalWs / countWs).toFixed(2));
        $('#project_plf').text((totalPower / 50400 * 100).toFixed(2));
    }, 1000);
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
  if(bp.turbine.length > 0) {
    app.loading(true);
    setTimeout(function () {
        var breakDownVal = $("#nosection").data("kendoDropDownList").value();
        var secDer = 360 / breakDownVal;
        wr.sectorDerajat(secDer);
        var param = {
            turbine: bp.turbine,
            project: bp.project,
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
    var majorUnit = 10;
    if(maxValue < 40) {
        majorUnit = 5;
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
                    template: "#= (series.data[0] || {}).WsCategoryDesc #",
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
                majorUnit: majorUnit,
                // max: maxValue,
                // min: 0
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

var interval = null;

function secondsToHms(d) {
    d = Number(d);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var res = (h > 0 ? (h < 10 ? "0" + h : h) : "00") + ":" + (m > 0 ? (m < 10 ? "0" + m : m) : "00") + ":" + (s > 0 ? s : "00")

    return res;
}

bp.ToIndividualTurbine = function(turbine) {
    setTimeout(function(){
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        document.cookie = "project="+bp.project.split("(")[0].trim()+";expires="+ newDateObj;
        document.cookie = "turbine="+turbine+";expires="+ newDateObj;
        window.location = viewModel.appName + "page/monitoringbyturbine";
    },300);
}

bp.ToAlarm = function(turbine) {


    var set = setTimeout(function(){
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        document.cookie = "project="+bp.project.split("(")[0].trim()+";expires="+ newDateObj;
        document.cookie = "turbine="+turbine+";expires="+ newDateObj;
    },300);

    console.log(document.cookie);

    $.when(set).done(function(){
         window.location = viewModel.appName + "page/monitoringalarm";
     });
}
bp.GenDetail = function(turbine){
    var param = {
        turbine: [turbine],
        project: bp.project
    };

    bp.selectedProject(bp.project);
    bp.selectedTurbine(turbine);    
    $("#modalDetail").on("shown.bs.modal", function () { 
        var param = {
            turbine: [bp.selectedTurbine()],
            project: bp.selectedProject()
        };
        var getDetail = toolkit.ajaxPost(viewModel.appName + "monitoring/getdetailchart", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            bp.chartWindSpeed(res.data.Data.ws);
            bp.chartProduction(res.data.Data.prod);
            bp.dataAvailChart(res.data.Data.avail);
            bp.dataChartLine(res.data.Data.line);
            bp.selectedMonitoring(res.data.Data.monitoring);
            if(bp.selectedMonitoring().winddirection <= -999999) {

            }
            if(bp.selectedMonitoring().pitchangle <= -999999) {

            }
        });
        var getEvent = toolkit.ajaxPost(viewModel.appName + "monitoring/getevent", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
           bp.detailEvent(res.data.Data)
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
        bp.turbine = [turbine];
        wr.GetData();
        bp.changeRotation();

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
};

bp.changeRotation = function(){
    $.each( $('.rotation'), function( key, value ) {
        var deg = $(value).attr("rotationval")
        $(value).attr("style", $(value).attr("style")+"-ms-transform: rotate("+deg+"deg);-webkit-transform: rotate("+deg+"deg);transform: rotate("+deg+"deg);");
    });
}

bp.chartWindSpeed = function(dataSource){
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

bp.chartProduction = function(dataSource){
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
bp.dataAvailChart = function(dataSource){
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
bp.dataChartLine = function (data) {
    $("#chartline").html("");
    $("#chartline").kendoChart({
        zoomable: true,
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
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
        },
         {
            type: "area",
            // style: "smooth",
            field: "avail",
            axis: "percentage",
            name: "Availability(%)",
            markers: {
                visible: false,
            },
            width: 3,
        }
        ],
        seriesColors: colorFields2,
        valueAxes: [{
            line: {
                visible: false
            },
            max: 100,
            // majorUnit: 20,
            labels: {
                format: "{0}%",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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

$(document).ready(function() {
    app.loading(true);
    // $("#monitoring").gridalicious({width: 225});
    setTimeout(function() {
        bp.GetData()
        window.setInterval(bp.GetData, 4000);
    }, 600);
});