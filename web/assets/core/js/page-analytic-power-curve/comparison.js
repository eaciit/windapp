'use strict';


viewModel.PCComparison = new Object();
var pc = viewModel.PCComparison;


vm.currentMenu('Comparison');
vm.currentTitle('Comparison');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Comparison',
    href: viewModel.appName + 'page/analyticpccomparison'
}]);


pc.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);

pc.turbineList = ko.observableArray([]);
pc.projectList = ko.observableArray([]);
pc.dateStart = ko.observable();
pc.dateEnd = ko.observable();
pc.project = ko.observable();
pc.sScater = ko.observable(false);

pc.rawturbine = ko.observableArray([]);
pc.rawproject = ko.observableArray([]);
pc.IDList = [];
pc.countList = 0;
pc.LastFilter;
pc.TableName;
pc.FieldList;
pc.ContentFilter;

pc.getPDF = function(selector){
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
      kendo.drawing.pdf.saveAs(group,  pc.project()+"PC Comparison.pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

pc.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var availDateAll = res.data;        
        var minDate  = (kendo.toString(moment.utc(availDateAll["Tejuva"]["ScadaData"][0]).format('DD-MMM-YYYY')));
        var maxDate = (kendo.toString(moment.utc(availDateAll["Tejuva"]["ScadaData"][1]).format('DD-MMM-YYYY')));

        $('#availabledatestartscada').html("from: <strong>" + minDate + "</strong> ");
        $('#availabledateendscada').html("until: <strong>" + maxDate + "</strong>");

        pc.generateElementFilter(null, "default1");
        pc.generateElementFilter(null, "default2");
    });
}
pc.populateTurbine = function () {
    if (pc.rawturbine().length == 0) {
        pc.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];

        $.each($("#projectList").data("kendoMultiSelect").value(), function(i, project){
            $.each(pc.rawturbine(), function(key, val){
                if(project == val.Project){
                    var data = {};
                    data.value = val.Value + "<>" + val.Project;
                    data.label = val.Turbine;
                    datavalue.push(data);
                }
            });
        });
        pc.turbineList(datavalue);
    }
    pc.IDList.forEach(function(id) {
        $('#turbineList-'+id).data('kendoDropDownList').setDataSource(new kendo.data.DataSource({ data: pc.turbineList() }));
        $('#turbineList-'+id).data('kendoDropDownList').select(0);
    });
};

pc.populateProject = function () {
    if (pc.rawproject().length == 0) {
        pc.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];        
        $.each(pc.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Value;
            datavalue.push(data);
        });
        pc.projectList(datavalue);
    }
    $("#projectList").data("kendoMultiSelect").setDataSource(pc.projectList());
    $('#projectList').data('kendoMultiSelect').value([pc.projectList()[0].value]);
    $("#projectList").data("kendoMultiSelect").trigger("change");
};

pc.getRandomId = function () {
    return pc.randomNumber() + pc.randomNumber() + pc.randomNumber() + pc.randomNumber();
}

pc.randomNumber = function () {
    return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
}

