'use strict';

viewModel.AnalyticWindDistribution = new Object();
var wd = viewModel.AnalyticWindDistribution;

wd.turbineList = ko.observableArray([]);
wd.turbine = ko.observableArray([]);
wd.valueCategory = ko.observableArray([
    { "value": "powerGeneration", "text": "Power Generation (MW)" },
    { "value": "machine", "text": "Machine Availability" },
    { "value": "scada", "text": "Scada Availability" },
    { "value": "grid", "text": "Grid Availability" },
]);

var color = ["#B71C1C", "#E57373", "#F44336", "#D81B60", "#F06292", "#880E4F",
    "#4A148C", "#7B1FA2", "#9C27B0", "#BA68C8", "#1A237E", "#5C6BC0",
    "#1E88E5", "#0277BD", "#0097A7", "#26A69A", "#4DD0E1", "#81C784",
    "#8BC34A", "#1B5E20", "#827717", "#C0CA33", "#DCE775", "#FF6F00", "#A1887F",
    "#FFEE58", "#004D40", "#212121", "#607D8B", "#BDBDBD", "#FF00CC", "#9999FF"
];

wd.populateTurbine = function(){
    wd.turbine([]);
    if(fa.turbine().length == 0){
        $.each(fa.turbineList(), function(i, val){
            if (i > 0){
                wd.turbine.push(val.text);
            }
        });
    }else{
        wd.turbine(fa.turbine());
    }

}

