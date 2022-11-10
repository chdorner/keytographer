{{ .markStart }}
const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {
{{- range $lidx, $layer := .config.Layers -}}
    {{- if $lidx }},{{ end }}
    [{{ $lidx }}] = LAYOUT_split_3x5_3(
        {{- range $kidx, $key := $layer.Keys -}}
        {{- if $kidx }}, {{ end -}}
        {{ or $key.Code "KC_NO" }}
        {{- end -}}
    )
{{- end }}
};
{{ .markStop }}
