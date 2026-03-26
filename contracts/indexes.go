/*
Copyright 2025 BubuStack.

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

package contracts

// Index field constants for controller-runtime field indexers.
// These constants define the virtual field paths used by SetupIndexers and
// consumed by controllers via client.MatchingFields lookups.
//
// IMPORTANT: Any change to these values must be coordinated between
// internal/setup/indexing.go (registration) and the controllers (consumption).
// Mismatched values will cause lookups to silently return empty results.
const (
	// Engram indexes
	// IndexEngramTemplateRef indexes Engrams by their referenced EngramTemplate name.
	// Used by: EngramTemplate controller to find dependent Engrams.
	IndexEngramTemplateRef = "spec.templateRef.name"

	// StepRun indexes
	// IndexStepRunEngramRef indexes StepRuns by their Engram reference (namespaced key).
	// Format: "<namespace>/<name>"
	// Used by: StepRun controller for Engram change fan-out.
	IndexStepRunEngramRef = "spec.engramRef.key"

	// IndexStepRunStoryRunRef indexes StepRuns by their parent StoryRun name.
	// Format: "<namespace>/<name>"
	// Note: Currently used for CRD printer columns; no active controller lookups.
	// Retained for potential future use in StoryRun cleanup or status aggregation.
	IndexStepRunStoryRunRef = "spec.storyRunRef.key"

	// StoryRun indexes
	// IndexStoryRunImpulseRef indexes StoryRuns by their triggering Impulse (namespaced key).
	// Used by: Impulse controller to find runs triggered by an Impulse.
	// Format: "<namespace>/<name>"
	IndexStoryRunImpulseRef = "spec.impulseRef.key"

	// IndexStoryRunStoryRefName indexes StoryRuns by their parent Story name only.
	// Used by: Story controller for name-based lookups within a namespace.
	IndexStoryRunStoryRefName = "spec.storyRef.name"

	// IndexStoryRunStoryRefKey indexes StoryRuns by their parent Story (namespaced key).
	// Format: "<namespace>/<name>"
	// Used by: Story controller to count active runs and check deletion eligibility.
	IndexStoryRunStoryRefKey = "spec.storyRef.key"

	// Impulse indexes
	// IndexImpulseTemplateRef indexes Impulses by their referenced ImpulseTemplate name.
	// Used by: ImpulseTemplate controller to find dependent Impulses.
	IndexImpulseTemplateRef = "spec.templateRef.name"

	// IndexImpulseStoryRef indexes Impulses by their referenced Story (namespaced key).
	// Used by: Story controller to find Impulses that trigger a Story.
	// Format: "<namespace>/<name>"
	IndexImpulseStoryRef = "spec.storyRef.key"

	// Story indexes
	// IndexStoryStepEngramRefs indexes Stories by the Engrams referenced in their steps.
	// Format: "<namespace>/<name>" for each Engram reference.
	// Used by: Engram controller to find Stories that use an Engram.
	IndexStoryStepEngramRefs = "spec.steps.ref.key"

	// IndexStoryStepStoryRefs indexes Stories by the Stories referenced in execute-story steps.
	// Format: "<namespace>/<name>" for each Story reference.
	// Used by: Story controller to find Stories that execute another Story.
	IndexStoryStepStoryRefs = "spec.steps.storyRef.key"

	// IndexStoryTransportRefs indexes Stories by their transport references.
	// Used by: Transport controller to find Stories using a Transport.
	IndexStoryTransportRefs = "spec.transports.transportRef"

	// TransportBinding indexes
	// IndexTransportBindingTransportRef indexes TransportBindings by their Transport reference.
	// Used by: Transport controller to find bindings for a Transport.
	IndexTransportBindingTransportRef = "spec.transportRef"

	// EngramTemplate indexes
	// IndexEngramTemplateDescription indexes EngramTemplates by their description.
	// Note: Currently not actively used in controller lookups.
	// Retained for potential future catalog search/filter functionality.
	IndexEngramTemplateDescription = "spec.description"

	// IndexEngramTemplateVersion indexes EngramTemplates by their version.
	// Note: Currently not actively used in controller lookups.
	// Retained for potential future version-based catalog queries.
	IndexEngramTemplateVersion = "spec.version"
)
