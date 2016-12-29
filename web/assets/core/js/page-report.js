'use strict';

var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) { return typeof obj; } : function (obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol ? "symbol" : typeof obj; };

// let menuLink = vm.menu()
// 	.find((d) => d.href == ('/' + document.URL.split('/').slice(3).join('/')))

vm.currentMenu('Report');
vm.currentTitle('Report');
vm.breadcrumb([{ title: 'Godrej', href: '#' }, { title: 'Report', href: '#' }]);

viewModel.report = new Object();
var rpt = viewModel.report;

rpt.filter = [{ _id: 'common', group: 'Base Filter', sub: [{ _id: 'Branch', from: 'Branch', title: 'Branch' }, { _id: 'BranchGroup', from: 'MasterBranchGroup', title: 'Branch Group' }, { _id: 'Brand', from: 'Brand', title: 'Brand' }, { _id: 'Channel', from: 'Channel', title: 'Channel' }, { _id: 'Region', from: 'Region', title: 'Region' }] }, { _id: 'geo', group: 'Geographical', sub: [{ _id: 'Zone', from: 'Zone', title: 'Zone' },
	// { _id: 'Region', from: 'Region', title: 'Region' },
	{ _id: 'Area', from: 'MasterArea', title: 'City' }] }, { _id: 'customer', group: 'Customer', sub: [
	// { _id: 'ChannelC', from: 'Channel', title: 'Channel' },
	{ _id: 'KeyAccount', from: 'KeyAccount', title: 'Key Account' }, { _id: 'CustomerGroup', from: 'CustomerGroup', title: 'Group' },
	// { _id: 'Customer', from: 'Customer', title: 'Outlet' },
	{ _id: 'Distributor', from: 'MasterDistributor', title: 'Distributor' }] }, { _id: 'product', group: 'Product', sub: [{ _id: 'HBrandCategory', from: 'HBrandCategory', title: 'Group' }] }];

rpt.pivotModel = [{ field: '_id', type: 'string', name: 'ID' }, { field: 'PC._id', type: 'string', name: 'Profit Center - ID' }, { field: 'PC.EntityID', type: 'string', name: 'Profit Center - Entity ID' }, { field: 'PC.Name', type: 'string', name: 'Profit Center - Name' }, { field: 'PC.BrandID', type: 'string', name: 'Profit Center - Brand ID' }, { field: 'PC.BrandCategoryID', type: 'string', name: 'Profit Center - Brand Category ID' }, { field: 'PC.BranchID', type: 'string', name: 'Profit Center - Branch ID' }, { field: 'PC.BranchType', type: 'int', name: 'Profit Center - Branch Type' }, { field: 'CC._id', type: 'string', name: 'Cost Center - ID' }, { field: 'CC.EntityID', type: 'string', name: 'Cost Center - Entity ID' }, { field: 'CC.Name', type: 'string', name: 'Cost Center - Name' }, { field: 'CC.CostGroup01', type: 'string', name: 'Cost Center - Cost Group 01' }, { field: 'CC.CostGroup02', type: 'string', name: 'Cost Center - Cost Group 02' }, { field: 'CC.CostGroup03', type: 'string', name: 'Cost Center - Cost Group 03' }, { field: 'CC.BranchID', type: 'string', name: 'Cost Center - Branch ID' }, { field: 'CC.BranchType', type: 'string', name: 'Cost Center - Branch Type' }, { field: 'CC.CCTypeID', type: 'string', name: 'Cost Center - Type' }, { field: 'CC.HCCGroupID', type: 'string', name: 'Cost Center - HCC Group ID' }, { field: 'CompanyCode', type: 'string', name: 'Company Code' }, { field: 'LedgerAccount', type: 'string', name: 'Ledger Account' }, { field: 'Customer._id', type: 'string', name: 'Customer - ID' }, { field: 'Customer.BranchID', type: 'string', name: 'Customer - Branch ID' }, { field: 'Customer.BranchName', type: 'string', name: 'Customer - branch Name' }, { field: 'Customer.Name', type: 'string', name: 'Customer - Name' }, { field: 'Customer.KeyAccount', type: 'string', name: 'Customer - Key Account' }, { field: 'Customer.ChannelID', type: 'string', name: 'Customer - Channel ID' }, { field: 'Customer.ChannelName', type: 'string', name: 'Customer - Channel Name' }, { field: 'Customer.CustomerGroup', type: 'string', name: 'Customer - Customer Group' }, { field: 'Customer.CustomerGroupName', type: 'string', name: 'Customer - Customer Group Name' }, { field: 'Customer.National', type: 'string', name: 'Customer - National' }, { field: 'Customer.Zone', type: 'string', name: 'Customer - Zone' }, { field: 'Customer.Region', type: 'string', name: 'Customer - Region' }, { field: 'Customer.Area', type: 'string', name: 'Customer - Area' }, { field: 'Product._id', type: 'string', name: 'Product - ID' }, { field: 'Product.Name', type: 'string', name: 'Product - Name' }, { field: 'Product.ProdCategory', type: 'string', name: 'Product - Category' }, { field: 'Product.Brand', type: 'string', name: 'Product - Brand' }, { field: 'Product.BrandCategoryID', type: 'string', name: 'Product - Brand Category ID' }, { field: 'Product.PCID', type: 'string', name: 'Product - PCID' }, { field: 'Product.ProdSubCategory', type: 'string', name: 'Product - Sub Category' }, { field: 'Product.ProdSubBrand', type: 'string', name: 'Product - Sub Brand' }, { field: 'Product.ProdVariant', type: 'string', name: 'Product - Variant' }, { field: 'Product.ProdDesignType', type: 'string', name: 'Product - Design Type' }, { field: 'Date.ID', type: 'string', name: 'Date - ID' }, { field: 'Date.Date', type: 'string', name: 'Date - Date' }, { field: 'Date.Month', type: 'string', name: 'Date - Month' }, { field: 'Date.Quarter', type: 'int', name: 'Date - Quarter' }, { field: 'Date.YearTxt', type: 'string', name: 'Date - YearTxt' }, { field: 'Date.QuarterTxt', type: 'string', name: 'Date - QuarterTxt' }, { field: 'Date.Year', type: 'int', name: 'Date - Year' }, { field: 'PLGroup1', type: 'string', name: 'PL Group 1' }, { field: 'PLGroup2', type: 'string', name: 'PL Group 2' }, { field: 'PLGroup3', type: 'string', name: 'PL Group 3' }, { field: 'PLGroup4', type: 'string', name: 'PL Group 4' }, { field: 'Value1', type: 'double', name: 'Value 1', as: 'dataPoints' }, { field: 'Value2', type: 'double', name: 'Value 2', as: 'dataPoints' }, { field: 'Value3', type: 'double', name: 'Value 3', as: 'dataPoints' }, { field: 'PCID', type: 'string', name: 'Profit Center ID' }, { field: 'CCID', type: 'string', name: 'Cost Center ID' }, { field: 'SKUID', type: 'string', name: 'SKU ID' }, { field: 'PLCode', type: 'string', name: 'PL Code' }, { field: 'Month', type: 'string', name: 'Month' }, { field: 'Year', type: 'string', name: 'Year' }];

