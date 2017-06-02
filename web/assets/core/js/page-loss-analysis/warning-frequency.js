'use strict';

viewModel.WarningFrequency = new Object();
var wf = viewModel.WarningFrequency;

wf.Warning = function(){
    var isValid = fa.LoadData();
    if(isValid) {
        app.loading(true);
        if(pg.isFirstWarning() === true){
            var param = {
                period: fa.period,
                Turbine: fa.turbine(),
                DateStart: fa.dateStart,
                DateEnd: fa.dateEnd,
                Project: fa.project
            };

            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getwarning", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                if (res.data.Data.length != 0) {
                    setTimeout(function(){
                        wf.generateGrid(res.data.Data);
                        app.loading(false);
                    },200);
                }else{
                    setTimeout(function(){
                        wf.generateGrid([]);
                        app.loading(false);
                    },200);
                }
            });
            $('#availabledatestart').html(pg.availabledatestartwarning());
            $('#availabledateend').html(pg.availabledateendwarning());
        }else{
            setTimeout(function(){
                $("#warningGrid").data("kendoGrid").refresh();
                $('#availabledatestart').html(pg.availabledatestartwarning());
                $('#availabledateend').html(pg.availabledateendwarning());
                app.loading(false);
            },200);
            
        }
    }
}

wf.generateGrid = function (dataSource) {
    var config = {
        dataSource: {
            data: dataSource,
            pageSize: 10
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Warning Description", field: "desc", attributes: { class: "align-left row-custom" }, width: 250, locked: true, filterable: false },
            { title: "Total", field: "total", attributes: { class: "align-center row-custom" }, width: 70, locked: true, filterable: false },
        ],
        dataBound: function(){
            setTimeout(function(){
                $("#warningGrid >.k-grid-header >.k-grid-header-locked > table > thead >tr").css("height","20px");
                $("#warningGrid >.k-grid-header >.k-grid-header-wrap > table > thead >tr").css("height","20px");
                // app.loading(false);
            },200);
        },
    };

    if (dataSource.length > 0){
        $.each(dataSource[0].turbines, function (i, val) {
            var column = {
                title: val.turbine.Turbine,
                field: "turbines["+i+"].count",
                attributes: { class: "align-center" },
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;'
                },
                width: 80
            }

            config.columns.push(column);
        });
    }else{
        var column = {
            title: "",
            attributes: { class: "align-center" },
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            width: 80
        }

        config.columns.push(column);
    }   

    $('#warningGrid').html("");
    $('#warningGrid').kendoGrid(config);
    $('#warningGrid').data('kendoGrid').refresh();

    // setTimeout(function() {
    //     app.loading(false);
    // }, 500);
}