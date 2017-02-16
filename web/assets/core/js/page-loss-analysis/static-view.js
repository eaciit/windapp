'use strict';

viewModel.StaticView = new Object();
var sv = viewModel.StaticView;

sv.type = ko.observable();

sv.StaticView = function(){
    fa.LoadData();
    app.loading(true);
    if(pg.isFirstStaticView() === true){
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine,
            project: fa.project,
        };

        toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarylist", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            var gData = res.data.Data

            $('#lossGrid').html("");
            $('#lossGrid').kendoGrid({
                dataSource: {
                    data: gData,
                    pageSize: 10,
                    aggregate: [
                        { field: "Production", aggregate: "sum" },
                        { field: "LossEnergy", aggregate: "sum" },
                        { field: "MachineDownHours", aggregate: "sum" },
                        { field: "GridDownHours", aggregate: "sum" },
                        { field: "EnergyyMD", aggregate: "sum" },
                        { field: "EnergyyGD", aggregate: "sum" },
                        { field: "ElectricLoss", aggregate: "sum" },
                        { field: "PCDeviation", aggregate: "sum" },
                        { field: "Others", aggregate: "sum" },
                    ]
                },
                groupable: false,
                sortable: true,
                filterable: false,
                // height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
                height: 399,
                pageable: {
                    pageSize: 10,
                    input: true, 
                },
                columns: [
                    { title: sv.type,field: "Id",width: 100,attributes: {style: "text-align:center;"},headerAttributes: {style: "text-align:center;"},footerTemplate: "<center>Total (All Turbines)</center>"}, 
                    { title: "Production (MWh)", headerAttributes: { tyle: "text-align:center;"}, field: "Production",width: 100,attributes: { class: "align-center" },format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" }, 
                    { title: "Lost Energy (MWh)",headerAttributes: {style: "text-align:center;"},field: "LossEnergy", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"},
                    {
                        title: "Downtime : Duration (Hrs)",
                        headerAttributes: {
                            style: 'font-weight: bold; text-align: center;'
                        },
                        columns: [
                            {
                                title: "Machine",
                                headerAttributes: { style: "text-align:center;" },
                                field: "MachineDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                            {
                                title: "Grid",
                                headerAttributes: { style: "text-align:center;" },
                                field: "GridDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                            },
                        ]
                    }, {
                        title: "Downtime : Energy Loss (MWh)",
                        headerAttributes: {
                            style: 'font-weight: bold; text-align: center;'
                        },
                        columns: [
                            {
                                title: "Machine",
                                headerAttributes: { style: "text-align:center;" },
                                field: "EnergyyMD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                            {
                                title: "Grid",
                                headerAttributes: { style: "text-align:center;" },
                                field: "EnergyyGD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                            },
                        ]
                    }, {
                        title: "Electrical Losses (MWh)",
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "ElectricLoss", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }, {
                        title: "Power Curve Deviation (MW)", //Sepertinya ini MW
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "PCDeviation", width: 120, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }, {
                        title: "Others (MWh)", //Sepertinya ini KWh
                        headerAttributes: {
                            style: "text-align:center;"
                        },
                        field: "Others", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    }
                ],
                dataBound: function(){
                     app.loading(false);
                     pg.isFirstStaticView(false);
                }
            })
        });
        $('#availabledatestart').html(pg.availabledatestartscada());
        $('#availabledateend').html(pg.availabledateendscada());
    }else{
        $("#lossGrid").data("kendoGrid").refresh();
        $('#availabledatestart').html(pg.availabledatestartscada());
        $('#availabledateend').html(pg.availabledateendscada());
    }
}