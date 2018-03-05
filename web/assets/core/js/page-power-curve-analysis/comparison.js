'use strict';
var cm = {};
cm.loadFirstTime = ko.observable(true);

cm.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);

cm.turbineList = ko.observableArray([]);
cm.projectList = ko.observableArray([]);
cm.dateStart = ko.observable();
cm.dateEnd = ko.observable();
cm.project = ko.observable();
cm.sScater = ko.observable(false);

cm.rawturbine = ko.observableArray([]);
cm.rawproject = ko.observableArray([]);
cm.IDList = [];
cm.countList = 0;
cm.LastFilter;
cm.TableName;
cm.FieldList;
cm.ContentFilter;


cm.reset = function(){
    cm.loadFirstTime(true);
}
cm.refresh = function() {
   if(this.loadFirstTime()) {
        cm.countList = 0;
        this.initElementEvents();
	    this.loadFirstTime(false);
	    this.setProjectTurbine();
		this.populateProject();
		this.loadData();
   	}
}

cm.internalRefresh = function(reloadData) {
    if(reloadData==undefined) {
        reloadData = true;
    }

    if(reloadData) {
       this.initChart();
    }
}

cm.loadData = function(){
    cm.IDList = [];
	$(".filter-part").html("");

    this.generateElementFilter(null, "default1");
    this.generateElementFilter(null, "default2");
}

cm.getPDF = function(selector){
    app.loading(true);

    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            margin: {
                left   : "5mm",
                top    : "5mm",
                right  : "10mm",
                bottom : "5mm"
            },
        });
      kendo.drawing.pdf.saveAs(group,"PC Comparison.pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}


cm.populateTurbine = function () {
    if (cm.rawturbine().length == 0) {
        cm.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];

        $.each($("#cm-project").data("kendoMultiSelect").value(), function(i, project){
            $.each(cm.rawturbine(), function(key, val){
                if(project == val.Project){
                    var data = {};
                    data.value = val.Value + "<>" + val.Project;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }
            });
        });
        cm.turbineList(datavalue);
    }

    var i = 0;

    cm.IDList.forEach(function(id) {
        $('#turbineList-'+id).data('kendoDropDownList').setDataSource(new kendo.data.DataSource({ data: cm.turbineList() }));
        $('#turbineList-'+id).data('kendoDropDownList').select(i);
        i++;
    });
};

cm.populateProject = function () {
    if (cm.rawproject().length == 0) {
        cm.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];        
        $.each(cm.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Value;
            datavalue.push(data);
        });
        cm.projectList(datavalue);
    }
    $("#cm-project").data("kendoMultiSelect").setDataSource(cm.projectList());
    $('#cm-project').data('kendoMultiSelect').value([fa.project]);
    $("#cm-project").data("kendoMultiSelect").trigger("change");
};

cm.getRandomId = function () {
    return cm.randomNumber() + cm.randomNumber() + cm.randomNumber() + cm.randomNumber();
}

cm.randomNumber = function () {
    return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
}

