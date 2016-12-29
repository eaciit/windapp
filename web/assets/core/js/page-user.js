'use strict';

vm.currentMenu('User Management');
vm.currentTitle("User Management");
vm.breadcrumb([{ title: 'Operational', href: '#' },{ title: 'User Management', href: viewModel.appName + 'page/user' }]);

viewModel.User = new Object();
var us = viewModel.User;

us.templateFilter = {
    search: ""
};
us.templateUser = {
    _id: "",
    LoginID: "",
    FullName: "",
    Email: "",
    Password: "",
    confPass: "",
    LoginType: 0,
    LoginConf: "",
    Enable: true,
    Groups: [],
    Grants: []
};
us.templateGrants = {
    Grant: "",
    AccessID: "",
    AccessValue: []
};
us.templateLdap = {
    Address: "",
    BaseDN: "",
    Filter: "",
    Username: "",
    Password: "",
    Attribute: []
};
us.templateSelectGrant = {
    _id: "",
    DataAccess: []
};

us.contentIsLoading = ko.observable(false);
us.TableColumns = ko.observableArray([{ headerTemplate: "<center><input type='checkbox' class='deletecheckall' onclick=\"us.checkDeleteData(this, 'deleteall', 'all')\"/></center>", attributes: { style: "text-align: center;" }, width: 40, template: function template(d) {
        return ["<input type='checkbox' class='deletecheck' idcheck='" + d._id + "' onclick=\"us.checkDeleteData(this, 'delete')\" />"].join(" ");
    } }, {
    field: "loginid",
    title: "Login Id"
}, {
    field: "fullname",
    title: "Fullame"
}, {
    field: "email",
    title: "Email"
}, {
    field: "enable",
    title: "Enable"
}, {
    headerTemplate: "<center>Action</center>", width: 100,
    template: function template(d) {
        return ["<button class='btn btn-sm btn-warning' onclick='us.editData(\"" + d._id + "\")'><span class='fa fa-pencil' ></span></button>"].join(" ");
    }
}]);

us.listGroupsColumns = ko.observableArray([{
    field: "_id",
    title: "List Group"
}]);
us.selectedGroupsColumns = ko.observableArray([{
    field: "_id",
    title: "Selected Group"
}]);

us.filter = ko.mapping.fromJS(us.templateFilter);
us.config = ko.mapping.fromJS(us.templateUser);
us.isNew = ko.observable(false);
us.tempCheckIdDelete = ko.observableArray([]);
us.selectedTableID = ko.observable("");
us.Logintype = ko.observableArray(["Ldap", "Google", "Facebook"]);
us.DataGroups = ko.observableArray([]);
us.listGroupsData = ko.observableArray([]);
us.selectedGroupsData = ko.observableArray([]);
us.selectGrant = ko.observableArray([]);

us.getGrant = function () {
    var data = [];
    toolkit.ajaxPost(viewModel.appName + "group/getgroup", {}, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        if (res.data == null) {
            res.data = "";
        }
        us.listGroupsData([]);
        for (var i in res.data) {
            us.DataGroups.push(res.data[i]._id);
            us.listGroupsData.push({ _id: res.data[i]._id });
        };
    });
};

us.selectlistGridGroups = function (e) {
    var selectedDataItem = e != null ? e.sender.dataItem(e.sender.select()) : null;
    us.listGroupsData.remove(function (item) {
        return item._id == selectedDataItem._id;
    });
    us.selectedGroupsData.push({ _id: selectedDataItem._id });
    us.config.Groups.push(selectedDataItem._id);
    us.findGrant(selectedDataItem._id);
};

