'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.detailDTTopTxt = ko.observable();
pg.isDetailDTTop = ko.observable(false);
pg.periodDesc = ko.observable();

pg.typeChart = ko.observable();
pg.dataSource = ko.observableArray();

pg.availabledatestartscada = ko.observable();
pg.availabledateendscada = ko.observable();
pg.availabledatestartscada2 = ko.observable();
pg.availabledateendscada2 = ko.observable();

pg.availabledatestartalarm = ko.observable();
pg.availabledateendalarm = ko.observable();

pg.availabledatestartscada3 = ko.observable();
pg.availabledateendscada3 = ko.observable();

pg.availabledatestartalarm2 = ko.observable();
pg.availabledateendalarm2 = ko.observable();

pg.availabledatestartwarning = ko.observable();
pg.availabledateendwarning = ko.observable();

pg.dtCompponentAlarm = ko.observable();
pg.labelAlarm = ko.observable("Downtime ");

var height = $(".content").width() * 0.125;

var SeriesDowntime = [{
    field: "AEBOK",
    name: "AEBOK"
}, {
    field: "ExternalStop",
    name: "External Stop"
}, {
    field: "GridDown",
    name: "Grid Down"
}, {
    field: "InternalGrid",
    name: "InternalGrid"
}, {
    field: "MachineDown",
    name: "Machine Down"
}, {
    field: "WeatherStop",
    name: "Weather Stop"
}, {
    field: "Unknown",
    name: "Unknown"
}]

var SeriesAlarm =  [{
                type: "pie",
                field: "result",
                categoryField: "_id",
            }]

pg.breakDownList = ko.observableArray([
    { "value": "dateinfo.dateid", "text": "Date" },
    { "value": "dateinfo.monthdesc", "text": "Month" },
    { "value": "dateinfo.year", "text": "Year" },
    { "value": "projectname", "text": "Project" },
    { "value": "turbine", "text": "Turbine" },
]);

pg.isFirstStaticView = ko.observable(true);
pg.isFirstDowntime = ko.observable(true);
pg.isFirstAvailability = ko.observable(true);
pg.isFirstLostEnergy = ko.observable(true);
pg.isFirstReliability = ko.observable(true);
pg.isFirstWindSpeed = ko.observable(true);
pg.isFirstWarning = ko.observable(true);
pg.isFirstComponentAlarm = ko.observable(true);


