package main

func main() {
	executor := &BDGExecutor{}
	executor.Execute("http://btdigg.pw/", "keyword", []string{"vs2008"})
	executor.Parse()
}
