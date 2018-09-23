<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=gb2312" />
<title>BGC Lobby</title>
<style type="text/css">
<!--
#Layer1 {
	position:absolute;
	left: 10%;
	right: 10%;
	height:50px;
	z-index:1;
}
#Layer2 {
	position:absolute;
	top: 120px;
	bottom: 0.5em;
	left: 10%;
	right: 10%;
	z-index:2;
}
#txtRoomNum {
	font-size:24px;
}
#btnJoin {
	font-size:24px;
}
#btnCreateRoom {
	font-size:24px;
}
#formJoin {
	font-size:24px;
}
#formCreate {
	font-size:24px;
}
#tblRooms {
	font-size:16px
}
#biggersize{
	font-size:18px;
}

-->
</style>
</head>

<body>
<div id="Layer1">
  <table width="100%" border="3" align="center" id="tblForms">
    <tr>
      <td>
	  <form id="formJoin" name="formJoin" method="post" action="/enter"  align="center">
	    Room number
	    <input name="txtRoomNum" type="text" id="txtRoomNum" size="8" maxlength="4" />
	    <input name="btnJoin" type="submit" id="btnJoin" value="Join" />
	  </form>
      </td>
      
      <td>
	  <form id="formCreate" name="formCreate" method="post" action="/create"  align="center">
		<input name="btnCreateRoom" type="submit" id="btnCreateRoom" value="Create room" />
		<a href="{{$}}/logout">logout</a>
	  </form>
      </td>
    </tr>
  </table>
  
</div>

<div id="Layer2">
  <table width="100%" border="3" align="center" id="tblRooms">
  	{{ range $key, $value := .Roomlist}}
  		<tr>
  		<td width="31%">{{$key}}</td>
  		<td width="69%">
  		{{ range $playerinroom, $mm := $value}}
  			{{$playerinroom}};
  		{{ end }}
  		</td>
  		</tr>
  	{{ end }}
  </table>
</div>
</body>
</html>
