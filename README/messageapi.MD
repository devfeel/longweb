#longweb实时通信框架接入文档：

**************************************************************************************************

##sendmessage:
api: http://yourwwwurl/sendmessage<br>
message:<br>
```
{
    "AppID":"10000",
    "GroupID":"1",
    "UserID":"",
    "MessageLevel":"2",//表示仅发送给经过鉴权的连接
    "ToAppID":"10000",
    "ToGroupID":"1",
    "ToUserID":"",
    "Content":"test from fiddler"
}
```

属性说明：<br>
MessageLevel:

* MessageLevel_All    = "0"  //"all"
* MessageLevel_Normal = "1"  //"normal"
* MessageLevel_Auth   = "2"  //"auth"

result:
```
{"RetCode":0,"RetMsg":"ok"}
```

RetCode说明：

* 0:ok
* -1:not allow ip
* -2:message is empty
* -10001:message format error
* -10002:this appid no have permission

**************************************************************************************************

##websocket:<br>
api-url: ws://yourwwwurl/ws/onsocket?appid=?&groupid=?&userid=?&token=?
<br>参数：

* appid：应用编号，统一申请
* groupid：用户组编号，应用自定义
* userid：用户编号，应用需保证userid在同Appid下的唯一性
* token：非必须，鉴权token，不传该参数或传入为空，则表示当前请求无需鉴权；当传入该参数时表示连接需要鉴权，需要应用提供鉴权API，接口定义参考checktoken接口描述

**************************************************************************************************

##longpoll:<br>
api-url: http://yourwwwurl/poll/onpolling?appid=?&groupid=?&userid=?&querykey=?&jsonpcallback=?&token=?
<br>参数：

* jsonpcallback: 应用定义jsonp-callback函数名，如果不传入默认为callback
* appid：应用编号，统一申请
* groupid：用户组编号，应用自定义
* userid：用户编号，应用需保证userid在同Appid下的唯一性
* querykey：透传key，会透传给应用messageapi，一般用于决定是否有需要马上返回的数据
* token：非必须，鉴权token，不传该参数或传入为空，则表示当前请求无需鉴权，当传入该参数时表示连接需要鉴权，需要应用提供鉴权API，接口定义参考checktoken接口描述

result:
```
type ResponseJson struct {
  RetCode int
  RetMsg  string
  Message string
}
```
RetCode描述：<br>
返回值：

* 0 成功
* -100001：appid、groupid、querykey不能为空
* -100002：指定appid不存在
* -100003：注册失败
* -100009：超时或其他异常，一般需重新发起请求
* -200009：应用返回异常
* -101001：no permission connect! => appid[" + appId + "] no exists
* -101002：no permission connect! => check token has an error =>
* -101003：no permission connect! => check token has an error =>
* -101004：no permission connect! => check token result =>
* -101005：no permission connect! => check token: appid|groupid|userid not match


**************************************************************************************************
**************************************************************************************************
**************************************************************************************************
需要业务方实现的api：<br>
##querymessage接口<br>
需要应用实现的查询消息接口<br>
url：http://XXXXXXXX/api/querymessage?appid=?&groupid=?&userid=?&querykey=?
<br>参数：

* appid：应用编号，统一申请
* groupid：用户组编号，应用自定义
* userid：用户编号，应用需保证userid在同Appid下的唯一性
* querykey：透传key，会透传给应用messageapi，一般用于决定是否有需要马上返回的数据

返回：<br>
为空，则表示无需马上返回的数据，请求会继续等待；<br>
不为空，请求马上返回，向客户端推送业务消息<br>

##checktoken接口
需要应用实现的鉴权api<br>
<br>url：http://XXXXXXXX/api/checktoken?appid=?&groupid=?&userid=?&token=?
参数：

* appid：应用编号，统一申请
* groupid：用户组编号，应用自定义
* userid：用户编号，应用需保证userid在同Appid下的唯一性
* token：鉴权token，用于和应用鉴权Api做鉴权
返回：
```
type AuthResponse struct {
  RetCode int
  RetMsg  string
  AppID   string
  GroupID string
  UserID  string
}
```
其中<br>
1）RetCode为0表示check成功<br>
2）AppID、GroupID、UserID在应用创建token时保存，会用于验证websocket入参合法性
