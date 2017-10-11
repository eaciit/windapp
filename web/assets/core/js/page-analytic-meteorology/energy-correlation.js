'use strict';

viewModel.EnergyCorrelation = new Object();
var ec = viewModel.EnergyCorrelation;
ec.Column = ko.observableArray([]);
ec.datas = ko.observableArray([]);

ec.newData = ko.observableArray([]);

ec.getCss = function(index, da){
    var color = 'white';
    var opacity = 1;
    var rgba = 'rgba(255,255,255)';
    var fontColor = "#333";
    var css = {"background":rgba, "color":fontColor};


    if(ec.newData().length != 0){
        if (da in ec.newData()[index]){
            color = ec.newData()[index][da].Color;
            opacity = ec.newData()[index][da].Opacity;

            // if(opacity > 0.5){
            //     fontColor = "#fff";
            // }

            if(color == "red") { 
                // rgba = 'rgba(255,0,0,'+opacity+')';
                rgba = 'rgba(248,109,111,'+opacity+')';
            }else if(color == "green"){
                 // rgba = 'rgba(0,128,0,'+opacity+')';
                 rgba = 'rgba(100,190,124,'+opacity+')';
            }else{
                 rgba = 'rgba(255,255,255,'+opacity+')';
            }

            // css = {"background":rgba, "font-weight":"bold","color":fontColor};
            css = {"background":rgba, "color":fontColor};
        }
    }
    
    return css;
}

// Turbine Correlation
ec.EnergyCorrelation = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        if(pm.isFirstEnergy() === true){
            app.loading(true);
            var param = {
                period: fa.period,
                dateStart: fa.dateStart,
                dateEnd: fa.dateEnd,
                turbine: fa.turbine(),
                project: fa.project
            };
            var dataSource;
            var columns;
            var heat;
            var turbineName;
            var judul;
             toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getenergycorrelation", param, function (res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                dataSource = res.data.Data;
                columns = res.data.Column;
                heat = res.data.Heat;
                turbineName = res.data.TurbineName;

                ec.datas(dataSource);
                ec.newData(heat);
                ec.Column(columns);

                var schemaModel = {};
                var columnArray = [];

                $.each(columns, function (index, da) {
                    schemaModel[da] = {type: (da == "Turbine" ? "string" : "int")};
                    judul = da
                    if(da != "Turbine" && da != "MetTower") {
                        judul = turbineName[da];
                    }

                    var column = {
                        title: judul,
                        field: da,
                        locked: (da == "Turbine" ? true : false),
                        headerAttributes: {
                            style: "text-align: center;",
                        },
                        attributes: {
                            style: "text-align:center",
                            turbine: da,
                            index: index,
                        },
                        width: 70,
                        template:( da != "Turbine" ? "#= kendo.toString("+da+", 'n2') #" : "#= kendo.toString("+da+") #")
                    }

                    columnArray.push(column);
                });

                var schemaModelNew = kendo.data.Model.define({
                    id: "Turbine",
                    fields: schemaModel,
                });

                var knownOutagesDataSource = new kendo.data.DataSource({
                    data: dataSource,
                    schema: {
                        model: schemaModelNew
                    }
                });
                $("#gridEnergyCorrelation").html("");
                $("#gridEnergyCorrelation").kendoGrid({
                    dataSource: knownOutagesDataSource,
                    columns: columnArray,
                    filterable: false,
                    sortable: false,
                    dataBound: function (e) {
                        
                        var ini = this.wrapper;
                        $.each(ec.Column(), function(i, col){
                            var columns = e.sender.columns;
                            var columnIndex = ini.find(".k-grid-header [data-field=" + col + "]").index();

                            // iterate the data items and apply row styles where necessary
                            var dataItems = e.sender.dataSource.view();
                            for (var j = 0; j < dataItems.length; j++) {

                                var units = dataItems[j].get(col);
          
                                var row = e.sender.tbody.find("[data-uid='" + dataItems[j].uid + "']");
                                var cell = row.children().eq(columnIndex);

                                cell.css(ec.getCss(j,col));
                            }
                        });


                        if (e.sender._data.length == 0) {
                            var mgs, col;
                            mgs = "No results found for";
                            col = 9;
                            var contentDiv = this.wrapper.children(".k-grid-content"),
                            dataTable = contentDiv.children("table");
                            if (!dataTable.find("tr").length) {
                                dataTable.children("tbody").append("<tr><td colspan='" + col + "'><div style='color:red;width:500px'>" + mgs + "</div></td></tr>");
                                if (navigator.userAgent.match(/MSIE ([0-9]+)\./)) {
                                    dataTable.width(this.wrapper.children(".k-grid-header").find("table").width());
                                    contentDiv.scrollLeft(1);
                                }
                            }
                        }
                        
                    },
                    pageable: false,
                    scrollable: true,
                    resizable: false,
                    height:390,
                });

                setTimeout(function(){
                    $("#gridEnergyCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
                    $("#gridEnergyCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
                    $("#gridEnergyCorrelation").data("kendoGrid").refresh(); 
                    app.loading(false);
                    pm.isFirstEnergy(false)    
                },200);
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

            });
        }else{
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
            setTimeout(function(){
                 app.loading(false);
                 $("#gridEnergyCorrelation").data("kendoGrid").refresh();
                 $("#gridEnergyCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
                 $("#gridEnergyCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
            }, 500);
        }
    }
}