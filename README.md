# healthz

Inspired by Kelsey Hightower's [talk at Monitorama](https://vimeo.com/173610242), the code accompanying the talk (https://github.com/kelseyhightower/app-healthz/) and the healthz implementation in Kubernetes(https://github.com/kubernetes/kubernetes/tree/master/pkg/healthz).

# Issuing requests

The following examples assume you have a healthz checker using the default `/healthz` prefix:


```
$ curl -i http://127.0.0.1:8080/healthz
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 25 Jun 2017 23:33:01 GMT
Content-Length: 122

{
    "app": {
        "metadata": {
            "key1": "value1"
        }
    },
    "app2": {
        "metadata": {
            "key2": "a different value"
        }
    }
}
```

Query a single health check:

```
$ curl -i http://127.0.0.1:8080/healthz/app
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 25 Jun 2017 23:31:18 GMT
Content-Length: 56

{
    "app": {
        "metadata": {
            "key1": "value1"
        }
    }
}
```

Query a non-existent check:

```
$ curl -i http://127.0.0.1:8080/healthz/app3
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sun, 25 Jun 2017 23:32:04 GMT
Content-Length: 22

No such healthz check
```

If you're only interested in the overall status, not the metadata or logs, use the `HEAD` verb:

```
$ curl -i -I http://127.0.0.1:8080/healthz
HTTP/1.1 200 OK
Date: Sun, 25 Jun 2017 23:34:23 GMT
Content-Type: text/plain; charset=utf-8
```
