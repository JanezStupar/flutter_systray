// +build linux

package flutter_systray

import (
	"fmt"
	"github.com/shurcooL/trayhost"
)

func (p *FlutterSystrayPlugin) updateMenu(actions []SystrayAction) {
	items, err := p.actionsToMenu(actions)
	if err != nil {
		fmt.Println("An error has occurred while registering actions", err)
	}

	go func() {
		trayhost.Exit()
		trayhost.UpdateMenu(items)
		trayhost.EnterLoop()
	}()
}

func initialize(title string, iconData []byte) {
	trayhost.Initialize(title, iconData, nil)
}
