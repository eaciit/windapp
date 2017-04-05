'use strict';


viewModel.IndividualTurbine = new Object();
var it = viewModel.IndividualTurbine;


vm.currentMenu('Individual Turbine');
vm.currentTitle('Individual Turbine');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Individual Turbine', href: viewModel.appName + 'page/monitoringindividualturbine' }]);
var intervalTurbine = null;
var chart;
var maxSamples = 20,count = 0;

it.projectList = ko.observableArray([]);
it.project = ko.observable();
it.turbineList = ko.observableArray([{}]);
it.allTurbineList = ko.observableArray([{}]);

it.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        it.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val.split("(")[0].trim();
                data.text = val;
                datavalue.push(data);
            });
        }
        it.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(0);
            it.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

it.populateTurbine = function (project, turbine, isChange) {
    if (turbine.length == 0) {
        turbine = [];;
        it.turbineList([{ Id: "", Lat: 0.0, Lon: 0.0 }]);
    } else {
        var turbinevalue = [];
        if (turbine.length > 0) {
            it.allTurbineList = turbine;
            $.each(turbine, function (key, val) {
                if(val.project == project[0].split("(")[0].trim()) {
                    var data = {};
                    data.Id = val.turbineid;
                    data.Lat = val.latitude;
                    data.Lon = val.longitude;
                    turbinevalue.push(data);
                }
            });
        }
        it.turbineList(turbinevalue);

        if(isChange) {
            setTimeout(function () {
                $("#turbine").data("kendoDropDownList").select(0);
            }, 100);
        }
    }
};

it.ChangeProject = function() {
    return function() {
        it.project = $("#projectList").data("kendoDropDownList").value();
        var projects = [];
        projects.push(it.project);
        it.populateTurbine(projects, it.allTurbineList, true);
        it.ShowData();
    };
};

it.ChangeSelection = function() {
    return function() {
        // clearInterval(intervalTurbine);
        it.isFirst(true);
        it.ShowData();
        // intervalTurbine = window.setInterval(it.ShowData, 3000);
    };
};

it.getTimestamp = function(param){
  var dateString = moment(param).format("DD-MM-YYYY HH:mm:ss"),
      dateTimeParts = dateString.split(' '),
      timeParts = dateTimeParts[1].split(':'),
      dateParts = dateTimeParts[0].split('-'),
      date;

      date = new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0], timeParts[0], timeParts[1]);

      return date.getTime();
}

it.GetData = function(project, turbine) {
    var param = { 
        Project: project,
        Turbine: turbine
    }
    var getDetail = toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataturbine", param, function (res) {
        // var time = (new Date).getTime();

        var time = it.getTimestamp(moment.utc(res.data["lastupdate"]));

        it.dataWindspeed([time, parseFloat(res.data["Wind speed Avg"].toFixed(2))]);
        it.dataPower([time, parseFloat(res.data["Power"].toFixed(2))]);


        if(it.isFirst() == false){
            chart.series[0].addPoint([time, parseFloat(res.data["Wind speed Avg"].toFixed(2))], true, (++count >= maxSamples));
            chart.series[1].addPoint([time, parseFloat(res.data["Power"].toFixed(2))], true, (++count >= maxSamples));
        }else{
            it.showWindspeedLiveChart();
        }

        it.PlotData(res.data);
        it.isFirst(false);
        app.loading(false);
    });
};

