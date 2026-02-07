// Package pipeline provides a fluent API for executing functions in a structured order
// with error handling, fallbacks, and cleanup.
//
// Execution order: Before → Then/Else → ThenCatch/ElseCatch → Error/NoError → After
//
// Example:
//
//	errCh := make(chan error)
//	go pipeline.New(ctx).
//	    Before(func() { /* setup */ }).
//	    Then(func(ctx context.Context) error {
//	        // main logic
//	        return nil
//	    }).
//	    Else(func(ctx context.Context) error {
//	        // fallback on error
//	        return nil
//	    }).
//	    After(func() { /* cleanup */ }).
//	    Run(func(err error) { errCh <- err })
package pipeline
