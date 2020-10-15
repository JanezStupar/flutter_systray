// +build darwin

package flutter_systray

import "github.com/shurcooL/trayhost"
import "fmt"

func (p *FlutterSystrayPlugin) updateMenu(actions []SystrayAction) {
	items, err := p.actionsToMenu(actions)
	if err != nil {
		fmt.Println("An error has occurred while registering actions", err)
	}

	trayhost.UpdateMenu(items)
}

func initialize(title string, iconData []byte) {
	trayhost.Initialize(title, iconData, nil)
}
