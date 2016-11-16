'use strict';

(function($) {

	$.fn.kendoChartToGrid = function() {
		var $this = $(this);

		var $replace = function (t, i, n) {
			return typeof i == 'string' && (i = [i]), i.forEach(function (i) {
				t = t.replace(new RegExp("\\" + i, "g"), n);
			}), t;
		};

		var $prefId = $replace($this.attr('id'), '-', '_');
		var $parent = $this.parent();
		var $cetegeNav = function() {
			var wrpNav = $('<div/>', {
				'class': 'ec-cetege-nav'
			});
			var navBtn = $('<div/>', {
				'class': 'btn-group',
				'data-toggle': 'buttons'
			})
			.append('<label class="btn btn-sm btn-primary active radio-filter"><input type="radio" name="btnSwitch'+$prefId+'" id="btnToChart'+ $prefId +'" autocomplete="off" value="chart" checked><i class="fa fa-bar-chart"></i></label><label class="btn btn-sm btn-primary radio-filter"><input type="radio" name="btnSwitch'+$prefId+'" id="btnToTable'+$prefId+'" autocomplete="off" value="table"><i class="fa fa-table"></i></label>')
			.appendTo(wrpNav);

			return wrpNav;
		};
		var $newColumn = function(colId, colName, format, datas) {
			var data = {
				colId: colId,
				colName: colName,
				colValues: datas,
				colFormat: format,
			};
			return data;
		}
		var $chartToGrid = function(objChart, objTarget) {
			var $chart = $(objChart).data("kendoChart");
			var $data = [];
			if($chart!=undefined) {
				var $opts = $chart.options;
				if($opts!=undefined) {
					var format = '';
					if($opts.hasOwnProperty('valueAxis')) {
						if($opts.valueAxis.hasOwnProperty('labels')) {
							if($opts.valueAxis.labels.hasOwnProperty('format'))
								format = $opts.valueAxis.labels.format;
						}
					}
					
					var cats = $opts.categoryAxis.categories;
					$data.push($newColumn("category", "Category", format, cats));
					
					var series = $opts.series;
					if(series.length > 0) {
						$.each(series, function(index, value){
							var colid = $replace(value.name.toLowerCase(), " ", "_");
							colid = $replace(colid, "(", "_");
							colid = $replace(colid, ")", "_");
							colid = $replace(colid, "/", "_");
							colid = $replace(colid, "-", "_");
							$data.push($newColumn(colid, value.name, format, value.data));
						});
					}
				}
			}

			var cols = [];
		    $.each($data, function(index, value){
		        var col = {
		            field: value.colId,
		            title: value.colName,
		            format: value.colFormat,
		        };
		        cols.push(col);
		    });

		    var gridDatas = [];
		    var totalValue = $data[0].colValues.length;
		    for(var i=0;i<totalValue;i++){
		        var gridDt = {};
		        $.each(cols, function(index, value){
		            gridDt[value.field] = $data[index].colValues[i];
		        });

		        gridDatas.push(gridDt);
		    }

		    $(objTarget).kendoGrid({
		        dataSource: {
		            data: gridDatas,
		            pageSize: 10
		        },
		        columns: cols,
		        pageable: true,
		        resizable: true,
		        reorderable: true,
		    });
		}

		$cetegeNav().appendTo($parent);
		$parent.append('<div id="cetege-grid'+$prefId+'"></div>');

		$('#cetege-grid'+$prefId).toggle();
		$('input[name="btnSwitch'+$prefId+'"]').on('change', function(){
			var value = $(this).val();
	        if(value=='table') {
	            $chartToGrid($($this), $('#cetege-grid'+$prefId));
	            $this.slideUp( "slow", function() {
			    	$('#cetege-grid'+$prefId).slideDown( "slow" );	
			  	});
	        } else {
	            $( "#cetege-grid"+$prefId ).slideUp( "slow", function() {
			    	$this.slideDown( "slow" );	
			  	});
	        }
	    });
	}

}(jQuery));