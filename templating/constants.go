package templating

const (
	// TemplateExprKey marks an object value as a template expression wrapper.
	TemplateExprKey = "$bubuTemplate"
	// TemplateVarsKey carries per-expression template variables.
	TemplateVarsKey = "$bubuTemplateVars"

	// StorageRefKey marks a value as offloaded storage-backed data.
	StorageRefKey = "$bubuStorageRef"
	// StoragePathKey identifies the nested path within an offloaded storage object.
	StoragePathKey = "$bubuStoragePath"
)

const (
	// RootInputs is the template root for resolved step inputs.
	RootInputs = "inputs"
	// RootSteps is the template root for upstream step outputs.
	RootSteps = "steps"
	// RootPacket is the template root for transport packet metadata.
	RootPacket = "packet"
)
