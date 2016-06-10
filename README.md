# gobundledhttp
Provides convenience functions to generate an *http.Client with bundled CA certificates.

Available Methods:
```go
sslclient := gobundledhttp.NewClient()         // Normal SSL-capable client
nosslclient := gobundledhttp.InsecureClient()  // Disables SSL certificate checking


context := gobundledhttp.CtxBundled()          // CtxBundled returns an oauth2 context with bundled http client
myX509pool := gobundledhttp.GetPool()          // Get the default cert pool to build your own client
```

To update certificates (needed only if the cacerts source changes):
```
cd updatecerts
go run update.go
```

