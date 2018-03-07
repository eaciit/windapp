'use strict';

viewModel.NacelleDistribution = new Object();
var nd = viewModel.NacelleDistribution;

nd.turbineList = ko.observableArray([]);
nd.turbine = ko.observableArray([]);
nd.project = ko.observable();
nd.dateStart = ko.observable();
nd.dateEnd = ko.observable();


nd.getPDF = function(selector){
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
      kendo.drawing.pdf.saveAs(group, project+"NDDistribution"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}


nd.InitRightTurbineList= function () {
    if (nd.turbineList().length > 1) {
        $("#showHideChkNacelle").html('<label>' +
            '<input type="checkbox" id="showHideAllNacelle" checked onclick="nd.showHideAllNacelleLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelNacelleShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChkNacelle").html("");
    }

    $("#right-turbine-list-nacelle").html("");

    $.each(nd.turbineList(), function (idx, val) {
        $("#right-turbine-list-nacelle").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="nd.showHideNacelleLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check fa-check-nacelledist" id="icon-nacelle-distribution' + (idx) + '"></i></button>' +
            '<input class="chk-option-dist" type="checkbox" name="' + val.turbine + '" checked id="chk-nacelle-distribution' + (idx) + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="nd.showHideNacelleLegend(' + (idx) + ')" type="button">' + val.turbine + '</button>' +
            '</div>');
    });
}

nd.ChartNacelleDistributon =  function () {
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
        breakdown: "nacelledeviation",
        color: color,
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwinddistribution/getlist", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        nd.turbineList([]);
        $.each(res.data.TurbineList, function (i, val) {
            var data = {
                color: color[i],
                turbine: val
            }

            nd.turbineList.push(data);
        });

        $('#nacelleDistribution').html("");
        var categories = res.data.Categories;
        var dataSeries = res.data.Data;
        var crossingValue = categories.length/2; /* supaya yAxis bisa pas di tengah */
        if (dataSeries.length == 0) {
            crossingValue -= 0.5; /* entah kenapa jika tidak ada data, maka akan bergeser 0.5 point makanya dikurangi 0.5 */
        }

        $("#nacelleDistribution").kendoChart({
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
            series: dataSeries,
            valueAxis: {
                labels: {
                    format: "{0:p0}",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                line: {
                    visible: true
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                }
            },
            categoryAxis: {
                categories: categories,
                majorGridLines: {
                    visible: false
                },
                axisCrossingValue: crossingValue,
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    // rotation: 25
                },
                majorTickType: "none"
            },
            tooltip: {
                visible: true,
                // template: "Contribution of #= series.name # : #= kendo.toString(value, 'n4')# % at #= category #",
                template: "#= series.name # : #= category #"+String.fromCharCode(176)+" (#= kendo.toString(value, 'p2')#)",
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
                pm.isFirstNacelleDis(false);
                nd.getLegendList();

            }
        });

        nd.InitRightTurbineList();
        $("#nacelleDistribution").data("kendoChart").refresh();
    });
}

nd.getLegendList = function(){
    var chart = $("#nacelleDistribution").data("kendoChart");
    var viewModel = kendo.observable({
      series: chart.options.series,
      markerColor: function(e) {
        return e.get("visible") ? e.color : "grey";
      }
    });

    kendo.bind($("#legendNd"), viewModel);
}
nd.showHideAllNacelleLegend = function (e) {
    nd.getLegendList();
    if (e.checked == true) {
        $('.fa-check-nacelledist').css("visibility", 'visible');
        $.each(nd.turbineList(), function (i, val) {
            if($("#nacelleDistribution").data("kendoChart").options.series[i] != undefined){
                $("#nacelleDistribution").data("kendoChart").options.series[i].visible = true;
            }
        });
        $('#labelNacelleShowHide b').text('Select All');
    } else {
        $.each(nd.turbineList(), function (i, val) {
            if($("#nacelleDistribution").data("kendoChart").options.series[i] != undefined){
                $("#nacelleDistribution").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check-nacelledist').css("visibility", 'hidden');
        $('#labelNacelleShowHide b').text('Select All');
    }
    $('.chk-option-dist').not(e).prop('checked', e.checked);

    $("#nacelleDistribution").data("kendoChart").redraw();
}

nd.showHideNacelleLegend = function (idx) {
    var stat = false;

    $('#chk-nacelle-distribution' + idx).trigger('click');
    var chart = $("#nacelleDistribution").data("kendoChart");
    var leTur = $('input[id*=chk-nacelle-distribution][type=checkbox]').length

    if ($('input[id*=chk-nacelle-distribution][type=checkbox]:checked').length == $('input[id*=chk-nacelle-distribution][type=checkbox]').length) {
        $('#showHideAllNacelle').prop('checked', true);
    } else {
        $('#showHideAllNacelle').prop('checked', false);
    }

    if ($('#chk-nacelle-distribution' + idx).is(':checked')) {
        $('#icon-nacelle-distribution' + idx).css("visibility", "visible");
    } else {
        $('#icon-nacelle-distribution' + idx).css("visibility", "hidden");
    }

    if ($('#chk-nacelle-distribution' + idx).is(':checked')) {
        $("#nacelleDistribution").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#nacelleDistribution").data("kendoChart").options.series[idx].visible = false
    }
    $("#nacelleDistribution").data("kendoChart").redraw();
    nd.getLegendList();
}

nd.NacelleDis = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        if(pm.isFirstNacelleDis() === true){
            app.loading(true);
            nd.ChartNacelleDistributon();
            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
            nd.project(project);
            nd.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            nd.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
        }else{
            app.loading(false);
            setTimeout(function () {
                $('#nacelleDistribution').data('kendoChart').refresh();
                var project = $('#projectList').data("kendoDropDownList").value();
                var dateStart = $('#dateStart').data('kendoDatePicker').value();
                var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
                nd.project(project);
                nd.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
                nd.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
                app.loading(false);
            }, 100);
        }
    }
}