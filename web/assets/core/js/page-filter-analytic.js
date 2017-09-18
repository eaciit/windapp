'use strict';

viewModel.FilterAnalytic = new Object();
var fa = viewModel.FilterAnalytic;

fa.turbineList = ko.observableArray([]);
fa.projectList = ko.observableArray([]);

fa.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
fa.periodType = ko.observable();
fa.dateStart = ko.observable();
fa.dateEnd = ko.observable();
fa.turbine = ko.observableArray([]);
fa.rawturbine = ko.observableArray([]);
fa.project = ko.observable();
fa.rawproject = ko.observableArray([]);
fa.period = ko.observable();
fa.infoPeriodRange = ko.observable();
fa.infoPeriodIcon = ko.observable(false);
fa.infoFiltersChanged = ko.observable(false);
fa.textFilterChanged = ko.observable("<strong>Filters have been changed !</strong><br>Previousfilter : <br>- Tejuva <br>- All Turbines <br>- Custom | 24 Apr 2017 - 30 Apr 2017");
fa.previousFilter = ko.observable(
    {
        "project":"",
        "turbine":[],
        "period":"",
        "startDate":"",
        "endDate":""
    }
);
fa.currentFilter = ko.observable(
    {
        "project":"",
        "turbine":[],
        "period":"",
        "startDate":"",
        "endDate":""
    }
);


var lastPeriod = "";
var turbineval = [];

/*fa.InitFirst = function () {
    $.when(
        app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                fa.turbineList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                if (res.data.length > 0) {
                    var allturbine = {}
                    $.each(res.data, function (key, val) {
                        turbineval.push(val);
                    });
                    allturbine.value = "All Turbine";
                    allturbine.text = "All Turbines";
                    datavalue.push(allturbine);
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                    });
                }
                fa.turbineList(datavalue);
            }
        }),
        app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                fa.projectList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                if (res.data.length > 0) {
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                    });
                }
                fa.projectList(datavalue);
            }
        })

    ).then(function () {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"])
        // override to set the value
        $("#projectList").data("kendoDropDownList").value("Tejuva");
        fa.project = $("#projectList").data("kendoDropDownList").value();
    });
}*/

/*fa.populateTurbine = function () {
    app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.length == 0) {
            res.data = [];;
            fa.turbineList([{ value: "", text: "" }]);
        } else {
            var datavalue = [];
            if (res.data.length > 0) {
                var allturbine = {}
                $.each(res.data, function (key, val) {
                    turbineval.push(val);
                });
                allturbine.value = "All Turbine";
                allturbine.text = "All Turbines";
                datavalue.push(allturbine);
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    datavalue.push(data);
                });
            }
            fa.turbineList(datavalue);
        }

        setTimeout(function () {
            $('#turbineList').data('kendoMultiSelect').value(["All Turbine"])
        }, 300);
    });
};

fa.populateProject = function () {
    app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.length == 0) {
            res.data = [];;
            fa.projectList([{ value: "", text: "" }]);
        } else {
            var datavalue = [];
            if (res.data.length > 0) {
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    datavalue.push(data);
                });
            }
            fa.projectList(datavalue);

            // override to set the value
            setTimeout(function () {
                $("#projectList").data("kendoDropDownList").value("Tejuva");
                fa.project = $("#projectList").data("kendoDropDownList").value();
            }, 300);
        }
    });
};*/

fa.populateTurbine = function (selected) {
    if (fa.rawturbine().length == 0) {
        fa.turbineList([{ value: "", label: "" }]);
    } else {
        var datavalue = [];        
        // $.each(fa.rawturbine(), function (key, val) {
        //     turbineval.push(val);
        // });
        var allturbine = {}
        // allturbine.value = "multiselect-all";
        // allturbine.label = "multiselect-all";
        // datavalue.push(allturbine);

        $.each(fa.rawturbine(), function (key, val) {
            if (selected == "") {
                var data = {};
                data.value = val.Value;
                data.label = val.Turbine;
                datavalue.push(data);
            }else if (selected == val.Project){
                var data = {};
                data.value = val.Value;
                data.label = val.Turbine;
                datavalue.push(data);
            }
        });

        var data = datavalue.sort(function(a, b){
            var a1= a.label.toLowerCase(), b1= b.label.toLowerCase();
            if(a1== b1) return 0;
            return a1> b1? 1: -1;
        });

        fa.turbineList(data);
    }

    setTimeout(function () {
         fa.setTurbine();
    }, 100);
};

fa.populateProject = function (selected) {
    if (fa.rawproject().length == 0) {
        fa.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];        
        $.each(fa.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Name;
            datavalue.push(data);
        });
        fa.projectList(datavalue);

        // override to set the value
        
        setTimeout(function () {
            if (selected != "") {
                $("#projectList").data("kendoDropDownList").value(selected);
            } else {
                $("#projectList").data("kendoDropDownList").select(0);
            }               
            fa.project = $("#projectList").data("kendoDropDownList").value();
            fa.populateTurbine(fa.project);
        }, 100);
    }
};

