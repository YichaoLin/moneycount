<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.0//EN" "http://www.wapforum.org/DTD/xhtml-mobile10.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" />
<title>GBC room</title>
<style type="text/css">
<!--
#Layer1 {
	position:absolute;
	width:100%;
	height:115px;
	z-index:1;
}

#LogArea {
	position:absolute;
	width:100%;
	height:249px;
	z-index:3;
	left: 12px;
	top: 200px;
}

#CalculatorButtons {
	width:200px;
}

#numOfTrans{
	font-size:24px;
}
#btnSend{
	font-size:20px;
	height:200%
	width:100%
}

#select{
	font-size:20px;
}

#RemainArea{
	font-size:20px;
}

-->
</style>

<script src="http://cdn.bootcss.com/jquery/2.0.3/jquery.min.js"></script>
<script type="text/javascript">
    $(function() {

    var conn;
    var msg = $("#numOfTrans");
    var log = $("#LogArea");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        var tmp = log.val()
        msg.prependTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form1").submit(function() {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            //return false;
        }
    	var targetuser = $("#select option:selected");
        conn.send(msg.val()+"&"+targetuser.text());
        msg.val("0");
        return false;
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://{{.TheHost}}/ws");
        conn.onclose = function(evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = function(evt) {
        	var markstr = "[_updateremain_]"
        	if(evt.data.indexOf(markstr) == "0"){
        		var numofmoney = evt.data.substr(markstr.length)
        		var remainmoney = $("#RemainArea");
        		remainmoney.html(numofmoney)
        	}
        	else{
            	appendLog($("<div/>").text(evt.data))
            }
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
    });
</script>

</head>

<body>
<div id="Layer1">
  <table width="90%" border="0" align="center">
  	<tr>
  	  <td width = "40%" colspan="2" align="center">City rise</td>
  	  <td align="center"><a href="/leaveroom">leave room</a></td>
    </tr>
    <tr align="center">
      <td>{{.PlayerName}}'s&nbsp;Balance:</td>
  	  <td name="RemainArea" id="RemainArea" align="center">UNKNOWN</td>
    </tr>
  </table>
  <form id="form1">
    <table width="90%" border="1" align="center">
      <tr>
        <td colspan="2" align="right" width="50px">
          <input name="numOfTrans" id="numOfTrans" type="text" style="width:170px" maxlength="6" />
        </td>
        <td rowspan="2" align="center">
          <input name="btnSend" id="btnSend" type="submit" value="Send"/>
        </td>
      </tr>
      <tr>
        <td width="40px">
          to:
        </td>
        <td align="right" width="50px">
          <select name="select" name="select" id="select" style="width:130px">
          <option value="0">BANK</option>
          {{ range $key, $value := .PlayerList}}
          	<option value="{{$value.PlayerId}}">{{$key}}</option>
  		  {{ end }}
          </select>
        </td>
      </tr>
    </table>
  </form>
</div>
<div name="LogArea" id="LogArea"></div>
</body>
</html>
