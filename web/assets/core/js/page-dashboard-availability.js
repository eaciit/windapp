'use strict';

// var monthNames = ["January", "February", "March", "April", "May", "June",
//     "July", "August", "September", "October", "November", "December"
// ];

viewModel.availability = {};
var avail = viewModel.availability;
var isFleetDetail = false;
var lastDataChartLevel1 = [];
var detailDTLostEnergyTxtLevel1 = '';
var detailDTLETxtForDDLLevel1 = '';
var detailDTLETxtForDDLLevel2 = '';

avail.isDetailDTLostEnergy = ko.observable(false);
avail.detailDTLostEnergyTxt = ko.observable();
avail.LEFleetByDown = ko.observable(false);
avail.detailDTTopTxt = ko.observable();
avail.isDetailDTTop = ko.observable(false);
avail.mdTypeList = ko.observableArray([]);
avail.projectItem = ko.observableArray([]);
avail.fleetMachAvailData = ko.observableArray([]);
avail.fleetGridAvailData = ko.observableArray([]);
avail.DTLostEnergyData = ko.observableArray([]);
avail.LossCategoriesData = ko.observableArray([]);
avail.LossCategoriesDataSeries = ko.observableArray([]);


avail.ChartLineCfg = function(data, color, legend, height) {
    return { 
        dataSource: {
            data: data,
            sort: {
                field: "OrderNo",
                dir: "asc"
            }
        },
        chartArea: {
            background: "transparent",
            height: height,
        },
        title: {
            visible: false
        },
        legend: {
            visible: false
        },
        seriesDefaults: {
            type: "line"
        },
        series: [{
            field: "Value",
            name: "",
            missingValues : "gap",
            style: "smooth",
            markers: {
                visible: false
            },
            color: color
        }],
        categoryAxis: {
            field: "Title",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                template: '#=  value.substring(0,3) #',
                // visible: legend
            },
            crosshair: {
                visible: false,
            },
            majorGridLines: {
                visible: false,
            },
            majorTicks: {
                visible: false,
            },
            visible: legend
        },
        valueAxis: {
            visible: false,
            crosshair: {
                visible: false
            },
            majorGridLines: {
                visible: true,
                step: 100
            },
            majorTicks: {
                visible: false,
            },
        },
        tooltip: {
            visible: true,
            template: "#: category # = #= kendo.toString(value * 100, 'n2') # %"
        }
    };
};

avail.loadData = function () {

    var project = $("#projectId").data("kendoDropDownList").value();
    var param = {};

    if (project == "Fleet"){
        param = { ProjectName: project, Date: maxdate };
    }else{
        param = { ProjectName: project, Date: maxdate, Type: "All Types" };
    }

    if (lgd.isAvailability()) {
        var availReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getmachgridavailability", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (project == "Fleet") {
                avail.fleetMachAvailData(res.data.machineAvailability);
                // avail.fleetMachAvail(res.data.machineAvailability); /*"#fleetChartMachAvail"*/
                avail.fleetMachAvail(res.data.machineAvailability);
                avail.fleetGridAvailData(res.data.gridAvailability);
                avail.fleetGridAvail(res.data.gridAvailability); /*"#fleetChartGridAvail"*/
            } else {
                avail.projectMachAvail(res.data.machineAvailability); /*"#projectChartMachAvail"*/
                avail.projectGridAvail(res.data.gridAvailability); /*"#projectChartGridAvail"*/
            }
        });
        var lostEnergyReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getlostenergy", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (project == "Fleet") {
                if (res.data.lostenergy != null) {
                    avail.DTLostEnergyData(res.data.lostenergy);
                    avail.DTLostEnergy(res.data.lostenergy); /*"#chartDTLostEnergy"*/
                }
                // avail.DTTurbines();
            } else {
                avail.LossEnergyByType(res.data.lostenergy) /*#"projectChartLossEnergy"*/
            }
        });
        var downtimeTopReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimetop", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (project == "Fleet") {
                // avail.TopTurbineByLoss(res.data.loss); /*"#fleetChartTopTurbineLoss"*/
            } else {
                avail.DTLoss(res.data.loss); /*#"projectChartTopTurbineLosses"*/
                avail.DTDuration(res.data.duration); /*#"projectChartDTDuration"*/
                avail.DTFrequency(res.data.frequency); /*#"projectChartDTFrequency"*/
            }
        });
        var lossCatReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getlosscategories", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            if (project == "Fleet") {
                avail.LossCategoriesData(res.data);
                avail.LossCategoriesDataSeries(res.data.dataseries);
                avail.TLossCat('fleetChartTopLossCatEnergyLoss',true,res.data.lossCatLoss, 'MWh', res.data.dataseries); /*"#fleetChartTopLossCatEnergyLoss"*/
                avail.TLossCat('fleetChartTopLossCatDuration',false,res.data.lossCatDuration, 'Hours', res.data.dataseries); /*"#fleetChartTopLossCatDuration"*/
                avail.TLossCat('fleetChartTopLossCatFreq',false, res.data.lossCatFrequency , 'Times', res.data.dataseries); /*"#fleetChartTopLossCatFreq"*/
            } else {
                avail.TLossCat('projectChartTopLossCatEnergyLoss',true, res.data.lossCatLoss, 'MWh'); /*#"projectChartTopLossCatEnergyLoss"*/
                avail.TLossCat('projectChartTopLossCatDuration',false, res.data.lossCatDuration, 'Hours'); /*#"projectChartTopLossCatDuration"*/
                avail.TLossCat('projectChartTopLossCatFreq',false, res.data.lossCatFrequency, 'Times'); /*#"projectChartTopLossCatFreq"*/
            }
        });
        if (avail.mdTypeList.length == 0) {
            avail.getMDTypeList();
        }
        $.when(availReq, lostEnergyReq, downtimeTopReq, lossCatReq).done(function(){
            setTimeout(function(){
                app.loading(false);
            },300)
        });
    }
};

