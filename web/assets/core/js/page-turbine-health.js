'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Turbine Health');
vm.currentTitle('Turbine Health');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Turbine Health', href: viewModel.appName + 'page/turbinehealth' }]);


$(document).ready(function () {
    fa.LoadData();
    $('#btnRefresh').on('click', function () {
        fa.LoadData();
    });

    setTimeout(function () {
        fa.LoadData();
    }, 1000);
    app.loading(false);

    $('#projectList').kendoDropDownList({
		change: function () {  
			var project = $('#projectList').data("kendoDropDownList").value();
			fa.populateTurbine(project);
		}
	});
});