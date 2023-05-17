
[![Go Reference](https://pkg.go.dev/badge/github.com/otyang/go-starter.svg)](https://pkg.go.dev/github.com/otyang/go-starter)


[![Go Report Card](https://goreportcard.com/badge/github.com/otyang/go-starter)](https://goreportcard.com/report/github.com/otyang/go-starter)

# Go Starter Kit (Boilerplate) 

This starter kit is designed to get you up and running with a project structure optimized for developing services (rest, cmd etc) in Go. It promotes the best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID) and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  It encourages writing clean and idiomatic Go code.



## Technologies Used
- Language:		[Golang](https://golang.org) 
- Docker:	[Docker](https://www.docker.com/)
- Logging: [Golang slog](golang.org/x/exp/slog) 
- Routing: [gofiber](https://github.com/gofiber/fiber)
- Tests:  [Testify](https://github.com/stretchr/testify)
- Database Access:	[Bun](https://github.com/uptrace/bun)
- Cache Access (Redis or Memory): [Redis/Rueids](https://github.com/redis/rueidis) OR [Patrick gocache](https://github.com/patrickmn/go-cache)
- Data validation: [gookit - validate](https://github.com/gookit/validate) OR [go-playground](https://github.com/go-playground/validator)


## User Guide
* [Technologies Used](#technologies-used)
* [User Guide](#user-guide)
* [Installations](#installations)
    * [Clone](#clone)  
    * [Via Go Run](#via-go-run)  
    * [Via Docker](#via-docker)  
* [Project Layout](#Project-layout)
* [Managing configurations](#managing-configurations)
* [Documentation](#documentation)
* [Running Tests](#running-tests)
* [Author](#author)
* [License](#license)


## Installations
To run this project, please refer to the below listed guide: 

### Clone
To clone this project to your local machine do this:
```bash
    git clone https://github.com/otyang/icd-10.git 
```

### Via Go Run
To run this project via go run:
- First Clone the git repo
- Then cd into the directory
- Next, do this
```bash
    git clone https://github.com/otyang/go-starter.git
    cd go-starter/
    go run cmd/zample/main.go --configFile="cmd/zample/config.toml"
``` 

### Via Docker
To run this project via Docker:
- First Clone the git repo
- Then cd into the directory 
- Next, build the docker image by running
- Finally run the docker image. 

```
    git clone https://github.com/otyang/go-starter.git
    cd go-starter/
    docker build -t go-starter-project -f Dockerfile.zample . 
    docker run -p 4000:4000 -t go-starter-project
```




## Project Layout
The following layout is used:
```
.
├── cmd/                main applications for this project
│   ├── cli             cli service entry (command line app)
│   └── zample          zample service entry (rest app)
├── config              configuration files/templates & default config
├── internal/           private application and library code
│   ├── event           entity definitions and domain logic for events
│   ├── middleware      middlewares libraries / codes
│   └── zample          application library for zample service
└── pkg/                public library code
    ├── _example        folder showing example to use pkg
    ├── config          configuration library
    ├── datastore       helpers for working with database
    ├── logger          structured and context-aware logger
    ├── pagination      library to handle pagination
    ├── response        handles http response, errors and request
    ├── utils           library for miscellanous things
    └── validators      helpers to efficiently handle validation 
```


The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `zample` directory contains the application logic related with the zample features. 

Within each feature package, code are organized in layers (API, entity, repository, handlers-for-http, handlers-for-events), following the dependency guidelines as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).


## Managing Configurations

The application configuration is represented in `config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path 
to the configuration file is specified via the `-configFile` command line argument which defaults to 
`./config/config.toml`. 

Configurations specified in environment variables should be named in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The  configuration files for each service are found under their respect `./cmd/service` entry folder.
for example the config file for `./cmd/zample` can be seen at `./cmd/zample/config.toml`



## Documentation
-   [Overall Documentation](/README.md)
-   [Go Documentation](http://godoc.org/github.com/otyang/go-starter)
-   [Pkg Documentation](/pkg/README.md)
-   [Rest Endpoint Documentation](/cmd/zample/README.md)

## Running Tests
To run tests, run the following command
```bash
go test -v ./...
```

To run integration tests
```bash
go test ./... -tags=integration ./...
```

To run coverage tests and generate the coverage report
```bash
go test ./... -v -coverpkg=./...
```

To run all test at a go use the Makefile command
```bash
make test
```


## Author
- O Yang [@otyang](https://www.github.com/otyang) 


## License 
This project is licensed & available for use under the terms of the [MIT license](/LICENSE).
