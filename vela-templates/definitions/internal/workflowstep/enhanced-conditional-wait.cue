import (
	"vela/op"
)

"enhanced-conditional-wait": {
	alias: ""
	attributes: {}
	description: "Wait for condition with configurable polling parameters"
	annotations: {
		"category": "Control Flow"
	}
	labels: {}
	type: "workflow-step"
}

template: {
	// Track attempts
	attempts: parameter.attempts

	// Main waiting logic with attempt counting
	waitLoop: {
		for i, _ in [ for x in parameter.maxAttempts {x} ] {
			"attempt-\(i)": op.#Steps & {
				// Check condition
				condition: parameter.continue

				// If condition met, succeed
				success: op.#Steps & {
					if condition {
						result: "success"
					}
				}

				// If condition not met and not max attempts, wait and continue
				continueWaiting: op.#Steps & {
					if !condition && i < parameter.maxAttempts-1 {
						wait: op.#ConditionalWait & {
							continue: false
							message: parameter.message
						}
					}
				}

				// If max attempts reached and condition still not met, fail
				failTimeout: op.#Steps & {
					if !condition && i == parameter.maxAttempts-1 {
						fail: op.#Fail & {
							message: "Condition not met after \(parameter.maxAttempts) attempts with interval \(parameter.interval)"
						}
					}
				}
			}
		}
	}

	parameter: {
		// Condition to wait for
		continue: bool

		// Maximum number of attempts
		maxAttempts: *10 | int

		// Interval between attempts (informational)
		interval: *"30s" | string

		// Optional message
		message?: string
	}
}
