'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Time Series Plots');
vm.currentTitle('Time Series Plots');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);

pg.availabledatestartscada = ko.observable();
pg.availabledateendscada = ko.observable();
var timeSeriesData = [];
var seriesOptions = [],
    seriesCounter = 0;

var yAxis = [];
var chart;
var legend = [];
var colors = ["#0066dd","#dc3912","#eee"];

pg.LoadData = function(){
	// fa.getProjectInfo();
    fa.LoadData();
    app.loading(true);

    var param = {
        period: fa.period,
        Turbine: fa.turbine,
        DateStart: fa.dateStart,
        DateEnd: fa.dateEnd,
        Project: fa.project
    };

    var requestData = toolkit.ajaxPost(viewModel.appName + "timeseries/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        // pg.chartWindSpeed(res.data.Data.windspeed);
        // pg.chartProduction(res.data.Data.production);
        timeSeriesData = res.data.Data;
        pg.createChart();
    });

    $.when(requestData).done(function(){
        setTimeout(function(){
            app.loading(false);
        },500);
    });
}

toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
    if (!app.isFine(res)) {
        return;
    }
    var minDatetemp = new Date(res.data.ScadaData[0]);
    var maxDatetemp = new Date(res.data.ScadaData[1]);

    pg.availabledatestartscada(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
    pg.availabledateendscada(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

    $('#availabledatestart').html(pg.availabledatestartscada());
    $('#availabledateend').html(pg.availabledateendscada());

})

pg.setSeries = function(name, axis, color, data){
  return {
    name: name,
    type: "line",
    field: "value",
    categoryField: "timestamp",
    axis: axis,
    color: color,
    data: data,
    aggregate: "sum",
    markers : {
        visible : false,
    },
  }
}

pg.setValueAxis = function(name, titleText, crossingVal){
  return {
    name: name,
    title: {
        text: titleText,
        visible: true,
        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
    },
    labels: {
        // format: "{0:p2}"
        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
        format: "{0}"
    },
    majorGridLines: {
        visible: true,
        color: "#eee",
        width: 0.8,
    },
    line: {
        visible: false
    },
    axisCrossingValue: crossingVal
  }
}

pg.createChart = function(){
  var seriesList = [];
  seriesList.push(pg.setSeries("Wind Speed", "windspeedAxis", "#337ab7", timeSeriesData.windspeed));
  seriesList.push(pg.setSeries("Production", "productionAxis", "#ea5b19", timeSeriesData.production));
  
  var valueAxisList = [];
  valueAxisList.push(pg.setValueAxis("windspeedAxis", "m/s", 0));
  valueAxisList.push(pg.setValueAxis("productionAxis", "MWh", 0));

  var naviSeriesList = [];
  var series1 = pg.setSeries("Wind Speed", "windspeedAxis", "#337ab7", timeSeriesData.windspeed);
  series1.shared = true;
  var series2 = pg.setSeries("Production", "productionAxis", "#ea5b19", timeSeriesData.production);
  series2.shared = true;
  naviSeriesList.push(series1);
  naviSeriesList.push(series2);

  $("#chartTimeSeries").kendoStockChart({
    title: {
        text: "Wind Speed & Production",
        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: true
      },
      theme: "flat",
      seriesDefaults: {
          area: {
              line: {
                  style: "smooth",
              }
          },
          width: 1.3,
      },
      series: seriesList,
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        series: naviSeriesList
      },
      valueAxis: valueAxisList,
      categoryAxis: {
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: false
            },
            majorTickType: "none",
            axisCrossingValues: [0, 1000],
            autoBaseUnitSteps: {
                // Would produce 31 groups
                // => Skip to weeks
                minutes: [10],
            }
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # MWh",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
            shared: true,
            sharedTemplate: kendo.template($("#template").html())
        }
    });
} 

