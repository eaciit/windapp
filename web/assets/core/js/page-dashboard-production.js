'use strict';

// var monthNames = ["January", "February", "March", "April", "May", "June",
//     "July", "August", "September", "October", "November", "December"
// ];

viewModel.production = {};
var prod = viewModel.production;

prod.loadData = function () {
    if (lgd.isProduction()) {
        prod.periodTypeProdChange();
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
        startDate = new Date(Date.UTC(endDateYear, endDateMonth, endDateDate - 7, 0, 0, 0, 0));
    } else if (type == "mtd") {
        startDate = new Date(Date.UTC(endDateYear, endDateMonth, 1, 0, 0, 0, 0));
    } else if (type == "ytd") {
        startDate = new Date(Date.UTC(endDateYear, 0, 1, 0, 0, 0, 0));
    }

    filters.push({ field: "dateinfo.dateid", operator: "gte", value: startDate });
    filters.push({ field: "dateinfo.dateid", operator: "lte", value: enddate });
    filters.push({ field: "type", operator: "eq", value: type });
    if (project != "Fleet") {
        filters.push({ field: "projectname", operator: "eq", value: project });
    }
    method = "getsummarydatadaily";

    var startDateStr = startDate.getUTCDate() + "-" + monthNames[startDate.getUTCMonth()] + "-" + startDate.getUTCFullYear();
    var endDateStr = enddate.getUTCDate() + "-" + monthNames[enddate.getUTCMonth()] + "-" + enddate.getUTCFullYear();

    $('#prodDateRangeStr').html(startDateStr + " to " + endDateStr);

    var filter = { filters: filters }
    var param = { filter: filter };

    var allcolumn = [
                { title: "Project Name", width:100, field: "name", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "No. of WTG", width:100,field: "noofwtg", format: "{0:n0}", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Production<br>(GWh)", width:100,field: "production", footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum/1000000, 'n2')#</div>", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "PLF<br>(%)", width:100,field: "plf", footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(plf*100, 'n2') #%" },
                { title: "Total Availability<br>(%)", width:100,footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", field: "totalavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(totalavail*100, 'n2') #%" },
                { title: "Lowest Machine Availability<br>(%)", width:100,field: "lowestmachineavail", headerAttributes: { style: "text-align: center" }, 
                    attributes: { class: "align-center" } ,
                    template : function(e){
                        if(e.lowestmachineavail == "-"){
                            return "<i class='fa fa-star txt-green'></i>"
                        }else{
                            return e.lowestmachineavail;
                        }
                    } 
                },
                { title: "Lowest PLF<br>(%)", width:100,field: "lowestplf", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Max. Lost Energy due to Downtime<br>(KWh)", width:100,field: "maxlossenergy", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Data Availability<br>(%)", width:100,footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", field: "dataavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(dataavail*100, 'n2') #%" },
            ]
    var allaggr = [
                    { field: "production", aggregate: "sum" },
                    { field: "plf", aggregate: "average" },
                    { field: "totalavail", aggregate: "average" },
                    { field: "dataavail", aggregate: "average" },
                ]

    if (project != "Fleet") {
        allcolumn = [
                { title: "Turbine Name", width:100, field: "name", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "Production<br>(GWh)", width:100,field: "production", footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum/1000000, 'n2')#</div>", template: "#= kendo.toString(production/1000000, 'n2') #", headerAttributes: { style: "text-align:center;" }, attributes: { style: "text-align:center;" } },
                { title: "PLF<br>(%)", width:100,field: "plf", footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(plf*100, 'n2') #%" },
                { title: "Total Availability<br>(%)", width:100,footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", field: "totalavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(totalavail*100, 'n2') #%" },
                { title: "Machine Availability<br>(%)", width:100,field: "machineavailfloat", footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(machineavailfloat*100, 'n2') #%" },
                { title: "Lost Energy due to Downtime<br>(KWh)", width:100,field: "maxlossenergy", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" } },
                { title: "Data Availability<br>(%)", width:100,footerTemplate: "<div style='text-align:center'>#=kendo.toString(average*100, 'n2')#%</div>", field: "dataavail", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#= kendo.toString(dataavail*100, 'n2') #%" },
            ]

        allaggr.push({ field: "machineavailfloat", aggregate: "average" })
    }

    var reqProdData = toolkit.ajaxPost(viewModel.appName + "dashboard/" + method, param, function (res) {
        $('#productionGrid').html("");
        $("#productionGrid").kendoGrid({
            dataSource: {
                data: res.data.Data,
                pageSize: 10,
                aggregate: allaggr
            /*,
            sort: [
                { field: '_id', dir: 'asc' },
            ],*/
            },
            groupable: false,
            /*serverPaging: true,
            serverSorting: true,*/
            pageable: {
                pageSize: 10,
                input: true, 
            },
            columns: allcolumn,
        });
    });

    $.when(reqProdData).done(function() {
        app.loading(false);
    })
}

$(function () {
    $('input[name="periodTypeProd"]').on('change', function () {
        prod.periodTypeProdChange();
    });
});