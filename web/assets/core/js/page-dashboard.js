'use strict';

viewModel.landing = {};
var lgd = viewModel.landing;

var monthNames = ["January", "February", "March", "April", "May", "June",
    "July", "August", "September", "October", "November", "December"
];

lgd.filter = ko.observableArray([
    { text: "Fleet", value: "1" },
    { text: "WindFarm-01", value: "2" },
    { text: "WindFarm-02", value: "3" }
]);

lgd.isFirst = ko.observable(true);

lgd.isSummary = ko.observable(false);
lgd.isProduction = ko.observable(false);
lgd.isAvailability = ko.observable(false);

lgd.isDetailProd = ko.observable(false);
lgd.isDetailProdByProject = ko.observable(false);

lgd.isDetailAvailLostEnergy = ko.observable(false);
lgd.detailDTLostEnergyTxt = ko.observable();

lgd.isDetailDTTop = ko.observable(false);
lgd.detailDTTopTxt = ko.observable();

lgd.projectList = ko.observableArray([]);
lgd.projectItem = ko.observableArray([]);
lgd.mdTypeList = ko.observableArray([]);
lgd.projectName = ko.observable();
lgd.isFleet = ko.observable(true);
lgd.isNonFleet = ko.observable(true);
lgd.FleetDTLEDownType = ko.observable();
lgd.LEFleetByDown = ko.observable(false);

lgd.prodDateRangeStr = ko.observable('');
lgd.projectAvailList = ko.observableArray([]);
lgd.projectAvailSelected = ko.observable('');

var lastParam = {};
var lastParamChart = {};
var lastParamLevel2 = {};
var dtType = '';
var monthDetailDT = '';
var projectSelected = '';
var projectSelectedLevel2 = '';
var maxDateData = new Date(app.getUTCDate(app.currentDateData));
var maxdate = maxDateData;
// var intervalMap = setInterval(function(){ sum.indiaMap(lgd.projectName())}, 36000);
var mapIndia;


lgd.populateProject = function (data) {
    if(data.length > 0) {
        var datavalue = [{ "value": "Fleet", "text": "Fleet" }];
        var projectAvail = {};
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val.Value;
                data.text = val.Name;
                datavalue.push(data);
                projectAvail = {text: val.Value, value: val.Value};
                lgd.projectAvailList.push(projectAvail);
                if(key===0) {
                    lgd.projectAvailSelected(val.Value);
                }
            });
        }
        lgd.projectList(datavalue);
        lgd.projectItem(datavalue);
    }
};

lgd.LoadData = function () {
    lgd.stop();
    app.loading(true);
    sum.scadaLastUpdate();

    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { ProjectName: project, Date: maxdate };

    if (project == "Fleet") {
        lgd.isFleet(true);
        lgd.isNonFleet(false);
        $("#div-windiness").hide();
        $("#div-winddistribution").show();
    } else {
        lgd.isFleet(false);
        lgd.isNonFleet(true);
        $("#div-windiness").show();
        $("#div-winddistribution").hide();
    }

    setTimeout(function () {
        prod.loadData();
        avail.loadData();
        sum.loadData();
    }, 600);
}


lgd.start = function() {  // use a one-off timer
    mapIndia =  setInterval(function() {
       var project =  $("#projectId").data("kendoDropDownList").value();
       sum.indiaMap(project);
    }, 5000);
};

lgd.stop = function(){
    clearInterval(mapIndia);
    return false;
};

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

vm.currentMenu('Dashboard');
vm.currentTitle('Dashboard');
vm.isDashboard(true);
vm.breadcrumb([{ title: 'Dashboard', href: viewModel.appName + 'page/landing' }, { title: 'Home', href: '#' }]);

$(function () {
    $(".prodToTable").on("change", function(){
        if($("#chartProduction").data("kendoGrid") != undefined){
            $("#chartProduction thead [data-field='category']").html("");
        }
    });

    console.log("First call!");

    $("#tabSummary").on("click", function () {
        // intervalMap = setInterval(function(){ sum.indiaMap(lgd.projectName())}, 4000);
        lgd.isSummary(true);
        lgd.isProduction(false);
        lgd.isAvailability(false);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $("#tabProduction").on("click", function () {
        lgd.stop();
        lgd.isSummary(false);
        lgd.isProduction(true);
        lgd.isAvailability(false);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $("#tabAvailability").on("click", function () {
        lgd.stop();
        lgd.isSummary(false);
        lgd.isProduction(false);
        lgd.isAvailability(true);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $('input[name="periodTypeAvail"]').on('change', function () {
        lgd.periodTypeAvailChange();
    });    

    console.log("Second call!");

    setTimeout(function(){
        lgd.isSummary(true);
        lgd.isProduction(false);
        lgd.isAvailability(false);
        lgd.projectName("Fleet");

        lgd.LoadData();
        google.maps.event.addDomListener(window, 'load', sum.initialize());

    },500);

    console.log("End call!");

});

// temporary to fired summary number left side map
$(document).ready(function() {
    setTimeout(function(){
        sum.loadData();
        console.log("Data call then!");
    },5000);
});