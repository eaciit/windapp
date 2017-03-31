'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;
var maxSelectedItems = 4;
var defaultHour = 5*24;

pg.tags = ko.observableArray();

if (pageType == "OEM") {
    vm.currentMenu('Time Series Plots');
    vm.currentTitle('Time Series Plots');

    vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);
    pg.tags = ko.observableArray([
        {text: "Wind Speed" , value:"windspeed"},        
        {text: "Power" , value:"power"},
        {text: "Production" , value:"production"}
    ]);
} else if (pageType == "HFD"){
    vm.currentMenu('Analysis');
    vm.currentTitle('Analysis');

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

pg.rangeData = ko.observable(true);
pg.errorValue = ko.observable(true);
pg.live = ko.observable(false);

pg.isFirst = ko.observable(true);

var timeSeriesData = [];
var seriesOptions = [],
    seriesCounter = 0;
var seriesOri = [];
var hourBefore = 0;
var date1Before, date2Before;
var breaks = [];    

var yAxis = [];
var newyAxis = [];
var chart;
var legend = [];
var colors = colorField;
var seriesSelectedColor = [];
var interval;

pg.periodList = ko.observableArray([]);



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

pg.hideRange = function(){
    var checked = $('[name=chk-column-range]:checked').length==1;
    $.each(yAxis, function(i, res){
        chart.yAxis[i].update({
            min: (checked ? res.min : null),
            max: (checked == 1 ? res.max : null),
        });

        i++;
    });
}

pg.hideErr = function(){
    var checked = $('[name=chk-column-error]:checked').length==1;
    $.each(chart.series, function(i, res){
        if(res.name.indexOf("_err") > 0){
            res.setVisible((checked ? true : false));
        }
    });
}

pg.getLocalSeries = function(startInt, endInt){
    var seriesOriTmp = JSON.parse(sessionStorage.seriesOri);
    $.each(seriesOriTmp, function(id, val){
        if (val != null){
            var len = val.length;
            var i = 0;
            var startIdx, endIdx = 0;

            while (i < len){
                var curr = val[i];

                if (curr[0]>=startInt && startInt==0) {
                    startIdx = i;
                } else if (curr[0]>=endInt && startIdx != 0) {
                    endIdx = i;
                    break;
                }

                i++;
            }

            chart.series[id].setData(val.slice(startIdx, endIdx), true, true, false);
        }
    });
}


pg.createStockChart = function(y){
    function afterSetExtremes(e) {
        if (pageType != "OEM") {
            var date1 = new Date(new Date(Math.round(e.min)).toUTCString())
            var date2 = new Date(new Date(Math.round(e.max)).toUTCString())
            
            var hours = Math.abs(date1 - date2) / 36e5;

            // console.log("hours: "+hours);

            if (hours <= defaultHour) {
                pg.dataType("SEC");
            }else{
                pg.dataType("MIN");
            }

            if (hourBefore == 0){
                hourBefore = hours;
                date1Before = date1;
                date2Before = date2;
            }

            if ((hours <= defaultHour && hourBefore>defaultHour) || ((date1<date1Before && hours <= defaultHour) || (date2>date2Before && hours <= defaultHour)) ) {
                chart.showLoading('Loading data from server...');
                var param = {
                    period: fa.period,
                    Turbine: [fa.turbine],
                    DateStart: date1,
                    DateEnd: date2,
                    Project: fa.project,
                    PageType: pg.pageType(),
                    DataType: pg.dataType() ,
                    TagList : pg.TagList(),
                    IsHour : true,
                };

                var url = "timeseries/getdatahfd";
                toolkit.ajaxPost(viewModel.appName + url, param, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }

                    var data = res.data.Data.Chart;
                    var periods = res.data.Data.PeriodList;

                    pg.generateSeriesOption(data, periods);
                    $.each(seriesOptions, function(id, val){
                        chart.series[id].setData(val.data, true, true, false);
                    });

                    chart.xAxis[0].update({
                        minRange: 5*1000,
                    })

                    chart.hideLoading();
                });
            }else if (hours > defaultHour) {
                // console.log(e.min+" - "+e.max);
                pg.getLocalSeries(e.min, e.max);
            }

            hourBefore = hours;
            date1Before = date1;
            date2Before = date2;
        }
    }


    $("#chartTimeSeries").html("");

    var minRange = 600 * 1000;
    if(pg.dataType() == 'SEC'){
        minRange = 5 * 1000;
    }

    Highcharts.setOptions({
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            zoomType: 'x',
        }
    });

    chart = Highcharts.stockChart('chartTimeSeries', {
        legend: {
            symbolHeight: 12,
            symbolWidth: 12,
            symbolRadius: 6,
            enabled: true,
            floating: false,
            align: 'center',
            verticalAlign: 'top',
            labelFormat: '<span style="color:{color}">{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
            borderWidth: 0,
            marginTop: -70,
        },
        rangeSelector: {
            buttons: [{
                type: 'hour',
                count: 1,
                text: '1h'
            }, {
                type: 'day',
                count: 1,
                text: '1d'
            }, {
                type: 'month',
                count: 1,
                text: '1m'
            }, {
                type: 'year',
                count: 1,
                text: '1y'
            }, {
                type: 'all',
                text: 'All'
            }],
            inputEnabled: true,
            selected: 2 ,// all,
            y: 50
        },
        navigator: {
            adaptToUpdatedData: false,
            series: {
                color: '#999',
                lineWidth: 1
            }
        },
        exporting: {
          enabled: false
        },
        xAxis: {
            events: {
                afterSetExtremes: afterSetExtremes
            },
            type: 'datetime',
            breaks: breaks,
            minRange: minRange,
        },
        yAxis: (y == undefined ? yAxis : y),
        plotOptions: {
            series: {
                lineWidth: 1,
                states: {
                    hover: {
                        enabled: true,
                        lineWidth: 2
                    }
                },
                events: {
                    legendItemClick: function () {
                        pg.hideLegendByName(this.name);
                        return false;
                    }
                }
            },
        },
        series: seriesOptions,
        tooltip:{
             formatter : function() {
                $("#dateInfo").html( Highcharts.dateFormat('%e %b %Y %H:%M:%S', this.x));
                return false ;
             }
        },
    });

    // seriesOri = chart.series;
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

