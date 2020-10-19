// You have generated a new plugin project without
// specifying the `--platforms` flag. A plugin project supports no platforms is generated.
// To add platforms, run `flutter create -t plugin --platforms <platforms> .` under the same
// directory. You can also find a detailed instruction on how to add platforms in the `pubspec.yaml` at https://flutter.dev/docs/development/packages-and-plugins/developing-packages#plugin-platforms.

import 'dart:async';
import 'dart:convert';

import 'package:flutter/services.dart';

enum ActionType {
  Quit, // Action will trigger application shutdown
  Focus, // Action will trigger GLFW `window.Show` and bring flutter app to front
  SystrayEvent // Action will trigger an event that will call a registered callback in flutter app
}

/*
* Root systray entry
* */
class MainEntry{
  final String title;
  final String iconPath;

  MainEntry({this.title="", this.iconPath=""});

  Map<String, String> serialize() {
    return <String, String>{
      "title": this.title,
      "iconPath": this.iconPath,
    };
  }
}

class SystrayAction {
  final ActionType actionType;
  final String name;
  final String label;

  SystrayAction({this.name, this.label, this.actionType});

  Map<String, String> serialize() {
    return <String, String>{
      "name": this.name,
      "label": this.label,
      "actionType": this.actionType.index.toString()
    };
  }
}

class FlutterSystray {
  static const MethodChannel _channel = const MethodChannel('plugins.flutter.io/flutter_systray');
  final _handlers = <String, Function>{};
  bool _initialized = false;

  FlutterSystray.init() {
    _channel.setMethodCallHandler((MethodCall methodCall) async {
      if (methodCall.method == "systrayEvent") {
        Function handler = _handlers[methodCall.arguments];
        if (handler != null) {
          handler();
        }
      }
    });
    _initialized = true;
  }

  /*
  * Show a systray icon
  * */
  static Future<String> initSystray(MainEntry main) async {
    String value = await _channel.invokeMethod('initSystray', jsonEncode(main.serialize()));
    return value;
  }

  static Future<String> updateMenu(List<SystrayAction> actions) async {
    List<Map<String, String>> map = _serializeActions(actions);
    String json = jsonEncode(map);
    String value = await _channel.invokeMethod('updateMenu', json);
    return value;
  }

  static List<Map<String, String>> _serializeActions(List<SystrayAction> actions) {
    var result = <Map<String, String>>[];

    actions.forEach((SystrayAction element) {
       result.add(element.serialize());
    });

    return result;
  }

  registerEventHandler(String handlerKey, Function handler) {
    if (_initialized == false) {
      throw Exception("not initialized, call init() before registering event handlers");
    }

    _handlers[handlerKey] = handler;
  }
}
