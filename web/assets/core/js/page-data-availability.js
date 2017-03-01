viewModel.AnalyticKpi = new Object();
var page = viewModel.DataAvailability;

vm.currentMenu('Data Availability');
vm.currentTitle('Data Availability');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'KPI Table', href: viewModel.appName + 'page/dataavailability' }]);

$(function () {
    setTimeout(function() {
        app.loading(false);    
    }, 200);
});
