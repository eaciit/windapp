'use strict';

var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) { return typeof obj; } : function (obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol ? "symbol" : typeof obj; };

viewModel.app = new Object();
var app = viewModel.app;

app.dev = ko.observable(true);
app.noop = function () {};
app.noob = {};
app.log = function () {
    if (!app.dev()) {
        return;
    }

    console.log.apply(console, [].slice.call(arguments));
};
app.error = function () {
    if (!app.dev()) {
        return;
    }

    console.error.apply(console, [].slice.call(arguments));
};
app.validateNumber = function (d) {
    var df = arguments.length <= 1 || arguments[1] === undefined ? 0 : arguments[1];
    return isNaN(d) || !isFinite(d) ? df : d;
};
app.isLastItem = function (o, d) {
    return o.indexOf(d) + 1 == o.length;
};
app.NaNable = function (o) {
    var dv = arguments.length <= 1 || arguments[1] === undefined ? 0 : arguments[1];
    return isNaN(o) ? dv : o;
};
app.nbspAble = function (o) {
    var dv = arguments.length <= 1 || arguments[1] === undefined ? '&nbsp;' : arguments[1];
    return $.trim(o) == '' ? dv : o;
};
app.allKeys = function (o) {
    var keys = [];
    for (var k in o) {
        if (o.hasOwnProperty(k)) {
            keys.push(String(k));
        }
    }
    return keys;
};
app.length = function (o) {
    if (o instanceof Object) {
        var i = 0;
        for (var k in o) {
            if (o.hasOwnProperty(k)) {
                i++;
            }
        }
        return i;
    }

    return o.length;
};
app.getPropVal = function (o, key) {
    var dv = arguments.length <= 2 || arguments[2] === undefined ? null : arguments[2];

    if (!o.hasOwnProperty(key)) {
        return dv;
    }

    return app.isUndefined(o[key]) ? dv : o[key];
};
app.isVoid = function (o) {
    if (app.isUndefined(o)) {
        return true;
    }
    if (o == null) {
        return true;
    }
    if (typeof o == 'string') {
        if ($.trim(o) == '') {
            return true;
        }
    }

    return false;
};
app.whenVoid = function (o) {
    var df = arguments.length <= 1 || arguments[1] === undefined ? null : arguments[1];
    return app.isVoid(o) ? df : o;
};
app.hasProp = function (o, key) {
    return o.hasOwnProperty(key);
};
app.ajaxPost = function (url) {
    var data = arguments.length <= 1 || arguments[1] === undefined ? {} : arguments[1];
    var callbackSuccess = arguments.length <= 2 || arguments[2] === undefined ? app.noop : arguments[2];
    var callbackError = arguments.length <= 3 || arguments[3] === undefined ? app.noop : arguments[3];
    var otherConfig = arguments.length <= 4 || arguments[4] === undefined ? app.noob : arguments[4];

    var startReq = moment();

    var params = ko.mapping.toJSON(app.noob);
    try {
        params = ko.mapping.toJSON(data);
    } catch (err) {}

    var cache = app.getPropVal(otherConfig, 'cache', false);
    if (cache !== false) {
        if (app.hasProp(localStorage, cache)) {
            var _data = JSON.parse(localStorage[cache]);
            callbackSuccess(_data);
            return;
        }
    }

    var config = {
        url: url.toLowerCase(),
        type: 'post',
        dataType: 'json',
        contentType: 'application/json charset=utf-8',
        data: params,
        success: function success(a) {
            if (cache !== '') {
                a.time = moment.now();
                localStorage[cache] = JSON.stringify(a);
            }

            callbackSuccess(a);
        },
        error: function error(a, b, c) {
            callbackError(a, b, c);
        }
    };

    if (data instanceof FormData) {
        delete config.config;
        config.data = data;
        config.async = false;
        config.cache = false;
        config.contentType = false;
        config.processData = false;
    }

    config = $.extend(true, config, otherConfig);
    return $.ajax(config);
};
app.o = function (raw) {
    return raw;
};
app.seriesColorsGodrej = ['#3498DB', '#28B463', '#F39C12', '#DB3434', '#34D3DB'];
app.randomRange = function (min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
};
app.capitalize = function (d) {
    return '' + d[0].toUpperCase() + d.slice(1);
};
app.typeIs = function (target, comparator) {
    return (typeof target === 'undefined' ? 'undefined' : _typeof(target)) === comparator;
};
app.is = function (observable, comparator) {
    var a = typeof observable === 'function' ? observable() : observable;
    var b = typeof comparator === 'function' ? comparator() : comparator;

    return a === b;
};
app.isNot = function (observable, comparator) {
    var a = typeof observable === 'function' ? observable() : observable;
    var b = typeof comparator === 'function' ? comparator() : comparator;

    return a !== b;
};
app.isDefined = function (o) {
    return !app.isUndefined(o);
};
app.isUndefined = function (o) {
    return typeof o === 'undefined';
};
app.showError = function (message) {
    return sweetAlert('Warning', message, 'error');
};
app.isFine = function (res) {
    if (!res.success && res.message.indexOf('expired') > -1) {
        app.isLoading(false);
        swal({
            title: "Warning",
            type: "warning",
            text: res.message,
        }, function () {
            setTimeout(function () {
                location.href = viewModel.appName + 'page/login';
            }, 200);
        });
        return false;
    }
    if (!res.success && res.message != "data is empty") {
        sweetAlert('Warning', res.message, 'error');
        return false;
    }

    return true;
};
app.isFormValid = function (selector) {
    app.resetValidation(selector);
    var $validator = $(selector).data('kendoValidator');
    return $validator.validate();
};
app.resetValidation = function (selectorID) {
    var $form = $(selectorID).data('kendoValidator');
    if (!$form) {
        $(selectorID).kendoValidator();
        $form = $(selectorID).data('kendoValidator');
    }

    try {
        $form.hideMessages();
    } catch (err) {}
};
app.resetForm = function ($o) {
    $o.trigger('reset');
};
app.prepareTooltipster = function ($o, argConfig) {
    var $tooltipster = typeof $o === 'undefined' ? $('.tooltipster') : $o;

    $tooltipster.each(function (i, e) {
        var position = 'top';

        if ($(e).attr('class').search('tooltipster-') > -1) {
            position = $(e).attr('class').split(' ').find(function (d) {
                return d.search('tooltipster-') > -1;
            }).replace(/tooltipster\-/g, '');
        }

        var config = {
            theme: 'tooltipster-val',
            animation: 'grow',
            delay: 0,
            offsetY: -5,
            touchDevices: false,
            trigger: 'hover',
            position: position,
            content: $('<div />').html($(e).attr('title')),
            multiple: true
        };
        if (typeof argConfig !== 'undefined') {
            config = $.extend(true, config, argConfig);
        }

        $(e).tooltipster(config);
    });
};
app.gridBoundTooltipster = function (selector) {
    return function () {
        app.prepareTooltipster($(selector).find(".tooltipster"));
    };
};
app.redefine = function (o, d) {
    return typeof o === 'undefined' ? d : o;
};
app.capitalize = function (s) {
    var isHardcore = arguments.length <= 1 || arguments[1] === undefined ? false : arguments[1];

    s = app.redefine(s, '');

    if (isHardcore) {
        s = s.toLowerCase();
    }

    if (s.length == 0) {
        return '';
    }

    var res = s.split(' ').map(function (d) {
        return d.length > 0 ? d[0].toUpperCase() + d.slice(1) : 0;
    }).join(' ');
    return res;
};
app.repeatAlphabetically = function (prefix) {
    return 'abcdefghijklmnopqrstuvwxyz'.split('').map(function (d) {
        return prefix + ' ' + d.toUpperCase();
    });
};
app.arrRemoveByIndex = function (arr, index) {
    arr.splice(index, 1);
};
app.arrRemoveByItem = function (arr, item) {
    var index = arr.indexOf(item);
    if (index > -1) {
        app.arrRemoveByIndex(arr, index);
    }
};
app.clone = function (o) {
    return $.extend(true, {}, o);
};
app.distinct = function (arr) {
    return arr.filter(function (v, i, self) {
        return self.indexOf(v) === i;
    });
};
app.forEach = function (d, callback) {
    if (d instanceof Array) {
        d.forEach(callback);
    }

    if (d instanceof Object) {
        for (var key in d) {
            if (d.hasOwnProperty(key)) {
                callback(key, d[key]);
            }
        }
    }
};

