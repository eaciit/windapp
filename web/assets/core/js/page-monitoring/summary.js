'use strict';

viewModel.Summary = new Object();
var summary = viewModel.Summary;

vm.currentMenu('Summary');
vm.currentTitle('Summary');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Summary', href: viewModel.appName + 'page/monitoringsummary' },
]);

summary.isFirstOverAll = ko.observable(true);
summary.isFirstAllFarms = ko.observable(true);


summary.LoadOverAll = function(){
	if(summary.isFirstOverAll() == true){
		app.loading(true);
	    $.when(bpc.getWeather()).done(function(){
	    	bpc.getData();	
	    	summary.isFirstOverAll(false);
	    	app.loading(false);
	    	bpc.refresh();
	    });
	}
	
}

summary.LoadAllFarms = function(){
	if(summary.isFirstAllFarms() == true){
		app.loading(true);
		page.getData();
    	summary.isFirstAllFarms(false);
    	setInterval(page.getData, 5000);
	}
}

$(function() {
	summary.LoadOverAll();
});

