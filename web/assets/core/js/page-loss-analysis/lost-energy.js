'use strict';

viewModel.LostEnergy = new Object();
var le = viewModel.LostEnergy;

le.LossEnergy = function(){
    var valid = fa.LoadData();
    if (valid) {
        pg.setAvailableDate(false);
        if(pg.isFirstLostEnergy() === true){
            app.loading(true);

            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));

            var paramdown = {
                Period: fa.period,
                DateStart: dateStart,
                DateEnd: dateEnd,
                Turbine: fa.turbine(),
                Project: fa.project
            };
            toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeloss", paramdown, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                setTimeout(function(){
                    le.DTLEbyType(res.data);
                },200)
            });
            
            var param = {
                period: fa.period,
                dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
                dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
                turbine: fa.turbine(),
                project: fa.project,
            }

            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getlostenergytab", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                setTimeout(function(){
                    le.TLossCat('chartLCByTEL', true, res.data.catloss, 'MWh');
                    le.TLossCat('chartLCByDuration', false, res.data.catlossduration, 'Hours');
                },300);
            });
            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getfrequencytab", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                setTimeout(function(){
                    le.TLossCat('chartLCByFreq', false, res.data.catlossfreq, 'Times');

                    app.loading(false);
                    pg.isFirstLostEnergy(false);
                },300);
            });
        }else{
            setTimeout(function(){
                $("#chartLCByTEL").data("kendoChart").refresh();
                $("#chartDTLEbyType").data("kendoChart").refresh();
                $("#chartLCByDuration").data("kendoChart").refresh();
                $("#chartLCByFreq").data("kendoChart").refresh();
            },200)
        }
    }
}

le.TLossCat = function (id, byTotalLostenergy, dataSource, measurement) {
    var gapVal = 1
    var templateLossCat = ''
    switch (dataSource.length) {
        case 1:
            gapVal = 5;
            break;
        case 2:
            gapVal = 3;
            break;
        case 3:
            gapVal = 1;
            break;
        case 4:
            gapVal = 1;
            break;
        case 5:
            gapVal = 1;
            break;
    } 

    if(measurement == "MWh") {
       templateLossCat = "<b>#: category # :</b> #: kendo.toString(value/1000, 'n1')# " + measurement
    } else if(measurement == "Hours") {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n1')# " + measurement
    } else {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n0')# "
    }

    $('#' + id).html("");
    $('#' + id).kendoChart({
        dataSource: {
            data: dataSource,
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
            height: 163,
        },
        seriesDefaults: {
            type: "column",
            gap: gapVal,
        },
        series: [{
            type: "column",
            field: "result",
        }],
        valueAxis: {
            labels: {
                step: 2,
                template: (byTotalLostenergy == true) ? "#= value / 1000 #" : "#= value#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            title: {
                text: measurement,
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

        },
        categoryAxis: {
            field: "title",
            title: {
                text: "Loss Categories",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: false
            },
            majorTickType: "none",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            template: templateLossCat,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });
}

le.DTLEbyType = function (dataSource) {
    $("#chartDTLEbyType").kendoChart({
        dataSource: {
            data: dataSource,
            group: [{ field: "_id.id3" }],
            sort: { field: "_id.id1", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209) - 120) / 2,
            height: 163,
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "powerlost",
            // opacity : 0.7,
            stacked: true,
            axis: "PowerLost"
        },
        {
            name: function () {
                return "Duration";
            },
            type: "line",
            field: "duration",
            axis: "Duration",
            markers: {
                visible: false
            }
        },
        {
            name: function () {
                return "Frequency";
            },
            type: "line",
            field: "frequency",
            axis: "Frequency",
            markers: {
                visible: false
            }
        }],
        seriesColors: colorField,
        valueAxis: [{
            name: "PowerLost",
            labels: {
                step: 2,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
        },
        {
            name: "Duration",
            title: { visible: false },
            visible: false,
        },
        {
            name: "Frequency",
            title: { visible: false },
            visible: false,
        }],
        categoryAxis: {
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            sharedTemplate: kendo.template($("#templateDTLE").html()),
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });
}