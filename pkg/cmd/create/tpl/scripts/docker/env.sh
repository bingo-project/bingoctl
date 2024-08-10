# App info
export app_name={[.AppName]}
export version=$(git describe --tags --match='v*' | sed 's/^v//' || echo '0.0.0')

# Build
export registry_prefix={[.AppName]}
export images=({[.AppName]}-apiserver {[.AppName]}-watcher {[.AppName]}-bot {[.AppName]}ctl)
export architecture=amd64
