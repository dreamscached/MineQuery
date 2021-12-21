#!/bin/sh
if [ "$EULA" = "true" ] && [ ! -e "/server/eula.txt" ]; then
  printf "%s" "eula=true" >"/server/eula.txt"
fi

# shellcheck disable=SC2164
cd "/server"
# shellcheck disable=SC2086
java $JAVA_OPTS -jar "/server/server.jar" "nogui" 2>&1 | while read -r line; do
  echo "$line"
  if printf "%s" "$line" | grep -Eq '\[(.*/)?INFO]:? Done'; then
    printf "%s" "1" >> "/server/ready"
  fi
done
