'use strict';

vm.currentMenu('Administration');
vm.currentTitle("Log");
vm.breadcrumb([{ title: 'Godrej', href: '#' }, { title: 'Administration', href: '#' }, { title: 'Log', href: '/log' }]);

viewModel.log = new Object();
var log = viewModel.log;

log.templateFilter = {
    search: ""
};

log.TableColumns = ko.observableArray([
    // { field: "status", title: "" },
    // { field: "loginid", title: "Username" },
    // { field: "created", title: "Created", template:'# if (created == "0001-01-01T00:00:00Z") {#-#} else {# #:moment(created).utc().format("DD-MMM-YYYY HH:mm:ss")# #}#' },
    // { field: "expired", title: "Expired", template:'# if (expired == "0001-01-01T00:00:00Z") {#-#} else {# #:moment(expired).utc().format("DD-MMM-YYYY HH:mm:ss")# #}#' },
    // { field: "duration", title: "Active In", template:'#= kendo.toString(duration, "n2")# H'},
    // { title: "Action", width: 80, attributes: { class: "align-center" }, template:"#if(status=='ACTIVE'){# <button data-value='#:_id #' onclick='ses.setexpired(\"#: _id #\", \"#: loginid #\")' name='expired' type='button' class='btn btn-sm btn-default btn-text-danger btn-stop tooltipster' title='Set Expired'><span class='fa fa-times'></span></button> #}else{# #}#" }
]);

log.contentIsLoading = ko.observable(false);
log.selectedTableID = ko.observable("");
log.filter = ko.mapping.fromJS(log.templateFilter);

log.refreshData = function () {
    log.contentIsLoading(true);
    $('.grid-log').data('kendoGrid').dataSource.read();
};

log.generateGrid = function () {
    $(".grid-log").html("");
    $('.grid-log').kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: "/log/getlog",
                    dataType: "json",
                    data: ko.mapping.toJS(log.filter),
                    type: "POST",
                    success: function success(data) {
                        $(".grid-group>.k-grid-content-locked").css("height", $(".grid-group").data("kendoGrid").table.height());
                    }
                }
            },
            schema: {
                data: function data(res) {
                    log.selectedTableID("show");
                    log.contentIsLoading(false);
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
        columns: log.TableColumns()
    });
};

$(function () {
    log.generateGrid();
});