package alien_invastion

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StreamParser struct {
}

func (s *StreamParser) ParseFile(filepath string) (ret *GameMap, errors []error) {
	// read file, one line by one line
	// parse line
	file, err := os.Open(filepath)
	if err != nil {
		return nil, append(errors, fmt.Errorf("failed to open file: %s", err))
	}
	defer func() {
		_ = file.Close()
	}()
	scanner := bufio.NewScanner(file)
	ret, errors = s.parseScannerResult(ret, scanner, errors)
	return
}

func (s *StreamParser) ParseString(str string) (ret *GameMap, errors []error) {
	// read file, one line by one line
	// parse line
	scanner := bufio.NewScanner(strings.NewReader(str))
	ret, errors = s.parseScannerResult(ret, scanner, errors)
	return
}

func (s *StreamParser) parseScannerResult(ret *GameMap, scanner *bufio.Scanner, errors []error) (*GameMap, []error) {
	ret = NewGameMap()
	for scanner.Scan() {
		line := scanner.Text()
		errs := s.parseSingleLine(line, ret)
		if errs != nil && len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	return ret, errors
}

func (s *StreamParser) parseSingleLine(line string, ret *GameMap) (errors []error) {
	// Line will looks like :
	// Foo north=Bar west=Baz south=Qu-ux
	// Bar south=Foo west=Bee
	elems := strings.Split(line, " ")
	for i, elem := range elems {
		if i == 0 {
			ret.UpsertCity(strings.Trim(elem, " "))
			continue
		}

		if strings.Contains(elem, "=") {
			//Regex might be another way but it's overkill
			dirCityPair := strings.Split(elem, "=")
			if len(dirCityPair) != 2 {
				continue
			}
			err := ret.UpdateCityWithNeighborhood(elems[0], DirectionFromString(dirCityPair[0]), strings.Trim(dirCityPair[1], " "))
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return
}
