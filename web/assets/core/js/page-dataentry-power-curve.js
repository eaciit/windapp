'use strict';

viewModel.AnalyticDgrScada = new Object();
var page = viewModel.AnalyticDgrScada;

page.CurrentData = ko.observable({ Id: '', Model: '', WindSpeed: 0, Power1: 0 });

var Data = {
	LoadData: function() {
		this.InitGrid();
	},
	InitGrid: function() {
		var param = {};
		$('#gridDataEntryPowerCurve').html("");
		$('#gridDataEntryPowerCurve').kendoGrid({
	      dataSource: {
	        serverPaging: true,
	        serverSorting: true,
	        transport: {
	          read: {
	            url: viewModel.appName + "dataentrypowercurve/getlist",
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
	          data: "Data",
	          total: "Total"
	        },
	        sort: [
				{ field: 'WindSpeed', dir: 'asc' }
			],
	      },
	      groupable: false,
	      sortable: true,
	      filterable: false,
	      pageable: {
                pageSize: 10,
                input: true, 
          },
	      scrollable: false,
	      columns: [
	        { title: "Model", field: "Model", width: 120 },
	        { title: "Avg. Wind Speed<br>(m/s)", field: "WindSpeed", headerAttributes: { style:"text-align: right" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120 },
	        { title: "Standard Power<br>(KW)", field: "Power1", headerAttributes: { style:"text-align: right" }, attributes:{ class:"align-center" }, format: "{0:n0}", width: 120 },
	        { title: "Control",headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" },template: "<div class='middles'><a href=\"javascript:Data.Edit('#: ID #')\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i>&nbsp;&nbsp;Edit</a></div>",width: 100}
	      ]
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

		toolkit.ajaxPost(url, param, function(data) {
			if (!app.isFine(data)) {
	            return;
	        }
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
	    	Id: data.ID,
	    	Model: data.Model,
	    	WindSpeed: parseFloat(data.WindSpeed),
	    	Power1: parseFloat(data.Power1)
	    };
	    var $this = this;

	    toolkit.ajaxPost(url, param, function(data) {
	    	if (!app.isFine(data)) {
	            return;
	        }
            if(data=="") {
		        $this.ShowForm('hide');
		        swal('Success', 'Data has been saved successfully!', 'success');
		        $this.LoadData();
		    } else {
		        swal("Warning", data, "error");
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

vm.currentMenu('Data Entry Power Curve');
vm.currentTitle('Data Entry Power Curve');
vm.breadcrumb([{ title: 'Data Entry', href: '#' },{ title: 'Power Curve', href: viewModel.appName + 'page/dataentrypowercurve' }]);


$(function (){
	Data.LoadData();
});

