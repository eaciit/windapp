'use strict';

viewModel.AnalyticComparison = new Object();
var page = viewModel.AnalyticComparison;

var keys = [
	{ "value": "ActualProduction", "text": "Production", "status": true, "unit": "MWh" },
	{ "value": "ActualPLF", "text": "PLF", "status": true, "unit": "%" },
	{ "value": "TotalAvailability", "text": "Total Avail.", "status": true, "unit": "%" },
	{ "value": "GridAvailability", "text": "Grid Avail.", "status": false, "unit": "%" },
	{ "value": "MachineAvailability", "text": "Machine Av.", "status": false, "unit": "%" },
	{ "value": "DataAvailability", "text": "Data Avail.", "status": false, "unit": "%" },
	// { "value": "MTTRMTTF", "text": "MTTR / MTTF", "status": false , "unit": "MWh"},
	{ "value": "P50Generation", "text": "P50 Gen.", "status": false, "unit": "MWh" },
	{ "value": "P50PLF", "text": "P50 PLF", "status": false, "unit": "%" },
	{ "value": "P75Generation", "text": "P75 Gen.", "status": false, "unit": "MWh" },
	{ "value": "P75PLF", "text": "P75 PLF", "status": false, "unit": "%" },
	{ "value": "P90Generation", "text": "P90 Gen.", "status": false, "unit": "MWh" },
	{ "value": "P90PLF", "text": "P90 PLF", "status": false, "unit": "%" },
	{ "value": "Revenue", "text": "Revenue", "status": true, "unit": "Rupee", "altUnit": "Lacs" },
];

page.turbineList = ko.observableArray([]);
page.projectList = ko.observableArray([]);
/*page.periodList = ko.observableArray([
	{ "value": "custom", "text": "Custom" },
	{ "value": "monthly", "text": "Monthly" },
	{ "value": "yearly", "text": "Yearly" },
	{ "value": "lastweek", "text": "Last Week" },
	{ "value": "lastmonth", "text": "Last Month" },
	{ "value": "lastthreemonth", "text": "Last 3 Months" },
	{ "value": "lastyear", "text": "Last Year" },
]);*/

page.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Last 24 hours" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);

page.keyComparison = ko.observableArray(keys);
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.turbine = ko.observableArray([]);
page.project = ko.observable();
page.period = ko.observable();
page.defaultId = ko.observable();
page.selectedKeys = ko.observableArray([]);

page.views = ko.observableArray([]);
page.viewList = ko.observableArray([]);
page.selectedView = ko.observable();
var filterList = [];
var paramViews = {};
var idList = [{}];
var lastPeriod = "";
var turbineval = [];

page.getData = function () {
	$.when(
		toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildate", {}, function (res) {
	        var minDatetemp = new Date(res.ScadaData[0]);
	        var maxDatetemp = new Date(res.ScadaData[1]);
	        $('#availabledatestartscada').html(kendo.toString(moment.utc(minDatetemp).format('DD-MMMM-YYYY')));
	        $('#availabledateendscada').html(kendo.toString(moment.utc(maxDatetemp).format('DD-MMMM-YYYY')));
	    }),
		app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
			if (!app.isFine(res)) {
				return;
			}
			if (res.data.length == 0) {
				res.data = [];;
				page.turbineList([{ value: "", text: "" }]);
			} else {
				var datavalue = [];
				if (res.data.length > 0) {
					var allturbine = {}
					$.each(res.data, function (key, val) {
						turbineval.push(val);
					});
					allturbine.value = "All Turbine";
					allturbine.text = "All Turbines";
					datavalue.push(allturbine);
					$.each(res.data, function (key, val) {
						var data = {};
						data.value = val;
						data.text = val;
						datavalue.push(data);
					});
				}
				page.turbineList(datavalue);
			}
		}),
		app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
			if (!app.isFine(res)) {
				return;
			}
			if (res.data.length == 0) {
				res.data = [];;
				page.projectList([{ value: "", text: "" }]);
			} else {
				var datavalue = [];
				if (res.data.length > 0) {
					$.each(res.data, function (key, val) {
						var data = {};
						data.value = val;
						data.text = val;
						datavalue.push(data);
					});
				}
				page.projectList(datavalue);
			}
		}), page.getViews()).then(function () {
			page.generateElementFilter(1, "first load", 0);
			$('.key-list input[type="checkbox"]').change();
		});
}

