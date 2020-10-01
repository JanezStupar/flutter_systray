# flutter_systray

Flutter Systray is a plugin for `go-flutter-desktop` that enables support for systray menu for desktop flutter apps.

Supports `MacOS`, `Windows` and `Linux` platforms via [`systray`](https://github.com/getlantern/systray)

This plugin implements limited support. There are no submenus, checkboxes and such. PRs are welcome.

Don't forget to check the example app!

## Getting Started

Install the plugin as is customary.

Import as:
```go
import "github.com/zephylac/title_bar"
```

Then add the following option to your go-flutter [application options](https://github.com/go-flutter-desktop/go-flutter/wiki/Plugin-info):
```go
flutter.AddPlugin(&title_bar.TitleBarPlugin{}),
```
## API

Below we initialize the systray menu and give it a focus action that bring the flutter app window to front.
```dart
MainEntry main = MainEntry(
  title: "title",
  iconPath: path,
);

FlutterSystray.initSystray(main, [
  SystrayAction(
      name: "focus",
      label: "Focus",
      tooltip: "Bring application window into focus",
      iconPath: '',
      actionType: ActionType.Focus),
]).then();
```
`MainEntry` - represents the root node of the systray menu. It can have an icon (Win, Linux, Mac) or/and a title and tooltip (Mac).
`[]SystrayAction` - a list of systray menu actions. Actions can have an icon, title and tooltip. Name serves as unique identifier. 


To add more actions we can call `addActions` function:
```dart
FlutterSystray.addActions([
SystrayAction(
    name: "counterEvent",
    label: "Counter event",
    tooltip: "Will notify the flutter app!",
    iconPath: '',
    actionType: ActionType.SystrayEvent),
SystrayAction(
    name: "systrayEvent2",
    label: "Event 2",
    tooltip: "Will notify the flutter app!",
    iconPath: '',
    actionType: ActionType.SystrayEvent),
SystrayAction(
    name: "quit", label: "Quit", tooltip: "Close the application", iconPath: '', actionType: ActionType.Quit)
]);
```

We can also register callback handlers for events triggered by systray:
```dart 
FlutterSystray systemTray = FlutterSystray.init();
systemTray.registerEventHandler("counterEvent", () {
  setState(() {
    _counter += 1;
  });
});
```
Flutter Systray matches callbacks by `SystrayAction.name` property.


## Available SystrayActions

At the moment Flutter Systray supports three kinds of actions, which allow you to call platform operation and  trigger custom events in your flutter app:
```dart
enum ActionType {
  Quit, // Action will trigger application shutdown
  Focus, // Action will trigger GLFW `window.Show` and bring flutter app to front
  SystrayEvent // Action will trigger an event that will call a registered callback in flutter app
}
```
