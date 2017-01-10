#! /bin/bash

rm -rf build/*
npm config set registry https://registry.npm.taobao.org
npm install --global gulp
npm install
npm install bower
/usr/src/app/node_modules/bower/bin/bower install --allow-root
gulp