fa.getProjectInfo = function () {
    var project = $("#projectList").data("kendoDropDownList").value();
    // var turbines = $('#turbineList').data('kendoMultiSelect').value();
     var turbines = $('#turbineList').val();

    if (turbines[0] == "multiselect-all") {
        turbines = [];
    }

    var param = {
        Project: project,
        Turbines: turbines
    }

    app.ajaxPost(viewModel.appName + "helper/getprojectinfo", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        $("#project-info").html($("#projectList").data("kendoDropDownList").value());
        $("#total-turbine-info").html('<i class="fa fa-flash tooltipster tooltipstered" aria-hidden="true" title="Total Turbine"></i>&nbsp;' + res.data.TotalTurbine);
        $("#total-capacity-info").html('<i class="fa fa-tachometer tooltipster tooltipstered" aria-hidden="true" title="Total Capacity"></i>&nbsp;' + res.data.TotalCapacity + "MW");
    });
};

fa.showHidePeriod = function (callback) {
    var period = $('#periodList').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.getDateMax(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));


    if (period == "custom") {
        $(".show_hide").show();
        $('#dateStart').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#dateEnd').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        // if (lastPeriod == "monthly") {
        //     $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        //     $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
        // } else if (lastPeriod == "annual") {
        //     $('#dateStart').data('kendoDatePicker').value(startYearDate);
        //     $('#dateEnd').data('kendoDatePicker').value(endYearDate);
        // }
        $('#dateStart').data('kendoDatePicker').value(fa.currentFilter().startDate !== "" ?fa.currentFilter().startDate : startMonthDate);
        $('#dateEnd').data('kendoDatePicker').value(fa.currentFilter().endMonthDate !== "" ?fa.currentFilter().endDate : endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM-yyyy",
            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM-yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startMonthDate);
            $('#dateEnd').data('kendoDatePicker').value(endMonthDate);

            $(".show_hide").show();
        } else if (period == "annual") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",

            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startYearDate);
            $('#dateEnd').data('kendoDatePicker').value(endYearDate);

            $(".show_hide").show();
        } else {
            if (period == 'last24hours') {
                $('#dateStart').data('kendoDatePicker').value(last24hours);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            } else if (period == 'last7days') {
                $('#dateStart').data('kendoDatePicker').value(lastweek);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            }
            $(".show_hide").hide();
        }
        lastPeriod = period;
    }

    setTimeout(function () {
        callback;
    }, 50);
}

fa.LoadData = function () {
    if ($("#turbineList").val() == "") {
        $('#turbineList').val("multiselect-all")
    }

    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    if (fa.dateStart - fa.dateEnd > 25200000) {
        toolkit.showError("Invalid Date Range Selection");
        return false;
    } else {
        fa.InitFilter();
        fa.checkCompleteDate();
        return true;
    }
}

fa.checkTurbine = function () {
    var arr = $('#turbineList').val();
    // var index = arr.indexOf("multiselect-all");
    // if (index == 0 && arr.length > 1) {
    //     arr.splice(index, 1);
    //     // $('#turbineList').data('kendoMultiSelect').value(arr)
    //     $('#turbineList').val(arr);
    // } else if (index > 0 && arr.length > 1) {
    //     // $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    //     $('#turbineList').val("multiselect-all")
    // } else if (arr.length == 0) {
    //     // $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    //     $('#turbineList').val("multiselect-all")
    // }
    
    if(arr == null){
        var $el = $("#turbineList");
        $('option', $el).each(function(element) {
          $el.multiselect('select', $(this).val());
        });
        arr = $('#turbineList').val();
    }
    fa.turbine(arr);
}

fa.InitFilter = function () {
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    fa.project = $("#projectList").data("kendoDropDownList").value();
    fa.period = $("#periodList").data("kendoDropDownList").value();
    fa.isDownTime = $("#isDownTime").is(":checked");

    // if ($("#turbineList").val().indexOf("All Turbine") >= 0) {
    //     fa.turbine = [];
    // } else {
    fa.turbine($("#turbineList").val());
    // }

    fa.periodType = $("#periodList").data("kendoDropDownList").value();

    fa.GetBreakDown();
}

