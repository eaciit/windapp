'use strict';


viewModel.XYanalysis = new Object();
var page = viewModel.XYanalysis;


vm.currentMenu('X/Y Analysis');
vm.currentTitle('X/Y Analysis');
vm.breadcrumb([{
    title: "Analysis Tool Box",
    href: '#'
},{
    title: 'X/Y Analysis',
    href: viewModel.appName + 'page/xyanalysis'
}]);


page.StateList = ko.observableArray([]);
page.ProjectList = ko.observableArray([]);
page.ProjectListState = ko.observableArray([]);
page.selectedProjectList = ko.observableArray([]);
page.TurbineList = ko.observableArray([]);
page.TurbineListByProject = ko.observableArray([]);
page.FieldList = ko.observableArray([]);
page.xAxis = ko.observableArray([]);
page.yAxis = ko.observableArray([]);
page.y2Axis = ko.observableArray([]);



page.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);


var lastPeriod = "";


page.SetValueFields = function(){
    setTimeout(function(){
        $("#xAxis").data("kendoDropDownList").select(0);
        $("#yAxis").data("kendoDropDownList").select(1);
        $("#y2Axis").data("kendoDropDownList").select(2);
        var date = new Date();
        date.setDate(date.getDate() - 1);
        $('#dateEnd').data("kendoDatePicker").value(date);

        var date2 = new Date();
        date2.setDate(date2.getDate() - 2);
        $('#dateStart').data("kendoDatePicker").value(date2);

    },300);
}

page.InitDefaultValue = function () {
    toolkit.ajaxPostDeffered(viewModel.appName + "xyanalysis/getxyfieldlist", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data;
        if(data !== null){
            setTimeout(function(){
                page.FieldList(data);
                page.xAxis(data);
                page.yAxis(data);
                page.y2Axis(data);

                $.when(page.SetValueFields()).done(function(){
                    page.checkKey();

                    $("#periodList").data("kendoDropDownList").value("custom");
                    $("#periodList").data("kendoDropDownList").trigger("change");
                    $("#StateList").data("kendoDropDownList").trigger("change");
                });
            },500)
        }   
    });

}

page.SetProjectList = function(){
    setTimeout(function(){
        var state = $("#StateList").data("kendoDropDownList").value();
        page.ProjectListState([]);

        $.each(page.ProjectList(), function(i , val){
            if(val.State == state){
                page.ProjectListState().push(val);
            }
        });
        page.selectedProjectList([page.ProjectListState()[0].ProjectId]);
        $("#ProjectList").data("kendoMultiSelect").setDataSource(page.ProjectListState());
        $('#ProjectList').data('kendoMultiSelect').value([page.selectedProjectList()]);
        $("#ProjectList").data("kendoMultiSelect").trigger("change");
    },300)

}

page.SetTurbineList = function(){
    page.TurbineListByProject([]);
    setTimeout(function(){
        var projectList = $("#ProjectList").data("kendoMultiSelect").value();

        $.each(projectList, function(i, val){
            $.each(page.TurbineList(), function(key, value){
                if(val == value.Project){
                    page.TurbineListByProject().push({label: value.Turbine, value: value.Value});
                }
            });
        });
        $("#turbineList").multiselect("dataprovider",page.TurbineListByProject());
        $('#turbineList').multiselect('select', page.TurbineListByProject()[0].value);
        $("#turbineList").multiselect("refresh");
    },300);
}

page.checkKey = function () {

    var key1 = $("#xAxis").data("kendoDropDownList").value();
    var key2 = $("#yAxis").data("kendoDropDownList").value();
    var key3 = $("#y2Axis").data("kendoDropDownList").value();

    page.xAxis([]);
    page.yAxis([]);
    page.y2Axis([]);

    var keys = page.FieldList();

    $.each(keys, function (i) {
        if (keys[i].Id == key2 || keys[i].Id == key3) {
            return true;
        }
        page.xAxis.push(keys[i]);
    });

    $.each(keys, function (i) {
        if (keys[i].Id == key1 || keys[i].Id == key3) {
            return true;
        }
        page.yAxis.push(keys[i]);
    });
    $.each(keys, function (i) {
        if (keys[i].Id == key1 || keys[i].Id == key2) {
            return true;
        }
        page.y2Axis.push(keys[i]);
    });

    $("#xAxis").data("kendoDropDownList").value(key1);
    $("#yAxis").data("kendoDropDownList").value(key2);
    $("#y2Axis").data("kendoDropDownList").value(key3);
}

