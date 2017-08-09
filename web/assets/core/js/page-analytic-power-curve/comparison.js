'use strict';


viewModel.PCComparison = new Object();
var pc = viewModel.PCComparison;


vm.currentMenu('Comparison');
vm.currentTitle('Comparison');
vm.breadcrumb([{
    title: "KPI's",
    href: '#'
}, {
    title: 'Power Curve',
    href: '#'
}, {
    title: 'Comparison',
    href: viewModel.appName + 'page/analyticpccomparison'
}]);


pc.periodList = ko.observableArray([
    { "value": "last24hours", "text": "Last 24 hours" },
    { "value": "last7days", "text": "Last 7 days" },
    { "value": "monthly", "text": "Monthly" },
    { "value": "annual", "text": "Annual" },
    { "value": "custom", "text": "Custom" },
]);

pc.turbineList = ko.observableArray([]);
pc.projectList = ko.observableArray([]);
pc.dateStart = ko.observable();
pc.dateEnd = ko.observable();
pc.turbine = ko.observableArray([]);
pc.project = ko.observable();

pc.rawturbine = ko.observableArray([]);
pc.rawproject = ko.observableArray([]);

var lastPeriod = "";
var turbineval = [];


/*pc.InitFirst = function () {
    $.when(
        app.ajaxPost(viewModel.appName + "/helper/getturbinelist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                pc.turbineList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                var dataturbine = [];
                if (res.data.length > 0) {
                    var allturbine = {}
                    $.each(res.data, function (key, val) {
                        turbineval.push(val);
                    });
                    // allturbine.value = "All Turbine";
                    // allturbine.text = "All Turbines";
                    // datavalue.push(allturbine);
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                        dataturbine.push(val);
                    });
                }
                pc.turbineList(datavalue);
                pc.turbine(dataturbine);
            }
        }),
        app.ajaxPost(viewModel.appName + "/helper/getprojectlist", {}, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            if (res.data.length == 0) {
                res.data = [];;
                pc.projectList([{ value: "", text: "" }]);
            } else {
                var datavalue = [];
                if (res.data.length > 0) {
                    $.each(res.data, function (key, val) {
                        var data = {};
                        data.value = val;
                        data.text = val;
                        datavalue.push(data);
                    });
                }
                pc.projectList(datavalue);
            }
        })

    ).then(function () {
        // $('#turbineList1').data('kendoDropDownList').value(["All Turbine"])
        // $('#turbineList2').data('kendoDropDownList').value(["All Turbine"])
        // override to set the value
        // $("#projectList1").data("kendoDropDownList").value("Tejuva");
        // $("#projectList2").data("kendoDropDownList").value("Tejuva");

        pc.project = $("#projectList").data("kendoDropDownList").value();
    });
}*/

pc.getAvailDate = function(){
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

        $('#dateStart').data('kendoDatePicker').value(startDate);
        $('#dateEnd').data('kendoDatePicker').value(maxDate);
        $('#dateStart2').data('kendoDatePicker').value(startDate);
        $('#dateEnd2').data('kendoDatePicker').value(maxDate);

        $('#availabledatestartscada').html(minDate);
        $('#availabledateendscada').html(maxDate);
    });
}
pc.populateTurbine = function (selected) {
    if (pc.rawturbine().length == 0) {
        pc.turbineList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        var dataturbine = [];
        // var allturbine = {}
        // $.each(pc.rawturbine(), function (key, val) {
        //     turbineval.push(val);
        // });
        // allturbine.value = "All Turbine";
        // allturbine.text = "All Turbines";
        // datavalue.push(allturbine);

        if (selected==""){
            selected = pc.rawproject()[0].Value;
        }
        
        $.each(pc.rawturbine(), function (key, val) {
            if (selected == val.Project){
                var data = {};
                data.value = val.Value;
                data.label = val.Turbine;
                datavalue.push(data);
                dataturbine.push(val);
            }
        });
        pc.turbineList(datavalue);
        pc.turbine(dataturbine);
    }

    setTimeout(function () {
        $('#turbineList1').data('kendoDropDownList').select(0);
        $('#turbineList2').data('kendoDropDownList').select(0);
    }, 50);
};

