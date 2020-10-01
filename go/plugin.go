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
	window  *glfw.Window
	channel *plugin.MethodChannel
}

type MainEntry struct {
	title    string
	tooltip  string
	iconPath string
}

type ActionEnumType int

type actionType struct {
	Quit         ActionEnumType
	Focus        ActionEnumType
	SystrayEvent ActionEnumType
}

var ActionType = &actionType{
	Quit:         0,
	Focus:        1,
	SystrayEvent: 2,
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
	p.channel = plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	p.channel.HandleFunc("initSystray", p.initSystrayHandler)
	p.channel.HandleFunc("addActions", p.addActionsHandler)
	return nil
}

func (p *FlutterSystrayPlugin) initSystrayHandler(arguments interface{}) (reply interface{}, err error) {
	// Convert the params into SystrayAction type list
	argsMap := arguments.(map[interface{}]interface{})
	var mainEntry MainEntry
	if argsMap["mainEntry"] != nil {
		mainEntry, err = parseMainEntry(argsMap["mainEntry"])
		if err != nil {
			fmt.Println("an error has occurred while parsing main entry parameters", err)
		}
		delete(argsMap, "mainEntry")
	}

	actions, err := parseActionParams(argsMap)
	if err != nil {
		fmt.Println("an error has occurred while parsing action parameters", err)
	}

	var readyFunc = func() {
		if len(mainEntry.iconPath) > 0 {
			var data []byte
			data, err := parseIcon(mainEntry.iconPath)
			if err != nil {
				fmt.Println(fmt.Sprintf("An error has occurred while parsing the icon: %s", err))
			}

			if data != nil {
				systray.SetIcon(data)
			}
		}

		if len(mainEntry.title) > 0 {
			println(mainEntry.title)
			systray.SetTitle(mainEntry.title)
		}

		if len(mainEntry.tooltip) > 0 {
			systray.SetTooltip(mainEntry.tooltip)
		}

		err = p.addActions(actions)
		if err != nil {
			fmt.Println(fmt.Sprintf("an error has occurred while registering actions: %s", err))
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
		fmt.Println("an error has occurred while parsing action parameters", err)
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
		} else if action.actionType == ActionType.SystrayEvent {
			mEvent := systray.AddMenuItem(action.label, action.tooltip)
			addIcon(action.iconPath, mEvent)
			// Set up a callback handler
			go func(reference SystrayAction) {
				for {
					<-mEvent.ClickedCh
					err := p.invokeSystrayEvent(reference)
					if err != nil {
						fmt.Println(fmt.Sprintf("An error has occurred while invoking SystrayEvent: %s", err))
					}
				}
			}(action)
		}
	}

	return nil
}

func (p *FlutterSystrayPlugin) invokeSystrayEvent(action SystrayAction) error {
	err := p.channel.InvokeMethod("systrayEvent", action.name)
	if err != nil {
		return err
	}

	return nil
}

func parseMainEntry(entry interface{}) (MainEntry, error) {
	m := entry.(map[interface{}]interface{})
	parsed := MainEntry{
		title:    m["title"].(string),
		tooltip:  m["tooltip"].(string),
		iconPath: m["iconPath"].(string),
	}
	return parsed, nil
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
			fmt.Println(fmt.Sprintf("An error has occurred while parsing the icon: %s", err))
		}
		item.SetIcon(data)
	}
}

/*
*	This function performs cleanup of systray menu
 */
func systrayOnExit() {
}
