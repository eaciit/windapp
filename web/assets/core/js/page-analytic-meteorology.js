'use strict';

viewModel.MeteorologyAnalysis = new Object();
var pm = viewModel.MeteorologyAnalysis;

vm.currentMenu('Meteorology');
vm.currentTitle('Meteorology');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Meteorology', href: viewModel.appName + 'page/analyticmeteorology' }]);

var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));
var availDateList = {};

pm.type = ko.observable();
pm.detailDTTopTxt = ko.observable();
pm.isDetailDTTop = ko.observable(false);
pm.periodDesc = ko.observable();
pm.breakDown = ko.observableArray([]);
pm.breakDownList = ko.observableArray([
    { "value": "dateinfo.dateid", "text": "Date" },
    { "value": "dateinfo.monthdesc", "text": "Month" },
    { "value": "dateinfo.year", "text": "Year" },
    { "value": "projectname", "text": "Project" },
    { "value": "turbine", "text": "Turbine" },
]);
pm.dataSourceAverage = ko.observableArray();
pm.dataWindrose = ko.observableArray([]);
pm.dataWindroseGrid = ko.observableArray([]);
pm.dataWindroseEachTurbine = ko.observableArray([]);
pm.sectorDerajat = ko.observable(0);
pm.sectionsBreakdownList = ko.observableArray([
    { "text": 36, "value": 36 },
    { "text": 24, "value": 24 },
    { "text": 12, "value": 12 },
]);

var maxValue = 0;
var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
var listOfChart = [];
var listOfButton = {};
var listOfCategory = [
    { "category": "0 to 4m/s", "color": colorFieldsWR[0] },
    { "category": "4 to 8m/s", "color": colorFieldsWR[1] },
    { "category": "8 to 12m/s", "color": colorFieldsWR[2] },
    { "category": "12 to 16m/s", "color": colorFieldsWR[3] },
    { "category": "16 to 20m/s", "color": colorFieldsWR[4] },
    { "category": "20m/s and above", "color": colorFieldsWR[5] },
];

pm.turbineList = ko.observableArray([]);
pm.turbine = ko.observableArray([]);
pm.valueCategory = ko.observableArray([
    { "value": "powerGeneration", "text": "Power Generation (MW)" },
    { "value": "machine", "text": "Machine Availability" },
    { "value": "scada", "text": "Scada Availability" },
    { "value": "grid", "text": "Grid Availability" },
]);

var color = ["#B71C1C", "#E57373", "#F44336", "#D81B60", "#F06292", "#880E4F",
    "#4A148C", "#7B1FA2", "#9C27B0", "#BA68C8", "#1A237E", "#5C6BC0",
    "#1E88E5", "#0277BD", "#0097A7", "#26A69A", "#4DD0E1", "#81C784",
    "#8BC34A", "#1B5E20", "#827717", "#C0CA33", "#DCE775", "#FF6F00", "#A1887F",
    "#FFEE58", "#004D40", "#212121", "#607D8B", "#BDBDBD", "#FF00CC", "#9999FF"
];

pm.dataSourceTable = ko.observableArray();

pm.MetTowerColumn = ko.observableArray([
    {value: true, text: "Wind Speed (m/s)", _id:"metWs", index:0 },
    {value: true, text: "Temp (°C)", _id:"metTemp", index: 1},

]);

pm.TurbineColumn = ko.observableArray([
    {_id: "turbineWs", text: "Wind Speed (m/s)", value:true , index:0},
    {_id: "turbineTemp", text: "Temp (°C)", value:true , index:1},
    {_id: "turbinePower", text: "Power (kWH)", value: true, index:2},
]);

pm.isMet = ko.observable(true);
pm.isFirstAverage = ko.observable(true);
pm.isFirstWindRose = ko.observable(true);
pm.isFirstWindDis = ko.observable(true);
pm.isFirstTurbulence = ko.observable(true);
pm.isFirstTemperature = ko.observable(true);
pm.isFirstTurbine = ko.observable(true);
pm.isFirstTwelve = ko.observable(true);


