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
bp.feeders = ko.observableArray();
bp.dataFeeders = ko.observableArray();
bp.newFeeders = ko.observableArray([]);

// var color = ["#4e6f90","#750c41","#009688","#1aa3a3","#de9c2b","#506642","#ee8d7d","#578897","#3f51b5","#5cbdaa"];
var color = ["#046293","#af1923","#66418c","#a8480c","#14717b","#4c792d","#880e4f","#9e7c21","#ac2258"]

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
    // app.loading(true);
    // bp.feeders([]);
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
        if(!app.isFine(res)) {
            app.loading(false);
            return;
        }
        var dataFeeder = [];
        var dataTurbine = {}; 

        var a = 0;
        var listAllTurbine = [];

        var feeders = [];


        bp.newFeeders([]);

        $.each(res.data.ListOfTurbine, function(i, val){
            var details = []
            $.each(val, function(idx, turbine){
                $.each(res.data.Detail, function(index, detail){
                    if(turbine == detail.Turbine){
                        details.push(detail);
                        detail.colorFeeder = color[a];
                        detail.feederName = i;
                        feeders.push(detail);
                    }
                });
            });
            dataFeeder.push({feederName : i, colorFeeder : color[a]});
            ++a;
        });

        var someArray = feeders;
        var groupSize = Math.floor(($(window).innerHeight() * (72 /100) - 24)/24);

        var groups = _.map(someArray, function(item, index){
          return index % groupSize === 0 ? someArray.slice(index, index + groupSize) : null; 
          })
          .filter(function(item){ return item; 
        });

        
        var iniData = [];
        
        $.each(groups, function(idx,val){
           var Detail = {details : val};
           iniData.push(Detail);
        });

        bp.newFeeders(iniData);

        bp.feeders(dataFeeder);
        bp.dataFeeders(listAllTurbine);

        $.when(bp.PlotData(res.data)).done(function(){
            setTimeout(function(){
                 app.loading(false);
            },200)
        })
    });
    bp.CheckWeather();
};

bp.PlotData = function(data) {
    var allData = data.Detail

    var lastUpdate = moment.utc(data.TimeStamp);
    $('#project_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));
    $('#project_turbine_active').text(24);
    $('#project_turbine_down').text(0);

    var totalPower = 0;
    var totalWs = 0;
    var countWs = 0;

    $.each(allData, function(idx, val){
        var turbine = val.Turbine;

        if(val.ActivePower > -999999) { 
            if(kendo.toString(val.ActivePower, 'n2')!=$('#power_'+ turbine).text()) {
                $('#power_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
             
            window.setTimeout(function(){ 
                $('#power_'+ turbine).css('background-color', 'transparent');
            }, 750);
        }
        if(val.WindSpeed > -999999) {
            if(kendo.toString(val.WindSpeed, 'n2')!=$('#wind_'+ turbine).text()) {
                $('#wind_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)'); 
            }
            
            window.setTimeout(function(){ 
                $('#wind_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.NacellePosition > -999999) {
            if(kendo.toString(val.NacellePosition, 'n2')!=$('#dir_'+ turbine).text()) {
                $('#dir_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
            }
            
            window.setTimeout(function(){ 
                $('#dir_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.RotorRPM > -999999) {
            if(kendo.toString(val.RotorRPM, 'n2')!=$('#rotor_'+ turbine).text()) {
                $('#rotor_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
            
            window.setTimeout(function(){ 
                $('#rotor_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.PitchAngle > -999999) {
            if(kendo.toString(val.PitchAngle, 'n2')!=$('#pitch_'+ turbine).text()) {
                $('#pitch_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
            }
            
            window.setTimeout(function(){ 
                $('#pitch_'+ turbine).css('background-color', 'transparent'); 
            }, 750);
        }
        if(val.Temperature > -999999) {
            if(kendo.toString(val.Temperature, 'n2')!=$('#temperature_'+ turbine).text()) {
                $('#temperature_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
            }
            
            window.setTimeout(function(){ 
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

    });


    $('#project_turbine_down').text(data.TurbineDown);
    $('#project_turbine_active').text(data.TurbineActive);

    window.setTimeout(function(){ 
        $('#project_generation').text(data.PowerGeneration.toFixed(2));
        $('#project_wind_speed').text(data.AvgWindSpeed.toFixed(2));
        $('#project_plf').text((data.PowerGeneration / 50400 * 100).toFixed(2));
    }, 1000);
    
};

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



$(document).ready(function() {
    app.loading(true);
    setTimeout(function() {
        bp.GetData()
        setInterval(bp.GetData, 4000);
    }, 600);
});