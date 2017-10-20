'use strict';


viewModel.summary = {};
var sum = viewModel.summary;

sum.isDetailProd = ko.observable(false);
sum.isDetailProdByProject = ko.observable(false);
sum.isMonthlyProject = ko.observable(false);

sum.isDetailLostEnergy = ko.observable(false);
sum.isDetailLostEnergyLevel2 = ko.observable(false);
sum.isSummaryDetail = ko.observable(true);
sum.isGridDetail = ko.observable(true);

sum.isDetailAvailability = ko.observable(false);

sum.isDetailLostEnergyPlot = ko.observable(false);
sum.isDetailLostEnergyPlotLevel2 = ko.observable(false);

sum.detailProdTxt = ko.observable('');
sum.detailProdMsTxt = ko.observable('');
sum.detailProdProjectTxt = ko.observable('');
sum.detailProdDateTxt = ko.observable('');
sum.titleDetailLevel1 = ko.observable('');

sum.noOfProjects = ko.observable();
sum.noOfProjectsExFleet = ko.observable();
sum.noOfTurbines = ko.observable();
sum.totalMaxCapacity = ko.observable();
sum.currentDown = ko.observable();
sum.twoDaysDown = ko.observable();
sum.dataSource = ko.observable();
sum.dataSourceScada = ko.observable();
sum.dataSourceScadaAvailability = ko.observable();
sum.dataSourceWindDistribution = ko.observable();
sum.windDistData = ko.observable();
sum.availData = ko.observableArray([]);
sum.availSeries = ko.observable([]);
sum.periodSelected = ko.observable('currentmonth');
sum.detailSummary = ko.observableArray([]);
sum.getScadaLastUpdate = ko.observableArray([]);
sum.detailProjectName = ko.observable();
sum.DetailAvailabilityData = ko.observableArray([]);
sum.DetailLostEnergyData = ko.observableArray([]);

sum.periodList = [
    // {"text": "Last 12 Months", "value": "last12months"},
    {"text": "Current Month", "value": "currentmonth"}
]
sum.paramPeriod = [];

sum.paramAvailPeriod = [];

var arrMarkers = [];
var turbines = [];
var map;

vm.dateAsOf(app.currentDateData);

sum.scadaLastUpdate = function(){
    var project = $("#projectId").data("kendoDropDownList").value();
    for(var i=0;i<sum.periodList.length;i++) {
        sum.paramPeriod.push(sum.periodList[i].value);
    }
    for(var i=0;i<lgd.projectAvailList().length;i++) {
        sum.paramAvailPeriod.push(lgd.projectAvailList()[i].value);
    }
    var param = { ProjectName: project, Date: maxdate};


    toolkit.ajaxPost(viewModel.appName + "dashboard/getscadalastupdate", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        sum.getScadaLastUpdate(res.data);

        if (sum.getScadaLastUpdate().length > 0){
            var lastUpdate = sum.getScadaLastUpdate()[0].LastUpdate;
            vm.dateAsOf(lastUpdate);
            
        }
        
    });
}
sum.loadData = function () {
    if (lgd.isSummary()) {
        var project = $("#projectId").data("kendoDropDownList").value();
        var param = { ProjectName: project, Date: maxdate};

        var getscadalastupdate = sum.getScadaLastUpdate();
        if (sum.getScadaLastUpdate().length > 0){
            sum.dataSource(getscadalastupdate[0]);
            sum.noOfProjects(getscadalastupdate[0].NoOfProjects);
            sum.noOfProjectsExFleet(getscadalastupdate[0].NoOfProjects);
            sum.noOfTurbines(getscadalastupdate[0].NoOfTurbines);
            sum.totalMaxCapacity((getscadalastupdate[0].TotalMaxCapacity / 1000) + " MW");
            // sum.currentDown(getscadalastupdate[0].CurrentDown);
            sum.twoDaysDown(getscadalastupdate[0].TwoDaysDown);

            

            // vm.dateAsOf(lastUpdate.addHours(-7));
            sum.ProductionChart(getscadalastupdate[0].Productions);
            sum.CumProduction(getscadalastupdate[0].CummulativeProductions);
           

            sum.SummaryData((project == 'Fleet'? 'gridSummaryDataFleet' : 'gridSummaryData'),project);

        } else {
            var projectStr = $("#projectId").data("kendoDropDownList").text();
            if (projectStr != "Fleet"){
                sum.noOfProjects(1);
                var split = (projectStr.split(" ("))[1].split("|");
                sum.noOfTurbines(split[0]);
                sum.totalMaxCapacity(split[1].slice(0, -1));
            }else{
                sum.noOfProjects($("#projectId").data("kendoDropDownList").dataSource.total()-1);
                sum.noOfTurbines("N/A");
                sum.totalMaxCapacity("N/A");
            }   

            sum.SummaryData((project == 'Fleet'? 'gridSummaryDataFleet' : 'gridSummaryData'),project);

            sum.dataSource(null);
            // sum.currentDown("N/A");
            sum.twoDaysDown("N/A");       
            sum.ProductionChart(null);
            sum.CumProduction(null);
        }

        param = { ProjectName: project, Date: maxdate, ProjectList: sum.paramAvailPeriod};
        var ajax2 = toolkit.ajaxPost(viewModel.appName + "dashboard/getscadasummarybymonth", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            sum.dataSourceScada(res.data["Data"]);
            
            sum.LostEnergy(res.data["Data"]);
            sum.Windiness(res.data["Data"]);
            
            var availabilityData = [];
            var availabilitySeries = [];
            if(project === "Fleet") {
                sum.dataSourceScadaAvailability(res.data["Data"]);
                sum.PLF('chartPLFFleet',res.data["Data"]);
                sum.ProdCurLast('chartCurrLastFleet',res.data["Data"]);
                sum.ProdMonth('chartProdMonthFleet',res.data["Data"]);
                // sum.ProdMonthFleet('chartProdMonthFleet',res.data["Availability"]);
                var availDatas = res.data["Availability"];
                var projectCount = 0;
                for(var key in availDatas){
                    var availData = availDatas[key];
                    var seriesObj = {};
                    for(var i=0;i<availData.length;i++){
                        if(projectCount < 1) {
                            var availObject = {
                                "DateInfo": availData[i].DateInfo
                            }
                            availObject[key] = availData[i].TrueAvail;
                            availabilityData.push(availObject);
                        } else {
                            var availObject = availabilityData[i]
                            availObject[key] = availData[i].TrueAvail;
                            availabilityData[i] = availObject;
                        }
                    }
                    seriesObj["name"] = key;
                    seriesObj["field"] = key;
                    seriesObj["color"] = colorFieldProject[projectCount];
                    seriesObj["missingValues"]= "gap";
                    availabilitySeries.push(seriesObj);
                    projectCount++;
                }

                sum.availData(availabilityData);
                sum.availSeries(availabilitySeries);

                sum.AvailabilityChart(availabilityData, availabilitySeries, "fleet");
            } else {
                sum.PLF('chartPLF',res.data["Data"]);
                sum.ProdCurLast('chartCurrLast',res.data["Data"]);
                sum.ProdMonth('chartProdMonth',res.data["Data"]);
                var availData = res.data["Data"];
                if(res.data["Data"] !== undefined) {
                    var seriesObj = {};
                    for(var i=0;i<availData.length;i++){
                        var availObject = {
                            "DateInfo": availData[i].DateInfo
                        }
                        availObject[project] = availData[i].TrueAvail;
                        availabilityData.push(availObject);
                    }
                    seriesObj["name"] = project;
                    seriesObj["field"] = project;
                    seriesObj["color"] = colorFieldProject[1];
                    seriesObj["missingValues"]= "gap";
                    availabilitySeries.push(seriesObj);

                    sum.availData(availabilityData);
                    sum.availSeries(availabilitySeries);

                    sum.AvailabilityChart(availabilityData, availabilitySeries, "project");
                } else {
                    sum.availData(availData);
                    sum.availSeries(availabilitySeries);
                    sum.AvailabilityChart(availData, availabilitySeries, "project");
                }
            }

            
            
            if (res.data != null ){
                sum.isDetailProd(false);
                sum.isDetailProdByProject(false);
                sum.isDetailLostEnergy(false);
                sum.isDetailLostEnergyLevel2(false);
            }
        });

        var ajax3, lostEnergyReq;

        if (project=="Fleet") {
            param = { ProjectName: project, Date: maxdate, PeriodList: sum.paramPeriod};
            ajax3 = toolkit.ajaxPost(viewModel.appName + "dashboard/getwinddistributionrev", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                sum.WindDistribution(res.data.Data[sum.periodSelected()]);
                sum.dataSourceWindDistribution(res.data.Data[sum.periodSelected()]);
                sum.windDistData(res.data.Data);
            });

            // lostEnergyReq = toolkit.ajaxPost(viewModel.appName + "dashboard/getlostenergy", { ProjectName: project, Date: maxdate}, function (res) {
            //     if (!app.isFine(res)) {
            //         return;
            //     }
            //     // if (project == "Fleet") {
            //     //     avail.DTTurbines();
            //     // } 
            // });
        }

        $.when(sum.indiaMap(project),ajax2, ajax3).done(function(){
            setTimeout(function(){
                if(project == "Fleet"){
                    map.setCenter({
                        lat : 23.334166,
                        lng : 75.037611
                    }); 
                    map.setZoom(5);
                    app.loading(false);
                }else{
                    map.setCenter({
                        lat : turbines[0].coords[0],
                        lng : turbines[0].coords[1]
                    }); 
                    map.setZoom(10);
                    
                }
                lgd.start();
                app.loading(false);
            },1000);
        });

    }

};

