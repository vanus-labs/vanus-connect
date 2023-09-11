package internal

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/pkg/errors"
)

func event2Member(event *ce.Event) (*Member, error) {
	var member Member
	err := json.Unmarshal(event.Data(), &member)
	if err != nil {
		return nil, errors.Wrap(err, "member unmarshal error")
	}
	if err = member.Validate(); err != nil {
		return nil, errors.Wrap(err, "member invalid")
	}
	return &member, nil
}

func emailHash(email string) string {
	h := md5.New()
	io.WriteString(h, strings.ToLower(email))
	return fmt.Sprintf("%x", h.Sum(nil))

}
