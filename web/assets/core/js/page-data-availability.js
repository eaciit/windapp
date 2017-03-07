viewModel.DataAvailability = new Object();
var page = viewModel.DataAvailability;

vm.currentMenu('Data Availability');
vm.currentTitle('Data Availability');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'KPI Table', href: viewModel.appName + 'page/dataavailability' }]);

page.isExpanded = ko.observable();
page.categoryHeader = ko.observableArray(["Jan", "Feb", "Mar", "Apr","May","Jun","Jul"]);

var colspan = page.categoryHeader().length;

page.widthColumn = ko.observable((90 / colspan) + "%");

page.dataDummy = ko.observableArray([
		{
			Category : "Scada Data", 
			Turbine	 : [{
					"TurbineName" : "HBR004",
					"details" : [{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "5%"
					},
					{
						"tooltip" : "11 jan 2017 until 1 feb",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "7%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "20%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "15%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "5%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "10%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "9%"
					}]
			},{
					"TurbineName" : "HBR005",
					"details" : [{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "5%"
					},
					{
						"tooltip" : "11 jan 2017 until 1 feb",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "7%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "20%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "15%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "5%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-red",
						"value"	  : "10%"
					},
					{
						"tooltip" : "2 jan 2017 until 10 jan",
						"class"	  : "progress-bar progress-bar-success",
						"value"	  : "9%"
					}]
			}],
			Data 	 : [
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "20%"
				},
				{
					"tooltip" : "11 jan 2017 until 1 feb",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "20%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "15%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "15%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "5%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},

			]
		},
		{
			Category : "Event Down", 
			Turbine	 : [],
			Data 	 : [
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},
				{
					"tooltip" : "11 jan 2017 until 1 feb",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "20%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "2%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "5%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "1%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-red",
					"value"	  : "10%"
				},
				{
					"tooltip" : "2 jan 2017 until 10 jan",
					"class"	  : "progress-bar progress-bar-success",
					"value"	  : "10%"
				},

			]
		}	
	]); 


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

        page.dataDummy(res.data.Data);

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

	$.each(page.dataDummy(), function(key, value){
		var progressData = "";
		$.each(value.Data, function(i, val){
			progressData += '<div aria-hidden="true" class="tooltipster tooltipstered '+val.class+'" style = "width:'+val.value+'"  title = "'+val.tooltip+'" role="progressbar"></div>'
			
		});

		var icon = "";
		if(value.Turbine.length > 1){
			icon = '<i class="fa fa-chevron-right"></i>';
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
	});
	app.loading(false);
}

$(function () {
	$('#btnRefresh').on('click', function () {
		app.loading(true);

        $.when(page.getData()).done(function(){
        	setTimeout(function(){
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
		$('.collapse').on('shown.bs.collapse', function(){
			$(this).parent().find(".fa-chevron-right").removeClass("fa-chevron-right").addClass("fa-chevron-down");
		}).on('hidden.bs.collapse', function(){
			$(this).parent().find(".fa-chevron-down").removeClass("fa-chevron-down").addClass("fa-chevron-right");
		});

        app.loading(false);    
    }, 500);
});
