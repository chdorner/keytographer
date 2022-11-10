{{ .markStart }}
const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {
{{- range $lidx, $layer := .layers -}}
    {{- if $lidx }},{{ end }}
    [{{ $lidx }}] = {{ or $.layoutMacro "LAYOUT" }}(
        {{- range $kidx, $key := $layer.Keys -}}
        {{- if $kidx }}, {{ end -}}
        {{ or $key.Code "KC_NO" }}
        {{- end -}}
    )
{{- end }}
};
{{ .markStop }}
