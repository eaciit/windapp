'use strict';

viewModel.Summary = new Object();
var summary = viewModel.Summary;

vm.currentMenu('Summary');
vm.currentTitle('Summary');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Summary', href: viewModel.appName + 'page/monitoringsummary' },
]);


