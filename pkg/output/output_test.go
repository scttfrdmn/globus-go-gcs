package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewFormatter(FormatJSON, buf)

	if formatter == nil {
		t.Fatal("NewFormatter() returned nil")
	}

	if formatter.format != FormatJSON {
		t.Errorf("NewFormatter() format = %v, want %v", formatter.format, FormatJSON)
	}

	if formatter.writer != buf {
		t.Error("NewFormatter() writer not set correctly")
	}
}

func TestFormatter_PrintJSON(t *testing.T) {
	tests := []struct {
		name    string
		format  Format
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "JSON format with simple object",
			format:  FormatJSON,
			data:    map[string]string{"key": "value"},
			want:    "{\n  \"key\": \"value\"\n}\n",
			wantErr: false,
		},
		{
			name:    "JSON format with array",
			format:  FormatJSON,
			data:    []string{"item1", "item2"},
			want:    "[\n  \"item1\",\n  \"item2\"\n]\n",
			wantErr: false,
		},
		{
			name:    "text format (should not output)",
			format:  FormatText,
			data:    map[string]string{"key": "value"},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			err := formatter.PrintJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("PrintJSON() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatter_PrintText(t *testing.T) {
	tests := []struct {
		name    string
		format  Format
		text    string
		args    []interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "text format with string",
			format:  FormatText,
			text:    "Hello, %s!",
			args:    []interface{}{"world"},
			want:    "Hello, world!",
			wantErr: false,
		},
		{
			name:    "text format with multiple args",
			format:  FormatText,
			text:    "Count: %d, Name: %s",
			args:    []interface{}{42, "test"},
			want:    "Count: 42, Name: test",
			wantErr: false,
		},
		{
			name:    "JSON format (should not output)",
			format:  FormatJSON,
			text:    "Should not appear",
			args:    nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			err := formatter.PrintText(tt.text, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("PrintText() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatter_Println(t *testing.T) {
	tests := []struct {
		name    string
		format  Format
		args    []interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "text format with single arg",
			format:  FormatText,
			args:    []interface{}{"Hello"},
			want:    "Hello\n",
			wantErr: false,
		},
		{
			name:    "text format with multiple args",
			format:  FormatText,
			args:    []interface{}{"Hello", "world"},
			want:    "Hello world\n",
			wantErr: false,
		},
		{
			name:    "JSON format (should not output)",
			format:  FormatJSON,
			args:    []interface{}{"Should not appear"},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			err := formatter.Println(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Println() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("Println() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatter_Print(t *testing.T) {
	tests := []struct {
		name    string
		format  Format
		data    interface{}
		wantErr bool
		check   func(string) bool
	}{
		{
			name:    "JSON format",
			format:  FormatJSON,
			data:    map[string]string{"key": "value"},
			wantErr: false,
			check: func(s string) bool {
				return strings.Contains(s, "\"key\"") && strings.Contains(s, "\"value\"")
			},
		},
		{
			name:    "text format",
			format:  FormatText,
			data:    "test message",
			wantErr: false,
			check: func(s string) bool {
				return strings.Contains(s, "test message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			err := formatter.Print(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.check(buf.String()) {
				t.Errorf("Print() output check failed, got: %q", buf.String())
			}
		})
	}
}

func TestFormatter_GetFormat(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "JSON format",
			format: FormatJSON,
		},
		{
			name:   "text format",
			format: FormatText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			if got := formatter.GetFormat(); got != tt.format {
				t.Errorf("GetFormat() = %v, want %v", got, tt.format)
			}
		})
	}
}

func TestFormatter_IsJSON(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   bool
	}{
		{
			name:   "JSON format",
			format: FormatJSON,
			want:   true,
		},
		{
			name:   "text format",
			format: FormatText,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			if got := formatter.IsJSON(); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatter_IsText(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   bool
	}{
		{
			name:   "text format",
			format: FormatText,
			want:   true,
		},
		{
			name:   "JSON format",
			format: FormatJSON,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(tt.format, buf)

			if got := formatter.IsText(); got != tt.want {
				t.Errorf("IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}
