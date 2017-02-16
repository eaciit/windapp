wfa.Turbine1Analysis = {
	LoadData: function() {
		var turbines = $('#turbine1List').data('kendoMultiSelect').value();
		if(turbines[0]=="All Turbines") {
			turbines = [];
		}
		var param = { 
			project:  $('#projectTurbine1List').data('kendoDropDownList').value(),
			turbines: turbines
		};
		//console.log(param.turbines);
		toolkit.ajaxPost(viewModel.appName + "windfarmanalysis/getdatabyturbine1", param, function (data) {
            wfa.Turbine1Analysis.GenerateGrid(data.data);
        });
	},
	GenerateGrid: function(data) {
		var $this = wfa.Turbine1Analysis;
		var cfg = {
			dataSource: {
				data: data.data,
				pageSize: wfa.Keys.length,
			},
			pageable: true,
			scrollable: true,
			columns: [
				{ title: "Turbine", field: "Turbine", headerAttributes: { class: "align-center" }, 
					attributes: { class: "align-left row-custom" }, 
					width: 180, locked: true, template: '<span class="grid-group">#: data.Turbine #</span>' },
				{ title: "Data Point", field: "Key", headerAttributes: { class: "align-center" }, attributes: { class: "align-left row-custom" }, 
					width: 180, locked: true, template: '<span class="cp-datapoint"></span>' },
				{ 
					title: data.header[0], field: "Roll12D", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12days" style="width: 160px; height:60px;"></div>'
				},
				{ title: data.header[1], field: "Roll12W", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12weeks" style="width: 160px; height:60px;"></div>' },
				{ title: data.header[2], field: "Roll12M", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12months" style="width: 160px; height:60px;"></div>' },
				{ title: data.header[3], field: "Roll12Q", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12qtrs" style="width: 160px; height:60px;"></div>' },
				// { title: "Custom View 1<br /><span class='k-info'>18-Nov-2016 to 25-Nov-2016</span>", field: "Custom1", 
				// 	headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
				// 	template:'<div class="cp-custom1" style="width: 160px; height:60px;"></div>' },
				// { title: "Custom View 2<br /><span class='k-info'>1-Nov-2016 to 25-Nov-2016</span>", field: "Custom2", 
				// 	headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
				// 	template:'<div class="cp-custom2" style="width: 160px; height:60px;"></div>' },
			],
			dataBound: function(arg) {
				$this.CreateRowChart(arg);
			}
		}

	    $('#gridTurbine1').html("");
	    $('#gridTurbine1').kendoGrid(cfg);
	    $('#gridTurbine1').data('kendoGrid').refresh();
	},
	CreateRowChart: function(arg) {
		var obj = $('#gridTurbine1').data('kendoGrid');
		if(obj.dataSource.data.length > 0) {
			var fieldGroup1Before = '';
			$('#gridTurbine1').find('tr').each(function(){
				var $this = $(this);
				var data = obj.dataItem($this);

				if(data!=null) {
					var fieldGroup1 = $(this).find('.grid-group');
					if(fieldGroup1.text()!=fieldGroup1Before) {
						fieldGroup1Before = fieldGroup1.text();
					} else {
						fieldGroup1.text('');
					}
					
					var title = $(this).find('.cp-datapoint');
					var chart1 = $(this).find('.cp-roll12days');
					var chart2 = $(this).find('.cp-roll12weeks');
					var chart3 = $(this).find('.cp-roll12months');
					var chart4 = $(this).find('.cp-roll12qtrs');
					// var chart5 = $(this).find('.cp-custom1');
					// var chart6 = $(this).find('.cp-custom2');

					var rows = wfa.GetKeyValues(data.Key);
					title.html(rows.text.replace('(GW','(MW'));
					chart1.kendoChart(wfa.GetChartConfig(data.Roll12Days.ValueAvg, rows));
					chart2.kendoChart(wfa.GetChartConfig(data.Roll12Weeks.ValueAvg, rows));
					chart3.kendoChart(wfa.GetChartConfig(data.Roll12Months.ValueAvg, rows));
					chart4.kendoChart(wfa.GetChartConfig(data.Roll12Quarters.ValueAvg, rows));
					// chart5.kendoSparkline({
					// 	type: data.RowType,
					// 	data: data.Custom1,
					// });
					// chart6.kendoSparkline({
					// 	type: data.RowType,
					// 	data: data.Custom2,
					// });
				}
			});
		}
		app.loading(false);
	},
}