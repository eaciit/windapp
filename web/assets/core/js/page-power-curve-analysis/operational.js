'use strict';
var op = {};
op.loadFirstTime = ko.observable(true);
op.scatterType = ko.observable('');
op.scatterList = ko.observableArray([
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "rotor", "text": "Rotor RPM" },
    { "value": "generatorrpm", "text": "Generator RPM" },
    { "value": "windspeed", "text": "Wind Speed" },
]);

op.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
op.turbineList = ko.observableArray([]);
op.turbineList2 = ko.observableArray([]);
op.projectList = ko.observableArray([]);
op.dateStart = ko.observable();
op.dateEnd = ko.observable();
op.project = ko.observable();
op.sScater = ko.observable(false);

op.rawturbine = ko.observableArray([]);
op.rawproject = ko.observableArray([]);

var lastPeriod = "";
var turbineval = [];
op.IDList = [];
op.countList = 0;
op.LastFilter;
op.TableName;
op.FieldList;
op.ContentFilter;


op.reset = function(){
    op.loadFirstTime(true);
}
op.refresh = function() {
   if(this.loadFirstTime()) {
        op.countList = 0;
        this.initElementEvents();
	    this.loadFirstTime(false);
	    this.setProjectTurbine();
		this.populateProject();
		this.firstLoad();
		// this.LoadData();
   	}
}

op.internalRefresh = function(reloadData) {
    if(reloadData==undefined) {
        reloadData = true;
    }

    if(reloadData) {
       // this.LoadData();
    }
}


op.firstLoad = function(){
    op.IDList = [];

	$(".op-filter-part").html("");

    this.generateElementFilter(null, "default1");
    this.generateElementFilter(null, "default2");
}

op.getPDF = function(selector){
    app.loading(true);

    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            scale: 0.5,
            margin: {
                left   : "5mm",
                top    : "5mm",
                right  : "10mm",
                bottom : "5mm"
            },
        });
      kendo.drawing.pdf.saveAs(group, "Operational Power Curve.pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

op.populateTurbine = function () {
    if (op.rawturbine().length == 0) {
        op.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];

        $.each($("#op-project").data("kendoMultiSelect").value(), function(i, project){
            $.each(op.rawturbine(), function(key, val){
                if(project == val.Project){
                    var data = {};
                    data.value = val.Value + "<>" + val.Project;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }
            });
        });
        op.turbineList(datavalue);
    }

    var a = 0;

    op.IDList.forEach(function(id) {
        $('#op-turbine-list'+id).data('kendoDropDownList').setDataSource(new kendo.data.DataSource({ data: op.turbineList() }));
        $('#op-turbine-list'+id).data('kendoDropDownList').select(a);
        a++;
    });
};

op.populateProject = function () {
    if (op.rawproject().length == 0) {
        op.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];        
        $.each(op.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Value;
            datavalue.push(data);
        });
        op.projectList(datavalue);
    }
    $("#op-project").data("kendoMultiSelect").setDataSource(op.projectList());
    $('#op-project').data('kendoMultiSelect').value([fa.project]);
    $("#op-project").data("kendoMultiSelect").trigger("change");
};

op.getRandomId = function () {
    return op.randomNumber() + op.randomNumber() + op.randomNumber() + op.randomNumber();
}

op.randomNumber = function () {
    return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
}

op.generateElementFilter = function (id_element, source) {
    op.countList++;
    var id = (id_element == null ? op.getRandomId() : id_element);
    var isDefault = false;
    if(source.indexOf("default") >= 0) {
        isDefault = true;
    }
    if(op.IDList.length == 5) {
        swal('Warning', 'You can only add 5 filters', 'warning');
        return;
    }
    op.IDList.push(id);
    var isLast = false;
    if(op.IDList.length == 5) {
        isLast = true;
        $(".op-button-add").hide();
    }

    var formFilter ='<div class="row op-dynamic-filter" id="filter-form-'+ id + '" data-count="'+ op.countList +'">' +
                        '<div class="col-md-3 no-padding">' +
                            '<select class="op-turbine-list" id="op-turbine-list' + id + '" name="table" multiple="multiple"></select>' +
                        '</div>' +
                        '<div class="col-md-9 no-padding">' +
                            '<div class="input-group mb-3">'+
                                '<select class="period-list" id="po-periodList-' + id + '" name="table"></select>' +
                                '<span class="custom-period" id="op-show-hide-'+ id +'">' +
                                    '<input type="text" id="op-dateStart-' + id + '"/>' +
                                    '<label>&nbsp;&nbsp;&nbsp;to&nbsp;&nbsp;&nbsp;</label>' +
                                    '<input type="text" id="op-dateEnd-' + id + '"/>' +
                                '</span>' +
                            '</div>'+
                        '</div>' +
                        '<button class="btn btn-sm btn-danger tooltipster tooltipstered remove-btn" onClick="op.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (isDefault ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
                    '</div>';
    var versusFilter = '<div class="versus-wrapper" data-count="'+ op.countList +'"><div class="versus">vs</div></div>';

    setTimeout(function () {
        $(".op-filter-part").append(formFilter);
        $(".op-filter-part").append(versusFilter);

        $("#op-turbine-list" + id).kendoDropDownList({
            dataValueField: 'value',
            dataTextField: 'label',
            suggest: true,
            dataSource: op.turbineList(),
            change: function(){
            	op.refreshChart();
            }
        });     
         $('#op-turbine-list'+id).data('kendoDropDownList').select(op.countList);

        $("#po-periodList-" + id).kendoDropDownList({
            dataSource: op.periodList(),
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { 
                op.showHidePeriod(id) 
            }
        });

        $('#op-dateStart-' + id).kendoDatePicker({
            value: new Date(),
            format: 'dd-MMM-yyyy',
            min: new Date("2013-01-01"),
            max:new Date(),
            change: function(){
            	op.refreshChart();
            }
        });

        $('#op-dateEnd-' + id).kendoDatePicker({
            value: new Date(),
            format: 'dd-MMM-yyyy',
            min: new Date("2013-01-01"),
            max:new Date(),
            change: function(){
            	op.refreshChart();
            }
        });
        op.InitDefaultValue(id);
        op.checkElementLast();
    }, 500);
}