rpt.getFilterValue = function () {
	var multiFiscalYear = arguments.length <= 0 || arguments[0] === undefined ? false : arguments[0];
	var fiscalField = arguments.length <= 1 || arguments[1] === undefined ? rpt.value.FiscalYear : arguments[1];

	var res = [{ 'Field': 'customer.branchname', 'Op': '$in', 'Value': rpt.value.Branch() }, { 'Field': 'customer.branchgroup', 'Op': '$in', 'Value': rpt.value.BranchGroup() }, { 'Field': 'product.brand', 'Op': '$in', 'Value': rpt.value.Brand().concat([]) }, { 'Field': 'customer.channelname', 'Op': '$in', 'Value': rpt.value.Channel().concat([]) }, { 'Field': 'customer.region', 'Op': '$in', 'Value': rpt.value.Region().concat([]) },
	// { 'Field': 'date.year', 'Op': '$gte', 'Value': rpt.value.From() },
	// { 'Field': 'date.year', 'Op': '$lte', 'Value': rpt.value.To() },

	{ 'Field': 'customer.zone', 'Op': '$in', 'Value': rpt.value.Zone() },
	// ---> Region OK
	{ 'Field': 'customer.areaname', 'Op': '$in', 'Value': rpt.value.Area() },

	// ---> Channel OK
	{ 'Field': 'customer.keyaccount', 'Op': '$in', 'Value': rpt.value.KeyAccount() }, { 'Field': 'customer.customergroup', 'Op': '$in', 'Value': rpt.value.CustomerGroup() },
	// { 'Field': 'customer.name', 'Op': '$in', 'Value': rpt.value.Customer() },
	{ 'Field': 'customer.reportsubchannel', 'Op': '$in', 'Value': rpt.value.Distributor() }, { 'Field': 'product.brandcategoryid', 'Op': '$in', Value: rpt.value.HBrandCategory() }];

	if (fiscalField !== false) {
		if (multiFiscalYear) {
			res.push({
				'Field': 'date.fiscal',
				'Op': '$in',
				'Value': fiscalField()
			});
		} else {
			rpt.saveFiscalYear(fiscalField());
			res.push({
				'Field': 'date.fiscal',
				'Op': '$eq',
				'Value': fiscalField()
			});
		}
	}

	res = res.filter(function (d) {
		if (d.Value instanceof Array) {
			return d.Value.length > 0;
		} else {
			return d.Value != '';
		}
	});

	return res;
};

rpt.getFiscalYear = function () {
	var fy = rpt.optionFiscalYears();

	var savedFY = toolkit.redefine(localStorage.fiscalYear, fy[1]);
	if (fy.indexOf(savedFY) == -1) {
		savedFY = fy[1];
	}

	return savedFY;
};

rpt.saveFiscalYear = function (fy) {
	localStorage.fiscalYear = fy;
};

rpt.pnlTableHeaderWidth = ko.observable('560px');
rpt.optionFiscalYears = ko.observableArray(['2014-2015', '2015-2016']);
rpt.analysisIdeas = ko.observableArray([]);
rpt.data = ko.observableArray([]);
rpt.percentOfTotal = ko.observable(true);
rpt.optionDimensions = ko.observableArray([
// { field: "", name: 'None', title: '' },
{ field: "customer.branchname", name: 'Branch/RD', title: 'customer_branchname' }, { field: "customer.branchgroup", name: 'Branch Group', title: 'customer_branchgroup' }, { field: "product.brand", name: 'Brand', title: 'product_brand' }, { field: 'customer.channelname', name: 'Channel', title: 'customer_channelname' },
// { field: 'customer.name', name: 'Outlet', title: 'customer_name' },
// { field: 'product.name', name: 'Product', title: 'product_name' },
// { field: 'customer.zone', name: 'Zone', title: 'customer_zone' },
{ field: "customer.areaname", name: "City", title: "customer_areaname" }, { field: 'customer.region', name: 'Region', title: 'customer_region' }, { field: "customer.zone", name: "Zone", title: "customer_zone" },
// { field: 'date.fiscal', name: 'Fiscal Year', title: 'date_fiscal' },
{ field: 'customer.keyaccount', name: 'Customer Group', title: 'customer_keyaccount' }]);
// rpt.optionDataPoints = ko.observableArray([
//     { field: 'value1', name: o['value1'] },
//     { field: 'value2', name: o['value2'] },
//     { field: 'value3', name: o['value3'] }
// ])
rpt.optionAggregates = ko.observableArray([{ aggr: 'sum', name: 'Sum' }, { aggr: 'avg', name: 'Avg' }, { aggr: 'max', name: 'Max' }, { aggr: 'min', name: 'Min' }]);
rpt.optionsChannels = ko.observableArray([{ _id: 'EXP', Name: 'Export' }, { _id: 'I2', Name: 'General Trade (GT)' }, { _id: 'I4', Name: 'Industrial Trade (IT)' }, { _id: 'I3', Name: 'Modern Trade (MT)' }, { _id: 'I6', Name: 'Motorist' }, { _id: 'I1', Name: 'Regional Distributor (RD)' }]);

rpt.date_quartertxt = ko.observableArray([{ FiscalYear: '2014-2015', _id: '2014-2015 Q1', Name: '2014-2015 Q1' }, { FiscalYear: '2014-2015', _id: '2014-2015 Q2', Name: '2014-2015 Q2' }, { FiscalYear: '2014-2015', _id: '2014-2015 Q3', Name: '2014-2015 Q3' }, { FiscalYear: '2014-2015', _id: '2014-2015 Q4', Name: '2014-2015 Q4' }, { FiscalYear: '2015-2016', _id: '2015-2016 Q1', Name: '2015-2016 Q1' }, { FiscalYear: '2015-2016', _id: '2015-2016 Q2', Name: '2015-2016 Q2' }, { FiscalYear: '2015-2016', _id: '2015-2016 Q3', Name: '2015-2016 Q3' }, { FiscalYear: '2015-2016', _id: '2015-2016 Q4', Name: '2015-2016 Q4' }]);
rpt.date_month = ko.observableArray(function () {
	var months = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];
	var res = [];
	months.forEach(function (e) {
		var month = e - 1 + 3;
		var text = moment(new Date(2015, month, 1)).format('MMMM');

		res.push({
			_id: e,
			Name: text
		});
	});

	return res;
}());
rpt.monthQuarter = ko.observable('');
rpt.optionMonthQuarters = ko.observableArray([{ _id: 'date_quartertxt', Name: 'Quarter' }, { _id: 'date_month', Name: 'Month' }]);
rpt.monthQuarterValues = ko.observableArray([]);
rpt.optionMonthsQuarterValues = function () {
	var fiscalYearObservable = arguments.length <= 0 || arguments[0] === undefined ? null : arguments[0];

	return ko.computed(function () {
		if (rpt.monthQuarter() == '') {
			return [];
		}

		if (fiscalYearObservable == null) {
			return rpt[rpt.monthQuarter()]();
		}

		var fiscalYears = fiscalYearObservable();
		if (!(fiscalYears instanceof Array)) {
			fiscalYears = [fiscalYears];
		}

		return rpt[rpt.monthQuarter()]().filter(function (d) {
			if (d.hasOwnProperty('FiscalYear')) {
				return fiscalYears.indexOf(d.FiscalYear) > -1;
			}

			return true;
		});
	}, rpt.monthQuarter);
};
rpt.resetMonthQuarter = function () {
	rpt.monthQuarter('');
	rpt.monthQuarterValues([]);
};
rpt.changeMonthQuarter = function () {
	rpt.monthQuarterValues([]);
};
rpt.injectMonthQuarterFilter = function (filters) {
	rpt.isUsingMonthQuarterFiter(true);

	if (rpt.monthQuarter() == '') {
		return filters;
	}

	if (rpt.monthQuarterValues().length == 0) {
		return filters;
	}

	var field = rpt.monthQuarter().replace(/_/g, '.');
	var prev = filters.find(function (d) {
		return d.Field == field;
	});

	if (toolkit.isDefined(prev)) {
		var prevValues = prev.Value;
		if (typeof prevValues === 'string') {
			prevValues = [prevValues];
		}

		prev.Op = '$in';
		prev.Value = prevValues.concat(rpt.monthQuarterValues());
	} else {
		filters.push({
			Field: field,
			Op: '$in',
			Value: rpt.monthQuarterValues()
		});
	}

	return filters;
};
rpt.isEnableMonthQuarterValues = ko.computed(function () {
	return rpt.monthQuarter() !== '';
});
rpt.isUsingMonthQuarterFiter = ko.observable(false);

