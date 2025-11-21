import (
	"vela/op"
	"vela/http"
	"encoding/json"
)

"http-get-wait": {
	alias: ""
	attributes: {}
	description: "Send HTTP GET request and wait for condition with configurable polling"
	annotations: {
		"category": "External Integration"
	}
	labels: {}
	type: "workflow-step"
}

template: {
	// Initial GET request (executed once)
	initialReq: http.#HTTPGet & {
		$params: {
			url: parameter.url
			request: {
				if parameter.header != _|_ {
					header: parameter.header
				}
				if parameter.timeout != _|_ {
					timeout: parameter.timeout
				}
			}
		}
	}

	// Parse initial response
	initialResponse: json.Unmarshal(initialReq.$returns.body)

	// Check if we should continue polling (based on initial response)
	shouldContinuePolling: initialResponse[parameter.statusField] != parameter.successValue

	// Polling logic (only executed if shouldContinuePolling is true)
	poll: op.#Steps & {
		if shouldContinuePolling {
			// Loop with max attempts and custom interval
			pollLoop: {
				for i, _ in [ for x in parameter.maxAttempts {x} ] {
					"poll-\(i)": op.#Steps & {
						// GET request for polling
						pollReq: http.#HTTPGet & {
							$params: {
								url: parameter.url
								request: {
									if parameter.header != _|_ {
										header: parameter.header
									}
									if parameter.timeout != _|_ {
										timeout: parameter.timeout
									}
								}
							}
						}

						// Parse polling response
						pollResponse: json.Unmarshal(pollReq.$returns.body)

						// Check if condition is met
						conditionMet: pollResponse[parameter.statusField] == parameter.successValue

						// If condition met, break out of loop
						breakLoop: op.#Steps & {
							if conditionMet {
								result: pollResponse
							}
						}

						// If not the last attempt and condition not met, wait
						waitNext: op.#Steps & {
							if i < parameter.maxAttempts-1 && !conditionMet {
								wait: op.#ConditionalWait & {
									continue: false  // Wait for specified interval
									message: "Waiting \(parameter.interval) before next poll (attempt \(i+2)/\(parameter.maxAttempts))"
								}
							}
						}

						// If this is the last attempt and still not successful, fail
						failMaxAttempts: op.#Steps & {
							if i == parameter.maxAttempts-1 && !conditionMet {
								fail: op.#Fail & {
									message: "Max polling attempts (\(parameter.maxAttempts)) exceeded. Final status: \(pollResponse[parameter.statusField])"
								}
							}
						}
					}
				}
			}
		}
	}

	// Determine final result
	result: op.#Steps & {
		if shouldContinuePolling {
			// Use result from polling
			if poll.pollLoop["poll-\(parameter.maxAttempts-1)"].breakLoop.result != _|_ {
				data: poll.pollLoop["poll-\(parameter.maxAttempts-1)"].breakLoop.result
			}
		}
		if !shouldContinuePolling {
			// Use initial response
			data: initialResponse
		}
	}

	// Check for HTTP errors
	httpCheck: op.#Steps & {
		if initialReq.$returns.statusCode >= 400 {
			fail: op.#Fail & {
				message: "HTTP request failed with status \(initialReq.$returns.statusCode)"
			}
		}
	}

	parameter: {
		// URL to poll
		url: string

		// HTTP headers
		header?: [string]: string

		// Request timeout
		timeout?: string

		// Field name to check in response JSON
		statusField: string

		// Expected value for success
		successValue: string

		// Maximum number of polling attempts
		maxAttempts: *10 | int

		// Polling interval (in ConditionalWait message, actual timing controlled by workflow engine)
		interval: *"30s" | string
	}
}