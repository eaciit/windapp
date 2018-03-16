'use strict';


viewModel.TLPlots = new Object();
var tlp = viewModel.TLPlots;


vm.currentMenu('Trend Line Plots');
vm.currentTitle('TrendLinePlots');
vm.breadcrumb([ {
    title: 'Analysis Tool Box',
    href: '#'
}, {
    title: 'Trend Line Plots',
    href: "#"//viewModel.appName + 'page/analytictrendlineplots'
}]);

tlp.turbineList = ko.observableArray([]);
tlp.temperatureList = ko.observableArray([]);

tlp.turbine = ko.observableArray([]);
tlp.compTemp = ko.observableArray([
    { "value": 1, "text": "Ambient Temp", "colname": "tempoutdoor" },
    // { "value": 2, "text": "Temp_GearBox_IMS_NDE", "colname": "temp_gearbox_ims_nde" }, 
    // { "value": 3, "text": "Temp_GearBox_HSS_NDE", "colname": "temp_gearbox_hss_nde"  },
    // { "value": 4, "text": "Temp_G1L1", "colname": "temp_g1l1"  },
    // { "value": 5, "text": "Temp_G1L2", "colname": "temp_g1l2"  },
    // { "value": 6, "text": "Temp_G1L3", "colname": "temp_g1l3"  },
    // { "value": 7, "text": "Temp_GearBox_HSS_DE", "colname": "temp_gearbox_hss_de"  },
    // { "value": 8, "text": "Temp_GearOilSump", "colname": "temp_gearoilsump"  },
    { "value": 9, "text": "Temp_GeneratorBearing_DE", "colname": "tempgeneratorbearingde"  },
    { "value": 10, "text": "Temp_GeneratorBearing_NDE", "colname": "tempgeneratorbearingnde"  },
    // { "value": 11, "text": "Temp_MainBearing", "colname": "temp_mainbearing"  },
    // { "value": 12, "text": "Temp_GearBox_IMS_DE", "colname": "temp_gearbox_ims_de"  },
    // { "value": 13, "text": "Converter-1,2 temp", "colname": ""  },
    { "value": 14, "text": "Nacelle Temp", "colname": "tempnacelle"  },
]);

tlp.deviation= ko.observable(2);
tlp.deviationList = ko.observableArray([1,2,3,4,5]);
tlp.isDeviation = ko.observable(false);
tlp.compTempVal = ko.observable(1);
tlp.project = ko.observable();
tlp.dateStart = ko.observable();
tlp.dateEnd = ko.observable();
var origWidth;
var origTransitions;
var seriesIndex;
var avgMaxValue = "";

function seriesHover(e){
    origWidth = e.series.width;
    origTransitions = e.sender.options.transitions;
    seriesIndex = e.series.index;
    e.sender.unbind("seriesHover", seriesHover);

    e.series.width = 5;
    e.sender.options.transitions = false;
    e.sender.redraw();

    setTimeout(function(){
        var chart = $("#charttlp").data("kendoChart");

        chart.options.series[seriesIndex].width = origWidth;
        chart.redraw();
        chart.options.transitions = origTransitions;
        chart.bind("seriesHover", seriesHover);
    },1000);
}

tlp.getPDF = function(selector){
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            margin: {
                left   : "5mm",
                top    : "5mm",
                right  : "5mm",
                bottom : "5mm"
            },
        });
      kendo.drawing.pdf.saveAs(group, project+"TrendlinePlots"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

tlp.getAvailDate = function(){
    return app.ajaxPost(viewModel.appName + "/analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var availDateAll = res.data;
        var namaproject = $('#projectList').data("kendoDropDownList").value();

        if(namaproject == "") {
            namaproject = "Tejuva";
        }

        // console.log(availDateAll[namaproject]);
        var startDate = kendo.toString(moment.utc(availDateAll[namaproject]["ScadaDataHFD"][0]).format('DD-MMM-YYYY'));
        var endDate = kendo.toString(moment.utc(availDateAll[namaproject]["ScadaDataHFD"][1]).format('DD-MMM-YYYY'));

        $('#availabledatestarttlp').html(startDate);
        $('#availabledateendtlp').html(endDate);

        var maxDateData = new Date(availDateAll[namaproject]["ScadaDataHFD"][1]);

        if(moment(maxDateData).get('year') !== 1){
            var startDatepicker = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

            $('#dateStart').data('kendoDatePicker').value(startDatepicker);
            $('#dateEnd').data('kendoDatePicker').value(endDate);
        }
    });
}

