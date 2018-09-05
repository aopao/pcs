;(function () {
	'use strict';

	// Placeholder
	var placeholderFunction = function() {
		$('input, textarea').placeholder({ customClass: 'my-placeholder' });
	};

	// Placeholder
	var contentWayPoint = function() {
		var i = 0;
		$('.animate-box').waypoint( function( direction ) {

			if( direction === 'down' && !$(this.element).hasClass('animated-fast') ) {

				i++;

				$(this.element).addClass('item-animate');
				setTimeout(function(){

					$('body .animate-box.item-animate').each(function(k){
						var el = $(this);
						setTimeout( function () {
							var effect = el.data('animate-effect');
							if ( effect === 'fadeIn') {
								el.addClass('fadeIn animated-fast');
							} else if ( effect === 'fadeInLeft') {
								el.addClass('fadeInLeft animated-fast');
							} else if ( effect === 'fadeInRight') {
								el.addClass('fadeInRight animated-fast');
							} else {
								el.addClass('fadeInUp animated-fast');
							}

							el.removeClass('item-animate');
						},  k * 200, 'easeInOutExpo' );
					});

				}, 100);

			}

		} , { offset: '85%' } );
	};

	var loadAddress = function () {
        var len = province.length;
        for (var i = 0; i < len; i++) {
            // var provOpt = document.createElement('option');
            // provOpt.innerText = province[i]['name'];
            // provOpt.value = i;
            // prov.appendChild(provOpt);
            sign.signVue.provinces.push(province[i]['name'])
        }
    };
	// On load
	$(function(){
		placeholderFunction();
		contentWayPoint();
		loadAddress();
	});

}());

/*根据id获取对象*/
// function $(str) {
//     return document.getElementById(str);
// }

// var addrShow = $('addr-show');
var prov = document.getElementById('prov');
var city = document.getElementById('city');
var country = document.getElementById('country');

/*用于保存当前所选的省市区*/
var current = {
    prov: '',
    city: '',
    country: ''
};

/*自动加载省份列表*/
(function showProv() {

})();

/*根据所选的省份来显示城市列表*/
function showCity(obj) {
    var val = obj.options[obj.selectedIndex].value;
    if (val != current.prov) {
        current.prov = val;
        // addrShow.value = '';
    }
    //console.log(val);
    if (val != null) {
        city.length = 1;
        var cityLen = province[val]["city"].length;
        sign.signVue.province = val;
        sign.signVue.cities = [];
        sign.signVue.city = -1;
        for (var j = 0; j < cityLen; j++) {
            // var cityOpt = document.createElement('option');
            // cityOpt.innerText = province[val]["city"][j].name;
            // cityOpt.value = j;
            // city.appendChild(cityOpt);
            // console.log(city);
            sign.signVue.cities.push(province[val]["city"][j].name)
        }
    }
}

/*根据所选的城市来显示县区列表*/
function showCountry(obj) {
    var val = obj.options[obj.selectedIndex].value;
    current.city = val;
    if (val != null) {
        country.length = 1; //清空之前的内容只留第一个默认选项
        var countryLen = province[current.prov]["city"][val].districtAndCounty.length;
        if(countryLen == 0){
            // addrShow.value = province[current.prov].name + '-' + province[current.prov]["city"][current.city].name;
            return;
        }
        sign.signVue.city = val;
        sign.signVue.countries = [];
        sign.signVue.country = -1;
        for (var n = 0; n < countryLen; n++) {
            // var countryOpt = document.createElement('option');
            // countryOpt.innerText = province[current.prov]["city"][val].districtAndCounty[n];
            // countryOpt.value = n;
            // country.appendChild(countryOpt);
            sign.signVue.countries.push(province[current.prov]["city"][val].districtAndCounty[n])
        }
    }
}

/*选择县区之后的处理函数*/
function selectCountry(obj) {
    current.country = obj.options[obj.selectedIndex].value;
    if ((current.city != null) && (current.country != null)) {
        // btn.disabled = false;
    }
}

/*点击确定按钮显示用户所选的地址*/
function showAddr() {
    // addrShow.value = province[current.prov].name + '-' + province[current.prov]["city"][current.city].name + '-' + province[current.prov]["city"][current.city].districtAndCounty[current.country];
}

