package commands

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newDiscountCodesUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	cmd.Flags().String("code", "", "Discount code")
	cmd.Flags().String("kind", "", "Code kind (percentage, fixed)")
	cmd.Flags().Int("value", 0, "Discount value")
	cmd.Flags().String("meeting-ids", "", "Comma-separated meeting SIDs")
	cmd.Flags().String("expires-at", "", "Expiration date (YYYY-MM-DD)")
	cmd.Flags().String("json-file", "", "Path to JSON file with discount code data")
	return cmd
}

func TestDiscountCodeUpdateBody_OmitsFlagsTheCallerNeverSet(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("code", "SAVE20"); err != nil {
		t.Fatalf("set code: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{"code": "SAVE20"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v; an unset --value must not send 0 and blank the discount", body, want)
	}
}

func TestDiscountCodeUpdateBody_SendsEveryFlagTheCallerSet(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	for flag, value := range map[string]string{
		"code":        "SAVE20",
		"kind":        "percentage",
		"value":       "20",
		"meeting-ids": "abc123x,def456y",
		"expires-at":  "2026-12-31",
	} {
		if err := cmd.Flags().Set(flag, value); err != nil {
			t.Fatalf("set %s: %v", flag, err)
		}
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	want := map[string]interface{}{
		"code":        "SAVE20",
		"kind":        "percentage",
		"value":       20,
		"meeting_ids": []string{"abc123x", "def456y"},
		"expires_at":  "2026-12-31",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestDiscountCodeUpdateBody_KeepsAnExplicitZeroValue(t *testing.T) {
	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("value", "0"); err != nil {
		t.Fatalf("set value: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}
	if body["value"] != 0 {
		t.Errorf("body[value] = %v, want 0 to reach the API so it can reject it", body["value"])
	}
}

func TestDiscountCodeUpdateBody_LetsFlagsOverrideTheJSONFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "code.json")
	contents := `{"code":"FROMFILE","redeems_limit":5}`
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("json-file", path); err != nil {
		t.Fatalf("set json-file: %v", err)
	}
	if err := cmd.Flags().Set("code", "FROMFLAG"); err != nil {
		t.Fatalf("set code: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	if body["code"] != "FROMFLAG" {
		t.Errorf("body[code] = %v, want the flag to win over the file", body["code"])
	}
	if body["redeems_limit"] != float64(5) {
		t.Errorf("body[redeems_limit] = %v, want the file to supply fields that have no flag", body["redeems_limit"])
	}
}

func newDiscountCodesListTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "list"}
	addPaginationFlags(cmd)
	cmd.Flags().String("code", "", "Filter by codes containing this text")
	cmd.Flags().String("kind", "", "Filter by code kind (percentage, fixed)")
	return cmd
}

func TestDiscountCodeFilterQuery_IsEmptyWithoutFilterFlags(t *testing.T) {
	if got := discountCodeFilterQuery(newDiscountCodesListTestCmd()); got != "" {
		t.Errorf("discountCodeFilterQuery() = %q, want empty so the request carries no filters key", got)
	}
}

func TestDiscountCodeFilterQuery_KeepsEachConditionsFieldsTogetherInOrder(t *testing.T) {
	cmd := newDiscountCodesListTestCmd()
	if err := cmd.Flags().Set("code", "SAVE 20"); err != nil {
		t.Fatalf("set code: %v", err)
	}
	if err := cmd.Flags().Set("kind", "fixed"); err != nil {
		t.Fatalf("set kind: %v", err)
	}

	got := discountCodeFilterQuery(cmd)

	want := strings.Join([]string{
		"node=code", "type=text", "rule=contains", "value=SAVE+20", "conditions_join_type=and",
		"node=kind", "type=text", "rule=is", "value=fixed", "conditions_join_type=and",
	}, "&")
	if simplifyFilterQuery(got) != want {
		t.Errorf("filter query = %q\nsimplified = %q\nwant       = %q\nRails reads these as an ordered nested query; a duplicate key is what starts the next condition, so order carries meaning",
			got, simplifyFilterQuery(got), want)
	}
}

func TestDiscountCodesListPath_AppendsPaginationAfterTheFilters(t *testing.T) {
	cmd := newDiscountCodesListTestCmd()
	if err := cmd.Flags().Set("kind", "percentage"); err != nil {
		t.Fatalf("set kind: %v", err)
	}
	if err := cmd.Flags().Set("page", "2"); err != nil {
		t.Fatalf("set page: %v", err)
	}

	path := discountCodesListPath(cmd)

	if !strings.HasPrefix(path, "/discount-codes?"+url.QueryEscape("filters[][conditions][][node]")) {
		t.Errorf("path = %q, want the ordered filter keys to lead the query", path)
	}
	if !strings.Contains(path, "page=2") {
		t.Errorf("path = %q, want the pagination keys kept", path)
	}
}

func TestDiscountCodesListPath_HasNoQueryWithoutFlags(t *testing.T) {
	if got := discountCodesListPath(newDiscountCodesListTestCmd()); got != "/discount-codes" {
		t.Errorf("discountCodesListPath() = %q, want a bare path", got)
	}
}

func simplifyFilterQuery(query string) string {
	var parts []string
	for _, pair := range strings.Split(query, "&") {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			continue
		}
		decoded, err := url.QueryUnescape(key)
		if err != nil {
			continue
		}
		field := strings.TrimSuffix(strings.TrimPrefix(decoded, "filters[][conditions][]["), "]")
		parts = append(parts, field+"="+value)
	}
	return strings.Join(parts, "&")
}

func TestDiscountCodeFilterQuery_HandlesOneFilterOnItsOwn(t *testing.T) {
	cmd := newDiscountCodesListTestCmd()
	if err := cmd.Flags().Set("kind", "percentage"); err != nil {
		t.Fatalf("set kind: %v", err)
	}

	want := strings.Join([]string{
		"node=kind", "type=text", "rule=is", "value=percentage", "conditions_join_type=and",
	}, "&")
	if got := simplifyFilterQuery(discountCodeFilterQuery(cmd)); got != want {
		t.Errorf("filter query = %q, want %q", got, want)
	}
}

func TestDiscountCodeUpdateBody_LetsAnEmptyFlagOverrideTheJSONFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "code.json")
	if err := os.WriteFile(path, []byte(`{"expires_at":"2026-12-31","value":10}`), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	cmd := newDiscountCodesUpdateTestCmd()
	if err := cmd.Flags().Set("json-file", path); err != nil {
		t.Fatalf("set json-file: %v", err)
	}
	if err := cmd.Flags().Set("expires-at", ""); err != nil {
		t.Fatalf("set expires-at: %v", err)
	}
	if err := cmd.Flags().Set("value", "25"); err != nil {
		t.Fatalf("set value: %v", err)
	}

	body, err := discountCodeUpdateBody(cmd)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	if body["expires_at"] != "" {
		t.Errorf("body[expires_at] = %v, want the explicit empty flag to win over the file", body["expires_at"])
	}
	if body["value"] != 25 {
		t.Errorf("body[value] = %v, want the flag to override the file", body["value"])
	}
}
