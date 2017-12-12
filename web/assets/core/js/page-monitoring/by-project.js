'use strict';


viewModel.ByProject = new Object();
var bp = viewModel.ByProject;


vm.currentMenu('By Project');
vm.currentTitle('By Project');
vm.isShowDataAvailability(false);
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

bp.isFirst = ko.observable(true);

// var count = 0;
var requests = [];

var audioElement = document.createElement('audio');
    audioElement.setAttribute('src', "../res/alarm/alarm.mp3");

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
   requests.push(toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
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
                    feedHeader.TotalProdDay = null;
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
    }));
}

bp.GetData = function(data) {

    var project = "";

    if(localStorage.getItem("projectname") !== null) {
        project = localStorage.getItem("projectname");
        $('#projectList').data('kendoDropDownList').value(project);
        app.resetLocalStorage();
    } else {
        project = $('#projectList').data('kendoDropDownList').value();
    }
    
    $.when(bp.CheckWeather()).done(function(){
        setTimeout(function(){
             bp.GetDataProject(project);
        }, 200);
    });

    // count++;
};

bp.PlotData = function(data) {
    audioElement.currentTime = 0;
    audioElement.pause();
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
        var oldTurbine = '';
        try {
            oldTurbine = oldData[idx].Turbine; // for a while to minimalize an error
        }
        catch(e) {}
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

                if(oldVal.ColorStatus == "lbl bg-green" && val.ColorStatus == "lbl bg-red"){
                    var playPromise = audioElement.play();
                    if (playPromise !== null){
                        playPromise.catch(() => { audioElement.play(); })
                    }
                }else{
                    audioElement.pause();
                }
                
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
            if(val.TotalProdDay > -999999) {
                if(kendo.toString(val.TotalProdDay, 'n2')!= kendo.toString(oldVal.TotalProdDay, 'n2')) {
                    $('#total_'+ turbine).css('background-color', 'rgba(255, 216, 0, 0.7)');    
                }
                
                window.setTimeout(function(){ 
                    $('#total_'+ turbine).css('background-color', 'transparent'); 
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
                window.setTimeout(function(){ 
                    $('#temperature_'+ turbine).css('background-color', 'transparent'); 
                }, 750);
            }

            // if(count == 5 || count == 6 || count == 7){
            //     console.log("yayay");
            //     if(turbine == "HBR004"){
            //         val.isbordered = true;
            //     }
            // }

            var colorTemperature = val.TemperatureColor;
            $('#temperature_'+ turbine).attr('class', colorTemperature);
            
            $('#temperaturecolor_'+ turbine).addClass(val.BulletColor);

            if (val.TemperatureInfo != "" && val.TemperatureInfo != undefined) {
                $('#temperaturecolor_'+ turbine).attr('data-original-title', val.TemperatureInfo);
            } else {
                $('#temperaturecolor_'+ turbine).attr('data-original-title', "");
            }
            
            /* TURBINE STATUS PART */
            if(val.AlarmDesc!="") {
                $('#alarmdesc_'+ turbine).text(val.AlarmCode);
                $('#alarmdesc_'+ turbine).attr('data-original-title', val.AlarmDesc);
            } else {
                $('#alarmdesc_'+ turbine).text("-");
                $('#alarmdesc_'+ turbine).attr('data-original-title', "This turbine already UP");
            }

            var colorStatus = val.ColorStatus;
            var defaultColorStatus = val.DefaultColorStatus;

            var comparison = 0;
            $('#statusturbinedefault_'+ turbine).addClass(defaultColorStatus);
            
            
            if((val.ActivePower / val.Capacity) > 0){
                comparison = (val.ActivePower / val.Capacity) * 65;
                var fixCom = (comparison > 64 ? 64 : comparison);
                $('#statusturbine_'+ turbine).attr('class', colorStatus);
                $('#statusturbine_'+ turbine).css('width', fixCom + 'px');
            }else{
                comparison = 0;
                $('#statusturbine_'+ turbine).attr('class', 'lbl');
            }

            if(colorStatus=="lbl bg-red"){
                $('#statusturbine_'+ turbine).attr('class', colorStatus);
                $('#statusturbine_'+ turbine).css('width',  65 +'px');
            }

            if(val.isbordered != undefined && val.isbordered == true){
                // $('#statusturbinedefault_'+turbine).addClass("bordered");
                $('#statusturbinedefault_'+turbine).find(".inner-triangle").show();
            }else{
                // $('#statusturbinedefault_'+turbine).removeClass("bordered");
                 $('#statusturbinedefault_'+turbine).find(".inner-triangle").hide();
            }

            if(val.IsReapeatedAlarm == true){
                $('#linkDetail_'+turbine).addClass("reapeat-alarm");
            }else{
                $('#linkDetail_'+turbine).removeClass("reapeat-alarm");
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


            if(val.IsRemark == true){
                 $('#iconTurbine_'+turbine).css("display", "block !important");
            }else{
                 $('#iconTurbine_'+turbine).hide();
            }


            $('#linkDetail_'+turbine).on("dblclick", function(e){
                TbCol.ResetData();

                var classIcon = 'txt-green';
                var classString = defaultColorStatus;
                if(classString !== undefined){
                    classIcon = 'txt'+ classString.substr(10);
                }
                
                TbCol.TurbineId(turbine);
                TbCol.TurbineName(val.Name);
                TbCol.UserId('');
                TbCol.UserName('');
                TbCol.Project(data.ProjectName);
                TbCol.Feeder(val.feederName);
                TbCol.IsTurbine(true);
                TbCol.OpenForm();
                TbCol.IconStatus(classIcon);
            });
        }

    });

    var $feederList = data.FeederRemarkList;

    $.each($feederList , function(key, val){
        if(val == true){
            $(".tableDetails").find('.icon-remark[data-id="iconFeeder_'+ key +'"]').css("display", "block !important");
        }else{
            $(".tableDetails").find('.icon-remark[data-id="iconFeeder_'+ key +'"]').hide();
        }

        $(".tableDetails").find('.feederRemark[data-id="linkFeeder_'+ key +'"]').on("click", function(e){
            TbCol.ResetData();
            TbCol.ProjectFeeder(data.ProjectName)
            TbCol.Feeder(key);
            TbCol.IsTurbine(false);
            TbCol.OpenForm();
        });
    });



    $('#project_turbine_down').text(data.TurbineDown);
    $('#project_turbine_active').text(data.TurbineActive);
    $('#project_turbine_na').text(data.TurbineNotAvail);
    $('#project_waiting_wind').text(data.TurbineWaitingWS);

    window.setTimeout(function(){ 
        $('#project_generation').text(data.PowerGeneration.toFixed(2));
        $('#project_wind_speed').text(data.AvgWindSpeed.toFixed(2));
        $('#project_plf').text((data.PLF).toFixed(2));
        $('#avg_tempout').text((data.AvgTempOutdoor).toFixed(2));
        bp.oldFeeders(allData);
    }, 1000);
    
};


