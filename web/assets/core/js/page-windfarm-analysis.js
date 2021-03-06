'use strict';

viewModel.WindFarmAnalysis = new Object();
var wfa = viewModel.WindFarmAnalysis;

vm.currentMenu('Wind Farm Analysis');
vm.currentTitle('Wind Farm Analysis');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Wind Farm Analysis', href: viewModel.appName + 'page/windfarmanalysis' }]);

function addOption(project) {
	return {
		text: project,
		value: project
	}
}

wfa.isProjectTab = ko.observable(false);
wfa.isTurbine1Tab = ko.observable(false);
wfa.isTurbine2Tab = ko.observable(false);

wfa.isProjectLoaded = ko.observable(true);
wfa.isTurbine1Loaded = ko.observable(false);
wfa.isTurbine2Loaded = ko.observable(false);
wfa.GridHeader = ko.observableArray([]);

wfa.ProjectList = ko.observableArray([]);
wfa.TurbineList = ko.observableArray([]);
wfa.Keys = [
	{ value: "Power", text: "Power (GW)", type: "column", color: "#EB8F1F", divider: 1000000 },
	{ value: "WindSpeed", text: "WindSpeed (m/s)", type: "area", color: "#37CAB7", divider: 1 },
	{ value: "Production", text: "Production (GWh)", type: "column", color: "#CDDC39", divider: 1000000 },
	{ value: "PLF", text: "PLF (%)", type: "line", color: "#9C28AF", divider: 1 },
	{ value: "TotalAvail", text: "Total Availability (%)", type: "line", color: "#F26A44", divider: 1 },
	{ value: "MachineAvail", text: "Machine Availability (%)", type: "line", color: "#EC1B4B", divider: 1 },
	{ value: "GridAvail", text: "Grid Availability (%)", type: "line", color: "#2E9598", divider: 1 },		
];


wfa.getDataAvailableInfo =  function(id){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var availDateAll = res.data;

        var namaproject,projectVal ;

        if(id == null){
            projectVal = $("#projectList").data("kendoDropDownList").value();
        }else{
            projectVal = $("#"+id).data("kendoDropDownList").value();
        }
       

        if( projectVal == undefined || projectVal == "") {
            namaproject = "Tejuva";
        }else{
            namaproject= projectVal;
        }

        
        var minDatetemp = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY')));
        var maxDatetemp = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY')));

        var maxDateData = new Date(availDateAll[namaproject]["ScadaData"][1]);

        var startDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

        $('#availabledatestartscada').html(minDatetemp);
        $('#availabledateendscada').html(maxDatetemp);

        // $('#availabledatestartscada').data('kendoDatePicker').value(startDate);
        // $('#availabledateendscada').data('kendoDatePicker').value(maxDatetemp);
    })
}

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
    app.loading(true);
    wfa.getDataAvailableInfo();
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

        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        $("#project-info").html($("#projectList").data("kendoDropDownList").value());
        $("#total-turbine-info").html('<i class="fa fa-flash tooltipster tooltipstered" aria-hidden="true" title="Total Turbine"></i>&nbsp;' + res.data.TotalTurbine);
        $("#total-capacity-info").html('<i class="fa fa-tachometer tooltipster tooltipstered" aria-hidden="true" title="Total Capacity"></i>&nbsp;' + res.data.TotalCapacity + "MW");

    });

    $.when(wfa.ProjectAnalysis.LoadData()).done(function(){
        setTimeout(function(){
            // app.loading(false);
        },500)
    })
}

