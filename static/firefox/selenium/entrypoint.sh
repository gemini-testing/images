#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

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
  if [ -n "$SELENIUM_PID" ]; then
    kill -TERM "$SELENIUM_PID"
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
SELENIUM_CMD="/usr/bin/java -Xmx256m -Djava.security.egd=file:/dev/./urandom \
-jar /usr/share/selenium/selenium-server-standalone.jar \
-port 4444 -browserTimeout 120"

if [ "$USE_FLUXBOX" = "true" ]; then
  /usr/bin/xvfb-run $XVFB_ARGS /usr/bin/fluxbox -display "$DISPLAY" -log /dev/null 2>/dev/null &
  XVFB_PID=$!

  DISPLAY="$DISPLAY" $SELENIUM_CMD &
  SELENIUM_PID=$!
else
  /usr/bin/xvfb-run $XVFB_ARGS $SELENIUM_CMD &
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
