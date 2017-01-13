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
    { "value": "temp", "text": "Temperature Analysis" },
    { "value": "pitch", "text": "Pitch Angle" },
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
    fa.LoadData();
    page.getPowerCurveScatter();
}

page.setAxis = function(name, title, crossVal, minVal, maxVal) {
    var result = {
        name: name,
        title: {
            text: title,
            font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            color: "#585555"
        },
        labels: {
            format: "N2",
        },
        // axisCrossingValue: crossVal,
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        min: minVal,
        max: maxVal,
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
    return result
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
        scatterType: page.scatterType,
    };
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var minDatetemp = new Date(res.data.ScadaData[0]);
        var maxDatetemp = new Date(res.data.ScadaData[1]);
        $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
    });

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurvescatter", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dtSeries = res.data.Data;
        
        var yAxes = [];
        var yAxis = page.setAxis("powerAxis", "Generation (KW)", -5, 0, 2500);
        yAxes.push(yAxis);
        switch(page.scatterType) {
            case "temp":
                var axis = page.setAxis("tempAxis", "Temperature (Celcius)", -5, 0, 90);
                yAxes.push(axis);
                break;
            case "pitch":
                var axis = page.setAxis("pitchAxis", "Angle (Degree)", -5, 0, 360);
                yAxes.push(axis);
                break;
        }

        $('#scatterChart').html("");
        $("#scatterChart").kendoChart({
            theme: "flat",
            renderAs: "canvas",
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
                },
                axisCrossingValues: [0, 30],
            },
            valueAxis: [{
                labels: {
                    format: "N2",
                }
            }],
            // valueAxes: yAxes,
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
            yAxis: yAxes,
            pannable: true,
            zoomable: true
        });
        app.loading(false);
    });
}

$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        page.LoadData();
    });

    setTimeout(function(){
        page.LoadData();
    }, 300);
});