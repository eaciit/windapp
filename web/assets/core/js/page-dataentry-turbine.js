'use strict';

viewModel.AnalyticDgrScada = new Object();
var page = viewModel.AnalyticDgrScada;

vm.currentMenu('Turbine Master');
vm.currentTitle('Turbine Master');
vm.breadcrumb([{ title: 'Data Entry', href: '#' },{ title: 'Turbine Master', href: viewModel.appName + 'page/dataentryturbine' }]);

page.CurrentData = ko.observable({ Id: '', Model: '', WindSpeed: 0, Power1: 0 });

var Data = {
	LoadData: function() {
		app.loading(true);
		var isValid = fa.LoadData();
        if(isValid) {
            this.InitGrid();
        }
	},
	InitGrid: function() {
		var param = {
			Project: fa.project,
			Turbine: fa.turbine(),
		};

		$('#gridDataEntryTurbine').html("");
		$('#gridDataEntryTurbine').kendoGrid({
	      dataSource: {
	        serverPaging: true,
	        serverSorting: true,
	        transport: {
	          read: {
	            url: viewModel.appName + "dataentryturbine/getlist",
	            type: "POST",
	            data: param,
	            dataType: "json",
	            contentType: "application/json; charset=utf-8"
	          },
	          parameterMap: function(options) {
	            return JSON.stringify(options);
	          }
	        },
	        pageSize: 10,
	        schema: {
	          data: function(res) {
                    return res.data.Data
                },
                total: function(res) {
	                if (!app.isFine(res)) {
	                    return;
	                }

	                return res.data.Total;
	            }
	        },
	        sort: [
	        	{ field: 'Project', dir: 'asc' },
	        	{ field: 'TurbineName', dir: 'asc' },
			],
			group: { field: "Project", dir: "desc" }
	      },
	      sortable: true,
	      filterable: false,
	      pageable: {
                pageSize: 10,
                input: true, 
          },
	      scrollable: false,
	      columns: [
	        { title: "Turbine", field: "TurbineName", width: 120, attributes:{ class:"align-center" },locked: true,},
	        { title: "Feeder", field: "Feeder", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, width: 120 },
	        { title: "Latitude", field: "Latitude", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120 },
	        { title: "Longitude", field: "Longitude", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120 },
	        { title: "Elevation", field: "Elevation", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120 },
	        { title: "Capacity (mW)", field: "Capacitymw", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120 }
	      ],
	      dataBound : function(){
	      		app.loading(false);
	      }
	    });

	    
	},
	New: function(){
		page.CurrentData({ Id: data.Id, Model: data.Model, WindSpeed: data.WindSpeed, Power1: data.Power1 });
	    this.ShowForm('show');
	    this.ResetValidation(".modal-body");
	},
	Edit: function(id) {
		this.ResetValidation(".modal-body");
	    var url = viewModel.appName + "dataentrypowercurve/getdata";
	    var param = { id: id };
	    var $this = this;

		toolkit.ajaxPost(url, param, function(res) {
			if (!app.isFine(res)) {
	            return;
	        }

	        var data = res.data;

            page.CurrentData({ ID: data.ID, Model: data.Model, WindSpeed: data.WindSpeed, Power1: data.Power1 });
			$this.ShowForm('show');
        });
	},
	Delete: function(id) {
	    var $this = this;
	    swal({
	      title: "Are you sure?",
	      text: "Your will not be able to recover this data",
	      type: "warning",
	      showCancelButton: true,
	      confirmButtonClass: "btn-danger",
	      confirmButtonText: "Yes, delete it!",
	      closeOnConfirm: false
	    },
	    function(res){
	      if(res) {
	       	var url = viewModel.appName + "dataentrypowercurve/delete";
	        var param = { id: id };

		    toolkit.ajaxPost(url, param, function(data) {
		    if (!app.isFine(data)) {
	            return;
	        }
	          swal('Success', 'This data has been deleted!', 'success');
	          $this.LoadData();
	          app.isLoading(false);
	        });
	      }
	    });
	},
	Save: function() {
		if (!this.CheckValidation(".modal-body")) {
	      return;
	    }

	    var url = viewModel.appName + "dataentrypowercurve/save";
	    var data = page.CurrentData();
	    var param = {
	    	ID: data.ID,
	    	Model: data.Model,
	    	WindSpeed: parseFloat(data.WindSpeed),
	    	Power1: parseFloat(data.Power1)
	    };
	    var $this = this;

	    toolkit.ajaxPost(url, param, function(data) {
	    	if (!app.isFine(data)) {
	            return;
	        }else{
	        	$this.ShowForm('hide');
		        swal('Success', 'Data has been saved successfully!', 'success');
		        $this.LoadData();
	        }
            
        });
	},
	ResetValidation: function (selectorID) {
	    var $form = $(selectorID).data("kendoValidator");
	    if ($form == undefined) {
	        $(selectorID).kendoValidator();
	        $form = $(selectorID).data("kendoValidator");
	    }

	    $form.hideMessages();
	},
	CheckValidation: function (selector) {
	    this.ResetValidation(selector);
	    var $validator = $(selector).data("kendoValidator");
	    return ($validator.validate());
	},
	ShowForm: function(showhide) {
		page.ShowModal('modalForm', showhide);
	}
};

page.ShowModal =  function(modalId, showhide) {
	if(showhide=='show') {
	  $('#'+modalId).appendTo("body").modal({
	          backdrop: 'static',
	          keyboard: false, 
	          show: showhide
	      });
	} else {
	  $('#'+modalId).modal('hide');
	}
}

page.hideElement = function(){
	$("#periodList").closest(".k-widget").hide();
    $("#dateStart").closest(".k-widget").hide();
    $("#dateEnd").closest(".k-widget").hide();
    $(".input-group").find("label").hide();
}



$(function (){
	page.hideElement();

	setTimeout(function(){
		Data.LoadData();
	},500);
	
	
	$('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        Data.LoadData();
    });

    $('#projectList').kendoDropDownList({
        change: function () {  
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
        }
    });

});

