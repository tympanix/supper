#!/usr/bin/env bash

useradd supper -r -s /bin/false || :

if [[ ! -d /var/log/supper ]]; then
  mkdir -p /var/log/supper || :
  chgrp supper /var/log/supper || :
  chmod g+s /var/log/supper || :
  touch /var/log/supper/supper.log || :
  chmod g+w /var/log/supper/supper.log || : 
fi
