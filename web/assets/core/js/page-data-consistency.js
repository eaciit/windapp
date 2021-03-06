'use strict';

viewModel.AvailabilityAnalysis = new Object();
var pg = viewModel.AvailabilityAnalysis;
var maxdate = new Date(Date.UTC(2016, 5, 30, 23, 59, 59, 0));

pg.type = ko.observable();
pg.detailDTTopTxt = ko.observable();
pg.isDetailDTTop = ko.observable(false);
pg.periodDesc = ko.observable();

pg.breakDownList = ko.observableArray([
    { "value": "dateinfo.dateid", "text": "Date" },
    { "value": "dateinfo.monthdesc", "text": "Month" },
    { "value": "dateinfo.year", "text": "Year" },
    { "value": "projectname", "text": "Project" },
    { "value": "turbine", "text": "Turbine" },
]);

pg.isFirstDataCon = ko.observable(true);
pg.isFirstVarience = ko.observable(true);
pg.isFirstVarienceOver = ko.observable(true);


pg.getAvailDate = function(){
    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var namaproject;
        
        var projectVal = $("#projectList").data("kendoDropDownList").value();
        if( projectVal == undefined || projectVal == "") {
            namaproject = "Tejuva";
        }else{
            namaproject= projectVal;
        }

        if(res.data[namaproject]["ScadaData"] != null) {
            $('#availabledatestartscada').html(kendo.toString(moment.utc(new Date(res.data[namaproject]["ScadaData"][0])).format('DD-MMM-YYYY')));
            $('#availabledateendscada').html(kendo.toString(moment.utc(new Date(res.data[namaproject]["ScadaData"][1])).format('DD-MMM-YYYY')));
        }
        // if(res.data[namaproject]["ScadaDataHFD"] != null) {
        //     $('#availabledatestartscadahfd').html(kendo.toString(moment.utc(new Date(res.data[namaproject]["ScadaDataHFD"][0])).format('DD-MMM-YYYY')));
        //     $('#availabledateendscadahfd').html(kendo.toString(moment.utc(new Date(res.data[namaproject]["ScadaDataHFD"][1])).format('DD-MMM-YYYY')));
        // }
        if(res.data[namaproject]["DGRData"] != null) {
            $('#availabledatestartdgr').html(kendo.toString(moment.utc(res.data[namaproject]["DGRData"][0]).format('DD-MMM-YYYY')));
            $('#availabledateendsdgr').html(kendo.toString(moment.utc(res.data[namaproject]["DGRData"][1]).format('DD-MMM-YYYY')));
        }


        var maxDateData = new Date(res.data[namaproject]["ScadaData"][1]);

        var startDate = new Date(Date.UTC(moment(maxDateData).get('year'), maxDateData.getMonth(), maxDateData.getDate() - 7, 0, 0, 0, 0));

        // $('#dateStart').data('kendoDatePicker').value(startDate);
        // $('#dateEnd').data('kendoDatePicker').value(kendo.toString(moment.utc(res.data[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY')));
    });
}
pg.loadData = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        app.loading(true);
        
        if (fa.project == "") {
            pg.type = "Project Name";
        } else {
            pg.type = "Turbine";
        }
    }
}

