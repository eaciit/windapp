'use strict';

vm.currentMenu('Administration');
vm.currentTitle("Admin Table");
vm.breadcrumb([{ title: 'Godrej', href: '#' }, { title: 'Administration', href: '#' }, { title: 'Admin Table', href: '/admintable' }]);

viewModel.admintable = new Object();
var at = viewModel.admintable;

at.templateFilter = {
    search: ""
};

at.TableColumns = ko.observableArray([{ title: "Table Name", template: function template(d) {
        var step = 6;
        var items = d.table.split("_");
        for (var i = 1; i < items.length / step; i++) {
            items.splice(i * step, 0, "<br />");
        }
        return items.join("_").replace(new RegExp("<br />_", "g"), "<br />");
    } }, { title: "Dimension", template: function template(d) {
        return d.dimensions.map(function (e) {
            return '<span class="tag bg-blue" style="margin-bottom: 3px; display: inline-block;">' + e.replace(/_/g, '.') + '</span>';
        }).join(' ');
    } }, { title: "Action", width: 80, attributes: { class: "align-center" }, template: "<button onclick='at.deletecollection(\"#:table #\")' class='btn btn-sm btn-danger tooltipster' title='Delete Collection'><span class='fa fa-remove'></span></button>" }]);

at.contentIsLoading = ko.observable(false);
at.selectedTableID = ko.observable("");
at.filter = ko.mapping.fromJS(at.templateFilter);

at.gridData = ko.observableArray([]);
at.gridConfig = {
    data: at.gridData,
    dataSource: {
        pageSize: 10
    },
    resizable: true,
    pageable: {
        pageSizes: 10
    },
    columns: at.TableColumns(),
    gridBound: toolkit.gridBoundTooltipster('.grid-at')
};

at.refreshData = function () {
    at.contentIsLoading(true);
    toolkit.ajaxPost(viewModel.appName + "report/getplcollections", {}, function (res) {
        if (!toolkit.isFine(res)) {
            return;
        }
        at.gridData(res.Data);
        at.contentIsLoading(false);
    });
};

at.clearcollection = function () {
    var allTables = at.gridData();

    swal({
        title: "Are you sure?",
        text: 'All PL* table will be deleted.',
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Delete",
        closeOnConfirm: true
    }, function () {
        setTimeout(function () {
            toolkit.ajaxPost(viewModel.appName + "report/deleteplcollection", { _id: allTables }, function (res) {
                if (!toolkit.isFine(res)) {
                    return;
                }
                at.refreshData();
            });
        }, 1000);
    });
};

at.deletecollection = function (idtable) {
    swal({
        title: "Are you sure?",
        text: 'Table will be deleted.',
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Delete",
        closeOnConfirm: true
    }, function () {
        setTimeout(function () {
            toolkit.ajaxPost(viewModel.appName + "report/deleteplcollection", { _id: [idtable] }, function (res) {
                if (!toolkit.isFine(res)) {
                    return;
                }
                at.refreshData();
            });
        }, 1000);
    });
};

$(function () {
    at.refreshData();
});