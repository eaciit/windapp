'use strict';

viewModel.AnalyticPerformanceIndex = {};
var page = viewModel.AnalyticPerformanceIndex;

vm.currentMenu('Performance Index');
vm.currentTitle('Performance Index');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Performance Index', href: viewModel.appName + 'page/analyticperformanceindex' }]);

var Data = {
    LoadData: function () {
        fa.LoadData();
        fa.getProjectInfo();
        app.loading(false);
        this.InitGrid();
    },
    InitGrid : function(){

        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine,
            project: fa.project
        };

        $("#performance-grid").html("");
        $("#performance-grid").kendoGrid({
            dataSource: {
                serverSorting: true,
                serverFiltering: true,
                transport: {
                    read: {
                        url: viewModel.appName + "analyticperformanceindex/getperformanceindex",
                        type: "POST",
                        data: param,
                        dataType: "json",
                        contentType: "application/json; charset=utf-8"
                    },
                    parameterMap: function(options) {
                        return JSON.stringify(options);
                    }
                },
                pageSize: 10,
                schema: {
                    data: function(res) {
                        app.loading(false);
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data
                    },
                    total: function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        return res.data.Data.length;
                    }
                },
            },
            scrollable: true,
            sortable: true,
            pageable: true,
            height        : 550,
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
                        { field: "PerformanceIndex", title: "Perf. Index <br> (%)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "Production", title: "Prod <br>(MWh)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "Power", title: "Act. Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "Perf. Index <br> (%)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" }, format: "{0:n2}"},
                        { field: "ProductionLast24Hours", title: "Prod <br>(MWh)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Act. Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "Perf. Index <br> (%)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionLastWeek", title: "Prod <br>(MWh)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Act. Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "Perf. Index <br> (%)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionMTD", title: "Prod <br>(MWh)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Act. Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "Perf. Index <br> (%)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionYTD", title: "Prod<br>(MWh)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Act. Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                { field: "StartDate", title: "Custom Range Selected", width: 200,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },template: "#= kendo.toString(moment.utc(StartDate).format('DD-MM-YYYY'), 'HH:mm:ss') # until  #= kendo.toString(moment.utc(EndDate).format('DD-MM-YYYY'), 'HH:mm:ss')#"},
            ]
        });
    },
    InitGridDetail: function (e) {
        $("<div/>").appendTo(e.detailCell).kendoGrid({
            selectable: "multiple",
            dataSource: {
                data: e.data.Details,
                sort: [
                    { field: 'Turbine', dir: 'asc' }
                ],
                 // pageSize: 10,
            },
            scrollable: false,
            sortable: false,
            pageable: false,
             columns       : [
                { field: "Turbine", title: "Name", width: 100,headerAttributes: { style: "text-align: center" },attributes: { style: "padding-left: 30px" }},
                {
                    title: "Performance",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndex", title: "Perf. Index", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "Production", title: "Prod", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "Power", title: "Act. Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "Perf. Index", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right"}, format: "{0:n2}"},
                        { field: "ProductionLast24Hours", title: "Prod", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Act. Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "Perf. Index", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionLastWeek", title: "Prod", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Act. Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "Perf. Index", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionMTD", title: "Prod", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Act. Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "Perf. Index", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "ProductionYTD", title: "Prod", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Act. Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: right" },format: "{0:n2}"},
                    ]
                },
                { field: "StartDate", title: "Custom Range Selected", width: 200,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },template: "#= kendo.toString(moment.utc(StartDate).format('DD-MM-YYYY'), 'HH:mm:ss') # until  #= kendo.toString(moment.utc(EndDate).format('DD-MM-YYYY'), 'HH:mm:ss')#"},
                
            ]
        });
        $(".k-grid tbody .k-grid .k-grid-header").hide();
    },
}
 

$(function(){
    $('#btnRefresh').on('click', function () {
        Data.LoadData();
    });
    setTimeout(function(){
        Data.LoadData();
    },100);
    
})  