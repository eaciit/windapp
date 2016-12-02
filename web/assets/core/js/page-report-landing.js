'use strict';

var monthNames = ["January", "February", "March", "April", "May", "June",
    "July", "August", "September", "October", "November", "December"
];

viewModel.landing = {};
var lgd = viewModel.landing;

lgd.filter = ko.observableArray([
    { text: "Fleet", value: "1" },
    { text: "WindFarm-01", value: "2" },
    { text: "WindFarm-02", value: "3" }
]);

lgd.isDetailProd = ko.observable(false);
lgd.isDetailProdByProject = ko.observable(false);

lgd.isSummary = ko.observable(false);
lgd.isDowntime = ko.observable(false);
lgd.isProduction = ko.observable(false);
lgd.isAvailability = ko.observable(false);

lgd.isDetailDTLostEnergy = ko.observable(false);
lgd.detailDTLostEnergyTxt = ko.observable();

lgd.isDetailDTTop = ko.observable(false);
lgd.detailDTTopTxt = ko.observable();
lgd.detailProdTxt = ko.observable('');
lgd.detailProdMsTxt = ko.observable('');
lgd.detailProdProjectTxt = ko.observable('');
lgd.detailProdDateTxt = ko.observable('');

lgd.noOfProjects = ko.observable();
lgd.noOfTurbines = ko.observable();
lgd.totalMaxCapacity = ko.observable();
lgd.currentDown = ko.observable();
lgd.twoDaysDown = ko.observable();

lgd.projectList = ko.observableArray([{ "value": "Fleet", "text": "Fleet" }]);
lgd.projectItem = ko.observableArray([]);
lgd.mdTypeList = ko.observableArray([]);
lgd.projectName = ko.observable();
lgd.isFleet = ko.observable(true);
lgd.isNonFleet = ko.observable(true);
lgd.FleetDTLEDownType = ko.observable();
lgd.LEFleetByDown = ko.observable(false);

lgd.prodDateRangeStr = ko.observable('');

var isFirst = true;
var isFleetDetail = false;
var lastParam = {};
var lastParamChart = {};
var dtType = '';
var monthDetailDT = '';
var projectSelected = '';
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

lgd.getProjectList = function () {
    app.ajaxPost(viewModel.appName + "/dashboard/getprojectlist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0) {
            res.data = [];

        } else {
            if (res.data.length > 0) {
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    lgd.projectList.push(data);
                    lgd.projectItem.push(data);
                });
            }
        }
    });
};

lgd.getMDTypeList = function () {
    app.ajaxPost(viewModel.appName + "/dashboard/getmdtypelist", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if (res.data.length == 0) {
            res.data = [];
        } else {
            if (res.data.length > 0) {
                /*var def = {};
                def.value = "All Type";
                def.text = "All Type";
                lgd.mdTypeList.push(def);*/
                $.each(res.data, function (key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    lgd.mdTypeList.push(data);
                });
            }
        }
    });
};

lgd.backToDashboard = function () {
    vm.isDashboard(true);
    lgd.isDetailProd(false);
    lgd.isDetailProdByProject(false);
}
lgd.toDetailProduction = function () {
    lgd.isDetailProd(true);
    lgd.isDetailProdByProject(false);
}

lgd.getDetailDT = function () {
    if (isFirst == false) {
        dtType = $("#mdTypeList").data("kendoDropDownList").value();

        if (dtType == "") {
            lgd.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - All Type");
        } else {
            lgd.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - " + dtType);
        }

        lgd.toDetailDTLostEnergy(null, true, "ddl");
    }
}

lgd.getDetailDTFromProject = function () {
    if (isFirst == false) {
        projectSelected = $("#projectList").data("kendoDropDownList").value();
        dtType = $("#mdTypeListFleet").data("kendoDropDownList").value();
        lgd.detailDTLostEnergyTxt("Lost Energy for Last 12 months - " + projectSelected);
        lgd.toDetailDTLostEnergy(null, false, "ddl");
    }
}

