'use strict';

viewModel.Table1224 = new Object();
var tb = viewModel.Table1224;

tb.isGraph = ko.observable(false);
tb.dataSourceTable = ko.observableArray();
tb.dataSourceGraphTurbine = ko.observableArray();
tb.dataSourceGraphMetTower = ko.observableArray();
tb.MetTowerColumn = ko.observableArray([
    {value: true, text: "Wind Speed (m/s)", _id:"metWs", index:0 },
    {value: true, text: "Temp (°C)", _id:"metTemp", index: 1},

]);

tb.TurbineColumn = ko.observableArray([
    {_id: "turbineWs", text: "Wind Speed (m/s)", value:true , index:0},
    {_id: "turbineTemp", text: "Temp (°C)", value:true , index:1},
    {_id: "turbinePower", text: "Power (kWH)", value: true, index:2},
]);

tb.generateGraph = function(serie){
    app.loading(true);
    tb.isGraph(true);

    var dataSource,buttonType,nameChk;

    if($("#met").is(':checked')) {
        dataSource = tb.dataSourceGraphMetTower();
        buttonType = tb.MetTowerColumn()[0]._id;
        nameChk = "chk-column-met"
    } else {
        dataSource = tb.dataSourceGraphTurbine();
        buttonType = tb.TurbineColumn()[0]._id;
        nameChk = "chk-column-turbine"
    }

    if(serie == undefined){
        if($("input[name="+nameChk+"]:checked").length > 1){
            $("input[name="+nameChk+"]").prop("checked",false);
            $('#'+buttonType).prop('checked', true);
        }
    }


    var series = $("input[name="+nameChk+"]:checked").attr("id");
    var seriesName = $("#met").is(':checked') == true ? series.substr(3).toLowerCase() : series.substr(7).toLowerCase();
    var idParent = $("input[name="+nameChk+"]:checked").closest('label').attr('id'); 


    var title = $("#"+idParent).find(".colRed").text();

    var dataGraph = new kendo.data.DataSource({
        data: dataSource,
        group: {
            field: "timeint"
        },

        sort: [{
            field: "hours",
            dir: "asc"
        },{
            field: "timeint",
            dir: "asc"
        }],
    });


    $("#chartTable1224").html("");
    $("#chartTable1224").kendoChart({
        title: { text: "" },
        dataSource: dataGraph,
        seriesColors: colorField,
        series: [{
            type: "line",
            style: "smooth",
            field: seriesName,
            categoryField: "hours",
            name: "#= group.value #",
            markers: {
                visible: false
            }
        }],
        legend: {
            position: "top",
            labels: {
                template: "#= moment(text).format('MMM YYYY') #",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
        },
        valueAxis: {
            labels: {
                format: "n0",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            line: {
                visible: false
            },
            title: { 
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                text: title
            },
            axisCrossingValue: -10,
            majorGridLines: {
                visible: true,
                color: "#eee",
                width: 0.8,
            }
        },
        categoryAxis: {
            field: "hours",
            majorGridLines: {
                visible: false
            },
            labels: {
                format: "MMM",
                font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            },
            majorTickType: "none"
        },
        tooltip: {
            visible: true,
            background: "rgb(255,255,255, 0.9)",
            template: "#= moment(series.name).format('MMM YYYY') # : #= kendo.toString(value,'n2')#",
            color: "#58666e",
            font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
            border: {
                color: "#eee",
                width: "2px",
            },

        },
        dataBound : function(){
            app.loading(false);
        }
    });

}
// 12/24 table 
tb.generateGridTable = function (datatype) {
    app.loading(true);
    tb.isGraph(false);
    $('#gridTable1224').html('');

    var dataSource = [];
    var total = [];

    if(datatype == undefined || datatype == ''){
        if($("#met").is(':checked')) {
            datatype = 'met';
        } else {
            datatype = 'turbine';
        }
    }else{
        datatype = datatype;
    }

    if(datatype == "turbine") {
        dataSource = tb.dataSourceTable().DataTurbine;
        total = tb.dataSourceTable().TotalTurbine;
    } else {
        dataSource = tb.dataSourceTable().DataMet;
        total = tb.dataSourceTable().TotalMet;
    }


    var config = {
        dataSource: {
            data: dataSource,
            pageSize: 10,
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Hours", field: "hours", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false,footerTemplate: "<center>Total</center>"},
        ],
         dataBound: function(){
            setTimeout(function(){
                $("#gridTable1224 >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","75px");
                // $("#gridTable1224 >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","75px");
                tb.refreshTable(datatype);
            },200);
        },
    };


    if(dataSource != null){
        $.each(dataSource[0].details, function (i, val) {
            var column = {
                title: val.time,
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                columns: []
            }
            var keyIndex = [];
            if(datatype == "turbine") {
                keyIndex = ["WS", "Temp", "Power"];
            } else {
                keyIndex = ["WS", "Temp"];
            }

            $.each(keyIndex, function(j, key){
                var title = "";
                var totalSum;
                if(key == "WS") {
                    title = key + " (m/s)";
                    totalSum = total[i].windspeed
                } else if(key == "Temp") {
                    title = key + " (" + String.fromCharCode(176) + "C)";
                    totalSum = total[i].temp;
                } else {
                    totalSum = total[i].power
                    title = key + " (MWH)";
                }

                var colChild = {
                    title: title,                
                    field: "details["+i+"].col."+ key,
                    attributes: { class: "align-center row-custom" },
                    width: 100,
                    headerAttributes: {
                        style: 'font-weight: bold; text-align: center;',
                    },
                    format: "{0:n2}",
                    filterable: false, 
                    footerTemplate: "<div style='text-align:center'>#= kendo.toString("+totalSum+",'n2') # </div>"
                };
                column.columns.push(colChild);

            });

            config.columns.push(column);
        });

        $('#gridTable1224').kendoGrid(config);
    }else{
        $('#gridTable1224').html("");
        app.loading(false);
    }
}

tb.hideShowColumn = function(i, type){
    if($("#gridDineural").is(':checked')) {
        var grid = $("#gridTable1224").data("kendoGrid");  
        var columns = grid.columns;

        if($('[name='+type+']:checked').length < 1){
            toolkit.showError("Grid must show at least one column");
            $('#'+i._id).prop('checked', true);
            return false;    
        }else{
            $.each(columns, function(index, val){
                if(index > 0){
                    var col = grid.columns[index].columns[i.index];
                    if (col.hidden) {
                      grid.showColumn(col.field);
                    } else {
                      grid.hideColumn(col.field);
                    } 
                } 
            });
        }
    }else{
        tb.checkboxGraph(i);
        tb.generateGraph(i._id);
    }
}

tb.checkboxGraph = function(i){
    var nameChk ="";
    if($("#met").is(':checked')) {
        nameChk = "chk-column-met"
    }else{
        nameChk = "chk-column-turbine"
    }

    $("input[name="+nameChk+"]").prop("checked",false);
    $('#'+i._id).prop('checked', true);
}

tb.getObjects = function(obj, key, val){
    var objects = [];
    for (var i in obj) {
        if (!obj.hasOwnProperty(i)) continue;
        if (typeof obj[i] == 'object') {
            objects = objects.concat(tb.getObjects(obj[i], key, val));
        } else if (i == key && obj[key] == val) {
            objects.push(obj);
        }
    }
    return objects;
}

tb.refreshTable = function(datatype){
    var grid = $("#gridTable1224").data("kendoGrid");  
    var columns = grid.columns;
    var data = (datatype == "met" ? tb.MetTowerColumn() : tb.TurbineColumn());
    var results = $.each($('[name="chk-column-'+datatype+'"]:not(:checked)'), function(i, val){
        var diff = tb.getObjects(data, "_id", val.id);
        $.each(diff, function(a, res){
             $.each(columns, function(e, value){
                if(e > 0){
                    var col = grid.columns[e].columns[res.index];
                    grid.hideColumn(col.field);
                } 
            });
             
        });
    });

    $.when(results).done(function(){
        setTimeout(function(){
            app.loading(false);
        },300);
    })
}

tb.Table = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        pm.hideFilter();
        var datatype = '';
        if(pm.isFirstTwelve() === true){
            app.loading(true);
            if(datatype == undefined || datatype == ''){
                if($("#met").is(':checked')) {
                    datatype = 'met';
                } else {
                    datatype = 'turbine';
                }
            }else{
                datatype = datatype;
            }
            

            var param = {
                Turbine: fa.turbine(),
                Project: fa.project,
            };

            var req1 = toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/table1224", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                
                tb.dataSourceTable(res.data);
                pm.isFirstTwelve(false); 
                if($("#met").is(':checked')) {
                    $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
                    $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
                } else {
                    $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                    $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
                }
            });


            var req2 = toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/graph1224", {Turbine: fa.turbine(),Project: fa.project,Data : "turbine"}, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                tb.dataSourceGraphTurbine(res.data);
            });

            var req3 = toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/graph1224", {Turbine: fa.turbine(),Project: fa.project,Data : "mettower"}, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                tb.dataSourceGraphMetTower(res.data);
            });

            $.when(req1, req2, req3).done(function(){
                if($("#gridDineural").is(':checked')) {
                    tb.generateGridTable(datatype);
                } else {
                    tb.generateGraph();
                }
            });

        }else{
            setTimeout(function(){
                if($("#met").is(':checked')) {
                    pm.isMet(true);
                    $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
                    $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
                } else {
                     pm.isMet(false);
                    $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                    $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
                }

                if($("#gridDineural").is(':checked')) {
                    tb.refreshTable(datatype);
                } else {
                    tb.generateGraph();
                }
            },300);
        }
    }
}
