'use strict';

viewModel.PCAnalysis = new Object();
var pg = viewModel.PCAnalysis; // for parent variable interaction
pg.currentObject = ko.observable({});

vm.currentMenu('Power Curve Analysis');
vm.currentTitle('Power Curve Analysis');
vm.breadcrumb([
    { title: 'Power Curve', href: '#' },
    { title: 'Power Curve Analysis', href: viewModel.appName + 'page/powercurveanalysis' },
]);

pg.setPc = function() {
    if(typeof pc !== 'undefined') {
        pg.currentObject(pc);
        pg.refresh();
    }
}
pg.setMs = function() {
    if(typeof ms !== 'undefined') {
        pg.currentObject(ms);
        pg.refresh();
    }
}
pg.setIm = function() {
    if(typeof im !== 'undefined') {
        pg.currentObject(im);
        pg.refresh();
    }
}
pg.setCm = function() {
    if(typeof cm !== 'undefined') {
        pg.currentObject(cm);
        pg.refresh();
    }
}
pg.setSc = function() {
    if(typeof sc !== 'undefined') {
        pg.currentObject(sc);
        pg.refresh();
    }
}
pg.setOp = function() {
    if(typeof op !== 'undefined') {
        pg.currentObject(op);
        pg.refresh();
    }
}
pg.resetAllTabs = function() {
    if(typeof pc !== 'undefined') {
        if(typeof pc.reset !== 'undefined') {
            pc.reset();
        }
    }
    if(typeof ms !== 'undefined') {
        if(typeof ms.reset !== 'undefined') {
            ms.reset();
        }
    }
    if(typeof im !== 'undefined') {
        if(typeof im.reset !== 'undefined') {
            im.reset();
        }
    }
    if(typeof cm !== 'undefined') {
        if(typeof cm.reset !== 'undefined') {
            cm.reset();
        }
    }
    if(typeof sc !== 'undefined') {
        if(typeof sc.reset !== 'undefined') {
            sc.reset();
        }
    }
    if(typeof op !== 'undefined') {
        if(typeof op.reset !== 'undefined') {
            op.reset();
        }
    }
}
pg.refresh = function() {
    var currentObj = pg.currentObject();
    if(typeof currentObj.refresh !== 'undefined') {
        currentObj.refresh();
    }
}

pg.initLoad = function() {
    window.setTimeout(function(){
        fa.LoadData();
        di.getAvailDate();
        pg.setPc();
        app.loading(false);
    }, 1000);
}

$(function(){
    $('#projectList').kendoDropDownList({
        data: fa.projectList,
        dataValueField: 'value',
        dataTextField: 'text',
        suggest: true,
        change: function () {
            setTimeout(function(){
                var project = $('#projectList').data("kendoDropDownList").value();
                fa.project = project;
                var projectName = $('#projectList').data("kendoDropDownList").dataItem().value;
                fa.populateEngine(projectName);
                fa.populateTurbine(project);
                di.getAvailDate();
            },500);
         }
    });
    
    $('#btnRefresh').on('click', function () {
        pg.resetAllTabs();
        pg.refresh();
    });

    $('#t-power-curve a').on('click', function(){
        pg.setPc();
    });
    $('#t-monthly-scatter a').on('click', function(){
        pg.setMs();
    });
    $('#t-individual-month a').on('click', function(){
        pg.setIm();
    });
    $('#t-comparison a').on('click', function(){
        pg.setCm();
    });
    $('#t-scatter a').on('click', function(){
        pg.setSc();
    });
    $('#t-operational a').on('click', function(){
        pg.setOp();
    });

    pg.initLoad();
});