app.koMap = ko.mapping.fromJS;
app.koUnmap = ko.mapping.toJS;
app.observ = ko.observable;
app.observArr = ko.observArr;

app.randomString = function () {
    var length = arguments.length <= 0 || arguments[0] === undefined ? 5 : arguments[0];
    return Math.random().toString(36).substring(2, length);
};

app.latLngIndonesia = { lat: -1.8504955, lng: 117.4004627 };
app.randomGeoLocations = function () {
    var center = arguments.length <= 0 || arguments[0] === undefined ? app.latLngIndonesia : arguments[0];
    var radius = arguments.length <= 1 || arguments[1] === undefined ? 1000000 : arguments[1];
    var count = arguments.length <= 2 || arguments[2] === undefined ? 100 : arguments[2];

    var generateRandomPoint = function generateRandomPoint(center, radius) {
        var x0 = center.lng;
        var y0 = center.lat;

        // Convert Radius from meters to degrees.
        var rd = radius / 111300;

        var u = Math.random();
        var v = Math.random();

        var w = rd * Math.sqrt(u);
        var t = 2 * Math.PI * v;
        var x = w * Math.cos(t);
        var y = w * Math.sin(t);

        var xp = x / Math.cos(y0);

        return {
            name: app.randomString(10),
            latlng: [y + y0, xp + x0]
        };
    };

    var points = [];
    for (var i = 0; i < count; i++) {
        points.push(generateRandomPoint(center, radius));
    }
    return points;
};

