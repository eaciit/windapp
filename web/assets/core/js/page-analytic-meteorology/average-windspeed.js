'use strict';

viewModel.AverageWindSpeed = new Object();
var aw = viewModel.AverageWindSpeed;
aw.dataSourceAverage = ko.observableArray();

aw.generateGridAverage = function () {
    var config = {
        dataSource: {
            data: aw.dataSourceAverage(),
            pageSize: 10
        },
        pageable: {
            pageSize: 10,
            input: true, 
        },
        scrollable: true,
        sortable: true,
        columns: [
            { title: "Turbine(s)", field: "turbine", attributes: { class: "align-center row-custom" }, width: 100, locked: true, filterable: false },
        ],
    };

    $.each(aw.dataSourceAverage()[0].details, function (i, val) {
        var wra = val.col.WRA;        
        var column = {
            title: val.time + " <br/> WRA "+wra+ " (m/s)",
            headerAttributes: {
                style: 'font-weight: bold; text-align: center;'
            },
            columns: [],
            width: 120
        }

        // var keyIndex = ["WRA", "Onsite"];
        var keyIndex = ["Onsite"];
        var j = 0;        

        $.each(keyIndex, function(j, key){
            // wra = 
            var colChild = {
                title: key + " (m/s)",                
                field: "details["+i+"].col."+ key ,
                width: 120,
                attributes: { class: "align-center row-custom" },
                headerAttributes: {
                    style: 'font-weight: bold; text-align: center;',
                },
                format: "{0:n2}",
                filterable: false
            };
            column.columns.push(colChild);
        });

        config.columns.push(column);
    });

    $('#gridAvgWs').html("");
    $('#gridAvgWs').kendoGrid(config);
    $('#gridAvgWs').data('kendoGrid').refresh();
}

aw.AverageWindSpeed = function() {
    app.loading(true);
    fa.LoadData();
    if(pm.isFirstAverage() === true){
        var param = {
            period: fa.period,
            Turbine: fa.turbine,
            DateStart: fa.dateStart,
            DateEnd: fa.dateEnd,
            Project: fa.project
        };

        toolkit.ajaxPost(viewModel.appName + "analyticmeteorology/averagewindspeed", param, function (res) {
            if (!app.isFine(res)) {
                return;
            }
            aw.dataSourceAverage(res.data.Data.turbine);
            aw.generateGridAverage();
            app.loading(false);
            pm.isFirstAverage(false);
            $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
            $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        });        
    }else{
        $('#availabledatestart').html('Data Available from: <strong>' + availDateList.availabledatestartscada + '</strong> until: ');
        $('#availabledateend').html('<strong>' + availDateList.availabledateendscada + '</strong>');
        setTimeout(function(){
            $("#gridAvgWs").data("kendoGrid").refresh();
            app.loading(false);
        }, 300);
    }

}