lgd.toDetailDTLostEnergy = function (e, isDetailFleet, source) {
    app.loading(true);
    vm.isDashboard(false);
    lgd.isDetailDTLostEnergy(true);

    var project = $("#projectId").data("kendoDropDownList").value();
    var dateStr = '';
    var type = '';
    var param = {};
    var paramChart = {};
    var method = "getdowntime";

    if (source == "chart" || source == "chartbytype") {
        monthDetailDT = e.category;
        lgd.detailDTLostEnergyTxt("Lost Energy for " + monthDetailDT + " - " + e.series.name);

        if (source == "chart") {
            dateStr = e.category;
        }
        type = e.series.name;
    } else if (source == "button") {
        lgd.detailDTLostEnergyTxt("Lost Energy for Last 12 months - " + lastParam.Type);
    }
    if (project == "Fleet" && isDetailFleet == false) { /*by type level 1*/
        $(".show_hide_downtime").hide();
        $(".show_hide_project").show();

        if (source == "button" || source == "ddl") {
            if (source == "button") {
                if (lgd.LEFleetByDown() == true) {
                    method = "getdowntimefleetbydown";
                }

                param = lastParam;
                paramChart = lastParamChart;
                $("#projectList").data("kendoDropDownList").value(param.Type);
            } else if (source == "ddl") {
                if (lgd.LEFleetByDown() == true) {
                    method = "getdowntimefleetbydown";
                    paramChart = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType };
                } else {
                    if (dtType == "") {
                        dtType = "All Types"
                    }
                    paramChart = { ProjectName: projectSelected, Date: lastParamChart.Date, Type: dtType };
                }

                param = { ProjectName: projectSelected, DateStr: lastParam.DateStr, Type: dtType };
                lastParam = param;
                lastParamChart = paramChart;
            }
            toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!toolkit.isFine(res)) {
                    return;
                }

                if (method == "getdowntimefleetbydown") {
                    lgd.DTLostEnergyByDown(res.data.lostenergy);
                } else {
                    lgd.DTLostEnergyManeh(res.data.lostenergy);
                }

                app.loading(false);
                lgd.setDownTimeSeriesCheck();
            });
        } else {
            $("#projectList").data("kendoDropDownList").select(0);
            projectSelected = $("#projectList").data("kendoDropDownList").value();

            if (source == "chart") {
                $("#mdTypeListFleet").data("kendoDropDownList").value(0);
                paramChart = { ProjectName: projectSelected, DateStr: e.category };
                param = { ProjectName: projectSelected, DateStr: e.category };

                method = "getdowntimefleetbydown";
            } else {
                $("#mdTypeListFleet").data("kendoDropDownList").value(e.category);
                paramChart = { ProjectName: projectSelected, Date: maxdate, Type: e.category };
                param = { ProjectName: projectSelected, DateStr: "fleet date", Type: e.category };
            }

            lastParam = param;
            lastParamChart = paramChart;

            toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, paramChart, function (res) {
                if (!toolkit.isFine(res)) {
                    return;
                }

                if (method == "getdowntimefleetbydown") {
                    lgd.DTLostEnergyByDown(res.data.lostenergy);
                } else {
                    lgd.DTLostEnergyManeh(res.data.lostenergy);
                }

                lgd.FleetDTLEDownType = e.category;
                lgd.setDownTimeSeriesCheck();
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
            if (!toolkit.isFine(res)) {
                return;
            }

            var dataSource = res.data;
            lgd.DTLostEnergyDetail(dataSource);

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
        pageable: true,
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

lgd.toDetailDTTop = function (e, type) {
    vm.isDashboard(false);

    if (type == "Times") {
        lgd.detailDTTopTxt("(" + e.category + ") - Frequency");
    } else {
        lgd.detailDTTopTxt("(" + e.category + ") - " + type);
    }
    lgd.isDetailDTTop(true);

    // get the data and push into the chart    
    lgd.DTTopDetail(e.category, type);
}

lgd.backToDownTime = function () {
    vm.isDashboard(true);

    lgd.isSummary(false);
    lgd.isProduction(false);
    lgd.isAvailability(false);
    lgd.isDowntime(true);

    lgd.isDetailDTLostEnergy(false);
    lgd.detailDTLostEnergyTxt("Lost Energy for Last 12 months");

    lgd.isDetailDTTop(false);
    lgd.detailDTTopTxt("");
}

lgd.setDownTimeSeriesCheck = function () {
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

lgd.backToDownTimeChart = function () {
    var project = $("#projectId").data("kendoDropDownList").value();
    if (project == "Fleet" && isFirst == false && isFleetDetail == true) {
        vm.isDashboard(false);
        lgd.isDetailDTLostEnergy(true);
        isFleetDetail = false;
        lgd.toDetailDTLostEnergy(null, false, "button");
        if ($("#projectList").data("kendoDropDownList") != null) {
            $("#projectList").data("kendoDropDownList").value(projectSelected);
        }
    } else {
        lgd.LEFleetByDown(false);
        vm.isDashboard(true);
        lgd.isSummary(false);
        lgd.isProduction(false);
        lgd.isAvailability(false);
        lgd.isDowntime(true);
        lgd.isDetailDTLostEnergy(false);
        lgd.detailDTLostEnergyTxt("Lost Energy for Last 12 months");
        lgd.isDetailDTTop(false);
        lgd.detailDTTopTxt("");
    }
}

lgd.LoadData = function () {
    app.loading(true);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { ProjectName: project, Date: maxdate };

    if (project == "Fleet") {
        // $("#lostcontrol").removeClass("col-md-12 col-sm-12").addClass("col-md-6 col-sm-6")
        lgd.isFleet(true);
        lgd.isNonFleet(false);
        $("#div-windiness").hide();


    } else {
        // $("#lostcontrol").removeClass("col-md-6 col-sm-6").addClass("col-md-12 col-sm-12")
        lgd.isFleet(false);
        lgd.isNonFleet(true);
        $("#div-windiness").show();
    }

    setTimeout(function () {

        toolkit.ajaxPost(viewModel.appName + "dashboard/getscadalastupdate", param, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }

            lgd.noOfProjects(res.data[0].NoOfProjects);
            lgd.noOfTurbines(res.data[0].NoOfTurbines);
            lgd.totalMaxCapacity(res.data[0].TotalMaxCapacity / 1000);
            lgd.currentDown(res.data[0].CurrentDown);
            lgd.twoDaysDown(res.data[0].TwoDaysDown);

            var lastUpdate = new Date(res.data[0].LastUpdate);

            vm.dateAsOf(lastUpdate.addHours(-7));

            lgd.ProductionChart(res.data[0].Productions);
            lgd.CumProduction(res.data[0].CummulativeProductions);
            lgd.SummaryData(project);
        });

        toolkit.ajaxPost(viewModel.appName + "dashboard/getscadasummarybymonth", param, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }

            lgd.PLF(res.data);
            lgd.LostEnergy(res.data);
            lgd.Windiness(res.data);
            lgd.ProdMonth(res.data);
            lgd.AvailabilityChart(res.data);
            lgd.ProdCurLast(res.data);
        });

        toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntime", param, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }

            lgd.DTLostEnergy(res.data.lostenergy);
            if (project == "Fleet") {
                lgd.DTLEbyType(res.data.lostenergybytype[0]);
            }
            lgd.DTDuration(res.data.duration);
            lgd.DTFrequency(res.data.frequency);
        });

        lgd.gridProduction(project, maxdate);
        // lgd.gridAvailability(project, maxdate);

        lgd.DTTurbines();
        lgd.indiaMap(project);
        app.loading(false);

        // $("#india-map").data("kendoMap").zoom((project == 'Fleet'? 3 : 10));        
    }, 600);
}


