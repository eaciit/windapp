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


var audioElement = document.createElement('audio');
    audioElement.setAttribute('src', "../res/alarm/alarm.mp3");


ko.bindingHandlers.singleOrDoubleClick = {
    init: function(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
        var singleHandler   = valueAccessor().click,
            doubleHandler   = valueAccessor().dblclick,
            delay           = valueAccessor().delay || 1000,
            clicks          = 0;

        $(element).bind('click',function(event) {
            clicks++;
            if (clicks === 1) {
                setTimeout(function() {
                    if( clicks === 1 ) {
                        // Call the single click handler - passing viewModel as this 'this' object
                        // you may want to pass 'this' explicitly
                        if (singleHandler !== undefined) { 
                            singleHandler.call(viewModel, bindingContext.$data, event); 
                        }
                    } else {
                        // Call the double click handler - passing viewModel as this 'this' object
                        // you may want to pass 'this' explicitly
                        if (doubleHandler !== undefined) { 
                            doubleHandler.call(viewModel, bindingContext.$data, event); 
                        }
                    }
                    clicks = 0;
                }, delay);
            }
            return false;
        });
    }
};

bpc.getShorterName = function(str){

    var strSplit =  str.split('(', 2);
    if(strSplit[0].length > 10) {
        str = str.substring(0,7) + ".. ("+strSplit[1] ;
    }else{
        str = str;
    }

    return str;
}

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
	audioElement.currentTime = 0;
    audioElement.pause();

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


	$elmDetail.find('#timemax_'+ project).text(moment.utc($data.TimeMax).format('DD MMM YYYY HH:mm:ss'));

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

			var $oldStatus = $elmupdate.attr('class').split(' ')[1];

			if($oldStatus !== undefined){
				if($oldStatus == "bg-default-green" && defaultColorStatus == "bg-default-red"){
                    var playPromise = audioElement.play();
                    if (playPromise !== null){
                        playPromise.catch(() => { audioElement.play(); })
                    }
                }else{
                    audioElement.pause();
                }
			}

			if(defaultColorStatus != "bg-default-green"){
				$elmupdate.prop('style', 'width: 100%');
			}else{
				$elmupdate.prop('style', 'width: '+ currPct.toFixed(0) + '%');
			}


			if(dt.IsReapeatedAlarm == true){
                $elmdetail.addClass("reapeat");
            }else{
                $elmdetail.removeClass("reapeat");
            }

			if(dt.IsRemark == true){
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="icon_'+ dt.Turbine +'"]').css("display", "block !important");
				// $("#cusmon-turbine-"+project).find(".turbine-detail").find('.inner-triangle[data-id="'+ dt.Turbine +'"]').show();
			}else{
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="icon_'+ dt.Turbine +'"]').hide();
				// $("#cusmon-turbine-"+project).find(".turbine-detail").find('.inner-triangle[data-id="'+ dt.Turbine +'"]').hide();
			}

			if(dt.BulletColor == "fa fa-circle txt-green" || dt.BulletColor == "fa fa-circle txt-grey" ){
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-temp[data-id="icontemp_'+ dt.Turbine +'"]').hide();
			}else{
				$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-temp[data-id="icontemp_'+ dt.Turbine +'"]').addClass(dt.BulletColor).attr('data-original-title', dt.TemperatureInfo);
			}

			

			$("#cusmon-turbine-"+project).find(".turbine-detail").find('.total-production[data-id="total_'+ dt.Turbine +'"]').attr("title","Gen. Today (Mwh) : "+ kendo.toString(dt.TotalProduction,'n1'));

			$elmupdate.prop('aria-valuenow', currPct);
			$elmupdate.attr("class" , "progress-bar " +defaultColorStatus);
			
		});
	}

	var $feederList = $data.FeederRemarkList;

	for(var key in $feederList){
		if($feederList[key] == true){
			$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="icon_'+ key +'"]').css("display", "block !important");
		}else{
			$("#cusmon-turbine-"+project).find(".turbine-detail").find('.icon-remark[data-id="icon_'+ key +'"]').hide();
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
			var classIcon = 'txt-green';
			var classString = $("div").find("[data-id='"+dt.Id+"']").children().attr('class').split(' ')[1];
			if(classString !== undefined){
				classIcon = 'txt'+ classString.substr(10);
			}
			
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
	    setTimeout(function(){
	        app.loading(true);
	        app.resetLocalStorage();
	        localStorage.setItem('turbine', data.Id);
	        localStorage.setItem('projectname', data.Project);
	        localStorage.setItem('isFromSummary', true);
	        localStorage.setItem('isFromByProject', false);

	        if(localStorage.getItem("turbine") !== null && localStorage.getItem("projectname")){
	        	window.location = viewModel.appName + "page/monitoringbyturbine";
	        }
	        
	    },1500);
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