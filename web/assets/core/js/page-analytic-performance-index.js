'use strict';

viewModel.AnalyticPerformanceIndex = {};
var page = viewModel.AnalyticPerformanceIndex;

vm.currentMenu('Performance Index');
vm.currentTitle('Performance Index');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Performance Index', href: viewModel.appName + 'page/analyticperformanceindex' }]);


page.data = ko.observableArray([
    {
        "Project":"Tejuva",
        "PerformanceIndexLast24Hours":82.26562913931056,
        "PerformanceIndexLastWeek":0,
        "PerformanceIndexMTD":0,
        "PerformanceIndexYTD":0,
        "PerformanceIndexDateRange":"",
        "ProductionLast24Hours":566784.266667,
        "ProductionLastWeek":0,
        "ProductionMTD":0,
        "ProductionYTD":0,
        "ProductionDateRange":"",
        "PowerLast24Hours":3400705.6,
        "PowerLastWeek":0,
        "PowerMTD":0,
        "PowerYTD":0,
        "PowerDateRange":""
    },{
        "Project":"Tejuva",
        "PerformanceIndexLast24Hours":82.26562913931056,
        "PerformanceIndexLastWeek":0,
        "PerformanceIndexMTD":0,
        "PerformanceIndexYTD":0,
        "PerformanceIndexDateRange":"",
        "ProductionLast24Hours":566784.266667,
        "ProductionLastWeek":0,
        "ProductionMTD":0,
        "ProductionYTD":0,
        "ProductionDateRange":"",
        "PowerLast24Hours":3400705.6,
        "PowerLastWeek":0,
        "PowerMTD":0,
        "PowerYTD":0,
        "PowerDateRange":""
    }])

var Data = {
    LoadData: function () {
        fa.LoadData();
        fa.getProjectInfo();
        app.loading(false);
        this.InitGrid();
    },
    InitGrid : function(){
        $("#performance-grid").html("");
        $("#performance-grid").kendoGrid({
            dataSource    : {
                data: page.data()
            },
            scrollable: true,
            sortable: true,
            pageable: true,
            height        : 430,
            detailInit : Data.InitGridDetail,
            dataBound: function () {
                this.expandRow(this.tbody.find("tr.k-master-row").first());
            },
            columns       : [
                { field: "Project", title: "Name", width: 110,headerAttributes: { style: "text-align: center" },attributes: { style: "padding-left: 20px" }},
                {
                    title: "Performance",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "Name", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                        { field: "Name", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                        { field: "Name", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" }, format: "{0:n2}"},
                        { field: "ProductionLast24Hours", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionLastWeek", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionMTD", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionYTD", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                { field: "CustomRangeSelected", title: "Custom Range Selected", width: 200,headerAttributes: { style: "text-align: center" }},
            ]
        });
    },
    InitGridDetail: function (e) {
        $("<div/>").appendTo(e.detailCell).kendoGrid({
            selectable: "multiple",
            dataSource: {
                data: page.data()
            },
            scrollable: false,
            sortable: false,
            pageable: false,
             columns       : [
                { field: "Project", title: "Name", width: 100,headerAttributes: { style: "text-align: center" },attributes: { style: "padding-left: 30px" }},
                {
                    title: "Performance",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "Name", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                        { field: "Name", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                        { field: "Name", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right"}, format: "{0:n2}"},
                        { field: "ProductionLast24Hours", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionLastWeek", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionMTD", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "Perf. Index", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionYTD", title: "Prod", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Act. Power", width: 90,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                { field: "CustomRangeSelected", title: "Custom Range Selected", width: 200,headerAttributes: { style: "text-align: center" }},
                
            ]
        });
        $(".k-grid tbody .k-grid .k-grid-header").hide();
    },
}
 

$(function(){
    setTimeout(function(){
        Data.LoadData();
    },100);
    
})  