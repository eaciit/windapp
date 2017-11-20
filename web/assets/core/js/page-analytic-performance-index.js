'use strict';

viewModel.AnalyticPerformanceIndex = {};
var page = viewModel.AnalyticPerformanceIndex;
page.summary = ko.observableArray();
vm.currentMenu('Performance Index');
vm.currentTitle('Performance Index');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Performance Index', href: viewModel.appName + 'page/analyticperformanceindex' }]);

var hideLast24 = false;
var hideLastWeek = false;
var hideMTD = false;
var hideYTD = false;

var Data = {
    LoadData: function () {
        app.loading(false);
        var isValid = fa.LoadData();
        if (isValid) {
            this.InitGrid();
        }
        // fa.getProjectInfo();
    },
    InitGrid : function(){

        var maxDateData = new Date(app.getUTCDate(app.currentDateData));
        var maxDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate(), 0, 0, 0, 0));
        var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
        var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
        var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), 1, 0, 0, 0, 0));
        var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
        var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 11, 31, 0, 0, 0, 0));
        if (fa.period == "annual") {
            fa.dateEnd = new Date(Date.UTC(moment(fa.dateEnd).get('year'), 11, 31, 0, 0, 0, 0));
        }
        
        if(fa.dateStart.getTime() === last24hours.getTime() && fa.dateEnd.getTime() === maxDate.getTime()) {
            hideLast24 = true;
        } else if(fa.dateStart.getTime() === lastweek.getTime() && fa.dateEnd.getTime() === maxDate.getTime()) {
            hideLastWeek = true;
        } else if(fa.dateStart.getTime() === startMonthDate.getTime() && fa.dateEnd.getTime() === maxDate.getTime()) {
            hideMTD = true;
        } else if(fa.dateStart.getTime() === startYearDate.getTime() && fa.dateEnd.getTime() === maxDate.getTime()) {
            hideYTD = true;
        } else if(fa.dateStart.getTime() === startYearDate.getTime() && fa.dateEnd.getTime() === endYearDate.getTime()) {
            hideYTD = true;
        }

        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));   

        var param = {
            period: fa.period,
            dateStart: dateStart,
            dateEnd: dateEnd,
            turbine: fa.turbine(),
            project: fa.project
        };

        $("#performance-grid").html(""); 
        // $("#performance-grid").kendoGrid().height($(".content-wrapper").height() - ($("#filter-analytic").height()+209));
        $("#performance-grid").kendoGrid({
            dataSource: {
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
                schema: {
                    data: function(res) {
                        if (res.data.Data == undefined) {
                            return;
                        }
                        var summary = res.data.Summary;

                        page.summary(summary);
                        return _.sortBy(res.data.Data, 'Turbine');
                    },
                    total: function (res) {
                        if (res.data.Data == undefined) {
                            return;
                        }

                        return res.data.Data.length;
                    },
                    model: {
                        fields: {
                            Project: { type: "string" },
                            Turbine: { type: "string" },
                            PerformanceIndex: { type: "number" },
                            PerformanceIndexLast24Hours: { type: "number" },
                            PerformanceIndexLastWeek: { type: "number" },
                            PerformanceIndexMTD: {type: "number"},
                            PerformanceIndexYTD: {type: "number"}
                        }
                    }
                },
                pageSize: 10,
                group: {
                    field: "Project", aggregates: [
                        { field: "Project", aggregate: "count" },
                        { field: "Turbine", aggregate: "count" },
                        { field: "PerformanceIndexLast24Hours", aggregate: "sum"},
                        { field: "PerformanceIndex", aggregate: "average" },
                        { field: "PerformanceIndexLastWeek", aggregate: "sum" }
                    ]
                },
                aggregate: [
                    { field: "Project", aggregate: "count" },
                    { field: "Turbine", aggregate: "count" },
                    { field: "PerformanceIndexLast24Hours", aggregate: "sum" },
                    { field: "PerformanceIndex", aggregate: "average" },
                    { field: "PerformanceIndexLastWeek", aggregate: "sum" }
                ]
            },
            scrollable: true,
            sortable: true,
            pageable: {
                pageSize: 10,
                input: true, 
            },
            columns: [
                { width: 150, locked: true, field: "Turbine", title: "Turbine" },
                {
                    title: "Performance<br>"+kendo.toString(moment.utc(fa.dateStart).format('DD-MM-YYYY '))+" to " + kendo.toString(moment.utc(fa.dateEnd).format('DD-MM-YYYY')) ,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { width: 110, field: "PerformanceIndex", title: "PI",headerAttributes: { style: "text-align: center" },format: "{0:n2}",attributes:{ style: "text-align: center" },},
                        { field: "PotentialPower", title: "Pot. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "Power", title: "Act. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "ProductionIndex", title: "Production Index", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    hidden: hideLast24,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "PI", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" }, format: "{0:n2}"},
                        { field: "PotentialPowerLast24Hours", title: "Pot. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Act. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "ProductionIndexLast24Hours", title: "Production Index", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    hidden: hideLastWeek,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "PI", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerLastWeek", title: "Pot. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Act. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "ProductionIndexLastWeek", title: "Production Index", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    hidden: hideMTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "PI", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerMTD", title: "Pot. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Act. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "ProductionIndexMTD", title: "Production Index", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    hidden: hideYTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "PI", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerYTD", title: "Pot. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Act. Power (MW)", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "ProductionIndexYTD", title: "Production Index", width: 110,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
            ],
            dataBound: function(){
                Data.InitHeader(hideLast24, hideLastWeek, hideMTD, hideYTD); 

            }
        });

    },
    InitHeader: function(hideLast24, hideLastWeek, hideMTD, hideYTD){
            $("#performance-grid").find(".k-grid-footer").remove();
            var contentProject   = $("#performance-grid").find(".k-grid-content-locked").find(".k-grouping-row>").find("p");
            var projects = []
            contentProject.each(function(i, obj){
                var project = $(obj).text().replace("Project: ", "");
                projects.push(project);
                
            });
            $.each(page.summary(), function(a, prj){
                if ( $.inArray(prj.Project,  projects) > -1 ){
                    var elem = $("#performance-grid").find($(".k-auto-scrollable")).find(".k-grouping-row")[a];

                    $(elem).find("td").remove();
                    $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PerformanceIndex, "n2")+'</td>');
                    $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PotentialPower, "n2")+'</td>');
                    $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.Power, "n2")+'</td>');
                    $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.ProductionIndex, "n2")+'</td>');

                    if(hideLast24 == false){
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PerformanceIndexLast24Hours, "n2")+'</td>');
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PotentialPowerLast24Hours, "n2")+'</td>'); 
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PowerLast24Hours, "n2")+'</td>');
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.ProductionIndexLast24Hours, "n2")+'</td>');                          
                    }


                    if(hideLastWeek == false){
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PerformanceIndexLastWeek, "n2")+'</td>');
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PotentialPowerLastWeek, "n2")+'</td>'); 
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PowerLastWeek, "n2")+'</td>');  
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.ProductionIndexLastWeek, "n2")+'</td>');  
                    }
     

                    if(hideMTD == false){
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PerformanceIndexMTD, "n2")+'</td>');
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PotentialPowerMTD, "n2")+'</td>'); 
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PowerMTD, "n2")+'</td>');  
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.ProductionIndexMTD, "n2")+'</td>');  
                    }
     
                    if(hideYTD == false){
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PerformanceIndexYTD, "n2")+'</td>');
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PotentialPowerYTD, "n2")+'</td>'); 
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.PowerYTD, "n2")+'</td>'); 
                        $(elem).append('<td aria-expanded="true" style="text-align:center">'+kendo.toString(prj.ProductionIndexYTD, "n2")+'</td>');  
                    }

                } else {
                    console.log("data not found");
                }
            })
    }
}
 

$(function(){
    di.getAvailDate();

    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        Data.LoadData();
    });

    $('#projectList').kendoDropDownList({
        change: function () { 
            setTimeout(function(){
                fa.currentFilter().project = this._old;
                fa.checkFilter();
                di.getAvailDate();
                var project = $('#projectList').data("kendoDropDownList").value();
                fa.populateTurbine(project);
            },500); 
        }
    });

    $( window ).resize(function() {
      vm.adjustLayout()
      $("#performance-grid").data("kendoGrid").resize();
    });
    setTimeout(function(){
        Data.LoadData();
        vm.adjustLayout();
    },500);
    
})  