/*
 *  properties.go
 *  Created on 08.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package build

import (
	"bufio"
	"os"
	"strings"
)

type Properties map[string]string

func ReadProperties(filename string) (Properties, error) {

	props := Properties{}

	if len(filename) == 0 {
		return props, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				props[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return props, nil
}
