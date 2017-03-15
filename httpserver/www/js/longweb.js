var Longweb = {
    _FailTimes: 0,
    _run: 1,

    ConnType:"",
    WsUrl:"",
    LongpollUrl:"",

    ConnType_Websocket:"websocket",
    ConnType_Longpoll:"longpoll",

    //主动停止longweb运行
    Stop: function () {
        Longweb._run = 0;
    },

    //初始化websocket连接
    initwebsocket: function (appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, Token) {
        if (Longweb._run != 1) { return; }

        var wsc; 
        try {
            if (typeof (Token) != "undefined") {
                wsc = new WebSocket(Longweb.WsUrl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&token=' + Token);
            } else {
                wsc = new WebSocket(Longweb.WsUrl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid);
            }
        } catch (e) {
            console.log("WebSocket Initialize Error");
        }
        
        wsc.onopen = OnOpen;
        wsc.onmessage = OnMessage;
        wsc.onerror = OnError;
        wsc.onclose = OnClose;
    },

    //初始化longpoll连接
    initlongpoll: function (appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, Token) {
        
        var pollingurl = Longweb.LongpollUrl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&querykey=' + escape(querykey) + '&_r='+Math.random();

        if (typeof (Token) != "undefined") {
            pollingurl = Longweb.LongpollUrl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&querykey=' + escape(querykey) + '&token=' + Token + '&_r=' + Math.random();
        }

        console.log(pollingurl);
        var ev = {};
        $.ajax({
            type: "GET",
            cache: false,
            url: pollingurl,
            dataType: "jsonp",
            jsonp: "jsonpcallback",
            success: function (data, textStatus, jqXHR) {
                if (data != null && data != "" && data != "null") {
                    if (data.RetCode == "0") {
                        eval('var obj=' + data.Message + ';');
                        if (typeof(obj)!="undefined") {
                            ev.data = obj;
                            OnMessage(ev);
                        }
                    } else {
                        OnError(data.RetCode);
                    }
                }
            },
            complete: function (XMLHttpRequest, textStatus) {
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                if (textStatus == "timeout") {
                    OnError("timeout");
                }
                
            }

        });

    },

    //连接
    Connect:function (appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, Token) {
        if (Longweb._run != 1) { return; }
        if (Longweb.ConnType == Longweb.ConnType_Websocket) {//支持WebSocket，启用WebSocket通讯
            Longweb.initwebsocket(Appid, Groupid, ClientId, OnMessage, OnClose, OnError, OnOpen, QueryKey, Token);
        } else {
            //启用LongPolling长轮询
            Longweb.initlongpoll(Appid, Groupid, ClientId, OnMessage, OnClose, OnError, OnOpen, QueryKey, Token);
        }
    },
};

function NewGuid() {
    var guid = "";
    for (var i = 1; i <= 32; i++) {
        var n = Math.floor(Math.random() * 16.0).toString(16);
        guid += n;
        if ((i == 8) || (i == 12) || (i == 16) || (i == 20))
            guid += "-";
    }
    return guid;
}

var browser = {
    versions: function () {
        var u = navigator.userAgent, app = navigator.appVersion;
        return {//移动终端浏览器版本信息
            trident: u.indexOf('Trident') > -1, //IE内核
            presto: u.indexOf('Presto') > -1, //opera内核
            webKit: u.indexOf('AppleWebKit') > -1, //苹果、谷歌内核
            gecko: u.indexOf('Gecko') > -1 && u.indexOf('KHTML') == -1, //火狐内核
            mobile: !!u.match(/AppleWebKit.*Mobile.*/) || !!u.match(/AppleWebKit/), //是否为移动终端
            ios: !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/), //ios终端
            android: u.indexOf('Android') > -1 || u.indexOf('Linux') > -1, //android终端或者uc浏览器
            iPhone: u.indexOf('iPhone') > -1 || u.indexOf('Mac') > -1, //是否为iPhone或者QQHD浏览器
            iPad: u.indexOf('iPad') > -1, //是否iPad
            webApp: u.indexOf('Safari') == -1 //是否web应该程序，没有头部与底部
        };
    }(),
    language: (navigator.browserLanguage || navigator.language).toLowerCase()
}