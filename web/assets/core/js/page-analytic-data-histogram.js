'use strict';

viewModel.KeyMetrics = new Object();
var km = viewModel.KeyMetrics;
km.MinValue = ko.observable(0);
km.MaxValue = ko.observable(2000);
km.BinValue = ko.observable(40);
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

vm.currentMenu('Histograms');
vm.currentTitle('Histograms');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Histograms', href: viewModel.appName + 'page/analyticdatahistogram' }]);

km.getData = function () {
    // fa.getProjectInfo();
    if(fa.LoadData()) {
        app.loading(true);

        di.getAvailDate();

        var paramFilter = {
            period: fa.period,
            Turbine: fa.turbine(),
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

        var parDataWS = {
            MinValue: parseFloat(km.MinValueWindSpeed()),
            MaxValue: parseFloat(km.MaxValueWindSpeed()),
            BinValue: parseInt(km.BinValueWindSpeed()),
            Filter: paramFilter
        };
        var requestHistogram = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/gethistogramdata", parDataWS, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data != null) {
                km.dsCategorywindspeed(res.data.categorywindspeed);
                km.dsValuewindSpeed(res.data.valuewindspeed);
                km.dsTotaldataWS(res.data.totaldata);
                // km.dsValuewindSpeed.push(0);
                // km.dsCategorywindspeed.push(km.dsCategorywindspeed()[km.dsCategorywindspeed().length - 1].split(' ~ ')[1]);
                km.createChart(res.data.turbinename);
            }
        });

        var parDataProd = {
            MinValue: parseFloat(km.MinValue()),
            MaxValue: parseFloat(km.MaxValue()),
            BinValue: parseInt(km.BinValue()),
            Filter: paramFilter
        };
        var requestProduction  = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getproductionhistogramdata", parDataProd, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data != null) {

                km.dsCategoryProduction(res.data.categoryproduction);
                km.dsValueProduction(res.data.valueproduction);
                km.dsTotaldataProduction(res.data.totaldata);
                // km.dsValueProduction.push(0);
                // km.dsCategoryProduction.push(km.dsCategoryProduction()[km.dsCategoryProduction().length - 1].split(' ~ ')[1]);
                km.createChartProduction(res.data.turbinename);
            }
        });

        $.when(requestHistogram, requestProduction).done(function(){
            setTimeout(function(){
                app.loading(false);
            },500);
        });
    }
}

km.SubmitValues = function () {
    km.getData();
}


$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        setTimeout(function () {
            km.getData();
        }, 300);
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
    }, 800);
});

$(document).bind("kendo:skinChange", km.createChart);