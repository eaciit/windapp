viewModel.AnalyticKpi = new Object();
var page = viewModel.AnalyticKpi;

var keys = [
    { "value": "Production", "text": "Production", "status": true },
    { "value": "PLF", "text": "PLF", "status": false },
    { "value": "Revenue", "text": "Revenue", "status": false },
    { "value": "TotalAvailability", "text": "Total Availability", "status": false },
    { "value": "GridAvailability", "text": "Grid Availability", "status": false },
    { "value": "MachineAvailability", "text": "Machine Availability", "status": false },
    { "value": "DataAvailability", "text": "Data Availability", "status": false },
    { "value": "P50Generation", "text": "P50 Generation", "status": false },
    { "value": "P50PLF", "text": "P50 PLF", "status": false },
    { "value": "P75Generation", "text": "P75 Generation", "status": false },
    { "value": "P75PLF", "text": "P75 PLF", "status": false },
    { "value": "P90Generation", "text": "P90 Generation", "status": false },
    { "value": "P90PLF", "text": "P90 PLF", "status": false },
];

page.columnsBreakdownList = ko.observableArray([
    { "text": "Daily", "value": "Daily" },
    { "text": "Monthly", "value": "Monthly" },
    { "text": "Yearly", "value": "Yearly" },
]);

page.rowsBreakdownList = ko.observableArray([
    { "text": "Project", "value": "Project" },
    { "text": "Turbine", "value": "Turbine" },
]);



page.projectList = ko.observableArray();
page.columnsBreakdown = ko.observable();
page.rowsBreakdown = ko.observable();
page.project = ko.observable();
page.headerColumns = ko.observableArray();
page.dataSource = ko.observableArray();
page.views = ko.observableArray([]);
page.viewList = ko.observableArray([]);
page.selectedView = ko.observable();

page.key1 = ko.observableArray([]);
page.key2 = ko.observableArray([]);
page.key3 = ko.observableArray([]);

page.key1(keys);
page.key2(keys);
page.key3(keys);

var isFirst = true;

var Data = {
    LoadData: function () {
        var isValid = fa.LoadData();
        if(isValid) {
            app.loading(true);

            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = $('#dateEnd').data('kendoDatePicker').value();   

            var columnBreakdown = $('#columnsBreakdown').val();
            var rowBreakdown = $('#rowsBreakdown').val();

            var keyA = $('#key1').val();
            var keyB = $('#key2').val();
            var keyC = $('#key3').val();

            var param = {
                period: fa.period,
                dateStart: new Date(moment(dateStart).format('YYYY-MM-DD')),
                dateEnd: new Date(moment(dateEnd).format('YYYY-MM-DD')),
                turbine: fa.turbine(),
                project: fa.project,
                columnBreakDown: columnBreakdown,
                rowBreakDown: rowBreakdown,
                keyA: keyA,
                keyB: keyB,
                keyC: keyC,
            };

            toolkit.ajaxPost(viewModel.appName + "analytickpi/getscadasummarylist", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                page.dataSource(res.data.Data);
                page.generateGrid();
            });
            fa.getDataAvailability();
        }
        // app.loading(false);
    }
}

page.getPdfGrid = function(){
    $("#gridKpiAnalysis").getKendoGrid().saveAsExcel();
     return false;
}

page.generateGrid = function () {

    var project = $("#projectList").data("kendoDropDownList").value();
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
    var title = project+"KPITable"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".xlsx";

    var config = {
        dataSource: {
            data: page.dataSource(),
            pageSize: 10,
            sort: ({ field: "Row", dir: "asc" }),
            aggregate : []
        },
        excel:{
            fileName:title,
            allPages:true, 
            filterable:true
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Breakdown", field: "Row", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false,footerTemplate: "Total",footerAttributes: {style: "text-align:center;"}},
        ],
    };

    var units = page.dataSource()[0].Unit;

    $.each(page.dataSource()[0].Column, function (i, val) {

        var column = {
            title: val.Name,
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            columns: []
        }

        

        var a = 0;
        var k = 3;
        var keyIndex = ["A", "B", "C"];

        for (var key in val) {

            var hidden = false;


            if (key == 'Name' || key == 'YearMonth' || key == 'TitleKeyA' || key == 'TitleKeyB' || key == 'TitleKeyC') {
                hidden = true;
            } else {
                hidden = false;
            }


            var unit = units[a] == "" ? units[a + 1] : units[a];

            if(unit == "%"){
                config.dataSource.aggregate.push({ field: "Column[" + i + "]." + key, aggregate: "average" },)
            }else{
                config.dataSource.aggregate.push({ field: "Column[" + i + "]." + key, aggregate: "sum" },)
            }

            var colChild = {
                title: key.replace(/([A-Z])/g, ' $1').trim() + "(" + unit + ")",
                encoded: true,
                hidden: hidden,
                field: "Column[" + i + "]." + key,
                width: 90,
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;white-space: normal !important',
                    class: "align-center column-" + key,
                },
                attributes: {
                    class: "align-center column-" + key
                },
                format: "{0:n2}",
                filterable: false,
                footerTemplate: (unit == '%') ? "#=kendo.toString(average, 'n2')#" : "#=kendo.toString(sum, 'n2')#",
                footerAttributes: {style: "text-align:center;"}

            };
            column.columns.push(colChild);


            a++;
            k++;


        }


        config.columns.push(column);
        
    });


    app.loading(false);
    $('#gridKpiAnalysis').html("");
    $('#gridKpiAnalysis').kendoGrid(config);
    $('#gridKpiAnalysis').data('kendoGrid').refresh();
}

