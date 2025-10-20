#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

if [ -z ${VERBOSE+x} ]; then VERBOSE=1; fi
DRIVER_ARGS=${DRIVER_ARGS:-""}
if [ -n "$VERBOSE" ]; then
    DRIVER_ARGS="$DRIVER_ARGS --verbose"
fi

clean() {
  if [ -n "$FILESERVER_PID" ]; then
    kill -TERM "$FILESERVER_PID"
  fi
  if [ -n "$XSELD_PID" ]; then
    kill -TERM "$XSELD_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    pkill -TERM -P "$XVFB_PID"
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$DRIVER_PID" ]; then
    kill -TERM "$DRIVER_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
  if [ -n "$DEVTOOLS_PID" ]; then
    kill -TERM "$DEVTOOLS_PID"
  fi
}

trap clean SIGINT SIGTERM

if env | grep -q ROOT_CA_; then
  mkdir -p $HOME/.pki/nssdb
  certutil -N --empty-password -d sql:$HOME/.pki/nssdb
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/cert.pem
    certutil -A -n ${certname} -t "TC,C,T" -i /tmp/cert.pem -d sql:$HOME/.pki/nssdb
    if cat tmp/cert.pem | grep -q "PRIVATE KEY"; then
      PRIVATE_KEY_PASS=${PRIVATE_KEY_PASS:-\'\'}
      openssl pkcs12 -export -in /tmp/cert.pem -clcerts -nodes -out /tmp/key.p12 -passout pass:${PRIVATE_KEY_PASS} -passin pass:${PRIVATE_KEY_PASS}
      pk12util -d sql:$HOME/.pki/nssdb -i /tmp/key.p12 -W ${PRIVATE_KEY_PASS}
      rm /tmp/key.p12
    fi
    rm /tmp/cert.pem
  done
fi

if env | grep -q CH_POLICY_; then
  for p in $(env | grep CH_POLICY_ | sed 's/CH_POLICY_//'); do
    jsonkey=$(echo $p | sed 's/=.*//')
    jsonvalue=$(echo $p | sed 's/^.*=//')
    cat <<< $(jq --arg key $jsonkey --argjson value $jsonvalue '.[$key] = $value' /etc/opt/chrome/policies/managed/policies.json) > /etc/opt/chrome/policies/managed/policies.json
  done
fi

/usr/bin/fileserver &
FILESERVER_PID=$!

/usr/bin/devtools &
DEVTOOLS_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

while ip addr | grep inet | grep -q tentative > /dev/null; do sleep 0.1; done

XVFB_ARGS="-l -n $DISPLAY_NUM -s \"-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp\""
CHROMEDRIVER_CMD="/usr/bin/chromedriver --port=4444 --allowed-ips='' --allowed-origins='*' ${DRIVER_ARGS}"

if [ "$USE_FLUXBOX" = "true" ]; then
    eval "/usr/bin/xvfb-run $XVFB_ARGS /usr/bin/fluxbox -display \"$DISPLAY\" -log /dev/null 2>/dev/null &"
    XVFB_PID=$!

    eval "DISPLAY=\"$DISPLAY\" $CHROMEDRIVER_CMD &"
    DRIVER_PID=$!
else
    eval "/usr/bin/xvfb-run $XVFB_ARGS $CHROMEDRIVER_CMD &"
    XVFB_PID=$!
fi

wait_for_x_server() {
  local cmd="$1"
  local msg="$2"
  local retcode=1

  until [ $retcode -eq 0 ]; do
    eval "$cmd"
    retcode=$?
    if [ $retcode -ne 0 ]; then
      echo "$msg"
      sleep 0.1
    fi
  done
}

if [ "$USE_FLUXBOX" = "true" ]; then
  wait_for_x_server "DISPLAY=\"$DISPLAY\" wmctrl -m >/dev/null 2>&1" "Waiting X server and window manager..."
else
  wait_for_x_server "xdpyinfo -display \"$DISPLAY\" >/dev/null 2>&1" "Waiting X server..."
fi

if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display "$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /dev/null &
    X11VNC_PID=$!
fi

wait
