<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Keytographer</title>
    <style>
      textarea#src {
        width: 100%;
        border: 0px;
        font-family: monospace;
        font-size: 1em;
      }
    </style>
  </head>
  <body>
    <div id="live" style="display: block"></div>
{{ if .debug }}
    <textarea id="src" disabled></textarea>
{{ end }}
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

{{ if .debug }}
          const src = document.querySelector('#src')
          src.value = event.data;
          src.style.height = `${src.scrollHeight}px`
{{ end }}
        }
    </script>
  </body>
</html>
