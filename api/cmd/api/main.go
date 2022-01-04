package main

import "api/internal/service"

func main() {
	(&service.Service{}).Run()
}