pg.getDataStockChart = function(param){
    fa.LoadData();
    app.loading(true);
    clearInterval(interval);
    if(param == "selectTags"){
       pg.TagList($("#TagList").val());
       $('.popover-markup>.trigger').popover("hide");
    }

    var IsHour = (pg.isFirst() == true ? false : true);

    var COOKIES = {};
    var cookieStr = document.cookie;
    var turbine = "";
    
    console.log(cookieStr);
    if(cookieStr.indexOf("turbine=") >= 0) {
        document.cookie = "turbine=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
        cookieStr.split(/; /).forEach(function(keyValuePair) {
            var cookieName = keyValuePair.replace(/=.*$/, "");
            var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
            COOKIES[cookieName] = cookieValue;
        });
        turbine = COOKIES["turbine"];
        $('#turbineList').data('kendoDropDownList').value(turbine);
    } else {
        turbine = $('#turbineList').data('kendoDropDownList').value();
    }

    var min = new Date(app.getUTCDate($('input.highcharts-range-selector:eq(0)').val()));
    var max = new Date(app.getUTCDate($('input.highcharts-range-selector:eq(1)').val()));

    var maxDate =  new Date(Date.UTC(max.getFullYear(), max.getMonth(), max.getDate(), 0, 0, 0));
    var minDate =  new Date(Date.UTC(min.getFullYear(), min.getMonth(), min.getDate(), 0, 0, 0));

    var now = new Date()

    if(pg.isFirst() == true){
      fa.period = "custom";
    }

    if(pg.pageType() == 'HFD'){
        fa.dateEnd = new Date();
        fa.dateStart  = new Date(now.setMonth(now.getMonth() - 24));
        
        date1Before = fa.dateStart;
        date2Before = fa.dateEnd;
        hourBefore = Math.abs(date1Before - date2Before) / 36e5;
    }
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

    // var IsHour = (param == 'detailPeriod' ? true : false);
    var paramX = {
        period: fa.period,
        Turbine: [turbine],
        DateStart: dateStart,
        DateEnd: dateEnd,
        Project: fa.project,
        PageType: pg.pageType(),
        DataType: pg.dataType() ,
        TagList : pg.TagList(),
        IsHour : IsHour,
    };

    var url = "timeseries/getdatahfd";
    if($('input[name="chk-column-live"]:checked').length > 0){
        pg.live(true);
        pg.rangeData(true);
        pg.errorValue(true);
    }else{
        pg.live(false);
        pg.rangeData(true);
        pg.errorValue(true);
    }


    var request;
    if(pg.live() == false){
        request = toolkit.ajaxPost(viewModel.appName + url, paramX, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            var data = res.data.Data.Chart;
            var periods = res.data.Data.PeriodList;
            breaks = res.data.Data.Breaks;

            pg.generateSeriesOption(data, periods);
            
            if (param=="first" || param=="refresh"){
                // if (sessionStorage.seriesOri){
                    sessionStorage.seriesOri = [];
                // }
                
                $.each(seriesOptions,function(idx, val){
                    if (val.data != null){
                        var valx = val.data.slice(idx);
                        seriesOri[idx] = valx;
                    }
                });

                // if (seriesOri != null) {
                    sessionStorage.seriesOri = JSON.stringify(seriesOri);
                    seriesOri = [];    
                // }
            }

            pg.createStockChart();
        });
    }else{
        pg.createLiveChart(IsHour);
    }


    $.when(request).done(function(){
        pg.isFirst(false);
        setTimeout(function(){
           app.loading(false);
         },200);

        // if(pg.dataType() == "SEC"){
        //   setTimeout(function(){
        //     pg.prepareScroll();
        //   },500);

        // }
    });
}

