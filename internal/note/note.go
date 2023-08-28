package note

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04"

type Note struct {
	Title      string
	Text       string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func New(title string, text string) *Note {
	return &Note{
		Title:      title,
		Text:       text,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}

func Load(path string) (*Note, error) {
	const op = "note.Load"

	f, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	chank := make([]byte, 1024)

	for {
		n, err := f.Read(chank)
		if n > 0 {
			buf.Write(chank)
		}
		if err == io.EOF {
			break
		}
	}

	n, err := Parse(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%s: error parsing: %w", op, err)
	}
	return n, nil
}

func (note *Note) Save(dirName string) error {
	const op = "note.Save"

	path := filepath.Join(dirName, note.Title)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: unable to create file: %w", op, err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("title: %s\n", note.Title)); err != nil {
		return fmt.Errorf("%s: unable to write title to file: %w", op, err)
	}

	if _, err := f.WriteString(fmt.Sprintf("created_at: %s\n", note.CreatedAt.Format(TimeFormat))); err != nil {
		return fmt.Errorf("%s: unable to write created_at to file: %w", op, err)
	}

	if _, err = f.WriteString(fmt.Sprintf("modified_at: %s\n", note.ModifiedAt.Format(TimeFormat))); err != nil {
		return fmt.Errorf("%s: unable to write modified_at to file: %w", op, err)
	}

	if _, err = f.WriteString(fmt.Sprintf("text: %s\n", note.Text)); err != nil {
		return fmt.Errorf("%s: unable to write text to file: %w", op, err)
	}

	return nil
}

func Parse(data []byte) (*Note, error) {
	const op = "note.Parse"

	buf := bytes.NewBuffer(data)
	n := Note{}
	var text string

	for stop := false; !stop; {
		line, err := buf.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("%s: error read string: %w", op, err)
		}

		sl := strings.Split(line, ": ")
		switch sl[0] {
		case "title":
			n.Title = strings.Trim(sl[1], "\n")
		case "created_at":
			t, err := time.Parse(TimeFormat, strings.Trim(sl[1], "\n"))
			if err != nil {
				return nil, fmt.Errorf("%s: error parsing created_at: %w", op, err)
			}
			n.CreatedAt = t
		case "modified_at":
			t, err := time.Parse(TimeFormat, strings.Trim(sl[1], "\n"))
			if err != nil {
				return nil, fmt.Errorf("%s: error parsing modified_at: %w", op, err)
			}
			n.ModifiedAt = t
		case "text":
			text = sl[1]
			stop = true
		}
	}

	strBuf := bytes.NewBufferString(text)
	for {
		line, err := buf.ReadString('\n')
		strBuf.WriteString(line)
		if err != nil && err == io.EOF {
			break
		}
	}

	n.Text = strBuf.String()

	if n.Title == "" {
		return nil, fmt.Errorf("%s: error parsing title", op)
	}
	if n.Text == "" {
		return nil, fmt.Errorf("%s: error parsing text", op)
	}

	return &n, nil
}
