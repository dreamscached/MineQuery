#!/bin/bash
version_k=("1.7.2" "1.6.2" "1.4.2" "1.2.5")
declare -A versions
versions=(
  ["1.7.2"]="https://launcher.mojang.com/mc/game/1.7.2/server/3716cac82982e7c2eb09f83028b555e9ea606002/server.jar"
  ["1.6.2"]="https://launcher.mojang.com/mc/game/1.6.2/server/01b6ea555c6978e6713e2a2dfd7fe19b1449ca54/server.jar"
  ["1.4.2"]="https://launcher.mojang.com/mc/game/1.4.2/server/5be700523a729bb78ef99206fb480a63dcd09825/server.jar"
  ["1.2.5"]="https://launcher.mojang.com/mc/game/1.2.5/server/d8321edc9470e56b8ad5c67bbd16beba25843336/server.jar"
)

port=25565
for ver in "${version_k[@]}"; do
  docker build -t "minecraft:$ver" --build-arg SERVER_URL="${versions[$ver]}" "test/docker"
  containers+=("$(docker run -d --rm -p "$port:25565" -e EULA=true "minecraft:$ver")")
  port="$((port + 1))"
done

cleanup() {
  for cont in "${containers[@]}"; do docker stop -t 0 "$cont"; done
}
trap cleanup EXIT

sleep 30
go test -v "./..."
