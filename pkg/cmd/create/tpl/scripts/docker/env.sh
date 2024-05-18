# App info
export app_name={[.AppName]}
export version=$(git describe --tags --match='v*' | sed 's/^v//')

# Build
export registry_prefix={[.AppName]}
export images=({[.AppName]}-apiserver {[.AppName]}-watcher {[.AppName]}-bot {[.AppName]}ctl)
export architecture=amd64
