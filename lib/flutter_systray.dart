// You have generated a new plugin project without
// specifying the `--platforms` flag. A plugin project supports no platforms is generated.
// To add platforms, run `flutter create -t plugin --platforms <platforms> .` under the same
// directory. You can also find a detailed instruction on how to add platforms in the `pubspec.yaml` at https://flutter.dev/docs/development/packages-and-plugins/developing-packages#plugin-platforms.

import 'dart:async';

import 'package:flutter/services.dart';

enum ActionType { Quit, Focus, SystrayEvent }

class SystrayAction {
  final ActionType actionType;
  final String name;
  final String label;
  final String tooltip;
  final String iconPath;

  SystrayAction({this.name, this.label, this.tooltip, this.iconPath, this.actionType});

  Map<String, String> serialize() {
    return <String, String>{
      "name": this.name,
      "label": this.label,
      "tooltip": this.tooltip,
      "iconPath": this.iconPath,
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
  static Future<String> initSystray(String iconPath, List<SystrayAction> actions) async {
    Map<String, Map<String, String>> map = _serializeActions(actions);
    map["mainIcon"] = <String, String>{
      "iconPath": iconPath,
    };

    String value = await _channel.invokeMethod('initSystray', map);
    return value;
  }

  static Future<String> addActions(List<SystrayAction> actions) async {
    Map<String, Map<String, String>> map = _serializeActions(actions);
    String value = await _channel.invokeMethod('addActions', map);
    return value;
  }

  static Map<String, Map<String, String>> _serializeActions(List<SystrayAction> actions) {
    var map = <String, Map<String, String>>{};

    actions.forEach((SystrayAction element) {
      map[element.name] = element.serialize();
    });

    return map;
  }

  registerEventHandler(String handlerKey, Function handler) {
    if (_initialized == false) {
      throw Exception("not initialized, call init() before registering event handlers");
    }

    _handlers[handlerKey] = handler;
  }
}
