module tengine/core

go 1.25.6

require (
	golang.org/x/term v0.22.0
	tengine/render v0.0.0
	tengine/structs v0.0.0
)

require golang.org/x/sys v0.22.0 // indirect

replace (
	tengine/render => ../render
	tengine/structs => ../structs
)
