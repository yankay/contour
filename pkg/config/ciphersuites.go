// Copyright Project Contour Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"strings"
)

// TLSCiphers holds a list of TLS ciphers
type TLSCiphers []string

// DefaultTLSCiphers contains the list of default ciphers used by Contour. A handful are
// commented out, as they're arguably less secure. They're also unnecessary
// - most of the clients that might need to use the commented ciphers are
// unable to connect without TLS 1.0, which contour never enables.
//
// This list is ignored if the client and server negotiate TLS 1.3.
//
// The commented ciphers are left in place to simplify updating this list for future
// versions of envoy.
var DefaultTLSCiphers = TLSCiphers([]string{
	"[ECDHE-ECDSA-AES128-GCM-SHA256|ECDHE-ECDSA-CHACHA20-POLY1305]",
	"[ECDHE-RSA-AES128-GCM-SHA256|ECDHE-RSA-CHACHA20-POLY1305]",
	// "ECDHE-ECDSA-AES128-SHA",
	// "ECDHE-RSA-AES128-SHA",
	// "AES128-GCM-SHA256",
	// "AES128-SHA",
	"ECDHE-ECDSA-AES256-GCM-SHA384",
	"ECDHE-RSA-AES256-GCM-SHA384",
	// "ECDHE-ECDSA-AES256-SHA",
	// "ECDHE-RSA-AES256-SHA",
	// "AES256-GCM-SHA384",
	// "AES256-SHA",
})

// validTLSCiphers contains the list of TLS ciphers that Envoy supports
// See: https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/transport_sockets/tls/v3/common.proto#extensions-transport-sockets-tls-v3-tlsparameters
// Note: This list is a superset of what is valid for stock Envoy builds and those using BoringSSL FIPS.
var validTLSCiphers = map[string]struct{}{
	"[ECDHE-ECDSA-AES128-GCM-SHA256|ECDHE-ECDSA-CHACHA20-POLY1305]": {},
	"[ECDHE-RSA-AES128-GCM-SHA256|ECDHE-RSA-CHACHA20-POLY1305]":     {},
	"ECDHE-ECDSA-AES128-GCM-SHA256":                                 {},
	"ECDHE-RSA-AES128-GCM-SHA256":                                   {},
	"ECDHE-ECDSA-AES128-SHA":                                        {},
	"ECDHE-RSA-AES128-SHA":                                          {},
	"AES128-GCM-SHA256":                                             {},
	"AES128-SHA":                                                    {},
	"ECDHE-ECDSA-AES256-GCM-SHA384":                                 {},
	"ECDHE-RSA-AES256-GCM-SHA384":                                   {},
	"ECDHE-ECDSA-AES256-SHA":                                        {},
	"ECDHE-RSA-AES256-SHA":                                          {},
	"AES256-GCM-SHA384":                                             {},
	"AES256-SHA":                                                    {},
}

// SanitizeCipherSuites trims a list of ciphers to remove whitespace and
// duplicates, returning the passed in default if the corrected list is empty.
// The ciphers argument should be a list of valid ciphers.
func SanitizeCipherSuites(ciphers []string) []string {
	if len(ciphers) == 0 {
		return DefaultTLSCiphers
	}

	uniqueCiphers := map[string]bool{}
	validatedCiphers := []string{}
	for _, v := range ciphers {
		cipher := strings.TrimSpace(v)
		if _, found := uniqueCiphers[cipher]; !found {
			uniqueCiphers[cipher] = true
			validatedCiphers = append(validatedCiphers, cipher)
		}
	}
	return validatedCiphers
}

// Validate ciphers. Returns error on unsupported cipher.
func (tlsCiphers TLSCiphers) Validate() error {
	invalidCiphers := []string{}
	for _, cipher := range tlsCiphers {
		trimmed := strings.TrimSpace(cipher)
		if _, ok := validTLSCiphers[trimmed]; !ok {
			invalidCiphers = append(invalidCiphers, trimmed)
		}
	}
	if len(invalidCiphers) > 0 {
		return fmt.Errorf("invalid ciphers: %s", strings.Join(invalidCiphers, ","))
	}
	return nil
}
