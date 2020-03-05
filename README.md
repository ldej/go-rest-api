# REST API example

## Contains
- A basic folder structure for a large monolithic REST API
- A couple of endpoints
- Protected and unprotected endpoints
- JWT
- Storage in a SQL database
- Separation of concerns
- Logging (go.uber.org/zap)
- Configuration for environments (github.com/spf13/viper)
- Clean Architecture
- 12 Factor

## TODO
- Automatic creation of tables if they don't exist
- OAuth
- Profiling
- Swagger/OpenAPI
- Server testing

## Requires
- Go 1.14
- A postgres database running on port 5432
- A table for users

# Sources

https://github.com/snowzach/gorestapi
https://github.com/golang-standards/project-layout
https://www.reddit.com/r/golang/comments/a35xfv/how_would_you_structure_a_large_rest_api_service/
https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1#.ds38va3pp
