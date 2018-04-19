#!/bin/sh
set -e

LOGDIR=/var/log/supper

useradd supper -r -s /bin/false || :

if [[ ! -d $LOGDIR ]]; then
  mkdir -p $LOGDIR || :
  chgrp supper $LOGDIR || :
  chmod g+s $LOGDIR || :
  touch $LOGDIR/supper.log || :
  chmod g+w $LOGDIR/supper.log || :
fi
