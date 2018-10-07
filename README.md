# A (really) quick and dirty process agent

This application will conditionally restart processes when specified directories have changed.

## Configuration

In order to properly setup the agent you will need to define:

1. a process alias (label)
2. a directory to monitor
3. the name of the process to manage
4. the arguments to supply the process
5. supply a config.yml where the above is defined. By default the mock agent assumes a `config.json` exists in the same
dir as the binary

### Example Configuration

        {
              "alias" : "wiremock",
              "monitor": "mappings",
              "command" : "java",
              "arguments": [
                "-jar", "wiremock-standalone-2.19.0.jar", "--port", "8082", "--global-response-templating"
              ]
        }

## Running

    ./wiremock-agent
    
## Building 

### For yourself

     go build

### For your linux peeps (penguins?)

    # 32 bit  
    GOOS=linux GOARCH=386 go build
    # 64 bit
    GOOS=linux GOARCH=amd64 go build

### For your Apple friends

    # 32 bit  
    GOOS=darwin GOARCH=386 go build
    # 64 bit
    GOOS=darwin GOARCH=amd64 go build

### For your Microsoft Friends

    # 32 bit  
    GOOS=windows GOARCH=386 go build
    # 64 bit
    GOOS=windows GOARCH=amd64 go build

  
## TODO

* Error handling is essentially absent. Will be addressed later  
* Rename App
* More Functionality?
* Better App Name