avail.refreshChart = function () {
    if(lgd.isAvailability() == true){
        if($("#projectId").data("kendoDropDownList").value() == 'Fleet'){
            if($("#fleetChartMachAvail").data("kendoChart") != undefined) {
                $("#fleetChartMachAvail").data("kendoChart").refresh();
            }
            if($("#fleetChartGridAvail").data("kendoChart") != undefined) {
                $("#fleetChartGridAvail").data("kendoChart").refresh();
            }
            if($("#fleetChartTopTurbineLoss").data("kendoChart") != undefined) {
                $("#fleetChartTopTurbineLoss").data("kendoChart").refresh();
            }
            if($("#fleetChartTopLossCatEnergyLoss").data("kendoChart") != undefined) {
                $("#fleetChartTopLossCatEnergyLoss").data("kendoChart").refresh();
            }
            if($("#fleetChartTopLossCatFreq").data("kendoChart") != undefined) {
                $("#fleetChartTopLossCatFreq").data("kendoChart").refresh();
            }
            if($("#fleetChartTopLossCatDuration").data("kendoChart") != undefined) {
                $("#fleetChartTopLossCatDuration").data("kendoChart").refresh();
            }
            if($("#chartDTLostEnergy").data("kendoChart") != undefined) {
                $("#chartDTLostEnergy").data("kendoChart").refresh();
            }
            // if($("#chartDTLEbyType").data("kendoChart") != undefined) {
            //     $("#chartDTLEbyType").data("kendoChart").refresh();
            // }
        }else{
            if($("#projectChartMachAvail").data("kendoChart") != undefined) {
                $("#projectChartMachAvail").data("kendoChart").refresh();
            }
            if($("#projectChartGridAvail").data("kendoChart") != undefined) {
                $("#projectChartGridAvail").data("kendoChart").refresh();
            }
            if($("#projectChartDTDuration").data("kendoChart") != undefined) {
                $("#projectChartDTDuration").data("kendoChart").refresh();
            }
            if($("#projectChartDTFrequency").data("kendoChart") != undefined) {
                $("#projectChartDTFrequency").data("kendoChart").refresh();
            }
            if($("#projectChartTopLossCatEnergyLoss").data("kendoChart") != undefined) {
                $("#projectChartTopLossCatEnergyLoss").data("kendoChart").refresh();
            }
            if($("#projectChartTopLossCatFreq").data("kendoChart") != undefined) {
                $("#projectChartTopLossCatFreq").data("kendoChart").refresh();
            }
            if($("#projectChartTopLossCatDuration").data("kendoChart") != undefined) {
                $("#projectChartTopLossCatDuration").data("kendoChart").refresh();
            }
        }
    }else{
        return;
    }
}