us.findGrant = function (e) {
    toolkit.ajaxPost(viewModel.appName + "group/getaccessgroup", {
        _id: e
    }, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        if (res.data == null) {
            res.data = "";
        }
        us.selectGrant.push({ _id: e, access: res.data });
        us.selectDataAccess(e, res.data);
        // for (var i = 0; i < res.data.length; i++) {
        //     us.config.Grants()[i].AccessID(res.data[i].AccessID)
        //     if (res.data[i].AccessValue.indexOf(1) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessCreate")
        //     if (res.data[i].AccessValue.indexOf(2) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessRead")
        //     if (res.data[i].AccessValue.indexOf(4) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessUpdate")
        //     if (res.data[i].AccessValue.indexOf(8) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessDelete")
        //     if (res.data[i].AccessValue.indexOf(16) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessSpecial1")
        //     if (res.data[i].AccessValue.indexOf(32) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessSpecial2")
        //     if (res.data[i].AccessValue.indexOf(64) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessSpecial3")
        //     if (res.data[i].AccessValue.indexOf(128) != -1)
        //         gr.config.Grants()[i].AccessValue.push("AccessSpecial4")
        // }
    });
};

us.selectDataAccess = function (id, data) {
    var datagrant = $.extend(true, {}, ko.mapping.toJS(us.config));
    var grant = {};
    for (var i = 0; i < data.length; i++) {
        grant = $.extend(true, {}, us.templateGrants);
        grant.AccessID = data[i].AccessID;
        grant.Grant = id;
        if (data[i].AccessValue.indexOf(1) != -1) grant.AccessValue.push("AccessCreate");
        if (data[i].AccessValue.indexOf(2) != -1) grant.AccessValue.push("AccessRead");
        if (data[i].AccessValue.indexOf(4) != -1) grant.AccessValue.push("AccessUpdate");
        if (data[i].AccessValue.indexOf(8) != -1) grant.AccessValue.push("AccessDelete");
        if (data[i].AccessValue.indexOf(16) != -1) grant.AccessValue.push("AccessSpecial1");
        if (data[i].AccessValue.indexOf(32) != -1) grant.AccessValue.push("AccessSpecial2");
        if (data[i].AccessValue.indexOf(64) != -1) grant.AccessValue.push("AccessSpecial3");
        if (data[i].AccessValue.indexOf(128) != -1) grant.AccessValue.push("AccessSpecial4");

        datagrant.Grants.push(grant);
        ko.mapping.fromJS(datagrant, us.config);
    }
};

us.removeselectGridGroups = function (e) {
    var selectedDataItem = e != null ? e.sender.dataItem(e.sender.select()) : null;
    us.selectedGroupsData.remove(function (item) {
        return item._id == selectedDataItem._id;
    });
    us.listGroupsData.push({ _id: selectedDataItem._id });
    us.selectGrant.remove(function (item) {
        return item._id == selectedDataItem._id;
    });
    us.config.Grants.remove(function (item) {
        return item.Grant != '';
    });
    us.config.Groups.remove(function (item) {
        return item != selectedDataItem._id;
    });
    for (var i in us.selectGrant()) {
        us.selectDataAccess(us.selectGrant()[i]._id, us.selectGrant()[i].access);
    }
};

us.checkDeleteData = function (elem, e) {
    if (e === 'delete') {
        if ($(elem).prop('checked') === true) us.tempCheckIdDelete.push($(elem).attr('idcheck'));else us.tempCheckIdDelete.remove(function (item) {
            return item === $(elem).attr('idcheck');
        });
    }
    if (e === 'deleteall') {
        if ($(elem).prop('checked') === true) {
            $('.deletecheck').each(function (index) {
                $(this).prop("checked", true);
                us.tempCheckIdDelete.push($(this).attr('idcheck'));
            });
        } else {
            (function () {
                var idtemp = '';
                $('.deletecheck').each(function (index) {
                    $(this).prop("checked", false);
                    idtemp = $(this).attr('idcheck');
                    us.tempCheckIdDelete.remove(function (item) {
                        return item === idtemp;
                    });
                });
            })();
        }
    }
};

us.newData = function () {
    us.isNew(true);
    $('#modalUpdate').modal('show');
    us.tempCheckIdDelete([]);
    ko.mapping.fromJS(us.templateUser, us.config);
};

us.addGrant = function () {
    var datagrant = $.extend(true, {}, ko.mapping.toJS(us.config));
    datagrant.Grants.push(us.templateGrants);
    ko.mapping.fromJS(datagrant, us.config);
};

