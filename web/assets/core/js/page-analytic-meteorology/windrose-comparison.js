'use strict';

viewModel.WindRoseComparison = new Object();
var wrb = viewModel.WindRoseComparison;

// WIND ROSE COMPARISON
wrb.sectorDerajatComparison = ko.observable(0);
wrb.dataWindroseComparison = ko.observableArray([]);
wrb.project = ko.observable();
wrb.dateStart = ko.observable();
wrb.dateEnd = ko.observable();
var listOfChartComparison = [];
var listOfButtonComparison = {};
/*wrb.showHideLegendComparison = function (index) {
    var idName = "btn" + index;
    listOfButtonComparison[idName] = !listOfButtonComparison[idName];
    if (listOfButtonComparison[idName] == false) {
        $("#" + idName).css({ 'background': '#8f8f8f', 'border-color': '#8f8f8f' });
    } else {
        $("#" + idName).css({ 'background': colorFieldsWR[index], 'border-color': colorFieldsWR[index] });
    }
    $.each(listOfChartComparison, function (idx, idChart) {
       if($(idChart).data("kendoChart").options.series.length - 1 >= index) {
          $(idChart).data("kendoChart").options.series[index].visible = listOfButtonComparison[idName];
          $(idChart).data("kendoChart").refresh();
        }
    });
}*/

wrb.getPDF = function(selector){
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
      kendo.drawing.pdf.saveAs(group, project+"WindRoseComparison"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".pdf");
        setTimeout(function(){
            app.loading(false);
        },2000)
    });
}

wrb.InitTurbineListCompare = function () {
    if (wrb.dataWindroseComparison().Data.length > 1) {
        $("#checkAllCompare").html('<label id="checkAllLabel">' +
            '<input type="checkbox" id="showHideAllCompare" checked onclick="wrb.showHideAllCompare(this)" >' +
            '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
            '<span id="labelShowHideCompare"><b>Select All</b></span>' +
            '</label>');
    } else {
        $("#checkAllCompare").html("");
    }

    $("#turbine-list-compare").html("");
    $.each(wrb.dataWindroseComparison().Data, function (idx, val) {
        $("#turbine-list-compare").append('<div class="btn-group">' +
            '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="wrb.showHideCompare(' + val.idxseries + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check fa-check-windcomp" id="icon-windrose-comparison' + val.idxseries + '"></i></button>' +
            '<input class="chk-option-windcomp" type="checkbox" name="' + val.name + '" checked id="chk-windrose-comparison' + val.idxseries + '" hidden>' +
            '<button class="btn btn-default btn-sm turbine-btn wbtn" onclick="wrb.showHideCompare(' + val.idxseries + ')" type="button">' + val.name + '</button>' +
            '</div>');
    });
}

wrb.showHideAllCompare = function (e) {

    if (e.checked == true) {
        $('.fa-check-windcomp').css("visibility", 'visible');
        $.each(wrb.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = true;
            }
        });
        $('#labelShowHideCompare b').text('Select All');
    } else {
        $.each(wrb.dataWindroseComparison().Data, function (i, val) {
            if($("#WRChartComparison").data("kendoChart").options.series[i] != undefined){
                $("#WRChartComparison").data("kendoChart").options.series[i].visible = false;
            }  
        });
        $('.fa-check-windcomp').css("visibility", 'hidden');
        $('#labelShowHideCompare b').text('Select All');
    }
    $('.chk-option-windcomp').not(e).prop('checked', e.checked);

    $("#WRChartComparison").data("kendoChart").redraw();
}

wrb.showHideCompare = function (idx) {
    var stat = false;

    $('#chk-windrose-comparison' + idx).trigger('click');
    var chart = $("#WRChartComparison").data("kendoChart");
    var leTur = $('input[id*=chk-windrose-comparison][type=checkbox]').length

    if ($('input[id*=chk-windrose-comparison][type=checkbox]:checked').length == $('input[id*=chk-windrose-comparison][type=checkbox]').length) {
        $('#showHideAllCompare').prop('checked', true);
    } else {
        $('#showHideAllCompare').prop('checked', false);
    }

    if ($('#chk-windrose-comparison' + idx).is(':checked')) {
        $('#icon-windrose-comparison' + idx).css("visibility", "visible");
    } else {
        $('#icon-windrose-comparison' + idx).css("visibility", "hidden");
    }

    if ($('#chk-windrose-comparison' + idx).is(':checked')) {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#WRChartComparison").data("kendoChart").options.series[idx].visible = false
    }
    $("#WRChartComparison").data("kendoChart").redraw();
}

