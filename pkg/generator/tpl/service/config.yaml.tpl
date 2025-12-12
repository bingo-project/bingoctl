{{- if .EnableHTTP}}
http:
  enabled: true
  addr: :8080
{{- end}}
{{- if .EnableGRPC}}

grpc:
  enabled: true
  addr: :9090
{{- end}}
{{- if .EnableWS}}

websocket:
  enabled: true
  addr: :8081
{{- end}}

log:
  level: info
  format: console
  output-paths:
    - stdout
  error-output-paths:
    - stderr
