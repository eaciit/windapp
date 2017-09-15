'use strict';

viewModel.ByProjectCustom = new Object();
var bpc = viewModel.ByProjectCustom;

// vm.currentMenu('All Farms');
// vm.currentTitle('All Farms');
// vm.breadcrumb([
//     { title: "Monitoring", href: '#' }, 
//     { title: 'All Farms', href: viewModel.appName + 'page/monitoringbyprojectcustom' },
// ]);

bpc.projectList = ko.observableArray(projectList);
bpc.feederList = ko.observableArray([]);
bpc.turbineList = ko.observableArray([]);

var requests = [];

ko.bindingHandlers.singleOrDoubleClick = {
    init: function(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
        var singleHandler   = valueAccessor().click,
            doubleHandler   = valueAccessor().dblclick,
            delay           = valueAccessor().delay || 1000,
            clicks          = 0;

        $(element).bind('click',function(event) {
            clicks++;
            console.log(clicks);
            if (clicks === 1) {
                setTimeout(function() {
                    if( clicks === 1 ) {
                        // Call the single click handler - passing viewModel as this 'this' object
                        // you may want to pass 'this' explicitly
                        if (singleHandler !== undefined) { 
                            singleHandler.call(viewModel, bindingContext.$data, event); 
                        }
                    } else {
                    	console.log(doubleHandler);
                        // Call the double click handler - passing viewModel as this 'this' object
                        // you may want to pass 'this' explicitly
                        if (doubleHandler !== undefined) { 
                            doubleHandler.call(viewModel, bindingContext.$data, event); 
                        }
                    }
                    clicks = 0;
                }, delay);
            }
        });
    }
};



// get the weather forecast
bpc.getWeather = function() {
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    var requests = [];
    if(projectList.length > 0) {
    	$.each(projectList, function(idx, p){
    		var param = { "q": p.City, "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };	
    		var $elm = $("#cusmon-project-"+ p.ProjectId);
    		requests.push($.ajax({
		        type: "GET",
		        url: surl,
		        data: param,
		        dataType: "jsonp",
		        success:function(data){
		            $elm.find('.img-temp-forecast').attr('src', 'http://openweathermap.org/img/w/'+ data.weather[0].icon +'.png');
		            $elm.find('.temp-forecast').html(parseFloat(data.main.temp).toFixed(1) + '&deg;C');
		            $elm.find('.ws-forecast').html(data.wind.speed + ' <small>m/s</small>');
		            
		            var winddeg = parseFloat(data.wind.deg);
		            $elm.find('.wd-forecast').html(winddeg.toFixed(1) + '&deg;');
		            $elm.find('.img-wd-forecast').rotate({
				    	angle: 0,
				    	animateTo: winddeg,
				    });
		        },
		        error:function(){
		            // do nothing
		        }  
		    }));	
    	});
    	$.when.apply(undefined, requests).then(function(){
    		// do nothing
    	});
    }
};
	
// get the data for all projects
bpc.getData = function() {
	if(projectList.length > 0) {
		$.each(projectList, function(idx, p){
			// param to get the data
			var param = {
		        Project: p.ProjectId,
		        LocationTemp: 30.0
		    };

		    // add to queue getting data for all projects
		    requests.push(
		    	toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdatafarm", param, function (res) {
					bpc.plotData(param.Project, res);
			    })
		    );
		});

		// applying all requests then prepare to plotting the data
		$.when.apply(undefined, requests).then(function(){
			// do nothing
		});
	}
};