cm.generateElementFilter = function (id_element, source) {
    cm.countList++;
    var id = (id_element == null ? cm.getRandomId() : id_element);
    var isDefault = false;
    if(source.indexOf("default") >= 0) {
        isDefault = true;
    }
    if(cm.IDList.length == 5) {
        return;
    }

    cm.IDList.push(id);
    var isLast = false;
    if(cm.IDList.length == 5) {
        isLast = true;
        $(".button-add").hide();
    }

	var formFilter =    '<div class="row dynamic-filter" id="filter-form-'+ id + '" data-count="'+ cm.countList +'">' +
	                        '<div class="mgb10">' +
	                            '<div class="col-md-3 no-padding">' +
	                                '<select class="cm-turbine-list" id="turbineList-' + id + '" name="table" multiple="multiple"></select>' +
	                            '</div>' +
	                            '<div class="col-md-9 no-padding">' +
	                                '<select class="period-list" id="periodList-' + id + '" name="table"></select>' +
	                                '<span class="custom-period" id="show_hide-' + id + '">' +
	                                    '<input type="text" id="dateStart-' + id + '"/>' +
	                                    '<label>&nbsp;&nbsp;&nbsp;to&nbsp;&nbsp;&nbsp;</label>' +
	                                    '<input type="text" id="dateEnd-' + id + '"/>' +
	                                '</span>' +
	                            '</div>' +
	                            '<button class="btn btn-sm btn-danger tooltipster tooltipstered remove-btn" onClick="cm.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (isDefault ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
	                        '</div>'
	                    '</div>';
	var versusFilter = '<div class="versus-wrapper" data-count="'+ cm.countList +'"><div class="versus">vs</div></div>';

    setTimeout(function () {
        $(".filter-part").append(formFilter);
        $(".filter-part").append(versusFilter);

        $("#turbineList-" + id).kendoDropDownList({
            dataValueField: 'value',
            dataTextField: 'label',
            suggest: true,
            dataSource: cm.turbineList(),
            change: function(){
                cm.initChart();
            }
        });     

        $('#turbineList-'+id).data('kendoDropDownList').select(cm.countList);

        $("#periodList-" + id).kendoDropDownList({
            dataSource: cm.periodList(),
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { 
                cm.showHidePeriod(id) 
            }
        });

        $('#dateStart-' + id).kendoDatePicker({
            value: new Date(),
            format: 'dd-MMM-yyyy',
            min: new Date("2013-01-01"),
            max:new Date(),
            change: function(){
                cm.initChart();
            }
        });

        $('#dateEnd-' + id).kendoDatePicker({
            value: new Date(),
            format: 'dd-MMM-yyyy',
            min: new Date("2013-01-01"),
            max:new Date(),
            change: function(){
                cm.initChart();
            }
        });

        cm.InitDefaultValue(id);

        if(source == "default2"){
            // setTimeout(function () {
            //     cm.initChart();                           
            // }, 500);
        }
        cm.checkElementLast();
    }, 500);
}

cm.removeFilter = function (id) {
    cm.countList--;
    $("#filter-form-" + id).remove();
    var tempList = [];
    cm.IDList.forEach(function(val){
        if (val !== id) {
            tempList.push(val);
        }
    });
    cm.IDList = tempList;
    cm.checkElementLast();

    if(cm.IDList.length < 5){
        $(".button-add").show()
    }else{
        $(".button-add").hide()
    }
}

cm.checkElementLast = function(){
    var elms = $('.dynamic-filter');
    $.each(elms, function(i, e){
        if(!$(e).hasClass('dynamic-filter-last')) {
            $(e).addClass('dynamic-filter-last');
        }
        var dataCount = parseInt($(e).attr('data-count'));
        if(dataCount < cm.countList) {
            $(e).removeClass('dynamic-filter-last');
        }
        var turbineElm = $(e).find('select.cm-turbine-list');
        var turbineElmId = turbineElm.attr('id');
        setTimeout(function(){
            $('#'+turbineElmId).data('kendoDropDownList').select(dataCount - 1);
        }, 100);
    });
    var elmvs = $('.versus-wrapper');
    $.each(elmvs, function(i, e){
        if(!$(e).hasClass('versus-last')) {
            $(e).addClass('versus-last');
        }
        var dataCount = parseInt($(e).attr('data-count'));
        if(dataCount < cm.countList) {
            $(e).removeClass('versus-last');
        }
    });
    setTimeout(function () {
        cm.initChart();                           
    }, 500);
}