pc.populateProject = function (selected) {
    if (pc.rawproject().length == 0) {
        pc.projectList([{ value: "", text: "" }]);
    } else {
        var datavalue = [];
        $.each(pc.rawproject(), function (key, val) {
            var data = {};
            data.value = val.Value;
            data.text = val.Name;
            datavalue.push(data);
        });
        pc.projectList(datavalue);

        setTimeout(function () {
            pc.populateTurbine(selected);
        }, 100);
    }
};


pc.showHidePeriod = function (callback) {
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

pc.showHidePeriod2 = function (callback) {
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


pc.InitDefaultValue = function () {
    // $("#projectList1").data("kendoDropDownList").value("Tejuva (24 | 50.4 MWh)")
    // $("#projectList2").data("kendoDropDownList").value("Tejuva (24 | 50.4 MWh)")
    $("#periodList").data("kendoDropDownList").value("custom");
    $("#periodList").data("kendoDropDownList").trigger("change");

    $("#periodList2").data("kendoDropDownList").value("custom");
    $("#periodList2").data("kendoDropDownList").trigger("change");

    var maxDateData = new Date(app.getUTCDate(app.currentDateData));
    var lastStartDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate()-30, 0, 0, 0, 0));
    var lastEndDate = new Date(app.getDateMax(maxDateData));
    // var dateEnd2 = new Date(Date.UTC(moment(lastStartDate).get('year'), lastStartDate.getMonth(), lastStartDate.getDate()-30, 0, 0, 0, 0));
    // var dateStart2 =new Date(Date.UTC(moment(dateEnd2).get('year'), dateEnd2.getMonth(), dateEnd2.getDate()-30, 0, 0, 0, 0));

    $('#dateEnd').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart').data('kendoDatePicker').value(lastStartDate);
    $('#dateEnd2').data('kendoDatePicker').value(lastEndDate);
    $('#dateStart2').data('kendoDatePicker').value(lastStartDate);
}
pc.initChart = function() {
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
            pc.getAvailDate();
            
            var link = "analyticpowercurve/getlistpowercurvecomparison"

            app.loading(true);
            var param = {
                PC1Period       : $('#periodList').data('kendoDropDownList').value(),
                PC1Project      :  $("#projectList1").data("kendoDropDownList").value(),
                PC1Turbine      :  $("#turbineList1").data('kendoDropDownList').value(),// == "All Turbine" || $("#turbineList1").data('kendoDropDownList').value() == undefined ? pc.turbine() : $("#turbineList1").data('kendoDropDownList').value(),
                PC1DateStart    : p1DateStart,
                PC1DateEnd      : p1DateEnd,

                PC2Period       : $('#periodList2').data('kendoDropDownList').value(),
                PC2Project      :  $("#projectList1").data("kendoDropDownList").value(),
                PC2Turbine      : $("#turbineList2").data('kendoDropDownList').value(),// == "All Turbine" || $("#turbineList2").data('kendoDropDownList').value() == undefined  ? pc.turbine() : $("#turbineList2").data('kendoDropDownList').value(),
                PC2DateStart    : p2DateStart,
                PC2DateEnd      : p2DateEnd

            };

            toolkit.ajaxPost(viewModel.appName + link, param, function(res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                
                var dataTurbine = res.data.Data;

                
                $('#chartPCcomparison').html("");
                $("#chartPCcomparison").kendoChart({
                    pdf: {
                      fileName: "DetailPowerCurve.pdf",
                    },
                    theme: "flat",
                    title: {
                        text: "Power Curves",
                        visible: false,
                        font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
                    },
                    legend: {
                        position: "bottom",
                        visible: true,
                        labels: {
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                        },
                    },
                    chartArea: {
                        height: 375,
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
                            step: 1,
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
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
                    tooltip: {
                        visible: true,
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
                    pannable: false,
                    zoomable: false
                });
                app.loading(false);
                $("#chartPCcomparison").data("kendoChart").refresh();

                
            });
        }
}

pc.setProjectTurbine = function(projects, turbines, selected){
	pc.rawproject(projects);
    pc.rawturbine(turbines);
	pc.populateProject(selected);
};

$(document).ready(function () {
    $('#btnRefresh').on('click', function() {
        setTimeout(function() {
            pc.initChart();
        }, 300);
    });

    $('#projectList1').kendoDropDownList({
        change: function () { 
            var project = $('#projectList1').data("kendoDropDownList").value();
            pc.getAvailDate();
            pc.populateTurbine(project);
         }
    });

    app.loading(true);
    pc.InitDefaultValue();
    setTimeout(function() {
        pc.initChart();
    }, 200);
});
