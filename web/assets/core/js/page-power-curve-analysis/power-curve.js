'use strict';
var pc = {};
pc.loadFirstTime = ko.observable(true);
pc.reset = function(){
    this.loadFirstTime(true);
}
pc.refresh = function() {
    if(this.loadFirstTime()) {
        this.loadFirstTime(false);
        this.initElementEvents();
        this.initChart();
    }
}
pc.internalRefresh = function(firsttime) {
    if(firsttime==undefined) {
        firsttime = false;
    }
    this.loadFirstTime(firsttime);
    this.refresh();
}
pc.initChart = function() {
    app.loading(true);

    var sUrl = "analyticpowercurve/getlistpowercurvescada"

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();   
    var turbineList = [];
    $.each(turbines, function(i, val) {
        if (fa.project == val.Project) {
            turbineList.push(val.Turbine);
        }
    });
    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
        turbine: $("#turbineList").val(),
        project: fa.project,
        isClean: true,
        isSpecific: true,
        isDeviation: true,
        isPower0: false,
        DeviationVal: "20",
        DeviationOpr: "0",
        ViewSession: "",
        Engine: fa.engine,
    };

    toolkit.ajaxPost(viewModel.appName + sUrl, param, function(res) {
        if (!app.isFine(res)) {
            app.loading(false);
            return;
        }

        var tempData = [];
        var powerCurveData;
        res.data.Data.forEach(function(val, idx){
            if(val.name != "Power Curve") {
                tempData.push(val);
            } else {
                powerCurveData = val;
            }
        });

        tempData = _.sortBy(tempData, 'name')
        tempData.forEach(function(val, idx){
            tempData[idx].idxseries = idx+1;
        });
        tempData.push(powerCurveData);
        res.data.Data = tempData;

        var dataTurbine = res.data.Data;
        localStorage.setItem("dataTurbine", JSON.stringify(dataTurbine));

        $('#powerCurve').html("");
        $("#powerCurve").kendoChart({
            pdf: {
              fileName: "DetailPowerCurve.pdf",
            },
            theme: "flat",
            title: {
                text: "Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("("))+""+$(".date-info").text(),
                visible: false,
                font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
            },
            legend: {
                position: "bottom",
                visible: false,
            },
            chartArea: {
                height: 425,
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
                    step: 1
                }
            },
            valueAxis: [{
                labels: {
                    format: "N0",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                },
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
                axisCrossingValue: -1000,
                min: 0,
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
                format: "{1}in {0} minutes",
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
            // zoomable: true,
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
            },
            dataBound : function(){
                //page.getLegendActive();
            },
        });
        $("#powerCurve").data("kendoChart").refresh();
        pc.initTurbineList();
        app.loading(false);
    });
}
pc.initTurbineList = function() {
    var dtTurbines = JSON.parse(localStorage.getItem("dataTurbine"));
    var turbineList = [];
    $.each(turbines, function(i, val) {
        if (fa.project == val.Project) {
            turbineList.push(val.Value);
        }
    });
    if (turbineList.length > 1) {
        $("#showHideChk").html('<label>' +
            '<input type="checkbox" id="showHideAll" checked onclick="page.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#showHideChk").html("");
    }

    var totalDataShoulBeInProject = 0;
    var totalDataAvailInProject = 0;
    $("#right-turbine-list").html("");
    $.each(dtTurbines, function(idx, val) {
        if(val.name != "Power Curve"){
            var nameTurbine = val.name;
             if ( fa.project == "Rajgarh" ) {
                nameTurbine = nameTurbine.replace("KH-", "-")
            }
            
            totalDataShoulBeInProject += val.totaldatashouldbe;
            totalDataAvailInProject += val.totaldata;
            $("#right-turbine-list").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="page.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + idx + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.turbineid + '" checked id="chk-' + idx + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="page.toDetail(\'' + val.turbineid + '\',\'' + val.turbineid + '\')" type="button">' + nameTurbine + ' <label id="dataavailpct-'+val.turbineid+'" class="label label-default pull-right" data-toggle="tooltip" title="Data available for turbine : '+ nameTurbine +'">'+ kendo.toString(val.dataavailpct, 'p1') +'</label></button>' +
            '</div>');
        }
    });
}
pc.initElementEvents = function() {
    
}