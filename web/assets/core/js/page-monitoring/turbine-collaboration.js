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
TbCol.ProjectFeeder = ko.observable('');
TbCol.IsTurbine = ko.observable(true);
TbCol.IconStatus = ko.observable('');
TbCol.haveRemark = ko.observable(false);

// variabel to set current data if any edit feature
TbCol.CurrentData = ko.observable({
	_id: '',
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
	IsTurbine: true,
});

// events for turbine collaboration page
TbCol.ShowModal = function(mode) {
	$('#mdlTbColab').modal(mode);
};
TbCol.OpenForm = function() {
	if(typeof summary !== 'undefined' ){
		summary.abortAll(requests);
	}

	if(typeof bp !== 'undefined'){
		bp.abortAll(requests);
	}	
	
	var project = TbCol.Project() == "" ? TbCol.ProjectFeeder() : TbCol.Project();
	$.when(TbCol.GenerateGrid(TbCol.TurbineId(),project,TbCol.Feeder())).done(function(){
		$('#mdlTbColab').modal('show');
	});
};
TbCol.CloseForm = function() {
	if(typeof summary !== 'undefined'){
		$allFarmsInterval = bpc.refresh();
	}

	if(typeof bp !== 'undefined'){
		setInterval(bp.GetData, 5000);
	}
	
	$('#mdlTbColab').modal('hide');
};


TbCol.ResetData = function(){
	TbCol.TurbineId('');
	TbCol.TurbineName('');
	TbCol.UserId('');
	TbCol.UserName('');
	TbCol.Project('');
	TbCol.Feeder('');
	TbCol.Status('');
	TbCol.IsTurbine(true);
	TbCol.IconStatus('');
}
TbCol.GenerateGrid = function(turbine, project,feeder){
	app.loading(true);
	var param = {
		Project : project,
		Turbine : turbine,
		Feeder : feeder,
		Take : 1,
	}
	toolkit.ajaxPost(viewModel.appName + 'turbinecollaboration/getlatest', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        var results = res.data;

        if(results !== null){
        	$.each(TbCol.CurrentData(), function(key, val){
				TbCol.CurrentData()[key] = results[key] == undefined ? '' : results[key];
			});

			// $("#date").data("kendoDateTimePicker").value(results.Date);
			$("#status").val(results.Status);
			$("#remark").val(results.Remark);
			
        }

        TbCol.haveRemark(results._id == "" ? false : true);

        app.loading(false);   
	});
}
TbCol.Save = function() {
	app.loading(true);
    var param = {
    		TurbineId : TbCol.TurbineId(),
			TurbineName : TbCol.TurbineName(),
			Feeder : TbCol.Feeder(),
			Project : (TbCol.Project() == '' ? TbCol.ProjectFeeder() : TbCol.Project()) ,
			Date : $("#date").data("kendoDateTimePicker").value(),
			Status :TbCol.CurrentData().Status,
			Remark : TbCol.CurrentData().Remark,
    }

    toolkit.ajaxPost(viewModel.appName + 'turbinecollaboration/save', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        swal({ title: "Saved", type: "success" });
       	TbCol.CloseForm();
       	app.loading(false);
    }, function (err) {
        toolkit.showError(err.responseText);
    });
	
}

TbCol.Delete = function() {
	app.loading(true);
    var param = {
    		TurbineId : TbCol.TurbineId(),
			TurbineName : TbCol.TurbineName(),
			Feeder : TbCol.Feeder(),
			Project : (TbCol.Project() == '' ? TbCol.ProjectFeeder() : TbCol.Project()) ,
			Date : $("#date").data("kendoDateTimePicker").value(),
			Status :TbCol.CurrentData().Status,
			Remark : TbCol.CurrentData().Remark,
			Id: TbCol.CurrentData()._id,
			IsDeleted: true,
    }

    toolkit.ajaxPost(viewModel.appName + 'turbinecollaboration/save', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        swal({ title: "Deleted", type: "success" });
       	TbCol.CloseForm();
       	app.loading(false);
    }, function (err) {
        toolkit.showError(err.responseText);
    });
	
}