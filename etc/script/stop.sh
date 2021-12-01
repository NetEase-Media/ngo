#! /bin/bash
process=$1
if [ "$process" = "" ]
then
  echo "please set the process keyword in first arg"
  exit 0
fi
to=$2
if [ "$to" = "" ]
then
  to="10s"
fi
pid=$(ps x | grep $process | grep -v 'stop.sh' | grep -v '.script' | grep -v grep | awk '{print $1}' | head -1)
if [ "$pid" = "" ]
then
  echo "can't find pid with keyword: $process"
  exit 0
fi
echo "$(date '+%Y-%m-%d %H:%M:%S') begin to send term to $pid ..."
kill -15 $pid
echo "$(date '+%Y-%m-%d %H:%M:%S') sign has sended, wait $pid ..."
timeout $to tail --pid=$pid -f /dev/null
status=$?
if [ $status -eq 124 ]
then
  echo "$(date '+%Y-%m-%d %H:%M:%S') timeout with $to"
fi
echo "$(date '+%Y-%m-%d %H:%M:%S') stop finished"
exit 0