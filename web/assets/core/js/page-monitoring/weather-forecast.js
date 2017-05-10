'use strict';


viewModel.WeatherForecast = new Object();
var wf = viewModel.WeatherForecast;


vm.currentMenu('Weather');
vm.currentTitle('Weather Forecast');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Weather Forecast', href: viewModel.appName + 'page/monitoringweather' }]);
var intervalTurbine = null;
wf.projectList = ko.observableArray([]);
wf.project = ko.observable();

wf.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        wf.projectList([{ value: "", text: "" }]);
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
        wf.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(0);
            wf.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

wf.GetData = function() {
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    var param = { "q": "Jaisalmer,in", "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };
    $.ajax({
        type: "GET",
        url: surl,
        data: param,
        dataType: "jsonp",
        success:function(data){
          wf.ParseData(data);
        },
        error:function(){
            // do nothing
        }  
    });
};
wf.GetDataChart = function() {
    var surl = 'http://api.openweathermap.org/data/2.5/forecast';
    var param = { "q": "Jaisalmer,in", "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };
    $.ajax({
        type: "GET",
        url: surl,
        data: param,
        dataType: "jsonp",
        success:function(data){
          // console.log(data);
          wf.CreateChart(data);
        },
        error:function(){
            // do nothing
        }  
    });
};

wf.ParseData = function(data) {
    $('#sea_level').text(data.main.sea_level + ' m');
    $('#city').text(data.name + ' ('+ data.sys.country +')');
    $('#img_weather').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
    // $('#project_img_weather').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
    $('#weather').text(data.weather[0].description);
    // $('#project_weather').text(data.weather[0].description);
    $('#temperature').text(data.main.temp);
    // $('#project_temperature').text(data.main.temp);
    $('#wind_speed').text(data.wind.speed);
    $('#wind_direction').text(data.wind.deg);
    $('#pressure').text(data.main.pressure);
    $('#humidity').text(data.main.humidity);
    $('#last_update').text(moment().format('DD-MM-YYYY HH:mm:ss'));
};

wf.RefreshChart = function() {
    wf.GetDataChart();
    window.setInterval(wf.GetDataChart, 60000);
};
wf.CreateChart = function(data) {
    var datas = data.list;
    var chartData = [];
    if(datas.length > 0) {
        var countData = 0;
        $.each(datas, function(idx, val){
            var dt = moment(val.dt_txt, 'YYYY-MM-DD HH:mm:ss')
            if(idx < datas.length - 1) {
                var dtNext = moment(datas[idx+1].dt_txt, 'YYYY-MM-DD HH:mm:ss')
                if(dt <= moment() && dtNext >= moment()) {
                    var chartItem = {
                        date: dt.format("DD MMM YYYY HH:mm"),
                        hour: dt.get('hour'),
                        windspeed: val.wind.speed,
                        temperature: val.main.temp,
                    }; 
                    chartData.push(chartItem);
                    countData++;    
                }
            }
            if(dt >= moment() && countData < 12) {
                var chartItem = {
                    date: dt.format("DD MMM YYYY HH:mm"),
                    hour: dt.get('hour'),
                    windspeed: val.wind.speed,
                    temperature: val.main.temp,
                }; 
                chartData.push(chartItem);
                countData++;
            }
        });
    }
    $('#weatherChart').html('');
    $('#weatherChart').kendoChart({
        chartArea: { background:"transparent" },
        dataSource: {
            data: chartData
        },
        title: {
            visible: false
        },
        legend: {
            position: "bottom"
        },
        defaultSeries: {
            type: "line",
        },
        series: [{
            type: "line",
            name: "Temperature (&deg;C)",
            color: "#ff1c1c",
            axis: "temp",
            field: "temperature",
        }, {
            type: "line",
            name: "Wind Speed (m/s)",
            color: "#73c100",
            axis: "wind",
            field: "windspeed",
        }],
        valueAxes: [{
            name: "temp",
            color: "#707070",
            majorTicks: {
                visible: false,
            },
            title: {
                text: "Temperature",
                font: "11px"
            }
        },{
            name: "wind",
            color: "#73c100",
            majorTicks: {
                visible: false,
            },
            title: {
                text: "Wind Speed",
                font: "11px"
            }
        }],
        categoryAxis: {
            field: "hour",
            axisCrossingValues: [0, 32],
            justified: true,
            crosshair: {
                visible: false,
            },
            majorGridLines: {
                visible: false,
            },
            majorTicks: {
                visible: false,
            },
        },
        tooltip: {
            visible: true,
            format: "{0}",
            template: "#= dataItem.date # : #= value #"
        }
    });
};

$(document).ready(function(){
    app.loading(true);
    $.when(wf.GetData(), wf.RefreshChart()).done(function () {
        setTimeout(function() {
            app.loading(false);
        }, 2000);
    });
    window.setInterval(wf.GetData, 60000);
});