tlp.initChart = function() {
    app.loading(true);
    var project = $('#projectList').data("kendoDropDownList").value();
    tlp.compTemp(tlp.temperatureList()[project]);    

    var compTemp =  $('#compTemp').data('kendoDropDownList').text()
    var ddldeviation = $('#deviationValue').val()
    var colnameTemp = _.find(tlp.compTemp(), function(num){ return num.text == compTemp; }).colname;
    if (avgMaxValue !== "") {
        colnameTemp += "_" + avgMaxValue
    }
    // var turb = $("#turbineList").data("kendoMultiSelect").value()[0] == "All Turbine" ? [] : $("#turbineList").data("kendoMultiSelect").value()
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 

    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: fa.turbine(), // $("#turbineList").data("kendoMultiSelect").value(),
        project: fa.project,
        colname: colnameTemp,
        deviationstatus:tlp.isDeviation(), // Param from checkbox
        deviation: parseFloat(ddldeviation)// Param from Dropdown
    };


    fa.dateStart = dateStart;
    fa.dateEnd = dateEnd ;

    var link = "trendlineplots/getlist"


    toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        var tempData = [];
        var firstData = [];

        if(res.data.Data.length > 1) {
            firstData.push(res.data.Data[0]);
            tempData = res.data.Data.slice(1);
            if (res.data.Data[1].name == "Met Tower") {
                firstData.push(res.data.Data[1]);
                tempData = res.data.Data.slice(2);
            }
            tempData = _.sortBy(tempData, 'name');
        }
        res.data.Data = firstData.concat(tempData);
        res.data.Data.forEach(function(val, idx){
            res.data.Data[idx].idxseries = idx;
        });


        var datatlp = res.data.Data;
        var categories = res.data.Categories;
        var catTitle = res.data.CatTitle;
        var minValue = res.data.Min;
        var maxValue = res.data.Max;
        var nullCount = 0
        datatlp.forEach( function(data, idxTlp) {
            if(data.data != undefined && data.data != null) {
                nullCount = 0
                data.data.forEach( function(element, idxData) {
                    if(element == 999999) {
                        nullCount++
                        datatlp[idxTlp].data[idxData] = null;
                    }
                });
                datatlp[idxTlp]["missingValues"] = "gap";
            }
        });

        localStorage.setItem("datatlp", JSON.stringify(datatlp));


        $('#charttlp').html("");
        
        $("#charttlp").kendoChart({
            pdf: {
                fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            // renderAs: "canvas",
            title: {
                text: "Trend Line Plots | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: false,
                labels : {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    // template: kendo.template($("#legendItemTemplate").html()),
                }
            },
            chartArea: {
                height: 400,
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
                    text: "Temperature",//compTemp,
                    visible: true,
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                labels: {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                line: {
                    visible: false
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#bdbdbd",
                    width: 0.8,
                },
                // majorUnit: 0.5,
                min: minValue,
                max: maxValue,
            },
            categoryAxis: {
                categories: categories,
                majorGridLines: {
                    visible: false
                },
                title: {
                    text: catTitle,//"Time",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                majorTickType: "none",
            }, 
            seriesHover : seriesHover,
            tooltip: {
                visible: true,
                format: "{0:n1}",
                background: "rgb(255,255,255, 0.9)",
                template: "#= series.name # : #= kendo.toString(value,'n2')#",
                shared: false,
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
            // pannable: true,
            // zoomable: true
        });

        app.loading(false);
        $("#charttlp").data("kendoChart").refresh();
        tlp.InitRightTurbineList(res.data.TurbineName);
        tlp.getActiveLegend();
        
    });
}

tlp.getActiveLegend = function(){
    var chart = $("#charttlp").data("kendoChart");
    var viewModel = kendo.observable({
      series: chart.options.series,
      markerColor: function(e) {
        return e.get("visible") ? e.color : "grey";
      }
    });

    kendo.bind($("#legend"), viewModel);
}
tlp.InitRightTurbineList = function(turbinename){
    tlp.turbineList([]);
    
    var dtTurbines = JSON.parse(localStorage.getItem("datatlp"));

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
        if(val.idxseries > 0){
            if(val.data == undefined || val.data == ""){
                $("#right-turbine-list").append('<div class="btn-group">' +
                '<button class="btn btn-default btn-sm turbine-chk" type="button" disabled="disabled" onclick="tlp.showHideLegend(' + val.idxseries + ')" style="border-color:#a0a0a0;background-color:#a0a0a0"></button>' +
                '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="disabled-' + val.idxseries + '" hidden disabled="disabled">' +
                '<button class="btn btn-default btn-sm turbine-btn wbtn" type="button" disabled="disabled" onclick="tlp.showHideLegend(' + val.idxseries  + ')">' + val.name + '</button>' +
                '</div>');
            }else{
                $("#right-turbine-list").append('<div class="btn-group">' +
                '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="tlp.showHideLegend(' + val.idxseries + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + val.idxseries + '"></i></button>' +
                '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="chk-' + val.idxseries + '" hidden>' +
                '<button class="btn btn-default btn-sm turbine-btn wbtn" type="button" onclick="tlp.showHideLegend(' + val.idxseries  + ')">' + val.name + '</button>' +
                '</div>'); 
            }
        }
    });
}

