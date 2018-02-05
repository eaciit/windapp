viewModel.DataAvailability = new Object();
var page = viewModel.DataAvailability;

vm.currentMenu('Data Availability');
vm.currentTitle('Data Availability');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Data Availability', href: viewModel.appName + 'page/dataavailability' }]);

page.isExpanded = ko.observable();
page.categoryHeader = ko.observableArray();
var categoryHeaderDay;

var colspan;
var colspanDay;

page.widthColumn = ko.observable();
var widthColumnDay;

page.dataAvail = ko.observableArray();
var dataAvailDay;


page.hideFilter = function(){
    $("#periodList").closest(".k-widget").hide();
    $("#dateStart").closest(".k-widget").hide();
    $("#dateEnd").closest(".k-widget").hide();
    $(".control-label:contains('Period')").hide();
    $(".control-label:contains('to')").hide();
}

page.getData = function(){
	fa.LoadData();
	di.getAvailDate();
	var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
    };

    var dataAvailReq = toolkit.ajaxPost(viewModel.appName + "dataavailability/getdataavailability", param, function (res) {
        if (!app.isFine(res)) {
                return;
            }

        page.dataAvail(res.data.Data);

        page.categoryHeader(res.data.Month);
		colspan = page.categoryHeader().length;
		page.widthColumn((90 / colspan) + "%");
		page.createView();
    });

    $.when(dataAvailReq).done(function(){
    	setTimeout(function(){
    		app.prepareTooltipster();
    		app.loading(false);
    	}, 100);
    });
}

page.monthDetail = function(month) {
	app.loading(true);
	var param = {
        period: month,
        breakdown: "daily",
        turbine: fa.turbine(),
        project: fa.project,
    };

    var dataAvailReq = toolkit.ajaxPost(viewModel.appName + "dataavailability/getdataavailability", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        if(res.data.Data[0].Data != undefined || res.data.Data[0].Data != null) {
        	dataAvailDay= res.data.Data;
	        categoryHeaderDay = res.data.Month;
			colspanDay = categoryHeaderDay.length;
			widthColumnDay = (90 / colspanDay) + "%";
			page.createViewDaily(month);
        } else {
        	swal({
	            title: "Warning",
	            type: "warning",
	            text: "Data is not available for "+ month,
	        }, function () {
	            app.loading(false);
	        });
        }
    });

    $.when(dataAvailReq).done(function(){
    	setTimeout(function(){
    		app.prepareTooltipster();
    		app.loading(false);
    	}, 100);
    });
}

page.createViewDaily = function(month){
	$("#tableContent").html("");
	$("#tableHeader").html("");
	$("#tableHeader").append('<td width="10%" class="border-right clickable" colspan="2" onclick="page.createView()"><a role="tab" data-toggle="tab" class="btn-back"><i class="fa fa-reply" aria-hidden="true"></i> Back </a> <strong style="float: right;">'+month.substring(0,3)+" "+month.slice(-4)+'</strong></td>');

	$.each(categoryHeaderDay, function(id, ress){
		var tdHeader = ' <td width="'+widthColumnDay+'" class="text-month" ><strong>'+ress+'</strong></td>';

		$("#tableHeader").append(tdHeader);
	});

	$.each(dataAvailDay, function(key, value){
		var progressData = "";
		if (value!=null){
			$.each(value.Data, function(i, val){
				progressData += '<div aria-hidden="true" class="tooltipster tooltipstered '+val.class+'" style = "width:'+val.value+';opacity:'+val.opacity+'"  title = "'+val.tooltip+' : '+ kendo.toString(val.floatval,'n2') + ' %" role="progressbar"></div>'
				
			});

			var icon = "";
			if(value.Turbine.length > 0){
				icon = '<i class="fa fa-chevron-right"></i><i class="fa fa-chevron-down" style="display:none;"></i>';
			}
			var master = '<tr class="clickable" data-toggle="collapse" data-target=".row'+key+'">'+
							'<td>'+icon+'</td>'+
							'<td class="border-right"><strong>'+value.Category+'</strong></span></td>'+
							'<td colspan='+colspanDay+'>'+
									'<div class="progress">'+progressData+'</div>'+
							'</td>'+
						'</tr>';
			
			$.each(value.Turbine, function(index, res){
				var progressDataDetails = "";

				$.each(res.details, function(idx, result){
					progressDataDetails += '<div aria-hidden="true" class="tooltipster tooltipstered '+result.class+'" style = "width:'+result.value+';opacity:'+result.opacity+'"  title = "'+result.tooltip+' : '+ kendo.toString(result.floatval,'n2') + ' %" role="progressbar"></div>'
					
				});

				var details = '<tr class="collapse details row'+key+'">'+
					'<td></td>'+
					'<td class="border-right" style="padding-left:30px">'+res.TurbineName+'</span></td>'+
					'<td colspan='+colspanDay+'>'+
							'<div class="progress">'+progressDataDetails+'</div>'+
					'</td>'+
				'</tr>';

				master += details;

			});

			$("#tableContent").append(master);
		}

		
	});

	$('.collapse').on('shown.bs.collapse', function(){
		$(this).parent().find(".fa-chevron-right").removeClass("fa-chevron-right").addClass("fa-chevron-down");
	}).on('hidden.bs.collapse', function(){
		$(this).parent().find(".fa-chevron-down").removeClass("fa-chevron-down").addClass("fa-chevron-right");
	});

}

