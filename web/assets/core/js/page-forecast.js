'use strict';

viewModel.Forecasting = new Object();
var pg = viewModel.Forecasting;
var heightSub = 200;

vm.currentMenu('Forecasting & Scheduling');
vm.currentTitle('Forecasting & Scheduling');
vm.breadcrumb([{ title: 'Forecasting & Scheduling', href: viewModel.appName + 'page/forecasting' }]);

pg.DataSource = ko.observableArray([]);
pg.CurrentTab = ko.observable('grid');
pg.TurbineDown = ko.observable(0);
pg.TurbineDownData = ko.observableArray([]);
pg.EditMode = ko.observable('saved');

pg.getData = function() {
    app.loading(true);
    var url = viewModel.appName + 'forecast/getlist';
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD')); 
    var param = {
        period: fa.period,
        dateStart: dateStart,
        dateEnd: dateEnd,
        turbine: fa.turbine(),
        project: fa.project,
    };
    var getdata = toolkit.ajaxPostDeffered(url, param, function(res) {});
    $.when(getdata).done(function(d){
        pg.DataSource(d.data);
        pg.TurbineDown(d.data[0].TurbineDown);
        if(pg.CurrentTab()=='grid')
            pg.genereateGrid();
        else
            pg.genereateChart();
        app.loading(false);
    });
}


pg.getPdfGrid = function(){
    // app.loading(true);
    $("#gridForecasting").data('kendoGrid').saveAsExcel();
    // app.loading(false);

    return false;
}

kendo.ui.Grid.fn.editCell = (function (editCell) {
    return function (cell) {
        cell = $(cell);

        var that = this,
            column = that.columns[that.cellIndex(cell)],
            model = that._modelForContainer(cell),
            event = {
                container: cell,
                model: model,
                preventDefault: function () {
                    this.isDefaultPrevented = true;
                }
            };

        if (model && typeof this.options.beforeEdit === "function") {
            this.options.beforeEdit.call(this, event);
            
            // don't edit if prevented in beforeEdit
            if (event.isDefaultPrevented) return;
        }
        
        editCell.call(this, cell);
    };
})(kendo.ui.Grid.fn.editCell);

