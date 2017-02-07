wfa.Turbine2Analysis = {
	LoadData: function() {
		var turbines = $('#turbine1List').data('kendoMultiSelect').value();
		if(turbines[0]=="All Turbines") {
			turbines = [];
		}
		var param = { 
			project:  $('#projectTurbine1List').data('kendoDropDownList').value(),
			turbines: turbines
		};
		toolkit.ajaxPost(viewModel.appName + "windfarmanalysis/getdatabyturbine2", param, function (data) {
           wfa.Turbine2Analysis.GenerateGrid(data.data);
        });
     //    var data = {
     //    	ChartData: [
	    //     	{
	    //     		Key: "Power",
	    //     		OrderNo: 0,
	    //     		ProjectName: "Tejuva",
	    //     		Roll12Days: [],
	    //     		Roll12Weeks: [],
	    //     		Roll12Months: [],
	    //     		Roll12Quarters: [],
	    //     	}
     //    	],
     //    	ChartSeries: [
     //    		{ field: "Average", color: "#ED1C24", name: "Average" },
     //    		{ field: "HBR004", color: "#A3238E", name: "HBR004" },
     //    		{ field: "HBR005", color: "#00A65D", name: "HBR005" },
     //    		{ field: "HBR006", color: "#F58220", name: "HBR006" },
     //    		{ field: "HBR007", color: "#0066B3", name: "HBR007" },
     //    	]
    	// };
     //    wfa.Turbine2Analysis.GenerateGrid(data);
	},
	GenerateGrid: function(data) {
		var chartSeries = data.ChartSeries;
		var $this = wfa.Turbine2Analysis;
		var titles = wfa.GridHeader();
		var cfg = {
			dataSource: {
				data: data.ChartData,
				pageSize: wfa.Keys.length,
			},
			pageable: true,
			scrollable: true,
			columns: [
				{ title: "Data Point", field: "Key", headerAttributes: { class: "align-center" }, attributes: { class: "align-left row-custom" }, 
					width: 180, locked: true, template: '<span class="cp-datapoint"></span>' },
				{ 
					title: titles[0], field: "Roll12D", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12days" style="width: 160px; height:120px;"></div>'
				},
				{ title: titles[1], field: "Roll12W", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12weeks" style="width: 160px; height:120px;"></div>' },
				{ title: titles[2], field: "Roll12M", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12months" style="width: 160px; height:120px;"></div>' },
				{ title: titles[3], field: "Roll12Q", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12qtrs" style="width: 160px; height:120px;"></div>' },
				// { title: "Custom View 1<br /><span class='k-info'>18-Nov-2016 to 25-Nov-2016</span>", field: "Custom1", 
				// 	headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
				// 	template:'<div class="cp-custom1" style="width: 160px; height:60px;"></div>' },
				// { title: "Custom View 2<br /><span class='k-info'>1-Nov-2016 to 25-Nov-2016</span>", field: "Custom2", 
				// 	headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
				// 	template:'<div class="cp-custom2" style="width: 160px; height:60px;"></div>' },
			],
			dataBound: function(arg) {
				$this.CreateRowChart(arg, chartSeries);
			}
		}

	    $('#gridTurbine2').html("");
	    $('#gridTurbine2').kendoGrid(cfg);
	    $('#gridTurbine2').data('kendoGrid').refresh();
	},
	CreateRowChart: function(arg, chartSeries) {
		var obj = $('#gridTurbine2').data('kendoGrid');
		if(obj.dataSource.data.length > 0) {
			$('#gridTurbine2').find('tr').each(function(){
				var $this = $(this);
				var data = obj.dataItem($this);

				if (data!=null) {
					var title = $(this).find('.cp-datapoint');
					var chart1 = $(this).find('.cp-roll12days');
					var chart2 = $(this).find('.cp-roll12weeks');
					var chart3 = $(this).find('.cp-roll12months');
					var chart4 = $(this).find('.cp-roll12qtrs');
					// var chart5 = $(this).find('.cp-custom1');
					// var chart6 = $(this).find('.cp-custom2');

					var rows = wfa.GetKeyValues(data.Key);
					title.html(rows.text); //.replace('(GW','(MW'));

					var rowSelected = {};
					rowSelected.text = rows.text;
					rowSelected.value = rows.value;
					rowSelected.color = rows.color;
					rowSelected.divider = rows.divider;
					rowSelected.type = "line";

					var datachart1 = new kendo.data.DataSource({
			            data: data.Roll12Days,
			            group: {
			                field: "Turbine"
			            },
			            sort: {
			                field: "OrderNo",
			                dir: "asc"
			            },
			        });

					var c1 = chart1.kendoChart(wfa.Turbine2Chart(datachart1, chartSeries));

					var datachart2 = new kendo.data.DataSource({
			            data: data.Roll12Weeks,
			            group: {
			                field: "Turbine"
			            },
			            sort: {
			                field: "OrderNo",
			                dir: "asc"
			            },
			        });

					var c2 = chart2.kendoChart(wfa.Turbine2Chart(datachart2, chartSeries));

					var datachart3 = new kendo.data.DataSource({
			            data: data.Roll12Months,
			            group: {
			                field: "Turbine"
			            },
			            sort: {
			                field: "OrderNo",
			                dir: "asc"
			            },
			        });

					var c3 = chart3.kendoChart(wfa.Turbine2Chart(datachart3, chartSeries));

					var datachart4 = new kendo.data.DataSource({
			            data: data.Roll12Quarters,
			            group: {
			                field: "Turbine"
			            },
			            sort: {
			                field: "OrderNo",
			                dir: "asc"
			            },
			        });

					var c4 = chart4.kendoChart(wfa.Turbine2Chart(datachart4, chartSeries));
					
				}
			});
		}
	},
}