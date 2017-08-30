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

var $overAllInterval = false, $allFarmsInterval = false, $intervalTime = 2500;

summary.LoadAllFarms = function(){
	//if(summary.isFirstAllFarms() == true){
		$.when(bpc.getWeather()).done(function(){
	    	bpc.getData();	
	    	summary.isFirstAllFarms(false);	
	    	$allFarmsInterval = bpc.refresh();
	    });
	//}
	
}

summary.LoadOverAll = function(){
	//if(summary.isFirstOverAll() == true){
		app.loading(true);
		page.getData();
    	summary.isFirstOverAll(false);
    	$overAllInterval = setInterval(page.getData, $intervalTime);
	//}
}

summary.SelectMode = function(type) {
	if(type	== 'overall') {
		clearInterval($allFarmsInterval);
		$allFarmsInterval = false;
		summary.LoadOverAll();
	} else {
		clearInterval($overAllInterval)
		$overAllInterval = false;
		summary.LoadAllFarms();
	}
}

$(function() {
	summary.SelectMode('overall');
});

