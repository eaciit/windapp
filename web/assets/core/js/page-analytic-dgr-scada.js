'use strict';

vm.currentMenu('Dgr vs Scada');
vm.currentTitle('Dgr vs Scada');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'DGR vs Scada', href: viewModel.appName + 'page/analyticdgrscada' }]);

viewModel.AnalyticDgrScada = new Object();
var ads = viewModel.AnalyticDgrScada;

ads.modelList = ko.observableArray([
	{ "value": "Regen", "text": "Regen" },
	{ "value": "Suzlon", "text": "Suzlon" },
]);

ads.model = ko.observable("Regen");
ads.project = ko.observable("Tejuva");
ads.turbine = ko.observable("All Turbine");

/*var dataSource = [
	{"desc" : 'Power (MW)' , "dgr" : 118 , "scada" : 471, "difference" : 253},
	{"desc" : 'Energy (MWH)' , "dgr" : 29 , "scada" : 118, "difference" : 88},
	{"desc" : 'Avg. Wind Speed (m/s)' , "dgr" : 7.08 , "scada" : 7.08, "difference" : 0.00},
	{"desc" : 'Downtime (Hrs)' , "dgr" : 0.42 , "scada" : 25.26	, "difference" : 24.84},
	{"desc" : 'PLF' , "dgr" : 23.38 , "scada" : 23.38, "difference" : 0.00},
	{"desc" : 'Grid Availability' , "dgr" : 0.00 , "scada" : 81.55, "difference" : 81.55},
	{"desc" : 'Machine Availability' , "dgr" : 84.46 , "scada" : 75.82, "difference" : 8.64},
	{"desc" : 'True Availability' , "dgr" : 84.46 , "scada" : 74.03, "difference" : 10.43},
];*/

var Data = {
	LoadData: function () {
		setTimeout(function () {
			fa.LoadData();
			Data.InitSummaryGrid();
			fa.getProjectInfo();
		}, 1000);
	},
	InitSummaryGrid: function () {
		var param = {
			period: fa.period,
			Turbine: fa.turbine,
			DateStart: fa.dateStart,
			DateEnd: fa.dateEnd,
			Project: fa.project
		};

		app.loading(true);
		$("#gridSummaryDgrScada").kendoGrid({
			theme: "flat",
			columns: [
				{ title: " ", field: "desc", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-left" }, width: 150 },
				{ title: "DGR", field: "dgr", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" }, template: "#if(desc== 'PLF'){# #: kendo.toString(dgr, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(dgr, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(dgr, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(dgr, 'N2') # % #}else {# #: kendo.toString(dgr, 'N2') # #}#" },
				{ title: "Scada", field: "scada", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" }, template: "#if(desc== 'PLF'){# #: kendo.toString(scada, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(scada, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(scada, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(scada, 'N2') # % #}else {# #: kendo.toString(scada, 'N2') # #}#" },
				{ title: "Difference", field: "difference", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-right" }, template: "#if(desc== 'PLF'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(difference, 'N2') # % #}else {# #: kendo.toString(difference, 'N2') # #}#" },
			],
	        /*dataSource: {
	            data : dataSource,
	        }*/
			dataSource: {
				serverPaging: false,
				serverSorting: false,
				serverFiltering: false,
				transport: {
					read: {
						url: viewModel.appName + "analyticdgrscada/getdata",
						type: "POST",
						data: param,
						dataType: "json",
						contentType: "application/json; charset=utf-8"
					},
					parameterMap: function (options) {
						return JSON.stringify(options);
					}
				},
				schema: {
					model: {
						fields: {
							AlarmOkTime: { type: "number" },
							OkTime: { type: "number" },
							Power: { type: "number" },
							PowerLost: { type: "number" },
						}
					},
					data: function (res) {
						app.loading(false);
						if (!app.isFine(res)) {
							return;
						}
						return res.data
					}
				},
			}
		});
		app.loading(false);
	},
};

$(function () {
	Data.LoadData();
	$('#btnRefresh').on('click', function () {
		app.loading(true);
		setTimeout(function () {
			Data.LoadData();
			app.loading(false);
		}, 200);
	});
});

