# beanstalk-cli

A  CLI for beanstalk work queue, built atop [go-beanstalkd](https://github.com/beanstalkd/go-beanstalk) 

Pre-built binaries can be obtained from the [releases page](https://github.com/1xyz/beanstalk-cli/releases)

# Example Usage

```
./bscli --help
usage: bs-client [--version] [--addr=<addr>] <command> [<args>...]
options:
   --addr=<addr>  Beanstalkd Address [default: :11300].
   -h, --help
The commands are:
   put        Put a job into a beanstalkd tube.
   reserve    Reserve a job from one or more tubes.

```

Explore individual commands

For example:
```
./beanstalk-cli put --help
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


# Development

Refer the makefile

```
$ make
``` 


