'use strict';


viewModel.summary = {};
var sum = viewModel.summary;

sum.isDetailProd = ko.observable(false);
sum.isDetailProdByProject = ko.observable(false);

sum.detailProdTxt = ko.observable('');
sum.detailProdMsTxt = ko.observable('');
sum.detailProdProjectTxt = ko.observable('');
sum.detailProdDateTxt = ko.observable('');

sum.noOfProjects = ko.observable();
sum.noOfProjectsExFleet = ko.observable();
sum.noOfTurbines = ko.observable();
sum.totalMaxCapacity = ko.observable();
sum.currentDown = ko.observable();
sum.twoDaysDown = ko.observable();
sum.dataSource = ko.observable();
sum.dataSourceScada = ko.observable();
sum.dataSourceWindDistribution = ko.observable();
sum.windDistData = ko.observable();
sum.availData = ko.observableArray([]);
sum.availSeries = ko.observable([]);
sum.periodSelected = ko.observable('currentmonth');
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
sum.loadData = function () {

    if (lgd.isSummary()) {
        var project = $("#projectId").data("kendoDropDownList").value();
        for(var i=0;i<sum.periodList.length;i++) {
            sum.paramPeriod.push(sum.periodList[i].value);
        }
        for(var i=0;i<lgd.projectAvailList().length;i++) {
            sum.paramAvailPeriod.push(lgd.projectAvailList()[i].value);
        }
        var param = { ProjectName: project, Date: maxdate};

        var ajax1 = toolkit.ajaxPost(viewModel.appName + "dashboard/getscadalastupdate", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }


            if (res.data.length > 0){
                sum.dataSource(res.data[0]);
                sum.noOfProjects(res.data[0].NoOfProjects);
                sum.noOfProjectsExFleet(res.data[0].NoOfProjects);
                sum.noOfTurbines(res.data[0].NoOfTurbines);
                sum.totalMaxCapacity((res.data[0].TotalMaxCapacity / 1000) + " MW");
                sum.currentDown(res.data[0].CurrentDown);
                sum.twoDaysDown(res.data[0].TwoDaysDown);

                var lastUpdate = res.data[0].LastUpdate;

                // vm.dateAsOf(lastUpdate.addHours(-7));
                sum.ProductionChart(res.data[0].Productions);
                sum.CumProduction(res.data[0].CummulativeProductions);
                vm.dateAsOf(lastUpdate);

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
                sum.currentDown("N/A");
                sum.twoDaysDown("N/A");       
                sum.ProductionChart(null);
                sum.CumProduction(null);
            }
            
        });

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
                sum.PLF('chartPLFFleet',res.data["Data"]);
                sum.ProdCurLast('chartCurrLastFleet',res.data["Data"]);
                sum.ProdMonth('chartProdMonthFleet',res.data["Data"]);
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
                    availabilitySeries.push(seriesObj);
                    projectCount++;
                }

                sum.availData(availabilityData);
                sum.availSeries(availabilitySeries);

                sum.AvailabilityChart(availabilityData, availabilitySeries);

                // sum.AvailabilityChart(res.data["Availability"][lgd.projectAvailSelected()]);
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
                    availabilitySeries.push(seriesObj);

                    sum.availData(availabilityData);
                    sum.availSeries(availabilitySeries);

                    sum.AvailabilityChart(availabilityData, availabilitySeries);
                } else {
                    sum.availData(availData);
                    sum.availSeries(availabilitySeries);
                    sum.AvailabilityChart(availData, availabilitySeries);
                }
                // sum.AvailabilityChart(res.data["Data"]);
            }

            
            
            if (res.data != null ){
                sum.isDetailProd(false);
                sum.isDetailProdByProject(false);
            }
        });

        var ajax3

        if (project=="Fleet") {
            param = { ProjectName: project, Date: maxdate, PeriodList: sum.paramPeriod};
            ajax3 = toolkit.ajaxPost(viewModel.appName + "dashboard/getwinddistribution", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }

                sum.WindDistribution(res.data.Data[sum.periodSelected()]);
                sum.dataSourceWindDistribution(res.data.Data[sum.periodSelected()]);
                sum.windDistData(res.data.Data);
            });
        }

        $.when(sum.indiaMap(project),ajax1, ajax2, ajax3).done(function(){
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
                { title: "Project Name", width:100, field: "name", headerAttributes: { style: "text-align:left;" }, attributes: { style: "text-align:left;" } },
                { title: "No. of WTG", width:90, field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Production<br>(GWh)", width:100, field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "PLF<br>(%)", width:100, field: "plf", format: "{0:n2}", template: "#= kendo.toString(plf*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Lost Energy<br>(MWh)", width:100,field: "lostenergy", template: "#= kendo.toString(lostenergy/1000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Downtime<br>(Hours)", width:120,field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Machine Availability<br>(%)", width:120, field: "machineavail", format: "{0:n2}", template: "#= kendo.toString(machineavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Total Availability<br>(%)", width:120, field: "trueavail", format: "{0:n2}", template: "#= kendo.toString(trueavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
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
            height: 160,
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
            height: 170,
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
            height: 160,
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
            height: 160,
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
        seriesColors: colorField,
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
            height: 150,
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
        series: [{
            name: "Budget",
            field: "Budget",
            // opacity : 0.7,
            color: "#21c4af",
        }, {
            name: "Production",
            field: "Production",
            // opacity : 0.7,
            color: "#ff9933",
        }],
        // seriesColors: colorField,
        seriesClick: function (e) {
            sum.DetailProd(e);
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
sum.AvailabilityChart = function (dataSource, dataSeries) {
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
            height: 165,
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
            height: 160,
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
        var jsonObj = res.data,i;

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
            height: 165,
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
            height: 170,
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
        series: [{
            name: "Budget",
            field: "CumBudget",
            // opacity : 0.5,
            color: "#21c4af",
        }, {
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
    var bulan = e.category;
    sum.detailProdTxt(bulan);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(true);
    sum.isDetailProdByProject(false);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { 'project': project, 'date': bulan };

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var dataSource = res.data;
        var measurement = " (" + dataSource[0].measurement + ") ";
        sum.detailProdMsTxt(measurement);

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
                sum.DetailProdByProject(e, bulan, dataSource);
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
                sum.DetailProdByProject(e, bulan, dataSource);
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
    });
    
    
}

sum.DetailProdByProject = function (e, month, data) {
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(true);
    sum.detailProdProjectTxt(e.category);
    sum.detailProdDateTxt(month);
    var dataSource;
    var measurement = '';

    $.each(data, function (i, val) {
        if (val.project == e.category) {
            dataSource = val.detail;
        }
        if (i == 0) {
            measurement = " (" + dataSource[0].measurement + ") "
        }
    })


    $("#chartDetailProdByProject").kendoChart({
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
            // opacity : 0.7,
        }, {
            name: "Lost Energy" + measurement,
            field: "lostenergy",
            // opacity : 0.7,
        }],
        seriesColors: colorField,
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
            }
        },
        tooltip: {
            visible: true,
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
        pageable: {
            pageSize: 10,
            input: true, 
        },
        columns: [
            { title: "Turbine Name", field: "turbine", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            { title: "Production<br>" + measurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            { title: "Lost Energy<br>" + measurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            { title: "Downtime<br>(Hours)", field: "downtime", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
        ],
        dataSource: {
            data: dataSource,
            sort: { field: "turbine", dir: 'asc' },
            pageSize: 10
        }
    });
}

sum.backToDashboard = function () {
    vm.isDashboard(true);
    lgd.isSummary(true);
    sum.isDetailProd(false);
    sum.isDetailProdByProject(false);
}
sum.toDetailProduction = function () {
    sum.isDetailProd(true);
    sum.isDetailProdByProject(false);
}