pg.createLiveChart = function(IsHour){
    var param = {
        period: fa.period,
        Turbine: [fa.turbine],
        Project: fa.project,
        PageType: "LIVE",
        DataType: pg.dataType() ,
        TagList : pg.TagList(),
        IsHour : IsHour,
    };

    var dateStart, dateEnd; 
    $("#chartTimeSeries").html("");
        toolkit.ajaxPost(viewModel.appName + "timeseries/getdatahfd", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }


        var data = res.data.Data.Chart;
        var periods = res.data.Data.PeriodList;

        dateStart = new Date(new Date(Math.round(data[0].data[0][0])).toUTCString());
        dateEnd = new Date(new Date(Math.round(data[0].data[0][0])).toUTCString());

        console.log(dateEnd);

        breaks = res.data.Data.Breaks;

        pg.generateSeriesOption(data, periods);
        
        if (param=="first" || param=="refresh"){
            if (sessionStorage.seriesOri){
                sessionStorage.seriesOri = null;
            }
            
            $.each(seriesOptions,function(idx, val){
                if (val.data != null){
                    var valx = val.data.slice(idx);
                    seriesOri[idx] = valx;
                }
            });

            sessionStorage.seriesOri = JSON.stringify(seriesOri);
            seriesOri = [];
        }

       
        $("#chartTimeSeries").html("");

        var minRange = 600 * 1000;
        if(pg.dataType() == 'SEC'){
            minRange = 5 * 1000;
        }


        chart = Highcharts.stockChart('chartTimeSeries', {
             chart: {
                tyle: {
                    fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                zoomType: 'x',
                type: 'spline',
                marginRight: 10,
                events: {
                    load: function () {
                        interval = setInterval(function () {
                            var paramX = {
                                period: fa.period,
                                Turbine: [fa.turbine],
                                DateStart: dateStart,
                                DateEnd: dateEnd,
                                Project: fa.project,
                                PageType: "LIVE",
                                DataType: pg.dataType() ,
                                TagList : pg.TagList(),
                                IsHour : IsHour,
                            };

                            toolkit.ajaxPost(viewModel.appName + "timeseries/getdatahfd", paramX, function (res) {
                                if (!app.isFine(res)) {
                                    return;
                                }

                                var seriesData = chart.series; 
                                var results = res.data.Data.Chart;
                                if(results.length > 0){
                                    dateStart = new Date(new Date(Math.round(results[0].data[0][0])).toUTCString());
                                    dateEnd = new Date(new Date(Math.round(results[0].data[0][0])).toUTCString());
                                    $.each(seriesOptions, function(i, res){
                                       $.each(results, function(id, val){
                                            if(res.name == val.name){
                                                var x = (new Date()).getTime();
                                                chart.series[i].addPoint(val.data, true, true);
                                            }
                                       })
                                    });
                                    console.log(dateEnd);
                                }else{
                                    console.log(dateEnd);
                                    return false;
                                }

                            });
                        }, 5000);
                    }
                }
            },
             legend: {
                symbolHeight: 12,
                symbolWidth: 12,
                symbolRadius: 6,
                enabled: true,
                floating: false,
                align: 'center',
                verticalAlign: 'top',
                labelFormat: '<span style="color:{color}">{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
                borderWidth: 0,
                marginTop: -70,
            },
             rangeSelector: {
                buttons: [{
                    type: 'hour',
                    count: 1,
                    text: '1h'
                }, {
                    type: 'day',
                    count: 1,
                    text: '1d'
                }, {
                    type: 'month',
                    count: 1,
                    text: '1m'
                }, {
                    type: 'year',
                    count: 1,
                    text: '1y'
                }, {
                    type: 'all',
                    text: 'All'
                }],
                inputEnabled: true,
                selected: 2 ,// all,
                y: 50
            },
            navigator: {
                adaptToUpdatedData: false,
                series: {
                    color: '#999',
                    lineWidth: 1
                }
            },
            exporting: {
              enabled: false
            },
            xAxis: {
                type: 'datetime',
                breaks: breaks,
                dateTimeLabelFormats : {
                    millisecond: '%H:%M:%S.%L',
                    second: '%H:%M:%S',
                    minute: '%H:%M',
                    hour: '%H:%M',
                    day: '%e. %b',
                    week: '%e. %b',
                    month: '%b \'%y',
                    year: '%Y'
                },
                // minRange: 5*1000,
            },
            yAxis: yAxis,
            plotOptions: {
            series: {
                    lineWidth: 1,
                    states: {
                        hover: {
                            enabled: true,
                            lineWidth: 2
                        }
                    },
                    events: {
                        legendItemClick: function () {
                            pg.hideLegendByName(this.name);
                            return false;
                        }
                    }
                }
            },
            series: seriesOptions,
            tooltip:{
                 formatter : function() {
                    $("#dateInfo").html( Highcharts.dateFormat('%e %b %Y %H:%M:%S', this.x));
                    return false ;
                 }
            },
        });
    });
    
}
pg.generateSeriesOption = function(data, periods){
    var IsHour = (pg.isFirst() == true ? false : true);
    var IsGroup = (pg.dataType() == "SEC" ? false : true);
    console.log("isgroup: "+IsGroup);

    if(!IsHour){
        pg.periodList(periods);             
    }

    yAxis = [];
    seriesOptions = [];

    var xCounter = 0;

    $.each(data, function(idx, val){
        var isOpposite = false;
        if (idx >= (maxSelectedItems/2) || (idx == 1 && data.length==2)) {
            isOpposite = true;
        }

        yAxis[xCounter] = {
            min: val.minval,
            max: val.maxval, 
            gridLineWidth: 1,
            endOnTick: false,
            startOnTick: false,
            showLastLabel: true,
            maxPadding: 0,
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
            yAxis: xCounter,
            id : "series"+idx,
            showInNavigator: true,
            tooltip: {
                valueSuffix: val.unit
            },
            dataGrouping:{
                enabled: IsGroup,
                // units: [["day",[1]],["weel",[1]],["month",[1]]],
            }
        }      

        seriesSelectedColor[idx] = val.name;

        legend[idx] = {
            name : val.name,
            unit : val.unit,
        }      

        xCounter+=1;

        yAxis[xCounter] = {
            min: 0,
            max: 100, 
            gridLineWidth: 1,
            labels: {
                format: '{value}',
            },
            title: {
                text: val.unit,
            },
            opposite: isOpposite,
            visible: false,
        }

        seriesOptions[xCounter] = {
            type: 'column',
            name: val.name+"_err",
            data: val.dataerr,
            color: colors[idx],
            pointWidth: 1,
            yAxis: xCounter,
            id : "series_col"+idx,
            showInLegend : false,
            // dataGrouping: {
            //     approximation: function () {
            //         return 100;
            //     },
            //     forced: true
            // },
            showInNavigator: true,
            onSeries: "series"+idx,                
        }

        xCounter+=1;

        seriesCounter += 1;
    });
}

/*pg.prepareScroll = function(){
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
}*/

$(document).ready(function () {
    newyAxis = yAxis;
    if(pg.pageType() === "HFD"){
        $("#periodList").closest(".k-widget").hide();
        $("#dateStart").closest(".k-widget").hide();
        $("#dateEnd").closest(".k-widget").hide();
        $(".label-filters:contains('Period')").hide();
        $(".label-filters:contains('to')").hide();
    }

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
        pg.rangeData(true);
        pg.errorValue(true);
    });

    setTimeout(function () {
        // pg.LoadData();
        pg.getDataStockChart("first");
        // pg.prepareScroll();
    }, 1000);
});
