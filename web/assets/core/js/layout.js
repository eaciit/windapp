'use strict';

var vm = viewModel;

vm.currentMenu = ko.observable('Dashboard');
vm.currentTitle = ko.observable('Dashboard');
vm.isDashboard = ko.observable(false);
vm.isShowDataAvailability = ko.observable(true);
vm.dateAsOf = ko.observable();
vm.menu = ko.observableArray([]); /*lek arep ngupdate lewat menu 'Operational / Menu Access 'yo*/
vm.breadcrumb = ko.observableArray([{ title: 'Windfarm', href: '#' }, { title: 'Dashboard', href: '#' }]);
vm.projectName = ko.observable("");
vm.dataAvailability = ko.observable("");

vm.getMenuList = function () {
    var isFine = function (res) {
        if (!res.success && (res.message.toLowerCase().indexOf("found") > -1
            || res.message.toLowerCase().indexOf("expired") > -1
            || res.message.toLowerCase().indexOf("failed") > -1)) {
            if (document.URL.indexOf(viewModel.appName + 'page/login') == -1) {
                swal({
                    title: "Warning",
                    type: "warning",
                    text: res.message,
                }, function () {
                    setTimeout(function () {
                        location.href = viewModel.appName + 'page/login';
                    }, 200);
                });

            }
            return false;
        }
        else if (!res.success && res.message.toLowerCase().indexOf("access to this page") > -1) {
            if (document.URL.indexOf(viewModel.appName + 'page/login') == -1) {
                swal({
                    title: "Warning",
                    type: "warning",
                    text: res.message,
                }, function () {
                    setTimeout(function () {
                        location.href = viewModel.appName + 'page/dashboard';
                    }, 200);
                });

            }
            return false;
        }
        return true;
    };

    toolkit.ajaxPost(viewModel.appName + "login/getmenulist", { url: document.URL }, function (res) {
        if (!isFine(res)) {
            return;
        }
        app.auth(true);
        vm.menu(res.data);
    });
}
vm.menuIcon = function (data) {
    return ko.computed(function () {
        return 'fa fa-' + data.icon;
    });
};

vm.prepareDropDownMenu = function () {
    $('ul.nav li.dd-hover').hover(function () {
        $(this).find('.dropdown-menu').stop(true, true).fadeIn(200);
    }, function () {
        $(this).find('.dropdown-menu').stop(true, true).fadeOut(200);
    });

};

vm.prepareFilterToggle = function () {
    $('.material-switch input[type="checkbox"]').on('change', function () {
        var show = $(this).is(':checked');
        var $target = $(this).closest('.panel').find('.panel-filter');
        if (show) {
            $target.show(200);
        } else {
            $target.hide(200);
        }
    }).trigger('click');
};

vm.setWidth = function(el){
    // $(el).find($(".dropdown-menu")).css("width",($(window).width() - $(el).offset().left));
    return false;
}

vm.setWidthUser = function(e){
    // $(e).find($(".dropdown-menu")).css("width",($(window).width() - $(e).offset().left) + 50);
    // // .navbar-nav.nav .dropdown-menu > li { float: left }
    // $(e).find($(".dropdown-menu li")).css("float" , "none");
    return false;
}

