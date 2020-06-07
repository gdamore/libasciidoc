package sgml

// ContextualPipeline carries the renderer context along with
// the pipeline data to process in a template or in a nested template
type ContextualPipeline struct {
	Context *Context
	// The actual pipeline
	Data interface{}
}