wrb.initChartWRC = function () {
    listOfChartComparison = [];
    var dataSeries = wrb.dataWindroseComparison().Data;
    // var categories = wrb.dataWindroseComparison().Categories;
    var nilaiMax = wrb.dataWindroseComparison().MaxValue;

    var majorUnit = 10;
    if(nilaiMax < 40) {
        majorUnit = 5;
    }

    $("#WRChartComparison").html("");
    setTimeout(function(){
        $("#WRChartComparison").kendoChart({
            theme: "flat",
            title: {
                text: "Wind Rose Comparison",
                font: '16px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                visible: false
            },
            chartArea : {
                height : 350
            },
            legend: {
                position: "bottom",
                visible: false,
                labels : {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    // template: kendo.template($("#legendItemTemplate").html()),
                }
            },
            series: dataSeries,
            xAxis: {
                majorUnit: 30,
                startAngle: 90,
                reverse: true,
                labels: {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
            },
            yAxis: {
                labels: {
                    template: kendo.template("#= kendo.toString(value, 'n0') #%"),
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
            },

            tooltip: {
                visible: true,
                template: "#= series.name # : #= dataItem.DirectionDesc #"+String.fromCharCode(176)+
                            " (#= kendo.toString(dataItem.Contribution, 'n2') #%) for #= kendo.toString(dataItem.Hours, 'n2') # Hours",
                background: "rgb(255,255,255, 0.9)",
                color: "#58666e",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                border: {
                    color: "#eee",
                    width: "2px",
                },
            }
        });
    }, 100);

    setTimeout(function(){
        var chart = $("#WRChartComparison").data("kendoChart");
        var viewModel = kendo.observable({
          series: chart.options.series,
          markerColor: function(e) {
            return e.get("visible") ? e.color : "grey";
          }
        });

        kendo.bind($("#legendWindrose"), viewModel);
    },300);

    wrb.InitTurbineListCompare();
    // $('#WRChartComparison').data('kendoChart').options.chartArea.width = $('#WRChartComparison').height() + ($('#WRChartComparison').height()/4);
    // $('#WRChartComparison').data('kendoChart').refresh();
}

wrb.WindRoseComparison = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        if(pm.isFirstWindRoseComparison() === true){
            app.loading(true);
            setTimeout(function () {
                // var breakDownVal = $("#nosectionComparison").data("kendoDropDownList").value();
                var breakDownVal = "36";
                var secDer = 360 / breakDownVal;
                // wrb.sectorDerajatComparison(secDer);
                var param = {
                    period: fa.period,
                    dateStart: fa.dateStart,
                    dateEnd: fa.dateEnd,
                    turbine: fa.turbine(),
                    project: fa.project,
                    breakDown: breakDownVal,
                };
                toolkit.ajaxPost(viewModel.appName + "analyticwindrose/getwindrosedata", param, function (res) {
                    if (!app.isFine(res)) {
                        app.loading(false);
                        return;
                    }
                    if (res.data != null) {
                        var metData;
                        var isMetExist = false;
                        var scadaData = res.data.Data;
                        if(res.data.Data[0].name == "Met Tower") {
                            isMetExist = true;
                            metData = res.data.Data[0];
                            scadaData = res.data.Data.slice(1);
                        }
                        var tempData = _.sortBy(scadaData, 'name');
                        if(isMetExist) {
                            tempData.unshift(metData);
                        }
                        res.data.Data = tempData;
                        res.data.Data.forEach(function(val, idx){
                            res.data.Data[idx].idxseries = idx;
                        });
                        wrb.dataWindroseComparison(res.data);
                        wrb.initChartWRC();
                    }

                    app.loading(false);
                    pm.isFirstWindRoseComparison(false);

                })
            }, 300);
            var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
            var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
            $('#availabledatestart').html(metDate);
            $('#availabledateend').html(scadaDate);

            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
            wrb.project(project);
            wrb.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            wrb.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

        }else{
            var metDate = 'Data Available (<strong>MET</strong>) from: <strong>' + availDateList.availabledatestartmet + '</strong> until: <strong>' + availDateList.availabledateendmet + '</strong>'
            var scadaDate = ' | (<strong>SCADA</strong>) from: <strong>' + availDateList.availabledatestartscada + '</strong> until: <strong>' + availDateList.availabledateendscada + '</strong>'
            $('#availabledatestart').html(metDate);
            $('#availabledateend').html(scadaDate);

            var project = $('#projectList').data("kendoDropDownList").value();
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
            wrb.project(project);
            wrb.dateStart(moment(new Date(dateStart)).format("DD-MMM-YYYY"));
            wrb.dateEnd(moment(new Date(dateEnd)).format("DD-MMM-YYYY"));

            setTimeout(function(){
                $.each(listOfChartComparison, function(idx, elem){
                    $(elem).data("kendoChart").refresh();
                });
                app.loading(false);
            }, 300);
        }
    }
}