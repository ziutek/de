<!doctype html>
<html lang=en>
	<head>
		<meta charset='utf-8'>
		<title>Differential evolution - 4d problem</title>
	</head>
	<body>
		<canvas id=img>Canvas not supported by your browser!</canvas><br>
		Status: <span id=status>working...</span>
		Iteration: <span id=iter>0</span>
		<script>
			var iter = 0;
			if ("HTMLCanvasElement" in window) {
				var canvas = document.getElementById("img");
				var ctx = canvas.getContext("2d");
				var img = new Image();
				var i = 0;
				var iter = document.getElementById("iter").firstChild;
				var stat = document.getElementById("status").firstChild;
				img.onload = function() {
					canvas.width = img.width;
					canvas.height = img.height;
					if ("WebSocket" in window) {
						var ws = new WebSocket("ws://$ListenOn/data");
						ws.onmessage = function(e) {
							var points = JSON.parse(e.data);
							ctx.drawImage(img, 0, 0);
							ctx.strokeStyle = "#f00";
							for (var k in points) {
								var p = points[k];
								ctx.strokeRect(p[0], p[1], 1, 1);
							}
							iter.nodeValue = ++i;
						}
						ws.onclose = function(e) {
							stat.nodeValue = "completed.";
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
