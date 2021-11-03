package keygenext

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"time"
)

// Artifact represents a Keygen artifact object.
type Artifact struct {
	ID            string    `json:"-"`
	Type          string    `json:"-"`
	Key           string    `json:"key"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Location      string    `json:"-"`
	ContentLength int64     `json:"-"`
	ContentType   string    `json:"-"`
}

func (a *Artifact) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Artifact) SetType(t string) error {
	a.Type = t
	return nil
}

func (a *Artifact) SetData(to func(target interface{}) error) error {
	return to(a)
}

func (a *Artifact) Upload(file *os.File) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", a.Location, bufio.NewReader(file))
	if err != nil {
		return err
	}

	// This must be set otherwise the Go http package sends a Transfer-Encoding
	// header, which S3 does not support.
	req.ContentLength = a.ContentLength

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to upload to storage provider")
	}

	return nil
}
