package all

import (
	"{[.RootPackage]}/internal/watcher/watcher"
	"{[.RootPackage]}/internal/watcher/watcher/user"
)

func init() {
	watcher.Register("user", &user.UserWatcher{})
}
