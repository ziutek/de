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
        var ready = false;
        var img = new Image();
        img.onload = function() {
            canvas.width = img.width;
            canvas.height = img.height;
            ready = true;
        }
        img.src = "/img";
        if ("WebSocket" in window) {
            var ws = new WebSocket("ws://$ListenOn/data");
            ws.onmessage = function(e) {
                if (ready) {
                    ctx.drawImage(img, 0, 0);
                    ctx.fillStyle = "#f00";
                    ctx.fillRect(e.data, e.data, 2, 2);
                }
            }
        } else {
            alert("Websocket not supported by your browser!");
        }
    }
</script>
</body>
</html>
