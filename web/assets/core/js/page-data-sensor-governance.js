'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Data and Sensor Governance');
vm.currentTitle('Data and Sensor Governance');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Data and Sensor Governance', href: viewModel.appName + 'page/datasensorgovernance' }]);


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