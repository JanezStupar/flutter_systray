// +build darwin

package flutter_systray

func updateMenu(items []trayhost.MenuItem) {
	trayhost.UpdateMenu(items)
}
