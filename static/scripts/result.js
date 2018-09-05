/**
 * Exam scripts
 *
 * @author Zhiguo.Chen
 */
$().ready(function () {
    // (function(){
    //
    // })();

    // $('ul').on("click", "li.option", function () {
    //     debugger;
    //     var examId = $(this).closest('.test_content_nr_main').closest('li').attr('id'); // 得到元素ID
    //     var cardLi = $('a[href=#' + examId + ']'); // 根据题目ID找到对应答题卡
    //     // 设置已答题
    //     if (!cardLi.hasClass('hasBeenAnswer')) {
    //         cardLi.addClass('hasBeenAnswer');
    //     }
    // });
    var data = {domId: "hollandPicture"};
    $.getJSON("/exam/getHollandAnswerResult", function (response) {
        if (response.success) {
            result.resultVue.hollandResults = response.mapData.results;
            data.value1 = response.mapData.scores;
            data.value2 = [17,33,21,13,41,45];
            drawPolar(data);
        } else {
            alert(response.message);
        }
    });

});

var result = new function () {
    'use strict';
    var storage = window.localStorage;

    var resultVueComp = new Vue({
        el: '#resultArea',
        data: {
            contestIndexTab: storage.getItem("contestIndexTab"),
            hollandResults: [],
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
                var answerString = "";
                var answerCount = 0;
                for (var objKey in examVue.answers) {
                    answerCount++;
                    if (answerString === "") {
                        answerString = objKey + "|" + examVue.answers[objKey];
                    } else {
                        answerString = answerString + "," + objKey + "|" + examVue.answers[objKey];
                    }
                }
                if (answerCount < examVue.hollandQuestions.length) {
                    alert("还有未作答的题目哦，具体请看右侧答题卡！");
                    return;
                }
                $.getJSON('/exam/doSubmit', {
                    answerStr: answerString
                }).done(function (response) {
                    if (response.success) {
                        console.log(response);
                        alert(response.message);
                        window.location.href = "/exam/result"
                    } else {
                        console.log(response.message);
                        alert(response.message);
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

    this.resultVue = resultVueComp;

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
