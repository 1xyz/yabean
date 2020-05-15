# yabean

Yet another CLI for beanstalk work queue, built atop [go-beanstalkd](https://github.com/beanstalkd/go-beanstalk) 

The CLI attempts to exercise most of the functionality.

Pre-built static binaries can be obtained from the [releases page](https://github.com/1xyz/yabean/releases)


# Example Usage

```
./yabean --help
usage: yabean [--version] [--addr=<addr>] <command> [<args>...]
options:
   --addr=<addr>  Beanstalkd Address [default: :11300].
   -h, --help
The commands are:
   del        Delete a specific job.
   kick       Kick a buried job (Note: see reserve command to bury a job).
   list       List tubes.
   peek       Peek at a specific job.
   peek-tube  Peek into a specific tube.
   put        Put a job into a beanstalkd tube.
   reserve    Reserve a job from one or more tubes.
   stats      Retrieve serve statistics.
   stats-job  Retrieve statistics for a specific job.
   stats-tube Retrieve statistics for a specific tube.
```

Explore individual commands

For example:
```
./yabean put --help
usage: put [--body=<body>] [--pri=<pri>] [--ttr=<ttr>] [--delay=<delay>] [--tube=<tube>]
options:
    -h, --help
    --body=<body>     body [default: hello]
    --pri=<pri>       job priority [default: 1]
    --ttr=<ttr>       ttr in seconds [default: 10]
    --delay=<delay    job delay in seconds [default: 0]
    --tube=<tube>     tube (topic) to put the job [default: default]

example:
    put --body "hello world"
    put --body "hello world" --tube foo
```

and..

```
yabean reserve --help
usage: reserve [--del|--bury|--release] [options]
options:
    -h, --help
    --timeout=<seconds>   reservation timeout in seconds [default: 0]
    --tubes=<tubes>       csv of tubes [default: default]
    --string              display job's body content as a string [default: false]
  Post reserve actions:
    --bury                bury the job once a job is reserved
    --del                 delete the job (similar to ACK) once a job is reserved
    --release             release the job (similar to NACK) once a job is reserved
  Post reserve action options:
    --pri=<int>           new priority if the job is buried or released [default: 1024]
    --delay=<seconds>     new delay if the job is release [default: 10]
  Other reserve options:
    --touch=<int>         touch (aka renew TTR) the reserved job n times prior to either burying,
                          deleting, releasing or timeout [default: 0]

example:
    watch for reservations on default tube (topic)
    reserve

    watch for reservations on tubes foo & bar with timeout of 10 seconds
    reserve --timeout 10 --tubes=foo,bar

    delete the job after it is reserved frpm the default tube
    reserve --del

    bury the job with a priority 123 after it is reserved from the foo tube
    reserve --tubes=foo --bury --priority 123

    release the job immediately after it is reserved from the bar tube
    reserve --tube=bar --release

    touch a job 5 times and bury it after it is reserved from the foobar tube
    reserve --tube=foobar --touch 5 --bury
```


# Development

Refer the makefile

```
  $  make
 target         ⾖ Description.
 -----------------------------------------------------------------
 build          generate a local build ⇨ bin/yabean
 clean          clean up bin/ & go test cache
 fmt            format go code files using go fmt
 release/darwin generate a darwin target build
 release/linux  generate a linux target build
 tidy           clean up go module file
 ------------------------------------------------------------------
``` 


