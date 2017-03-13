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

it.PlotData = function(data) {

    var lastUpdate = moment.utc(data["lastupdate"]);
    $('#turbine_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));

    /*WIND SPEED PART*/
    if(data["Wind speed Avg"] != -999999)
        $('#windspeed_avg').text(data["Wind speed Avg"].toFixed(2));
    else $('#windspeed_avg').text('N/A');
    if(data["Wind speed 1"] != -999999)
        $('#windspeed1').text(data["Wind speed 1"].toFixed(2));
    else $('#windspeed1').text('N/A');
    if(data["Wind speed 2"] != -999999)
        $('#windspeed2').text(data["Wind speed 2"].toFixed(2));
    else $('#windspeed2').text('N/A');

    /*WIND DIRECTION PART*/
    if(data["Wind Direction"] != -999999)
        $('#wind_dir').text(data["Wind Direction"].toFixed(2));
    else $('#wind_dir').text('N/A');
    if(data["Vane 1 wind direction"] != -999999)
        $('#wind_dir1').text(data["Vane 1 wind direction"].toFixed(2));
    else $('#wind_dir1').text('N/A');
    if(data["Vane 2 wind direction"] != -999999)
        $('#wind_dir2').text(data["Vane 2 wind direction"].toFixed(2));
    else $('#wind_dir2').text('N/A');

    /*NACELLE POSITION PART*/
    if(data["Nacelle Direction"] != -999999)
        $('#nacel_dir').text(data["Nacelle Direction"].toFixed(2));
    else $('#nacel_dir').text('N/A');
    if(data["Rotor RPM"] != -999999)
        $('#rotor_rpm').text(data["Rotor RPM"].toFixed(2));
    else $('#rotor_rpm').text('N/A');
    if(data["Generator RPM"] != -999999)
        $('#gen_rpm').text(data["Generator RPM"].toFixed(2));
    else $('#gen_rpm').text('N/A');
    if(data["DFIG speed generator encoder"] != -999999)
        $('#speed_gen').text(data["DFIG speed generator encoder"].toFixed(2));
    else $('#speed_gen').text('N/A');

    /*BLADE ANGLE PART*/
    if(data["Blade Angle 1"] != -999999)
        $('#blade_angle1').text(data["Blade Angle 1"].toFixed(2));
    else $('#blade_angle1').text('N/A');
    if(data["Blade Angle 2"] != -999999)
        $('#blade_angle2').text(data["Blade Angle 2"].toFixed(2));
    else $('#blade_angle2').text('N/A');
    if(data["Blade Angle 3"] != -999999)
        $('#blade_angle3').text(data["Blade Angle 3"].toFixed(2));
    else $('#blade_angle3').text('N/A');

    /*VOLT BATTERY BLADE PART*/
    if(data["Volt. Battery - blade 1"] != -999999)
        $('#volt_battery1').text(data["Volt. Battery - blade 1"].toFixed(2));
    else $('#volt_battery1').text('N/A');
    if(data["Volt. Battery - blade 2"] != -999999)
        $('#volt_battery2').text(data["Volt. Battery - blade 2"].toFixed(2));
    else $('#volt_battery2').text('N/A');
    if(data["Volt. Battery - blade 3"] != -999999)
        $('#volt_battery3').text(data["Volt. Battery - blade 3"].toFixed(2));
    else $('#volt_battery3').text('N/A');

    /*CURRENT PITCH MOTOR PART*/
    if(data["Current 1 Pitch Motor"] != -999999)
        $('#cur_pitch_motor1').text(data["Current 1 Pitch Motor"].toFixed(2));
    else $('#cur_pitch_motor1').text('N/A');
    if(data["Current 2 Pitch Motor"] != -999999)
        $('#cur_pitch_motor2').text(data["Current 2 Pitch Motor"].toFixed(2));
    else $('#cur_pitch_motor2').text('N/A');
    if(data["Current 3 Pitch Motor"] != -999999)
        $('#cur_pitch_motor3').text(data["Current 3 Pitch Motor"].toFixed(2));
    else $('#cur_pitch_motor3').text('N/A');

    /*PITCH MOTOR TEMPERATURE PART*/
    if(data["Pitch motor temperature - Blade 1"] != -999999)
        $('#pitch_motor_temp1').text(data["Pitch motor temperature - Blade 1"].toFixed(2));
    else $('#pitch_motor_temp1').text('N/A');
    if(data["Pitch motor temperature - Blade 2"] != -999999)
        $('#pitch_motor_temp2').text(data["Pitch motor temperature - Blade 2"].toFixed(2));
    else $('#pitch_motor_temp2').text('N/A');
    if(data["Pitch motor temperature - Blade 3"] != -999999)
        $('#pitch_motor_temp3').text(data["Pitch motor temperature - Blade 3"].toFixed(2));
    else $('#pitch_motor_temp3').text('N/A');

    /*PHASE VOLTAGE PART*/
    if(data["Phase 1 voltage"] != -999999)
        $('#phase_volt1').text(data["Phase 1 voltage"].toFixed(2));
    else $('#phase_volt1').text('N/A');
    if(data["Phase 2 voltage"] != -999999)
        $('#phase_volt2').text(data["Phase 2 voltage"].toFixed(2));
    else $('#phase_volt2').text('N/A');
    if(data["Phase 3 voltage"] != -999999)
        $('#phase_volt3').text(data["Phase 3 voltage"].toFixed(2));
    else $('#phase_volt3').text('N/A');

    /*PHASE CURRENT PART*/
    if(data["Phase 1 current"] != -999999)
        $('#phase_cur1').text(data["Phase 1 current"].toFixed(2));
    else $('#phase_cur1').text('N/A');
    if(data["Phase 2 current"] != -999999)
        $('#phase_cur2').text(data["Phase 2 current"].toFixed(2));
    else $('#phase_cur2').text('N/A');
    if(data["Phase 3 current"] != -999999)
        $('#phase_cur3').text(data["Phase 3 current"].toFixed(2));
    else $('#phase_cur3').text('N/A');

    /*POWER PART*/
    if(data["Power"] != -999999)
        $('#power').text(data["Power"].toFixed(2));
    else $('#power').text('N/A');
    if(data["Power Reactive"] != -999999)
        $('#power_react').text(data["Power Reactive"].toFixed(2));
    else $('#power_react').text('N/A');
    if(data["Freq. Grid"] != -999999)
        $('#freq_grid').text(data["Freq. Grid"].toFixed(2));
    else $('#freq_grid').text('N/A');

    /*DFIG PART*/
    if(data["DFIG active power"] != -999999)
        $('#dfig_act_power').text(data["DFIG active power"].toFixed(2));
    else $('#dfig_act_power').text('N/A');
    if(data["DFIG reactive power"] != -999999)
        $('#dfig_react_power').text(data["DFIG reactive power"].toFixed(2));
    else $('#dfig_react_power').text('N/A');
    if(data["DFIG mains Frequency"] != -999999)
        $('#dfig_main_freq').text(data["DFIG mains Frequency"].toFixed(2));
    else $('#dfig_main_freq').text('N/A');

    if(data["DFIG main voltage"] != -999999)
        $('#dfig_main_volt').text(data["DFIG main voltage"].toFixed(2));
    else $('#dfig_main_volt').text('N/A');
    if(data["DFIG main current"] != -999999)
        $('#dfig_main_cur').text(data["DFIG main current"].toFixed(2));
    else $('#dfig_main_cur').text('N/A');
    if(data["DFIG DC link voltage"] != -999999)
        $('#dfig_link_volt').text(data["DFIG DC link voltage"].toFixed(2));
    else $('#dfig_link_volt').text('N/A');

    /*ROTOR CURRENT PART*/
    if(data["Rotor R current"] != -999999)
        $('#rotor_r_cur').text(data["Rotor R current"].toFixed(2));
    else $('#rotor_r_cur').text('N/A');
    if(data["Roter Y current"] != -999999)
        $('#rotor_y_cur').text(data["Roter Y current"].toFixed(2));
    else $('#rotor_y_cur').text('N/A');
    if(data["Roter B current"] != -999999)
        $('#rotor_b_cur').text(data["Roter B current"].toFixed(2));
    else $('#rotor_b_cur').text('N/A');

    /*PRODUCTION PART*/
    if(data["Production"] != -999999)
        $('#production').text(data["Production"].toFixed(2));
    else $('#production').text('N/A');
    if(data["Cos Phi"] != -999999)
        $('#cos_phi').text(data["Cos Phi"].toFixed(2));
    else $('#cos_phi').text('N/A');

    /*TEMP GENERATOR PART*/
    if(data["Temp. generator 1 phase 1 coil"] != -999999)
        $('#temp_gen_coil1').text(data["Temp. generator 1 phase 1 coil"].toFixed(2));
    else $('#temp_gen_coil1').text('N/A');
    if(data["Temp. generator 1 phase 2 coil"] != -999999)
        $('#temp_gen_coil2').text(data["Temp. generator 1 phase 2 coil"].toFixed(2));
    else $('#temp_gen_coil2').text('N/A');
    if(data["Temp. generator 1 phase 3 coil"] != -999999)
        $('#temp_gen_coil3').text(data["Temp. generator 1 phase 3 coil"].toFixed(2));
    else $('#temp_gen_coil3').text('N/A');
    if(data["Temp. generator bearing driven End"] != -999999)
        $('#temp_gen_bearing_driven').text(data["Temp. generator bearing driven End"].toFixed(2));
    else $('#temp_gen_bearing_driven').text('N/A');
    if(data["Temp. generator bearing non-driven End"] != -999999)
        $('#temp_gen_bearing_non_driven').text(data["Temp. generator bearing non-driven End"].toFixed(2));
    else $('#temp_gen_bearing_non_driven').text('N/A');

    /*TEMP GEARBOX PART*/
    if(data["Temp. Gearbox driven end"] != -999999)
        $('#temp_gear_driven').text(data["Temp. Gearbox driven end"].toFixed(2));
    else $('#temp_gear_driven').text('N/A');
    if(data["Temp. Gearbox non-driven end"] != -999999)
        $('#temp_gear_non_driven').text(data["Temp. Gearbox non-driven end"].toFixed(2));
    else $('#temp_gear_non_driven').text('N/A');
    if(data["Temp. Gearbox inter. driven end"] != -999999)
        $('#temp_gear_inter_driven').text(data["Temp. Gearbox inter. driven end"].toFixed(2));
    else $('#temp_gear_inter_driven').text('N/A');
    if(data["Temp. Gearbox inter. non-driven end"] != -999999)
        $('#temp_gear_inter_non_driven').text(data["Temp. Gearbox inter. non-driven end"].toFixed(2));
    else $('#temp_gear_inter_non_driven').text('N/A');
    if(data["Pressure Gear box oil"] != -999999)
        $('#press_gear_oil').text(data["Pressure Gear box oil"].toFixed(2));
    else $('#press_gear_oil').text('N/A');
    if(data["Temp. Gear box oil"] != -999999)
        $('#temp_gear_oil').text(data["Temp. Gear box oil"].toFixed(2));
    else $('#temp_gear_oil').text('N/A');

    /*NACELLE PART*/
    if(data["Temp. Nacelle"] != -999999)
        $('#temp_nacelle').text(data["Temp. Nacelle"].toFixed(2));
    else $('#temp_nacelle').text('N/A');
    if(data["Temp. Ambient"] != -999999)
        $('#temp_ambient').text(data["Temp. Ambient"].toFixed(2));
    else $('#temp_ambient').text('N/A');
    if(data["Temp. Main bearing"] != -999999)
        $('#temp_main_bear').text(data["Temp. Main bearing"].toFixed(2));
    else $('#temp_main_bear').text('N/A');

    /*VIBRATION PART*/
    if(data["Damper Oscillation mag."] != -999999)
        $('#damper_osci_mag').text(data["Damper Oscillation mag."].toFixed(2));
    else $('#damper_osci_mag').text('N/A');
    if(data["Drive train vibration"] != -999999)
        $('#drive_train_vibra').text(data["Drive train vibration"].toFixed(2));
    else $('#drive_train_vibra').text('N/A');
    if(data["Tower vibration"] != -999999)
        $('#tower_vibra').text(data["Tower vibration"].toFixed(2));
    else $('#tower_vibra').text('N/A');
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