(function () {
	const CHART_DATA = 'chart2grid-chart'

	jQuery.fn.kendoChart2Grid =  function (options, dateFormat) {

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

		var Template = "#= kendo.toString(kendo.parseDate(category, 'yyyy-MM-dd'), '" + dateFormat+ "') #"
		
		var chart = self.data('kendoChart')
		var isUsingDataSource = false
		var data = []

		var addwidth = false
		if (chart.options.categoryAxis.categories.length > 4) {
			addwidth = true
		}

		var columns = [{
			field: 'category',
			template: dateFormat == undefined ? "#= category #" : Template,
			title: " "
		}]

		if (addwidth) {
			columns[0].width = "100px"
		}

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
					var tColumn = {
						field: columnField,
						title: (e.hasOwnProperty('name') ? e.name : e.field),
						attributes: { style: "text-align:center;" },
						headerAttributes: { style: "text-align:center;" },
						//template: "#: kendo.toString(kendo.parseDate(new Date('DateId')), 'dd-MM-yyyy')#"
					}

					if (addwidth) {
						tColumn.width = "60px"
					}

					columns.push(tColumn)
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
		// console.log(columns)
	    self.empty()
	    self.kendoGrid({
	    	dataSource: {
	    		data: data
	    	},
	    	columns: columns
	    }) 
	}
})()