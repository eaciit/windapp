
pg.dtMt = ko.observable();

var mt = {
	RefreshData: function() {
	    app.loading(true);
	    pg.showFilter();
	    fa.LoadData();
	    if(pg.isFirstMTBF() === true){
	        mt.Refreshchartmt();
	        $('#availabledatestart').html('Data Available from: <strong>' + availDateListLoss.startScadaOEM + '</strong> until: ');
			$('#availabledateend').html('<strong>' + availDateListLoss.endScadaOEM + '</strong>');

	    }else{
	        $('#availabledatestart').html('Data Available from: <strong>' + availDateListLoss.startScadaOEM + '</strong> until: ');
			$('#availabledateend').html('<strong>' + availDateListLoss.endScadaOEM + '</strong>');
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
			pg.isFirstMTBF(false);

			pg.dtMt(data)
			var width = $(".main-header").width()
			var Height = width*0.2

			$('#gridmtbf').html('');

			$("#gridmtbf").kendoGrid({
				theme: "flat",
	            dataSource: {
	                data: pg.dtMt().data,
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
	                field: "mtbf",
	                title: "MTBF (Hrs)",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                attributes: { class: "align-center"},
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "mttf",
	                title: "MTTF (Hrs)",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                attributes: { class: "align-center"},
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "mttr",
	                title: "MTTR (Hrs)",
	                attributes: { class: "align-center"},
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "totoptime",
	                title: "Total Operation Time (Hrs)",
	                attributes: { class: "align-center"},
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "totnooffailure",
	                title: "Total Number Of Failures",
	                attributes: { class: "align-center"},
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "totdowntime",
	                title: "Total Downtime (Hrs)",
	                attributes: { class: "align-center"},
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            }  ],
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