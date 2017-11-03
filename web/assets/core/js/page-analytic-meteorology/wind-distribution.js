'use strict';

viewModel.WindDistribution = new Object();
var wd = viewModel.WindDistribution;

wd.turbineList = ko.observableArray([]);
wd.turbine = ko.observableArray([]);

wd.populateTurbine = function(){
    wd.turbine([]);
    if(fa.turbine().length == 0){
        $.each(fa.turbineList(), function(i, val){
            if (i > 0){
                wd.turbine.push(val.text);
            }
        });
    }else{
        wd.turbine(fa.turbine());
    }

}

wd.InitRightTurbineList= function () {
    if (wd.turbine().length > 0) {
        wd.turbineList([]);
        $.each(wd.turbine(), function (i, val) {
            var data = {
                color: color[i],
                turbine: val
            }

            wd.turbineList.push(data);
        });
    }

    if (wd.turbineList().length > 1) {
        $("#showHideChk").html('<label>' +
            '<input type="checkbox" id="showHideAll" checked onclick="wd.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChk").html("");
    }

    $("#right-turbine-list").html("");

    $.each(wd.turbineList(), function (idx, val) {
        $("#right-turbine-list").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="wd.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check fa-check-winddist" id="icon-wind-distribution' + (idx) + '"></i></button>' +
            '<input class="chk-option-winddist" type="checkbox" name="' + val.turbine + '" checked id="chk-wind-distribution' + (idx) + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="wd.showHideLegend(' + (idx) + ')" type="button">' + val.turbine + '</button>' +
            '</div>');
    });
}

wd.ChartWindDistributon =  function () {
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
        breakdown: "avgwindspeed",
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwinddistribution/getlist", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        if (wd.turbine().length == 0) {
            wd.turbineList([]);

            $.each(res.data.TurbineList, function (i, val) {
                var data = {
                    color: color[i],
                    turbine: val
                }

                wd.turbineList.push(data);
            });


        }

        $('#windDistribution').html("");
        var data = res.data.Data;

        $("#windDistribution").kendoChart({
            dataSource: {
                data: data,
                group: { field: "Turbine" },
                sort: { field: "Category", dir: 'asc' }
            },
            theme: "flat",
            title: {
                text: ""
            },
            legend: {
                position: "right",
                visible: false
            },
            chartArea: {
                height: 360
            },
            series: [{
                type: "line",
                style: "smooth",
                field: "Contribute",
                // opacity : 0.7,
                markers: {
                    visible: false,
                    size: 3,
                }
            }],
            seriesColors: color,
            valueAxis: {
                labels: {
                    format: "{0:p0}",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                line: {
                    visible: true
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                }
            },
            categoryAxis: {
                field: "Category",
                majorGridLines: {
                    visible: false
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    // rotation: 25
                },
                majorTickType: "none"
            },
            tooltip: {
                visible: true,
                // template: "Contribution of #= series.name # : #= kendo.toString(value, 'n4')# % at #= category #",
                template: "#= kendo.toString(value, 'p2')#",
                // shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },

            },
            dataBound: function(){
                app.loading(false);
                pm.isFirstWindDis(false);
            }
        });

        wd.InitRightTurbineList();

        /* hanya untuk mengembalikan Met Tower ke urutan pertama
        entah kenapa setelah di grouping oleh kendo otomatis melakukan sorting based on series name nya */
        var seriesCurrent = $("#windDistribution").data("kendoChart").options.series;
        var seriesMet = [];
        var seriesScada = [];
        var colorList = [];

        // ambil warnanya terlebih dahulu
        seriesCurrent.forEach(function(val, idx){
            colorList.push(val.color);
        });
        seriesCurrent.forEach(function(val, idx){
            if(val.name == "Met Tower") {
                seriesMet.push(val);
            } else {
                seriesScada.push(val);
            }
        });
        var seriesNew = seriesMet.concat(seriesScada);
        seriesNew.forEach(function(val, idx){
            seriesNew[idx].color = colorList[idx];
        });

        $("#windDistribution").data("kendoChart").options.series = seriesNew;


        // app.loading(false);
        $("#windDistribution").data("kendoChart").refresh();
    });
}

wd.showHideAllLegend = function (e) {

    if (e.checked == true) {
        $('.fa-check-winddist').css("visibility", 'visible');
        $.each(wd.turbineList(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = true;
            }
        });
        /*$('#labelShowHide b').text('Untick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(wd.turbineList(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check-winddist').css("visibility", 'hidden');
        /*$('#labelShowHide b').text('Tick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option-winddist').not(e).prop('checked', e.checked);

    $("#windDistribution").data("kendoChart").redraw();
}

wd.showHideLegend = function (idx) {
    var stat = false;

    $('#chk-wind-distribution' + idx).trigger('click');
    var chart = $("#windDistribution").data("kendoChart");
    var leTur = $('input[id*=chk-wind-distribution][type=checkbox]').length

    if ($('input[id*=chk-wind-distribution][type=checkbox]:checked').length == $('input[id*=chk-wind-distribution][type=checkbox]').length) {
        $('#showHideAll').prop('checked', true);
    } else {
        $('#showHideAll').prop('checked', false);
    }

    if ($('#chk-wind-distribution' + idx).is(':checked')) {
        $('#icon-wind-distribution' + idx).css("visibility", "visible");
    } else {
        $('#icon-wind-distribution' + idx).css("visibility", "hidden");
    }

    if ($('#chk-wind-distribution' + idx).is(':checked')) {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = false
    }
    $("#windDistribution").data("kendoChart").redraw();
}

wd.WindDis = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        if(pm.isFirstWindDis() === true){
            app.loading(true);
            // wd.populateTurbine();
            wd.ChartWindDistributon();
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        }else{
            app.loading(false);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
            setTimeout(function () {
                $('#windDistribution').data('kendoChart').refresh();
                app.loading(false);
            }, 100);
        }
    }
}