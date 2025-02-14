SET CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
@REM nps server
go build -ldflags "-s -w -extldflags -static -extldflags -static" -o z_output/linux-amd64/nps ./cmd/nps/nps.go
@REM scu ttu arm
set GOARCH=arm
go build -ldflags "-s -w -extldflags -static -extldflags -static" -o z_output/linux-arm/npc ./cmd/npc/npc.go
@REM desktop
SET GOARCH=amd64
SET GOOS=windows
cd gui/desktop
fyne package -os windows -icon icon.png -executable ../../z_output/windows/npc_desktop.exe

