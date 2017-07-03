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
bp.oldFeeders = ko.observableArray([]);
bp.fullscreen = ko.observable(false);
bp.currentTempLocation = ko.observable();
// var color = ["#4e6f90","#750c41","#009688","#1aa3a3","#de9c2b","#506642","#ee8d7d","#578897","#3f51b5","#5cbdaa"];
var color = ["#046293","#af1923","#66418c","#a8480c","#14717b","#4c792d","#880e4f","#9e7c21","#ac2258"]

bp.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        bp.projectList([{ value: "", text: "" , city: ""}]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val.Value;
                data.text = val.Name;
                data.city = val.City;
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
    var city = bp.projectList()[$('#projectList').data('kendoDropDownList').select()].city;
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    $("#citytxt").html(city);
    var param = { "q": city, "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };
    $.ajax({
        type: "GET",
        url: surl,
        data: param,
        dataType: "jsonp",
        success:function(data){
            $('#project_img_weather').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
            $('#project_weather').text(data.weather[0].description);
            $('#project_temperature').text(data.main.temp);
            bp.currentTempLocation(data.main.temp);
        },
        error:function(){
            // do nothing
        }  
    });
};

bp.GetDataProject = function(project) {
    var param = {
        Project: project,
        LocationTemp: parseFloat($('#project_temperature').text())
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
                var feedHeader = {}
                if(idx == 0){
                    feedHeader.ActivePower = null
                    feedHeader.ActivePowerColor = null
                    feedHeader.AlarmCode = null
                    feedHeader.AlarmDesc = null
                    feedHeader.AlarmUpdate = null
                    feedHeader.DataComing = null
                    feedHeader.IconStatus = null
                    feedHeader.IsWarning = null
                    feedHeader.NacellePosition = null
                    feedHeader.PitchAngle = null
                    feedHeader.RotorRPM = null
                    feedHeader.Status = null
                    feedHeader.Temperature = null 
                    feedHeader.TemperatureColor = null
                    feedHeader.Turbine = null
                    feedHeader.Name = null
                    feedHeader.WindDirection = null
                    feedHeader.WindSpeed = null
                    feedHeader.WindSpeedColor = null
                    feedHeader.isHeader = true;
                    feedHeader.feederName = i;
                    feedHeader.index = idx;
                    feeders.push(feedHeader);
                }
                $.each(res.data.Detail, function(index, detail){
                    if(turbine == detail.Turbine){
                        details.push(detail);
                        detail.colorFeeder = color[a];
                        detail.feederName = i;
                        detail.index = idx;
                        detail.isHeader = false;
                        feeders.push(detail);
                    }
                });
            });
            dataFeeder.push({feederName : i, colorFeeder : color[a]});
            ++a;
        });

        var someArray = feeders;
        var groupSize = (bp.fullscreen() == true ? Math.floor(($(window).innerHeight() - 105 - 24)/24) : Math.floor(($(window).innerHeight() - 175 - 24)/24));

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
            },200);
        });
    });
}

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
    
    $.when(bp.CheckWeather()).done(function(){
        setTimeout(function(){
             bp.GetDataProject(project);
        }, 200);
    });
};

