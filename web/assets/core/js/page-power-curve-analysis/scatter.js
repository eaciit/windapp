'use strict';
var sc = {};
sc.loadFirstTime = ko.observable(true);
sc.reset = function(){
    sc.loadFirstTime(true);
}
sc.refresh = function() {
    sc.loadFirstTime(false);
    console.log('scatter refresh');
}