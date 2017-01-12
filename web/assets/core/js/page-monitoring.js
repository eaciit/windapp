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

vm.dateAsOf(app.currentDateData);
monitoring.createGauge = function(id){
	$("#gauge"+id).html("");
	$("#gauge"+id).kendoLinearGauge({
		theme: "flat",
        pointer: {
            value: 65,
            shape: "arrow"
        },
        gaugeArea: {
        	margin: {
        		bottom: -40
        	}
        },
        scale: {
            majorUnit: 20,
            minorUnit: 5,
            max: 180,
            vertical: false,
            labels: {
            	visible: false,
            	padding: 0,
            },
            ranges: [
                {
                    from: 80,
                    to: 120,
                    color: "#ffc700"
                }, {
                    from: 120,
                    to: 150,
                    color: "#ff7a00"
                }, {
                    from: 150,
                    to: 180,
                    color: "#c20000"
                }
            ]
        }
    });
}


monitoring.getData = function(){
    toolkit.ajaxPost(viewModel.appName + "monitoring/getdata", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        console.log(res.data.Data);
    });
}

$(function () {
	for(var i = 0 ; i < 5 ; i++){
		monitoring.createGauge(i);
	}

    $("#restore-screen").hide();

    $("#max-screen").click(function(){
        $("html").addClass("maximize-mode");
        $("#max-screen").hide();
        $("#restore-screen").show();  
    });

    $("#restore-screen").click(function(){
        $("html").removeClass("maximize-mode");
        $("#max-screen").show();  
        $("#restore-screen").hide();  
    });

    setTimeout(function() {
        app.loading(false);
    }, 500);
});