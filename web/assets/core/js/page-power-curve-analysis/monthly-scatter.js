'use strict';
var ms = {};

ms.dataPCEachTurbine = ko.observableArray([]);
ms.PrintPdf = ko.observable(false);
ms.dataAvail = ko.observable();
ms.project = ko.observable();
ms.dateStart = ko.observable();
ms.dateEnd = ko.observable();
ms.isValid = ko.observable(true);
ms.LastFilter;
ms.TableName;
ms.FieldList;
ms.ContentFilter;

ms.loadFirstTime = ko.observable(true);
ms.reset = function(){
    ms.loadFirstTime(true);
}
ms.refresh = function() {
    if(ms.loadFirstTime()) {
        var currTime = $('#dateStart').data('kendoDatePicker').value();
        var dateToFilter = moment(currTime).add(-1, 'months');
        var newFilterValue = new Date(Date.UTC(dateToFilter.year(), dateToFilter.month(), 1, 0, 0, 0, 0));
        $('#ms-period').kendoDatePicker({
            start: "year",
            depth: "year",
            format: "MMM-yyyy",
            dateInput: true,
            value: newFilterValue,
            max: $('#dateEnd').data('kendoDatePicker').value(),
        });

        $("input[name=isValid]").on("change", function() {
            ms.isValid(true);
            if(this.id == "alldata"){
                ms.isValid(false);
            }
            ms.LoadData();
        });

        ms.LoadData();
        ms.loadFirstTime(false);
    }
}
ms.internalRefresh = function() {
    ms.LoadData();
}

var msListOfChart = [];
var msListOfButton = {};
var msListOfCategory = [];

ms.ExportIndividualMonthPdf = function() {
    kendo.drawing.drawDOM($(".individual-month"))
    .then(function(group) {
        // Render the result as a PDF file
        return kendo.drawing.exportPDF(group, {
            paperSize: "auto",
            margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" }
        });
    })
    .done(function() {
        // Save the PDF file
        kendo.saveAs({
            dataURI: data,
            fileName: "Individual-Month.pdf",
        });
    });
}

ms.getPDF = function(selector){
    var title = fa.project + " | "+ kendo.toString($('#dateStart').data('kendoDatePicker').value(), 'MMM-yyyy');
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    app.loading(true);
    ms.PrintPdf(true);
    $("#illusion-month").append($(".individual-month").html());
    $("#pdf-title").text(title);
    var dateStart = moment($('#dateStart').data('kendoDatePicker').value()).format("DD MMM YYYY");
    var project = $("#projectList").data("kendoDropDownList").value();
    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            margin: {
                left   : "10mm",
                top    : "10mm",
                right  : "10mm",
                bottom : "10mm"
            },
            multiPages: true,
        });
     kendo.drawing.pdf.saveAs(group, project+"PCMonthlyScatter"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            $("#illusion-month").empty();
            ms.PrintPdf(false);
            app.loading(false);
        },2000)
    });
}

ms.showHideLegend = function (index) {
    var idName = "btn" + index;
    msListOfButton[idName] = !msListOfButton[idName];
    if (msListOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': msListOfCategory[index].color, 'border-color': msListOfCategory[index].color });
    }
    $.each(msListOfChart, function (idx, idChart) {
        $(idChart).data("kendoChart").options.series[index].visible = msListOfButton[idName];
        $(idChart).data("kendoChart").refresh();
    });
}

ms.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: ms.LastFilter,
        FieldList: ms.FieldList,
        Tablename: ms.TableName,
        TypeExcel: namaFile,
        ContentFilter: ms.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: ms.LastFilterDetails,
            FieldList: ms.FieldListDetails,
            Tablename: ms.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: ms.ContentFilterDetails,
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

ms.LoadData = function() {
    app.loading(true);
    setTimeout(function () {
        var param = {
            engine : fa.engine,
            turbine: fa.turbine(),
            project: fa.project,
            datestart: $('#ms-period').data('kendoDatePicker').value(),
            isclean: ms.isValid()
        };
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getlistpowercurvemonthlyscatter", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            ms.LastFilter = res.data.LastFilter;
            ms.FieldList = res.data.FieldList;
            ms.TableName = res.data.TableName;
            ms.ContentFilter = res.data.ContentFilter;
            if (res.data.Data != null) {
                ms.dataPCEachTurbine(_.sortBy(res.data.Data, 'Name'));
                ms.InitLinePowerCurve();
            }
            app.loading(false);
        });
        
    }, 500);
}

ms.InitLinePowerCurve = function() {

    msListOfChart = [];
    $.each(ms.dataPCEachTurbine(), function (i, dataTurbine) {
        var name = dataTurbine.Name
        var idDataAvailability = "#ms-dataAv-" + dataTurbine.Name
        var idTotalAvailability = "#ms-totalAv-" + dataTurbine.Name
        var idChart = "#ms-chart-" + dataTurbine.Name
        msListOfChart.push(idChart);
        var rotation = 300;
        var heightVal = 250;
        var isPannable = false;
        var isZoomable = false;
        var isTitle = true;
        var titleFont = '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var titleAxisFont = '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var labelAxisFont = '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var tooltipAxisFont = '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';

        if (fa.turbine().length == 1) {
            rotation = 0;
            heightVal = 400;
            isPannable = false;
            isZoomable = false;
            isTitle = false;
            var titleFont = '20px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var titleAxisFont = '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var labelAxisFont = '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var tooltipAxisFont = 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            $(".ms-power-curve-item").removeClass("col-md-4");
        }
        
        $(idDataAvailability).html(kendo.toString(dataTurbine.DataAvailability,'p1'));        
        $(idTotalAvailability).html(kendo.toString(dataTurbine.DataTotalAvailability,'p1'));

        $(idChart).html("");
        $(idChart).kendoChart({
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            title: {
                text: name,
                font: titleFont,
                visible: isTitle
            },
            legend: {
                position: "bottom",
                visible: false,
            },
            chartArea: {
                height: heightVal
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
            series: dataTurbine.Data,
            categoryAxis: {
                labels: {
                    step: 1
                }
            },
            valueAxis: [{
                labels: {
                    format: "N0",
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: titleAxisFont,
                    color: "#585555",
                    visible: true,
                },
                labels: {
                    font: labelAxisFont,
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: tooltipAxisFont,
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
                    font: titleAxisFont,
                    color: "#585555"
                },
                labels: {
                    format: "N0",
                    font: labelAxisFont,
                    rotation: rotation,
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
                        font: tooltipAxisFont,
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
                font: tooltipAxisFont,
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
            pannable: isPannable,
            zoomable: isZoomable
        });
        $(idChart).data("kendoChart").refresh();
    });
}

function sticky_relocate() {
    var window_top = $(window).scrollTop();
    var div_top = $('#legend-anchor').offset().top;
    if (window_top > div_top) {
        $('#legend-list').addClass('legend');
        $('#legend-anchor').height($('#legend-list').outerHeight());
    } else {
        $('#legend-list').removeClass('legend');
        $('#legend-anchor').height(0);
    }
}