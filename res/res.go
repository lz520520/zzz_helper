package res

import _ "embed"

//go:embed engines.yml
var Engines []byte

//go:embed agents.yml
var Agents []byte

//go:embed driver.yml
var Drivers []byte

//go:embed VERSION.txt
var VERSION string

//go:embed TITLE.txt
var TITLE string
