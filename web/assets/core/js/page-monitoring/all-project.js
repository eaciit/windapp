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


page.getData = function(){
    var param = {Project: "", LocationTemp: 0};

    toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
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

        app.loading(false);
    });
}

page.generateView = function(){
    var datasource = page.DataDetails();
   

    var setView = $.each(datasource, function(idx, val){
         var oldData = page.OldData()[idx];

        for(var key in val){
            var id = '#'+key+'_'+ val.Project
            if(oldData.hasOwnProperty(key) && (oldData.Project == val.Project)){
                if(val[key] != oldData[key] ) {
                    // $(id).css('background-color', 'rgba(255, 216, 0, 0.7)'); 
                    $(id).animate( { backgroundColor: 'rgba(255, 216, 0, 0.7)' }, 500).animate( { backgroundColor: 'transparent' }, 500); 
                }

            }
        }       

        var comparison = 0;
        var defaultColorStatus = "bg-default-green"
        var colorStatus = "lbl bg-green"
        $('#statusprojectdefault_'+ val.Project).addClass(defaultColorStatus);
        
        
        // if((val.PowerGeneration / val.Capacity) > 0){
        //     comparison = (val.ActivePower / val.Capacity) * 70;
        //     $('#statusproject_'+ val.Project).attr('class', colorStatus);
        //     $('#statusproject_'+ val.Project).css('width', comparison + 'px');
        // }else{
        //     comparison = 0;
        //     $('#statusproject_'+ val.Project).attr('class', 'lbl');
        // }

        if(((val.PowerGeneration/1000) / val.Capacity) > 0){
            comparison = ((val.PowerGeneration/1000) / val.Capacity) * 100;
            // console.log(comparison);
            $('#statusproject_'+ val.Project).attr('class', colorStatus);
            $('#statusproject_'+ val.Project).css('width', comparison + '%');
        }else{
            comparison = 0;
            $('#statusproject_'+ val.Project).attr('class', 'lbl');
        }
    });

    $.when(setView).done(function(){
        page.OldData(datasource);
    });
}