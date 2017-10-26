'use strict';

viewModel.KeyMetrics = new Object();
var km = viewModel.KeyMetrics;
km.MinValue = ko.observable(0);
km.MaxValue = ko.observable(2000);
km.BinValue = ko.observable(40);
km.MinValueTemp = ko.observable(0);    //
km.MaxValueTemp = ko.observable(100); // Jare mak-e static sek ae, mengko selanjutnya monggo dibuat dynamic. backend sudah ada
km.BinValueTemp = ko.observable(30);
km.MinValueWindSpeed = ko.observable(0);
km.MaxValueWindSpeed = ko.observable(24);
km.BinValueWindSpeed = ko.observable(24);
km.CategoryProduction = ko.observableArray([]);
km.ValueProduction = ko.observableArray([]);
km.dsCategorywindspeed = ko.observableArray();
km.dsValuewindSpeed = ko.observableArray();
km.dsTotaldataWS = ko.observable();
km.dsCategoryProduction = ko.observableArray();
km.dsValueProduction = ko.observableArray();
km.dsTotaldataProduction = ko.observable();
km.dsCategoryTemp = ko.observableArray();
km.dsValueTemp = ko.observableArray();
km.dsTotaldataTemp = ko.observable();

km.dsWindTurbinename = ko.observableArray([]);
km.dsTempTurbinename = ko.observableArray([]);
km.dsProdTurbinename = ko.observableArray([]);


km.tempTagsDs = ko.observableArray();
km.tempTagsList = ko.observableArray();

km.MaxValueTempList = ko.observableArray([]);

km.histogramCols = ko.observableArray([
    { text: "Production", value: "production" },
    { text: "Wind Speed", value: "windspeed" },
    { text: "Temperature", value: "temperature" },
]);

km.pageView = ko.observable("windspeed");

km.ExportKeyMetrics = function () {
    var chart = $("#dh-chart").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "AnalyticDataHistogram.pdf",
        });
    });
}

