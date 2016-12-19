'use strict';

// var monthNames = ["January", "February", "March", "April", "May", "June",
//     "July", "August", "September", "October", "November", "December"
// ];

viewModel.production = {};
var prod = viewModel.production;

prod.loadData = function () {
    if (lgd.isProduction()) {
        $.when(prod.periodTypeProdChange()).done(function(){
            setTimeout(function(){
                app.loading(false);
            },500);
        })
    }
};

prod.periodTypeProdChange = function () {
    prod.gridProduction($("#projectId").data("kendoDropDownList").value(), maxdate);
}


prod.gridProduction = function (project, enddate) {
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

$(function () {
    $('input[name="periodTypeProd"]').on('change', function () {
        prod.periodTypeProdChange();
    });
});