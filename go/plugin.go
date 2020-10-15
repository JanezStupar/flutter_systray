package flutter_systray

import (
	"encoding/json"
	"errors"
	"fmt"
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
	p.channel.HandleFuncSync("initSystray", p.initSystrayHandler)
	p.channel.HandleFunc("updateMenu", p.updateMenuHandler)
	return nil
}

func (p *FlutterSystrayPlugin) initSystrayHandler(arguments interface{}) (reply interface{}, err error) {
	// Convert the params into SystrayAction type list
	var argsMap map[string]interface{}
	err = json.Unmarshal([]byte(arguments.(string)), &argsMap)
	println(arguments)
	if err != nil {
		fmt.Println("Failed to get config json file: ", err)
		return nil, errors.New("failed to parse json")
	}

	var mainEntry MainEntry
	if argsMap != nil {
		mainEntry, err = parseMainEntry(argsMap)
		if err != nil {
			fmt.Println("an error has occurred while parsing main entry parameters", err)
		}
	}

	var iconData []byte
	var title string

	if len(mainEntry.iconPath) > 0 {
		var data []byte
		data, err := parseIcon(mainEntry.iconPath)
		if err != nil {
			fmt.Println(fmt.Sprintf("An error has occurred while parsing the icon: %s", err))
		}

		if data != nil {
			iconData = data
		}
	}

	if len(mainEntry.title) > 0 {
		title = mainEntry.title
	}

	initialize(title, iconData)

	return "ok", nil
}

func (p *FlutterSystrayPlugin) updateMenuHandler(arguments interface{}) (reply interface{}, err error) {
	var argsMap map[string]interface{}
	err = json.Unmarshal([]byte(arguments.(string)), &argsMap)
	if err != nil {
		fmt.Println("Failed to get config json file: ", err)
		return nil, errors.New("failed to parse json")
	}

	actions, err := parseActionParams(argsMap)
	if err != nil {
		fmt.Println("an error has occurred while parsing action parameters", err)
	}

	p.updateMenu(actions)

	return "ok", nil
}

func (p *FlutterSystrayPlugin) focusHandler(action *SystrayAction) func() {
	return func() {
		p.window.Show()
	}
}

func (p *FlutterSystrayPlugin) closeHandler(action *SystrayAction) func() {
	return func() {
		p.window.SetShouldClose(true)
	}
}

func (p *FlutterSystrayPlugin) eventHandler(action *SystrayAction) func() {
	return func() {
		err := p.invokeSystrayEvent(action)
		if err != nil {
			fmt.Println(fmt.Sprintf("An error has occurred while invoking SystrayEvent: %s", err))
		}
	}
}

func (p *FlutterSystrayPlugin) invokeSystrayEvent(action *SystrayAction) error {
	err := p.channel.InvokeMethod("systrayEvent", action.name)
	if err != nil {
		return err
	}

	return nil
}

func parseMainEntry(entry interface{}) (MainEntry, error) {
	m := entry.(map[string]interface{})
	println(fmt.Sprintf("parse main entry: %s", entry))
	parsed := MainEntry{
		title:    m["title"].(string),
		iconPath: m["iconPath"].(string),
	}
	return parsed, nil
}

func parseActionParams(argsMap map[string]interface{}) ([]SystrayAction, error) {
	var actions []SystrayAction
	for _, v := range argsMap {
		valsMap := v.(map[string]interface{})

		number, _ := strconv.Atoi(valsMap["actionType"].(string))
		action := SystrayAction{
			name:       valsMap["name"].(string),
			label:      valsMap["label"].(string),
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
