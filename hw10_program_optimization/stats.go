package hw10programoptimization

import (
	"bufio"
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

		slitted := strings.SplitN(user.Email, "@", 2)
		if len(slitted) != 2 {
			continue
		}
		key := strings.ToLower(slitted[1])
		val, ok := emails[key]

		if !ok {
			matched := strings.HasSuffix(strings.ToLower(user.Email), "@"+key)

			if !matched {
				continue
			}
		}

		emails[key] = val + 1
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return
}