tlp.showHideAllLegend = function(e){
    var dtTurbines = JSON.parse(localStorage.getItem("datatlp"));

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#charttlp").data("kendoChart").options.series[val.idxseries].visible = true;
            }
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#charttlp").data("kendoChart").options.series[val.idxseries].visible = false;
            }
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    $("#charttlp").data("kendoChart").redraw();
    tlp.getActiveLegend();
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
    tlp.getActiveLegend();
}


$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        setTimeout(function() {
            if(fa.LoadData()) {
                var project = $('#projectList').data("kendoDropDownList").value();
                var dateStart = $('#dateStart').data('kendoDatePicker').value();
                var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
                tlp.project(project);
                tlp.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
                tlp.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
                tlp.initChart();
            }
        }, 300);
    });

    $("input[name=isAvg]").on("change", function() {
        avgMaxValue = this.value;
        fa.checkTurbine();
        setTimeout(function() {
            if(fa.LoadData()) {
                // tlp.initChart();
                $.when(tlp.initChart()).done(function(){
                    $("#charttlp").data("kendoChart").bind("seriesHover", seriesHover);
                });
            }
        }, 300);
    })

     $('#compTemp').on("change", function() {
        fa.checkTurbine();
        setTimeout(function() {
            if(fa.LoadData()) {
                // tlp.initChart();
                $.when(tlp.initChart()).done(function(){
                    $("#charttlp").data("kendoChart").bind("seriesHover", seriesHover);
                });
            }
        }, 300);
    });

    $('#projectList').kendoDropDownList({
        change: function () {
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateEngine(project, true);
            // fa.populateTurbine(project);
            fa.project = project;
            tlp.getAvailDate();
            tlp.compTemp(tlp.temperatureList()[project]);
        }
    });

    $.when(tlp.getAvailDate()).done(function(){
        fa.checkTurbine();
        if(fa.LoadData()) {
            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
            tlp.project(project);
            tlp.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            tlp.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
            tlp.initChart();
        }
    })
});