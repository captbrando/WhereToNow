package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

// URLs as a grouping of URL structs
type URLs []URL

// Data as read from the config file
var Data URLs

// URL is our struct that has time and url
type URL struct {
	Time int    `json:"time,string"` // some magic happening here, for readability, we will convert string to int.
	URL  string `json:"url"`
}

// WeDebugging grabs the Bool from debug flag and sets it here. We are doing it this way to ensure the
// usage dump includes the -debug flag up front, we can then reference its value insde the ServeHTTP func.
var WeDebugging bool

func main() {
	// Get command line args first.
	confPtr := flag.String("conf", "config/config.json", "/path/to/config")
	netPtr := flag.String("network", ":8080", "Your listening port, see the Readme for usage")
	debugPtr := flag.Bool("debug", false, "To turn on debugging, pass true")
	flag.Parse()            // the -h flag will print usage.
	WeDebugging = *debugPtr // set this so we can see it in our helper function

	// Read in our config file. This becomes important with Docker images.
	file, err := ioutil.ReadFile(*confPtr)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	// Decode the JSON data and read it into our cache.
	err = json.Unmarshal(file, &Data)
	if err != nil {
		fmt.Printf("JSON Unmarshal error: %v\n", err)
		os.Exit(2)
	}

	// I was lazy in my initial JSON config file creation. I figured others might make mistakes
	// The Readme is VERY specific here, don't be a dodo like me. Be OCD and do it in order.
	sort.Sort(byTime(Data))

	// Do some error checking of the config file. Two error conditions we can't have:
	// 	1) Times larger than 2359
	if Data[len(Data)-1].Time >= 2400 {
		fmt.Printf("Config file cannot contain times larger than 2359.")
		os.Exit(3)
	}
	// 	2) And times smaller than 0.
	if Data[0].Time < 0 {
		fmt.Printf("Config file cannot contain times smaller than 0.")
		os.Exit(4)
	}

	// Check for dups by starting with index 0 of our slice of structs, and iterate through to see if we find one.
	dupDetect := Data[0].Time
	for i := 1; i < len(Data); i++ {
		if dupDetect == Data[i].Time {
			fmt.Printf("Duplicate Time entry '%v' found in JSON. Time entries must be unique.", dupDetect)
			os.Exit(5)
		} else {
			dupDetect = Data[i].Time
		}
	}

	// Now that we have all our stuff set up and we've checked our config, let's get it on!
	// Start that server and wait.
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(*netPtr, nil)
}

// ServeHTTP to handle our request.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// debugPtr := flag.Bool("debug", false, "To turn on debugging, pass true")
	// We need the time in a HHMM format as a string first, then converted to an integer value. The line below
	// does this by using the Format function and passing the Reference Time values of 1504 (HHMM). See the
	// Time.Format() reference for more.
	currTime, _ := strconv.Atoi(time.Now().Format("1504"))

	// The logic here is to step backwards through the list to find a point where the current time
	// is larger than the time we have slated in the config. So if the current time is 0130, but
	// the smallest time listed is 0230, we will proceed all the way through the loop with no hits.
	// So the default value there is at the end of the for loop, which would be the largest time
	// in the config file (because we would roll over 2359-2400/0000). 2400 is not a valid time.
	for i := len(Data) - 1; i >= 0; i-- {
		if currTime >= Data[i].Time {
			if WeDebugging {
				fmt.Printf("Current time is %v, currTime is %v, and the Data[%v] time is %v\n", time.Now(), currTime, i, Data[i].Time) // debugging.
			}
			http.Redirect(w, r, Data[i].URL, 303)
			return
		}
	}
	// Whoops, we found nothing. Must be the last slot of the day that carries over.
	http.Redirect(w, r, Data[len(Data)-1].URL, 303)
	return
}

// Some helper functions. This  one is to ensure we don't have duplicate times in our config file
// Sort our times
type byTime []URL

func (a byTime) Len() int           { return len(a) }
func (a byTime) Less(i, j int) bool { return a[i].Time < a[j].Time }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
