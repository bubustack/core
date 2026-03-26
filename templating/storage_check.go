package templating

// HasStorageRefsInStepVars parses a template string, extracts referenced step
// names, and checks whether any of them contain $bubuStorageRef in their output.
func HasStorageRefsInStepVars(text string, vars map[string]any) bool {
	refs, _ := ExtractStepReferencesWithError(text)
	if len(refs) == 0 {
		return false
	}
	steps, ok := lookupMapValue(vars, RootSteps)
	if !ok {
		return false
	}
	for _, ref := range refs {
		stepData, ok := lookupMapValue(steps, ref)
		if !ok {
			continue
		}
		if containsStorageRef(stepData) {
			return true
		}
	}
	return false
}

// containsStorageRef recursively checks if a value contains $bubuStorageRef at any depth.
func containsStorageRef(data any) bool {
	return containsStorageRefImpl(data, make(map[containerVisit]struct{}))
}

func containsStorageRefImpl(data any, visiting map[containerVisit]struct{}) bool {
	if visit, ok := containerVisitFor(data); ok {
		if _, exists := visiting[visit]; exists {
			return false
		}
		visiting[visit] = struct{}{}
		defer delete(visiting, visit)
	}

	switch v := data.(type) {
	case map[string]any:
		if _, ok := v[StorageRefKey]; ok {
			return true
		}
		for _, entry := range v {
			if containsStorageRefImpl(entry, visiting) {
				return true
			}
		}
	case map[any]any:
		if _, ok := v[StorageRefKey]; ok {
			return true
		}
		for _, entry := range v {
			if containsStorageRefImpl(entry, visiting) {
				return true
			}
		}
	case []any:
		for _, entry := range v {
			if containsStorageRefImpl(entry, visiting) {
				return true
			}
		}
	}
	return false
}

func lookupMapValue(container any, key string) (any, bool) {
	switch typed := container.(type) {
	case map[string]any:
		value, ok := typed[key]
		return value, ok
	case map[any]any:
		value, ok := typed[key]
		return value, ok
	default:
		return nil, false
	}
}