page.createView = function(){
	$("#tableContent").html("");
	$("#tableHeader").html("");
	$("#tableHeader").append('<td width="10%" class="border-right" colspan="2">&nbsp;</td>');

	$.each(page.categoryHeader(), function(id, ress){
		var month = ress.split(" ")[0].slice(0, 3);
		var tdHeader = ' <td width="'+page.widthColumn()+'" class="text-month label label-primary clickable" onclick="page.monthDetail(\'' + ress + '\')" ><strong>'+
		month+'</strong></td>';

		$("#tableHeader").append(tdHeader);
	});

	$.each(page.dataAvail(), function(key, value){
		var progressData = "";
		if (value!=null){
			$.each(value.Data, function(i, val){
				progressData += '<div aria-hidden="true" class="tooltipster tooltipstered '+val.class+'" style = "width:'+val.value+'"  title = "'+val.tooltip+'" role="progressbar"></div>'
				
			});

			var icon = "";
			if(value.Turbine.length > 0){
				icon = '<i class="fa fa-chevron-right"></i><i class="fa fa-chevron-down" style="display:none;"></i>';
			} 
			var master = '<tr class="clickable" data-toggle="collapse" data-target=".row'+key+'">'+
							'<td>'+icon+'</td>'+
							'<td class="border-right"><strong>'+value.Category+'</strong></span></td>'+
							'<td colspan='+colspan+'>'+
									'<div class="progress">'+progressData+'</div>'+
							'</td>'+
						'</tr>';


			$.each(value.Turbine, function(index, res){
				var progressDataDetails = "";

				$.each(res.details, function(idx, result){
					progressDataDetails += '<div aria-hidden="true" class="tooltipster tooltipstered '+result.class+'" style = "width:'+result.value+'"  title = "'+result.tooltip+'" role="progressbar"></div>'
					
				});

				var details = '<tr class="collapse details row'+key+'">'+
					'<td></td>'+
					'<td class="border-right" style="padding-left:30px">'+res.TurbineName+'</span></td>'+
					'<td colspan='+colspan+'>'+
							'<div class="progress">'+progressDataDetails+'</div>'+
					'</td>'+
				'</tr>';

				master += details;

			});
			$("#tableContent").append(master);
		}

		
	});

	$('.collapse').on('shown.bs.collapse', function(){
		$(this).parent().find(".fa-chevron-right").removeClass("fa-chevron-right").addClass("fa-chevron-down");
	}).on('hidden.bs.collapse', function(){
		$(this).parent().find(".fa-chevron-down").removeClass("fa-chevron-down").addClass("fa-chevron-right");
	});

}

function sticky_relocate() {
    var window_top = $(window).scrollTop();

    $("#header-fixed").html("");
    var tableOffset = $("#table-dataavailability").offset().top;
	var $header = $("#table-dataavailability > thead").clone();
	var $fixedHeader = $("#header-fixed").append($header);

    var offset = $(window).scrollTop();

    if (offset >= tableOffset && $fixedHeader.is(":hidden")) {
        $fixedHeader.show();
    }
    else if (offset < tableOffset) {
        $fixedHeader.hide();
    }
}

$(function () {
    $(window).scroll(sticky_relocate);
    sticky_relocate();
    
	$('#btnRefresh').on('click', function () {
		app.loading(true);
		fa.checkTurbine();
		page.getData();
    });

	page.hideFilter();
	setTimeout(function() {
		fa.LoadData();
		page.getData();
	},200);

	$('#projectList').kendoDropDownList({
		change: function () {  
			di.getAvailDate();
			var project = $('#projectList').data("kendoDropDownList").value();
			fa.populateTurbine(project);
		}
	});

	setTimeout(function() {
		$('#projectList').data("kendoDropDownList").select(0);
	}, 100);
});
