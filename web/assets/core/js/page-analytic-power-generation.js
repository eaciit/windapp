'use strict';

viewModel.ClusterWiseGeneration = new Object();
var page = viewModel.ClusterWiseGeneration;

page.dataSource = ko.observableArray([]);

page.LoadData = function(){
    app.loading(true);

    var project = $('#projectList').data('kendoDropDownList').value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));

    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: fa.turbine(),
        project: project
    };

    toolkit.ajaxPost(viewModel.appName + "clusterwisegeneration/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var data = res.data.data;
        page.dataSource(data);
        page.generateChart(data);
        app.loading(false);
    });

}

page.generateChart = function(dataSource){

    var categoryTurbine = [];
    var categoryCluster = [];
    var datas = [];
    var series = [];
    $.each(dataSource, function(key, val){
        var data = {
            turbine : val.turbine, 
            cluster : val.cluster,
            sumGeneration : kendo.toString(val.sumGeneration.value , 'n2'),
            averageGa: kendo.toString(val.averageGa.value, 'n2'),
            averageMa: kendo.toString(val.averageMa.value, 'n2'),
            
        }
        datas.push(data);
    });

    datas =  _.sortBy(datas, ['cluster', 'turbine']);

    $("#cw-chart").html("");
    $("#cw-chart").kendoChart({
        theme: "flat",
        dataSource : {
            data : datas
        },
        title: {
            text: ""
        },
        legend: {
            position: "top",
            visible: true,
            labels: {
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            }
        },
        chartArea: {
            height : 370,
            padding: 10,
        },
        seriesDefaults: {
            type: "column"
        },
        series: [{
            name: "Sum of Controller Generation",
            axis : "generation",
            categoryField: "turbine",
            field : "sumGeneration",
            type: "column",
            color : "#3d8dbd",
        },{
            name: "Average of MA (%)",
            axis : "avail",
            style: "smooth",
            categoryField: "turbine",
            field : "averageMa",
            type: "line",
            width: 3,
            color : "#ffca28",
            markers: {
                visible: false,
            },
        },{
            name: "Average of Ext.GA (%)",
            axis : "avail",
            categoryField: "turbine",
            field : "averageGa",
            type: "line",
            color: "#ff7043",
            width: 3,
            markers: {
                visible: false,
            },
        }],
        valueAxes: [{
            name: "generation",
            title: {
                text: "Generation (kWh)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                step: 2,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            line: {
                visible: false
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
        },
        {
            name: "avail",
            title: {
                text: "Avail (%)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            visible: true,
            labels: {
                format : "{0:p0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            max: 1,
        }],
        categoryAxis: {
            majorGridLines: {
                visible: false
            },
            title: {
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                rotation : "auto",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none",
            axisCrossingValues: [0, 1000],
        },
        tooltip: {
            visible: true,
            background: "rgb(255,255,255, 0.9)",
            shared: true,
            sharedTemplate: kendo.template($("#template").html()),
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        }
    });
}

$(function(){
    app.loading(true);
    $('#btnRefresh').on('click', function () {
        setTimeout(function () {
            page.LoadData();
        }, 200);
    });

    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            var project = this._old;
            di.getAvailDate();
            fa.populateTurbine(project);
        }
    });

    setTimeout(function(){
        di.getAvailDate();
        fa.LoadData();
        page.LoadData();
    },300);
});