wfa.ProjectAnalysis = {
	LoadData: function() {
		var param = { Project: $('#projectList').data('kendoDropDownList').value() };
		toolkit.ajaxPost(viewModel.appName + "windfarmanalysis/getdatabyproject", param, function (data) {
            wfa.ProjectAnalysis.GenerateGrid(data.data);
            $.each(data.data.header, function(idx, val){
            	wfa.GridHeader.push(val);
            });
        });
	},
	GenerateGrid: function(data) {
		var $this = wfa.ProjectAnalysis;
		var cfg = {
			dataSource: {
				data: data.data,
				pageSize: 10,
			},
			pageable: false,
			scrollable: true,
			columns: [
				{ title: "Data Point", field: "Key", headerAttributes: { class: "align-center" }, attributes: { class: "align-left row-custom" }, 
					width: 180, locked: true, template: '<span class="cp-datapoint"></span>' },
				{ 
					title: data.header[0], field: "Roll12D", 
					headerAttributes: { class: "align-center" }, attributes: { class: "align-center row-custom" }, width: 180,
					template:'<div class="cp-roll12days" style="width: 160px; height:60px;"></div>',
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

	    $('#gridProject').html("");
	    $('#gridProject').kendoGrid(cfg);
	    // $('#gridProject').data('kendoGrid').refresh();
	},
	CreateRowChart: function(arg) {
		var obj = $('#gridProject').data('kendoGrid');
		if(obj.dataSource.data.length > 0) {
			$('#gridProject').find('tr').each(function(){
				var $this = $(this);
				var data = obj.dataItem($this);
				
				var title = $(this).find('.cp-datapoint');
				var chart1 = $(this).find('.cp-roll12days');
				var chart2 = $(this).find('.cp-roll12weeks');
				var chart3 = $(this).find('.cp-roll12months');
				var chart4 = $(this).find('.cp-roll12qtrs');
				// var chart5 = $(this).find('.cp-custom1');
				// var chart6 = $(this).find('.cp-custom2');

				var rows = wfa.GetKeyValues(data.Key);
				title.html(rows.text);
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
			});
		}
	},
};

