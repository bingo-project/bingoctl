server:
{{- if .EnableHTTP}}
  addr: :8080
  mode: release
{{- end}}
{{- if .EnableGRPC}}

grpc-server:
  addr: :9090
{{- end}}

log:
  level: info
  format: console
  output-paths:
    - stdout
  error-output-paths:
    - stderr