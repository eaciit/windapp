'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.dataPCEachTurbine = ko.observableArray([]);
page.PrintPdf = ko.observable(false);
page.dataAvail = ko.observable();
page.project = ko.observable();
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.LastFilter;
page.TableName;
page.FieldList;
page.ContentFilter;

var listOfChart = [];
var listOfButton = {};
var listOfCategory = [];

page.ExportIndividualMonthPdf = function() {
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

vm.currentMenu('Monthly Scatter');
vm.currentTitle('Monthly Scatter');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Monthly Scatter',
    href: viewModel.appName + 'page/analyticpcmonthlyscatter'
}]);


page.getPDF = function(selector){
    var title = fa.project + " | "+ kendo.toString($('#dateStart').data('kendoDatePicker').value(), 'MMM-yyyy');
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    app.loading(true);
    page.PrintPdf(true);
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
            page.PrintPdf(false);
            app.loading(false);
        },2000)
    });
}

page.showHideLegend = function (index) {
    var idName = "btn" + index;
    listOfButton[idName] = !listOfButton[idName];
    if (listOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': listOfCategory[index].color, 'border-color': listOfCategory[index].color });
    }
    $.each(listOfChart, function (idx, idChart) {
        $(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
        $(idChart).data("kendoChart").refresh();
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

page.LoadData = function() {
    app.loading(true);
    setTimeout(function () {
        var param = {
            engine : fa.engine,
            turbine: fa.turbine(),
            project: fa.project,
            datestart: $('#dateStart').data('kendoDatePicker').value(),
        };
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getlistpowercurvemonthlyscatter", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            page.LastFilter = res.data.LastFilter;
            page.FieldList = res.data.FieldList;
            page.TableName = res.data.TableName;
            page.ContentFilter = res.data.ContentFilter;
            if (res.data.Data != null) {
                page.dataPCEachTurbine(_.sortBy(res.data.Data, 'Name'));
                page.InitLinePowerCurve();
            }
            app.loading(false);
        });
        
    }, 500);
}

page.InitLinePowerCurve = function() {

    listOfChart = [];
    $.each(page.dataPCEachTurbine(), function (i, dataTurbine) {
        var name = dataTurbine.Name
        var idDataAvailability = "#dataAv-" + dataTurbine.Name
        var idTotalAvailability = "#totalAv-" + dataTurbine.Name
        var idChart = "#chart-" + dataTurbine.Name
        listOfChart.push(idChart);
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
            $(".power-curve-item").removeClass("col-md-4");
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

$(function() {
    $(window).scroll(sticky_relocate);
    sticky_relocate();

    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        page.LoadData();
    });

    $(".period-list").hide();
    // $(".filter-date-start").hide();
    $(".filter-date-end").hide();
    $(".label-filter")[3].remove();
    // $(".label-filter")[2].remove();

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            fa.currentFilter().project = this._old;
            fa.checkFilter();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.project = project;
            fa.populateEngine(project);
            di.getAvailDate();
         }
    });

    di.getAvailDate();
    fa.LoadData();
    setTimeout(function(){
        $("#periodList").data("kendoDropDownList").value("monthly");
        $("#periodList").data("kendoDropDownList").trigger("change");
        page.LoadData();
    },700);
    
});
