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
	// app.loading(true);
	$.when(bpc.getWeather()).done(function(){
    	bpc.getData();	
    	summary.isFirstAllFarms(false);	
    	$allFarmsInterval = bpc.refresh();
    });
	
}

summary.LoadOverAll = function(){
	app.loading(true);
	page.getData();
	summary.isFirstOverAll(false);
	$overAllInterval = setInterval(page.getData, $intervalTime);
}

summary.SelectMode = function(type) {
	summary.abortAll(requests);
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

summary.abortAll = function(requests) {
     var length = requests.length;
     while(length--) {
         requests[length].abort && requests[length].abort();  // the if is for the first case mostly, where array is still empty, so no abort method exists.
     }
}

$(function() {
	summary.SelectMode('overall');
});