page.populateTurbine = function (id) {
	setTimeout(function () {
		$('#turbineList-' + id).data('kendoMultiSelect').value(["All Turbine"])
	}, 500);
};
page.checkTurbine = function (id) {
	var arr = $('#turbineList-' + id).data('kendoMultiSelect').value();
	var index = arr.indexOf("All Turbine");
	if (index == 0 && arr.length > 1) {
		arr.splice(index, 1);
		$('#turbineList-' + id).data('kendoMultiSelect').value(arr)
	} else if (index > 0 && arr.length > 1) {
		$("#turbineList-" + id).data("kendoMultiSelect").value(["All Turbine"]);
	} else if (arr.length == 0) {
		$("#turbineList-" + id).data("kendoMultiSelect").value(["All Turbine"]);
	}
}

page.populateProject = function (id) {
	setTimeout(function () {
		$("#projectList-" + id).data("kendoDropDownList").value("Tejuva");
		page.project = $("#projectList-" + id).data("kendoDropDownList").value();
	}, 500);
};

/*page.showHidePeriod = function (idx) {

	var id = (idx == null ? 1 : idx);
	var period = $('#periodList-' + id).data('kendoDropDownList').value();

	var tempDate = new Date(app.currentDateData);
	var endMonthDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth() - 1,
		1, 0, 0, 0, 0));
	var startMonthDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth(),
		tempDate.getDate() - 1, 0, 0, 0, 0));
	var endDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth(),
		tempDate.getDate() - 1, 0, 0, 0, 0));
	var startDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth(),
		tempDate.getDate() - 1, 0, 0, 0, 0));

	if (period == "custom") {
		$('#label-start-' + id).html("Start date");
		$('#label-end-' + id).html("End date");
		$("#show_hide" + id).show();
		$('#dateStart-' + id).data('kendoDatePicker').setOptions({
			start: "month",
			depth: "month",
			format: 'dd-MMM-yyyy'
		});
		$('#dateEnd-' + id).data('kendoDatePicker').setOptions({
			start: "month",
			depth: "month",
			format: 'dd-MMM-yyyy'
		});
		if (lastPeriod == "monthly") {
			$('#dateStart-' + id).data('kendoDatePicker').value(endMonthDate);
			$('#dateEnd-' + id).data('kendoDatePicker').value(startMonthDate);
		} else if (lastPeriod == "yearly") {
			$('#dateStart-' + id).data('kendoDatePicker').value(new Date(Date.UTC(tempDate.getFullYear(), 0, 1, 0, 0, 0, 0)));
			$('#dateEnd-' + id).data('kendoDatePicker').value(new Date(Date.UTC(tempDate.getFullYear(), 11, 31, 0, 0, 0, 0)));
		}
	} else if (period == "monthly") {
		lastPeriod = "monthly";
		$('#label-start-' + id).html("Start month");
		$('#label-end-' + id).html("End month");
		$('#dateStart-' + id).data('kendoDatePicker').setOptions({
			start: "year",
			depth: "year",
			format: "MMMM yyyy"
		});
		$('#dateEnd-' + id).data('kendoDatePicker').setOptions({
			start: "year",
			depth: "year",
			format: "MMMM yyyy"
		});

		$('#dateStart-' + id).data('kendoDatePicker').value(endMonthDate);
		$('#dateEnd-' + id).data('kendoDatePicker').value(startMonthDate);

		$("#show_hide" + id).show();
	} else if (period == "yearly") {
		lastPeriod = "yearly";
		$('#dateStart-' + id).data('kendoDatePicker').setOptions({
			start: "decade",
			depth: "decade",
			format: "yyyy"
		});
		$('#dateEnd-' + id).data('kendoDatePicker').setOptions({
			start: "decade",
			depth: "decade",
			format: "yyyy"
		});
		$('#label-start-' + id).html("Start year");
		$('#label-end-' + id).html("End year");
		$('#dateStart-' + id).data('kendoDatePicker').value(new Date(Date.UTC(tempDate.getFullYear(), 0, 1, 0, 0, 0, 0)));
		$('#dateEnd-' + id).data('kendoDatePicker').value(new Date(Date.UTC(tempDate.getFullYear(), 11, 31, 0, 0, 0, 0)));
		$("#show_hide" + id).show();
	} else {
		if (period == "lastweek") {
			startDate.setDate(endDate.getDate() - 7);
		} else if (period == "lastmonth") {
			startDate.setDate(endDate.getDate() - 30);
		} else if (period == "lastthreemonth") {
			startDate.setDate(endDate.getDate() - 90);
		} else if (period == "lastyear") {
			startDate.setDate(endDate.getDate() - 365);
		}

		$('#dateStart-' + id).data('kendoDatePicker').value(startDate);
		$('#dateEnd-' + id).data('kendoDatePicker').value(endDate);
		$("#show_hide" + id).hide();
	}
}*/
page.checkCompleteDate = function(id){
    var period = $('#periodList-' + id).data('kendoDropDownList').value();

    var monthNames = moment.months();

    var currentDateData = moment(app.currentDateData).format('YYYY-MM-DD');
    var startDate = $('#dateStart-' + id).data('kendoDatePicker').value();
	var endDate = $('#dateEnd-' + id).data('kendoDatePicker').value();
    var today = moment().format('YYYY-MM-DD');
    var thisMonth = moment().get('month');
    var firstDayMonth = moment(new Date(startDate.getFullYear(), startDate.getMonth(), 1)).format("YYYY-MM-DD");
    var lastDayMonth = moment(new Date(endDate.getFullYear(), endDate.getMonth() + 1, 0)).format("YYYY-MM-DD"); 
    var firstDayYear = moment().startOf('year').format('YYYY-MM-DD');
    var endDayYear = moment().endOf('year').format('YYYY-MM-DD');

    var dateStart = moment(startDate).format('YYYY-MM-DD');
    var dateEnd = moment(endDate).format('YYYY-MM-DD'); 

    $('#period-info-'+id).html("");

    if(period === 'custom'){
        if((dateEnd > currentDateData) && (dateStart  > currentDateData)){
            $('#period-info-'+id).html("<span>* Incomplete period range on start date and end date</span>");
        }else if(dateStart  > currentDateData){
            $('#period-info-'+id).html("<span>* Incomplete period range on start date</span>");
        }else if(dateEnd  > currentDateData){
            $('#period-info-'+id).html("<span>* Incomplete period range on end date</span>");
        }else{
            $('#period-info-'+id).html("");
        }
    }else if(period === 'annual'){
        if((moment(endDate).get('year') == moment(app.currentDateData).get('year')) && (currentDateData < today)){
             $('#period-info-'+id).html("<span>*Incomplete period range in end year</span>");
        }else{
             $('#period-info-'+id).html("");
        }
    }else if(period === 'monthly'){
        if((dateEnd > currentDateData) && (dateStart  > currentDateData)){
            $('#period-info-'+id).html("<span>*Incomplete period range in start month and start month</span>");
        }else if(dateStart > currentDateData){
            $('#period-info-'+id).hrml("<span>*Incomplete period range in start month</span>");
        }else if(dateEnd > currentDateData){
             $('#period-info-'+id).html("<span>*Incomplete period range in end month</span>");
        }else{
             $('#period-info-'+id).html("");
        }
    }else{
         $('#period-info-'+id).html("");
    }


}
page.showHidePeriod = function (idx) {

	var id = (idx == null ? 1 : idx);
	var period = $('#periodList-' + id).data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),1,0,0,0,0));
    var endMonthDate = new Date(app.toUTC(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0,1,0,0,0,0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0,1,0,0,0,0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-1,0,0,0,0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));

	if (period == "custom") {
		$('#label-start-' + id).html("Start date");
		$('#label-end-' + id).html("End date");
		$("#show_hide" + id).show();
		$('#dateStart-' + id).data('kendoDatePicker').setOptions({
			start: "month",
			depth: "month",
			format: 'dd-MMM-yyyy'
		});
		$('#dateEnd-' + id).data('kendoDatePicker').setOptions({
			start: "month",
			depth: "month",
			format: 'dd-MMM-yyyy'
		});
		$('#dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
		// if (lastPeriod == "monthly") {
		// 	$('#dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
		// 	$('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
		// } else if (lastPeriod == "annual") {
  //           $('#dateStart-' + id).data('kendoDatePicker').value(startYearDate);
  //           $('#dateEnd-' + id).data('kendoDatePicker').value(endYearDate);
  //       }
	} else if (period == "monthly") {
		lastPeriod = "monthly";
		$('#label-start-' + id).html("Start month");
		$('#label-end-' + id).html("End month");
		$('#dateStart-' + id).data('kendoDatePicker').setOptions({
			start: "year",
			depth: "year",
			format: "MMMM yyyy"
		});
		$('#dateEnd-' + id).data('kendoDatePicker').setOptions({
			start: "year",
			depth: "year",
			format: "MMMM yyyy"
		});

		$('#dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
		$('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);

		$("#show_hide" + id).show();
	} else if (period == "annual") {
		$('#label-start-' + id).html("Start year");
		$('#label-end-' + id).html("End year");
		$("#show_hide" + id).show();

        $('#dateStart-' + id).data('kendoDatePicker').setOptions({
            start: "decade",
            depth: "decade",
            format: "yyyy"
        });
        $('#dateEnd-' + id).data('kendoDatePicker').setOptions({
            start: "decade",
            depth: "decade",
            format: "yyyy"
        });

       $('#dateStart-' + id).data('kendoDatePicker').value(startYearDate);
       $('#dateEnd-' + id).data('kendoDatePicker').value(endYearDate);

        $(".show_hide").show();
    } else {
	    if(period == 'last24hours'){
             $('#dateStart-' + id).data('kendoDatePicker').value(last24hours);
             $('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
        }else if(period == 'last7days'){
             $('#dateStart-' + id).data('kendoDatePicker').value(lastweek);
             $('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
        }
		$("#show_hide" + id).hide();
	}
	 lastPeriod = period;
}

page.ChangeKeyList = function (e) {
	var title = $(e).attr("text");
	var val = $(e).val();
	var value = title,
		$list = $(".content-selected-key"),
		$list2 = $(".result-content-selected-key");
	if (e.checked) {
		page.selectedKeys.push(val);
		$list.find('div[data-value="default"]').remove();
		$list.append("<div data-value='" + val + "' class='selected-key'>" + value + "</div>");
		$list2.find('div[data-value="default"]').remove();
		$list2.append("<div data-value='" + val + "' class='result-selected-key'> &nbsp;<span style='float:right'> 0 </span></div>");
	} else {
		page.selectedKeys.remove(val);
		$list.find('div[data-value="' + val + '"]').slideUp("fast", function () {
			$(this).remove();
			if ($('.content-selected-key').children().length == 0) {
				$list.append("<div data-value='default' class='selected-key data-default'><i>* No key selected</i></div>");
			}
		});
		$list2.find('div[data-value="' + val + '"]').slideUp("fast", function () {
			$(this).remove();
			if ($('.result-content-selected-key').children().length == 0) {
				$list2.append("<div data-value='default' class='result-selected-key data-default'><i>* No result</i></div>");
			}
		});
	}
	page.refreshAll();
}

page.getRandomId = function () {
	return page.randomNumber() + page.randomNumber() + page.randomNumber() + page.randomNumber();
}

page.randomNumber = function () {
	return Math.floor((1 + Math.random()) * 0x10000)
		.toString(16)
		.substring(1);
}

page.InitDefaultValue = function (id) {
	$("#periodList-" + id).data("kendoDropDownList").value("custom");
	$("#periodList-" + id).data("kendoDropDownList").trigger("change");

	var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));
    var lastEndDate = new Date(app.toUTC(maxDateData));

	$('#dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
	$('#dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
}

page.InitViewsValue = function (id, data) {
	$("#periodList-" + id).data("kendoDropDownList").value(data.Period);
	$("#periodList-" + id).data("kendoDropDownList").trigger("change");
	var tempDate = new Date(data.DateStart);
	var lastStartDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth(),
		tempDate.getDate(), 0, 0, 0, 0));
	$('#dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
	tempDate = new Date(data.DateEnd);
	var lastEndDate = new Date(Date.UTC(tempDate.getFullYear(), tempDate.getMonth(),
		tempDate.getDate(), 0, 0, 0, 0));
	$('#dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
	$('#projectList-' + id).data('kendoDropDownList').value(data.Project);
	$("#projectList-" + id).data("kendoDropDownList").trigger("change");
	$("#turbineList-" + id).data("kendoMultiSelect").value(data.Turbine);

	var kunci = [{}];
	kunci = keys;
	$.each(kunci, function (i, key) { /*unchecked all*/
		var cek = $('#chk-key-' + i).prop('checked', false);
		page.ChangeKeyByViews(key.value, key.text, false);

	});
	$.each(data.Keys, function (iViews, val) {
		$.each(kunci, function (i, key) {
			if (val == key.value) {
				$('#chk-key-' + i).prop('checked', true);
				page.ChangeKeyByViews(key.value, key.text, true);
			}
		});
	});
	page.refreshAll();
}

page.ChangeKeyByViews = function (keyValue, text, isChecked) {
	var title = text;
	var val = keyValue;
	var value = title,
		$list = $(".content-selected-key"),
		$list2 = $(".result-content-selected-key");
	if (isChecked) {
		page.selectedKeys.push(val);
		$list.find('div[data-value="default"]').remove();
		$list.append("<div data-value='" + val + "' class='selected-key'>" + value + "</div>");
		$list2.find('div[data-value="default"]').remove();
		$list2.append("<div data-value='" + val + "' class='result-selected-key'> &nbsp;<span style='float:right'> 0 </span></div>");
	} else {
		page.selectedKeys.remove(val);
		$list.find('div[data-value="' + val + '"]').slideUp("fast", function () {
			$(this).remove();
			if ($('.content-selected-key').children().length == 0) {
				$list.append("<div data-value='default' class='selected-key data-default'><i>* No key selected</i></div>");
			}
		});
		$list2.find('div[data-value="' + val + '"]').slideUp("fast", function () {
			$(this).remove();
			if ($('.result-content-selected-key').children().length == 0) {
				$list2.append("<div data-value='default' class='result-selected-key data-default'><i>* No result</i></div>");
			}
		});
	}
}

page.checkKeyFromSavedView = function (val) {
	var data = {};
	$.each(idList, function (i, id) {
		$(id.filter).remove();
		$(id.result).remove();
	});
	idList = [{}];

	$.each(val.Filters, function (i, value) {
		data = { "Keys": [], "DateStart": "", "DateEnd": "", "Project": "", "Turbine": [], "Period": "" };
		data.DateStart = value.DateStart;
		data.DateEnd = value.DateEnd;
		data.Project = value.Project;
		data.Turbine = value.Turbine;
		data.Period = value.Period;
		data.Keys = val.Keys;

		if (i == 0) {
			page.generateElementFilter(1, "views", data);
		} else {
			page.generateElementFilter(null, "views", data);
		}
	});
}

page.generateElementFilter = function (id_element, source, dataViews) {
	var id = (id_element == null ? page.getRandomId() : id_element);
	var ids = {};
	ids = { "filter": "#td-form-filter-" + id, "result": "#td-filter-result-" + id };
	idList.push(ids);
	page.defaultId = id;
	var formFilter = '<td class="column-filter-form" id="td-form-filter-' + id + '">' +
		'<div style="overflow-y:auto; height:100%">' +
		'<label class="col-md-1 col-sm-1 control-label" style="width:70px;">Project</label>' +
		'<select class="col-md-1 col-sm-1" id="projectList-' + id + '" name="" style="width:110px"></select>' +
		'<div class="clearfix">&nbsp;</div>' +
		'<label class="col-md-1 col-sm-1 control-label" style="width:70px;">Turbine</label>' +
		'<select class="col-md-1 col-sm-1" id="turbineList-' + id + '" name="" style="width:140px"></select>' +
		'<div class="clearfix">&nbsp;</div>' +
		'<label class="control-label col-md-1 col-sm-1" style="width:70px;">Period</label>&nbsp;' +
		'<select class="col-md-1 col-sm-1" id="periodList-' + id + '" name="" style="width:110px"></select>' +
		'<div class="clearfix">&nbsp;</div>' +
		'<span id="show_hide' + id + '">' +
		'<label class="col-md-1 col-sm-1 control-label" id="label-start-' + id + '" style="width:120px;margin-right:-50px">Start date</label>' +
		'<input class="col-md-1 col-sm-1" type="text" id="dateStart-' + id + '"  />' +
		'<div class="clearfix">&nbsp;</div>' +
		'<label class="col-md-1 col-sm-1 control-label" id="label-end-' + id + '" style="width:120px;margin-right:-50px">End date</label>' +
		'<input class="col-md-1 col-sm-1" type="text" id="dateEnd-' + id + '" " />&nbsp;' +
		'&nbsp;' +
		'</span>' +
		'<div class="clearfix">&nbsp;</div>' +
		'<div class="col-md-12">'+
			'<div class="col-md-7">'+
				'<div class="pull-left period-info-label" id="period-info-'+id+'"></div>'+
			'</div>'+
			'<div class="col-md-5">'+
				'<div class="pull-right" style="margin-top:10px;">' +
					'<button class="btn btn-sm btn-primary tooltipster tooltipstered" id="btn-refresh-' + id + '" onClick="page.refreshFilter(\'' + id + '\')" title="Refresh Filter" data-loading-text="<i class=\'fa fa-circle-o-notch fa-spin\'></i> Loading"><i class="fa fa-refresh"></i></button>&nbsp;' +
					'<button class="btn btn-sm btn-danger tooltipster tooltipstered" onClick="page.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (id == 1 ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
				'</div>' +
			'</div>'+
		'</div>'+
		'</div>' +
		'</td>';


	setTimeout(function () {
		var $filterResults = '<td class="column-result" style="border-right: 1.5px solid #ddd;" id="td-filter-result-' + id + '">' +
			'<div class="result-content-selected-key" style="padding: 5px 5px 5px 10px;"></div>' +
			'</td>';

		$('.filter-results').append($filterResults);

		if ($('.content-selected-key').children().length > 0) {

			$.each($('.content-selected-key').children(), function (i, val) {
				var value = $(val).attr("data-value");
				if (val.className == "selected-key data-default") {
					var $list = '<div class="result-selected-key" data-value="default"><i>* No result</i></div>';
				} else {
					var $list = '<div data-value="' + value + '" class="result-selected-key"> &nbsp;<span style="float:right"> 0 </span></div>';
				}
				$("#td-filter-result-" + id).find($('.result-content-selected-key')).append($list);
			});

		} else {
			var $list = '<div class="result-selected-key" data-value="default"><i>* No result</i></div>';
			$("#td-filter-result-" + id).find($('.result-content-selected-key')).append($list);
		}

		$(".filter-form").append(formFilter);

		$("#projectList-" + id).kendoDropDownList({
			dataValueField: 'value',
			dataTextField: 'text',
			optionLabel: 'Please Select',
			suggest: true,
			dataSource: page.projectList(),
		});
		$("#turbineList-" + id).kendoMultiSelect({
			dataSource: page.turbineList(),
			dataValueField: 'value',
			dataTextField: 'text',
			suggest: true,
			change: function () { page.checkTurbine(id) }
		});
		$("#periodList-" + id).kendoDropDownList({
			dataSource: page.periodList(),
			dataValueField: 'value',
			dataTextField: 'text',
			suggest: true,
			change: function () { page.showHidePeriod(id) }
		});

		$('#dateStart-' + id).kendoDatePicker({
			value: new Date(),
			format: 'dd-MMM-yyyy',
			min: new Date("2013-01-01"),
			max:new Date(),
		});

		$('#dateEnd-' + id).kendoDatePicker({
			value: new Date(),
			format: 'dd-MMM-yyyy',
			min: new Date("2013-01-01"),
			max:new Date(),
		});

		if (source != "views") {
			page.populateTurbine(id);
			page.populateProject(id);
			page.InitDefaultValue(id);
		} else {
			page.InitViewsValue(id, dataViews);
		}

		if(source == "first load"){
			setTimeout(function () {
				page.refreshFilter('1');								
		    }, 100);
		}

		
	}, 300);
}

page.removeFilter = function (id) {
	$("#td-form-filter-" + id).remove();
	$("#td-filter-result-" + id).remove();
}

page.refreshFilter = function (id) {
	$('#btn-refresh-'+id).button("loading");
	setTimeout(function(){
		var startdate = $('#dateStart-' + id).data('kendoDatePicker').value();
		var enddata = $('#dateEnd-' + id).data('kendoDatePicker').value();
		var period = $('#periodList-' + id).data('kendoDropDownList').value();
		if (startdate > enddata) {
			toolkit.showError("Invalid Date Range Selection");
			return;
		} else {
			var turbine = [];
			var isAllTurbine = false;

			if ($("#turbineList-" + id).data("kendoMultiSelect").value().indexOf("All Turbine") >= 0) {
				isAllTurbine = true;
				// $.each($("#turbineList-" + id).data("kendoMultiSelect").dataSource.options.data, function (idx, val) {
				// 	if (val.value != "All Turbine") {
				// 		turbine.push(val.value);
				// 	}
				// })
				turbine = [];
			} else {
				turbine = $("#turbineList-" + id).data("kendoMultiSelect").value();
			}

			var param = {
				period: period,
				DateStart: $('#dateStart-' + id).data('kendoDatePicker').value(),
				DateEnd: $('#dateEnd-' + id).data('kendoDatePicker').value(),
				Turbine: turbine,
				Project: $("#projectList-" + id).data("kendoDropDownList").value(),
				Keys: page.selectedKeys(),
			};

			toolkit.ajaxPost(viewModel.appName + "analyticcomparison/getdata", param, function (res) {
				if (!app.isFine(res)) {
					return;
				}

				$.each(res.data, function (idx, val) {
					$.each(keys, function (iSel, valSel) {
						if (valSel.value == idx) {
							if (valSel.value == "Revenue") {
								if (val < 100000) {
									$('#td-filter-result-' + id + ' > div > div[data-value="' + idx + '"]  > span').html(kendo.toString(val, "n") + " " + valSel.unit);
								} else {
									val = val / 100000
									$('#td-filter-result-' + id + ' > div > div[data-value="' + idx + '"]  > span').html(kendo.toString(val, "n") + " " + valSel.altUnit);
								}
							} else {
								$('#td-filter-result-' + id + ' > div > div[data-value="' + idx + '"]  > span').html(kendo.toString(val, "n") + " " + valSel.unit);
							}
						}
					});
				});
				$('#btn-refresh-'+id).button("reset");
			});
			if (isAllTurbine) {
				param.Turbine = ["All Turbine"];
			}
			param.Period = $("#periodList-" + id).data("kendoDropDownList").value();
			filterList.push(param);
			page.checkCompleteDate(id);
		}
	});
}

page.refreshAll = function () {
	$(".refresh-all").button('loading');
	setTimeout(function(){
		filterList = [];
		$.each($('button[id^=btn-refresh-]'), function (idx, btn) {
			btn.click();
		})
		paramViews = {
			OldName: page.selectedView.Name,
			Name: $("#inputViewName").val(),
			Keys: page.selectedKeys(),
			Filters: filterList
		}
		$(".refresh-all").button('reset');
	},500);
	
}

page.loadView = function () {
	var selectedVal = $("#savedViews").data("kendoDropDownList").value();
	if (selectedVal != "") {
		page.selectedView = null;
		$.each(page.views(), function (i, val) {
			if (val.Name == selectedVal) {
				page.selectedView = val;
				page.checkKeyFromSavedView(val);
			}
		});
	}
}

page.getViews = function () {
	page.viewList = [];
	page.viewList.push({
		value: "",
		text: "Please Select"
	})

	app.ajaxPost(viewModel.appName + "userpreferences/getanalysisstudioviews", "", function (res) {
		if (!app.isFine(res)) {
			return;
		}

		page.views(res.data);
		$.each(page.views(), function (i, val) {
			page.viewList.push({
				value: val.Name,
				text: val.Name
			})
		});

		$("#savedViews").data("kendoDropDownList").dataSource.data(page.viewList);
		$("#savedViews").data("kendoDropDownList").dataSource.query();
		if ($("#savedViews").data("kendoDropDownList").value() == "") {
			$("#savedViews").data("kendoDropDownList").select(0);
		}
	});
}

page.modalSaveView = function () {
	var selectedVal = $("#savedViews").data("kendoDropDownList").value();
	$("#inputViewName").val(selectedVal);


	if (page.viewList.length == 4 && selectedVal == "") {
		toolkit.showError("Only 3 views are allowed");
	} else if (selectedVal != "") {
		page.refreshAll();
		$("#modal-views-title").html("Update View");
		page.ShowModal('modalForm', 'show');
	} else {
		page.refreshAll();
		$("#modal-views-title").html("Create New View");
		page.ShowModal('modalForm', 'show');
	}
}

page.saveView = function () {
	page.ShowModal('modalForm', 'hide');
	var selectedVal = $("#savedViews").data("kendoDropDownList").value();
	if (selectedVal == "") {
		paramViews.OldName = "";
	}
	paramViews.Name = $("#inputViewName").val();

	app.ajaxPost(viewModel.appName + "userpreferences/saveanalysisstudioviews", paramViews, function (res) {
		if (!app.isFine(res) || res.data == null) {
			toolkit.showError("Error Occur when save the KPI");
			return;
		}

		page.viewList = [];
		page.viewList.push({
			value: "",
			text: "Please Select"
		})

		swal({
			title: "Info",
			type: "success",
			text: "Data Successfully Saved",
		}, function () { });

		// page.views(res.data);

		$.each(page.views(), function (i, val) {
			page.viewList.push({
				value: val.Name,
				text: val.Name
			})
		});

		var idx = $("#savedViews").data("kendoDropDownList").select();

		$("#savedViews").data("kendoDropDownList").dataSource.data(page.viewList);
		$("#savedViews").data("kendoDropDownList").dataSource.query();

		page.getViews();

		setTimeout(function () {
			$("#savedViews").data("kendoDropDownList").select(idx);
		}, 100);
	});
}

page.ShowModal = function (modalId, showhide) {
	if (showhide == 'show') {
		$('#' + modalId).appendTo("body").modal({
			backdrop: 'static',
			keyboard: false,
			show: showhide
		});
	} else {
		$('#' + modalId).modal('hide');
	}
}

vm.currentMenu('Analytics Studio');
vm.currentTitle('Analytics Studio');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Analytics Studio', href: viewModel.appName + 'page/analyticavailability' }]);

$(document).ready(function () {
	$('#btnSaveView').on('click', function () {
		page.modalSaveView();
	});
	$('#savedViews').kendoDropDownList({
		data: [],
		dataValueField: 'value',
		dataTextField: 'text',
		change: function () { page.loadView() },
	});
	page.getData();

});