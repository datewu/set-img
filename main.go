package main

func main() {
	parseFlag()
	panicIfErr(initKey)
	panicIfErr(initK8s)
	server()
}