sum.SummaryData = function (id,project) {
    var param = {project: project};
    var ajax1 = toolkit.ajaxPost(viewModel.appName + "dashboard/getsummarydata", param, function (result) {
        $('#'+id).html("");
        $('#'+id).kendoGrid({
            height: 155,
            theme: "flat",
            dataSource: {
                // serverPaging: true,
                // serverSorting: true,
                data: result,
                // pageSize: 2,
                schema: {
                    data: function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Total;
                    }
                },
            },
            serverPaging: true,
            serverSorting: true,
            pageable: {
                pageSize: 2,
                input: true, 
            },
            columns: [
                { title: "Project Name", width:100, field: "name", headerAttributes: { style: "text-align:left;" }, attributes: { style: "text-align:center;" } },
                { title: "No. of WTG", width:90, field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Controller Generation<br>(GWh)", width:120, field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "PLF<br>(%)", width:100, field: "plf", format: "{0:n2}", template: "#= kendo.toString(plf*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Lost Energy<br>(MWh)", width:100,field: "lostenergy", template: "#= kendo.toString(lostenergy/1000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Downtime<br>(Hours)", width:120,field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Machine Availability<br>(%)", width:120, field: "machineavail", format: "{0:n2}", template: "#= kendo.toString(machineavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Total Availability<br>(%)", width:120, field: "trueavail", format: "{0:n2}", template: "#= kendo.toString(trueavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Data Availability<br>(%)", width:120, field: "dataavail", format: "{0:n2}", template: "#= kendo.toString(dataavail, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            ]
        });
    });

    $.when(ajax1).done(function(){
        setTimeout(function(){
            var grid = $('#'+id).data("kendoGrid");
            if (project == "Fleet") {
                $("#"+id+" th[data-field=name]").html("Project Name")
                grid.showColumn("noofwtg");
            } else {
                $("#"+id+" th[data-field=name]").html("Turbine Name")
                grid.hideColumn("noofwtg");
            }
            var dataSource = grid.dataSource.data();
            $.each(dataSource, function (i, row) {
                $('tr[data-uid="' + row.uid + '"]').css("border-bottom", "1pt solid black");
            });
            $("#"+id).data("kendoGrid").refresh();
        }, 100);        
    })
}

sum.PLF = function (id,dataSource) {
    $("#"+id).replaceWith('<div id='+id+'></div>');
    $("#"+id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: false,
        },
        chartArea: {
            height: 185,
            background: "transparent",
            padding: 0,
        },
        seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        series: [{
            field: "PLF",
            // opacity : 0.7,
        }],
        seriesColors: colorField,
        valueAxis: {
            // majorUnit: 25,
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
        },
        categoryAxis: {
            field: "DateInfo.MonthDesc",
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
        },
        plotAreaClick: function(e) {
            if(id === "chartPLFFleet") {
                if (e.originalEvent.type === "contextmenu") {
                  // Disable browser context menu
                  e.originalEvent.preventDefault();
                }
                sum.MonthlyProject(e, "plf");
            }
        },
        tooltip: {
            visible: true,
            // template: "PLF comparison for #: category # : #: kendo.toString(value, 'n0')# % ",
            template: "#: category #: #: kendo.toString(value, 'n0')# % ",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        }
    });
}