wfa.RefreshGrid = function() {
	// $('.grid-custom').each(function(){
	// 	var grid = $(this).data("kendoGrid");
	// 	grid.refresh();
	// });
    var ids = $(".panel-body").find(".nav-tabs").find(".active")[0].id
    switch (ids) {
        case "tTurbine1Analysis":
            $("#gridTurbine1").data("kendoGrid").refresh();
            break;
        case "tTurbine2Analysis":
            $("#gridTurbine2").data("kendoGrid").refresh();
            break;
        default :
            $("#gridProject").data("kendoGrid").refresh();
            break;
    }
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

wfa.setTurbines = function (id,elemid) {
    var datavalue = [];

    var allturbine = {}
    allturbine.value = "All Turbines";
    allturbine.text = "All Turbines";
    datavalue.push(allturbine);

    var project = $("#"+id).data("kendoDropDownList").value();
    $.each(turbines, function(idx, val) {
        if (project == "") {
            var data = {};
            data.value = val.Value;
            data.text = val.Turbine;
            datavalue.push(data);
        }else if (project == val.Project){
            var data = {};
            data.value = val.Value;
            data.text = val.Turbine;
            datavalue.push(data);
        }
    });

    wfa.TurbineList(datavalue);
    $("#"+elemid).data('kendoMultiSelect').setDataSource(new kendo.data.DataSource({ data: wfa.TurbineList() }));
    $("#"+elemid).data('kendoMultiSelect').value(["All Turbines"]);

}
wfa.resetStatus = function(){
    wfa.isProjectLoaded(false);
    wfa.isTurbine1Loaded(false);
    wfa.isTurbine2Loaded(false);
}
wfa.showFilter = function(project, turbine1, turbine2, id){
    wfa.isProjectTab(project);
    wfa.isTurbine1Tab(turbine1);
    wfa.isTurbine2Tab(turbine2);

    switch (id) {
        case "tProjectAnalysis":
            if(!wfa.isProjectLoaded()) {
                wfa.getDataAvailableInfo();
                app.loading(true);
                wfa.isProjectLoaded(true);
                wfa.LoadData();
            }
            break;
        case "tTurbine1Analysis":
            if(!wfa.isTurbine1Loaded()) {
                wfa.getDataAvailableInfo('projectTurbine1List');
                app.loading(true);
                wfa.isTurbine1Loaded(true);
                wfa.Turbine1Analysis.LoadData();
            }
            break;
        case "tTurbine2Analysis":
            if(!wfa.isTurbine2Loaded()) {
                wfa.getDataAvailableInfo('projectTurbine2List');
                app.loading(true);
                wfa.isTurbine2Loaded(true);
                wfa.Turbine2Analysis.LoadData();
            }
            break;
        default :
            //project udah otomatis ke load
            break;
    }
}

// initiate value for projects & turbines
$.each(projects, function(idx, val) {
    var data = {};
    data.value = val.Value;
    data.text = val.Name;
    wfa.ProjectList.push(data);
	// wfa.ProjectList.push(addOption(val));
});
$.each(turbines, function(idx, val) {
    var data = {};
    data.value = val.Value;
    data.text = val.Turbine;
    wfa.TurbineList.push(data);
	// wfa.TurbineList.push(addOption(val));
});

wfa.Turbine2Chart = function(data, chartSeries) {
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
            visible: false,
            position: "bottom"
        },
        seriesColors: colors,
        seriesDefaults: {
            type: "line",
            style: "smooth",
            markers: {
                visible: false
            }
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
};

$(function(){
    wfa.isProjectTab(true);

	$(window).on("resize orientationchange", function () {        
	    wfa.RefreshGrid();
	});



    $('#btnRefresh').on('click', function () {
        wfa.resetStatus();
        $('.nav').find('li.active').find('a').trigger("click");
    });

    $("#turbine1List").kendoMultiSelect({
        dataSource: wfa.TurbineList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {wfa.checkTurbine('turbine1List')}, 
        suggest: true 
    }); 

    $("#turbine2List").kendoMultiSelect({
        dataSource: wfa.TurbineList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {wfa.checkTurbine('turbine2List')}, 
        suggest: true
    }); 

    wfa.setTurbines('projectTurbine1List','turbine1List');
    wfa.setTurbines('projectTurbine2List','turbine2List');


	wfa.LoadData();

    
});



