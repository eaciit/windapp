"use strict";

viewModel.Administration = new Object();
var adm = viewModel.Administration;

adm.Access = ko.observableArray([]);
adm.Grant = ko.observableArray([]);
adm.dataAccess = ko.observableArray([]);
adm.dataAccessUser = ko.observableArray([]);
adm.dropdownAccess = ko.observableArray([]);

adm.unselectedAccessGroup = function (index) {
    return ko.computed(function () {
        var result = [];

        adm.dataAccess().forEach(function (d) {
            var isFound = false;

            gr.config.Grants().forEach(function (f, i) {
                if (f.AccessID() == d && i != index) {
                    isFound = true;
                }
            });

            if (!isFound) {
                result.push(d);
            }
        });
        return result;
    }, adm);
};

adm.getAccess = function () {
    var data = [];
    toolkit.ajaxPost(viewModel.appName + "access/getaccessdropdown", {}, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data == null) {
            res.data = "";
        }
        for (var i in res.data) {
            adm.Access.push(res.data[i]._id);
            adm.dataAccess.push(res.data[i]._id);
            adm.dataAccessUser.push(res.data[i]._id);

            data.push({
                text: res.data[i]._id,
                value: res.data[i].title
            });
        };
        // adm.createGridPrivilege(res.data)
        adm.dropdownAccess(res.data);
    });
};