pm.loadData = function () {
    setTimeout(function () {
        if (fa.project == "") {
            pm.type = "Project Name";
        } else {
            pm.type = "Turbine";
        }

    }, 100);
    
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateList.availabledatestartscada = kendo.toString(moment.utc(res.data.ScadaData[0]).format('DD-MMMM-YYYY'));
        availDateList.availabledateendscada = kendo.toString(moment.utc(res.data.ScadaData[1]).format('DD-MMMM-YYYY'));

        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

        availDateList.startScadaHFD = kendo.toString(moment.utc(res.data.ScadaDataHFD[0]).format('DD-MMMM-YYYY'));
        availDateList.endScadaHFD = kendo.toString(moment.utc(res.data.ScadaDataHFD[1]).format('DD-MMMM-YYYY'));

        availDateList.availabledatestartmet = kendo.toString(moment.utc(res.data.MET[0]).format('DD-MMMM-YYYY'));
        availDateList.availabledateendmet = kendo.toString(moment.utc(res.data.MET[1]).format('DD-MMMM-YYYY'));
    })
}

// AVERAGE WINDSPEED
pm.generateGridAverage = function () {
    var config = {
        dataSource: {
            data: pm.dataSourceAverage(),
            pageSize: 10
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Turbine(s)", field: "turbine", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false },
        ],
    };

    $.each(pm.dataSourceAverage()[0].details, function (i, val) {
        var wra = val.col.WRA;        
        var column = {
            title: val.time + " <br/> WRA "+wra+ " (m/s)",
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            columns: [],
            width: 120
        }

        // var keyIndex = ["WRA", "Onsite"];
        var keyIndex = ["Onsite"];
        var j = 0;        

        $.each(keyIndex, function(j, key){
            // wra = 
            var colChild = {
                title: key + " (m/s)",                
                field: "details["+i+"].col."+ key ,
                width: 120,
                attributes: { class: "align-center row-custom" },
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;',
                },
                format: "{0:n2}",
                filterable: false
            };
            column.columns.push(colChild);
        });

        config.columns.push(column);
    });

    $('#gridAvgWs').html("");
    $('#gridAvgWs').kendoGrid(config);
    $('#gridAvgWs').data('kendoGrid').refresh();
}

pm.AverageWindSpeed = function() {
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstAverage() === true){
        var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/averagewindspeed", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            pm.dataSourceAverage(res.data.Data.turbine);
            pm.generateGridAverage();
            app.loading(false);
            pm.isFirstAverage(false);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        });        
    }else{
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        setTimeout(function(){
            $("#gridAvgWs").data("kendoGrid").refresh();
            app.loading(false);
        }, 300);
    }

}

