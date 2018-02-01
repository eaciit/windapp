'use strict';

viewModel.DateInfo = new Object();
var di = viewModel.DateInfo;


var availDateAll;

di.minDatetemp = ko.observable([]);
di.maxDatetemp = ko.observable([]);

di.getAvailDate = function (param) {

    var reqDate = toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getavaildateall", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        availDateAll = res.data;

        var namaproject;
        
        var projectVal = $("#projectList").data("kendoDropDownList").value();
        if( projectVal == undefined || projectVal == "") {
            namaproject = "Tejuva";
        }else{
            namaproject= projectVal;
        }

        if(param != undefined){
            di.minDatetemp(kendo.toString(moment.utc(availDateAll[namaproject][param][0]).format('DD-MMM-YYYY')));
            di.maxDatetemp(kendo.toString(moment.utc(availDateAll[namaproject][param][1]).format('DD-MMM-YYYY')));   
            var maxDateData = new Date(availDateAll[namaproject][param][1]);    
        }else{
            di.minDatetemp(kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][0]).format('DD-MMM-YYYY')));
            di.maxDatetemp(kendo.toString(moment.utc(availDateAll[namaproject]["ScadaData"][1]).format('DD-MMM-YYYY')));   
            var maxDateData = new Date(availDateAll[namaproject]["ScadaData"][1]);         
        }
        
         
        // $('#dateStart').data('kendoDatePicker').value(new Date(Date.UTC(moment(maxDateData).get('year'),  maxDateData.getMonth(),  maxDateData.getDate() - 7, 0, 0, 0, 0)));
        // $('#dateEnd').data('kendoDatePicker').value(di.maxDatetemp());
    })
    return reqDate;
};