rpt.parseGroups = function (what) {
	return what;

	if (what.indexOf('customer.branchname') > -1) {
		what.push('customer.branchid');
	}
	// if (what.indexOf('customer.channelname') > -1) {
	// 	what.push('customer.channelid')
	// }
	// if (what.indexOf('customer.customergroupname') > -1) {
	// 	what.push('customer.customergroup')
	// }

	return what;
};
rpt.rowHeaderHeight = ko.observable(34);
rpt.rowContentHeight = ko.observable(26);
rpt.mode = ko.observable('render');
rpt.refreshView = ko.observable('');
rpt.modecustom = ko.observable(false);
rpt.idanalysisreport = ko.observable();
rpt.valueMasterData = {};
rpt.masterData = {
	geographi: ko.observableArray([]),
	subchannel: ko.observableArray([])
};
rpt.showPercentOfTotal = ko.observable(true);
rpt.enableHolder = {};
rpt.eventChange = {};
rpt.value = {
	HQ: ko.observable(false),
	From: ko.observable(new Date(2014, 0, 1)),
	To: ko.observable(new Date(2016, 11, 31)),
	FiscalYear: ko.observable(rpt.getFiscalYear()),
	FiscalYears: ko.observableArray([rpt.getFiscalYear()])
};
rpt.masterData.Type = ko.observableArray([{ value: 'Mfg', text: 'Mfg' }, { value: 'Branch', text: 'Branch' }]);
rpt.masterData.HQ = ko.observableArray([{ value: true, text: 'True' }, { value: false, text: 'False' }]);
rpt.filter.forEach(function (d) {
	d.sub.forEach(function (e) {
		if (rpt.masterData.hasOwnProperty(e._id)) {
			return;
		}

		if (!rpt.value.hasOwnProperty(e._id)) {
			rpt.value[e._id] = ko.observableArray([]);
		}
		rpt.valueMasterData[e._id] = ko.observableArray([]);
		rpt.masterData[e._id] = ko.observableArray([]);
		rpt.enableHolder[e._id] = ko.observable(true);
		rpt.eventChange[e._id] = function () {
			var self = this;
			var value = self.value();

			setTimeout(function () {
				var vZone = rpt.valueMasterData['Zone']();
				var vRegion = rpt.valueMasterData['Region']();
				var vArea = rpt.valueMasterData['Area']();

				if (e._id == 'Zone') {
					var raw = Lazy(rpt.masterData.geographi()).filter(function (f) {
						return vZone.length == 0 ? true : vZone.indexOf(f.Zone) > -1;
					}).toArray();

					rpt.masterData.Region(rpt.groupGeoBy(raw, 'Region'));
					rpt.masterData.Area(rpt.groupGeoBy(raw, 'Area'));
				} else if (e._id == 'Region') {
					var _raw = Lazy(rpt.masterData.geographi()).filter(function (f) {
						return vZone.length == 0 ? true : vZone.indexOf(f.Zone) > -1;
					}).filter(function (f) {
						return vRegion.length == 0 ? true : vRegion.indexOf(f.Region) > -1;
					}).toArray();

					rpt.masterData.Area(rpt.groupGeoBy(_raw, 'Area'));
					rpt.enableHolder['Zone'](vRegion.length == 0);
				} else if (e._id == 'Area') {
					var _raw2 = Lazy(rpt.masterData.geographi()).filter(function (f) {
						return vZone.length == 0 ? true : vZone.indexOf(f.Zone) > -1;
					}).filter(function (f) {
						return vRegion.length == 0 ? true : vRegion.indexOf(f.Region) > -1;
					}).toArray();

					rpt.enableHolder['Region'](vArea.length == 0);
					rpt.enableHolder['Zone'](vRegion.length == 0);
				}

				// change value event goes here
				toolkit.log(e._id, value);
			}, 100);
		};
	});
});

rpt.getOtherMasterData = function () {
	toolkit.ajaxPost(viewModel.appName + 'report/getdatasubchannel', {}, function (res) {
		if (!toolkit.isFine(res)) {
            return;
        }

		rpt.masterData.subchannel(res.data);
	});
};

rpt.isEmptyData = function (res) {
	if (res == null || typeof res === 'undefined') {
		return true;
	}

	if (res.Data == null) {
		return true;
	}

	if (res.Data.length == 0) {
		toolkit.showError('The UI data is not ready');
		return true;
	}

	return false;
};

rpt.groupGeoBy = function (raw, category) {
	var groupKey = category == 'Area' ? '_id' : category;
	var data = Lazy(raw).groupBy(function (f) {
		return f[groupKey];
	}).map(function (k, v) {
		return { _id: v, Name: v };
	}).toArray();

	return data;
};

rpt.filterMultiSelect = function (d) {
	var config = {
		filter: 'contains',
		placeholder: 'Choose items ...',
		// change: rpt.eventChange[d._id],
		value: rpt.value[d._id]
	};

	if (['HQ', 'Type'].indexOf(d.from) > -1) {
		config = $.extend(true, config, {
			data: rpt.masterData[d._id],
			dataValueField: 'value',
			dataTextField: 'text',
			value: rpt.value[d._id]
		});
	} else if (['Customer'].indexOf(d.from) > -1) {
		config = $.extend(true, config, {
			autoBind: false,
			minLength: 1,
			placeholder: 'Type min 1 chars',
			dataValueField: 'Name',
			dataTextField: 'Name',
			template: function template(d) {
				return d._id + ' - ' + d.Name;
			},
			enabled: rpt.enableHolder[d._id],
			dataSource: {
				serverFiltering: true,
				transport: {
					read: {
						url: '/report/getdata' + d.from.toLowerCase()
					},
					parameterMap: function parameterMap(data, type) {
						var keyword = data.filter.filters[0].value;
						return { keyword: keyword };
					}
				},
				schema: {
					data: 'data'
				}
			},
			value: rpt.value[d._id]
		});
	} else if (['Branch', 'Brand', 'MasterArea', 'MasterDistributor', 'MasterBranchGroup', 'HCostCenterGroup', 'Entity', 'Channel', 'HBrandCategory', 'Product', 'Type', 'KeyAccount', 'CustomerGroup', 'LedgerAccount'].indexOf(d.from) > -1) {
		config = $.extend(true, config, {
			data: rpt.masterData[d._id],
			dataValueField: '_id',
			dataTextField: 'Name',
			enabled: rpt.enableHolder[d._id],
			value: rpt.value[d._id]
		});

		if (['Branch', 'Brand'].indexOf(d.from) > -1) {
			config.dataValueField = 'Name';
		} else if (d.from == 'Product') {
			config = $.extend(true, config, {
				dataValueField: 'Name',
				minLength: 1,
				placeholder: 'Type min 1 chars'
			});
		} else if (['Channel', 'KeyAccount'].indexOf(d.from) > -1) {
			config.dataValueField = '_id';
		}

		toolkit.ajaxPost(viewModel.appName + ('report/getdata' + d.from.toLowerCase()), {}, function (res) {
			if (!toolkit.isFine(res)) {
	            return;
	        }

			rpt.masterData[d._id](_.sortBy(res.data, function (d) {
				return d.Name;
			}));

			if (['KeyAccount', 'Brand', 'Branch', 'BranchGroup'].indexOf(d.from) > -1) {
				rpt.masterData[d._id].push({ _id: "OTHER", Name: "OTHER" });
			}
		});
	} else if (['Region', /* 'Area', */'Zone'].indexOf(d.from) > -1) {
		config = $.extend(true, config, {
			data: rpt.masterData[d._id],
			dataValueField: '_id',
			dataTextField: 'Name',
			enabled: rpt.enableHolder[d._id],
			value: rpt.value[d._id]
		});

		if (d.from == 'Region') {
			toolkit.ajaxPost(viewModel.appName + 'report/getdatahgeographi', {}, function (res) {
				if (!toolkit.isFine(res)) {
		            return;
		        }

				rpt.masterData.geographi(_.sortBy(res.data, function (d) {
					return d.Name;
				}));

				['Region', /* 'Area', */'Zone'].forEach(function (e) {
					var res = rpt.groupGeoBy(rpt.masterData.geographi(), e);
					rpt.masterData[e](_.sortBy(res, function (d) {
						return d.Name;
					}));
					rpt.masterData[e].push({ _id: "OTHER", Name: "OTHER" });
				});

				// rpt.masterData.RegionC(rpt.masterData.Region())
			});
		}
	} else {
		config.data = rpt.masterData[d._id]();
	}

	return config;
};

