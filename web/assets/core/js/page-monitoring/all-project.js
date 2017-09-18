'use strict';


viewModel.AllProject = new Object();
var page = viewModel.AllProject;


// vm.currentMenu('Overall');
// vm.currentTitle('Overall');
// vm.isShowDataAvailability(false);
// vm.breadcrumb([
//     { title: "Monitoring", href: '#' }, 
//     { title: 'Overall', href: viewModel.appName + 'page/monitoringallproject' }]);



page.DataDetails = ko.observableArray([]);
page.OldData = ko.observableArray([]);
page.TimeMax = ko.observable("");
page.projectList = ko.observableArray(projectList);
var requests = [];

// get the weather forecast
page.getDirection = function() {
    var surl = 'http://api.openweathermap.org/data/2.5/weather';
    var requests = [];
    if(projectList.length > 0) {
        $.each(projectList, function(idx, p){
            var param = { "q": p.City, "appid": "88f806b961b1057c0df02b5e7df8ae2b", "units": "metric" };    
            var $elm = $("#detail-"+ p.ProjectId);
            requests.push($.ajax({
                type: "GET",
                url: surl,
                data: param,
                dataType: "jsonp",
                success:function(data){
                    var winddeg = parseFloat(data.wind.deg);
                    $elm.find('.fa-location-arrow').rotate({
                        angle: 0,
                        animateTo: winddeg - 45,
                    });
                },
                error:function(){
                    // do nothing
                }  
            }));    
        });
    }
};


page.getData = function(){
    var param = {Project: "", LocationTemp: 0};

    requests.push(toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
        if(!app.isFine(res)) {
            app.loading(false);
            return;
        }
        
        if(page.OldData().length == 0){
             page.DataDetails(res.data.Detail);
             page.OldData(res.data.Detail);
        }else{
             page.DataDetails(res.data.Detail);
        }
       
        page.generateView();
        page.TimeMax(moment.utc(res.data.TimeMax).format("DD MMM YYYY HH:mm:ss"));

        page.getDirection();

        app.loading(false);
    }));
}

page.generateView = function(){
    var datasource = page.DataDetails();
    var setView = $.each(datasource, function(idx, val){
         var oldData = page.OldData()[idx];

        for(var key in val){
            var id = '#'+key+'_'+ val.Project
            if(oldData.hasOwnProperty(key) && (oldData.Project == val.Project)){
                if(val[key] != oldData[key] ) {
                    $(id).animate( { backgroundColor: 'rgba(255, 216, 0, 0.7)' }, 500).animate( { backgroundColor: 'transparent' }, 500); 
                }

            }
        }       

        var comparison = 0;
        var defaultColorStatus = val.DefaultColorStatus
        var colorStatus = val.ColorStatus

        $('#statusprojectdefault_'+ val.Project).addClass(defaultColorStatus);

        if(((val.PowerGeneration/1000) / val.Capacity) > 0){
            comparison = ((val.PowerGeneration/1000) / val.Capacity) * 100;
            var fixCom = (comparison > 100 ? 100 : comparison);
            // console.log(comparison);
            $('#statusproject_'+ val.Project).attr('class', colorStatus);
            $('#statusproject_'+ val.Project).css('width', fixCom + '%');
        }else{
            comparison = 0;
            $('#statusproject_'+ val.Project).attr('class', 'lbl');
        }
    });

    $.when(setView).done(function(){
        page.OldData(datasource);
    });
}

page.ToByProject = function(project) {
    setTimeout(function(){
        app.loading(true);
        app.resetLocalStorage();
        localStorage.setItem('projectname', project);
        if(localStorage.getItem('projectname') !== null){
            window.location = viewModel.appName + "page/monitoringbyproject";
        }
    },1500);
}