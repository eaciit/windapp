 'use strict';

viewModel.AnalyticDgrReport = new Object();
var page = viewModel.AnalyticDgrReport;

vm.currentMenu('DGR Report');
vm.currentTitle('DGR Report');
vm.breadcrumb([{ title: "KPI's", href: '#' },{ title: 'DGR Report', href: viewModel.appName + 'page/dgrreport' }]);

page.CurrentData = ko.observable({ Id: '', Model: '', WindSpeed: 0, Power1: 0 });
page.turbineList = ko.observableArray([]);


page.getPdfGrid = function(){
    $("#gridDGRReport").getKendoGrid().saveAsExcel();
     return false;
}

page.fancyTimeFormat = function(time)
{   
	var hours = Math.floor(time / 3600);
	time -= hours * 3600;

	var minutes = Math.floor(time / 60);
	time -= minutes * 60;

	var seconds = parseInt(time % 60, 10);

	return hours + ':' + (minutes < 10 ? '0' + minutes : minutes) + ':' + (seconds < 10 ? '0' + seconds : seconds);
}

var Data = {
	LoadData: function() {
		app.loading(true);
		var isValid = fa.LoadData();
        if(isValid) {
            this.InitGrid();
        }
	},
	InitGrid: function() {

	    var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();   


		var param = {
            period: fa.period,
            dateStart: new Date(moment(dateStart).format('YYYY-MM-DD')),
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: fa.turbine(),
            project: fa.project,
		};


	    var title = param.project+"DGRReport"+kendo.toString(param.dateStart, "dd/MM/yyyy")+"to"+kendo.toString(param.dateEnd, "dd/MM/yyyy")+".xlsx";

		$('#gridDGRReport').html("");
		$('#gridDGRReport').kendoGrid({
	      dataSource: {
	        // serverPaging: true,
	        // serverSorting: true,
	        // serverFiltering: true,
	        transport: {
	          read: {
	            url: viewModel.appName + "analyticdgrreport/getlist",
	            type: "POST",
	            data: param,
	            dataType: "json",
	            contentType: "application/json; charset=utf-8",
	          },
	          parameterMap: function(options) {
	            return JSON.stringify(options);
	          }
	        },
	        pageSize: 10,
	        schema: {
	          data: function(res) {
                    return res.data.Data
                },
                total: function(res) {
	                if (!app.isFine(res)) {
	                    return;
	                }

	                return res.data.Total;
	            }
	        },
	        aggregate: [
                { field: "Production", aggregate: "sum" },
                { field: "PLF", aggregate: "average" },
                { field: "ScadaAvail", aggregate: "average" },
                { field: "OkTime", aggregate: "average" },
            ],
	      },
	      sortable: true,
	      filterable: false,
	      excel:{
            fileName:title,
            allPages:true, 
            filterable:true
	      },
	      pageable: {
                pageSize: 10,
                input: true, 
          },
	      scrollable: false,
	      columns: [
	        { title: "Date", field: "DateInfo.DateId", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-left" }, width: 120, template: "#= kendo.toString(moment.utc(DateInfo.DateId).format('DD-MMM-YYYY'), 'dd-MMM-yyyy') #",footerTemplate: "Total",footerAttributes: {style: "text-align:center;"}}, 
	        { title: "Turbine", field: "Turbine", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120},
	        { title: "Generation (MWh)", field: "Production", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120,template: "#= kendo.toString(Production / 1000, 'n2') #",footerTemplate: "#=kendo.toString(sum/1000, 'n2')#",footerAttributes: {style: "text-align:center;"}},
	        { title: "PLF (%)", field: "PLF", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:p2}", width: 120,footerTemplate: "#=kendo.toString(average, 'p2')#",footerAttributes: {style: "text-align:center;"}},
	        { title: "Opr. Hrs. (HH:mm:ss)", field: "OkTime", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, width: 120, template: "#= page.fancyTimeFormat(OkTime) #" ,footerTemplate: "#=page.fancyTimeFormat(average)#",footerAttributes: {style: "text-align:center;"}},
	        { title: "Lull Hrs. (HH:mm:ss)", field: "LULL", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120,template: "#= page.fancyTimeFormat(0) #" ,footerTemplate: "#=page.fancyTimeFormat(0)#",footerAttributes: {style: "text-align:center;"}},
	        { title: "Breakdown Hrs. (HH:mm:ss)", field: "LULL", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:n2}", width: 120,template: "#= page.fancyTimeFormat(0) #" ,footerTemplate: "#=page.fancyTimeFormat(0)#",footerAttributes: {style: "text-align:center;"}},
	        { title: "Data Avail. (%)", field: "ScadaAvail", headerAttributes: { style:"text-align: center" }, attributes:{ class:"align-center" }, format: "{0:p2}", width: 120,footerTemplate: "#=kendo.toString(average, 'p2')#",footerAttributes: {style: "text-align:center;" }},
	      ],
	      dataBound : function(){
	      		app.loading(false);
	      }
	    });

	    
	},
};



$(function (){
	setTimeout(function(){
		Data.LoadData();
	},500);
	
	
	$('#btnRefresh').on('click', function () {
        // fa.checkTurbine();
        Data.LoadData();
    });

     $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
        }
    });

});

