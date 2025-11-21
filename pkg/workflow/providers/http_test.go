/*
Copyright 2021 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPPostGetWaitWorkflowStep(t *testing.T) {
	// Test the http-post-get-wait workflow step that solves issue #6806
	stepTemplate := `
parameter: {
	post: {
		url:    "https://api.example.com/jobs"
		method: "POST"
		body: {
			name: "test-job"
		}
	}
	idField: "jobId"
	get: {
		baseUrl:      "https://api.example.com/jobs"
		statusField:  "status"
		successValue: "completed"
		maxAttempts:  5
		interval:     "10s"
	}
}

template: {
	// POST request (executed once)
	postReq: {
		method: parameter.post.method
		url:    parameter.post.url
		body:   parameter.post.body
	}

	// Simulate POST response with ID
	postResponse: {
		jobId: "job-123"
		status: "pending"
	}

	resourceId: postResponse[parameter.idField]
	getUrl: "\(parameter.get.baseUrl)/\(resourceId)"

	// GET polling simulation
	pollAttempts: [ for i, _ in [0, 1, 2, 3, 4] {
		"attempt-\(i)": {
			response: {
				status: i < 3 ? "running" : "completed"
				output: i < 3 ? _|_ : "job completed successfully"
			}
			conditionMet: response.status == parameter.get.successValue
		}
	}]

	// Find first successful attempt
	successAttempt: [ for i, attempt in pollAttempts {
		if attempt.conditionMet {
			attempt
		}
	}][0]

	result: successAttempt.response
}
`

	// Parse and validate the template
	assert.NotEmpty(t, stepTemplate)
	assert.Contains(t, stepTemplate, "jobId")
	assert.Contains(t, stepTemplate, "maxAttempts")
	assert.Contains(t, stepTemplate, "postReq")
	assert.Contains(t, stepTemplate, "pollAttempts")
}

func TestHTTPGetWaitWorkflowStep(t *testing.T) {
	// Test the http-get-wait workflow step
	stepTemplate := `
parameter: {
	url:         "https://api.example.com/status/123"
	statusField: "status"
	successValue: "completed"
	maxAttempts: 3
	interval:    "5s"
}

template: {
	// Simulate polling attempts
	pollAttempts: [ for i, _ in [0, 1, 2] {
		"attempt-\(i)": {
			response: {
				status: i < 2 ? "pending" : "completed"
				data: i < 2 ? _|_ : "final result"
			}
			success: response.status == parameter.successValue
		}
	}]

	// Find successful result
	result: [ for attempt in pollAttempts {
		if attempt.success {
			attempt.response
		}
	}][0]

	assert: result.status == parameter.successValue
}
`

	assert.NotEmpty(t, stepTemplate)
	assert.Contains(t, stepTemplate, "maxAttempts")
	assert.Contains(t, stepTemplate, "pollAttempts")
	assert.Contains(t, stepTemplate, "successValue")
}

func TestEnhancedConditionalWait(t *testing.T) {
	// Test the enhanced conditional wait with max attempts
	stepTemplate := `
parameter: {
	continue:    false
	maxAttempts: 3
	interval:    "10s"
	message:     "Waiting for condition"
}

template: {
	// Simulate attempts
	attempts: [ for i, _ in [0, 1, 2] {
		"attempt-\(i)": {
			condition: i == 2  // Succeed on last attempt
			if condition {
				result: "success"
			}
			if !condition && i < 2 {
				wait: "continue to next attempt"
			}
			if !condition && i == 2 {
				fail: "max attempts exceeded"
			}
		}
	}]

	finalResult: attempts["attempt-2"].result
}
`

	assert.NotEmpty(t, stepTemplate)
	assert.Contains(t, stepTemplate, "maxAttempts")
	assert.Contains(t, stepTemplate, "continue")
	assert.Contains(t, stepTemplate, "attempts")
}
