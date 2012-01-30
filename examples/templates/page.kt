<!doctype html>
<html lang=en>
<head>
  <meta charset='utf-8'>
  <title>DE test</title>
</head>
<body>
<section>
  Received from server: <span id=data>?</span>
</section>
<script>
    if ("WebSocket" in window) {
        var ws = new WebSocket("ws://$HostPort/data");
        var data = document.getElementById("data").firstChild
        ws.onmessage = function(e) { data.nodeValue = e.data; }
    } else {
        alert("Websocket NOT supported by your browser!");
    }
</script>
</body>
</html>
