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
	IsTurbine: true,
});

// events for turbine collaboration page
TbCol.ShowModal = function(mode) {
	$('#mdlTbColab').modal(mode);
};
TbCol.OpenForm = function() {
	setTimeout(function(){
		// TbCol.GenerateGrid(TbCol.TurbineId(),TbCol.Project());
		$('#mdlTbColab').modal('show');
	},200);

};
TbCol.CloseForm = function() {
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
}
TbCol.GenerateGrid = function(turbine, project){
	app.loading(true);
	var param = {
		Project : project,
		Turbine : turbine,
		Take : null,
	}
	toolkit.ajaxPost(viewModel.appName + 'turbinecollaboration/getlatest', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        $("#gridHistory").html("");
		$("#gridHistory").kendoGrid({
	        dataSource: {
                data: res.data,
                pageSize: 10,
	            schema: {
	                data: function(res) {
	                    return res.data
	                },
	            },
            },
	        filterable: false,
	        sortable: true,
	        pageable: true,
	        columns: [{
	                field: "Remark",
	                title: "Remark",
	                headerAttributes: { style: 'text-align: center;' },
	            }, {
	                field: "Date",
	                title: "Timestamp",
	                template: "#= kendo.toString(new Date(Date), 'dd-MMM-yyyy HH:mm') #",
	            }
	        ],
	        dataBound: function(){
	        	app.loading(false);
	        }
	    });
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
			Status : TbCol.Status(),
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