avail.TLossCat = function(id, byTotalLostenergy,dataSource,measurement, dataseries){
    var isStack = false;
    var catLossSeries = [{
            type: "column",
            field: "result",
        }];

    var colors = colorField;

    if(id.indexOf("fleet") >= 0) {
        isStack = true;
        catLossSeries = dataseries;
        colors = colorFieldProject;
    }   

    var templateLossCat = ''
    if(measurement == "MWh") {
       templateLossCat = "<b>#: category # :</b> #: kendo.toString(value/1000, 'n1')# " + measurement
    } else if(measurement == "Hours") {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n1')# " + measurement
    } else {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n0')# " + measurement
    }

    var margin = 0;

    if(dataseries != undefined){
        margin = -10;
    }

    $('#'+id).html("");
    $("#"+id).replaceWith('<div id='+id+'></div>');
    $('#'+id).kendoChart({
        dataSource: {
            data: dataSource,
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: id == "fleetChartTopLossCatEnergyLoss" ? 175 : 185,
            background: "transparent",
            padding: 0,
            margin: {
                top: margin
            }
        },
        seriesDefaults: {
            type: "column",
            stack: isStack,
        },
        series: catLossSeries,
        seriesColors: colors,
        valueAxis: {
            labels: {
                step: 2,
                template: (byTotalLostenergy == true) ? "#= value / 1000 #" : "#= value#",
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
            },

        },
        categoryAxis: {
            field: "_id.id2",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
        seriesClick: function (e) {
            if(id == 'fleetChartTopLossCatEnergyLoss'){
                avail.toDetailDTLELevel1(e, "fleetChartTopLossCatEnergyLoss");

            }
        }
    });
}
avail.fleetMachAvail = function (dataSource) {
    $("#fleetChartMachAvail").html("");
    $(".legend-machine-project").html("");

    var length = dataSource.length;
    var height = 160/length;
    $.each(dataSource, function(i, val){

        var chartProject = '<div class="col-md-12" id="chartMachine_'+i+'"></div>';
        var legendProject = '<span style="padding-left:5px;padding-right:5px;"><i class="fa fa-square" style="color:'+colorFieldProject[i]+';font-size:10px;padding-right:5px;"></i>'+val.Project+'</span>';
        var legend = false;

        if (i === dataSource.length - 1){ 
            legend = true;
            height = 200/length;
        }

        $("#fleetChartMachAvail").append(chartProject);
        $(".legend-machine-project").append(legendProject);
        $('#chartMachine_'+i).kendoChart(avail.ChartLineCfg(val.Details, colorFieldProject[i],legend,height));
    });
}

avail.fleetGridAvail = function (dataSource) {
    $("#fleetChartGridAvail").html("");
    $(".legend-grid-project").html("");

    var length = dataSource.length;
    var height = 160/length;
    $.each(dataSource, function(i, val){

        var chartProject = '<div class="col-md-12" id="chartGrid_'+i+'"></div>';
        var legendProject = '<span style="padding-left:5px;padding-right:5px;"><i class="fa fa-square" style="color:'+colorFieldProject[i]+';font-size:10px;padding-right:5px;"></i>'+val.Project+'</span>';
        var legend = false;

        if (i === dataSource.length - 1){ 
            legend = true;
            height = 200/length;
        }

        $("#fleetChartGridAvail").append(chartProject);
        $(".legend-grid-project").append(legendProject);
        $('#chartGrid_'+i).kendoChart(avail.ChartLineCfg(val.Details, colorFieldProject[i],legend,height));
    });
}

avail.DTLEbyType = function (dataSource) {
    var series = [];
    var categories = [];
    var dataLegends = [{id: "powerlost", name:"Power Lost"}, {id:"frequency",name:"Frequency"}, {id:"duration", name: "Duration"}];

    $.each(dataSource, function(idx, datas){
        $.each(datas.source, function(i, data){
            if($.inArray(data._id.id1, categories) === -1) categories.push(data._id.id1);
        });
    });

    categories = categories.sort();

    $.each(dataSource, function(idx, datas){
        $.each(dataLegends, function(key, legend){
            var serie = {
                name : legend.name, 
                visibleInLegend: false
            }
            var value = [];
            var idLegend = legend.id;
            $.each(datas.source, function(i, data){
                serie.stack = data._id.id3
                
                if(categories[i] == data._id.id1){
                    value.push(data[idLegend]);
                }else{
                    value.push(0);
                }

            });
            serie.data = value;
            series.push(serie);
        });
    });

    $("#chartDTLEbyType").kendoChart({
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: false,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 160
        },
        seriesDefaults: {
            type: "column",
            stacked: true
        },
        series: series, 
        seriesColors: ["#e65100", "#ff9800", "#ffb74d", 
                       "#00796b", "#4db6ac", "#80cbc4"],
        valueAxis: [{
            name: "PowerLost",
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
            categories: categories,
            majorGridLines: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            sharedTemplate: kendo.template($("#templateDTLEbyType").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        seriesClick: function (e) {
            avail.toDetailDTLELevel1(e, "chartbytype");
        }
    });
}

avail.DTLostEnergy = function (dataSource) {
    $("#chartDTLostEnergy").replaceWith('<div id="chartDTLostEnergy"></div>');
    $("#chartDTLostEnergy").kendoChart({
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 170,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7,
            stacked: false
        }],
        seriesColors: colorFieldProject,
        valueAxis: {
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
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            // template : "#: series.name # for #: category # : #:  kendo.toString(value, 'n0') #",
            // sharedTemplate: kendo.template($("#templateDowntimeLostEnergy").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        // seriesHover: function(e) {
        //   console.log(e);
        //   var positionX = e.originalEvent.clientX,
        //       positionY = e.originalEvent.clientY,
        //       value = e.value;
        //   $("#chartDTLostEnergyCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateDowntimeLostEnergy").html())({ e:e }));             
        // },
        seriesClick: function (e) {
            avail.toDetailLossEnergyLevel1(e, "chart");
        }
    });
    //  $("#chartDTLostEnergy").mouseleave(function(e){
    //    $("#chartDTLostEnergyCustomTooltip").hide();
    // })
}

avail.TopTurbineByLoss = function (dataSource) {
    $("#fleetChartTopTurbineLoss").kendoChart({
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
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 190,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [
        // {
        //     field: "AEBOK",
        //     name: "AEBOK"
        // }, 
        // {
        //     field: "ExternalStop",
        //     name: "External Stop"
        // }, 
        {
            field: "GridDown",
            name: "Grid Down"
        }, 
        // {
        //     field: "InternalGrid",
        //     name: "InternalGrid"
        // }, 
        {
            field: "MachineDown",
            name: "Machine Down"
        }, 
        // {
        //     field: "WeatherStop",
        //     name: "Weather Stop"
        // }, 
        {
            field: "Unknown",
            name: "Unknown"
        }],
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
            title: {
                text: "MWh",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: false, 
            },
            labels: {
                step: 2,
                template: "#: kendo.toString(value, 'n0') #",
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
            field: "_id",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n1') # MWh",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        seriesClick: function (e) {
            avail.toDetailDTTop(e, "MWh");
        }
    });
}

avail.DTLostEnergyFleet = function (dataSource) {
    $("#chartDTLostEnergyDetail").kendoChart({
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 160
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7,
            stacked: true
        }],
        seriesColors: colorFieldProject,
        valueAxis: {
            title : {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                text: "MWh",
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
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            // template : "#: series.name # for #: category # : #:  kendo.toString(value, 'n0') #",
            sharedTemplate: kendo.template($("#templateDowntimeLostEnergy").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        // seriesHover: function(e) {
        //   console.log(e);
        //   var positionX = e.originalEvent.clientX,
        //       positionY = e.originalEvent.clientY,
        //       value = e.value;
        //   $("#chartDTLostEnergyManehCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateDowntimeLostEnergy").html())({ e:e }));             
        // },
        seriesClick: function (e) {
            avail.toDetailDTLELevel2(e, "chart");
        }
    });
    // $("#chartDTLostEnergyDetail").mouseleave(function(e){
    //    $("#chartDTLostEnergyManehCustomTooltip").hide();
    // })
    setTimeout(function () {
        $("#chartDTLostEnergyDetail").data("kendoChart").refresh();
    }, 100);
}

avail.DTLostEnergyByDown = function (dataSource) {
    avail.LEFleetByDown(true)
    $("#chartDTLostEnergyDetail").kendoChart({
        dataSource: {
            data: dataSource,
            group: [{ field: "type" }],
            // sort: { field: "_id.id1", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 160
        },
        seriesDefaults: {
            type: "column"
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7
        }],
        seriesColors: colorField,
        valueAxis: {
            title: {
                text : "MWh",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
            categories: [lastParam.DateStr],
            majorGridLines: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            template: "#: series.name # : #:  kendo.toString(value, 'n1') # MWh",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        seriesClick: function (e) {
            avail.toDetailLossEnergyLevel2(e, "chart");
        }
    });

    setTimeout(function () {
        $("#chartDTLostEnergyDetail").data("kendoChart").refresh();
    }, 100);
}

avail.projectMachAvail = function (dataSource) {
    $("#projectChartMachAvail").kendoChart({
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 185,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7,
            stacked: true
        }],
        seriesColors: colorFieldProject,
        valueAxis: {
            labels: {
                // step: 2,
                template: '#=  value * 100 #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                // format: "{0:p1}",
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
            max : 1,
        },
        categoryAxis: {
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            // template : "#: series.name # for #: category # : #:  kendo.toString(value, 'n0') #",
            sharedTemplate: kendo.template($("#templateAvailPercentage").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });
}

avail.projectGridAvail = function (dataSource) {
    $("#projectChartGridAvail").kendoChart({
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 185,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7,
            stacked: true
        }],
        seriesColors: colorFieldProject,
        valueAxis: {
            labels: {
                // step: 2,
                template: '#=  value * 100 #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                // format: "{0:p1}",
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
            max : 1
        },
        categoryAxis: {
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            // template : "#: series.name # for #: category # : #:  kendo.toString(value, 'n0') #",
            sharedTemplate: kendo.template($("#templateAvailPercentage").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
    });
}

avail.LossEnergyByType = function (dataSource) {
    $('#projectChartLossEnergy').html("");
    $("#projectChartLossEnergy").kendoChart({
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
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 180,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true
        },
        series: [{
            type: "column",
            field: "result",
            // opacity : 0.7,
            stacked: true
        }],
        seriesColors: colorField,
        valueAxis: {
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
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            // template : "#: series.name # for #: category # : #:  kendo.toString(value, 'n0') #",
            sharedTemplate: kendo.template($("#templateDowntimeLostEnergy").html()),
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        // seriesHover: function(e) {
        //   console.log(e);
        //   var positionX = e.originalEvent.clientX,
        //       positionY = e.originalEvent.clientY,
        //       value = e.value;
        //   $("#chartDTLostEnergyManehCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateDowntimeLostEnergy").html())({ e:e }));             
        // },
        seriesClick: function (e) {
            projectSelected = $("#projectId").data("kendoDropDownList").value();
            avail.toDetailLossEnergyLevel2(e, "chart")
        }
    });

    setTimeout(function () {
        $("#projectChartLossEnergy").data("kendoChart").refresh();
    }, 100);
}

avail.DTDuration = function (dataSource) {
    $("#projectChartDTDuration").kendoChart({
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
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 190,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [
        // {
        //     field: "AEBOK",
        //     name: "AEBOK"
        // }, 
        // {
        //     field: "ExternalStop",
        //     name: "External Stop"
        // }, 
        {
            field: "GridDown",
            name: "Grid Down"
        }, 
        // {
        //     field: "InternalGrid",
        //     name: "InternalGrid"
        // }, 
        {
            field: "MachineDown",
            name: "Machine Down"
        }, 
        // {
        //     field: "WeatherStop",
        //     name: "Weather Stop"
        // }, 
        {
            field: "Unknown",
            name: "Unknown"
        }],
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
            title: {
                text: "Hours",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: false
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
            field: "_id",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n1') # Hours",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        seriesClick: function (e) {
            avail.toDetailDTTop(e, "Hours");
        }
    });
}

avail.DTLoss = function (dataSource) {
    $("#projectChartTopTurbineLosses").kendoChart({
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
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 190,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [
        // {
        //     field: "AEBOK",
        //     name: "AEBOK"
        // }, 
        // {
        //     field: "ExternalStop",
        //     name: "External Stop"
        // }, 
        {
            field: "GridDown",
            name: "Grid Down"
        }, 
        // {
        //     field: "InternalGrid",
        //     name: "InternalGrid"
        // }, 
        {
            field: "MachineDown",
            name: "Machine Down"
        }, 
        // {
        //     field: "WeatherStop",
        //     name: "Weather Stop"
        // }, 
        {
            field: "Unknown",
            name: "Unknown"
        }],
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
            title: {
                text: "MWh",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: false,
            },
            labels: {
                step: 2,
                template: "#: kendo.toString(value, 'n0') #",
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
            field: "_id",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            template: "#: category #: #: kendo.toString(value, 'n1') # MWh",
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        seriesClick: function (e) {
            avail.toDetailDTTop(e, "MWh");
        }
    });
}

avail.DTFrequency = function (dataSource) {
    $("#projectChartDTFrequency").kendoChart({
        dataSource: {
            data: dataSource,
            // group: [{field: "_id.id4"}],
            sort: { field: "Total", dir: 'desc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 190,
            background: "transparent",
            padding: 0,
            margin: {
                top: -10
            }
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [
        // {
        //     field: "AEBOK",
        //     name: "AEBOK"
        // }, 
        // {
        //     field: "ExternalStop",
        //     name: "External Stop"
        // }, 
        {
            field: "GridDown",
            name: "Grid Down"
        }, 
        // {
        //     field: "InternalGrid",
        //     name: "InternalGrid"
        // }, 
        {
            field: "MachineDown",
            name: "Machine Down"
        }, 
        // {
        //     field: "WeatherStop",
        //     name: "Weather Stop"
        // }, 
        {
            field: "Unknown",
            name: "Unknown"
        }],
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: "Times",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: false, 
            },
            name: "result",
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
            field: "_id",
            dir: "desc",
            majorGridLines: {
                visible: false
            },
            labels: {
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            // format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            // template: "#: category #: #: kendo.toString(value, 'n1') # Hours",
            template: "#: category #: #: value # ",
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        seriesClick: function (e) {
            avail.toDetailDTTop(e, "Times");
        }
    });
}

avail.DTLostEnergyDetail = function (dataSource) {
    $("#chartDTLostEnergyDetail").kendoChart({
        dataSource: {
            data: dataSource,
            // sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "Flat",
        chartArea: {
            height: 160
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        series: [{
            type: "column",
            field: "powerlost",
            // opacity : 0.7,
            axis: "EnergyLost",
            name: "Lost Energy (KWh)"
        }, {
            type: "line",
            field: "duration",
            axis: "duration",
            name: "Duration (Hours)",
            markers: {
                visible: false
            },
        }],
        seriesColors: colorField,
        valueAxes: [{
            name: "EnergyLost",
            title: { visible: true , text : "KWh", font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
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
        }, {
            name: "duration",
            title: { visible: false },
            visible: false
        }],
        categoryAxis: {
            field: "_id",
            majorGridLines: {
                visible: false
            },
            labels: {
                // template: '#=  value.substring(0,3) #'
                rotation: -330,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none",
            axisCrossingValues: [0, 30],
            justified: true
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            // template: "#: category #: #: kendo.toString(value, 'n1') #",
            shared: true,
            border: {
                color: "#eee",
                width: "2px",
            },

        }
    });

    setTimeout(function () {
        $("#chartDTLostEnergyDetail").data("kendoChart").refresh();
    }, 100);
}

avail.DTTopDetail = function (e, type) {
    app.loading(true);
    var project = $("#projectId").data("kendoDropDownList").value();
    var date = maxdate;
    var param = { ProjectName: project, Date: date, Type: e.series.field+"_"+type, Turbine: e.category };

    var templateTooltip = "#: category # : #:  kendo.toString(value, 'n1') #"+type
    if (type == 'Times') {
        templateTooltip = "#: category # : #:  kendo.toString(value, 'n0') #"+type
    } else if (type == 'MWh') {
        templateTooltip = "#: category # : #:  kendo.toString(value/1000, 'n1') #"+type
    }

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimetopdetail", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var dataSource = res.data

        $("#chartDTTopDetail").kendoChart({
            dataSource: {
                data: dataSource,
            },
            theme: "flat",
            title: {
                text: ""
            },
            legend: {
                position: "top",
                visible: false,
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            chartArea: {
                height: 160
            },
            seriesDefaults: {
                area: {
                    line: {
                        style: "smooth"
                    }
                }
            },
            series: [{
                // name : "Lost Energy",
                field: "result",
                // opacity : 0.7,
            }],
            seriesColors: colorField,
            valueAxis: {
                //majorUnit: 100,
                title : {
                    text : type,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                labels: {
                    step: 2,
                    template: (type == 'MWh' ? "#:  kendo.toString(value/1000, 'n0') #" : "#:  kendo.toString(value, 'n0') #"),
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
                field: "_id.id2",
                majorGridLines: {
                    visible: false
                },
                labels: {
                    template: '#=  value.substring(0,3) #',
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                majorTickType: "none"
            },
            tooltip: {
                visible: true,
                template: templateTooltip,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },

            }
        });

        app.loading(false);
        $("#chartDTTopDetail").data("kendoChart").refresh();
    });

    $('#gridDTTopDetail').html("");
    $('#gridDTTopDetail').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "dashboard/getdowntimetopdetailtable",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function (options) {
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.data
                },
                total: function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.total;
                }
            },
            sort: [
                { field: 'StartDate', dir: 'asc' }
            ],
        },
        groupable: false,
        sortable: true,
        filterable: false,
        pageable: {
            pageSize: 10,
            input: true, 
        },
        //resizable: true,
        columns: [
            { title: "Date", field: "StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80 },
            { title: "Start Time", field: "StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 75, attributes: { style: "text-align:center;" } },
            { title: "End Date", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80 },
            { title: "End Time", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 70, attributes: { style: "text-align:center;" } },
            { title: "Alert Description", field: "AlertDescription", width: 200 },
            { title: "External Stop", field: "ExternalStop", width: 80, sortable: false, template: '# if (ExternalStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Grid Down", field: "GridDown", width: 80, sortable: false, template: '# if (GridDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Internal Grid", field: "InternalGrid", width: 80, sortable: false, template: '# if (InternalGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Machine Down", field: "MachineDown", width: 80, sortable: false, template: '# if (MachineDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "AEbOK", field: "AEbOK", width: 80, sortable: false, template: '# if (AEbOK == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Unknown", field: "Unknown", width: 80, sortable: false, template: '# if (Unknown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Weather Stop", field: "WeatherStop", width: 80, sortable: false, template: '# if (WeatherStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
        ]
    });

}

avail.DTTurbines = function (turbineDownList) {
    $("#dtturbines").html("");
    $("#dtturbinesAvail").html("");
    if (turbineDownList.length == 0) {
        $("#dtturbines").html("<center><h2>NONE</h2></center>");
        $("#dtturbinesAvail").html("<center><h2>NONE</h2></center>");
    } else {
        $.each(turbineDownList, function (idx, val) {

            var btn = "btn-grey";
            var turbine = val._id;
            var value = val.result.toFixed(2);

            if (val.color == "red") {
                btn = "btn-danger";
            }

            $("#dtturbines").append('<div class="btn-group" role="group" style="margin: 3px;">' +
                '<button type="button" class="btn btn-sm ' + btn + '" style="width: 70px !important;">' + turbine + '</button>' +
                '<button type="button" class="btn btn-sm btn-warning" style="width: 40px !important;">' + value + '</button>' +
                '</div>');
            $("#dtturbinesAvail").append('<div class="btn-group" role="group">' +
                '<button type="button" class="btn btn-sm ' + btn + '" style="width: 70px !important;">' + turbine + '</button>' +
                '<button type="button" class="btn btn-sm btn-warning" style="width: 40px !important;">' + value + '</button>' +
                '</div>');
        });
    }
}

avail.toDetailLossEnergyLevel1 = function (e, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isAvailability(false);
    avail.isDetailDTLostEnergy(true);

    var param = {}; /*buat parameter tabel*/
    var paramChart = {}; /*buat parameter chart*/
    var method = "getdowntimefleetbydown"; /*nama controller untuk generate chart*/

    $(".show_hide_downtime").hide();
    $(".show_hide_project").show();

    if (source == "button" || source == "ddl") {
        if (source == "button") {
            /*set title label*/
            avail.detailDTLostEnergyTxt(detailDTLostEnergyTxtLevel1);

            /*set parameter for table*/
            param = lastParam;
            $("#projectList").data("kendoDropDownList").value(param.Type);

            /*create chart & table*/
            avail.DTLostEnergyByDown(lastDataChartLevel1);
           
            setTimeout(function(){
                app.loading(false);
            },50);
        } else if (source == "ddl") {
            projectSelected = $("#projectList").data("kendoDropDownList").value();
            /*set title label*/
            avail.detailDTLostEnergyTxt(detailDTLETxtForDDLLevel1 + projectSelected);
            detailDTLostEnergyTxtLevel1 = avail.detailDTLostEnergyTxt(); /*digunakan ketika tombol back dari level 2 ditekan*/

            /*set parameter for chart and table*/
            param = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType };
            paramChart = { ProjectName: projectSelected, DateStr: lastParamChart.DateStr, Type: dtType , IsDetail: true };
            lastParam = param;
            lastParamChart = paramChart;                

            /*create chart & table*/
            var chartRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                avail.DTLostEnergyByDown(res.data.lostenergy);
                lastDataChartLevel1 = res.data.lostenergy; /*data yang digunakan ketika tombol back dari level 2 ditekan*/
            });
            $.when(chartRequest).done(function(){
                setTimeout(function(){
                    app.loading(false);
                },50);
            });            
        }
    } else {
        monthDetailDT = e.category;
        $("#projectList").data("kendoDropDownList").value(e.series.name);
        $("#mdTypeListFleet").data("kendoDropDownList").value(0); /*karena yang di select ada project maka by default 'All Types'*/
        projectSelected = $("#projectList").data("kendoDropDownList").value();

        /*set title label*/
        avail.detailDTLostEnergyTxt("Lost Energy (MWh) for " + monthDetailDT + " - " + e.series.name);
        detailDTLostEnergyTxtLevel1 = avail.detailDTLostEnergyTxt(); /*digunakan ketika tombol back dari level 2 ditekan*/
        detailDTLETxtForDDLLevel1 = "Lost Energy (MWh) for " + monthDetailDT + " - "; /*digunakan ketika ingin ganti2 title via ddl */

        /*set parameter for chart and table*/
        paramChart = { ProjectName: projectSelected, DateStr: monthDetailDT, IsDetail: true };
        param = { ProjectName: projectSelected, DateStr: monthDetailDT };            
        lastParam = param;
        lastParamChart = paramChart;

        /*create chart & table*/
        var chartRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            avail.DTLostEnergyByDown(res.data.lostenergy);
            lastDataChartLevel1 = res.data.lostenergy; /*data yang digunakan ketika tombol back dari level 2 ditekan*/
        });
        // var tableRequest = avail.toDetailDTLETTable(param);

        // $.when(chartRequest, tableRequest).done(function(){
        $.when(chartRequest).done(function(){
            setTimeout(function(){
                app.loading(false);
            },50);
        });
    }
}

avail.toDetailLossEnergyLevel2 = function (e, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isAvailability(false);
    avail.isDetailDTLostEnergy(true);

    var project = $("#projectId").data("kendoDropDownList").value();
    var dateStr = '';
    var param = {};/*param grid & param chart is similar*/

    $(".show_hide_downtime").show();
    $(".show_hide_project").hide();
    if (project == "Fleet") {
        isFleetDetail = true;
    }

    if(source != "ddl") { /*karena kalo dari ddl, nilai variabel 'e' adalah null, tidak ada button karena level 2*/
        monthDetailDT = e.category;
        dateStr = e.category;
    }

    if (source == "chart") {
        if(projectSelected == "") { /*jika non fleet karena tidak melalui level 1*/
            projectSelected = project;
        }
        dtType = e.series.name;
        avail.detailDTLostEnergyTxt("Lost Energy (KWh) for " + monthDetailDT + " - " + dtType + " (" + projectSelected + ")");
        detailDTLETxtForDDLLevel2 = "Lost Energy (KWh) for " + monthDetailDT + " - ";
    } else if (source == "ddl") {
        var ddlVal = dtType;
        if (ddlVal == "") {
            ddlVal = "All Types";
        }
        avail.detailDTLostEnergyTxt(detailDTLETxtForDDLLevel2 + ddlVal + " (" + projectSelected + ")");
    } else if (source == "chartperproject") {
        dtType = $("#mdTypeListFleet").data("kendoDropDownList").value();
        projectSelected = e.series.name;
        avail.detailDTLostEnergyTxt("Lost Energy (KWh) for " + monthDetailDT + " - " + dtType + " (" + projectSelected + ")");
        detailDTLETxtForDDLLevel2 = "Lost Energy (KWh) for " + monthDetailDT + " - ";
    }

    if (project == "Fleet") {
        if (source == "chart") {
            param = { ProjectName: lastParamChart.ProjectName, DateStr: dateStr, Type: dtType };
            $("#mdTypeList").data("kendoDropDownList").value(dtType);
        } else if (source == "ddl") {
            param = { ProjectName: lastParamChart.ProjectName, DateStr: lastParamChart.DateStr, Type: dtType };
            $("#mdTypeList").data("kendoDropDownList").value(dtType);
        }
    } else {
        if (source == "chart") {
            param = { ProjectName: project, DateStr: dateStr, Type: dtType };
            lastParamLevel2 = param; /*butuh last param karena tidak berasal dari level 1 melainkan langsung ke level 2*/
            avail.LEFleetByDown(true); /*numpang variabel biar pas ganti ddl bisa balik lagi ke fungsi ini*/
            $("#mdTypeList").data("kendoDropDownList").value(dtType);
        } else if (source == "ddl") {
            param = { ProjectName: project, DateStr: lastParamLevel2.DateStr, Type: dtType };
            $("#mdTypeList").data("kendoDropDownList").value(dtType);
        }
    }

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimelostenergydetail", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var dataSource = res.data;
        avail.DTLostEnergyDetail(dataSource);

        app.loading(false);
    });

    // avail.toDetailDTLETTable(param);
   
    if(project == 'Fleet'){
        avail.toDetailDTLETTable(param);
    }
}

avail.toDetailDTLELevel1 = function (e, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isAvailability(false);
    avail.isDetailDTLostEnergy(true);

    var param = {}; /*buat parameter tabel*/
    var paramChart = {}; /*buat parameter chart*/
    var method = "getdowntime"; /*nama controller untuk generate chart*/

    $(".show_hide_downtime").hide();
    $(".show_hide_project").show();

    if (source == "button" || source == "ddl") {
        if (source == "button") {
            $("#projectList").data("kendoDropDownList").value(projectSelected);
            /*set title label*/
            avail.detailDTLostEnergyTxt(detailDTLostEnergyTxtLevel1);

            /*set parameter for table*/
            param = lastParam;
            param.DateStr = "fleet date";

            /*create chart & table*/
            avail.DTLostEnergyFleet(lastDataChartLevel1);
            setTimeout(function(){
                app.loading(false);
            },50);
        } else if (source == "ddl") {
            if (dtType == "") {
                dtType = "All Types"
            }
            projectSelected = $("#projectList").data("kendoDropDownList").value();
            /*set title label*/
            avail.detailDTLostEnergyTxt(detailDTLETxtForDDLLevel1 + dtType);
            detailDTLostEnergyTxtLevel1 = avail.detailDTLostEnergyTxt(); /*digunakan ketika tombol back dari level 2 ditekan*/

            /*set parameter for chart and table*/
            paramChart = { ProjectName: projectSelected, Date: lastParamChart.Date, Type: dtType , IsDetail: true };
            param = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType };
            lastParam = param;
            lastParamChart = paramChart;

            /*create chart & table*/
            var chartRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                avail.DTLostEnergyFleet(res.data.lostenergy);
                lastDataChartLevel1 = res.data.lostenergy; /*data yang digunakan ketika tombol back dari level 2 ditekan*/
            });
            // var tableRequest = avail.toDetailDTLETTable(param);
            // $.when(chartRequest, tableRequest).done(function(){
            $.when(chartRequest).done(function(){
                setTimeout(function(){
                    app.loading(false);
                },50);
            });
        }
    } else { /*chart*/
        monthDetailDT = e.category;

        if(source == "fleetChartTopLossCatEnergyLoss"){
            $("#projectList").data("kendoDropDownList").value(e.series.name);
        }else{
            $("#projectList").data("kendoDropDownList").select(0);
        }
        
        $("#mdTypeListFleet").data("kendoDropDownList").value(monthDetailDT);
        projectSelected = $("#projectList").data("kendoDropDownList").value();

        /*set title label*/
        avail.detailDTLostEnergyTxt("Lost Energy (MWh) for Last 12 months - " + monthDetailDT);
        detailDTLostEnergyTxtLevel1 = avail.detailDTLostEnergyTxt(); /*digunakan ketika tombol back dari level 2 ditekan*/
        detailDTLETxtForDDLLevel1 = "Lost Energy (MWh) for Last 12 months - "; /*digunakan ketika ingin ganti2 title via ddl */

        /*set parameter for chart and table*/
        paramChart = { ProjectName: projectSelected, Date: maxdate, Type: monthDetailDT, IsDetail: true };
        param = { ProjectName: projectSelected, DateStr: "fleet date", Type: monthDetailDT };
        lastParam = param;
        lastParamChart = paramChart;

        /*create chart & table*/
        var chartRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            avail.DTLostEnergyFleet(res.data.lostenergy);
            lastDataChartLevel1 = res.data.lostenergy; /*data yang digunakan ketika tombol back dari level 2 ditekan*/
        });
        $.when(chartRequest).done(function(){
            setTimeout(function(){
                app.loading(false);
            },50);
        });
    }
}

