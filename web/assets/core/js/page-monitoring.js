'use strict';

viewModel.monitoring = {};
var monitoring = viewModel.monitoring;


vm.currentMenu('Monitoring');
vm.currentTitle('Monitoring');
vm.breadcrumb([{
    title: 'Monitoring',
    href: viewModel.appName + 'page/monitoring'
}, {
    title: 'Monitoring',
    href: '#'
}]);

monitoring.turbineList = ko.observableArray([]);
monitoring.projectList = ko.observableArray([]);
monitoring.turbine = ko.observableArray([]);
monitoring.project = ko.observable();
monitoring.data = ko.observableArray([]);
monitoring.event = ko.observableArray([]);
monitoring.last_minute = ko.observable();
monitoring.last_date = ko.observable();
var turbineval = [];


monitoring.populateTurbine = function (data) {
    if (data.length == 0) {
        data = [];
        monitoring.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            var allturbine = {}
            $.each(data, function (key, val) {
                turbineval.push(val);
            });
            allturbine.value = "All Turbine";
            allturbine.text = "All Turbines";
            datavalue.push(allturbine);
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        monitoring.turbineList(datavalue);
    }

    setTimeout(function () {
        $('#turbineList').data('kendoMultiSelect').value(["All Turbine"]);
    }, 300);
};

monitoring.populateProject = function (data) {
    if (data.length == 0) {
        data = [];;
        monitoring.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        if (data.length > 0) {
            $.each(data, function (key, val) {
                var data = {};
                data.value = val;
                data.text = val;
                datavalue.push(data);
            });
        }
        monitoring.projectList(datavalue);

        // override to set the value
        setTimeout(function () {
            $("#projectList").data("kendoDropDownList").select(1);
            monitoring.project = $("#projectList").data("kendoDropDownList").value();
        }, 300);
    }
};

monitoring.checkTurbine = function () {
    var arr = $('#turbineList').data('kendoMultiSelect').value();
    var index = arr.indexOf("All Turbine");
    if (index == 0 && arr.length > 1) {
        arr.splice(index, 1);
        $('#turbineList').data('kendoMultiSelect').value(arr)
    } else if (index > 0 && arr.length > 1) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    } else if (arr.length == 0) {
        $("#turbineList").data("kendoMultiSelect").value(["All Turbine"]);
    }
}

monitoring.getData = function(){
    // app.loading(true);

    var turbine = $("#turbineList").data("kendoMultiSelect").value()
    var project = $("#projectList").data("kendoDropDownList").value()
    var param = {
        turbine: (turbine == "All Turbine" ? []: turbine),
        project: project
    };

    var request = toolkit.ajaxPost(viewModel.appName + "monitoring/getdata", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        monitoring.last_minute(res.data.Data.timestamp.minute);
        monitoring.last_date(res.data.Data.timestamp.date);

        monitoring.data([]);
        $.each(res.data.Data.data, function (index, item) {   
            monitoring.data.push(item);                    
        });
    });

    var requestEvent = toolkit.ajaxPost(viewModel.appName + "monitoring/getevent", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
       monitoring.event(res.data.Data)
    });


    $.when(request, requestEvent).done(function(){
        $(".red").addClass("flash");
        setTimeout(function(){
            // app.loading(false);
            app.prepareTooltipster();            
        },500);

        setTimeout(function() {
            $(".red").removeClass("flash");
        }, 2500);
    });
}

$(function () {

    $("#restore-screen").hide();

    $("#max-screen").click(function(){
        $("html").addClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 80);
        $(".multicol").height($(window).innerHeight() - 80 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 80-50);
        $("#max-screen").hide();
        $("#restore-screen").show();  
    });

    $("#restore-screen").click(function(){
        $("html").removeClass("maximize-mode");
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150-50);
        $("#max-screen").show();  
        $("#restore-screen").hide();  
    });

    $('#btnRefresh').on('click', function() {
        monitoring.getData();
    });

    setInterval(function(){monitoring.getData()},1000*120);

    setTimeout(function() {
        $(".multicol-div").height($(window).innerHeight() - 150);
        $(".multicol").height($(window).innerHeight() - 150 - 25);
        $(".control-sidebar").height($(window).innerHeight() - 150 - 50);
        monitoring.getData();
    }, 500);
});