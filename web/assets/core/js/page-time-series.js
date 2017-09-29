'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;
var maxSelectedItems = 5;
var defaultHour = 5*24;

pg.tags = ko.observableArray();

if (pageType == "OEM") {
    vm.currentMenu('Time Series Plots');
    vm.currentTitle('Time Series Plots');

    vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);
    pg.tags = ko.observableArray([
        {text: "Wind Speed" , value:"windspeed"},        
        {text: "Power" , value:"power"},
        {text: "Production" , value:"production"}, 
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
        {text: "Pitch Cabinet Temp Blade 1" , value:"PitchCabinetTempBlade1"},
        {text: "Pitch Cabinet Temp Blade 2" , value:"PitchCabinetTempBlade2"},
        {text: "Pitch Cabinet Temp Blade 3" , value:"PitchCabinetTempBlade3"},
        {text: "Pitch Conv Internal Temp Blade 1" , value:"PitchConvInternalTempBlade1"},
        {text: "Pitch Conv Internal Temp Blade 2" , value:"PitchConvInternalTempBlade2"},
        {text: "Pitch Conv Internal Temp Blade 3" , value:"PitchConvInternalTempBlade3"},
        {text: "Temp G1L1" , value:"TempG1L1"},
        {text: "Temp G1L2" , value:"TempG1L2"},
        {text: "Temp G1L3" , value:"TempG1L3"},
        {text: "Temp Generator Bearing DE" , value:"TempGeneratorBearingDE"},
        {text: "Temp Generator Bearing NDE" , value:"TempGeneratorBearingNDE"},
        {text: "Temp Gear Box Oil Sump" , value:"TempGearBoxOilSump"},
        {text: "Temp Hub Bearing" , value:"TempHubBearing"},
        {text: "Temp Generator Choke" , value:"TempGeneratorChoke"},
        {text: "Temp Grid Choke" , value:"TempGridChoke"},
        {text: "Temp Generator Cooling Unit" , value:"TempGeneratorCoolingUnit"},
        {text: "Temp Conv Cabinet 2" , value:"TempConvCabinet2"},
        {text: "Temp Outdoor" , value:"TempOutdoor"},
        {text: "Temp Slip Ring" , value:"TempSlipRing"},
        {text: "Transformer Winding Temp 1" , value:"TransformerWindingTemp1"},
        {text: "Transformer Winding Temp 2" , value:"TransformerWindingTemp2"},
        {text: "Transformer Winding Temp 3" , value:"TransformerWindingTemp3"},
        {text: "Temp Shaft Bearing 1" , value:"TempShaftBearing1"},
        {text: "Temp Shaft Bearing 2" , value:"TempShaftBearing2"},
        {text: "Temp Shaft Bearing 3" , value:"TempShaftBearing3"},
        {text: "Temp Gear Box IMSDE" , value:"TempGearBoxIMSDE"},
        {text: "Temp Bottom Control Section" , value:"TempBottomControlSection"},
        {text: "Temp Bottom Power Section" , value:"TempBottomPowerSection"},
        {text: "Temp Cabinet Top Box" , value:"TempCabinetTopBox"},
        {text: "Temp Nacelle" , value:"TempNacelle"},
        {text: "Voltage L1" , value:"VoltageL1"},
        {text: "Voltage L2" , value:"VoltageL2"},
        {text: "Voltage L3" , value:"VoltageL3"},
        {text: "Blade Voltage V1" , value:"PitchAccuV1"},
        {text: "Blade Voltage V2" , value:"PitchAccuV2"},
        {text: "Blade Voltage V3" , value:"PitchAccuV3"},      
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
var seriesOriX = [];
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
var minXAxis, maxXAxis;
var isSelected = false;
var dateTimeLabelFormats = {
    millisecond: '%H:%M:%S',
    second: '%H:%M:%S',
    minute: '%H:%M',
    hour: '%H:%M',
    day: '%e. %b',
    week: '%e. %b',
    month: '%b \'%y',
    year: '%Y'
};
pg.periodList = ko.observableArray([]);



pg.hideLegend = function(idx){
  var series = chart.series[idx];
  if (series.visible) {
      series.hide();
  } else {
      series.show();
  }
}

pg.hideLegendByName = function(name){
    $.each(chart.series, function(i, series) {
        if (series.name === name || series.name === "_err"){
            if (series.color == "rgba(0, 0, 255, 0)") {
                series.options.color =  colors[i];
                series.update(series.options);
            } else {
                series.options.color = 'rgba(0, 0, 255, 0)';
                series.update(series.options);
            }
        }
    });
}

pg.hideRange = function(){
    var checked = $('#option1:checked').length==1;
    $.each(yAxis, function(i, res){
        if(chart.series[i].name != "_err"){
            chart.yAxis[i].update({
                min: (checked ? res.min : null),
                max: (checked ? res.max : null),
                tickInterval: (checked ? res.max/5 : null),
                alignTicks: (checked ? false : true),
            });
        }
    });
}

pg.hideErr = function(){
    var checked = $('#option2:checked').length==1;
    $.each(chart.series, function(i, res){
        if(res.name == "_err"){
            res.setVisible(checked);
        }
    });
}

pg.getLocalSeries = function(startInt, endInt){
    // var seriesOriTmp = JSON.parse(seriesOriX);
    $.each(JSON.parse(seriesOriX), function(id, val){
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
                $('.highcharts-range-selector-buttons').hide();
                var param = {
                    period: fa.period,
                    Turbine: fa.turbine,
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
                    var outliers = res.data.Data.Outliers;

                    pg.generateSeriesOption(data, periods);
                    pg.generateOutliers(outliers);

                    $.each(seriesOptions, function(id, val){
                        chart.series[id].setData(val.data, true, true, false);
                    });

                    chart.xAxis[0].update({
                        minRange: 5*1000,
                    })

                    chart.hideLoading();
                    $('.highcharts-range-selector-buttons').show();
                });
            }else if (hours > defaultHour) {
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

    var rangeSelected = 3;
    if(pg.pageType() == 'OEM'){
        rangeSelected = 5;
    }

    var chartOptions = {
        chart: {
            style: {
                fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            zoomType: 'x',
            events: {
                selection: function(event) {
                    if(event.xAxis != null){
                        isSelected = true;
                        minXAxis = event.xAxis[0].min;
                        maxXAxis = event.xAxis[0].max;
                    }
                },
            },
        },
        legend: {
            symbolHeight: 12,
            symbolWidth: 12,
            symbolRadius: 6,
            enabled: true,
            floating: false,
            align: 'center',
            verticalAlign: 'top',
            // labelFormat: '<span>{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
            labelFormatter: function() {
                
                if(this.point.y == undefined){
                     return '<span style="color:' + this.color + '"> ' + this.name + ' </span> : <span style="min-width:50px"><b>  -  </b> '+this.tooltipOptions.valueSuffix+'</n>';
                }
                else{
                    return '<span style="color:'+ this.color +'"> ' + this.name + ' </span> : <span style="min-width:50px"><b> '+ kendo.toString(this.point.y,'n2')+' </b></span> <b>'+this.tooltipOptions.valueSuffix+'</b><br/>'
                }
               
            },
            borderWidth: 0,
            marginTop: -70,
        },
        scrollbar: {
            barBackgroundColor: 'gray',
            barBorderRadius: 7,
            barBorderWidth: 0,
            buttonBackgroundColor: 'gray',
            buttonBorderWidth: 0,
            buttonBorderRadius: 7,
            trackBackgroundColor: 'none',
            trackBorderWidth: 1,
            trackBorderRadius: 8,
            trackBorderColor: '#CCC'
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
            },{
                type: 'day',
                count: 5,
                text: '5d'
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
            },{
                type: '+',
                text: '+'
            }, {
                type: '-',
                text: '-'
            }
            ],
            inputEnabled: true,
            selected: rangeSelected,
            y: 50
        },
        navigator: {
            adaptToUpdatedData: false,
            series: {
                color: '#999',
                lineWidth: 0.7,
            },
            // margin: 2,
            xAxis: {
                dateTimeLabelFormats: dateTimeLabelFormats,
                labels: {
                    style: {
                        color: '#585555',
                        fontWeight: 'bold',
                    },
                }
            }
        },
        exporting: {
          enabled: true,
          buttons: {
                contextButton:{
                    enabled: false,
                },
                optionsButton: {
                    id: '_idoption',
                    text: 'Options',
                    symbol:'menu',
                    onclick: function () {
                        pg.options();
                    }
                },
                liveButton: {
                    id: '_idlive',
                    text: 'Live',
                    symbol: 'circle',
                    onclick: function () {
                        // alert('You pressed the button!');
                        pg.live(!pg.live());
                        pg.getDataStockChart();
                    },
                    enabled: pg.pageType() == 'HFD',
                }
            }
        },
        xAxis: {
            events: {
                afterSetExtremes: afterSetExtremes
            },
            type: 'datetime',
            // breaks: breaks,
            minRange: minRange,
            // ordinal: false,
            dateTimeLabelFormats : dateTimeLabelFormats,
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
             },
             shared: true,
        },
    };

    var chartCallback = function(e){
        setTimeout(function() {
            $.each($('.highcharts-button'), function(i, res){
                if (i==6){
                    $(this).attr("id", "");
                    $(this).attr("id", "zoomin");
                    $(this).attr("onclick", "pg.ZoomIn()");
                } else if (i==7){
                    $(this).attr("id", "");
                    $(this).attr("id", "zoomout");
                    $(this).attr("onclick", "pg.ZoomOut()");
                }
            });

        }, 200);
    };

    chart = new Highcharts.StockChart('chartTimeSeries', chartOptions, chartCallback);

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

pg.options = function(){
    $("#modalDetail").on("shown.bs.modal", function () { 
        $("#selectTagsDiv").html("");
        $("#selectTagsDiv").html('<select id="TagList"></select>');
        $('#TagList').kendoMultiSelect({
            dataSource: pg.tags(), 
            value: pg.TagList() ,
            dataValueField : 'value', 
            dataTextField: 'text',
            suggest: true, 
            minSelectedItems: 1,
            change: function(e) {
            
                var amountOfSelectedItems = this.dataItems().length;
               
                if (this.value().length == 0) {
                    this.value("windspeed")
                }
                if (amountOfSelectedItems > maxSelectedItems){
                    var values = [];
                    var data = (this.dataItems()).slice(0,-1);
                    $.each(data, function(key, val){
                        values.push(val.value);
                    });

                    this.value(values);
                    app.showError("Max. 5 Tags");
                }

            }
        });
        var multiselect= $("#TagList").data("kendoMultiSelect");
        multiselect.ul.addClass('hide-selected');
    }).modal('show');
}

// pg.hidePopover = function(){
//     $('.popover-markup>.trigger').popover('hide');
// }

pg.getDataStockChart = function(param){
    if(fa.LoadData()) {
        // if (param == "refresh") {
        $("#modalDetail").modal("hide");
        pg.dataType("MIN");
        // }
        app.loading(true);
        clearInterval(interval);
        if(param == "selectTags"){
           pg.TagList($("#TagList").val());
        //    $('.popover-markup>.trigger').popover("hide");
        }

        var IsHour = (pg.isFirst() == true ? false : true);

        var COOKIES = {};
        var cookieStr = document.cookie;
        var turbine = "";
        var project = "";
        
        // console.log(cookieStr);
        if(cookieStr.indexOf("turbine=") >= 0) {

            document.cookie = "turbine=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
            cookieStr.split(/; /).forEach(function(keyValuePair) {
                var cookieName = keyValuePair.replace(/=.*$/, "");
                var cookieValue = keyValuePair.replace(/^[^=]*\=/, "");
                COOKIES[cookieName] = cookieValue;
            });
            turbine = COOKIES["turbine"];
            project = COOKIES["project"];
        

            setTimeout(function(){
                $('#projectList').data('kendoDropDownList').value(project);
                var change = $("#projectList").data("kendoDropDownList").trigger("change");
                setTimeout(function(){
                    $('#turbineList').data('kendoDropDownList').value(turbine);
                },200);
                fa.project = project;
                fa.turbine = turbine;
            },500);

        } else {
            turbine = $('#turbineList').data('kendoDropDownList').value();
            project = $('#projectList').data('kendoDropDownList').value();
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
        
        var paramX = {
            period: fa.period,
            Turbine: turbine,
            DateStart: dateStart,
            DateEnd: dateEnd,
            Project: project,
            PageType: pg.pageType(),
            DataType: pg.dataType() ,
            TagList : pg.TagList(),
            IsHour : IsHour,
        };

        var url = "timeseries/getdatahfd";
        
        var request;
        if(pg.live() == false){
            request = toolkit.ajaxPost(viewModel.appName + url, paramX, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                var data = res.data.Data.Chart;
                var periods = res.data.Data.PeriodList;
                breaks = res.data.Data.Breaks;
                var outliers = res.data.Data.Outliers;

                pg.generateSeriesOption(data, periods);
                pg.generateOutliers(outliers);

                if (param=="first" || param=="refresh" || param=="selectTags"){
                    // if (seriesOriX){
                        seriesOriX = [];
                    // }
                    
                    $.each(seriesOptions,function(idx, val){
                        if (val.data != null){
                            var valx = val.data.slice(idx);
                            seriesOri[idx] = valx;
                        }
                    });

                    // if (seriesOri != null) {
                        seriesOriX = JSON.stringify(seriesOri);
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
        });
    }
}

pg.createLiveChart = function(IsHour){

        var dateStart, dateEnd; 
        // initiate first data

        seriesOriX = JSON.stringify(seriesOri);
        seriesOri = [];

        $.each(seriesOptions, function(i, s){
            s.data = [];
        });

        var seriesData = seriesOptions;
        var paramX = {
            period: fa.period,
            Turbine: fa.turbine,
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

            var data = res.data.Data.Chart;
            var periods = res.data.Data.PeriodList;
            var outliers = res.data.Data.Outliers;

            dateStart = new Date(new Date(Math.round(data[0].data[0][0])).toUTCString());

            pg.generateSeriesOption(data, periods);
            pg.generateOutliers(outliers);

            seriesOriX = JSON.stringify(seriesOri);
            seriesOri = [];

            $("#chartTimeSeries").html("");

            chart = Highcharts.stockChart('chartTimeSeries', {
                chart: {
                    style: {
                        fontFamily: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                    },
                    zoomType: 'x',
                    type: 'spline',
                    marginRight: 10,
                    events: {
                        load: function () {
                            var seriesData = this.series;
                            interval = setInterval(function () {
                                var paramX = {
                                    period: fa.period,
                                    Turbine: fa.turbine,
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

                                    var results = res.data.Data.Chart;
                                    var outliers = res.data.Data.Outliers;

                                    if(results.length > 0){
                                        $.each(results, function(id, tag){
                                            $.each(seriesData, function(i, series){
                                                if(series.name == tag.name){
                                                    // console.log(tag.name+" >> "+tag.data);
                                                    if (tag.data!=null){
                                                        $.each(tag.data, function(ix, dt){
                                                            dateStart = new Date(new Date(Math.round(dt[0])).toUTCString());
                                                            seriesData[i].addPoint(dt, true);
                                                        });
                                                    }
                                                }
                                            }); 
                                        });
                                    }else{
                                        return false;
                                    }

                                    chart.redraw(); 
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
                    // labelFormat: '<span>{name}</span> : <span style="min-width:50px"><b>{point.y:.2f} </b></span> <b>{tooltipOptions.valueSuffix}</b><br/>',
                    labelFormatter: function() {
                        if(this.point.y == undefined){
                             return '<span style="color:' + this.color + '"> ' + this.name + ' </span> : <span style="min-width:50px"><b>  -  </b> '+this.tooltipOptions.valueSuffix+'</n>';
                        }
                        else{
                            return '<span style="color:'+ this.color +'"> ' + this.name + ' </span> : <span style="min-width:50px"><b> '+ kendo.toString(this.point.y,'n2')+' </b></span> <b>'+this.tooltipOptions.valueSuffix+'</b><br/>'
                        }
                       
                    },
                    borderWidth: 0,
                    marginTop: -70,
                },
                rangeSelector: {
                    buttons: [{
                        type: 'minute',
                        count: 1,
                        text: '1"'
                    },{
                        type: 'minute',
                        count: 5,
                        text: '5"'
                    }, {
                        type: 'hour',
                        count: 1,
                        text: '1h'
                    }, {
                        type: 'all',
                        text: 'All'
                    }],
                    inputEnabled: false,
                    selected: 1 ,
                    // y: 50
                },
                navigator: {
                    adaptToUpdatedData: true,
                    series: {
                        color: '#999',
                        lineWidth: 0.7
                    },
                    xAxis: {
                        dateTimeLabelFormats: {
                            day: '%Y',
                            week: '%Y',
                            month: '%b %Y',
                            year: '%Y'
                        },
                        labels: {
                            style: {
                                color: '#585555',
                                fontWeight: 'bold',
                            },
                        }
                    }
                },
                exporting: {
                    enabled: true,
                    buttons: {
                            contextButton:{
                                enabled: false,
                            },
                            optionsButton: {
                                id: '_idoption',
                                text: 'Options',
                                symbol:'menu',
                                onclick: function () {
                                    pg.options();
                                    // alert('You pressed the button!');
                                    // $('.popover-markup>.trigger').popover('toggle');
                                }
                            },
                            liveButton: {
                                id: '_idlive',
                                text: 'Live',
                                symbol: 'circle',
                                symbolFill: '#31B445',
                                symbolStroke: '#31B445',
                                onclick: function () {
                                    pg.live(!pg.live());
                                    pg.getDataStockChart();
                                }
                            }
                        }
                },
                xAxis: {
                    type: 'datetime',
                    dateTimeLabelFormats : dateTimeLabelFormats,
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
                    },
                    shared: true,
                },
            });
        });

        

        // -----

        
    // });
    
}
pg.generateSeriesOption = function(data, periods){
    var IsHour = (pg.isFirst() == true ? false : true);
    var IsGroup = (pg.dataType() == "SEC" ? false : true);
    // console.log("isgroup: "+IsGroup);

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
            tickInterval: val.maxval/5,
            alignTicks: false,
            gridLineWidth: 0.75,
            endOnTick: false,
            startOnTick: false,
            showLastLabel: true,
            showFirstLabel: true,
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
            },
            // zones: [{
            //     value: val.minval,
            //     color: '#CD4B5B'
            // }, {
            //     color: colors[idx]
            // }, {
            //     value: val.maxval,
            //     color: '#CD4B5B'
            // }]
        }      

        seriesSelectedColor[idx] = val.name;

        legend[idx] = {
            name : val.name,
            unit : val.unit,
        }      

        xCounter+=1;

        seriesCounter += 1;
    });
}

pg.generateOutliers = function(data){
    if (data.length>0){
        var counter = yAxis.length;
        yAxis[counter] = {
            min: 0,
            max: 100, 
            gridLineWidth: 0,
            tickInterval: 100/5,
            alignTicks: false,
            endOnTick: false,
            startOnTick: false,
            labels: {
                format: '{value}',
            },
            visible: false,
        }

        seriesOptions[counter] = {
            type: 'column',
            name: "_err",
            data: data,
            color: '#ff0000',
            pointWidth: 1,
            yAxis: counter,
            id : "series_col"+counter,
            showInLegend : false,
            showInNavigator: false,
        }
    }
}

pg.ZoomIn = function(){
    var minXAxis = chart.xAxis[0].getExtremes().min;
    var maxXAxis = chart.xAxis[0].getExtremes().max;
    var newMin = (minXAxis + 12 * 3600 * 1000),  newMax = (maxXAxis - 12 * 3600 * 1000); 

    chart.xAxis[0].setExtremes(newMin,newMax);
}

pg.ZoomOut = function(){
    var minXAxis = chart.xAxis[0].getExtremes().min;
    var maxXAxis = chart.xAxis[0].getExtremes().max;
    var newMin = (minXAxis - 12 * 3600 * 1000), newMax = (maxXAxis + 12 * 3600 * 1000); 

    chart.xAxis[0].setExtremes(newMin,newMax);
}

pg.ToByProject = function(){
    setTimeout(function(){
        app.loading(true);
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        var project =  $('#projectList').data('kendoDropDownList').value();
        document.cookie = "project="+project.split("(")[0].trim()+";expires="+ newDateObj;
        window.location = viewModel.appName + "page/monitoringbyproject";
    },1500);
}


$(document).ready(function () {
    di.getAvailDate();
    newyAxis = yAxis;
    if(pg.pageType() === "HFD"){
        $("#periodList").closest(".k-widget").hide();
        $("#dateStart").closest(".k-widget").hide();
        $("#dateEnd").closest(".k-widget").hide();
        $(".label-filters:contains('Period')").hide();
        $(".label-filters:contains('to')").hide();
    }

    $('#btnRefresh').on('click', function () {
        $("#option1").prop("checked", true);
        $("#option2").prop("checked", true);
        pg.getDataStockChart("refresh");
    });

    setTimeout(function () {
        // pg.LoadData();
        pg.getDataStockChart("first");
        // pg.prepareScroll();
        // pg.hideRange();
        // pg.hideErr();

    }, 1000);
});