fa.InitDefaultValue = function () {
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-7, 0, 0, 0, 0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
    setTimeout(function(){
        $("#turbineList").multiselect({
            includeSelectAllOption: true,
            enableCaseInsensitiveFiltering: true,
            enableFiltering: true,
            maxHeight: 200,
            dropRight: false,
            onDropdownHide: function(event) {
                fa.checkTurbine();
                fa.currentFilter().turbine = this.$select.val();
                fa.checkFilter();
            },
        });
        $("#turbineList").multiselect("dataprovider",fa.turbineList());
        fa.checkTurbine();
        $.when(fa.setPreviousFilter(true)).done(function(){
            var prevFilter = fa.previousFilter();
            
            $.each(prevFilter, function(key,val){
                fa.currentFilter()[key] = val
            }); 
        });
    },200);
}

fa.GetBreakDown = function () {
    fa.periodType = $("#periodList").data("kendoDropDownList").value();
    fa.project = $("#projectList").data("kendoDropDownList").value();
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    fa.dateStart = new Date(Date.UTC(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), fa.dateStart.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), 0, 0, 0));

    var result = [];
    var monthyearStart = 0;
    var monthyearEnd = 0;

    monthyearStart = fa.dateStart.getFullYear() + fa.dateStart.getMonth();
    monthyearEnd = fa.dateEnd.getFullYear() + fa.dateEnd.getMonth();

    if (fa.periodType == "last24hours" || fa.periodType == "last7days") {
        result.push({ "value": "Date", "text": "Date" });
    } else if (fa.periodType == "monthly" || fa.periodType == "annual") {
        result.push({ "value": "Month", "text": "Month" });
        result.push({ "value": "Year", "text": "Year" });
    } else if (fa.periodType == "custom") {
        if (fa.dateStart.getDate() == 1){
            if (monthyearStart == monthyearEnd){
                result.push({ "value": "Date", "text": "Date" });
            }
        }else{
            if ((fa.dateEnd - fa.dateStart) / 86400000 + 1 <= 31) {
                result.push({ "value": "Date", "text": "Date" });
            }
        }
        
        result.push({ "value": "Month", "text": "Month" });
        result.push({ "value": "Year", "text": "Year" });
    }

    if (fa.project == "") {
        result.push({ "value": "Project", "text": "Project" });
    } else {
        result.push({ "value": "Turbine", "text": "Turbine" });
    }

    return result;
}

fa.DateChange = function () {
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    fa.dateStart = new Date(Date.UTC(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), fa.dateStart.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), 0, 0, 0));
}

fa.checkCompleteDate = function () {
    var period = $('#periodList').data('kendoDropDownList').value();

    var monthNames = moment.months();

    var currentDateData = moment(app.currentDateData).format('YYYY-MM-DD');
    var today = moment().format('YYYY-MM-DD');
    var thisMonth = moment().get('month');
    var firstDayMonth = moment(new Date(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), 1)).format("YYYY-MM-DD");
    var lastDayMonth = moment(new Date(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth() + 1, 0)).format("YYYY-MM-DD");
    var firstDayYear = moment().startOf('year').format('YYYY-MM-DD');
    var endDayYear = moment().endOf('year').format('YYYY-MM-DD');

    var dateStart = moment(fa.dateStart).format('YYYY-MM-DD');
    var dateEnd = moment(fa.dateEnd).format('YYYY-MM-DD');

    if (period === 'custom') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range on start date and end date");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodRange("* Incomplete period data range on start date");
            fa.infoPeriodIconmozilla(true);
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range on end date");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'annual') {
        if ((moment(fa.dateEnd).get('year') == moment(app.currentDateData).get('year')) && (currentDateData < today)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("* Incomplete period data range in end year");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else if (period === 'monthly') {
        if ((dateEnd > currentDateData) && (dateStart > currentDateData)) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in start month and start month");
        } else if (dateStart > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in start month");
        } else if (dateEnd > currentDateData) {
            fa.infoPeriodIcon(true);
            fa.infoPeriodRange("*Incomplete period data range in end month");
        } else {
            fa.infoPeriodIcon(false);
            fa.infoPeriodRange("");
        }
    } else {
        fa.infoPeriodRange("");
        fa.infoPeriodIcon(false);
    }


}

fa.disableRefreshButton = function(param){
    if(param == true){
        $("#btnRefresh").attr("disabled",true);
    }else{
        $("#btnRefresh").removeAttr("disabled");
    }
}

fa.setTurbine = function(){
    setTimeout(function(){
        $("#turbineList").multiselect("dataprovider",fa.turbineList());
        fa.checkTurbine();
    },200);
}
fa.setProjectTurbine = function(projects, turbines, selected){
	fa.rawproject(projects);
    fa.rawturbine(turbines);
	fa.populateProject(selected);
};

