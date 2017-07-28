'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.turbineList = ko.observableArray([]);
page.downList = ko.observableArray([]);
page.dtLineChart = ko.observableArray([]);
page.projectList = ko.observableArray([{
    "value": 1,
    "text": "WindFarm-01"
}, {
    "value": 2,
    "text": "WindFarm-02"
}, ]);

page.isMain = ko.observable(true);
page.isDetail = ko.observable(false);
page.detailTitle = ko.observable("");
page.detailStartDate = ko.observable("");
page.detailEndDate = ko.observable("");

page.isClean = ko.observable(true);
page.idName = ko.observable("");
page.isDeviation = ko.observable(true);
page.sScater = ko.observable(false);
page.showDownTime = ko.observable(false);
page.deviationVal = ko.observable("20");
page.viewSession = ko.observable("");
page.turbine = ko.observableArray([]);
page.powerCurveOptions = ko.observable();
page.currProject = ko.observable();
page.project = ko.observable();

page.backToMain = function() {
    page.isMain(true);
    page.isDetail(false);
}
page.toDetail = function(turbineid,turbinename) {
    var isValid = fa.LoadData();
    if (isValid) {
        page.isMain(false);
        page.isDetail(true);
        Data.InitCurveDetail(turbineid,turbinename);
    }
}
page.populateTurbine = function() {
    page.turbine([]);
    if (page.turbine().length == 0) {
        $.each(fa.turbineList(), function(i, val) {
            if (i > 0) {
                page.turbine.push(val.text);
            }
        });
    } else {
        page.turbine(fa.turbine());
    }
}
page.ExportPowerCurvePdf = function() {
    var chart = $("#powerCurve").getKendoChart();
    var container = $('<div />').css({
        position: 'absolute',
        top: 0,
        left: -1500
      }).appendTo('body');


      var options = chart.options;

      var exportOptions ={
            // Custom settings for export
            legend: {
              visible: true
            },
            title:{
                visible: true,
            },
            chartArea: {
                height: 500,
            },
            transitions: false,

            // Cleanup
            render: function(e){
              setTimeout(function(){
                    e.sender.saveAsPDF();
                    container.remove();
              }, 500);
            }
      }

      var options2 = $.extend(true, options, exportOptions);
      container.kendoChart(options2);

      $("#powerCurve").kendoChart($.extend(true, options, {legend: {visible: false},title:{visible: false},chartArea: { height: 375}, render: function(e){return false}}));
}
page.ExportPowerCurveDetailPdf = function() {
        var chart = $("#powerCurveDetail").getKendoChart();
        var container = $('<div />').css({
            position: 'absolute',
            top: 0,
            left: -1500
          }).appendTo('body');


          var options = chart.options;

          var exportOptions ={
                // Custom settings for export
                legend: {
                  visible: true
                },
                title:{
                    visible: true,
                },
                chartArea: {
                    height: 500,
                },
                transitions: false,

                // Cleanup
                render: function(e){
                  setTimeout(function(){
                        e.sender.saveAsPDF();
                        container.remove();
                  }, 500);
                }
          }

          var options2 = $.extend(true, options, exportOptions);
          container.kendoChart(options2);

          $("#powerCurveDetail").kendoChart($.extend(true, options, {legend: {visible: false},title:{visible: false},chartArea: { height: 375}, render: function(e){return false}}));
    }

vm.currentMenu('Power Curve');
vm.currentTitle('Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Power Curve',
    href: viewModel.appName + 'page/analyticpowercurve'
}]);

var dataPowerCurve
var dataTurbine

