# typechat-go

typechat-go is a library that makes it easy to build natural language interfaces using types.

This project is a Go implementation of the original project [TypeChat](https://github.com/microsoft/TypeChat) by Microsoft.

Visit https://microsoft.github.io/TypeChat for more information on what it enables you do.

This implementation loosely follows what version 0.10 of TypeChat does with slightly different ergonomics more appropiate for Go.

Some of the key differences are that this library has less opinions about how you communicate with the LLM, and so as long as you provide a valid client you can easily use this.

## How to use

### Prompt + Return Type

This functionality allows you to pass in a natural language prompt and the expected result type you wish the LLM to use when replying. For example:
```go
type Classifier struct {
    Sentiment string
}

ctx := context.Background()

// Provide a model client that implements the required interface
// i.e. Do(ctx context.Context, prompt string) (response []byte, err error)
// This model can call to OpenAPI, Azure or any LLM. You control the transport.
model := ... 

prompt := typechat.NewPrompt[Classifier](model, "Today is a good day!")
result, err := prompt.Execute(ctx)
if err != nil {
    ...
}

fmt.Println(result.Sentiment) // provided by the LLM
```

You'll notice that this implementation is using Generics, so the result you get from the LLM is fully typed and able to be uused by the rest of your application.

### Prompt + Program

This functionality allows you to pass in a natural language prompt along with an interface of behavior that your application supports. The library will have the LLM generate a sequence of steps it deems necessary to accomplish a given task.

```go
type API interface {
    CreateTweet(message string)
    CreateLinkedInMessage(message string)
}

ctx := context.Background()

// Provide a model client that implements the required interface
// i.e. Do(ctx context.Context, prompt string) (response []byte, err error)
// This model can call to OpenAPI, Azure or any LLM. You control the transport.
model := ... 

prompt := typechat.NewPrompt[API](model, "I really need to tweet and post on my LinkedIN that I've been promoted!")
program, err := prompt.CreateProgram(ctx)
if err != nil {
    ...
}

// Program will contain the necessary invocations your application has to do with the provided API to accomplish the task as idenfitied by the LLM.
program.Steps[0].Name == "CreateTweet"
program.Steps[0].Args == []any{"I have been promoted!"}

program.Steps[1].Name == "CreateLinkedInMessage"
program.Steps[1].Args == []any{"I have been promoted!"}

// You can build a program executor on top of this structure.
```

### Error Handling Example

When working with external services or APIs, it's crucial to handle errors gracefully. Below is an example of how to handle errors when using the `Execute` method of the `Prompt` struct.

```go
ctx := context.Background()
model := ... // your model client

prompt := typechat.NewPrompt[Classifier](model, "Analyze the sentiment of this text.")
result, err := prompt.Execute(ctx)
if err != nil {
    fmt.Printf("An error occurred: %s\n", err)
    return
}

fmt.Printf("Sentiment analysis result: %s\n", result.Sentiment)
```

This example demonstrates catching and handling errors returned by the `Execute` method, ensuring that your application can respond appropriately to failures.

### Custom Adapter Example

To use a custom adapter with the library, you need to create an adapter that implements the `client` interface. Below is an example of how to create a custom adapter and use it with `NewPrompt`.

```go
type MyCustomAdapter struct {
    // Custom fields and methods
}

func (m *MyCustomAdapter) Do(ctx context.Context, prompt []Message) (string, error) {
    // Implement the logic to send the prompt to your service and return the response
    return "response from your service", nil
}

ctx := context.Background()
adapter := &MyCustomAdapter{}

prompt := typechat.NewPrompt[Classifier](adapter, "Analyze this text with my custom adapter.")
result, err := prompt.Execute(ctx)
if err != nil {
    fmt.Printf("An error occurred: %s\n", err)
    return
}

fmt.Printf("Custom adapter result: %s\n", result.Sentiment)
```

This example demonstrates creating a custom adapter that implements the `client` interface and using it with `NewPrompt` to send prompts to your custom service.

## Contributing

This library is under development and still requires more work to solidify the provided APIs so use with caution. A release will be done at some point in the near future.

### TODOs

- More tests
- Provide an adapter package to create out of the box with OpenAI and Azure clients
- Other use cases such as conversations
- Figure out best way to stay in sync with the original TypeChat project
- Setup CI

Get involved! Lots todo!
