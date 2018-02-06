"use strict";

var _typeof2 = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) { return typeof obj; } : function (obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol ? "symbol" : typeof obj; };

window.hasOwnProperty("ko") || console.error("knockoutjs is not installed"), window.hasOwnProperty("jQuery") || console.error("jQuery is not installed"), window.hasOwnProperty("moment") || console.error("momentjs is not installed"), window.hasOwnProperty("kendo") || console.error("kendoui is not installed");var toolkit = {};toolkit.dev = ko.observable(!0), toolkit.noop = function () {}, toolkit.noob = {};
"use strict";toolkit.isLast = function (t, n) {
  return t.indexOf(n) + 1 == t.length;
}, toolkit.merge = function (t, n) {
  return $.extend(!0, t, n);
}, toolkit.clone = function (t) {
  return toolkit.merge({}, t);
}, toolkit["return"] = function (t) {
  return t;
}, toolkit.distinct = function (t) {
  return t.filter(function (t, n, o) {
    return o.indexOf(t) === n;
  });
}, toolkit.hasProperty = function (t, n) {
  return t.hasOwnProperty(n);
}, toolkit.allKeys = function (t) {
  var n = [];for (var o in t) {
    t.hasOwnProperty(o) && n.push(String(o));
  }return n;
}, toolkit.length = function (t) {
  return t instanceof Object ? toolkit.allKeys(t).length : t.length;
}, toolkit.forEach = function (t, n) {
  if (t instanceof Object) for (var o in t) {
    t.hasOwnProperty(o) && n(o, t[o]);
  } else for (var r = 0; t.length; r++) {
    n(t[r], r);
  }
}, toolkit.split = function (t) {
  var n = arguments.length <= 1 || void 0 === arguments[1] ? "" : arguments[1],
      o = arguments.length <= 2 || void 0 === arguments[2] ? 0 : arguments[2];if (0 > o && (o = 0), 0 == o) return t.split(n);var r = [],
      i = [];return t.split(n).forEach(function (t, n) {
    return o > n ? void r.push(t) : void i.push(t);
  }), r.concat(i.join(n));
};
"use strict";toolkit.log = function () {
  toolkit.dev() && console.log.apply(console, [].slice.call(arguments));
}, toolkit.error = function () {
  toolkit.dev() && console.error.apply(console, [].slice.call(arguments));
};
"use strict";toolkit.resetForm = function (t) {
  var o = toolkit.$(t);o.trigger("reset");
}, toolkit.newEl = function (t) {
  return $("<" + t + " />");
};
"use strict";toolkit.ajaxPost = function (o) {
  var t = arguments.length <= 1 || void 0 === arguments[1] ? {} : arguments[1],
      n = arguments.length <= 2 || void 0 === arguments[2] ? toolkit.noop : arguments[2],
      a = arguments.length <= 3 || void 0 === arguments[3] ? toolkit.noop : arguments[3],
      e = arguments.length <= 4 || void 0 === arguments[4] ? toolkit.noob : arguments[4],
      i = (moment(), ko.mapping.toJSON(toolkit.noob));try {
    i = ko.mapping.toJSON(t);
  } catch (c) {}var s = { url: o.toLowerCase(), type: "post", dataType: "json", contentType: "application/json charset=utf-8", data: i, success: function success(o) {
      n(o);
    }, error: function error(o, t, n) {
      a(o, t, n);
    } };return t instanceof FormData && (delete s.config, s.data = t, s.async = !1, s.cache = !1, s.contentType = !1, s.processData = !1), s = $.extend(!0, s, e), $.ajax(s);
};
"use strict";toolkit.ajaxPostDeffered = function (o) {
  var D = $.Deferred();
  var t = arguments.length <= 1 || void 0 === arguments[1] ? {} : arguments[1],
      n = arguments.length <= 2 || void 0 === arguments[2] ? toolkit.noop : arguments[2],
      a = arguments.length <= 3 || void 0 === arguments[3] ? toolkit.noop : arguments[3],
      e = arguments.length <= 4 || void 0 === arguments[4] ? toolkit.noob : arguments[4],
      i = (moment(), ko.mapping.toJSON(toolkit.noob));try {
    i = ko.mapping.toJSON(t);
  } catch (c) {}var s = { url: o.toLowerCase(), type: "post", dataType: "json", contentType: "application/json charset=utf-8", data: i, success: function success(o) {
      n(o);
      D.resolve();
    }, error: function error(o, t, n) {
      a(o, t, n);
      D.resolve();
    } };return t instanceof FormData && (delete s.config, s.data = t, s.async = !1, s.cache = !1, s.contentType = !1, s.processData = !1), s = $.extend(!0, s, e), $.ajax(s);
    return D.promise();
};
"use strict";toolkit.$ = function (t) {
  return toolkit.typeIs(t, "string") ? $(t) : t;
};
"use strict";toolkit.resetValidation = function (t) {
  var a = toolkit.$(t),
      o = $(a).data("kendoValidator");o || (o = $(a).kendoValidator().data("kendoValidator"));try {
    o.hideMessages();
  } catch (i) {}
}, toolkit.isFormValid = function (t) {
  var a = toolkit.$(t);return toolkit.resetValidation(a), $(t).data("kendoValidator").validate();
};
"use strict";toolkit.koMap = ko.mapping.fromJS, toolkit.koUnmap = ko.mapping.toJS, toolkit.observ = ko.observable, toolkit.observArr = ko.observArr;
"use strict";var _typeof = "function" == typeof Symbol && "symbol" == _typeof2(Symbol.iterator) ? function (t) {
  return typeof t === "undefined" ? "undefined" : _typeof2(t);
} : function (t) {
  return t && "function" == typeof Symbol && t.constructor === Symbol ? "symbol" : typeof t === "undefined" ? "undefined" : _typeof2(t);
};toolkit.number = function (t) {
  var i = arguments.length <= 1 || void 0 === arguments[1] ? 0 : arguments[1];return isNaN(t) || !isFinite(t) ? i : t;
}, toolkit.noNaN = function (t) {
  var i = arguments.length <= 1 || void 0 === arguments[1] ? 0 : arguments[1];return isNaN(t) ? i : t;
}, toolkit.noInfinity = function (t) {
  var i = arguments.length <= 1 || void 0 === arguments[1] ? 0 : arguments[1];return isFinite(t) ? t : i;
}, toolkit.safeDivide = function (t, i) {
  return toolkit.number(t) / toolkit.number(i);
}, toolkit.redefine = function (t, i) {
  return toolkit.isUndefined(t) ? i : t;
}, toolkit.typeIs = function (t, i) {
  return ("undefined" == typeof t ? "undefined" : _typeof(t)) === i;
}, toolkit.isUndefined = function (t) {
  return toolkit.typeIs(t, "undefined");
}, toolkit.isDefined = function (t) {
  return !app.isUndefined(t);
}, toolkit.trim = function (t) {
  return $.trim(t);
}, toolkit.isEmptyString = function (t) {
  return "" == toolkit.trim(toolkit.redefine(t, ""));
}, toolkit.whenEmptyString = function (t) {
  var i = arguments.length <= 1 || void 0 === arguments[1] ? "" : arguments[1];return toolkit.isEmptyString(t) ? i : t;
}, toolkit.capitalize = function (t) {
  return 0 == String(t).length ? "" : "" + t[0].toUpperCase() + string.slice(1);
}, toolkit.isVoid = function (t) {
  return app.isUndefined(t) ? !0 : null == t ? !0 : !(!toolkit.typeIs(t, "string") || "" != toolkit.trim(t));
}, toolkit.capitalize = function (t) {
  return toolkit.isEmptyString(t) ? "" : t.toLowerCase().split(" ").map(function (t) {
    return t.length > 0 ? t[0].toUpperCase() + t.slice(1) : 0;
  }).join(" ");
}, toolkit.replace = function (t, i, n) {
  return toolkit.typeIs(i, "string") && (i = [i]), i.forEach(function (i) {
    t = t.replace(new RegExp("\\" + i, "g"), n);
  }), t;
};