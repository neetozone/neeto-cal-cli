package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newAddOnServiceTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	addAddOnServiceFlags(cmd)
	return cmd
}

func TestAddOnServiceBody_SendsOnlyChangedFlags(t *testing.T) {
	cmd := newAddOnServiceTestCmd()
	setFlags(t, cmd, map[string]string{"name": "Recording", "duration": "15", "price": "25.00"})

	body := addOnServiceBody(cmd)

	want := map[string]interface{}{"name": "Recording", "duration": 15, "price": "25.00"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func TestAddOnServiceBody_SendsPositionAndDescription(t *testing.T) {
	cmd := newAddOnServiceTestCmd()
	setFlags(t, cmd, map[string]string{"description": "Session recording", "position": "2"})

	body := addOnServiceBody(cmd)

	want := map[string]interface{}{"description": "Session recording", "position": 2}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}
