#!/bin/bash

go build conf/*.go
go build model/*.go
go build dao/*.go
go build api/*.go
go build main.go
go run .
