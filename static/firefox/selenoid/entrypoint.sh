#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

cp /home/selenium/browsers.json /tmp

VERBOSE=${VERBOSE:-""}
DRIVER_ARGS=${DRIVER_ARGS:-""}
if [ -n "$VERBOSE" ]; then
    DRIVER_ARGS="$DRIVER_ARGS, \"--log\", \"debug\""
fi
sed -i "s|@@DRIVER_ARGS@@|$DRIVER_ARGS|g" /tmp/browsers.json

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
  if [ -n "$SELENOID_PID" ]; then
    kill -TERM "$SELENOID_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/fileserver &
FILESERVER_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

XVFB_ARGS="-l -n $DISPLAY_NUM -s \"-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp\""
SELENOID_CMD="/usr/bin/selenoid -conf /tmp/browsers.json -disable-docker -timeout 1h -max-timeout 24h -enable-file-upload -capture-driver-logs"

if [ "$USE_FLUXBOX" = "true" ]; then
  /usr/bin/xvfb-run $XVFB_ARGS /usr/bin/fluxbox -display "$DISPLAY" -log /dev/null 2>/dev/null &
  XVFB_PID=$!

  DISPLAY="$DISPLAY" $SELENOID_CMD &
  SELENOID_PID=$!
else
  /usr/bin/xvfb-run $XVFB_ARGS $SELENOID_CMD &
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

if env | grep -q ROOT_CA_; then
  while true; do
    if certDB=$(ls -d /tmp/rust_mozprofile*/cert9.db 2>/dev/null); then
      break
    else
      sleep 0.1
    fi
  done
  certdir=$(dirname ${certDB})
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/cert.pem
    certutil -A -n ${certname} -t "TC,C,T" -i /tmp/cert.pem -d sql:${certdir}
    if cat tmp/cert.pem | grep -q "PRIVATE KEY"; then
      PRIVATE_KEY_PASS=${PRIVATE_KEY_PASS:-\'\'}
      openssl pkcs12 -export -in /tmp/cert.pem -clcerts -nodes -out /tmp/key.p12 -passout pass:${PRIVATE_KEY_PASS}  -passin pass:${PRIVATE_KEY_PASS}
      pk12util -d sql:${certdir} -i /tmp/key.p12 -W ${PRIVATE_KEY_PASS}
      rm /tmp/key.p12
    fi
    rm /tmp/cert.pem
  done
fi

wait
