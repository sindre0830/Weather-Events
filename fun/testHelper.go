package fun

import "os"

func GoToRoot() (string, error) {
	var path string
	if _, err := os.Stat("./main.go"); err != nil {
		os.Chdir("./../../")
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	return path, nil
}
