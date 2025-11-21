// Test cases for HTTP workflow steps that solve issue #6806

// Test http-get-wait step
httpGetWaitTest: {
	template: {
		// This demonstrates the solution to issue #6806
		// Only the GET request is polled, POST is not re-executed

		pollResult: {
			"http-get-wait": {
				parameter: {
					url:         "https://api.example.com/status/123"
					header: {
						"Authorization": "Bearer token123"
						"Content-Type":  "application/json"
					}
					timeout:     "30s"
					statusField: "status"
					successValue: "completed"
					maxAttempts: 5
					interval:    "10s"
				}
			}
		}

		// Validate the result
		result: pollResult["http-get-wait"].result
		assert: result.status == "completed"
	}
}

// Test http-post-get-wait step (complete solution for issue #6806)
httpPostGetWaitTest: {
	template: {
		// This is the complete solution for the issue described in #6806
		// POST is executed once, then only GET is polled

		asyncOperation: {
			"http-post-get-wait": {
				parameter: {
					post: {
						url:    "https://api.example.com/jobs"
						method: "POST"
						body: {
							name:    "data-processing-job"
							type:    "async"
							priority: "high"
						}
						header: {
							"Authorization": "Bearer token123"
							"Content-Type":  "application/json"
						}
						timeout: "30s"
					}
					idField: "jobId"
					get: {
						baseUrl:     "https://api.example.com/jobs"
						header: {
							"Authorization": "Bearer token123"
							"Content-Type":  "application/json"
						}
						timeout:     "30s"
						statusField: "status"
						successValue: "completed"
						maxAttempts: 10
						interval:    "30s"
					}
				}
			}
		}

		// Validate the final result
		finalResult: asyncOperation["http-post-get-wait"].output.data
		assert: finalResult.status == "completed"
		assert: finalResult.output != _|_
	}
}

// Test enhanced conditional wait
enhancedConditionalWaitTest: {
	template: {
		waitForCondition: {
			"enhanced-conditional-wait": {
				parameter: {
					continue:    false  // Wait condition
					maxAttempts: 5
					interval:    "10s"
					message:     "Waiting for external condition"
				}
			}
		}
	}
}

// Example usage that solves the original issue #6806
issue6806Solution: {
	// Original problematic workflow (re-executes POST on every poll)
	originalWorkflow: {
		post: {
			// POST executed on every polling cycle - PROBLEM!
			httpPost: {
				// ... POST logic
			}
		}
		get: {
			// GET polling
			wait: {
				// ... polling logic that re-executes entire workflow
			}
		}
	}

	// Fixed workflow using new http-post-get-wait step
	fixedWorkflow: {
		asyncOperation: {
			"http-post-get-wait": {
				parameter: {
					post: {
						url:    "https://api.example.com/submit-job"
						method: "POST"
						body: {
							name: "my-job"
						}
					}
					idField: "id"
					get: {
						baseUrl:      "https://api.example.com/job"
						statusField:  "status"
						successValue: "done"
						maxAttempts:  20
						interval:     "30s"
					}
				}
			}
		}

		// Use the result
		result: asyncOperation["http-post-get-wait"].output.data
	}
}
