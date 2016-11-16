$.fn.ecDataBrowser = function (method) {
	if (methodsDataBrowser[method]) {
		return methodsDataBrowser[method].apply(this, Array.prototype.slice.call(arguments, 1));
	} else {
		methodsDataBrowser['init'].apply(this,arguments);
	}
}
// Format : Integer, Double, Float, Currency, Date, DateTime and Date Format
var Setting_DataBrowser = {
	title: "",
	widthPerColumn: 4,
	widthKeyFilter: 3,
	showFilter: "Simple",
	dataSource: {
		data:[],
	},
	metadata: [],
	dataSimple: [],
	dataAdvance: [],
};
var Setting_ColumnGridDB = {
	Field: "",
    Label: "",
    DataType: "",
    Format: "",
    Align: "",
    ShowIndex: 1,
    Sortable: true,
    SimpleFilter: true,
    AdvanceFilter: true,
    Aggregate: ""
};
var SettingDataSourceBrowser = {
	data: [],
	url: "",
	type: "",
	fieldTotal: "",
	fieldData: "",
	serverPaging: true,
	pageSize: 10,
	serverSorting: true,
	callOK: function(a){

	}, 
	callFail: function(a,b,c){

	},
}
var Setting_TypeData = {
	number: ['integer', 'int', 'double', 'float', 'n0'],
	date: ['date','datetime'],
}