lgd.SummaryData = function (project) {
    var filters = [
        { field: "_id", operator: "eq", value: project },
    ];
    var filter = { filters: filters }
    var param = { filter: filter };
    $('#gridSummaryData').html("");
    $("#gridSummaryData").kendoGrid({
        height: 155,
        theme: "flat",
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "dashboard/getsummarydata",
                    type: "POST",
                    data: param,
                    dataType: "json",
                    contentType: "application/json; charset=utf-8"
                },
                parameterMap: function (options) {
                    return JSON.stringify(options);
                }
            },
            pageSize: 2,
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
            sort: [
                { field: 'name', dir: 'asc' },
            ],
        },
        /*serverPaging: true,
        serverSorting: true,*/
        pageable: true,
        columns: [
            { title: "Project Name", field: "name", headerAttributes: { style: "text-align:left;" }, attributes: { style: "text-align:left;" } },
            { title: "No. of WTG", field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Production<br>(GWh)", field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "PLF<br>(%)", field: "plf", width: 80, format: "{0:n2}", template: "#= kendo.toString(plf*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Lost Energy<br>(MWh)", field: "lostenergy", template: "#= kendo.toString(lostenergy/1000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Downtime Hours", field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Machine Availability<br>(%)", field: "machineavail", format: "{0:n2}", template: "#= kendo.toString(machineavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Total Availability<br>(%)", field: "trueavail", format: "{0:n2}", template: "#= kendo.toString(trueavail*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
        ]
    });

    setTimeout(function () {
        var grid = $("#gridSummaryData").data("kendoGrid");
        if (project == "Fleet") {
            $("#gridSummaryData th[data-field=name]").html("Project Name")
            grid.showColumn("noofwtg");
        } else {
            $("#gridSummaryData th[data-field=name]").html("Turbine Name")
            grid.hideColumn("noofwtg");
        }
        var dataSource = grid.dataSource.data();
        $.each(dataSource, function (i, row) {
            $('tr[data-uid="' + row.uid + '"]').css("border-bottom", "1pt solid black");
        });
        $("#gridSummaryData").data("kendoGrid").refresh();
    }, 100);

    var result = parseFloat(parseInt(41.5, 10) * 100) / parseInt(lgd.totalMaxCapacity(), 10);

    lgd.createDonutChart({ id: 'sumProductionChart', value: result, title: "Production" });
    lgd.createDonutChart({ id: 'sumTotalAvailChart', value: 95, title: "Total Availability" });
    lgd.createDonutChart({ id: 'sumPerfBudgetChart', value: 0, title: "Performance vs Budget" });
    lgd.createDonutChart({ id: 'sumAchievmentAnnualChart', value: 0, title: "Achievement vs Annual Budget" });
}

