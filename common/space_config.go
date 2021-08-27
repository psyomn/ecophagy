package common

import (
	"bytes"
)

// ReadSpaceConfig reads a text file, which each line has two
// continuous strings, delimited by one blank space so the configuration
// would look something like this:
//
// blahblah 1234
// uhoh 3214
//
// This is used by `randparty`, and the reason I'm not using a json
// file for configuration is because I wanted something that is more
// user friendly (this was being used by some non technical people).
// I might revisit, remove this, and use the CSV encoder instead.
func ReadSpaceConfig(filename string) (map[string]string, error) {
	ret := make(map[string]string)

	bs, err := FileToBytes(filename)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(bs, []byte{byte(Newline)})

	for i := range lines {
		kv := bytes.Split(lines[i], []byte{' '})
		if len(kv) < 2 {
			break
		}
		ret[string(kv[0])] = string(kv[1])
	}

	return ret, nil
}
