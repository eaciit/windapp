
pm.turbineListturbulence = ko.observableArray([]);
pm.turbineturbulence = ko.observableArray([]);
pm.ChartSeriesturbulence = ko.observableArray([]);
pm.DtTurbulence = ko.observableArray([]);
pm.ShowScatter = ko.observable(false);
pm.project = ko.observable();
pm.dateStart = ko.observable();
pm.dateEnd = ko.observable();

pm.ShowScatter.subscribe(function(){ 
	ti.ShowScatter(); 
});


pm.getPDF = function(selector){
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            margin: {
                left   : "5mm",
                top    : "5mm",
                right  : "5mm",
                bottom : "5mm"
            },
        });
      kendo.drawing.pdf.saveAs(group, project+"TurbulenceInstensity"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}


var ti = {
	// RefreshchartTI: function(){
	// 	setTimeout(function() {
	// 		$("#chartTI").data("kendoChart").refresh();
	// 	}, 100);
	// },
	RefreshData: function() {	  
	    var isValid = fa.LoadData();
	    if(isValid) {
	    	pm.showFilter();
		    if(pm.isFirstTurbulence() === true){
		    	app.loading(true);
		        ti.RefreshchartTI();
		        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
		        var scadaDate = ' | (<strong>SCADA HFD</strong>) from: <strong>' + availDateList.startScadaHFD + '</strong> until: <strong>' + availDateList.endScadaHFD + '</strong>'
		        $('#availabledatestart').html(metDate);
		        $('#availabledateend').html(scadaDate);
		        var project = $('#projectList').data("kendoDropDownList").value();
	            var dateStart = $('#dateStart').data('kendoDatePicker').value();
	            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
	            pm.project(project);
	            pm.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
	            pm.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

		    }else{
		        var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
		        var scadaDate = ' | (<strong>SCADA HFD</strong>) from: <strong>' + availDateList.startScadaHFD + '</strong> until: <strong>' + availDateList.endScadaHFD + '</strong>'
		        $('#availabledatestart').html(metDate);
		        $('#availabledateend').html(scadaDate);
		        var project = $('#projectList').data("kendoDropDownList").value();
	            var dateStart = $('#dateStart').data('kendoDatePicker').value();
	            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
	            pm.project(project);
	            pm.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
	            pm.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
		        setTimeout(function(){
		            $("#chartTI").data("kendoChart").refresh();
		            app.loading(false);
		        }, 300);
		    }		      
	    }
	},
	RefreshchartTI: function() {
		var param = {
            period: fa.period,
            Turbine: fa.turbine(),
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };
		toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getturbulenceintensity", param, function (data) {
            pm.isFirstTurbulence(false);

            if(data.ChartSeries == null){
            	app.loading(false);
            }

			var width = $(".main-header").width()
	        // var cfg = ti.ChartConfig(data.Data, data.ChartSeries);
	        var tempData = _.sortBy(data.ChartSeries, 'name');
	        tempData.forEach(function(val, idx){
	        	tempData[idx].idxseries = idx;
	        });
	        data.ChartSeries = tempData;
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
                    },
                }],
                xAxis: {
                    // majorUnit: 1,
                    title: {
                        text: "Wind Speed (m/s)",
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
                        text: "Turbulence Intensity",
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
	            dataBound: function(){
	                app.loading(false);
	                pm.isFirstNacelleDis(false);

	                var chart = $("#chartTI").data("kendoChart");
	                var viewModel = kendo.observable({
	                  series: chart.options.series,
	                  markerColor: function(e) {
	                    return e.get("visible") ? e.color : "grey";
	                  }
	                });

	                kendo.bind($("#legendTurbulence"), viewModel);
	            }
			});
			setTimeout(function() {
				ti.InitRightList();
				app.loading(false);
 				$("#chartTI").data("kendoChart").refresh();
			}, 100);

			$('#wCbScatter').show();
			// $('#chartTI').width($('#chartTI').parent().parent().width());
			// $('#chartTI').width(width * 0.8);
			// $('#chartTI').height(width * 0.3);
			// $('#chartTI').kendoChart(cfg);
            // $("#chartTI").data("kendoChart").refresh();
		});
	},
	showHideAllLegend: function(e){
	    if (e.checked == true) {
	        $('.fa-check-turbulence').css("visibility", 'visible');
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
	        $('.fa-check-turbulence').css("visibility", 'hidden');
	        /*$('#labelShowHide b').text('Tick All Turbines');*/
	        $('#labelShowHide b').text('Select All');
	    }
	    $('.chk-option-turbulence').not(e).prop('checked', e.checked);

	    $("#chartTI").data("kendoChart").redraw();
	},
	showHideLegend: function(idx){
	    $('#chk-turbulence' + idx).trigger('click');
	    var chart = $("#chartTI").data("kendoChart");
	    var leTur = $('input[id*=chk-turbulence][type=checkbox]').length;

	    if ($('input[id*=chk-turbulence][type=checkbox]:checked').length == $('input[id*=chk-turbulence][type=checkbox]').length) {
	        $('#showHideAllturbulence').prop('checked', true);
	    } else {
	        $('#showHideAllturbulence').prop('checked', false);
	    }

	    if ($('#chk-turbulence' + idx).is(':checked')) {
	        $('#icon-turbulence' + idx).css("visibility", "visible");
	    } else {
	        $('#icon-turbulence' + idx).css("visibility", "hidden");
	    }
	    if(!pm.ShowScatter()) {
		    if ($('#chk-turbulence' + idx).is(':checked')) {
		        $("#chartTI").data("kendoChart").options.series[idx].visible = true
		    } else {
		        $("#chartTI").data("kendoChart").options.series[idx].visible = false
		    }
		    
	    	$("#chartTI").data("kendoChart").redraw();	
	    }
    	else {
	    	ti.GetScatter(idx);
    	}
	},
	InitRightList: function(){
		$("#right-turbine-turbulence").html("");

		if (pm.ChartSeriesturbulence()!=null){
			if (pm.ChartSeriesturbulence().length > 1) {
				$("#showHideChkturbulence").html('<label style="padding-left: 1%;">' +
					'<input type="checkbox" id="showHideAllturbulence" checked onclick="ti.showHideAllLegend(this)" >' +
					'<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
					'<span id="labelShowHide"><b>Select All</b></span>' +
					'</label>');
			} else {
				$("#showHideChk").html("");
			}

			$.each(pm.ChartSeriesturbulence(), function (idx, val) {
	        $("#right-turbine-turbulence").append('<div class="btn-group">' +
				'<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="ti.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check fa-check-turbulence" id="icon-turbulence' + (idx) + '"></i></button>' +
				'<input class="chk-option-turbulence" type="checkbox" name="' + val.turbineid + '" checked id="chk-turbulence' + (idx) + '" hidden>' +
				'<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="ti.showHideLegend(' + (idx) + ')" type="button">' + val.name + '</button>' +
				'</div>');
			});
		} else {
			$("#showHideChk").html("");
		}
	},
	/// nyuwun amit, sepurane yo, iki mung tak remark, soale kok ra diceluk, fungsi iki mestine digawe mergo enek alasan pada saat itu, 
	/// menowone mben digawe maneh, lek sampe mben ra digawe, yo sing ngelanjutne iso remove remark an iki, ok?! suwun
	// ChartConfig: function(data, chartSeries) {
	// 	var colors = [];
	// 	$.each(chartSeries, function(idx,val){
	// 		colors.push(val);
	// 	});
	// 	return { 
	// 		dataSource: data,
	//         chartArea: {
	// 		    background: "transparent"
	// 		},
	// 		title: {
	// 	        visible: false
	// 	    },
	// 	    legend: {
	// 	        // visible: false,
	// 	        position: "bottom"
	// 	    },
	// 	    seriesDefaults: {
	//             type: "scatterLine",
	//             style: "smooth",
	//         },
	// 	    series: chartSeries,
	// 	    categoryAxis: {
	// 	        // field: "turbine",
	//             labels: {
	//                 visible: false,
	//             },
	//             crosshair: {
	//                 visible: false,
	//             },
	// 	        majorGridLines: {
	// 	            visible: false,
	// 	        },
	// 	        majorTicks: {
	// 	            visible: false,
	// 	        },
	// 	    },
	// 	    valueAxis: {
	// 	        crosshair: {
	//                 visible: false
	//             },
	// 	        majorGridLines: {
	// 	            visible: false,
	// 	        },
	// 	        majorTicks: {
	// 	            visible: false,
	// 	        },
	// 	    },
	//         tooltip: {
	//             visible: true,
	//             template: "#: category # = #= kendo.format('{0:N2}',value) #"
	//         },
	//         dataBound: function(e) {
	//         	var series = e.sender.options.series;
	//         	$.each(series, function(idx, s){
	//         		if(s.name=='Average') {
	//         			s.dashType = 'dash';
	//         		}
	//         	});
	//         },
	// 	};
	// },
	ShowScatter: function() {
		if(pm.ShowScatter()) {
			$('#showHideAllturbulence').prop('checked', false);
			$('#showHideAllturbulence').prop('disabled', true);
			var cbs = $('#right-turbine-turbulence').find('input[type=checkbox]');
			$.each(cbs, function(idx, elm){
				if(idx > 0) {
					$(elm).parent().find('button').find('i').css('visibility','hidden');
					$(elm).removeAttr('checked');
					$("#chartTI").data("kendoChart").options.series[idx].visible = false;
				}
			});
			ti.GetScatter(-1);
		} else {
			// $('#showHideAllturbulence').prop('checked', true);
			// $('#showHideAllturbulence').prop('disabled', false);
			// var cbs = $('#right-turbine-turbulence').find('input[type=checkbox]');
			// $.each(cbs, function(idx, elm){
			// 	$(elm).parent().find('button').find('i').removeAttr('style');
			// 	$(elm).prop('checked', true);
			// 	$("#chartTI").data("kendoChart").options.series[idx].visible = true;
			// });
			app.loading(true);
			ti.RefreshchartTI();
		}
	},
	GetScatter: function(index) {
	    app.loading(true);
		var cbsChecked = $('input[id*=chk-][type=checkbox]:checked');
        if (cbsChecked.length > 3) {
        	var cbs = $('#right-turbine-turbulence').find('input[type=checkbox]');
			$.each(cbs, function(idx, elm){
				if(idx==index) {
					$(elm).parent().find('button').find('i').css('visibility','hidden');
					$(elm).removeAttr('checked');
					//$("#chartTI").data("kendoChart").options.series[index].visible = false;
				}
			});
			app.loading(false);
        	swal('Warning', 'You can only select 3 turbines !', 'warning');
            return
        }

        if(cbsChecked.length > 0) {
        	var dtLine = [], turbines = [], colors = [];
	        var dtLineSrc = pm.DtTurbulence().ChartSeries;
        	$.each(cbsChecked, function(idx, elm){
        		var turbineName = $(elm).prop('name');
        		var dtChartSeries = _.find(dtLineSrc, function(nm) {
	                return nm.turbineid == turbineName
	            });
	            dtLine.push(dtChartSeries);
        		turbines.push(turbineName);
        		colors.push(dtChartSeries.color);
        	}); 

	        var dateStart = $('#dateStart').data('kendoDatePicker').value();
	        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

	        var param = {
	            period: fa.period,
	            dateStart: dateStart,
	            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
	            turbine: turbines,
	            project: fa.project,
	            Color: colors,
	            isClean: false,
	            isSpecific: false,
	            isDeviation: false,
	            isPower0: false,
	            deviationVal: '0',
	            DeviationOpr: '0',
	            IsDownTime: false,
	            ViewSession: ''
	        };
	        
	        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getturbulenceintensityscatter", param, function(res) {
	            if (!app.isFine(res)) {
	                app.loading(false);
	                return;
				}

	            var dataPowerCurves = res.data.Data;
	            $.each(dataPowerCurves, function(idx, val){
	            	var xnew = val;
	            	//xnew.yAxis = "power";
	            	dataPowerCurves[idx] = xnew;
	            });

	            var dtSeries = new Array();
	            if (dataPowerCurves != null) {
	                if (dataPowerCurves.length > 0) {
	                    dtSeries = dtLine.concat(dataPowerCurves);
	                }
	            } else {
	                dtSeries = dtLine;
	            }

	            $('#chartTI').html("");
	            $("#chartTI").kendoChart({
	                theme: "flat",
	                // renderAs: "canvas",
	                pdf: {
	                  fileName: "DetailPowerCurve.pdf",
	                },
	                title: {
	                    text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
	                    visible: false,
	                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
	                },
	                legend: {
	                    visible: false,
	                    position: "bottom"
	                },
	                seriesDefaults: {
	                    type: "scatterLine",
	                    style: "smooth",
	                },
	                series: dtSeries,
	                categoryAxis: {
	                    labels: {
	                        step: 1
	                    }
	                },
	                valueAxis: [{
	                    labels: {
	                        format: "N0",
	                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	                    }
	                }],
	                xAxis: {
	                    majorUnit: 1,
	                    title: {
	                        text: "Wind Speed (m/s)",
	                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	                        color: "#585555",
	                        visible: true,
	                    },
	                    labels: {
	                        format: "N0",
	                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
	                    },
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
	                yAxis: [
	                	{
	                		title: {
		                        text: "Turbulence Intensity",
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
	                	{
	                		name: "tivalue",
		                    title: {
		                        text: "Generation (KW)",
		                        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		                        color: "#585555"
		                    },
		                    labels: {
		                        format: "N0",
		                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		                    },
		                    axisCrossingValue: 5,
		                    majorGridLines: {
		                        visible: true,
		                        color: "#eee",
		                        width: 0.8,
		                    },
		                    // crosshair: {
		                    //     visible: true,
		                    //     tooltip: {
		                    //         visible: true,
		                    //         format: "N1",
		                    //         background: "rgb(255,255,255, 0.9)",
		                    //         color: "#58666e",
		                    //         font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
		                    //         border: {
		                    //             color: "#eee",
		                    //             width: "2px",
		                    //         },
		                    //     }
		                    // },
							yAxisType: 'Secondary',
							visible: false,
		                }
	                ],
	                pannable: {
	                    lock: "y"
	                },
	                zoomable: {
	                    mousewheel: {
	                        lock: "y"
	                    },
	                    selection: {
	                        lock: "y",
	                        key: "none",
	                    }
	                },
	               	dataBound: function(){
		                app.loading(false);
		                pm.isFirstNacelleDis(false);

		                var chart = $("#chartTI").data("kendoChart");
		                var viewModel = kendo.observable({
		                  series: chart.options.series,
		                  markerColor: function(e) {
		                    return e.get("visible") ? e.color : "grey";
		                  }
		                });

		                kendo.bind($("#legendTurbulence"), viewModel);
		            }
	            });

				// var series = $("#chartTI").data("kendoChart").options.series;
				// $.each(series, function(idx, elm){
				// 	$.each(cbsChecked, function(idxy, elmy) {
				// 		if(elm.name.indexOf('Scatter') < 0) {
				// 			if(elm.name!=$(elmy).prop('name')) {
				// 				$("#chartTI").data("kendoChart").options.series[idx].visible = false;
				// 			}
				// 		}
				// 	});
				// });

	            app.loading(false);
	        });

			// $("#chartTI").data("kendoChart").redraw();
		}
	},
};

$(document).ready(function() {
	// ti.LoadData();
	$('#wCbScatter').hide();
	$('#lCbScatter').on('click', function(){
		var cb = $(this).find('input[type=checkbox]');
		if(cb!=undefined) {
			pm.ShowScatter(cb.is(':checked'));
			// setTimeout(function(){
			// 	ti.ShowScatter();
			// }, 1000);
		}
	});  
});