vm.adjustLayout = function () {
    var height = window.innerHeight - $('.main-header').height();
    $('.content-wrapper').css('min-height', height);
};
vm.showFilterCallback = toolkit.noop;
vm.showFilter = function () {
    var btnToggleFilter = $('.btn-toggle-filter');
    var panelFilterContainer = $('.panel-filter').parent();

    panelFilterContainer.removeClass('minimized');
    btnToggleFilter.find('.fa').removeClass('color-blue').addClass('color-orange').removeClass('fa-angle-double-right').addClass('fa-angle-double-left');

    $('.panel-filter').show(300);
    $('.panel-content').animate({ 'width': 'auto' }, 300, vm.showFilterCallback);
};
vm.hideFilterCallback = toolkit.noop;
vm.hideFilter = function () {
    var btnToggleFilter = $('.btn-toggle-filter');
    var panelFilterContainer = $('.panel-filter').parent();

    panelFilterContainer.addClass('minimized');
    btnToggleFilter.find('.fa').removeClass('color-orange').addClass('color-blue').removeClass('fa-angle-double-left').addClass('fa-angle-double-right');

    $('.panel-filter').hide(300);
    $('.panel-content').animate({ 'width': '100%' }, 300, vm.hideFilterCallback);
};
vm.prepareToggleFilter = function () {
    var btnToggleFilter = $('.btn-toggle-filter');
    var panelFilterContainer = $('.panel-filter').parent();

    $('<i class="fa fa-angle-double-left tooltipster align-center color-orange" title="Toggle filter pane visibility"></i>').appendTo(btnToggleFilter);
    toolkit.prepareTooltipster($(btnToggleFilter).find('.fa'));

    btnToggleFilter.on('click', function () {
        if (panelFilterContainer.hasClass('minimized')) {
            vm.showFilter();
        } else {
            vm.hideFilter();
        }
    });
};
vm.prepareLoader = function () {
    $('.loader canvas').each(function (i, cvs) {
        var ctx = cvs.getContext("2d");
        var sA = Math.PI / 180 * 45;
        var sE = Math.PI / 180 * 90;
        var ca = canvas.width;
        var ch = canvas.height;

        ctx.clearRect(0, 0, ca, ch);
        ctx.lineWidth = 15;

        ctx.beginPath();
        ctx.strokeStyle = "#ffffff";
        ctx.shadowColor = "#eeeeee";
        ctx.shadowOffsetX = 2;
        ctx.shadowOffsetY = 2;
        ctx.shadowBlur = 5;
        ctx.arc(50, 50, 25, 0, 360, false);
        ctx.stroke();
        ctx.closePath();

        sE += 0.05;
        sA += 0.05;

        ctx.beginPath();
        ctx.strokeStyle = "#aaaaaa";
        ctx.arc(50, 50, 25, sA, sE, false);
        ctx.stroke();
        ctx.closePath();
    });
};
vm.logout = function () {
    swal({
        title: "Are you sure want to logout?",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "btn-success",
        confirmButtonText: "Yes",
    },
        function (isconfirm) {
            if (isconfirm) {
                toolkit.ajaxPost(viewModel.appName + 'login/logout', {}, function (res) {
                    if (!app.isFine(res)) {
                        return;
                    }
                });
                setTimeout(function () {
                    location.href = viewModel.appName + 'page/login';
                }, 200);
            }
        });
    // toolkit.ajaxPost(viewModel.appName + 'login/logout', {}, function (res) {
    // 	if (!app.isFine(res)) {
    // 		return;
    // 	}
    // 	swal({
    // 		title: 'Logout Success',
    // 		text: 'Will automatically redirect to login page in 3 seconds',
    // 		type: 'success',
    // 		timer: 3000,
    // 		showConfirmButton: false
    // 	}, function () {
    // 		location.href = viewModel.appName + 'page/login';
    // 	});
    // });
};

$(function () {
    if (document.URL.indexOf(viewModel.appName + 'page/login') == -1) {
        $( '.navbar' ).append( '<span class="nav-bg"></span>' );
        jQuery(".dropdown-menu li.dropdown").on("mouseenter", function () {
            
             jQuery(this).find( ".dropdown-menu").css({ "display": "inline-block" }, "slow" );
             
            }).on("mouseleave", function () {
            
            jQuery(this).find( ".dropdown-menu").css({ "display": "none" }, "slow" );
             
        });
        // vm.getMenuList();
        // vm.prepareDropDownMenu();
        localStorage.clear();
        vm.prepareFilterToggle();
        vm.adjustLayout();
        vm.prepareToggleFilter();
        toolkit.prepareTooltipster();
        vm.prepareLoader();
    }
});

Date.prototype.getMonthName = function (lang) {
    lang = lang && (lang in Date.locale) ? lang : 'en';
    return Date.locale[lang].month_names[this.getMonth()];
};

Date.prototype.getMonthNameShort = function (lang) {
    lang = lang && (lang in Date.locale) ? lang : 'en';
    return Date.locale[lang].month_names_short[this.getMonth()];
};

Date.locale = {
    en: {
        month_names: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
        month_names_short: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    }
};

Date.prototype.addHours = function (h) {
    this.setTime(this.getTime() + (h * 60 * 60 * 1000));
    return this;
}