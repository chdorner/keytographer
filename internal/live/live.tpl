 <!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Keytographer</title>
    <link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre.min.css">
    <link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre-exp.min.css">
    <link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre-icons.min.css">

    <style>
      .container { padding: 8px; }
      textarea#src {
        width: 100%;
        border: 0px;
        font-family: monospace;
        font-size: 1em;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <header class="navbar">
        <section class="navbar-section">
          <div class="columns">
            <div class="column">
              <a href="/" class="label label-primary text-uppercase">Keytographer</a>
            </div>
            <div class="divider-vert"></div>
            <div id="header-info" class="column">
              <div id="loading" class="loading"></div>
            </div>
          </div>
        </section>
      </header>
      <div id="renders" class="columns">
        <div class="column col-12">
          <ul id="tabs" class="tab tab-block">
          </ul>
          <div id="layouts"></div>
        </div>
      </div>
    </div>

    <template id="infotmpl">
      <div id="info">
        <span id="info" class="label label-secondary">
          <span id="info-name"></span><sub id="info-kb"></sub>
        </span>
      </div>
    </template>
    <template id="layertab">
      <li class="tab-item">
        <a href="#" class="layer-name"></a>
      </li>
    </template>
    <template id="layercontent">
      <div class="layer-content d-none">
        <div class="svg"></div>
{{ if .debug }}
        <pre class="code" data-lang="SVG"></pre>
{{ end }}
      </div>
    </template>

    <script>
        let socket = new WebSocket("{{ .websocketURL }}")

{{ if .debug }}
        socket.onclose = event => {
          setTimeout(() => {
            location.reload();
          }, 2000);
        };
{{ end }}

        socket.onmessage = event => {
          toggleElement("#loading", false);

          const msg = JSON.parse(event.data);
          renderInfo(msg.name, msg.keyboard);
          renderLayers(msg.layers);

          // const render = document.createRange().createContextualFragment(event.data);
          // document.querySelector('#live').innerHTML = '';
          // document.querySelector('#live').appendChild(render);

{{ if .debug }}
          // const src = document.querySelector('#src')
          // src.value = event.data;
          // src.style.height = `${src.scrollHeight}px`
{{ end }}
        }

        function selectLayerTab(layerId) {
          const tabs = document.querySelector("#tabs");
          const layouts = document.querySelector("#layouts");

          const active = tabs.querySelector(".active");
          console.log("active", active);
          if (active != null) {
            active.classList.remove("active");
            layouts.querySelector(`.layer-content[data-layer-id="${active.dataset.layerId}"]`).classList.add("d-none");
          }

          const tab = tabs.querySelector(`li[data-layer-id="${layerId}"`)
          tab.classList.add("active");
          layouts.querySelector(`.layer-content[data-layer-id="${layerId}"]`).classList.remove("d-none");
        }

        function renderLayers(layers) {
          const tabs = document.querySelector("#tabs");
          const layouts = document.querySelector("#layouts");

          tabs.innerHTML = '';
          layouts.innerHTML = '';

          var initialLayerId = null;
          for (layer of layers) {
            const layerId = crypto.randomUUID();
            if (initialLayerId == null) {
              initialLayerId = layerId;
            }

            const tab = layertab.content.cloneNode(true);
            tab.querySelector('.tab-item').dataset.layerId = layerId;
            tab.querySelector("a.layer-name").innerHTML = layer.name;
            tab.querySelector("a.layer-name").onclick = function() { selectLayerTab(layerId); };
            tabs.append(tab);

            const contentDoc = layercontent.content.cloneNode(true);
            const content = contentDoc.querySelector(".layer-content")
            content.dataset.layerId = layerId;
            content.querySelector(".svg").innerHTML = layer.svg;
            content.querySelector(".code").innerHTML = xmlEncode(layer.svg);

            layouts.append(content);
          }

          selectLayerTab(initialLayerId);
        }

        function renderInfo(name, keyboard) {
          const headerInfo = document.querySelector("#header-info");

          const existing = headerInfo.querySelector("#info");
          if (existing != null) {
            existing.remove();
          }

          const clone = infotmpl.content.cloneNode(true);
          clone.querySelector("#info-name").innerHTML = name;
          clone.querySelector("#info-kb").innerHTML = keyboard;
          headerInfo.append(clone);
        }

        function toggleElement(selector, show) {
          var show = (typeof show !== 'undefined') ? show : true;

          const el = document.querySelector(selector)
          if (show) {
            el.classList.remove("d-none");
          } else {
            if (!el.classList.contains("d-none")) {
              el.classList.add("d-none");
            }
          }
        }

        function xmlEncode(e) {
          return e.replace(/[\<\>\"\^]/g, function(e) {
	          return "&#" + e.charCodeAt(0) + ";";
          });
        }
    </script>
  </body>
</html>
