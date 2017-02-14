
pm.dtMt = ko.observable();

var mt = {
	RefreshData: function() {
	    app.loading(true);
	    pm.showFilter();
	    fa.LoadData();
	    if(pm.isFirstMTBF() === true){
	        mt.Refreshchartmt();
	        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.startScadaOEM + '</strong> until: ');
			$('#availabledateend').html('<strong>' + availDateList.endScadaOEM + '</strong>');

	    }else{
	        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.startScadaOEM + '</strong> until: ');
			$('#availabledateend').html('<strong>' + availDateList.endScadaOEM + '</strong>');
	        setTimeout(function(){
	            $("#chartTI").data("kendoChart").refresh();
	            app.loading(false);
	        }, 300);
	    }

	},
	
	Refreshchartmt: function() {
		app.loading(true);
		var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

		toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/GetListMtbf", param, function (data) {
			pm.isFirstMTBF(false);

			pm.dtMt(data)
			var width = $(".main-header").width()
			var Height = width*0.2

			$('#gridmtbf').html('');

			$("#gridmtbf").kendoGrid({
				theme: "flat",
	            dataSource: {
	                data: pm.dtMt().data,
	                pageSize: 10,
	                sort: [{
	                    field: 'id',
	                    dir: 'asc'
	                }],
	            },
	            pageable: {
		            pageSize: 10,
		            input: true, 
		        },
	            scrollable: true,
	            sortable: true,
	            columns: [{
	                field: "id",
	                title: "Turbine",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                attributes: { class: "align-center"},
	                width: 120,
	            },{
	                field: "avgmtbf",
	                title: "Avg. MTBF",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                attributes: { class: "align-center"},
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "avgmttf",
	                title: "Avg. MTTF",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                attributes: { class: "align-center"},
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "avgmttr",
	                title: "Avg. MTTR",
	                attributes: { class: "align-center"},
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            } ],
			});
			setTimeout(function() {
				app.loading(false);
 				$("#gridmtbf").data("kendoGrid").refresh();
			}, 100);



		});
	},
};

$(document).ready(function() {
	// ti.LoadData();
});