fa.setPreviousFilter = function(isFirst){
    fa.infoFiltersChanged(false);

    var data = {
        project: $("#projectList").data("kendoDropDownList").value(), 
        turbine: $("#turbineList").val(), 
        period: $('#periodList').data('kendoDropDownList').value(), 
        startDate: kendo.toString(new Date($('#dateStart').data('kendoDatePicker').value()), "dd-MMM-yyyy"),
        endDate:  kendo.toString(new Date($('#dateEnd').data('kendoDatePicker').value()), "dd-MMM-yyyy"),
    }

    fa.previousFilter(data);

    if(isFirst == null){
        fa.currentFilter(data);
    }

    if($("#btnRefresh").hasClass("tooltipstered") == true){
        $(".filter-changed").tooltipster('destroy');
    }

}
fa.resetFilter = function(){
    fa.infoFiltersChanged(false);

    var data = fa.previousFilter();
    fa.populateTurbine(data.project);

    setTimeout(function(){
        $("#projectList").data("kendoDropDownList").value(data.project); 

        $("#turbineList").val("");
        $("#turbineList").multiselect("select",data.turbine);
        $("#turbineList").multiselect("refresh");

        $('#periodList').data('kendoDropDownList').value(data.period);
        fa.showHidePeriod();
        $('#dateStart').data('kendoDatePicker').value(data.startDate);
        $('#dateEnd').data('kendoDatePicker').value(data.endDate);

        $(".filter-changed").tooltipster('destroy');
    },500);
}

fa.checkFilter = function(){
    // fa.infoFiltersChanged(false);
    // $.each(fa.currentFilter(), function(key, val){
    //     if(key == "turbine"){
    //         var diff = [];
    //         if(fa.previousFilter().turbine.length > fa.currentFilter().turbine.length){
    //             $.grep(fa.previousFilter().turbine, function(el) {
    //                     if ($.inArray(el, fa.currentFilter().turbine) == -1) diff.push(el);
    //             });  
    //         }else{
    //             $.grep(fa.currentFilter().turbine, function(el) {
    //                     if ($.inArray(el, fa.previousFilter().turbine) == -1) diff.push(el);
    //             }); 
    //         }

    //         if(diff.length > 0){
    //             fa.infoFiltersChanged(true);
    //         }
    //     }else{
    //         if(fa.previousFilter()[key] !== val){
    //             fa.infoFiltersChanged(true);
    //         }
    //     }
    // });

    // if(fa.infoFiltersChanged() == true){
    //      fa.setTextFilterChanged();
    // }else{
    //     if($("#btnRefresh").hasClass("tooltipstered") == true){
    //         $(".filter-changed").tooltipster('destroy');
    //     }
    // }

    return false;
   

}

fa.setTextFilterChanged = function(){

    var textTurbine = "";
    var data = fa.previousFilter();
    if(data.turbine.length <= 4){
        textTurbine = data.turbine;
    }else{
        textTurbine = data.turbine.length == fa.turbineList().length ? "All Turbines" :  data.turbine.length+" Turbines Selected";
    }

    if($("#btnRefresh").hasClass("tooltipstered") == true){
        $(".filter-changed").tooltipster('destroy');
    }
    
    $("#btnRefresh").addClass("flash-button");

    $("#btnRefresh").tooltipster({
        theme: 'tooltipster-val',
        animation: 'grow',
        delay: 0,
        offsetY: -5,
        touchDevices: false,
        trigger: 'hover',
        interactive : true,
        position: "bottom",
        content: "<strong>Filters have been changed !</strong><br>Previous filter : <br>- "+data.project+" <br>- "+textTurbine+" <br>- "+data.period+" | "+data.startDate+ " - "+data.endDate +"<br><span class='pull-right btn btn-danger btn-xs' onClick='fa.resetFilter()'><i class='fa fa-times'></i> Reset</span>",
        multiple: false,
        contentAsHTML : true,
    });
}


fa.getDataAvailability = function(){

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value(); 
    var param = {
        period: fa.period,
        dateStart: new Date(moment(dateStart).format('YYYY-MM-DD')),
        dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
        turbine: fa.turbine(),
        project: fa.project,
    };

    toolkit.ajaxPost(viewModel.appName + "dataavailability/getcurrentdataavailability", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        vm.projectName(param.project);
        vm.dataAvailability(kendo.toString((res.data * 100), 'n2') + " %");
    })
}
$(document).ready(function () {
    app.loading(true);
    fa.showHidePeriod();
    fa.InitDefaultValue();
    setTimeout(function(){
        $("#dateStart").attr("readonly", true);
        $("#dateEnd").attr("readonly", true);
        $(".multiselect-native-select").find(".btn-group").find(".multiselect-filter").find(".input-group").addClass("input-group-sm");
        $(".multiselect-native-select").find(".btn-group").find(".multiselect-filter").find(".input-group").find(".input-group-addon").remove();
        $(".multiselect-native-select").find(".btn-group").find(".multiselect-filter").find(".input-group").find(".input-group-btn").remove();
    },1000);
    
});
