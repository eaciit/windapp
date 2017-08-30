'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.dataPCEachTurbine = ko.observableArray([]);
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
    .done(function(data) {
        // Save the PDF file
        kendo.saveAs({
            dataURI: data,
            fileName: "Individual-Month.pdf",
        });
    });
}

vm.currentMenu('Individual Month');
vm.currentTitle('Individual Month');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Individual Month',
    href: viewModel.appName + 'page/analyticpcmonthly'
}]);

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

page.LoadData = function() {
    fa.LoadData();
    app.loading(true);
    setTimeout(function () {
        var param = {
            turbine: fa.turbine(),
            project: fa.project,
        };
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getlistpowercurvemonthly", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            if (res.data.Data != null) {
                page.dataPCEachTurbine(res.data.Data);
                page.InitLinePowerCurve();
            }
            if (res.data.Category != null) {
                listOfCategory = res.data.Category;
                $("#legend-list").html("");
                listOfButton = {};
                $.each(listOfCategory, function (idx, val) {
                    if(idx > 0) {
                        var idName = "btn" + idx;
                        listOfButton[idName] = true;
                        $("#legend-list").append(
                            '<button id="' + idName + 
                            '" class="btn btn-default btn-sm btn-legend" type="button" onclick="page.showHideLegend(' + idx + ')" style="border-color:' + 
                            val.color + ';background-color:' + val.color + ';"></button>' +
                            '<span class="span-legend">' + val.category + '</span>'
                        );
                    }
                });
            }
            app.loading(false);
        })
    }, 300);
}

page.InitLinePowerCurve = function() {
    fa.LoadData();
    listOfChart = [];

    di.getAvailDate();


    $.each(page.dataPCEachTurbine(), function (i, dataTurbine) {
        var name = dataTurbine.Name
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

$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        fa.checkTurbine();
        page.LoadData();
    });
    $(".period-list").hide();
    $(".filter-date-start").hide();
    $(".filter-date-end").hide();
    $(".label-filter")[3].remove();
    $(".label-filter")[2].remove();

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            fa.currentFilter().project = this._old;
            fa.checkFilter();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
            di.getAvailDate();
         }
    });

    page.LoadData();
});