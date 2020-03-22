# Contribute

## Prerequisites

- [GO](https://golang.org/dl/)
- Kubernetes configuration

## UI

In the `cmd` directory, it has the `mock` command to mocking resources. The aim of this command is to simulate your codes of UI, writing test codes of UI is not easy, so that you can check it works as you want.

When you run the command like `go run cmd/mock/main.go`, you can see the process starts with resource which is defined in the `testdata` directory and also the context of namespace is same as your Kubernetes configuration (you can switch by the command `:ns`). 
