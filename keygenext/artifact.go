package keygenext

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"
)

// Artifact represents a Keygen artifact object.
type Artifact struct {
	ID       string    `json:"-"`
	Type     string    `json:"-"`
	Key      string    `json:"key"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Location string    `json:"-"`
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

	info, err := file.Stat()
	if err != nil {
		return err
	}

	size := info.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	req, err := http.NewRequest("PUT", a.Location, bytes.NewReader(buffer))
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to upload to storage provider")
	}

	return nil
}
