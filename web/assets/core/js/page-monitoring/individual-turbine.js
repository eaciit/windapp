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
        it.PlotData(res);
    });
};

it.PlotData = function(data) {
    if(data.ActivePower!=-999999)
        $('#turbine_generation').text(data.ActivePower.toFixed(2));
    else $('#turbine_generation').text('N/A');
    if(data.ActivePower!=-999999)
        $('#turbine_production').text((data.ActivePower/6).toFixed(2));
    else $('#turbine_production').text('N/A');
    if(data.ActivePower!=-999999)
        $('#turbine_plf').text((data.ActivePower/2100*100).toFixed(2));
    else $('#turbine_plf').text('N/A');
    if(data.WindSpeed!=-999999)
        $('#turbine_wind_speed').text(data.WindSpeed.toFixed(2));
    else $('#turbine_wind_speed').text('N/A');
    if(data.NacellePosition!=-999999)
        $('#turbine_nacelle_position').text(data.NacellePosition.toFixed(2));
    else $('#turbine_nacelle_position').text('N/A');
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
    var turbine = $('#turbine').data('kendoDropDownList').value();
    it.LoadData(turbine);
    it.GetData(turbine);
};


$(document).ready(function(){
    it.ShowData();
    intervalTurbine = window.setInterval(it.ShowData, 6000);
});