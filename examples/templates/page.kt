<!doctype html>
<html lang=en>
<head>
  <meta charset='utf-8'>
  <title>DE test</title>
</head>
<body>
<section>
    <canvas id=img>Canvas not supported by your browser!</canvas>
</section>
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
                    ctx.drawImage(img, 0, 0);
                    ctx.strokeStyle = "#f00";
                    ctx.strokeRect(e.data, e.data, 1, 1);
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
