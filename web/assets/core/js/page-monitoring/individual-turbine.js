'use strict';


viewModel.IndividualTurbine = new Object();
var it = viewModel.IndividualTurbine;


vm.currentMenu('Individual Turbine');
vm.currentTitle('Individual Turbine');
vm.isShowDataAvailability(false);
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Individual Turbine', href: viewModel.appName + 'page/monitoringindividualturbine' }]);
var intervalTurbine = null;
var chart;

var maxSamples = 50,count = 0;

it.projectList = ko.observableArray([]);
it.project = ko.observable();
it.turbineList = ko.observableArray([{}]);
it.allTurbineList = ko.observableArray([{}]);
it.isFromSummary = ko.observable(false);
it.isFromByProject = ko.observable(false);
it.isShowVibration = ko.observable(false);
it.feeder = ko.observableArray();
it.colorStatus = ko.observable();



// it.getFeeder = function(project, turbine){
//     setTimeout(function(){
//         var turbineList = it.allTurbineList;
//         if(turbineList.length > 0) {
//             $.each(it.allTurbineList, function(key, val){
//                 if(val.Project == project && val.Turbine == turbine){
//                     it.feeder(val.Feeder);
//                 }
//             });
//         }
//     },500);
// }

it.checkVisible = function(){
    if(it.project() == 'Amba' || it.project() == 'Sattigeri' || it.project() == 'Nimbagallu' || it.project() == 'Taralkatti' ){
        return true;
    }else{
        return false;
    }
}


it.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        it.projectList([{ value: "", text: "" }]);
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
        it.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(0);
            it.project($("#projectList").data("kendoDropDownList").value());
        }, 300);
    }
};

it.populateTurbine = function (project, turbine, isChange) {
    if (turbine.length == 0) {
        turbine = [];
        it.turbineList([{ Id: "", Lat: 0.0, Lon: 0.0 }]);
    } else {
        var turbinevalue = [];
        if (turbine.length > 0) {
            it.allTurbineList = turbine;
            $.each(turbine, function (key, val) {
                if(val.Project == project) {
                    var data = {};
                    data.value = val.Value;
                    data.label = val.Turbine;
                    // data.Id = val.Turbine;
                    // data.Lat = val.latitude;
                    // data.Lon = val.longitude;
                    turbinevalue.push(data);
                }
            });
        }

        var data = turbinevalue.sort(function(a, b){
            var a1= a.label.toLowerCase(), b1= b.label.toLowerCase();
            if(a1== b1) return 0;
            return a1> b1? 1: -1;
        });

        it.turbineList(data);

        if(isChange) {
            setTimeout(function () {
                $("#turbine").data("kendoDropDownList").select(0);
            }, 100);
        }
    }
};

it.ChangeProject = function() {
    return function() {
        it.project($("#projectList").data("kendoDropDownList").value());
        var projects = [];
        projects.push(it.project());
        
        it.populateTurbine(projects, it.allTurbineList, true);
        setTimeout(function(){
            it.isFirst(true);
            it.ShowData();
            it.isShowVibration(it.project() == "Lahori" ? true : false); 
        },300);
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

      date = new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0], timeParts[0], timeParts[1], timeParts[2]);

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
        var wsVal = parseFloat(res.data["Wind speed Avg"].toFixed(2));
        var pwrVal = parseFloat(res.data["Power"].toFixed(2));

        if(it.isFirst() == false){
            if (wsVal == -999999 ) {
                chart.series[0].addPoint([time, null], true, chart.series[0].data.length>maxSamples ? true:false);
            }else{
                chart.series[0].addPoint([time, wsVal], true, chart.series[0].data.length>maxSamples ? true:false);
            }

            if (pwrVal == -999999 ) {
                chart.series[1].addPoint([time, null], true, chart.series[0].data.length>maxSamples ? true:false);
            }else{
                chart.series[1].addPoint([time, pwrVal], true, chart.series[0].data.length>maxSamples ? true:false);
            }
        }else{

            setTimeout(function(){
                if (wsVal == -999999 ) {
                    it.dataWindspeed([time, null]);
                }else{
                    it.dataWindspeed([time, wsVal]);
                }

                if (pwrVal == -999999 ) {
                    it.dataPower([time, null]);
                }else{
                    it.dataPower([time, pwrVal]);
                }

                it.showWindspeedLiveChart();
            },500);
        }


        it.PlotData(res.data);
        it.isFirst(false);
    });
};

