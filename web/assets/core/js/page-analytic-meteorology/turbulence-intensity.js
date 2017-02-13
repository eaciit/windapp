
pm.turbineListturbulence = ko.observableArray([]);
pm.turbineturbulence = ko.observableArray([]);
pm.ChartSeriesturbulence = ko.observableArray([]);
pm.DtTurbulence = ko.observableArray([]);

var ti = {
	// RefreshchartTI: function(){
	// 	setTimeout(function() {
	// 		$("#chartTI").data("kendoChart").refresh();
	// 	}, 100);
	// },
	RefreshData: function() {
	    app.loading(true);
	    fa.LoadData();
	    pm.showFilter();
	    if(pm.isFirstTurbulence() === true){
	        ti.RefreshchartTI();
	        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
	        var scadaDate = ' | (<strong>SCADA HFD</strong>) from: <strong>' + availDateList.startScadaHFD + '</strong> until: <strong>' + availDateList.endScadaHFD + '</strong>'
	        $('#availabledatestart').html(metDate);
	        $('#availabledateend').html(scadaDate);

	    }else{
	        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
	        var scadaDate = ' | (<strong>SCADA HFD</strong>) from: <strong>' + availDateList.startScadaHFD + '</strong> until: <strong>' + availDateList.endScadaHFD + '</strong>'
	        $('#availabledatestart').html(metDate);
	        $('#availabledateend').html(scadaDate);
	        setTimeout(function(){
	            $("#chartTI").data("kendoChart").refresh();
	            app.loading(false);
	        }, 300);
	    }

	},
	RefreshchartTI: function() {
		var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };
		toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getturbulenceintensity", param, function (data) {
            pm.isFirstTurbulence(false);

			var width = $(".main-header").width()
	        // var cfg = ti.ChartConfig(data.Data, data.ChartSeries);
	        pm.ChartSeriesturbulence(data.ChartSeries)
	        pm.DtTurbulence(data)

			$('#chartTI').html('');

			$("#chartTI").kendoChart({
			    // pdf: {
			    //   fileName: "DetailPowerCurve.pdf",
			    // },
			    theme: "flat",
			    chartArea: {
			        background: "transparent"
			    },
			    title: {
			        visible: false
			    },
			    legend: {
			        visible: false,
			        position: "bottom"
			    },
			    seriesDefaults: {
			        type: "scatterLine",
			        style: "smooth",
			    },
			    series: data.ChartSeries,
			    valueAxis: [{
                    labels: {
                        format: "N2",
                    }
                }],
                xAxis: {
                    // majorUnit: 1,
                    title: {
                        // text: "Wind Speed (m/s)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555",
                        visible: true,
                    },
                    crosshair: {
                        visible: true,
                        tooltip: {
                            visible: true,
                            format: "N2",
                            background: "rgb(255,255,255, 0.9)",
                            color: "#58666e",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            border: {
                                color: "#eee",
                                width: "2px",
                            },
                        }
                    },
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    // max: 25
                },
                yAxis: {
                    title: {
                        // text: "Generation (KW)",
                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        color: "#585555"
                    },
                    labels: {
                        format: "N2",
                    },
                    axisCrossingValue: -5,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    },
                    crosshair: {
                        visible: true,
                        tooltip: {
                            visible: true,
                            format: "N2",
                            background: "rgb(255,255,255, 0.9)",
                            color: "#58666e",
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                            border: {
                                color: "#eee",
                                width: "2px",
                            },
                        }
                    },
                },
                tooltip: {
                    visible: true,
                    // format: "{1}in {0} minutes",
                    template: "#= series.name #",
                    shared: true,
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                },
			});
			setTimeout(function() {
				ti.InitRightList();
				app.loading(false);
 				$("#chartTI").data("kendoChart").refresh();
			}, 100);



			// $('#chartTI').width($('#chartTI').parent().parent().width());
			// $('#chartTI').width(width * 0.8);
			// $('#chartTI').height(width * 0.3);
			// $('#chartTI').kendoChart(cfg);
            // $("#chartTI").data("kendoChart").refresh();
		});
	},
	showHideAllLegend: function(e){
	    if (e.checked == true) {
	        $('.fa-check').css("visibility", 'visible');
	        $.each(pm.ChartSeriesturbulence(), function (i, val) {
	            if($("#chartTI").data("kendoChart").options.series[val.idxseries] != undefined){
	                $("#chartTI").data("kendoChart").options.series[val.idxseries].visible = true;
	            }
	        });
	        /*$('#labelShowHide b').text('Untick All Turbines');*/
	        $('#labelShowHide b').text('Select All');
	    } else {
	        $.each(pm.ChartSeriesturbulence(), function (i, val) {
	            if($("#chartTI").data("kendoChart").options.series[val.idxseries] != undefined){
	                $("#chartTI").data("kendoChart").options.series[val.idxseries].visible = false;
	            }  
	        });
	        $('.fa-check').css("visibility", 'hidden');
	        /*$('#labelShowHide b').text('Tick All Turbines');*/
	        $('#labelShowHide b').text('Select All');
	    }
	    $('.chk-option').not(e).prop('checked', e.checked);

	    $("#chartTI").data("kendoChart").redraw();
	},
	showHideLegend: function(idx){
	    $('#chk-' + idx).trigger('click');
	    var chart = $("#chartTI").data("kendoChart");
	    var leTur = $('input[id*=chk-][type=checkbox]').length

	    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
	        $('#showHideAllturbulence').prop('checked', true);
	    } else {
	        $('#showHideAllturbulence').prop('checked', false);
	    }

	    if ($('#chk-' + idx).is(':checked')) {
	        $('#icon-' + idx).css("visibility", "visible");
	    } else {
	        $('#icon-' + idx).css("visibility", "hidden");
	    }

	    if ($('#chk-' + idx).is(':checked')) {
	        $("#chartTI").data("kendoChart").options.series[idx].visible = true
	    } else {
	        $("#chartTI").data("kendoChart").options.series[idx].visible = false
	    }
	    $("#chartTI").data("kendoChart").redraw();
	},
	InitRightList: function(){

	    if (pm.ChartSeriesturbulence().length > 1) {
	        $("#showHideChkturbulence").html('<label style="padding-left: 1%;">' +
	            '<input type="checkbox" id="showHideAllturbulence" checked onclick="ti.showHideAllLegend(this)" >' +
	            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
	            '<span id="labelShowHide"><b>Select All</b></span>' +
	            '</label>');
	    } else {
	        $("#showHideChk").html("");
	    }

	    $("#right-turbine-turbulence").html("");
	    $.each(pm.ChartSeriesturbulence(), function (idx, val) {
	        $("#right-turbine-turbulence").append('<div class="btn-group">' +
	            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="ti.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + (idx) + '"></i></button>' +
	            '<input class="chk-option" type="checkbox" name="' + val.name + '" checked id="chk-' + (idx) + '" hidden>' +
	            '<button class="btn btn-default btn-sm turbine-btn" onclick="ti.showHideLegend(' + (idx) + ')" type="button" style="width:70px">' + val.name + '</button>' +
	            '</div>');
	    });
	},
	ChartConfig: function(data, chartSeries) {
		var colors = [];
		$.each(chartSeries, function(idx,val){
			colors.push(val);
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
		    seriesDefaults: {
	            type: "scatterLine",
	            style: "smooth",
	        },
		    series: chartSeries,
		    categoryAxis: {
		        // field: "turbine",
	            labels: {
	                visible: false,
	            },
	            crosshair: {
	                visible: false,
	            },
		        majorGridLines: {
		            visible: false,
		        },
		        majorTicks: {
		            visible: false,
		        },
		    },
		    valueAxis: {
		        crosshair: {
	                visible: false
	            },
		        majorGridLines: {
		            visible: false,
		        },
		        majorTicks: {
		            visible: false,
		        },
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
	// ti.LoadData();
});