us.removeGrant = function (data) {
    us.config.Grants.remove(data);
};

us.editData = function (id) {
    us.isNew(false);
    toolkit.ajaxPost(viewModel.appName + 'user/edituser', { _id: id }, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        // for (var i in res.data.Grants){
        //     res.data.Grants[i].AccessValue = []
        //     res.data.Grants[i]["Grant"] = ""
        // }
        res.data.Grants = [];
        ko.mapping.fromJS(res.data, us.config);
        us.getGrant();
        us.selectedGroupsData([]);
        setTimeout(function() {
            for (var i in res.data.Groups) {
                us.listGroupsData.remove(function (item) {
                    return item._id == res.data.Groups[i];
                });
                us.selectedGroupsData.push({ _id: res.data.Groups[i] });
                us.findGrant(res.data.Groups[i]);
            }
        }, 100);
        $('#modalUpdate').modal('show');
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

us.saveChanges = function () {
    if (!toolkit.isFormValid(".form-user")) {
        return;
    }
    var parm = ko.mapping.toJS(us.config);
    var parmGrant = [];
    for (var i in parm.Grants) {
        parmGrant.push({
            AccessID: parm.Grants[i].AccessID,
            AccessValue: parm.Grants[i].AccessValue
        });
    }
    var postparm = {
        grants: parmGrant,
        user: {
            _id: parm._id,
            LoginID: parm.LoginID,
            FullName: parm.FullName,
            Email: parm.Email,
            Password: parm.Password,
            Enable: true,
            LoginType: parm.LoginType,
            Groups: parm.Groups
        }
    };
    toolkit.ajaxPost(viewModel.appName + 'user/saveuser', postparm, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }

        $('#modalUpdate').modal('hide');
        us.refreshData();
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

us.refreshData = function () {
    us.contentIsLoading(true);
    us.generateGrid();
    $('.grid-user').data('kendoGrid').dataSource.read();
    ko.mapping.fromJS(us.templateUser, us.config);
};

us.deleteuser = function () {
    if (us.tempCheckIdDelete().length === 0) {
        swal({
            title: "",
            text: 'You havent choose any user to delete',
            type: "warning",
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "OK",
            closeOnConfirm: true
        });
    } else {
        swal({
            title: "Are you sure?",
            text: 'Data user(s) ' + us.tempCheckIdDelete().toString() + ' will be deleted',
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Delete",
            closeOnConfirm: true
        }, function () {
            setTimeout(function () {
                toolkit.ajaxPost(viewModel.appName + "user/deleteuser", { _id: us.tempCheckIdDelete() }, function (res) {
                    if (!toolkit.isFine(res)) {
                        return;
                    }
                    us.refreshData();
                    swal({ title: "Data user(s) successfully deleted", type: "success" });
                });
            }, 1000);
        });
    }
};

us.generateGrid = function () {
    $(".grid-user").html("");
    $('.grid-user').kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: viewModel.appName + "user/getuser",
                    type: "POST",
                    data: ko.mapping.toJS(us.filter),
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                    success: function success(data) {
                        $(".grid-user>.k-grid-content-locked").css("height", $(".grid-user").data("kendoGrid").table.height());
                    }
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            schema: {
                data: function data(res) {
                    us.selectedTableID("show");
                    us.contentIsLoading(false);
                    toolkit.isFine(res);
                    return res.data.Data;
                },
                total: "data.total"
            },

            pageSize: 10,
            serverPaging: true, // enable server paging
            serverSorting: true
        },
        // selectable: "multiple, row",
        // change: ac.selectGridAccess,
        resizable: true,
        scrollable: true,
        // sortable: true,
        // filterable: true,
        pageable: {
            refresh: false,
            pageSizes: 10,
            buttonCount: 5
        },
        columns: us.TableColumns()
    });
};

$(function () {
    us.generateGrid();
    adm.getAccess();
    us.getGrant();
    $("#modalUpdate").insertAfter("body");
});