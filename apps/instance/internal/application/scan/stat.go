package scan

import "os"

// statFile follows symlinks, so a library can link to files on another drive.
var statFile = os.Stat