// ploting the data
var defaultValue = -999999;
bpc.plotData = function(project, data) {
	var $data = data.data;
	var $elm = $("#cusmon-project-"+ project);
	var $elmDetail= $("#cusmon-detail-"+ project);
	var $elmTurbine= $("#cusmon-turbine-"+ project);

	var totalpower = $data.PowerGeneration;
	var avgwindspd = $data.AvgWindSpeed;
	
	// set project updates
	$elm.find('.power[data-id="'+ project +'"]').text((parseFloat(totalpower)/1000).toFixed(1));
	$elm.find('.ws[data-id="'+ project +'"]').text(parseFloat(avgwindspd).toFixed(2));	
	$elm.find('.t-up[data-id="'+ project +'"]').text($data.TurbineActive);
	$elm.find('.t-down[data-id="'+ project +'"]').text($data.TurbineDown);
	$elm.find('.t-wait[data-id="'+ project +'"]').text($data.TurbineWaitingWS);
	$elm.find('.t-na[data-id="'+ project +'"]').text($data.TurbineNotAvail);

	if($data.IsRemark == true){
		$elmDetail.find('.project-remark[data-id="'+ project +'"]').css("display", "block !important");
	}else{
		$elmDetail.find('.project-remark[data-id="'+ project +'"]').hide();
	}
	

	// set turbine updates
	var $detail = $data.Detail;
	if($detail.length > 0) {
		$.each($detail, function(idx, dt){
			var $elmdetail = $('.progress[data-id="'+ dt.Turbine +'"]');
			var currPct = 0;
			var power = dt.ActivePower;

			var defaultColorStatus = dt.DefaultColorStatus;


			if(dt.ActivePower > 0 && dt.Capacity > 0) {
				currPct = (dt.ActivePower/dt.Capacity)*100;
			}
			var $elmupdate = $elmdetail.find('.progress-bar[role="progressbar"]');

			if(defaultColorStatus != "bg-default-green"){
				$elmupdate.prop('style', 'width: 100%');
			}else{
				$elmupdate.prop('style', 'width: '+ currPct.toFixed(0) + '%');
			}


			if(dt.IsRemark == true){
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="'+ dt.Turbine +'"]').css("display", "block !important");
			}else{
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="'+ dt.Turbine +'"]').hide();
			}

			$elmupdate.prop('aria-valuenow', currPct);
			$elmupdate.attr("class" , "progress-bar " +defaultColorStatus);
			
		});
	}

	var $feederList = $data.FeederRemarkList;

	for(var key in $feederList){
		if($feederList[key] == true){
			$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="'+ key +'"]').css("display", "block !important");
		}else{
			$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="'+ key +'"]').hide();
		}
	}


}

// getting data every interval time
bpc.refresh = function() {
	return setInterval(bpc.getData, $intervalTime);
}


// turbine collaboration open
bpc.OpenTurbineCollaboration = function(dt) {
	return function(dt) {
		TbCol.ResetData();
		if(dt.IsTurbine) {
			var classString = $("div").find("[data-id='"+dt.Id+"']").children().attr('class').split(' ')[1];
			var classIcon = 'txt'+ classString.substr(10);
			
			TbCol.TurbineId(dt.Id);
			TbCol.TurbineName(dt.Name);
			TbCol.UserId('');
			TbCol.UserName('');
			TbCol.Project(dt.Project);
			TbCol.Feeder(dt.Feeder);
			TbCol.Status(dt.Status);
			TbCol.IsTurbine(true);
			TbCol.OpenForm();
			TbCol.IconStatus(classIcon);
		}else{
			TbCol.ProjectFeeder(dt.Project)
			TbCol.Feeder(dt.Name);
			TbCol.IsTurbine(false);
			TbCol.OpenForm();
		}
	}
};

bpc.OpenModal = function(data){
	TbCol.ResetData();
	if(data.isProject){
		TbCol.Project(data.Project);
		TbCol.IsTurbine(false);
		TbCol.OpenForm();
	}
}

bpc.ToIndividualTurbine = function(data) {
	return function(data) {
		if(data.IsTurbine) {
		    app.loading(true);
		    var oldDateObj = new Date();
		    var newDateObj = moment(oldDateObj).add(3, 'm');
		    document.cookie = "projectname="+data.Project.split("(")[0].trim()+";expires="+ newDateObj;
		    document.cookie = "turbine="+data.Id+";expires="+ newDateObj;
		    document.cookie = "isFromSummary=true;expires="+ newDateObj;
		    document.cookie = "isFromByProject=false;expires="+ newDateObj;
		    if(document.cookie.indexOf("projectname=") >= 0 && document.cookie.indexOf("turbine=") >= 0 && document.cookie.indexOf("isFromSummary=") >= 0 && document.cookie.indexOf("isFromByProject=") >= 0) {
		        window.location = viewModel.appName + "page/monitoringbyturbine";
		    } else {
		        app.loading(false);
		    }
		}
	}
}

// init page
$(function() {
	$('#savedViews').kendoDropDownList({
		data: [],
		dataValueField: 'value',
		dataTextField: 'text',
		change: function () {  },
	});

	$( "#sortable" ).sortable();
    $( "#sortable" ).disableSelection();

    // refresh the data every second
    // bpc.refresh();
});