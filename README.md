# GoApi122

With the introduction of enhanced routing patterns in Go 1.22, there’s no longer a need to use the external mux package.

- https://tip.golang.org/doc/go1.22#enhanced_routing_patterns
- https://go.dev/blog/routing-enhancements

```sh {"id":"01J7T85S9WB6MD6ZZHG75Y8X7Y"}
go run *.go
```

```sh {"id":"01J7T85S9YZ574NZM1E203K9ZC"}
# health -> OK
curl http://localhost:8080/health

# access unknown endpoint -> 404 page not found
curl http://localhost:8080/user/1  

# user id 100 set name to "John Doe"
curl http://localhost:8080/api/v1/user/name -H "Authorization: Bearer 100" -H "Content-Type: application/json" -X POST -d '{"name": "John Doe"}'

# user id 100 get name -> "John Doe"
curl http://localhost:8080/api/v1/user/name -H "Authorization: Bearer 100"

# user id 200 set name to "Jane Doe"
curl http://localhost:8080/api/v1/user/name -H "Authorization: Bearer 200" -H "Content-Type: application/json" -X POST -d '{"name": "Jane Doe"}'

# user id 200 get name -> "Jane Doe"
curl http://localhost:8080/api/v1/user/name -H "Authorization: Bearer 200"

# user id 100 try to access name of user 200 with GET /api/v1/user/id/{userID} that need super user role -> Unauthorized
curl http://localhost:8080/api/v1/user/id/200 -H "Authorization: Bearer 100"

# super user id 1 try to access name of user 200 with GET /api/v1/user/id/{userID} -> "Jane Doe"
curl http://localhost:8080/api/v1/user/id/200 -H "Authorization: Bearer 1"

# test api/v2 user id 100 get name -> "John Doe"
curl http://localhost:8080/api/v2/user/name -H "Authorization: Bearer 100"

curl http://localhost:8080/api/v2/version -H "Authorization: Bearer 100"

```

## Mux Patterns

Ref: https://pkg.go.dev/net/http@master#hdr-Patterns-ServeMux

In general, a pattern looks like `[METHOD ][HOST]/[PATH]`

Some examples:

- `/index.html` matches the path "/index.html" for any host and method.
- `GET /static/` matches a GET request whose path begins with "/static/".
- `example.com/` matches any request to the host "example.com".
- `example.com/{$}` matches requests with host "example.com" and path "/".
- `/b/{bucket}/o/{objectname...}` matches paths whose first segment is "b" and whose third segment is "o". The name "bucket" - denotes the second segment and "objectname" denotes the remainder of the path.

Note that:

- A pattern with no method matches every method.
- A pattern with the method GET matches both GET and HEAD requests.
- Otherwise, the method must match exactly.
- A pattern with no host matches every host.
