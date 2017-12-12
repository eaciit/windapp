'use strict';

viewModel.StaticView = new Object();
var sv = viewModel.StaticView;

sv.type = ko.observable();
sv.isGrid = ko.observable(true);


sv.refreshView = function(view){
    // app.loading(true);
    setTimeout(function(){
        if(view == "gridStatic"){
            $('#lossGrid').data("kendoGrid").refresh();
        }else{
             // $("#staticProductionChart").data("kendoChart").refresh();
             // $("#staticDowntimeChart").data("kendoChart").refresh();
             // $("#staticLossChart").data("kendoChart").refresh();
        }
        // app.loading(false);
    },500);
}

sv.getChartView = function(gData){
    if (fa.project==""){
        sv.type = "Projects";
    }else{
        sv.type = "Turbines";
    }

    var rotation = 0
    if (gData.length >= 25) {
        rotation = 30
    }

    $("#staticProductionChart").kendoChart({
        theme:"flat",
        dataSource: {
            data: gData,
            sort:{
                field: "Id",
                dir: "asc"
            }
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea :{
            height: 250, 
            margin : 0,
            padding: 0,
            width: (screen.width * 0.89),
            background: "transparent",
        },
        plotArea: { margin: 0, padding: 0, height: 200, width: (screen.width * 0.89) },
        series: [{
            field: "Production",
            name: "Production",
        }],
        seriesColors: colorField,
        valueAxes: [{
            visible: true,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            title: { text: "Production (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            field: "Id",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                rotation : "auto",
            },
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: "#= series.name # at #= category # : #= kendo.toString(value, 'n2')# MWh",
            shared: false,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });

    $("#staticDowntimeChart").kendoChart({
        theme:"flat",
        dataSource: {
            data: gData,
            sort:{
                field: "Id",
                dir: "asc"
            }
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea :{
            height: 250, 
            margin : 0,
            padding: 0,
            background: "transparent",
            width: (screen.width * 0.89)
        },
        plotArea: { margin: 0, padding: 0, height: 200, width: (screen.width * 0.89) },
        seriesDefaults:{
            stack: true,
        },
        series: [{
            field: "MachineDownHours",
            name: "Machine",
            color:"#21c4af",
        },{
            field: "GridDownHours",
            name: "Grid",
            color: "#0097a4"
        }],
        // seriesColors: colorField,
        valueAxes: [{
            visible: true,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            title: { text: "Downtime Duration (Hrs)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            field: "Id",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                rotation : "auto",
            },
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: "#= series.name # at #= category # : #= kendo.toString(value, 'n2')# Hrs",
            shared: false,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });

    $("#staticLossChart").kendoChart({
        theme:"flat",
        dataSource: {
            data: gData,
            sort:{
                field: "Id",
                dir: "asc"
            }
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea :{
            height: 250, 
            margin : 0,
            padding: 0,
            background: "transparent",
            width: (screen.width * 0.89),
        },
        plotArea: { margin: 0, padding: 0, height: 200, width: (screen.width * 0.89) },
        seriesDefaults:{
            stack: true,
        },
        series: [{
            field: "EnergyyMD",
            name: "Machine",
            color:"#4589b0",
        },{
            field: "EnergyyGD",
            name: "Grid",
            color:"#80deea",
        }],
        // seriesColors: colorField,
        valueAxes: [{
            visible: true,
            line: {
                visible: false
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            name: "availline",
            title: { text: "Energy Loss (MWh)", font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif' },
        }],
        categoryAxis: {
            field: "Id",
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                rotation : "auto",
            },
            justified: true,
            majorGridLines: {
                visible: false
            },
        },
        tooltip: {
            visible: true,
            template: "#= series.name # at #= category # : #= kendo.toString(value, 'n2')# MWh",
            shared: false,
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
    });

    app.loading(false);
}

sv.getGridView = function(gData){
    $('#lossGrid').html("");
    $('#lossGrid').kendoGrid({
        dataSource: {
            data: gData,
            aggregate: [
                { field: "Production", aggregate: "sum" },
                { field: "LossEnergy", aggregate: "sum" },
                { field: "MachineDownHours", aggregate: "sum" },
                { field: "GridDownHours", aggregate: "sum" },
                { field: "EnergyyMD", aggregate: "sum" },
                { field: "EnergyyGD", aggregate: "sum" },
                { field: "EnergyyOD", aggregate: "sum" },
                { field: "ElectricLoss", aggregate: "sum" },
                { field: "PCDeviation", aggregate: "sum" },
                { field: "OtherDownHours", aggregate: "sum" },
                { field: "TotalAvail", aggregate: "average" },
                { field: "OKTime", aggregate: "sum" },
            ],
            serverPaging: false,
            sort: [ { field: "Id", dir: "asc" }],
        },
        scrollable: true,
        groupable: false,
        sortable: true,
        filterable: false,
        // height: $(".content-wrapper").height() - ($("#filter-analytic").height()+209),
        height: 350,
        pageable: false,
        // pageable: {
        //     pageSize: 24,
        //     input: true, 
        // },
        columns: [
            { title: sv.type,field: "Id",width: 130,attributes: {style: "text-align:center;"},headerAttributes: {style: "text-align:center;"},footerTemplate: "<center>Total</center>"}, 
            { title: "Production (MWh)", headerAttributes: { tyle: "text-align:center;"}, field: "Production",width: 100,attributes: { class: "align-center" },format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" }, 
            { title: "Lost Energy (MWh)",headerAttributes: {style: "text-align:center;"},field: "LossEnergy", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"},
            {
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
                    {
                        title: "Others",
                        headerAttributes: { style: "text-align:center;" },
                        field: "EnergyyOD", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>" 
                    },
                ]
            },
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
                    {
                        title: "Others",
                        headerAttributes: { style: "text-align:center;" },
                        field: "OtherDownHours", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    },
                ]
            }, 
            {
                title: "Availability",
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: [
                    {
                        title: "Total Avail. (%)",
                        headerAttributes: { style: "text-align:center;" },
                        field: "TotalAvail", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(average, 'n2')#</div>" 
                    },
                    {
                        title: "Uptime (Hrs)",
                        headerAttributes: { style: "text-align:center;" },
                        field: "OKTime", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate: "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
                    },
                ]
            },
            {
                title: "Electrical Losses (MWh)",
                headerAttributes: {
                    style: "text-align:center;"
                },
                field: "ElectricLoss", width: 100, attributes: { class: "align-center" }, format: "{0:n2}",footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
            }, 
            {
                title: "Power Curve Deviation (MW)", //Sepertinya ini MW
                headerAttributes: {
                    style: "text-align:center;"
                },
                field: "PCDeviation", width: 120, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
            }
            // , {
            //     title: "Others (MWh)", //Sepertinya ini KWh
            //     headerAttributes: {
            //         style: "text-align:center;"
            //     },
            //     field: "Others", width: 100, attributes: { class: "align-center" }, format: "{0:n2}", footerTemplate:"<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"
            // }
        ],
        dataBound: function(){
             app.loading(false);
             pg.isFirstStaticView(false);
        }
    });
}

sv.StaticView = function(){
    var valid = fa.LoadData();
    app.loading(true);    
    if (valid) {
        pg.setAvailableDate(false);
        if(pg.isFirstStaticView() === true){
            setTimeout(function(){
                var dateStart = $('#dateStart').data('kendoDatePicker').value();
                var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));  

                var param = {
                    period: fa.period,
                    dateStart: dateStart,
                    dateEnd: dateEnd,
                    turbine: fa.turbine(),
                    project: fa.project,
                };

                toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getscadasummarylist", param, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }

                    var gData = res.data.Data

                    sv.getChartView(gData);
                    sv.getGridView(gData);

                });
            },500);
        }else{
            sv.refreshView($('input[name=convertStatic]:checked').val());
            app.loading(false);
        }
    }

}