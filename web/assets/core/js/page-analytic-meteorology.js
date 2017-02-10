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
pm.sectionsBreakdownList = ko.observableArray([
    { "text": 36, "value": 36 },
    { "text": 24, "value": 24 },
    { "text": 12, "value": 12 },
]);
var colorFieldsWR = ["#000292", "#005AFD", "#25FEDF", "#EBFE14", "#FF4908", "#9E0000", "#ff0000"];
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

pm.newData = ko.observableArray([]);
pm.Column = ko.observableArray([]);
pm.datas = ko.observableArray([]);

pm.isMet = ko.observable(true);
pm.isFirstAverage = ko.observable(true);
pm.isFirstWindRose = ko.observable(true);
pm.isFirstWindDis = ko.observable(true);
pm.isFirstTurbulence = ko.observable(true);
pm.isFirstTemperature = ko.observable(true);
pm.isFirstTurbine = ko.observable(true);
pm.isFirstTwelve = ko.observable(true);
pm.isFirstWindRoseComparison = ko.observable(true);


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

pm.getBackground = function(index, da){
    var color = 'white';
    var opacity = 1;

    var rgba = 'rgba(255,255,255)';
    if(pm.newData().length != 0){
        if (da in pm.newData()[index]){
            color = pm.newData()[index][da].Color;
            opacity = pm.newData()[index][da].Opacity
            if(color == "red") { 
                rgba = 'rgba(255,0,0,'+opacity+')';
            }else if(color == "green"){
                 rgba = 'rgba(0,128,0,'+opacity+')';
            }else{
                 rgba = 'rgba(255,255,255,'+opacity+')';
            }
        }
    }
    
    return rgba;
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
        var heat;
         toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getwindcorrelation", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            dataSource = res.data.Data;
            columns = res.data.Column;
            heat = res.data.Heat;

            pm.datas(dataSource);
            pm.newData(heat);
            pm.Column(columns);



            var schemaModel = {};
            var columnArray = [];

            $.each(columns, function (index, da) {
                schemaModel[da] = {type: (da == "Turbine" ? "string" : "int")};

                var column = {
                    title: da,
                    field: da,
                    locked: (da == "Turbine" ? true : false),
                    headerAttributes: {
                        style: "text-align: center;",
                    },
                    attributes: {
                        style: "text-align:center",
                        turbine: da,
                        index: index,
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
                    
                    var ini = this.wrapper;
                    $.each(pm.Column(), function(i, col){
                        var columns = e.sender.columns;
                        var columnIndex = ini.find(".k-grid-header [data-field=" + col + "]").index();

                        // iterate the data items and apply row styles where necessary
                        var dataItems = e.sender.dataSource.view();
                        for (var j = 0; j < dataItems.length; j++) {

                            var units = dataItems[j].get(col);
      
                            var row = e.sender.tbody.find("[data-uid='" + dataItems[j].uid + "']");
                            var cell = row.children().eq(columnIndex);

                            cell.css({"background": pm.getBackground(j,col)});
                        }
                    });


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
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

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
        tb.generateGridTable(this.id);
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
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="wr.showHideLegendWR(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
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