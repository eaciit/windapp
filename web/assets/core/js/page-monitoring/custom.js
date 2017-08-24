'use strict';

viewModel.ByProjectCustom = new Object();
var bpc = viewModel.ByProjectCustom;

vm.currentMenu('By Project');
vm.currentTitle('By Project');
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'By Project', href: viewModel.appName + 'page/monitoringbyproject' },
]);

bpc.projectList = ko.observableArray(projectList);
bpc.feederList = ko.observableArray([]);
bpc.turbineList = ko.observableArray([]);

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
	var requests = [];
	if(projectList.length > 0) {
		$.each(projectList, function(idx, p){
			// param to get the data
			var param = {
		        Project: p.ProjectId,
		        LocationTemp: 30.0
		    };

		    // add to queue getting data for all projects
		    requests.push(
		    	toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
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

	var totalpower = $data.PowerGeneration;
	var avgwindspd = $data.AvgWindSpeed;
	
	// set project updates
	$elm.find('.power[data-id="'+ project +'"]').text((parseFloat(totalpower)/1000).toFixed(1));
	$elm.find('.ws[data-id="'+ project +'"]').text(parseFloat(avgwindspd).toFixed(2));	
	$elm.find('.t-up[data-id="'+ project +'"]').text($data.TurbineActive);
	$elm.find('.t-down[data-id="'+ project +'"]').text($data.TurbineDown);
	$elm.find('.t-wait[data-id="'+ project +'"]').text($data.TurbineWaitingWS);
	$elm.find('.t-na[data-id="'+ project +'"]').text($data.TurbineNotAvail);

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
			$elmupdate.prop('style', 'width: '+ currPct.toFixed(0) + '%');
			$elmupdate.prop('aria-valuenow', currPct);
			$elmupdate.addClass(defaultColorStatus);
		});
	}
}

// getting data every interval time
bpc.refresh = function() {
	setInterval(bpc.getData, 3000);
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

    // load get weather first then load the data
    $.when(bpc.getWeather()).done(function(){
    	bpc.getData();	
    });

    // refresh the data every second
    bpc.refresh();
});