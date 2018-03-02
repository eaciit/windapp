'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

vm.currentMenu('Operational Power Curve');
vm.currentTitle('Operational Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Operational Power Curve',
    href: viewModel.appName + 'page/analyticpcscatteroperational'
}]);


page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "rotor", "text": "Rotor RPM" },
    { "value": "generatorrpm", "text": "Generator RPM" },
    { "value": "windspeed", "text": "Wind Speed" },
]);

page.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
page.turbineList = ko.observableArray([]);
page.turbineList2 = ko.observableArray([]);
page.projectList = ko.observableArray([]);
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.project = ko.observable();
page.sScater = ko.observable(false);

page.rawturbine = ko.observableArray([]);
page.rawproject = ko.observableArray([]);

var lastPeriod = "";
var turbineval = [];
page.IDList = [];
page.countList = 0;
page.LastFilter;
page.TableName;
page.FieldList;
page.ContentFilter;

page.getPDF = function(selector){
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

page.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var availDateAll = res.data;        
        var minDate  = (kendo.toString(moment.utc(availDateAll["Tejuva"]["ScadaData"][0]).format('DD-MMM-YYYY')));
        var maxDate = (kendo.toString(moment.utc(availDateAll["Tejuva"]["ScadaData"][1]).format('DD-MMM-YYYY')));

        $('#availabledatestartscada').html("from: <strong>" + minDate + "</strong> ");
        $('#availabledateendscada').html("until: <strong>" + maxDate + "</strong>");

        page.generateElementFilter(null, "default1");
        page.generateElementFilter(null, "default2");
    });
}

page.populateTurbine = function () {
    if (page.rawturbine().length == 0) {
        page.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];

        $.each($("#projectList").data("kendoMultiSelect").value(), function(i, project){
            $.each(page.rawturbine(), function(key, val){
                if(project == val.Project){
                    var data = {};
                    data.value = val.Value + "<>" + val.Project;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }
            });
        });
        page.turbineList(datavalue);
    }
    page.IDList.forEach(function(id) {
        $('#turbineList-'+id).data('kendoDropDownList').setDataSource(new kendo.data.DataSource({ data: page.turbineList() }));
        $('#turbineList-'+id).data('kendoDropDownList').select(0);
    });
};

page.populateProject = function () {
    if (page.rawproject().length == 0) {
        page.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];        
        $.each(page.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Value;
            datavalue.push(data);
        });
        page.projectList(datavalue);
    }
    $("#projectList").data("kendoMultiSelect").setDataSource(page.projectList());
    $('#projectList').data('kendoMultiSelect').value([page.projectList()[0].value]);
    $("#projectList").data("kendoMultiSelect").trigger("change");
};

page.getRandomId = function () {
    return page.randomNumber() + page.randomNumber() + page.randomNumber() + page.randomNumber();
}

page.randomNumber = function () {
    return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
}

