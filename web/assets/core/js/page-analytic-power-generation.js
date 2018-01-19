'use strict';

viewModel.ClusterWiseGeneration = new Object();
var page = viewModel.ClusterWiseGeneration;

vm.currentMenu('DGR Cluster');
vm.currentTitle('DGR Cluster');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'DGR Cluster', href: viewModel.appName + 'page/clusterwisegeneration' }]);

page.dataSource = ko.observableArray([]);
page.GenerationDetails = {
    Id          : "ClusterWise",
    Title       : "Cluster Wise Generation",
    DrildownUrl : "/dashboard/getjmrdetailspersite",
    IsLoading   : ko.observable(true),
};

function toObject(arr, heads) {
    var tmp = {};
    tmp["heads"] = heads;
    _.each(arr, function(val, i){
      tmp["data"+i] = val;
    })
    return tmp;
}

page.AdjustGridCol = function(selectorGrid){
    $(selectorGrid+' table').css('table-layout', 'fixed');
    var prevEl = null;
    $(selectorGrid+' g:first > g > g > path').each(function (i, d) {
      if (i == 0) {
          prevEl = d;
          return;
      }
      var currentX = d.getBBox().x;
      var prevX = prevEl.getBBox().x;
      var width = currentX - prevX;

      if (i == 1) {
        var col = $(selectorGrid+' colgroup col:eq(' + (i) + ')');
        if (width > 0) col.attr('width', width);
      }

      var col = $(selectorGrid+' colgroup col:eq(' + (i + 1) + ')');
      if (width > 0) col.attr('width', width);
      prevEl = d;
    });
}

page.InitGraph = function(){
    var tmp = {
        theme: "flat",
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
        seriesDefaults: {},
        series: [],
        valueAxes: [],
        categoryAxis: {
          labels:{
            font: '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            rotation: 'auto',
          },
          majorGridLines: {
            visible: false
          },
          axisCrossingValues: [0, 1000],
        },
    };
    return tmp;
}

page.InitGrid = function(){
    var tmp = {         
        groupable: false,
        sortable: true,
        filterable: false,
        pageable: false,
        scrollable: false,
        columns: [],
    };
    return tmp;
}


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

    toolkit.ajaxPost(viewModel.appName + "clusterwisegeneration/getdatadgr", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data.data;

        var categoryTurbine = [];
        var categoryCluster = [];
        var datas = [];
        var series = [];
        $.each(data, function(key, val){
            var data = {
                turbine : val.turbine, 
                cluster : val.cluster,
                sumGeneration : val.sumGeneration.value,
                averageGa: val.averageGa.value,
                averageMa: val.averageMa.value,
                averageRa: val.averageRa.value,
                
            }
            datas.push(data);
        });

        datas =  _.sortBy(datas, ['cluster', 'turbine']);

        console.log(datas);
        
        page.dataSource(data);
        // page.generateChart(data);
        page.RenderGenerationWidget(datas);
        app.loading(false);
    });

}

page.GetGraphWidth = function(selectorGraph){
    var headerWidth = $(selectorGraph+" > svg > g > path:nth-child(2)")[0].getBoundingClientRect().left - $(selectorGraph+" > svg > g > path:nth-child(1)")[0].getBoundingClientRect().left;
    var chartWidth = $(selectorGraph+" > svg > g > path:nth-child(2)")[0].getBoundingClientRect().width;
    var tmp = {
        header: headerWidth,
        chart: chartWidth,
    };  

    return tmp;
}
    
page.RenderGenerationWidget = function(master, isDetail, site){

    var conf = page.GenerationDetails;
    conf.IsLoading(true);

    var selectorGraph = "#"+conf.Id+"Chart";
    var selectorGrid = "#"+conf.Id+"Grid";

    $(selectorGraph).html("");
    $(selectorGrid).html("");

    var cluster = _.map(master, function(x){
        return kendo.toString(x.cluster,"n0");
    });

    var gridData = [];
    gridData.push(toObject(cluster, ""));


    var columns = [{
        title: " ",
        field: "heads",
        headerAttributes: { style: "text-align: center"},
        attributes:{style:"font :12px bold Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif ;text-align: right"},
    }];

    _.each(master, function(val, i){
        var tmp  = "data"+i
        columns.push({
          title: " ",
          template:"#=kendo.toString("+tmp+", 'N0') #",
          headerAttributes:{style:"text-align:center"},
          attributes:{style:"font :12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif ;text-align: right ; font-weight:bold;vertical-align: bottom;"},
        })
    });

    var grid = page.InitGrid();
    grid.dataSource = gridData;
    grid.columns = columns;
    grid.dataBound = function(){
            $(selectorGrid+' .k-grid-header').css("display", "none");
            page.AdjustGridCol(selectorGrid);
            conf.IsLoading(false);
    };
        

    var categoryTurbine = [];
    var categoryCluster = [];
    $.each(page.dataSource(), function(key, val){
        categoryCluster.push(val.cluster)
    });

    var chart  = page.InitGraph();
    chart.title.text = (site !== undefined) ? site +" "+conf.Title : conf.Title;
    chart.dataSource = master;
    chart.series = [
        {
            name: "Sum of Controller Generation",
            axis : "generation",
            field : "sumGeneration",
            type: "column",
            color : "#3d8dbd",
        },{
            name: "Average of MA (%)",
            axis : "avail",
            style: "smooth",
            field : "averageMa",
            type: "line",
            width: 3,
            color : "#ffca28",
            markers: {
                visible: false,
            },
        },{
            name: "Average of GA (%)",
            axis : "avail",
            field : "averageGa",
            type: "line",
            color: "#ff7043",
            width: 3,
            markers: {
                visible: false,
            },
        },{
            name: "Average of RA (%)",
            axis : "avail",
            field : "averageRa",
            type: "line",
            color: "#9c9c9c",
            width: 4,
            markers: {
                visible: false,
            },
        }
    ];
    chart.valueAxes = [{
            name: "generation",
            title: {
                text: "Generation (MWh)",
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
            min: 0.8,
        }],

    chart.categoryAxis.field = "turbine";
    chart.tooltip = {
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
    };

    chart.render  = function(){
        var graphWidth = page.GetGraphWidth(selectorGraph);
        columns[0].width = graphWidth.header;
        $(selectorGrid).html("");
        $(selectorGrid).css("width", (graphWidth.header+graphWidth.chart)+"px");
        $(selectorGrid).kendoGrid(grid);
    }
    
    chart.legendItemClick = function(e) {
        setTimeout(function(){
            page.redrawCategory(categoryCluster);
        },300);
    }

    $(selectorGraph).kendoChart(chart);
    var chart = $("#ClusterWiseChart").data("kendoChart");
    chart.redraw();

    page.redrawCategory(categoryCluster);
}


page.redrawCategory = function(categoryCluster){
    categoryCluster = $.unique(categoryCluster);
    categoryCluster.sort(function(a,b){return a-b});

    $.each(categoryCluster, function(key, val){
        var tableSelector = $('#ClusterWiseGrid >table>tbody>tr>td').filter(function() {
                                return $(this).text() === val.toString();
                            })
        var length = tableSelector.length;
        $.each(tableSelector, function(i, val){
            if(i < length-1){
                tableSelector.eq(i).remove();
            }else{
                tableSelector.attr('colspan',length);
            }
        });
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
            di.getAvailDate("DGRData");
            fa.populateTurbine(project);
        }
    });

    setTimeout(function(){
        // di.getAvailDate();
        $.when(di.getAvailDate("DGRData")).done(function(){
            fa.LoadData();
            page.LoadData();
        })
    },300);
});