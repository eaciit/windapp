'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

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
      kendo.drawing.pdf.saveAs(group, project+"PCScatter"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

page.plotWith = ko.observable();
page.scatterList = ko.observableArray([]);

vm.currentMenu('Scatter');
vm.currentTitle('Scatter');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Scatter',
    href: viewModel.appName + 'page/analyticpcscatter'
}]);

page.project = ko.observable();
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.LastFilter;
page.TableName;
page.FieldList;
page.ContentFilter;

page.LoadData = function() {
    page.getPowerCurveScatter();
}

page.setAxis = function(name, title) {
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
        result.crosshair.tooltip.padding = {left:5};
    } else {
        switch(page.scatterType) {
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

page.refreshChart = function() {
    page.LoadData();
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.plotWith = $.grep(page.scatterList(), function(e){ return e.Id == $("#scatterType").data("kendoDropDownList").value(); })[0];

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));   

    var param = {
        engine: fa.engine,
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: $("#turbineList").val(),
        project: $('#projectList').data("kendoDropDownList").value(),
        plotWith: page.plotWith,
    };

    // console.log(plotWith.Text)

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurvescatterrev", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        page.LastFilter = res.data.LastFilter;
        page.FieldList = res.data.FieldList;
        page.TableName = res.data.TableName;
        page.ContentFilter = res.data.ContentFilter;

        $("#turbineName").html($("#turbineList option:selected").text());
        
        var dtSeries = res.data.Data;
        var yAxes = [];
        var yAxis = page.setAxis("powerAxis", "Generation (KW)");
        yAxes.push(yAxis);
        var axis = page.setAxis("PlotWith", page.plotWith.Text);
        yAxes.push(axis);
        // switch(page.scatterType) {
        //     case "temp":
        //         var axis = page.setAxis("tempAxis", "Temperature (Celsius)");
        //         yAxes.push(axis);
        //         break;
        //     case "deviation":
        //         var axis = page.setAxis("deviationAxis", "Wind Direction (Degree)");
        //         yAxes.push(axis);
        //         break;
        //     case "pitch":
        //         var axis = page.setAxis("pitchAxis", "Angle (Degree)");
        //         yAxes.push(axis);
        //         break;
        //     case "ambient":
        //         var axis = page.setAxis("ambientAxis", "Temperature (Celcius)");
        //         yAxes.push(axis);
        //         break;
        //     case "windspeed_dev":
        //         var axis = page.setAxis("windspeed_dev", "Wind Speed Std. Dev.");
        //         yAxes.push(axis);
        //         break;
        //     case "windspeed_ti":
        //         var axis = page.setAxis("windspeed_ti", "TI Wind Speed");
        //         yAxes.push(axis);
        //         break;
        // }

        $('#scatterChart').html("");
        $("#scatterChart").kendoChart({
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

page.getPowerCurveScatterFieldList = function(){
    // var param  = {project : $('#ProjectList').data('kendoMultiSelect').value()}
    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatterfieldlist", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var data = res.data;
        if(data !== null){
            setTimeout(function(){
                page.scatterList(data);
                $("#scatterType").data("kendoDropDownList").select(0);
            },300)
        }   
    });

    return page.scatterList[0]
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

            $.when(page.getPowerCurveScatterFieldList()).done(function() {
                setTimeout(function(){
                    page.LoadData();
                },600);}
            )
        },1000);
        
    });
});