pg.DataCon = function(){
    app.loading(true);
    // if(pg.isFirstDataCon() === true){

            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 
            var project = $('#projectList').data('kendoDropDownList').value();

            var param = {
                period: fa.period,
                Turbine: fa.turbine(),
                DateStart: dateStart,
                DateEnd: dateEnd,
                Project: project
            };

            var dataAvailDGR = "";
            var dataAvailScada = "";
            $("#gridSummaryDgrScada").kendoGrid({
                theme: "flat",
                /*dataSource: {
                    data : dataSource,
                }*/
                dataSource: {
                    serverPaging: false,
                    serverSorting: false,
                    serverFiltering: false,
                    transport: {
                        read: {
                            url: viewModel.appName + "analyticdgrscada/getdatarev",
                            type: "POST",
                            data: param,
                            dataType: "json",
                            contentType: "application/json; charset=utf-8"
                        },
                        parameterMap: function (options) {
                            return JSON.stringify(options);
                        }
                    },
                    schema: {
                        model: {
                            fields: {
                                AlarmOkTime: { type: "number" },
                                OkTime: { type: "number" },
                                Power: { type: "number" },
                                PowerLost: { type: "number" },
                            }
                        },
                        data: function (res) {
                            app.loading(false);
                            if (!app.isFine(res)) {
                                return;
                            }
                            
                            var availDate = res.data.availdate;

                            var dateStartDGR = (availDate["dgr"] == "N/A") ?  "N/A" : kendo.toString(moment.utc(new Date(availDate["dgr"][0])).format('DD-MMM-YYYY'));
                            var dateEndDGR = (availDate["dgr"] == "N/A") ?  "N/A" : kendo.toString(moment.utc(availDate["dgr"][1]).format('DD-MMM-YYYY'));

                            dataAvailDGR =  dateStartDGR+ " to " +dateEndDGR ;
                            dataAvailScada = dataAvailDGR = kendo.toString(moment.utc(new Date(availDate["scada"][0])).format('DD-MMM-YYYY')) + " to " + kendo.toString(moment.utc(new Date(availDate["scada"][1])).format('DD-MMM-YYYY'));

                             $("#dateDgr").html(dateStartDGR+ " to " +dateEndDGR);
                             $("#dateScada").html(dataAvailScada);

                            return res.data.data
                        }
                    },
                },
                columns: [
                    { title: " ", field: "desc", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-left" }, width: 150 },
                    // { 
                    //     headerTemplate: "<span>DGR</span><br><span id='dateDgr'></span>",
                    //     width: 120, 
                    //     field: "dgr", 
                    //     headerAttributes: { style: "text-align: center" }, 
                    //     attributes: { class: "align-center" }, 
                    //     template: "#if(desc== 'PLF'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc== 'Grid Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc== 'Machine Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc=='True Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else {# #: kendo.toString(dgr, 'N2') # #}#" },
                    { 
                        headerTemplate: "<span>DGR</span><br><span id='dateDgr'></span>",
                        width: 120, 
                        field: "dgr", 
                        headerAttributes: { style: "text-align: center" }, 
                        attributes: { class: "align-center" }, 
                        template: "#if(desc== 'PLF'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc== 'Grid Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc== 'Machine Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else if(desc=='True Availability'){# #: '-' # #if(dgr!= 'N/A'){#  #}}else if(desc== 'Data Availability'){# #: kendo.toString(dgr, 'N2') # #if(dgr!= 'N/A'){# % #}}else {# #: kendo.toString(dgr, 'N2') # #}#" },
                    {  headerTemplate: "<span>Scada</span><br><span id='dateScada'></span>", width: 120,field: "scada", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(scada, 'N2') # #if(scada!= 'N/A'){# % #}}else if(desc== 'Grid Availability'){# #: kendo.toString(scada, 'N2') # #if(scada!= 'N/A'){# % #}}else if(desc== 'Machine Availability'){# #: kendo.toString(scada, 'N2') # #if(scada!= 'N/A'){# % #}}else if(desc=='True Availability'){# #: kendo.toString(scada, 'N2') # #if(scada!= 'N/A'){# % #}}else if(desc=='Data Availability'){# #: kendo.toString(scada, 'N2') # #if(scada!= 'N/A'){# % #}}else {# #: kendo.toString(scada, 'N2') # #}#" },
                    // { title: "Realtime", width: 120,field: "ScadaHFD", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(ScadaHFD, 'N2') # #if(ScadaHFD!= 'N/A'){# % #}}else if(desc== 'Grid Availability'){# #: kendo.toString(ScadaHFD, 'N2') # #if(ScadaHFD!= 'N/A'){# % #}}else if(desc== 'Machine Availability'){# #: kendo.toString(ScadaHFD, 'N2') # #if(ScadaHFD!= 'N/A'){# % #}}else if(desc=='True Availability'){# #: kendo.toString(ScadaHFD, 'N2') # #if(ScadaHFD!= 'N/A'){# % #}}else {# #: kendo.toString(ScadaHFD, 'N2') # #}#" },
                    // {
                    //     title: "Difference",
                    //     headerAttributes: { style: 'font-weight: bold; text-align: center;' },
                    //     columns: [
                    //         { title: "DGR vs Scada", width: 120,field: "difference", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(difference, 'N2') # % #}else {# #: kendo.toString(difference, 'N2') # #}#" },
                    //         { title: "DGR vs Scada HFD", width: 120,field: "diffdgrhfd", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(diffdgrhfd, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(diffdgrhfd, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(diffdgrhfd, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(diffdgrhfd, 'N2') # % #}else {# #: kendo.toString(diffdgrhfd, 'N2') # #}#" },
                    //      ]
                    // },
                    // { title: "DGR vs Scada", width: 120,field: "difference", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc=='True Availability'){# #: kendo.toString(difference, 'N2') # % #}else {# #: kendo.toString(difference, 'N2') # #}#" },
                    { title: "DGR vs Scada", width: 120,field: "difference", headerAttributes: { style: "text-align: center" }, attributes: { class: "align-center" }, template: "#if(desc== 'PLF'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Grid Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc== 'Machine Availability'){# #: kendo.toString(difference, 'N2') # % #}else if(desc=='True Availability'){# #: '-' # #}else {# #: kendo.toString(difference, 'N2') # #}#" },

                ],
                dataBound: function(){
                    app.loading(false);
                    pg.isFirstDataCon(false);
                }
            });
    // }
    // else{
    //      setTimeout(function(){
    //         $("#gridSummaryDgrScada").data("kendoGrid").refresh();
    //         app.loading(false);
    //     },200);
    // }
}
pg.Variance = function(){

}
pg.VarianceOver = function(){
    
}
pg.resetStatus = function(){
    pg.isFirstDataCon(true);
    pg.isFirstVarience(true);
    pg.isFirstVarienceOver(true);
}

vm.currentMenu('Data Consistency');
vm.currentTitle('Data Consistency');
vm.breadcrumb([{ title: "KPI's", href: '#' }, { title: 'Data Consistency', href: viewModel.appName + 'page/analyticdataconsistency' }]);

$(document).ready(function () {
    pg.getAvailDate();
    fa.LoadData();

    setTimeout(function(){
        pg.loadData();
        pg.DataCon();
    },500);

    $('#btnRefresh').on('click', function () {
        fa.checkTurbine();
        pg.resetStatus();

        if (fa.LoadData()) {
            pg.DataCon();
        }
        // $('.nav').find('li.active').find('a').trigger( "click" );
    });

    $('#projectList').kendoDropDownList({
        change: function () {  
            pg.getAvailDate();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.populateTurbine(project);
        }
    });
});
