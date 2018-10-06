# A (really) quick and dirty wiremock agent

This application manages the execution of "wiremock-standalone-2.19.0.jar" in the path relative to the current binary
It will restart this process when it detects changes made the directory provided as the first argument. 

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
* Include CLI Output

## Example

    # will execute "java -jar ./wiremock-standalone-2.19.0.jar --port 8082 --global-response-templating"
    # and restart it as changes happen in mappings/
    ./wire-mock-agent mappings/