#! /bin/bash

rm -rf build/*
npm config set registry https://registry.npm.taobao.org
npm install --global gulp
npm install
bower install --allow-root
gulp
