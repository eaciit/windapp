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
sum.noOfTurbines = ko.observable();
sum.totalMaxCapacity = ko.observable();
sum.currentDown = ko.observable();
sum.twoDaysDown = ko.observable();
sum.dataSource = ko.observable();
sum.dataSourceScada = ko.observable();

vm.dateAsOf(app.currentDateData);
sum.loadData = function () {
    if (lgd.isSummary()) {
        var project = $("#projectId").data("kendoDropDownList").value();
        var param = { ProjectName: project, Date: maxdate };

        var ajax1 = toolkit.ajaxPost(viewModel.appName + "dashboard/getscadalastupdate", param, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }

            sum.dataSource(res.data[0]);

            sum.noOfProjects(res.data[0].NoOfProjects);
            sum.noOfTurbines(res.data[0].NoOfTurbines);
            sum.totalMaxCapacity(res.data[0].TotalMaxCapacity / 1000);
            sum.currentDown(res.data[0].CurrentDown);
            sum.twoDaysDown(res.data[0].TwoDaysDown);

            var lastUpdate = new Date(res.data[0].LastUpdate);

            // vm.dateAsOf(lastUpdate.addHours(-7));
            sum.ProductionChart(res.data[0].Productions);
            sum.CumProduction(res.data[0].CummulativeProductions);
            sum.SummaryData(project);
        });

        var ajax2 = toolkit.ajaxPost(viewModel.appName + "dashboard/getscadasummarybymonth", param, function (res) {
            if (!toolkit.isFine(res)) {
                return;
            }
            sum.dataSourceScada(res.data);
            sum.PLF(res.data);
            sum.LostEnergy(res.data);
            sum.Windiness(res.data);
            sum.ProdMonth(res.data);
            sum.AvailabilityChart(res.data);
            sum.ProdCurLast(res.data);
            sum.indiaMap(project);
            sum.isDetailProd(false);
            sum.isDetailProdByProject(false);
        });



        $.when(ajax1, ajax2).done(function(){
            setTimeout(function(){
                app.loading(false);
            },200);        
        })
    }
};

sum.SummaryData = function (project) {
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
        pageable: {
            pageSize: 2,
            input: true, 
        },
        columns: [
            { title: "Project Name", field: "name", headerAttributes: { style: "text-align:left;" }, attributes: { style: "text-align:left;" } },
            { title: "No. of WTG", field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Production<br>(GWh)", field: "production", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "PLF<br>(%)", field: "plf", width: 80, format: "{0:n2}", template: "#= kendo.toString(plf*100, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Lost Energy<br>(MWh)", field: "lostenergy", template: "#= kendo.toString(lostenergy/1000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
            { title: "Downtime<br>(Hours)", field: "downtimehours", format: "{0:n2}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
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
}

sum.PLF = function (dataSource) {
    $("#chartPLF").replaceWith('<div id="chartPLF"></div>');
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
    });
}

sum.ProdMonth = function (dataSource) {
    $("#chartProdMonth").replaceWith('<div id="chartProdMonth"></div>');
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

sum.AvailabilityChart = function (dataSource) {
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
            name: "DBA",
            field: "ScadaAvail",
            // opacity : 0.5,
            color: "#21c4af"
        }, {
            name: "TBA",
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

sum.ProdCurLast = function (dataSource) {
    $("#chartCurrLast").replaceWith('<div id="chartCurrLast"></div>');
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

sum.indiaMap = function (project) {
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

sum.ProductionChart = function (dataSource) {
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


sum.DetailProd = function (e) {
    var bulan = e.category;
    sum.detailProdTxt(bulan);
    vm.isDashboard(false);
    lgd.isSummary(false);
    sum.isDetailProd(true);
    sum.isDetailProdByProject(false);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { 'project': project, 'date': bulan };

    toolkit.ajaxPost(viewModel.appName + "dashboard/getdetailprod", param, function (res) {
        if (!toolkit.isFine(res)) {
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
                sum.DetailProdByProject(e, bulan, dataSource);
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