// WIND ROSE
pm.ExportWindRose = function () {
    var chart = $("#wr-chart").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "WindRose.pdf",
        });
    });
}
pm.initChart = function () {
    listOfChart = [];
    var breakDownVal = $("#nosection").data("kendoDropDownList").value();
    var stepNum = 1
    var gapNum = 1
    if (breakDownVal == 36) {
        stepNum = 3
        gapNum = 0
    } else if (breakDownVal == 24) {
        stepNum = 2
        gapNum = 0
    } else if (breakDownVal == 12) {
        stepNum = 1
        gapNum = 0
    }

    $.each(pm.dataWindroseEachTurbine(), function (i, val) {
        var name = val.Name
        if (name == "MetTower") {
            name = "Met Tower"
        }

        var idChart = "#chart-" + val.Name
        listOfChart.push(idChart);
        // var pWidth = $('body').width() * 0.235;//$('body').width() * ($(idChart).closest('div.windrose-item').width() - 2) / 100;
        var pWidth = 290;

        $(idChart).kendoChart({
            theme: "nova",
            chartArea: {
                width: pWidth,
                height: pWidth,
                padding: 25
            },

            title: {
                text: name,
                font: '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            legend: {
                position: "bottom",
                labels: {
                    template: "#= (series.data[0] || {}).WsCategoryDesc #"
                },
                visible: false,
            },
            dataSource: {
                data: val.Data,
                group: {
                    field: "WsCategoryNo",
                    dir: "asc"
                },
                sort: {
                    field: "DirectionNo",
                    dir: "asc"
                }
            },
            seriesColors: colorFieldsWR,
            series: [{
                type: "radarColumn",
                stack: true,
                field: "Contribution",
                gap: gapNum,
                border: {
                    width: 1,
                    color: "#7f7f7f",
                    opacity: 0.5
                },
            }],
            categoryAxis: {
                field: "DirectionDesc",
                visible: true,
                majorGridLines: {
                    visible: true,
                    step: stepNum
                },
                labels: {
                    font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    visible: true,
                    step: stepNum
                }
            },
            valueAxis: {
                labels: {
                    template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                    font: '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                majorUnit: 10,
                max: maxValue,
                min: 0
            },
            tooltip: {
                visible: true,
                template: "#= category #"+String.fromCharCode(176)+" (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });
    });
}

pm.ChangeSector = function(){
    pm.isFirstWindRose(true);
    pm.WindRose();
}

pm.WindRose = function(){
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstWindRose() === true){
        setTimeout(function () {
            var breakDownVal = $("#nosection").data("kendoDropDownList").value();
            var secDer = 360 / breakDownVal;
            pm.sectorDerajat(secDer);
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine,
                project: fa.project,
                breakDown: breakDownVal,
            };
            toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getflexidataeachturbine", param, function (res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                if (res.data.WindRose != null) {
                    var metData = res.data.WindRose;
                    maxValue = res.data.MaxValue;
                    pm.dataWindroseEachTurbine(metData);
                    pm.initChart();
                }

                app.loading(false);
                pm.isFirstWindRose(false);

            })
        }, 300);
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
    }else{
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
        setTimeout(function(){
            $.each(listOfChart, function(idx, elem){
                $(elem).data("kendoChart").refresh();
            });
            app.loading(false);
        }, 300);
    }
}

// WIND ROSE COMPARISON
pm.isFirstWindRoseComparison = ko.observable(true);
pm.sectorDerajatComparison = ko.observable(0);
pm.dataWindroseComparison = ko.observableArray([]);
var listOfChartComparison = [];
var listOfButtonComparison = {};
pm.showHideLegendComparison = function (index) {
    var idName = "btn" + index;
    listOfButtonComparison[idName] = !listOfButtonComparison[idName];
    if (listOfButtonComparison[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChartComparison, function (idx, idChart) {
       if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
          $(idChart).data("kendoChart").options.series[index].visible = listOfButtonComparison[idName];
          $(idChart).data("kendoChart").refresh();
        }
    });
}

pm.InitTurbineListCompare = function () {
    if (pm.dataWindroseComparison().Data.length > 1) {
        $("#checkAllCompare").html('<label id="checkAllLabel">' +
            '<input type="checkbox" id="showHideAllCompare" checked onclick="pm.showHideAllCompare(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHideCompare"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#checkAllCompare").html("");
    }

    $("#turbine-list-compare").html("");
    $.each(pm.dataWindroseComparison().Data, function (idx, val) {
        $("#turbine-list-compare").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="pm.showHideCompare(' + val.idxseries + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + val.idxseries + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="chk-' + val.idxseries + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="pm.showHideCompare(' + val.idxseries + ')" type="button">' + val.name + '</button>' +
            '</div>');
    });
}

pm.showHideAllCompare = function (e) {

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(pm.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = true;
            }
        });
        $('#labelShowHideCompare b').text('Select All');
    } else {
        $.each(pm.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHideCompare b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);

    $("#WRChartComparison").data("kendoChart").redraw();
}

pm.showHideCompare = function (idx) {
    var stat = false;

    $('#chk-' + idx).trigger('click');
    var chart = $("#WRChartComparison").data("kendoChart");
    var leTur = $('input[id*=chk-][type=checkbox]').length

    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAllCompare').prop('checked', true);
    } else {
        $('#showHideAllCompare').prop('checked', false);
    }

    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }

    if ($('#chk-' + idx).is(':checked')) {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = false
    }
    $("#WRChartComparison").data("kendoChart").redraw();
}

