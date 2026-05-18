package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Message
		wantErr bool
	}{
		{
			input:   "PING :gallium.libera.chat",
			want:    Message{"", "PING", []string{"gallium.libera.chat"}},
			wantErr: false,
		},
		{
			input:   ":gallium.libera.chat 001 lpall :Welcome to the Libera.Chat Internet Relay Chat Network lpall",
			want:    Message{"gallium.libera.chat", "001", []string{"lpall", "Welcome to the Libera.Chat Internet Relay Chat Network lpall"}},
			wantErr: false,
		},
		{
			input:   ":lpall!~lpall@45.130.200.119 JOIN #test",
			want:    Message{"lpall!~lpall@45.130.200.119", "JOIN", []string{"#test"}},
			wantErr: false,
		},
		{
			input:   ":gallium.libera.chat 353 lpall = #test :lpall milaextract x1bncwn nitrix_",
			want:    Message{"gallium.libera.chat", "353", []string{"lpall", "=", "#test", "lpall milaextract x1bncwn nitrix_"}},
			wantErr: false,
		},
		{
			input:   ":lpall MODE lpall :+iw",
			want:    Message{"lpall", "MODE", []string{"lpall", "+iw"}},
			wantErr: false,
		},
		{
			input:   "",
			want:    Message{},
			wantErr: true,
		},
		{
			input:   ":prefixonly",
			want:    Message{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		got, err := parse(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parse(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parse(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
