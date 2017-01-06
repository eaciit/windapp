'use strict';

vm.currentMenu('Group');
vm.currentTitle("Group");
vm.breadcrumb([{ title: 'Operational', href: '#' },{ title: 'Group', href: viewModel.appName + 'page/group' }]);

viewModel.Group = new Object();
var gr = viewModel.Group;

gr.templateGroup = {
    _id: "",
    Title: "",
    Enable: true,
    Owner: "",
    GroupType: 0,
    Grants: [],
    Filter: "",
    LoginID: "",
    Fullname: "",
    Email: ""
};
gr.templateGrants = {
    AccessID: "",
    AccessValue: []
};
gr.templateFilter = {
    search: ""
};
gr.contentIsLoading = ko.observable(false);
gr.TableColumns = ko.observableArray([{ headerTemplate: "<center><input type='checkbox' class='deletecheckall' onclick=\"gr.checkDeleteData(this, 'deleteall', 'all')\"/></center>", attributes: { style: "text-align: center;" }, width: 40, template: function template(d) {
        return ["<input type='checkbox' class='deletecheck' idcheck='" + d._id + "' onclick=\"gr.checkDeleteData(this, 'delete')\" />"].join(" ");
    } }, {
    field: "_id",
    title: "ID",
    headerAttributes: { style: "text-align: center;" }
}, {
    field: "title",
    title: "Title",
    headerAttributes: { style: "text-align: center;" }
}, {
    field: "owner",
    title: "Owner",
    headerAttributes: { style: "text-align: center;" }
}, {
    field: "enable",
    title: "Enable",
    headerAttributes: { style: "text-align: center;" },
    attributes: { style: "text-align: center;" }
}, {
    headerTemplate: "<center>Action</center>", width: 100,
    template: function template(d) {
        return ["<button class='btn btn-sm btn-warning' onclick='gr.editData(\"" + d._id + "\")'><span class='fa fa-pencil' ></span></button>"].join(" ");
    }
}]);

gr.filter = ko.mapping.fromJS(gr.templateFilter);
gr.config = ko.mapping.fromJS(gr.templateGroup);
gr.isNew = ko.observable(false);
gr.tempCheckIdDelete = ko.observableArray([]);
gr.selectedTableID = ko.observable("");

gr.checkDeleteData = function (elem, e) {
    if (e === 'delete') {
        if ($(elem).prop('checked') === true) gr.tempCheckIdDelete.push($(elem).attr('idcheck'));else gr.tempCheckIdDelete.remove(function (item) {
            return item === $(elem).attr('idcheck');
        });
    }
    if (e === 'deleteall') {
        if ($(elem).prop('checked') === true) {
            $('.deletecheck').each(function (index) {
                $(this).prop("checked", true);
                gr.tempCheckIdDelete.push($(this).attr('idcheck'));
            });
        } else {
            (function () {
                var idtemp = '';
                $('.deletecheck').each(function (index) {
                    $(this).prop("checked", false);
                    idtemp = $(this).attr('idcheck');
                    gr.tempCheckIdDelete.remove(function (item) {
                        return item === idtemp;
                    });
                });
            })();
        }
    }
};

gr.newData = function () {
    gr.isNew(true);
    $('#modalUpdate').modal('show');
    ko.mapping.fromJS(gr.templateGroup, gr.config);
};

gr.addGrant = function () {
    var datagrant = $.extend(true, {}, ko.mapping.toJS(gr.config));
    datagrant.Grants.push(gr.templateGrants);
    ko.mapping.fromJS(datagrant, gr.config);
};

gr.removeGrant = function (data) {
    gr.config.Grants.remove(data);
};

