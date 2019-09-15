# Marvel API martian querystring modifier


This package is a [martian](https://github.com/google/martian) modifier for authentication against the Marvel comics API.


The procedure extracted from the Marvel API documentation:
https://developer.marvel.com/documentation/authorization

```
Authentication for Server-Side Applications
Server-side applications must pass two parameters in addition to the apikey parameter:
ts - a timestamp (or other long string which can change on a request-by-request basis)
hash - a md5 digest of the ts parameter, your private key and your public key (e.g. md5(ts+privateKey+publicKey)
For example, a user with a public key of "1234" and a private key of "abcd" could construct a valid call as follows: http://gateway.marvel.com/v1/public/comics?ts=1&apikey=1234&hash=ffd275c5130566a2916217b101f26150 (the hash value is the md5 digest of 1abcd1234)
```

## Adding the custom modifier in KrakenD

To add the custom modifier to KrakenD it's only needed to create a file like this one in the root of the [KrakenD-CE](https://github.com/devopsfaith/krakend-ce) repository:

*marvel.go:*
```
package main

import _ "github.com/taik0/marvelapi-martian_querystring"
```
and just follow the regular build steps (`make prepare` and `make`).

To configure KrakenD I needed to add an `extra_config` using the krakend-martian `Namespace` and the custom configuration for the modifier.
Check also the `whitelist` and `mapping` sections where I set a filter of some fields.

#### KrakenD configuration example
```
"endpoints": [
    {
        "backend": [
            {
                "extra_config": {
                    "github.com/devopsfaith/krakend-martian": {
                        "fifo.Group": {
                            "aggregateErrors": true,
                            "modifiers": [
                                {
                                    "querystring.MarvelModifier": {
                                        "private": "private_api_key",
                                        "public": "public_api_key",
                                        "scope": [
                                            "request"
                                        ]
                                    }
                                }
                            ],
                            "scope": [
                                "request",
                                "response"
                            ]
                        }
                    }
                },
                "host": [
                    "http://gateway.marvel.com"
                ],
                "url_pattern": "/v1/public/characters?name={name}",
                "whitelist": [
                  "data.results",
                  "attributionHTML"
                ],
                "mapping": { "data": "characters"}
            }
        ],
        "endpoint": "/character/{name}",
        "method": "GET"
    }
  ],
  "extra_config": {
      "github_com/devopsfaith/krakend-gologging": {
          "level": "INFO",
          "prefix": "[KRAKEND]",
          "stdout": true,
          "syslog": false
      }
  },
  "name": "My lovely gateway",
  "port": 8080,
  "timeout": "3s",
  "version": 2
```


