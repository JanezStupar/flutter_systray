package flutter_systray

import (
	"errors"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"strconv"
)

const channelName = "plugins.flutter.io/flutter_systray"

// FlutterSystrayPlugin implements flutter.Plugin and handles method.
type FlutterSystrayPlugin struct {
	window *glfw.Window
}

type ActionEnumType int

type actionType struct {
	Quit     ActionEnumType
	Focus    ActionEnumType
	Callback ActionEnumType
}

var ActionType = &actionType{
	Quit:     0,
	Focus:    1,
	Callback: 2,
}

type SystrayAction struct {
	name       string
	label      string
	tooltip    string
	iconPath   string
	actionType ActionEnumType
}

var _ flutter.Plugin = &FlutterSystrayPlugin{} // compile-time type check

// InitPluginGLFW initializes the GLFW
func (p *FlutterSystrayPlugin) InitPluginGLFW(window *glfw.Window) error {
	p.window = window
	return nil
}

// InitPlugin initializes the plugin.
func (p *FlutterSystrayPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("initSystray", p.initSystrayHandler)
	channel.HandleFunc("addActions", p.addActionsHandler)
	return nil
}

func (p *FlutterSystrayPlugin) initSystrayHandler(arguments interface{}) (reply interface{}, err error) {
	// Convert the params into SystrayAction type list
	argsMap := arguments.(map[interface{}]interface{})
	var mainIcon string
	if argsMap["mainIcon"] != nil {
		mainIcon = argsMap["mainIcon"].(map[interface{}]interface{})["iconPath"].(string)
		delete(argsMap, "mainIcon")
	}

	actions, err := parseActionParams(argsMap)
	if err != nil {
		fmt.Println("An error has occurred while parsing action parameters", err)
	}

	var readyFunc = func() {
		var data []byte
		data, err := parseIcon(mainIcon)
		if err != nil {
			fmt.Println("An error has occurred while parsing the icon: %S", err)
		}

		if data != nil {
			systray.SetIcon(data)
		}

		err = p.addActions(actions)
		if err != nil {
			fmt.Println("An error has occurred while registering actions", err)
		}
	}

	go func() {
		systray.Run(readyFunc, systrayOnExit)
	}()

	return "ok", nil
}

func (p *FlutterSystrayPlugin) addActionsHandler(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})

	actions, err := parseActionParams(argsMap)
	if err != nil {
		fmt.Println("An error has occurred while parsing action parameters", err)
	}

	err = p.addActions(actions)
	if err != nil {
		fmt.Println("An error has occurred while registering actions", err)
	}

	return "ok", nil
}

func (p *FlutterSystrayPlugin) addActions(actions []SystrayAction) error {
	for _, action := range actions {
		// Adds a GLFW `window.Show` (https://godoc.org/github.com/go-gl/glfw/v3.2/glfw#Window.Show) operation to the
		// systray menu. It is used to bring window to front.
		if action.actionType == ActionType.Focus {
			mShow := systray.AddMenuItem(action.label, action.tooltip)
			addIcon(action.iconPath, mShow)
			go func() {
				for {
					<-mShow.ClickedCh
					p.window.Show()
				}
			}()
		} else if action.actionType == ActionType.Quit {
			// Set up a handler to close the window
			mQuit := systray.AddMenuItem(action.label, action.tooltip)
			addIcon(action.iconPath, mQuit)
			go func() {
				<-mQuit.ClickedCh
				p.window.SetShouldClose(true)
			}()
		} else if action.actionType == ActionType.Callback {
			// Set up a callback handler
		}
	}

	return nil
}

func parseActionParams(argsMap map[interface{}]interface{}) ([]SystrayAction, error) {
	var actions []SystrayAction
	for _, v := range argsMap {
		valsMap := v.(map[interface{}]interface{})

		number, _ := strconv.Atoi(valsMap["actionType"].(string))
		action := SystrayAction{
			name:       valsMap["name"].(string),
			label:      valsMap["label"].(string),
			tooltip:    valsMap["tooltip"].(string),
			iconPath:   valsMap["iconPath"].(string),
			actionType: ActionEnumType(number),
		}
		actions = append(actions, action)
	}

	return actions, nil
}

func parseIcon(absPath string) ([]byte, error) {
	// Parse the icon if available
	var data []byte
	if len(absPath) > 0 {
		data, err := ioutil.ReadFile(absPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("file reading error: %s", err))
		}
		return data, nil
	}
	return data, nil
}

func addIcon(iconPath string, item *systray.MenuItem) {
	if len(iconPath) > 0 {
		data, err := parseIcon(iconPath)
		if err != nil {
			fmt.Println("An error has occurred while parsing the icon: %S", err)
		}
		item.SetIcon(data)
	}
}

/*
*	This function performs cleanup of systray menu
 */
func systrayOnExit() {
}