gr.editData = function (id) {
    gr.isNew(false);
    toolkit.ajaxPost(viewModel.appName + 'group/editgroup', { _id: id }, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        for (var i in res.data.Grants) {
            res.data.Grants[i].AccessValue = [];
        }
        ko.mapping.fromJS(res.data, gr.config);
        gr.displayAccess(res.data._id);
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

gr.displayAccess = function (e) {
    toolkit.ajaxPost(viewModel.appName + "group/getaccessgroup", {
        _id: e
    }, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data == null) {
            res.data = "";
        }
        for (var i = 0; i < res.data.length; i++) {
            gr.config.Grants()[i].AccessID(res.data[i].AccessID);
            if (res.data[i].AccessValue.indexOf(1) != -1) gr.config.Grants()[i].AccessValue.push("AccessCreate");
            if (res.data[i].AccessValue.indexOf(2) != -1) gr.config.Grants()[i].AccessValue.push("AccessRead");
            if (res.data[i].AccessValue.indexOf(4) != -1) gr.config.Grants()[i].AccessValue.push("AccessUpdate");
            if (res.data[i].AccessValue.indexOf(8) != -1) gr.config.Grants()[i].AccessValue.push("AccessDelete");
            if (res.data[i].AccessValue.indexOf(16) != -1) gr.config.Grants()[i].AccessValue.push("AccessSpecial1");
            if (res.data[i].AccessValue.indexOf(32) != -1) gr.config.Grants()[i].AccessValue.push("AccessSpecial2");
            if (res.data[i].AccessValue.indexOf(64) != -1) gr.config.Grants()[i].AccessValue.push("AccessSpecial3");
            if (res.data[i].AccessValue.indexOf(128) != -1) gr.config.Grants()[i].AccessValue.push("AccessSpecial4");
        }
        $('#modalUpdate').modal('show');
    });
};

gr.getAccess = function () {};

gr.saveChanges = function () {
    if (!toolkit.isFormValid(".form-group")) {
        return;
    }
    var parm = ko.mapping.toJS(gr.config);
    // parm.GroupType = parseInt(parm.GroupType)
    var postparm = {
        grants: parm.Grants,
        group: {
            _id: parm._id,
            Title: parm.Title,
            Owner: parm.Owner,
            Enable: parm.Enable,
            GroupType: parm.GroupType
        }
    };
    toolkit.ajaxPost(viewModel.appName + 'group/savegroup', postparm, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        $('#modalUpdate').modal('hide');
        gr.refreshData();
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

gr.refreshData = function () {
    gr.contentIsLoading(true);
    gr.generateGrid();
    $('.grid-group').data('kendoGrid').dataSource.read();
    gr.tempCheckIdDelete([]);
    ko.mapping.fromJS(gr.templateGroup, gr.config);
};

gr.deletegroup = function () {
    if (gr.tempCheckIdDelete().length === 0) {
        swal({
            title: "",
            text: 'You havent choose any group to delete',
            type: "warning",
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "OK",
            closeOnConfirm: true
        });
    } else {
        swal({
            title: "Are you sure?",
            text: 'Data group(s) ' + gr.tempCheckIdDelete().toString() + ' will be deleted',
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Delete",
            closeOnConfirm: true
        }, function () {
            setTimeout(function () {
                toolkit.ajaxPost(viewModel.appName + "group/deletegroup", { _id: gr.tempCheckIdDelete() }, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    gr.refreshData();
                    swal({ title: "Data group(s) successfully deleted", type: "success" });
                });
            }, 1000);
        });
    }
};

gr.generateGrid = function () {
    $(".grid-group").html("");
    $('.grid-group').kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: viewModel.appName + "group/search",
                    type: "POST",
                    data: ko.mapping.toJS(gr.filter),
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                    success: function success(data) {
                        $(".grid-group>.k-grid-content-locked").css("height", $(".grid-group").data("kendoGrid").table.height());
                    }
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            schema: {
                data: function data(res) {
                    gr.selectedTableID("show");
                    gr.contentIsLoading(false);
                    app.isFine(res);
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
        columns: gr.TableColumns()
    });
};

$(function () {
    $("#modalUpdate").insertAfter("body");
    gr.generateGrid();
    adm.getAccess();
    // adm.getGrant()
});