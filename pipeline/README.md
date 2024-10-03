# pipeline

Pipeline executes sequence of functions with simple syntax. It acts as a wrapper for calling functions.

Example:

```go
errCh := make(chan error)
go pipeline.New(context.Background()).
    Before(func() { fmt.Println("1. before") }).
    Then(func(context.Context) error {
        fmt.Println("2. then")
        return errors.New("sample error")
    }).
    ThenCatch(func(err error) error {
        fmt.Println("3. then catch")
        return err
    }).
    Else(func(context.Context) error {
        fmt.Println("4. else")
        return errors.New("sample error")
    }).
    ElseCatch(func(err error) error {
        fmt.Println("5. else catch")
        return err
    }).
    After(func() { fmt.Println("6. after") }).
    Run(func(err error) {
        fmt.Println("7. run")
        errCh <- err
    })
fmt.Println(<-errCh)
```