it.rotor_rpm = ko.observable('');
it.rotor_r_cur = ko.observable('');
it.rotor_y_cur = ko.observable('');
it.windspeed_avg = ko.observable('');
it.windspeed1 = ko.observable('');
it.windspeed2 = ko.observable('');
it.wind_dir = ko.observable('');
it.wind_dir1 = ko.observable('');
it.wind_dir2 = ko.observable('');
it.nacel_dir = ko.observable('');
it.gen_rpm = ko.observable('');
it.speed_gen = ko.observable('');
it.blade_angle1 = ko.observable('');
it.blade_angle2 = ko.observable('');
it.blade_angle3 = ko.observable('');
it.volt_battery1 = ko.observable('');
it.volt_battery2 = ko.observable('');
it.volt_battery3 = ko.observable('');
it.cur_pitch_motor1 = ko.observable('');
it.cur_pitch_motor2 = ko.observable('');
it.cur_pitch_motor3 = ko.observable('');
it.pitch_motor_temp1 = ko.observable('');
it.pitch_motor_temp2 = ko.observable('');
it.pitch_motor_temp3 = ko.observable('');
it.phase_volt1 = ko.observable('');
it.phase_volt2 = ko.observable('');
it.phase_volt3 = ko.observable('');
it.phase_cur1 = ko.observable('');
it.phase_cur2 = ko.observable('');
it.phase_cur3 = ko.observable('');
it.power = ko.observable('');
it.power_react = ko.observable('');
it.freq_grid = ko.observable('');
it.dfig_act_power = ko.observable('');
it.dfig_react_power = ko.observable('');
it.dfig_main_freq = ko.observable('');
it.dfig_main_volt = ko.observable('');
it.dfig_main_cur = ko.observable('');
it.dfig_link_volt = ko.observable('');
it.rotor_b_cur = ko.observable('');
it.production = ko.observable('');
it.cos_phi = ko.observable('');
it.temp_gen_coil1 = ko.observable('');
it.temp_gen_coil2 = ko.observable('');
it.temp_gen_coil3 = ko.observable('');
it.temp_gen_bearing_driven = ko.observable('');
it.temp_gen_bearing_non_driven = ko.observable('');
it.temp_gear_driven = ko.observable('');
it.temp_gear_non_driven = ko.observable('');
it.temp_gear_inter_driven = ko.observable('');
it.temp_gear_inter_non_driven = ko.observable('');
it.press_gear_oil = ko.observable('');
it.temp_gear_oil = ko.observable('');
it.temp_nacelle = ko.observable('');
it.temp_ambient = ko.observable('');
it.temp_main_bear = ko.observable('');
it.damper_osci_mag = ko.observable('');
it.drive_train_vibra = ko.observable('');
it.tower_vibra = ko.observable('');
it.isFirst = ko.observable(true);

it.dataWindspeed = ko.observableArray([]);
it.dataPower = ko.observableArray([]);

it.rotor_status = ko.observable("N/A");

