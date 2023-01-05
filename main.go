package main

func main() {
	server := NewServer("127.0.0.1", 9999)
	server.Start()
}
