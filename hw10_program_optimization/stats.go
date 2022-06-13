package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	br := bufio.NewReader(r)

	result := make(DomainStat)

	exp, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	isEOF := false

	for {
		line, err := br.ReadBytes(byte('\n'))
		if err != nil {
			if errors.Is(err, io.EOF) {
				isEOF = true
			} else {
				return result, err
			}
		}

		var user User
		if err = easyjson.Unmarshal(line, &user); err != nil && !errors.Is(err, io.EOF) {
			return result, err
		}

		matched := exp.MatchString(user.Email)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}

		if isEOF {
			break
		}
	}

	return result, nil
}
