'use strict';

vm.currentMenu('Email Management');
vm.currentTitle("Email Management");
vm.breadcrumb([{ title: 'Operational', href: '#' },{ title: 'Email Management', href: viewModel.appName + 'page/emailmanagement' }]);

viewModel.Email = new Object();
var em = viewModel.Email;

em.templateEmail = {
    _id: "",
    subject: "",
    category: "", // refer to ref_emailCategory
    receivers: [], // list of user ACL
    alarmcodes: [], // list of alarm code from AlarmBrake > alarmname
    intervaltime: 0, // in minutes
    template: "",
    enable: true,
    createddate: '',
    lastupdate: '',
    createdby: '',
    updatedby: ''
};

em.CategoryMailList = ko.observableArray([]);
em.UserMailList = ko.observableArray([]);
em.AlarmCodesMailList = ko.observableArray([]);
em.TemplateMailList = ko.observable();
em.isAlarmCode = ko.observable(false);
em.isInterval = ko.observable(false);

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
}, 
{
    headerTemplate: "<center>Action</center>", width: 100,
    template: function template(d) {
        return ["<button class='btn btn-sm btn-warning' onclick='em.editData(\"" + d._id + "\")'><span class='fa fa-pencil' ></span></button>"].join(" ");
    },
    attributes: { style: "text-align: center;" }
}
]);

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

em.checkCategory = function() {
    em.showHide($('#categoryList').data('kendoDropDownList').value());
}

em.showHide = function(category) {
    var resObj = em.CategoryMailList().filter(function(obj) {
        return obj.value == category;
    });
    var condition = resObj[0].condition.split(",");
    em.isAlarmCode(false);
    em.isInterval(false);
    $.each(condition, function(idx, val){
        if(val.indexOf("isAlarmCode") >= 0) {
            em.isAlarmCode(true);
        } else if(val.indexOf("isInterval") >= 0) {
            em.isInterval(true);
        }
    });

    var catVal = $('#categoryList').data('kendoDropDownList').value();
    if(catVal == "alarm01") {
        $('#templateMail').html(em.TemplateMailList().alarmTemplate)
    } else {
        $('#templateMail').html(em.TemplateMailList().dataTemplate)
    }
}

em.resetDDL = function() {
    $('#categoryList').data('kendoDropDownList').select(0);
    $('#userList').data('kendoMultiSelect').value([]);
    $('#alarmcodesList').data('kendoMultiSelect').value([]);
}

em.setDDL = function(data) {
    $('#categoryList').data('kendoDropDownList').value(data.category);
    $('#userList').data('kendoMultiSelect').value(data.receivers);
    $('#alarmcodesList').data('kendoMultiSelect').value(data.alarmcodes);

    em.showHide(data.category);
}
em.newData = function () {
    em.isNew(true);
    ko.mapping.fromJS(em.templateEmail, em.config);
    $('#editor').data('kendoEditor').value("");
    em.resetDDL();
    em.checkCategory();

    setTimeout(function(){
        $('#modalUpdate').modal('show');
    }, 100);
};

em.editData = function (id) {
    em.isNew(false);
    toolkit.ajaxPost(viewModel.appName + 'email/editemail', { _id: id }, function (res) {
        if (!app.isFine(res)) {
            return;
        }
        ko.mapping.fromJS(res.data, em.config);
        em.setDDL(res.data);
        $('#editor').data('kendoEditor').value(res.data.template);

        setTimeout(function(){
            $('#modalUpdate').modal('show');
        }, 100);
    });
};

em.setEditor = function() {
    $("#editor").html("");
    $("#editor").kendoEditor({ 
        resizable: {
            content: true,
            toolbar: true,
        },
        messages: {
            // fontName: "Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif"
            fontNameInherit: "Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif",
            fontSize: 12
        }
    });
}

em.saveChanges = function () {
    if (!toolkit.isFormValid(".form-group")) {
        return;
    }
    var param = ko.mapping.toJS(em.config);
    param.id = param._id;
    param.intervaltime = parseInt(param.intervaltime);
    param.category = $('#categoryList').data('kendoDropDownList').value();
    param.receivers = $('#userList').data('kendoMultiSelect').value();
    param.alarmcodes = $('#alarmcodesList').data('kendoMultiSelect').value();
    param.template = $('#editor').data('kendoEditor').value();
    param.lastupdate = new Date();
    if(em.isNew()) {
        param.createddate = new Date();
    }
    toolkit.ajaxPost(viewModel.appName + 'email/saveemail', param, function (res) {
        if (!app.isFine(res)) {
            return;
        }

        var dataEmail = res.data;
        var resCreate = em.UserMailList().filter(function(obj) {
            return obj.value == dataEmail.createdby;
        });
        dataEmail.createdby = resCreate[0].text;
        var resUpdate = em.UserMailList().filter(function(obj) {
            return obj.value == dataEmail.updatedby;
        });
        dataEmail.updatedby = resUpdate[0].text;

        var ajaxToFile = $.ajax({
            url: "http://ostrowfm-realtime.eaciitapp.com/email/mailtofile",
            data: dataEmail,
            contentType: false,
            dataType: "json",
            type: 'GET',
            success: function (data) {         
            }
        });

        $('#modalUpdate').modal('hide');
        em.refreshData();        
        swal({ title: res.message, type: "success" });        
    }, function (err) {
        toolkit.showError(err.responseText);
    });
};

em.refreshData = function () {
    em.contentIsLoading(true);
    em.generateGrid();
    $('.grid-email').data('kendoGrid').dataSource.read();
    em.tempCheckIdDelete([]);
    ko.mapping.fromJS(em.templateEmail, em.config);
};

em.deleteemail = function () {
    if (em.tempCheckIdDelete().length === 0) {
        swal({
            title: "",
            text: 'You havent choose any email to delete',
            type: "warning",
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "OK",
            closeOnConfirm: true
        });
    } else {
        swal({
            title: "Are you sure?",
            text: 'Data email(s) ' + em.tempCheckIdDelete().toString() + ' will be deleted',
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
        columns: em.TableColumns(),
        /*dataBound: function(e){
            var that = this;
            $(that.tbody).on("click", "tr", function (e) {
                var rowData = that.dataItem(this);
                em.editData(rowData._id);
            });
        }*/
    });
};

$(function () {
    $("#modalUpdate").insertAfter("body");
    em.generateGrid();
    em.setEditor();
});