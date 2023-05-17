#!/bin/sh

GIT_ROOT="$(git rev-parse --show-toplevel)"

log_file=$GIT_ROOT/run.log

(sleep 3 && fortune | tee -a $GIT_ROOT/tmp/fortune3.txt) 2>&1 >> $log_file &
(sleep 4 && fortune | tee -a $GIT_ROOT/tmp/fortune4.txt) 2>&1 >> $log_file &
(sleep 5 && fortune | tee -a $GIT_ROOT/tmp/fortune5.txt) 2>&1 >> $log_file &
(sleep 6 && fortune | tee -a $GIT_ROOT/tmp/fortune6.txt) 2>&1 >> $log_file &

go run ./cmd/rebounce/main.go -dir $GIT_ROOT/tmp
