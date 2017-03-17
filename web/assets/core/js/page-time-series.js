'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;
var maxSelectedItems = 4;

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

var breaks = [];    

var yAxis = [];
var newyAxis = [];
var chart;
var legend = [];
var colors = colorField;
var seriesSelectedColor = [];

pg.periodList = ko.observableArray([]);


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
    var yA = [];
    $.each(yAxis, function(i, res){
        var y = {
            min: ($('[name=chk-column-range]:checked').length == 1 ? res.min : null),
            max: ($('[name=chk-column-range]:checked').length == 1 ? res.max : null),
            gridLineWidth: 1,
            labels: {
                format: '{value}',
            },
            title: {
                text: res.title.text,
            },
            opposite: res.opposite,
            visible: res.visible,
        }

        yA.push(y);
    });

    newyAxis = yA;

    pg.createStockChart(newyAxis);
}

pg.hideErr = function(){
    $.each(seriesOptions, function(i, res){
          if(res.name.indexOf("_err") > 0){
              res.visible = ($('[name=chk-column-error]:checked').length == 1 ? true : false);
          }
    });
    if($('[name=chk-column-range]:checked').length == 1){
        pg.createStockChart();
    }else{
        pg.createStockChart(newyAxis);
    }
    
}


pg.createStockChart = function(y){
    function afterSetExtremes(e) {
        var date1 = new Date(new Date(Math.round(e.min)).toUTCString())
        var date2 = new Date(new Date(Math.round(e.max)).toUTCString())

        var hours = Math.abs(date1 - date2) / 36e5;
        if (hours <= 24) {
            pg.dataType("SEC");
        }else{
            pg.dataType("MIN");
        }

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
        var request = toolkit.ajaxPost(viewModel.appName + url, param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            var data = res.data.Data.Chart;
            var periods = res.data.Data.PeriodList;

            pg.generateSeriesOption(data, periods);
            // chart.addSeries(seriesOptions);
            $.each(seriesOptions, function(id, val){
                chart.series[id].setData(val.data);
            });

            // chart.series = seriesOptions;
            // chart.yAxis = yAxis;

            chart.hideLoading();
        });
    }


    $("#chartTimeSeries").html("");

    // var minRange = 600 * 1000;
    // if(pg.dataType() == 'SEC'){
        var minRange = 5 * 1000;
    // }

    Highcharts.setOptions({
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            zoomType: 'x'
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
            inputEnabled: (pg.isSecond() == true ? false : true),
            selected: 4 ,// all,
            y: 50
        },
        navigator: {
            adaptToUpdatedData: false,
            series: {
                color: '#999',
                lineWidth: 2
            }
        },
        exporting: {
          enabled: false
        },
        xAxis: {
            // events: {
            //     afterSetExtremes: afterSetExtremes
            // },
            type: 'datetime',
            breaks: breaks,
            minRange: minRange,
        },
        yAxis: (y == undefined ? yAxis : y),
        plotOptions: {
        series: {
                lineWidth: 2,
                states: {
                    hover: {
                        enabled: true,
                        lineWidth: 1
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


    if(pg.isFirst() == true){
      fa.period = "custom";
      var now = new Date()
      fa.dateEnd = new Date();
      fa.dateStart  = new Date(now.setMonth(now.getMonth() - 6));
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

    if(param == "refresh"){
        dateStart = fa.dateStart; 
        dateEnd = fa.dateEnd;
    }

    // var IsHour = (param == 'detailPeriod' ? true : false);
    var IsHour = (pg.isFirst() == true ? false : true);

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

    var url = "timeseries/getdatahfd";

    var request = toolkit.ajaxPost(viewModel.appName + url, param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data.Data.Chart;
        var periods = res.data.Data.PeriodList;
        // breaks = res.data.Data.Breaks;

        pg.generateSeriesOption(data, periods);        
        pg.createStockChart();
    });

    $.when(request).done(function(){
        pg.isFirst(false);
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

pg.generateSeriesOption = function(data, periods){
    var IsHour = (pg.isFirst() == true ? false : true);

    if(!IsHour){
        pg.periodList(periods);             
    }

    yAxis = [];
    seriesOptions = [];

    var xCounter = 0;

    $.each(data, function(idx, val){
        var isOpposite = false;
        if (idx >= (maxSelectedItems/2)) {
            isOpposite = true;
        }

        yAxis[xCounter] = {
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
            yAxis: xCounter,
            id : "series"+idx,
            showInNavigator: true,
            tooltip: {
                valueSuffix: val.unit
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
            pointWidth: 2,
            yAxis: xCounter,
            id : "series_col"+idx,
            showInLegend : false,
            // showInNavigator: true,
            // onSeries: "series"+idx,                
        }

        xCounter+=1;

        seriesCounter += 1;
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
    });

    setTimeout(function () {
        // pg.LoadData();
        pg.getDataStockChart();
        // pg.prepareScroll();
    }, 1000);
});
