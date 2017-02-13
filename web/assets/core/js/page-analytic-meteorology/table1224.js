'use strict';

viewModel.Table1224 = new Object();
var tb = viewModel.Table1224;

tb.dataSourceTable = ko.observableArray();
tb.MetTowerColumn = ko.observableArray([
    {value: true, text: "Wind Speed (m/s)", _id:"metWs", index:0 },
    {value: true, text: "Temp (°C)", _id:"metTemp", index: 1},

]);

tb.TurbineColumn = ko.observableArray([
    {_id: "turbineWs", text: "Wind Speed (m/s)", value:true , index:0},
    {_id: "turbineTemp", text: "Temp (°C)", value:true , index:1},
    {_id: "turbinePower", text: "Power (kWH)", value: true, index:2},
]);

// 12/24 table 
tb.generateGridTable = function (datatype) {
    app.loading(true);
    $('#gridTable1224').html('');

    var dataSource = [];
    if(datatype == "turbine") {
        dataSource = tb.dataSourceTable().DataTurbine;
    } else {
        dataSource = tb.dataSourceTable().DataMet;
    }

    // var aggregates = [];

    // $.each(dataSource[0].details[0].col, function (key, val) {
    //     var aggregate = { field: key , aggregate: (key == "Power" ? "sum" : "average")};
    //     aggregates.push(aggregate);
    // });

    var config = {
        dataSource: {
            data: dataSource,
            pageSize: 10,
            // aggregate: aggregates
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Hours", field: "hours", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false},
        ],
         dataBound: function(){
            setTimeout(function(){
                $("#gridTable1224 >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","75px");
                // $("#gridTable1224 >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","75px");
                tb.refreshTable(datatype);
            },200);
        },
    };
    console.log(config);
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
            if(key == "WS") {
                title = key + " (m/s)";
            } else if(key == "Temp") {
                title = key + " (" + String.fromCharCode(176) + "C)";
            } else {
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
                // filterable: false, 
                // footerTemplate: (key == "Power" ? "<div style='text-align:center'>#=kendo.toString(sum, 'n2')#</div>"  : "<div style='text-align:center'>#=kendo.toString(average, 'n2')#</div>"  )
                // footerTemplate: "<div>#= kendo.toString(data.details["+i+"].col."+ key+".sum,'n2') # </div>"
            };
            column.columns.push(colChild);
        });

        config.columns.push(column);
    });
    
    
    $('#gridTable1224').kendoGrid(config);
}

tb.hideShowColumn = function(i, type){
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
    app.loading(true);
    fa.LoadData();
    var datatype = '';
    if(pm.isFirstTwelve() === true){
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
            Turbine: fa.turbine,
            Project: fa.project,
        };

        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/table1224", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            tb.dataSourceTable(res.data);
            tb.generateGridTable(datatype);
            pm.isFirstTwelve(false); 
            if($("#met").is(':checked')) {
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
            } else {
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
            }
        });
    }else{
        setTimeout(function(){
            tb.refreshTable(datatype);
            if($("#met").is(':checked')) {
                pm.isMet(true);
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartmet + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendmet + '</strong>');
            } else {
                 pm.isMet(false);
                $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
                $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
            }
        },300);
    }
}