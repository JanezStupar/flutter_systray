// +build linux windows

package flutter_systray

import "github.com/shurcooL/trayhost"

func updateMenu(items []trayhost.MenuItem) {
	go func() {
		trayhost.Exit()
		trayhost.UpdateMenu(items)
		trayhost.EnterLoop()
	}()
}