var Data = {
    LoadData: function() {
        var isValid = fa.LoadData();
        // fa.getProjectInfo();
        if(isValid) {
            page.populateTurbine();
            this.InitLinePowerCurve();
            this.InitRightTurbineList();
        }
    },
    InitLinePowerCurve: function() {
        var isValid = fa.LoadData();
        page.getSelectedFilter();
        if(isValid) {
            page.deviationVal($("#deviationValue").val());

            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function(res) {
                if (!app.isFine(res)) {
                    return;
                }
                var minDatetemp = new Date(res.data.ScadaData[0]);
                var maxDatetemp = new Date(res.data.ScadaData[1]);
                $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
                $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
            })

            var link = "analyticpowercurve/getlistpowercurvescada"

            app.loading(true);
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine(),
                project: fa.project,
                isClean: page.isClean,
                isDeviation: page.isDeviation,
                DeviationVal: page.deviationVal,
                ViewSession: page.viewSession
            };

            toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }

                dataTurbine = res.data.Data;
                localStorage.setItem("dataTurbine", JSON.stringify(res.data.Data));
                page.dtLineChart(res.data.Data);
                
                $('#powerCurve').html("");
                $("#powerCurve").kendoChart({
                    pdf: {
                      fileName: "DetailPowerCurve.pdf",
                    },
                    theme: "flat",
                    title: {
                        text: "Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
                        visible: false,
                        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                    },
                    legend: {
                        position: "bottom",
                        visible: false,
                    },
                    chartArea: {
                        height: 375,
                    },
                    seriesDefaults: {
                        type: "scatterLine",
                        style: "smooth",
                        dashType: "longDash",
                        markers: {
                            visible: false,
                            size: 4,
                        },
                    },
                    seriesColors: colorField,
                    series: dataTurbine,
                    categoryAxis: {
                        labels: {
                            step: 1
                        }
                    },
                    valueAxis: [{
                        labels: {
                            format: "N0",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        }
                    }],
                    xAxis: {
                        majorUnit: 1,
                        title: {
                            text: "Wind Speed (m/s)",
                            font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            color: "#585555",
                            visible: true,
                        },
                        labels: {
                            format: "N0",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        },
                        crosshair: {
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
                        },
                        majorGridLines: {
                            visible: true,
                            color: "#eee",
                            width: 0.8,
                        },
                        max: 25
                    },
                    yAxis: {
                        title: {
                            text: "Generation (KW)",
                            font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            color: "#585555"
                        },
                        labels: {
                            format: "N0",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        },
                        axisCrossingValue: -5,
                        majorGridLines: {
                            visible: true,
                            color: "#eee",
                            width: 0.8,
                        },
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
                    },
                    tooltip: {
                        visible: true,
                        format: "{1}in {0} minutes",
                        template: "#= series.name #",
                        shared: true,
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    },
                    // zoomable: true,
                    zoomable: {
                        selection: {
                            lock: "y",
                        }
                    }
                });
                app.loading(false);
                $("#powerCurve").data("kendoChart").refresh();
                
                if (page.sScater()) {
                    $('#showDownTime').removeAttr("disabled");
                } else {
                    Data.InitRightTurbineList();
                    $('#showDownTime').attr('checked', false);
                    $('#showDownTime').attr("disabled", "disabled");
                    $('#downtime-list').hide();
                    page.showDownTime(false);
                }
                if (page.sScater()) {
                    Data.getPowerCurve();
                }
                page.powerCurveOptions($("#powerCurve").getKendoChart().options);
                page.ShowHideAfterInitChart();
            });
        }
    },
    getPowerCurve: function() {
        page.deviationVal($("#deviationValue").val());
        var turbineList = [];
        var kolor = [];
        // var kolorDeg = [];
        var dataTurbine = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'turbineid');

        var len = $('input[id*=chk-][type=checkbox]:checked').length;

        for (var a = 0; a < len; a++) {
            var chk = $('input[id*=chk-][type=checkbox]:checked')[a].name;
            turbineList.push(chk);
            var even = _.find(dataTurbine, function(nm) {
                return nm.turbineid == chk
            });
            kolor.push(even.color);
            var indOf = 0;
            for (var i = 0; i < colorField.length; i++) {
                if(colorField[i] === even.color) {
                    indOf = i
                }
            }
            // var indOf = colorField.indexOf(even.color);
            // kolorDeg.push(colorDegField[indOf]);
        }

        var dtLine = JSON.parse(localStorage.getItem("dataTurbine"));

        app.loading(true);
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: turbineList,
            project: fa.project,
            Color: kolor,
            // ColorDeg: kolorDeg,
            isDeviation: page.isDeviation,
            deviationVal: page.deviationVal,
            IsDownTime: page.showDownTime(),
            ViewSession: page.viewSession()
        };

        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurve", param, function(res) {
            if (!app.isFine(res)) {
                return;
            }

            var dataPowerCurves = res.data.Data;
            var dtSeries = new Array();
            if (dataPowerCurves != null) {
                if (dataPowerCurves.length > 0) {
                    dtSeries = dtLine.concat(dataPowerCurves);
                }
            } else {
                dtSeries = dtLine;
            }

            $('#powerCurve').html("");
            $("#powerCurve").kendoChart({
                theme: "flat",
                // renderAs: "canvas",
                pdf: {
                  fileName: "DetailPowerCurve.pdf",
                },
                title: {
                    text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
                    visible: false,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                legend: {
                    visible: false,
                    position: "bottom"
                },
                seriesDefaults: {
                    type: "scatterLine",
                    style: "smooth",
                },
                series: dtSeries,
                categoryAxis: {
                    labels: {
                        step: 1
                    }
                },
                valueAxis: [{
                    labels: {
                        format: "N0",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    }
                }],
                xAxis: {
                    majorUnit: 1,
                    title: {
                        text: "Wind Speed (m/s)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555",
                        visible: true,
                    },
                    labels: {
                        format: "N0",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    },
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    crosshair: {
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
                    },
                    max: 25
                },
                yAxis: {
                    title: {
                        text: "Generation (KW)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555"
                    },
                    labels: {
                        format: "N0",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    },
                    axisCrossingValue: -5,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
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
                },
                pannable: true,
                zoomable: true
            });

            app.loading(false);
            if (page.showDownTime()) {
                $('#downtime-list').show();
            } else {
                $('#downtime-list').hide();
            }
            page.ShowHideAfterInitChart();
        });
    },
    InitCurveDetail: function(turbineid,turbinename) {
        app.loading(true);
        page.detailTitle(turbinename);
        page.detailStartDate(fa.dateStart.getUTCDate() + "-" + fa.dateStart.getMonthNameShort() + "-" + fa.dateStart.getUTCFullYear());
        page.detailEndDate(fa.dateEnd.getUTCDate() + "-" + fa.dateStart.getMonthNameShort() + "-" + fa.dateEnd.getUTCFullYear());

        var colorDetail = [];

        var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'turbineid');
        var colD = _.find(dtTurbines, function(num) {
            return num.turbineid == turbineid;
        }).color;
        if (colD != undefined) {
            colorDetail.push(colD);
        }

        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: [turbineid],
            project: fa.project,
            Color: colorDetail
        };

        var dataTurbineDetail

        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getdetails", param, function(res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }

            dataTurbineDetail = res.data.Data;

            $('#powerCurveDetail').html("");
            $("#powerCurveDetail").kendoChart({
                pdf: {
                  fileName: "DetailPowerCurve.pdf",
                },
                theme: "flat",
                renderAs: "canvas",
                title: {
                    text: "Detail Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
                    visible: false,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                legend: {
                    visible: false,
                    position: "bottom"
                },
                seriesDefaults: {
                    type: "scatter",
                    style: "smooth",
                },
                series: dataTurbineDetail,
                categoryAxis: {
                    labels: {
                        step: 1
                    }
                },
                valueAxis: [{
                    labels: {
                        format: "N2",
                    }
                }],
                xAxis: {
                    majorUnit: 1,
                    title: {
                        text: "Wind Speed (m/s)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555"
                    },
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    crosshair: {
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
                    },
                    max: 25
                },
                yAxis: {
                    title: {
                        text: "Generation (KW)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555"
                    },
                    labels: {
                        format: "N2",
                    },
                    axisCrossingValue: -5,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    crosshair: {
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
                    },
                }
            });
            app.loading(false);

            $("#powerCurveDetail").data("kendoChart").refresh();

        });
    },
    InitRightTurbineList: function() {
        page.turbineList([]);
        var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'name');

        if (page.turbine().length > 1) {
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
            if(val.name != "Power Curve"){
                $("#right-turbine-list").append('<div class="btn-group">' +
                '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="page.showHideLegend(' + val.idxseries + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + val.idxseries + '"></i></button>' +
                '<input class="chk-option" type="checkbox" name="' + val.turbineid + '" checked id="chk-' + val.idxseries + '" hidden>' +
                '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="page.toDetail(\'' + val.turbineid + '\',\'' + val.turbineid + '\')" type="button">' + val.name + '</button>' +
                '</div>');
            }
        });
    },
    InitDownList: function() {
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getdownlist", "", function(res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            page.downList(res.data);

            $("#downtime-list").html("");
            $.each(page.downList(), function(idx, val) {
                $("#downtime-list").append('<div class="btn-group">' +
                    '<button class="btn btn-default btn-sm down-chk" id="down-' + val.down + '" type="button" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-down-' + val.down + '"></i></button>' +
                    '<input class="chk-option" type="checkbox" name="' + val.down + '" checked id="down-check-' + val.down + '" hidden>' +
                    '<label class="btn btn-default btn-sm turbine-btn">&nbsp;&nbsp;' + val.label + '&nbsp;&nbsp;</label>' +
                    '</div>'
                );
            });

            $('#downtime-list').hide();
        });
    },
    SetDownOnClick: function() {
        $.each($("#powerCurve").data("kendoChart").options.series, function(idx, val) {
            $.each($("#downtime-list").find('button[id^="down-"]'), function(idx2, val2) {
                if (("down-" + val.name) == val2.id) {
                    $(val2).attr('onclick', 'page.showHideDown(' + idx + ', "' + val.name + '")');
                    $(val2).find('i').css("visibility", "visible")
                    $('#down-check-' + val.name).prop('checked', true);
                }
            });
        });
    },

};