sum.LostEnergy = function (dataSource) {
    $("#chartLostEnergy").replaceWith('<div id="chartLostEnergy"></div>');
    $("#chartLostEnergy").kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: false,
        },
        chartArea: {
            height: 175,
            background: "transparent",
            padding: 0,
        },
        seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        series: [{
            field: "LostEnergy",
            // opacity : 0.7,
        }],
        seriesColors: colorField,
        seriesClick: function (e) {
            sum.DetailLostEnergy(e);
        },
        plotAreaClick: function(e) {
            if (e.originalEvent.type === "contextmenu") {
              // Disable browser context menu
              e.originalEvent.preventDefault();
            }
            sum.DetailLTPlot("Fleet", e);
        },
        valueAxis: {
            // labels: {
            //     step : 2,
            //     format: "n0"
            // },
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
            field: "DateInfo.MonthDesc",
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
            // template: "Lost Energy for #: category # : #: kendo.toString(value, 'n1')# GWh ",
            template: "#: category #: #: kendo.toString(value, 'n2')# GWh ",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        }
    });
}

sum.Windiness = function (dataSource) {
    $("#chartWindiness").replaceWith('<div id="chartWindiness"></div>');
    $("#chartWindiness").kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
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
        series: [{
            type: "line",
            style: "smooth",
            name: "Avg. Wind Speed",
            field: "AvgWindSpeed",
            // opacity : 0.7,
            markers: {
                visible: false
            }
        }, {
            type: "line",
            style: "smooth",
            name: "Avg. Expected Wind Speed",
            field: "ExpWindSpeed",
            // opacity : 0.7,
            markers: {
                visible: false
            }
        }],
        seriesColors: colorField,
        valueAxis: {
            labels: {
                step: 2,
                format: "n0",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            max: 10,
            majorUnit: 2,
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
            field: "DateInfo.MonthDesc",
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
            format: "{0:n1}",
            shared: true,
            sharedTemplate: kendo.template($("#templateWindiness").html()),
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
sum.UpdateWindDist = function() {
    setTimeout(function() {
        sum.WindDistribution(sum.windDistData()[sum.periodSelected()]);
    }, 300);
}
sum.WindDistribution = function (dataSource) {
    $("#chartWindDistribution").replaceWith('<div id="chartWindDistribution"></div>');
    $("#chartWindDistribution").kendoChart({
        dataSource: {
            data: dataSource,
            group: { field: "Project" },
            sort: { field: "Category", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
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
        series: [{
            type: "line",
            style: "smooth",
            field: "Contribute",
            // opacity : 0.7,
            markers: {
                visible: false,
            }
        }],
        seriesColors: colorFieldProject,
        valueAxis: {
            labels: {
                step: 2,
                format: "{0:p0}",
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
            field: "Category",
            majorGridLines: {
                visible: false
            },
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                step: 2,
                // rotation: 25
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            shared: true,
            sharedTemplate: kendo.template($("#templateDistribution").html()),
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

sum.ProdMonthFleet = function (id, dataSource) {
    var series = [];
    
    var category = [];

    var i = 0;
    $.each(sum.dataSourceScadaAvailability(), function(idx, val){
        var serie = {
            name : idx,
            color: colorFieldProject[i],
            data: []
        }
        $.each(val, function(index, data){
            if(i == 0){
                category.push(data.DateInfo.MonthDesc);
            }
            serie.data.push(data.Production);
        });

        series.push(serie);

        i++;
    });


    $("#"+id).replaceWith('<div id='+id+'></div>');
    $("#"+id).kendoChart({
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 175,
            background: "transparent",
            padding: 0,
            margin: {
                top: -5
            }
        },
        seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            },
            stack: true,
        },
        series: series,
        // seriesColors: colorField,
        seriesClick: function (e) {
            setTimeout(function(){
                lgd.stop();
                sum.DetailProdByProject(e);
            },500);
            
        },
        valueAxes: [{
            name: "Production",
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            labels: {
                step: 2,
                format: "n0",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        }],
        categoryAxis: {
            categories: category,
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
            format: "{0:n1}",
            shared: false,
            template: "<strong>#= series.name # </strong>: #= kendo.toString(value , 'n2') # Gwh",
            // sharedTemplate: kendo.template($("#templateProdMonth").html()),
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        }
    });
}

sum.ProdMonth = function (id, dataSource) {
    $("#"+id).replaceWith('<div id='+id+'></div>');
    $("#"+id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height: 175,
            background: "transparent",
            padding: 0,
            margin: {
                top: -5
            }
        },
        seriesDefaults: {
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        series: [
        // {
        //     name: "Budget",
        //     field: "Budget",
        //     // opacity : 0.7,
        //     color: "#21c4af",
        // }, 
        {
            name: "Production",
            field: "Production",
            // opacity : 0.7,
            color: "#ff9933",
        }],
        // seriesColors: colorField,
        seriesClick: function (e) {
            sum.DetailProd(e);
        },
        plotAreaClick: function(e) {
            if (e.originalEvent.type === "contextmenu") {
              // Disable browser context menu
              e.originalEvent.preventDefault();
            }
            sum.MonthlyProject(e, "production");
        },
        valueAxes: [{
            name: "Production",
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            labels: {
                step: 2,
                format: "n0",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        }],
        categoryAxis: {
            field: "DateInfo.MonthDesc",
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
            format: "{0:n1}",
            shared: true,
            sharedTemplate: kendo.template($("#templateProdMonth").html()),
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        }
    });
}
// sum.UpdateAvailability = function() {
//     setTimeout(function() {
//         sum.AvailabilityChart(sum.availabilityData()[lgd.projectAvailSelected()]);
//     }, 300);
// }
sum.AvailabilityChart = function (dataSource, dataSeries, tipe) {
    $("#chartAbility").replaceWith('<div id="chartAbility"></div>');
    $("#chartAbility").kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
            type: "line",
            style: "smooth",
            // area: {
            //     line: {
            //         style: "smooth"
            //     }
            // }
            markers: {
                visible: false,
            }
        },
        // series: [{
        //     name: "Tejuva",
        //     field: "ScadaAvail",
        //     // opacity : 0.5,
        //     color: "#21c4af"
        // }, {
        //     name: "Lahori",
        //     field: "TrueAvail",
        //     // opacity : 0.5,
        //     color: "#ff880e",
        // }],
        series: dataSeries,
        seriesColors: colorFieldProject,
        valueAxis: {
            max: 100,
            majorUnit: 25,
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            labels: {
                // format: "{0}%"
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
        },
        categoryAxis: {
            field: "DateInfo.MonthDesc",
            majorGridLines: {
                visible: false
            },
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        plotAreaClick: function(e) {
            if(tipe === "fleet") {
                if (e.originalEvent.type === "contextmenu") {
                  // Disable browser context menu
                  e.originalEvent.preventDefault();
                }
                // sum.DetailAvailability("Fleet", e, "availability");
                sum.MonthlyProject(e, "availability");
            }
        },
        tooltip: {
            visible: true,
            format: "{0:n1}",
            sharedTemplate: kendo.template($("#templateAvail").html()),
            shared: true,
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

sum.ProdCurLast = function (id,dataSource) {

    $("#"+id).replaceWith('<div id='+id+'></div>');
    $("#"+id).kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
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
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        seriesColors: colorField,
        series: [{
            type: "column",
            style: "smooth",
            name: "Current",
            field: "Production",
            // opacity : 0.7,
            axis: "production",
            // color: "#21c4af",
        }, {
            type: "column",
            style: "smooth",
            name: "Last",
            field: "ProductionLastYear",
            // opacity : 0.7,
            axis: "production",
            // color: "#ff880e",
        }, {
            type: "line",
            style: "smooth",
            name: "Variance(%)",
            field: "Variance",
            axis: "variance",
            // color: "#ff7663",
            markers: {
                visible: false
            }
        }],
        valueAxes: [{
            line: {
                visible: false
            },
            labels: {
                step: 2,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "production",
        }, {
            line: {
                visible: false
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            labels: {
                format: "{0}%",
                visible: false,
            },
            name: "variance",
        }],
        categoryAxis: {
            field: "DateInfo.MonthDesc",
            axisCrossingValues: [0, 1000],
            majorGridLines: {
                visible: false
            },
            labels: {
                // step: 2,
                template: '#=  value.substring(0,3) #'
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            format: "{0:n2}",
            // template: "#= value #",
            shared: true,
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

// INDIA MAPS

sum.setMarkers = function(map, turbineInfos,project) {
    turbineInfos.forEach(function (obj, idx) {
        
        var imgUrl ="../res/img/turb-"+obj.status+".png";


        var marker = new google.maps.Marker({
            position: new google.maps.LatLng(obj.coords[0], obj.coords[1]),
            map: map,
            title: obj.name,
            icon: {
                url: imgUrl, // url
                scaledSize: new google.maps.Size(20, 20), // scaled size
            }
        });

        arrMarkers.push(marker);

        var infowindow = new google.maps.InfoWindow({
            content: ""
        });

        google.maps.event.addListener(marker, 'click', function () {
            var project = $("#projectId").data("kendoDropDownList").value();
            if(project == "Fleet"){
                // sum.ToMonitoringProject(obj.name);
                setTimeout(function(){
                    $("#projectId").data('kendoDropDownList').value(obj.name);
                    lgd.LoadData();
                }, 200);
            }else{
                sum.ToMonitoringIndividual(project, obj.value);
            }
        });
    });

}

sum.initialize = function() {
    $("#india-map").html("");
    var mapOptions = {
        types: ['(region)'],
        componentRestrictions: {country: "in"},
        // center: (projectname == 'Fleet' ? new google.maps.LatLng(22.460533, 79.650879) : center),
        // center: center,
        center: new google.maps.LatLng(23.334166, 75.037611) ,
        // zoom: (project == 'Fleet' ? 4 : 10),
        zoom: 5,
        styles: [
          {
            "featureType": "administrative.country",
            "elementType": "geometry",
            "stylers": [
              {
                "color": "#ff2631"
              },
              {
                "weight": 2
              }
            ]
          }
        ],
        mapTypeId: google.maps.MapTypeId.HYBRID,
        mapTypeControl: true,
        mapTypeControlOptions: {
            position: google.maps.ControlPosition.LEFT_BOTTOM
        },
        zoomControl: true,
        zoomControlOptions: {
            position: google.maps.ControlPosition.RIGHT_BOTTOM
        },
        scaleControl: true,
        streetViewControl: true,
        streetViewControlOptions: {
            position: google.maps.ControlPosition.RIGHT_BOTTOM
        },
        fullscreenControl: true,
        fullscreenControlOptions: {
            position: google.maps.ControlPosition.RIGHT_BOTTOM
        },
    }
    map = new google.maps.Map(document.getElementById("india-map"), mapOptions);

    var project =  $("#projectId").data("kendoDropDownList").value();

    $.when(sum.indiaMap(project)).done(function(){
        setTimeout(function(){
             sum.setMarkers(map, turbines,project);
        },500);
    })
   
}

sum.removeMarkers = function(){
    var i;
    for(i=0;i<arrMarkers.length;i++){
        arrMarkers[i].setMap(null);
    }
    arrMarkers = [];

}
sum.indiaMap = function (project) {
    var param = { projectname: project }
    toolkit.ajaxPost(viewModel.appName + "dashboard/getmapdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        sum.removeMarkers();
        var jsonObj = res.data.resultMap;

        if(project === "Fleet") {
            sum.currentDown(res.data.totalDownFleet);
            avail.DTTurbines(res.data.turbineDownList);
        } else {
            sum.currentDown(res.data.downPerProject[project]);
        }

        turbines =[];//Erasing the beaches array

        turbines = jsonObj;
        //Adding the new ones
        // for(i=0;i < jsonObj.turbines.length; i++) {
        //     turbines.push(jsonObj.turbines[i]);
        // }

        //Adding them to the map
        sum.setMarkers(map, turbines,project);

    })
}

sum.ToMonitoringProject = function(project) {
    setTimeout(function(){
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');
        document.cookie = "project="+project.split("(")[0].trim()+";expires="+ newDateObj;
        window.location = viewModel.appName + "page/monitoringbyproject";
    },300);
}

sum.ToMonitoringIndividual = function(project, turbine) {
    setTimeout(function(){
        app.loading(true);
        var oldDateObj = new Date();
        var newDateObj = moment(oldDateObj).add(3, 'm');

        document.cookie = "projectname="+project+";expires="+ newDateObj;
        document.cookie = "turbine="+turbine+";expires="+ newDateObj;

        if(document.cookie.indexOf("projectname=") >= 0 && document.cookie.indexOf("turbine=") >= 0) {
            window.location = viewModel.appName + "page/monitoringbyturbine";
        } else {
            app.loading(false);
        }
    },1500);
}

sum.ProductionChart = function (dataSource) {
    var dataFormat = "n2";
    if (dataSource != null){
        if (dataSource.length > 0) {
            var totalPotential = 0;
            for (var i = 0; i < dataSource.length; i++) {
                totalPotential += dataSource[i].PotentialKwh;
            }
            if (totalPotential > 10) {
                dataFormat = "n0";
            }
        }
    }

    $("#chartProduction").replaceWith('<div id="chartProduction"></div>');
    $("#chartProduction").kendoChart({
        dataSource: {
            data: dataSource,
            // sort: { field: "Hour", dir: "asc"}
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
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
            type: "area",
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        series: [{
            name: "Potential Power",
            field: "PotentialKwh",
            // opacity : 0.5,
            color: "#21c4af",
        }, {
            name: "Production",
            field: "EnergyKwh",
            // opacity : 0.5,
            color: "#ff9933"
        }],
        valueAxis: {
            labels: {
                step: 2,
                format: dataFormat,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            line: {
                visible: false
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            }
            //  majorUnit: 3,
            //     line: {
            //         visible: false
            //     },
            //     axisCrossingValue: -10,
            //     majorGridLines: {
            //         visible: true,
            //         color: "#eee",
            //         width: 0.8,
            //     },
        },
        categoryAxis: {
            field: "Hour",
            majorGridLines: {
                visible: false
            },
            labels: {
                step: 2,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                // template : "#: Number(kendo.toString(kendo.parseDate(value), 'HH')) #"
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            sharedTemplate: kendo.template($("#templateProd").html()),
            shared: true,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
            format: dataFormat,
        },
    });
}

sum.CumProduction = function (dataSource) {
    $("#chartCumProduction").replaceWith('<div id="chartCumProduction"></div>');
    $("#chartCumProduction").kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DayNo", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
            type: "line",
            style: "smooth",
            // opacity : 0.7,
            markers: {
                visible: false,
            }
        },
        series: [
        {
            name: "Budget",
            field: "CumBudget",
            // opacity : 0.5,
            color: "#21c4af",
        }, 
        {
            name: "Production",
            field: "CumProduction",
            // opacity : 0.5,
            color: "#ff9933"
        }],
        // seriesColors: colorField,
        valueAxis: {
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            labels: {
                step: 2,
                format: "n0",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        categoryAxis: {
            field: "DateId",
            majorGridLines: {
                visible: false
            },
            labels: {
                step: 3,
                template: "#: Number(kendo.toString(kendo.parseDate(value), 'dd'))#",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            sharedTemplate: kendo.template($("#templateCum").html()),
            shared: true,
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


sum.DetailProd = function (e) {
    app.loading(true);

    $("#chartDetailProduction").html("");
    $("#chartDetailLostEnergy").html("");
    $("#gridDetailProduction").html("");

    var bulan = e.category;
    sum.detailProdTxt(bulan);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(true);
    sum.isDetailProdByProject(false);
    sum.isDetailLostEnergy(false);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { 'project': project, 'date': bulan };
    var dataSource;
    var measurement;

    var reqDetail = toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprodlevel1", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        dataSource = res.data;
        measurement = " (" + dataSource[0].measurement + ") ";
        sum.detailProdMsTxt(measurement);
    });
    
    $.when(reqDetail).done(function() {
        $("#chartDetailProduction").kendoChart({
            theme: "material",
            dataSource: {
                data: dataSource
            },
            title: {
                text: ""
            },
            legend: {
                position: "top",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            chartArea: {
                height: 200,
                padding: 0,
                margin: 0,
            },
            seriesDefaults: {
                type: "column",
                area: {
                    line: {
                        style: "smooth"
                    }
                }
            },
            series: [{
                name: "Production" + measurement,
                field: "production",
                gap: 3,
                // opacity : 0.7,
            }],
            seriesColors: colorField,
            seriesClick: function (e) {
                sum.DetailProdByProject(e);
            },
            valueAxis: {
                // majorUnit : 2000,
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
                },
            },
            categoryAxis: {
                field: "project",
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
                template: "#: series.name # : #:  kendo.toString(value, 'n2') #",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });

        $("#chartDetailLostEnergy").kendoChart({
            theme: "material",
            dataSource: {
                data: dataSource
            },
            title: {
                text: ""
            },
            legend: {
                position: "top",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            chartArea: {
                height: 200,
                padding: 0,
                margin: 0,
            },
            seriesDefaults: {
                type: "column",
                area: {
                    line: {
                        style: "smooth"
                    }
                }
            },
            series: [{
                name: "Lost Energy" + measurement,
                field: "lostenergy",
                gap: 3,
                color: colorField[1],
                // opacity : 0.7,
            }],
            // seriesColors: colorField[1],
            seriesClick: function (e) {
                sum.DetailProdByProject(e);
            },
            valueAxis: {
                // majorUnit : 2000,
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
                },
            },
            categoryAxis: {
                field: "project",
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
                template: "#: series.name # : #:  kendo.toString(value, 'n2') #",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });
        $("#gridDetailProduction").kendoGrid({
            theme: "flat",
            pageable: {
                pageSize: 5,
                input: true, 
            },
            columns: [
                { title: "Project Name", field: "project", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "No. of WTG", field: "wtg", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
                { title: "Production <br>" + measurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
                { title: "Lost Energy <br>" + measurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            ],
            dataSource: {
                data: dataSource,
                pageSize: 5
            },
            dataBound: function(){
                app.loading(false);
            }
        });
    })
}

// show monthly project, level 2 of generation summary / production monthly from dashboard
sum.MonthlyProject = function (e, tipe) {
    app.loading(true);

    $('.monthlyProjectChart').html("");

    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isMonthlyProject(true);

    var param = {
        Projects: ["Tejuva", "Lahori", "Amba"],
    };
    var dataSource;
    var dataRequest = toolkit.ajaxPost(viewModel.appName + "dashboard/getmonthlyproject", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        dataSource = res.data;
    });
    var dataSeries = [{
        name: "Production (MWh)",
        field: "production",
        axis: "production"
    }, {
        name: "Lost Energy (MWh) ",
        field: "lostenergy",
        axis: "lostenergy"
    }];

    var valueAxesData = [{
        name: "production",
        title: { 
            margin: {
                right: 0
            },
            text: "Production (MWh)",font: "10px"
        },
        line: {
            visible: false
        },
        labels:{
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
        },
        axisCrossingValue: -10,
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        
    }, {
        name: "lostenergy",
        title: { 
            margin: {
                left: 0
            },
            text: "Lost Energy (MWh)",font: "10px"},
        line: {
            visible: false
        },
        labels:{
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
        },
        axisCrossingValue: -10,
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
    }];
    sum.titleDetailLevel1('Controller Generation (GWh) - Last 12 Months By Project');

    if(tipe == "plf") {
        sum.titleDetailLevel1('PLF (%) - Last 12 Months By Project');
        dataSeries = [{
            name: "PLF ",
            field: "plf",
            axis: "plf"
        }];
        valueAxesData = [{
            name: "plf",
            title: { 
                margin: {
                    left: 0
                },
                text: "PLF (%)",font: "11px"},
            line: {
                visible: false
            },
            labels:{
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
        }];
    }

    if(tipe == "availability"){
        sum.titleDetailLevel1('Availability (%) - Last 12 Months By Project');
        dataSeries = [{
            name: "Availability",
            field: "availability",
            axis: "availability"
        }];
        valueAxesData = [{
            name: "availability",
            title: { 
                margin: {
                    left: 0
                },
                text: "Availability (%)",font: "11px"},
            line: {
                visible: false
            },
            labels:{
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
        }];
    }

    var chartProperties = {
        theme: "flat",
        dataSource: {
            data: null
        },
        title: {
            text: ""
        },
        legend: {
            position: "top",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            padding: 10,
            margin: 5,
            height: 200,
        },
        seriesDefaults: {
            type: "column",
            area: {
                line: {
                    style: "smooth"
                }
            }
        },
        series: dataSeries,
        seriesColors: colorField,
        valueAxes: valueAxesData,
        // valueAxis: {
        //     line: {
        //         visible: false
        //     },
        //     labels:{
        //         font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
        //     },
        //     axisCrossingValue: -10,
        //     majorGridLines: {
        //         visible: true,
        //         color: "#eee",
        //         width: 0.8,
        //     },
        // },
        categoryAxis: {
            field: "monthdesc",
            majorGridLines: {
                visible: false
            },
            majorTickType: "none",
            labels: {
                template: '#=  value.substring(0,3) #',
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValues: [0, 1000],
        },
        tooltip: {
            visible: true,
            shared: true,
            template: "#= kendo.toString(value, 'n2') #",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        },
        seriesClick: function (e) {
            sum.DetailProdByProject(e);
        },
    };
    
    $.when(dataRequest).done(function() {
        var charts = $('.monthlyProjectChart');
        if(charts.length > 0) {
            $.each(charts, function(idx, elm){
                var project = $(elm).attr('data-project');
                if(project!='Fleet') {
                    chartProperties.dataSource = eval("dataSource."+ project);
                    chartProperties.seriesClick = function(e) {
                        sum.DetailProdByProjectDetail(project, e, tipe);
                    };
                    $(elm).kendoChart(chartProperties);
                }
            });
        }
        app.loading(false);
    });
} 

sum.DetailProdByProject = function (e) {
    $('#btn-back-prod-summary').removeAttr('onclick');

    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailProdByProject(true);
    $("#chartDetailProdByProject").html("");
    $("#gridDetailProdByProject").html("");

    // var project = e.series.name;
    var param = { 'project': e.category, 'date': sum.detailProdTxt() };

    sum.detailProdProjectTxt(e.category);
    sum.detailProdDateTxt(sum.detailProdTxt());

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data[0];

        var measurement = " (" + dataSource.measurement + ") ";

        sum.detailSummary(dataSource);
        sum.detailProdMsTxt(measurement);
        $("#chartDetailProdByProject").kendoChart({
            theme: "material",
            dataSource: {
                data: dataSource.detail
            },
            title: {
                text: ""
            },
            legend: {
                position: "top",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            chartArea: {
                padding: 0,
                margin: 0,
                height: 350,
            },
            seriesDefaults: {
                type: "column",
                area: {
                    line: {
                        style: "smooth"
                    }
                }
            },
            series: [{
                name: "Production " +measurement,
                field: "production",
                axis: "production"
                // opacity : 0.7,
            }, {
                name: "Lost Energy " +measurement,
                field: "lostenergy",
                axis: "lostenergy"
                // opacity : 0.7,
            }],
            seriesColors: colorField,
            valueAxes: [{
                name: "production",
                title: { text: "Production " +measurement ,font: "11px"},
                
            }, {
                name: "lostenergy",
                title: { text: "Lost Energy " + measurement,font: "11px"},
            }],
            valueAxis: {
                line: {
                    visible: false
                },
                labels:{
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            },
            categoryAxis: {
                field: "turbine",
                majorGridLines: {
                    visible: false
                },
                majorTickType: "none",
                labels: {
                    rotation: 45,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValues: [0, 1000],
            },
            tooltip: {
                visible: true,
                shared: true,
                template: "#= kendo.toString(value, 'n2') #",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });

        
        $("#gridDetailProdByProject").kendoGrid({
            theme: "flat",
            height: 200,
            pageable: {
                pageSize: 10,
                input: true, 
            },
            columns: [
                { title: "Turbine Name", field: "turbine", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Production<br>" + measurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Lost Energy<br>" + measurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Downtime<br>(Hours)", field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            ],
            dataSource: {
                data: dataSource.detail,
                sort: { field: "turbine", dir: 'asc' },
                pageSize: 10
            },
            dataBound : function(){
                setTimeout(function(){
                    app.loading(false);
                },200);
            }
        });

    });


}

sum.DetailProdByProjectDetail = function (project, e, tipe) {
    $('#btn-back-prod-summary').attr('onclick', 'sum.backToMonthlyProject()');

    var bulan = e.category;
    sum.detailProdTxt(bulan);

    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isMonthlyProject(false);
    sum.isDetailProd(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailProdByProject(true);
    sum.isSummaryDetail(true);
    sum.isGridDetail(true);

    if(tipe === "plf" || tipe === "availability") {
        sum.isSummaryDetail(false);
    }

    if(tipe == "availability"){
        sum.isGridDetail(false);
    }

    $("#chartDetailProdByProject").html("");
    $("#gridDetailProdByProject").html("");

    // var project = e.series.name;
    var param = { 'project': project, 'date': sum.detailProdTxt() };

    sum.detailProdProjectTxt(project);
    sum.detailProdDateTxt(sum.detailProdTxt());

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data[0];

        var measurement = " (" + dataSource.measurement + ") ";
        var gridMeasurement = " (" + dataSource.measurement + ") ";
        var seriesData = [{
            name: "Production " +measurement,
            field: "production",
            axis: "production"
        }, {
            name: "Lost Energy " +measurement,
            field: "lostenergy",
            axis: "lostenergy"
        }];
        var isHidden = true;
        var valueAxesData = [{
            name: "production",
            title: { text: "Production " +measurement ,font: "11px"},
            
        }, {
            name: "lostenergy",
            title: { text: "Lost Energy " + measurement,font: "11px"},
        }];
        if(tipe === "plf") {
            measurement = " (%) ";
            seriesData = [{
                name: "PLF " +measurement,
                field: "plf",
                axis: "plf"
            }]
            valueAxesData = [{
                name: "plf",
                title: { text: "PLF " +measurement ,font: "11px"},
                line: {
                    visible: false
                },
                labels:{
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            }];
            isHidden = false;
        }

        if(tipe === "availability") {
            measurement = " (%) ";
            seriesData = [{
                name: "Availability " +measurement,
                field: "trueavail",
                axis: "trueavail"
            }]
            valueAxesData = [{
                name: "trueavail",
                title: { text: "Availability " +measurement ,font: "11px"},
                line: {
                    visible: false
                },
                labels:{
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            }];
            isHidden = false;
        }

        sum.detailSummary(dataSource);
        sum.detailProdMsTxt(measurement);

        $("#chartDetailProdByProject").kendoChart({
            theme: "material",
            dataSource: {
                data: dataSource.detail
            },
            title: {
                text: ""
            },
            legend: {
                position: "top",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            chartArea: {
                padding: 0,
                margin: 0,
                height: 350,
            },
            seriesDefaults: {
                type: "column",
                area: {
                    line: {
                        style: "smooth"
                    }
                }
            },
            series: seriesData,
            seriesColors: colorField,
            valueAxes: valueAxesData,
            valueAxis: {
                line: {
                    visible: false
                },
                labels:{
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -10,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
            },
            categoryAxis: {
                field: "turbine",
                majorGridLines: {
                    visible: false
                },
                majorTickType: "none",
                labels: {
                    rotation: 45,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValues: [0, 1000],
            },
            tooltip: {
                visible: true,
                shared: true,
                template: "#= kendo.toString(value, 'n2') #",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });

        
        $("#gridDetailProdByProject").kendoGrid({
            theme: "flat",
            height: 200,
            pageable: {
                pageSize: 10,
                input: true, 
            },
            columns: [
                { title: "Turbine Name", field: "turbine", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Production<br>" + gridMeasurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Lost Energy<br>" + gridMeasurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { hidden: !isHidden, title: "Downtime<br>(Hours)", field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { hidden: isHidden, title: "PLF<br>(%)", field: "plf", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            ],
            dataSource: {
                data: dataSource.detail,
                sort: { field: "turbine", dir: 'asc' },
                pageSize: 10
            },
            dataBound : function(){
                setTimeout(function(){
                    app.loading(false);
                },200);
            }
        });

    });


}

sum.DetailAvailability = function (project, e, tipe) {
    $('#btn-back-availability').attr('onclick', 'sum.backToDashboard()');
    $('.detailAvailability').html("");

    var bulan = e.category;
    sum.detailProdTxt(bulan);

    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isMonthlyProject(false);
    sum.isDetailProd(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailProdByProject(false);
    sum.isSummaryDetail(false);
    sum.isGridDetail(false);
    sum.isDetailAvailability(true);

    var param = { 'project': project, 'date': sum.detailProdTxt() };

    sum.detailProdProjectTxt(project);
    sum.detailProdDateTxt(sum.detailProdTxt());

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data;
        sum.DetailAvailabilityData(dataSource);

        $.each(dataSource, function(key, val){
            var seriesData = [{
                name: "Availability (%)",
                field: "trueavail",
                axis: "availability",
                color: customColorProject[val.project],
            }];

            var isHidden = true;

            $("#chartDetailAvail-"+val.project).kendoChart({
                theme: "material",
                dataSource: {
                    data: val.detail
                },
                title: { 
                    text: val.project,
                    font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                legend: {
                    visible: false,
                    position: "top",
                    labels: {
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    }
                },
                chartArea: {
                    padding: 10,
                    margin: 10,
                    height: 250,
                },
                seriesDefaults: {
                    type: "column",
                    area: {
                        line: {
                            style: "smooth"
                        }
                    }
                },
                series: seriesData,
                valueAxis: {
                    name: "availability",
                    title: { text: "Availability  (%)",font: "11px"},
                    line: {
                        visible: false
                    },
                    labels:{
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    },
                    axisCrossingValue: -10,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                },
                categoryAxis: {
                    field: "turbine",
                    majorGridLines: {
                        visible: false
                    },
                    majorTickType: "none",
                    labels: {
                        rotation: 45,
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    },
                    axisCrossingValues: [0, 1000],
                },
                tooltip: {
                    visible: true,
                    shared: true,
                    template: "#= kendo.toString(value, 'n2') #",
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

        app.loading(false);
    });
}

sum.DetailLostEnergy = function(e){
    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(false);
    sum.isDetailLostEnergy(true);
    sum.detailProdTxt(e.category);

    var project = $("#projectId").data("kendoDropDownList").value();

    var param = { ProjectName: project, Date: maxdate, DateStr : e.category,Type: "project"};
    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimetop", param, function (res) {

        if (!app.isFine(res)) {
            return;
        }

        var dataSource = res.data.project;

        $("#chartDetailLost").kendoChart({
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
            },
            series: [
            {
                field: "GridDown",
                name: "Grid Down"
            }, 
            {
                field: "MachineDown",
                name: "Machine Down"
            }, 
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
                    visible: true, 
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
                sum.DetailLostEnergyLevel2(e);
            },
            dataBound : function(){
                setTimeout(function(){
                    app.loading(false);
                },200);
            }
        });
    });

    // $("#chartDetailLost").data("kendoChart").refresh();

}
sum.DetailLostEnergyLevel2 = function(e){
    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailLostEnergyLevel2(true);
    sum.detailProjectName(e.category);

    var param = { ProjectName: e.category, DateStr: sum.detailProdTxt(), Type: "", IsDetail:true};
    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimefleetbydown", param, function (res) {

        if (!app.isFine(res)) {
            return;
        }

        var dataSource = res.data.lostenergy;

        $("#chartDetailLostLevel2").kendoChart({
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
                    text: "MWh",
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    visible: true, 
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
            dataBound: function(){
                setTimeout(function(){
                    app.loading(false);
                },200);
            }
        });
    });
}


sum.DetailLTPlot = function(project, e){

    $('.detaillostenergy').html("");

    var bulan = e.category;
    sum.detailProdTxt(bulan);

    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isMonthlyProject(false);
    sum.isDetailProd(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailProdByProject(false);
    sum.isSummaryDetail(false);
    sum.isGridDetail(false);
    sum.isDetailAvailability(false);
    sum.isDetailLostEnergyPlotLevel2(false);
    sum.isDetailLostEnergyPlot(true);


    var param = { 'ProjectName': project, 'Date': new Date()};

    sum.detailProdProjectTxt(project);
    sum.detailProdDateTxt(sum.detailProdTxt());

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetaillosslevel1", param, function (res) {
        
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data;

        setTimeout(function(){
            sum.DetailLostEnergyData(dataSource.datachart);
            $.each(dataSource.datachart, function(key, val){
                $("#chartLostEnergyByMonth-"+key).kendoChart({
                    dataSource: {
                        data: dataSource.datachart[key],
                    },
                    theme: "flat",
                    title: { 
                        text: key,
                        font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                    },
                    seriesDefaults: {
                        type: "column",
                        stack: true
                    },
                    series: [{
                        name: "Grid Down",
                        type: "column",
                        field: "GridDown",
                    },{
                        name: "Machine Down",
                        type: "column",
                        field: "MachineDown",
                    },{
                        name: "Unknown",
                        type: "column",
                        field: "Unknown",
                    }],
                    seriesColors: colorFieldProject,
                    valueAxis: {
                        title: {
                            text: "MWh",
                            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            visible: true, 
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
                            template: '#=  value.substring(0,3) #',
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        },
                        majorTickType: "none"
                    },
                    tooltip: {
                        visible: true,
                        format: "{0:n1}",
                        background: "rgb(255,255,255, 0.9)",
                        // template: "#= series.name # :  #= series.value #",
                        shared: true,
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    },
                    seriesClick: function (e) {
                        sum.DetailLTPlotLevel2(key,e);
                    },
                });
            });

            $.each(dataSource.datapie, function(key, val){

                var dataPie = dataSource.datapie[key].sort(function (a, b) {
                    return a.name.localeCompare( b.name );
                });

                $("#chartLostEnergyByType-"+key).kendoChart({
                    theme: "flat",
                    title: {
                        text: ""
                    },
                    chartArea: {
                        width: 300,
                        height: 200
                    },
                    legend: {
                        position: "right",
                        visible: true,
                        labels: {
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        }
                    },
                    dataSource: {
                        data: dataPie
                    },
                    series: [{
                        type: "pie",
                        field: "value",
                        categoryField: "name",
                    }],
                    seriesColors: colorFieldProject,
                    tooltip: {
                        visible: true,
                        format: "{0:n1}",
                        background: "rgb(255,255,255, 0.9)",
                        template: "${ category } : ${ kendo.toString(value, 'n2') }",
                        shared: true,
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    },
                });
            });


        },500);

        app.loading(false);
    });

}
sum.DetailLTPlotLevel2 = function(project, e){

    $('.detaillostenergylvl2').html("");

    sum.detailProdTxt(e.category);

    app.loading(true);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isMonthlyProject(false);
    sum.isDetailProd(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailProdByProject(false);
    sum.isSummaryDetail(false);
    sum.isGridDetail(false);
    sum.isDetailAvailability(false);
    sum.isDetailLostEnergyPlotLevel2(true);
    sum.isDetailLostEnergyPlot(false);


    var param = { 'project': project, 'date': e.category};

    sum.detailProdProjectTxt(project);
    sum.detailProdDateTxt(sum.detailProdTxt());

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetaillosslevel2", param, function (res) {
        
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data;

        setTimeout(function(){
            $(".detaillostenergylvl2").kendoChart({
                dataSource: {
                    data: dataSource,
                },
                theme: "flat",
                title: { 
                    text: "",
                    font: 'bold 12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                legend: {
                    position: "top",
                    visible: true,
                    labels: {
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    }
                },
                chartArea: {
                    height: 350,
                    background: "transparent",
                    padding: 0,
                },
                seriesDefaults: {
                    type: "column",
                    stack: true
                },
                series: [{
                    name: "Grid Down",
                    type: "column",
                    field: "GridDown",
                },{
                    name: "Machine Down",
                    type: "column",
                    field: "MachineDown",
                },{
                    name: "Unknown",
                    type: "column",
                    field: "Unknown",
                }],
                seriesColors: colorFieldProject,
                valueAxis: {
                    // min: 0,
                    title: {
                        text: "MWh",
                        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        visible: true, 
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
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        rotation: "auto",
                    },
                    majorTickType: "none"
                },
                tooltip: {
                    visible: true,
                    format: "{0:n1}",
                    background: "rgb(255,255,255, 0.9)",
                    // template: "#= series.name # :  #= series.value #",
                    shared: true,
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                },
            });
        },500);

        app.loading(false);
    });

}
sum.backToDashboard = function () {
    lgd.start();
    vm.isDashboard(true);
    lgd.isSummary(true);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(false);
    sum.isDetailLostEnergy(false);
    sum.isDetailLostEnergyLevel2(false);
    sum.isMonthlyProject(false);
    sum.isDetailAvailability(false);
    sum.isDetailLostEnergyPlot(false);
    sum.isDetailLostEnergyPlotLevel2(false);
}

sum.toDetailLostEnergyLvl1 = function(){
    sum.isDetailLostEnergyPlot(true);
    sum.isDetailLostEnergyPlotLevel2(false);
}

sum.toDetailProduction = function () {
    sum.isDetailProd(true);
    sum.isDetailProdByProject(false);
}

sum.backToMonthlyProject = function() {
    sum.isMonthlyProject(true);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(false);
}

sum.backToLostEnegery = function(){
    sum.isDetailLostEnergy(true);
    sum.isDetailLostEnergyLevel2(false);
}