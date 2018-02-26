'use strict';
var sf = {};
sf.loadFirstTime = ko.observable(true);
sf.reset = function(){
    sf.loadFirstTime(true);
}
sf.refresh = function() {
    sf.loadFirstTime(false);
    console.log('scatter with filter refresh');
}