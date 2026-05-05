package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  Message
	}{
		{
			input: "PING :gallium.libera.chat",
			want:  Message{"", "PING", []string{"gallium.libera.chat"}},
		},
		{
			input: ":gallium.libera.chat 001 lpall :Welcome to the Libera.Chat Internet Relay Chat Network lpall",
			want:  Message{"gallium.libera.chat", "001", []string{"lpall", "Welcome to the Libera.Chat Internet Relay Chat Network lpall"}},
		},
		{
			input: ":lpall!~lpall@45.130.200.119 JOIN #test",
			want:  Message{"lpall!~lpall@45.130.200.119", "JOIN", []string{"#test"}},
		},
		{
			input: ":gallium.libera.chat 353 lpall = #test :lpall milaextract x1bncwn nitrix_",
			want:  Message{"gallium.libera.chat", "353", []string{"lpall", "=", "#test", "lpall milaextract x1bncwn nitrix_"}},
		},
		{
			input: ":lpall MODE lpall :+iw",
			want:  Message{"lpall", "MODE", []string{"lpall", "+iw"}},
		},
	}

	for _, tt := range tests {
		got := parse(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
