'use strict';


viewModel.IndividualTurbine = new Object();
var it = viewModel.IndividualTurbine;


vm.currentMenu('Individual Turbine');
vm.currentTitle('Individual Turbine');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Individual Turbine', href: viewModel.appName + 'page/monitoringindividualturbine' }]);
var intervalTurbine = null;
it.projectList = ko.observableArray([]);
it.project = ko.observable();

it.turbineList = ko.observableArray([
    { Id: "HBR004", Lat: 27.1314613912264, Lon: 70.6185588069878 },
    { Id: "HBR005", Lat: 27.1270877985472, Lon: 70.622884146285 },
    { Id: "HBR006", Lat: 27.1236484571162, Lon: 70.6283625931435 },
    { Id: "HBR007", Lat: 27.1207371982508, Lon: 70.631881314612 },
    { Id: "SSE001", Lat: 27.01307895, Lon: 70.49387633 },
    { Id: "SSE002", Lat: 27.010279756987, Lon: 70.4973565360181 },
    { Id: "SSE006", Lat: 27.0210933331988, Lon: 70.5024788963212 },
    { Id: "SSE007", Lat: 27.0248279694312, Lon: 70.5001700425682 },
    { Id: "SSE011", Lat: 27.0402772539165, Lon: 70.5014338714604 },
    { Id: "SSE012", Lat: 27.0337931496561, Lon: 70.5042405803792 },
    { Id: "SSE015", Lat: 27.0388148011425, Lon: 70.5123920837843 },
    { Id: "SSE017", Lat: 27.0548841404124, Lon: 70.5080202799257 },
    { Id: "SSE018", Lat: 27.0528220756451, Lon: 70.5125899817518 },
    { Id: "SSE019", Lat: 27.0495446014842, Lon: 70.5168508007021 },
    { Id: "SSE020", Lat: 27.044182589382, Lon: 70.5210429390468 },
    { Id: "TJ013", Lat: 27.1176472129651, Lon: 70.6360125744664 },
    { Id: "TJ016", Lat: 27.1132054610326, Lon: 70.6391757345863 },
    { Id: "TJ021", Lat: 27.0972435993558, Lon: 70.6640052338389 },
    { Id: "TJ022", Lat: 27.0987203759723, Lon: 70.6566747880341 },
    { Id: "TJ023", Lat: 27.10074814474, Lon: 70.6515811382493 },
    { Id: "TJ024", Lat: 27.1030071512702, Lon: 70.6459964756472 },
    { Id: "TJ025", Lat: 27.1067003493502, Lon: 70.6412494768147 },
    { Id: "HBR038", Lat: 27.0920863398176, Lon: 70.6257778736997 },
    { Id: "TJW024", Lat: 27.0590442370984, Lon: 70.6309369071587 },
]);

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

it.ChangeSelection = function() {
    return function() {
        // clearInterval(intervalTurbine);
        it.ShowData();
        // intervalTurbine = window.setInterval(it.ShowData, 3000);
    };
};

it.GetData = function(turbine) {
    var param = { Turbine: turbine }
    var getDetail = toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataturbine", param, function (res) {
        it.PlotData(res.data);
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
    var COOKIES = {};
    var cookieStr = document.cookie;
    var turbine = "";
    
    if(cookieStr.indexOf("turbine=") >= 0) {
        document.cookie = "turbine=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        cookieStr.split(/; /).forEach(function(keyValuePair) {
            var cookieName = keyValuePair.replace(/=.*$/, "");
            var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
            COOKIES[cookieName] = cookieValue;
        });
        turbine = COOKIES["turbine"];
        $('#turbine').data('kendoDropDownList').value(turbine);
    } else {
        turbine = $('#turbine').data('kendoDropDownList').value();
    }
    
    it.LoadData(turbine);
    it.GetData(turbine);
};


$(document).ready(function(){
    app.loading(true);
    $.when(it.ShowData()).done(function () {
        setTimeout(function() {
            app.loading(false);
        }, 1000);
    });
    intervalTurbine = window.setInterval(it.ShowData, 6000);
});