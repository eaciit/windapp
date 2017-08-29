'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "rotor", "text": "Rotor RPM" },
]);

vm.currentMenu('Operational Power Curve');
vm.currentTitle('Operational Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Operational Power Curve',
    href: viewModel.appName + 'page/analyticpcscatteroperational'
}]);

page.LoadData = function() {
    var isValid = fa.LoadData();
    if(isValid) {
        page.getPowerCurveScatter();
    }
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
    
    

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatteroperational", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dtSeries = res.data.Data;
        
        var minAxisY = res.data.MinAxisY;
        var maxAxisY = res.data.MaxAxisY;
        var minAxisX = res.data.MinAxisX;
        var maxAxisX = res.data.MaxAxisX;
        var name = '';
        var title = '';
        var xAxis = {};
        var measurement = '';
        var format = 'N0'
        if(maxAxisX - minAxisX < 7) {
            format = 'N2'
        }
        switch(page.scatterType) {
            case "pitch":
                name = 'pitchAxis'
                title = 'Angle (Degree)'
                measurement = String.fromCharCode(176)
                break;
            case "rotor":
                name = "rotorAxis"
                title = "Revolutions per Minute (RPM)";
                measurement = 'rpm'
                break;
        }
        xAxis = {
            name: name,
            title: {
                text: title,
                visible: true,
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                format: format
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
            crosshair: {
                visible: true,
                tooltip: {
                    visible: true,
                    format: "N2",
                    template: "#= kendo.toString(value, 'n2') # " + measurement,
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                }
            },
            // majorUnit: 0.5,
            min: minAxisX,
            max: maxAxisX,
        }
        var yAxis = {};
        yAxis = {
            name: "powerAxis",
            title: {
                text: "Generation (KW)",
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
                color: "#eee",
                width: 0.8,
            },
            crosshair: {
                visible: true,
                tooltip: {
                    visible: true,
                    format: "N2",
                    template: "#= kendo.toString(value, 'n2') # kWh",
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                }
            },
            // majorUnit: 0.5,
            min: minAxisY,
            max: maxAxisY,
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
                    step: 1
                },
            },
            valueAxis: [{
                labels: {
                    format: "N2",
                }
            }],
            xAxis: xAxis,
            yAxes: yAxis
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