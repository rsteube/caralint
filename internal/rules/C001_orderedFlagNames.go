package rules

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/owenrumney/go-sarif/sarif"
)

type Rule struct {
	Id         string
	Desciption string
	Check      func(path string) (*Result, error)
}

type Result struct {
	Message string
	Region  *sarif.Region
}

var C001 = Rule{
	Id:         "C001",
	Desciption: "flags should be ordered",
	Check: func(path string) (*Result, error) {
		file, err := os.Open(path)
		if err != nil {
			// TODO add as another result?
			//return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		r := regexp.MustCompile(`\.Flags\(\)\.(Bool|String)[^(]*\("(?P<name>[^"]+)"`)

		previous := ""
		line := 0
		for scanner.Scan() {
			if !r.MatchString(scanner.Text()) {
				previous = ""
			} else {
				matches := r.FindStringSubmatch(scanner.Text())
				current := matches[2]
				if previous != "" && previous > current {
					return &Result{
						Message: fmt.Sprintf("flag '%v' should be before '%v'", current, previous),
						Region:  sarif.NewSimpleRegion(10, 11),
					}, nil
				}
				previous = current
			}
			line += 1
		}
		return nil, nil
	},
}
