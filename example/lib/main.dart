import 'package:flutter/material.dart';
import 'package:flutter_systray/flutter_systray.dart';
import 'package:path/path.dart' as p;

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  String path = p.absolute('go/assets', 'icon.png');
  FlutterSystray.initSystray(path, [
    SystrayAction(
        name: "quit", label: "Quit", tooltip: "Close the application", iconPath: '', actionType: ActionType.Quit)
  ]).then((value) async {
    FlutterSystray.addActions([
      SystrayAction(
          name: "counterEvent",
          label: "Counter event",
          tooltip: "Will notify the flutter app!",
          iconPath: '',
          actionType: ActionType.SystrayEvent),
      SystrayAction(
          name: "focus",
          label: "Focus",
          tooltip: "Bring application window into focus",
          iconPath: '',
          actionType: ActionType.Focus),
      SystrayAction(
          name: "systrayEvent2",
          label: "Event 2",
          tooltip: "Will notify the flutter app!",
          iconPath: '',
          actionType: ActionType.SystrayEvent),
    ]);
  });
  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  // Register an event handler
  final FlutterSystray systemTray = FlutterSystray.init();

  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  int _counter = 0;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Plugin example app'),
        ),
        body: Center(
          child: Text(
              'There should be a menu with a Hover icon in the systray.\n Number of times that the counter was triggered: $_counter '),
        ),
      ),
    );
  }

  @override
  void initState() {
    widget.systemTray.registerEventHandler("counterEvent", () {
       setState(() {
         _counter += 1;
       });
    });

    super.initState();
  }
}
