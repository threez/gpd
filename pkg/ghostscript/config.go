// Copyright (c) 2018 Vincent Landgraf

package ghostscript

// Config for the ghostscript executable
type Config struct {
	Command   string
	Arguments []string
}

// DefaultConfig the default configuration uses /usr/bin/env and the first
// argument gs to identify the location of the ghostscript executable.
// If not applicable and if the path will be specified directly, the first
// argument needs to be removed as well
var DefaultConfig = Config{
	"/usr/bin/env",
	[]string{"gs", "-dQUIET", "-dBATCH", "-dSAFER", "-sOutputFile=-"},
}
