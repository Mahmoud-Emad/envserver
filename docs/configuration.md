# Project Configuration

This guide will help you set up the config.toml configuration file for your project.

## Configuration File (`config.toml`)

Create a new file named `config.toml` in the root directory and populate it with the following contents:

```toml
[database]
host = "<database_host>"
port = <database_port>
user = "<database_user>"
password = "<database_password>"
name = "<database_name>"

[server]
host = <server_host>
port = <server_port>
jwtSecretKey = <jwtSecretKey>
shutdownTimeout = <shutdownTimeout>
```

Replace the placeholder values `<database_host>`, `<database_port>`, `<database_user>`, `<database_password>`, `<database_name>`, and `<server_port>` with the appropriate values see the [config.toml.template](../config.toml.template) .

- `<database_host>`         : Replace with the host address of your database server (e.g., "localhost").
- `<database_port>`         : Replace with the port number of your database server (e.g., 5432).
- `<database_user>`         : Replace with the username for accessing your database (e.g., "postgres").
- `<database_password>`     : Replace with the password for accessing your database (e.g., "postgres").
- `<database_name>`         : Replace with the name of your database (e.g., "postgres").
- `<server_host>`           : Replace with the host address of your server (e.g., "localhost").
- `<server_port>`           : Replace with the desired port number for your server (e.g., 8080).
- `<jwtSecretKey?>`         : Replace with simple text used as secret key for the jwt token.
- `<shutdownTimeout?>`?     : To shut down the server in time, replace the value with a simple number, it's optional.

Make sure to save the config.toml file after updating the values.
