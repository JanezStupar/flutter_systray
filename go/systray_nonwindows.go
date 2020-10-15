// +build linux darwin !windows

package flutter_systray

import "github.com/shurcooL/trayhost"

func (p *FlutterSystrayPlugin) actionsToMenu(actions []SystrayAction) ([]trayhost.MenuItem, error) {
	var items []trayhost.MenuItem

	for _, action := range actions {
		localAction := action
		if localAction.actionType == ActionType.Focus {
			// Adds a GLFW `window.Show` (https://godoc.org/github.com/go-gl/glfw/v3.2/glfw#Window.Show) operation to the
			// systray menu. It is used to bring window to front.
			mShow := trayhost.MenuItem{
				Title:   localAction.label,
				Enabled: nil,
				Handler: p.focusHandler(&localAction),
			}
			items = append(items, mShow)
		} else if localAction.actionType == ActionType.Quit {
			// Set up a handler to close the window
			mQuit := trayhost.MenuItem{
				Title:   localAction.label,
				Enabled: nil,
				Handler: p.closeHandler(&localAction),
			}
			items = append(items, mQuit)
		} else if localAction.actionType == ActionType.SystrayEvent {
			mEvent := trayhost.MenuItem{
				Title:   localAction.label,
				Enabled: nil,
				Handler: p.eventHandler(&localAction),
			}
			items = append(items, mEvent)
		}
	}

	return items, nil
}
