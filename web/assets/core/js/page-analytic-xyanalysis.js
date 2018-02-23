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
page.Data = ko.observableArray([]);



page.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);


var lastPeriod = "";
var colors = ["#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"];
var colorsDeg = ["#FFD6AD", "#A6E7DF", "#FFC8C0", "#FFE2B8", "#D9F2BA", "#A4D8E7", "#FFC0DB", "#FAB3AE", "#C3EDF5", "#CFC8DC", "#D6A0E0", "#A8E6CC", "#F5B9BD", "#E7D8B5", "#FFBBD5", "#E7A89D", "#EDC7BE", "#FFA9EF", "#ADDDD0", "#9FE0F7", "#99B7C9", "#FF99AF", "#B9CADA", "#FFC1C1", "#FFEEC1", "#C6DDFF", "#C9BBB5", "#AFADC6", "#C3B5D4", "#E5E7BA", "#DDBCA3", "#FBDFF3", "#CAADA1", "#99ABF8", "#D1C7A3", "#A5CF9B", "#FFD699", "#D7A8DF", "#C4CBD0", "#EFB1A4", "#BDF5D9", "#F099ED", "#DFDBC5", "#CBADEA", "#D9F6FB", "#D7E2FA", "#D8D8D8", "#9BBC9B", "#9AB2CD", "#A2D3CB", "#AAC9C9", "#DBCA9D", "#A7B3B3", "#A3DEC9", "#C9E89F", "#C7C5BF", "#B8C7A5", "#C2ADD7", "#A4FAE1", "#F2AA9C", "#EFEBDD", "#C5B8B9", "#FAA2CC", "#C9C9F7", "#FBDDF5", "#B1E1DE", "#BE9BDF", "#A1ECE3", "#D9BDBD", "#A5BFB0", "#B79CCC", "#C0D5BF", "#A4D2AB", "#DFE99D", "#C1A8D4", "#F39B9E", "#A6A7F9"];

page.hexToRgb = function(hex) {
    // Expand shorthand form (e.g. "03F") to full form (e.g. "0033FF")
    var shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
    hex = hex.replace(shorthandRegex, function(m, r, g, b) {
        return r + r + g + g + b + b;
    });

    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);

    return "rgba("+parseInt(result[1], 16) + "," +parseInt(result[2], 16) + "," + parseInt(result[3], 16)+",0.6)";
}

page.SetValueFields = function(){
    setTimeout(function(){
        $("#xAxis").data("kendoDropDownList").select(1);
        $("#yAxis").data("kendoDropDownList").select(2);
        $("#y2Axis").data("kendoDropDownList").select(0);
        var date = new Date();
        date.setDate(date.getDate() - 1);
        $('#dateEnd').data("kendoDatePicker").value(date);

        var date2 = new Date();
        date2.setDate(date2.getDate() - 2);
        $('#dateStart').data("kendoDatePicker").value(date2);
    },500);
}

page.checkProject = function (elmId) {
    var arr = $('#'+elmId).data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbines");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#'+elmId).data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $('#'+elmId).data("kendoMultiSelect").value(["All Project"]);
    } else if (arr.length == 0) {
        $('#'+elmId).data("kendoMultiSelect").value(["All Project"]);
    }
}
page.GetFieldList = function(){
    var param  = {project : $('#ProjectList').data('kendoMultiSelect').value()}
    toolkit.ajaxPostDeffered(viewModel.appName + "xyanalysis/getxyfieldlist", param, function(res) {
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

                $("#xAxis").data("kendoDropDownList").select(1);
                $("#yAxis").data("kendoDropDownList").select(2);
                $("#y2Axis").data("kendoDropDownList").select(0);
            },300)
        }   
    });
}

page.InitDefaultValue = function () {
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");
    $("#StateList").data("kendoDropDownList").trigger("change");
    setTimeout(function(){
        $.when(page.SetValueFields()).done(function(){
            page.checkKey();
        });
    },500)

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
    },500)

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
    },500);
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
        page.GenerateTurbineList();
        page.checkKey();
        app.loading(false);
    });
}