app.split = function (arr) {
    var separator = arguments.length <= 1 || arguments[1] === undefined ? '' : arguments[1];
    var length = arguments.length <= 2 || arguments[2] === undefined ? 0 : arguments[2];

    if (length == 0) {
        return arr.split(separator);
    }

    var res = [];
    var resJoin = [];

    arr.split(separator).forEach(function (d, i) {
        if (i < length) {
            res.push(d);
            return;
        }

        resJoin.push(d);
    });

    res = res.concat(resJoin.join(separator));
    return res;
};

app.extend = function (which, klass) {
    app.forEach(klass, function (key, val) {
        if (app.typeIs(val, 'function')) {
            var body = { value: val };

            if (app.typeIs(which, 'string')) {
                Object.defineProperty(window[which].prototype, key, body);
            } else {
                Object.defineProperty(target.prototype, key, body);
            }
        }
    });
};
app.newEl = function (s) {
    return $('<' + s + ' />');
};
app.idAble = function (s) {
    return s.replace(/\./g, '_').replace(/\-/g, '_').replace(/\//g, '_').replace(/ /g, '_');
};
app.logAble = function () {
    var args = [].slice.call(arguments);
    app.log(args);
    return args[0];
};
app.htmlDecode = function (s) {
    var elem = document.createElement('textarea');
    elem.innerHTML = s;
    return elem.value;
};
app.runAfter = function () {
    for (var _len = arguments.length, jobs = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
        jobs[_key - 1] = arguments[_key];
    }

    var delay = arguments.length <= 0 || arguments[0] === undefined ? 0 : arguments[0];

    var doWork = function doWork() {
        jobs.forEach(function (job) {
            job();
        });
    };

    var timeout = setTimeout(function () {
        return doWork;
    }, delay);
    return timeout;
};

app.resetLocalStorage = function(){
    var SummaryMode = localStorage.getItem('SummaryMode') ; 

    localStorage.clear();
    if(SummaryMode != null){
        localStorage.setItem('SummaryMode', SummaryMode);
    }
    
}

viewModel.StringExt = new Object();
var s = viewModel.StringExt;

s.toObject = function () {
    var source = String(this);
    try {
        return JSON.parse(source);
    } catch (err) {
        console.error(err);
        return {};
    }
};

app.isContentShow = ko.observable(true);
app.isLoading = ko.observable(false);
app.loading = function(status){
    if (status == true || status == false){
        app.isLoading(status);
        app.isContentShow(true);
    }    
}

app.isAuth = ko.observable(false);
app.auth = function(status){
    if (status == true || status == false){
        app.isAuth(status);
    }    
}

app.getUTCDate = function(strdate){
    var d = moment.utc(strdate);
    return d;
}

app.getDateMax = function(strdate){
    var d = moment.utc(strdate);
    return new Date(d.year(), d.month(), d.date(), 0, 0, 0)
}
app.toUTC = function(d){
    var year = d.getFullYear();
    var month = d.getMonth();
    var date = d.getDate();
    var hours = d.getHours();
    var minutes = d.getMinutes();
    var seconds = d.getSeconds();
    return moment(Date.UTC(year, month, date, hours, minutes, seconds)).toISOString();
}

app.getUrlParameter = function(sParam){
    var sPageURL = decodeURIComponent(window.location.search.substring(1)),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : sParameterName[1];
        }
    }
}


