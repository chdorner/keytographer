<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Go WebSocket Tutorial</title>
  </head>
  <body>
    <div id="live"></div>
    <script>
        let socket = new WebSocket("{{ .websocketURL }}")

{{ if .debug }}      
        socket.onclose = event => {
          setTimeout(() => {
            location.reload();
          }, 3000);
        };
{{ end }}

        socket.onmessage = event => {
          const render = document.createRange().createContextualFragment(event.data);
          document.querySelector('#live').innerHTML = '';
          document.querySelector('#live').appendChild(render);
        }
    </script>
  </body>
</html> 