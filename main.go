package main

import "bufio"
import "fmt"
import "github.com/mitchellh/go-homedir"
import "github.com/urfave/cli"
import "io"
import "log"
import "os"
import "os/user"
import "path/filepath"
import "sort"

var out = "~/.pin"

func init() {
	u, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	out = filepath.Join(u.HomeDir, ".pin")
}

type list []string

func (l list) add(s string) list {
	for _, v := range l {
		if v == s {
			return l
		}
	}
	return append(l, s)
}

func main() {
	app := cli.NewApp()
	app.Name = "pin"
	app.Version = "1.0.0"
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	if !c.Args().Present() {
		show()
		return nil
	}
	err := add(c.Args().First())
	if err != nil {
		return err
	}
	return nil
}

func add(path string) error {
	path, err := format(path)
	if err != nil {
		return nil
	}
	l, err := read()
	if err != nil {
		return nil
	}
	l = l.add(path)
	w, err := open(out)
	defer w.Close()
	if err != nil {
		return err
	}
	sort.Strings(l)
	err = write(w, l)
	if err != nil {
		return err
	}
	return nil
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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func show() error {
	l, err := read()
	if err != nil {
		return nil
	}
	for _, v := range l {
		fmt.Println(v)
	}
	return nil
}

func format(path string) (string, error) {
	path, err := expand(path)
	if err != nil {
		return "", err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		path = filepath.Clean(path) + string(filepath.Separator)
	}
	if _, err = os.Stat(path); err != nil {
		return "", err
	}
	return path, err
}

func read() (l list, err error) {
	l = make(list, 0)
	file, err := open(out)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l = l.add(scanner.Text())
	}
	return
}

func write(w io.Writer, s []string) error {
	for _, v := range s {
		_, err := fmt.Fprintln(w, v)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
