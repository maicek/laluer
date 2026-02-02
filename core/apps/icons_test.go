package apps_test

import (
	"fmt"
	"testing"

	"github.com/maicek/laluer/core/apps"
)

func TestFindIcon(t *testing.T) {
	result := apps.FindIcon("firefox.png")
	// fmt.Println(os.ReadFile(result))
	fmt.Println(result)
}