var Data = {
    LoadData: function () {
        fa.getProjectInfo();
        fa.LoadData();
         wd.populateTurbine();
        this.ChartWindDistributon();
    },
    ChartWindDistributon: function () {
        // var res = {"Seconds":0,"Data":{"Total":619,"Data":[{"Turbine":"B1","Category":0,"Contribute":0.012868},{"Turbine":"B1","Category":1,"Contribute":0.000212},{"Turbine":"B1","Category":1.5,"Contribute":0.000254},{"Turbine":"B1","Category":2,"Contribute":0.00055},{"Turbine":"B1","Category":2.5,"Contribute":0.000339},{"Turbine":"B1","Category":3,"Contribute":0.000339},{"Turbine":"B1","Category":3.5,"Contribute":0.000339},{"Turbine":"B1","Category":4,"Contribute":0.000847},{"Turbine":"B1","Category":4.5,"Contribute":0.00072},{"Turbine":"B1","Category":5,"Contribute":0.001608},{"Turbine":"B1","Category":5.5,"Contribute":0.001989},{"Turbine":"B1","Category":6,"Contribute":0.00237},{"Turbine":"B1","Category":6.5,"Contribute":0.002455},{"Turbine":"B1","Category":7,"Contribute":0.001947},{"Turbine":"B1","Category":7.5,"Contribute":0.002201},{"Turbine":"B1","Category":8,"Contribute":0.002074},{"Turbine":"B1","Category":8.5,"Contribute":0.002921},{"Turbine":"B1","Category":9,"Contribute":0.002074},{"Turbine":"B1","Category":9.5,"Contribute":0.002116},{"Turbine":"B1","Category":10,"Contribute":0.001312},{"Turbine":"B1","Category":10.5,"Contribute":0.000296},{"Turbine":"B1","Category":11,"Contribute":0.000127},{"Turbine":"B1","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B16","Category":0,"Contribute":0.012317},{"Turbine":"B16","Category":1,"Contribute":0.000127},{"Turbine":"B16","Category":1.5,"Contribute":0.000296},{"Turbine":"B16","Category":2,"Contribute":0.000339},{"Turbine":"B16","Category":2.5,"Contribute":0.00055},{"Turbine":"B16","Category":3,"Contribute":0.000254},{"Turbine":"B16","Category":3.5,"Contribute":0.000381},{"Turbine":"B16","Category":4,"Contribute":0.000804},{"Turbine":"B16","Category":4.5,"Contribute":0.001101},{"Turbine":"B16","Category":5,"Contribute":0.00127},{"Turbine":"B16","Category":5.5,"Contribute":0.001905},{"Turbine":"B16","Category":6,"Contribute":0.002328},{"Turbine":"B16","Category":6.5,"Contribute":0.002286},{"Turbine":"B16","Category":7,"Contribute":0.00182},{"Turbine":"B16","Category":7.5,"Contribute":0.002328},{"Turbine":"B16","Category":8,"Contribute":0.002286},{"Turbine":"B16","Category":8.5,"Contribute":0.002286},{"Turbine":"B16","Category":9,"Contribute":0.003005},{"Turbine":"B16","Category":9.5,"Contribute":0.001905},{"Turbine":"B16","Category":10,"Contribute":0.001566},{"Turbine":"B16","Category":10.5,"Contribute":0.000381},{"Turbine":"B16","Category":11,"Contribute":0.000212},{"Turbine":"B16","Category":11.5,"Contribute":0.000212},{"Turbine":"B16","Category":12,"Contribute":4.2E-05},{"Turbine":"B33","Category":0,"Contribute":0.011683},{"Turbine":"B33","Category":0.5,"Contribute":4.2E-05},{"Turbine":"B33","Category":1,"Contribute":0.000296},{"Turbine":"B33","Category":1.5,"Contribute":0.000296},{"Turbine":"B33","Category":2,"Contribute":0.000296},{"Turbine":"B33","Category":2.5,"Contribute":0.000593},{"Turbine":"B33","Category":3,"Contribute":0.000423},{"Turbine":"B33","Category":3.5,"Contribute":0.000339},{"Turbine":"B33","Category":4,"Contribute":0.000889},{"Turbine":"B33","Category":4.5,"Contribute":0.001185},{"Turbine":"B33","Category":5,"Contribute":0.001143},{"Turbine":"B33","Category":5.5,"Contribute":0.001735},{"Turbine":"B33","Category":6,"Contribute":0.002243},{"Turbine":"B33","Category":6.5,"Contribute":0.003132},{"Turbine":"B33","Category":7,"Contribute":0.002328},{"Turbine":"B33","Category":7.5,"Contribute":0.001862},{"Turbine":"B33","Category":8,"Contribute":0.002497},{"Turbine":"B33","Category":8.5,"Contribute":0.00237},{"Turbine":"B33","Category":9,"Contribute":0.001566},{"Turbine":"B33","Category":9.5,"Contribute":0.001947},{"Turbine":"B33","Category":10,"Contribute":0.001312},{"Turbine":"B33","Category":10.5,"Contribute":0.001101},{"Turbine":"B33","Category":11,"Contribute":0.000466},{"Turbine":"B33","Category":11.5,"Contribute":0.000127},{"Turbine":"B33","Category":12,"Contribute":8.5E-05},{"Turbine":"B33","Category":13,"Contribute":4.2E-05},{"Turbine":"B38","Category":0,"Contribute":0.012571},{"Turbine":"B38","Category":1,"Contribute":0.000169},{"Turbine":"B38","Category":1.5,"Contribute":0.000339},{"Turbine":"B38","Category":2,"Contribute":0.000381},{"Turbine":"B38","Category":2.5,"Contribute":0.000423},{"Turbine":"B38","Category":3,"Contribute":0.000212},{"Turbine":"B38","Category":3.5,"Contribute":0.000423},{"Turbine":"B38","Category":4,"Contribute":0.000508},{"Turbine":"B38","Category":4.5,"Contribute":0.001058},{"Turbine":"B38","Category":5,"Contribute":0.001439},{"Turbine":"B38","Category":5.5,"Contribute":0.001524},{"Turbine":"B38","Category":6,"Contribute":0.001989},{"Turbine":"B38","Category":6.5,"Contribute":0.002201},{"Turbine":"B38","Category":7,"Contribute":0.002836},{"Turbine":"B38","Category":7.5,"Contribute":0.002497},{"Turbine":"B38","Category":8,"Contribute":0.002963},{"Turbine":"B38","Category":8.5,"Contribute":0.002497},{"Turbine":"B38","Category":9,"Contribute":0.002159},{"Turbine":"B38","Category":9.5,"Contribute":0.001608},{"Turbine":"B38","Category":10,"Contribute":0.001312},{"Turbine":"B38","Category":10.5,"Contribute":0.000508},{"Turbine":"B38","Category":11,"Contribute":0.000296},{"Turbine":"B38","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B38","Category":12,"Contribute":4.2E-05},{"Turbine":"B4","Category":0,"Contribute":0.012148},{"Turbine":"B4","Category":0.5,"Contribute":4.2E-05},{"Turbine":"B4","Category":1,"Contribute":0.000127},{"Turbine":"B4","Category":1.5,"Contribute":0.000381},{"Turbine":"B4","Category":2,"Contribute":0.000339},{"Turbine":"B4","Category":2.5,"Contribute":0.000508},{"Turbine":"B4","Category":3,"Contribute":0.000169},{"Turbine":"B4","Category":3.5,"Contribute":0.00055},{"Turbine":"B4","Category":4,"Contribute":0.000974},{"Turbine":"B4","Category":4.5,"Contribute":0.000593},{"Turbine":"B4","Category":5,"Contribute":0.001862},{"Turbine":"B4","Category":5.5,"Contribute":0.001905},{"Turbine":"B4","Category":6,"Contribute":0.002751},{"Turbine":"B4","Category":6.5,"Contribute":0.002921},{"Turbine":"B4","Category":7,"Contribute":0.001778},{"Turbine":"B4","Category":7.5,"Contribute":0.002667},{"Turbine":"B4","Category":8,"Contribute":0.003175},{"Turbine":"B4","Category":8.5,"Contribute":0.002328},{"Turbine":"B4","Category":9,"Contribute":0.001608},{"Turbine":"B4","Category":9.5,"Contribute":0.001566},{"Turbine":"B4","Category":10,"Contribute":0.001228},{"Turbine":"B4","Category":10.5,"Contribute":0.000254},{"Turbine":"B4","Category":11,"Contribute":0.000127},{"Turbine":"B71","Category":0,"Contribute":0.012275},{"Turbine":"B71","Category":1,"Contribute":0.000423},{"Turbine":"B71","Category":1.5,"Contribute":0.000466},{"Turbine":"B71","Category":2,"Contribute":0.000169},{"Turbine":"B71","Category":2.5,"Contribute":0.000339},{"Turbine":"B71","Category":3,"Contribute":0.000296},{"Turbine":"B71","Category":3.5,"Contribute":0.000381},{"Turbine":"B71","Category":4,"Contribute":0.000423},{"Turbine":"B71","Category":4.5,"Contribute":0.00072},{"Turbine":"B71","Category":5,"Contribute":0.001651},{"Turbine":"B71","Category":5.5,"Contribute":0.001439},{"Turbine":"B71","Category":6,"Contribute":0.001947},{"Turbine":"B71","Category":6.5,"Contribute":0.002243},{"Turbine":"B71","Category":7,"Contribute":0.001989},{"Turbine":"B71","Category":7.5,"Contribute":0.001862},{"Turbine":"B71","Category":8,"Contribute":0.003175},{"Turbine":"B71","Category":8.5,"Contribute":0.002624},{"Turbine":"B71","Category":9,"Contribute":0.002328},{"Turbine":"B71","Category":9.5,"Contribute":0.001947},{"Turbine":"B71","Category":10,"Contribute":0.000931},{"Turbine":"B71","Category":10.5,"Contribute":0.001228},{"Turbine":"B71","Category":11,"Contribute":0.000508},{"Turbine":"B71","Category":11.5,"Contribute":0.000339},{"Turbine":"B71","Category":12,"Contribute":0.000212},{"Turbine":"B71","Category":12.5,"Contribute":8.5E-05},{"Turbine":"B72","Category":0,"Contribute":0.012063},{"Turbine":"B72","Category":1,"Contribute":0.000169},{"Turbine":"B72","Category":1.5,"Contribute":0.000423},{"Turbine":"B72","Category":2,"Contribute":0.000339},{"Turbine":"B72","Category":2.5,"Contribute":0.000466},{"Turbine":"B72","Category":3,"Contribute":0.000339},{"Turbine":"B72","Category":3.5,"Contribute":0.000254},{"Turbine":"B72","Category":4,"Contribute":0.000804},{"Turbine":"B72","Category":4.5,"Contribute":0.00055},{"Turbine":"B72","Category":5,"Contribute":0.000931},{"Turbine":"B72","Category":5.5,"Contribute":0.002159},{"Turbine":"B72","Category":6,"Contribute":0.002243},{"Turbine":"B72","Category":6.5,"Contribute":0.002328},{"Turbine":"B72","Category":7,"Contribute":0.002836},{"Turbine":"B72","Category":7.5,"Contribute":0.001862},{"Turbine":"B72","Category":8,"Contribute":0.002328},{"Turbine":"B72","Category":8.5,"Contribute":0.002963},{"Turbine":"B72","Category":9,"Contribute":0.002497},{"Turbine":"B72","Category":9.5,"Contribute":0.002328},{"Turbine":"B72","Category":10,"Contribute":0.001058},{"Turbine":"B72","Category":10.5,"Contribute":0.00072},{"Turbine":"B72","Category":11,"Contribute":0.000254},{"Turbine":"B72","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B72","Category":12,"Contribute":4.2E-05},{"Turbine":"B73","Category":0,"Contribute":0.011725},{"Turbine":"B73","Category":1,"Contribute":0.000127},{"Turbine":"B73","Category":1.5,"Contribute":0.000254},{"Turbine":"B73","Category":2,"Contribute":0.000381},{"Turbine":"B73","Category":2.5,"Contribute":0.000381},{"Turbine":"B73","Category":3,"Contribute":0.000677},{"Turbine":"B73","Category":3.5,"Contribute":0.000339},{"Turbine":"B73","Category":4,"Contribute":0.000423},{"Turbine":"B73","Category":4.5,"Contribute":0.000847},{"Turbine":"B73","Category":5,"Contribute":0.00127},{"Turbine":"B73","Category":5.5,"Contribute":0.001439},{"Turbine":"B73","Category":6,"Contribute":0.001566},{"Turbine":"B73","Category":6.5,"Contribute":0.002074},{"Turbine":"B73","Category":7,"Contribute":0.002582},{"Turbine":"B73","Category":7.5,"Contribute":0.002201},{"Turbine":"B73","Category":8,"Contribute":0.002709},{"Turbine":"B73","Category":8.5,"Contribute":0.003005},{"Turbine":"B73","Category":9,"Contribute":0.002582},{"Turbine":"B73","Category":9.5,"Contribute":0.002159},{"Turbine":"B73","Category":10,"Contribute":0.001905},{"Turbine":"B73","Category":10.5,"Contribute":0.000635},{"Turbine":"B73","Category":11,"Contribute":0.000381},{"Turbine":"B73","Category":11.5,"Contribute":0.000212},{"Turbine":"B73","Category":12,"Contribute":4.2E-05},{"Turbine":"B73","Category":12.5,"Contribute":4.2E-05},{"Turbine":"B73","Category":14.5,"Contribute":4.2E-05},{"Turbine":"B75","Category":0,"Contribute":0.012233},{"Turbine":"B75","Category":1,"Contribute":0.000296},{"Turbine":"B75","Category":1.5,"Contribute":8.5E-05},{"Turbine":"B75","Category":2,"Contribute":0.00055},{"Turbine":"B75","Category":2.5,"Contribute":0.000339},{"Turbine":"B75","Category":3,"Contribute":0.000381},{"Turbine":"B75","Category":3.5,"Contribute":0.00055},{"Turbine":"B75","Category":4,"Contribute":0.000593},{"Turbine":"B75","Category":4.5,"Contribute":0.000974},{"Turbine":"B75","Category":5,"Contribute":0.001439},{"Turbine":"B75","Category":5.5,"Contribute":0.001397},{"Turbine":"B75","Category":6,"Contribute":0.002116},{"Turbine":"B75","Category":6.5,"Contribute":0.002497},{"Turbine":"B75","Category":7,"Contribute":0.002286},{"Turbine":"B75","Category":7.5,"Contribute":0.002624},{"Turbine":"B75","Category":8,"Contribute":0.002328},{"Turbine":"B75","Category":8.5,"Contribute":0.002455},{"Turbine":"B75","Category":9,"Contribute":0.001989},{"Turbine":"B75","Category":9.5,"Contribute":0.002243},{"Turbine":"B75","Category":10,"Contribute":0.001397},{"Turbine":"B75","Category":10.5,"Contribute":0.000593},{"Turbine":"B75","Category":11,"Contribute":0.000296},{"Turbine":"B75","Category":11.5,"Contribute":0.000254},{"Turbine":"B75","Category":12,"Contribute":8.5E-05},{"Turbine":"B77","Category":0,"Contribute":0.012275},{"Turbine":"B77","Category":1,"Contribute":0.000127},{"Turbine":"B77","Category":1.5,"Contribute":0.00055},{"Turbine":"B77","Category":2,"Contribute":0.00055},{"Turbine":"B77","Category":2.5,"Contribute":0.000254},{"Turbine":"B77","Category":3,"Contribute":0.000339},{"Turbine":"B77","Category":3.5,"Contribute":0.000339},{"Turbine":"B77","Category":4,"Contribute":0.000339},{"Turbine":"B77","Category":4.5,"Contribute":0.001058},{"Turbine":"B77","Category":5,"Contribute":0.001185},{"Turbine":"B77","Category":5.5,"Contribute":0.001101},{"Turbine":"B77","Category":6,"Contribute":0.002116},{"Turbine":"B77","Category":6.5,"Contribute":0.00237},{"Turbine":"B77","Category":7,"Contribute":0.002413},{"Turbine":"B77","Category":7.5,"Contribute":0.002116},{"Turbine":"B77","Category":8,"Contribute":0.002413},{"Turbine":"B77","Category":8.5,"Contribute":0.002243},{"Turbine":"B77","Category":9,"Contribute":0.002286},{"Turbine":"B77","Category":9.5,"Contribute":0.001862},{"Turbine":"B77","Category":10,"Contribute":0.00182},{"Turbine":"B77","Category":10.5,"Contribute":0.001354},{"Turbine":"B77","Category":11,"Contribute":0.000593},{"Turbine":"B77","Category":11.5,"Contribute":0.000296},{"Turbine":"B78","Category":0,"Contribute":0.012275},{"Turbine":"B78","Category":1,"Contribute":4.2E-05},{"Turbine":"B78","Category":1.5,"Contribute":0.000296},{"Turbine":"B78","Category":2,"Contribute":0.000381},{"Turbine":"B78","Category":2.5,"Contribute":0.000593},{"Turbine":"B78","Category":3,"Contribute":0.000508},{"Turbine":"B78","Category":3.5,"Contribute":0.000466},{"Turbine":"B78","Category":4,"Contribute":0.000762},{"Turbine":"B78","Category":4.5,"Contribute":0.001016},{"Turbine":"B78","Category":5,"Contribute":0.001101},{"Turbine":"B78","Category":5.5,"Contribute":0.001693},{"Turbine":"B78","Category":6,"Contribute":0.001524},{"Turbine":"B78","Category":6.5,"Contribute":0.001905},{"Turbine":"B78","Category":7,"Contribute":0.002921},{"Turbine":"B78","Category":7.5,"Contribute":0.003513},{"Turbine":"B78","Category":8,"Contribute":0.002582},{"Turbine":"B78","Category":8.5,"Contribute":0.002751},{"Turbine":"B78","Category":9,"Contribute":0.001693},{"Turbine":"B78","Category":9.5,"Contribute":0.002159},{"Turbine":"B78","Category":10,"Contribute":0.001185},{"Turbine":"B78","Category":10.5,"Contribute":0.000466},{"Turbine":"B78","Category":11,"Contribute":8.5E-05},{"Turbine":"B78","Category":11.5,"Contribute":8.5E-05},{"Turbine":"B79","Category":0,"Contribute":0.011767},{"Turbine":"B79","Category":0.5,"Contribute":8.5E-05},{"Turbine":"B79","Category":1,"Contribute":0.000212},{"Turbine":"B79","Category":1.5,"Contribute":0.000339},{"Turbine":"B79","Category":2,"Contribute":0.000296},{"Turbine":"B79","Category":2.5,"Contribute":0.000466},{"Turbine":"B79","Category":3,"Contribute":0.000296},{"Turbine":"B79","Category":3.5,"Contribute":0.000508},{"Turbine":"B79","Category":4,"Contribute":0.000762},{"Turbine":"B79","Category":4.5,"Contribute":0.000593},{"Turbine":"B79","Category":5,"Contribute":0.001354},{"Turbine":"B79","Category":5.5,"Contribute":0.001566},{"Turbine":"B79","Category":6,"Contribute":0.002032},{"Turbine":"B79","Category":6.5,"Contribute":0.002709},{"Turbine":"B79","Category":7,"Contribute":0.002201},{"Turbine":"B79","Category":7.5,"Contribute":0.002116},{"Turbine":"B79","Category":8,"Contribute":0.002286},{"Turbine":"B79","Category":8.5,"Contribute":0.002328},{"Turbine":"B79","Category":9,"Contribute":0.002582},{"Turbine":"B79","Category":9.5,"Contribute":0.002201},{"Turbine":"B79","Category":10,"Contribute":0.001481},{"Turbine":"B79","Category":10.5,"Contribute":0.000847},{"Turbine":"B79","Category":11,"Contribute":0.000296},{"Turbine":"B79","Category":11.5,"Contribute":0.000339},{"Turbine":"B79","Category":12,"Contribute":0.000169},{"Turbine":"B79","Category":13,"Contribute":4.2E-05},{"Turbine":"B79","Category":13.5,"Contribute":4.2E-05},{"Turbine":"B79","Category":14.5,"Contribute":4.2E-05},{"Turbine":"B79","Category":15,"Contribute":4.2E-05},{"Turbine":"B80","Category":0,"Contribute":0.012148},{"Turbine":"B80","Category":0.5,"Contribute":4.2E-05},{"Turbine":"B80","Category":1,"Contribute":0.000127},{"Turbine":"B80","Category":1.5,"Contribute":0.000508},{"Turbine":"B80","Category":2,"Contribute":0.000254},{"Turbine":"B80","Category":2.5,"Contribute":0.00055},{"Turbine":"B80","Category":3,"Contribute":0.000254},{"Turbine":"B80","Category":3.5,"Contribute":0.000254},{"Turbine":"B80","Category":4,"Contribute":0.000593},{"Turbine":"B80","Category":4.5,"Contribute":0.000762},{"Turbine":"B80","Category":5,"Contribute":0.001185},{"Turbine":"B80","Category":5.5,"Contribute":0.001481},{"Turbine":"B80","Category":6,"Contribute":0.001524},{"Turbine":"B80","Category":6.5,"Contribute":0.002074},{"Turbine":"B80","Category":7,"Contribute":0.00237},{"Turbine":"B80","Category":7.5,"Contribute":0.002328},{"Turbine":"B80","Category":8,"Contribute":0.002286},{"Turbine":"B80","Category":8.5,"Contribute":0.002497},{"Turbine":"B80","Category":9,"Contribute":0.002243},{"Turbine":"B80","Category":9.5,"Contribute":0.003132},{"Turbine":"B80","Category":10,"Contribute":0.002116},{"Turbine":"B80","Category":10.5,"Contribute":0.00072},{"Turbine":"B80","Category":11,"Contribute":0.000508},{"Turbine":"B80","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B82","Category":0,"Contribute":0.012233},{"Turbine":"B82","Category":1,"Contribute":0.000254},{"Turbine":"B82","Category":1.5,"Contribute":0.000423},{"Turbine":"B82","Category":2,"Contribute":0.000296},{"Turbine":"B82","Category":2.5,"Contribute":0.000508},{"Turbine":"B82","Category":3,"Contribute":0.000339},{"Turbine":"B82","Category":3.5,"Contribute":0.000423},{"Turbine":"B82","Category":4,"Contribute":0.000381},{"Turbine":"B82","Category":4.5,"Contribute":0.001354},{"Turbine":"B82","Category":5,"Contribute":0.001397},{"Turbine":"B82","Category":5.5,"Contribute":0.001651},{"Turbine":"B82","Category":6,"Contribute":0.00237},{"Turbine":"B82","Category":6.5,"Contribute":0.002032},{"Turbine":"B82","Category":7,"Contribute":0.002921},{"Turbine":"B82","Category":7.5,"Contribute":0.002328},{"Turbine":"B82","Category":8,"Contribute":0.002201},{"Turbine":"B82","Category":8.5,"Contribute":0.003259},{"Turbine":"B82","Category":9,"Contribute":0.002624},{"Turbine":"B82","Category":9.5,"Contribute":0.001185},{"Turbine":"B82","Category":10,"Contribute":0.000889},{"Turbine":"B82","Category":10.5,"Contribute":0.00055},{"Turbine":"B82","Category":11,"Contribute":0.000296},{"Turbine":"B82","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B82","Category":12.5,"Contribute":4.2E-05},{"Turbine":"B83","Category":0,"Contribute":0.011598},{"Turbine":"B83","Category":1,"Contribute":0.000169},{"Turbine":"B83","Category":1.5,"Contribute":0.000254},{"Turbine":"B83","Category":2,"Contribute":0.000508},{"Turbine":"B83","Category":2.5,"Contribute":0.000381},{"Turbine":"B83","Category":3,"Contribute":0.00055},{"Turbine":"B83","Category":3.5,"Contribute":0.000466},{"Turbine":"B83","Category":4,"Contribute":0.00072},{"Turbine":"B83","Category":4.5,"Contribute":0.000847},{"Turbine":"B83","Category":5,"Contribute":0.001058},{"Turbine":"B83","Category":5.5,"Contribute":0.001354},{"Turbine":"B83","Category":6,"Contribute":0.002328},{"Turbine":"B83","Category":6.5,"Contribute":0.002413},{"Turbine":"B83","Category":7,"Contribute":0.002794},{"Turbine":"B83","Category":7.5,"Contribute":0.002455},{"Turbine":"B83","Category":8,"Contribute":0.002582},{"Turbine":"B83","Category":8.5,"Contribute":0.002751},{"Turbine":"B83","Category":9,"Contribute":0.002413},{"Turbine":"B83","Category":9.5,"Contribute":0.001439},{"Turbine":"B83","Category":10,"Contribute":0.001185},{"Turbine":"B83","Category":10.5,"Contribute":0.000466},{"Turbine":"B83","Category":11,"Contribute":0.000508},{"Turbine":"B83","Category":11.5,"Contribute":0.000254},{"Turbine":"B83","Category":12,"Contribute":0.000254},{"Turbine":"B83","Category":12.5,"Contribute":0.000169},{"Turbine":"B83","Category":13,"Contribute":4.2E-05},{"Turbine":"B83","Category":15,"Contribute":4.2E-05},{"Turbine":"B84","Category":0,"Contribute":0.012148},{"Turbine":"B84","Category":1,"Contribute":0.000127},{"Turbine":"B84","Category":1.5,"Contribute":0.000254},{"Turbine":"B84","Category":2,"Contribute":0.000593},{"Turbine":"B84","Category":2.5,"Contribute":0.000508},{"Turbine":"B84","Category":3,"Contribute":0.000466},{"Turbine":"B84","Category":3.5,"Contribute":0.000381},{"Turbine":"B84","Category":4,"Contribute":0.000847},{"Turbine":"B84","Category":4.5,"Contribute":0.001016},{"Turbine":"B84","Category":5,"Contribute":0.001397},{"Turbine":"B84","Category":5.5,"Contribute":0.001693},{"Turbine":"B84","Category":6,"Contribute":0.001989},{"Turbine":"B84","Category":6.5,"Contribute":0.003048},{"Turbine":"B84","Category":7,"Contribute":0.00309},{"Turbine":"B84","Category":7.5,"Contribute":0.002751},{"Turbine":"B84","Category":8,"Contribute":0.002455},{"Turbine":"B84","Category":8.5,"Contribute":0.002878},{"Turbine":"B84","Category":9,"Contribute":0.001354},{"Turbine":"B84","Category":9.5,"Contribute":0.001439},{"Turbine":"B84","Category":10,"Contribute":0.000847},{"Turbine":"B84","Category":10.5,"Contribute":0.000508},{"Turbine":"B84","Category":11,"Contribute":8.5E-05},{"Turbine":"B84","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B84","Category":12,"Contribute":8.5E-05},{"Turbine":"B85","Category":0,"Contribute":0.011683},{"Turbine":"B85","Category":1,"Contribute":0.000212},{"Turbine":"B85","Category":1.5,"Contribute":0.000296},{"Turbine":"B85","Category":2,"Contribute":0.000508},{"Turbine":"B85","Category":2.5,"Contribute":0.000466},{"Turbine":"B85","Category":3,"Contribute":0.000169},{"Turbine":"B85","Category":3.5,"Contribute":0.000296},{"Turbine":"B85","Category":4,"Contribute":0.000635},{"Turbine":"B85","Category":4.5,"Contribute":0.000931},{"Turbine":"B85","Category":5,"Contribute":0.001101},{"Turbine":"B85","Category":5.5,"Contribute":0.001143},{"Turbine":"B85","Category":6,"Contribute":0.00182},{"Turbine":"B85","Category":6.5,"Contribute":0.001905},{"Turbine":"B85","Category":7,"Contribute":0.002328},{"Turbine":"B85","Category":7.5,"Contribute":0.002286},{"Turbine":"B85","Category":8,"Contribute":0.002455},{"Turbine":"B85","Category":8.5,"Contribute":0.002286},{"Turbine":"B85","Category":9,"Contribute":0.002582},{"Turbine":"B85","Category":9.5,"Contribute":0.002328},{"Turbine":"B85","Category":10,"Contribute":0.001439},{"Turbine":"B85","Category":10.5,"Contribute":0.001397},{"Turbine":"B85","Category":11,"Contribute":0.000677},{"Turbine":"B85","Category":11.5,"Contribute":0.000508},{"Turbine":"B85","Category":12,"Contribute":8.5E-05},{"Turbine":"B85","Category":12.5,"Contribute":0.000127},{"Turbine":"B85","Category":13,"Contribute":8.5E-05},{"Turbine":"B85","Category":13.5,"Contribute":0.000127},{"Turbine":"B85","Category":14,"Contribute":4.2E-05},{"Turbine":"B85","Category":14.5,"Contribute":8.5E-05},{"Turbine":"B86","Category":0,"Contribute":0.012021},{"Turbine":"B86","Category":1,"Contribute":8.5E-05},{"Turbine":"B86","Category":1.5,"Contribute":0.00055},{"Turbine":"B86","Category":2,"Contribute":0.000212},{"Turbine":"B86","Category":2.5,"Contribute":0.000466},{"Turbine":"B86","Category":3,"Contribute":0.000466},{"Turbine":"B86","Category":3.5,"Contribute":0.000508},{"Turbine":"B86","Category":4,"Contribute":0.000508},{"Turbine":"B86","Category":4.5,"Contribute":0.000762},{"Turbine":"B86","Category":5,"Contribute":0.001016},{"Turbine":"B86","Category":5.5,"Contribute":0.001481},{"Turbine":"B86","Category":6,"Contribute":0.001566},{"Turbine":"B86","Category":6.5,"Contribute":0.002878},{"Turbine":"B86","Category":7,"Contribute":0.002159},{"Turbine":"B86","Category":7.5,"Contribute":0.001989},{"Turbine":"B86","Category":8,"Contribute":0.00237},{"Turbine":"B86","Category":8.5,"Contribute":0.002582},{"Turbine":"B86","Category":9,"Contribute":0.002794},{"Turbine":"B86","Category":9.5,"Contribute":0.002709},{"Turbine":"B86","Category":10,"Contribute":0.001947},{"Turbine":"B86","Category":10.5,"Contribute":0.00072},{"Turbine":"B86","Category":11,"Contribute":8.5E-05},{"Turbine":"B86","Category":11.5,"Contribute":4.2E-05},{"Turbine":"B86","Category":12,"Contribute":8.5E-05},{"Turbine":"B87","Category":0,"Contribute":0.012106},{"Turbine":"B87","Category":0.5,"Contribute":8.5E-05},{"Turbine":"B87","Category":1,"Contribute":0.000254},{"Turbine":"B87","Category":1.5,"Contribute":0.000381},{"Turbine":"B87","Category":2,"Contribute":0.000254},{"Turbine":"B87","Category":2.5,"Contribute":0.000296},{"Turbine":"B87","Category":3,"Contribute":0.000381},{"Turbine":"B87","Category":3.5,"Contribute":0.000339},{"Turbine":"B87","Category":4,"Contribute":0.00055},{"Turbine":"B87","Category":4.5,"Contribute":0.000931},{"Turbine":"B87","Category":5,"Contribute":0.001143},{"Turbine":"B87","Category":5.5,"Contribute":0.001566},{"Turbine":"B87","Category":6,"Contribute":0.002328},{"Turbine":"B87","Category":6.5,"Contribute":0.00237},{"Turbine":"B87","Category":7,"Contribute":0.001735},{"Turbine":"B87","Category":7.5,"Contribute":0.002878},{"Turbine":"B87","Category":8,"Contribute":0.00182},{"Turbine":"B87","Category":8.5,"Contribute":0.001778},{"Turbine":"B87","Category":9,"Contribute":0.002582},{"Turbine":"B87","Category":9.5,"Contribute":0.002328},{"Turbine":"B87","Category":10,"Contribute":0.001524},{"Turbine":"B87","Category":10.5,"Contribute":0.001101},{"Turbine":"B87","Category":11,"Contribute":0.000677},{"Turbine":"B87","Category":11.5,"Contribute":0.000339},{"Turbine":"B87","Category":12,"Contribute":0.000127},{"Turbine":"B87","Category":12.5,"Contribute":8.5E-05},{"Turbine":"B87","Category":14,"Contribute":4.2E-05},{"Turbine":"B89","Category":0,"Contribute":0.012952},{"Turbine":"B89","Category":1.5,"Contribute":0.000296},{"Turbine":"B89","Category":2,"Contribute":4.2E-05},{"Turbine":"B89","Category":2.5,"Contribute":0.000339},{"Turbine":"B89","Category":3,"Contribute":0.000169},{"Turbine":"B89","Category":3.5,"Contribute":0.000212},{"Turbine":"B89","Category":4,"Contribute":0.000593},{"Turbine":"B89","Category":4.5,"Contribute":0.001101},{"Turbine":"B89","Category":5,"Contribute":0.001185},{"Turbine":"B89","Category":5.5,"Contribute":0.001566},{"Turbine":"B89","Category":6,"Contribute":0.003048},{"Turbine":"B89","Category":6.5,"Contribute":0.002836},{"Turbine":"B89","Category":7,"Contribute":0.001778},{"Turbine":"B89","Category":7.5,"Contribute":0.001651},{"Turbine":"B89","Category":8,"Contribute":0.00237},{"Turbine":"B89","Category":8.5,"Contribute":0.002963},{"Turbine":"B89","Category":9,"Contribute":0.002243},{"Turbine":"B89","Category":9.5,"Contribute":0.00182},{"Turbine":"B89","Category":10,"Contribute":0.001524},{"Turbine":"B89","Category":10.5,"Contribute":0.000847},{"Turbine":"B89","Category":11,"Contribute":0.000254},{"Turbine":"B89","Category":11.5,"Contribute":0.000169},{"Turbine":"B89","Category":13,"Contribute":4.2E-05},{"Turbine":"B90","Category":0,"Contribute":0.01236},{"Turbine":"B90","Category":1,"Contribute":0.000169},{"Turbine":"B90","Category":1.5,"Contribute":0.000339},{"Turbine":"B90","Category":2,"Contribute":0.000339},{"Turbine":"B90","Category":2.5,"Contribute":0.000508},{"Turbine":"B90","Category":3,"Contribute":0.000212},{"Turbine":"B90","Category":3.5,"Contribute":0.000212},{"Turbine":"B90","Category":4,"Contribute":0.000677},{"Turbine":"B90","Category":4.5,"Contribute":0.001016},{"Turbine":"B90","Category":5,"Contribute":0.001312},{"Turbine":"B90","Category":5.5,"Contribute":0.002328},{"Turbine":"B90","Category":6,"Contribute":0.001989},{"Turbine":"B90","Category":6.5,"Contribute":0.002413},{"Turbine":"B90","Category":7,"Contribute":0.002667},{"Turbine":"B90","Category":7.5,"Contribute":0.002624},{"Turbine":"B90","Category":8,"Contribute":0.002751},{"Turbine":"B90","Category":8.5,"Contribute":0.003132},{"Turbine":"B90","Category":9,"Contribute":0.001778},{"Turbine":"B90","Category":9.5,"Contribute":0.001735},{"Turbine":"B90","Category":10,"Contribute":0.000762},{"Turbine":"B90","Category":10.5,"Contribute":0.000296},{"Turbine":"B90","Category":11,"Contribute":0.000127},{"Turbine":"B90","Category":11.5,"Contribute":0.000169},{"Turbine":"B90","Category":12,"Contribute":8.5E-05},{"Turbine":"B91","Category":0,"Contribute":0.012148},{"Turbine":"B91","Category":1,"Contribute":0.000127},{"Turbine":"B91","Category":1.5,"Contribute":0.000423},{"Turbine":"B91","Category":2,"Contribute":0.000296},{"Turbine":"B91","Category":2.5,"Contribute":0.00055},{"Turbine":"B91","Category":3,"Contribute":0.000254},{"Turbine":"B91","Category":3.5,"Contribute":0.000296},{"Turbine":"B91","Category":4,"Contribute":0.00055},{"Turbine":"B91","Category":4.5,"Contribute":0.000931},{"Turbine":"B91","Category":5,"Contribute":0.001862},{"Turbine":"B91","Category":5.5,"Contribute":0.001439},{"Turbine":"B91","Category":6,"Contribute":0.001228},{"Turbine":"B91","Category":6.5,"Contribute":0.002074},{"Turbine":"B91","Category":7,"Contribute":0.002328},{"Turbine":"B91","Category":7.5,"Contribute":0.00237},{"Turbine":"B91","Category":8,"Contribute":0.00182},{"Turbine":"B91","Category":8.5,"Contribute":0.002878},{"Turbine":"B91","Category":9,"Contribute":0.00254},{"Turbine":"B91","Category":9.5,"Contribute":0.002159},{"Turbine":"B91","Category":10,"Contribute":0.001989},{"Turbine":"B91","Category":10.5,"Contribute":0.00127},{"Turbine":"B91","Category":11,"Contribute":0.000169},{"Turbine":"B91","Category":11.5,"Contribute":0.000127},{"Turbine":"B91","Category":12,"Contribute":0.000169},{"Turbine":"B92","Category":0,"Contribute":0.01236},{"Turbine":"B92","Category":1,"Contribute":0.000127},{"Turbine":"B92","Category":1.5,"Contribute":0.000508},{"Turbine":"B92","Category":2,"Contribute":0.000508},{"Turbine":"B92","Category":2.5,"Contribute":0.000127},{"Turbine":"B92","Category":3,"Contribute":0.000254},{"Turbine":"B92","Category":3.5,"Contribute":0.000339},{"Turbine":"B92","Category":4,"Contribute":0.000762},{"Turbine":"B92","Category":4.5,"Contribute":0.000804},{"Turbine":"B92","Category":5,"Contribute":0.000889},{"Turbine":"B92","Category":5.5,"Contribute":0.00182},{"Turbine":"B92","Category":6,"Contribute":0.002328},{"Turbine":"B92","Category":6.5,"Contribute":0.002074},{"Turbine":"B92","Category":7,"Contribute":0.002286},{"Turbine":"B92","Category":7.5,"Contribute":0.002201},{"Turbine":"B92","Category":8,"Contribute":0.002243},{"Turbine":"B92","Category":8.5,"Contribute":0.002455},{"Turbine":"B92","Category":9,"Contribute":0.002159},{"Turbine":"B92","Category":9.5,"Contribute":0.00237},{"Turbine":"B92","Category":10,"Contribute":0.001862},{"Turbine":"B92","Category":10.5,"Contribute":0.001016},{"Turbine":"B92","Category":11,"Contribute":0.000212},{"Turbine":"B92","Category":11.5,"Contribute":0.000169},{"Turbine":"B92","Category":12,"Contribute":4.2E-05},{"Turbine":"B92","Category":12.5,"Contribute":8.5E-05},{"Turbine":"T1","Category":0,"Contribute":0.013841},{"Turbine":"T1","Category":1,"Contribute":0.000127},{"Turbine":"T1","Category":1.5,"Contribute":0.000212},{"Turbine":"T1","Category":2,"Contribute":8.5E-05},{"Turbine":"T1","Category":2.5,"Contribute":0.000169},{"Turbine":"T1","Category":3,"Contribute":0.000296},{"Turbine":"T1","Category":3.5,"Contribute":0.000169},{"Turbine":"T1","Category":4,"Contribute":0.000381},{"Turbine":"T1","Category":4.5,"Contribute":0.000635},{"Turbine":"T1","Category":5,"Contribute":0.000762},{"Turbine":"T1","Category":5.5,"Contribute":0.001651},{"Turbine":"T1","Category":6,"Contribute":0.001862},{"Turbine":"T1","Category":6.5,"Contribute":0.002328},{"Turbine":"T1","Category":7,"Contribute":0.002116},{"Turbine":"T1","Category":7.5,"Contribute":0.001947},{"Turbine":"T1","Category":8,"Contribute":0.002286},{"Turbine":"T1","Category":8.5,"Contribute":0.002243},{"Turbine":"T1","Category":9,"Contribute":0.002582},{"Turbine":"T1","Category":9.5,"Contribute":0.002116},{"Turbine":"T1","Category":10,"Contribute":0.002159},{"Turbine":"T1","Category":10.5,"Contribute":0.000847},{"Turbine":"T1","Category":11,"Contribute":0.000635},{"Turbine":"T1","Category":11.5,"Contribute":0.000296},{"Turbine":"T1","Category":12,"Contribute":0.000169},{"Turbine":"T1","Category":12.5,"Contribute":4.2E-05},{"Turbine":"T1","Category":13,"Contribute":4.2E-05},{"Turbine":"T2","Category":0,"Contribute":0.012275},{"Turbine":"T2","Category":1,"Contribute":0.000296},{"Turbine":"T2","Category":1.5,"Contribute":0.000381},{"Turbine":"T2","Category":2,"Contribute":0.000423},{"Turbine":"T2","Category":2.5,"Contribute":0.000296},{"Turbine":"T2","Category":3,"Contribute":0.000339},{"Turbine":"T2","Category":3.5,"Contribute":0.000423},{"Turbine":"T2","Category":4,"Contribute":0.000423},{"Turbine":"T2","Category":4.5,"Contribute":0.000593},{"Turbine":"T2","Category":5,"Contribute":0.001101},{"Turbine":"T2","Category":5.5,"Contribute":0.001312},{"Turbine":"T2","Category":6,"Contribute":0.002286},{"Turbine":"T2","Category":6.5,"Contribute":0.00237},{"Turbine":"T2","Category":7,"Contribute":0.001651},{"Turbine":"T2","Category":7.5,"Contribute":0.002159},{"Turbine":"T2","Category":8,"Contribute":0.003344},{"Turbine":"T2","Category":8.5,"Contribute":0.002751},{"Turbine":"T2","Category":9,"Contribute":0.002582},{"Turbine":"T2","Category":9.5,"Contribute":0.002286},{"Turbine":"T2","Category":10,"Contribute":0.001397},{"Turbine":"T2","Category":10.5,"Contribute":0.000931},{"Turbine":"T2","Category":11,"Contribute":0.000254},{"Turbine":"T2","Category":11.5,"Contribute":8.5E-05},{"Turbine":"T2","Category":12,"Contribute":4.2E-05}]},"Result":"OK","Message":null,"Trace":null};

        app.loading(true);
        var param = {
            period: fa.period,
            dateStart: fa.dateStart,
            dateEnd: fa.dateEnd,
            turbine: fa.turbine(),
            project: fa.project
        };

        toolkit.ajaxPost(viewModel.appName + "analyticwinddistribution/getlist", param, function (res) {
            // console.log(res.data.Data);
            if (!app.isFine(res)) {
                app.loading(false);
                return;
            }

            if (wd.turbine().length == 0) {
                var turbine = []
                for (var i=0;i<res.data.Data.length;i++) {
                    if ($.inArray( res.data.Data[i].Turbine, turbine ) == -1){
                        turbine.push(res.data.Data[i].Turbine);
                    }
                }

                $.each(turbine, function (i, val) {
                    var data = {
                        color: color[i],
                        turbine: val
                    }

                    wd.turbineList.push(data);
                });
            }

            $('#windDistribution').html("");
            var data = res.data.Data;

            $("#windDistribution").kendoChart({
                dataSource: {
                    data: data,
                    group: { field: "Turbine" },
                    sort: { field: "Category", dir: 'asc' }
                },
                theme: "flat",
                title: {
                    text: ""
                },
                legend: {
                    position: "right",
                    visible: false
                },
                chartArea: {
                    height: 500
                },
                series: [{
                    type: "line",
                    style: "smooth",
                    field: "Contribute",
                    // opacity : 0.7,
                    markers: {
                        visible: true,
                        size: 3,
                    }
                }],
                seriesColors: color,
                valueAxis: {
                    labels: {
                        format: "{0:p0}",
                    },
                    line: {
                        visible: true
                    },
                    axisCrossingValue: -10,
                    majorGridLines: {
                        visible: true,
                        color: "#eee",
                        width: 0.8,
                    }
                },
                categoryAxis: {
                    field: "Category",
                    majorGridLines: {
                        visible: false
                    },
                    labels: {
                        // rotation: 25
                    },
                    majorTickType: "none"
                },
                tooltip: {
                    visible: true,
                    // template: "Contribution of #= series.name # : #= kendo.toString(value, 'n4')# % at #= category #",
                    template: "#= kendo.toString(value, 'p2')#",
                    // shared: true,
                    background: "rgb(255,255,255, 0.9)",
                    color: "#58666e",
                    font: 'Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
                    border: {
                        color: "#eee",
                        width: "2px",
                    },

                },
            });

            Data.InitRightTurbineList();

            app.loading(false);
            $("#windDistribution").data("kendoChart").refresh();
        });
    },
    InitRightTurbineList: function () {
        if (wd.turbine().length > 0) {
            wd.turbineList([]);
            $.each(wd.turbine(), function (i, val) {
                var data = {
                    color: color[i],
                    turbine: val
                }

                wd.turbineList.push(data);
            });
        }

        if (wd.turbineList().length > 1) {
            $("#showHideChk").html('<label>' +
                '<input type="checkbox" id="showHideAll" checked onclick="wd.showHideAllLegend(this)" >' +
                '<span class="cr"><i class="cr-icon glyphicon glyphicon-ok"></i></span>' +
                '<span id="labelShowHide"><b>Select All</b></span>' +
                '</label>');
        } else {
            $("#showHideChk").html("");
        }

        $("#right-turbine-list").html("");
        $.each(wd.turbineList(), function (idx, val) {
            $("#right-turbine-list").append('<div class="btn-group">' +
                '<button class="btn btn-default btn-sm turbine-chk" type="button" onclick="wd.showHideLegend(' + (idx) + ')" style="border-color:' + val.color + ';background-color:' + val.color + '"><i class="fa fa-check" id="icon-' + (idx) + '"></i></button>' +
                '<input class="chk-option" type="checkbox" name="' + val.turbine + '" checked id="chk-' + (idx) + '" hidden>' +
                '<button class="btn btn-default btn-sm turbine-btn" onclick="wd.showHideLegend(' + (idx) + ')" type="button" style="width:70px">' + val.turbine + '</button>' +
                '</div>');
        });
    }
}

