'use strict';

// var monthNames = ["January", "February", "March", "April", "May", "June",
//     "July", "August", "September", "October", "November", "December"
// ];

viewModel.availability = {};
var avail = viewModel.availability;
var isFleetDetail = false;

avail.isDetailDTLostEnergy = ko.observable(false);
avail.detailDTLostEnergyTxt = ko.observable();
avail.LEFleetByDown = ko.observable(false);
avail.detailDTTopTxt = ko.observable();
avail.isDetailDTTop = ko.observable(false);
avail.mdTypeList = ko.observableArray([]);
avail.projectItem = ko.observableArray([]);

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
                avail.fleetMachAvail(res.data.machineAvailability); /*"#fleetChartMachAvail"*/
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
                avail.DTLEbyType(res.data.lostenergybytype[0]); /*"#chartDTLEbyType"*/
                avail.DTLostEnergy(res.data.lostenergy); /*"#chartDTLostEnergy"*/
            } else {
                avail.LossEnergyByType(res.data.lostenergy) /*#"projectChartLossEnergy"*/
            }
        });
        var downtimeTopReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimetop", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (project == "Fleet") {
                avail.TopTurbineByLoss(res.data.loss); /*"#fleetChartTopTurbineLoss"*/
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
                avail.TLossCat('fleetChartTopLossCatEnergyLoss',true,res.data.lossCatLoss, 'MWh'); /*"#fleetChartTopLossCatEnergyLoss"*/
                avail.TLossCat('fleetChartTopLossCatDuration',false,res.data.lossCatDuration, 'Hours'); /*"#fleetChartTopLossCatDuration"*/
                avail.TLossCat('fleetChartTopLossCatFreq',false, res.data.lossCatFrequency , 'Times'); /*"#fleetChartTopLossCatFreq"*/
            } else {
                avail.TLossCat('projectChartTopLossCatEnergyLoss',true, res.data.lossCatLoss, 'MWh'); /*#"projectChartTopLossCatEnergyLoss"*/
                avail.TLossCat('projectChartTopLossCatDuration',false, res.data.lossCatDuration, 'Hours'); /*#"projectChartTopLossCatDuration"*/
                avail.TLossCat('projectChartTopLossCatFreq',false, res.data.lossCatFrequency, 'Times'); /*#"projectChartTopLossCatFreq"*/
            }
        });
        avail.DTTurbines();
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
            if($("#chartDTLEbyType").data("kendoChart") != undefined) {
                $("#chartDTLEbyType").data("kendoChart").refresh();
            }
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



