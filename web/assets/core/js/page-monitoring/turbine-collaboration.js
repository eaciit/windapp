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
TbCol.Status = ko.observable('');

// variabel to set current data if any edit feature
TbCol.CurrentData = ko.observable({
	Id: '',
	TurbineId: '',
	TurbineName: '',
	Project: '',
	Feeder: '',
	Date: new Date(),
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
    var param = {
    		TurbineId : TbCol.TurbineId(),
			TurbineName : TbCol.TurbineName(),
			Feeder : TbCol.Feeder(),
			Project : TbCol.Project(),
			Date : TbCol.CurrentData().Date,
			Status : '',
			Remark : TbCol.CurrentData().Remark,
    }

    toolkit.ajaxPost(viewModel.appName + 'turbinecollaboration/save', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        swal({ title: "Saved", type: "success" });
       	TbCol.CloseForm();
    }, function (err) {
        toolkit.showError(err.responseText);
    });
	
}