lgd.createDonutChart = function (param) {
    $('#' + param.id).attr("data-percent", param.value);
    $('#' + param.id).pieChart({
        barColor: '#ea5b19',
        trackColor: '#fff',
        lineCap: 'round',
        lineWidth: 4,
        size: 65,
        onStep: function (from, to, percent) {
            $(this.element).find('.pie-value').html(Math.round(percent) + '%');
        }
    });
}

lgd.ProductionChart = function (dataSource) {
    var dataFormat = "n2";
    if (dataSource.length > 0) {
        var totalPotential = 0;
        for (var i = 0; i < dataSource.length; i++) {
            totalPotential += dataSource[i].PotentialKwh;
        }
        if (totalPotential > 10) {
            dataFormat = "n0";
        }
    }
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
            position: "top"
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
            color: "#ff880e"
        }],
        valueAxis: {
            labels: {
                step: 2,
                format: dataFormat
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
        // seriesHover: function(e) {
        //   var positionX = e.originalEvent.clientX,
        //       positionY = e.originalEvent.clientY,
        //       value = e.value;
        //   $("#chartProductionCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateProd").html())({ e:e }));;                  
        // },
    });

    // $("#chartProduction").mouseleave(function(e){
    //    $("#chartProductionCustomTooltip").hide();
    // })
}

