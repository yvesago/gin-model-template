package models

/*
  Shared functions for models tests
*/

import (
	"fmt"
	"os"
)

func deleteFile(file string) {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
