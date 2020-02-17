package gapt

import "strings"

var root = &dnode{dnType: dnTypeDirectory}

func Import(filename string, data []byte) {
	segments := strings.Split(filename,"/")
	for _, segment := range segments{

	}
}
