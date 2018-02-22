'use strict';
var op = {};
op.loadFirstTime = ko.observable(true);
op.reset = function(){
    op.loadFirstTime(true);
}
op.refresh = function() {
    op.loadFirstTime(false);
    console.log('operational refresh');
}