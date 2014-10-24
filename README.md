mackerel-client-go
==================

mackerel-client-go is a Go client library for [mackerel.io API](http://help-ja.mackerel.io/entry/spec/api/v0).

# Usage

```go
import mkr "github.com/mackerelio/mackerel-client-go"
```

```go
client = mkr.NewClient("<Put your API key>")
hosts, err := client.FindHosts(&mkr.FindHostsParam{Service: "My-Service", Role: "proxy", Status: 'working'})
params := []*mkr.MetricParam{
        &mkr.MetricParam{name: 'proxy.access_log.latency', time: '1414146305', value: '250'},
}
ret, err := client.PostServiceMetric("My-Service", params)
```
