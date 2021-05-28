mackerel-client-go
==================

[![Build Status](https://github.com/mackerelio/mackerel-client-go/workflows/Build/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/mackerelio/mackerel-client-go/badge.svg?branch=master)][coveralls]
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/mackerelio/mackerel-client-go)][pkg.go.dev]

[actions]: https://github.com/mackerelio/mackerel-client-go/actions?workflow=Build
[coveralls]: https://coveralls.io/r/mackerelio/mackerel-client-go?branch=master
[pkg.go.dev]: https://pkg.go.dev/github.com/mackerelio/mackerel-client-go

mackerel-client-go is a Go client library for [mackerel.io API](https://mackerel.io/api-docs/).

# Usage

```go
import "github.com/mackerelio/mackerel-client-go"
```

```go
client := mackerel.NewClient("<Put your API key>")

hosts, err := client.FindHosts(&mackerel.FindHostsParam{
        Service: "My-Service",
        Roles: []string{"proxy"},
        Statuses: []string{mackerel.HostStatusWorking},
})

err := client.PostServiceMetricValues("My-Service", []*mackerel.MetricValue{
        &mackerel.MetricValue{
              Name:  "proxy.access_log.latency",
              Time:  123456789,
              Value: 500,
        },
})
```

# CAUTION

Now, mackerel-client-go is an ALPHA version. In the future release, it may change it's interface.

# CONTRIBUTION

1. Fork ([https://github.com/mackerelio/mackerel-client-go/fork](https://github.com/mackerelio/mackerel-client-go/fork))
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
