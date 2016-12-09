'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Issue Tracking');
vm.currentTitle('Issue Tracking');
vm.breadcrumb([{ title: 'O&M', href: '#' }, { title: 'Issue Tracking', href: viewModel.appName + 'page/issuetracking' }]);


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