rpt.toggleFilterCallback = toolkit.noop;
rpt.toggleFilter = function () {
	var panelFilter = $('.panel-filter');
	var panelContent = $('.panel-content');

	if (panelFilter.is(':visible')) {
		panelFilter.hide();
		panelContent.attr('class', 'col-md-12 col-sm-12 ez panel-content');
		$('.breakdown-filter').removeAttr('style');
		$('.scroll-grid-bottom.arrow-right').css('right', 23);
	} else {
		panelFilter.show();
		panelContent.attr('class', 'col-md-9 col-sm-9 ez panel-content');
		$('.breakdown-filter').css('width', '60%');
		$('.scroll-grid-bottom.arrow-right').css('right', 320);
	}

	$('.k-grid').each(function (i, d) {
		try {
			$(d).data('kendoGrid').refresh();
		} catch (err) {}
	});

	$('.k-pivot').each(function (i, d) {
		$(d).data('kendoPivotGrid').refresh();
	});

	$('.k-chart').each(function (i, d) {
		$(d).data('kendoChart').redraw();
	});

	$('.table-content').trigger('scroll');
};

// rpt.getIdeas = () => {
// 	toolkit.ajaxPost(viewModel.appName + 'report/getdataanalysisidea', { }, (res) => {
// 		if (!toolkit.isFine(res)) {
// 			return
// 		}

// 		rpt.idanalysisreport('')
// 		rpt.analysisIdeas(_.sortBy(res.data, (d) => d.order))
// 		let idreport = _.find(rpt.analysisIdeas(), function(a) { return a._id == o.ID })
// 		if (idreport != undefined) {
// 			rpt.idanalysisreport(idreport.name)
// 			vm.currentTitle("Report " + rpt.idanalysisreport())
// 		}
// 	})
// }

rpt.wrapParam = function () {
	var dimensions = arguments.length <= 0 || arguments[0] === undefined ? [] : arguments[0];
	var dataPoints = arguments.length <= 1 || arguments[1] === undefined ? [] : arguments[1];

	return {
		dimensions: dimensions,
		dataPoints: dataPoints,
		filters: rpt.getFilterValue(),
		which: o.ID
	};
};

rpt.setName = function (data, options) {
	return function () {
		setTimeout(function () {
			var row = options().find(function (d) {
				return d.field == data.field();
			});
			if (toolkit.isDefined(row)) {
				data.name(row.name);
			}

			console.log(toolkit.koUnmap(data), options());
		}, 150);
	};
};
rpt.refresh = function () {
	['pvt', 'tbl', 'crt', 'sct', 'bkd'].forEach(function (d, i) {
		setTimeout(function () {
			if (toolkit.isDefined(window[d])) {
				window[d].refresh();
			}
		}, 1000 * i);
	});
};
rpt.refreshAll = function () {
	switch (rpt.refreshView()) {
		case 'analysis':
			bkd.changeBreakdown();
			bkd.refresh();
			rs.refresh();
			ccr.refresh();
			break;
		case 'dashboard':
			dsbrd.changeBreakdown();
			dsbrd.refresh();
			rank.refresh();
			sd.refresh();
			break;
		case 'reportwidget':
			pvt.refresh();
			crt.refresh();
			break;
	}
};

rpt.tabbedContent = function () {
	$('.app-title h2').html('&nbsp;');
	$('.tab-content').parent().find('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
		var $tabContent = $(".tab-content " + $(e.target).attr("href"));

		setTimeout(function () {
			$tabContent.find('.k-chart').each(function (i, e) {
				try {
					$(e).data('kendoChart').redraw();
				} catch (err) {
					console.log(err);
				}
			});
			$tabContent.find('.k-grid').each(function (i, e) {
				try {
					$(e).data('kendoGrid').refresh();
				} catch (err) {
					console.log(err);
				}
			});
		}, 200);
	});
};

rpt.plmodels = ko.observableArray([]);
rpt.allowedPL = ko.computed(function () {
	var plmodels = ["PL0", "PL6A", "PL7A", "PL8", "PL7", "PL8A", "PL6", "PL2", "PL1", "PL14A", "PL14", "PL9", "PL74A", "PL74", "PL21", "PL74B", "PL74C", "PL23", "PL26A", "PL25", "PL32A", "PL31", "PL31E", "PL31D", "PL31C", "PL31B", "PL31A", "PL29A", "PL29A32", "PL29A31", "PL29A30", "PL29A29", "PL29A27", "PL29A26", "PL29A25", "PL29A24", "PL29A23", "PL29A22", "PL29A20", "PL29A19", "PL29A18", "PL29A17", "PL29A16", "PL29A15", "PL29A14", "PL29A13", "PL29A12", "PL29A11", "PL29A10", "PL29A9", "PL29A8", "PL29A6", "PL29A5", "PL29A4", "PL29A3", "PL29A2", "PL28", "PL28I", "PL28G", "PL28F", "PL28E", "PL28D", "PL28C", "PL28B", "PL28A", "PL32B", "PL94A", "PL94B", "PL44B", "PL44", "PL43", "PL42", "PL44C", "PL44E", "PL44D", "PL44F", "PL33", "PL34", "PL35", "PL74D"];
	return rpt.plmodels().filter(function (d) {
		return plmodels.indexOf(d._id) > -1;
	});
}, rpt.plmodels);
rpt.idarrayhide = ko.observableArray(['PL44A']);

rpt.prepareEvents = function () {
	$('.breakdown-view').parent().off('mouseover').on('mouseover', 'tr', function () {
		var rowID = $(this).attr('data-row');

		var elh = $('.breakdown-view .table-header tr[data-row="' + rowID + '"]').addClass('hover');
		var elc = $('.breakdown-view .table-content tr[data-row="' + rowID + '"]').addClass('hover');
	});
	$('.breakdown-view').parent().off('mouseleave').on('mouseleave', 'tr', function () {
		$('.breakdown-view tr.hover').removeClass('hover');
	});
};

rpt.hardcodePLGA = function (data, plmodels) {
	if (data.length == 0) {
		return { Data: data, PLModels: plmodels };
	}

	if (document.URL.indexOf('cogsanalysis') > -1) {
		return { Data: data, PLModels: plmodels };
	}

	var sgaLvl1 = ['Direct', 'Allocated'];
	var sgaLvl2 = [{ _id: 'PL33', header: 'Personnel Exp - Office' }, { _id: 'PL34', header: 'General Exp - Office' }, { _id: 'PL35', header: 'Depr & A Exp - Office' }];
	var sgaLvl3 = ['RD', 'Sales', 'General Service', 'General Management', 'Manufacturing', 'Finance', 'Marketing', 'Logistic Overhead', 'Human Resource', 'OTHER'];

	var gnaExpenses = plmodels.find(function (d) {
		return d._id == 'PL94A';
	});
	var plprev = gnaExpenses._id;
	var gnaChildren = [];
	var wrap = function wrap(what, i) {
		if (i == 0) {
			return what;
		}

		return what + ' ';
	};

	sgaLvl1.forEach(function (d, i) {
		var plTop1 = {};
		plTop1._id = [gnaExpenses._id, d].join('_').replace(/\ /g, '_');
		plTop1.OrderIndex = [gnaExpenses.OrderIndex, i].join('');
		plTop1.PLHeader1 = ['G&A Expenses', d].join(' - ');
		plTop1.PLHeader2 = plTop1.PLHeader1;
		plTop1.PLHeader3 = plTop1.PLHeader1;
		plTop1.parent = '';
		plTop1.prev = plprev;

		gnaChildren.push(plTop1);
		plprev = plTop1._id;

		sgaLvl2.forEach(function (e, j) {
			var plTop2 = {};
			plTop2._id = [e._id, d].join('_').replace(/\ /g, '_');
			plTop2.OrderIndex = [gnaExpenses.OrderIndex, i, j].join(''); // [plTop1.OrderIndex, j].join('')
			plTop2.PLHeader1 = ['G&A Expenses', d].join(' - ');
			plTop2.PLHeader2 = wrap([d, e.header].join('_'), i);
			plTop2.PLHeader3 = plTop2.PLHeader2;
			plTop2.parent = plTop1._id;
			plTop2.prev = plprev;

			gnaChildren.push(plTop2);
			plprev = plTop2._id;

			sgaLvl3.forEach(function (f, k) {
				var plTop3 = {};
				plTop3._id = [e._id, d, f].join('_').replace(/\ /g, '_');
				plTop3.OrderIndex = [gnaExpenses.OrderIndex, i, j, k].join(''); // [plTop1.OrderIndex, k].join('')
				plTop3.PLHeader1 = ['G&A Expenses', d].join(' - ');
				plTop3.PLHeader2 = plTop2.PLHeader2;
				plTop3.PLHeader3 = wrap([d, e._id, f].join('_'), i);
				plTop3.parent = plTop2._id;
				plTop3.prev = plprev;

				gnaChildren.push(plTop3);
				plprev = plTop3._id;
			});
		});
	});

	console.log('gnaChildren', gnaChildren);
	plmodels = plmodels.concat(gnaChildren);

	rpt.arrChangeParent(rpt.arrChangeParentOriginal().slice(0));
	var direct = gnaChildren.filter(function (d) {
		return d.PLHeader1.indexOf('G&A Expenses - Direct') > -1;
	});
	direct.forEach(function (d, i) {
		var o = {};
		o.idfrom = d._id;
		o.idto = d.parent;
		o.after = d.prev;

		rpt.arrChangeParent.push(o);
	});

	var allocated = gnaChildren.filter(function (d) {
		return d.PLHeader1.indexOf('G&A Expenses - Allocated') > -1;
	});
	allocated.forEach(function (d, i) {
		var o = {};
		o.idfrom = d._id;
		o.idto = d.parent;
		o.after = d.prev;
		rpt.arrChangeParent.push(o);
	});

	return { Data: data, PLModels: plmodels };
};

