'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('DIY View');
vm.currentTitle('DIY View');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'DIY View', href: viewModel.appName + 'page/diyview' }]);


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