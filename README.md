# Golang Embedded Workflow

Simple library for embedded workflow mechanism. Support serial and parallel exectution.

Code abstracted and then extracted from a complex application. I havent checked other available libraries as this one was built from a growing application.


## TODO

- [ ] battle-tested, test errors / job config
- [ ] persistent storage for stateful workflows
- [ ] pass only data to processor and find another a way to handle error (by function returning) and spawning (by enclosing processors ?)