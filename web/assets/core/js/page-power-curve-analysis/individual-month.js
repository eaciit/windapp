'use strict';
var im = {};
im.loadFirstTime = ko.observable(true);
im.reset = function(){
    im.loadFirstTime(true);
}
im.refresh = function() {
    im.loadFirstTime(false);
    console.log('individual month refresh');
}