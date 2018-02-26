'use strict';
var im = {};
im.loadFirstTime = ko.observable(true);
im.reset = function(){
    im.loadFirstTime(true);
}
im.refresh = function() {
    if(im.loadFirstTime()) {
        im.LoadData();
        $(window).scroll(im.sticky_relocate);
        im.sticky_relocate();
        im.loadFirstTime(false);
    }
}

im.dataPCEachTurbine = ko.observableArray([]);
var msListOfChart = [];
var msListOfButton = {};
var msListOfCategory = [];
im.LastFilter;
im.TableName;
im.FieldList;
im.ContentFilter;

im.PrintPdf = ko.observable(false);
im.getPDF = function(selector){
    
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    kendo.drawing.drawDOM(selector, {
        paperSize: "A3",
        margin: {
            bottom: 80,
            left: 20,
            right: 20,
            top: 50
        },
        landscape: true,
        scale: 0.5,
        template: kendo.template($("#page-template").html())(
        {
            project: project,
            dateStart: moment($('#dateStart').data('kendoDatePicker').value()).format("DD-MMM-YYYY"),
            dateEnd: moment($('#dateEnd').data('kendoDatePicker').value()).format("DD-MMM-YYYY"),
            legend : msListOfCategory,
        })
    }).then(function(group){
        console.log(msListOfCategory);
        kendo.drawing.pdf.saveAs(group, project+"PCIndividualMonth"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },400)
    });
}

im.PowerCurveExporttoExcel = function(tipe, isSplittedSheet, isMultipleProject) {
    app.loading(true);
    var namaFile = tipe;
    if (!isSplittedSheet) {
        namaFile = fa.project + " " + tipe;
    }

    var param = {
        Filters: im.LastFilter,
        FieldList: im.FieldList,
        Tablename: im.TableName,
        TypeExcel: namaFile,
        ContentFilter: im.ContentFilter,
        IsSplittedSheet: isSplittedSheet,
        IsMultipleProject: isMultipleProject,
    };
    if (tipe.indexOf("Details") > 0) {
        var param = {
            Filters: im.LastFilterDetails,
            FieldList: im.FieldListDetails,
            Tablename: im.TableNameDetails,
            TypeExcel: namaFile,
            ContentFilter: im.ContentFilterDetails,
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

im.getColor = function(color) {
    var c;
    if(/^#([A-Fa-f0-9]{3}){1,2}$/.test(color)){
        c= color.substring(1).split('');
        if(c.length== 3){
            c= [c[0], c[0], c[1], c[1], c[2], c[2]];
        }
        c= '0x'+c.join('');
        return 'rgba('+[(c>>16)&255, (c>>8)&255, c&255].join(',')+',1)';
    }
    throw new Error('Bad Hex');
}


im.ExportIndividualMonthPdf = function() {
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
            fileName: "Individual-Month.pdf",
        });
    });
}

im.showHideLegend = function (index) {
    var idName = "btn" + index;
    msListOfButton[idName] = !msListOfButton[idName];
    if (msListOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#E0E0E0', 'border-color': '#E0E0E0' });
    } else {
        $("#" + idName).css({ 'background': msListOfCategory[index].color, 'border-color': msListOfCategory[index].color });
    }
    $.each(msListOfChart, function (idx, idChart) {
        $(idChart).data("kendoChart").options.series[index].visible = msListOfButton[idName];
        $(idChart).data("kendoChart").refresh();
    });
}

im.LoadData = function() {
    app.loading(true);
    setTimeout(function () {
        var param = {
            turbine: fa.turbine(),
            project: fa.project,
            engine: fa.engine,
        };
        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getlistpowercurvemonthly", param, function (res) {
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }
            im.LastFilter = res.data.LastFilter;
            im.FieldList = res.data.FieldList;
            im.TableName = res.data.TableName;
            im.ContentFilter = res.data.ContentFilter;
            if (res.data.Data != null) {
                im.dataPCEachTurbine(_.sortBy(res.data.Data, 'Name'));
                im.InitLinePowerCurve();
            }
            if (res.data.Category != null) {
                msListOfCategory = res.data.Category;
                $("#legend-list").html("");
                msListOfButton = {};
                $.each(msListOfCategory, function (idx, val) {
                    if(idx > 0) {
                        var idName = "btn" + idx;
                        msListOfButton[idName] = true;
                        $("#legend-list").append(
                            '<button id="' + idName + 
                            '" class="btn btn-default btn-sm btn-legend" type="button" onclick="im.showHideLegend(' + idx + ')" style="border-color:' + 
                            val.color + ';background-color:' + val.color + ';"></button>' +
                            '<span class="span-legend">' + val.category + '</span>'
                        );
                    }
                });
            }
            app.loading(false);
        })
    }, 300);
}

im.InitLinePowerCurve = function() {
    msListOfChart = [];

    $.each(im.dataPCEachTurbine(), function (i, dataTurbine) {
        var name = dataTurbine.Name
        var idChart = "#im-chart-" + dataTurbine.Name
        msListOfChart.push(idChart);
        var rotation = 300;
        var heightVal = 250;
        var isPannable = false;
        var isZoomable = false;
        var isTitle = true;
        var titleFont = '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var titleAxisFont = '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var labelAxisFont = '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
        var tooltipAxisFont = '10px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';

        if (fa.turbine().length == 1) {
            rotation = 0;
            heightVal = 400;
            isPannable = false;
            isZoomable = false;
            isTitle = false;
            var titleFont = '20px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var titleAxisFont = '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var labelAxisFont = '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            var tooltipAxisFont = 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif';
            $(".power-curve-item").removeClass("col-md-4");
        }
        
        $(idChart).html("");
        $(idChart).kendoChart({
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            title: {
                text: name,
                font: titleFont,
                visible: isTitle
            },
            legend: {
                position: "bottom",
                visible: false,
            },
            chartArea: {
                height: heightVal
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
            series: dataTurbine.Data,
            categoryAxis: {
                labels: {
                    step: 1
                }
            },
            valueAxis: [{
                labels: {
                    format: "N0",
                }
            }],
            xAxis: {
                majorUnit: 1,
                title: {
                    text: "Wind Speed (m/s)",
                    font: titleAxisFont,
                    color: "#585555",
                    visible: true,
                },
                labels: {
                    font: labelAxisFont,
                },
                crosshair: {
                    visible: true,
                    tooltip: {
                        visible: true,
                        format: "N2",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: tooltipAxisFont,
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
                    font: titleAxisFont,
                    color: "#585555"
                },
                labels: {
                    format: "N0",
                    font: labelAxisFont,
                    rotation: rotation,
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
                        font: tooltipAxisFont,
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
            },
            tooltip: {
                visible: true,
                format: "{1}in {0} minutes",
                template: "#= series.name #",
                shared: true,
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: tooltipAxisFont,
                border: {
                    color: "#eee",
                    width: "2px",
                },
            },
            pannable: isPannable,
            zoomable: isZoomable
        });
        $(idChart).data("kendoChart").refresh();
    });
}

im.sticky_relocate = function() {
    var window_top = $(window).scrollTop();
    var div_top = $('#legend-anchor').offset().top;
    if (window_top > div_top) {
        $('#legend-list').addClass('legend');
        $('#legend-anchor').height($('#legend-list').outerHeight());
    } else {
        $('#legend-list').removeClass('legend');
        $('#legend-anchor').height(0);
    }
}