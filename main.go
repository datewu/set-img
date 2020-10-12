package main

func main() {
	parseFlag()
	go initK8s()
	server()
}
