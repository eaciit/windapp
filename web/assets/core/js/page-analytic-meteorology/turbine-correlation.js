'use strict';

viewModel.TurbineCorrelation = new Object();
var tc = viewModel.TurbineCorrelation;
tc.Column = ko.observableArray([]);
tc.datas = ko.observableArray([]);

tc.newData = ko.observableArray([]);

tc.getCss = function(index, da){
    var color = 'white';
    var opacity = 1;
    var rgba = 'rgba(255,255,255)';
    var fontColor = "#333";
    var css = {"background":rgba, "color":fontColor};


    if(tc.newData().length != 0){
        if (da in tc.newData()[index]){
            color = tc.newData()[index][da].Color;
            opacity = tc.newData()[index][da].Opacity;

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
tc.TurbineCorrelation = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.showFilter();
        if(pm.isFirstTurbine() === true){
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
             toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/getwindcorrelation", param, function (res) {
                if (!app.isFine(res)) {
                    app.loading(false);
                    return;
                }
                dataSource = res.data.Data;
                columns = res.data.Column;
                heat = res.data.Heat;
                turbineName = res.data.TurbineName;

                tc.datas(dataSource);
                tc.newData(heat);
                tc.Column(columns);

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
                $("#gridTurbineCorrelation").html("");
                $("#gridTurbineCorrelation").kendoGrid({
                    dataSource: knownOutagesDataSource,
                    columns: columnArray,
                    filterable: false,
                    sortable: false,
                    dataBound: function (e) {
                        
                        var ini = this.wrapper;
                        $.each(tc.Column(), function(i, col){
                            var columns = e.sender.columns;
                            var columnIndex = ini.find(".k-grid-header [data-field=" + col + "]").index();

                            // iterate the data items and apply row styles where necessary
                            var dataItems = e.sender.dataSource.view();
                            for (var j = 0; j < dataItems.length; j++) {

                                var units = dataItems[j].get(col);
          
                                var row = e.sender.tbody.find("[data-uid='" + dataItems[j].uid + "']");
                                var cell = row.children().eq(columnIndex);

                                cell.css(tc.getCss(j,col));
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
                    $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
                    $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
                    $("#gridTurbineCorrelation").data("kendoGrid").refresh(); 
                    app.loading(false);
                    pm.isFirstTurbine(false)    
                },200);
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');

            });
        }else{
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
            setTimeout(function(){
                 app.loading(false);
                 $("#gridTurbineCorrelation").data("kendoGrid").refresh();
                 $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","37px");
                 $("#gridTurbineCorrelation >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","37px");
            }, 500);
        }
    }
}