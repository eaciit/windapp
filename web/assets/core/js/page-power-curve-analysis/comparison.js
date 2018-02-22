'use strict';
var cm = {};
cm.loadFirstTime = ko.observable(true);
cm.reset = function(){
    cm.loadFirstTime(true);
}
cm.refresh = function() {
    cm.loadFirstTime(false);
    console.log('comparison refresh');
}