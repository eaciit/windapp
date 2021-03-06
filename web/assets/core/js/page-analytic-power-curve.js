'use strict';


viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.turbineList = ko.observableArray([]);
page.downList = ko.observableArray([]);
page.dtLineChart = ko.observableArray([]);
page.projectList = ko.observableArray([]);

page.isMain = ko.observable(true);
page.isDetail = ko.observable(false);
page.detailTitle = ko.observable("");
page.detailStartDate = ko.observable("");
page.detailEndDate = ko.observable("");

page.isSpecific = ko.observable(true);
page.isClean = ko.observable(true);
page.isPower0 = ko.observable(false); // to show all data, even power less than 0
page.idName = ko.observable("");
page.isDeviation = ko.observable(true);
page.sScater = ko.observable(false);
page.showDownTime = ko.observable(false);
page.deviationVal = ko.observable("20");

page.isDensity = ko.observable(false);
page.dataAvail = ko.observable(0.0);
page.dataAvailAll = ko.observable(0.0);
page.totalAvail = ko.observable(0.0);
page.totalAvailAll = ko.observable(0.0);
page.viewName = ko.observable();
page.LastFilter;
page.TableName;
page.FieldList;
page.ContentFilter;
page.LastFilterDetails;
page.TableNameDetails;
page.FieldListDetails;
page.ContentFilterDetails;

page.totalAvailTurbines = ko.observableArray([]);

// add by ams Aug 11, 2017
page.deviationOpts = ko.observableArray([
    { "value": 0, "text": "<" },
    { "value": 1, "text": ">" },
])
page.deviationOpr = ko.observable(0);

page.viewSession = ko.observable("");
page.turbine = ko.observableArray([]);
page.powerCurveOptions = ko.observable();
page.currProject = ko.observable();
page.project = ko.observable();
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.ss_airdensity = ko.observable(0.0);
page.std_airdensity = ko.observable(0.0);
var lastParam;
var lastParamDetail;

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

page.getPDF = function(selector, detail){
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
    var title = project+"PowerCurve"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf";
    if(detail == true){
        title = project+"_"+page.detailTitle()+"DetailPowerCurve"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf";
    }
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
      kendo.drawing.pdf.saveAs(group, title);
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

page.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: page.LastFilter,
        FieldList: page.FieldList,
        Tablename: page.TableName,
        TypeExcel: namaFile,
        ContentFilter: page.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: page.LastFilterDetails,
            FieldList: page.FieldListDetails,
            Tablename: page.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: page.ContentFilterDetails,
            IsSplittedSheet: isSplittedSheet,
            IsMultipleProject: isMultipleProject,
        };
    }

    var urlName = viewModel.appName + "analyticpowercurve/genexcelpowercurve";
    app.ajaxPost(urlName, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        window.location = viewModel.appName + "/".concat(res.data);
        app.loading(false);
    });
}


