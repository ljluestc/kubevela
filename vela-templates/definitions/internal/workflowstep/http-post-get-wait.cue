import (
	"vela/op"
	"vela/http"
	"encoding/json"
)

"http-post-get-wait": {
	alias: ""
	attributes: {}
	description: "Send HTTP POST request followed by GET polling (solves issue #6806)"
	annotations: {
		"category": "External Integration"
	}
	labels: {}
	type: "workflow-step"
}

template: {
	// POST request (executed once)
	postReq: http.#HTTPDo & {
		$params: {
			method: parameter.post.method
			url:    parameter.post.url
			request: {
				if parameter.post.body != _|_ {
					body: json.Marshal(parameter.post.body)
				}
				if parameter.post.header != _|_ {
					header: parameter.post.header
				}
				if parameter.post.timeout != _|_ {
					timeout: parameter.post.timeout
				}
			}
		}
	}

	// Parse POST response
	postResponse: json.Unmarshal(postReq.$returns.body)

	// Extract ID from POST response for GET polling
	resourceId: postResponse[parameter.idField]

	// Fail if POST failed
	postCheck: op.#Steps & {
		if postReq.$returns.statusCode >= 400 {
			fail: op.#Fail & {
				message: "POST request failed with status \(postReq.$returns.statusCode): \(postReq.$returns.body)"
			}
		}
	}

	// Build GET URL
	getUrl: "\(parameter.get.baseUrl)/\(resourceId)"

	// GET polling (only polls GET, doesn't re-execute POST)
	poll: op.#Steps & {
		// Initial GET request
		initialGet: http.#HTTPGet & {
			$params: {
				url: getUrl
				request: {
					if parameter.get.header != _|_ {
						header: parameter.get.header
					}
					if parameter.get.timeout != _|_ {
						timeout: parameter.get.timeout
					}
				}
			}
		}

		// Parse initial GET response
		initialGetResponse: json.Unmarshal(initialGet.$returns.body)

		// Check if already complete
		alreadyComplete: initialGetResponse[parameter.get.statusField] == parameter.get.successValue

		// Polling loop (only if not already complete)
		polling: op.#Steps & {
			if !alreadyComplete {
				pollLoop: {
					for i, _ in [ for x in parameter.get.maxAttempts {x} ] {
						"poll-\(i)": op.#Steps & {
							// GET request for this polling attempt
							pollReq: http.#HTTPGet & {
								$params: {
									url: getUrl
									request: {
										if parameter.get.header != _|_ {
											header: parameter.get.header
										}
										if parameter.get.timeout != _|_ {
											timeout: parameter.get.timeout
										}
									}
								}
							}

							// Parse polling response
							pollResponse: json.Unmarshal(pollReq.$returns.body)

							// Check if condition met
							conditionMet: pollResponse[parameter.get.statusField] == parameter.get.successValue

							// If condition met, store result and break
							successResult: op.#Steps & {
								if conditionMet {
									result: pollResponse
								}
							}

							// If not met and not last attempt, wait
							waitNext: op.#Steps & {
								if !conditionMet && i < parameter.get.maxAttempts-1 {
									wait: op.#ConditionalWait & {
										continue: false
										message: "Waiting for \(parameter.get.interval) before polling again (attempt \(i+2)/\(parameter.get.maxAttempts))"
									}
								}
							}

							// If last attempt and still not successful, fail
							failTimeout: op.#Steps & {
								if !conditionMet && i == parameter.get.maxAttempts-1 {
									fail: op.#Fail & {
										message: "Operation did not complete after \(parameter.get.maxAttempts) polling attempts. Final status: \(pollResponse[parameter.get.statusField])"
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Determine final output
	output: op.#Steps & {
		if alreadyComplete {
			// Use initial response
			data: initialGetResponse
		}
		if !alreadyComplete {
			// Use result from polling
			data: poll.polling.pollLoop["poll-\(parameter.get.maxAttempts-1)"].successResult.result
		}
	}

	parameter: {
		// POST request configuration
		post: {
			url:     string
			method: *"POST" | string
			body?:   {...}
			header?: [string]: string
			timeout?: string
		}

		// Field name containing ID in POST response
		idField: string

		// GET polling configuration
		get: {
			baseUrl:      string
			header?:      [string]: string
			timeout?:     string
			statusField:  string
			successValue: string
			maxAttempts: *10 | int
			interval:    *"30s" | string
		}
	}
}
