# gobundledhttp
Provides convenience functions to generate an *http.Client with bundled CA certificates.

To update certificates (needed only if the cacerts source changes):
```
cd updatecerts
go run update.go
```