page.showHidePeriod = function (callback) {
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

        $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
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

page.GetValueField = function(id){
    var value = $.grep(page.FieldList(), function(e){ return e.Id == $("#"+id).data("kendoDropDownList").value(); });
    return value[0];
}
page.LoadData = function(){
    app.loading(true);

    var url = viewModel.appName + 'xyanalysis/getdata';
    var param = {
        period : $('#periodList').data('kendoDropDownList').value(),
        project : $('#ProjectList').data('kendoMultiSelect').value(),
        engine : "",
        turbine :  $("#turbineList").val(),
        dateStart : $('#dateStart').data('kendoDatePicker').value(),
        dateEnd : $('#dateEnd').data('kendoDatePicker').value(),
        xAxis : page.GetValueField("xAxis"),
        y1Axis : page.GetValueField("yAxis"),
        y2Axis : page.GetValueField("y2Axis")
    }

    var getdata = toolkit.ajaxPostDeffered(url, param, function(res) {});
    $.when(getdata).done(function(d){
        page.GenerateChart(d.data);
        app.loading(false);
    });
}

page.GenerateChart = function(dataSource) {

    var dataSourceCustom = [];
    var byTurbine = [];
    $.each(dataSource.data, function(i, val){
        var datas = {
            turbine : val.turbine,
            details : [],
        }
        $.each(val.detail, function(index, value){
            var seriesData = {};
            $.each(dataSource.axisinfo, function(j, axisValue){
                seriesData[axisValue.Id] = value[j];
            });
            datas.details.push(seriesData);
            dataSourceCustom.push(seriesData);
        });
        byTurbine.push(datas);
    });


    var xAxis = {};
    var yAxes = [];
    var series = [];

    $.each(dataSource.axisinfo , function(i, value){
        var xField = dataSource.axisinfo[0].Id;
        var yAxis = dataSource.axisinfo[2].Id;
        if(i == 0){
            xAxis.axisCrossingValues =  [0, 10000];
            xAxis.title = {
                text: value.Text,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                color: "#585555",
                visible: true,
            };
            xAxis.majorGridLines = {
                            visible: true,
                            color: "#eee",
                            width: 0.8,
            };

        }
        else{
            var yAx = {
                name: value.Id, 
                title : {
                    text: value.Text,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555",
                    visible: true,
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            }

            var seriesData = {};

            if(i == 1){
                seriesData = {
                    name : value.Text, 
                    xField : xField,
                    yField : value.Id,
                    tooltip: {
                        format: value.Text + " : {0:N2}"
                    },
                }
            }else{
                seriesData = {
                    name : value.Text, 
                    xField : xField,
                    yField : value.Id, 
                    yAxis : yAxis,
                    tooltip: {
                        format: value.Text + " : {0:N2}"
                    },
                }
            }

            yAxes.push(yAx);
            series.push(seriesData);
        }
    });

    


    $('#chartxyAnalysis').html("");
    $("#chartxyAnalysis").kendoChart({
        theme: "flat",
        legend: {
            visible: false
        },
        title: {
            text: "X/Y Analysis",
            visible: false,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        dataSource: {
            data: dataSourceCustom
        },
        seriesDefaults: {
            type: "scatter",
            style: "smooth",
            markers:{size: 2}
        },
        series: series,
        xAxis: xAxis,
        yAxes: yAxes,
        tooltip: {
            visible: true
        }
    });
    $("#chartxyAnalysis").data("kendoChart").refresh();
}
$(function() {
    

    $("#ProjectList").kendoMultiSelect({
        dataSource: page.ProjectListState(), 
        dataValueField: 'ProjectId', 
        dataTextField: 'ProjectId', 
        change: function() {page.SetTurbineList()}, 
        suggest: true
    });

    $("#turbineList").multiselect({
        includeSelectAllOption: true,
        enableCaseInsensitiveFiltering: true,
        enableFiltering: true,
        maxHeight: 200,
        dropRight: false,
        onDropdownHide: function(event) {
            // fa.checkTurbine();
            // fa.currentFilter().turbine = this.$select.val();
            // fa.checkFilter();
        },
    });

    $('#btnRefresh').on('click', function() {
        setTimeout(function() {
            page.LoadData();
        }, 300);
    });


    app.loading(true);
    page.InitDefaultValue();
    setTimeout(function() {
        page.LoadData();
    }, 1500);

});