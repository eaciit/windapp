(function () {
	const CHART_DATA = 'chart2grid-chart'

	jQuery.fn.kendoChart2Grid =  function (options) {
		var self = this

		if (self.data('kendoGrid') != undefined) {
			return
		}
		var decimal = ''
		if(options != undefined) {
			decimal = options;
		} else {
			decimal = 'n2';
		}
		
		var chart = self.data('kendoChart')
		var isUsingDataSource = false
		var data = []
		var columns = [{
			field: 'category',
			title: 'Category Axis'
		}]

		if (chart.options.series.length > 0) {
			if (chart.options.series[0].data[0] instanceof Object) {
				isUsingDataSource = true
			}
		}

		chart.options.categoryAxis.categories.forEach(function (d, i) {
			var o = {}
			o[columns[0].field] = d

			chart.options.series.forEach(function (e, j) {
				var columnField = 'column' + j
				o[columnField] = d[e.field]

				if (i == 0) {
					columns.push({
						field: columnField,
						title: (e.hasOwnProperty('name') ? e.name : e.field),
						attributes: { style: "text-align:center;" },
						headerAttributes: { style: "text-align:center;" }
					})
				}
			})

			data.push(o)
		})

		chart.options.series.forEach(function (d, i) {
			var columnField = 'column' + i

			d.data.forEach(function (e, j) {
				let value = e

				if (isUsingDataSource) {
					value = e[d.field]
				}
				data[j][columnField] = kendo.toString(value, decimal);
			})
		})

		self.data(CHART_DATA, self.data('kendoChart'))

		if (self.data('kendoChart') != undefined) {
			self.data('kendoChart').destroy()
		}
	    self.empty()
	    self.kendoGrid({
	    	dataSource: {
	    		data: data
	    	},
	    	columns: columns
	    }) 
	}
})()