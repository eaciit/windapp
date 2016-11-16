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

lgd.projectList = ko.observableArray([{ "value": "Fleet", "text": "Fleet" }]);
lgd.projectItem = ko.observableArray([]);
lgd.mdTypeList = ko.observableArray([]);
lgd.projectName = ko.observable();
lgd.isFleet = ko.observable(true);
lgd.isNonFleet = ko.observable(true);
lgd.FleetDTLEDownType = ko.observable();
lgd.LEFleetByDown = ko.observable(false);

lgd.prodDateRangeStr = ko.observable('');

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

lgd.LoadData = function () {
    app.loading(true);
    var project = $("#projectId").data("kendoDropDownList").value();
    var param = { ProjectName: project, Date: maxdate };

    if (project == "Fleet") {
        lgd.isFleet(true);
        lgd.isNonFleet(false);
        $("#div-windiness").hide();
    } else {
        lgd.isFleet(false);
        lgd.isNonFleet(true);
        $("#div-windiness").show();
    }

    setTimeout(function () {
        sum.loadData();
        prod.loadData();
        avail.loadData();
    }, 600);
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

vm.currentMenu('Dashboard');
vm.currentTitle('Dashboard');
vm.isDashboard(true);
vm.breadcrumb([{ title: 'Dashboard', href: viewModel.appName + 'page/landing' }, { title: 'Home', href: '#' }]);

$(function () {
    lgd.isSummary(true);
    lgd.isProduction(false);
    lgd.isAvailability(false);
    lgd.getProjectList();
    lgd.projectName("Fleet");

    lgd.LoadData();

    $("#tabSummary").on("click", function () {
        lgd.isSummary(true);
        lgd.isProduction(false);
        lgd.isAvailability(false);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $("#tabProduction").on("click", function () {
        lgd.isSummary(false);
        lgd.isProduction(true);
        lgd.isAvailability(false);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $("#tabAvailability").on("click", function () {
        lgd.isSummary(false);
        lgd.isProduction(false);
        lgd.isAvailability(true);
        lgd.LoadData();
        lgd.isFirst(false);
    });

    $('input[name="periodTypeAvail"]').on('change', function () {
        lgd.periodTypeAvailChange();
    });
});