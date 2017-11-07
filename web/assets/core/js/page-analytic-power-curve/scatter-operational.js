'use strict';

viewModel.AnalyticPowerCurve = new Object();
var page = viewModel.AnalyticPowerCurve;

vm.currentMenu('Operational Power Curve');
vm.currentTitle('Operational Power Curve');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Operational Power Curve',
    href: viewModel.appName + 'page/analyticpcscatteroperational'
}]);


page.scatterType = ko.observable('');
page.scatterList = ko.observableArray([
    { "value": "pitch", "text": "Pitch Angle" },
    { "value": "rotor", "text": "Rotor RPM" },
    { "value": "generatorrpm", "text": "Generator RPM" },
    { "value": "windspeed", "text": "Wind Speed" },
]);

page.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Today" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);
page.turbineList = ko.observableArray([]);
page.projectList = ko.observableArray([]);
page.dateStart = ko.observable();
page.dateEnd = ko.observable();
page.turbine = ko.observableArray([]);
page.project = ko.observable();
page.sScater = ko.observable(false);


page.rawturbine = ko.observableArray([]);
page.rawproject = ko.observableArray([]);

var lastPeriod = "";
var turbineval = [];

page.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var availDateAll = res.data;
        var projectVal = $("#projectList1").data("kendoDropDownList").value();

        var namaproject = "";

        if( projectVal == undefined || projectVal == "") {
            namaproject = "Tejuva";
        }else{
            namaproject= projectVal;
        }

        
        var minDate  = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY')));
        var maxDate = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY')));

        var maxDateData = new Date(availDateAll[namaproject]["ScadaData"][1]);
        var startDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

        $("#periodList").data("kendoDropDownList").value("custom");
        $("#periodList").data("kendoDropDownList").value("custom");

        $('#dateStart').data('kendoDatePicker').value(startDate);
        $('#dateEnd').data('kendoDatePicker').value(maxDate);
        $('#dateStart2').data('kendoDatePicker').value(startDate);
        $('#dateEnd2').data('kendoDatePicker').value(maxDate);

        $('#availabledatestartscada').html(minDate);
        $('#availabledateendscada').html(maxDate);
    });
}
page.populateTurbine = function (selected) {
    if (page.rawturbine().length == 0) {
        page.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        var dataturbine = [];
        // var allturbine = {}
        // $.each(page.rawturbine(), function (key, val) {
        //     turbineval.push(val);
        // });
        // allturbine.value = "All Turbine";
        // allturbine.text = "All Turbines";
        // datavalue.push(allturbine);

        if (selected==""){
            selected = page.rawproject()[0].Value;
        }
        
        $.each(page.rawturbine(), function (key, val) {
            if (selected == val.Project){
                var data = {};
                data.value = val.Value;
                data.label = val.Turbine;
                datavalue.push(data);
                dataturbine.push(val);
            }
        });
        page.turbineList(datavalue);
        page.turbine(dataturbine);
    }

    setTimeout(function () {
        $('#turbineList1').data('kendoDropDownList').select(0);
        $('#turbineList2').data('kendoDropDownList').select(1);
    }, 50);
};

page.populateProject = function (selected) {
    if (page.rawproject().length == 0) {
        page.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        $.each(page.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Name;
            datavalue.push(data);
        });
        page.projectList(datavalue);

        setTimeout(function () {
            page.populateTurbine(selected);
        }, 100);
    }
};


page.showHidePeriod = function (callback) {
    var period = $('#periodList').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.getDateMax(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
    
    if (period == "custom") {
        $(".show_hide").show();
        $('#dateStart').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#dateEnd').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });

        $('#dateStart').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startMonthDate);
            $('#dateEnd').data('kendoDatePicker').value(endMonthDate);

            $(".show_hide").show();
        } else if (period == "annual") {
            $('#dateStart').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",

            });
            $('#dateEnd').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",
            });

            $('#dateStart').data('kendoDatePicker').value(startYearDate);
            $('#dateEnd').data('kendoDatePicker').value(endYearDate);

            $(".show_hide").show();
        } else {
            if (period == 'last24hours') {
                $('#dateStart').data('kendoDatePicker').value(last24hours);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            } else if (period == 'last7days') {
                $('#dateStart').data('kendoDatePicker').value(lastweek);
                $('#dateEnd').data('kendoDatePicker').value(endMonthDate);
            }
            $(".show_hide").hide();
        }
        lastPeriod = period;
    }

    setTimeout(function () {
        callback;
    }, 50);
}

