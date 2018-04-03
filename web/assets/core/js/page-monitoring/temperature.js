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
mt.logEntries = ko.observableArray([]);
mt.getMode = ko.observable('heatmap');

var $tableInterval = false, $heatmapInterval = false, $intervalTime = 5000;

var requests = [];
var $temperatureInterval = false, $intervalTime = 5000;

ko.subscribable.fn.subscribeChanged = function (callback) {
    var oldValue;
    this.subscribe(function (_oldValue) {
        var value = ko.utils.unwrapObservable(this);
        if (value != null && value.constructor == Array){
            oldValue = _oldValue.slice();
        } else {
            oldValue = _oldValue;
        }
    }, this, 'beforeChange');

    this.subscribe(function (newValue) {
        callback(newValue, oldValue);
    });
};

mt.remove = function(str){ 
    return str.replace(/[\. ,:-]+/g, "");
}  

mt.SelectMode = function(type) {
    mt.abortAll(requests);

    if(type == 'table') {
        clearInterval($heatmapInterval);
        $heatmapInterval = false;
        mt.getMode("table");
        $tableInterval = setInterval(function() { mt.GetData("table"); }, $intervalTime);
    } else {
        clearInterval($tableInterval);
        $tableInterval = false;
        mt.getMode("heatmap");
        $heatmapInterval = setInterval(function() { mt.GetData("heatmap"); }, $intervalTime);
    }


}

mt.Details.subscribeChanged(function(newValue,oldValue){
    if(oldValue.length >0){
        $.each(newValue, function(idx, value){
            $.each(value.turbines, function(i, val){
                for(var key in val){
                   if(key !== "Turbine"){
                        var detailsNewVal = oldValue[idx].turbines;
                        detailsNewVal = detailsNewVal[i][key];
                        if(kendo.toString(detailsNewVal,'n2') != kendo.toString(val[key],'n2')){  
                            val[key+'_Change'] = 1; 
                        }else{
                            val[key+'_Change'] = 0;
                        }
                   }  
                }
            });

            if(mt.getMode() == 'table'){
                window.setTimeout(function(){ 
                    $('table').find($('.blinkYellow')).css('background-color', 'transparent'); 
                    $('table').find($('.blinkYellow')).css('transition' , 'background-color 0.5s ease;');

                }, 1200);
            }



        });
    }
});

mt.GetData = function(type) {
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
         mt.GetDataProject(project, type);
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
    setTimeout(function(){
        app.loading(true);
        mt.Columns([]);
        mt.Details([]);
        mt.logEntries([]);

         var length = requests.length;
         while(length--) {
             requests[length].abort && requests[length].abort();  // the if is for the first case mostly, where array is still empty, so no abort method exists.
         }
    },200)

}

mt.GetDataProject = function(project,type) {
    var param = {
        Project: project,
    };

    var url = (type == "table" ? "monitoringrealtime/getdatatemperature" : "monitoringrealtime/gettemperatureheatmap");

   requests.push(toolkit.ajaxPost(viewModel.appName + url, param, function (res) {
        if(!app.isFine(res)) {
            app.loading(false);
            return;
        }
        var details  = res.data[0].Details;
        var columnList = res.data[1].ColumnList;
        columnList.unshift({title : "Turbine", desc : ""});



        var width = ($(".feeder-column > table").innerWidth() - (columnList.length * 11 + 75)) / (columnList.length - 1);

        $.each(columnList, function(i, val){
            if(val.title == "Turbine"){
                val.Width = "70px";
            }else{
                val.Width = width +"px";
            }
        });

        mt.Details(res.data[0].Details);
        mt.Columns(res.data[1].ColumnList);
        mt.Columns(columnList);

        app.loading(false);
        
   }));
}


$(function() {
    app.loading(true);
    // mt.GetData();

    if(mt.getMode() == null){
        mt.SelectMode('heatmap');
    }else{
        $('.nav-pills a[href="#'+mt.getMode()+'"]').trigger("click");
    }

    $('#projectList').kendoDropDownList({
        data: mt.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            setTimeout(function(){
                mt.abortAll(requests);
                mt.GetData(mt.getMode());
            },1500);
            
         }
    });

    
});