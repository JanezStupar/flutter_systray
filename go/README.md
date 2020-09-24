# flutter_systray

This Go package implements the host-side of the Flutter [flutter_systray](https://github.com/JanezStupar/flutter_systray) plugin.

## Usage

Import as:

```go
import flutter_systray "github.com/JanezStupar/flutter_systray/go"
```

Then add the following option to your go-flutter [application options](https://github.com/go-flutter-desktop/go-flutter/wiki/Plugin-info):

```go
flutter.AddPlugin(&flutter_systray.FlutterSystrayPlugin{}),
```