page.generateElementFilter = function (id_element, source) {
    page.countList++;
    var id = (id_element == null ? page.getRandomId() : id_element);
    var isDefault = false;
    if(source.indexOf("default") >= 0) {
        isDefault = true;
    }
    if(page.IDList.length == 5) {
        swal('Warning', 'You can only add 5 filters', 'warning');
        return;
    }
    page.IDList.push(id);
    var isLast = false;
    if(page.IDList.length == 5) {
        isLast = true;
    }

    var formFilter ='<div class="row dynamic-filter" id="filter-form-'+ id + '" data-count="'+ page.countList +'">' +
                        '<div class="col-md-3 no-padding">' +
                            '<select class="turbine-list" id="turbineList-' + id + '" name="table" multiple="multiple"></select>' +
                        '</div>' +
                        '<div class="col-md-9 no-padding">' +
                            '<div class="input-group mb-3">'+
                                '<select class="period-list" id="periodList-' + id + '" name="table"></select>' +
                                '<span class="show_hide custom-period">' +
                                    '<input type="text" id="dateStart-' + id + '"/>' +
                                    '<label>&nbsp;&nbsp;&nbsp;to&nbsp;&nbsp;&nbsp;</label>' +
                                    '<input type="text" id="dateEnd-' + id + '"/>' +
                                '</span>' +
                            '</div>'+
                        '</div>' +
                        '<button class="btn btn-sm btn-danger tooltipster tooltipstered remove-btn" onClick="page.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (isDefault ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
                    '</div>';
    var versusFilter = '<div class="versus-wrapper" data-count="'+ page.countList +'"><div class="versus">vs</div></div>';

    setTimeout(function () {
        $(".filter-part").append(formFilter);
        $(".filter-part").append(versusFilter);

        $("#turbineList-" + id).kendoDropDownList({
            dataValueField: 'value',
            dataTextField: 'label',
            suggest: true,
            dataSource: page.turbineList(),
        });     

        $("#periodList-" + id).kendoDropDownList({
            dataSource: page.periodList(),
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { 
                page.showHidePeriod(id) 
            }
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
        page.InitDefaultValue(id);
        page.checkElementLast();
    }, 500);
}

page.removeFilter = function (id) {
    page.countList--;
    $("#filter-form-" + id).remove();
    var tempList = [];
    page.IDList.forEach(function(val){
        if (val !== id) {
            tempList.push(val);
        }
    });
    page.IDList = tempList;
    page.checkElementLast();
}

page.checkElementLast = function(){
    var elms = $('.dynamic-filter');
    $.each(elms, function(i, e){
        if(!$(e).hasClass('dynamic-filter-last')) {
            $(e).addClass('dynamic-filter-last');
        }
        var dataCount = parseInt($(e).attr('data-count'));
        if(dataCount < page.countList) {
            $(e).removeClass('dynamic-filter-last');
        }
        var turbineElm = $(e).find('select.turbine-list');
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
        if(dataCount < page.countList) {
            $(e).removeClass('versus-last');
        }
    });
    setTimeout(function () {
        page.LoadData();                           
    }, 500);
}

page.showHidePeriod = function (idx) {
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
}

page.InitDefaultValue = function (id) {
    $("#periodList-" + id).data("kendoDropDownList").value("custom");
    $("#periodList-" + id).data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
}

page.LoadData = function() {
    page.getPowerCurveScatter();
}

page.refreshChart = function() {
    page.LoadData();
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();

    var mostDateStart;
    var mostDateEnd;
    var projectList = $("#projectList").data("kendoMultiSelect").value();
    var turbineList = [];
    var details = [];

    page.IDList.forEach(function(id){
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
            ScatterType  : page.scatterType,
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
        page.LastFilter = res.data.LastFilter;
        page.FieldList = res.data.FieldList;
        page.TableName = res.data.TableName;
        page.ContentFilter = res.data.ContentFilter;
        
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
        switch(page.scatterType) {
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

        $('#scatterChart').html("");
        $("#scatterChart").kendoChart({
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

page.setProjectTurbine = function(projects, turbines){
    page.rawproject(projects);
    page.rawturbine(turbines);
    var sortedTurbine = page.rawturbine().sort(function(a, b){
        var a1= a.Turbine.toLowerCase(), b1= b.Turbine.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    var sortedProject = page.rawproject().sort(function(a, b){
        var a1= a.Value.toLowerCase(), b1= b.Value.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    page.rawturbine(sortedTurbine);
    page.rawproject(sortedProject);
};

page.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: page.LastFilter,
        FieldList: page.FieldList,
        Tablename: page.TableName,
        TypeExcel: namaFile,
        ContentFilter: page.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: page.LastFilterDetails,
            FieldList: page.FieldListDetails,
            Tablename: page.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: page.ContentFilterDetails,
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


$(document).ready(function() {

    $('#btnRefresh').on('click', function() {
        setTimeout(function() {
            page.LoadData();
        }, 300);
    });

    $("#projectList").kendoMultiSelect({
        dataSource: page.projectList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {
            if ($("#projectList").data("kendoMultiSelect").value().length > 0) {
                page.populateTurbine();
            } else {
                $('#projectList').data('kendoMultiSelect').value([page.projectList()[0].value]);
                $("#projectList").data("kendoMultiSelect").trigger("change");
            }
        }, 
        suggest: true
    });
    page.populateProject();

    setTimeout(function() {
        page.getAvailDate();
    }, 700);
});