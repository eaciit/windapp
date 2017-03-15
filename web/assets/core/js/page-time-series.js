'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;
var maxSelectedItems = 4;

vm.currentMenu('Time Series Plots');
vm.currentTitle('Time Series Plots');

pg.tags = ko.observableArray();

if (pageType == "OEM") {
    vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);
    pg.tags = ko.observableArray([
        {text: "Wind Speed" , value:"windspeed"},        
        {text: "Power" , value:"power"},
        {text: "Production" , value:"production"}
    ]);
} else if (pageType == "HFD"){
    vm.breadcrumb([{ title: 'Monitoring', href: '#' }, { title: 'Analysis', href: viewModel.appName + 'page/timeserieshfd' }]);
    pg.tags = ko.observableArray([
        {text: "Wind Speed" , value:"windspeed"},
        {text: "Power" , value:"power"},
        {text: "Wind Direction" , value:"winddirection"},
        {text: "Nacelle Direction" , value:"nacellepos"},
        {text: "Rotor RPM" , value:"rotorrpm"},
        {text: "Pitch Angle" , value:"pitchangle"},        
    ]);
}


pg.availabledatestartscada = ko.observable();
pg.availabledateendscada = ko.observable();
pg.pageType = ko.observable(pageType);
pg.dataType = ko.observable("MIN");

pg.isSecond = ko.observable(false);
pg.TagList = ko.observableArray(["windspeed","power"]);

pg.startTime = ko.observable();
pg.endTime = ko.observable();

var timeSeriesData = [];
var seriesOptions = [],
    seriesCounter = 0;

var breaks = [];    

var yAxis = [];
var chart;
var legend = [];
// var colors = ["#0066dd","#dc3912","#eee"];
var colors = colorField;
var seriesSelectedColor = [];

pg.periodList = ko.observableArray([]);

// pg.periodList = ko.observableArray([{"endtime":"2017-02-17T15:00:00Z","starttime":"2017-02-17T12:00:00Z"},{"endtime":"2017-02-17T18:00:00Z","starttime":"2017-02-17T15:00:00Z"},{"endtime":"2017-02-17T21:00:00Z","starttime":"2017-02-17T18:00:00Z"},{"endtime":"2017-02-18T00:00:00Z","starttime":"2017-02-17T21:00:00Z"},{"endtime":"2017-02-18T03:00:00Z","starttime":"2017-02-18T00:00:00Z"},{"endtime":"2017-02-18T06:00:00Z","starttime":"2017-02-18T03:00:00Z"},{"endtime":"2017-02-18T09:00:00Z","starttime":"2017-02-18T06:00:00Z"},{"endtime":"2017-02-18T12:00:00Z","starttime":"2017-02-18T09:00:00Z"},{"endtime":"2017-02-18T15:00:00Z","starttime":"2017-02-18T12:00:00Z"},{"endtime":"2017-02-18T18:00:00Z","starttime":"2017-02-18T15:00:00Z"},{"endtime":"2017-02-18T21:00:00Z","starttime":"2017-02-18T18:00:00Z"},{"endtime":"2017-02-19T00:00:00Z","starttime":"2017-02-18T21:00:00Z"},{"endtime":"2017-02-19T03:00:00Z","starttime":"2017-02-19T00:00:00Z"},{"endtime":"2017-02-19T06:00:00Z","starttime":"2017-02-19T03:00:00Z"},{"endtime":"2017-02-19T09:00:00Z","starttime":"2017-02-19T06:00:00Z"},{"endtime":"2017-02-19T12:00:00Z","starttime":"2017-02-19T09:00:00Z"},{"endtime":"2017-02-19T15:00:00Z","starttime":"2017-02-19T12:00:00Z"},{"endtime":"2017-02-19T18:00:00Z","starttime":"2017-02-19T15:00:00Z"},{"endtime":"2017-02-19T21:00:00Z","starttime":"2017-02-19T18:00:00Z"},{"endtime":"2017-02-20T00:00:00Z","starttime":"2017-02-19T21:00:00Z"},{"endtime":"2017-02-20T03:00:00Z","starttime":"2017-02-20T00:00:00Z"},{"endtime":"2017-02-20T06:00:00Z","starttime":"2017-02-20T03:00:00Z"},{"endtime":"2017-02-20T09:00:00Z","starttime":"2017-02-20T06:00:00Z"},{"endtime":"2017-02-20T12:00:00Z","starttime":"2017-02-20T09:00:00Z"},{"endtime":"2017-02-20T15:00:00Z","starttime":"2017-02-20T12:00:00Z"},{"endtime":"2017-02-20T18:00:00Z","starttime":"2017-02-20T15:00:00Z"},{"endtime":"2017-02-20T21:00:00Z","starttime":"2017-02-20T18:00:00Z"},{"endtime":"2017-02-21T00:00:00Z","starttime":"2017-02-20T21:00:00Z"}])

