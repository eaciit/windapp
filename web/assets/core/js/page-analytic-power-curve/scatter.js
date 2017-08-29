'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

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
            fileName: "Power-Curve-Scatter.pdf",
        });
    });
}

page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "temp", "text": "Nacelle Temperature" },
    { "value": "deviation", "text": "Nacelle Deviation" },
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "ambient", "text": "Ambient Temperature" },
    // { "value": "mainbearing", "text": "Temp Main Bearing" },
    // { "value": "gearbox", "text": "Temp  Gearbox HSS De" },
]);

vm.currentMenu('Scatter');
vm.currentTitle('Scatter');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Scatter',
    href: viewModel.appName + 'page/analyticpcscatter'
}]);

page.LoadData = function() {
    var isValid = fa.LoadData();
    if (isValid) {
        page.getPowerCurveScatter();
    }
}

page.setAxis = function(name, title) {
    var result = {
        name: name,
        title: {
            text: title,
            font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            color: "#585555"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        labels: {
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
    }
    if(name == "powerAxis") {
        result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # kW";
        result.crosshair.tooltip.padding = {left:5};
    } else {
        switch(page.scatterType) {
            case "temp":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176) + "C";
                break;
            case "deviation":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176);
                break;
            case "pitch":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176);
                break;
            case "ambient":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176) + "C";
                break;
        }
    }
    return result
}

page.refreshChart = function() {
    page.LoadData();
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));   

    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        scatterType: page.scatterType,
    };

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurvescatter", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dtSeries = res.data.Data;
        
        var yAxes = [];
        var yAxis = page.setAxis("powerAxis", "Generation (KW)");
        yAxes.push(yAxis);
        switch(page.scatterType) {
            case "temp":
                var axis = page.setAxis("tempAxis", "Temperature (Celsius)");
                yAxes.push(axis);
                break;
            case "deviation":
                var axis = page.setAxis("deviationAxis", "Wind Direction (Degree)");
                yAxes.push(axis);
                break;
            case "pitch":
                var axis = page.setAxis("pitchAxis", "Angle (Degree)");
                yAxes.push(axis);
                break;
            case "ambient":
                var axis = page.setAxis("ambientAxis", "Temperature (Celcius)");
                yAxes.push(axis);
                break;
        }

        $('#scatterChart').html("");
        $("#scatterChart").kendoChart({
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
                position: "bottom",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            seriesDefaults: {
                type: "scatterLine",
                style: "smooth",
            },
            series: dtSeries,
            categoryAxis: {
                labels: {
                    step: 1,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
            },
            valueAxis: [{
                labels: {
                    format: "N2",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555"
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                axisCrossingValues: [0, 30],
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        template: "#= kendo.toString(value, 'n2') # m/s",
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
            yAxes: yAxes
        });
        app.loading(false);
    });
}

$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        setTimeout(function(){
            page.LoadData();
        }, 300);
    });

    $.when(di.getAvailDate()).done(function() {
        page.LoadData();
    });
});