var sign = new function () {
    'use strict';
    var storage = window.localStorage;

    var signVue = new Vue({
        el: '#signArea',
        data: {
            hasError: false,
            errorMessage: "",
            provinces: [],
            cities: [],
            countries: [],
            province: -1,
            city: -1,
            country: -1,
            school: "",
            username: "",
            name: "",
            password: "",
            password2: "",
            token: "",
            tel: "",
            email: "",
            qq: "",

            contestIndexTab: storage.getItem("contestIndexTab"),
            hollandQuestions: [],
            answers: {}
        },
        methods: {
            /**
             * 点击作答
             * @param id
             * @param answer
             */
            clickAnswer: function (id, answer) {
                examVue.answers[id] = answer;
            },
            /**
             * 答案提交
             */
            doSubmit: function () {
                signVue.errorMessage = "";
                signVue.hasError = false;
                if (signVue.username === "") {
                    signVue.errorMessage = "用户名不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.name === "") {
                    signVue.errorMessage = "姓名不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.password === "") {
                    signVue.errorMessage = "密码不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.password !== signVue.password2) {
                    signVue.errorMessage = "输入的两次密码不一致！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.token === "") {
                    signVue.errorMessage = "卡密不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.tel === "") {
                    signVue.errorMessage = "电话不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.email === "") {
                    signVue.errorMessage = "邮箱不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.province === -1) {
                    signVue.errorMessage = "地区（省）不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.city === -1) {
                    signVue.errorMessage = "地区（市）不可为空！";
                    signVue.hasError = true;
                    return;
                }
                if (signVue.country === -1) {
                    signVue.errorMessage = "地区（区|县）不可为空！";
                    signVue.hasError = true;
                    return;
                }

                $.getJSON('/doSign', {
                    username: signVue.username,
                    name: signVue.name,
                    password: signVue.password,
                    token: signVue.token,
                    tel: signVue.tel,
                    email: signVue.email,
                    qq: signVue.qq,
                    province: province[signVue.province].name,
                    city: province[signVue.province]["city"][signVue.city].name,
                    country: province[signVue.province]["city"][signVue.city].districtAndCounty[signVue.country],
                    school: signVue.school
                }).done(function (response) {
                    if (response.success) {
                        console.log(response);
                        alert(response.message);
                        window.location.href = "/login"
                    } else {
                        console.log(response.message);
                        signVue.errorMessage = response.message;
                        signVue.hasError = true;
                    }
                });
            },

            /**
             * 查看代码
             * @param versionId
             */
            showSourceCode: function (strategyId, versionId, languageType) {
                // contestHomeVue.historyVersionId = versionId;
                if (arguments.length != 3) {
                    alert("缺少必要查询参数！");
                    return false;
                }
                $.getJSON('/strategy/getStrategyOldVersion', {
                    strategyId: strategyId,
                    versionId: versionId
                }).done(function (response) {
                    if (response.success) {
                        var editor = ace.edit("showSourceCode");
                        editor.setTheme("ace/theme/xcode");
                        if (languageType == 'java') {
                            editor.getSession().setMode("ace/mode/java");
                        } else {
                            editor.getSession().setMode("ace/mode/python");
                        }
                        editor.setValue(response.sourceCode, -1);
                        editor.setReadOnly(true);
                        editor.gotoLine(1);
                        popUp.showLayer($('#sourceCodeBox'));
                    } else {
                        console.log(response.message);
                        alert("操作失败，请稍后再试!");
                    }
                });
            }
        }
    });

    this.signVue = signVue;

    /**
     * 关闭参加竞赛确认窗口
     */
    this.closeSubmitBox = function () {
        popUp.hideLayer($('#submitComplate'));
    };

    /**
     * 删除策略
     */
    this.doDelete = function () {
        $.getJSON('/contest/removeContestStrategy', {
            contestStrategyId: $('#contestStrategyId').val()
        }).done(function (response) {
            if (response.success) {
                contestHomeVue.getContestList();
                $(".sel-lv-list [data-value='" + $('#delStrategyId').val() + "']").removeClass("add-s-round");
            } else {
                console.log(response.message);
                alert("操作失败，请稍后再试!");
            }
        });
        this.closeDeleteBox();
    };

    /**
     * 关闭删除窗口
     */
    this.closeDeleteBox = function () {
        popUp.hideLayer($('#deleteContestStrategyBox'));
    };

    /**
     * 参赛者职业选择
     *
     * @param me
     */
    this.selectCareerType = function (me) {
        $('#occupation').data("value", $(me).data("value"));
    };

    /**
     * 竞赛报名用户信息提交
     */
    this.submitUserInfo = function () {
        if ($("#phone").val() && $("#email").val()) {
            var telReg = /^1\d{10}$/;
            if (!telReg.test($("#phone").val())) {
                alert("手机号填写不正确！");
                return;
            }
            if ($("#referralPhone").val() && !telReg.test($("#referralPhone").val())) {
                alert("推荐人手机号填写不正确！");
                return;
            }
            var emailReg = /^(\w-*\.*)+@(\w-?)+(\.\w{2,})+$/;
            if (!emailReg.test($("#email").val())) {
                alert("邮箱格式填写不正确！");
                return;
            }
            var strategyIdArray = [];
            $(".sel-lv-list li.add-s-round").each(function () {
                strategyIdArray.push($(this).attr("data-value"));
            });
            $.getJSON('/contest/submitStrategy', {
                contestId: $('#contestId').val(),
                strategyIds: strategyIdArray.join(",")
            }).done(function (response) {
                if (response.success) {
                    $.getJSON('/contest/insertUserInfo', {
                        contestId: $('#contestId').val(),
                        phone: $('#phone').val(),
                        email: $('#email').val(),
                        occupation: $('#occupation').data("value"),
                        referralName: $('#referralName').val(),
                        referralPhone: $('#referralPhone').val()
                    }).done(function (response) {
                        if (response.success) {
                            $("#hasInsertUserInfo").val("true");
                            popUp.hideLayer($('#useInfoDiv'));
                            popUp.showLayer($('#submitComplate'));
                        } else {
                            console.log(response.message);
                            alert("信息提交失败，请稍后再试!");
                        }
                    });
                    $(".compe-title").addClass("add-bg");
                    $("#homeTab").addClass("li-cpe-add").siblings().removeClass("li-cpe-add");
                    $('.cpt-content .com-cpt').hide().eq($('#homeTab').index()).show();
                    contestIndex.contestHomeVue.getContestList();
                    contestIndex.chooseStrategyVue.hasSelected = false;
                } else {
                    alert(response.message);
                }
            });
        } else {
            alert("请输入必填项(*标示输入框)!")
        }
    };
};