page.ExportPowerCurvePdf = function() {
    var chart = $("#powerCurve").getKendoChart();
    var container = $('<div />').css({
        position: 'absolute',
        top: 0,
        left: -1500,
      }).appendTo('body');
    
      var dateStart = moment(lastParam.dateStart).format("DD MMM YYYY");
      var dateEnd = moment(lastParam.dateEnd).format("DD MMM YYYY");

      var options = chart.options;

      var exportOptions ={
            // Custom settings for export            
            legend: {
              visible: true
            },
            title:{
                text: "Power Curve | " + dateStart + " until " + dateEnd + " | " + lastParam.project,
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
      
      $("#powerCurve").kendoChart($.extend(true, options, {legend: {visible: false},title:{visible: false},chartArea: { height: 425 }, render: function(e){return false}}));
}
page.ExportPowerCurveDetailPdf = function() {
        var chart = $("#powerCurveDetail").getKendoChart();
        var container = $('<div />').css({
            position: 'absolute',
            top: 0,
            left: -1500
          }).appendTo('body');


        var options = chart.options;
        var dateStart = moment(lastParamDetail.dateStart).format("DD MMM YYYY");
        var dateEnd = moment(lastParamDetail.dateEnd).format("DD MMM YYYY");

          var exportOptions ={
                // Custom settings for export
                legend: {
                  visible: true
                },
                title:{
                    text: "Power Curve Detail | " + dateStart + " until " + dateEnd + " | " + lastParamDetail.project,
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

          $("#powerCurveDetail").kendoChart($.extend(true, options, {legend: {visible: false},title:{visible: false},chartArea: { height: 425 }, render: function(e){return false}}));
    }

vm.currentMenu('Power Curve');
vm.currentTitle('Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
},{
    title: "Power Curve",
    href: '#'
},{
    title: 'Power Curve',
    href: viewModel.appName + 'page/analyticpowercurve'
}]);

var dataPowerCurve
var dataTurbine

var Data = {
    LoadData: function() {
        page.populateTurbine();
        this.InitLinePowerCurve();
        this.InitRightTurbineList();
    },
    InitLinePowerCurve: function() {
        page.getSelectedFilter();
        
        page.deviationOpr($("#deviationOpr").val());
        page.deviationVal($("#deviationValue").val());

        var link = "analyticpowercurve/getlistpowercurvescada"

        app.loading(true);

        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();   

        var param = {
            period: fa.period,
            dateStart: dateStart,
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: $("#turbineList").val(),
            project: fa.project,
            isClean: page.isClean,
            isSpecific: page.isSpecific,
            isDeviation: page.isDeviation,
            isPower0: page.isPower0,
            DeviationVal: page.deviationVal,
            DeviationOpr: page.deviationOpr,
            ViewSession: page.viewSession,
            Engine: fa.engine,
        };
        lastParam = param;

        toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }

            page.totalAvail(res.data.TotalDataAvail);
            page.totalAvailAll(res.data.TotalDataAvail);
            page.totalAvailTurbines(res.data.TotalPerTurbine);
            page.LastFilter = res.data.LastFilter;
            page.FieldList = res.data.FieldList;
            page.TableName = res.data.TableName;
            page.ContentFilter = res.data.ContentFilter;

            var tempData = [];
            var powerCurveData;
            res.data.Data.forEach(function(val, idx){
                if(val.name != "Power Curve") {
                    tempData.push(val);
                } else {
                    powerCurveData = val;
                }
            });

            tempData = _.sortBy(tempData, 'name')
            tempData.forEach(function(val, idx){
                tempData[idx].idxseries = idx+1;
            });
            tempData.push(powerCurveData);
            res.data.Data = tempData;


            dataTurbine = res.data.Data;
            localStorage.setItem("dataTurbine", JSON.stringify(dataTurbine));
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
                    height: 425,
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
                pannable: {
                    lock: "y"
                },
                zoomable: {
                    mousewheel: {
                        lock: "y"
                    },
                    selection: {
                        lock: "y",
                        key: "none",
                    }
                },
                dataBound : function(){
                    page.getLegendActive();
                },
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
            Data.InitRefreshValueAvailability();
        });
        
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

        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

        var param = {
            period: fa.period,
            dateStart: dateStart,
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: turbineList,
            project: fa.project,
            Color: kolor,
            // ColorDeg: kolorDeg,
            isClean: page.isClean,
            isSpecific: page.isSpecific,
            isDeviation: page.isDeviation,
            isPower0: page.isPower0,
            deviationVal: page.deviationVal,
            DeviationOpr: page.deviationOpr,
            IsDownTime: page.showDownTime(),
            ViewSession: page.viewSession()
        };
        // var param = {
        //     period: fa.period,
        //     dateStart: dateStart,
        //     dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
        //     turbine: fa.turbine(),
        //     project: fa.project,
        //     Color: kolor,
        //     isClean: page.isClean,
        //     isSpecific: page.isSpecific,
        //     isDeviation: page.isDeviation,
        //     isPower0: page.isPower0,
        //     deviationVal: page.deviationVal,
        //     DeviationOpr: page.deviationOpr,
        //     IsDownTime: page.showDownTime(),
        //     ViewSession: page.viewSession
        // }
        lastParam = param;

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
                    align: "center",
                    position: "bottom",

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
                pannable: {
                    lock: "y"
                },
                zoomable: {
                    mousewheel: {
                        lock: "y"
                    },
                    selection: {
                        lock: "y",
                        key: "none",
                    }
                },                
                dataBound : function(){
                    page.getLegendActive();
                },
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

        var dateStart = lastParam.dateStart;
        var dateEnd = lastParam.dateEnd;

        page.detailStartDate(dateStart.getUTCDate() + "-" + dateStart.getMonthNameShort() + "-" + dateStart.getUTCFullYear());
        page.detailEndDate(dateEnd.getUTCDate() + "-" + dateStart.getMonthNameShort() + "-" + dateEnd.getUTCFullYear());

        var colorDetail = [];

        var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'turbineid');
        var colD = _.find(dtTurbines, function(num) {
            return num.turbineid == turbineid;
        }).color;
        if (colD != undefined) {
            colorDetail.push(colD);
        }

        var param = {
            period: lastParam.period,
            dateStart: dateStart,
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: [turbineid],
            project: lastParam.project,
            Color: colorDetail
        };
        lastParamDetail = param;

        var dataTurbineDetail;

        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getdetails", param, function(res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }

            dataTurbineDetail = res.data.Data;
            page.LastFilterDetails = res.data.LastFilter;
            page.FieldListDetails = res.data.FieldList;
            page.TableNameDetails = res.data.TableName;
            page.ContentFilterDetails = res.data.ContentFilter;

            $('#powerCurveDetail').html("");
            $("#powerCurveDetail").kendoChart({
                pdf: {
                  fileName: "DetailPowerCurve.pdf",
                },
                theme: "flat",
                renderAs: "canvas",
                title: {
                    // text: "Detail Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
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
                },
                dataBound : function(){
                    var chart = $("#powerCurveDetail").data("kendoChart");
                    var viewModel = kendo.observable({
                      series: chart.options.series,
                      markerColor: function(e) {
                        return e.get("visible") ? e.color : "grey";
                      }
                    });

                    kendo.bind($("#legendPowerCurveDetail"), viewModel);
                }
            });
            app.loading(false);

            $("#powerCurveDetail").data("kendoChart").refresh();

        });
    },
    InitRightTurbineList: function() {
        page.turbineList([]);
        var dtTurbines = JSON.parse(localStorage.getItem("dataTurbine"));

        if (page.turbine().length > 1) {
            $("#showHideChk").html('<label>' +
                '<input type="checkbox" id="showHideAll" checked onclick="page.showHideAllLegend(this)" >' +
                '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
                '<span id="labelShowHide"><b>Select All</b></span>' +
                '</label>');
        } else {
            $("#showHideChk").html("");
        }

        var totalDataShoulBeInProject = 0;
        var totalDataAvailInProject = 0;
        $("#right-turbine-list").html("");
        $.each(dtTurbines, function(idx, val) {
            if(val.name != "Power Curve"){
                var nameTurbine = val.name;
                 if ( fa.project == "Rajgarh" ) {
                    nameTurbine = nameTurbine.replace("KH-", "-")
                }
                
                totalDataShoulBeInProject += val.totaldatashouldbe;
                totalDataAvailInProject += val.totaldata;
                $("#right-turbine-list").append('<div class="btn-group">' +
                '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="page.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + idx + '"></i></button>' +
                '<input class="chk-option" type="checkbox" name="' + val.turbineid + '" checked id="chk-' + idx + '" hidden>' +
                '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="page.toDetail(\'' + val.turbineid + '\',\'' + val.turbineid + '\')" type="button">' + nameTurbine + ' <label id="dataavailpct-'+val.turbineid+'" class="label label-default pull-right" data-toggle="tooltip" title="Data available for turbine : '+ nameTurbine +'">'+ kendo.toString(val.dataavailpct, 'p1') +'</label></button>' +
                '</div>');
            }
        });
        page.dataAvail((totalDataAvailInProject / totalDataShoulBeInProject));
        page.dataAvailAll((totalDataAvailInProject / totalDataShoulBeInProject));
    },
    InitRefreshValueAvailability : function(){
        var dtTurbines = JSON.parse(localStorage.getItem("dataTurbine"));
        $.each(dtTurbines, function(idx, val) {
            var elm = $("#right-turbine-list").find($("#dataavailpct-"+val.turbineid));

            if(elm.length > 0){
                elm.text(kendo.toString(val.dataavailpct, 'p1'));
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

page.getLegendActive = function(){
    var chart = $("#powerCurve").data("kendoChart");
    var viewModel = kendo.observable({
      series: chart.options.series,
      markerColor: function(e) {
        return e.get("visible") ? e.color : "grey";
      }
    });

    kendo.bind($("#legend"), viewModel);
}
page.showHideAllLegend = function(e) {
    var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'name');
    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            val.idxseries = val.idxseries - 1;
            if(val.name !== "Power Curve"){
                $("#powerCurve").data("kendoChart").options.series[val.idxseries].visible = true;
            }
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            val.idxseries = val.idxseries - 1;
            if(val.name !== "Power Curve"){
                $("#powerCurve").data("kendoChart").options.series[val.idxseries].visible = false;
            }
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    $("#powerCurve").data("kendoChart").redraw();
    page.getLegendActive();
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

    // check if turbines not all checked
    if (!$('#showHideAll').is(':checked')) {
        var chks = $('input[id*=chk-][type=checkbox]:checked');
        var totalavail = 0;
        var totalCount = 0;
        var sampleAvail = 0;
        var sampleCount = 0;
        
        $.each(chks, function(idx, elm){
            var tbName = $(elm).attr('name');
            var tbAvail = page.totalAvailTurbines()[tbName];
            totalavail += tbAvail.avail;
            totalCount++;

            var elmAvail = $(elm).parent().find('button.wbtn').find('label').text().replace(' %', '');
            var currElmAvail = parseFloat(elmAvail);
            sampleAvail += currElmAvail; 
            sampleCount++;
        });

        if(totalCount > 0) {
            var selectedAvail = totalavail / totalCount;
            page.totalAvail(selectedAvail);

            var selectedSampleAvail = (sampleAvail / 100) / sampleCount;
            page.dataAvail(selectedSampleAvail);
        } else {
            page.totalAvail(page.totalAvailAll());
            page.dataAvail(page.dataAvailAll());

        }
    }

    chart._legendItemClick(idx);
    page.getLegendActive();
}

page.ShowHideAfterInitChart = function() {
    var len = $('input[id*=chk-][type=checkbox]').length;
    var chart = $("#powerCurve").data("kendoChart");
    for (var i = 0; i < len; i++) {
        if (!$('#chk-' + i).is(':checked')) {
            // console.log(chart.options);
            chart.options.series[i].visible = false;
        }
    }
    $("#powerCurve").data("kendoChart").redraw();
    page.getLegendActive();
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
    page.isSpecific(true);
    page.isDeviation(true);
    page.isPower0(false);
    page.sScater(false);
    page.showDownTime(false);
    page.deviationVal("20");
    page.isDensity(false);
    $('#isClean').prop('checked',true);
    $('#isSpecific').prop('checked',false);
    $('#isPower0').prop('checked',false);
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
        var delim = "";
        $('input[name="filter"]:checked').each(function() {
            if(this.value == "Deviation"){
                $("#selectedFilter").append(delim + this.value + ($("#deviationOpr").val()=="0"?" < ":" > ") + deviationVal + " % ");
            }else if(this.value == "Specific"){
                 var value = "Site Specific PW"
                 $("#selectedFilter").append(delim + value + " ");
            }else{
                $("#selectedFilter").append(delim + this.value + " ");
            }

            delim = "| ";
        });      
        if(delim == "") {
            $("#selectedFilter").append("No Filter Selected.");
        }  
    },200);
}

page.CheckDeviationValue = function(elm) {
    var value = $(elm).val();
    if(value=='') {
        $(elm).val(0);
    }
}

$(document).ready(function() {
    di.getAvailDate();
    page.getSelectedFilter();

    $('#pc-filter-density').hide();
    $('#pc-filter-downtime').hide();
    
    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        setTimeout(function() {   
            $("#selectedFilter").empty();
            page.getSelectedFilter();
            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

            var isValid = fa.LoadData();
            if(isValid) {
                app.loading(true);
                page.resetFilter();
                Data.InitLinePowerCurve();
            }

            page.project(project);
            page.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            page.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

            var getAd = _.find(page.projectList(), function(p) {
                return p.ProjectId == project
            });
            if(getAd!=undefined) {
                page.ss_airdensity(getAd.SS_AirDensity);
                page.std_airdensity(getAd.STD_AirDensity);
            }
        }, 1000);
    });

    setTimeout(function() {
        // $(".label-filter:contains('Turbine')" ).hide();
        // $('.multiselect-native-select').hide();
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

        page.currProject(fa.project);

        page.project(fa.project);
        page.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
        page.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
        page.viewName($('input[name=isAvg]:checked').parent('label').text());

        var getAd = _.find(page.projectList(), function(p) {
            return p.ProjectId == fa.project
        });
        if(getAd!=undefined) {
            page.ss_airdensity(getAd.SS_AirDensity);
            page.std_airdensity(getAd.STD_AirDensity);
        }
        fa.LoadData();
        Data.LoadData();
    }, 1500);

    $('.keep-open').click(function (e) {
      e.stopPropagation()
    });

    $("input[name=isAvg]").on("change", function() {
        // if(this.id == "density"){
        //     $('#isSpecific').attr("disabled", "disabled");
        //     $('#isSpecific').prop('checked',false);
        // }else{
        //     $('#isSpecific').removeAttr('disabled');
        // }
        page.viewName($('input[name=isAvg]:checked').parent('label').text());
        page.isSpecific(true);
        $('#pc-filter-density').toggle();
        if(this.id == "sitespesific"){
            $('#isDensity').attr("disabled", "disabled");
            $('#isDensity').prop('checked',false);
        }else{
            page.isSpecific(false);
            $('#isDensity').removeAttr('disabled');
        }

        page.viewSession(this.id);
        Data.InitLinePowerCurve();
    });

    $('#isClean').on('click', function() {
        var isClean = $('#isClean').prop('checked');
        page.isClean(isClean);
        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#isPower0').on('click', function() {
        var isPower0 = $('#isPower0').is(':checked');
        page.isPower0(isPower0);
        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#isSpecific').on('click', function() {
        var isSpecific = $('#isSpecific').prop('checked');
        page.isSpecific(isSpecific);
        page.getSelectedFilter();
        Data.InitLinePowerCurve();
    });

    $('#isDensity').on('click', function() {
        var isDensity = $('#isDensity').prop('checked');
        page.isDensity(isDensity);
        if(isDensity) {
            page.viewSession("density");
        } else {
            page.viewSession("standardpc");
        }
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

        $('#pc-filter-downtime').toggle();

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
            setTimeout(function(){
                fa.currentFilter().project = this._old;
                fa.checkFilter();
                fa.disableRefreshButton(true);
                var project = $('#projectList').data("kendoDropDownList").value();
                var lastProject = page.currProject();

                // fa.populateTurbine(project);
                var projectName = $('#projectList').data("kendoDropDownList").dataItem().value;

                fa.populateEngine(projectName);

                di.getAvailDate();
                if(project != lastProject){
                    page.project(lastProject)
                    page.currProject(project);
                }else{
                    page.project(project)
                    page.currProject(project);
                }

               fa.disableRefreshButton(false);
            },500);
         }
    });

    $('#showDownTime').attr("disabled", "disabled");
    Data.InitDownList();
});
page.deviationOpr.subscribe(Data.InitLinePowerCurve, this);
