package exitcode

import "os"

const Success = 0
const Failure = 1

var Exit = func(code int) { os.Exit(code) }
