mackerel-client-go
==================

mackerel-client-go is a Go client library for [mackerel.io API](http://help-ja.mackerel.io/entry/spec/api/v0).

# Usage

```go
import mkr "github.com/mackerelio/mackerel-client-go"
```

```go
client = mkr.NewClient("<Put your API key>")

hosts, err := client.FindHosts(&mkr.FindHostsParam{
        Service: "My-Service",
        Roles: []string{"proxy"},
        Statuses: []string{"working"}
})

err := client.PostServiceMetricValues("My-Service", []*mkr.ServiceMetricValue{
        &mkr.ServiceMetricValue{
              Name:  "proxy.access_log.latency",
              Time:  123456789,
              Value: 500,
        },
})
```