cm.showHidePeriod = function (idx) {
    var id = (idx == null ? 1 : idx);
    var period = $('#periodList-' + id).data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.getDateMax(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

    if (period == "custom") {
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
    } else if (period == "monthly") {
        $('#dateStart-' + id).data('kendoDatePicker').setOptions({
            start: "year",
            depth: "year",
            format: "MMM yyyy"
        });
        $('#dateEnd-' + id).data('kendoDatePicker').setOptions({
            start: "year",
            depth: "year",
            format: "MMM yyyy"
        });

        $('#dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);

        $("#show_hide" + id).show();
    } else if (period == "annual") {
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

        $("#show_hide").show();
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
}

cm.InitDefaultValue = function (id) {
    $("#periodList-" + id).data("kendoDropDownList").value("custom");
    $("#periodList-" + id).data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
}

cm.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = tipe;
    }

    var param = {
        Filters: cm.LastFilter,
        FieldList: cm.FieldList,
        Tablename: cm.TableName,
        TypeExcel: namaFile,
        ContentFilter: cm.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: cm.LastFilterDetails,
            FieldList: cm.FieldListDetails,
            Tablename: cm.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: cm.ContentFilterDetails,
            IsSplittedSheet: isSplittedSheet,
            IsMultipleProject: isMultipleProject,
        };
    }

    var urlName = viewModel.appName + "analyticpowercurve/genexcelpowercurve";
    app.ajaxPost(urlName, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        window.location = viewModel.appName + "/".concat(res.data);
        app.loading(false);
    });
}