lgd.CumProduction = function (dataSource) {
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
            position: "top"
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
            type: "area",
            area: {
                line: {
                    style: "smooth"
                }
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
            color: "#ff880e"
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
                format: "n0"
            },
        },
        categoryAxis: {
            field: "DateId",
            majorGridLines: {
                visible: false
            },
            labels: {
                step: 3,
                template: "#: Number(kendo.toString(kendo.parseDate(value), 'dd'))#"
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

lgd.PLF = function (dataSource) {
    $("#chartPLF").kendoChart({
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
                format: "{0}"
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
                template: '#=  value.substring(0,3) #'
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
lgd.LostEnergy = function (dataSource) {
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
                template: '#=  value.substring(0,3) #'
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

lgd.Windiness = function (dataSource) {
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
            position: "top"
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
                format: "n0"
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
                template: '#=  value.substring(0,3) #'
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
        // seriesHover: function(e) {
        //   console.log(e)
        //   var positionX = e.originalEvent.offsetX,
        //       positionY = e.originalEvent.offsetY,
        //       value = e.value;
        //   $("#chartWindinessCustomTooltip").show().css('position', 'absolute').css("top", positionY).css("left", positionX).html(kendo.template($("#templateWindiness").html())({ e:e }));;                  
        // },
    });

    //  $("#chartWindiness").mouseleave(function(e){
    //    $("#chartWindinessCustomTooltip").hide();
    // })
}

lgd.ProdMonth = function (dataSource) {
    $("#chartProdMonth").kendoChart({
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
            color: "#ff880e",
        }],
        // seriesColors: colorField,
        seriesClick: function (e) {
            lgd.DetailProd(e);
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
                format: "n0"
            },
        }],
        categoryAxis: {
            field: "DateInfo.MonthDesc",
            majorGridLines: {
                visible: false
            },
            labels: {
                template: '#=  value.substring(0,3) #'
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

lgd.AvailabilityChart = function (dataSource) {
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
            position: "top"
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
            name: "TBA",
            field: "ScadaAvail",
            // opacity : 0.5,
            color: "#21c4af"
        }, {
            name: "DBA",
            field: "TrueAvail",
            // opacity : 0.5,
            color: "#ff880e",
        }],
        // seriesColors: colorField,
        valueAxis: {
            max: 100,
            majorUnit: 25,
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            labels: {
                // format: "{0}%"
                format: "{0}"
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
                template: '#=  value.substring(0,3) #'
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

lgd.ProdCurLast = function (dataSource) {
    $("#chartCurrLast").kendoChart({
        dataSource: {
            data: dataSource,
            sort: { field: "DateInfo.MonthId", dir: 'asc' }
        },
        theme: "flat",
        title: {
            text: ""
        },
        legend: {
            position: "top"
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
        series: [{
            type: "column",
            style: "smooth",
            name: "Current",
            field: "Production",
            // opacity : 0.7,
            axis: "production",
            color: "#21c4af",
        }, {
            type: "column",
            style: "smooth",
            name: "Last",
            field: "ProductionLastYear",
            // opacity : 0.7,
            axis: "production",
            color: "#ff880e",
        }, {
            type: "line",
            style: "smooth",
            name: "Variance(%)",
            field: "Variance",
            axis: "variance",
            color: "#ff7663",
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

lgd.indiaMap = function (project) {
    $("#india-map").html("");
    var param = { projectname: project }

    toolkit.ajaxPost(viewModel.appName + "dashboard/getmapdata", param, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }

        var turbineInfos = res.data;
        var center = turbineInfos[0].coords[0] + "," + turbineInfos[0].coords[1];
        var mapProp = {
            center: (param.projectname == 'Fleet' ? new google.maps.LatLng(22.460533, 79.650879) : new google.maps.LatLng(27.131461, 70.618559)),
            zoom: (param.projectname == 'Fleet' ? 4 : 10),
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
        };
        var map = new google.maps.Map(document.getElementById("india-map"), mapProp);

        var markers = new Array();

        turbineInfos.forEach(function (obj, idx) {
            var marker = new google.maps.Marker({
                position: new google.maps.LatLng(obj.coords[0], obj.coords[1]),
                map: map,
                title: obj.name,
                icon: {
                    url: "../res/img/wind-turbine.png", // url
                    scaledSize: new google.maps.Size(30, 30), // scaled size
                }
            });

            var infowindow = new google.maps.InfoWindow({
                content: ""
            });

            google.maps.event.addListener(marker, 'click', function () {
                map.panTo(this.getPosition());
                map.setZoom(20);
            });
        });
    });

}

lgd.DetailProd = function (e) {
    var bulan = e.category;
    lgd.detailProdTxt(bulan);
    vm.isDashboard(false);
    lgd.isDetailProd(true);
    lgd.isDetailProdByProject(false);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { 'project': project, 'date': bulan };

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        var dataSource = res.data;
        var measurement = " (" + dataSource[0].measurement + ") ";
        lgd.detailProdMsTxt(measurement);

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
            }, {
                name: "Lost Energy" + measurement,
                field: "lostenergy",
                gap: 3
                // opacity : 0.7,
            }],
            seriesColors: colorField,
            seriesClick: function (e) {
                lgd.DetailProdByProject(e, bulan, dataSource);
            },
            valueAxis: {
                // majorUnit : 2000,
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
                field: "project",
                majorGridLines: {
                    visible: false
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
            pageable: true,
            columns: [
                { title: "Project Name", field: "project", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "No. of WTG", field: "wtg", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
                { title: "Production <br>" + measurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
                { title: "Lost Energy <br>" + measurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            ],
            dataSource: {
                data: dataSource,
                pageSize: 5
            }
        });
    });
}

lgd.DetailProdByProject = function (e, month, data) {
    vm.isDashboard(false);
    lgd.isDetailProd(false);
    lgd.isDetailProdByProject(true);
    lgd.detailProdProjectTxt(e.category);
    lgd.detailProdDateTxt(month);
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
                rotation: 45
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
        pageable: true,
        columns: [
            { title: "Turbine Name", field: "turbine", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            { title: "Production<br>" + measurement, field: "production", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            { title: "Lost Energy<br>" + measurement, field: "lostenergy", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
            { title: "Downtime<br>(Hrs)", field: "downtime", format: "{0:n2}", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" } },
        ],
        dataSource: {
            data: dataSource,
            sort: { field: "turbine", dir: 'asc' },
            pageSize: 10
        }
    });
}

lgd.DTLostEnergy = function (dataSource) {
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
            field: "_id.id2",
            majorGridLines: {
                visible: false
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
            lgd.toDetailDTLostEnergy(e, false, "chart");
        }
    });
    //  $("#chartDTLostEnergy").mouseleave(function(e){
    //    $("#chartDTLostEnergyCustomTooltip").hide();
    // })
}

lgd.DTLEbyType = function (dataSource) {
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
                rotation: -330
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
            lgd.toDetailDTLostEnergy(e, false, "chartbytype");
        }
    });
    // $("#chartDTLEbyType").mouseleave(function(e){
    //    $("#chartDTLEbyTypeCustomTooltip").hide();
    // })
}

lgd.DTLostEnergyManeh = function (dataSource) {
    lgd.detailDTLostEnergyTxt("Lost Energy for Last 12 months");
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
            field: "_id.id2",
            majorGridLines: {
                visible: false
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
            lgd.toDetailDTLostEnergy(e, true, "chart");
        }
    });
    // $("#chartDTLostEnergyDetail").mouseleave(function(e){
    //    $("#chartDTLostEnergyManehCustomTooltip").hide();
    // })
    setTimeout(function () {
        $("#chartDTLostEnergyDetail").data("kendoChart").refresh();
    }, 100);
}

lgd.DTLostEnergyByDown = function (dataSource) {
    lgd.LEFleetByDown(true)
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
            categories: [lastParam.DateStr],
            majorGridLines: {
                visible: false
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
            lgd.toDetailDTLostEnergy(e, true, "chart");
        }
    });

    setTimeout(function () {
        $("#chartDTLostEnergyDetail").data("kendoChart").refresh();
    }, 100);
}

lgd.DTDuration = function (dataSource) {
    $("#chartDTDuration").kendoChart({
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
        },
        chartArea: {
            height: 160
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [{
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
        }],
        seriesColors: colorField,
        valueAxis: {
            //majorUnit: 100,
            title: {
                text: "Hours",
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
            field: "_id",
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
            lgd.toDetailDTTop(e, "Hours");
        }
    });
}

lgd.DTFrequency = function (dataSource) {
    $("#chartDTFrequency").kendoChart({
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
        },
        chartArea: {
            height: 160
        },
        seriesDefaults: {
            type: "column",
            stack: true,
            // opacity : 0.7
        },
        series: [{
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
        }],
        seriesColors: colorField,
        valueAxis: {
            title: {
                text: "Times",
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            name: "result",
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
            field: "_id",
            dir: "desc",
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
            lgd.toDetailDTTop(e, "Times");
        }
    });
}

lgd.DTLostEnergyDetail = function (dataSource) {
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
                rotation: -330
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

lgd.DTTopDetail = function (turbine, type) {
    app.loading(true);
    var project = $("#projectId").data("kendoDropDownList").value();
    var date = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));
    var param = { ProjectName: project, Date: date, Type: type, Turbine: turbine };

    var template = (type == 'Hours' ? "#: category # : #:  kendo.toString(value, 'n1') #" : "#: category # : #:  kendo.toString(value, 'n0') #")
    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimetopdetail", param, function (res) {
        if (!toolkit.isFine(res)) {
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
                field: "_id.id2",
                majorGridLines: {
                    visible: false
                },
                labels: {
                    template: '#=  value.substring(0,3) #'
                },
                majorTickType: "none"
            },
            tooltip: {
                visible: true,
                template: template,
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
        pageable: true,
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

lgd.DTTurbines = function () {
    var project = $("#projectId").data("kendoDropDownList").value();
    var date = new Date(Date.UTC(2016, 5, 30, 23, 50, 0, 0));
    var param = { ProjectName: project, Date: date };

    $("#dtturbines").html("");

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdowntimeturbines", param, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }

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
    });
}

lgd.periodTypeProdChange = function () {
    lgd.gridProduction($("#projectId").data("kendoDropDownList").value(), maxdate);
}

// lgd.periodTypeAvailChange = function () {
//     lgd.gridAvailability($("#projectId").data("kendoDropDownList").value(), maxdate);
// }

lgd.gridProduction = function (project, enddate) {
    var filters = [];
    var type = $('input[name="periodTypeProd"]:checked').val();
    var method, startDate;

    var endDateMonth = enddate.getUTCMonth();
    var endDateYear = enddate.getUTCFullYear();
    var endDateDate = enddate.getUTCDate();

    if (type == "lw") {
        // startDate = new Date(Date.UTC(2016,5,23,0,0,0,0));
        startDate = new Date(Date.UTC(endDateYear, endDateMonth, endDateDate - 7, 0, 0, 0, 0));
        filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
        filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
        filters.push({ field: "type", operator: "eq", value: type });
        if (project != "Fleet") {
            filters.push({ field: "projectname", operator: "eq", value: project });
        }
        method = "getsummarydatadaily";
    } else if (type == "mtd") {
        // startDate = new Date(Date.UTC(2016,5,1,0,0,0,0));
        startDate = new Date(Date.UTC(endDateYear, endDateMonth, 1, 0, 0, 0, 0));
        filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
        filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
        filters.push({ field: "type", operator: "eq", value: type });
        if (project != "Fleet") {
            filters.push({ field: "projectname", operator: "eq", value: project });
        }
        method = "getsummarydatadaily";
    } else if (type == "ytd") {
        // filters.push({ field: "_id", operator: "eq", value: project });
        // method = "getsummarydata";

        // startDate = new Date(Date.UTC(2015,6,1,0,0,0,0));
        startDate = new Date(Date.UTC(endDateYear, 0, 1, 0, 0, 0, 0));
        filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
        filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
        filters.push({ field: "type", operator: "eq", value: type });
        if (project != "Fleet") {
            filters.push({ field: "projectname", operator: "eq", value: project });
        }
        method = "getsummarydatadaily";
    }

    var startDateStr = startDate.getUTCDate() + "-" + monthNames[startDate.getUTCMonth()] + "-" + startDate.getUTCFullYear();
    var endDateStr = enddate.getUTCDate() + "-" + monthNames[enddate.getUTCMonth()] + "-" + enddate.getUTCFullYear();

    $('#prodDateRangeStr').html(startDateStr + " to " + endDateStr);

    var filter = { filters: filters }
    var param = { filter: filter };
    $('#productionGrid').html("");
    $("#productionGrid").kendoGrid({
        dataSource: {
            serverPaging: true,
            serverSorting: true,
            transport: {
                read: {
                    url: viewModel.appName + "dashboard/" + method,
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
                    return res.data.Data
                },
                total: function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    return res.data.Total;
                }
            }/*,
        sort: [
            { field: '_id', dir: 'asc' },
        ],*/
        },
        /*serverPaging: true,
        serverSorting: true,*/
        pageable: true,
        columns: [
            { title: "Project Name", field: "name", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:left;" } },
            { title: "No. of WTG", width: 80, field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Max Capacity<br>(GWh)", field: "maxcapacity", template: "#= kendo.toString(maxcapacity/1000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Production<br>(GWh)", field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "PLF<br>(%)", field: "plf", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(plf*100, 'n2') #%" },
            { title: "Total Availability<br>(%)", field: "totalavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(totalavail*100, 'n2') #%" },
            // { title: "Production Ratio", field: "lostEnergy",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-right" }},
            // { title: "Worst Single Production Ratio", field: "lostEnergy",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-right" }},
            { title: "Lowest Machine Availability<br>(%)", field: "lowestmachineavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            { title: "Lowest PLF<br>(%)", field: "lowestplf", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
            // { title: "Max. Lost Energy to Effeciency", field: "lostenergy",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-right" },template: "#= kendo.toString(lostenergy, 'n2') #"},
            { title: "Max. Lost Energy due to Downtime<br>(KWh)", field: "maxlossenergy", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
        ],
    });

    setTimeout(function () {
        var grid = $("#productionGrid").data("kendoGrid");
        if (project == "Fleet") {
            $("#productionGrid th[data-field=name]").html("Project Name")
            grid.showColumn("noofwtg");
        } else {
            $("#productionGrid th[data-field=name]").html("Turbine Name")
            grid.hideColumn("noofwtg");
        }
        $("#productionGrid").data("kendoGrid").refresh();
    }, 100);
}

// lgd.gridAvailability = function(project, enddate){
//     var filters = [];
//     var type = $('input[name="periodTypeAvail"]:checked').val();
//     var method, startDate;

//     var endDateMonth = enddate.getUTCMonth();
//     var endDateYear = enddate.getUTCFullYear();
//     var endDateDate = enddate.getUTCDate();

//     if (type=="lw") {
//         // startDate = new Date(Date.UTC(2016,5,23,0,0,0,0));
//         startDate = new Date(Date.UTC(endDateYear, endDateMonth, endDateDate - 7, 0,0,0,0));
//         filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
//         filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
//         filters.push({ field: "type", operator: "eq", value: type });
//         if (project != "Fleet"){
//             filters.push({ field: "projectname", operator: "eq", value: project });
//         }
//         method = "getsummarydatadaily";
//     }else if (type=="mtd") {
//         // startDate = new Date(Date.UTC(2016,5,1,0,0,0,0));
//         startDate = new Date(Date.UTC(endDateYear, endDateMonth, 1, 0,0,0,0));
//         filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
//         filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
//         filters.push({ field: "type", operator: "eq", value: type });
//         if (project != "Fleet"){
//             filters.push({ field: "projectname", operator: "eq", value: project });
//         }
//         method = "getsummarydatadaily";
//     }else if (type=="ytd") {
//         // filters.push({ field: "_id", operator: "eq", value: project });
//         // method = "getsummarydata";

//         // startDate = new Date(Date.UTC(2015,6,1,0,0,0,0));
//         startDate = new Date(Date.UTC(endDateYear, 0, 1, 0,0,0,0));
//         filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
//         filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
//         filters.push({ field: "type", operator: "eq", value: type });
//         if (project != "Fleet"){
//             filters.push({ field: "projectname", operator: "eq", value: project });
//         }
//         method = "getsummarydatadaily";
//     }

//     var startDateStr = startDate.getUTCDate() +"-"+monthNames[startDate.getUTCMonth()]+"-"+startDate.getUTCFullYear();
//     var endDateStr = enddate.getUTCDate() +"-"+monthNames[enddate.getUTCMonth()]+"-"+enddate.getUTCFullYear();

//     $('#availDateRangeStr').html(startDateStr+" to "+endDateStr);

//     var filter = {filters : filters}
//     var param = {filter : filter};
//     $('#availabilityGrid').html("");
//     $("#availabilityGrid").kendoGrid({
//       dataSource: {
//         serverPaging: true,
//         serverSorting: true,
//         transport: {
//           read: {
//             url: viewModel.appName + "dashboard/"+method,
//             type: "POST",
//             data: param,
//             dataType: "json",
//             contentType: "application/json; charset=utf-8"
//           },
//           parameterMap: function(options) {
//             return JSON.stringify(options);
//           }
//         },
//         pageSize: 10,
//         schema: {
//           data: function(res){
//             if (!app.isFine(res)) {
//                 return;
//             }
//             return res.data.Data
//           },
//           total: function(res){
//             if (!app.isFine(res)) {
//                 return;
//             }
//             return res.data.Total;
//           }
//         }/*,
//         sort: [
//             { field: 'name', dir: 'asc' },
//         ],*/
//       },  
//       pageable: true,
//       columns: [
//             { title: "Project Name", field: "name", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:left;"} },
//             { title: "Production<br>(GWh)", field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes:{style:"text-align:center;"}, attributes:{style:"text-align:center;"} },
//             { title: "PLF<br>(%)", field: "plf",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" },template: "#= kendo.toString(plf*100, 'n2') #"},
//             { title: "Machine Availability<br>(%)", field: "machineavail",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" },template: "#= kendo.toString(machineavail*100, 'n2') #"},
//             { title: "Grid Availability<br>(%)", field: "gridavail",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" },template: "#= kendo.toString(gridavail*100, 'n2') #"},
//             { title: "Total Availability<br>(%)", field: "totalavail",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" },template: "#= kendo.toString(totalavail*100, 'n2') #"},
//         ],
//     });

//     setTimeout(function() {
//         var grid = $("#availabilityGrid").data("kendoGrid");
//         if(project == "Fleet") {
//             $("#availabilityGrid th[data-field=name]").html("Project Name")    
//             grid.showColumn("noofwtg");
//         } else {
//             $("#availabilityGrid th[data-field=name]").html("Turbine Name")
//             grid.hideColumn("noofwtg");
//         }
//         $("#availabilityGrid").data("kendoGrid").refresh();
//     }, 100);

// }


vm.currentMenu('Dashboard');
vm.currentTitle('Dashboard');
vm.isDashboard(true);
vm.breadcrumb([{ title: 'Dashboard', href: viewModel.appName + 'page/landing' }, { title: 'Home', href: '#' }]);

$(function () {
    lgd.isSummary(true);
    lgd.isProduction(false);
    lgd.isAvailability(false);
    lgd.isDowntime(false);
    lgd.getProjectList();
    lgd.getMDTypeList();
    lgd.projectName("Fleet");
    lgd.isDetailProd(false);
    lgd.isDetailProdByProject(false);
    lgd.isDetailDTLostEnergy(false);
    lgd.isDetailDTTop(false);

    lgd.LoadData();

    $("#tabSummary").on("click", function () {
        lgd.LoadData();
        lgd.isSummary(true);
        lgd.isDowntime(false);
        lgd.isProduction(false);
        lgd.isAvailability(false);
        isFirst = true;

    });

    $("#tabProduction").on("click", function () {
        lgd.isSummary(false);
        lgd.isDowntime(false);
        lgd.isProduction(true);
        lgd.isAvailability(false);
        isFirst = false;

    });

    $("#tabAvailability").on("click", function () {
        lgd.isSummary(false);
        lgd.isDowntime(false);
        lgd.isProduction(false);
        lgd.isAvailability(true);
        isFirst = false;

    });

    $("#tabDowntime").on("click", function () {
        lgd.backToDownTime();
        isFirst = false;

        $("#chartDTLostEnergy").data("kendoChart").refresh();
        $("#chartDTLEbyType").data("kendoChart").refresh();
        $("#chartDTDuration").data("kendoChart").refresh();
        $("#chartDTFrequency").data("kendoChart").refresh();
    });

    $('input[name="periodTypeProd"]').on('change', function () {
        lgd.periodTypeProdChange();
    });

    $('input[name="periodTypeAvail"]').on('change', function () {
        lgd.periodTypeAvailChange();
    });
});