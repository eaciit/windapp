'use strict';
var ms = {};
ms.loadFirstTime = ko.observable(true);
ms.reset = function(){
    ms.loadFirstTime(true);
}
ms.refresh = function() {
    ms.loadFirstTime(false);
    console.log('monthly scatter refresh');
}