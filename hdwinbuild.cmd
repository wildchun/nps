@REM desktop
SET GOARCH=amd64
SET GOOS=windows
cd gui/desktop
fyne package -os windows -icon icon.png -executable ../../z_output/windows/npc_desktop.exe



