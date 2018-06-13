package mongo

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	ot "github.com/opentracing/opentracing-go"
	otExt "github.com/opentracing/opentracing-go/ext"
	otLog "github.com/opentracing/opentracing-go/log"
)

type dbTrace struct {
	// The Name of the struct who
	Manager    string
	Method     string
	Selector   interface{}
	Query      interface{}
	CustomTags []ot.Tag
}

// traceMongoCall provides an abstraction from opentracing to obtain a span
// with relevant details when tracing call time to mongoDB
func traceMongoCall(ctx context.Context, trace dbTrace) (ot.Span, context.Context) {
	// Build a new OpenTracing Child span to track how long it takes for mongo
	// to complete the operation.
	opName := fmt.Sprintf("storage.mongo.%s.%s", trace.Manager, trace.Method)
	span, ctx := ot.StartSpanFromContext(ctx, opName)

	// Tag component details
	otExt.Component.Set(span, "storage")
	otExt.DBType.Set(span, "mongo")

	// Set the DB selector if provided.
	// Generally useful for mongo updates where a selector is applied, then the
	// payload supplied updates the given selected document. For example, the
	// selector could end up selecting an inner document to be updated.
	if trace.Selector != nil {
		span.SetTag("db.selector", fmt.Sprintf("%#+v", trace.Selector))
	}

	// Set the DB query if provided.
	// Generally speaking, the query may not be needed, but may be helpful in
	// debugging errors, therefore it is better advised to log the query out if
	// an error occurs.
	if trace.Query != nil {
		otExt.DBStatement.Set(span, fmt.Sprintf("%#+v", trace.Query))
	}

	// Set the custom tags if provided
	for _, tag := range trace.CustomTags {
		tag.Set(span)
	}

	return span, ctx
}

// otLogQuery given a span and a query,
func otLogQuery(span ot.Span, query interface{}) {
	otExt.DBStatement.Set(span, fmt.Sprintf("%#+v", query))
}

// otLogErr given a span, logs out the error
func otLogErr(span ot.Span, err error) {
	span.LogFields(otLog.Error(err))
}
