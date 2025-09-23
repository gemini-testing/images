#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

QUIET=${QUIET:-""}
DRIVER_ARGS=""
if [ -z "$QUIET" ]; then
    DRIVER_ARGS="--verbose"
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
  if [ -n "$PRISM_PID" ]; then
    kill -TERM "$PRISM_PID"
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

if env | grep -q ROOT_CA_; then
  mkdir -p /tmp/ca-certificates
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/ca-certificates/${certname}.crt
  done
  update-ca-certificates --localcertsdir /tmp/ca-certificates
fi

XVFB_ARGS="-l -n $DISPLAY_NUM -s \"-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp\""
WEBKITDRIVER_CMD="/opt/webkit/bin/WebKitWebDriver --port=5555 --host=0.0.0.0 ${DRIVER_ARGS}"

if [ "$USE_FLUXBOX" = "true" ]; then
    /usr/bin/xvfb-run $XVFB_ARGS /usr/bin/fluxbox -display "$DISPLAY" -log /dev/null 2>/dev/null &
    XVFB_PID=$!

    DISPLAY="$DISPLAY" $WEBKITDRIVER_CMD &
    DRIVER_PID=$!
else
    /usr/bin/xvfb-run $XVFB_ARGS $WEBKITDRIVER_CMD &
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

/usr/bin/prism  &
PRISM_PID=$!

wait
