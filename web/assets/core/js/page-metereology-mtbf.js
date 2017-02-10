
pm.dtMt = ko.observable();

var mt = {
	
	Refreshchartmt: function() {
		app.loading(true);
		var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

		app.ajaxPost(viewModel.appName + "/databrowsernew/getscadadataoemavaildate", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }

            //Scada Data
            if (res.data.ScadaDataOEM.length == 0) {
                res.data.ScadaDataOEM = [];
            } else {
                if (res.data.ScadaDataOEM.length > 0) {
                    var minDatetemp = new Date(res.data.ScadaDataOEM[0]);
                    var maxDatetemp = new Date(res.data.ScadaDataOEM[1]);
					$('#availabledatestart').html('Data Available from: <strong>' + kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')) + '</strong> until: ');
					$('#availabledateend').html('<strong>' + kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')) + '</strong>');
                }
            }         
        });

		toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/GetListMtbf", param, function (data) {

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
	                width: 120,
	            },{
	                field: "avgmtbf",
	                title: "Avg. MTBF",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "avgmttf",
	                title: "Avg. MTTF",
	                headerAttributes: {
	                    style: "text-align: center"
	                },
	                format: "{0:n2}",
	                width: 120,
	            },{
	                field: "avgmttr",
	                title: "Avg. MTTR",
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