avail.toDetailDTLELevel2 = function (e, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isAvailability(false);
    avail.isDetailDTLostEnergy(true);

    var dateStr = '';
    var param = {};
    var paramChart = {};

    $(".show_hide_downtime").show();
    $(".show_hide_project").hide();
    isFleetDetail = true;

    if(source != "ddl") { /*karena kalo dari ddl, nilai variabel 'e' adalah null*/
        monthDetailDT = e.category;
        dateStr = e.category;
    }
    
    if (source == "chart") {
        dtType = $("#mdTypeListFleet").data("kendoDropDownList").value();
        /*set title label*/
        
        var ddlVal = dtType;
        if(projectSelected == "Fleet") {
            projectSelectedLevel2 = e.series.name;
            if (ddlVal == "") {
                ddlVal = "All Types";
            }
        } else {
            projectSelectedLevel2 = projectSelected
            ddlVal = e.series.name;
            dtType = ddlVal;
        }
        avail.detailDTLostEnergyTxt("Lost Energy (KWh) for " + monthDetailDT + " - " + ddlVal + " (" + projectSelectedLevel2 + ")");
        detailDTLETxtForDDLLevel2 = "Lost Energy (KWh) for " + monthDetailDT + " - ";

        $("#mdTypeList").data("kendoDropDownList").value(dtType);

        /*set param chart & table*/
        param = { ProjectName: lastParamChart.ProjectName, DateStr: dateStr, Type: dtType };
        lastParam = param;
    } else if (source == "ddl") {
        var ddlVal = dtType;
        if (ddlVal == "") {
            ddlVal = "All Types";
        }
        /*set label*/
        avail.detailDTLostEnergyTxt(detailDTLETxtForDDLLevel2 + ddlVal + " (" + projectSelectedLevel2 + ")");

        param = { ProjectName: lastParamChart.ProjectName, DateStr: lastParam.DateStr, Type: dtType };
        $("#mdTypeList").data("kendoDropDownList").value(dtType);
    }

    /*create chart & table*/
    var chartRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimelostenergydetail", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var dataSource = res.data;
        avail.DTLostEnergyDetail(dataSource);

        app.loading(false);
    });
    var tableRequest = avail.toDetailDTLETTable(param);

    $.when(chartRequest,tableRequest).done(function(){
        setTimeout(function(){
            app.loading(false);
        },50);
    });
}

