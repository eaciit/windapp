'use strict';


viewModel.MonitoringAlarm = new Object();
var ma = viewModel.MonitoringAlarm;
ma.minDatetemp = new Date();
ma.maxDatetemp = new Date();

ma.minDateRet = new Date();
ma.maxDateRet = new Date();


vm.currentMenu('Alarm Data');
vm.currentTitle('Alarm Data');
vm.isShowDataAvailability(false);
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Alarm Data', href: viewModel.appName + 'page/monitoringalarm' }]);
var intervalTurbine = null;

ma.UpdateProjectList = function(project) {
    setTimeout(function(){
        $('#projectList').data('kendoDropDownList').value(project);
        $("#projectList").data("kendoDropDownList").trigger("change");
    },1000)
}
ma.UpdateTurbineList = function(turbineList) {
    if(turbineList.length == 0){
         $("#turbineList").multiselect('selectAll', false).multiselect("refresh");
    }else{
        $('#turbineList').multiselect("deselectAll", false).multiselect("refresh");
        $('#turbineList').multiselect('select', turbineList); 
    }

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

        if(localStorage.getItem("projectname") !==  null && localStorage.getItem("turbine") !== null) {
            var locTurbine = localStorage.getItem("turbine");
            param.turbine = locTurbine == "" ? [] : [locTurbine];
            param.project = localStorage.getItem("projectname");

            var tabActive = localStorage.getItem("tabActive");

            $.when(ma.UpdateProjectList(param.project)).done(function () {
                setTimeout(function(){
                    ma.UpdateTurbineList(param.turbine);

                    setTimeout(function() {
                        ma.LoadDataAvail(param.project, gridType);
                        if(tabActive !== null){
                            if(tabActive == "alarmRaw" ){
                                $("#alarmrawTab a:first-child").trigger('click'); 
                            }else{
                                ma.CreateGridAlarm(gridType, param);
                            }
                        }
                    },500);
                app.resetLocalStorage();
                }, 1500);
            });
        } else {
            param.turbine = fa.turbine();
            param.project = fa.project;
            ma.CreateGridAlarm(gridType, param);
        }

    });
}
ma.buildParentFilter = function(filters, additionalFilter) {
    $.each(filters, function(idx, val){
        if(val.filter !== undefined) {
            ma.buildParentFilter(val.filter.filters, additionalFilter)
        }
        additionalFilter.push(val);
    });
}

