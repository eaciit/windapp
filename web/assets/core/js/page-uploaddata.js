'use strict';

vm.currentMenu('Data Manager');
vm.currentTitle('Upload Data');
vm.breadcrumb([{ title: 'Godrej', href: '#' }, { title: 'Data Manager', href: '#' }, { title: 'Upload Data', href: '/uploaddata' }]);

viewModel.uploadData = new Object();
var ud = viewModel.uploadData;

ud.inputDescription = ko.observable('');
ud.inputModel = ko.observable('');

ud.dataUploadedFiles = ko.observableArray([]);
ud.masterDataBrowser = ko.observableArray([]);
ud.dropDownModel = {
	data: ko.computed(function () {
		return ud.masterDataBrowser().map(function (d) {
			return { text: d.TableNames, value: d.TableNames };
		});
	}),
	dataValueField: 'value',
	dataTextField: 'text',
	value: ud.inputModel,
	optionLabel: 'Select one'
};
ud.gridUploadedFiles = {
	data: ud.dataUploadedFiles,
	dataSource: {
		pageSize: 10
	},
	columns: [{ title: '&nbsp;', width: 40, attributes: { class: 'align-center' }, template: function template(d) {
			return '<input type="checkbox" />';
		} }, { title: 'File Name', field: 'Filename', attributes: { class: 'bold' }, template: function template(d) {
			return '\n\t\t\t\t<div class=\'tooltipster\' title=\'File: ' + d.Filename + '<br />Description: ' + d.Desc + '\'>\n\t\t\t\t\t' + d.Filename + '\n\t\t\t\t</div>\n\t\t\t';
		} }, { headerTemplate: '<center>Model</center>', field: 'DocName', width: 90, template: function template(d) {
			return '<center>\n\t\t\t\t<span class="tag bg-green">' + d.DocName + '</span>\n\t\t\t</center>';
		} }, { headerTemplate: '<center>Date</center>', width: 120, template: function template(d) {
			return moment(d.date).format('DD-MM-YYYY HH:mm:ss');
		}
	}, { headerTemplate: '<center>Action</center>', width: 100, template: function template(d) {
			switch (d.Status) {
				case 'ready':
					return '\n\t\t\t\t\t<button class="btn btn-xs btn-warning tooltipster" title="Ready" onclick="ud.processData(`' + d._id + '`,this)">\n\t\t\t\t\t\t<i class="fa fa-play"></i> Run process\n\t\t\t\t\t</button>\n\t\t\t\t';
				case 'rollback':
					return '\n\t\t\t\t\t<button class="btn btn-xs btn-warning tooltipster" title="Ready"">\n\t\t\t\t\t\t<i class="fa fa-refresh"></i> Rollback\n\t\t\t\t\t</button>\n\t\t\t\t';
				case 'done':
					return '<span class=\'tag bg-green\'>Done</span>';
				case 'failed':
					return '<span class=\'tag bg-green\'>Failed</span>';
				case 'onprocess':
					return '<span class=\'tag bg-green\'>On Process</span>';
			}

			return '';
		} }],
	filterable: false,
	sortable: false,
	resizable: false,
	dataBound: toolkit.gridBoundTooltipster('.grid-uploadData')
};
ud.processData = function (data) {
	toolkit.ajaxPost(viewModel.appName + 'uploaddata/processdata', { _id: data }, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}

		ud.getUploadedFiles();
	}, function (err) {
		toolkit.showError(err.responseText);
	}, {
		timeout: 5000
	});
};

ud.getMasterDataBrowser = function () {
	ud.masterDataBrowser([]);

	toolkit.ajaxPost(viewModel.appName + 'databrowser/getdatabrowsers', {}, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}

		ud.masterDataBrowser(res.data);
	}, function (err) {
		toolkit.showError(err.responseText);
	}, {
		timeout: 5000
	});
};
ud.getUploadedFiles = function () {
	ud.dataUploadedFiles([]);

	toolkit.ajaxPost(viewModel.appName + 'uploaddata/getuploadedfiles', {}, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}

		ud.dataUploadedFiles(res.data);
	}, {
		timeout: 5000
	});
};
ud.doUpload = function () {
	if (!toolkit.isFormValid('.form-upload-file')) {
		return;
	}

	var payload = new FormData();
	payload.append('model', ud.inputModel());
	payload.append('desc', ud.inputDescription());
	payload.append('userfile', $('[name=file]')[0].files[0]);

	toolkit.ajaxPost(viewModel.appName + 'uploaddata/uploadfile', payload, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}

		ud.getUploadedFiles();
	}, function (err) {
		toolkit.showError(err.responseText);
	}, {
		timeout: 5000
	});
};

ud.init = function () {
	ud.getMasterDataBrowser();
	ud.getUploadedFiles();
};

$(function () {
	ud.init();
});