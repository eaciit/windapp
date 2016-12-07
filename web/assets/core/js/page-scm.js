'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Inventory / SCM Management');
vm.currentTitle('Inventory / SCM Management');
vm.breadcrumb([{ title: 'O&M', href: '#' }, { title: 'Inventory / SCM Management', href: viewModel.appName + 'page/scmmanagement' }]);


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