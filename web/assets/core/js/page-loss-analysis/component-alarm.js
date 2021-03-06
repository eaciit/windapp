'use strict';

viewModel.ComponentAlarm = new Object();
var ca = viewModel.ComponentAlarm;

ca.dtCompponentAlarm = ko.observable();

ca.Component = function(){    
    var valid = fa.LoadData();
    if (valid) {
        app.loading(true);
        pg.setAvailableDate(false);
        if(pg.isFirstComponentAlarm() === true){
            var dateStart = $('#dateStart').data('kendoDatePicker').value();
            var dateEnd = new Date(moment($('#dateEnd').data('kendoDatePicker').value()).format('YYYY-MM-DD'));

            var param = {
                period: fa.period,
                dateStart: dateStart,
                dateEnd: dateEnd,
                turbine: fa.turbine(),
                project: fa.project,
            }

            toolkit.ajaxPost(viewModel.appName + "analyticlossanalysis/getcomponentalarmtab", param, function (res) {
                if (!app.isFine(res)) {
                    return;
                }
                setTimeout(function(){
                    ca.dtCompponentAlarm(res.data)
                    var HAlarm = $('#filter-analytic').width() * 0.235
                    var wAll = $('#filter-analytic').width() * 0.275
                    var componentduration = _.sortBy(ca.dtCompponentAlarm().componentduration, '_id');
                    var componentfrequency = _.sortBy(ca.dtCompponentAlarm().componentfrequency, '_id');
                    var componentloss = _.sortBy(ca.dtCompponentAlarm().componentloss, '_id');

                    var id = $("#downtimeGroup .active").attr('id')

                    if(id == 'lblComp'){
                        /*Component / Alarm Type Tab*/
                        dt.GenChartDownAlarmComponent("component",componentduration,'chartCADuration',SeriesAlarm,true, "", "Hours",false,-90,HAlarm,wAll,"N1");
                        dt.GenChartDownAlarmComponent("component",componentfrequency,'chartCAFrequency',SeriesAlarm,true, "", "Times",false,-90,HAlarm,wAll,"N0");
                        dt.GenChartDownAlarmComponent("component",componentloss,'chartCATurbineLoss',SeriesAlarm,true, "", "MWh",false,-90,HAlarm,wAll,"N1");
                    }else{                    
                        dt.GenChartDownAlarmComponent("alarm",ca.dtCompponentAlarm().alarmduration,'chartCADuration',SeriesAlarm,false, "", "Hours",false,-90,HAlarm,wAll,"N1");
                        dt.GenChartDownAlarmComponent("alarm",ca.dtCompponentAlarm().alarmfrequency,'chartCAFrequency',SeriesAlarm,false, "", "Times",false,-90,HAlarm,wAll,"N0");
                        dt.GenChartDownAlarmComponent("alarm",ca.dtCompponentAlarm().alarmloss,'chartCATurbineLoss',SeriesAlarm,false, "", "MWh",false,-90,HAlarm,wAll,"N1");
                    }

                    app.loading(false);
                    pg.isFirstComponentAlarm(false);
                },300);
            }); 
        }else{
            setTimeout(function(){
                $("#chartCADuration").data("kendoChart").refresh();
                $("#chartCAFrequency").data("kendoChart").refresh();
                $("#chartCATurbineLoss").data("kendoChart").refresh();
                app.loading(false);
            },200); 
        }
    }
}