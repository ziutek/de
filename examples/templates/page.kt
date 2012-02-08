<!doctype html>
<html lang=en>
	<head>
		<meta charset='utf-8'>
		<title>DE test</title>
	</head>
	<body>
		<table>
		<canvas id=img>Canvas not supported by your browser!</canvas><br>
		Status: <span id=status>working...</span>
		<script>
			if ("HTMLCanvasElement" in window) {
				var canvas = document.getElementById("img");
				var ctx = canvas.getContext("2d");
				var img = new Image();
				img.onload = function() {
					canvas.width = img.width;
					canvas.height = img.height;
					if ("WebSocket" in window) {
						var ws = new WebSocket("ws://$ListenOn/data");
						ws.onmessage = function(e) {
							var points = JSON.parse(e.data);
							ctx.drawImage(img, 0, 0);
							for (var i in points) {
								var p = points[i];
								ctx.strokeStyle = "#f00";
								ctx.strokeRect(p[0], p[1], 1, 1);
							}
						}
						ws.onclose = function(e) {
							var s = document.getElementById("status");
							s.firstChild.nodeValue = "completed.";
						}
						} else {
						alert("Websocket not supported by your browser!");
					}
				}
				img.src = "/img?i=" + new Date().getTime();
			}
		</script>
	</body>
</html>
