SET CGO_ENABLED=0
set GOARCH=arm
set GOOS=linux
go build -ldflags="-s -w" -o z_output/npc cmd/npc/npc.go