page.checkKeyFromSavedView = function () {
    var key1 = page.selectedView.KeyA;
    var key2 = page.selectedView.KeyB;
    var key3 = page.selectedView.KeyC;

    page.key1([]);
    page.key2([]);
    page.key3([]);

    $.each(keys, function (i) {
        if (keys[i].value != key2 || keys[i].value != key3) {
            page.key1.push(keys[i]);
        }
        if (keys[i].value != key1 || keys[i].value != key3) {
            page.key2.push(keys[i]);
        }
        if (keys[i].value != key1 || keys[i].value != key2) {
            page.key3.push(keys[i]);
        }
    });

    page.key2.unshift({ "value": "None", "text": "None", "status": true });
    page.key3.unshift({ "value": "None", "text": "None", "status": true });

    setTimeout(function () {
        $("#key1").data("kendoDropDownList").value(key1);
        $("#key2").data("kendoDropDownList").value(key2);
        $("#key3").data("kendoDropDownList").value(key3);

        $("#columnsBreakdown").data("kendoDropDownList").value(page.selectedView.ColumnBreakdown);
        $("#rowsBreakdown").data("kendoDropDownList").value(page.selectedView.RowBreakdown);
        page.checkKey();
        if (isFirst == false) {
            Data.LoadData();
        }
    }, 50);
}

page.checkKey = function () {
    var key1 = $("#key1").data("kendoDropDownList").value();
    var key2 = $("#key2").data("kendoDropDownList").value();
    var key3 = $("#key3").data("kendoDropDownList").value();

    page.key1([]);
    page.key2([]);
    page.key3([]);

    $.each(keys, function (i) {
        if (keys[i].value == key2 || keys[i].value == key3) {
            return true;
        }
        page.key1.push(keys[i]);
    });

    $.each(keys, function (i) {
        if (keys[i].value == key1 || keys[i].value == key3) {
            return true;
        }
        page.key2.push(keys[i]);
    });
    $.each(keys, function (i) {
        if (keys[i].value == key1 || keys[i].value == key2) {
            return true;
        }
        page.key3.push(keys[i]);
    });

    page.key2.unshift({ "value": "None", "text": "None", "status": true });
    page.key3.unshift({ "value": "None", "text": "None", "status": true });

    $("#key1").data("kendoDropDownList").value(key1);
    $("#key2").data("kendoDropDownList").value(key2);
    $("#key3").data("kendoDropDownList").value(key3);

    if (isFirst == false) {
        Data.LoadData();
    }
}

page.loadView = function () {
    var selectedVal = $("#savedViews").data("kendoDropDownList").value();
    if (selectedVal != "") {
        page.selectedView = null;
        $.each(page.views(), function (i, val) {
            if (val.Name == selectedVal) {
                page.selectedView = val;
                page.checkKeyFromSavedView();
            }
        });
    }
}

page.getViews = function () {
    page.viewList = [];
    page.viewList.push({
        value: "",
        text: "Please Select"
    })

    app.ajaxPost(viewModel.appName + "userpreferences/getsavedkpianalysis", "", function (res) {
        if (!app.isFine(res)) {
            return;
        }

        page.views(res.data);
        $.each(page.views(), function (i, val) {
            page.viewList.push({
                value: val.Name,
                text: val.Name
            })
        });

        $("#savedViews").data("kendoDropDownList").dataSource.data(page.viewList);
        $("#savedViews").data("kendoDropDownList").dataSource.query();
        if ($("#savedViews").data("kendoDropDownList").value() == "") {
            $("#savedViews").data("kendoDropDownList").select(0);
        }
    });
}

page.modalSaveView = function () {
    var selectedVal = $("#savedViews").data("kendoDropDownList").value();
    $("#inputViewName").val(selectedVal);
    // $("#inputKeyA").html($("#key1").data("kendoDropDownList").text());
    // $("#inputKeyB").html($("#key2").data("kendoDropDownList").text());
    // $("#inputKeyC").html($("#key3").data("kendoDropDownList").text());
    // $("#inputColBreakdown").html($("#columnsBreakdown").data("kendoDropDownList").text());
    // $("#inputRowBreakdown").html($("#rowsBreakdown").data("kendoDropDownList").text());

    if (page.viewList.length == 4 && selectedVal == "") {
        toolkit.showError("Only 3 views are allowed");
    } else if (selectedVal != "") {
        $("#modal-views-title").html("Update View");
        page.ShowModal('modalForm', 'show');
    } else {
        $("#modal-views-title").html("Create New View");
        page.ShowModal('modalForm', 'show');
    }
}

