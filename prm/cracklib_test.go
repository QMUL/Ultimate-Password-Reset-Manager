package prm

import (
        "testing"
)

// Setup test cases
var test_map = map[string]string {
    "adoiah1223423sfdiunIOH": "GOOD",
    "short": "it is too short",
    "repreprep": "it does not contain enough DIFFERENT characters",
    "proletariat": "it is based on a dictionary word",
    "12345678": "it is too simplistic/systematic",
}

// Test password validation
func TestCracklibPassword(t *testing.T) {
    for key, value := range test_map {
        ret := TestPassword(key)
        if ret != value {
            t.Error("For:", key, "got:", ret, "expected:", value)
        }
    }
}
