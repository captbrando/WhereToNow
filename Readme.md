# Where To Now

**WhereToNow** is a small project designed to help me learn Go and Docker Containers. It's designed to vary the webpage that might show up when a user inside a corporation hits "Home" on their default browser. The idea being that "Home" might be a different target depending on the time of day. For reference, this was written during the COVID-19 pandemic. 

When workers fire up their laptops at home in the morning and open their browser, content managers could first send users to an internal COVID-19 status page (or perhaps something from the CDC or WHO). When the local stock markets open, perhaps route them to websites that read out stock data. In the afternoon to promote wellness content managers would want the home page to be a "Get up and stretch" message.

All of this is possible to be done in a number of ways, but I wanted to make a simple traffic cop style HTTP Listener in Go that would forward (via a HTTP redirect 303 message) the browser to a desired target depending on the time of day. It's extremely lightweight as there is no proxy, just a redirect. With the Docker instance, you can run it from anywhere. Just configure the new default home page to this address, and it will do the hard work for you. 

## Getting Started

These instructions will get you a copy of the project up and running on your local machine or via Docker.

### Prerequisites

You can run this as a container or you can run it standalone. Either run the executable directly or compile it via Go. The network port is a configurable so if you want to change the default of port 8080, pass it in as below. For a container, this should not matter as you will expose port 80 (or whatever you want) on the node.

#### Runtime Options
There are three total flags available to be passed at runtime. Those are:

* -conf
	* A string value
	* default "config/config.json"
	* -conf=/path/to/config 
* -debug
	* A boolean value that turns on some basic debugging info (more for me less for you).
	* default false
	* -debug 
* -network
	* A string value
	* default ":8080"
	* Examples:
		* -network=":8080" // Listen on port 8080 for any IPv4/6 address we have (this is the default)
		* -network="127.0.0.1:3000" // Listen on the loopback, port 3000
		* -network="10.0.1.4:9080" // Listen on IP 10.0.1.4, port 9080

#### config.json
The `config/config.json` file feeds the magic into the binary. This is a standard JSON formatted file with a list of times (HHMM format on the 24hour clock) and a URL. Be sure that both are strings, the decoder will handle the data types. You are welcome to change the location of the config if you like, just pass a -conf=/path/to/config.json on execution.

Sample config:

	[
	    {
	        "time": "1400",
	        "url": "https://www.google.com/"
	    },
	    {
	        "time": "0000",
	        "url": "https://www.flower-mound.com/covid19"
	    },
	    {
	        "time": "0800",
	        "url": "https://www.marketwatch.com"
	    },
	    {
	        "time": "1100",
	        "url": "https://www.grubhub.com"
	    }
	]

Here are a few notes about the config file:

* All elements are strings. Time must be in HHMM on the 24hr clock.
* You can drop elements into your file unordered like above, but it can get confusing to read. The program will automatically sort the list, but for your own sanity, you should probably write the JSON config in order from 0000-2359. In this case, moving the first entry to the bottom would make it much easier to read.
* No duplicate times.
* A time value of less than 0 or greater than 2359 is not allowed.
* Encode URLs if they have funky characters in them.


### Installing

Provided you have the Go compiler installed, you can just:

```
go build
./WhereToNow &
```

If you want to run this as a service on your respective operating system, just build a service around it. The above works fine for backgrounding the process in Linux. You will need to build your own process monitoring to make sure the OS didn't kill the process.

Building it as a service will come in hand if you want it to run on startup.

For Docker, use this:

```
docker build . 
docker run -p 80:8080 -d <imageid>
```

Where you bind whatever port you want to access from a user's browser (the part before the `:`) to port 8080 which is the default port the app listens on.

OR using docker-compose:

```
docker-compose up -d
```

## Exit Codes

There are a number of reasons why XXX could exit. Each one will print out an error to STDOUT as well as give an exit code. Those codes are:

* 1: Error reading the config file. This would likely be due to passing in a file that does not exist.
* 2: JSON file parsing error. Check your formatting.
* 3: A time in the JSON file was larger than 2359, meaning you might not be on Earth.
* 4: A time in the JSON file was smaller than 0, meaning you might not be in this dimension.
* 5: There's a duplicate time entry in your JSON.

Should be easy to solve.

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Go](https://golang.org/) - Underlying code.
* [Alpine Linux](https://www.alpinelinux.org/) - Alpine Linux for the container
* [Docker](https://docker.com/) - Docker

## Authors

* **Branden Williams** - [BrandoLabs](https://github.com/captbrando)


## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

## Acknowledgments

* [Mike Van Sickle](https://www.pluralsight.com/authors/mike-vansickle) for having a ton of great courses on PluralSight.
* [Matt Springfield](https://www.12feet.com) for dealing with my incessant questions and answering them with stamina and grace. And confusion.
