'use strict';

viewModel.DateInfo = new Object();
var di = viewModel.DateInfo;


var availDateAll;

di.minDatetemp = ko.observable([]);
di.maxDatetemp = ko.observable([]);

di.getAvailDate = function () {

    toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateAll = res.data;

        var namaproject;
        
        if(fa.project == undefined || fa.project == "") {
            namaproject = "Tejuva";
        }else{
            namaproject=  fa.project;
        }

        di.minDatetemp(kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][0]).format('DD-MMMM-YYYY')));
        di.maxDatetemp(kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][1]).format('DD-MMMM-YYYY')));
    })
};

