'use strict';

viewModel.ByProjectCustom = new Object();
var bpc = viewModel.ByProjectCustom;


vm.currentMenu('By Project');
vm.currentTitle('By Project');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'By Project', href: viewModel.appName + 'page/monitoringbyproject' }]);


bpc.projectList = ko.observableArray([]);


bpc.getData = function(){
	
}

$(function() {

	$('#savedViews').kendoDropDownList({
		data: [],
		dataValueField: 'value',
		dataTextField: 'text',
		change: function () {  },
	});

	$( "#sortable" ).sortable();
    $( "#sortable" ).disableSelection();
});