# go-consul-client
A consul library that can load json into consul KV and read a consul KV namespace into a config cache.  
This repo provides both an application capable of importing JSON into consul key values as well as a library
for using consul as a configuration store for golang applications.

The library provides a configuration cache as well as some helper functions to convert consul configuration into
strings, ints, booleans, and durations.

## Usage

### Application
The easiest way to use this application is as a docker container which has been made available through [docker hub](https://hub.docker.com/r/divideandconquer/go-consul-client/):

```bash
docker run -v /path/to/json/file:/config.json divideandconquer/go-consul-client  -file /config.json -namespace testing/fun -consul 172.17.8.101:8500
```

You can also build this application yourself with the provide build script in `build/build.sh` and run the application binary directly.

### Library
You can import this library into you golang application and then use it to access your consul configuration:

```golang

import "github.com/divideandconquer/go-consul-client/src/client"


func main() {
	consulAddress := "172.17.8.101:8500"
	environment := "dev"
	appNamespace := environment + "/" + "my-app"

	// create a cached loader
	conf, err := client.NewCachedLoader(appNamespace, consulAddress)
	if err != nil {
		panic(err)
	}

	// initialize the cache
	err = conf.Initialize()
	if err != nil {
		panic(err)
	}

	//fetch data from the cache
	myConfigString := conf.MustGetString("config_key")
	myConfigBool := conf.MustGetBool("config_key")
	myConfigInt := conf.MustGetInt("config_key")
	myConfigDuration := conf.MustGetDuration("config_key")

	...
}

```


