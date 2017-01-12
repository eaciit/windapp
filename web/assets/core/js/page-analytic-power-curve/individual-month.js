'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

page.idName = ko.observable("");
page.powerCurveOptions = ko.observable();

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

vm.currentMenu('Power Curve');
vm.currentTitle('Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Individual Monthly',
    href: viewModel.appName + 'page/analyticpcmonthly'
}]);

var dataPowerCurve
var dataTurbine

page.dataPCEachTurbine = ko.observableArray([]);
var listOfChart = [];

page.LoadData = function() {
    fa.LoadData();
    app.loading(true);
    setTimeout(function () {
        var param = {
            turbine: fa.turbine,
            project: fa.project,
        };
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getlistpowercurvemonthly", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            if (res.data.Data != null) {
                localStorage.setItem("dataTurbine", JSON.stringify(res.data.Data));
                page.dataPCEachTurbine(res.data.Data);
                page.InitLinePowerCurve();
            }

            app.loading(false);

        })
    }, 300);
}

page.InitLinePowerCurve = function() {
    fa.LoadData();
    listOfChart = [];

    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var minDatetemp = new Date(res.data.ScadaData[0]);
        var maxDatetemp = new Date(res.data.ScadaData[1]);
        $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
    })

    $.each(page.dataPCEachTurbine(), function (i, dataTurbine) {
        var name = dataTurbine.Name
        var idChart = "#chart-" + dataTurbine.Name
        listOfChart.push(idChart);
        
        $(idChart).html("");
        $(idChart).kendoChart({
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            renderAs: "canvas",
            title: {
                text: name,
                font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: false,
            },
            chartArea: {
                width: 300,
                height: 200
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
                    font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555",
                    visible: true,
                },
                labels: {
                    font: '8px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N1",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                    font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555"
                },
                labels: {
                    format: "N0",
                    font: '8px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    rotation: 300,
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
                        font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
            pannable: true,
            zoomable: true
        });
        $(idChart).data("kendoChart").refresh();

        page.powerCurveOptions($(idChart).getKendoChart().options);
    });
}

$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        page.LoadData();
    });
    $(".period-list").hide();
    $(".filter-date-start").hide();
    $(".filter-date-end").hide();
    $(".label-filter")[3].remove();
    $(".label-filter")[2].remove();

    page.LoadData();
});