it.ShowRemark = function(){
    TbCol.ResetData();
    var turbineName = $('#turbine').data('kendoDropDownList').text();
    var turbineId = $('#turbine').data('kendoDropDownList').value();
    var result = $.grep(turbineList, function(e){ return e.Project == it.project() && e.Value == turbineId})[0];


    TbCol.TurbineId(turbineId);
    TbCol.TurbineName(turbineName);
    TbCol.UserId('');
    TbCol.UserName('');
    TbCol.Project(it.project());
    TbCol.Feeder(result.Feeder);
    TbCol.IsTurbine(true);
    TbCol.OpenForm();
    TbCol.IconStatus(it.colorStatus());
}

// it.rotor_r_cur = ko.observable('');
// it.rotor_y_cur = ko.observable('');
// it.rotor_b_cur = ko.observable('');
// it.windspeed1 = ko.observable('');
// it.windspeed2 = ko.observable('');
// it.wind_dir1 = ko.observable('');
// it.wind_dir2 = ko.observable('');
// it.speed_gen = ko.observable('');
// it.dfig_act_power = ko.observable('');
// it.dfig_react_power = ko.observable('');
// it.dfig_main_freq = ko.observable('');
// it.dfig_main_volt = ko.observable('');
// it.dfig_main_cur = ko.observable('');
// it.dfig_link_volt = ko.observable('');
// it.damper_osci_mag = ko.observable('');
// it.press_gear_oil = ko.observable('');
// it.tower_vibra = ko.observable('');

it.rotor_rpm = ko.observable('');
it.windspeed_avg = ko.observable('');
it.wind_dir = ko.observable('');
it.nacel_dir = ko.observable('');
it.gen_rpm = ko.observable('');
it.pitch_angle = ko.observable('');
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
it.pitch_conv_tempblade1 = ko.observableArray('');
it.pitch_conv_tempblade2 = ko.observableArray('');
it.pitch_conv_tempblade3 = ko.observableArray('');
it.phase_volt1 = ko.observable('');
it.phase_volt2 = ko.observable('');
it.phase_volt3 = ko.observable('');
it.phase_cur1 = ko.observable('');
it.phase_cur2 = ko.observable('');
it.phase_cur3 = ko.observable('');
it.power = ko.observable('');
it.power_react = ko.observable('');
it.freq_grid = ko.observable('');
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
it.temp_nacelle = ko.observable('');
it.temp_ambient = ko.observable('');
it.temp_main_bear = ko.observable('');
it.temp_gear_oil = ko.observable('');
it.drive_train_vibra = ko.observable('');
it.transformer_winding_temp1 = ko.observable('');
it.transformer_winding_temp2 = ko.observable('');
it.transformer_winding_temp3 = ko.observable('');
it.transformerWindingTemp1 = ko.observable('');
it.transformerWindingTemp2 = ko.observable('');
it.transformerWindingTemp3 = ko.observable('');
it.temp_slip_ring = ko.observable('');
it.hydraulic_pressure = ko.observable('');
it.hydraulic_temp = ko.observable('');
it.AccXDir = ko.observable('');
it.AccYDir = ko.observable('');

it.isFirst = ko.observable(true);

it.grid_voltage = ko.observable('');
it.grid_current = ko.observable('');
it.rotor_current = ko.observable('');
it.stator_current = ko.observable('');
it.radiator_temp1 = ko.observable('');
it.radiator_temp2 = ko.observable('');
it.dataWindspeed = ko.observableArray([]);
it.dataPower = ko.observableArray([]);