pg.LoadData = function(){
	// fa.getProjectInfo();
    fa.LoadData();
    app.loading(true);

    var param = {
        period: fa.period,
        Turbine: fa.turbine,
        DateStart: fa.dateStart,
        DateEnd: fa.dateEnd,
        Project: fa.project,
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

pg.hideLegendByName = function(name){
    $.each(chart.series, function(i, series) {
        if (series.name === name || series.name === (name+"_err")){
            if (series.visible) {
                series.hide();
            } else {
                series.show();
            }
        }
    });
}

pg.createStockChart = function(){
    $("#chartTimeSeries").html("");

    Highcharts.setOptions({
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            marginTop: 0,
            zoomType: 'x'
        }
    });

    chart = Highcharts.stockChart('chartTimeSeries', {
        legend: {
            layout: 'horizontal',
            // padding: 3,
            verticalAlign: 'top',
            borderWidth: 0,
            enabled: true,
            margin : 5,
            enabled: false
        },
        rangeSelector: {
            selected: 1,
            inputEnabled: (pg.isSecond() == true ? false : true)
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
        xAxis: {
            type: 'datetime',
            breaks: breaks,
        },
        yAxis: yAxis,
        plotOptions: {
        series: {
                lineWidth: 2,
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
                // console.log("-----------------------");
                $.each(this.points, function(i, point) {
                    if (typeof legend[i] !== "undefined"){
                        // console.log(point.series.name);
                        if (point.series.name.indexOf("_err") < 0){                            
                            var color = "";

                            $.each(seriesSelectedColor, function(ic, n){
                                if (n==point.series.name) {
                                    color = colors[ic];
                                }
                            });


                            s.push('<span style="color:'+color+';font-weight:bold;cursor:pointer" id="btn-'+i+'" onClick="pg.hideLegendByName(\''+point.series.name+'\')"><i class="fa fa-circle"></i> &nbsp;</span><span style="color:#585555;font-weight:bold;">'+ point.series.name +' : '+kendo.toString(point.y , "n2")+" " +legend[i].unit+'<span>');
                        }
                    }
                });
                $("#legendTooltip").html(s.join("&nbsp;"));
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
pg.hidePopover = function(){
  $('.popover-markup>.trigger').popover('hide');
}
pg.getDataStockChart = function(param, idBtn){
    fa.LoadData();
    app.loading(true);

    if(param == "selectTags"){
       pg.TagList($("#TagList").val());

       $('.popover-markup>.trigger').popover("hide");

    }

    var min = new Date(app.getUTCDate($('input.highcharts-range-selector:eq(0)').val()));
    var max = new Date(app.getUTCDate($('input.highcharts-range-selector:eq(1)').val()));

    var maxDate =  new Date(Date.UTC(max.getFullYear(), max.getMonth(), max.getDate(), 0, 0, 0));
    var minDate =  new Date(Date.UTC(min.getFullYear(), min.getMonth(), min.getDate(), 0, 0, 0));

    var dateStart = fa.dateStart; 
    var dateEnd = fa.dateEnd;

    if(pg.dataType() == 'SEC'){
      dateStart = minDate;
      dateEnd = maxDate;
      if(param == 'detailPeriod'){
          dateStart = new Date(pg.startTime());
          dateEnd = new Date(pg.endTime());
      }
    }

    if(param == "refresh"){
        dateStart = fa.dateStart; 
        dateEnd = fa.dateEnd;
    }

    var IsHour = (param == 'detailPeriod' ? true : false);

    var param = {
        period: fa.period,
        Turbine: [fa.turbine],
        DateStart: dateStart,
        DateEnd: dateEnd,
        Project: fa.project,
        PageType: pg.pageType(),
        DataType: pg.dataType() ,
        TagList : pg.TagList(),
        IsHour : IsHour,
    };

    // var url = (pg.pageType() == "HFD"? "timeseries/getdatahfd" : "timeseries/getdatahfd" )
    var url = "timeseries/getdatahfd";

    var request = toolkit.ajaxPost(viewModel.appName + url, param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data.Data.Chart;
        var periods = res.data.Data.PeriodList;
        // breaks = res.data.Data.Breaks;

        // console.log(breaks);

        if(!IsHour){
            pg.periodList(periods);             
        }

        // console.log(pg.periodList);
        // console.log(data);

        var xCounter = 0;

        $.each(data, function(idx, val){
            var isOpposite = false;
            if (idx >= (maxSelectedItems/2)) {
                isOpposite = true;
            }

            yAxis [idx] = { 
                // startOnTick: false,
                // endOnTick: false,
                // maxPadding: 0.5,
                // minPadding: 0.5,
                min: val.minval,
                max: val.maxval, 
                gridLineWidth: 1,
                labels: {
                    format: '{value}',
                },
                title: {
                    text: val.unit,
                },
                opposite: isOpposite
            }
            seriesOptions[xCounter] = {
                name : val.name, 
                data : val.data,
                color: colors[idx],
                type: 'line',
                yAxis: idx,
                id : "series"+idx,
            }      

            seriesSelectedColor[idx] = val.name;

            legend[idx] = {
                name : val.name,
                unit : val.unit,
                // isCol : false,
            }      

            xCounter+=1;

            seriesOptions[xCounter] = {
                type: 'column',
                name: val.name+"_err",
                data: val.dataerr,
                color: colors[idx],
                pointWidth: 2,
                yAxis: idx,
                id : "series_col"+idx,
                // onSeries: "series"+idx,                
            }

            xCounter+=1;

            seriesCounter += 1;

        //   if (seriesCounter === data.length) {
        //       pg.createStockChart();
        //   }
        });

        pg.createStockChart();
    });

    $.when(request).done(function(){
        setTimeout(function(){
           app.loading(false);
         },200);

        if(pg.dataType() == "SEC"){
          setTimeout(function(){
            pg.prepareScroll();
          },500);

        }
    });
}

pg.prepareScroll = function(){
        var $frame  = $('#basic');
        var $slidee = $frame.children('ul').eq(0);
        var $wrap   = $frame.parent();

        // Call Sly on frame
        $frame.sly({
          horizontal: 1,
          itemNav: 'basic',
          smart: 1,
          activateOn: 'click',
          mouseDragging: 1,
          touchDragging: 1,
          releaseSwing: 1,
          // startAt: 3,
          scrollBar: $wrap.find('.scrollbar'),
          scrollBy: 1,
          pagesBar: $wrap.find('.pages'),
          activatePageOn: 'click',
          speed: 300,
          elasticBounds: 1,
          easing: 'easeOutExpo',
          dragHandle: 1,
          dynamicHandle: 1,
          clickBar: 1,

          // Buttons
          forward: $wrap.find('.forward'),
          backward: $wrap.find('.backward'),
          prev: $wrap.find('.prev'),
          next: $wrap.find('.next'),
          prevPage: $wrap.find('.prevPage'),
          nextPage: $wrap.find('.nextPage')
        });

        // To Start button
        $wrap.find('.toStart').on('click', function () {
          var item = $(this).data('item');
          // Animate a particular item to the start of the frame.
          // If no item is provided, the whole content will be animated.
          $frame.sly('toStart', item);
        });

        // To Center button
        $wrap.find('.toCenter').on('click', function () {
          var item = $(this).data('item');
          // Animate a particular item to the center of the frame.
          // If no item is provided, the whole content will be animated.
          $frame.sly('toCenter', item);
        });

        // To End button
        $wrap.find('.toEnd').on('click', function () {
          var item = $(this).data('item');
          // Animate a particular item to the end of the frame.
          // If no item is provided, the whole content will be animated.
          $frame.sly('toEnd', item);
        });

        // Add item
        $wrap.find('.add').on('click', function () {
          $frame.sly('add', '<li>' + $slidee.children().length + '</li>');
        });

        // Remove item
        $wrap.find('.remove').on('click', function () {
          $frame.sly('remove', -1);
        });
}

$(document).ready(function () {
    $('.popover-markup>.trigger').popover({
        animation: true,
        html: true,
        title: function () {
            return $(this).parent().find('.head').html();
        },
        content: function () {
            return $(this).parent().find('.content').html();
        }
    }).on('click',function () {
            $("#selectTagsDiv").html("");
            $("#selectTagsDiv").html('<select id="TagList"></select>');
            $('#TagList').kendoMultiSelect({
              dataSource: pg.tags(), 
              value: pg.TagList() , 
              dataValueField : 'value', 
              dataTextField: 'text',
              suggest: true, 
              maxSelectedItems: maxSelectedItems, 
              minSelectedItems: 1,
              change: function(e) {
                  if (this.value().length == 0) {
                      this.value("windspeed")
                  }
              }
      });
    });

    $('#btnRefresh').on('click', function () {
        pg.getDataStockChart("refresh");
    });

    setTimeout(function () {
        // pg.LoadData();
        pg.getDataStockChart();
        // pg.prepareScroll();
    }, 1000);
});
