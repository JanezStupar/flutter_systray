// +build darwin

package flutter_systray

import "github.com/shurcooL/trayhost"

func updateMenu(items []trayhost.MenuItem) {
	trayhost.UpdateMenu(items)
}
