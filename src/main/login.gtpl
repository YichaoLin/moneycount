<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=gb2312" />
<title>BGC login</title>
<style type="text/css">
<!--
#body {
	width:100%;
	height:100%;
	z-index:1;
	background: gray;
}
#Layer1 {
	position:absolute;
	top: 25%;
	left: 25%;
	right: 25%;
	bottom: 25%;
	z-index:1;
	height: 189px;
}
#username {
	font-size:24px;
}
#password {
	font-size:24px;
}
#btnsubmit {
	font-size:24px;
}
#Layer2 {
	position:absolute;
	width:100%;
	z-index:2;
}
#ErrorTips {
	font-size: 18px
}
-->
</style>
</head>

<body>
<div id="Layer1">
  <form id="form1" name="form1" method="post" action="/login">
    <p align="center">
      <label>Username
        <input name="tfname" type="text" id="username" maxlength="10" />
      </label>
    </p>
    <p align="center">
      <label>Password
      	<input name="tfpswd" type="password" id="password" />
      </label>
	</p>
    <div id="Layer2">
		<div id="ErrorTips" style="width:100%" align="center">{{ .Tips}}</div>
		<p align="center">
		  <input name="Submit" type="submit" id="btnsubmit" value="Join" />
		</p>
    </div>
  </form>
</div>
</body>
</html>
