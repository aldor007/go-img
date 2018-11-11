#!/usr/bin/env bash

pkill -f transformerTest

PORT=8082
export PORT

go run main.go -transformerTest &
pid=$!
sleep 1

./node_modules/.bin/mocha test-int/*.Spec.js
TEST_RESULT=$?
echo
kill -9  $pid
exit ${TEST_RESULT}