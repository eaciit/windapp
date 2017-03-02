'use strict';

vm.currentMenu('Email Management');
vm.currentTitle("Email Management");
vm.breadcrumb([{ title: 'Operational', href: '#' },{ title: 'Email Management', href: viewModel.appName + 'page/emailmanagement' }]);

viewModel.Email = new Object();
var em = viewModel.Email;

em.templateEmail = {
    _id: "",
    Subject: "",
    Category: "", // refer to ref_emailCategory
    Receivers: [], // list of user ACL
    AlarmCodes: [], // list of alarm code from AlarmBrake > alarmname
    IntervalTime: 0, // in minutes
    Template: "",
    Enable: true
};
em.templateFilter = {
    search: ""
};
em.contentIsLoading = ko.observable(false);
em.TableColumns = ko.observableArray([{ headerTemplate: "<center><input type='checkbox' class='deletecheckall' onclick=\"em.checkDeleteData(this, 'deleteall', 'all')\"/></center>", attributes: { style: "text-align: center;" }, width: 40, template: function template(d) {
        return ["<input type='checkbox' class='deletecheck' idcheck='" + d._id + "' onclick=\"em.checkDeleteData(this, 'delete')\" />"].join(" ");
    } }, {
    field: "_id",
    title: "ID",
    headerAttributes: { style: "text-align: center;" }
}, {
    field: "subject",
    title: "Subject",
    headerAttributes: { style: "text-align: center;" }
}, {
    field: "category",
    title: "Category",
    headerAttributes: { style: "text-align: center;" },
    attributes: { style: "text-align: center;" }
}, {
    field: "enable",
    title: "Enable",
    headerAttributes: { style: "text-align: center;" },
    attributes: { style: "text-align: center;" }
}, {
    headerTemplate: "<center>Action</center>", width: 100,
    template: function template(d) {
        return ["<button class='btn btn-sm btn-warning' onclick='em.editData(\"" + d._id + "\")'><span class='fa fa-pencil' ></span></button>"].join(" ");
    }
}]);

em.filter = ko.mapping.fromJS(em.templateFilter);
em.config = ko.mapping.fromJS(em.templateEmail);
em.isNew = ko.observable(false);
em.tempCheckIdDelete = ko.observableArray([]);
em.selectedTableID = ko.observable("");

em.checkDeleteData = function (elem, e) {
    if (e === 'delete') {
        if ($(elem).prop('checked') === true) em.tempCheckIdDelete.push($(elem).attr('idcheck'));else em.tempCheckIdDelete.remove(function (item) {
            return item === $(elem).attr('idcheck');
        });
    }
    if (e === 'deleteall') {
        if ($(elem).prop('checked') === true) {
            $('.deletecheck').each(function (index) {
                $(this).prop("checked", true);
                em.tempCheckIdDelete.push($(this).attr('idcheck'));
            });
        } else {
            (function () {
                var idtemp = '';
                $('.deletecheck').each(function (index) {
                    $(this).prop("checked", false);
                    idtemp = $(this).attr('idcheck');
                    em.tempCheckIdDelete.remove(function (item) {
                        return item === idtemp;
                    });
                });
            })();
        }
    }
};

em.newData = function () {
    em.isNew(true);
    $('#modalUpdate').modal('show');
    ko.mapping.fromJS(em.templateEmail, em.config);
};

em.editData = function (id) {
    em.isNew(false);
    toolkit.ajaxPost(viewModel.appName + 'email/editemail', { _id: id }, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        ko.mapping.fromJS(res.data, em.config);
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

em.saveChanges = function () {
    if (!toolkit.isFormValid(".form-group")) {
        return;
    }
    var parm = ko.mapping.toJS(em.config);
    toolkit.ajaxPost(viewModel.appName + 'email/saveemail', parm, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        $('#modalUpdate').modal('hide');
        em.refreshData();
    }, function (err) {
        toolkit.showError(err.responseText);
    }, {
        timeout: 5000
    });
};

em.refreshData = function () {
    em.contentIsLoading(true);
    em.generateGrid();
    $('.grid-email').data('kendoGrid').dataSource.read();
    em.tempCheckIdDelete([]);
    ko.mapping.fromJS(em.templateEmail, em.config);
};

em.deletegroup = function () {
    if (em.tempCheckIdDelete().length === 0) {
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
            text: 'Data group(s) ' + em.tempCheckIdDelete().toString() + ' will be deleted',
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Delete",
            closeOnConfirm: true
        }, function () {
            setTimeout(function () {
                toolkit.ajaxPost(viewModel.appName + "email/deleteemail", { _id: em.tempCheckIdDelete() }, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                    em.refreshData();
                    swal({ title: "Email(s) successfully deleted", type: "success" });
                });
            }, 1000);
        });
    }
};

em.generateGrid = function () {
    $(".grid-email").html("");
    $('.grid-email').kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: viewModel.appName + "email/search",
                    type: "POST",
                    data: ko.mapping.toJS(em.filter),
                    dataType: "json",
                    contentType: "application/json; charset=utf-8",
                    success: function success(data) {
                        $(".grid-email>.k-grid-content-locked").css("height", $(".grid-email").data("kendoGrid").table.height());
                    }
                },
                parameterMap: function(options) {
                    return JSON.stringify(options);
                }
            },
            schema: {
                data: function data(res) {
                    em.selectedTableID("show");
                    em.contentIsLoading(false);
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
        columns: em.TableColumns()
    });
};

$(function () {
    $("#modalUpdate").insertAfter("body");
    em.generateGrid();
    // adm.getGrant()
});