avail.toDetailDTLETTable = function(param) {
    $('#gridDTLostEnergyDetail').html("");
    $('#gridDTLostEnergyDetail').kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "dashboard/getdowntimelostenergydetailtable",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function (options) {
                    return JSON.stringify(options);
                }
            },
            pageSize: 10,
            schema: {
                data: function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.data
                },
                total: function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.total;
                }
            },
            sort: [
                { field: 'StartDate', dir: 'asc' }
            ],
        },
        groupable: false,
        sortable: true,
        filterable: false,
        pageable: {
            pageSize: 10,
            input: true, 
        },
        //resizable: true,
        columns: [
            { title: "Date", field: "detail.StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80 },
            { title: "Turbine", field: "Turbine", width: 90, attributes: { style: "text-align:center;" } },
            { title: "Start Time", field: "detail.StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 75, attributes: { style: "text-align:center;" } },
            { title: "End Date", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80 },
            { title: "End Time", field: "EndDate", template: "#= kendo.toString(moment.utc(EndDate).format('HH:mm:ss'), 'HH:mm:ss') #", width: 70, attributes: { style: "text-align:center;" } },
            { title: "Alert Description", field: "AlertDescription", width: 200 },
            // { title: "External Stop", field: "ExternalStop", width: 80, sortable: false, template: '# if (ExternalStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Grid Down", field: "GridDown", width: 80, sortable: false, template: '# if (GridDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            // { title: "Internal Grid", field: "InternalGrid", width: 80, sortable: false, template: '# if (InternalGrid == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Machine Down", field: "MachineDown", width: 80, sortable: false, template: '# if (MachineDown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            // { title: "AEbOK", field: "AEbOK", width: 80, sortable: false, template: '# if (AEbOK == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            { title: "Unknown", field: "Unknown", width: 80, sortable: false, template: '# if (Unknown == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
            // { title: "Weather Stop", field: "WeatherStop", width: 80, sortable: false, template: '# if (WeatherStop == true ) { # <img src="../res/img/red-dot.png" /> # } else {# #}#', headerAttributes: { style: "text-align: center" }, attributes: { style: "text-align:center;" } },
        ]
    });
}

