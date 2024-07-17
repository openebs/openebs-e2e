# Mayastor E2E fio test pod
## Introduction
Derived from `dmonakhov/alpine-fio`

### Command line options
 * `sleep <sleep seconds> ;`
   * sleep for `<sleep seconds>`, execution of sleep is synchronous, further options processing paused for
     the duration of the sleep
 * `segfault <delay seconds> ;`
   * wait `<delay seconds>` and signal running child processes with SIGKILL and raise SIGSEGV
 * `sigterm <delay seconds> ;`
   * wait `<delay seconds>` and signal running child processes with SIGTERM, handle child process termination
 * `exitv <exit value> ;`
   * override exit code value
 * `makefile <filesystempath> <relativefilepath> <availblockspercent|availblockslessby|bytes> <integer value> ;`
   * create a file using fallocate with size based on available block count on the filesystem
    * `availblockspercent`: file size is percentage of available blocks, range 1 - 100
    * `availblockslessby`: file size is available block minus `<integer value>`
    * `bytes`: file size is `<integer value>` bytes
    * if the file already exists it is deleted and created again.
    * execution fails if the file could not be created
 * `liveness <target> <read interval> <timeout> ;`
   * perform liveness checks on target. Block 1 is read from target every `<read_interval>` seconds,
     execution is terminated if no read was successful for `<timeout>` seconds.
 * `-- <fio args...> <;|&&>` ( Note deprecated, use `--- fio ....` instead)
   * fork and run `fio`
   * delimited by `&&` (wait for process to complete), or `;` (do not wait for process to complete) at the end
    * if end of arguments is reached then do not wait for process
   * fork and run fio, fio arguments are between  delimiters mentiond above
   * multiple occurrences of this sequence are supported, a new separate process is created for each occurrence which runs concurrently
 * `--- <executable> <args....> <;|&&>`
   * fork and run `<executable>`
   * delimited by `---` at the start and `&&` (wait for process to complete), or `;` (do not wait for process to complete) at the end
    * if end of arguments is reached then do not wait for process
   * each executable is run as a forked process and if delimited by `&` asynchronously
   * multiple occurrences of this sequence are supported, a new separate process is created for each occurrence.
 * `exitv <v>` override exit value - to simulate (pod) failure.
 * `sessionId` set session ID, used as distinguishing string for artefacts, etc.
 * `postopsleep <sleep seconds>` sleep after all operations are completed. This is a developer debug option, keeps the pod alive for the defined period.
Thus allows developer to exec a shell on the pod and examine the pod.

Execution will only complete after all forked processes (if any) have completed as well as inline sleep and signal generation actions.

For legacy compatibility where implicit start of a new option is detected

### Exit value
 * If `exitv` is specified that is *always* returned.
 * 0 if all instances of forked processes ran successfully exit value is 0
 * If any one instance of forked processes fails, the exit value is the exit value of the failed forked process
 * If multiple instances of forked processes fail, the exit value is the exit value of one of the failed forked process
 * If multiple forked process fail the exit value is the exit value of a failing forked process

### Examples
 * time based run
   * `--- fio .... ; sigterm 960`
    - fork and run fio, sleep for 960 seconds and send sigterm to fio
      if fio has not faulted prior to this the executable will complete successfully
 * fio file system run
   * `makefile path_to_fs rel_path_to_file availblockslessby 10 --- fio .... --filename=path_to_fs/rel_path_to_file ... ;`
    - create file `path_to_fs/rel_path_to_file` of size `available blocks - 10` && run fio without `--size=nnn` option
    - execution terminates if the file could not be created
 * sleep before and after fio run
   * `sleep 100 --- fio ...... && sleep 200`
    - sleep 100 seconds, run fio with `......` then sleep 100 seconds (then execution completes)
 * multiple fio runs
   * `--- fio ...... && --- fio ......`
    - run fio twice (with different arguments for example)
    - execution terminates if and when first fio run fails
   * `--- fio ...... && sleep 10 && --- fio ......`
    - run fio twice (with different arguments for example), sleeping 10
    - execution terminates if and when first fio run fails
 * liveness
  * `liveness /dev/sdm 1 60 ; fio --filename=/dev/sdm .....`
   - run fio on `/dev/sdm` whilst checking that `/dev/sdm` is readble every second, fail if it is not for 60 seconds.

## building
Run `./build.sh`

This builds the image `openebs/e2e-fio`
