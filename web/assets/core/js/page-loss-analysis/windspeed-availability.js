'use strict';

viewModel.WindspeedAvailability = new Object();
var wa = viewModel.WindspeedAvailability;

wa.WindSpeed = function(){
    fa.LoadData()
    if(pg.isFirstWindSpeed() === true){
        app.loading(true);
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
                    height: 400,
                },
                legend: {
                    position: "top",
                    visible: true,
                    labels: {
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    }
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
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                    labels: {
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    }
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
                dataBound: function(){
                    pg.isFirstWindSpeed(false);
                    app.loading(false);
                }
            });
        });
		$('#availabledatestart').html(pg.availabledatestartscada2());
        $('#availabledateend').html(pg.availabledateendscada2());
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartscada2());
            $('#availabledateend').html(pg.availabledateendscada2());
            $("#windAvailabilityChart").data("kendoChart").refresh();
            // app.loading(false);
        },200);
    } 
}