cm.initChart = function() {
    app.loading(true);

    var link = "analyticpowercurve/getlistpowercurvecomparison";
    var mostDateStart;
    var mostDateEnd;
    var projectList = $("#cm-project").data("kendoMultiSelect").value();
    var turbineList = [];
    var details = [];

    cm.IDList.forEach(function(id){
        var dateStart = $('#dateStart-'+id).data('kendoDatePicker').value();
            dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));

        var dateEnd  = $('#dateEnd-'+id).data('kendoDatePicker').value();
            dateEnd = new Date(Date.UTC(dateEnd.getFullYear(), dateEnd.getMonth(), dateEnd.getDate(), 0, 0, 0));
        if (mostDateStart !== undefined) {
            if (dateStart < mostDateStart) {
                mostDateStart = dateStart
            }
            if (dateEnd > mostDateEnd) {
                mostDateEnd = dateEnd;
            }
        } else {
            mostDateStart = dateStart;
            mostDateEnd = dateEnd
        }
        var splitTurbineVal = $("#turbineList-"+id).data("kendoDropDownList").value().split("<>");
        turbineList.push(splitTurbineVal[0]);

        var detail = {
            Period       : $('#periodList-'+id).data('kendoDropDownList').value(),
            Project      : splitTurbineVal[1],
            Turbine      : splitTurbineVal[0],
            DateStart    : dateStart,
            DateEnd      : dateEnd,
        };

        details.push(detail);
    });
    var param = {
        ProjectList: projectList,
        TurbineList: turbineList,
        Details:     details,
        MostDateStart: mostDateStart,
        MostDateEnd: mostDateEnd,
    }

    toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }
        var dataTurbine = res.data.Data;
        cm.LastFilter = res.data.LastFilter;
        cm.FieldList = res.data.FieldList;
        cm.TableName = res.data.TableName;
        cm.ContentFilter = res.data.ContentFilter;
        
        $('#cm-chart').html("");
        $("#cm-chart").kendoChart({
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            title: {
                text: "Power Curves",
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: true,
                align: "start",
                offsetX : 55,
                labels: {
                    margin: {
                        right : 0,
                    },
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
            },
            chartArea: {
                height: 400,
                background: 'transparent',
            },
            seriesDefaults: {
                type: "scatterLine",
                style: "smooth",
                dashType: "longDash",
                markers: {
                    visible: false,
                    size: 4,
                },
            },
            seriesColors: colorField,
            series: dataTurbine,
            categoryAxis: {
                labels: {
                    step: 1,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            valueAxis: [{
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555",
                    visible: true,
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N1",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                max: 25
            },
            yAxis: {
                title: {
                    text: "Generation (KW)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555"
                },
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -5,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N1",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
            },
            tooltip: {
                visible: true,
                template: "#= series.name #",
                shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
            pannable: {
                lock: "y"
            },
            zoomable: {
                mousewheel: {
                    lock: "y"
                },
                selection: {
                    lock: "y",
                    key: "none",
                }
            }
        });
        app.loading(false);
        if (cm.sScater()) {
            cm.getScatter(details, dataTurbine, projectList.length);
        }
        $("#cm-chart").data("kendoChart").refresh();                
    });
}

cm.getScatter = function(paramLine, dtLine, startColorIdx) {
    var turbineList = [];
    var kolor = [];
    var idx;
    app.loading(true);
    var paramList = [];
    paramLine.forEach(function(data){
        var dateStart = data.DateStart;
        var dateEnd = data.DateEnd;
        var param = {
            period: data.Period,
            dateStart: dateStart,
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: data.Turbine,
            project: data.Project,
            Color: dtLine[startColorIdx].color,
        };
        paramList.push(param);
        startColorIdx++;
    });
    var dataPowerCurves = [];
    var reqScatter = toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getscattercomparison", paramList, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var resData = res.data.Data;
        if (resData != null) {
            if (resData.length > 0) {
                dataPowerCurves = resData
            }
        }
    });
    $.when(reqScatter).done(function() {
        var dtSeries = new Array();
        if (dataPowerCurves != null) {
            if (dataPowerCurves.length > 0) {
                dtSeries = dtLine.concat(dataPowerCurves);
            }
        } else {
            dtSeries = dtLine;
        }

        $('#cm-chart').html("");
        $("#cm-chart").kendoChart({
            theme: "flat",
            // renderAs: "canvas",
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            // title: {
            //     text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
            //     visible: false,
            //     font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            // },
            // legend: {
            //     visible: false,
            //     position: "bottom"
            // },
            chartArea: {
                height: 400,
                background: 'transparent',
            },
            legend: {
                position: "bottom",
                visible: true,
                align: "start",
                offsetX : 55,
                labels: {
                    margin: {
                        right : 0,
                    },
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
            },
            seriesDefaults: {
                type: "scatterLine",
                style: "smooth",
            },
            series: dtSeries,
            categoryAxis: {
                labels: {
                    step: 1
                }
            },
            valueAxis: [{
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555",
                    visible: true,
                },
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
                max: 25
            },
            yAxis: {
                title: {
                    text: "Generation (KW)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555"
                },
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                axisCrossingValue: -5,
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N1",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
            },
            pannable: true,
            zoomable: true
        });

        app.loading(false);
    });
}

cm.setProjectTurbine = function(){
	var projects = fa.rawproject();
	var turbines = fa.rawturbine();

	cm.rawproject(projects);
    cm.rawturbine(turbines);
    var sortedTurbine = cm.rawturbine().sort(function(a, b){
        var a1= a.Turbine.toLowerCase(), b1= b.Turbine.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    var sortedProject = cm.rawproject().sort(function(a, b){
        var a1= a.Value.toLowerCase(), b1= b.Value.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    cm.rawturbine(sortedTurbine);
    cm.rawproject(sortedProject)
};

cm.initElementEvents = function() {
    $('#cm-btn-refresh').on('click', function() {
        setTimeout(function() {
            cm.initChart();
        }, 300);
    });
    $('#sScater').on('click', function() {
        var sScater = $('#sScater').prop('checked');
        cm.sScater(sScater);
        cm.initChart();
    });

}
$(document).ready(function() {
    $("#cm-project").kendoMultiSelect({
        dataSource: cm.projectList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {
            if ($("#cm-project").data("kendoMultiSelect").value().length > 0) {
                cm.populateTurbine();
            } else {
                $('#cm-project').data('kendoMultiSelect').value([fa.project]);
                $("#cm-project").data("kendoMultiSelect").trigger("change");
            }
        }, 
        suggest: true
    });   
});


