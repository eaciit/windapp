'use strict';


viewModel.TLPlots = new Object();
var tlp = viewModel.TLPlots;


vm.currentMenu('TrendLinePlots');
vm.currentTitle('TrendLinePlots');
vm.breadcrumb([ {
    title: 'Analysis Tool Box',
    href: '#'
}, {
    title: 'Trend Line Plots',
    href: "#"//viewModel.appName + 'page/analytictrendlineplots'
}]);

tlp.turbineList = ko.observableArray([]);

tlp.turbine = ko.observableArray([]);
tlp.compTemp = ko.observableArray([
    { "value": 2, "text": "Ambient Temp", "colname": "temp_yawbrake_1" },
    { "value": 5, "text": "Temp_GearBox_IMS_NDE", "colname": "temp_gearbox_ims_nde" }, 
    { "value": 5, "text": "Temp_GearBox_HSS_NDE", "colname": "temp_gearbox_hss_nde"  },
    { "value": 5, "text": "Temp_G1L1", "colname": "temp_g1l1"  },
    { "value": 5, "text": "Temp_G1L2", "colname": "temp_g1l2"  },
    { "value": 5, "text": "Temp_G1L3", "colname": "temp_g1l3"  },
    { "value": 5, "text": "Temp_GearBox_HSS_DE", "colname": "temp_gearbox_hss_de"  },
    { "value": 5, "text": "Temp_GearOilSump", "colname": "temp_gearoilsump"  },
    { "value": 5, "text": "Temp_GeneratorBearing_DE", "colname": "temp_generatorbearing_de"  },
    { "value": 5, "text": "Temp_GeneratorBearing_NDE", "colname": "temp_generatorbearing_nde"  },
    { "value": 5, "text": "Temp_MainBearing", "colname": "temp_mainbearing"  },
    { "value": 5, "text": "Temp_GearBox_IMS_DE", "colname": "temp_gearbox_ims_de"  },
    { "value": 5, "text": "Converter-1,2 temp", "colname": ""  },
    { "value": 5, "text": "Nacelle Temp", "colname": "temp_nacelle"  },
]);

tlp.compTempVal = ko.observable("2");

// function Categories(start,end){

// }


tlp.initChart = function() {
    app.loading(true);

    app.ajaxPost(viewModel.appName + "/trendlineplots/getscadaoemavaildate", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        //Scada Data
        if (res.data.ScadaOemAvailDate.length == 0) {
            res.data.ScadaOemAvailDate = [];
        } else {
            if (res.data.ScadaOemAvailDate.length > 0) {
                var minDatetemp = new Date(res.data.ScadaOemAvailDate[0]);
                var maxDatetemp = new Date(res.data.ScadaOemAvailDate[1]);
                $('#availabledatestarttlp').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                $('#availabledateendtlp').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
            }
        }
    });

    var compTemp =  $('#compTemp').data('kendoDropDownList').text()
    var colnameTemp = _.find(tlp.compTemp(), function(num){ return num.text == compTemp; }).colname;
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        colname: colnameTemp,
        breakdown:""
    };

    var link = "trendlineplots/getlist"


    toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        var datatlp = res.data.Data;
        var categories = res.data.Categories;
        var catTitle = res.data.CatTitle;

        localStorage.setItem("datatlp", JSON.stringify(datatlp));

        $('#charttlp').html("");
        
        $("#charttlp").kendoChart({
            pdf: {
                fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            renderAs: "canvas",
            title: {
                text: "Trend Line Plots | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: false,
            },
            chartArea: {
                height: 325,
            },
            seriesDefaults: {
                type: "line",
                style: "smooth",
                dashType: "longDash",
                markers: {
                    visible: false,
                    size: 4,
                },
            },
            // seriesColors: colorField,
            series: datatlp,
            valueAxis: {
                name: compTemp,
                title: {
                    text: compTemp,
                    visible: true,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                labels: {
                    step: 2
                },
                line: {
                    visible: false
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            },
            categoryAxis: {
                categories: categories,
                majorGridLines: {
                    visible: false
                },
                title: {
                    text: catTitle,
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                majorTickType: "none",
            },
            tooltip: {
                visible: true,
                format: "{0:n1}",
                background: "rgb(255,255,255, 0.9)",
                shared: true,
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
            // pannable: true,
            // zoomable: true
        });

        app.loading(false);
        $("#charttlp").data("kendoChart").refresh();

        tlp.InitRightTurbineList();
        
    });
}

tlp.InitRightTurbineList = function(){
    tlp.turbineList([]);
    var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("datatlp")).sort(name), 'name');

    if (dtTurbines.length > 1) {
        $("#showHideChk").html('<label>' +
            '<input type="checkbox" id="showHideAll" checked onclick="tlp.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChk").html("");
    }

    $("#right-turbine-list").html("");
    $.each(dtTurbines, function(idx, val) {
        if(val.name != "Power Curve"){
            $("#right-turbine-list").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="tlp.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + idx + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="chk-' + idx + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" type="button" onclick="tlp.showHideLegend(' + idx  + ')">' + val.name + '</button>' +
            '</div>');
        }
    });
}

tlp.showHideAllLegend = function(e){
    var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("datatlp")).sort(name), 'name');
                console.log(dtTurbines)
    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#charttlp").data("kendoChart").options.series[val.idxseries - 1].visible = true;
            }
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#charttlp").data("kendoChart").options.series[val.idxseries - 1].visible = false;
            }
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    $("#charttlp").data("kendoChart").redraw();
}
tlp.showHideLegend = function(idx){
    $('#chk-' + idx).trigger('click');
    var chart = $("#charttlp").data("kendoChart");

    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAll').prop('checked', true);
    } else {
        $('#showHideAll').prop('checked', false);
    }

    // var idxHide = idx+1
    // console.log(idxHide)
    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }
    // if (idx == $('input[id*=chk-][type=checkbox]').length) {
    //     idx == 0
    // }

    chart._legendItemClick(idx);
}



$(document).ready(function() {
    setTimeout(function() {
        fa.LoadData();
        tlp.initChart();
    }, 300);

    $('#btnRefresh').on('click', function() {
        setTimeout(function() {
            tlp.initChart();
        }, 300);
    });

     $('#compTemp').on("change", function() {
        setTimeout(function() {
            tlp.initChart();
        }, 300);
    });


});


