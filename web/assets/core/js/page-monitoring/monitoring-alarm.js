'use strict';


viewModel.MonitoringAlarm = new Object();
var ma = viewModel.MonitoringAlarm;


vm.currentMenu('Alarm Data');
vm.currentTitle('Alarm Data');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Alarm Data', href: viewModel.appName + 'page/monitoringalarm' }]);
var intervalTurbine = null;

ma.CreateGrid = function() {
    app.loading(true);

    var param = {
    };

    $('#alarmGrid').html('');
    $('#alarmGrid').kendoGrid({
        dataSource: {
            serverPaging: true,
            transport: {
                read: {
                    url: viewModel.appName + "monitoringrealtime/getdataalarm",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            schema: {
                data: function data(res) {
                    app.loading(false);

                    var totalFreq = 0;
                    var totalHour = 0;
                    $.each(res.data.Data, function(idx, val) {
                        totalFreq++;
                        totalHour += val.Duration;
                    });
                    $('#alarm_duration').text((totalHour/3600).toFixed(2));
                    $('#alarm_frequency').text(totalFreq);
                    
                    return res.data.Data;
                },
                total: function data(res) {
                    return res.data.Total;
                }
            },
            pageSize: 10,
            sort: [
                { field: "TimeStart", dir: "desc" },
                { field: "TimeEnd", dir: "asc" }
            ],
        },
        sortable: true,
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        columns: [{
            field: "Turbine",
            title: "Turbine",
            width: 160
        }, {
            field: "TimeStart",
            title: "Time Start",
            template: "#= moment.utc(data.TimeStart).format('DD-MMM-YYYY HH:mm:ss') #"
        }, {
            field: "TimeEnd",
            title: "Time End",
            template: "#= (moment.utc(data.TimeEnd).format('DD-MM-YYYY')=='01-01-0001'?'Not yet finished':moment.utc(data.TimeEnd).format('DD-MMM-YYYY HH:mm:ss')) #"
        }, {
            field: "Duration",
            title: "Duration (hh:mm:ss)",
            template: "#= time(data.Duration) #"
        }, {
            field: "AlarmCode",
            title: "Alarm Code",
            width: 60
        }, {
            field: "AlarmDesc",
            title: "Description",
            width: 240
        }]
    });
};

function time(s) {
    return new Date(s * 1e3).toISOString().slice(-13, -5);
}


$(document).ready(function(){
    ma.CreateGrid();
    // ma.LoadData();
});