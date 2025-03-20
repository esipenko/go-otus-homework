package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getDomains(r, domain)
}

func getDomains(r io.Reader, domain string) (emails DomainStat, err error) {
	scanner := bufio.NewScanner(r)
	emails = make(DomainStat)

	for scanner.Scan() {
		line := scanner.Bytes()

		var user User
		if err = json.Unmarshal(line, &user); err != nil {
			return
		}

		if !strings.HasSuffix(user.Email, "."+domain) {
			continue
		}

		key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		val, ok := emails[key]

		if !ok {
			matched := strings.HasSuffix(strings.ToLower(user.Email), "@"+key)

			if matched {
				emails[key] = 1
			}
		} else {
			emails[key] = val + 1
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return
}
