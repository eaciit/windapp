'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Reporting');
vm.currentTitle('Reporting');
vm.breadcrumb([{ title: 'Reporting', href: viewModel.appName + 'page/reporting' }]);


$(document).ready(function () {
    fa.LoadData();
    $('#btnRefresh').on('click', function () {
        fa.LoadData();
    });

    setTimeout(function () {
        fa.LoadData();
    }, 1000);
    app.loading(false);
});