op.removeFilter = function (id) {
    op.countList--;
    $("#filter-form-" + id).remove();
    var tempList = [];
    op.IDList.forEach(function(val){
        if (val !== id) {
            tempList.push(val);
        }
    });
    op.IDList = tempList;
    op.checkElementLast();
}

op.checkElementLast = function(){
    var elms = $('.op-dynamic-filter');
    $.each(elms, function(i, e){
        if(!$(e).hasClass('op-dynamic-filter-last')) {
            $(e).addClass('op-dynamic-filter-last');
        }
        var dataCount = parseInt($(e).attr('data-count'));
        if(dataCount < op.countList) {
            $(e).removeClass('op-dynamic-filter-last');
        }
        var turbineElm = $(e).find('select.op-turbine-list');
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
        if(dataCount < op.countList) {
            $(e).removeClass('versus-last');
        }
    });
    setTimeout(function () {
        op.LoadData();                           
    }, 500);
}

op.showHidePeriod = function (idx) {
    var id = (idx == null ? 1 : idx);
    var period = $('#po-periodList-' + id).data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.getDateMax(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

    if (period == "custom") {
        $("#op-show-hide-" + id).show();
        $('#op-dateStart-' + id).data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#op-dateEnd-' + id).data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#op-dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
        $('#op-dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
    } else if (period == "monthly") {
        $('#op-dateStart-' + id).data('kendoDatePicker').setOptions({
            start: "year",
            depth: "year",
            format: "MMM yyyy"
        });
        $('#op-dateEnd-' + id).data('kendoDatePicker').setOptions({
            start: "year",
            depth: "year",
            format: "MMM yyyy"
        });

        $('#op-dateStart-' + id).data('kendoDatePicker').value(startMonthDate);
        $('#op-dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);

        $("#op-show-hide-" + id).show();
    } else if (period == "annual") {
        $("#op-show-hide-" + id).show();

        $('#op-dateStart-' + id).data('kendoDatePicker').setOptions({
            start: "decade",
            depth: "decade",
            format: "yyyy"
        });
        $('#op-dateEnd-' + id).data('kendoDatePicker').setOptions({
            start: "decade",
            depth: "decade",
            format: "yyyy"
        });

       $('#op-dateStart-' + id).data('kendoDatePicker').value(startYearDate);
       $('#op-dateEnd-' + id).data('kendoDatePicker').value(endYearDate);

        $("#op-show-hide-" + id).show();
    } else {
        if(period == 'last24hours'){
             $('#op-dateStart-' + id).data('kendoDatePicker').value(last24hours);
             $('#op-dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
        }else if(period == 'last7days'){
             $('#op-dateStart-' + id).data('kendoDatePicker').value(lastweek);
             $('#op-dateEnd-' + id).data('kendoDatePicker').value(endMonthDate);
        }
        $("#op-show-hide-" + id).hide();
    }
}

op.InitDefaultValue = function (id) {
    $("#po-periodList-" + id).data("kendoDropDownList").value("custom");
    $("#po-periodList-" + id).data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#op-dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
    $('#op-dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
}

op.LoadData = function() {
    op.getPowerCurveScatter();
}

op.refreshChart = function() {
    op.LoadData();
}

op.getPowerCurveScatter = function() {
    app.loading(true);
    op.scatterType = $("#op-scatterType").data('kendoDropDownList').value();

    var mostDateStart;
    var mostDateEnd;
    var projectList = $("#op-project").data("kendoMultiSelect").value();
    var turbineList = [];
    var details = [];

    op.IDList.forEach(function(id){
        var dateStart = $('#op-dateStart-'+id).data('kendoDatePicker').value();
            dateStart = new Date(Date.UTC(dateStart.getFullYear(), dateStart.getMonth(), dateStart.getDate(), 0, 0, 0));

        var dateEnd  = $('#op-dateEnd-'+id).data('kendoDatePicker').value();
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
        var splitTurbineVal = $("#op-turbine-list"+id).data("kendoDropDownList").value().split("<>");
        turbineList.push(splitTurbineVal[0]);

        var detail = {
            Period       : $('#po-periodList-'+id).data('kendoDropDownList').value(),
            Project      : splitTurbineVal[1],
            Turbine      : splitTurbineVal[0],
            DateStart    : dateStart,
            DateEnd      : dateEnd,
            ScatterType  : op.scatterType,
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

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatteroperational", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dtSeries = res.data.Data;
        op.LastFilter = res.data.LastFilter;
        op.FieldList = res.data.FieldList;
        op.TableName = res.data.TableName;
        op.ContentFilter = res.data.ContentFilter;
        
        var minAxisY = res.data.MinAxisY;
        var maxAxisY = res.data.MaxAxisY;
        var minAxisX = res.data.MinAxisX;
        var maxAxisX = res.data.MaxAxisX;
        var name = '';
        var title = '';
        var xAxis = {};
        var measurement = '';
        var format = 'N0'
        if(maxAxisX - minAxisX < 7) {
            format = 'N2'
        }
        switch(op.scatterType) {
            case "pitch":
                name = 'pitchAxis'
                title = 'Angle (Degree)'
                measurement = String.fromCharCode(176)
                break;
            case "rotor":
                name = "rotorAxis"
                title = "Revolutions per Minute (RPM)";
                measurement = 'rpm'
                break;
            case "generatorrpm":
                name = "generatorAxis"
                title = "Generator per Minute (RPM)";
                measurement = 'rpm'
                break;
            case "windspeed":
                name = "windspeedAxis"
                title = "Avg. Wind Speed (m/s)";
                measurement = 'm/s'
                break;
        }
        xAxis = {
            name: name,
            title: {
                text: title,
                visible: true,
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                format: format
            },
            line: {
                visible: false
            },
            axisCrossingValue: -10,
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
                    template: "#= kendo.toString(value, 'n2') # " + measurement,
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                }
            },
            // majorUnit: 0.5,
            min: minAxisX,
            max: maxAxisX,
        }
        var yAxis = {};
        yAxis = {
            name: "powerAxis",
            title: {
                text: "Generation (KW)",
                visible: true,
                font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            line: {
                visible: false
            },
            axisCrossingValue: -10,
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
                    template: "#= kendo.toString(value, 'n2') # kWh",
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },
                }
            },
            // majorUnit: 0.5,
            min: minAxisY,
            max: maxAxisY,
        }

        $('#op-chart').html("");
        $("#op-chart").kendoChart({
            theme: "flat",
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            chartArea: {
                background: "transparent",
                padding: 0,
            },
            title: {
                // text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
                 text: "Scatter Power Curves",
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: true,
                align: "center",
                offsetX : 50,
                labels: {
                    margin: {
                        right : 20
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
                },
            },
            valueAxis: [{
                labels: {
                    format: "N2",
                }
            }],
            xAxis: xAxis,
            yAxes: yAxis,
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
    });
    
}

op.setProjectTurbine = function(projects, turbines){

	var projects = fa.rawproject();
	var turbines = fa.rawturbine();

    op.rawproject(projects);
    op.rawturbine(turbines);
    var sortedTurbine = op.rawturbine().sort(function(a, b){
        var a1= a.Turbine.toLowerCase(), b1= b.Turbine.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    var sortedProject = op.rawproject().sort(function(a, b){
        var a1= a.Value.toLowerCase(), b1= b.Value.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    op.rawturbine(sortedTurbine);
    op.rawproject(sortedProject);
};

op.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = tipe;
    }

    var param = {
        Filters: op.LastFilter,
        FieldList: op.FieldList,
        Tablename: op.TableName,
        TypeExcel: namaFile,
        ContentFilter: op.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: op.LastFilterDetails,
            FieldList: op.FieldListDetails,
            Tablename: op.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: op.ContentFilterDetails,
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


op.initElementEvents = function() {
    $('#op-btn-refresh').on('click', function() {
        setTimeout(function() {
            op.LoadData();
        }, 300);
    });
}


$(document).ready(function() {
    $("#op-project").kendoMultiSelect({
        dataSource: op.projectList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {
            if ($("#op-project").data("kendoMultiSelect").value().length > 0) {
                op.populateTurbine();
            } else {
                $('#op-project').data('kendoMultiSelect').value([fa.project]);
                $("#op-project").data("kendoMultiSelect").trigger("change");
            }
        }, 
        suggest: true
    });   
});

