'use strict';

viewModel.WindAvailabilityAnalysis = new Object();
var avb = viewModel.WindAvailabilityAnalysis;

avb.ChartAvailability = function () {
    fa.getProjectInfo();
    fa.LoadData();
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwindavailability/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var data = res.data;

        $("#windAvailabilityChart").html("");
        $("#windAvailabilityChart").kendoChart({
            dataSource: {
                data: data,
                sort: { field: "WindSpeed", dir: 'asc' }
            },
            theme: "Flat",
            chartArea: {
                height: 500,
            },
            legend: {
                position: "top",
                visible: true,
            },
            series: [{
                type: "column",
                field: "TotalAvail",
                axis: "windPercentage",
                name: "Total Availability [%]",
                opacity: 0.6
            }, {
                type: "line",
                style: "smooth",
                field: "Time",
                axis: "windPercentage",
                name: "Cumulative % of Time",
                markers: {
                    visible: false,
                },
                width: 3,
            }, {
                type: "line",
                style: "smooth",
                field: "Energy",
                axis: "cumProd",
                name: "Cumulative % of Energy Delivered",
                markers: {
                    visible: false,
                },
                width: 3,
            }],
            seriesColors: colorFields2,
            valueAxes: [{
                line: {
                    visible: false
                },
                max: 100,
                majorUnit: 20,
                labels: {
                    format: "{0}%",
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "windPercentage",
                title: { text: "Availability (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }, {
                line: {
                    visible: false
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                max: 100,
                labels: {
                    format: "{0}%",
                },
                name: "cumProd",
                title: { text: "Cumulative Production (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
            }],
            categoryAxis: {
                field: "WindSpeed",
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                axisCrossingValues: [0, 1000],
                justified: true,
                majorGridLines: {
                    visible: false
                },
            },
            tooltip: {
                visible: true,
                shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                // template: "#= series.name # : #= kendo.toString(value, 'n2')# at #= category #",
                template: "#= kendo.toString(value, 'n2')#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
        });

        app.loading(false);
        $("#windAvailabilityChart").data("kendoChart").refresh();
    });
};
vm.currentMenu('Wind Speed vs Availability');
vm.currentTitle('Wind Speed vs Availability');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind', href: '#' }, { title: 'Wind Speed vs Availability', href: viewModel.appName + 'page/analyticloss' }]);

$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        avb.ChartAvailability();
    });

    app.loading(true);
    setTimeout(function () {
        avb.ChartAvailability();
    }, 1000);
});