avail.toDetailDTTop = function (e, type) {
    vm.isDashboard(false);
    lgd.isAvailability(false);

    avail.detailDTTopTxt("(" + e.category + ") - " + type);
    avail.isDetailDTTop(true);

    // get the data and push into the chart    
    avail.DTTopDetail(e, type);
}

avail.backToDownTime = function () {
    vm.isDashboard(true);

    lgd.isSummary(false);
    lgd.isProduction(false);
    lgd.isAvailability(true);

    avail.isDetailDTLostEnergy(false);
    avail.detailDTLostEnergyTxt("Lost Energy for Last 12 months");

    avail.isDetailDTTop(false);
    avail.detailDTTopTxt("");
}

avail.setDownTimeSeriesCheck = function () {
    /*if (lgd.FleetDTLEDownType!=null){
        var chart = $('#chartDTLostEnergyDetail').data("kendoChart");
        var idx = 0;
        var found = -1;
        $.each(lgd.mdTypeList(),function(idx, val){
            if (val.value==lgd.FleetDTLEDownType){
                found=idx;
            }
        });

        if (found != -1) {
            $.each(lgd.mdTypeList(),function(idx, val){
                if (val.value!=lgd.FleetDTLEDownType){
                    chart._legendItemClick(idx);
                }
                idx++;
            });
        }
    }*/
}