wd.showHideAllLegend = function (e) {

    if (e.checked == true) {
        $('.fa-check').css("visibility", 'visible');
        $.each(wd.turbine(), function (i, val) {
            $("#windDistribution").data("kendoChart").options.series[i].visible = true;
        });
        /*$('#labelShowHide b').text('Untick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    } else {
        $.each(wd.turbine(), function (i, val) {
            $("#windDistribution").data("kendoChart").options.series[i].visible = false;
        });
        $('.fa-check').css("visibility", 'hidden');
        /*$('#labelShowHide b').text('Tick All Turbines');*/
        $('#labelShowHide b').text('Select All');
    }
    $('.chk-option').not(e).prop('checked', e.checked);

    $("#windDistribution").data("kendoChart").redraw();
}

wd.showHideLegend = function (idx) {
    var stat = false;

    $('#chk-' + idx).trigger('click');
    var chart = $("#windDistribution").data("kendoChart");
    var leTur = $('input[id*=chk-][type=checkbox]').length

    if ($('input[id*=chk-][type=checkbox]:checked').length == $('input[id*=chk-][type=checkbox]').length) {
        $('#showHideAll').prop('checked', true);
    } else {
        $('#showHideAll').prop('checked', false);
    }

    if ($('#chk-' + idx).is(':checked')) {
        $('#icon-' + idx).css("visibility", "visible");
    } else {
        $('#icon-' + idx).css("visibility", "hidden");
    }

    if ($('#chk-' + idx).is(':checked')) {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = true
    } else {
        $("#windDistribution").data("kendoChart").options.series[idx].visible = false
    }
    $("#windDistribution").data("kendoChart").redraw();
}



vm.currentMenu('Wind Distribution');
vm.currentTitle('Wind Distribution');
vm.breadcrumb([{ title: 'Analysis', href: '#' }, { title: 'Wind Analysis', href: '#' }, { title: 'Wind Distribution', href: viewModel.appName + 'page/analyticwinddistribution' }]);

$(document).ready(function () {
    $('#btnRefresh').on('click', function () {
        app.loading(true);
        setTimeout(function () {
            Data.LoadData();
        }, 1000);
    });

    app.loading(true);
    setTimeout(function () {
        Data.LoadData();
    }, 1000);

});

// $(function (){
//     Data.LoadData();
// });