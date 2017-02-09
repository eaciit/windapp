'use strict';

viewModel.WindRoseComparison = new Object();
var wrb = viewModel.WindRoseComparison;

// WIND ROSE COMPARISON
wrb.sectorDerajatComparison = ko.observable(0);
wrb.dataWindroseComparison = ko.observableArray([]);
var listOfChartComparison = [];
var listOfButtonComparison = {};
/*wrb.showHideLegendComparison = function (index) {
    var idName = "btn" + index;
    listOfButtonComparison[idName] = !listOfButtonComparison[idName];
    if (listOfButtonComparison[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChartComparison, function (idx, idChart) {
       if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
          $(idChart).data("kendoChart").options.series[index].visible = listOfButtonComparison[idName];
          $(idChart).data("kendoChart").refresh();
        }
    });
}*/

wrb.InitTurbineListCompare = function () {
    if (wrb.dataWindroseComparison().Data.length > 1) {
        $("#checkAllCompare").html('<label id="checkAllLabel">' +
            '<input type="checkbox" id="showHideAllCompare" checked onclick="wrb.showHideAllCompare(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHideCompare"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#checkAllCompare").html("");
    }

    $("#turbine-list-compare").html("");
    $.each(wrb.dataWindroseComparison().Data, function (idx, val) {
        $("#turbine-list-compare").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="wrb.showHideCompare(' + val.idxseries + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + val.idxseries + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="chk-' + val.idxseries + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="wrb.showHideCompare(' + val.idxseries + ')" type="button">' + val.name + '</button>' +
            '</div>');
    });
}

wrb.showHideAllCompare = function (e) {

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(wrb.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = true;
            }
        });
        $('#labelShowHideCompare b').text('Select All');
    } else {
        $.each(wrb.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHideCompare b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);

    $("#WRChartComparison").data("kendoChart").redraw();
}

wrb.showHideCompare = function (idx) {
    var stat = false;

    $('#chk-' + idx).trigger('click');
    var chart = $("#WRChartComparison").data("kendoChart");
    var leTur = $('input[id*=chk-][type=checkbox]').length

    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAllCompare').prop('checked', true);
    } else {
        $('#showHideAllCompare').prop('checked', false);
    }

    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }

    if ($('#chk-' + idx).is(':checked')) {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = false
    }
    $("#WRChartComparison").data("kendoChart").redraw();
}

wrb.initChartWRC = function () {
    listOfChartComparison = [];
    var dataSeries = wrb.dataWindroseComparison().Data;
    var categories = wrb.dataWindroseComparison().Categories;
    var nilaiMax = wrb.dataWindroseComparison().MaxValue;

    // var paddingTitle = Math.floor($('.windrose-part').width() / 16.16) * -1;
    // var offsetLegend = Math.floor($('.windrose-part').width() / 3.46) * -1;
    var majorUnit = 10;
    if(nilaiMax < 40) {
        majorUnit = 5;
    }

    $("#WRChartComparison").kendoChart({
        theme: "flat",
        title: {
            // padding: {
            //   left: paddingTitle
            // },
            text: "Wind Rose Comparison",
            font: '16px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            visible: false
        },
        legend: {
            position: "right",
            // offsetX: offsetLegend,
            labels: {
                font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
            },
            visible: false
        },
        dataSource: {
            sort: {
                field: "DirectionNo",
                dir: "asc"
            }
        },
        series: dataSeries,
        categoryAxis: {
            categories: categories,
            visible: true,
            majorGridLines: {
                visible: false
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
            }
        },
        valueAxis: {
            labels: {
                template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorUnit: majorUnit,
            // max: nilaiMax,
            // min: 0
        },
        tooltip: {
            visible: true,
            // template: "#= category # #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
            template: "#= series.name # : #= category #"+String.fromCharCode(176)+" (#= kendo.toString(value, 'n2') #%)",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        }
    });
    wrb.InitTurbineListCompare();
    // $('#WRChartComparison').data('kendoChart').options.chartArea.width = $('#WRChartComparison').height() + ($('#WRChartComparison').height()/4);
    // $('#WRChartComparison').data('kendoChart').refresh();
}

wrb.WindRoseComparison = function(){
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstWindRoseComparison() === true){
        setTimeout(function () {
            // var breakDownVal = $("#nosectionComparison").data("kendoDropDownList").value();
            var breakDownVal = "36";
            var secDer = 360 / breakDownVal;
            // wrb.sectorDerajatComparison(secDer);
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine,
                project: fa.project,
                breakDown: breakDownVal,
            };
            toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getwindrosedata", param, function (res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                if (res.data != null) {
                    wrb.dataWindroseComparison(res.data);
                    wrb.initChartWRC();
                }

                app.loading(false);
                pm.isFirstWindRoseComparison(false);

            })
        }, 300);
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
    }else{
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
        setTimeout(function(){
            $.each(listOfChartComparison, function(idx, elem){
                $(elem).data("kendoChart").refresh();
            });
            app.loading(false);
        }, 300);
    }
}