it.PlotData = function(data) {
    var lastUpdate = moment.utc(data["lastupdate"]);
    $('#turbine_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));

    if (it.rotor_status() === "N/A") {
        $('#rotorPic').css({"fill": "#b8b9bb"});
    } else if (it.rotor_status() == "up") {
        $('#rotorPic').css({"fill": "#31b445"});
    } else {
        $('#rotorPic').css({"fill": "#db1e1e"});
    }

    /*WIND SPEED PART*/
    if(data["Wind speed Avg"] != -999999)
        it.windspeed_avg(data["Wind speed Avg"].toFixed(2));
    else it.windspeed_avg('N/A');
    if(data["Wind speed 1"] != -999999)
        it.windspeed1(data["Wind speed 1"].toFixed(2));
    else it.windspeed1('N/A');
    if(data["Wind speed 2"] != -999999)
        it.windspeed2(data["Wind speed 2"].toFixed(2));
    else it.windspeed2('N/A');

    it.showWindspeedColumnChart();

    /*WIND DIRECTION PART*/
    if(data["Wind Direction"] != -999999)
        it.wind_dir(data["Wind Direction"].toFixed(2));
    else it.wind_dir('N/A');
    if(data["Vane 1 wind direction"] != -999999)
        it.wind_dir1(data["Vane 1 wind direction"].toFixed(2));
    else it.wind_dir1('N/A');
    if(data["Vane 2 wind direction"] != -999999)
        it.wind_dir2(data["Vane 2 wind direction"].toFixed(2));
    else it.wind_dir2('N/A');

    /*NACELLE POSITION PART*/
    if(data["Nacelle Direction"] != -999999)
        it.nacel_dir(data["Nacelle Direction"].toFixed(2));
    else it.nacel_dir('N/A');
    if(data["Rotor RPM"] != -999999)
        it.rotor_rpm(data["Rotor RPM"].toFixed(2));
    else it.rotor_rpm('N/A');
    if(data["Generator RPM"] != -999999)
        it.gen_rpm(data["Generator RPM"].toFixed(2));
    else it.gen_rpm('N/A');
    if(data["DFIG speed generator encoder"] != -999999)
        it.speed_gen(data["DFIG speed generator encoder"].toFixed(2));
    else it.speed_gen('N/A');

    it.showRotor();

    /*BLADE ANGLE PART*/
    if(data["Blade Angle 1"] != -999999)
        it.blade_angle1(data["Blade Angle 1"].toFixed(2));
    else it.blade_angle1('N/A');
    if(data["Blade Angle 2"] != -999999)
        it.blade_angle2(data["Blade Angle 2"].toFixed(2));
    else it.blade_angle2('N/A');
    if(data["Blade Angle 3"] != -999999)
        it.blade_angle3(data["Blade Angle 3"].toFixed(2));
    else it.blade_angle3('N/A');

    /*VOLT BATTERY BLADE PART*/
    if(data["Volt. Battery - blade 1"] != -999999)
        it.volt_battery1(data["Volt. Battery - blade 1"].toFixed(2));
    else it.volt_battery1('N/A');
    if(data["Volt. Battery - blade 2"] != -999999)
        it.volt_battery2(data["Volt. Battery - blade 2"].toFixed(2));
    else it.volt_battery2('N/A');
    if(data["Volt. Battery - blade 3"] != -999999)
        it.volt_battery3(data["Volt. Battery - blade 3"].toFixed(2));
    else it.volt_battery3('N/A');

    /*CURRENT PITCH MOTOR PART*/
    if(data["Current 1 Pitch Motor"] != -999999)
        it.cur_pitch_motor1(data["Current 1 Pitch Motor"].toFixed(2));
    else it.cur_pitch_motor1('N/A');
    if(data["Current 2 Pitch Motor"] != -999999)
        it.cur_pitch_motor2(data["Current 2 Pitch Motor"].toFixed(2));
    else it.cur_pitch_motor2('N/A');
    if(data["Current 3 Pitch Motor"] != -999999)
        it.cur_pitch_motor3(data["Current 3 Pitch Motor"].toFixed(2));
    else it.cur_pitch_motor3('N/A');

    /*PITCH MOTOR TEMPERATURE PART*/
    if(data["Pitch motor temperature - Blade 1"] != -999999)
        it.pitch_motor_temp1(data["Pitch motor temperature - Blade 1"].toFixed(2));
    else it.pitch_motor_temp1('N/A');
    if(data["Pitch motor temperature - Blade 2"] != -999999)
        it.pitch_motor_temp2(data["Pitch motor temperature - Blade 2"].toFixed(2));
    else it.pitch_motor_temp2('N/A');
    if(data["Pitch motor temperature - Blade 3"] != -999999)
        it.pitch_motor_temp3(data["Pitch motor temperature - Blade 3"].toFixed(2));
    else it.pitch_motor_temp3('N/A');

    /*PHASE VOLTAGE PART*/
    if(data["Phase 1 voltage"] != -999999)
        it.phase_volt1(data["Phase 1 voltage"].toFixed(2));
    else it.phase_volt1('N/A');
    if(data["Phase 2 voltage"] != -999999)
        it.phase_volt2(data["Phase 2 voltage"].toFixed(2));
    else it.phase_volt2('N/A');
    if(data["Phase 3 voltage"] != -999999)
        it.phase_volt3(data["Phase 3 voltage"].toFixed(2));
    else it.phase_volt3('N/A');

    /*PHASE CURRENT PART*/
    if(data["Phase 1 current"] != -999999)
        it.phase_cur1(data["Phase 1 current"].toFixed(2));
    else it.phase_cur1('N/A');
    if(data["Phase 2 current"] != -999999)
        it.phase_cur2(data["Phase 2 current"].toFixed(2));
    else it.phase_cur2('N/A');
    if(data["Phase 3 current"] != -999999)
        it.phase_cur3(data["Phase 3 current"].toFixed(2));
    else it.phase_cur3('N/A');

    /*POWER PART*/
    if(data["Power"] != -999999)
        it.power(data["Power"].toFixed(2));
    else it.power('N/A');
    if(data["Power Reactive"] != -999999)
        it.power_react(data["Power Reactive"].toFixed(2));
    else it.power_react('N/A');
    if(data["Freq. Grid"] != -999999)
        it.freq_grid(data["Freq. Grid"].toFixed(2));
    else it.freq_grid('N/A');

    /*DFIG PART*/
    if(data["DFIG active power"] != -999999)
        it.dfig_act_power(data["DFIG active power"].toFixed(2));
    else it.dfig_act_power('N/A');
    if(data["DFIG reactive power"] != -999999)
        it.dfig_react_power(data["DFIG reactive power"].toFixed(2));
    else it.dfig_react_power('N/A');
    if(data["DFIG mains Frequency"] != -999999)
        it.dfig_main_freq(data["DFIG mains Frequency"].toFixed(2));
    else it.dfig_main_freq('N/A');

    if(data["DFIG main voltage"] != -999999)
        it.dfig_main_volt(data["DFIG main voltage"].toFixed(2));
    else it.dfig_main_volt('N/A');
    if(data["DFIG main current"] != -999999)
        it.dfig_main_cur(data["DFIG main current"].toFixed(2));
    else it.dfig_main_cur('N/A');
    if(data["DFIG DC link voltage"] != -999999)
        it.dfig_link_volt(data["DFIG DC link voltage"].toFixed(2));
    else it.dfig_link_volt('N/A');

    /*ROTOR CURRENT PART*/
    if(data["Rotor R current"] != -999999) {
        it.rotor_r_cur(data["Rotor R current"].toFixed(2));
    }
    else { 
        it.rotor_r_cur('N/A');
    }
    if(data["Roter Y current"] != -999999) {
        it.rotor_y_cur(data["Roter Y current"].toFixed(2));
    }
    else { 
        it.rotor_y_cur('N/A');
    }
    if(data["Roter B current"] != -999999)
        it.rotor_b_cur(data["Roter B current"].toFixed(2));
    else it.rotor_b_cur('N/A');


    /*PRODUCTION PART*/
    if(data["Production"] != -999999)
        it.production(data["Production"].toFixed(2));
    else it.production('N/A');
    if(data["Cos Phi"] != -999999)
        it.cos_phi(data["Cos Phi"].toFixed(2));
    else it.cos_phi('N/A');

    /*TEMP GENERATOR PART*/
    if(data["Temp. generator 1 phase 1 coil"] != -999999)
        it.temp_gen_coil1(data["Temp. generator 1 phase 1 coil"].toFixed(2));
    else it.temp_gen_coil1('N/A');
    if(data["Temp. generator 1 phase 2 coil"] != -999999)
        it.temp_gen_coil2(data["Temp. generator 1 phase 2 coil"].toFixed(2));
    else it.temp_gen_coil2('N/A');
    if(data["Temp. generator 1 phase 3 coil"] != -999999)
        it.temp_gen_coil3(data["Temp. generator 1 phase 3 coil"].toFixed(2));
    else it.temp_gen_coil3('N/A');
    if(data["Temp. generator bearing driven End"] != -999999)
        it.temp_gen_bearing_driven(data["Temp. generator bearing driven End"].toFixed(2));
    else it.temp_gen_bearing_driven('N/A');
    if(data["Temp. generator bearing non-driven End"] != -999999)
        it.temp_gen_bearing_non_driven(data["Temp. generator bearing non-driven End"].toFixed(2));
    else it.temp_gen_bearing_non_driven('N/A');

    /*TEMP GEARBOX PART*/
    if(data["Temp. Gearbox driven end"] != -999999)
        it.temp_gear_driven(data["Temp. Gearbox driven end"].toFixed(2));
    else it.temp_gear_driven('N/A');
    if(data["Temp. Gearbox non-driven end"] != -999999)
        it.temp_gear_non_driven(data["Temp. Gearbox non-driven end"].toFixed(2));
    else it.temp_gear_non_driven('N/A');
    if(data["Temp. Gearbox inter. driven end"] != -999999)
        it.temp_gear_inter_driven(data["Temp. Gearbox inter. driven end"].toFixed(2));
    else it.temp_gear_inter_driven('N/A');
    if(data["Temp. Gearbox inter. non-driven end"] != -999999)
        it.temp_gear_inter_non_driven(data["Temp. Gearbox inter. non-driven end"].toFixed(2));
    else it.temp_gear_inter_non_driven('N/A');
    if(data["Pressure Gear box oil"] != -999999)
        it.press_gear_oil(data["Pressure Gear box oil"].toFixed(2));
    else it.press_gear_oil('N/A');
    if(data["Temp. Gear box oil"] != -999999)
        it.temp_gear_oil(data["Temp. Gear box oil"].toFixed(2));
    else it.temp_gear_oil('N/A');

    /*NACELLE PART*/
    if(data["Temp. Nacelle"] != -999999)
        it.temp_nacelle(data["Temp. Nacelle"].toFixed(2));
    else it.temp_nacelle('N/A');
    if(data["Temp. Ambient"] != -999999)
        it.temp_ambient(data["Temp. Ambient"].toFixed(2));
    else it.temp_ambient('N/A');
    if(data["Temp. Main bearing"] != -999999)
        it.temp_main_bear(data["Temp. Main bearing"].toFixed(2));
    else it.temp_main_bear('N/A');

    /*VIBRATION PART*/
    if(data["Damper Oscillation mag."] != -999999)
        it.damper_osci_mag(data["Damper Oscillation mag."].toFixed(2));
    else it.damper_osci_mag('N/A');
    if(data["Drive train vibration"] != -999999)
        it.drive_train_vibra(data["Drive train vibration"].toFixed(2));
    else it.drive_train_vibra('N/A');
    if(data["Tower vibration"] != -999999)
        it.tower_vibra(data["Tower vibration"].toFixed(2));
    else it.tower_vibra('N/A');


    it.changeRotation();
    it.changeColor();
};

it.LoadData = function(turbine) {
    var selLat = 0;
    var selLon = 0;
    $.each(it.turbineList(), function(idx, val){
        if(val.Id == turbine) {
            selLon = val.Lon;
            selLat = val.Lat;
        }
    });
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    var param = { "lat": selLat, "lon": selLon, "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };
    $.ajax({
        type: "GET",
        url: surl,
        data: param,
        dataType: "jsonp",
        success:function(data){
          it.ParseWeather(data);
        },
        error:function(){
            // do nothing
        }  
    });
};

it.ParseWeather = function(data) {
    $('#turbine_lat').text(data.coord.lat);
    $('#turbine_lon').text(data.coord.lon);
    $('#turbine_location').text(data.name + ' ('+ data.sys.country +')');
    $('#turbine_img_weather').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
    $('#weather').text(data.weather[0].description);
    $('#turbine_weather').text(data.weather[0].description);
    $('#turbine_temperature').text(data.main.temp);
};

it.ShowData = function() {
    if(it.isFirst() == true){
        app.loading(true);
    }
    var COOKIES = {};
    var cookieStr = document.cookie;
    var turbine = "";
    var project = "";
    
    if(cookieStr.indexOf("turbine=") >= 0 && cookieStr.indexOf("project=") >= 0) {
        document.cookie = "project=;expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        document.cookie = "turbine=;expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        cookieStr.split(/; /).forEach(function(keyValuePair) {
            var cookieName = keyValuePair.replace(/=.*$/, "");
            var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
            COOKIES[cookieName] = cookieValue;
        });
        turbine = COOKIES["turbine"];
        project = COOKIES["project"];
        $('#turbine').data('kendoDropDownList').value(turbine);
        $('#projectList').data('kendoDropDownList').value(project);
    } else {
        turbine = $('#turbine').data('kendoDropDownList').value();
        project = $('#projectList').data('kendoDropDownList').value();
    }
    
    it.LoadData(turbine);
    it.GetData(project, turbine);
    it.showWindRoseChart();
};

it.showWindspeedColumnChart = function(){

    $("#compareWindChart").kendoLinearGauge({
        theme: "flat",
        gaugeArea: {
          height : 125
        },
        pointer: {
            value: it.windspeed_avg(),
            shape: "arrow"
        },
        scale: {
            majorUnit: 10,
            minorUnit: 5,
            // max: 180,
            ranges: [
                {
                    from: 30,
                    to: 50,
                    color : "#8dcb2a"
                },
                {
                    from: 20,
                    to: 30,
                    color: "#ffc700"
                }, {
                    from: 10,
                    to: 20,
                    color: "#ff7a00"
                }, {
                    from: 0,
                    to: 10,
                    color: "#c20000"
                }
            ]
        }
    });
    $("#comparePowerChart").kendoLinearGauge({
        theme: "flat",
        gaugeArea: {
          height : 125
        },
        pointer: {
            value: it.power(),
            shape: "arrow"
        },
        scale: {
            minorUnit: 5,
            // max: 180,
            majorUnit: 10,
            ranges: [
                {
                    from: 30,
                    to: 50,
                    color : "#8dcb2a"
                },
                {
                    from: 20,
                    to: 30,
                    color: "#ffc700"
                }, {
                    from: 10,
                    to: 20,
                    color: "#ff7a00"
                }, {
                    from: 0,
                    to: 10,
                    color: "#c20000"
                }
            ]
        }
    });

}
it.showWindspeedLiveChart = function(){

    Highcharts.setOptions({
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                fontSize: '12px',
                fontWeight: 'normal',
            },
        },
        global: {
            useUTC: false
        }
    });

    chart = Highcharts.stockChart('container', {
        chart: {
            marginTop: 50,
            height: 180,
            width: 500,
        },
        credits: {
              enabled: false
        },
        legend: {
            enabled: true,
            verticalAlign: 'top',
            layout: "horizontal",
            labelFormat: '<span style="color:{color}">{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
        },
        rangeSelector: {
            buttons: [{
                count: 1,
                type: 'minute',
                text: '1M'
            }, {
                count: 5,
                type: 'minute',
                text: '5M'
            }, {
                type: 'all',
                text: 'All'
            }],
            inputEnabled: false,
            selected: 0,
            enabled:false,
        },
            scrollbar: {
                enabled: false
            },
        exporting: {
            enabled: false
        },
        yAxis:{
          labels:
          {
            enabled: false
          },
          gridLineWidth: 0,
          minorGridLineWidth: 0,
          opposite:false
        },
        xAxis: {
           lineWidth: 0,
           minorGridLineWidth: 0,
           lineColor: 'transparent',
           labels: {
               enabled: false
           },
           minorTickLength: 0,
           tickLength: 0
        },

            navigator: {
            enabled: false,
        },
        plotOptions: {
            series: {
                lineWidth: 1,
                marker: {
                    enabled: true,
                    radius: 3
                },
            },
        },
        tooltip:{
             formatter : function() {
                $("#dateInfo").html( Highcharts.dateFormat('%e %b %Y %H:%M:%S', this.x));
                return false ;
             }
        },
        series: [{
            color: colorField[0],
            name: 'Wind Speed',
            data: [it.dataWindspeed()],
            tooltip: {
                valueSuffix: 'm/s'
            },
        },{
            name: 'Power',
            color: colorField[1],
            data: [it.dataPower()],
            tooltip: {
                valueSuffix: 'kW'
            },
        }]
    });

}
it.showRotor = function(){
    $("#rotorChart").kendoRadialGauge({
        title: "Rotor RPM",
        theme: "flat",
        pointer: {
            value: it.rotor_rpm()
        },
        gaugeArea: {
          height : 125,
        },
        scale: {
            minorUnit: 5,
            majorUnit: 40,
            startAngle: -30,
            endAngle: 210,
            max: 180,
            labels: {
                position: "inside"
            },
            ranges: [
                 {
                    from: 100,
                    to: 160,
                    color : "#8dcb2a"
                },
                {
                    from: 60,
                    to: 100,
                    color: "#ffc700"
                }, {
                    from: 20,
                    to: 60,
                    color: "#ff7a00"
                }, {
                    from: 0,
                    to: 20,
                    color: "#c20000"
                }
            ]
        }
    });
}

