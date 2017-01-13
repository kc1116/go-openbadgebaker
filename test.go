package main

import "os"
import "io/ioutil"
import "fmt"
import "strings"

func main() {
    bytes, err := ioutil.ReadAll(os.Stdin)
    if err == nil{
        fmt.Printf(string(bytes));
    }
    
    var combinedList []int16{}
    
    for i, num := range bytes {
        if i == 0 || i == 2{
            var temp [i]int16
            strToInt(&temp, bytes[i+1])
        }
    }
}

func strToInt(temp *[]int16, )