var methodsDataBrowser = {
	init: function(options){
		var databrowser = $.extend({}, SettingDataSourceBrowser, options.dataSource || {});
		var settings = $.extend({}, Setting_DataBrowser, options || {});
		settings.dataSource = databrowser;
		var sortMeta = settings.metadata.sort(function(a, b) {
		    return parseFloat(a.ShowIndex) - parseFloat(b.ShowIndex);
		});
		settings.metadata = sortMeta;
		// var settingDataSources = $.extend({}, Setting_DataBrowser, settings['dataSource'] || {});
		return this.each(function () {
			// $(this).data("ecDataSource", settingDataSources);
			$(this).data("ecDataBrowser", new $.ecDataBrowserSetting(this, settings));
			methodsDataBrowser.createElement(this, settings);
		});
	},
	createElement: function(element, options){
		$(element).html("");
		var $o = $(element), settingFilter = {}, widthfilter = 0, dataSimple= [], dataAdvance= [];

		$divFilterSimple = $('<div class="col-md-12 ecdatabrowser-filtersimple"></div>');
		$divFilterSimple.appendTo($o);
		$divFilterAdvance = $('<div class="col-md-12 ecdatabrowser-filteradvance"></div>');
		$divFilterAdvance.appendTo($o);
		for (var key in options.metadata){
			settingFilter = $.extend({}, Setting_ColumnGridDB, options.metadata[key] || {});
			widthfilter = 12-options.widthKeyFilter;
			if (settingFilter.SimpleFilter){
				$divFilter = $('<div class="col-md-'+options.widthPerColumn+' filter-'+key+'"></div>');
				$divFilter.appendTo($divFilterSimple);
				$labelFilter = $("<label class='col-md-"+options.widthKeyFilter+" ecdatabrowser-filter'>"+settingFilter.Label+"</label>");
				$labelFilter.appendTo($divFilter);
				$divContentFilter = $('<div class="col-md-'+widthfilter+' filter-form"></div>');
				$divContentFilter.appendTo($divFilter);
				methodsDataBrowser.createElementFilter(settingFilter, 'simple', key, $divContentFilter, $o);
				dataSimple.push('filter-simple-'+key);
			}
			if (settingFilter.AdvanceFilter){
				$divFilter = $('<div class="col-md-'+options.widthPerColumn+' filter-'+key+'"></div>');
				$divFilter.appendTo($divFilterAdvance);
				$labelFilter = $("<label class='col-md-"+options.widthKeyFilter+" ecdatabrowser-filter'>"+settingFilter.Label+"</label>");
				$labelFilter.appendTo($divFilter);
				$divContentFilter = $('<div class="col-md-'+widthfilter+' filter-form"></div>');
				$divContentFilter.appendTo($divFilter);
				methodsDataBrowser.createElementFilter(settingFilter, 'advance', key, $divContentFilter, $o);
				dataAdvance.push('filter-advance-'+key);
			}
		}
		$(element).data("ecDataBrowser").dataSimple = dataSimple;
		$(element).data("ecDataBrowser").dataAdvance = dataAdvance;

		$divContainerGrid = $('<div class="ecdatabrowser-gridview" style="width: 100%;"></div>');
		$divContainerGrid.appendTo($o);

		$divGrid = $('<div class="ecdatabrowser-grid"></div>');
		$divGrid.appendTo($divContainerGrid);

		methodsDataBrowser.createGrid($divGrid, options, $o);

		$("<div class='clearfix'></div>").appendTo($o);
		$(element).data("ecDataBrowser").ChangeViewFilter(options.showFilter);

	},
	createElementFilter: function(settingFilter, filterchoose, index, element, id){
		var $divElementFilter;
		if (settingFilter.DataType.toLowerCase() == 'integer' || settingFilter.DataType.toLowerCase() == "float32" || settingFilter.DataType.toLowerCase() == 'int' || settingFilter.DataType.toLowerCase() == 'float64' || settingFilter.DataType.toLowerCase() == 'date'){
			$divElementFilter = $('<input type="checkbox" class="ecdatabrowser-ckcrange"/>');
			$divElementFilter.bind('click').click(function(){
				if ($(this).prop("checked")){
					$(this).parent().find('.ecdatabrowser-spacerange').show();
					$(this).parent().find('.ecdatabrowser-filterto').css('display','inline-table');
				}else{
					$(this).parent().find('.ecdatabrowser-spacerange').hide();
					$(this).parent().find('.ecdatabrowser-filterto').hide();
				}
			});
			$divElementFilter.appendTo(element);
		}
		if (settingFilter.DataType.toLowerCase() == 'integer' || settingFilter.DataType.toLowerCase() == "float32" || settingFilter.DataType.toLowerCase() == 'int' || settingFilter.DataType.toLowerCase() == 'float64'){
			$divElementFilter = $('<input class="ecdatabrowser-filterfrom" idfilter="filter-'+filterchoose+'-'+index+'" typedata="'+settingFilter.DataType.toLowerCase()+'" fielddata="'+ settingFilter.Field +'"/>');
			$divElementFilter.appendTo(element);
			$divElementFilter = $('<span class="ecdatabrowser-spacerange"> - </span><input class="ecdatabrowser-filterto" idfilter="filter-'+filterchoose+'-'+index+'" typedata="'+settingFilter.DataType.toLowerCase()+'" fielddata="'+ settingFilter.Field +'"/>');
			$divElementFilter.appendTo(element);
			id.find('input[idfilter=filter-'+filterchoose+'-'+index+']').kendoNumericTextBox();
			return '';
		}
		else if (settingFilter.DataType.toLowerCase() == 'date'){
			$divElementFilter = $('<input class="ecdatabrowser-filterfrom" style="width: 100px;" idfilter="filter-'+filterchoose+'-'+index+'" typedata="date" fielddata="'+ settingFilter.Field +'"/>');
			$divElementFilter.appendTo(element);
			$divElementFilter = $('<span class="ecdatabrowser-spacerange"> - </span><input class="ecdatabrowser-filterto" style="width: 100px;" idfilter="filter-'+filterchoose+'-'+index+'" typedata="date" fielddata="'+ settingFilter.Field +'"/>');
			$divElementFilter.appendTo(element);
			id.find('input[idfilter=filter-'+filterchoose+'-'+index+']').kendoDatePicker({
				format: settingFilter.Format,
			});
			return '';
		} else if (settingFilter.DataType.toLowerCase() == 'bool') {
			$divElementFilter = $('<input type="checkbox" idfilter="filter-'+filterchoose+'-'+index+'" typedata="bool" fielddata="'+ settingFilter.Field +'"/>');
			$divElementFilter.appendTo(element);
			return '';
		}
		else {
			if (settingFilter.Lookup == false){
				$divElementFilter = $('<input type="text" class="form-control input-sm" idfilter="filter-'+filterchoose+'-'+index+'" typedata="string" fielddata="'+ settingFilter.Field +'" haslookup="false"/>');
				$divElementFilter.appendTo(element);
			} else {
				$divElementFilter = $('<input type="text" class="form-control input-sm" idfilter="filter-'+filterchoose+'-'+index+'" typedata="string" fielddata="'+ settingFilter.Field +'" haslookup="true"/>');
				$divElementFilter.appendTo(element);
				var callData = {};
				callData['browserid'] = id.data('ecDataBrowser').mapdatabrowser.dataSource.callData.browserid;
				callData['take'] = 10;
				callData['skip'] = 0;
				callData['page'] = 1;
				callData['pageSize'] = 10;
				callData['haslookup'] = true;
				callData['tablename'] = id.data('ecDataBrowser').mapdatabrowser.dataSource.callData.tablename;
				$('input[idfilter=filter-'+filterchoose+'-'+index+']').ecLookupDD({
					dataSource:{
						url: id.data('ecDataBrowser').mapdatabrowser.dataSource.url,
						call: 'post',
						callData: callData,
						resultData: function(a){
							return a.data.DataValue;
						}
					}, 
					inputType: 'multiple', 
					inputSearch: settingFilter.Field, 
					idField: settingFilter.Field, 
					idText: settingFilter.Field, 
					displayFields: settingFilter.Field, 
				});
			}
			return '';
		}
	},
	createGrid: function(element, options, id){
		var colums = [], format="", aggr= {}, footerText = "", column = {};
		for(var key in options.metadata){
			if ((options.metadata[key].DataType.toLowerCase() == 'integer' || options.metadata[key].DataType.toLowerCase() == "float32" || options.metadata[key].DataType.toLowerCase() == 'int' || options.metadata[key].DataType.toLowerCase() == 'float64') && options.metadata[key].Format != "" ){
				format = "{0:"+options.metadata[key].Format+"}"
			} else {
				format = "";
			}
			// aggr = JSON.parse("{\"avg\":\"220000.0000\",\"sum\":\"1100000\"}");
			aggr= {};
			if (options.metadata[key].Aggregate != '')
				aggr = JSON.parse(options.metadata[key].Aggregate);
			footerText = "";
			$.each( aggr, function( key, value ) {
				footerText+= key + ' : ' + value + '<br/>';
			});
			if (options.metadata[key].HiddenField != true){
				if (options.metadata[key].DataType.toLowerCase() == 'date'){
					if (options.metadata[key].Format != '')
						format = "moment(Date.parse("+options.metadata[key].Field+")).format('"+options.metadata[key].Format.toUpperCase()+"')";
					else 
						format = options.metadata[key].Field;
					column = {
						field: options.metadata[key].Field,
						title: options.metadata[key].Label,
						// format: format,
						sortable: options.metadata[key].Sortable,
						attributes: {
							style: "text-align: "+options.metadata[key].Align+";",
						},
						headerAttributes: {
							style: "text-align: "+options.metadata[key].Align+";",
						},
						aggregates: aggr,
						footerTemplate: footerText,
						template: "#:"+format+"#",
					};
				} else {
					column = {
						field: options.metadata[key].Field,
						title: options.metadata[key].Label,
						format: format,
						sortable: options.metadata[key].Sortable,
						attributes: {
							style: "text-align: "+options.metadata[key].Align+";",
						},
						headerAttributes: {
							style: "text-align: "+options.metadata[key].Align+";",
						},
						aggregates: aggr,
						footerTemplate: footerText,
					};
				}
				colums.push(column);
			}
		}
		column = {
			headerTemplate: "<a class='k-link align-center' href='#'>Action</a>", width: 80, attributes: { style: "text-align: center; cursor: pointer;"}, 
			headerAttributes: { style: "font-weight: bold;"},
			template: function (d) {
	    		return [
	    			"<button class='btn btn-xs btn-warning tooltipster' title='Edit data' onclick='db.editData("+JSON.stringify(d)+")'><span class='glyphicon glyphicon-pencil'></span></button>",
	    		].join(" ");
			}
		}
		colums.push(column);

		// colums = Lazy(colums).map(function (e, i) {
		// 	if (colums.length > 5) {
		// 		e.width = 150;

		// 		if (i == 0) {
		// 			e.width = 200;
		// 			e.locked = true;
		// 		}
		// 	}

		// 	return e;
		// }).toArray();

		$divElementGrid = $('<div idfilter="gridFilterBrowser"></div>');
		$divElementGrid.appendTo(element);
		if (options.dataSource.data.length > 0){
			id.find('div[idfilter=gridFilterBrowser]').kendoGrid({
				dataSource: {data: options.dataSource.data},
				sortable: true,
				columns: colums
			});
		} else {
			if (colums.length > 4) {
				let columnsLocked = []

				let firstColumn = colums.find((d) => d.field == '_id')
				if (app.isDefined(firstColumn)) {
					columnsLocked.push($.extend(true, firstColumn, {
						locked: true,
						width: 110
					}))
				}

				let lastColumn = colums.find((d) => app.isUndefined(d.field))
				if (app.isDefined(lastColumn)) {
					columnsLocked.push($.extend(true, lastColumn, {
						locked: true,
						width: 60
					}))
				}

				colums.filter((d) => d.field != '_id' && app.isDefined(d.field)).forEach((d) => {
					columnsLocked.push($.extend(true, d, {
						width: 150
					}))
				})

				colums = columnsLocked
			}

				// console.log(colums)

			id.find('div[idfilter=gridFilterBrowser]').kendoGrid({
				dataSource: {
					transport: {
	                    read: function(yo){
							var callData = $(id).data('ecDataBrowser').GetDataFilter(), $parentElem = id;
							// var callData = {}, $parentElem = id;
							$.each( options.dataSource.callData, function( key, value ) {
								callData[key] = value;
							});
							if (yo.data["sort"] == "")
  								yo.data["sort"] = undefined;
      						for(var i in yo.data){
		                        callData[i] = yo.data[i];
	                        }
				            app.ajaxPost($parentElem.data('ecDataBrowser').mapdatabrowser.dataSource.url,callData, function (res){
				            	yo.success(res.data);
				            	$parentElem.data('ecDataBrowser').mapdatabrowser.dataSource.callOK(res.data);
	                        });
	                    }
	                },
	                schema: {
	                    data: function(res){
	                    	if (res.dataresult.TableNames == 'salespls'){
	                    		res.DataValue.forEach((d) => {
	                    			colums.forEach((col) => {
	                    				if (!col.hasOwnProperty('field')) {
	                    					return
	                    				}
	                    				if (col.field.indexOf("pldatas") == -1) {
	                    					return
	                    				}

                    					let plcode = col.field.split('.')[1]
                    					if (!d.pldatas.hasOwnProperty(plcode)) {
                    						d.pldatas[plcode] = {
                    							plcode: "",
                    							plorder: "",
                    							group1: "",
                    							group2: "",
                    							group3: "",
                    							amount: 0
                    						}
                    					}
                    				})
	                    		})
		                    }
	                    	return res[options.dataSource.fieldData];
	                    },
	                    total: options.dataSource.fieldTotal
	                },
	                pageSize: options.dataSource.pageSize,
	                serverPaging: options.dataSource.serverPaging, // enable server paging
	                serverSorting: options.dataSource.serverSorting,
					// serverFiltering: true,
				},
				sortable: true,
	            pageable: true,
	            scrollable: true,
				columns: colums,
				dataBound: app.gridBoundTooltipster('div[idfilter=gridFilterBrowser]')
			});
		}
	},
	setShowFilter: function(res){
		$(this).data("ecDataBrowser").ChangeViewFilter(res);
	},
	getDataFilter: function(){
		var res = $(this).data('ecDataBrowser').GetDataFilter();
		return res;
	},
	postDataFilter: function(){
		var $dataBrowser = $(this).data('ecDataBrowser');
		if (typeof $dataBrowser !== "undefined") {
			$dataBrowser.refreshDataGrid(); 
		}
	},
	setDataGrid: function(res){
		// var mapNewGrid = $.extend({}, $(this).data("ecDataBrowser").mapdatabrowser, res || {});
		// var mapNewGrid = $(this).data("ecDataBrowser").mapdatabrowser
	}
}

