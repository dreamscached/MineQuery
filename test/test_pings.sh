#!/bin/bash
version_k=("1.7.2" "1.6.2" "1.4.2" "1.2.5")
declare -A versions
versions=(
  ["1.7.2"]="https://launcher.mojang.com/mc/game/1.7.2/server/3716cac82982e7c2eb09f83028b555e9ea606002/server.jar"
  ["1.6.2"]="https://launcher.mojang.com/mc/game/1.6.2/server/01b6ea555c6978e6713e2a2dfd7fe19b1449ca54/server.jar"
  ["1.4.2"]="https://launcher.mojang.com/mc/game/1.4.2/server/5be700523a729bb78ef99206fb480a63dcd09825/server.jar"
  ["1.2.5"]="https://launcher.mojang.com/mc/game/1.2.5/server/d8321edc9470e56b8ad5c67bbd16beba25843336/server.jar"
)

echo "Building server Docker images."
for ver in "${version_k[@]}"; do
  docker build -t "minecraft:$ver" --build-arg SERVER_URL="${versions[$ver]}" "test/docker" >/dev/null 2>&1 &
done
wait

echo "Starting Docker containers."
mkfifo -m a+rw /tmp/mcready
port=25565
for ver in "${version_k[@]}"; do
  containers+=("$(docker run -d --rm -p "$port:25565" -e EULA=true -v "/tmp/mcready:/server/ready" "minecraft:$ver")")
  port="$((port + 1))"
done

cleanup() {
  rm "/tmp/mcready"
  for cont in "${containers[@]}"; do { docker stop -t 0 "$cont" >/dev/null 2>&1 & } done
  for ver in "${version_k[@]}"; do { docker image rm "minecraft:$ver" >/dev/null 2>&1 & } done
}
trap cleanup EXIT

echo "Waiting for test servers to start."
nf=
while [ "${#nf}" -lt "${#version_k[@]}" ]; do
  nf="$nf$(cat "/tmp/mcready")"
done

go test -v "./..."
