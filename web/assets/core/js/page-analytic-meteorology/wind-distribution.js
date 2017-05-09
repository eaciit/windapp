'use strict';

viewModel.WindDistribution = new Object();
var wd = viewModel.WindDistribution;

wd.turbineList = ko.observableArray([]);
wd.turbine = ko.observableArray([]);

wd.populateTurbine = function(){
    wd.turbine([]);
    if(fa.turbine == ""){
        $.each(fa.turbineList(), function(i, val){
            if (i > 0){
                wd.turbine.push(val.text);
            }
        });
    }else{
        wd.turbine(fa.turbine);
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
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="wd.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + (idx) + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.turbine + '" checked id="chk-' + (idx) + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="wd.showHideLegend(' + (idx) + ')" type="button">' + val.turbine + '</button>' +
            '</div>');
    });
}

wd.ChartWindDistributon =  function () {
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwinddistribution/getlist", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        if (wd.turbine().length == 0) {
            var turbine = []
            for (var i=0;i<res.data.Data.length;i++) {
                if ($.inArray( res.data.Data[i].Turbine, turbine ) == -1){
                    turbine.push(res.data.Data[i].Turbine);
                }
            }

            $.each(turbine, function (i, val) {
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
                    visible: true,
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

        // app.loading(false);
        $("#windDistribution").data("kendoChart").refresh();
    });
}

wd.showHideAllLegend = function (e) {

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(wd.turbine(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = true;
            }
        });
        /*$('#labelShowHide b').text('Untick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(wd.turbine(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check').css("visibility", 'hidden');
        /*$('#labelShowHide b').text('Tick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);

    $("#windDistribution").data("kendoChart").redraw();
}

wd.showHideLegend = function (idx) {
    var stat = false;

    $('#chk-' + idx).trigger('click');
    var chart = $("#windDistribution").data("kendoChart");
    var leTur = $('input[id*=chk-][type=checkbox]').length

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

    if ($('#chk-' + idx).is(':checked')) {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = false
    }
    $("#windDistribution").data("kendoChart").redraw();
}

wd.WindDis = function(){
    
    fa.LoadData();
    pm.showFilter();
    if(pm.isFirstWindDis() === true){
        app.loading(true);
        wd.populateTurbine();
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