pm.initChartWRC = function () {
    listOfChartComparison = [];
    var dataSeries = pm.dataWindroseComparison().Data;
    var categories = pm.dataWindroseComparison().Categories;
    var nilaiMax = pm.dataWindroseComparison().MaxValue;

    // var paddingTitle = Math.floor($('.windrose-part').width() / 16.16) * -1;
    // var offsetLegend = Math.floor($('.windrose-part').width() / 3.46) * -1;
    var majorUnit = 10;
    if(nilaiMax < 40) {
        majorUnit = 5;
    }

    $("#WRChartComparison").kendoChart({
        theme: "flat",
        title: {
            // padding: {
            //   left: paddingTitle
            // },
            text: "Wind Rose Comparison",
            font: '16px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            visible: false
        },
        legend: {
            position: "right",
            // offsetX: offsetLegend,
            labels: {
                font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
            },
            visible: false
        },
        dataSource: {
            sort: {
                field: "DirectionNo",
                dir: "asc"
            }
        },
        series: dataSeries,
        categoryAxis: {
            categories: categories,
            visible: true,
            majorGridLines: {
                visible: false
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
            }
        },
        valueAxis: {
            labels: {
                template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            majorUnit: majorUnit,
            max: nilaiMax,
            min: 0
        },
        tooltip: {
            visible: true,
            // template: "#= category # #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
            template: "#= series.name # : #= category #"+String.fromCharCode(176)+" (#= kendo.toString(value, 'n2') #%)",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        }
    });
    pm.InitTurbineListCompare();
    // $('#WRChartComparison').data('kendoChart').options.chartArea.width = $('#WRChartComparison').height() + ($('#WRChartComparison').height()/4);
    // $('#WRChartComparison').data('kendoChart').refresh();
}

pm.WindRoseComparison = function(){
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstWindRoseComparison() === true){
        setTimeout(function () {
            // var breakDownVal = $("#nosectionComparison").data("kendoDropDownList").value();
            var breakDownVal = "36";
            var secDer = 360 / breakDownVal;
            // pm.sectorDerajatComparison(secDer);
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine,
                project: fa.project,
                breakDown: breakDownVal,
            };
            toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getwindrosedata", param, function (res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                if (res.data != null) {
                    pm.dataWindroseComparison(res.data);
                    pm.initChartWRC();
                }

                app.loading(false);
                pm.isFirstWindRoseComparison(false);

            })
        }, 300);
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
    }else{
        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
        var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
        $('#availabledatestart').html(metDate);
        $('#availabledateend').html(scadaDate);
        setTimeout(function(){
            $.each(listOfChartComparison, function(idx, elem){
                $(elem).data("kendoChart").refresh();
            });
            app.loading(false);
        }, 300);
    }
}

// WIND DISTRIBUTION
pm.populateTurbine = function(){
    pm.turbine([]);
    if(fa.turbine == ""){
        $.each(fa.turbineList(), function(i, val){
            if (i > 0){
                pm.turbine.push(val.text);
            }
        });
    }else{
        pm.turbine(fa.turbine);
    }

}

pm.InitRightTurbineList= function () {
    if (pm.turbine().length > 0) {
        pm.turbineList([]);
        $.each(pm.turbine(), function (i, val) {
            var data = {
                color: color[i],
                turbine: val
            }

            pm.turbineList.push(data);
        });
    }

    if (pm.turbineList().length > 1) {
        $("#showHideChk").html('<label>' +
            '<input type="checkbox" id="showHideAll" checked onclick="pm.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChk").html("");
    }

    $("#right-turbine-list").html("");
    $.each(pm.turbineList(), function (idx, val) {
        $("#right-turbine-list").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="pm.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + (idx) + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.turbine + '" checked id="chk-' + (idx) + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="pm.showHideLegend(' + (idx) + ')" type="button">' + val.turbine + '</button>' +
            '</div>');
    });
}

