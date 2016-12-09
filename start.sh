#!/bin/bash

CONF_JS_FILE=/dataman/frontend/src/conf.js

if [ $UI_BACKEND ]
then
  sed -i "s#UI_BACKEND#$UI_BACKEND"g $CONF_JS_FILE
else
  echo "ENV UI_BACKEND is required, exit."
  exit
fi

./mola
