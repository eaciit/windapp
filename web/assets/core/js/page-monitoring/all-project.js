'use strict';


viewModel.AllProject = new Object();
var page = viewModel.AllProject;


vm.currentMenu('By Project');
vm.currentTitle('By Project');
vm.isShowDataAvailability(false);
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Summary By Project', href: viewModel.appName + 'page/monitoringallproject' }]);



page.DataDetails = ko.observableArray([]);
page.TimeMax = ko.observable("");


page.getData = function(){
    var param = {Project: "", LocationTemp: 0};

    toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdataproject", param, function (res) {
        if(!app.isFine(res)) {
            app.loading(false);
            return;
        }
        
        page.DataDetails(res.data.Detail);
        page.generateView(res.data.Detail);
        page.TimeMax(moment.utc(res.data.TimeMax).format("DD MMM YYYY HH:mm:ss"));

        app.loading(false);
    });
}

page.generateView = function(datasource){
    
    $.each(datasource, function(idx, val){

        var comparison = 0;
        var defaultColorStatus = "bg-default-green"
        var colorStatus = "lbl bg-green"
        $('#statusprojectdefault_'+ val.Project).addClass(defaultColorStatus);
        
        
        if((val.PowerGeneration / val.Capacity) > 0){
            comparison = (val.ActivePower / val.Capacity) * 70;
            $('#statusproject_'+ val.Project).attr('class', colorStatus);
            $('#statusproject_'+ val.Project).css('width', comparison + 'px');
        }else{
            comparison = 0;
            $('#statusproject_'+ val.Project).attr('class', 'lbl');
        }
    });
}

$(function() {
    app.loading(true);
    page.getData();
    setInterval(page.getData, 5000);
});