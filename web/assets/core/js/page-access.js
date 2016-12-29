'use strict';

vm.currentMenu('Menu Access');
vm.currentTitle("Menu Access");
vm.breadcrumb([{ title: 'Operational', href: '#' },{ title: 'Menu Access', href: viewModel.appName + 'page/access' }]);

viewModel.Access = new Object();
var ac = viewModel.Access;

ac.index = [
    { "value": 1, "text": 1 },
    { "value": 2, "text": 2 },
    { "value": 3, "text": 3 },
    { "value": 4, "text": 4 },
    { "value": 5, "text": 5 },
    { "value": 6, "text": 6 },
    { "value": 7, "text": 7 },
    { "value": 8, "text": 8 },
    { "value": 9, "text": 9 },
    { "value": 10, "text": 10 },
    { "value": 11, "text": 11 },
    { "value": 12, "text": 12 },
    { "value": 13, "text": 13 },
    { "value": 14, "text": 14 },
    { "value": 15, "text": 15 },
];

ac.parentIDList = ko.observableArray([]);
ac.templateAccess = {
    _id: "",
    Title: "",
    Icon: "",
    ParentId: "",
    Index:1,
    Url: "",
    Enable: true
};
ac.templateFilter = {
    search: ""
};
ac.contentIsLoading = ko.observable(false);
ac.AccessColumns = ko.observableArray([{ headerTemplate: "<center><input type='checkbox' class='deletecheckall' onclick=\"ac.checkDeleteData(this, 'deleteall', 'all')\"/></center>", attributes: { style: "text-align: center;" }, width: 40, template: function template(d) {
        return ["<input type='checkbox' class='deletecheck' idcheck='" + d._id + "' onclick=\"ac.checkDeleteData(this, 'delete')\" />"].join(" ");
    } }, {
    field: "_id",
    title: "ID"
}, {
    field: "title",
    title: "Title"
}, {
    field: "icon",
    title: "Icon"
}, {
    field: "parentid",
    title: "Parent ID"
}, {
    field: "url",
    title: "URL"
}, {
    field: "index",
    title: "Index"
},{
    field: "enable",
    title: "Enable"
}, {
    headerTemplate: "<center>Action</center>", width: 100,
    template: function template(d) {
        return ["<button class='btn btn-sm btn-warning' onclick='ac.editData(\"" + d._id + "\")'><span class='fa fa-pencil'></span></button>"].join(" ");
    }
}]);

ac.filter = ko.mapping.fromJS(ac.templateFilter);
ac.config = ko.mapping.fromJS(ac.templateAccess);
ac.selectedTableID = ko.observable("");
ac.tempCheckIdDelete = ko.observableArray([]);
ac.isNew = ko.observable(false);


ac.populateParent = function() {
    app.ajaxPost(viewModel.appName + "/access/getparentid", {}, function(res) {
        if (!app.isFine(res)) {
            return;
        }
        if (res.data.length == 0) {
            res.data = [];;
            // ac.parentIDList([{ value: "", text: "" }]);
        } else {
            var datavalue = [{ value: "", text: "" }];
            if (res.data.length > 0) {
                $.each(res.data, function(key, val) {
                    var data = {};
                    data.value = val;
                    data.text = val;
                    datavalue.push(data);
                });
            }
            ac.parentIDList(datavalue);
        }
    });
};

ac.checkDeleteData = function (elem, e) {
    if (e === 'delete') {
        if ($(elem).prop('checked') === true) ac.tempCheckIdDelete.push($(elem).attr('idcheck'));else ac.tempCheckIdDelete.remove(function (item) {
            return item === $(elem).attr('idcheck');
        });
    }
    if (e === 'deleteall') {
        if ($(elem).prop('checked') === true) {
            $('.deletecheck').each(function (index) {
                $(this).prop("checked", true);
                ac.tempCheckIdDelete.push($(this).attr('idcheck'));
            });
        } else {
            (function () {
                var idtemp = '';
                $('.deletecheck').each(function (index) {
                    $(this).prop("checked", false);
                    idtemp = $(this).attr('idcheck');
                    ac.tempCheckIdDelete.remove(function (item) {
                        return item === idtemp;
                    });
                });
            })();
        }
    }
};

ac.newData = function () {
    ac.isNew(true);
    $('#modalUpdate').modal('show');
    ko.mapping.fromJS(ac.templateAccess, ac.config);
    ac.populateParent();
    setTimeout(function() {
        $("#ddlIndex").data("kendoDropDownList").value(1);
        $("#ddlParent").data("kendoDropDownList").value('');
    }, 300);
};

ac.editData = function (id) {
    ac.isNew(false);
    toolkit.ajaxPost(viewModel.appName + 'access/editaccess', { _id: id }, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        $('#modalUpdate').modal('show');
        ko.mapping.fromJS(res.data, ac.config);
        ac.populateParent();
        setTimeout(function() {
            $("#ddlIndex").data("kendoDropDownList").value(res.data.Index);
            $("#ddlParent").data("kendoDropDownList").value(res.data.ParentId);
        }, 300);
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

ac.saveChanges = function () {
    if (!toolkit.isFormValid(".form-access")) {
        return;
    }
    ac.config.Index( $("#ddlIndex").data("kendoDropDownList").value());
    ac.config.ParentId( $("#ddlParent").data("kendoDropDownList").value());
    toolkit.ajaxPost(viewModel.appName + 'access/saveaccess', ko.mapping.toJS(ac.config), function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }

        $('#modalUpdate').modal('hide');
        ac.refreshDataBrowser();
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

ac.refreshDataBrowser = function () {
    ac.contentIsLoading(true);
    ac.generateGrid();
    ac.tempCheckIdDelete([]);
    ko.mapping.fromJS(ac.templateAccess, ac.config);
};

ac.deleteaccess = function () {
    if (ac.tempCheckIdDelete().length === 0) {
        swal({
            title: "",
            text: 'You havent choose any access to delete',
            type: "warning",
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "OK",
            closeOnConfirm: true
        });
    } else {
        swal({
            title: "Are you sure?",
            text: 'Data access(s) ' + ac.tempCheckIdDelete().toString() + ' will be deleted',
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Delete",
            closeOnConfirm: true
        }, function () {
            setTimeout(function () {
                toolkit.ajaxPost(viewModel.appName + "access/deleteaccess", { _id: ac.tempCheckIdDelete() }, function (res) {
                    if (!toolkit.isFine(res)) {
                        return;
                    }
                    ac.refreshDataBrowser();
                    swal({ title: "Data access(s) successfully deleted", type: "success" });
                });
            }, 1000);
        });
    }
};

ac.generateGrid = function () {
    $(".grid-access").html("");
    $('.grid-access').kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: viewModel.appName + "access/getaccess",
                    type: "POST",
                    data: ko.mapping.toJS(ac.filter),
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                    success: function success(data) {
                        $(".grid-access>.k-grid-content-locked").css("height", $(".grid-access").data("kendoGrid").table.height());
                    }
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            schema: {
                data: function data(res) {
                    toolkit.isFine(res);
                    ac.selectedTableID("show");
                    ac.contentIsLoading(false);
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
        columns: ac.AccessColumns()
    });
};

$(function () {
    ac.generateGrid();
    $("#modalUpdate").insertAfter("body");
});