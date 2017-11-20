'use strict';

viewModel.FilterScatter = new Object();
var fa = viewModel.FilterScatter;

fa.turbineList = ko.observableArray([]);
fa.projectList = ko.observableArray([]);
fa.engineList = ko.observableArray([]);


fa.showEngine = ko.observable(false);
fa.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Last 24 hours" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
fa.periodType = ko.observable();
fa.dateStart = ko.observable();
fa.dateEnd = ko.observable();
fa.turbine = ko.observable();
fa.project = ko.observable();
fa.period = ko.observable();
fa.engine = ko.observable("");
fa.infoPeriodRange = ko.observable();
fa.infoPeriodIcon = ko.observable(false);
fa.rawproject = ko.observableArray([]);
fa.rawturbine = ko.observableArray([]);

var lastPeriod = "";

fa.populateTurbine = function (selected, engine) {
    if (fa.rawturbine().length == 0) {
        fa.turbineList([{ value: "", label: "" }]);
    } else {
        var datavalue = [];  

        $.each(fa.rawturbine(), function (key, val) {
            if(engine !== undefined){
                if (selected == "" && val.Engine == engine) {
                    var data = {};
                    data.value = val.Value;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }else if (selected == val.Project && val.Engine == engine){
                    var data = {};
                    data.value = val.Value;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }
            }else{
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
            }

        });

        var data = datavalue.sort(function(a, b){
            var a1= a.label.toLowerCase(), b1= b.label.toLowerCase();
            if(a1== b1) return 0;
            return a1> b1? 1: -1;
        });

        fa.turbineList(data);
        setTimeout(function () {
            $('#turbineList').data('kendoDropDownList').select(0);
        }, 100);
    }
};
fa.populateEngine = function(selected){
    var list = [];
    $.each(fa.rawproject(), function(i, val){
        if(val.ProjectId == selected){
          if(val.Engine.length > 0){
            $.each(val.Engine, function(id, engine){
                var data = {text : engine, value : engine};
                list.push(data);
            })
          }
        }
    });

    fa.engineList(list);

    setTimeout(function () {
        if(fa.engineList().length > 0){
            fa.showEngine(true);
            $('#engineList').kendoDropDownList({
                data: fa.engineList(),
                dataValueField: 'value',
                dataTextField: 'text',
                suggest: true,
                change: function () { 
                    fa.populateTurbine(selected,this._old);
                }
            });

            $("#engineList").data("kendoDropDownList").select(0);               
            fa.engine = $("#engineList").data("kendoDropDownList").value();
            fa.populateTurbine(selected, fa.engine);
            
        }else{
            fa.engine = "";
            fa.showEngine(false);
            fa.populateTurbine(selected);
        }
    }, 500);
}


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

fa.showHidePeriod = function (callback) {
    var period = $('#periodList').data('kendoDropDownList').value();

    var date = new Date(); 
    
    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    // var endMonthDate;

    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
    var startDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-7, 0, 0, 0, 0));

    if(moment(new Date()).get('year') ==  moment(maxDateData).get('year') && new Date(date).getMonthName() == new Date(maxDateData).getMonthName()){
        var endMonthDate = new Date(app.getDateMax(maxDateData));
    }else{
        var endMonthDate = new Date(date.getFullYear(), date.getMonth() + 1, 0);
    }


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

        fa.dateStart = startDate;
        fa.dateEnd = endMonthDate;

        $('#dateStart').data('kendoDatePicker').value(startDate);
        $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
        
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


            fa.dateStart = startMonthDate;
            fa.dateEnd = endMonthDate;

            $('#dateStart').data('kendoDatePicker').value(fa.dateStart);
            $('#dateEnd').data('kendoDatePicker').value(fa.dateEnd);

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
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();

    if (fa.dateStart - fa.dateEnd > 25200000) {
        toolkit.showError("Invalid Date Range Selection");
        return false;
    } else {
        fa.InitFilter();
        fa.checkCompleteDate();
        var period = $('#periodList').data('kendoDropDownList').value();
        return true;
    }
}

fa.InitFilter = function () {
    fa.dateStart = $('#dateStart').data('kendoDatePicker').value();
    fa.dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    fa.dateStart = new Date(Date.UTC(fa.dateStart.getFullYear(), fa.dateStart.getMonth(), fa.dateStart.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(fa.dateEnd.getFullYear(), fa.dateEnd.getMonth(), fa.dateEnd.getDate(), 0, 0, 0));
    fa.project = $("#projectList").data("kendoDropDownList").value();
    fa.period = $("#periodList").data("kendoDropDownList").value();
    fa.isDownTime = $("#isDownTime").is(":checked");
    fa.turbine = $("#turbineList").data("kendoDropDownList").value();
    fa.periodType = $("#periodList").data("kendoDropDownList").value();
}

fa.InitDefaultValue = function () {
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-7, 0, 0, 0, 0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
}

fa.DateChange = function () {
    var period = $('#periodList').data('kendoDropDownList').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var endMonthDate;

    if(period == "monthly" && new Date(dateEnd).getMonthName() == new Date(maxDateData).getMonthName()){
         $('#dateEnd').data('kendoDatePicker').value(new Date(app.getDateMax(maxDateData)))
    }else{
        $('#dateEnd').data('kendoDatePicker').value(new Date(dateEnd.getFullYear(), dateEnd.getMonth() + 1, 0));
    }

    var start = $('#dateStart').data('kendoDatePicker').value();
    var end = $('#dateEnd').data('kendoDatePicker').value();

    fa.dateStart = new Date(Date.UTC(start.getFullYear(), start.getMonth(), start.getDate(), 0, 0, 0));
    fa.dateEnd = new Date(Date.UTC(end.getFullYear(), end.getMonth(), end.getDate(), 0, 0, 0));
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

fa.changeEndDate = function(){
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEndMax = $('#dateEnd').data('kendoDatePicker').max();

    var dateEndPicker  = $('#dateEnd').data('kendoDatePicker');

    
    setTimeout(function(){
        if(moment(dateStart).format('MMM-Y') == moment(new Date()).format('MMM-Y')){
            dateEndPicker.value(new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateEndMax.getDate(), 0, 0, 0)));
        }else{
            dateEndPicker.value(new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth()+1, 0, 0, 0, 0)));
        }
    },200);
    
    dateEndPicker.min(new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0)));

}

fa.setProjectTurbine = function(projects, turbines, selected){
	fa.rawproject(projects);
    fa.rawturbine(turbines);
	fa.populateProject(selected);
};

$(document).ready(function () {
    app.loading(true);
    $('#projectList').kendoDropDownList({
        change: function () { 
            var project = $('#projectList').data("kendoDropDownList").value();
            di.getAvailDate();
            fa.populateEngine(project);
         }
    });
    fa.showHidePeriod();
    fa.InitDefaultValue();
});