rpt.showExpandAll = function (a) {
	if (a == true) {
		$('tr.dd').find('i').removeClass('fa-chevron-right');
		$('tr.dd[idheaderpl=\'PL0\']').find('i').addClass('fa fa-chevron-up');
		$('tr.dd[idheaderpl!=\'PL0\']').find('i').addClass('fa fa-chevron-down');
		$('tr[idparent]').css('display', '');
		$('tr[idcontparent]').css('display', '');
		$('tr[statusvaltemp=hide]').css('display', 'none');
	} else {
		$('tr.dd').find('i').removeClass('fa-chevron-up');
		$('tr.dd').find('i').removeClass('fa-chevron-down');
		$('tr.dd').find('i').addClass('fa fa-chevron-right');
		$('tr[idparent]').css('display', 'none');
		$('tr[idcontparent]').css('display', 'none');
		$('tr[statusvaltemp=hide]').css('display', 'none');
	}
};

rpt.zeroValue = ko.observable(false);
rpt.showZeroValue = function (a) {
	rpt.zeroValue(a);
	if (a == true) {
		$(".table-header tbody>tr").each(function (i) {
			if (i > 0) {
				$(this).attr('statusvaltemp', 'show');
				$('tr[idpl=' + $(this).attr('idheaderpl') + ']').attr('statusvaltemp', 'show');
				if (!$(this).attr('idparent')) {
					$(this).show();
					$('tr[idpl=' + $(this).attr('idheaderpl') + ']').show();
				}
			}
		});
	} else {
		$(".table-header tbody>tr").each(function (i) {
			if (i > 0) {
				$(this).attr('statusvaltemp', $(this).attr('statusval'));
				$('tr[idpl=' + $(this).attr('idheaderpl') + ']').attr('statusvaltemp', $(this).attr('statusval'));
			}
		});
	}

	rpt.showExpandAll(false);
	if (a == false) {
		(function () {
			var countchild = 0,
			    hidechild = 0;
			$(".table-header tbody>tr.dd").each(function (i) {
				if (i > 0) {
					countchild = $('.table-header tr[idparent=' + $(this).attr('idheaderpl') + ']').length;
					hidechild = $('.table-header tr[idparent=' + $(this).attr('idheaderpl') + '][statusvaltemp=hide]').length;
					if (countchild > 0) {
						if (countchild == hidechild) {
							$(this).find('td:eq(0)>i').removeClass().css('margin-right', '0px');
							if ($(this).attr('idparent') == undefined) $(this).find('td:eq(0)').css('padding-left', '20px');
						}
					}
				}
			});
		})();
	}
};

rpt.arrChangeParentOriginal = ko.observableArray([{ idfrom: 'PL1', idto: 'PL0', after: 'PL0' }, { idfrom: 'PL7', idto: 'PL0', after: 'PL1' }, { idfrom: 'PL2', idto: 'PL0', after: 'PL7' }, { idfrom: 'PL8', idto: 'PL0', after: 'PL2' }, { idfrom: 'PL6', idto: 'PL0', after: 'PL8' }, { idfrom: 'PL7A', idto: '', after: 'PL6' }]);
rpt.arrChangeParent = ko.observableArray(rpt.arrChangeParentOriginal().slice(0));

// rpt.arrFormulaPL = ko.observableArray([
// 	{ id: "PL0", formula: ["PL1","PL2","PL3","PL4","PL5","PL6"], cal: "sum"},
// 	{ id: "PL6A", formula: ["PL7","PL8","PL7A"], cal: "sum"},
// ])
// rpt.arrFormulaPL = ko.observableArray([
// 	{ id: "PL1", formula: ["PL7"], cal: "sum"},
// 	{ id: "PL2", formula: ["PL8"], cal: "sum"},
// ])

rpt.arrFormulaPL = ko.observableArray([{ id: "PL2", formula: ["PL2", "PL8"], cal: "sum" }, { id: "PL1", formula: ["PL8A", "PL2", "PL6"], cal: "min" }]);

rpt.changeParent = function (elemheader, elemcontent, PLCode) {
	var change = _.find(rpt.arrChangeParent(), function (a) {
		return a.idfrom == PLCode;
	});
	if (change != undefined) {
		if (change.idto != '') {
			elemheader.attr('idparent', change.idto);
			elemcontent.attr('idcontparent', change.idto);
		} else {
			elemheader.removeAttr('idparent');
			elemheader.find('td:eq(0)').css('padding-left', '8px');
			elemcontent.removeAttr('idcontparent');
		}
		return change.after;
	} else {
		return "";
	}
};

rpt.fixRowValue = function (data) {
	return;

	data.forEach(function (e, a) {
		rpt.arrFormulaPL().forEach(function (d) {
			// let total = toolkit.sum(d.formula, (f) => e[f])
			var total = 0;
			var isNotNumber = false;

			d.formula.forEach(function (f, l) {
				var eachValue = e[f];

				if (!toolkit.typeIs(eachValue, "number")) {
					eachValue = toolkit.getNumberFromString(eachValue);
					isNotNumber = true;
				}

				if (l == 0) {
					total = eachValue;
				} else {
					if (d.cal == 'sum') {
						total += eachValue;
					} else {
						total -= eachValue;
					}
				}
			});

			if (isNotNumber) {
				total = kendo.toString(total, 'n2') + ' %';
			}

			data[a][d.id] = total;
		});
	});
	// console.log(data)
};

