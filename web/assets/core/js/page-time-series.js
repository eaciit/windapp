'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Time Series Plots');
vm.currentTitle('Time Series Plots');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);


pg.LoadData = function(){
	fa.getProjectInfo();
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
        pg.chartWindSpeed(res.data.Data.windspeed);
        pg.chartProduction(res.data.Data.production);
    });

    $.when(requestData).done(function(){
        setTimeout(function(){
            app.loading(false);
        },500);
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
            template: "#= kendo.toString(value,'n2') #",
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
            template: "#= kendo.toString(value,'n2') #",
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
$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        pg.LoadData();
    });

    setTimeout(function () {
        pg.LoadData();
    }, 1000);
});
