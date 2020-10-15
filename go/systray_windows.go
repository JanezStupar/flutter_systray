// +build windows !linux !darwin

package flutter_systray

import (
	"github.com/getlantern/systray"
)

func (p *FlutterSystrayPlugin) updateMenu(actions []SystrayAction) {
	for _, action := range actions {
		// Adds a GLFW `window.Show` (https://godoc.org/github.com/go-gl/glfw/v3.2/glfw#Window.Show) operation to the
		// systray menu. It is used to bring window to front.
		if action.actionType == ActionType.Focus {
			mShow := systray.AddMenuItem(action.label, "")
			go func(reference SystrayAction) {
				for {
					<-mShow.ClickedCh
					p.focusHandler(&reference)()
				}
			}(action)
		} else if action.actionType == ActionType.Quit {
			// Set up a handler to close the window
			mQuit := systray.AddMenuItem(action.label, "")
			go func(reference SystrayAction) {
				<-mQuit.ClickedCh
				p.closeHandler(&reference)()
			}(action)
		} else if action.actionType == ActionType.SystrayEvent {
			mEvent := systray.AddMenuItem(action.label, "")
			// Set up a callback handler
			go func(reference SystrayAction) {
				for {
					<-mEvent.ClickedCh
					p.eventHandler(&reference)()
				}
			}(action)
		}
	}
}

func initialize(title string, iconData []byte) {
	readyFunc := func() {
		if iconData != nil {
			systray.SetIcon(iconData)
		}

		if len(title) > 0 {
			systray.SetTitle(title)
		}
	}

	go func() {
		systray.Run(readyFunc, nil)
	}()
}