app.currentDateData = ko.observable;

app.extend('String', s);

// var colorField = ["#ff880e","#21c4af","#ff7663","#ffb74f","#a2df53","#1c9ec4","#ff63a5","#f44336","#D91E18","#8877A9","#9A12B3","#26C281","#E7505A","#C49F47","#ff5597","#c3260c","#d4735e","#ff2ad7","#34ac8b","#11b2eb","#f35838","#ff0037","#507ca3","#ff6565","#ffd664","#72aaff","#795548"];
var colorField = ["#ff9933","#21c4af","#4589b0", "#e4cc37","#008571", "#3ca5bf",
    "#003399", "#cccccc", "#9C27B0", "#BA68C8", "#1A237E",
    "#5C6BC0", "#1E88E5", "#0277BD", "#0097A7", "#26A69A",
    "#81C784", "#8BC34A", "#24752A", "#827717", "#004D40",
    "#C0CA33", "#FF6F00", "#D6C847", "#FFB300", "#BA8914",
    "#9999FF"];
var colorDegField = ["#ffcf9e","#a6e7df","#ffc8c0","#ffe2b8","#d9f2ba","#a4d8e7","#ffc0db","#fab3ae","#efa5a2","#cfc8dc","#d6a0e0","#a8e6cc","#f5b9bd","#e7d8b5","#ffbbd5","#e7a89d","#edc7be","#ffa9ef","#adddd0","#9fe0f7","#fabcaf","#ff99af","#b9cada","#ffc1c1","#ffeec1","#c6ddff","#c9bbb5"];
var colorFields2 =  ["#9e9e9e","#337ab7","#ff0000"];
var colorFieldProject = ["#398a6b","#f3a41b","#2f81b7","#c8be00","#30b4c9","#d90057","#f0c55d","#800080","#e67c52","#9098ad","#9e7d54","#a468a6"];
var colorAnalysis =["#ee6d54", "#ba21f3", "#5665cc", "#9f40f6", "#eff421", "#9e74f7", "#ed5ed9", "#842b29", "#dcd5aa", "#f0cda7", "#c869ca", "#d77f1d", "#475a8d", "#81b0c5", "#27aa94", "#29ac5f", "#791952", "#affc60", "#20fdaa", "#2296df", "#4574bc", "#c79549", "#5f41e1", "#ce3332", "#2cb3d0", "#9fe0cc", "#526042", "#ac7454", "#d2bec1", "#9f4f6f", "#5b82e5", "#8f46ae", "#21c377", "#b03c98", "#b0fbe4", "#9d636f", "#a9f834", "#c1bb71", "#41f0d5", "#448a88", "#e4373c", "#89c5f6", "#338b5a", "#b24bc1", "#d2e244", "#80cdd0", "#e76caf", "#7d91a1", "#9e5969", "#b1846b", "#85888b", "#977a42", "#db3c84", "#e2efa6", "#d977bd", "#94b4a8", "#92883b", "#6d7c80", "#79ddc7", "#8fb6b1"]

var customColorProject = {
    Amba : "#21c4af",
    Lahori : "#2f81b7",
    Tejuva : "#f0c231",
}