rpt.buildGridLevels = function (rows) {
	var plmodels = rpt.plmodels();

	var grouppl1 = _.map(_.groupBy(plmodels, function (d) {
		return d.PLHeader1;
	}), function (k, v) {
		return { data: k, key: v };
	});
	var grouppl2 = _.map(_.groupBy(plmodels, function (d) {
		return d.PLHeader2;
	}), function (k, v) {
		return { data: k, key: v };
	});
	var grouppl3 = _.map(_.groupBy(plmodels, function (d) {
		return d.PLHeader3;
	}), function (k, v) {
		return { data: k, key: v };
	});

	var $trElem = void 0,
	    $columnElem = void 0;
	var resg1 = void 0,
	    resg2 = void 0,
	    resg3 = void 0,
	    PLyo = void 0,
	    PLyo2 = void 0,
	    child = 0,
	    parenttr = 0,
	    textPL = void 0;
	$(".table-header tbody>tr").each(function (i) {
		if (i > 0) {
			$trElem = $(this);
			resg1 = _.find(grouppl1, function (o) {
				return o.key == $trElem.find('td:eq(0)').text();
			});
			resg2 = _.find(grouppl2, function (o) {
				return o.key == $trElem.find('td:eq(0)').text();
			});
			resg3 = _.find(grouppl3, function (o) {
				return o.key == $trElem.find('td:eq(0)').text();
			});

			var idplyo = _.find(rpt.idarrayhide(), function (a) {
				return a == $trElem.attr("idheaderpl");
			});
			if (idplyo != undefined) {
				$trElem.remove();
				$('.table-content tr.column' + $trElem.attr("idheaderpl")).remove();
			}

			var idplyo2 = _.find(rpt.idarrayhide(), function (a) {
				return a == $trElem.attr("idparent");
			});
			if (resg1 == undefined && idplyo2 == undefined) {
				if (resg2 != undefined) {
					textPL = _.find(resg2.data, function (o) {
						return o._id == $trElem.attr("idheaderpl");
					});
					PLyo = _.find(rows, function (o) {
						return o.PNL == textPL.PLHeader1;
					});
					PLyo2 = _.find(rows, function (o) {
						return o.PLCode == textPL._id;
					});
					$trElem.find('td:eq(0)').css('padding-left', '40px');
					$trElem.attr('idparent', PLyo.PLCode);
					child = $('tr[idparent=' + PLyo.PLCode + ']').length;
					$columnElem = $('.table-content tr.column' + PLyo2.PLCode);
					$columnElem.attr('idcontparent', PLyo.PLCode);
					var PLCodeChange = rpt.changeParent($trElem, $columnElem, $columnElem.attr('idpl'));
					if (PLCodeChange != "") PLyo.PLCode = PLCodeChange;
					if (child > 1) {
						var $parenttr = $('tr[idheaderpl=' + PLyo.PLCode + ']');
						var $parenttrcontent = $('tr[idpl=' + PLyo.PLCode + ']');
						// $trElem.insertAfter($(`tr[idparent=${PLyo.PLCode}]:eq(${(child-1)})`))
						// $columnElem.insertAfter($(`tr[idcontparent=${PLyo.PLCode}]:eq(${(child-1)})`))
						$trElem.insertAfter($parenttr);
						$columnElem.insertAfter($parenttrcontent);
					} else {
						$trElem.insertAfter($('tr[idheaderpl=' + PLyo.PLCode + ']'));
						$columnElem.insertAfter($('tr[idpl=' + PLyo.PLCode + ']'));
					}
				} else if (resg2 == undefined) {
					if (resg3 != undefined) {
						PLyo = _.find(rows, function (o) {
							return o.PNL == resg3.data[0].PLHeader2;
						});
						PLyo2 = _.find(rows, function (o) {
							return o.PNL == resg3.data[0].PLHeader3;
						});
						$trElem.find('td:eq(0)').css('padding-left', '70px');
						if (PLyo == undefined) {
							PLyo = _.find(rows, function (o) {
								return o.PNL == resg3.data[0].PLHeader1;
							});
							if (PLyo != undefined) $trElem.find('td:eq(0)').css('padding-left', '40px');
						}
						$trElem.attr('idparent', PLyo.PLCode);
						child = $('tr[idparent=' + PLyo.PLCode + ']').length;
						$columnElem = $('.table-content tr.column' + PLyo2.PLCode);
						$columnElem.attr('idcontparent', PLyo.PLCode);
						var _PLCodeChange = rpt.changeParent($trElem, $columnElem, $columnElem.attr('idpl'));
						if (_PLCodeChange != "") PLyo.PLCode = _PLCodeChange;
						if (child > 1) {
							// $trElem.insertAfter($(`tr[idparent=${PLyo.PLCode}]:eq(${(child-1)})`))
							// $columnElem.insertAfter($(`tr[idcontparent=${PLyo.PLCode}]:eq(${(child-1)})`))
							$trElem.insertAfter($('tr[idheaderpl=' + PLyo.PLCode + ']'));
							$columnElem.insertAfter($('tr[idpl=' + PLyo.PLCode + ']'));
						} else {
							$trElem.insertAfter($('tr[idheaderpl=' + PLyo.PLCode + ']'));
							$columnElem.insertAfter($('tr[idpl=' + PLyo.PLCode + ']'));
						}

						if ($trElem.attr('idparent') == "PL33" || $trElem.attr('idparent') == "PL34" || $trElem.attr('idparent') == "PL35" || $trElem.attr('idparent') == 'PL33_allocated' || $trElem.attr('idparent') == 'PL34_allocated' || $trElem.attr('idparent') == 'PL35_allocated') {
							var texthtml = $trElem.find('td:eq(0)').text();
							$trElem.find('td:eq(0)').text(texthtml.substring(5, texthtml.length));
						}
					}
				}
			}

			if (idplyo2 != undefined) {
				$trElem.removeAttr('idparent');
				$trElem.addClass('bold');
				$trElem.css('display', 'inline-grid');
				$('.table-content tr.column' + $trElem.attr("idheaderpl")).removeAttr("idcontparent");
				$('.table-content tr.column' + $trElem.attr("idheaderpl")).attr('statusval', 'show');
				$('.table-content tr.column' + $trElem.attr("idheaderpl")).attr('statusvaltemp', 'show');
				$('.table-content tr.column' + $trElem.attr("idheaderpl")).css('display', 'inline-grid');
			}

			toolkit.try(function () {
				$trElem.find('td:eq(0)')[0].childNodes.forEach(function (g) {
					if (g.nodeName == '#text') {
						if (g.nodeValue.indexOf('_') > -1) {
							g.nodeValue = g.nodeValue.split('_').reverse()[0];
						}
					}
				});
			});
		}
	});

	var countChild = '';
	$(".table-header tbody>tr").each(function (i) {
		$trElem = $(this);
		parenttr = $('tr[idparent=' + $trElem.attr('idheaderpl') + ']').length;
		if (parenttr > 0) {
			$trElem.addClass('dd');
			$trElem.find('td:eq(0)>i').addClass('fa fa-chevron-right').css('margin-right', '5px');
			$('tr[idparent=' + $trElem.attr('idheaderpl') + ']').css('display', 'none');
			$('tr[idcontparent=' + $trElem.attr('idheaderpl') + ']').css('display', 'none');
			$('tr[idparent=' + $trElem.attr('idheaderpl') + ']').each(function (a, e) {
				if ($(e).attr('statusval') == 'show') {
					$('tr[idheaderpl=' + $trElem.attr('idheaderpl') + ']').attr('statusval', 'show');
					$('tr[idpl=' + $trElem.attr('idheaderpl') + ']').attr('statusval', 'show');
					if ($('tr[idheaderpl=' + $trElem.attr('idheaderpl') + ']').attr('idparent') == undefined) {
						$('tr[idpl=' + $trElem.attr('idheaderpl') + ']').css('display', '');
						$('tr[idheaderpl=' + $trElem.attr('idheaderpl') + ']').css('display', '');
					}
				}
			});
		} else {
			countChild = $trElem.attr('idparent');
			if (countChild == '' || countChild == undefined) $trElem.find('td:eq(0)').css('padding-left', '20px');
		}
	});

	var prev = rpt.arrChangeParent()[rpt.arrChangeParent().length - 1].idfrom;
	if (prev.indexOf('_Allocated') > -1 || prev.indexOf('_Direct') > -1) {
		var under = ['PL94B', 'PL44B', 'PL44C', 'PL44E', 'PL44D', 'PL44F'];
		under.forEach(function (d) {
			var headerFrom = $('[idheaderpl="' + d + '"]');
			var headerTo = $('[idheaderpl="' + prev + '"]');

			if (headerFrom.size() == 0 || headerTo.size() == 0) {
				return;
			}

			headerFrom.insertAfter(headerTo);

			var contentFrom = $('[idpl="' + d + '"]');
			var contentTo = $('[idpl="' + prev + '"]');
			contentFrom.insertAfter(contentTo);

			prev = d;
		});
	}

	rpt.showZeroValue(false);
	rpt.hideSubGrowthValue();
	$(".pivot-pnl .table-header tr:not([idparent]):not([idcontparent])").addClass('bold');
	rpt.refreshHeight();
	rpt.addScrollBottom();
};

