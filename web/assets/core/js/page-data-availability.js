viewModel.DataAvailability = new Object();
var page = viewModel.DataAvailability;

vm.currentMenu('Data Availability');
vm.currentTitle('Data Availability');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Data Availability', href: viewModel.appName + 'page/dataavailability' }]);

page.isExpanded = ko.observable();
page.categoryHeader = ko.observableArray(["Jan", "Feb", "Mar", "Apr","May","Jun","Jul"]);

var colspan = page.categoryHeader().length;

page.widthColumn = ko.observable((90 / colspan) + "%");
page.data = ko.observableArray([]);
page.hideFilter = function(){
    $("#periodList").closest(".k-widget").hide();
    $("#dateStart").closest(".k-widget").hide();
    $("#dateEnd").closest(".k-widget").hide();
    $(".control-label:contains('Period')").hide();
    $(".control-label:contains('to')").hide();
}

page.getData = function(){
	fa.LoadData();
	var param = {
        period: fa.period,
        dateStart: fa.dateStart,
        dateEnd: fa.dateEnd,
        turbine: fa.turbine,
        project: fa.project,
    };

    toolkit.ajaxPost(viewModel.appName + "dataavailability/getdataavailability", param, function (res) {
        if (!app.isFine(res)) {
                return;
            }

        page.data(res.data.Data);
        page.categoryHeader(res.data.Month);
		colspan = page.categoryHeader().length;
		page.widthColumn((90 / colspan) + "%");
		page.createView();
    });
}

page.createView = function(){
	$("#tableContent").html("");
	$("#tableHeader").html("");
	$("#tableHeader").append('<td width="10%" class="border-right" colspan="2">&nbsp;</td>');

	$.each(page.categoryHeader(), function(id, ress){
		var tdHeader = ' <td width="'+page.widthColumn()+'" class="border-right"><strong>'+ress+'</strong></td>';
		$("#tableHeader").append(tdHeader);
	});

	$.each(page.data(), function(key, value){
		if (value != null){
			var progressData = "";
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

$(function () {
	$('#btnRefresh').on('click', function () {
		app.loading(true);

        $.when(page.getData()).done(function(){
        	setTimeout(function(){
        		app.loading(false);
        		app.prepareTooltipster();
        	},1000);	
        });
    });

	page.hideFilter();
	setTimeout(function() {
		fa.LoadData();
		page.getData();
	},200);
	
    setTimeout(function() {
		app.prepareTooltipster();
        app.loading(false);    
    }, 1000);
});
