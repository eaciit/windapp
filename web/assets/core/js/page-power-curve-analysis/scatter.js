'use strict';
var sc = {};
sc.loadFirstTime = ko.observable(true);

sc.plotWith = ko.observable();
sc.scatterList = ko.observableArray([]);
sc.turbineList = ko.observableArray([]);

sc.project = ko.observable();
sc.turbine = ko.observable();

sc.dateStart = ko.observable();
sc.dateEnd = ko.observable();
sc.LastFilter;
sc.TableName;
sc.FieldList;
sc.ContentFilter;

sc.reset = function(){
	this.getTurbineList();
    this.loadFirstTime(true);
}
sc.refresh = function() {
	if(this.loadFirstTime()) {
		this.getTurbineList();
	    this.loadFirstTime(false);
	    var project = $('#projectList').data("kendoDropDownList").value();
	    var dateStart = $('#dateStart').data('kendoDatePicker').value();
	    var dateEnd = $('#dateEnd').data('kendoDatePicker').value(); 

	    this.project(project);
	    this.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
	    this.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));
	    this.LoadData();
   	}
}

sc.internalRefresh = function(reloadData) {
    if(reloadData==undefined) {
        reloadData = true;
    }

    this.buildFilterInfo();
    if(reloadData) {
        this.LoadData();
    }
}

sc.getPowerCurveScatterFieldList = toolkit.ajaxPostDeffered(viewModel.appName + "analyticpowercurve/getpcscatterfieldlist", {}, function(res) {
    if (!app.isFine(res)) {
        return;
    }

    var data = res.data;
    if(data !== null){
        sc.scatterList(data);
        $("#sc-scatter-type").data("kendoDropDownList").select(0);
    }   
});

sc.getTurbineList = function(){
    var turbineList = $("#turbineList").val();

    sc.turbineList(turbineList);
    sc.turbine(turbineList[0]);
}

sc.getPDF = function(selector){
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    kendo.drawing.drawDOM($(selector)).then(function(group){
        group.options.set("pdf", {
            paperSize: "auto",
            margin: {
                left   : "5mm",
                top    : "5mm",
                right  : "5mm",
                bottom : "5mm"
            },
        });
      kendo.drawing.pdf.saveAs(group, project+"PCScatter"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}


sc.LoadData = function() {
    sc.getPowerCurveScatter();
}

sc.setAxis = function(name, title) {
    var result = {
        name: name,
        title: {
            text: title,
            font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            color: "#585555"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        labels: {
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
    }
    if(name == "powerAxis") {
        result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # kW";
        result.crosshair.tooltip.padding = {left:35};
    } else {
        switch(sc.scatterType) {
            case "temp":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176) + "C";
                break;
            case "deviation":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176);
                break;
            case "pitch":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176);
                break;
            case "ambient":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + String.fromCharCode(176) + "C";
                break;
            case "windspeed_dev":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # " + "m/s";
                break;
            case "windspeed_ti":
                result.crosshair.tooltip.template = "#= kendo.toString(value, 'n2') # ";
                break;
        }
    }
    return result
}

sc.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: sc.LastFilter,
        FieldList: sc.FieldList,
        Tablename: sc.TableName,
        TypeExcel: namaFile,
        ContentFilter: sc.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: sc.LastFilterDetails,
            FieldList: sc.FieldListDetails,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               
            Tablename: sc.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: sc.ContentFilterDetails,
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

    return false;
}

sc.refreshChart = function() {
	sc.turbine($("#sc-turbine").data("kendoDropDownList").value());
    sc.LoadData();
}

sc.getPowerCurveScatter = function() {
    app.loading(true);
    sc.plotWith = $.grep(sc.scatterList(), function(e){ return e.Name == $("#sc-scatter-type").data("kendoDropDownList").value(); })[0];

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));   

    var param = {
        engine: fa.engine,
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: sc.turbine(),
        project: $('#projectList').data("kendoDropDownList").value(),
        plotWith: sc.plotWith,
    };

    // console.log(plotWith.Text)

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurvescatterrev", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        sc.LastFilter = res.data.LastFilter;
        sc.FieldList = res.data.FieldList;
        sc.TableName = res.data.TableName;
        sc.ContentFilter = res.data.ContentFilter;

        $("#turbineName").html($("#turbineList option:selected").text());
        
        var dtSeries = res.data.Data;
        var yAxes = [];
        var yAxis = sc.setAxis("powerAxis", "Generation (KW)");
        yAxes.push(yAxis);
        var axis = sc.setAxis("PlotWith", sc.plotWith.Text);
        yAxes.push(axis);

        $('#sc-chart').html("");
        $("#sc-chart").kendoChart({
            theme: "flat",
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            title: {
                text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            },
            seriesDefaults: {
                type: "scatterLine",
                style: "smooth",
            },
            series: dtSeries,
            categoryAxis: {
                labels: {
                    step: 1,
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
            },
            valueAxis: [{
                labels: {
                    format: "N2",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    color: "#585555"
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                axisCrossingValues: [0, 30],
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        template: "#= kendo.toString(value, 'n2') # m/s",
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
            yAxes: yAxes,
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