pg.isExport = false;
pg.genereateGrid = function(){
    app.loading(true);
    setTimeout(function(){ 
        var project = $("#projectList").data("kendoDropDownList").value();
        var dateStart = $('#dateStart').data('kendoDatePicker').value();
        var dateEnd = $('#dateEnd').data('kendoDatePicker').value();  
        var title = project+"Forecasting"+kendo.toString(dateStart, "dd/MM/yyyy")+"to"+kendo.toString(dateEnd, "dd/MM/yyyy")+".xlsx";

        $("#gridForecasting").html('');
        $("#gridForecasting").kendoGrid({
            toolbar: ["save"],
            dataSource: new kendo.data.DataSource({
                data: pg.DataSource(),
                pageSize: 15,
                schema: {
                    model: {
                        id: "ID",
                        fields: {
                            ID: { editable: false },
                            Date: { editable: false },
                            TimeBlock: { editable: false },
                            AvaCap: { editable: false },
                            Forecast: { editable: false },
                            Actual: { editable: false },
                            ExpProd: { editable: false },
                            FcastWs: { editable: false },
                            DevFcast: { editable: false },
                            ActualWs: { editable: false },
                            Deviation: { editable: false },
                            DevSchAct: { editable: false },
                            DSMPenalty: { editable: false },
                            SchFcast: { type: "number" },
                        }
                    }
                },
                batch: true,
            }),
            height: 520, //$('body').height() - heightSub + 30,
            // scrollable: true,
            sortable: true,
            filterable: false,
            pageable: {
                input: true,
                numeric: false
            },
            cellClose:  function(e) {
                console.log('close');
                console.log(e.type);
            },
            saveChanges: function(e) {
                swal({
                    title: 'Save all changes now?',
                    text: "You won't be able to revert this!",
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#3085d6',
                    cancelButtonColor: '#d33',
                    confirmButtonText: 'Yes, save it!'
                }, function() {
                    var data = $("#gridForecasting").data("kendoGrid").dataSource.data();
                    var dirty = $.grep(data, function(item) {
                        return item.dirty
                    });
                    var updatedData = [];
                    console.log(dirty);
                    $.each(dirty, function(i, v){
                        var item = {
                            id: v.ID,
                            value: v.SchFcast,
                        };
                        updatedData.push(item);
                    });
                    var param = {
                        project: fa.project,
                        values: updatedData,
                    };
                    app.loading(true);
                    var url = viewModel.appName + 'forecast/updatesldc';
                    toolkit.ajaxPost(url, param, function(res) {
                        app.loading(false);
                    });
                    return;
                });
                e.preventDefault();
            },
            excel:{
                fileName:title,
                allPages:true, 
                filterable:true
            },
            excelExport: function(e) {
                app.loading(true);
                    //console.log(e);
                    // make all column is wrap : true
                    var ecols = e.workbook.sheets[0].columns;
                    $.each(ecols, function(i, c){
                        if(i>1 && (c.width==NaN || c.width==null || c.autoWidth) && i<ecols.length - 1) {
                            e.workbook.sheets[0].columns[i].autoWidth = false;
                            e.workbook.sheets[0].columns[i].width = 80;
                        }
                        e.workbook.sheets[0].columns[i].wrap = true;
                    });

                    var sheet = e.workbook.sheets[0];

                    for (var rowIndex = 0; rowIndex < sheet.rows.length; rowIndex++) {
                        var row = sheet.rows[rowIndex];
                        if(rowIndex == 0){
                             for (var idx = 0; idx < row.cells.length; idx ++) {
                                 var title = row.cells[idx].value
                                 row.cells[idx].value = title.replace(/<br>/g,'\r');
                             }
                        }else{
                            for (var cellIndex = 0; cellIndex < row.cells.length; cellIndex ++) {
                                if(cellIndex>1) {
                                    if(cellIndex==2) {
                                         row.cells[cellIndex].format = '#,##0';
                                    }
                                    if(cellIndex==9 || cellIndex==10) {
                                         row.cells[cellIndex].format = '#,##0.0#%';
                                    } else {
                                         row.cells[cellIndex].format = '#,##0.0#';
                                    }
                                }
                            }
                        }
                    }
                app.loading(false);
            },
            columns: [
                { field: "Date", title: "Date"},
                { field: "TimeBlock", title: "Time Block", width: 100, },
                // { field: "TimeBlockInt", title: "Time<br/>Block", width: 50, },
                { field: "AvaCap", title: "Avg. Cap.<br>(MW)", template : "#: (AvaCap==null?'-':kendo.toString(AvaCap, 'n0')) #", format: '{0:n0}' },
                { field: "Forecast", title: "Forecast<br>(MW)", template : "#: (Forecast==null?'-':kendo.toString(Forecast, 'n2')) #", format: '{0:n2}'},
                { title: "Sch Fcast /<br>SLDC (MW)", field: "SchFcast", template : "#: (SchFcast==null?'-':kendo.toString(SchFcast, 'n2')) #", format: '{0:n2}' },
                { field: "Actual", title: "Actual Prod<br>(MW)", template : "#: (Actual==null?'-':kendo.toString(Actual, 'n2')) #", format: '{0:n2}' },
                { title: "Exp Prod /<br>Pwr Curv (MW)", width: 120, field: "ExpProd", template : "#: (ExpProd==null?'-':kendo.toString(ExpProd, 'n2')) #", format: '{0:n2}' },
                { field: "FcastWs", title: "Fcast ws<br>(m/s)", template : "#: (FcastWs==null?'-':kendo.toString(FcastWs, 'n2')) #", format: '{0:n2}' },
                { field: "ActualWs", title: "Actual ws<br>(m/s)", template : "#: (ActualWs==null?'-':kendo.toString(ActualWs, 'n2')) #", format: '{0:n2}' },
                { field: "DevFcast", title: "% Error<br>Act / Fcst", template : "#: (DevFcast==null?'-':kendo.toString(DevFcast, 'p2')) #", format: '#,##0.0#%' },
                { field: "DevSchAct", title: "% Error<br>Act / Schd", template : "#: (DevSchAct==null?'-':kendo.toString(DevSchAct, 'p2')) #", format: '{0:p2}' },
                { field: "Deviation", title: "Deviation<br>(MW)", template : "#: (Deviation==null?'-':kendo.toString(Deviation, 'n2')) #", format: '{0:n2}' },
                { field: "DSMPenalty", title: "DSM Penalty"},
            ],
            editable: true,
            beforeEdit: function(e) {
                var data = e.model;
                var schFcast = data.SchFcast;
                var timeStamp = data.TimeStamp;
                var timefilter = moment.utc().add(6.5, 'hour');
                var isAllowed = (schFcast != null && moment(timeStamp).isAfter(timefilter));
                if(!isAllowed) {
                    e.preventDefault();
                } 
            },
        });
        $("#gridForecasting").data("kendoGrid").refresh();
        app.loading(false);
    }, 300);
}

