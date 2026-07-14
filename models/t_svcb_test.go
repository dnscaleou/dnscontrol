package models

import (
	"testing"

	dnsv2 "codeberg.org/miekg/dns"
)

func TestPopulateFromStringHTTPSAndSVCB(t *testing.T) {
	tests := []struct {
		name       string
		rtype      string
		typeNum    uint16
		content    string
		priority   uint16
		target     string
		params     string
		comparable string
	}{
		{
			name:       "HTTPS alias form",
			rtype:      "HTTPS",
			typeNum:    dnsv2.TypeHTTPS,
			content:    "0 test.com.",
			priority:   0,
			target:     "test.com.",
			comparable: "0 test.com.",
		},
		{
			name:       "HTTPS with params",
			rtype:      "HTTPS",
			typeNum:    dnsv2.TypeHTTPS,
			content:    "1 test.com. alpn=h2,h3 port=999",
			priority:   1,
			target:     "test.com.",
			params:     "alpn=h2,h3 port=999",
			comparable: `1 test.com. alpn="h2,h3" port="999"`,
		},
		{
			name:       "HTTPS with ECH",
			rtype:      "HTTPS",
			typeNum:    dnsv2.TypeHTTPS,
			content:    `3 example.com. alpn=h2,h3 port=999 ech="some+base64+encoded+value///"`,
			priority:   3,
			target:     "example.com.",
			params:     "alpn=h2,h3 port=999 ech=some+base64+encoded+value///",
			comparable: `3 example.com. alpn="h2,h3" port="999" ech="some+base64+encoded+value///"`,
		},
		{
			name:       "SVCB alias form",
			rtype:      "SVCB",
			typeNum:    dnsv2.TypeSVCB,
			content:    "0 test.com.",
			priority:   0,
			target:     "test.com.",
			comparable: "0 test.com.",
		},
		{
			name:       "SVCB with params",
			rtype:      "SVCB",
			typeNum:    dnsv2.TypeSVCB,
			content:    "1 test.com. alpn=h2,h3 port=999",
			priority:   1,
			target:     "test.com.",
			params:     "alpn=h2,h3 port=999",
			comparable: `1 test.com. alpn="h2,h3" port="999"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RecordConfig{}
			rc.SetLabel("@", "example.com")

			if err := rc.PopulateFromString(tt.rtype, tt.content, "example.com"); err != nil {
				t.Fatalf("PopulateFromString() returned error: %v", err)
			}

			if rc.Type != tt.rtype {
				t.Fatalf("Type expected %q, got %q", tt.rtype, rc.Type)
			}
			if rc.TypeNum != tt.typeNum {
				t.Fatalf("TypeNum expected %d, got %d", tt.typeNum, rc.TypeNum)
			}
			if rc.SvcPriority != tt.priority {
				t.Fatalf("SvcPriority expected %d, got %d", tt.priority, rc.SvcPriority)
			}
			if rc.GetTargetField() != tt.target {
				t.Fatalf("target expected %q, got %q", tt.target, rc.GetTargetField())
			}
			if rc.SvcParams != tt.params {
				t.Fatalf("SvcParams expected %q, got %q", tt.params, rc.SvcParams)
			}
			if got := rc.ToComparableNoTTL(); got != tt.comparable {
				t.Fatalf("ToComparableNoTTL expected %q, got %q", tt.comparable, got)
			}
			if rc.ComparableV3 != tt.comparable {
				t.Fatalf("ComparableV3 expected %q, got %q", tt.comparable, rc.ComparableV3)
			}
		})
	}
}

func TestSetTargetHTTPSStringDefaultsType(t *testing.T) {
	rc := &RecordConfig{}

	if err := rc.SetTargetHTTPSString("example.com", "0 test.com."); err != nil {
		t.Fatalf("SetTargetHTTPSString() returned error: %v", err)
	}

	if rc.Type != "HTTPS" {
		t.Fatalf("Type expected HTTPS, got %q", rc.Type)
	}
	if rc.TypeNum != dnsv2.TypeHTTPS {
		t.Fatalf("TypeNum expected %d, got %d", dnsv2.TypeHTTPS, rc.TypeNum)
	}
	if rc.SvcPriority != 0 {
		t.Fatalf("SvcPriority expected 0, got %d", rc.SvcPriority)
	}
	if rc.GetTargetField() != "test.com." {
		t.Fatalf("target expected test.com., got %q", rc.GetTargetField())
	}
}