pc.generateElementFilter = function (id_element, source) {
    pc.countList++;
    var id = (id_element == null ? pc.getRandomId() : id_element);
    var isDefault = false;
    if(source.indexOf("default") >= 0) {
        isDefault = true;
    }
    if(pc.IDList.length == 5) {
        return;
    }
    pc.IDList.push(id);
    var isLast = false;
    if(pc.IDList.length == 5) {
        isLast = true;
    }

    // var formFilterOri = '<div class="col-md-12 dynamic-filter" id="filter-form-'+ id + '">' +
    //                         '<div class="row mgb10">' +
    //                             '<div class="col-md-2 no-padding">' +
    //                                 '<label class="control-label">Turbine</label>' +
    //                             '</div>' +
    //                             '<div class="col-md-9 no-padding">' +
    //                                 '<select class="turbine-list" id="turbineList-' + id + '" name="table" multiple="multiple"></select>' +
    //                             '</div>' +
    //                             '<div class="col-md-1 no-padding">' +
    //                                 '<button class="btn btn-sm btn-danger tooltipster tooltipstered remove-btn" onClick="pc.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (isDefault ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
    //                             '</div>' +
    //                         '</div>' +
    //                         '<div class="row mgb10">' +
    //                             '<div class="col-md-2 no-padding">' +
    //                                 '<label class="control-label lbl-period">Period</label>' +
    //                             '</div>' +
    //                             '<div class="col-md-10 no-padding">' +
    //                                 '<select class="period-list" id="periodList-' + id + '" name="table"></select>' +
    //                                 '<span class="show_hide custom-period">' +
    //                                     '<input type="text" id="dateStart-' + id + '"/>' +
    //                                     '<label class="control-label label-to">&nbsp;&nbsp;to</label>' +
    //                                     '<div class="period-nwline">&nbsp;</div>' +
    //                                     '<input type="text" id="dateEnd-' + id + '"/>' +
    //                                 '</span>' +
    //                             '</div>' +
    //                         '</div>' +
    //                         '<div class="row">' +
    //                             '<hr class="horizontal-line" style="display:' + (isLast ? 'none' : 'inherit') + '">'+
    //                         '</div>' +
    //                     '</div>';
    
var formFilter =    '<div class="row dynamic-filter" id="filter-form-'+ id + '" data-count="'+ pc.countList +'">' +
                        '<div class="mgb10">' +
                            '<div class="col-md-3 no-padding">' +
                                '<select class="turbine-list" id="turbineList-' + id + '" name="table" multiple="multiple"></select>' +
                            '</div>' +
                            '<div class="col-md-9 no-padding">' +
                                '<select class="period-list" id="periodList-' + id + '" name="table"></select>' +
                                '<span class="show_hide custom-period">' +
                                    '<input type="text" id="dateStart-' + id + '"/>' +
                                    '<label>&nbsp;&nbsp;&nbsp;to&nbsp;&nbsp;&nbsp;</label>' +
                                    '<input type="text" id="dateEnd-' + id + '"/>' +
                                '</span>' +
                            '</div>' +
                            '<button class="btn btn-sm btn-danger tooltipster tooltipstered remove-btn" onClick="pc.removeFilter(\'' + id + '\')" id="btn-remove-' + id + '" title="Remove Filter" style="display:' + (isDefault ? 'none' : 'inline') + '"><i class="fa fa-times"></i></button>' +
                        '</div>'
                    '</div>';
var versusFilter = '<div class="versus-wrapper" data-count="'+ pc.countList +'"><div class="versus">vs</div></div>';

    setTimeout(function () {
        $(".filter-part").append(formFilter);
        $(".filter-part").append(versusFilter);

        $("#turbineList-" + id).kendoDropDownList({
            dataValueField: 'value',
            dataTextField: 'label',
            suggest: true,
            dataSource: pc.turbineList(),
        });     

        $("#periodList-" + id).kendoDropDownList({
            dataSource: pc.periodList(),
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { 
                pc.showHidePeriod(id) 
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
        // setTimeout(function () {
        //     if (source !== "default2") {
        //         $('#turbineList-'+id).data('kendoDropDownList').select(pc.countList - 1);
        //     } else {
        //         $('#turbineList-'+id).data('kendoDropDownList').select(1);
        //     }
        // }, 100);
        pc.InitDefaultValue(id);

        if(source == "default2"){
            // setTimeout(function () {
            //     pc.initChart();                           
            // }, 500);
        }
        pc.checkElementLast();
    }, 500);
}

pc.removeFilter = function (id) {
    pc.countList--;
    $("#filter-form-" + id).remove();
    var tempList = [];
    pc.IDList.forEach(function(val){
        if (val !== id) {
            tempList.push(val);
        }
    });
    pc.IDList = tempList;
    pc.checkElementLast();
}

pc.checkElementLast = function(){
    var elms = $('.dynamic-filter');
    $.each(elms, function(i, e){
        if(!$(e).hasClass('dynamic-filter-last')) {
            $(e).addClass('dynamic-filter-last');
        }
        var dataCount = parseInt($(e).attr('data-count'));
        if(dataCount < pc.countList) {
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
        if(dataCount < pc.countList) {
            $(e).removeClass('versus-last');
        }
    });
    setTimeout(function () {
        pc.initChart();                           
    }, 500);
}

pc.showHidePeriod = function (idx) {
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

pc.InitDefaultValue = function (id) {
    $("#periodList-" + id).data("kendoDropDownList").value("custom");
    $("#periodList-" + id).data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(),maxDateData.getDate()-7,0,0,0,0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));

    $('#dateStart-' + id).data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd-' + id).data('kendoDatePicker').value(lastEndDate);
}

pc.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: pc.LastFilter,
        FieldList: pc.FieldList,
        Tablename: pc.TableName,
        TypeExcel: namaFile,
        ContentFilter: pc.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: pc.LastFilterDetails,
            FieldList: pc.FieldListDetails,
            Tablename: pc.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: pc.ContentFilterDetails,
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

pc.initChart = function() {
    app.loading(true);

    var link = "analyticpowercurve/getlistpowercurvecomparison";
    var mostDateStart;
    var mostDateEnd;
    var projectList = $("#projectList").data("kendoMultiSelect").value();
    var turbineList = [];
    var details = [];

    pc.IDList.forEach(function(id){
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
        pc.LastFilter = res.data.LastFilter;
        pc.FieldList = res.data.FieldList;
        pc.TableName = res.data.TableName;
        pc.ContentFilter = res.data.ContentFilter;
        
        $('#chartPCcomparison').html("");
        $("#chartPCcomparison").kendoChart({
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
        if (pc.sScater()) {
            pc.getScatter(details, dataTurbine, projectList.length);
        }
        $("#chartPCcomparison").data("kendoChart").refresh();                
    });
}

pc.getScatter = function(paramLine, dtLine, startColorIdx) {
    var turbineList = [];
    var kolor = [];
    var idx;
    app.loading(true);
    var paramList = [];
    paramLine.forEach(function(data){
        turbineList = [];
        kolor = [];
        kolor.push(dtLine[startColorIdx].color);
        turbineList.push(data.Turbine);
        var dateStart = data.DateStart;
        var dateEnd = data.DateEnd;
        var param = {
            period: data.Period,
            dateStart: dateStart,
            dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
            turbine: turbineList,
            project: data.Project,
            Color: kolor,
            isDeviation: false,
            deviationVal: "",
            DeviationOpr: "0",
            IsDownTime: false,
            ViewSession: "",
            isPower0: false,
        };
        paramList.push(param);
        startColorIdx++;
    });
    var dataPowerCurves = [];
    var reqScatter1 = toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurve", paramList[0], function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dataPowerCurves1 = res.data.Data;
        if (dataPowerCurves1 != null) {
            if (dataPowerCurves1.length > 0) {
                dataPowerCurves.push(dataPowerCurves1[0]);
            }
        }
    });
    var reqScatter2 = toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurve", paramList[1], function(res) {
        if (!app.isFine(res)) {
            return;
        }
        var dataPowerCurves2 = res.data.Data;
        if (dataPowerCurves2 != null) {
            if (dataPowerCurves2.length > 0) {
                dataPowerCurves.push(dataPowerCurves2[0]);
            }
        }
    });
    $.when(reqScatter1, reqScatter2).done(function() {
        var dtSeries = new Array();
        if (dataPowerCurves != null) {
            if (dataPowerCurves.length > 0) {
                dtSeries = dtLine.concat(dataPowerCurves);
            }
        } else {
            dtSeries = dtLine;
        }

        $('#chartPCcomparison').html("");
        $("#chartPCcomparison").kendoChart({
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

        var chart = $("#chartPCcomparison").data("kendoChart");
        var series = chart.options.series;
        for (var i = 0; i < series.length; i++) {
            if(i >= series.length-2) {
                series[i].visibleInLegend = false;
            }
        };
        chart.redraw();

        app.loading(false);
    });
}

pc.setProjectTurbine = function(projects, turbines){
	pc.rawproject(projects);
    pc.rawturbine(turbines);
    var sortedTurbine = pc.rawturbine().sort(function(a, b){
        var a1= a.Turbine.toLowerCase(), b1= b.Turbine.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    var sortedProject = pc.rawproject().sort(function(a, b){
        var a1= a.Value.toLowerCase(), b1= b.Value.toLowerCase();
        if(a1== b1) return 0;
        return a1> b1? 1: -1;
    });
    pc.rawturbine(sortedTurbine);
    pc.rawproject(sortedProject)
};

$(document).ready(function () {
    
    $('#btnRefresh').on('click', function() {
        setTimeout(function() {
            pc.initChart();
        }, 300);
    });
    $('#sScater').on('click', function() {
        var sScater = $('#sScater').prop('checked');
        pc.sScater(sScater);
        pc.initChart();
    });
    $("#projectList").kendoMultiSelect({
        dataSource: pc.projectList(), 
        dataValueField: 'value', 
        dataTextField: 'text', 
        change: function() {
            if ($("#projectList").data("kendoMultiSelect").value().length > 0) {
                pc.populateTurbine();
            } else {
                $('#projectList').data('kendoMultiSelect').value([pc.projectList()[0].value]);
                $("#projectList").data("kendoMultiSelect").trigger("change");
            }
        }, 
        suggest: true
    });
    pc.populateProject();

    app.loading(true);
    setTimeout(function() {
        pc.getAvailDate();
    }, 700);
});