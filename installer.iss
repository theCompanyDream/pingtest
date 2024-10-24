; pingserver_installer.iss - Inno Setup Script for Ping Server

[Setup]
AppName=Ping Server
AppVersion=1.0.0
DefaultDirName={pf}\PingServer
DefaultGroupName=Ping Server
OutputBaseFilename=PingServerInstaller
Compression=lzma
SolidCompression=yes
SetupIconFile=server.ico

[Files]
; Include the pre-built executable
Source: "*"; DestDir: "{app}"; Flags: ignoreversion


[Icons]
; Create a Start Menu shortcut
Name: "{group}\Ping Server"; Filename: "{app}\PingTest.exe"

; Create a Desktop shortcut (optional)
Name: "{userdesktop}\Ping Server"; Filename: "{app}\PingTest.exe"; Tasks: desktopicon

[Tasks]
; Option to create a desktop shortcut
Name: "desktopicon"; Description: "Create a &desktop icon"; GroupDescription: "Additional icons:"; Flags: unchecked

[Run]
; Launch Ping Server after installation
Filename: "{app}\PingTest.exe"; Description: "Launch Ping Server"; Flags: nowait postinstall skipifsilent
