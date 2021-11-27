## Krakend martian custom plugin

To add the custom modifier to KrakenD it’s only needed to create a file like this one in the KrakenD-CE repository `/cmd/krakend-ce`:

```
package main

import _ "github.com/tarsidi-danesh/martian-custom-plugin"
```

and just re-build with command `make build` or `make build_on_doker` if your machine is not using linux OS

### Sample Configuration

```json
{
"endpoints": [
    {
      "endpoint": "/tix-flight-master-discovery/airline/partners",
      "method": "GET",
      "output_encoding": "no-op",
      "backend": [
        {
          "url_pattern": "/tix-flight-master-discovery/airline/partners",
          "encoding": "no-op",
          "method": "GET",
          "sd": "static",
          "host": [
            "https://flight-master-discovery-be-svc.test-flight-cluster.tiket.com/"
          ],
          "disable_host_sanitize": true,
          "extra_config": {
            "github.com/devopsfaith/krakend-martian": {
              "fifo.Group": {
                "aggregateErrors": true,
                "modifiers": [
                  {
                    "header.MandatoryModifier": {
                      "headers": [
                        {
                          "name": "serviceId",
                          "value": {
                            "type": "STATIC",
                            "generator": "NONE",
                            "staticValue": "GATEWAY"
                          }
                        },
                        {
                          "name": "requestId",
                          "value": {
                            "type": "DYNAMIC",
                            "generator": "UUID",
                            "staticValue": null
                          }
                        },
                        {
                          "name": "channelId",
                          "value": {
                            "type": "STATIC",
                            "generator": "NONE",
                            "staticValue": "DESKTOP"
                          }
                        },
                        {
                          "name": "username",
                          "value": {
                            "type": "STATIC",
                            "generator": "NONE",
                            "staticValue": "GUEST"
                          }
                        },
                        {
                          "name": "storeId",
                          "value": {
                            "type": "STATIC",
                            "generator": "NONE",
                            "staticValue": "TIKETCOM"
                          }
                        }
                      ]
                    }
                  }
                ],
                "scope": [
                  "request"
                ]
              }
            }
          }
        }
      ]
    }
  ]
}
```