page.showHidePeriod2 = function (callback) {
    var period = $('#periodList2').data('kendoDropDownList').value();

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var startMonthDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth()-1, 1, 0, 0, 0, 0));
    var endMonthDate = new Date(app.getDateMax(maxDateData));
    var startYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var endYearDate = new Date(Date.UTC(moment(maxDateData).get('year'), 0, 1, 0, 0, 0, 0));
    var last24hours = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 1, 0, 0, 0, 0));
    var lastweek = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));
    if (period == "custom") {
        $(".show_hide2").show();
        $('#dateStart2').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });
        $('#dateEnd2').data('kendoDatePicker').setOptions({
            start: "month",
            depth: "month",
            format: 'dd-MMM-yyyy'
        });

        $('#dateStart2').data('kendoDatePicker').value(startMonthDate);
        $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
    } else {
        var today = new Date();
        if (period == "monthly") {
            $('#dateStart2').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });
            $('#dateEnd2').data('kendoDatePicker').setOptions({
                start: "year",
                depth: "year",
                format: "MMM yyyy",
            });

            $('#dateStart2').data('kendoDatePicker').value(startMonthDate);
            $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);

            $(".show_hide2").show();
        } else if (period == "annual") {
            $('#dateStart2').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",

            });
            $('#dateEnd2').data('kendoDatePicker').setOptions({
                start: "decade",
                depth: "decade",
                format: "yyyy",
            });

            $('#dateStart2').data('kendoDatePicker').value(startYearDate);
            $('#dateEnd2').data('kendoDatePicker').value(endYearDate);

            $(".show_hide2").show();
        } else {
            if (period == 'last24hours') {
                $('#dateStart2').data('kendoDatePicker').value(last24hours);
                $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
            } else if (period == 'last7days') {
                $('#dateStart2').data('kendoDatePicker').value(lastweek);
                $('#dateEnd2').data('kendoDatePicker').value(endMonthDate);
            }
            $(".show_hide2").hide();
        }
        lastPeriod = period;
    }

    setTimeout(function () {
        callback;
    }, 50);
}
page.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }

        var availDateAll = res.data;
        var projectVal = $("#projectList1").data("kendoDropDownList").value();

        var namaproject = "";

        if( projectVal == undefined || projectVal == "") {
            namaproject = "Tejuva";
        }else{
            namaproject= projectVal;
        }

        
        var minDate  = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY')));
        var maxDate = (kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY')));

        var maxDateData = new Date(availDateAll[namaproject]["ScadaData"][1]);
        var startDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

        $("#periodList").data("kendoDropDownList").value("custom");
        $("#periodList").data("kendoDropDownList").value("custom");

        $('#dateStart').data('kendoDatePicker').value(startDate);
        $('#dateEnd').data('kendoDatePicker').value(maxDate);
        $('#dateStart2').data('kendoDatePicker').value(startDate);
        $('#dateEnd2').data('kendoDatePicker').value(maxDate);

        $('#availabledatestartscada').html(minDate);
        $('#availabledateendscada').html(maxDate);
    });
}

page.InitDefaultValue = function () {
    page.getAvailDate();
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    $("#periodList2").data("kendoDropDownList").value("custom");
}

page.LoadData = function() {
    page.getPowerCurveScatter();
}

page.refreshChart = function() {
    page.LoadData();
}

