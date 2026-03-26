package templating

import "fmt"

// ErrEvaluationBlocked indicates evaluation is waiting on upstream data.
type ErrEvaluationBlocked struct {
	Reason string
}

func (e *ErrEvaluationBlocked) Error() string {
	if e == nil {
		return "template evaluation blocked"
	}
	if e.Reason == "" {
		return "template evaluation blocked"
	}
	return fmt.Sprintf("template evaluation blocked: %s", e.Reason)
}

// ErrOffloadedDataUsage indicates a template requires offloaded payload data.
type ErrOffloadedDataUsage struct {
	Reason string
}

func (e *ErrOffloadedDataUsage) Error() string {
	if e == nil {
		return "offloaded data usage detected"
	}
	if e.Reason == "" {
		return "offloaded data usage detected"
	}
	return fmt.Sprintf("offloaded data usage detected: %s", e.Reason)
}
