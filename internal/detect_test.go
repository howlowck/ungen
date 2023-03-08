package internal

import (
	"fmt"
	"testing"
)

type DetectTestCase struct {
	Line              string
	ExpectedDetected  bool
	ExpectedExtracted string
}

func TestDetect(t *testing.T) {
	testCases := []DetectTestCase{
		{
			Line:              "// UNGEN: copy ln 10 to cb.keyVault",
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: copy ln 10 to cb.keyVault",
		},
		{
			Line:              "    /* UNGEN: copy ln 10 to cb.keyVault */ ",
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: copy ln 10 to cb.keyVault",
		},
		{
			Line:              "	# UNGEN: copy ln 10 to cb.keyVault ",
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: copy ln 10 to cb.keyVault ",
		},
		{
			Line:              "[//]: # 'UNGEN: replace \"Hello World\" with var.appName'",
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: replace \"Hello World\" with var.appName",
		},
		{
			Line:              `  <!-- UNGEN: replace "Hello World" with var.appName -->`,
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: replace \"Hello World\" with var.appName",
		},
		{
			Line:              `  {/* UNGEN: replace "Hello World" with var.appName */}`,
			ExpectedDetected:  true,
			ExpectedExtracted: "UNGEN: replace \"Hello World\" with var.appName",
		},
	}

	for _, c := range testCases {
		line := c.Line
		actualDetected, actualExtracted := Detect(line)

		if actualDetected != c.ExpectedDetected {
			t.Errorf("Test case \n%s\n failed for detection", c.Line)
			break
		}
		if actualExtracted != c.ExpectedExtracted {
			t.Errorf("Test case \n%s\n failed for extraction", c.Line)
			fmt.Printf(`
actual:   > %s <
expected: > %s <
`, actualExtracted, c.ExpectedExtracted)
			break
		}
		fmt.Println("✅", c.Line)
	}
}