page.getPowerCurveScatter = function() {
    app.loading(true);
    page.scatterType = $("#scatterType").data('kendoDropDownList').value();

    var p1DateStart = $('#dateStart').data('kendoDatePicker').value();
        p1DateStart = new Date(Date.UTC(p1DateStart.getFullYear(), p1DateStart.getMonth(), p1DateStart.getDate(), 0, 0, 0));

    var p1DateEnd  = $('#dateEnd').data('kendoDatePicker').value();
        p1DateEnd = new Date(Date.UTC(p1DateEnd.getFullYear(), p1DateEnd.getMonth(), p1DateEnd.getDate(), 0, 0, 0));

    var p2DateStart = $('#dateStart2').data('kendoDatePicker').value();
        p2DateStart = new Date(Date.UTC(p2DateStart.getFullYear(), p2DateStart.getMonth(), p2DateStart.getDate(), 0, 0, 0));

    var p2DateEnd  = $('#dateEnd2').data('kendoDatePicker').value();
        p2DateEnd = new Date(Date.UTC(p2DateEnd.getFullYear(), p2DateEnd.getMonth(), p2DateEnd.getDate(), 0, 0, 0));

    if (p1DateStart - p1DateEnd > 25200000) {
        toolkit.showError("Invalid Date Range Selection for Filter 1");
    } else if(p2DateStart - p2DateEnd > 25200000) {
        toolkit.showError("Invalid Date Range Selection for Filter 2");
    } else {

        var param1 = {
            Period       : $('#periodList').data('kendoDropDownList').value(),
            Project      : $("#projectList1").data("kendoDropDownList").value(),
            Turbine      : $("#turbineList1").data('kendoDropDownList').value(),
            DateStart    : p1DateStart,
            DateEnd      : p1DateEnd,
            ScatterType  : page.scatterType,
        };
        var param2 = {
            Period       : $('#periodList2').data('kendoDropDownList').value(),
            Project      : $("#projectList1").data("kendoDropDownList").value(),
            Turbine      : $("#turbineList2").data('kendoDropDownList').value(),
            DateStart    : p2DateStart,
            DateEnd      : p2DateEnd,
            ScatterType  : page.scatterType,
        };
        var param = [param1, param2];        

        toolkit.ajaxPost(viewModel.appName + "analyticpowercurve/getpcscatteroperational", param, function(res) {
            if (!app.isFine(res)) {
                return;
            }
            var dtSeries = res.data.Data;
            
            var minAxisY = res.data.MinAxisY;
            var maxAxisY = res.data.MaxAxisY;
            var minAxisX = res.data.MinAxisX;
            var maxAxisX = res.data.MaxAxisX;
            var name = '';
            var title = '';
            var xAxis = {};
            var measurement = '';
            var format = 'N0'
            if(maxAxisX - minAxisX < 7) {
                format = 'N2'
            }
            switch(page.scatterType) {
                case "pitch":
                    name = 'pitchAxis'
                    title = 'Angle (Degree)'
                    measurement = String.fromCharCode(176)
                    break;
                case "rotor":
                    name = "rotorAxis"
                    title = "Revolutions per Minute (RPM)";
                    measurement = 'rpm'
                    break;
                case "generatorrpm":
                    name = "generatorAxis"
                    title = "Generator per Minute (RPM)";
                    measurement = 'rpm'
                    break;
                case "windspeed":
                    name = "windspeedAxis"
                    title = "Avg. Wind Speed (m/s)";
                    measurement = 'm/s'
                    break;
            }
            xAxis = {
                name: name,
                title: {
                    text: title,
                    visible: true,
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                labels: {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    format: format
                },
                line: {
                    visible: false
                },
                axisCrossingValue: -10,
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
                        template: "#= kendo.toString(value, 'n2') # " + measurement,
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
                // majorUnit: 0.5,
                min: minAxisX,
                max: maxAxisX,
            }
            var yAxis = {};
            yAxis = {
                name: "powerAxis",
                title: {
                    text: "Generation (KW)",
                    visible: true,
                    font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                labels: {
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                line: {
                    visible: false
                },
                axisCrossingValue: -10,
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
                        template: "#= kendo.toString(value, 'n2') # kWh",
                        background: "rgb(255,255,255, 0.9)",
                        color: "#58666e",
                        font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        border: {
                            color: "#eee",
                            width: "2px",
                        },
                    }
                },
                // majorUnit: 0.5,
                min: minAxisY,
                max: maxAxisY,
            }

            $('#scatterChart').html("");
            $("#scatterChart").kendoChart({
                theme: "flat",
                pdf: {
                  fileName: "DetailPowerCurve.pdf",
                },
                title: {
                    // text: "Scatter Power Curves | Project : "+fa.project.substring(0,fa.project.indexOf("(")).project+""+$(".date-info").text(),
                     text: "Scatter Power Curves",
                    visible: false,
                    font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                },
                legend: {
                    position: "bottom",
                    offsetX: 40,
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
                    }
                }],
                xAxis: xAxis,
                yAxes: yAxis
            });
            app.loading(false);
        });
    }
}

page.setProjectTurbine = function(projects, turbines, selected){
    page.rawproject(projects);
    page.rawturbine(turbines);
    page.populateProject(selected);
};


$(document).ready(function() {
    $('#btnRefresh').on('click', function() {
        setTimeout(function(){
            page.LoadData();
        }, 300);
    });

    $('#projectList1').kendoDropDownList({
        change: function () { 
            var project = $('#projectList1').data("kendoDropDownList").value();
            page.getAvailDate();
            page.populateTurbine(project);
         }
    });

    $.when(page.InitDefaultValue()).done(function() {
        setTimeout(function(){
            page.LoadData();
        },200);
       
    });
});