'use strict';

// initiate mvvm for turbine collaboration
viewModel.TurbineCollaboration = new Object();
var TbCol = viewModel.TurbineCollaboration;

// initiate any variables for turbine collaboration mvvm
TbCol.TurbineId = ko.observable('');
TbCol.TurbineName = ko.observable('');
TbCol.UserId = ko.observable('');
TbCol.UserName = ko.observable('');
TbCol.Project = ko.observable('');
TbCol.Feeder = ko.observable('');

// variabel to set current data if any edit feature
TbCol.CurrentData = ko.observable({
	Id: '',
	TurbineId: '',
	TurbineName: '',
	Project: '',
	Feeder: '',
	Date: '',
	Time: '',
	UserId: '',
	UserName: '',
	Status: '',
	Remark: '',
});

// events for turbine collaboration page
TbCol.ShowModal = function(mode) {
	$('#mdlTbColab').modal(mode);
};
TbCol.OpenForm = function() {
	$('#mdlTbColab').modal('show');
};
TbCol.CloseForm = function() {
	$('#mdlTbColab').modal('hide');
};
TbCol.Save = function() {

};