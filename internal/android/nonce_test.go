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
	"testing"

	"github.com/google/exposure-notifications-server/internal/model"
)

// Data from this test was generated by the Android reference application.
func TestGetNonce(t *testing.T) {
	testData := []struct {
		Data     model.Publish
		Expected string
	}{
		{
			Data: model.Publish{
				Keys: []model.ExposureKey{
					{Key: "x21Goi8X9m/glOZ0+wz8fA", IntervalNumber: 263123, IntervalCount: 144},
					{Key: "2mvFSmRsFmJR5r07dxGSjg", IntervalNumber: 263267, IntervalCount: 144},
					{Key: "6bAd3dv7p+VEuaJVkVItaQ", IntervalNumber: 263411, IntervalCount: 27},
				},
				Regions:                   []string{"GB", "US"},
				AppPackageName:            appPackage,
				TransmissionRisk:          4,
				VerificationAuthorityName: "QRTH-ROWO-LOLO-FOOB",
			},
			Expected: "xH8QNR09EKuCCuNitam1RgjPaGHO/9p54VikqFdirVY=",
		},
		{
			Data: model.Publish{
				Keys: []model.ExposureKey{
					{Key: "zdCW5HrOKbirxmQVc0L/eA", IntervalNumber: 263123, IntervalCount: 144},
					{Key: "t+k51ifogJo9jq3GH9LWGQ", IntervalNumber: 263267, IntervalCount: 144},
					{Key: "3uXRrSlcv1+OMI3oFtdaUw", IntervalNumber: 263411, IntervalCount: 27},
				},
				Regions:                   []string{"gB", "us"},
				AppPackageName:            appPackage,
				TransmissionRisk:          7,
				VerificationAuthorityName: "BREA-KMEO-FFAP-IECE",
			},
			Expected: "LHSwWAjTf3nMVTk7LBwMx9Wg7jEPRjEJf1zRtoxQI64=",
		},
	}

	for _, input := range testData {
		data := NewNonce(&input.Data)
		nonce := data.Nonce()
		if nonce != input.Expected {
			t.Errorf("getNonce: got '%v', want '%v'", nonce, input.Expected)
		}
	}
}
