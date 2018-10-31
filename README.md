# Golang Embedded Workflow

Simple library for embedded workflow mechanism. Support serial and parallel exectution.

Code abstracted and then extracted from a complex application. I havent checked other available libraries as this one was built from a growing application.


## TODO

- [ ] battle-tested, test errors / job config
- [ ] persistent storage for stateful workflows
- [X] pass only data to processor and find another a way to handle error (by function returning)
- [X] pass only data to processor and find another a way for spawning (by enclosing processors ?) : enclose method with workflow instance and pass it to the closure to spawn a new Job `workflow.NewJob()` so workflow to pass data output from 1 process to next one + test
- [X] caller to provide processor to manage result instead of workflow to provide a channel for return

## Changelog
- 2018-10-30: processor can spawn new job
- 2018-10-30: caller to provide processor to manage result instead of workflow to provide a channel for return.