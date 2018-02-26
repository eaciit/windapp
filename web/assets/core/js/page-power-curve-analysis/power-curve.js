'use strict';
var pc = {};
pc.loadFirstTime = ko.observable(true);
pc.deviationOpts = ko.observableArray([
    { "value": 0, "text": "<" },
    { "value": 1, "text": ">" },
]);
pc.deviationOpr = ko.observable(0);
pc.deviationVal = ko.observable(20);
pc.isDeviation = ko.observable(true);
pc.selectedFilter = ko.observable('');
pc.chartTypeValue = {
    line: 'line',
    scatter: 'scatter',
};
pc.chartType = ko.observable(pc.chartTypeValue.line); // line, scatter
pc.isSpecific = ko.observable(true);
pc.isShowDowntime = ko.observable(false);

pc.buildFilterInfo = function() {
    var isValid = $('#pc-is-valid').is(':checked');
    var isDeviation = pc.isDeviation();
    var infos = [];
    if(isValid) {
        infos.push('Valid');
    }
    if(isDeviation) {
        var devOpr = '<';
        if(pc.deviationOpr()==1) {
            devOpr = '>';
        }
        infos.push('Deviation '+ devOpr +' '+ pc.deviationVal() +'%');
    }
    if(infos.length <= 0) {
        infos.push('No filter selected.');
    }
    var info = '', delim = '';
    $.each(infos, function(i,v){
        info += delim + v;
        delim = ' | ';
    });
    pc.selectedFilter(info);
}
// this is required function for this object which accessed by main page
pc.reset = function(){
    this.loadFirstTime(true);
}
// this is required function for this object which accessed by main page
pc.refresh = function() {
    if(this.loadFirstTime()) {
        this.loadFirstTime(false);
        this.initElementEvents();
        this.internalRefresh();
    }
}
pc.internalRefresh = function(reloadData) {
    if(reloadData==undefined) {
        reloadData = true;
    }

    this.buildFilterInfo();
    if(reloadData) {
        if(pc.chartType()==pc.chartTypeValue.line) {
            this.initLineChart();
        } else {
            this.initScatterChart();
        }
    }
}
pc.initLineChart = function() {
    app.loading(true);

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();   
    var turbineList = [];
    $.each(turbines, function(i, val) {
        if (fa.project == val.Project) {
            turbineList.push(val.Turbine);
        }
    });

    var sUrl = "analyticpowercurve/getlistpowercurvescada"
    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
        turbine: $("#turbineList").val(),
        project: fa.project,
        isClean: $('#pc-is-valid').is(':checked'),
        isSpecific: pc.isSpecific(),
        isDeviation: $('#pc-is-deviation').is(':checked'),
        isPower0: false,
        DeviationVal: parseInt(pc.deviationVal()).toString(),
        DeviationOpr: parseInt(pc.deviationOpr()).toString(),
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

        $('#pc-chart').html("");
        $("#pc-chart").kendoChart({
            pdf: {
              fileName: "PowerCurve.pdf",
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
                background: "transparent",
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
        $("#pc-chart").data("kendoChart").refresh();
        pc.initTurbineList();
        app.loading(false);
    });
}
pc.initScatterChart = function() {
    var turbineList = [];
    var kolor = [];
    var dataTurbine = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'turbineid');

    var len = $('input[id*=chk-][type=checkbox]:checked').length;

    for (var a = 0; a < len; a++) {
        var chk = $('input[id*=chk-][type=checkbox]:checked')[a].name;
        turbineList.push(chk);
        var even = _.find(dataTurbine, function(nm) {
            return nm.turbineid == chk
        });
        kolor.push(even.color);
        var indOf = 0;
        for (var i = 0; i < colorField.length; i++) {
            if(colorField[i] === even.color) {
                indOf = i
            }
        }
        // var indOf = colorField.indexOf(even.color);
        // kolorDeg.push(colorDegField[indOf]);
    }

    var dtLine = JSON.parse(localStorage.getItem("dataTurbine"));

    app.loading(true);

    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  

    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
        turbine: turbineList,
        project: fa.project,
        isClean: $('#pc-is-valid').is(':checked'),
        isSpecific: pc.isSpecific(),
        isDeviation: $('#pc-is-deviation').is(':checked'),
        isPower0: false,
        DeviationVal: parseInt(pc.deviationVal()).toString(),
        DeviationOpr: parseInt(pc.deviationOpr()).toString(),
        ViewSession: "",
        Engine: fa.engine,
        Color: kolor,
        IsDownTime: pc.isShowDowntime(),
        ViewSession: "",
    };
    //lastParam = param;

    toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpowercurve", param, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var dataPowerCurves = res.data.Data;
        var dtSeries = new Array();
        if (dataPowerCurves != null) {
            if (dataPowerCurves.length > 0) {
                dtSeries = dtLine.concat(dataPowerCurves);
            }
        } else {
            dtSeries = dtLine;
        }

        $('#pc-chart').html("");
        $("#pc-chart").kendoChart({
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
                visible: false,
                align: "center",
                position: "bottom",

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
                // page.getLegendActive();
            },
        });

        app.loading(false);
        // if (page.showDownTime()) {
        //     $('#downtime-list').show();
        // } else {
        //     $('#downtime-list').hide();
        // }
        // page.ShowHideAfterInitChart();
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
        $("#pc-show-hide-check").html('<label>' +
            '<input type="checkbox" id="pc-show-hide-all" checked onclick="pc.showHideAllLegend(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHide"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#pc-show-hide-check").html("");
    }

    var totalDataShoulBeInProject = 0;
    var totalDataAvailInProject = 0;
    $("#pc-right-turbine-list").html("");
    $.each(dtTurbines, function(idx, val) {
        if(val.name != "Power Curve"){
            var nameTurbine = val.name;
             if ( fa.project == "Rajgarh" ) {
                nameTurbine = nameTurbine.replace("KH-", "-")
            }
            
            totalDataShoulBeInProject += val.totaldatashouldbe;
            totalDataAvailInProject += val.totaldata;
            $("#pc-right-turbine-list").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm pc-turbine-chk" type="button" onclick="pc.showHideLegend(' + idx + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="pc-icon-' + idx + '"></i></button>' +
            '<input class="chk-option" type="checkbox" name="' + val.turbineid + '" checked id="pc-chk-' + idx + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="pc.showDetail(\'' + val.turbineid + '\',\'' + val.turbineid + '\')" type="button">' + nameTurbine + ' <label id="dataavailpct-'+val.turbineid+'" class="label label-default pull-right" data-toggle="tooltip" title="Data available for turbine : '+ nameTurbine +'">'+ kendo.toString(val.dataavailpct, 'p1') +'</label></button>' +
            '</div>');
        }
    });
}
pc.showHideLegend = function(idx) {
    $('#pc-chk-' + idx).trigger('click');
    var chart = $("#pc-chart").data("kendoChart");

    if ($('input[id*=pc-chk-][type=checkbox]:checked').length == $('input[id*=pc-chk-][type=checkbox]').length) {
        $('#pc-show-hide-all').prop('checked', true);
    } else {
        $('#pc-show-hide-all').prop('checked', false);
    }

    if ($('#pc-chk-' + idx).is(':checked')) {
        $('#pc-icon-' + idx).css("visibility", "visible");
    } else {
        $('#pc-icon-' + idx).css("visibility", "hidden");
    }
    if (idx == $('input[id*=pc-chk-][type=checkbox]').length) {
        idx == 0
    }

    // check if turbines not all checked
    if (!$('#pc-show-hide-all').is(':checked')) {
        var chks = $('input[id*=pc-chk-][type=checkbox]:checked');
        var totalavail = 0;
        var totalCount = 0;
        var sampleAvail = 0;
        var sampleCount = 0;
        
        $.each(chks, function(idx, elm){
            var tbName = $(elm).attr('name');
            var tbAvail = 0;//page.totalAvailTurbines()[tbName];
            totalavail += tbAvail.avail;
            totalCount++;

            var elmAvail = $(elm).parent().find('button.wbtn').find('label').text().replace(' %', '');
            var currElmAvail = parseFloat(elmAvail);
            sampleAvail += currElmAvail; 
            sampleCount++;
        });

        if(totalCount > 0) {
            var selectedAvail = totalavail / totalCount;
            // page.totalAvail(selectedAvail);

            var selectedSampleAvail = (sampleAvail / 100) / sampleCount;
            // page.dataAvail(selectedSampleAvail);
        } else {
            // page.totalAvail(page.totalAvailAll());
            // page.dataAvail(page.dataAvailAll());
        }
    }

    chart._legendItemClick(idx);
    // page.getLegendActive();
}
pc.showHideAllLegend = function(e) {
    var chart = $("#pc-chart").data("kendoChart");
    var dtTurbines = _.sortBy(JSON.parse(localStorage.getItem("dataTurbine")), 'name');
    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(dtTurbines, function(i, val) {
            val.idxseries = val.idxseries - 1;
            if(val.name !== "Power Curve"){
                chart.options.series[val.idxseries].visible = true;
            }
        });
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(dtTurbines, function(i, val) {
            val.idxseries = val.idxseries - 1;
            if(val.name !== "Power Curve"){
                chart.options.series[val.idxseries].visible = false;
            }
        });
        $('.fa-check').css("visibility", 'hidden');
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);
    chart.redraw();
    // page.getLegendActive();
}
pc.initElementEvents = function() {
    $('#pc-deviation-opr').on('change', function(){
        if(pc.isDeviation()) {
            pc.internalRefresh();
        }
    });
    $('#pc-deviation-value').on('change', function(){
        var value = $(this).val();
        if(value=='') {
            $(this).val(0);
        }
        if(pc.isDeviation()) {
            pc.internalRefresh();
        }
    });
    $('#pc-is-valid').on('click', function(){
        pc.internalRefresh();
    });
    $('#pc-is-deviation').on('click', function(){
        pc.isDeviation($(this).is(':checked'));
        pc.internalRefresh();
    });
    $('#pc-show-scatter').on('click', function(){
        var showScatter = $(this).is(':checked');
        if(showScatter) {
            pc.chartType(pc.chartTypeValue.scatter);
        } else {
            pc.chartType(pc.chartTypeValue.line);
        }
        pc.internalRefresh();
    }); 
    $('#pc-sitespesific').on('change', function(){
        pc.isSpecific(true);
        pc.internalRefresh();
    });
    $('#pc-standardpc').on('change', function(){
        pc.isSpecific(false);
        pc.internalRefresh();
    });
}