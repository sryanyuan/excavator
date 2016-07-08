package main

type IExecutor interface {
	Execute(executeUrl string, key string, values []string) error
	GetResult() interface{}
}
