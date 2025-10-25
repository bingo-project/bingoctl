# App info
export app_name={[.AppName]}
export version=$(git describe --tags --match='v*' | sed 's/^v//' || echo '0.0.0')

# Build
export registry_prefix={[.AppName]}
export images=({[.AppName]}-apiserver {[.AppName]}-admserver {[.AppName]}-scheduler {[.AppName]}-bot {[.AppName]}ctl)
export architecture=amd64
export build_from="bin" # image/bin
