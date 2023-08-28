package note

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fromString(format string, value string) time.Time {
	t, _ := time.Parse(format, value)
	return t
}

func TestParse(t *testing.T) {
	cases := []struct {
		name string
		arg  []byte
		want Note
		err  bool
	}{
		{
			name: "Ok",
			arg:  []byte("title: Title\ncreated_at: 2023-08-15 16:51\nmodified_at: 2023-08-15 16:51\ntext: lorem lorem\nLorem 123!"),
			want: Note{
				Title:      "Title",
				CreatedAt:  fromString(TimeFormat, "2023-08-15 16:51"),
				ModifiedAt: fromString(TimeFormat, "2023-08-15 16:51"),
				Text:       "lorem lorem\nLorem 123!",
			},
			err: false,
		},
		{
			name: "faile",
			arg:  []byte("title:Titlecreated_at: 2023-08-15 16:51\nmodified_at: 2023-08-15 16:51\ntext: lorem lorem\nLorem 123!"),
			want: Note{},
			err:  true,
		},
	}

	for _, tc := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			n, err := Parse(tc.arg)
			if tc.err {
				assert.Error(t, err)
				assert.Nil(t, n)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, n)
				assert.Equal(t, tc.want, *n)
			}
		})
	}
}
