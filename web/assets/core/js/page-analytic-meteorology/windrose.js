'use strict';

viewModel.WindRose = new Object();
var wr = viewModel.WindRose;

wr.dataWindroseEachTurbine = ko.observableArray([]);
var maxValue = 0;
wr.sectorDerajat = ko.observable(0);
var listOfChart = [];
var listOfButton = {};
var listOfButtonZoom = {};

wr.ExportWindRose = function () {
    var chart = $("#wr-chart").getKendoChart();
    chart.exportPDF({ paperSize: "auto", margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" } }).done(function (data) {
        kendo.saveAs({
            dataURI: data,
            fileName: "WindRose.pdf",
        });
    });
}

wr.getColor = function(color, index) {
    var c;
    var idName = "btn" + index;
    if (!listOfButton[idName]) {
        color = "#8f8f8f";
    }
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

wr.getPDF = function(selector){
    app.loading(true);
    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  


    kendo.drawing.drawDOM(selector, {
        paperSize: "A3",
        margin: {
            bottom: 50,
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
            legend : listOfCategory,
        })
    }).then(function(group){
        kendo.drawing.pdf.saveAs(group, project+"WindRose"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}


wr.showHideLegendWR = function (index) {
    var idName = "btn" + index;
    listOfButton[idName] = !listOfButton[idName];
    if (listOfButton[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChart, function (idx, idChart) {
       if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
          $(idChart).data("kendoChart").options.series[index].visible = listOfButton[idName];
          $(idChart).data("kendoChart").refresh();
        }
    });
}

wr.showHideLegendZoom = function (idxLegend) {
    var idName = "btnZoom" + idxLegend;
    listOfButtonZoom[idName] = !listOfButtonZoom[idName];
    if (listOfButtonZoom[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[idxLegend], 'border-color': colorFieldsWR[idxLegend] });
    }
    // $.each(listOfChart, function (idx, idChart) {
    //    if($(idChart).data("kendoChart").options.series.length - 1 >= idxLegend) {
    //       $(idChart).data("kendoChart").options.series[idxLegend].visible = listOfButton[idName];
    //       $(idChart).data("kendoChart").refresh();
    //     }
    // });
    if($("#windroseZoom").data("kendoChart").options.series.length - 1 >= idxLegend) {
        $("#windroseZoom").data("kendoChart").options.series[idxLegend].visible = listOfButtonZoom[idName];
        $("#windroseZoom").data("kendoChart").refresh();
    }
}

wr.ZoomChart = function(divID){
    $("#modalDetail").on("shown.bs.modal", function () { 
        /*WINDROSE LEGEND INITIAL*/
        var idxChart = "#"+divID.replace("btn-zoom-", "chart-");
        var indexChart = 0;
        $.each(listOfChart, function (idx, idChart) {
            if(idChart == idxChart){
                indexChart = idx;
            }
        });
        var titleZoom = divID.replace("btn-zoom-", "");
        if(titleZoom.indexOf("MetTower") >= 0) {
            titleZoom = "Chart Met Tower";
        }
        $('#titleWRZoom').html('<strong>' + titleZoom + '</strong>');
        
        $("#legend-list-zoom").html("");
        $.each(listOfCategory, function (idx, val) {
            var idName = "btnZoom" + idx;
            listOfButtonZoom[idName] = true;
            $("#legend-list-zoom").append(
                '<button id="' + idName + '" class="btn btn-default btn-sm btn-legend" type="button" onclick="wr.showHideLegendZoom(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + ';"></button>' +
                '<span class="span-legend">' + val.category + '</span>'
            );
        });
        wr.initZoomChart(wr.dataWindroseEachTurbine()[indexChart]);       

    }).modal('show');

    $('#modalDetail').on('hidden.bs.modal', function (e) {
        $('#modalDetail').off();
    });
}

wr.initZoomChart = function (dataSource) {
    var breakDownVal = $("#nosection").data("kendoDropDownList").value();
    var stepNum = 1
    var gapNum = 1
    if (breakDownVal == 36) {
        stepNum = 3
        gapNum = 0
    } else if (breakDownVal == 24) {
        stepNum = 2
        gapNum = 0
    } else if (breakDownVal == 12) {
        stepNum = 1
        gapNum = 0
    }

    var majorUnit = 10;
    if(maxValue < 40) {
        majorUnit = 5;
    }

    var name = dataSource.Name
    if (name == "MetTower") {
        name = "Met Tower"
    }

    $("#windroseZoom").kendoChart({
        theme: "nova",
        chartArea: {
            height: 500,
            margin: 0,
            padding: 0
        },

        title: {
            visible: false
        },
        legend: {
            visible: false,
        },
        dataSource: {
            data: dataSource.Data,
            group: {
                field: "WsCategoryNo",
                dir: "asc"
            },
            sort: {
                field: "DirectionNo",
                dir: "asc"
            }
        },
        seriesColors: colorFieldsWR,
        series: [{
            type: "radarColumn",
            stack: true,
            field: "Contribution",
            gap: gapNum,
            border: {
                width: 1,
                color: "#7f7f7f",
                opacity: 0.5
            },
        }],
        categoryAxis: {
            field: "DirectionDesc",
            visible: true,
            majorGridLines: {
                visible: true,
                step: stepNum
            },
            labels: {
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: true,
                step: stepNum
            }
        },
        valueAxis: {
            labels: {
                template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            // majorUnit: majorUnit,
            // max: maxValue,
            // min: 0
        },
        tooltip: {
            visible: true,
            template: "#= category #"+String.fromCharCode(176)+" (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
            background: "rgb(255,255,255, 0.9)",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },
        }
    });
    $("#windroseZoom").data('kendoChart').refresh();
}

wr.initChart = function () {
    listOfChart = [];
    var breakDownVal = $("#nosection").data("kendoDropDownList").value();
    var stepNum = 1
    var gapNum = 1
    if (breakDownVal == 36) {
        stepNum = 3
        gapNum = 0
    } else if (breakDownVal == 24) {
        stepNum = 2
        gapNum = 0
    } else if (breakDownVal == 12) {
        stepNum = 1
        gapNum = 0
    }

    var majorUnit = 10;
    if(maxValue < 40) {
        majorUnit = 5;
    }

    $.each(wr.dataWindroseEachTurbine(), function (i, val) {
        var name = val.Name
        if (name == "MetTower") {
            name = "Met Tower"
        }

        var idChart = "#chart-" + val.Name
        listOfChart.push(idChart);
        // var pWidth = $('body').width() * 0.235;//$('body').width() * ($(idChart).closest('div.windrose-item').width() - 2) / 100;
        var pWidth = 290;

        $(idChart).kendoChart({
            theme: "nova",
            chartArea: {
                width: pWidth,
                height: pWidth,
                padding: 25
            },

            title: {
                text: name,
                font: '13px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            legend: {
                position: "bottom",
                labels: {
                    template: "#= (series.data[0] || {}).WsCategoryDesc #"
                },
                visible: false,
            },
            dataSource: {
                data: val.Data,
                group: {
                    field: "WsCategoryNo",
                    dir: "asc"
                },
                sort: {
                    field: "DirectionNo",
                    dir: "asc"
                }
            },
            seriesColors: colorFieldsWR,
            series: [{
                type: "radarColumn",
                stack: true,
                field: "Contribution",
                gap: gapNum,
                border: {
                    width: 1,
                    color: "#7f7f7f",
                    opacity: 0.5
                },
            }],
            categoryAxis: {
                field: "DirectionDesc",
                visible: true,
                majorGridLines: {
                    visible: true,
                    step: stepNum
                },
                labels: {
                    font: '11px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    visible: true,
                    step: stepNum
                }
            },
            valueAxis: {
                labels: {
                    template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                    font: '9px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                majorUnit: majorUnit,
                max: maxValue,
                min: 0
            },
            tooltip: {
                visible: true,
                template: "#= category #"+String.fromCharCode(176)+" (#= dataItem.WsCategoryDesc #) #= kendo.toString(value, 'n2') #% for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });
    });
}

wr.ChangeSector = function(){
    pm.isFirstWindRose(true);
    wr.WindRose();
}

wr.WindRose = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        if(pm.isFirstWindRose() === true){
            app.loading(true);
            setTimeout(function () {
                var breakDownVal = $("#nosection").data("kendoDropDownList").value();
                var secDer = 360 / breakDownVal;
                wr.sectorDerajat(secDer);
                var param = {
                    period: fa.period,
                    dateStart: fa.dateStart,
                    dateEnd: fa.dateEnd,
                    turbine: fa.turbine(),
                    project: fa.project,
                    breakDown: breakDownVal,
                };
                toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getflexidataeachturbine", param, function (res) {
                    if (!app.isFine(res)) {
                        app.loading(false);
                        return;
                    }
                    if (res.data.WindRose != null) {
                        maxValue = res.data.MaxValue;
                        
                        var metData;
                        var isMetExist = false;
                        var scadaData = res.data.WindRose;
                        if(res.data.WindRose[0].Name == "MetTower") {
                            isMetExist = true;
                            metData = res.data.WindRose[0];
                            scadaData = res.data.WindRose.slice(1);
                        }
                        var tempData = _.sortBy(scadaData, 'Name');
                        if(isMetExist) {
                            tempData.unshift(metData);
                        }
                        res.data.WindRose = tempData;
                        res.data.WindRose.forEach(function(val, idx){
                            res.data.WindRose[idx].idxseries = idx;
                        });

                        wr.dataWindroseEachTurbine(res.data.WindRose);
                        wr.initChart();
                    }

                    app.loading(false);
                    pm.isFirstWindRose(false);

                })
            }, 300);
            var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
            var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
            $('#availabledatestart').html(metDate);
            $('#availabledateend').html(scadaDate);


        }else{
            var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
            var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
            $('#availabledatestart').html(metDate);
            $('#availabledateend').html(scadaDate);
            setTimeout(function(){
                $.each(listOfChart, function(idx, elem){
                    $(elem).data("kendoChart").refresh();
                });
                app.loading(false);
            }, 300);
        }
    }
}

function sticky_relocate() {
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

$(function() {
    $(window).scroll(sticky_relocate);
    sticky_relocate();
});