avail.getMDTypeList = function () {
    app.ajaxPost(viewModel.appName + "/dashboard/getmdtypelist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0) {
            res.data = [];
        } else {
            avail.mdTypeList([]);

            if (res.data.length > 0) {
                /*var def = {};
                def.value = "All Type";
                def.text = "All Type";
                lgd.mdTypeList.push(def);*/
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    avail.mdTypeList.push(data);
                });
            }
        }
    });
};

avail.getDetailDT = function () { /*ddl change in level 2*/
    if (!lgd.isFirst()) {
        dtType = $("#mdTypeList").data("kendoDropDownList").value();

        // if (dtType == "") {
        //     avail.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - All Types");
        // } else {
        //     avail.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - " + dtType);
        // }
        if (avail.LEFleetByDown()) {
            avail.toDetailLossEnergyLevel2(null, "ddl");
        } else {
            avail.toDetailDTLELevel2(null, "ddl");
        }
    }
}

avail.getDetailDTFromProject = function () { /*ddl change & project in level 1*/
    if (!lgd.isFirst()) {
        projectSelected = $("#projectList").data("kendoDropDownList").value();
        dtType = $("#mdTypeListFleet").data("kendoDropDownList").value();
        if (avail.LEFleetByDown()) {
            avail.toDetailLossEnergyLevel1(null, "ddl");
        } else {
            avail.toDetailDTLELevel1(null, "ddl");
        }
    }
}

avail.backToDownTimeChart = function () {
    var project = $("#projectId").data("kendoDropDownList").value();
    if (project == "Fleet" && !lgd.isFirst() && isFleetDetail == true) {
        vm.isDashboard(false);
        avail.isDetailDTLostEnergy(true);
        isFleetDetail = false;
        if (avail.LEFleetByDown()) {
            avail.toDetailLossEnergyLevel1(null, "button");
        } else {
            avail.toDetailDTLELevel1(null, "button");
        }
        if ($("#projectList").data("kendoDropDownList") != null) {
            $("#projectList").data("kendoDropDownList").value(projectSelected);
        }
        $("#gridDTLostEnergyDetail").html("");
    } else {
        avail.LEFleetByDown(false);
        vm.isDashboard(true);
        lgd.isSummary(false);
        lgd.isProduction(false);
        lgd.isAvailability(true);
        avail.isDetailDTLostEnergy(false);
        avail.detailDTLostEnergyTxt("Lost Energy (MWh) for Last 12 months");
        avail.isDetailDTTop(false);
        avail.detailDTTopTxt("");
    }
}
$( window ).resize(function() {
    avail.refreshChart();
});
