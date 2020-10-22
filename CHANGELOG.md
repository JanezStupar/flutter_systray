## 0.3.5
Clean up go dependencies

## 0.3.4
Replaced getlantern systray with r10v fork to improve menu behavior on windows (menu still doesn't always appear)

## 0.3.3
Fix import in windows implementation

## 0.3.2
Fix an issue where systray menu would sometimes not display in windows after startup.

## 0.3.1
Systray menu entries now display in order declared by Flutter app.

## 0.3.0
Made init call synchronous and made other changes that allow the plugin to work on MacOS correctly.

## 0.2.0
Replaced `systray` with `trayhost`
Changed IPC format from serialized map to JSON
Removed tooltips and icons from flutter_systray menu items as underlying library doesn't support them

## 0.1.1
Support for basic systray functionalities

## 0.0.1
A proof of concept