$.ecDataBrowserSetting = function(element,options){
	this.mapdatabrowser = options;
	this.ChangeViewFilter = function(res){
		if (res.toLowerCase() == 'simple'){
			$(element).find('div.ecdatabrowser-filtersimple').show();
			$(element).find('div.ecdatabrowser-filteradvance').hide();
			this.mapdatabrowser.showFilter = "Simple";
		} else {
			$(element).find('div.ecdatabrowser-filtersimple').hide();
			$(element).find('div.ecdatabrowser-filteradvance').show();
			this.mapdatabrowser.showFilter = "Advance";
		}
	};
	this.CheckRangeData = function(findElem, typeData){
		$elemfrom = $(findElem+'.ecdatabrowser-filterfrom');
		$elemto = $(findElem+'.ecdatabrowser-filterto');
		if ($elemfrom.closest('.filter-form').find('.ecdatabrowser-ckcrange').prop("checked")){
			return $elemfrom.val() + '..' + $elemto.val();
		} else {
			var res = $elemfrom.val();
			if (typeData == 'float')
				return parseFloat(res);
			else if (typeData == 'int')
				return parseInt(res);
			else
				return res;
		}
	}
	this.GetDataFilter = function(){
		var resFilter = {}, dataTemp = [], $elem = '', valtype = '', lookupdata = [];
		if (this.mapdatabrowser.showFilter.toLowerCase() == "simple"){
			dataTemp = $(element).data('ecDataBrowser').dataSimple;
		} else {
			dataTemp = $(element).data('ecDataBrowser').dataAdvance;
		}
		for (var i in dataTemp){
			$elem = $('input[idfilter='+dataTemp[i]+']');
			field = $elem.attr('fielddata');
			if ($elem.val() != '' || $elem.attr('haslookup') == "true"){
				if ($elem.attr("typedata") == "integer" || $elem.attr("typedata") == "int" || $elem.attr("typedata") == "number"){
					// valtype = parseInt($elem.val());
					valtype = this.CheckRangeData('input[idfilter='+dataTemp[i]+']', 'int');
				} else if ($elem.attr("typedata") == "float32" || $elem.attr("typedata") == "float64"){
					// valtype = parseFloat($elem.val());
					valtype = this.CheckRangeData('input[idfilter='+dataTemp[i]+']', 'float');
				} else if ($elem.attr("typedata") == "bool"){
					valtype = $('input[idfilter='+dataTemp[i]+']')[0].checked;
				} else if ($elem.attr("typedata") == "date"){
					valtype = this.CheckRangeData('input[idfilter='+dataTemp[i]+']', 'date');
				} else {
					if ($elem.attr('haslookup') == "false")
						valtype = $elem.val();
					else {
						lookupdata = [];
						for(var a in $elem.ecLookupDD('get')){
							lookupdata.push($elem.ecLookupDD('get')[a][$elem.attr('fielddata')]);
						}
						if (lookupdata.length > 0)
							valtype = lookupdata;
						else if (lookupdata.length == 0 && $elem.ecLookupDD('gettext') != '')
							valtype = $elem.ecLookupDD('gettext');
						else
							valtype = '';
					}
				}
				if (valtype != '' || valtype.length > 0)
					resFilter[field] = valtype;
			}
		}
		return resFilter;
	};
	this.refreshDataGrid = function(){
		$('div[idfilter=gridFilterBrowser]').data('kendoGrid').dataSource.read();
		$('div[idfilter=gridFilterBrowser]').data('kendoGrid').refresh();
	}
}

// ecLookupDropdown