page.GenerateChart = function(dataSource) {

    var seriesData = [];
    var color = 0;
    $.each(dataSource.data, function(i, val){
        $.each(dataSource.axisinfo, function(e, axis){
            if(e > 0){

                var valueColor = (e == 1 ? colors[color] : colorsDeg[color-1]) ;
                var series =   {
                    colorField: "valueColor",
                    data: [],
                    "markers": {
                      "size": 2
                    },
                    name: (e == 1 ? 'Y' : 'Y1') + " ("+val.turbine+")",
                    type: "scatter",
                    xField: dataSource.axisinfo[0].Id,
                    yField: dataSource.axisinfo[e].Id,
                    turbineid : val.turbine,
                    color : valueColor,
                }

                if(e == 2){
                    series.yAxis = dataSource.axisinfo[2].Id;
                }
                
                $.each(val.detail, function(index, value){
                    var seriesData = {};
                    $.each(dataSource.axisinfo, function(j, axisValue){
                        seriesData[axisValue.Id] = value[j];
                        seriesData["valueColor"] = valueColor;
                    });
                    series.data.push(seriesData);
                });
                seriesData.push(series);
            }
            color++;
        })
    });

    page.Data(seriesData);

    var xAxis = {};
    var yAxes = [];

    $.each(dataSource.axisinfo , function(i, value){
        var xField = dataSource.axisinfo[0].Id;
        var yAxis = dataSource.axisinfo[2].Id;
        if(i == 0){
            xAxis.axisCrossingValues =  [-10000, 10000];
            xAxis.title = {
                text: value.Text,
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
            };
            xAxis.majorGridLines = {
                            visible: true,
                            color: "#eee",
                            width: 0.8,
            };
            xAxis.crosshair = {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                };

        }
        else{
            var yAx = {
                name: value.Id, 
                title : {
                    text: value.Text,
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    visible: true,
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                axisCrossingValues : [-10000, 10000],
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N1",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
            }

            yAxes.push(yAx);
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
        seriesDefaults: {
            type: "scatter",
            style: "smooth",
            markers:{size: 2}
        },
        xAxis: xAxis,
        yAxis: yAxes,
        tooltip: {
            visible: true, 
            fomat : "N2"
        },
        categoryAxis: {
            line: {
                visible: false
            },
            labels: {
                padding: {top: 135}
            }
        },
    });

    $("#chartxyAnalysis").data("kendoChart").options.series = page.Data();
    $("#chartxyAnalysis").data("kendoChart").redraw();
}

page.GenerateTurbineList = function() {
    var dtTurbines = page.Data();

    if (dtTurbines.length > 1) {
        $("#showHideChk").html('<label>' +
            '<input type="checkbox" id="showHideAll" checked onclick="page.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChk").html("");
    }

    $("#right-turbine-list").html("");
    $.each(dtTurbines, function(idx, val) {
        var nameTurbine = val.name;

        $("#right-turbine-list").append('<div class="btn-group">' +
        '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="page.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + idx + '"></i></button>' +
        '<input class="chk-option" type="checkbox" name="' + val.turbineid + '" checked id="chk-' + idx + '" hidden>' +
        '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="page.toDetail(\'' + val.turbineid + '\',\'' + val.turbineid + '\')" type="button">' + nameTurbine +'</button>' +
        '</div>');
    });
}

page.showHideAllLegend = function(e) {
    var dtTurbines = page.Data();
    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            $("#chartxyAnalysis").data("kendoChart").options.series[val.index].visible = true;
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            $("#chartxyAnalysis").data("kendoChart").options.series[val.index].visible = false;
        });

        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    $("#chartxyAnalysis").data("kendoChart").redraw();
}


page.showHideLegend = function(idx) {
    $('#chk-' + idx).trigger('click');
    var chart = $("#chartxyAnalysis").data("kendoChart");


    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAll').prop('checked', true);
    } else {
        $('#showHideAll').prop('checked', false);
    }

    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }
    if (idx == $('input[id*=chk-][type=checkbox]').length) {
        idx == 0
    }

    chart._legendItemClick(idx);
}

$(function() {
    

    $("#ProjectList").kendoMultiSelect({
        dataSource: page.ProjectListState(), 
        dataValueField: 'ProjectId', 
        dataTextField: 'ProjectId', 
        change: function() {
            page.SetTurbineList();
            page.GetFieldList();
        }, 
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