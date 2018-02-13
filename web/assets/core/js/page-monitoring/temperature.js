'use strict';


viewModel.Temperature = new Object();
var mt = viewModel.Temperature;


vm.currentMenu('Temperature');
vm.currentTitle('Temperature');
vm.isShowDataAvailability(false);
vm.breadcrumb([
    { title: "Monitoring", href: '#' }, 
    { title: 'Temperature', href: viewModel.appName + 'page/monitoringtemperature' }]);

mt.projectList = ko.observableArray([]);
mt.project = ko.observable();
mt.Columns = ko.observableArray([]);
mt.Details = ko.observableArray([]);

var requests = [];
var $temperatureInterval = false, $intervalTime = 5000;


mt.GetData = function(data) {
    // app.loading(true);
    var project = "";

    if(localStorage.getItem("projectname") !== null) {
        project = localStorage.getItem("projectname");
        $('#projectList').data('kendoDropDownList').value(project);
        app.resetLocalStorage();
    } else {
        project = $('#projectList').data('kendoDropDownList').value();
    }
    
    setTimeout(function(){
         mt.GetDataProject(project);
    }, 200);

    // count++;
};

mt.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        mt.projectList([{ value: "", text: "" , city: ""}]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val.Value;
                data.text = val.Name;
                data.city = val.City;
                datavalue.push(data);
            });
        }
        mt.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(0);
            mt.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

mt.abortAll = function(requests) {
     var length = requests.length;
     while(length--) {
         requests[length].abort && requests[length].abort();  // the if is for the first case mostly, where array is still empty, so no abort method exists.
     }
}

mt.GetDataProject = function(project) {
    var param = {
        Project: project,
    };
   requests.push(toolkit.ajaxPost(viewModel.appName + "monitoringrealtime/getdatatemperature", param, function (res) {
        if(!app.isFine(res)) {
            app.loading(false);
            return;
        }
        var details  = res.data[0].Details;

        mt.Details(res.data[0].Details);
        mt.Columns(res.data[1].ColumnList);
        mt.Columns.unshift({title : "Turbine", desc : ""});




        app.loading(false);
        
   }));
}



$(function() {
    app.loading(true);
    // mt.GetData()

    $('#projectList').kendoDropDownList({
        data: mt.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            setTimeout(function(){
                mt.abortAll(requests);
                mt.GetData();
            },1500);
            
         }
    });

    setInterval(mt.GetData, 5000);
});