page.showHideAllLegend = function(e) {
    var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'name');
    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#powerCurve").data("kendoChart").options.series[val.idxseries].visible = true;
            }
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            if(val.idxseries > 0){
                $("#powerCurve").data("kendoChart").options.series[val.idxseries].visible = false;
            }
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    $("#powerCurve").data("kendoChart").redraw();
}

page.showHideLegend = function(idx) {
    $('#chk-' + idx).trigger('click');
    var chart = $("#powerCurve").data("kendoChart");
    // var datas = $("#powerCurve").data("kendoChart").options.series;
    // var Nama = $("#powerCurve").data("kendoChart").options.series[idx].name;

    if (page.sScater()) {
        var len = $('input[id*=chk-][type=checkbox]:checked').length;
        if (len > 3) {
            $('#chk-' + idx).prop('checked', false);
            swal('Warning', 'You can only select 3 turbines !', 'warning');
            return
        }
        Data.InitLinePowerCurve();
        // var scatterIndex = _.find(datas, function(num){ return num.name == 'Scatter-' + Nama; }).index;
        // chart._legendItemClick(scatterIndex);
    }

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

page.ShowHideAfterInitChart = function() {
    var len = $('input[id*=chk-][type=checkbox]').length;
    var chart = $("#powerCurve").data("kendoChart");
    for (var i = 1; i <= len; i++) {
        if (!$('#chk-' + i).is(':checked')) {
            // console.log(chart.options);
            chart.options.series[i].visible = false;
        }
    }
    $("#powerCurve").data("kendoChart").redraw();
}

page.hideAll = function() {
    // var chart = $("#powerCurve").data("kendoChart");
    var len = $('input[id*=chk-][type=checkbox]').length;
    for (var i = 1; i <= len; i++) {
        $('#icon-' + i).css("visibility", "hidden");
        $('#chk-' + i).prop('checked', false);
    }
}

page.resetFilter = function(){
    page.isClean(true);
    page.isDeviation(true);
    page.sScater(false);
    page.showDownTime(false);
    page.deviationVal(20);
    $('#isClean').prop('checked',true);
    $('#isDeviation').prop('checked',true);
    $('#sScater').prop('checked',false);
    $('#showDownTime').prop('checked',false);
}

page.HideforScatter = function() {
    var len = $('input[id*=chk-][type=checkbox]:checked').length;

    var sScater = page.sScater();
    if (sScater) {
        // $('#showHideChk').hide();
        $('#showHideAll').attr("disabled", true);
        $('#showHideAll').prop('checked', false);
        if (len > 3) {
            page.hideAll();
            $('#icon-1').css("visibility", "visible");
            $('#chk-1').prop('checked', true);
        }
    } else {
        // $('#showHideChk').show();
        $('#showHideAll').removeAttr("disabled");
        $('#showHideAll').prop('checked', false); /*can be hardcoded because max turbine is 3*/
    }
}

page.getSelectedFilter = function(){

    setTimeout(function(){
       $("#selectedFilter").empty();
        var deviationVal = $("#deviationValue").val();
        $('input[name="filter"]:checked').each(function() {
            if(this.value == "Deviation"){
                $("#selectedFilter").append(this.value + " < " + deviationVal + " % | ");
            }else{
                $("#selectedFilter").append(this.value + " | ");
            }
        });        
    },200);
}

$(document).ready(function() {

    page.getSelectedFilter();
    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        setTimeout(function() {   
            $("#selectedFilter").empty();
            page.getSelectedFilter();
            var project = $('#projectList').data("kendoDropDownList").value();
            var isValid = fa.LoadData();
            if(isValid) {
                app.loading(true);
                page.resetFilter();
                Data.InitLinePowerCurve();
            }
            page.project(project);
        }, 1000);
    });

    setTimeout(function() {
        $(".label-filter:contains('Turbine')" ).hide();
        $('.multiselect-native-select').hide();
        page.currProject(fa.project);
        page.project(fa.project);

        Data.LoadData();
    }, 1000);

    $('.keep-open').click(function (e) {
      e.stopPropagation()
    });

    $("input[name=isAvg]").on("change", function() {
        page.viewSession(this.id);
        Data.InitLinePowerCurve();
    });

    $('#isClean').on('click', function() {
        var isClean = $('#isClean').prop('checked');
        page.isClean(isClean);
        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#isDeviation').on('click', function() {
        var isDeviation = $('#isDeviation').prop('checked');
        page.isDeviation(isDeviation);

        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#sScater').on('click', function() {
        var sScater = $('#sScater').prop('checked');
        page.sScater(sScater);

        page.getSelectedFilter();
        page.HideforScatter();
        Data.InitLinePowerCurve();
    });

    $('#showDownTime').on('click', function() {
        var isShow = $('#showDownTime').prop('checked');
        page.showDownTime(isShow);

        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            var project = $('#projectList').data("kendoDropDownList").value();
            var lastProject = page.currProject();
            if(project != lastProject){
                page.project(lastProject)
                page.currProject(project);
            }else{
                page.project(project)
                page.currProject(project);
            }

            fa.populateTurbine(project);
         }
    });

    $('#showDownTime').attr("disabled", "disabled");
    Data.InitDownList();
});
