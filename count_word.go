package main


import ("fmt"
	"log"
	"io/ioutil"
	"os"
	"bufio"
	"sync"
	"strings"
	)


func reader(filename string) []string {
    file, err := os.Open(filename)
    if err != nil {
           log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var txtlines []string

    for scanner.Scan() {
        txtlines = append(txtlines, scanner.Text())
    }
    defer file.Close()
    return txtlines
}

func counter(contents []string) map[string]int {
    map_m := make(map[string]int)
    for _, eachline := range contents {
	str := strings.Replace(eachline, ".", "", -1)
	str = strings.Replace(eachline, ",", "", -1)
	for _,char := range strings.Split(str, " ") {
		if char != " " {
			map_m[char] += 1
		}
	}
    }
    return map_m
}

func executer(i chan int, filename string, chan_main chan map[string]int, wg *sync.WaitGroup){
    contents := reader(filename)
    map_m := counter(contents)
    chan_main <- map_m
    <-i
    wg.Done()
}


func main() {
    files, err := ioutil.ReadDir("./data")
    if err != nil {
	    log.Fatal(err)
    }
    i := make(chan int, 2)
    chan_main := make(chan map[string]int, 10)
    var wg sync.WaitGroup
    for _, f := range files {
	  wg.Add(1)
	  i <- 1
          go executer(i, "./data/" + f.Name(), chan_main, &wg)

    }
    log.Println("waiting")
    wg.Wait()
    log.Println("done waiting")
    close(chan_main)

    map_result := make(map[string]int)
    for map_m := range chan_main {
	    for key, val := range map_m{
		    map_result[key] += val
	    }
    }
    for key, val := range map_result{
	    fmt.Println(key,": ", val)
    }
}