pm.ChartWindDistributon =  function () {
    var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project
    };

    toolkit.ajaxPost(viewModel.appName + "analyticwinddistribution/getlist", param, function (res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        if (pm.turbine().length == 0) {
            var turbine = []
            for (var i=0;i<res.data.Data.length;i++) {
                if ($.inArray( res.data.Data[i].Turbine, turbine ) == -1){
                    turbine.push(res.data.Data[i].Turbine);
                }
            }

            $.each(turbine, function (i, val) {
                var data = {
                    color: color[i],
                    turbine: val
                }

                pm.turbineList.push(data);
            });
        }

        $('#windDistribution').html("");
        var data = res.data.Data;

        $("#windDistribution").kendoChart({
            dataSource: {
                data: data,
                group: { field: "Turbine" },
                sort: { field: "Category", dir: 'asc' }
            },
            theme: "flat",
            title: {
                text: ""
            },
            legend: {
                position: "right",
                visible: false
            },
            chartArea: {
                height: 360
            },
            series: [{
                type: "line",
                style: "smooth",
                field: "Contribute",
                // opacity : 0.7,
                markers: {
                    visible: true,
                    size: 3,
                }
            }],
            seriesColors: color,
            valueAxis: {
                labels: {
                    format: "{0:p0}",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                line: {
                    visible: true
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                }
            },
            categoryAxis: {
                field: "Category",
                majorGridLines: {
                    visible: false
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    // rotation: 25
                },
                majorTickType: "none"
            },
            tooltip: {
                visible: true,
                // template: "Contribution of #= series.name # : #= kendo.toString(value, 'n4')# % at #= category #",
                template: "#= kendo.toString(value, 'p2')#",
                // shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },

            },
            dataBound: function(){
                app.loading(false);
                pm.isFirstWindDis(false);
            }
        });

        pm.InitRightTurbineList();

        // app.loading(false);
        $("#windDistribution").data("kendoChart").refresh();
    });
}

