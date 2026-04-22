#!/usr/bin/env bash
set -euo pipefail

# fio --name=<name> --ioengine=<engine> --iodepth=<depth> --rw=<mode> --bs=<size> --size=<total_size> --numjobs=<jobs> --runtime=<secs>
fio --name=rand-rw --ioengine=libaio --iodepth=64 --rw=randrw --bs=4k --size=256m --numjobs=4 --runtime=60 --group_reporting