avail.TLossCat = function(id, byTotalLostenergy,dataSource,measurement){
    var templateLossCat = ''
    if(measurement == "MWh") {
       templateLossCat = "<b>#: category # :</b> #: kendo.toString(value/1000, 'n1')# " + measurement
    } else if(measurement == "Hours") {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n1')# " + measurement
    } else {
        templateLossCat = "<b>#: category # :</b> #: kendo.toString(value, 'n0')# "
    }

    $('#'+id).html("");
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
            height: 160
        },
        seriesDefaults: {
            type: "column",
        },
        series: [{
            type: "column",
            field: "result",
        }],
        seriesColor: colorField,
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
    });
}
avail.fleetMachAvail = function (dataSource) {
    $("#fleetChartMachAvail").kendoChart({
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
            height: 195
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
                template: '#=  value * 100 #'
                // format: "{0:p0}",
            },
            line: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #'
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:p1}",
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

avail.fleetGridAvail = function (dataSource) {
    $("#fleetChartGridAvail").kendoChart({
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
            labels:{
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 195
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
            }
        },
        categoryAxis: {
            field: "_id.id2",
            majorGridLines: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            /*labels:{
                template: '#=  value.substring(0,3)+" "+value.substring((value.length-4),value.length) #'
            },*/
            labels: {
                template: '#=  value.substring(0,3) #'
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:p1}",
            // template : "#:  kendo.toString(value, 'n0') #",
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

avail.DTLEbyType = function (dataSource) {
    $("#chartDTLEbyType").kendoChart({
        dataSource: {
            data: dataSource.source,
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
            // min: dataSource.minPowerLost,
            // max: dataSource.maxPowerLost
        },
        {
            name: "Duration",
            title: { visible: false },
            visible: false,
            // min: dataSource.minDuration,
            // max: dataSource.maxDuration
        },
        {
            name: "Frequency",
            title: { visible: false },
            visible: false,
            // min: dataSource.minFreq,
            // max: dataSource.maxFreq
        }],
        categoryAxis: {
            field: "_id.id2",
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
        // seriesHover: function(e) {
        //   console.log(e);
        //   var positionX = e.originalEvent.clientX,
        //       positionY = e.originalEvent.clientY,
        //       value = e.value;
        //   $("#chartDTLEbyTypeCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateDowntimeLostEnergy").html())({ e:e }));             
        // },
        seriesClick: function (e) {
            avail.toDetailDTLostEnergy(e, false, "chartbytype");
        }
    });
    // $("#chartDTLEbyType").mouseleave(function(e){
    //    $("#chartDTLEbyTypeCustomTooltip").hide();
    // })
}

avail.DTLostEnergy = function (dataSource) {
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
        //   $("#chartDTLostEnergyCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateDowntimeLostEnergy").html())({ e:e }));             
        // },
        seriesClick: function (e) {
            avail.toDetailDTLostEnergy(e, false, "chart");
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
            height: 160
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
                template: "#: kendo.toString(value/1000, 'n0') #",
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
            template: "#: category #: #: kendo.toString(value/1000, 'n1') # MWh",
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

avail.DTLostEnergyManeh = function (dataSource) {
    avail.detailDTLostEnergyTxt("Lost Energy for Last 12 months");
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
            avail.toDetailDTLostEnergy(e, true, "chart");
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
            avail.toDetailDTLostEnergy(e, true, "chart");
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
            height: 195
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
            height: 195
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
            avail.toDetailDTLostEnergy(e, true, "chart");
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
            height: 160
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
            height: 160
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
                template: "#: kendo.toString(value/1000, 'n0') #",
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
            template: "#: category #: #: kendo.toString(value/1000, 'n1') # MWh",
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
            height: 160
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
            title: { visible: false },
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

avail.DTTopDetail = function (turbine, type) {
    app.loading(true);
    var project = $("#projectId").data("kendoDropDownList").value();
    var date = maxdate;
    var param = { ProjectName: project, Date: date, Type: type, Turbine: turbine };

    var templateTooltip = "#: category # : #:  kendo.toString(value, 'n1') #"
    if (type == 'Times') {
        templateTooltip = "#: category # : #:  kendo.toString(value, 'n0') #"
    } else if (type == 'MWh') {
        templateTooltip = "#: category # : #:  kendo.toString(value/1000, 'n1') #"
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

avail.DTTurbines = function () {
    var project = $("#projectId").data("kendoDropDownList").value();
    var date = maxdate;//new Date(Date.UTC(2016, 5, 30, 23, 50, 0, 0));
    var param = { ProjectName: project, Date: date };

    $("#dtturbines").html("");

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeturbines", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0){
            $("#dtturbines").html("<center><h2>NONE</h2></center>");
        }else{
            $.each(res.data, function (idx, val) {
            var btn = "btn-success";
            var turbine = val._id;
            var value = val.result.toFixed(2);

            if (val.isdown == true) {
                btn = "btn-danger";
            }

            $("#dtturbines").append('<div class="btn-group" role="group">' +
                '<button type="button" class="btn btn-sm ' + btn + '">' + turbine + '</button>' +
                '<button type="button" class="btn btn-sm btn-warning">' + value + '</button>' +
                '</div>');
            });
        }        
    });
}


avail.toDetailDTLostEnergy = function (e, isDetailFleet, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isAvailability(false);
    avail.isDetailDTLostEnergy(true);

    var project = $("#projectId").data("kendoDropDownList").value();
    var dateStr = '';
    var type = '';
    var param = {};
    var paramChart = {};
    var method = "getdowntime";

    if (source == "chart" || source == "chartbytype") {
        monthDetailDT = e.category;
        avail.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - " + e.series.name);

        if (source == "chart") {
            dateStr = e.category;
        }
        type = e.series.name;
    } else if (source == "button") {
        avail.detailDTLostEnergyTxt("Lost Energy for Last 12 months - " + lastParam.Type);
    }
    if (project == "Fleet" && isDetailFleet == false) { /*by type level 1*/
        $(".show_hide_downtime").hide();
        $(".show_hide_project").show();

        if (source == "button" || source == "ddl") {
            if (source == "button") {
                if (avail.LEFleetByDown() == true) {
                    method = "getdowntimefleetbydown";
                }

                param = lastParam;
                paramChart = lastParamChart;
                $("#projectList").data("kendoDropDownList").value(param.Type);
            } else if (source == "ddl") {
                if (avail.LEFleetByDown() == true) {
                    method = "getdowntimefleetbydown";
                    paramChart = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType , IsDetail: true };
                } else {
                    if (dtType == "") {
                        dtType = "All Types"
                    }
                    paramChart = { ProjectName: projectSelected, Date: lastParamChart.Date, Type: dtType , IsDetail: true };
                }

                param = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType };
                lastParam = param;
                lastParamChart = paramChart;
            }
            toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                if (method == "getdowntimefleetbydown") {
                    avail.DTLostEnergyByDown(res.data.lostenergy);
                } else {
                    avail.DTLostEnergyManeh(res.data.lostenergy);
                }

                app.loading(false);
                avail.setDownTimeSeriesCheck();
            });
        } else {
            $("#projectList").data("kendoDropDownList").select(0);
            projectSelected = $("#projectList").data("kendoDropDownList").value();

            if (source == "chart") {
                $("#mdTypeListFleet").data("kendoDropDownList").value(0);
                paramChart = { ProjectName: projectSelected, DateStr: e.category, IsDetail: true };
                param = { ProjectName: projectSelected, DateStr: e.category };

                method = "getdowntimefleetbydown";
            } else {
                $("#mdTypeListFleet").data("kendoDropDownList").value(e.category);
                paramChart = { ProjectName: projectSelected, Date: maxdate, Type: e.category, IsDetail: true };
                param = { ProjectName: projectSelected, DateStr: "fleet date", Type: e.category };
            }

            lastParam = param;
            lastParamChart = paramChart;

            toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                if (method == "getdowntimefleetbydown") {
                    avail.DTLostEnergyByDown(res.data.lostenergy);
                } else {
                    avail.DTLostEnergyManeh(res.data.lostenergy);
                }

                avail.FleetDTLEDownType = e.category;
                avail.setDownTimeSeriesCheck();
                app.loading(false);
            });
        }
    } else { /*bagian detail (level 2)*/
        $(".show_hide_downtime").show();
        $(".show_hide_project").hide();
        $("#projectList").data("kendoDropDownList").value(projectSelected);
        if (project == "Fleet" && isDetailFleet == true) {
            isFleetDetail = true;
        }

        // dtType = $("#mdTypeList").data("kendoDropDownList").value();

        if (dtType == "All Types") {
            dtType = "";
        }

        if (project == "Fleet") {
            if (source == "chart") {
                param = { ProjectName: lastParamChart.ProjectName, DateStr: dateStr, Type: type };
                lastParam = param;
                $("#mdTypeList").data("kendoDropDownList").value(type);
            } else if (source == "ddl") {
                param = { ProjectName: lastParamChart.ProjectName, DateStr: lastParam.DateStr, Type: dtType };
                $("#mdTypeList").data("kendoDropDownList").value(dtType);
            }
        } else {
            if (source == "chart") {
                param = { ProjectName: project, DateStr: dateStr, Type: type };
                lastParam = param;
                $("#mdTypeList").data("kendoDropDownList").value(type);
            } else if (source == "ddl") {
                param = { ProjectName: project, DateStr: lastParam.DateStr, Type: dtType };
                lastParam = param;
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
    }

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
            { title: "Date", field: "StartDate", template: "#= kendo.toString(moment.utc(StartDate).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #", width: 80 },
            { title: "Turbine", field: "Turbine", width: 90, attributes: { style: "text-align:center;" } },
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

avail.toDetailDTTop = function (e, type) {
    vm.isDashboard(false);
    lgd.isAvailability(false);

    if (type == "Times") {
        avail.detailDTTopTxt("(" + e.category + ") - Frequency");
    } else {
        avail.detailDTTopTxt("(" + e.category + ") - " + type);
    }
    avail.isDetailDTTop(true);

    // get the data and push into the chart    
    avail.DTTopDetail(e.category, type);
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

avail.getDetailDT = function () {
    if (!lgd.isFirst()) {
        dtType = $("#mdTypeList").data("kendoDropDownList").value();

        if (dtType == "") {
            avail.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - All Type");
        } else {
            avail.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - " + dtType);
        }

        avail.toDetailDTLostEnergy(null, true, "ddl");
    }
}

avail.getDetailDTFromProject = function () {
    if (!lgd.isFirst()) {
        projectSelected = $("#projectList").data("kendoDropDownList").value();
        dtType = $("#mdTypeListFleet").data("kendoDropDownList").value();
        avail.detailDTLostEnergyTxt("Lost Energy for Last 12 months - " + projectSelected);
        avail.toDetailDTLostEnergy(null, false, "ddl");
    }
}

avail.backToDownTimeChart = function () {
    var project = $("#projectId").data("kendoDropDownList").value();
    if (project == "Fleet" && !lgd.isFirst() && isFleetDetail == true) {
        vm.isDashboard(false);
        avail.isDetailDTLostEnergy(true);
        isFleetDetail = false;
        avail.toDetailDTLostEnergy(null, false, "button");
        if ($("#projectList").data("kendoDropDownList") != null) {
            $("#projectList").data("kendoDropDownList").value(projectSelected);
        }
    } else {
        avail.LEFleetByDown(false);
        vm.isDashboard(true);
        lgd.isSummary(false);
        lgd.isProduction(false);
        lgd.isAvailability(true);
        avail.isDetailDTLostEnergy(false);
        avail.detailDTLostEnergyTxt("Lost Energy for Last 12 months");
        avail.isDetailDTTop(false);
        avail.detailDTTopTxt("");
    }
}
$( window ).resize(function() {
    avail.refreshChart();
});