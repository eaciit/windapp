'use strict';
var sf = {};


sf.loadFirstTime = ko.observable(true);

sf.colorPalette = ko.observable("websafe");

sf.turbineList = ko.observableArray([]);
sf.turbine = ko.observable();
sf.project = ko.observable();
sf.scatterType = ko.observable('');
sf.scatterList = ko.observableArray([
    { "value": "deviation", "text": "Nacelle Deviation" },
    { "value": "temp", "text": "Nacelle Temperature" },
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "ambient", "text": "Ambient Temperature" },
    { "value": "windspeed_dev", "text": "Wind Speed Std. Dev." },
    { "value": "windspeed_ti", "text": "TI Wind Speed" },
]);

sf.lessSelectedColour = ko.observable("#4589b0");
sf.greaterSelectedColour = ko.observable("#e4cc37");
sf.markerStyleList = ko.observableArray([
    {value:"circle",text:"Circle"},
    {value:"square",text:"Square"},
    {value:"triangle",text:"Triangle"},
    {value:"cross",text:"Cross"}]);

sf.lessValue = ko.observable(20);
sf.greaterValue= ko.observable(20);
sf.lessSelectedMarker = ko.observable("circle");
sf.greaterSelectedMarker = ko.observable("circle");
sf.dtSeries = ko.observableArray([]);
sf.dateStart = ko.observable();
sf.dateEnd = ko.observable();

sf.LastFilter;
sf.TableName;
sf.FieldList;
sf.ContentFilter;

sf.reset = function(){
	this.getTurbineList();
    this.loadFirstTime(true);
}
sf.refresh = function() {
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

sf.internalRefresh = function(reloadData) {
    if(reloadData==undefined) {
        reloadData = true;
    }

    this.buildFilterInfo();
    if(reloadData) {
        this.LoadData();
    }
}



sf.getTurbineList = function(){
    var turbineList = $("#turbineList").val();

    sf.turbineList(turbineList);
    sf.turbine(turbineList[0]);
}

sf.getPDF = function(selector){
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
      kendo.drawing.pdf.saveAs(group, project+"PCScatterWithFilter"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

sf.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: sf.LastFilter,
        FieldList: sf.FieldList,
        Tablename: sf.TableName,
        TypeExcel: namaFile,
        ContentFilter: sf.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: sf.LastFilterDetails,
            FieldList: sf.FieldListDetails,
            Tablename: sf.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: sf.ContentFilterDetails,
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


sf.LoadData = function() {
    sf.getPowerCurveScatter();
}

sf.refreshChart = function() {
	sf.turbine($("#sf-turbine").data("kendoDropDownList").value());
    sf.scatterType = $("#sf-scatter-type").data('kendoDropDownList').value();

    switch(sf.scatterType) {
        case 'pitch':
            $("#txtLessVal").val(0);
            $('#txtGreaterVal').val(0);
            break;
        case 'ambient':
            $("#txtLessVal").val(20);
            $('#txtGreaterVal').val(20);
            break;
        case 'temp':
            $("#txtLessVal").val(35);
            $('#txtGreaterVal').val(35);
            break;
        case 'deviation':
            $("#txtLessVal").val(5);
            $('#txtGreaterVal').val(5);
            break;
        case 'windspeed_dev':
            $("#txtLessVal").val(1);
            $('#txtGreaterVal').val(1);
            break;
        case 'windspeed_ti':
            $("#txtLessVal").val(0.2);
            $('#txtGreaterVal').val(0.2);
            break;
        default:
            $("#txtLessVal").val(0);
            $('#txtGreaterVal').val(0);
            break;
    }

    sf.LoadData();
}

sf.getPowerCurveScatter = function() {
    app.loading(true);
    sf.scatterType = $("#sf-scatter-type").data('kendoDropDownList').value();

    var lessValue = $("#txtLessVal").val();
    var greaterValue = $('#txtGreaterVal').val();

    var lessColor = $("#lessColor").data("kendoColorPicker").value();
    var greaterColor = $("#greaterColor").data("kendoColorPicker").value();
    var lessMarker = $("#lessMarker").data("kendoDropDownList").value();
    var greaterMarker = $("#greaterMarker").data("kendoDropDownList").value();


    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));   

    var param = {
        engine: fa.engine,
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: sf.turbine(),
        project: $('#projectList').data('kendoDropDownList').value(), // temporary changed
        scatterType: sf.scatterType,
        lessValue: parseFloat(lessValue),
        greaterValue: parseFloat(greaterValue),
        lessColor: lessColor,
        greaterColor: greaterColor,
        lessMarker: lessMarker, 
        greaterMarker: greaterMarker,
        
    };
    
    
    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatteranalysis", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        sf.LastFilter = res.data.LastFilter;
        sf.FieldList = res.data.FieldList;
        sf.TableName = res.data.TableName;
        sf.ContentFilter = res.data.ContentFilter;

        var result = res.data.Data;

        sf.dtSeries(result);
        sf.createChart(sf.dtSeries());

    
        app.loading(false);
    });
}


sf.changeView=function(param){

    var lessColor = $("#lessColor").data("kendoColorPicker").value();
    var greaterColor = $("#greaterColor").data("kendoColorPicker").value();
    var lessMarker = $("#lessMarker").data("kendoDropDownList").value();
    var greaterMarker = $("#greaterMarker").data("kendoDropDownList").value();

    $.each(sf.dtSeries(), function(index, value){
        if(value.name.indexOf(param) !== -1){
            sf.dtSeries()[index].color = (param == ">" ? greaterColor : lessColor);
            sf.dtSeries()[index].markers = {
                size : 2,
                type : (param == ">" ? greaterMarker : lessMarker),
                background : (param == ">" ? greaterColor : lessColor),
            }
        }
    });

    sf.createChart(sf.dtSeries());
}

sf.createChart = function(dtSeries){
        $('#sf-chart').html("");
        $("#sf-chart").kendoChart({
            theme: "flat",
            pdf: {
              fileName: "ScatterWithFilter.pdf",
            },
            title: {
                text: "Scatter with Filter | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
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
            yAxis: {
                name: "powerAxis",
                title: {
                    text: "Generation (kW)",
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
}
