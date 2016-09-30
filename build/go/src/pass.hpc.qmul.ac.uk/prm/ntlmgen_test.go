package prm

import (
	"strings"
	"testing"
)

// TestNTLMGen - see if the returned hash is actually correct.
// Based on https://www.tobtu.com/lmntlm.php

func TestNTLMGen(t *testing.T) {
	hash := Ntlmgen("n4klxui")
	hash = strings.ToUpper(hash)
	if hash != "734B8DB35667632442FE2A5ABF408DA5" {
		t.Error("hash is incorrect. Got :", hash)
	}
}
