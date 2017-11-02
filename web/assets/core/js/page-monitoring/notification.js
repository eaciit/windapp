'use strict';


viewModel.MonitoringNotification = new Object();
var mn = viewModel.MonitoringNotification;
mn.minDatetemp = new Date();
mn.maxDatetemp = new Date();

mn.minDateRet = new Date();
mn.maxDateRet = new Date();


vm.currentMenu('Notification');
vm.currentTitle('Notification');
vm.isShowDataAvailability(false);
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Notification', href: viewModel.appName + 'page/monitoringnotification' }]);

var intervalTurbine = null;

mn.typeList = ko.observableArray([]);
mn.typeListData = ko.observableArray([]);

mn.UpdateProjectList = function(project) {
    setTimeout(function(){
        $('#projectList').data('kendoDropDownList').value(project);
        $("#projectList").data("kendoDropDownList").trigger("change");
    },1000)
}
mn.UpdateTurbineList = function(turbineList) {
    if(turbineList.length == 0){
         $("#turbineList").multiselect('selectAll', false).multiselect("refresh");
    }else{
        $('#turbineList').multiselect("deselectAll", false).multiselect("refresh");
        $('#turbineList').multiselect('select', turbineList); 
    }

}
mn.LoadData = function() {
    app.loading(true);

    $.when(fa.LoadData()).done(function () {
        var project = fa.project
        if(project == "") {
            project = "Tejuva";
        }
        mn.typeList(mn.typeListData()[project]);
        var colnameType = $('#typeList').data('kendoDropDownList').value();

        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine(),
            project: fa.project,
            tipe: colnameType,
        };

        mn.LoadDataAvail(param.project, "alarm");
        mn.generateGrid(param);
    });
}


mn.generateGrid = function(param){
    var nameFile = "Monitoring Notification_"+ moment(new Date()).format("Y-M-D")+"_"+time;

    $("#notificationGrid").html('');
    $("#notificationGrid").kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "monitoringrealtime/getdatanotification",
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

                    mn.minDateRet = new Date(res.data.mindate);
                    mn.maxDateRet = new Date(res.data.maxdate);

                    mn.checkCompleteDate()

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
            sort: [ { field: "TimeStart", dir: "desc" }, { field: "TimeEnd", dir: "asc" } ],
        },
        sortable: true,
        toolbar: ["excel"],
        excel: {
            fileName: nameFile+".xlsx",
            filterable: true,
            allPages: true
        },
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        columns: [{
            field: "turbine",
            title: "Turbine",
            attributes: {
                style: "text-align:center;"
            },
            width: 90
        }, {
            field: "timestart",
            title: "Time Start",
            width: 130,
            attributes: {
                style: "text-align:center;"
            },
            template: "#= moment.utc(data.timestart).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timestart).format('HH:mm:ss')#"
        }, {
            field: "timeend",
            title: "Time End",
            width: 130,
            attributes: {
                style: "text-align:center;"
            },
            template: "#= (moment.utc(data.timeend).format('DD-MM-YYYY') == '01-01-0001'?'Not yet finished' : (moment.utc(data.timeend).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timeend).format('HH:mm:ss')))#"
        }, {
            field: "duration",
            title: "Duration (hh:mm:ss)",
             width: 120,
             attributes: {
                style: "text-align:center;"
            },
            template: "#= time(data.duration) #"
        }, {
            field: "startcond",
            title: "Start Value",
            attributes: {
                style: "text-align:center;"
            },
            width: 130,
        }, {
            field: "endcond",
            title: "End Value",
            attributes: {
                style: "text-align:center;"
            },
            width: 130,
        },{
            field: "description",
            title: "Description",
            width: 190,
            attributes: {
                style: "text-align:center;"
            },
        }]
    });

}


function time(s) {
    return new Date(s * 1e3).toISOString().slice(-13, -5);
}

mn.InitDateValue = function () {
    var maxDateData = new Date();

    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0, 0));
    var lastEndDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0, 0));

    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
}

mn.LoadDataAvail = function(projectname, gridType){
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
                mn.minDatetemp = new Date(res.data.Data[0]);
                mn.maxDatetemp = new Date(res.data.Data[1]);
                app.currentDateData = new Date(res.data.Data[1]);

                $('#availabledatestart').html(kendo.toString(moment.utc(mn.minDatetemp).format('DD-MMMM-YYYY')));
                $('#availabledateend').html(kendo.toString(moment.utc(mn.maxDatetemp).format('DD-MMMM-YYYY')));

                mn.checkCompleteDate()
            }
        }         
    });
}

mn.checkCompleteDate = function () {
    var currentDateData = moment.utc(mn.maxDatetemp).format('YYYY-MM-DD');
    var prevDateData = moment.utc(mn.minDatetemp).format('YYYY-MM-DD');

    var dateStart = moment.utc(mn.minDateRet).format('YYYY-MM-DD');
    var dateEnd = moment.utc(mn.maxDateRet).format('YYYY-MM-DD');

    if ((dateEnd > currentDateData) || (dateStart > currentDateData)) {
        fa.infoPeriodIcon(true);
    } else if ((dateEnd < prevDateData) || (dateStart < prevDateData)) {
        fa.infoPeriodIcon(true);
    } else {
        fa.infoPeriodIcon(false);
        fa.infoPeriodRange("");
    }
}

// mn.ToByProject = function(){    
//     setTimeout(function(){
//         app.loading(true);
//         app.resetLocalStorage();
//         var project =  $('#projectList').data('kendoDropDownList').value();
//         localStorage.setItem('projectname', project);
//         if(localStorage.getItem("projectname")){
//             window.location = viewModel.appName + "page/monitoringbyproject";
//         }
//     },1500);
// }

$(document).ready(function(){
    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        mn.LoadData();
    });

    $.when(mn.InitDateValue()).done(function () { 
        setTimeout(function() {
            mn.LoadData();
        }, 100);
    });
    
    setTimeout(function () {
        $("#typeList").data("kendoDropDownList").value("alltypes");
        var dropdownlist = $("#typeList").data("kendoDropDownList");
        dropdownlist.list.width("auto");
    }, 300);

    $('#projectList').kendoDropDownList({
        change: function () {  
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
            mn.typeList(mn.typeListData()[project]);
        }
    });

    $('#typeList').on("change", function() {
        fa.checkTurbine();
        setTimeout(function() {
            mn.LoadData();
        }, 300);
    });
});
