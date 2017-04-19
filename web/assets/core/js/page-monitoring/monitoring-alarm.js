'use strict';


viewModel.MonitoringAlarm = new Object();
var ma = viewModel.MonitoringAlarm;
ma.minDatetemp = new Date();
ma.maxDatetemp = new Date();

ma.minDateRet = new Date();
ma.maxDateRet = new Date();


vm.currentMenu('Alarm Data');
vm.currentTitle('Alarm Data');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Alarm Data', href: viewModel.appName + 'page/monitoringalarm' }]);
var intervalTurbine = null;

ma.CreateGrid = function() {
    app.loading(true);

    fa.LoadData();

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
        $('#turbineList').data('kendoMultiSelect').value([turbine]);
        $('#projectList').data('kendoDropDownList').value(project);
        console.log("tet");
    } else {
        turbine = $('#turbineList').data('kendoMultiSelect').value();
        project = $('#projectList').data('kendoDropDownList').value();
        console.log("yoy");
    }


    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: turbine,
        project: project
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
                    var totalFreq = res.data.Total;
                    var totalHour = res.data.Duration;

                    ma.minDateRet = new Date(res.data.mindate);
                    ma.maxDateRet = new Date(res.data.maxdate);

                    ma.checkCompleteDate()

                    $('#alarm_duration').text((totalHour/3600).toFixed(2));
                    $('#alarm_frequency').text(totalFreq);

                    setTimeout(function(){
                        app.loading(false);
                    }, 300)
                    
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

ma.InitDateValue = function () {
    var maxDateData = new Date();

    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0, 0));
    var lastEndDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0, 0));

    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);

    ma.LoadDataAvail()
}

ma.LoadDataAvail = function(){
    fa.LoadData();
    var payload = {
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataalarmavaildate", payload, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.Data.length == 0) {
            res.data.Data = [];
        } else {
            if (res.data.Data.length > 0) {
                ma.minDatetemp = new Date(res.data.Data[0]);
                ma.maxDatetemp = new Date(res.data.Data[1]);
                app.currentDateData = new Date(res.data.Data[1]);

                $('#availabledatestart').html(kendo.toString(moment.utc(ma.minDatetemp).format('DD-MMMM-YYYY')));
                $('#availabledateend').html(kendo.toString(moment.utc(ma.maxDatetemp).format('DD-MMMM-YYYY')));

                ma.checkCompleteDate()
            }
        }         
    });
}

ma.checkCompleteDate = function () {
    var currentDateData = moment.utc(ma.maxDatetemp).format('YYYY-MM-DD');
    var prevDateData = moment.utc(ma.minDatetemp).format('YYYY-MM-DD');

    var dateStart = moment.utc(ma.minDateRet).format('YYYY-MM-DD');
    var dateEnd = moment.utc(ma.maxDateRet).format('YYYY-MM-DD');

    if ((dateEnd > currentDateData) || (dateStart > currentDateData)) {
        fa.infoPeriodIcon(true);
    } else if ((dateEnd < prevDateData) || (dateStart < prevDateData)) {
        fa.infoPeriodIcon(true);
    } else {
        fa.infoPeriodIcon(false);
        fa.infoPeriodRange("");
    }
}

$(document).ready(function(){
    $('#btnRefresh').on('click', function () {
        ma.CreateGrid();
    });

    setTimeout(function() {
        $.when(ma.InitDateValue()).done(function () {
            setTimeout(function() {
                ma.CreateGrid();
            }, 100);
        });
    }, 300);
});