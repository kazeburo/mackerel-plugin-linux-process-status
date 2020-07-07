# mackerel-plugin-linux-process-status

mackerel metric plugin for CPU/FDs/Memory usage of each linux process.

```
$ ./mackerel-plugin-linux-process-status --key-prefix test -p 54321
process-status.fds_test.count   117     1593759145
process-status.fds_test.max     40960   1593759145
process-status.fds_usage_test.percentage        0.285645        1593759145
process-status.cpu_test.percentage      0.000000        1593759145
process-status.mem_test.used    25716518912     1593759145
process-status.mem_test.max     64215793664     1593759145
process-status.mem_usage_test.percentage        40.047031       1593759145
```

