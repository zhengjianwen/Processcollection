package funcs

import "os"

func Writefile(data []byte) error {

	file6, error := os.OpenFile("./1.txt", os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
	}
	file6.Write(data)
	file6.Close()
	return nil
}