km.createChart = function (turbinename) {
    $("#totalCountData").html('(Total Count Data: ' + km.dsTotaldataWS() + ')');
    var turbineData = '';
    if(fa.turbine().length == 0) {
        turbineData = 'All Turbines';
    }else if($(".multiselect-native-select").find($(".multiselect-item.multiselect-all.active")).length == 1){
        turbineData = 'All Turbines';
    } else {
        var turbineName;
        for(var i=0; i<fa.turbine().length; i++) {
            if(i==0) {
                turbineName = turbinename[fa.turbine()[i]];
            } else {
                turbineName += ", " + turbinename[fa.turbine()[i]];
            }
        }
        turbineData = turbineName;
    }
    $("#turbineListTitle").html('for ' + turbineData);
    $("#dh-chart").replaceWith('<div id="dh-chart"></div>');

    $("#dh-chart").kendoChart({
        theme: "flat",
        legend: {
            position: "top",
            visible: false
        },
        chartArea: {
            background: "transparent",
            margin: 0,
            padding: 0
        },
        seriesDefaults: {
            type: "column",
            gap: 0,
            border: 1
        },
        series: [{
            name: "Total Count of Wind Speed (m/s)",
            data: km.dsValuewindSpeed(),
            color: "#337ab7"
        }],
        valueAxis: {
            title: {
                text: "Percentage of Wind Speed (%)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                // format: "{0:p2}"
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                format: "{0}"
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            line: {
                visible: false
            },
            axisCrossingValue: 0
        },
        categoryAxis: {
            title: {
                text: "Wind Speed (m/s)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            categories: km.dsCategorywindspeed(),
            majorGridLines: {
                visible: false
            },
            line: {
                visible: false
            },
            labels: {
                // padding: { 
                //     left: 600 / valuewindspeed.length
                // },
                // margin: {
                //     left: -600 / km.dsValuewindSpeed().length
                // },
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                template: "#: (value.split('~'))[0] #"
            },
            axisCrossingValue: [0]
        },
        tooltip: {
            format: "{0:n0}%",
            visible: true,
            template: "#= category # : #= value #%",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });

}

km.createChartProduction = function (turbinename) {
    $("#totalCountProd").html('(Total Count Data: ' + km.dsTotaldataProduction() + ')');
    var turbineData = '';
    if(fa.turbine().length == 0) {
        turbineData = 'All Turbines';
    }else if($(".multiselect-native-select").find($(".multiselect-item.multiselect-all.active")).length == 1){
        turbineData = 'All Turbines';
    } else {
        var turbineName;
        for(var i=0; i<fa.turbine().length; i++) {
            if(i==0) {
                turbineName = turbinename[fa.turbine()[i]];
            } else {
                turbineName += ", " + turbinename[fa.turbine()[i]];
            }
        }
        turbineData = turbineName;
    }
    var _rotationlabel = 0
    if (km.BinValue() > 20) {
        _rotationlabel = 68
    }
    $("#turbineListProd").html('for ' + turbineData);
    $("#dhprod-chart").replaceWith("<div id='dhprod-chart'></div>");
    $("#dhprod-chart").kendoChart({
        theme: "flat",
        legend: {
            position: "top",
            visible: false
        },
        chartArea: {
            background: "transparent",
            margin: 0,
            padding: 0
        },
        seriesDefaults: {
            type: "column",
            gap: 0,
            border: 1
        },
        series: [{
            name: "Production (MWh)",
            data: km.dsValueProduction(),
            color: "#ea5b19"
        }],
        valueAxis: {
            title: {
                text: "Percentage of Production (%)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            line: {
                visible: false
            },
            min: 0,
            // max: 100,
            axisCrossingValue: 0
        },
        categoryAxis: {
            title: {
                text: "Production (MWh)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            categories: km.dsCategoryProduction(),
            majorGridLines: {
                visible: false
            },
            line: {
                visible: false
            },
            labels: {
                rotation : _rotationlabel,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                // padding: { 
                //     left: 600 / categoryproduction.length
                // },
                // margin: {
                //     left: -600 / km.dsCategoryProduction().length
                // },
                template: "#: ((value.split('~'))[0]) #",
                format: "{0:n0}"
            }
        },
        tooltip: {
            visible: true,
            format: "{0:n0}%",
            template: "#= category # : #= value #%",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });

}

km.createChartTemp = function (turbinename) {
    $("#totalCountTemp").html('(Total Count Data: ' + km.dsTotaldataTemp() + ')');
    var turbineData = '';
    if(fa.turbine().length == 0) {
        turbineData = 'All Turbines';
    }else if($(".multiselect-native-select").find($(".multiselect-item.multiselect-all.active")).length == 1){
        turbineData = 'All Turbines';
    } else {
        var turbineName;
        for(var i=0; i<fa.turbine().length; i++) {
            if(i==0) {
                turbineName = turbinename[fa.turbine()[i]];
            } else {
                turbineName += ", " + turbinename[fa.turbine()[i]];
            }
        }
        turbineData = turbineName;
    }
    var _rotationlabel = 0
    if (km.BinValue() > 20) {
        _rotationlabel = 68
    }
    $("#turbineListTemp").html('for ' + turbineData);
    $("#dhtemp-chart").replaceWith("<div id='dhtemp-chart'></div>");
    $("#dhtemp-chart").kendoChart({
        theme: "flat",
        legend: {
            position: "top",
            visible: false
        },
        chartArea: {
            background: "transparent",
            margin: 0,
            padding: 0
        },
        seriesDefaults: {
            type: "column",
            gap: 0,
            border: 1
        },
        series: [{
            name: "Temperature",
            data: km.dsValueTemp(),
            color: "#ea5b19"
        }],
        valueAxis: {
            title: {
                text: "Percentage of Temperature (%)",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                format: "{0}",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            },
            line: {
                visible: false
            },
            min: 0,
            axisCrossingValue: 0
        },
        categoryAxis: {
            title: {
                text: "Temperature",
                visible: true,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            categories: km.dsCategoryTemp(),
            majorGridLines: {
                visible: false
            },
            line: {
                visible: false
            },
            labels: {
                rotation : _rotationlabel,
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                template: "#: ((value.split('~'))[0]) #",
                format: "{0:n0}"
            }
        },
        tooltip: {
            visible: true,
            format: "{0:n0}%",
            template: "#= category # : #= value #%",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            }
        }
    });

}

vm.currentMenu('Histograms');
vm.currentTitle('Histograms');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Histograms', href: viewModel.appName + 'page/analyticdatahistogram' }]);

km.getData = function () {
    // fa.getProjectInfo();
    if(fa.LoadData()) {
        app.loading(true);
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 

        var tagList = km.tempTagsList();
        km.tempTagsDs(eval("tagList."+ fa.project));

        $('#sTempTags').kendoDropDownList({
            dataSource: km.tempTagsDs(),
            dataValueField: 'colname', 
            dataTextField: 'text',
            change: function () {  
                km.setMaxValue();
            }
        });

        var paramFilter = {
            period: fa.period,
            Turbine: fa.turbine(),
            DateStart: dateStart,
            DateEnd: dateEnd,
            Project: fa.project
        };
        // var request;
        switch(km.pageView()) {
            case "windspeed":
                var parDataWS = {
                    MinValue: parseFloat(km.MinValueWindSpeed()),
                    MaxValue: parseFloat(km.MaxValueWindSpeed()),
                    BinValue: parseInt(km.BinValueWindSpeed()),
                    Filter: paramFilter
                };
                toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gethistogramdata", parDataWS, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    if (res.data != null) {
                        km.dsWindTurbinename(res.data.turbinename);
                        km.dsCategorywindspeed(res.data.categorywindspeed);
                        km.dsValuewindSpeed(res.data.valuewindspeed);
                        km.dsTotaldataWS(res.data.totaldata);
                        // km.dsValuewindSpeed.push(0);
                        // km.dsCategorywindspeed.push(km.dsCategorywindspeed()[km.dsCategorywindspeed().length - 1].split(' ~ ')[1]);
                        km.createChart(res.data.turbinename);

                        app.loading(false);
                    }
                });
                break;
            case "production":
                var parDataProd = {
                    MinValue: parseFloat(km.MinValue()),
                    MaxValue: parseFloat(km.MaxValue()),
                    BinValue: parseInt(km.BinValue()),
                    Filter: paramFilter
                };
                toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getproductionhistogramdata", parDataProd, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    if (res.data != null) {
                        km.dsProdTurbinename(res.data.turbinename);
                        km.dsCategoryProduction(res.data.categoryproduction);
                        km.dsValueProduction(res.data.valueproduction);
                        km.dsTotaldataProduction(res.data.totaldata);
                        // km.dsValueProduction.push(0);
                        // km.dsCategoryProduction.push(km.dsCategoryProduction()[km.dsCategoryProduction().length - 1].split(' ~ ')[1]);
                        km.createChartProduction(res.data.turbinename);

                        app.loading(false);
                    }
                });
                break;
            case "temperature":
                    var parDataTemp = {
                        MinValue: parseFloat(km.MinValueTemp()),
                        MaxValue: parseFloat(km.MaxValueTemp()),
                        BinValue: parseInt(km.BinValueTemp()),
                        FieldName: $('#sTempTags').data('kendoDropDownList').value(),
                        Filter: paramFilter,
                    };
                    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettemphistogramdata", parDataTemp, function (res) {
                        if (!app.isFine(res)) {
                            return;
                        }
                        if (res.data != null) {
                            km.dsTempTurbinename(res.data.turbinename);
                            km.dsCategoryTemp(res.data.category);
                            km.dsValueTemp(res.data.value);
                            km.dsTotaldataTemp(res.data.totaldata);
                            km.createChartTemp(res.data.turbinename);

                            app.loading(false);
                        }
                    });
                break;
        }

        // $.when(requestHistogram, requestProduction, requestHistogramTemp).done(function(){
       // $.when(request).done(function(){
       //      setTimeout(function(){
       //          app.loading(false);
       //      },500);
       //  });
    }
}

