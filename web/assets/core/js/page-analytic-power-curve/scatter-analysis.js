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

page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "deviation", "text": "Nacelle Deviation" },
    { "value": "pitch", "text": "Pitch Angle" },
    /*{ "value": "power", "text": "Temperature Analysis" },
    { "value": "grid", "text": "Temperature Analysis" },*/
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
    if(page.scatterType == 'pitch') {
        $("#txtLessVal").val(0);
        $('#txtGreaterVal').val(0);
    } else {
        $("#txtLessVal").val(20);
        $('#txtGreaterVal').val(20);
    }
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
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: turbine,
        project: fa.project,
        scatterType: page.scatterType,
        lessValue: parseInt(lessValue,10),
        greaterValue: parseInt(greaterValue,10),
        lessColor: lessColor,
        greaterColor: greaterColor,
        lessMarker: lessMarker, 
        greaterMarker: greaterMarker
    };
    
    
    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatteranalysis", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }

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
            }
        });
}

$(document).ready(function() {

    $('#btnRefresh').on('click', function() {
        setTimeout(function(){
           page.LoadData();
        }, 300);
    });

    $.when(di.getAvailDate()).done(function() {
        setTimeout(function(){
            fa.LoadData();
            page.LoadData();
        },500);
       
    });
});