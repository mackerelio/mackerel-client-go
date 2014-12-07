mackerel-client-go [![Build Status](https://travis-ci.org/mackerelio/mackerel-client-go.svg?branch=master)](https://travis-ci.org/mackerelio/mackerel-client-go)
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

err := client.PostServiceMetricValues("My-Service", []*mkr.MetricValue{
        &mkr.MetricValue{
              Name:  "proxy.access_log.latency",
              Time:  123456789,
              Value: 500,
        },
})
```

# CONTRIBUTION

1. Fork ([https://github.com/mackerelio/gomkr/fork](https://github.com/mackerelio/mackerel-client-go/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

License
----------

Copyright 2014 Hatena Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
