// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package android

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/exposure-notifications-server/internal/model"
)

const (
	appPackage = "com.google.android.apps.exposurenotification"
)

type emptyNonce struct{}

func (e *emptyNonce) Nonce() string {
	return ""
}

var (
	// **This is not a secret value.**
	// This is a SafetyNet Attestation payload, https://developer.android.com/training/safetynet/attestation#overview
	// generated on a test device without an apk certificate diges.
	// The content of this is the nonce (generated on the other test data) and
	// this particular payload is fixed in time (in the past) doesn't pass
	// basic integrity checks and couldn't be used in a real system.
	// **This comment applies to all test data in this file.**
	payload = `eyJhbGciOiJSUzI1NiIsIng1YyI6WyJNSUlGa3pDQ0JIdWdBd0lCQWdJUkFOY1NramRzNW42K0NBQUFBQUFwYTBjd0RRWUpLb1pJaHZjTkFRRUxCUUF3UWpFTE1Ba0dBMVVFQmhNQ1ZWTXhIakFjQmdOVkJBb1RGVWR2YjJkc1pTQlVjblZ6ZENCVFpYSjJhV05sY3pFVE1CRUdBMVVFQXhNS1IxUlRJRU5CSURGUE1UQWVGdzB5TURBeE1UTXhNVFF4TkRsYUZ3MHlNVEF4TVRFeE1UUXhORGxhTUd3eEN6QUpCZ05WQkFZVEFsVlRNUk13RVFZRFZRUUlFd3BEWVd4cFptOXlibWxoTVJZd0ZBWURWUVFIRXcxTmIzVnVkR0ZwYmlCV2FXVjNNUk13RVFZRFZRUUtFd3BIYjI5bmJHVWdURXhETVJzd0dRWURWUVFERXhKaGRIUmxjM1F1WVc1a2NtOXBaQzVqYjIwd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUNXRXJCUVRHWkdOMWlaYk45ZWhSZ2lmV0J4cWkyUGRneHcwM1A3VHlKWmZNeGpwNUw3ajFHTmVQSzVIemRyVW9JZDF5Q0l5Qk15eHFnYXpxZ3RwWDVXcHNYVzRWZk1oSmJOMVkwOXF6cXA2SkQrMlBaZG9UVTFrRlJBTVdmTC9VdVp0azdwbVJYZ0dtNWpLRHJaOU54ZTA0dk1ZUXI4OE5xd1cva2ZaMWdUT05JVVQwV3NMVC80NTIyQlJXeGZ3eGMzUUUxK1RLV2tMQ3J2ZWs2V2xJcXlhQzUyVzdNRFI4TXBGZWJ5bVNLVHZ3Zk1Sd3lLUUxUMDNVTDR2dDQ4eUVjOHNwN3dUQUhNL1dEZzhRb3RhcmY4T0JIa25vWjkyWGl2aWFWNnRRcWhST0hDZmdtbkNYaXhmVzB3RVhDdnFpTFRiUXRVYkxzUy84SVJ0ZFhrcFFCOUFnTUJBQUdqZ2dKWU1JSUNWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUhBd0V3REFZRFZSMFRBUUgvQkFJd0FEQWRCZ05WSFE0RUZnUVU2REhCd3NBdmI1M2cvQzA3cHJUdnZ3TlFRTFl3SHdZRFZSMGpCQmd3Rm9BVW1OSDRiaERyejV2c1lKOFlrQnVnNjMwSi9Tc3daQVlJS3dZQkJRVUhBUUVFV0RCV01DY0dDQ3NHQVFVRkJ6QUJoaHRvZEhSd09pOHZiMk56Y0M1d2Eya3VaMjl2Wnk5bmRITXhiekV3S3dZSUt3WUJCUVVITUFLR0gyaDBkSEE2THk5d2Eya3VaMjl2Wnk5bmMzSXlMMGRVVXpGUE1TNWpjblF3SFFZRFZSMFJCQll3RklJU1lYUjBaWE4wTG1GdVpISnZhV1F1WTI5dE1DRUdBMVVkSUFRYU1CZ3dDQVlHWjRFTUFRSUNNQXdHQ2lzR0FRUUIxbmtDQlFNd0x3WURWUjBmQkNnd0pqQWtvQ0tnSUlZZWFIUjBjRG92TDJOeWJDNXdhMmt1WjI5dlp5OUhWRk14VHpFdVkzSnNNSUlCQkFZS0t3WUJCQUhXZVFJRUFnU0I5UVNCOGdEd0FIY0E5bHlVTDlGM01DSVVWQmdJTUpSV2p1Tk5FeGt6djk4TUx5QUx6RTd4Wk9NQUFBRnZudXkwWndBQUJBTUFTREJHQWlFQTdlLzBZUnUzd0FGbVdIMjdNMnZiVmNaL21ycCs0cmZZYy81SVBKMjlGNmdDSVFDbktDQ0FhY1ZOZVlaOENDZllkR3BCMkdzSHh1TU9Ia2EvTzQxaldlRit6Z0IxQUVTVVpTNnc3czZ2eEVBSDJLaitLTURhNW9LKzJNc3h0VC9UTTVhMXRvR29BQUFCYjU3c3RKTUFBQVFEQUVZd1JBSWdFWGJpb1BiSnA5cUMwRGoyNThERkdTUk1BVStaQjFFaVZFYmJiLzRVdk5FQ0lCaEhrQnQxOHZSbjl6RHZ5cmZ4eXVkY0hUT1NsM2dUYVlBLzd5VC9CaUg0TUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFESUFjUUJsbWQ4TUVnTGRycnJNYkJUQ3ZwTVhzdDUrd3gyRGxmYWpKTkpVUDRqWUZqWVVROUIzWDRFMnpmNDluWDNBeXVaRnhBcU9SbmJqLzVqa1k3YThxTUowajE5ekZPQitxZXJ4ZWMwbmhtOGdZbExiUW02c0tZN1AwZXhmcjdIdUszTWtQMXBlYzE0d0ZFVWFHcUR3VWJHZ2wvb2l6MzhGWENFK0NXOEUxUUFFVWZ2YlFQVFliS3hZait0Q05sc3MwYlRTb0wyWjJkL2ozQnBMM01GdzB5eFNLL1VUcXlrTHIyQS9NZGhKUW14aStHK01LUlNzUXI2MkFuWmF1OXE2WUZvaSs5QUVIK0E0OFh0SXlzaEx5Q1RVM0h0K2FLb2hHbnhBNXVsMVhSbXFwOEh2Y0F0MzlQOTVGWkdGSmUwdXZseWpPd0F6WHVNdTdNK1BXUmMiLCJNSUlFU2pDQ0F6S2dBd0lCQWdJTkFlTzBtcUdOaXFtQkpXbFF1REFOQmdrcWhraUc5dzBCQVFzRkFEQk1NU0F3SGdZRFZRUUxFeGRIYkc5aVlXeFRhV2R1SUZKdmIzUWdRMEVnTFNCU01qRVRNQkVHQTFVRUNoTUtSMnh2WW1Gc1UybG5iakVUTUJFR0ExVUVBeE1LUjJ4dlltRnNVMmxuYmpBZUZ3MHhOekEyTVRVd01EQXdOREphRncweU1URXlNVFV3TURBd05ESmFNRUl4Q3pBSkJnTlZCQVlUQWxWVE1SNHdIQVlEVlFRS0V4VkhiMjluYkdVZ1ZISjFjM1FnVTJWeWRtbGpaWE14RXpBUkJnTlZCQU1UQ2tkVVV5QkRRU0F4VHpFd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURRR005RjFJdk4wNXprUU85K3ROMXBJUnZKenp5T1RIVzVEekVaaEQyZVBDbnZVQTBRazI4RmdJQ2ZLcUM5RWtzQzRUMmZXQllrL2pDZkMzUjNWWk1kUy9kTjRaS0NFUFpSckF6RHNpS1VEelJybUJCSjV3dWRnem5kSU1ZY0xlL1JHR0ZsNXlPRElLZ2pFdi9TSkgvVUwrZEVhbHROMTFCbXNLK2VRbU1GKytBY3hHTmhyNTlxTS85aWw3MUkyZE44RkdmY2Rkd3VhZWo0YlhocDBMY1FCYmp4TWNJN0pQMGFNM1Q0SStEc2F4bUtGc2JqemFUTkM5dXpwRmxnT0lnN3JSMjV4b3luVXh2OHZObWtxN3pkUEdIWGt4V1k3b0c5aitKa1J5QkFCazdYckpmb3VjQlpFcUZKSlNQazdYQTBMS1cwWTN6NW96MkQwYzF0Skt3SEFnTUJBQUdqZ2dFek1JSUJMekFPQmdOVkhROEJBZjhFQkFNQ0FZWXdIUVlEVlIwbEJCWXdGQVlJS3dZQkJRVUhBd0VHQ0NzR0FRVUZCd01DTUJJR0ExVWRFd0VCL3dRSU1BWUJBZjhDQVFBd0hRWURWUjBPQkJZRUZKalIrRzRRNjgrYjdHQ2ZHSkFib090OUNmMHJNQjhHQTFVZEl3UVlNQmFBRkp2aUIxZG5IQjdBYWdiZVdiU2FMZC9jR1lZdU1EVUdDQ3NHQVFVRkJ3RUJCQ2t3SnpBbEJnZ3JCZ0VGQlFjd0FZWVphSFIwY0RvdkwyOWpjM0F1Y0d0cExtZHZiMmN2WjNOeU1qQXlCZ05WSFI4RUt6QXBNQ2VnSmFBamhpRm9kSFJ3T2k4dlkzSnNMbkJyYVM1bmIyOW5MMmR6Y2pJdlozTnlNaTVqY213d1B3WURWUjBnQkRnd05qQTBCZ1puZ1F3QkFnSXdLakFvQmdnckJnRUZCUWNDQVJZY2FIUjBjSE02THk5d2Eya3VaMjl2Wnk5eVpYQnZjMmwwYjNKNUx6QU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFHb0ErTm5uNzh5NnBSamQ5WGxRV05hN0hUZ2laL3IzUk5Ha21VbVlIUFFxNlNjdGk5UEVhanZ3UlQyaVdUSFFyMDJmZXNxT3FCWTJFVFV3Z1pRK2xsdG9ORnZoc085dHZCQ09JYXpwc3dXQzlhSjl4anU0dFdEUUg4TlZVNllaWi9YdGVEU0dVOVl6SnFQalk4cTNNRHhyem1xZXBCQ2Y1bzhtdy93SjRhMkc2eHpVcjZGYjZUOE1jRE8yMlBMUkw2dTNNNFR6czNBMk0xajZieWtKWWk4d1dJUmRBdktMV1p1L2F4QlZielltcW13a201ekxTRFc1bklBSmJFTENRQ1p3TUg1NnQyRHZxb2Z4czZCQmNDRklaVVNweHU2eDZ0ZDBWN1N2SkNDb3NpclNtSWF0ai85ZFNTVkRRaWJldDhxLzdVSzR2NFpVTjgwYXRuWnoxeWc9PSJdfQ.eyJub25jZSI6ImRuUmhhR1p6VEhSWlJFVkpiWE5wYm1WeVdsSnhheXR3T0ZCMVdHWnZlamhvYlZCRGMyaHNVM3BuZHowPSIsInRpbWVzdGFtcE1zIjoxNTg5MTU0MDA2NDk1LCJhcGtQYWNrYWdlTmFtZSI6ImNvbS5nb29nbGUuYW5kcm9pZC5hcHBzLmV4cG9zdXJlbm90aWZpY2F0aW9uIiwiYXBrRGlnZXN0U2hhMjU2IjoiVFlLdXRVT3hVekpjZmxFb0x3My9Ya2xoNXFoK08wanF4WTU2c3YycEZVYz0iLCJjdHNQcm9maWxlTWF0Y2giOnRydWUsImFwa0NlcnRpZmljYXRlRGlnZXN0U2hhMjU2IjpbImpxbVlFcWk5cVV2cFVlMTFxTWYzdjJvNlZFUU0rNU5EZWUyYnoweGR6V2M9Il0sImJhc2ljSW50ZWdyaXR5Ijp0cnVlLCJldmFsdWF0aW9uVHlwZSI6IkJBU0lDIn0.TDOJdtDgDybKwzmvEWQuKxiiktBSPrlF0UrbtBWEovpCa_36orW5gMqmB6_nwviC1PK65dNoXpmbQ4mJa9YwN3Tm99DWNdawbswVRP37yYSQuL2-m_NPfxeXy9qH38nfyrXwWceZbrjOAxOKZ-NZqpACiahy0PfGezLr8WXGeOT_a41xtP_pjVvG4fblcFso_3kFIF0t9L2uA5W-wZwNUiS_tl0kpbU-wJMHWwyW_aM_6Oyb0e5I1HaLh8v2JLsoFgr9oX0uiAifgm41mNJdsXNHg4k6JZeiFSQQvSq30TJ_VjmYz8hioeDBXQPb3w9aoIwq3_yCq297lJ-r8oMcLQ`

	publish = &model.Publish{
		Keys: []model.ExposureKey{
			{Key: "HKXVlIO+vDmQNJ2M1MVtHQ==", IntervalNumber: 2647872, IntervalCount: 144},
			{Key: "JjEtCT9Lcyw5oPiaNcWC/Q==", IntervalNumber: 2648016, IntervalCount: 144},
			{Key: "cLTwDu9onEv/N6FMV3Uy4Q==", IntervalNumber: 2648160, IntervalCount: 144},
			{Key: "ko6TsgPP8Wvu+ijpSLbY3A==", IntervalNumber: 2648304, IntervalCount: 144},
			{Key: "9kMgBy7qdG3o6eh3vAD3mQ==", IntervalNumber: 2648448, IntervalCount: 144},
		},
		TransmissionRisk:          1,
		AppPackageName:            appPackage,
		VerificationAuthorityName: "PUBLIC_HEALTH_AUTHORITY",
		Regions:                   []string{"GB"},
	}
)

