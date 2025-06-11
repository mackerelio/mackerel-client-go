module github.com/mackerelio/mackerel-client-go

go 1.21.0

require github.com/kylelemons/godebug v1.1.0

retract [v0.27.0, v0.37.0] // Path encoding for FindHostByCustomIdentifier is incorrect.

retract v0.33.0 // API endpoint for DeleteGraphDef has changed.
