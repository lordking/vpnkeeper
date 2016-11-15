package vpn

import (
	"fmt"
	"testing"
)

func init() {

}

func TestFetch(t *testing.T) {

	services, err := Fetch()

	if err != nil {
		fmt.Print("Find vpn error:", err)
		t.Error(err)
	}

	if len(services) == 0 {
		t.Error("Not found vpn")
	}

	for _, service := range services[0:] {
		fmt.Printf("%v, %s, %s\n", service.Status, service.ID, service.Name)
	}

}

func TestStatus(t *testing.T) {

	services, err := Fetch()

	if err != nil {
		fmt.Print("Find vpn error:", err)
		t.Error(err)
	}

	if len(services) == 0 {
		t.Error("Not found vpn")
	}

	status, err := std.Status(services[0])
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("status: %d\n", status)
}
