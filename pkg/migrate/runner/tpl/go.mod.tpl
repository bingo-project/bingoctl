module bingoctl-migrator

go {{.GoVersion}}

require (
	{{.UserModule}} v0.0.0
	github.com/bingo-project/bingoctl {{.BingoctlVersion}}
)

replace {{.UserModule}} => {{.UserModulePath}}
