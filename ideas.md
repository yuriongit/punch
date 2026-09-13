# Ideas

## Punch Web-Application

1. Possibly add /health/ready for establishing whether
Redis, and other external dependencies are
currently healthy ```r.GET("/ready", health.Ready)```.

## Punch CLI

### Features

1. There should be a dedicated API
strictly for checking whether the
`punch.json` in the provided directory
actually exists. If it exists, the API will 
then go on to validate the config.
This will be the Punch Validator API. I want punch 
to be called just like `docker-compose up`: If there's 
no configuration file in that current directory, 
it'll exit with telling you that there's no 
configuration found. I want the same functionality.
3. When a test is run, it will start with initializing
logs and once its done with those, the load-test will
start. Although, right before it starts, I want there
to be line breaks right after the initializing logs
and once the load-test logs end too for clarity.

### Commands

_punch run_
  - ```punch run {directory} | punch run ./lilify```
  - ```punch run ./ --store-logs | punch run ./ --sl```

_punch logs_
  - ```punch logs {test-name} | punch logs lilify```

### Configuration (punch.json)
  
4. Include different protocols within 
configuration file: ws, sockio, gql
    - Figure out the the best shape for the 
    configuration which will give the ability to
    include these other protocols with ease.

### Extras

  1. Also, i'd actually like this to be a
  real tool that's installed into /bin. 
  2. I'd love to introduce containerization
  to the CLI tool with Docker :)

### Notes
Make use of `switch` statements!