pm.showHideAllLegend = function (e) {

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(pm.turbine(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = true;
            }
        });
        /*$('#labelShowHide b').text('Untick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(pm.turbine(), function (i, val) {
            if($("#windDistribution").data("kendoChart").options.series[i] != undefined){
                $("#windDistribution").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check').css("visibility", 'hidden');
        /*$('#labelShowHide b').text('Tick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);

    $("#windDistribution").data("kendoChart").redraw();
}

pm.showHideLegendWR = function (index) {
    var idName = "btn" + index;
    listOfButton[idName] = !listOfButton[idName];
    if (listOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChart, function (idx, idChart) {
       if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
          $(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
          $(idChart).data("kendoChart").refresh();
        }
    });
}

pm.showHideLegend = function (idx) {
    var stat = false;

    $('#chk-' + idx).trigger('click');
    var chart = $("#windDistribution").data("kendoChart");
    var leTur = $('input[id*=chk-][type=checkbox]').length

    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAll').prop('checked', true);
    } else {
        $('#showHideAll').prop('checked', false);
    }

    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }

    if ($('#chk-' + idx).is(':checked')) {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = false
    }
    $("#windDistribution").data("kendoChart").redraw();
}

pm.WindDis = function(){
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstWindDis() === true){
        pm.populateTurbine();
        pm.ChartWindDistributon();
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
    }else{
        app.loading(false);
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        setTimeout(function () {
            $('#windDistribution').data('kendoChart').refresh();
            app.loading(false);
        }, 100);
    }
}

// Turbulence
pm.Turbulence = function(){

}

// Temperature and Season Plots
pm.Temperature = function(){

}

// Turbine Correlation
pm.TurbineCorrelation = function(){
    app.loading(true);
    fa.LoadData();

    if(pm.isFirstTurbine() === true){
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine,
            project: fa.project
        };
        var dataSource;
        var columns;
        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getwindcorrelation", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            dataSource = res.data.Data;
            columns = res.data.Column
            
            var schemaModel = {};
            var columnArray = [];

            $.each(columns, function (index, da) {
                schemaModel[da] = {type: (da == "Turbine" ? "string" : "int")};

                var column = {
                    title: da,
                    field: da,
                    locked: (da == "Turbine" ? true : false),
                    headerAttributes: {
                        style: "text-align: center"
                    },
                    attributes: {
                        style: "text-align:center;"
                    },
                    width: 70,
                    template:( da != "Turbine" ? "#= kendo.toString("+da+", 'n2') #" : "#= kendo.toString("+da+") #")
                }

                columnArray.push(column);
            });

            var schemaModelNew = kendo.data.Model.define({
                id: "Turbine",
                fields: schemaModel,
            });

            var knownOutagesDataSource = new kendo.data.DataSource({
                data: dataSource,
                schema: {
                    model: schemaModelNew
                }
            });
            $("#gridTurbineCorrelation").html("");
            $("#gridTurbineCorrelation").kendoGrid({
                dataSource: knownOutagesDataSource,
                columns: columnArray,
                filterable: false,
                sortable: false,
                dataBound: function (e) {
                    if (e.sender._data.length == 0) {
                        var mgs, col;
                        mgs = "No results found for";
                        col = 9;
                        var contentDiv = this.wrapper.children(".k-grid-content"),
                     dataTable = contentDiv.children("table");
                        if (!dataTable.find("tr").length) {
                            dataTable.children("tbody").append("<tr><td colspan='" + col + "'><div style='color:red;width:500px'>" + mgs + "</div></td></tr>");
                            if (navigator.userAgent.match(/MSIE ([0-9]+)\./)) {
                                dataTable.width(this.wrapper.children(".k-grid-header").find("table").width());
                                contentDiv.scrollLeft(1);
                            }
                        }
                    }  
                },
                pageable: false,
                scrollable: true,
                resizable: false,
                height:390,
            });

            setTimeout(function(){
                $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
                $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
                $("#gridTurbineCorrelation").data("kendoGrid").refresh(); 
                app.loading(false);
                pm.isFirstTurbine(false)    
            },200);

        });
    }else{
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        setTimeout(function(){
             app.loading(false);
             $("#gridTurbineCorrelation").data("kendoGrid").refresh();
             $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
             $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
        }, 500);
    }
}

// 12/24 table 
pm.generateGridTable = function (datatype) {
    app.loading(true);
    $('#gridTable1224').html('');

    var dataSource = [];
    if(datatype == "turbine") {
        dataSource = pm.dataSourceTable().DataTurbine;
    } else {
        dataSource = pm.dataSourceTable().DataMet;
    }
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
            { title: "Hours", field: "hours", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false },
        ],
         dataBound: function(){
            setTimeout(function(){
                $("#gridTable1224 >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","75px");
                // $("#gridTable1224 >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","75px");
                pm.refreshTable(datatype);
            },200);
        },
    };

    $.each(dataSource[0].details, function (i, val) {
        var column = {
            title: val.time,
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            columns: []
        }
        var keyIndex = [];
        if(datatype == "turbine") {
            keyIndex = ["WS", "Temp", "Power"];
        } else {
            keyIndex = ["WS", "Temp"];
        }

        $.each(keyIndex, function(j, key){
            var title = "";
            if(key == "WS") {
                title = key + " (m/s)";
            } else if(key == "Temp") {
                title = key + " (" + String.fromCharCode(176) + "C)";
            } else {
                title = key + " (kWH)";
            }
            var colChild = {
                title: title,                
                field: "details["+i+"].col."+ key,
                attributes: { class: "align-center row-custom" },
                width: 100,
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;',
                },
                format: "{0:n2}",
                filterable: false, 
            };
            column.columns.push(colChild);
        });

        config.columns.push(column);
    });
    
    
    $('#gridTable1224').kendoGrid(config);
}

pm.hideShowColumn = function(i, type){
    var grid = $("#gridTable1224").data("kendoGrid");  
    var columns = grid.columns;

    if($('[name='+type+']:checked').length < 1){
        toolkit.showError("Grid must show at least one column");
        $('#'+i._id).prop('checked', true);
        return false;    
    }else{
        $.each(columns, function(index, val){
            if(index > 0){
                var col = grid.columns[index].columns[i.index];
                if (col.hidden) {
                  grid.showColumn(col.field);
                } else {
                  grid.hideColumn(col.field);
                } 
            } 
        });
    }

}

pm.getObjects = function(obj, key, val){
    var objects = [];
    for (var i in obj) {
        if (!obj.hasOwnProperty(i)) continue;
        if (typeof obj[i] == 'object') {
            objects = objects.concat(pm.getObjects(obj[i], key, val));
        } else if (i == key && obj[key] == val) {
            objects.push(obj);
        }
    }
    return objects;
}

pm.refreshTable = function(datatype){
    var grid = $("#gridTable1224").data("kendoGrid");  
    var columns = grid.columns;
    var data = (datatype == "met" ? pm.MetTowerColumn() : pm.TurbineColumn());
    var results = $.each($('[name="chk-column-'+datatype+'"]:not(:checked)'), function(i, val){
        var diff = pm.getObjects(data, "_id", val.id);
        $.each(diff, function(a, res){
             $.each(columns, function(e, value){
                if(e > 0){
                    var col = grid.columns[e].columns[res.index];
                    grid.hideColumn(col.field);
                } 
            });
             
        });
    });

    $.when(results).done(function(){
        setTimeout(function(){
            app.loading(false);
        },300);
    })
}

pm.Table = function(datatype){
    app.loading(true);
    fa.LoadData();

    if(pm.isFirstTwelve() === true){
        if(datatype == undefined || datatype == ''){
            if($("#met").is(':checked')) {
                datatype = 'met';
            } else {
                datatype = 'turbine';
            }
        }else{
            datatype = datatype;
        }
        

        var param = {
            Turbine: fa.turbine,
            Project: fa.project,
        };

        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/table1224", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            pm.dataSourceTable(res.data);
            pm.generateGridTable(datatype);
            pm.isFirstTwelve(false); 
        });
    }else{
        if($("#met").is(':checked')) {
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
        } else {
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        }
        setTimeout(function(){
            pm.generateGridTable(datatype);
        }, 300);
    }
}

pm.resetStatus= function(){
    pm.isFirstAverage(true);
    pm.isFirstWindRose(true);
    pm.isFirstWindDis(true);
    pm.isFirstTurbulence(true);
    pm.isFirstTemperature(true);
    pm.isFirstTurbine(true);
    pm.isFirstTwelve(true);
    pm.isFirstWindRoseComparison(true);
}

$(function(){
    pm.loadData();
    pm.AverageWindSpeed();

    $('#btnRefresh').on('click', function () {
        pm.resetStatus();
        $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $("input[name=isMet]").on("change", function() {
        pm.generateGridTable(this.id);
        if($("#met").is(':checked')) {
            pm.isMet(true);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
        } else {
             pm.isMet(false);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        }
    });

    setTimeout(function () {
        $("#legend-list").html("");
        $.each(listOfCategory, function (idx, val) {
            var idName = "btn" + idx;
            listOfButton[idName] = true;
            $("#legend-list").append(
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="pm.showHideLegendWR(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
                '<span class="span-legend">' + val.category + '</span>'
            );
        });
        $("#nosection").data("kendoDropDownList").value(12);
    }, 300);

    /*setTimeout(function () {
        $("#legend-list-comparison").html("");
        $.each(listOfCategory, function (idx, val) {
            var idName = "btn" + idx;
            listOfButtonComparison[idName] = true;
            $("#legend-list-comparison").append(
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="pm.showHideLegendComparison(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
                '<span class="span-legend">' + val.category + '</span>'
            );
        });
        $("#nosectionComparison").data("kendoDropDownList").value(12);
    }, 300);*/
});