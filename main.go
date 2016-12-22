package main

import "bufio"
import "fmt"
import "github.com/mitchellh/go-homedir"
import "github.com/urfave/cli"
import "os"
import "path/filepath"

const defaultStoreFile = "~/.pin"

func main() {
	app := cli.NewApp()
	app.Name = "pin"
	app.Usage = "pin directory."
	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			list()
			return nil
		}
		err := appendList(c.Args().First())
		if err != nil {
			return err
		}
		return nil
	}
	app.Run(os.Args)
}

func expand(path string) (string, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func open(path string) (*os.File, error) {
	path, err := expand(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func list() error {
	file, err := open(defaultStoreFile)
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func appendList(path string) error {
	path, err := expand(path)
	if err != nil {
		return err
	}
	if _, err = os.Stat(path); err != nil {
		return err
	}
	file, err := open(defaultStoreFile)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		v := scanner.Text()
		if v == path {
			// duplicate.
			return nil
		}
	}
	_, err = file.WriteString(path + "\n")
	if err != nil {
		return err
	}
	return nil
}
