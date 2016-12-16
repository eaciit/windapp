'use strict';

viewModel.AnalyticPerformanceIndex = {};
var page = viewModel.AnalyticPerformanceIndex;

vm.currentMenu('Performance Index');
vm.currentTitle('Performance Index');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Performance Index', href: viewModel.appName + 'page/analyticperformanceindex' }]);

var Data = {
    LoadData: function () {
        fa.LoadData();
        fa.getProjectInfo();
        app.loading(false);
        this.InitGrid();
    },
    InitGrid : function(){
        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
            var minDatetemp = new Date(res.ScadaData[0]);
            var maxDatetemp = new Date(res.ScadaData[1]);
            $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
            $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
        })
        var hideLast24 = false;
        var hideLastWeek = false;
        var hideMTD = false;
        var hideYTD = false;
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
            height        : $(".content-wrapper").height() * 0.48,
            detailInit : Data.InitGridDetail,
            dataBound: function () {
                this.expandRow(this.tbody.find("tr.k-master-row").first());
            },
            columns       : [
                { field: "Project", title: "Name", width: 110,headerAttributes: { style: "text-align: center" },attributes: { style: "padding-left: 20px;border-right:1px solid rgba(128, 128, 128, 0.26)" }},
                {
                    title: "Performance<br>"+kendo.toString(moment.utc(fa.dateStart).format('DD-MM-YYYY '))+" to " + kendo.toString(moment.utc(fa.dateEnd).format('DD-MM-YYYY')) ,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndex", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPower", title: "Potential Power <br>(MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "Power", title: "Actual Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    hidden: hideLast24,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" }, format: "{0:n2}"},
                        { field: "PotentialPowerLast24Hours", title: "Potential Power <br>(MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Actual Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    hidden: hideLastWeek,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerLastWeek", title: "Potential Power <br>(MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Actual Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    hidden: hideMTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerMTD", title: "Potential Power <br>(MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Actual Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    hidden: hideYTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerYTD", title: "Potential Power<br>(MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Actual Power <br> (MW)", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
            ]
        });
    },
    InitGridDetail: function (e) {
        var hideLast24 = false;
        var hideLastWeek = false;
        var hideMTD = false;
        var hideYTD = false;
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
                { field: "Turbine", title: "Name", width: 100,headerAttributes: { style: "text-align: center" },attributes: { style: "padding-left: 30px;border-right:1px solid rgba(128, 128, 128, 0.26)" }},
                {
                    title: "Performance",
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndex", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPower", title: "Potential Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "Power", title: "Actual Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last 24 Hours",
                    hidden: hideLast24,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLast24Hours", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center"}, format: "{0:n2}"},
                        { field: "PotentialPowerLast24Hours", title: "Potential Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLast24Hours", title: "Actual Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "Last Week",
                    hidden: hideLastWeek,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexLastWeek", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerLastWeek", title: "Potential Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerLastWeek", title: "Actual Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "MTD",
                    hidden: hideMTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexMTD", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerMTD", title: "Potential Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerMTD", title: "Actual Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
                {
                    title: "YTD",
                    hidden: hideYTD,
                    headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    columns: [
                        { field: "PerformanceIndexYTD", title: "PI", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PotentialPowerYTD", title: "Potential Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center" },format: "{0:n2}"},
                        { field: "PowerYTD", title: "Actual Power", width: 100,headerAttributes: { style: "text-align: center" },attributes:{ style: "text-align: center;border-right:1px solid rgba(128, 128, 128, 0.26)" },format: "{0:n2}"},
                    ]
                },
            ]
        });
        $(".k-grid tbody .k-grid .k-grid-header").hide();
    },
}
 

$(function(){
    $('#btnRefresh').on('click', function () {
        Data.LoadData();
    });

    $( window ).resize(function() {
      $('#performance-grid').data('kendoGrid').refresh();
    });

    setTimeout(function(){
        Data.LoadData();
    },100);
    
})  