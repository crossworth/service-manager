## Simple and lightweight service manager for simple microservices

*Simple* is the key, service manager created to be used inside docker containers
to "share" the endpoint of different services to a "gateway" or "parent" service.

There is no authentication or complex endpoints, the only thing the client has to do
is periodically call an endpoint with the name, endpoint and (optionality) a value.


The client can use `GET` or `POST`
```
http://service-manager:8080/register?service=MyServiceName&endpoint=http://10.0.0.155:3000&value=OK
```

You can get the registered services on the endpoint `/services`.

That's it.

If you really wants, thats webhook support.
You can provide a comma separated list of url's to be called when an service changes.


Example `docker-compose.yml`
```yml
version: "3.7"

services:
    service-manager:
        restart: unless-stopped
        image: crossworth/service-manager
        environment:
          - CHANGES_WEBHOOK=http://my-gateway/check-services,http://auth-api:9000/check-services
```


There is a simple client as well.

```go

import (
 	"log"
	"github.com/crossworth/service-manager/client"
)

const service = "My-Service-Name"

func init() {
	if os.Getenv("SERVICE_MANAGER") != "" {
		log.Println("ServiceManager notifications enable")

		localIp, err := client.GetLocalIP()
		if err != nil {
			log.Fatal("could not determine local ip", err)
		}


		var opts client.Options
		opts.Name = service
		opts.ManagerEndpoint = os.Getenv("SERVICE_MANAGER")
		opts.Endpoint = localIp + ":8080"
		opts.Value = func() string {
			return "my-value"
		}

		client.Register(opts)
	}
}
```