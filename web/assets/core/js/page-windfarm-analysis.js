'use strict';

viewModel.WindFarmAnalysis = new Object();
var wfa = viewModel.WindFarmAnalysis;

function addOption(project) {
	return {
		text: project,
		value: project
	}
}

wfa.ProjectList = [];
wfa.TurbineList = [];
wfa.Keys = [
	{ value: "Power", text: "Power (GW)", type: "column", color: "#EB8F1F", divider: 1000000 },
	{ value: "WindSpeed", text: "WindSpeed (m/s)", type: "area", color: "#37CAB7", divider: 1 },
	{ value: "Production", text: "Production (GWh)", type: "column", color: "#CDDC39", divider: 1000000 },
	{ value: "PLF", text: "PLF (%)", type: "line", color: "#9C28AF", divider: 1 },
	{ value: "TotalAvail", text: "Total Availability (%)", type: "line", color: "#F26A44", divider: 1 },
	{ value: "MachineAvail", text: "Machine Availability (%)", type: "line", color: "#EC1B4B", divider: 1 },
	{ value: "GridAvail", text: "Grid Availability (%)", type: "line", color: "#2E9598", divider: 1 },		
];

wfa.GetChartConfig = function(data, dataKey) {
	var ret = wfa.ChartLineCfg(data, dataKey);
	switch(dataKey.type) {
		case "column":
			ret = wfa.ChartColCfg(data, dataKey);
			break;
		case "area":
			ret = wfa.ChartAreaCfg(data, dataKey);
			break;
	}

	return ret;
};

wfa.ChartColCfg = function(data, dataKey) {
	return {
        dataSource: {
            data: data,
            sort: {
                field: "OrderNo",
                dir: "asc"
            }
        },
        chartArea: {
		    background: "transparent"
		},
        title: {
            visible: false
        },
        legend: {
            visible: false
        },
        seriesDefaults: {
            type: "column"
        },
        series:
        [{
            field: "Value",
            name: "",
            color: dataKey.color
        }],
        categoryAxis: {
            field: "Title",
            labels: {
                rotation: -90,
                visible: false
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
            visible: false
        },
        valueAxis: {
            labels: {
                //format: "N0"
            },
            //majorUnit: 10000,
            line: {
                visible: false
            },
            visible: false,
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
        }
    };
};

wfa.ChartLineCfg = function(data, dataKey) {
	return { 
		dataSource: {
            data: data,
            sort: {
                field: "OrderNo",
                dir: "asc"
            }
        },
        chartArea: {
		    background: "transparent"
		},
		title: {
	        visible: false
	    },
	    legend: {
	        visible: false
	    },
	    seriesDefaults: {
            type: "line"
        },
	    series: [{
	        field: "Value",
            name: "",
	        style: "smooth",
	        markers: {
	            visible: false
	        },
            color: dataKey.color
	    }],
	    categoryAxis: {
	        field: "Title",
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
	        visible: false
	    },
	    valueAxis: {
	        visible: false,
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
        }
	};
};

wfa.ChartAreaCfg = function(data, dataKey) {
	return {
        dataSource: {
            data: data,
            sort: {
                field: "OrderNo",
                dir: "asc"
            }
        },
        chartArea: {
		    background: "transparent"
		},
        title: {
            visible: false
        },
        legend: {
            visible: false
        },
        seriesDefaults: {
            type: "area",
            line: {
				style: "smooth"
			},
			labels: {
		     	format: "N0"
		    }
        },
        series: [{
                field: "Value",
                name: "",
            	color: dataKey.color
            }],
        categoryAxis: {
            field: "Title",
            labels: {
                rotation: -90
            },
            crosshair: {
                visible: false
            },
            majorGridLines: {
	            visible: false,
	        },
	        majorTicks: {
	            visible: false,
	        },
            visible: false
        },
        valueAxis: {
            visible: false,
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
        }
    };
};

wfa.GetKeyValues = function(objValue) {
	var ret = {};
	$.each(wfa.Keys, function(idx, val){
		if(val.value==objValue) {
			ret = val;
			return;
		}
	});

	return ret;
};

wfa.LoadData = function() {
	if($('#turbine1List').data('kendoMultiSelect').value().length == 0)
		$('#turbine1List').data('kendoMultiSelect').value(["All Turbines"]);
	if($('#turbine2List').data('kendoMultiSelect').value().length == 0)
    	$('#turbine2List').data('kendoMultiSelect').value(["All Turbines"]);

	var project = $("#projectList").data("kendoDropDownList").value();
    var turbines = [];

    var param = {
        Project: project,
        Turbines: turbines
    }

    toolkit.ajaxPost(viewModel.appName + "helper/getprojectinfo", param, function (res) {
        app.loading(true);

        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        $("#project-info").html($("#projectList").data("kendoDropDownList").value());
        $("#total-turbine-info").html('<i class="fa fa-flash tooltipster tooltipstered" aria-hidden="true" title="Total Turbine"></i>&nbsp;' + res.data.TotalTurbine);
        $("#total-capacity-info").html('<i class="fa fa-tachometer tooltipster tooltipstered" aria-hidden="true" title="Total Capacity"></i>&nbsp;' + res.data.TotalCapacity + "MW");

        app.loading(false);
    });

    var minDatetemp = new Date(availableDate.ScadaData[0]);
    var maxDatetemp = new Date(availableDate.ScadaData[1]);
    $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
    $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));

    wfa.ProjectAnalysis.LoadData();
    wfa.Turbine1Analysis.LoadData();
    wfa.Turbine2Analysis.LoadData();
}

wfa.RefreshGrid = function() {
	$('.grid-custom').each(function(){
		var grid = $(this).data("kendoGrid");
		grid.refresh();
	});
};

wfa.checkTurbine = function (elmId) {
    var arr = $('#'+elmId).data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbines");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#'+elmId).data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $('#'+elmId).data("kendoMultiSelect").value(["All Turbines"]);
    } else if (arr.length == 0) {
        $('#'+elmId).data("kendoMultiSelect").value(["All Turbines"]);
    }
}

// initiate value for projects & turbines
$.each(projects, function(idx, val) {
	wfa.ProjectList.push(addOption(val));
});
$.each(turbines, function(idx, val) {
	wfa.TurbineList.push(addOption(val));
});

$(document).ready(function(){
	$(window).on("resize orientationchange", function () {        
	    wfa.RefreshGrid();
	});

	$('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
	    wfa.RefreshGrid();
	});

	wfa.LoadData();
});