# Queue thing

## Given task
- Read a line
- Write to service
- Recieve line
- Write to file
repeat

## Idea
The task states to use only standard lib for the queue service. So I decided
to use golangs builtin channels for this along with the built in http service.
With a bit of creative use of HTTP status codes I can model the process of
"I'm going to send a file now" and "Here's a line" and "I'm Done".

Python Asyncio could be a way to async both sending and receiving the file.

So that's what I've done here.

## To run it on a Linux with `make(1)`, Go 1.10 and Python 3.6 installed
```
$ make build-go
$ ./queue &
$ ./client file1 outputfile
```

## To Run it with make and docker installed
```
$ make build
$ ./dockerqueue
$ ./dockerclient file1 file_output
```
When running inside of docker you can't reference files not nested within
`PWD`.

## Missing

- No tests of neither client or service
- If a queue get's half a file and you then cancel you need to restart the
  queue service before you can rerun the client on that specific filename
  - no resume (you'll get 409 Conflict)
- No deployment of service
- No logging in service
- No stats/monitoring metrics exposed
- It's not extremely fast - I would like a profiler attached to the client to
  see what is taking the time. If I leave in print statements it looks very
  much like python interleaves GET and POSTs more or less one to one.
- Hardly any comments in the code