rpt.refreshchildadd = function (e) {
	var $columnElem = void 0,
	    $trElem = void 0;
	$('.table-header tbody>tr[idparent=' + e + ']').each(function (i) {
		$trElem = $(this);
		$columnElem = $('.table-content tbody>tr[idpl=' + $trElem.attr('idheaderpl') + ']');
		if (e == 'PL0') {
			$trElem.insertBefore($('tr[idheaderpl=' + $trElem.attr('idparent') + ']'));
			$columnElem.insertBefore($('tr[idpl=' + $trElem.attr('idparent') + ']'));
		} else {
			$trElem.insertAfter($('tr[idheaderpl=' + $trElem.attr('idparent') + ']'));
			$columnElem.insertAfter($('tr[idpl=' + $trElem.attr('idparent') + ']'));
		}
	});
};

rpt.hideSubGrowthValue = function () {
	toolkit.repeat(8, function (i) {
		$('[idheaderpl="PL' + (i + 1) + '"] td:contains("%")').html('&nbsp;');
		$('[idpl="PL' + (i + 1) + '"] td:contains("%")').html('&nbsp;');
	});
};

rpt.hideAllChild = function (PLCode) {
	$('.table-header tbody>tr[idparent=' + PLCode + ']').each(function (i) {
		var $trElem = $(this);
		var child = $('tr[idparent=' + $trElem.attr('idheaderpl') + ']').length;
		if (child > 0) {
			var $c = $('tr[idheaderpl=' + $trElem.attr('idheaderpl') + ']');
			$($c).find('i').removeClass('fa-chevron-up');
			$($c).find('i').addClass('fa fa-chevron-right');
			$('tr[idparent=' + $c.attr('idheaderpl') + ']').css('display', 'none');
			$('tr[idcontparent=' + $c.attr('idheaderpl') + ']').css('display', 'none');
			rpt.hideAllChild($c.attr('idheaderpl'));
		}
	});
};

rpt.putStatusVal = function (trHeader, trContent) {
	var boolStatus = false;
	trContent.find('td').each(function (a, e) {
		var text = $(e).html().replace(/&nbsp;/g, ' ');
		if (text != '0' && text != '0.00 %') {
			boolStatus = true;
		}
	});

	if (boolStatus) {
		trContent.attr('statusval', 'show');
		trHeader.attr('statusval', 'show');
	} else {
		trContent.attr('statusval', 'hide');
		trHeader.attr('statusval', 'hide');
	}
};

rpt.refreshHeight = function (PLCode) {
	$('.table-header tbody>tr[idparent=' + PLCode + ']').each(function (i) {
		var $trElem = $(this);
		$('tr[idcontparent=' + $trElem.attr('idheaderpl') + ']').css('height', $trElem.height());
	});
};

rpt.showExport = ko.observable(false);
rpt.export = function (target, title, mode) {
	target = toolkit.$(target);

	if (mode == 'kendo') {
		var workbook;

		var _ret2 = function () {
			var rowdata = [],
			    cellval = {},
			    cells = [],
			    headertype = "";
			var tableHeaderLock = target.find('.k-grid-header-locked');
			var tableHeader = target.find('.k-grid-header-wrap');
			var tableContentLock = target.find('.k-grid-content-locked');
			var tableContent = target.find('.k-grid-content');
			tableHeaderLock.find('tr').each(function (i, e) {
				cells = [];
				$(e).find('th').each(function (i, e) {
					cellval = {};
					cellval['value'] = $(e).attr('data-title');
					if ($(e).attr('rowspan')) {
						if (title == 'Distribution Analysis') cellval['rowSpan'] = parseInt($(e).attr('rowspan')) + 2;else if (title == 'Summary P&L Analysis') cellval['rowSpan'] = parseInt($(e).attr('rowspan')) + 1;else cellval['rowSpan'] = parseInt($(e).attr('rowspan'));
					}
					if ($(e).attr('colspan')) cellval['colSpan'] = parseInt($(e).attr('colspan'));
					cells.push(cellval);
				});
				rowdata.push({ cells: cells });
			});
			tableHeader.find('tr').each(function (a, e) {
				cells = [];
				$(e).find('th').each(function (i, e) {
					cellval = {};
					cellval['value'] = $(e).attr('data-title');
					if ($(e).attr('rowspan')) cellval['rowSpan'] = parseInt($(e).attr('rowspan'));
					if ($(e).attr('colspan')) cellval['colSpan'] = parseInt($(e).attr('colspan'));
					if (rowdata[a]) rowdata[a].cells.push(cellval);else cells.push(cellval);
				});
				if (cells.length > 0) rowdata.push({ cells: cells });
			});
			tableContentLock.find('tr').each(function (i, e) {
				cells = [];
				$(e).find('td').each(function (i, e) {
					cellval = {};
					headertype = parseFloat($(e).html().replace(/,/g, ""));
					if (isNaN(parseFloat(headertype)) == true) headertype = $(e).html();
					cellval['value'] = headertype;
					cells.push(cellval);
				});
				tableContent.find('tr:eq(' + i + ') td').each(function (i, e) {
					cellval = {};
					cellval['value'] = parseFloat($(e).html().replace(/,/g, ""));
					cells.push(cellval);
				});
				rowdata.push({ cells: cells });
			});
			// console.log(rowdata)
			workbook = new kendo.ooxml.Workbook({
				sheets: [{
					columns: [{ autoWidth: true }, { autoWidth: true }],
					title: title,
					rows: rowdata
				}]
			});

			kendo.saveAs({
				dataURI: workbook.toDataURL(),
				fileName: title + ".xlsx"
			});
			return {
				v: void 0
			};
		}();

		if ((typeof _ret2 === 'undefined' ? 'undefined' : _typeof(_ret2)) === "object") return _ret2.v;
	} else if (mode == "kendonormal") {
		var workbook;

		var _ret3 = function () {
			var rowdata = [],
			    cellval = {},
			    cells = [],
			    headertype = "";
			var tableHeader = target.find('.k-grid-header-wrap');
			var tableContent = target.find('.k-grid-content');
			var tableFooter = target.find('.k-grid-footer');
			tableHeader.find('tr').each(function (i, e) {
				cells = [];
				$(e).find('th').each(function (i, e) {
					cellval = {};
					cellval['value'] = $(e).attr('data-title');
					if ($(e).attr('rowspan')) {
						if (title == 'Distribution Analysis') cellval['rowSpan'] = parseInt($(e).attr('rowspan')) + 2;else if (title == 'Summary P&L Analysis') cellval['rowSpan'] = parseInt($(e).attr('rowspan')) + 1;else cellval['rowSpan'] = parseInt($(e).attr('rowspan'));
					}
					if ($(e).attr('colspan')) cellval['colSpan'] = parseInt($(e).attr('colspan'));
					cells.push(cellval);
				});
				rowdata.push({ cells: cells });
			});
			tableContent.find('tr').each(function (i, e) {
				cells = [];
				$(e).find('td').each(function (i, e) {
					cellval = {};
					headertype = parseFloat($(e).html().replace(/,/g, ""));
					if (isNaN(parseFloat(headertype)) == true) headertype = $(e).html();
					cellval['value'] = headertype;
					cells.push(cellval);
				});
				if (cells.length > 0) rowdata.push({ cells: cells });
			});
			tableFooter.find('tr').each(function (i, e) {
				cells = [];
				$(e).find('td').each(function (i, e) {
					cellval = {};
					if (i == 0) headertype = parseFloat($(e).html().replace(/,/g, ""));else {
						if ($(e).find('div').length > 0) headertype = parseFloat($(e).find('div').html().replace(/,/g, ""));else headertype = parseFloat($(e).html().replace(/,/g, ""));
					}

					if (isNaN(parseFloat(headertype)) == true) {
						headertype = $(e).text();
					}
					cellval['value'] = headertype;
					cells.push(cellval);
				});
				if (cells.length > 0) rowdata.push({ cells: cells });
			});
			// console.log(rowdata)
			workbook = new kendo.ooxml.Workbook({
				sheets: [{
					columns: [{ autoWidth: true }, { autoWidth: true }],
					title: title,
					rows: rowdata
				}]
			});

			kendo.saveAs({
				dataURI: workbook.toDataURL(),
				fileName: title + ".xlsx"
			});
			return {
				v: void 0
			};
		}();

		if ((typeof _ret3 === 'undefined' ? 'undefined' : _typeof(_ret3)) === "object") return _ret3.v;
	} else if (mode == 'normal') {
		(function () {
			$('#fake-table').remove();

			var body = $('body');
			var fakeTable = $('<table />').attr('id', 'fake-table').appendTo(body);

			if (target.attr('name') != 'table') {
				target = target.find('table:eq(0)');
			}

			target.clone(true).appendTo(fakeTable);

			var downloader = $('<a />').attr('href', '#').attr('download', title + '.xls').attr('onclick', 'return ExcellentExport.excel(this, \'fake-table\', \'sheet1\')').html('export').appendTo(body);

			fakeTable.find('td').css('height', 'inherit');
			downloader[0].click();

			setTimeout(function () {
				fakeTable.remove();
				downloader.remove();
			}, 400);
		})();
	} else if (mode == 'header-content') {
		(function () {
			$('#fake-table').remove();

			var body = $('body');
			var fakeTable = $('<table />').attr('id', 'fake-table').appendTo(body);

			var tableHeader = target.find('.table-header');
			if (tableHeader.attr('name') != 'table') {
				tableHeader = tableHeader.find('table');
			}

			var tableContent = target.find('.table-content');
			if (tableContent.attr('name') != 'table') {
				tableContent = tableContent.find('table');
			}

			tableHeader.find('tr').each(function (i, e) {
				$(e).css('height', '');

				if (i == 0) {
					var rowspan = parseInt($(e).find('td,th').attr('data-rowspan'), 10);
					if (isNaN(rowspan)) rowspan = 1;

					for (var j = 0; j < rowspan; j++) {
						$(e).clone(true).appendTo(fakeTable);
					}
					return;
				}

				$(e).clone(true).appendTo(fakeTable);
			});

			tableContent.find('tr').each(function (i, e) {
				$(e).css('height', '');

				var rowTarget = fakeTable.find('tr:eq(' + i + ')');
				$(e).find('td,th').each(function (j, f) {
					$(f).clone(true).appendTo(rowTarget);
				});
			});

			var downloader = $('<a />').attr('href', '#').attr('download', title + '.xls').attr('onclick', 'return ExcellentExport.excel(this, \'fake-table\', \'sheet1\')').html('export').appendTo(body);

			fakeTable.find('tr:hidden').show();
			fakeTable.find('td,th').css('height', 'inherit');
			fakeTable.find('td .fa-chevron-right').remove();

			downloader[0].click();

			setTimeout(function () {
				fakeTable.remove();
				downloader.remove();
			}, 400);
		})();
	}
};

