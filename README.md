# Reverse proxy based on a header's value 

This middleware can be used to reverse proxy a request based on a headers value, e.g. based on the value of `X-Forwarded-User` from something like [thomseddon/traefik-forward-auth](https://github.com/thomseddon/traefik-forward-auth).

## Configuration

### Static: 
Production:

```
[pilot]
  token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

[experimental]
  [experimental.plugins]
    [experimental.plugins.my-plugin-name]
      moduleName = "github.com/vidosits/header-pattern-proxy"
      version = "v1.1.0"
```

or if you're using [devMode](https://doc.traefik.io/traefik-pilot/plugins/plugin-dev/#developer-mode):

```
[pilot]
  token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  
[experimental.devPlugin]
  goPath = "/plugins/go"
  moduleName = "github.com/vidosits/header-pattern-proxy"
  # Plugin will be loaded from '/plugins/go/src/github.com/vidosits/header-pattern-proxy'
```

### Dynamic:

Production:

```
[http]
  [http.middlewares]
    [http.middlewares.my-middleware-name.plugin.my-plugin-name]
      header  = "X-Forwarded-User"
      [http.middlewares.my-middleware-name.plugin.my.plugin-name.mapping]
        "user1@gmail.com" = "http://nginx"
        "user2@gmail.com" = "http://whoami"

  [http.routers]
    [http.routers.my-router-name]
      entryPoints = ["websecure"]
      rule = "Host(`my-service-name.domain.tld`)"
      middlewares = ["traefik-forward-auth@docker", "my-middleware-name@file"]
      
      # if no matches are found this is the service that we forward the request to
      service = "noop@internal"
      
      [http.routers.my-router-name.tls]
        certResolver = "letsencrypt"
```

or if you're using [devMode](https://doc.traefik.io/traefik-pilot/plugins/plugin-dev/#developer-mode):

```
[http]
  [http.middlewares]
    [http.middlewares.my-middleware-name.plugin.dev]
      header  = "X-Forwarded-User"
      [http.middlewares.my-middleware-name.plugin.dev.mapping]
        "user1@gmail.com" = "http://nginx"
        "user2@gmail.com" = "http://whoami"

  [http.routers]
    [http.routers.my-router-name]
      entryPoints = ["websecure"]
      rule = "Host(`my-service-name.domain.tld`)"
      middlewares = ["traefik-forward-auth@docker", "my-middleware-name@file"]
      
      # if no matches are found this is the service that we forward the request to
      service = "noop@internal"
      
      [http.routers.my-router-name.tls]
        certResolver = "letsencrypt"
```

## License
This software is released under the Apache 2.0 License.