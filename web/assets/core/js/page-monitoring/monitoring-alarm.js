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

ma.UpdateProjectList = function(project) {
    $('#projectList').data('kendoDropDownList').value(project);
    $("#projectList").data("kendoDropDownList").trigger("change");
}
ma.UpdateTurbineList = function(turbineList) {
    $('#turbineList').multiselect("deselectAll", false).multiselect("refresh");
    $('#turbineList').multiselect('select', turbineList);
}
ma.CreateGrid = function(gridType) {
    app.loading(true);
    $.when(fa.LoadData()).done(function () {
        var COOKIES = {};
        var cookieStr = document.cookie;
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: [],
            project: "",
            tipe: gridType,
        };
        
        if(cookieStr.indexOf("turbine=") >= 0 && cookieStr.indexOf("projectname=") >= 0) {
            document.cookie = "projectname=;expires=Thu, 01 Jan 1970 00:00:00 UTC;";
            document.cookie = "turbine=;expires=Thu, 01 Jan 1970 00:00:00 UTC;";
            cookieStr.split(/; /).forEach(function(keyValuePair) {
                var cookieName = keyValuePair.replace(/=.*$/, "");
                var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
                COOKIES[cookieName] = cookieValue;
            });
            param.turbine = [COOKIES["turbine"]];
            param.project = COOKIES["projectname"];

            $.when(ma.UpdateProjectList(param.project)).done(function () {
                setTimeout(function(){
                    $.when(ma.UpdateTurbineList(param.turbine)).done(function () {
                        ma.LoadDataAvail(param.project, gridType);
                        ma.CreateGridAlarm(gridType, param);
                    });
                }, 700);
            });

        } else {
            param.turbine = fa.turbine();
            param.project = fa.project;
            ma.LoadDataAvail(param.project, gridType);
            ma.CreateGridAlarm(gridType, param);
        }
    });
}
ma.CreateGridAlarm = function(gridType, param) {
    var gridName = "#alarmGrid"
    if(gridType == "warning") {
        gridName = "#warningGrid"
    }
    $(gridName).html('');
    $(gridName).kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
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
                    if (!app.isFine(res)) {
                        return;
                    }
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
            attributes: {
                style: "text-align:center;"
            },
            width: 90
        }, {
            field: "TimeStart",
            title: "Time Start",
            width: 170,
            attributes: {
                style: "text-align:center;"
            },
            template: "#= moment.utc(data.TimeStart).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.TimeStart).format('HH:mm:ss')#"
        }, {
            field: "TimeEnd",
            title: "Time End",
            width: 170,
            attributes: {
                style: "text-align:center;"
            },
            template: "#= (moment.utc(data.TimeEnd).format('DD-MM-YYYY') == '01-01-0001'?'Not yet finished' : (moment.utc(data.TimeEnd).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.TimeEnd).format('HH:mm:ss')))#"
        }, {
            field: "Duration",
            title: "Duration (hh:mm:ss)",
             width: 120,
             attributes: {
                style: "text-align:center;"
            },
            template: "#= time(data.Duration) #"
        }, {
            field: "AlarmCode",
            title: "Alarm Code",
            attributes: {
                style: "text-align:center;"
            },
            width: 90,
        }, {
            field: "AlarmDesc",
            title: "Description",
            width: 330
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
}

ma.LoadDataAvail = function(projectname, gridType){
    //fa.LoadData();

    var payload = {
        project: projectname,
        tipe: gridType
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

ma.ToByProject = function(){
    setTimeout(function(){
        app.loading(true);
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        var project =  $('#projectList').data('kendoDropDownList').value();
        document.cookie = "project="+project.split("(")[0].trim()+";expires="+ newDateObj;
        window.location = viewModel.appName + "page/monitoringbyproject";
    },1500);
}

$(document).ready(function(){
    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        if($('.nav').find('li.active').find('a.tab-custom').text() == "Alarm Down") {
            ma.CreateGrid("alarm");
        } else {
            ma.CreateGrid("warning");
        }
    });

    setTimeout(function() {
        $.when(ma.InitDateValue()).done(function () {
            setTimeout(function() {
                ma.CreateGrid("alarm");
            }, 100);
        });
    }, 300);
    
    $('#projectList').kendoDropDownList({
        change: function () {  
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
        }
    });
});