rpt.addScrollBottom = function (container) {
	$('.scroll-grid-bottom-yo').remove();
	$('.scroll-grid-bottom').remove();

	if (container == undefined) container = $(".breakdown-view:visible");

	toolkit.newEl('div').addClass('scroll-grid-bottom-yo').appendTo(container.find(".pivot-pnl"));

	var tableContent = toolkit.newEl('div').addClass('scroll-grid-bottom').appendTo(container.find(".pivot-pnl"));

	var arrowLeft = toolkit.newEl('div').addClass('scroll-grid-bottom arrow arrow-left viewscrollfix').html('<i class="fa fa-arrow-left"></i>').appendTo(container.find(".pivot-pnl"));

	var arrowRight = toolkit.newEl('div').addClass('scroll-grid-bottom arrow arrow-right viewscrollfix').html('<i class="fa fa-arrow-right"></i>').appendTo(container.find(".pivot-pnl"));

	toolkit.newEl('div').addClass('content-grid-bottom')
	// .css("min-width", container.find('.table-content>.table').width() - 48)
	.html("&nbsp;").appendTo(tableContent);

	var target = container.find(".scroll-grid-bottom")[0];
	var target2 = container.find(".table-content")[0];
	container.find(".table-content").scroll(function () {
		target.scrollLeft = this.scrollLeft;
	});
	container.find(".scroll-grid-bottom").scroll(function () {
		target2.scrollLeft = this.scrollLeft;
	});

	var walkLength = 50;

	arrowLeft.on('click', function () {
		var newVal = target.scrollLeft - walkLength;
		if (newVal < 0) {
			newVal = 0;
		}

		target.scrollLeft = newVal;
	});
	arrowRight.on('click', function () {
		var newVal = target.scrollLeft + walkLength;
		if (newVal < 0) {
			newVal = 0;
		}

		target.scrollLeft = newVal;
	});

	rpt.panel_scrollrelocated();
};

rpt.panel_scrollrelocated = function () {
	$(".scroll-grid-bottom").each(function (i) {
		$(this).find('.content-grid-bottom').css("min-width", $(this).parent().find('.table-content>.table').width());
		if ($(this).parent().find('.scroll-grid-bottom-yo').size() == 0) {
			return;
		}

		var window_top = $(window).scrollTop() + $(window).innerHeight();
		var div_top = $(this).parent().find('.scroll-grid-bottom-yo').offset().top;
		if (parseInt(div_top, 10) < parseInt(window_top, 10)) {
			$(this).removeClass('viewscrollfix');
			$(this).hide();
			$(this).css("width", "100%");
		} else {
			$(this).show();
			$(this).css("width", $(this).parent().find('.table-content').width());
			if (!$(this).hasClass('viewscrollfix')) $(this)[0].scrollLeft = $(this).parent().find(".table-content")[0].scrollLeft;
			$(this).addClass('viewscrollfix');
		}
	});
};

rpt.orderByChannel = function (what, defaultValue) {
	toolkit.try(function () {
		what = what.replace(/&nbsp;/g, ' ');
	});

	switch (what) {
		case 'Branch Modern Trade':
		case 'Modern Trade':
		case 'MT':
			return 2000000000005;
		case 'Branch General Trade':
		case 'General Trade':
		case 'GT':
			return 2000000000004;
		case 'Regional Distributor':
		case 'RD':
			return 2000000000003;
		case 'Industrial Trade':
		case 'Industrial':
		case 'IT':
			return 2000000000002;
		case 'Motorist':
			return 2000000000001;
		case 'Export':
			return 2000000000000;
	}

	return defaultValue;
};

rpt.clickExpand = function (container, e) {
	var right = $(e).find('i.fa-chevron-right').length,
	    down = 0;
	if (e.attr('idheaderpl') == 'PL0') down = $(e).find('i.fa-chevron-up').length;else down = $(e).find('i.fa-chevron-down').length;
	if (right > 0) {
		if (['PL28', 'PL29A', 'PL31'].indexOf($(e).attr('idheaderpl')) > -1) {
			container.find('.pivot-pnl .table-header').css('width', rpt.pnlTableHeaderWidth());
			container.find('.pivot-pnl .table-content').css('margin-left', rpt.pnlTableHeaderWidth());
		}

		$(e).find('i').removeClass('fa-chevron-right');
		if (e.attr('idheaderpl') == 'PL0') $(e).find('i').addClass('fa fa-chevron-up');else $(e).find('i').addClass('fa fa-chevron-down');
		container.find('tr[idparent=' + e.attr('idheaderpl') + ']').css('display', '');
		container.find('tr[idcontparent=' + e.attr('idheaderpl') + ']').css('display', '');
		container.find('tr[statusvaltemp=hide]').css('display', 'none');
		rpt.refreshHeight(e.attr('idheaderpl'));
		rpt.refreshchildadd(e.attr('idheaderpl'));
	}
	if (down > 0) {
		if (['PL28', 'PL29A', 'PL31'].indexOf($(e).attr('idheaderpl')) > -1) {
			$('.pivot-pnl .table-header').css('width', '');
			$('.pivot-pnl .table-content').css('margin-left', '');
		}

		$(e).find('i').removeClass('fa-chevron-up');
		$(e).find('i').removeClass('fa-chevron-down');
		$(e).find('i').addClass('fa fa-chevron-right');
		container.find('tr[idparent=' + e.attr('idheaderpl') + ']').css('display', 'none');
		container.find('tr[idcontparent=' + e.attr('idheaderpl') + ']').css('display', 'none');
		rpt.hideAllChild(e.attr('idheaderpl'));
	}
};

$(function () {
	$(window).on('scroll', function () {
		rpt.panel_scrollrelocated();
	});

	rpt.getOtherMasterData();
});