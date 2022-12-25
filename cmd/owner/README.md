# Step 1 - create a DAPR service

Use `gitlab.com/dapr/go-sdk/service/http` 

**NOTE**: At the time of writing (Dapr 1.9.5), the Go SDK does not support GRPC actors. Therefore, *do not* use `gitlab.com/dapr/go-sdk/service/grpc`.

# Step 2 - create an Actor implementation factory

An actor implementation factory returns instances of `actor.Server`, as per example below. The factory method allows developers to encapsulate the initialization logic.

```go
func ActorServerFactory() actor.Server {
	return &OwnerActor{}
}
```

# Step 3 - add the Type() method

The Type() method allows the actor to be bound to the /actors/\<name\> token.

# Step 4 - what are the valid method signatures?

```go
func Method(ctx context.Context, request any) (any, error)
```