bp.PlotData = function(data) {
    var allData = data.Detail
    var oldData = (bp.oldFeeders().length == 0 ? allData : bp.oldFeeders());
    var lastUpdate = moment.utc(data.TimeMax);

    $('#project_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));
    $('#project_turbine_active').text(24);
    $('#project_turbine_down').text(0);

    var totalPower = 0;
    var totalWs = 0;
    var countWs = 0;

    $.each(allData, function(idx, val){
        var turbine = val.Turbine;
        var oldTurbine = oldData[idx].Turbine;
        var oldVal = oldData[idx];
        if(oldTurbine == turbine){
            if(val.ActivePower > -999999) { 
                if(kendo.toString(val.ActivePower, 'n2')!= kendo.toString(oldVal.ActivePower, 'n2')) {
                    $('#power_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
                }
                $('#powerdata_'+ turbine).text(val.ActivePower.toFixed(2));
                window.setTimeout(function(){ 
                    $('#power_'+ turbine).css('background-color', 'transparent');
                }, 750);
                
            }
            if(val.WindSpeed > -999999) {
                if(kendo.toString(val.WindSpeed, 'n2')!= kendo.toString(oldVal.WindSpeed, 'n2')) {
                    $('#wind_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)'); 
                }
                
                window.setTimeout(function(){ 
                    $('#wind_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }
            if(val.NacellePosition > -999999) {
                if(kendo.toString(val.NacellePosition, 'n2')!=kendo.toString(oldVal.NacellePosition, 'n2')) {
                    $('#dir_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
                }
                
                window.setTimeout(function(){ 
                    $('#dir_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }
            if(val.RotorRPM > -999999) {
                if(kendo.toString(val.RotorRPM, 'n2')!= kendo.toString(oldVal.RotorRPM, 'n2')) {
                    $('#rotor_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
                }
                
                window.setTimeout(function(){ 
                    $('#rotor_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }
            if(val.PitchAngle > -999999) {
                if(kendo.toString(val.PitchAngle, 'n2')!=kendo.toString(oldVal.PitchAngle, 'n2')) {
                    $('#pitch_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
                }
                
                window.setTimeout(function(){ 
                    $('#pitch_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }
            if(val.Temperature > -999999) {
                if(kendo.toString(val.Temperature, 'n2')!=kendo.toString(oldVal.Temperature, 'n2')) {
                    $('#temperature_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');  
                }

                var curTempBA = bp.currentTempLocation() + 4;
                var curTempBB = bp.currentTempLocation() - 4;
                var turbineTemp = val.Temperature ;

                // console.log(curTempBA + " -- " + curTempBB + "--" + turbineTemp);
                if(turbineTemp < curTempBA && turbineTemp > curTempBB){
                    $('#temperaturecolor_'+ turbine).attr('class','fa fa-circle txt-green');
                }else if(turbineTemp > curTempBA){
                    $('#temperaturecolor_'+ turbine).attr('class','fa fa-circle txt-orange');
                }else if(turbineTemp < curTempBB){
                    $('#temperaturecolor_'+ turbine).attr('class','fa fa-circle txt-red');
                }

                
                window.setTimeout(function(){ 
                    $('#temperature_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }else{
                $('#temperaturecolor_'+ turbine).attr('class','fa fa-circle txt-grey');
            }


            var colorTemperature = val.TemperatureColor;
            $('#temperature_'+ turbine).attr('class', colorTemperature);
            
            /* TURBINE STATUS PART */
            if(val.AlarmDesc!="") {
                $('#alarmdesc_'+ turbine).text(val.AlarmCode);
                $('#alarmdesc_'+ turbine).attr('data-original-title', val.AlarmDesc);
            } else {
                $('#alarmdesc_'+ turbine).text("-");
                $('#alarmdesc_'+ turbine).attr('data-original-title', "This turbine already UP");
            }

            var colorStatus = "lbl bg-green";
            var defaultColorStatus = "bg-default-green";

            if(val.Status==0) {
                colorStatus = "lbl bg-red"; // faa-flash animated
                defaultColorStatus = "bg-default-red";
            } else if(val.Status === 1 && val.IsWarning === true) {
                colorStatus = "lbl bg-orange";
                defaultColorStatus = "bg-default-orange";
            }
            if(val.DataComing==0) {
                colorStatus = "lbl bg-grey";
                defaultColorStatus = "bg-default-grey";
            }

            var comparison = 0;
            $('#statusturbinedefault_'+ turbine).addClass(defaultColorStatus);
            if((val.ActivePower / val.Capacity) >= 0){
                comparison = (val.ActivePower / val.Capacity) * 70;
                
                $('#statusturbine_'+ turbine).attr('class', colorStatus);
                $('#statusturbine_'+ turbine).css('width', comparison + 'px');
            }else{
                comparison = 0;
                $('#statusturbine_'+ turbine).attr('class', 'lbl');
            }




            $('#statusturbinedefault_'+turbine).popover({
                placement: 'bottom',
                html: 'true',
                content : '<a class="btn btn-xs btn-primary individual"><i class="fa fa-line-chart"></i>&nbsp;Individual Turbine</a> &nbsp; <a class="btn btn-xs btn-primary alarm"><i class="fa fa-chevron-right"></i>&nbsp;View Alarm</a>'
            }).on('shown.bs.popover', function () {
                var $popup = $(this);
                $(this).next('.popover').find('a.individual').click(function (e) {
                    bp.ToIndividualTurbine(turbine);
                });
                $(this).next('.popover').find('a.alarm').click(function (e) {
                    bp.ToAlarm(turbine);
                });
            });
        }

    });




    $('#project_turbine_down').text(data.TurbineDown);
    $('#project_turbine_active').text(data.TurbineActive);
    $('#project_turbine_na').text(data.TurbineNotAvail);

    window.setTimeout(function(){ 
        $('#project_generation').text(data.PowerGeneration.toFixed(2));
        $('#project_wind_speed').text(data.AvgWindSpeed.toFixed(2));
        $('#project_plf').text((data.PLF).toFixed(2));
        bp.oldFeeders(allData);
    }, 1000);
    
};

bp.ToIndividualTurbine = function(turbine) {
    app.loading(true);
    var oldDateObj = new Date();
    var newDateObj = moment(oldDateObj).add(3, 'm');
    var project =  $('#projectList').data('kendoDropDownList').value();
    document.cookie = "projectname="+project.split("(")[0].trim()+";expires="+ newDateObj;
    document.cookie = "turbine="+turbine+";expires="+ newDateObj;
    if(document.cookie.indexOf("projectname=") >= 0 && document.cookie.indexOf("turbine=") >= 0) {
        window.location = viewModel.appName + "page/monitoringbyturbine";
    } else {
        app.loading(false);
    }
}

bp.ToAlarm = function(turbine) {
    app.loading(true);
    var oldDateObj = new Date();
    var newDateObj = moment(oldDateObj).add(3, 'm');
    var project =  $('#projectList').data('kendoDropDownList').value();
    
    document.cookie = "projectname="+project.split("(")[0].trim()+";expires="+ newDateObj;
    document.cookie = "turbine="+turbine+";expires="+ newDateObj;
    if(document.cookie.indexOf("projectname=") >= 0 && document.cookie.indexOf("turbine=") >= 0) {
        window.location = viewModel.appName + "page/monitoringalarm";
    } else {
        app.loading(false);
    }
}

bp.resetFeeders = function(){
    bp.oldFeeders([])
}

$(function() {
    app.loading(true);

    $("#restore-screen").hide();

    $('#projectList').kendoDropDownList({
        data: bp.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            setTimeout(function(){
                $.when(bp.resetFeeders()).done(function(){
                     bp.GetData();
                });
            },1500);
            
         }
    });

    $("#max-screen").click(function(){
        bp.fullscreen(true);
        $("html").addClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 80);
        $(".multicol").height($(window).innerHeight() - 80 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 80-50);
        $("#max-screen").hide();
        $("#restore-screen").show();  
    });

    $("#restore-screen").click(function(){
        bp.fullscreen(false);
        $("html").removeClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150-50);
        $("#max-screen").show();  
        $("#restore-screen").hide();  
    });

    $(document).on("click", ".popover .close" , function(){
        $(this).parents(".popover").popover('hide');
    });

    setTimeout(function() {
        bp.GetData()
        setInterval(bp.GetData, 4000);
    }, 600);
});