km.SubmitValues = function () {
    km.getData();
}

km.getTempTags = function() {
    var param = {};
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gettemptags", param, function (res) {
        if (res.data != null) {
            km.tempTagsList(res.data);
        }
    });
}


km.getMaxMinValueTemp = function(isRefresh){
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 
    var project = $('#projectList').data('kendoDropDownList').value();
    var tagList = [];

    $.each(km.tempTagsList()[project] , function(key, val){
        tagList.push(val.colname);
    });

    var paramFilter = {
        period: fa.period,
        Turbine: fa.turbine(),
        DateStart: dateStart,
        DateEnd: dateEnd,
        Project: project
    };

    var parDataTemp = {
        FieldList: tagList,
        MinValue: parseFloat(km.MinValueTemp()),
        MaxValue: parseFloat(km.MaxValueTemp()),
        BinValue: parseInt(km.BinValueTemp()),
        Filter: paramFilter,
    };
 
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getmaxvaltemptags", parDataTemp, function (res) {
        if (res.data != null) {
            setTimeout(function(){
                km.MaxValueTempList(res.data);
                if(isRefresh == true){
                    km.setMaxValue(true);
                }
            },500);
        }
    });

}

km.setMaxValue = function(isRefresh){
    setTimeout(function(){
        var tagTemp = $('#sTempTags').data('kendoDropDownList').value();
        var val = km.MaxValueTempList()[tagTemp];

        console.log(tagTemp);
        var maxValue = (val !== null) ? kendo.toString(val , 'n0') : 100;


        if(km.MaxValueTemp(maxValue) && isRefresh == true){
           km.getData();
        }else{
           km.MaxValueTemp(maxValue);
        }

    },500);
}