bp.ToIndividualTurbine = function(turbine){
    setTimeout(function(){
        app.loading(true);
        app.resetLocalStorage();
        var project =  $('#projectList').data('kendoDropDownList').value();
        localStorage.setItem('turbine', turbine);
        localStorage.setItem('projectname', project);
        localStorage.setItem('isFromSummary', false);
        localStorage.setItem('isFromByProject', true);
        if(localStorage.getItem("turbine") !== null && localStorage.getItem("projectname") !== null){
            window.location = viewModel.appName + "page/monitoringbyturbine";
        }
    },1500);
}

bp.ToAlarm = function(turbine) {
    setTimeout(function(){
        app.loading(true);
        app.resetLocalStorage();
        var project =  $('#projectList').data('kendoDropDownList').value();
        localStorage.setItem('turbine', turbine == [] ? null : turbine);
        localStorage.setItem('projectname', project);
        localStorage.setItem('tabActive', "default");
        if(localStorage.getItem("turbine") !== null && localStorage.getItem("projectname") !== null){
            window.location = viewModel.appName + "page/monitoringalarm";
        }
    },1500);
}

bp.ToSummary = function(){
    window.location = viewModel.appName + "page/monitoringsummary";
}

bp.abortAll = function(requests) {
     // count = 0;
     bp.resetFeeders();
     var length = requests.length;
     while(length--) {
         requests[length].abort && requests[length].abort();  // the if is for the first case mostly, where array is still empty, so no abort method exists.
     }
}

bp.resetFeeders = function(){
    bp.oldFeeders([])
}

bp.setFullScreen = function(){
    if ((document.fullScreenElement && document.fullScreenElement !== null) || (!document.mozFullScreen && !document.webkitIsFullScreen)) {
        bp.fullscreen(true);
        if (document.documentElement.requestFullScreen) {
            document.documentElement.requestFullScreen();
        } else if (document.documentElement.mozRequestFullScreen) {
            document.documentElement.mozRequestFullScreen();
        } else if (document.documentElement.webkitRequestFullScreen) {
            document.documentElement.webkitRequestFullScreen(Element.ALLOW_KEYBOARD_INPUT);
        }
    } else {

        bp.fullscreen(false);
        if (document.cancelFullScreen) {
            document.cancelFullScreen();
        } else if (document.mozCancelFullScreen) {
            document.mozCancelFullScreen();
        } else if (document.webkitCancelFullScreen) {
            document.webkitCancelFullScreen();
        }
    }
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
                bp.isFirst(true);
                bp.abortAll(requests);
                bp.GetData();
            },1500);
            
         }
    });

    $("#max-screen").click(function(){
        bp.setFullScreen();
        $("html").addClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 80);
        $(".multicol").height($(window).innerHeight() - 80 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 80-50);
        $("#max-screen").hide();
        $("#restore-screen").show();  
    });

    $("#restore-screen").click(function(){
        bp.setFullScreen();
        $("html").removeClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150-50);
        $("#max-screen").show();  
        $("#restore-screen").hide();  
    });

    $('.bstooltip').mouseenter(function(){
        var that = $(this)
        that.tooltip('show');
        setTimeout(function(){
            that.tooltip('hide');
        }, 7000);
    });

    $('.bstooltip').mouseleave(function(){
        $(this).tooltip('hide');
    });

    $(document).on("click", ".popover .close" , function(){
        $(this).parents(".popover").popover('hide');
    });

    setTimeout(function() {
        bp.GetData()
        setInterval(bp.GetData, 5000);
    }, 600);
});