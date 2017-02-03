'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.type = ko.observable();
pg.detailDTTopTxt = ko.observable();
pg.isDetailDTTop = ko.observable(false);
pg.periodDesc = ko.observable();

pg.breakDown = ko.observableArray([]);
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
pg.SetBreakDown = function () {
    fa.disableRefreshButton(true);
    pg.breakDown = [];

    setTimeout(function () {
        $.each(fa.GetBreakDown(), function (i, val) {
            if (val.value == "Turbine" || val.value == "Project") {
                return false;
            } else {
               pg.breakDown.push(val);
            }
        });

        $("#breakdownlistavail").data("kendoDropDownList").dataSource.data(pg.breakDown);
        $("#breakdownlistavail").data("kendoDropDownList").dataSource.query();
        $("#breakdownlistavail").data("kendoDropDownList").select(0);


        fa.disableRefreshButton(false);
    }, 500);
}
pg.LoadData = function(){
    fa.LoadData();
    if (fa.project == "") {
        pg.type = "Project Name";
    } else {
        pg.type = "Turbine";
    }
    pg.getDataAvailableInfo();
}
pg.GenChartDownAlarmComponent = function (dataSource,id,Series,legend,name,axisLabel, vislabel,rotate,heightParam,wParam) {

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
                        template: "#= category #: \n #= kendo.format('{0:N1}', value)# " + axisLabel,
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
                step: 2
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
                rotation: rotate
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n1') # " + axisLabel,
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
pg.TLossCat = function (id, byTotalLostenergy, dataSource, measurement) {
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
            visible: true
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
        seriesColor: colorField,
        valueAxis: {
            labels: {
                step: 2,
                template: (byTotalLostenergy == true) ? "#= value / 1000 #" : "#= value#"
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
            field: "_id.id2",
            title: {
                text: "Loss Categories",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
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
pg.DTLEbyType = function (dataSource) {
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
        seriesColor: colorField,
        valueAxis: [{
            name: "PowerLost",
            labels: {
                step: 2
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
                rotation: -330
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
pg.createChartAvailability = function (dataSource) {
    var series = dataSource.SeriesAvail;
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#944dff";

    $("#availabilityChart").replaceWith('<div id="availabilityChart"></div>');
    $("#availabilityChart").kendoChart({
        theme: "Flat",
        legend: {
            position: "top",
            visible: true,
        },
        series: series,
        seriesColors: colorField,
        chartArea :{
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209 + 100) ) / 2,
            height: 200,
            padding: 0,
            background: "transparent",
        },
        valueAxes: [{
            line: {
                visible: false
            },
            max: 100,
            min: 0,
            labels: {
                format: "{0}",
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availpercentage",
            title: { text: "Availability (%)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }, {
            visible: false,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availline",
            title: { text: "Production (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            categories: categories,
            title: {
                text: $("#breakdownlistavail").data("kendoDropDownList").value(),
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            axisCrossingValues: [0, 1],
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            shared: true,
            sharedTemplate: kendo.template($("#template").html()),
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });
    $("#availabilityChart").css("top",-60);
    $("#availabilityChart").css("margin-bottom",-25);
}
pg.createChartProduction = function (dataSource) {
    var seriesProd = dataSource.SeriesProd;
    var categories = dataSource.Categories;
    var max = dataSource.Max;
    var min = dataSource.Min;
    colorField[0] = "#ff880e";

    $("#productionChart").replaceWith('<div id="productionChart"></div>');
    $("#productionChart").height(height);
    $("#productionChart").kendoChart({
        height: "150px",
        theme: "Flat",
        legend: {
            position: "top",
            visible: true,
        },
        chartArea :{
            // height: ($(".content-wrapper").height() - ($("#filter-analytic").height()+209 + 100) ) / 2,
            height: 200, 
            margin : 0,
            padding: 0,
            background: "transparent",
        },
        series: seriesProd,
        seriesColors: colorField,
        valueAxes: [{
            max: max,
            visible: true,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availline",
            title: { text: "Production (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            categories: categories,
            title: {
                text: $("#breakdownlistavail").data("kendoDropDownList").value(),
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: "#= series.name # at #= category # : #= kendo.toString(value, 'n2')# MWh",
            shared: false,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });
     $("#productionChart").css("top",-60);
     $("#productionChart").css("margin-bottom",-60);
}
pg.StaticView = function(){
    fa.LoadData();
    app.loading(true);
    if(pg.isFirstStaticView() === true){
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine,
            project: fa.project,
        };

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarylist", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            var gData = res.data.Data

            $('#lossGrid').html("");
            $('#lossGrid').kendoGrid({
                dataSource: {
                    data: gData,
                    pageSize: 10,
                    aggregate: [
                        { field: "Production", aggregate: "sum" },
                        { field: "LossEnergy", aggregate: "sum" },
                        { field: "MachineDownHours", aggregate: "sum" },
                        { field: "GridDownHours", aggregate: "sum" },
                        { field: "EnergyyMD", aggregate: "sum" },
                        { field: "EnergyyGD", aggregate: "sum" },
                        { field: "ElectricLoss", aggregate: "sum" },
                        { field: "PCDeviation", aggregate: "sum" },
                        { field: "Others", aggregate: "sum" },
                    ]
                },
                groupable: false,
                sortable: true,
                filterable: false,
                // height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
                height: 399,
                pageable: {
                    pageSize: 10,
                    input: true, 
                },
                columns: [
                    { title: pg.type,field: "Id",width: 100,attributes: {style: "text-align:center;"},headerAttributes: {style: "text-align:center;"},footerTemplate: "<center>Total (All Turbines)</center>"}, 
                    { title: "Production (MWh)", headerAttributes: { tyle: "text-align:center;"}, field: "Production",width: 100,attributes: { class: "align-center" },format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" }, 
                    { title: "Lost Energy (MWh)",headerAttributes: {style: "text-align:center;"},field: "LossEnergy", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"},
                    {
                        title: "Downtime : Duration (Hrs)",
                        headerAttributes: {
                            style: 'font-weight: bold; text-align: center;'
                        },
                        columns: [
                            {
                                title: "Machine",
                                headerAttributes: { style: "text-align:center;" },
                                field: "MachineDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                            {
                                title: "Grid",
                                headerAttributes: { style: "text-align:center;" },
                                field: "GridDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                            },
                        ]
                    }, {
                        title: "Downtime : Energy Loss (MWh)",
                        headerAttributes: {
                            style: 'font-weight: bold; text-align: center;'
                        },
                        columns: [
                            {
                                title: "Machine",
                                headerAttributes: { style: "text-align:center;" },
                                field: "EnergyyMD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                            {
                                title: "Grid",
                                headerAttributes: { style: "text-align:center;" },
                                field: "EnergyyGD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                        ]
                    }, {
                        title: "Electrical Losses (MWh)",
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "ElectricLoss", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }, {
                        title: "Power Curve Deviation (MW)", //Sepertinya ini MW
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "PCDeviation", width: 120, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }, {
                        title: "Others (MWh)", //Sepertinya ini KWh
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "Others", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }
                ],
                dataBound: function(){
                     app.loading(false);
                     pg.isFirstStaticView(false);
                }
            })
        });
    }else{
        $("#lossGrid").data("kendoGrid").refresh();
        $('#availabledatestart').html(pg.availabledatestartscada());
        $('#availabledateend').html(pg.availabledateendscada());
    }
}
pg.Downtime = function(){
    app.loading(true);
    fa.LoadData();
    if(pg.isFirstDowntime() === true){
        var param = {
            period: fa.period,
            dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
            dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
            turbine: fa.turbine,
            project: fa.project,
        }

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getdowntimetab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                var HDowntime = $('#filter-analytic').width() * 0.2
                var wAll = $('#filter-analytic').width() * 0.275

                /*Downtime Tab*/
                pg.GenChartDownAlarmComponent(res.data.duration,'chartDTDuration',SeriesDowntime,true,"Turbine", "Hours",false,-330,HDowntime,wAll);
                pg.GenChartDownAlarmComponent(res.data.frequency,'chartDTFrequency',SeriesDowntime,true,"Turbine", "Times",false,-330,HDowntime,wAll);
                pg.GenChartDownAlarmComponent(res.data.loss,'chartTopTurbineLoss',SeriesDowntime,true,"Turbine","MWh",false,-330,HDowntime,wAll);

                pg.isFirstDowntime(false);
                app.loading(false);
            },300);
           
        }); 
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartalarm2());
            $('#availabledateend').html(pg.availabledateendalarm2());
            $("#chartDTDuration").data("kendoChart").refresh();
            $("#chartDTFrequency").data("kendoChart").refresh();
            $("#chartTopTurbineLoss").data("kendoChart").refresh();
            app.loading(false);
        },300);
    }
}
pg.Availability = function(){
    app.loading(true);
    fa.LoadData();
    if(pg.isFirstAvailability() === true){
        pg.breakDownVal = $("#breakdownlistavail").data("kendoDropDownList").value();
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine,
            project: fa.project,
            breakDown: pg.breakDownVal,
        };
        toolkit.ajaxPost(viewModel.appName + "analyticavailability/getdata", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            setTimeout(function(){
                pg.dataSource(res.data);
                pg.createChartAvailability(pg.dataSource());
                pg.createChartProduction(pg.dataSource());
                pg.isFirstAvailability(false);
                app.loading(false);
            },200);
        });
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartscada3());
            $('#availabledateend').html(pg.availabledateendscada3());
            $("#availabilityChart").data("kendoChart").refresh();
            $("#productionChart").data("kendoChart").refresh();
            app.loading(false);
        },200);
    }
}
pg.LossEnergy = function(){
    app.loading(true);
    fa.LoadData();
    if(pg.isFirstLostEnergy() === true){
        var paramdown = {
            Period: fa.period,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Turbine: fa.turbine,
            Project: fa.project
        };
        toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeloss", paramdown, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                pg.DTLEbyType(res.data);
            },200)
        });
        
        var param = {
            period: fa.period,
            dateStart: moment(Date.UTC((fa.dateStart).getFullYear(), (fa.dateStart).getMonth(), (fa.dateStart).getDate(), 0, 0, 0)).toISOString(),
            dateEnd: moment(Date.UTC((fa.dateEnd).getFullYear(), (fa.dateEnd).getMonth(), (fa.dateEnd).getDate(), 0, 0, 0)).toISOString(),
            turbine: fa.turbine,
            project: fa.project,
        }

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getlostenergytab", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            setTimeout(function(){
                pg.TLossCat('chartLCByTEL', true, res.data.catloss, 'MWh');
                pg.TLossCat('chartLCByDuration', false, res.data.catlossduration, 'Hours');
                pg.TLossCat('chartLCByFreq', false, res.data.catlossfreq, 'Times');

                app.loading(false);
                pg.isFirstLostEnergy(false);
            },300);
        });
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartalarm());
            $('#availabledateend').html(pg.availabledateendalarm());
            $("#chartLCByTEL").data("kendoChart").refresh();
            $("#chartDTLEbyType").data("kendoChart").refresh();
            $("#chartLCByDuration").data("kendoChart").refresh();
            $("#chartLCByFreq").data("kendoChart").refresh();
            app.loading(false);
        },200)
    }
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
pg.WindSpeed = function(){
    app.loading(true);
    fa.LoadData()
    if(pg.isFirstWindSpeed() === true){
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
                dataBound: function(){
                    pg.isFirstWindSpeed(false);
                    app.loading(false);
                }
            });
        });
    }else{
        setTimeout(function(){
            $('#availabledatestart').html(pg.availabledatestartscada2());
            $('#availabledateend').html(pg.availabledateendscada2());
            $("#windAvailabilityChart").data("kendoChart").refresh();
            app.loading(false);
        },200);
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
                    pg.GenChartDownAlarmComponent(componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll);
                    pg.GenChartDownAlarmComponent(componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll);
                    pg.GenChartDownAlarmComponent(componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll);
                }else{                    
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll);
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll);
                    pg.GenChartDownAlarmComponent(pg.dtCompponentAlarm().alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll);
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


$(function(){
    setTimeout(function(){
        pg.LoadData();
        pg.StaticView();
    },200);

    $('#btnRefresh').on('click', function () {
        pg.resetStatus();
        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $('#breakdownlistavail').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { pg.isFirstAvailability(true); pg.Availability(); },
    });

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { fa.showHidePeriod(pg.SetBreakDown()) }
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { pg.SetBreakDown() }
    });

    $("#dateStart").change(function () { fa.DateChange(pg.SetBreakDown()) });
    $("#dateEnd").change(function () { fa.DateChange(pg.SetBreakDown()) });

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
            pg.GenChartDownAlarmComponent(data.alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll);
            pg.GenChartDownAlarmComponent(data.alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll);
            pg.GenChartDownAlarmComponent(data.alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll);

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
            pg.GenChartDownAlarmComponent(componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll);
            pg.GenChartDownAlarmComponent(componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll);
            pg.GenChartDownAlarmComponent(componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll);

            pg.labelAlarm(" Downtime")
        }
    });

    /*$(window).resize(function() {
        $("#chartCADuration").data("kendoChart").refresh();
    });*/

})
