package main

import (
	"os"
)

func dbhost() string {
	return os.Getenv("dbhost")
}

func dbname() string {
	return os.Getenv("dbname")
}

func dbusername() string {
	return os.Getenv("dbusername")
}

func dbpassword() string {
	return os.Getenv("dbpassword")
}