ma.CreateGridAlarm = function(gridType, param) {
    var gridName = "#alarmGrid"
    var dt = new Date();
    var time = dt.getHours() + "" + dt.getMinutes() + "" + dt.getSeconds();
    var nameFile = "Monitoring Alarm Down_"+ moment(new Date()).format("Y-M-D")+"_"+time;
    var defaultsort = [ { field: "TimeStart", dir: "desc" }, { field: "TimeEnd", dir: "asc" } ]
    var url = viewModel.appName + "monitoringrealtime/getdataalarm";

    if(gridType == "warning") {
        gridName = "#warningGrid"
        nameFile = "Monitoring Alarm Warning";
    }

    if(gridType == "alarmraw"){
        gridName = "#alarmRawGrid"
        nameFile = "Monitoring Alarm Raw";
        defaultsort = [ { field: "TimeStamp", dir: "desc" } ]
    }

    var columns = [{
            field: "turbine",
            title: "Turbine",
            attributes: {
                style: "text-align:center;"
            },
            filterable: false,
            width: 90
        }, {
            field: "timestart",
            title: "Time Start",
            width: 170,
            filterable: false,
            attributes: {
                style: "text-align:center;"
            },
            template: "#= moment.utc(data.timestart).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timestart).format('HH:mm:ss')#"
        }, {
            field: "timeend",
            title: "Time End",
            width: 170,
            filterable: false,
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
            filterable: false,
            template: "#= time(data.duration) #"
        }, {
            field: "alarmcode",
            title: "Alarm Code",
            attributes: {
                style: "text-align:center;"
            },
            filterable: {
                ui: function(element) {
                    element.closest("form").find(".k-filter-help-text:first").remove();
                    element.closest("form").find("select").remove();
                }
            },
            width: 90,
        },{
            field: "alarmdesc",
            title: "Description",
            width: 330,
        }];

    if(gridType == "alarm"){
        columns.push({
            field: "reduceavailability",
            title: "Reduce Avail.",
            attributes: {
                style: "text-align:center;"
            },
            filterable: false,
            width: 90,
        });
    }

    if(gridType == "alarmraw"){
        columns = [{
            field: "turbine",
            title: "Turbine",
            attributes: {
                style: "text-align:center;"
            },
            filterable: false,
            width: 70
        }, {
            field: "timestamp",
            title: "Timestamp",
            width: 120,
            attributes: {
                style: "text-align:center;"
            },
            filterable: false,
            template: "#= moment.utc(data.timestamp).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timestamp).format('HH:mm:ss')#"
        }, {
            field: "tag",
            title: "Tag",
            width: 120,
            filterable: false,
            attributes: {
                style: "text-align:center;"
            },
        }, {
            field: "value",
            title: "Value",
            attributes: {
                style: "text-align:center;"
            },
            filterable: false,
            width: 70,
            // template: "#= kendo.toString(data.Timestamp,'n2') #"
        }, {
            field: "description",
            title: "Description",
            width: 200
        }, {
            field: "addinfo",
            title: "Note",
            filterable: false,
            width: 250
        }];
    }


    $(gridName).html('');
    $(gridName).kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            serverFiltering: true,
            transport: {
                read: {
                    url: url,
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                },
                parameterMap: function(options) {
                    var additionalFilter = [];
                    if(options.filter !== undefined) {
                        ma.buildParentFilter(options.filter.filters, additionalFilter)
                    }
                    if (additionalFilter.length > 0) {
                        options["filter"] = additionalFilter;
                    }

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
            sort: defaultsort,
        },
        // toolbar: ["excel"],
        excel: {
            fileName: nameFile+".xlsx",
            filterable: true,
            allPages: true
        },
        // pdf: {
        //     fileName: nameFile+".pdf",
        // },
        sortable: true,
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        filterable: {
            extra: false,
            operators: {
                string: {
                    contains: "Contains",
                    eq: "Is equal to"
                },
            }
        },
        columns: columns
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

                // $('#dateStart').data('kendoDatePicker').value( new Date(Date.UTC(moment( ma.maxDatetemp).get('year'),  ma.maxDatetemp.getMonth(),  ma.maxDatetemp.getDate() - 7, 0, 0, 0, 0)));
                // $('#dateEnd').data('kendoDatePicker').value(kendo.toString(moment.utc(res.data.Data[1]).format('DD-MMM-YYYY')));

                ma.checkCompleteDate();
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
        app.resetLocalStorage();
        var project =  $('#projectList').data('kendoDropDownList').value();
        localStorage.setItem('projectname', project);
        if(localStorage.getItem("projectname")){
            window.location = viewModel.appName + "page/monitoringbyproject";
        }
    },1500);
}

ma.exportToExcel = function(id){
    $("#"+id).getKendoGrid().saveAsExcel();
}

$(document).ready(function(){
    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        if($('.nav').find('li.active').find('a.tab-custom').text() == "Alarm Down") {
            ma.CreateGrid("alarm");
        } else if($('.nav').find('li.active').find('a.tab-custom').text() == "Alarm Warning") {
            ma.CreateGrid("warning");
        }else{
            ma.CreateGrid("alarmraw");
        }
    });

    //setTimeout(function() {
        $.when(ma.InitDateValue()).done(function () { 
            setTimeout(function() {
                ma.CreateGrid("alarm");
                ma.LoadDataAvail(fa.project, "alarm");
            }, 100);
        });
    //}, 300);
    
    $('#projectList').kendoDropDownList({
        change: function () {  
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
            ma.LoadDataAvail(project, "alarm");
        }
    });
});
