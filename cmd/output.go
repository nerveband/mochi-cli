package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
)

// printJSON prints data as formatted JSON
func printJSON(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		if jsonErrors {
			fmt.Fprintf(os.Stderr, `{"error": "failed to marshal JSON: %s"}`, err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %s\n", err.Error())
		}
		return
	}
	fmt.Println(string(output))
}

// printCompactJSON prints data as compact JSON (single line)
func printCompactJSON(data interface{}) {
	output, err := json.Marshal(data)
	if err != nil {
		if jsonErrors {
			fmt.Fprintf(os.Stderr, `{"error": "failed to marshal JSON: %s"}`, err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %s\n", err.Error())
		}
		return
	}
	fmt.Println(string(output))
}

// extractField extracts a specific field from data using reflection
func extractField(data interface{}, field string) interface{} {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	// Handle slice
	if v.Kind() == reflect.Slice {
		var result []interface{}
		for i := 0; i < v.Len(); i++ {
			item := extractField(v.Index(i).Interface(), field)
			if item != nil {
				result = append(result, item)
			}
		}
		return result
	}

	// Handle struct/map
	if v.Kind() == reflect.Map {
		key := reflect.ValueOf(field)
		val := v.MapIndex(key)
		if val.IsValid() {
			return val.Interface()
		}
		return nil
	}

	if v.Kind() == reflect.Struct {
		// Try exact match first
		f := v.FieldByName(field)
		if f.IsValid() {
			return f.Interface()
		}

		// Try with common variations
		variations := []string{
			field,
			strings.Title(field),
			strings.ToUpper(field[:1]) + field[1:],
		}

		for _, varName := range variations {
			f = v.FieldByName(varName)
			if f.IsValid() {
				return f.Interface()
			}
		}
	}

	return nil
}

// printTable prints data as a table
func printTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		if !quiet {
			fmt.Println("No results")
		}
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print headers
	if !noHeaders {
		for i, h := range headers {
			if i > 0 {
				fmt.Print("  ")
			}
			fmt.Printf("%-*s", widths[i], h)
		}
		fmt.Println()

		// Print separator
		for i := range headers {
			if i > 0 {
				fmt.Print("  ")
			}
			fmt.Print(strings.Repeat("-", widths[i]))
		}
		fmt.Println()
	}

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Print("  ")
			}
			fmt.Printf("%-*s", widths[i], cell)
		}
		fmt.Println()
	}
}

// truncateString truncates a string to max length
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// formatTime formats a time value
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

// readStdin reads all content from stdin
func readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// printSuccess prints a success message (unless quiet)
func printSuccess(msg string) {
	if !quiet {
		color.Green(msg)
	}
}

// printWarning prints a warning message (unless quiet)
func printWarning(msg string) {
	if !quiet {
		color.Yellow(msg)
	}
}

// printError prints an error message
func printError(msg string) {
	if jsonErrors {
		fmt.Fprintf(os.Stderr, `{"error": "%s"}`, msg)
	} else {
		color.Red(msg)
	}
}

// printInfo prints an info message (unless quiet)
func printInfo(msg string) {
	if !quiet {
		fmt.Println(msg)
	}
}
