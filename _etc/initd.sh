#!/bin/sh

# scp _etc/initd.sh root@37.59.112.124:/etc/init.d/ahouhpuc && ssh -T root@37.59.112.124 "chmod +x /etc/init.d/ahouhpuc"
# update-rc.d ahouhpuc defaults

### BEGIN INIT INFO
# Provides:        ahouhpuc
# Required-Start:  $network
# Required-Stop:
# Default-Start:   2 3 4 5
# Default-Stop:    0 1 6
### END INIT INFO

AOP_ROOT=/home/martin/ahouhpuc

SERVER=$AOP_ROOT/server
SERVER_NAME=ahouhpuc
SERVER_PID=$AOP_ROOT/server.pid
SERVER_LOGFILE=$AOP_ROOT/server.log

. /lib/lsb/init-functions
. "$AOP_ROOT/env"

case "$1" in
	start)
		start-stop-daemon --start \
      --chuid martin \
			--background \
			--no-close \
			--exec "$SERVER" \
			--pidfile "$SERVER_PID" \
			--make-pidfile \
			-- >> "$SERVER_LOGFILE" 2>&1
		;;

	stop)
		start-stop-daemon --stop --pidfile "$SERVER_PID"
		;;

	restart|force-reload)
		pid=`cat "$SERVER_PID" 2>/dev/null`
		[ -n "$pid" ] \
			&& ps -p $pid > /dev/null 2>&1 \
			&& $0 stop
		$0 start
		;;

	status)
		status_of_proc "$SERVER" "$SERVER_NAME"
		;;

	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
		;;
esac

exit 0