it.PlotData = function(data) {

    var lastUpdate = moment.utc(data["lastupdate"]);
    $('#turbine_last_update').text(lastUpdate.format("DD MMM YYYY HH:mm:ss"));

    if (data["Turbine Status"] === -999) {
        $('#rotorPic').css({"fill": "#b8b9bb"});
        it.colorStatus("txt-grey");
    } else if (data["Turbine Status"] == 1) {
        $('#rotorPic').css({"fill": "#31b445"});
        it.colorStatus("txt-green");
    } else {
        $('#rotorPic').css({"fill": "#db1e1e"});
        it.colorStatus("txt-red");
    }

    /*WIND SPEED PART*/
    if(data["Wind speed Avg"] != -999999)
        it.windspeed_avg(data["Wind speed Avg"].toFixed(2));
    else it.windspeed_avg('N/A');
    /*if(data["Wind speed 1"] != -999999)
        it.windspeed1(data["Wind speed 1"].toFixed(2));
    else it.windspeed1('N/A');
    if(data["Wind speed 2"] != -999999)
        it.windspeed2(data["Wind speed 2"].toFixed(2));
    else it.windspeed2('N/A');*/

    it.showWindspeedColumnChart();

    /*WIND DIRECTION PART*/
    if(data["Wind Direction"] != -999999)
        it.wind_dir(data["Wind Direction"].toFixed(2));
    else it.wind_dir('N/A');
    /*if(data["Vane 1 wind direction"] != -999999)
        it.wind_dir1(data["Vane 1 wind direction"].toFixed(2));
    else it.wind_dir1('N/A');
    if(data["Vane 2 wind direction"] != -999999)
        it.wind_dir2(data["Vane 2 wind direction"].toFixed(2));
    else it.wind_dir2('N/A');*/

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
    // if(data["DFIG speed generator encoder"] != -999999)
    //     it.speed_gen(data["DFIG speed generator encoder"].toFixed(2));
    // else it.speed_gen('N/A');

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

    if(data["Pitch Angle"] != -999999)
        it.pitch_angle("PITCH SYSTEM ("+data["Pitch Angle"].toFixed(2)+")");
    else it.pitch_angle('PITCH SYSTEM');

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

    /*PITCH CONV INTERNAL TEMPERATURE PART*/
    if(data["Pitch Conv Internal Temp Blade1"] != -999999)
        it.pitch_conv_tempblade1(data["Pitch Conv Internal Temp Blade1"].toFixed(2));
    else it.pitch_conv_tempblade1('N/A');

    if(data["Pitch Conv Internal Temp Blade2"] != -999999)
        it.pitch_conv_tempblade2(data["Pitch Conv Internal Temp Blade2"].toFixed(2));
    else it.pitch_conv_tempblade2('N/A');
    
    if(data["Pitch Conv Internal Temp Blade3"] != -999999)
        it.pitch_conv_tempblade3(data["Pitch Conv Internal Temp Blade3"].toFixed(2));
    else it.pitch_conv_tempblade3('N/A');

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
    // if(data["DFIG active power"] != -999999)
    //     it.dfig_act_power(data["DFIG active power"].toFixed(2));
    // else it.dfig_act_power('N/A');
    // if(data["DFIG reactive power"] != -999999)
    //     it.dfig_react_power(data["DFIG reactive power"].toFixed(2));
    // else it.dfig_react_power('N/A');
    // if(data["DFIG mains Frequency"] != -999999)
    //     it.dfig_main_freq(data["DFIG mains Frequency"].toFixed(2));
    // else it.dfig_main_freq('N/A');

    // if(data["DFIG main voltage"] != -999999)
    //     it.dfig_main_volt(data["DFIG main voltage"].toFixed(2));
    // else it.dfig_main_volt('N/A');
    // if(data["DFIG main current"] != -999999)
    //     it.dfig_main_cur(data["DFIG main current"].toFixed(2));
    // else it.dfig_main_cur('N/A');
    // if(data["DFIG DC link voltage"] != -999999)
    //     it.dfig_link_volt(data["DFIG DC link voltage"].toFixed(2));
    // else it.dfig_link_volt('N/A');

    /*ROTOR CURRENT PART*/
    // if(data["Rotor R current"] != -999999) {
    //     it.rotor_r_cur(data["Rotor R current"].toFixed(2));
    // }
    // else { 
    //     it.rotor_r_cur('N/A');
    // }
    // if(data["Roter Y current"] != -999999) {
    //     it.rotor_y_cur(data["Roter Y current"].toFixed(2));
    // }
    // else { 
    //     it.rotor_y_cur('N/A');
    // }
    // if(data["Roter B current"] != -999999)
    //     it.rotor_b_cur(data["Roter B current"].toFixed(2));
    // else it.rotor_b_cur('N/A');


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



    /*TEMP WINDING PART*/
    if(data["Transformer Winding Temp1"] != -999999)
        it.transformer_winding_temp1(data["Transformer Winding Temp1"].toFixed(2));
    else it.transformer_winding_temp1('N/A');
    if(data["Transformer Winding Temp2"] != -999999)
        it.transformer_winding_temp2(data["Transformer Winding Temp2"].toFixed(2));
    else it.transformer_winding_temp2('N/A');
    if(data["Transformer Winding Temp3"] != -999999)
        it.transformer_winding_temp3(data["Transformer Winding Temp3"].toFixed(2));
    else it.transformer_winding_temp3('N/A');
    if(data["Temp Slip Ring"] != -999999)
        it.temp_slip_ring(data["Temp Slip Ring"].toFixed(2));
    else it.temp_slip_ring('N/A');
    if(data["Hydraulic Pressure"] != -999999)
        it.hydraulic_pressure(data["Hydraulic Pressure"].toFixed(2));
    else it.hydraulic_pressure('N/A');
    if(data["Hydraulic Temp"] != -999999)
        it.hydraulic_temp(data["Hydraulic Temp"].toFixed(2));
    else it.hydraulic_temp('N/A');





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
    // if(data["Pressure Gear box oil"] != -999999)
    //     it.press_gear_oil(data["Pressure Gear box oil"].toFixed(2));
    // else it.press_gear_oil('N/A');
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
    // if(data["Damper Oscillation mag."] != -999999)
    //     it.damper_osci_mag(data["Damper Oscillation mag."].toFixed(2));
    // else it.damper_osci_mag('N/A');
    if(data["Drive train vibration"] != -999999)
        it.drive_train_vibra(data["Drive train vibration"].toFixed(2));
    else it.drive_train_vibra('N/A');

    if(data["AccXDir"] != -999999)
        it.AccXDir(data["AccXDir"].toFixed(2));
    else it.AccXDir('N/A');

    if(data["AccYDir"] != -999999)
        it.AccYDir(data["AccYDir"].toFixed(2));
    else it.AccYDir('N/A');

    /* AMBA PART */
    if(data["Grid Current"] != -999999)
        it.grid_current(data["Grid Current"].toFixed(2));
    else it.grid_current('N/A');
    if(data["Grid Voltage"] != -999999)
        it.grid_voltage(data["Grid Voltage"].toFixed(2));
    else it.grid_voltage('N/A');
    if(data["Rotor Current"] != -999999)
        it.rotor_current(data["Rotor Current"].toFixed(2));
    else it.rotor_current('N/A');
    if(data["Stator Current"] != -999999)
        it.stator_current(data["Stator Current"].toFixed(2));
    else it.stator_current('N/A');
    if(data["Ref Radiator Temp1"] != -999999)
        it.radiator_temp1(data["Ref Radiator Temp1"].toFixed(2));
    else it.radiator_temp1('N/A');
    if(data["Ref Radiator Temp2"] != -999999)
        it.radiator_temp2(data["Ref Radiator Temp2"].toFixed(2));
    else it.radiator_temp2('N/A');


    // if(data["Tower vibration"] != -999999)
    //     it.tower_vibra(data["Tower vibration"].toFixed(2));
    // else it.tower_vibra('N/A');

    if(data["isRemark"] == true){
        $(".icon-remark").show();
    }else{
         $(".icon-remark").hide();
    }

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

    var turbine = "";
    var project = "";

    if(localStorage.getItem("projectname") !== null && localStorage.getItem("turbine") !== null) {
        turbine = localStorage.getItem("turbine")
        project = localStorage.getItem("projectname");

        if(localStorage.getItem("isFromSummary") !== undefined && localStorage.getItem("isFromSummary") == "true"){
            it.isFromSummary(true);
        }

        if(localStorage.getItem("isFromByProject") !== undefined && localStorage.getItem("isFromByProject") == "true"){
            it.isFromByProject(true);
       }

      
        setTimeout(function(){
            $('#projectList').data('kendoDropDownList').value(project);
            var change = $("#projectList").data("kendoDropDownList").trigger("change");
            setTimeout(function(){
                $('#turbine').data('kendoDropDownList').value(turbine);
            },200);
            app.resetLocalStorage();
        },500);
        setTimeout(function(){
            it.isFirst(true);
            it.ShowData();
        },300);

    } else {
        turbine = $('#turbine').data('kendoDropDownList').value();
        project = $('#projectList').data('kendoDropDownList').value();
        it.GetData(project, turbine);
    }

    
    $.when(it.showWindRoseChart()).done(function () {
        setTimeout(function() {
            app.loading(false);
        }, 1500);
    });


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
            minorUnit: (2100/4)/2,
            max: 2100,
            // min : -200,
            reverse: false,
            majorUnit: 2100/4,
            ranges: [
                {
                    from: 1525,
                    to: 2100,
                    color : "#8dcb2a"
                },
                {
                    from: 950,
                    to: 1525,
                    color: "#ffc700"
                }, {
                    from: 375,
                    to: 950,
                    color: "#ff7a00"
                }, {
                    from: -200,
                    to: 375,
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
            width: 390,
        },
        credits: {
              enabled: false
        },
        legend: {
            enabled: true,
            verticalAlign: 'top',
            layout: "horizontal",
            labelFormatter: function() {
                
                if(this.point.y == undefined){
                     return '<span style="color:' + this.color + '"> ' + this.name + ' </span> : <span style="min-width:50px"><b>  -  </b> '+this.tooltipOptions.valueSuffix+'</n>';
                }
                else{
                    return '<span style="color:'+ this.color +'"> ' + this.name + ' </span> : <span style="min-width:50px"><b> '+ kendo.toString(this.point.y,'n2')+' </b></span> <b>'+this.tooltipOptions.valueSuffix+'</b><br/>'
                }
               
            },
            // labelFormat: '<span style="color:{color}">{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
        },
        rangeSelector: {
            enabled:false,
        },
        scrollbar: {
            enabled: false
        },
        exporting: {
            enabled: false
        },
        // yAxis:{
        //   labels:
        //   {
        //     enabled: false
        //   },
        //   gridLineWidth: 1,
        //   minorGridLineWidth: 0,
        //   opposite:false
        // },
        yAxis : [{ // Primary yAxis
            gridLineWidth: 1,
            minorGridLineWidth: 0,
            labels: {
                enabled: false
            },
            title: {
                text: 'Wind Speed',
                enabled: false
            },
            height: '50%',
            // opposite: true

        }, { // Secondary yAxis
            gridLineWidth: 0,
            minorGridLineWidth: 0,
            title: {
                text: 'Power',
                enabled: false
            },
            labels: {
                enabled: false
            },
            top: '55%',
            height: '45%',
            offset: 0,

        }],
        xAxis: {
           type: 'datetime',
           dateTimeLabelFormats : {
                millisecond: '%H:%M:%S',
                second: '%H:%M:%S',
                minute: '%H:%M',
                hour: '%H:%M',
                day: '%e. %b',
                week: '%e. %b',
                month: '%b \'%y',
                year: '%Y'
           },
           lineWidth: 1,
           minorGridLineWidth: 1,
           lineColor: 'transparent',
           labels: {
               enabled: true
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
                    enabled: false,
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
            yAxis: 0,
            tooltip: {
                valueSuffix: 'm/s'
            },
        },{
            name: 'Power',
            color: colorField[1],
            data: [it.dataPower()],
            yAxis: 1,
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
            var windRoseData = res.data["WindRose"][0].Data;

            $("#windRoseChart").kendoChart({
                theme: "flat",
                chartArea: {
                    height: 215,
                    width: 200,
                    margin: 0,
                    marginLeft: -30,
                    padding: 0,
                },
                dataSource: {
                    data: windRoseData,
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
                tooltip: {
                    visible: true,
                    template: "#= category #"+String.fromCharCode(176)+" (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                }
            });
        }
    })
}

it.ToTimeSeriesHfd = function() {
    setTimeout(function(){
        app.loading(true);
        var turbine = $("#turbine").val();
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        var project =  $('#projectList').data('kendoDropDownList').value();
        document.cookie = "project="+project.split("(")[0].trim()+";expires="+ newDateObj;
        document.cookie = "turbine="+turbine+"; expires="+ newDateObj;
        window.location = viewModel.appName + "page/timeserieshfd";
    },1500);
}

it.ToByProject = function(){    
    setTimeout(function(){
        app.loading(true);
        app.resetLocalStorage();
        var project =  $('#projectList').data('kendoDropDownList').value();
        localStorage.setItem('projectname', project);
        if(localStorage.getItem("projectname")){
            window.location = viewModel.appName + "page/monitoringbyproject";
        }
    },1500);
}

it.ToSummary = function(){
    window.location = viewModel.appName + "page/monitoringsummary";
}


it.backToProject = function(){
    if(it.isFromSummary()){
        it.ToSummary();
    }

    if(it.isFromByProject()){
        it.ToByProject();
    }
}


it.ToAlarm = function() {
    setTimeout(function(){
        app.loading(true);
        app.resetLocalStorage();
        var turbine = $("#turbine").val();
        var project =  $('#projectList').data('kendoDropDownList').value();
        localStorage.setItem('turbine', turbine == [] ? null : turbine);
        localStorage.setItem('projectname', project);
        localStorage.setItem('tabActive', "alarmRaw");
        if(localStorage.getItem("turbine") !== null && localStorage.getItem("projectname")){
            window.location = viewModel.appName + "page/monitoringalarm";
        }
    },1500);
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
    setTimeout(function() {
        it.ShowData();
        intervalTurbine = window.setInterval(it.ShowData, 6000);
    }, 600);
});