pg.genereateChart = function(){
    app.loading(true);
    var date1 = $('#dateStart').data('kendoDatePicker').value();
    var date2 = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));
    var timeDiff = Math.abs(date2.getTime() - date1.getTime());
    var diffDays = Math.ceil(timeDiff / (1000 * 3600 * 24)); 
    var mindays = 2;
    setTimeout(function(){
        $("#chartForecasting").html("");
        $("#chartForecasting").kendoChart({
            dataSource: {
                data: pg.DataSource(),
            },
            title: {
                text: "Forecasting and Scheduling",
                font: '18px bold Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                margin: {
                    top: 0,
                }
            },
            chartArea : {
                height : 500, // $('body').height() - heightSub + 30,
                background: "transparent",
            },
            legend: {
                visible: true,
                position : "top",
                offsetY: -10,
                item: {
                    visual: function (e) {
                        var color = e.options.markers.background;
                        var dashType = e.series.dashType;
                        var labelColor = e.options.labels.color;
                        var rect = new kendo.geometry.Rect([0, 0], [200, 50]);
                        var layout = new kendo.drawing.Layout(rect, {
                          spacing: 5,
                          alignItems: "center"
                        });
              
                        var svgPath = "M0 5.5 L 28 5.5 28 8 0 8Z";
                        if(dashType!="solid") {
                            svgPath = "M0 5.5 L 8 5.5 8 8 0 8Z M11 5.5 L 19 5.5 19 8 11 8Z M22 5.5 L 30 5.5 30 8 22 8Z";
                        }
                        var path = kendo.drawing.Path.parse(svgPath, {
                            stroke: {
                                color: color,
                                width: 0
                            },
                            fill: {
                                color: color,
                            },
                            cursor: "pointer",
                        });
              
                        var label = new kendo.drawing.Text(e.series.name, [0, 0], {
                          fill: {
                            color: '#232323',
                          },
                          cursor: "pointer",
                          font: "12px 'Source Sans Pro',Lato,'Open Sans','Helvetica Neue',Arial,sans-serif!important",
                          opacity: 0.8,
                        });
              
                        layout.append(path, label);
                        layout.reflow()
              
                        return layout;
                    }
                }
            },
            seriesDefaults: {
                type: "line",
                labels: {
                    visible: false,
                    background: "transparent"
                },
                style: "smooth",
            },
            axisDefaults: {
                crosshair: {
                    visible: true,
                    opacity: 0.175,
                    width: 0.7,
                },
                majorTicks: {
                    visible: false,
                    step: 3,
                    width: 1,
                    size: 2,
                },
            },
            series: [{
                field: "Forecast",
                name: "Forecast (MW)",
                markers : {
                    visible : false
                },
                color: "#9c27b0",
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "SchFcast",
                name: "Sch Fcast / SLDC (MW)",
                markers : {
                    visible : false
                },
                color: "#e91e63",
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "Actual",
                name: "Actual Prod (MW)",
                markers : {
                    visible : false
                },
                color: "#3d8dbd",
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "ExpProd",
                name: "Exp Prod / Pwr Curv (MW)",
                markers : {
                    visible : false
                },
                color: "#8bc34a", 
                dashType: "solid",
                axis: "dynamic",
            },{
                field: "FcastWs",
                name: "Fcast ws (m/s)",
                markers : {
                    visible : false
                },
                color: "#00bcd4",
                dashType: "longDash",
                axis: "forecast",
            },{
                field: "ActualWs",
                name: "Actual ws (m/s)",
                markers : {
                    visible : false
                },
                color: "#ff9800",
                dashType: "longDash",
                axis: "forecast",
            }],
            valueAxes: [{
                line: {
                    visible: false
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n1') #",
                    visible: true,
                },
                name: "dynamic",
                title: {
                    text: "MW",
                },
                axisCrossingValue: [-10],
            },{
                line: {
                    visible: false
                }, 
                labels: {
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    template: "#= kendo.toString(value, 'n1') #",
                    margin: {
                        left: 10,
                    },
                },
                majorGridLines: {
                    visible: true,
                    color: "#eee",
                    width: 0.8,
                },
                name: "forecast",
                title: {
                    text: "m/s",
                },
            }],
            categoryAxis: {
                field: (diffDays>mindays?'Date':'TimeBlock'),
                axisCrossingValues: [0, 10000],
                majorGridLines: {
                    visible: false
                },
                majorTickType: "none",
                labels: {
                  rotation: (diffDays>mindays?45:'auto'),
                  step: (diffDays>mindays?96:4),
                }
            },
            tooltip: {
                visible: true,
                template: "${series.name} on #= moment.utc(dataItem.TimeStamp).format('DD-MM-YYYY HH:mm') # = <b>#= kendo.toString(value, 'n2') #</b>"
            },
        });
        $("#chartForecasting").data("kendoChart").refresh();
        app.loading(false);
    },200);
}

