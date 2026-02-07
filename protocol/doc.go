// Package protocol defines common interfaces used throughout the core package.
//
// These interfaces enable dependency injection and mocking for testing.
//
// Example - implementing Lifecycle:
//
//	type MyService struct{}
//
//	func (s *MyService) Start(ctx context.Context) error { return nil }
//	func (s *MyService) Stop(ctx context.Context) error { return nil }
package protocol