page.saveView = function () {
    page.ShowModal('modalForm', 'hide');

    var param = {
        OldName: page.selectedView.Name,
        Name: $("#inputViewName").val(),
        KeyA: $("#key1").val(),
        KeyB: $("#key2").val(),
        KeyC: $("#key3").val(),
        ColumnBreakdown: $("#columnsBreakdown").val(),
        RowBreakdown: $("#rowsBreakdown").val()
    }

    app.ajaxPost(viewModel.appName + "userpreferences/savekpi", param, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data == null) {
            toolkit.showError("Error Occur when save the KPI");
            return;
        }

        page.viewList = [];
        page.viewList.push({
            value: "",
            text: "Please Select"
        })

        swal({
            title: "Info",
            type: "success",
            text: "Data Successfully Saved",
        }, function () { });

        // page.views(res.data);

        $.each(page.views(), function (i, val) {
            page.viewList.push({
                value: val.Name,
                text: val.Name
            })
        });

        var idx = $("#savedViews").data("kendoDropDownList").select();

        $("#savedViews").data("kendoDropDownList").dataSource.data(page.viewList);
        $("#savedViews").data("kendoDropDownList").dataSource.query();

        page.getViews();

        setTimeout(function () {
            $("#savedViews").data("kendoDropDownList").select(idx);
        }, 100);
    });
}

page.ShowModal = function (modalId, showhide) {
    if (showhide == 'show') {
        $('#' + modalId).appendTo("body").modal({
            backdrop: 'static',
            keyboard: false,
            show: showhide
        });
    } else {
        $('#' + modalId).modal('hide');
    }
}

page.setBreakDown = function () {
    fa.disableRefreshButton(true);
    page.columnsBreakdownList = [];
    page.rowsBreakdownList = [];
    setTimeout(function () {
        var project = $('#projectList').data("kendoDropDownList").value();
        fa.populateTurbine(project);

        $.each(fa.GetBreakDown(), function (i, val) {
            if (val.value == "Turbine" || val.value == "Project") {
                // page.rowBreakdown = val.value
                page.rowsBreakdownList.push(val);
            } else {
                page.columnsBreakdownList.push(val);
            }
        });

        $("#columnsBreakdown").data("kendoDropDownList").dataSource.data(page.columnsBreakdownList);
        $("#columnsBreakdown").data("kendoDropDownList").dataSource.query();
        $("#columnsBreakdown").data("kendoDropDownList").select(0);

        $("#rowsBreakdown").data("kendoDropDownList").dataSource.data(page.rowsBreakdownList);
        $("#rowsBreakdown").data("kendoDropDownList").dataSource.query();
        
        if (project == "") {
            $("#rowsBreakdown").data("kendoDropDownList").value("Project");
        } else {
            $("#rowsBreakdown").data("kendoDropDownList").value("Turbine");
        }
        fa.disableRefreshButton(false);
    }, 500);
}

vm.currentMenu('KPI Table');
vm.currentTitle('KPI Table');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'KPI Table', href: viewModel.appName + 'page/analytickpi' }]);

$(function () {
    di.getAvailDate();
    page.columnsBreakdown('Daily');
    page.rowsBreakdown('Turbine');

    $("#key1").data("kendoDropDownList").value("Production");
    $("#key2").data("kendoDropDownList").value("PLF");
    $("#key3").data("kendoDropDownList").value("Revenue");

    $('#savedViews').kendoDropDownList({
        data: [],
        dataValueField: 'value',
        dataTextField: 'text',
        change: function () { page.loadView() },
    });

    // page.checkKey();
    setTimeout(function () {
        isFirst = true;

        app.loading(true);
        Data.LoadData();
        page.checkKey();
        isFirst = false;
        page.getViews();
    }, 1500);

    $('#btnRefresh').on('click', function () {
        // page.columnsBreakdownList = fa.GetBreakDown();
        fa.setPreviousFilter();
        fa.checkTurbine();
        Data.LoadData();
    });

    $('#btnSaveView').on('click', function () {
        page.modalSaveView();
    });


    // smart filter :)

    $('#periodList').kendoDropDownList({
        data: fa.periodList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () { 
            if(isFirst !== true){
                fa.currentFilter().period = this._old;
                fa.checkFilter();
            }
            fa.showHidePeriod(page.setBreakDown());
            // fa.checkFilter();
        }
    });

    setTimeout(function () {
        $('#projectList').kendoDropDownList({
            data: fa.projectList,
            dataValueField: 'value',
            dataTextField: 'text',
            suggest: true,
            change: function () { 
                fa.currentFilter().project = this._old;
                page.setBreakDown(); 
                setTimeout(function(){
                   fa.setTurbine();
                   di.getAvailDate();
                   fa.checkFilter();
                },500);
            }
        });

        $("#dateStart").change(function () { 
            fa.DateChange(page.setBreakDown());
            fa.currentFilter().startDate = this.value;
            fa.checkFilter(); 

        });
        $("#dateEnd").change(function () { 
            fa.DateChange(page.setBreakDown()); 
            fa.currentFilter().endDate = this.value;
            fa.checkFilter(); 
        });

    }, 1500);
});


