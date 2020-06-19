package main

func main() {
	go watch()
	server := &Server{
		Ip:   "127.0.0.1",
		Port: 8882,
	}
	server.Run()
}
