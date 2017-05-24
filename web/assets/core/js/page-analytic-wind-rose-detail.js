
'use strict';

vm.currentMenu('Wind Rose');
vm.currentTitle('Wind Rose');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind Rose', href: viewModel.appName + 'page/analyticwindrosedetail' }]);
vm.Turbine = ko.observableArray(["MetTower"]);
vm.Series = ko.observableArray([]);
// vm.Turbine(fa.turbine);

vm.dummyDataMetHour = ko.observableArray([]);



viewModel.WindRose = new Object();
var wrd = viewModel.WindRose;

wrd.getData = function(){
	fa.LoadData();
	app.loading(true);
    var param = {
    	period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
    };
    toolkit.ajaxPost(viewModel.appName + "analyticwindrosedetail/getwsdata", param, function (res) {
		if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        var metData = res.data.WindRose;
        vm.dummyDataMetHour(metData.reverse());

        wrd.initChart();
    
        app.loading(false)
    });
}
	 
wrd.initChart = function() {
	$.each(vm.dummyDataMetHour(), function(i, val){
		var name = val.Name
		if(name =="MetTower"){
			name = "Met Tower"
		}

   		var idChart = "#chart-" + val.Name
		$(idChart).kendoChart({
			theme : "nova",
		    title: {
		        text: name,
		        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		    },
		    legend: {
		        position: "bottom",
		        labels: {
		            template: "#= (series.data[0] || {}).wscategorydesc #"
		        },
		    },
		    dataSource: {
		    	data : val.Data,
		        group: {
		            field: "wscategoryno",
		            dir: "asc"
		        },
		        sort: {
		            field: "directionno",
		            dir: "asc"
		        }
		    },
		    seriesColors: colorField,
		    series: [{
		        type: "radarColumn",
		        stack: true,
		        field: "contribute"
		    }],
		    categoryAxis: {
		        field: "directiondesc"
		    },
		    valueAxis: {
		        visible: false
		    },
			tooltip: {
		        visible: true,
		        template: "#= category # (#= dataItem.wscategorydesc #) #= kendo.toString(value * 100, 'n2') #% for #= kendo.toString(dataItem.hours, 'n2') # hours",
		        background: "rgb(255,255,255, 0.9)",
		        color : "#58666e",
		        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		        border : {
		            color : "#eee",
		            width : "2px",
		        },
		    }
		});
		setTimeout(function() {
            if ($(idChart).data("kendoChart") != null) {
                $(idChart).data("kendoChart").refresh();
            }
        }, 10);
        
    });
}
$(document).ready(function(){
	
	$('#btnRefresh').on('click', function(){
		app.loading(true);
		setTimeout(function(){ 
			fa.LoadData();
			wrd.getData();
		},200);
	});
	
	app.loading(true);
	setTimeout(function(){ 
		var lastStartDate = new Date(Date.UTC(2016, 6, 7, 0, 0, 0, 0));
	    var lastEndDate = new Date(Date.UTC(2016, 6, 1, 0, 0, 0, 0));
	    $('#dateEnd').data('kendoDatePicker').value(lastStartDate);
	    $('#dateStart').data('kendoDatePicker').value(lastEndDate);
	},500);

	setTimeout(function(){ 
		fa.LoadData();
		wrd.getData();
	},1000);
});