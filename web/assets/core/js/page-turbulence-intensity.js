var ti = {
	LoadData: function() {
		var data = [];
		var chartSeries = [];
		var cfg = ti.ChartConfig(data, chartSeries);

		$('#chartTI').html('');
		$('#chartTI').width($('#chartTI').parent().parent().width());
		$('#chartTI').height(500);
		$('#chartTI').kendoChart();
	},
	ChartConfig: function(data, chartSeries) {
		var colors = [];
		$.each(chartSeries, function(idx,val){
			colors.push(val.color);
		});
		return { 
			dataSource: data,
	        chartArea: {
			    background: "transparent"
			},
			title: {
		        visible: false
		    },
		    legend: {
		        // visible: false,
		        position: "bottom"
		    },
		    seriesColors: colors,
		    seriesDefaults: {
	            type: "line",
	            style: "smooth",
	        },
		    // series: chartSeries,
		    series: [{
	            field: "Value",
	            name: "#= group.value #",
	            dashType: "solid"
	        }],
		    categoryAxis: {
		        field: "Title",
	            labels: {
	                visible: false,
	            },
	         //    crosshair: {
	         //        visible: false,
	         //    },
		        // majorGridLines: {
		        //     visible: false,
		        // },
		        // majorTicks: {
		        //     visible: false,
		        // },
		        // visible: false
		    },
		    valueAxis: {
		        // visible: false,
		        // crosshair: {
	         //        visible: false
	         //    },
		        // majorGridLines: {
		        //     visible: false,
		        // },
		        // majorTicks: {
		        //     visible: false,
		        // },
		    },
	        tooltip: {
	            visible: true,
	            template: "#: category # = #= kendo.format('{0:N2}',value) #"
	        },
	        dataBound: function(e) {
	        	var series = e.sender.options.series;
	        	$.each(series, function(idx, s){
	        		if(s.name=='Average') {
	        			s.dashType = 'dash';
	        		}
	        	});
	        },
		};
	},
};

$(document).ready(function() {
	ti.LoadData();
});