pg.chartWindSpeed = function(dataSource){
	$("#chartWindSpeed").kendoStockChart({
	  title: {
        text: "Wind Speed",
        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      theme: "flat",
      seriesDefaults: {
	        area: {
	            line: {
	                style: "smooth"
	            }
	        }
	    },
      dateField: "timestamp",
      series: [{
        type: "area",
        field: "value",
        aggregate: "avg", 
        color: "#337ab7",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        series: [{
          type: "area",
          field: "value",
          aggregate: "avg",
          color: "#337ab7",
        }]
      },
      valueAxis: {
        title: {
            text: "m/s",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
	  categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # m/s",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
} 

pg.chartProduction = function(dataSource){
	$("#chartProduction").kendoStockChart({
	  title: {
        text: "Production",
        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: dataSource
      },
      theme: "flat",
      seriesDefaults: {
	        area: {
	            line: {
	                style: "smooth"
	            }
	        }
	    },
      dateField: "timestamp",
      series: [{
        type: "area",
        field: "value",
        aggregate: "sum", 
        color: "#ea5b19",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        series: [{
          type: "area",
          field: "value",
          aggregate: "sum",
          color: "#ea5b19",
        }]
      },
      valueAxis: {
        title: {
            text: "MWh",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
	  categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
      tooltip: {
            visible: true,
            template: "#= kendo.toString(value,'n2') # MWh",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
} 


pg.hideLegend = function(idx){
  var series = chart.series[idx];
  if (series.visible) {
      series.hide();
      // $button.html('Show series');
  } else {
      series.show();
      // $button.html('Hide series');
  }
}
pg.createStockChart = function(){
    $("#chartTimeSeries").html("");

    Highcharts.setOptions({
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            marginTop: 0
        }
    });

    chart = Highcharts.stockChart('chartTimeSeries', {
        legend: {
                layout: 'horizontal',
                padding: 3,
                verticalAlign: 'top',
                borderWidth: 0,
                enabled: true,
                margin : 5,
                enabled: false
        },
        rangeSelector: {
            selected: 1
        },
         navigator: {
            series: {
                color: '#999',
                lineWidth: 2
            }
        },
        exporting: {
          enabled: false
        },
        yAxis: yAxis,
        plotOptions: {
        series: {
                lineWidth: 1,
                states: {
                    hover: {
                        enabled: true,
                        lineWidth: 1
                    }
                }
            }
        },
        tooltip:{         
          formatter : function() {
            var s = [];
              $.each(this.points, function(i, point) {
                  s.push('<span style="color:'+colors[i]+';font-weight:bold;cursor:pointer" id="btn-'+i+'" onClick="pg.hideLegend('+i+')"><i class="fa fa-circle"></i> &nbsp;</span><span style="color:#585555;font-weight:bold;">'+ point.series.name +' : '+
                      kendo.toString(point.y , "n2")+" " +legend[i].unit+'<span>');
              });
              
               $("#testTooltip").html(s.join("&nbsp;"));
               return false;
           }
        },
        series: seriesOptions,

    });
}


pg.getTimestamp = function(param){
  var dateString = moment(param).format("DD-MM-YYYY HH:mm:ss"),
      dateTimeParts = dateString.split(' '),
      timeParts = dateTimeParts[1].split(':'),
      dateParts = dateTimeParts[0].split('-'),
      date;

      date = new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0], timeParts[0], timeParts[1]);

      return date.getTime();
}
pg.getDataStockChart = function(){
    fa.LoadData();
    app.loading(true);

    var param = {
        period: fa.period,
        Turbine: fa.turbine,
        DateStart: fa.dateStart,
        DateEnd: fa.dateEnd,
        Project: fa.project
    };

    var request = toolkit.ajaxPost(viewModel.appName + "timeseries/getdatahfd", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data.Data;

        $.each(data, function(idx, val){
             yAxis [idx] = { 
                startOnTick: false,
                gridLineWidth: 1,
                labels: {
                    format: '{value}',
                },
                title: {
                    text: val.unit,
                },
                opposite: (val.name == "Production" ? true : false)

            }
            seriesOptions[idx] = {
                  name : val.name, 
                  data : val.data,
                  color: colors[idx],
                  type: 'line',
                  yAxis: idx,
                  tooltip: {
                      valueSuffix: val.unit,
                  }
            }

            

            legend[idx] = {
                name : val.name,
                unit : val.unit
            }

          seriesCounter += 1;

          if (seriesCounter === data.length) {
              pg.createStockChart();
          }
        });
        console.log(legend);
    });

    $.when(request).done(function(){
        setTimeout(function(){
           app.loading(false);
         },200);
    });
}

$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        pg.getDataStockChart();
    });

    setTimeout(function () {
        // pg.LoadData();
        pg.getDataStockChart();
    }, 1000);
});
