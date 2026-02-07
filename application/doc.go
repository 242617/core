// Package application provides component lifecycle management with graceful shutdown.
//
// It manages multiple components as a cohesive unit, handling startup and shutdown
// with configurable timeouts. Components implement the protocol.Lifecycle interface.
//
// Example:
//
//	app, err := application.New(
//	    application.WithName("my-service"),
//	    application.WithComponents(
//	        application.NewLifecycleComponent("db", db),
//	    ),
//	)
//	app.Run(context.Background())
package application
