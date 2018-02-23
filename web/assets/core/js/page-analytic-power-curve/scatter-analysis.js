'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;
page.colorPalette = ko.observable("websafe");
page.lessSelectedColour = ko.observable("#4589b0");
page.greaterSelectedColour = ko.observable("#e4cc37");
page.markerStyleList = ko.observableArray([
    {value:"circle",text:"Circle"},
    {value:"square",text:"Square"},
    {value:"triangle",text:"Triangle"},
    {value:"cross",text:"Cross"}]);

page.lessValue = ko.observable(20);
page.greaterValue= ko.observable(20);
page.lessSelectedMarker = ko.observable("circle");
page.greaterSelectedMarker = ko.observable("circle");
page.dtSeries = ko.observableArray([]);
page.project = ko.observable();
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.LastFilter;
page.TableName;
page.FieldList;
page.ContentFilter;


page.ExportIndividualMonthPdf = function() {
    kendo.drawing.drawDOM($(".individual-month"))
    .then(function(group) {
        // Render the result as a PDF file
        return kendo.drawing.exportPDF(group, {
            paperSize: "auto",
            margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" }
        });
    })
    .done(function(data) {
        // Save the PDF file
        kendo.saveAs({
            dataURI: data,
            fileName: "Power-Curve-Scatter.pdf",
        });
    });
}

page.getPDF = function(selector){
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


page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "deviation", "text": "Nacelle Deviation" },
    { "value": "temp", "text": "Nacelle Temperature" },
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "ambient", "text": "Ambient Temperature" },
    { "value": "windspeed_dev", "text": "Wind Speed Std. Dev." },
    { "value": "windspeed_ti", "text": "TI Wind Speed" },
]);

vm.currentMenu('Scatter with Filter');
vm.currentTitle('Scatter with Filter');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Scatter with Filter',
    href: viewModel.appName + 'page/analyticpcscatteranalysis'
}]);

page.LoadData = function() {
    page.getPowerCurveScatter();
}

page.refreshChart = function() {
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();
    switch(page.scatterType) {
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
    // if(page.scatterType == 'pitch') {
    //     $("#txtLessVal").val(0);
    //     $('#txtGreaterVal').val(0);
    // } else {
    //     $("#txtLessVal").val(20);
    //     $('#txtGreaterVal').val(20);
    // }
    page.LoadData();
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();
    var turbine = $("#turbineList").data('kendoDropDownList').value();

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
        turbine: turbine,
        // project: fa.project, // has a bug, always selected Tejuva instead selected project
        project: $('#projectList').data('kendoDropDownList').value(), // temporary changed
        scatterType: page.scatterType,
        // lessValue: parseInt(lessValue,10),
        // greaterValue: parseInt(greaterValue,10),
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
        page.LastFilter = res.data.LastFilter;
        page.FieldList = res.data.FieldList;
        page.TableName = res.data.TableName;
        page.ContentFilter = res.data.ContentFilter;

        var result = res.data.Data;

        page.dtSeries(result);
        page.createChart(page.dtSeries());

    
        app.loading(false);
    });
}


page.changeView=function(param){

    var lessColor = $("#lessColor").data("kendoColorPicker").value();
    var greaterColor = $("#greaterColor").data("kendoColorPicker").value();
    var lessMarker = $("#lessMarker").data("kendoDropDownList").value();
    var greaterMarker = $("#greaterMarker").data("kendoDropDownList").value();

    $.each(page.dtSeries(), function(index, value){
        if(value.name.indexOf(param) !== -1){
            page.dtSeries()[index].color = (param == ">" ? greaterColor : lessColor);
            page.dtSeries()[index].markers = {
                size : 2,
                type : (param == ">" ? greaterMarker : lessMarker),
                background : (param == ">" ? greaterColor : lessColor),
            }
        }
    });

    page.createChart(page.dtSeries());
}

page.createChart = function(dtSeries){
        $('#scatterChart').html("");
        $("#scatterChart").kendoChart({
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

$(document).ready(function() {

    $('#btnRefresh').on('click', function() {
        setTimeout(function(){
            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value(); 


            page.project(project);
            page.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            page.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

           page.LoadData();
        }, 300);
    });

    $.when(di.getAvailDate()).done(function() {
        setTimeout(function(){
            fa.LoadData();


            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

            page.project(fa.project);
            page.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            page.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

            page.LoadData();
        },600);
       
    });
});