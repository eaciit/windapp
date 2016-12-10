"use strict";

viewModel.login = new Object();
var lg = viewModel.login;

lg.templateConfigLogin = {
	username: "",
	password: ""
};

lg.templateForgotLogin = {
	email: "",
	baseurl: ""
};

lg.templateConfirmReset = {
	new_pass: "",
	confirm_pass: ""
};

lg.templateUrlParam = {
	userid: "",
	tokenid: "",
	newpassword: ""

};

lg.configLogin = ko.mapping.fromJS(lg.templateConfigLogin);
lg.forgetLogin = ko.mapping.fromJS(lg.templateForgotLogin);
lg.confirmReset = ko.mapping.fromJS(lg.templateConfirmReset);
lg.dataMenu = ko.observableArray([]);
lg.ErrorMessage = ko.observable('');
lg.getConfirReset = ko.mapping.fromJS(lg.templateUrlParam);

lg.getLogin = function () {
	if (!toolkit.isFormValid("#login-form")) {
		return;
	}

	lg.showLoader(true);

	var param = ko.mapping.toJS(lg.configLogin);
	toolkit.ajaxPost(viewModel.appName + "login/processlogin", param, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}

		lg.ErrorMessage(res.message);

		if (res.message == "Login Success") {
			window.location = viewModel.appName + "page/dashboard";
		}

		lg.showLoader(false);
	});
};

lg.showAccesReset = function () {
	$('#modalForgot').modal({ show: 'true' });
	lg.forgetLogin.email('');
};

lg.getForgetLogin = function () {
	if (!toolkit.isFormValid("#email-form")) {
		$('#modalForgot').modal({
			backdrop: 'static',
			keyboard: false
		});
		return;
	}
	var url = lg.forgetLogin.baseurl(location.origin);
	var param = ko.mapping.toJS(lg.forgetLogin);

	toolkit.ajaxPost(viewModel.appName + "login/resetpassword", param, function (res) {
		if (!toolkit.isFine(res)) {
			return;
		}
	});

	$('#modalForgot').modal('hide');
};

lg.getUrlParam = function () {
	var url = new RegExp('[\?&]' + param + '=([^&#]*)').exec(window.location.href);
	return url[1] || 0;
};

lg.getConfirmReset = function () {
	if (!toolkit.isFormValid("#form-reset")) {
		return;
	}

	if (lg.confirmReset.confirm_pass() != lg.confirmReset.new_pass()) {
		lg.ErrorMessage('Your confirm new password not match');
	} else {
		lg.getConfirReset.userid(lg.getUrlParam('1'));
		lg.getConfirReset.tokenid(lg.getUrlParam('2'));
		lg.getConfirReset.newpassword(lg.confirmReset.confirm_pass());
		var param = ko.mapping.toJS(lg.getConfirReset);
		toolkit.ajaxPost(viewModel.appName + "login/savepassword", param, function (res) {
			if (!toolkit.isFine(res)) {
				return;
			}

			window.location = viewModel.appName + "page/login";
		});
	}
};

lg.checkSession = function () {
	toolkit.ajaxPost(viewModel.appName + 'login/checkcurrentsession', {}, function (res) {
		if (res.data == true && res.message == "active") {
			window.location = viewModel.appName + "page/dashboard";
		}
	});
}

lg.showLoader = function(visible) {
	if(visible) {
		$('div.loader').show();
		$('div.login-form-bg').hide();
	} else {
		$('div.loader').hide();
		$('div.login-form-bg').show();
	}
}

$(function () {
	lg.showLoader(false);
	// lg.checkSession();
	// setTimeout(function() {
	//        $('#loginForm').height($(window).height());
	//    }, 300);
});