pg.initLoad = function() {
    window.setTimeout(function(){
        fa.LoadData();
        di.getAvailDate();
        app.loading(false);

        pg.refresh();
    }, 200);
}

pg.refresh = function() {
    fa.checkTurbine();

    pg.getData();
}

pg.showTurbineDown = function() {
    $('#modalTurbineDown').modal('show');
}

pg.toTime = function(s) {
    return new Date(s * 1e3).toISOString().slice(-13, -5);
};
pg.generateGridTurbineDown = function() {
    $("#grid-turbine-down").html('');
    $("#grid-turbine-down").kendoGrid({
        dataSource: {
            data: pg.TurbineDownData(),
            pageSize: 10
        },
        height: 360, //$('body').height() - heightSub + 30,
        // scrollable: true,
        sortable: true,
        filterable: false,
        pageable: {
            input: true,
            numeric: false
        },
        columns: [
            { field: "turbine", title: "Turbine", width: 80, attributes: { style: "text-align:center;" }, },
            { field: "timestart", title: "Time Start", template: "#= moment.utc(data.timestart).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timestart).format('HH:mm:ss')#", width: 140, attributes: { style: "text-align:center;" }, },
            { field: "timeend", title: "Time End", template: "#= (moment.utc(data.timeend).format('DD-MM-YYYY') == '01-01-0001'?'Not yet finished' : (moment.utc(data.timeend).format('DD-MMM-YYYY') # &nbsp; &nbsp; &nbsp; #=moment.utc(data.timeend).format('HH:mm:ss')))#", width: 140, attributes: { style: "text-align:center;" },},
            { field: "duration", title: "Duration (hh:mm:ss)", template: "#= pg.toTime(data.duration) #", width: 140, attributes: { style: "text-align:center;" }, },
            { field: "alarmcode", title: "Alarm Code", width: 90, attributes: { style: "text-align:center;" }, },
            { field: "alarmdesc", title: "Alarm Description", width: 270 },
        ]
    });
    $("#grid-turbine-down").data("kendoGrid").refresh();
}
pg.editData = function() {
    pg.EditMode('edited');
    
}
pg.saveData = function() {
    pg.EditMode('saved');
}

$(function(){
    $('#projectList').kendoDropDownList({
        change: function () {  
            di.getAvailDate();
            var project = $('#projectList').data("kendoDropDownList").value();
            fa.project = project;
            fa.populateTurbine(project);
        }
    });
    $('#btnRefresh').on('click', function () {
        pg.refresh();
    });

    $('.modal-draggable .modal-content').draggable();

    $('#modalTurbineDown').on('shown.bs.modal', function(e){
        app.loading(true);
        var url = viewModel.appName + 'forecast/getlistturbinedown';
        var param = {
            project: fa.project,
        };
        var getdata = toolkit.ajaxPostDeffered(url, param, function(res) {});
        $.when(getdata).done(function(d){
            pg.TurbineDownData(d.data);
            pg.generateGridTurbineDown();
            app.loading(false);
        });
    });

    pg.initLoad();
})
