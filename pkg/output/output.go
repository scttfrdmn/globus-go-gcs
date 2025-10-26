// Package output provides output formatting for the Globus Connect Server CLI.
//
// Supports multiple output formats:
//   - text: Human-readable table format (default)
//   - json: Machine-readable JSON format
//
// Example usage:
//
//	formatter := output.NewFormatter(format, os.Stdout)
//	formatter.PrintJSON(data)
package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format represents the output format type.
type Format string

const (
	// FormatText is human-readable text output (default).
	FormatText Format = "text"

	// FormatJSON is machine-readable JSON output.
	FormatJSON Format = "json"
)

// Formatter handles output formatting for different formats.
type Formatter struct {
	format Format
	writer io.Writer
}

// NewFormatter creates a new output formatter.
//
// Parameters:
//   - format: Output format (text or json)
//   - writer: Destination for output (typically os.Stdout)
func NewFormatter(format Format, writer io.Writer) *Formatter {
	return &Formatter{
		format: format,
		writer: writer,
	}
}

// PrintJSON outputs data in JSON format.
//
// If the formatter is set to JSON format, outputs pretty-printed JSON.
// Otherwise, does nothing (text format should use PrintText instead).
func (f *Formatter) PrintJSON(data interface{}) error {
	if f.format != FormatJSON {
		return nil
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}

// PrintText outputs a text message.
//
// If the formatter is set to text format, outputs the message.
// Otherwise, does nothing (JSON format should use PrintJSON instead).
func (f *Formatter) PrintText(format string, args ...interface{}) error {
	if f.format != FormatText {
		return nil
	}

	if _, err := fmt.Fprintf(f.writer, format, args...); err != nil {
		return fmt.Errorf("write text: %w", err)
	}

	return nil
}

// Println outputs a text line.
//
// If the formatter is set to text format, outputs the line with newline.
// Otherwise, does nothing (JSON format should use PrintJSON instead).
func (f *Formatter) Println(args ...interface{}) error {
	if f.format != FormatText {
		return nil
	}

	if _, err := fmt.Fprintln(f.writer, args...); err != nil {
		return fmt.Errorf("write text line: %w", err)
	}

	return nil
}

// Print outputs data in the appropriate format.
//
// This is a convenience method that automatically chooses PrintJSON or PrintText
// based on the formatter's format setting.
func (f *Formatter) Print(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.PrintJSON(data)
	case FormatText:
		// For text format, try to convert to string
		return f.PrintText("%v\n", data)
	default:
		return fmt.Errorf("unsupported format: %s", f.format)
	}
}

// GetFormat returns the current output format.
func (f *Formatter) GetFormat() Format {
	return f.format
}

// IsJSON returns true if the formatter is set to JSON format.
func (f *Formatter) IsJSON() bool {
	return f.format == FormatJSON
}

// IsText returns true if the formatter is set to text format.
func (f *Formatter) IsText() bool {
	return f.format == FormatText
}