func TestVerifyAttestation(t *testing.T) {
	ctx := context.Background()
	claims, err := verifyAttestation(ctx, payload)
	if err != nil {
		t.Fatalf("error verifying attestation %v", err)
	}

	expectedNonce := NewNonce(publish).Nonce()
	actualBytes, err := base64.StdEncoding.DecodeString(claims["nonce"].(string))
	actualNonce := string(actualBytes)
	if err != nil {
		t.Fatalf("unable to decode nonce from attestation: %v", err)
	}
	if actualNonce != expectedNonce {
		t.Errorf("attestation nonce, want: %v got %v", expectedNonce, actualNonce)
	}
}

func TestValidateAttestation(t *testing.T) {
	nonce := NewNonce(publish)

	apkDigest := "jqmYEqi9qUvpUe11qMf3v2o6VEQM+5NDee2bz0xdzWc="
	generateTimeS := int64(1589154006495) / 1000
	generateTime := time.Unix(generateTimeS, 0)
	maxValidTime := generateTime.Add(1 * time.Minute)
	minValidTime := generateTime.Add(-1 * time.Minute)

	tests := []struct {
		Opts  VerifyOpts
		Valid bool
		Error string
	}{
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       apkDigest,
				Nonce:           nonce,
				CTSProfileMatch: false,
				BasicIntegrity:  false,
				MinValidTime:    minValidTime,
				MaxValidTime:    maxValidTime,
			},
			true,
			"",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       "",
				Nonce:           nil,
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    minValidTime,
				MaxValidTime:    maxValidTime,
			},
			false,
			"missing nonce",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       "",
				Nonce:           &emptyNonce{},
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    minValidTime,
				MaxValidTime:    maxValidTime,
			},
			false,
			"missing nonce",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       "",
				Nonce:           nonce,
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    time.Time{},
				MaxValidTime:    time.Time{},
			},
			false,
			"missing timestamp bounds for attestation",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       "",
				Nonce:           nonce,
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    time.Time{},
				MaxValidTime:    maxValidTime,
			},
			false,
			"missing timestamp bounds for attestation",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       apkDigest,
				Nonce:           nonce,
				CTSProfileMatch: true,
				BasicIntegrity:  true,
				MinValidTime:    minValidTime,
				MaxValidTime:    time.Time{},
			},
			false,
			"missing timestamp bounds for attestation",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       apkDigest,
				Nonce:           nonce,
				CTSProfileMatch: false,
				BasicIntegrity:  false,
				MinValidTime:    maxValidTime,
				MaxValidTime:    maxValidTime,
			},
			false,
			"attestation is too old, must be newer than 1589154066, was 1589154006",
		},
		{
			VerifyOpts{
				AppPkgName:      appPackage,
				APKDigest:       "",
				Nonce:           nonce,
				CTSProfileMatch: false,
				BasicIntegrity:  false,
				MinValidTime:    minValidTime,
				MaxValidTime:    minValidTime,
			},
			false,
			"attestation is in the future, must be older than 1589153946, was 1589154006",
		},
	}

	for i, test := range tests {
		ctx := context.Background()
		err := ValidateAttestation(ctx, payload, test.Opts)
		if test.Valid && err != nil {
			t.Errorf("test %v, wanted valid, got %v", i, err)
		} else if !test.Valid {
			if err == nil {
				t.Errorf("test %v, expected error, want %v, got: nil", i, test.Error)
			} else if err.Error() != test.Error {
				t.Errorf("test %v, wrong error, want %v, got %v", i, test.Error, err)
			}
		}
	}
}