km.changePageView = function() {
    km.pageView($('#select-page-view').val());



    setTimeout(function () {
        var $el = $("#turbineList");
        $('option', $el).each(function(element) {
          $el.multiselect('deselect', $(this).val());
        });

        if(fa.turbineList().length > 1){
            $('#turbineList').multiselect('select', fa.turbineList()[0].value);
        }

        if(km.pageView() == "temperature"){
            km.setMaxValue(true);
        }else{
            km.getData();
        }
        
    }, 300);
};

$(document).ready(function () {
    di.getAvailDate();
    
    km.getTempTags();
   
    $('#btnRefresh').on('click', function () {
        app.loading(true);
        fa.checkTurbine();
        $("#sTempTags").data("kendoDropDownList").setDataSource(km.tempTagsDs());
        $("#sTempTags").data("kendoDropDownList").select(0);
        km.getMaxMinValueTemp(true);
    });

    $('#exportXlsx').on('click', function (e) {
        window.open('data:application/vnd.ms-excel,' + encodeURIComponent($('div[id$=dhprod-chart]').html()));
        e.preventDefault();
    });

    $('#projectList').kendoDropDownList({
		change: function () {  
			var project = $('#projectList').data("kendoDropDownList").value();
			fa.populateTurbine(project);
            setTimeout(function() {
                di.getAvailDate();
                $('#turbineList').multiselect('select', fa.turbineList()[0].value);
            }, 100);
		}
	});

    setTimeout(function () {
        var $el = $("#turbineList");
        $('option', $el).each(function(element) {
          $el.multiselect('deselect', $(this).val());
        });

        if(fa.turbineList().length > 1){
            $('#turbineList').multiselect('select', fa.turbineList()[0].value);
        }

        km.getData();
        km.getMaxMinValueTemp();
    }, 800);
});

$(document).bind("kendo:skinChange", km.createChart);