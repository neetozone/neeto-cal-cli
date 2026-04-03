package output

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Breadcrumb struct {
	Label   string `json:"label"`
	Command string `json:"command"`
}

type Envelope struct {
	Data        json.RawMessage `json:"data"`
	Breadcrumbs []Breadcrumb    `json:"breadcrumbs,omitempty"`
	Pagination  json.RawMessage `json:"pagination,omitempty"`
}

var (
	ForceJSON bool
	QuietMode bool
)

func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func UseJSON() bool {
	return ForceJSON || QuietMode || !IsTTY()
}

func Print(data json.RawMessage, breadcrumbs []Breadcrumb) {
	if QuietMode {
		fmt.Println(string(data))
		return
	}

	if UseJSON() {
		printEnvelope(data, breadcrumbs, nil)
		return
	}

	printPretty(data)
}

func PrintWithPagination(data json.RawMessage, pagination json.RawMessage, breadcrumbs []Breadcrumb) {
	if QuietMode {
		fmt.Println(string(data))
		return
	}

	if UseJSON() {
		printEnvelope(data, breadcrumbs, pagination)
		return
	}

	printPretty(data)
	printPaginationSummary(pagination)
}

func PrintMessage(msg string) {
	if UseJSON() {
		envelope := map[string]string{"message": msg}
		data, _ := json.Marshal(envelope)
		fmt.Println(string(data))
	} else {
		fmt.Println(msg)
	}
}

func printEnvelope(data json.RawMessage, breadcrumbs []Breadcrumb, pagination json.RawMessage) {
	envelope := Envelope{
		Data:        data,
		Breadcrumbs: breadcrumbs,
		Pagination:  pagination,
	}
	out, _ := json.MarshalIndent(envelope, "", "  ")
	fmt.Println(string(out))
}

func printPretty(data json.RawMessage) {
	var out []byte
	out, err := json.MarshalIndent(json.RawMessage(data), "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(out))
}

func printPaginationSummary(pagination json.RawMessage) {
	if pagination == nil {
		return
	}

	var p struct {
		CurrentPageNumber int `json:"current_page_number"`
		TotalPages        int `json:"total_pages"`
		TotalRecords      int `json:"total_records"`
	}
	if err := json.Unmarshal(pagination, &p); err != nil {
		return
	}

	fmt.Printf("\nPage %d of %d (%d total records)\n", p.CurrentPageNumber, p.TotalPages, p.TotalRecords)
}
