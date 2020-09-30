import 'package:flutter/material.dart';
import 'package:flutter_systray/flutter_systray.dart';
import 'package:path/path.dart' as p;

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  String path = p.absolute('go/assets', 'icon.png');
  FlutterSystray.initSystray(path, [
    SystrayAction(
        name: "quit",
        label: "Quit",
        tooltip: "Close the application",
        iconPath: '',
        actionType: ActionType.Quit
    )
  ]).then((value) async {
    FlutterSystray.addActions([
      SystrayAction(
          name: "callback",
          label: "Callback",
          tooltip: "Will notify the flutter app!",
          iconPath: '',
          actionType: ActionType.Callback),
      SystrayAction(
          name: "focus",
          label: "Focus",
          tooltip: "Bring application window into focus",
          iconPath: '',
          actionType: ActionType.Focus),
    ]);
  });
  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Plugin example app'),
        ),
        body: Center(
          child: Text('There should be a menu with a Hover icon in the systray.\n'),
        ),
      ),
    );
  }
}
