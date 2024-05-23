# typechat-go

typechat-go es una biblioteca que facilita la construcción de interfaces de lenguaje natural utilizando tipos.

Este proyecto es una implementación en Go del proyecto original [TypeChat](https://github.com/microsoft/TypeChat) de Microsoft.

Visite https://microsoft.github.io/TypeChat para obtener más información sobre lo que le permite hacer.

Esta implementación sigue de manera flexible lo que hace la versión 0.10 de TypeChat con ergonomía ligeramente diferente más apropiada para Go.

Algunas de las diferencias clave son que esta biblioteca tiene menos opiniones sobre cómo se comunica con el LLM, y siempre que proporcione un cliente válido, puede usarlo fácilmente.

## Cómo usar

### Prompt + Tipo de Retorno

Esta funcionalidad le permite pasar un prompt de lenguaje natural y el tipo de resultado esperado que desea que el LLM use al responder. Por ejemplo:
```go
type Classifier struct {
    Sentiment string
}

ctx := context.Background()

// Proporcione un cliente de modelo que implemente la interfaz requerida
// es decir, Do(ctx context.Context, prompt string) (response []byte, err error)
// Este modelo puede llamar a OpenAPI, Azure o cualquier LLM. Usted controla el transporte.
model := ... 

prompt := typechat.NewPrompt[Classifier](model, "¡Hoy es un buen día!")
result, err := prompt.Execute(ctx)
if err != nil {
    ...
}

fmt.Println(result.Sentiment) // proporcionado por el LLM
```

Notará que esta implementación está utilizando Generics, por lo que el resultado que obtiene del LLM está completamente tipificado y puede ser utilizado por el resto de su aplicación.

### Prompt + Programa

Esta funcionalidad le permite pasar un prompt de lenguaje natural junto con una interfaz de comportamiento que su aplicación admite. La biblioteca hará que el LLM genere una secuencia de pasos que considera necesarios para realizar una tarea dada.

```go
type API interface {
    CreateTweet(message string)
    CreateLinkedInMessage(message string)
}

ctx := context.Background()

// Proporcione un cliente de modelo que implemente la interfaz requerida
// es decir, Do(ctx context.Context, prompt string) (response []byte, err error)
// Este modelo puede llamar a OpenAPI, Azure o cualquier LLM. Usted controla el transporte.
model := ... 

prompt := typechat.NewPrompt[API](model, "¡Realmente necesito twittear y publicar en mi LinkedIn que he sido promovido!")
program, err := prompt.CreateProgram(ctx)
if err != nil {
    ...
}

// El programa contendrá las invocaciones necesarias que su aplicación tiene que hacer con la API proporcionada para realizar la tarea identificada por el LLM.
program.Steps[0].Name == "CreateTweet"
program.Steps[0].Args == []any{"¡He sido promovido!"}

program.Steps[1].Name == "CreateLinkedInMessage"
program.Steps[1].Args == []any{"¡He sido promovido!"}

// Puede construir un ejecutor de programas sobre esta estructura.
```

### Ejemplo de Manejo de Errores

Al trabajar con servicios o API externos, es crucial manejar los errores de manera elegante. A continuación, se muestra un ejemplo de cómo manejar errores al usar el método `Execute` de la estructura `Prompt`.

```go
ctx := context.Background()
model := ... // su cliente de modelo

prompt := typechat.NewPrompt[Classifier](model, "Analizar el sentimiento de este texto.")
result, err := prompt.Execute(ctx)
if err != nil {
    fmt.Printf("Ocurrió un error: %s\n", err)
    return
}

fmt.Printf("Resultado del análisis de sentimiento: %s\n", result.Sentiment)
```

Este ejemplo demuestra cómo capturar y manejar errores devueltos por el método `Execute`, asegurando que su aplicación pueda responder adecuadamente a los fallos.

### Ejemplo de Adaptador Personalizado

Para usar un adaptador personalizado con la biblioteca, necesita crear un adaptador que implemente la interfaz `client`. A continuación, se muestra un ejemplo de cómo crear un adaptador personalizado y usarlo con `NewPrompt`.

```go
type MyCustomAdapter struct {
    // Campos y métodos personalizados
}

func (m *MyCustomAdapter) Do(ctx context.Context, prompt []Message) (string, error) {
    // Implemente la lógica para enviar el prompt a su servicio y devolver la respuesta
    return "respuesta de su servicio", nil
}

ctx := context.Background()
adapter := &MyCustomAdapter{}

prompt := typechat.NewPrompt[Classifier](adapter, "Analizar este texto con mi adaptador personalizado.")
result, err := prompt.Execute(ctx)
if err != nil {
    fmt.Printf("Ocurrió un error: %s\n", err)
    return
}

fmt.Printf("Resultado del adaptador personalizado: %s\n", result.Sentiment)
```

Este ejemplo demuestra cómo crear un adaptador personalizado que implementa la interfaz `client` y usarlo con `NewPrompt` para enviar prompts a su servicio personalizado.

## Contribuyendo

Esta biblioteca está en desarrollo y aún requiere más trabajo para solidificar las API proporcionadas, así que úsela con precaución. Se realizará un lanzamiento en un futuro cercano.

### Tareas Pendientes

- Más pruebas
- Proporcionar un paquete de adaptador para crear de manera predeterminada con clientes de OpenAI y Azure
- Otros casos de uso como conversaciones
- Descubrir la mejor manera de mantenerse sincronizado con el proyecto original de TypeChat
- Configurar CI

¡Participe! ¡Mucho por hacer!

[Read this in English](README.md)
