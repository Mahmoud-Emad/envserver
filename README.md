# Envserver

## Project Description

envserver is a server application that allows users to store and manage environment keys for their projects. It provides a centralized platform for securely storing and accessing environment variables, similar to the functionality provided by tools like `flagsmeth`. With envserver, users can easily store and retrieve environment keys for their projects, enhancing their development workflow.

The server application includes a database with two main tables: "User" and "Project". The "User" table handles user registration and login functionality, with each user being assigned a unique token for authentication. The "Project" table is used to manage projects, their teams, and their associated environment variables.

## CLI Tool

envserver provides a command-line interface (CLI) tool to facilitate key management for users. The CLI tool offers several commands to interact with the server and manage environment keys effectively. The available commands are:

- pull: Pulls the latest changes from the server and creates or updates the local Config file.
- push: Pushes the local changes to the server, updating the environment keys.
- add: Adds new environment keys to the local Config file.
- commit: Commits the changes to the local Config file, providing a commit message. The commit message can be customized and will be updated if conflicts occur.

### Please note that, all of these commands are still under implementation

## Project Config

For detailed information on configuring the envserver project, refer to the [Project Config](./docs/Config.md) document. This document provides instructions on setting up the config.toml Config file, which includes important settings such as database connection details and server port.

## Makefile Commands

- `build`: This command builds the project by compiling the `cmd/envserver.go` file.

```sh
make build
```

- `run`: This command first builds the project by invoking the build command, and then runs the executable with the specified config file using `./envserver -config ${config_file}`, Running make run config=config.toml will execute the run command with the specified config file.

```sh
make run config='<path_to_config>' # it will exec the make-build also
```

- `test`: This command first builds the project by invoking the build command, and then it runs all the tests in the project using the go test command.

```sh
make test
```

- `clean`: This command will remove the executable file `./envserver`.

```sh
make clean # remove the build file.
```

## Contributing

If you would like to contribute to the envserver project, please refer to the [Contributing Guidelines](./docs/contributing.md) document. It outlines the steps to contribute, including guidelines for reporting issues, suggesting improvements, and submitting pull requests.