pg.getDataAvailableInfo =  function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var minDatetemp = new Date(res.data.ScadaData[0]);
        var maxDatetemp = new Date(res.data.ScadaData[1]);

        pg.availabledatestartscada(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartscada2(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada2(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartalarm(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendalarm(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

        pg.availabledatestartscada3(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
        pg.availabledateendscada3(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

        pg.availabledatestartalarm2(kendo.toString(moment.utc(res.data.Alarm[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendalarm2(kendo.toString(moment.utc(res.data.Alarm[1]).format('DD-MMMM-YYYY')));

        pg.availabledatestartwarning(kendo.toString(moment.utc(res.data.Warning[0]).format('DD-MMMM-YYYY')));
        pg.availabledateendwarning(kendo.toString(moment.utc(res.data.Warning[1]).format('DD-MMMM-YYYY')));

        $('#availabledatestart').html(pg.availabledatestartscada());
        $('#availabledateend').html(pg.availabledateendscada());
    })
}
pg.backToDownTime = function () {
    pg.isDetailDTTop(false);
    pg.detailDTTopTxt("");
}

pg.LoadData = function(){
    fa.LoadData();
    if (fa.project == "") {
        sv.type = "Project Name";
    } else {
        sv.type = "Turbine";
    }
    pg.getDataAvailableInfo();
}
pg.GenChartDownAlarmComponent = function (dataSource,id,Series,legend,name,axisLabel, vislabel,rotate,heightParam,wParam,format) {

    $("#" + id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: [
                { "field": "Total", "dir": "desc" }
            ],
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: legend,
            labels: {              
                template: "#: kendo.toString(replaceString(text))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        chartArea: {
            height: heightParam, 
            width: wParam, 
            padding: 0,
            margin: 0
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            labels: {
                        visible: vislabel,
                        background: "transparent",
                        template: "#= category #: \n #= kendo.format('{0:" + format + "}', value)# " + axisLabel,
                    }
        },
        series: Series,
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: axisLabel,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
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
            }
        },
        categoryAxis: {
            // visible: legend,
            field: "_id",
            title: {
                text: name,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: rotate,                
                template: "#: kendo.toString(replaceString(value))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none",
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: kendo.toString(replaceString(category)) #: #: kendo.toString(value, '" + format + "') # " + axisLabel,
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });

    setTimeout(function () {
        if ($("#" + id).data("kendoChart") != null) {
            $("#" + id).data("kendoChart").refresh();
        }
    }, 100);
}

pg.generateGrid = function (dataSource) {
    var config = {
        dataSource: {
            data: dataSource,
            pageSize: 10
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Warning Description", field: "desc", attributes: { class: "align-left row-custom" }, width: 200, locked: true, filterable: false },
            { title: "Total", field: "total", attributes: { class: "align-center row-custom" }, width: 50, locked: true, filterable: false },
        ],
        dataBound: function(){
            setTimeout(function(){
                $("#warningGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","20px");
                $("#warningGrid >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","20px");
                // app.loading(false);
            },200);
        },
    };

    if (dataSource.length > 0){
        $.each(dataSource[0].turbines, function (i, val) {
            var column = {
                title: val.turbine,
                field: "turbines["+i+"].count",
                attributes: { class: "align-center" },
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                width: 80
            }

            config.columns.push(column);
        });
    }else{
        var column = {
            title: "",
            attributes: { class: "align-center" },
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            width: 80
        }

        config.columns.push(column);
    }   

    $('#warningGrid').html("");
    $('#warningGrid').kendoGrid(config);
    $('#warningGrid').data('kendoGrid').refresh();

    // setTimeout(function() {
    //     app.loading(false);
    // }, 500);
}

pg.Reliability = function(){
    if(pg.isFirstReliability() === true){
        $('#availabledatestart').html(pg.availabledatestartscada2());
        $('#availabledateend').html(pg.availabledateendscada2());
    }else{
        $('#availabledatestart').html(pg.availabledatestartscada2());
        $('#availabledateend').html(pg.availabledateendscada2());
    }
}

pg.Warning = function(){
    fa.LoadData()
    app.loading(true);
    if(pg.isFirstWarning() === true){
        var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getwarning", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.Data.length != 0) {
                setTimeout(function(){
                    pg.generateGrid(res.data.Data);
                    app.loading(false);
                },200);
            }else{
                setTimeout(function(){
                    pg.generateGrid([]);
                    app.loading(false);
                },200);
            }
        });
    }else{
        setTimeout(function(){
            $("#warningGrid").data("kendoGrid").refresh();
            $('#availabledatestart').html(pg.availabledatestartwarning());
            $('#availabledateend').html(pg.availabledateendwarning());
            app.loading(false);
        },200);
        
    }
}
pg.Component = function(){
    app.loading(true)
    fa.LoadData();
    if(pg.isFirstComponentAlarm() === true){
        var param = {
            period: fa.period,
            dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
            dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
            turbine: fa.turbine,
            project: fa.project,
        }

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getcomponentalarmtab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                pg.dtCompponentAlarm(res.data)
                var HAlarm = $('#filter-analytic').width() * 0.235
                var wAll = $('#filter-analytic').width() * 0.275
                var componentduration = _.sortBy(pg.dtCompponentAlarm().componentduration, '_id');
                var componentfrequency = _.sortBy(pg.dtCompponentAlarm().componentfrequency, '_id');
                var componentloss = _.sortBy(pg.dtCompponentAlarm().componentloss, '_id');

                var id = $("#downtimeGroup .active").attr('id')

                if(id == 'lblComp'){
                    /*Component / Alarm Type Tab*/
                    pg.GenChartDownAlarmComponent(componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll,"N1");
                    pg.GenChartDownAlarmComponent(componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll,"N0");
                    pg.GenChartDownAlarmComponent(componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll,"N1");
                }else{                    
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll,"N1");
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll,"N0");
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll,"N1");
                }

                app.loading(false);
                pg.isFirstComponentAlarm(false);
            },300);
        }); 
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartalarm());
            $('#availabledateend').html(pg.availabledateendalarm());
            $("#chartCADuration").data("kendoChart").refresh();
            $("#chartCAFrequency").data("kendoChart").refresh();
            $("#chartCATurbineLoss").data("kendoChart").refresh();
            app.loading(false);
        },200); 
    }
}
pg.resetStatus = function(){
    pg.isFirstStaticView(true);
    pg.isFirstDowntime(true);
    pg.isFirstAvailability(true);
    pg.isFirstLostEnergy(true);
    pg.isFirstReliability(true);
    pg.isFirstWindSpeed(true);
    pg.isFirstWarning(true);
    pg.isFirstComponentAlarm(true);
}
vm.currentMenu('Losses and Efficiency');
vm.currentTitle('Losses and Efficiency');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Losses and Efficiency', href: viewModel.appName + 'page/analyticloss' }]);

function replaceString(value) {
    return value.replace(/_/gi, "  ");
}

$(function(){
    setTimeout(function(){
        pg.LoadData();
        sv.StaticView();
    },200);

    $('#btnRefresh').on('click', function () {
        pg.resetStatus();
        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $('#breakdownlistavail').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.isFirstAvailability(true); av.Availability(); },
    });

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { fa.showHidePeriod(av.SetBreakDown()) }
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { av.SetBreakDown() }
    });

    $("#dateStart").change(function () { fa.DateChange(av.SetBreakDown()) });
    $("#dateEnd").change(function () { fa.DateChange(av.SetBreakDown()) });

    $("input[name=IsAlarm]").on("change", function() {
        var HAlarm = $('#filter-analytic').width() * 0.235
        var wAll = $('#filter-analytic').width() * 0.275
    
        var data = pg.dtCompponentAlarm()
        if(this.id == "alarm"){   
            SeriesAlarm =  [{
                field: "result",
                name: "Downtime"
            }]             
            // ===== Alarm =====
            pg.GenChartDownAlarmComponent(data.alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll,"N1");
            pg.GenChartDownAlarmComponent(data.alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll,"N0");
            pg.GenChartDownAlarmComponent(data.alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll,"N1");

            pg.labelAlarm(" Top 10 Downtime")
        }else{     
             SeriesAlarm = [{
                type: "pie",
                field: "result",
                categoryField: "_id",
            }]           
            // ===== Component =====
            var componentduration = _.sortBy(data.componentduration, '_id');
            var componentfrequency = _.sortBy(data.componentfrequency, '_id');
            var componentloss = _.sortBy(data.componentloss, '_id');
            pg.GenChartDownAlarmComponent(componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll,"N1");
            pg.GenChartDownAlarmComponent(componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll,"N0");
            pg.GenChartDownAlarmComponent(componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll,"N1");

            pg.labelAlarm(" Downtime")
        }
    });

    /*$(window).resize(function() {
        $("#chartCADuration").data("kendoChart").refresh();
    });*/

})