it.showWindRoseChart = function(){
    var param = {
        turbine: $("#turbine").data("kendoDropDownList").value(),
        project: $("#projectList").data("kendoDropDownList").value(),
        breakDown: 12,
    };
    toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getwindrosemonitoring", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        if (res.data != null) {
            var nacelleData = res.data["nacelle"]["WindRose"][0].Data;
            var winddirData = res.data["winddir"]["WindRose"][0].Data;

            $("#windRoseChart").kendoChart({
                theme: "flat",
                chartArea: {
                    height: 200,
                    width: 250,
                    margin: 0,
                    padding: 0,
                },
                dataSource: {
                    data: nacelleData,
                    group: {
                        field: "WsCategoryNo",
                        dir: "asc"
                    },
                    sort: {
                        field: "DirectionNo",
                        dir: "asc"
                    }
                },
                legend: {
                    visible: false,
                },
                series: [{
                    type: "radarColumn",
                    stack: true,
                    field: "Contribution",
                    gap: 0,
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
                        step: 1
                    },
                    labels: {
                        font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        visible: true,
                        step: 1
                    }
                },
                valueAxis: {
                    labels: {
                        template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                        font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                    }
                },
            });

            $("#windDirectionChart").kendoChart({
                theme: "flat",
                chartArea: {
                    height: 200,
                    width: 250,
                    margin: 0,
                    padding: 0,
                },
                dataSource: {
                    data: winddirData,
                    group: {
                        field: "WsCategoryNo",
                        dir: "asc"
                    },
                    sort: {
                        field: "DirectionNo",
                        dir: "asc"
                    }
                },
                legend: {
                    visible: false,
                },
                series: [{
                    type: "radarColumn",
                    stack: true,
                    field: "Contribution",
                    gap: 0,
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
                        step: 1
                    },
                    labels: {
                        font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        visible: true,
                        step: 1
                    }
                },
                valueAxis: {
                    labels: {
                        template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                        font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                    }
                },
            });
        }
    })
}

it.ToTimeSeriesHfd = function() {
    var turbine = $("#turbine").val();
    var oldDateObj = new Date();
    var newDateObj = moment(oldDateObj).add(3, 'm');
    document.cookie = "turbine="+turbine+"; expires="+ newDateObj;
    window.location = viewModel.appName + "page/timeserieshfd";
}

it.changeRotation = function(){
    $.each( $('.rotation'), function( key, value ) {
        var deg = $(value).attr("rotationval")
        $(value).attr("style", $(value).attr("style")+"-ms-transform: rotate("+deg+"deg);-webkit-transform: rotate("+deg+"deg);transform: rotate("+deg+"deg);");
    });
}

it.changeColor = function(){
    $('span').each(function(){
      if($(this).html() == 'N/A'){
         $( this ).css( "color", "red" );
      }else{
        $( this ).css( "color", "#585555" );
      }
    });
}

$(document).ready(function(){
    $.when(it.ShowData()).done(function () {
        setTimeout(function() {
            it.isFirst(false);
            app.loading(false);
        }, 1000);
    });
    intervalTurbine = window.setInterval(it.ShowData, 6000);
});