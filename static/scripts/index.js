$(function() {
	$(".outer").mouseover(function() {
		$(".outer .left").css("display","block");
		$(".outer .right").css("display","block");
	})
	$(".outer").mouseout(function() {
		$(".outer .left").css("display","none");
		$(".outer .right").css("display","none");
	})
})
$(function (){
	var imgW=$('.inner img').width();
	var x=1;
	var anm=true;
	var timer1=null;
	
	$('.M').scrollLeft(imgW);
	var fir=$('.inner img:first').clone();
	var las=$('.inner img:last').clone();
	$('.inner').append(fir);
	$('.inner').prepend(las);

	$('.left').click(function (){
		clearInterval(timer1);
		if (anm) {
			anm=false;
			x--;
			if (x<0) {
				x=$('.inner img').length-1;
				$('.M').scrollLeft(imgW*(x+1));
			};
			move();
		};
		autoMove();
	});

	$('.right').click(function (){
		clearInterval(timer1);
		if (anm) {
			anm=false;
			x++;
			if (x>=$('.inner img').length) {
				x=2;
				$('.M').scrollLeft(imgW);
			};
			move();
		};
		autoMove();
	});

	function autoMove(){
		timer1=setInterval(function (){
			x++;
			if (x>=$('.inner img').length) {
				x=2;
				$('.M').scrollLeft(imgW);
			};
			move();
		},3500)
	}
	autoMove();

	function move(){
		$('.M').stop().animate({scrollLeft:imgW*x},function (){
			anm=true;
		});
	}
})

function SetHome(obj, vrl) {
        try {
            obj.style.behavior = 'url(#default#homepage)';
            obj.setHomePage(vrl);
        }
        catch (e) {
            if (window.netscape) {
                try {
                    netscape.security.PrivilegeManager.enablePrivilege("UniversalXPConnect");
                }
                catch (e) {
                    alert("此操作被浏览器拒绝！\n请在浏览器地址栏输入“about:config”并回车\n然后将 [signed.applets.codebase_principal_support]的值设置为'true',双击即可。");
                }
                var prefs = Components.classes['@mozilla.org/preferences-service;1'].getService(Components.interfaces.nsIPrefBranch);
                prefs.setCharPref('browser.startup.homepage', vrl);
            } else {
                alert("您的浏览器不支持，请按照下面步骤操作：1.打开浏览器设置。2.点击设置网页。3.输入：" + vrl + "点击确定。");
            }
        }
    }
    // 加入收藏 兼容360和IE6
    function shoucang(sTitle, sURL) {
        try {
            window.external.addFavorite(sURL, sTitle);
        }
        catch (e) {
            try {
                window.sidebar.addPanel(sTitle, sURL, "");
            }
            catch (e) {
                alert("加入收藏失败，请使用Ctrl+D进行添加");
            }
        }
    }