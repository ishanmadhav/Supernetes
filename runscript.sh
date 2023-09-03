#!/bin/bash

# Bash commands go here
gnome-terminal --title="SuperCache" -- go run cmd/cache/main.go
gnome-terminal --title="SuperAPIServer" -- go run cmd/superapiserver/main.go
gnome-terminal --title="Superlet" -- go run cmd/superlet/main.go
gnome-terminal --title="SuperController" -- go run cmd/supercontroller/main.go
