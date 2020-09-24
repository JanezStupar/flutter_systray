package flutter_systray

import (
	"errors"
	"fmt"
	"github.com/getlantern/systray"
	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"io/ioutil"
)

const channelName = "plugins.flutter.io/flutter_systray"

// FlutterSystrayPlugin implements flutter.Plugin and handles method.
type FlutterSystrayPlugin struct{}

var _ flutter.Plugin = &FlutterSystrayPlugin{} // compile-time type check

// InitPlugin initializes the plugin.
func (p *FlutterSystrayPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("getPlatformVersion", p.handlePlatformVersion)
	channel.HandleFunc("showSystrayIcon", p.showSystrayIcon)
	channel.HandleFunc("clearSystrayIcon", p.clearSystrayIcon)
	return nil
}

func (p *FlutterSystrayPlugin) handlePlatformVersion(arguments interface{}) (reply interface{}, err error) {
	return "go-flutter " + flutter.PlatformVersion, nil
}

func (p *FlutterSystrayPlugin) showSystrayIcon(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})
	path := argsMap["path"].(string)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("file reading error: %s", err))
	}

	systrayOnReady(data)
	return "ok", nil
}

func (p *FlutterSystrayPlugin) clearSystrayIcon(arguments interface{}) (reply interface{}, err error) {
	systrayOnExit()
	return "ok", nil
}

func systrayOnReady(data []byte) {
	systray.SetIcon(data)
}

func systrayOnExit() {
	systray.Quit()
}
