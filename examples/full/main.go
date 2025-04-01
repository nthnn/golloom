package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nthnn/golloom"
)

func main() {
	// Define and parse command-line flags
	baseURL := flag.String("url", "http://localhost:11434", "Base URL for the Ollama server")
	timeout := flag.Int("timeout", 5, "HTTP client timeout in minutes")

	// Parse the flags provided in the command line
	flag.Parse()
	// Check if any command is provided in the arguments
	if flag.NArg() < 1 {
		printUsage() // Print usage details if no command is provided
		os.Exit(1)
	}

	// Get the first command-line argument as the command to execute
	command := flag.Arg(0)

	// Initialize a new client with the base URL and timeout
	client, err := golloom.NewClient(*baseURL, time.Duration(*timeout))

	if err != nil {
		// Handle error if the client initialization fails
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Set up a context for the HTTP request (to manage cancellations and timeouts)
	ctx := context.Background()

	// Switch statement to determine which function to call based on the provided command
	switch command {
	case "version":
		// Fetch and display the version of the Ollama server
		doVersion(ctx, client)

	case "list":
		// List all available models on the server
		doList(ctx, client)

	case "chat":
		// Start a chat session with the server using the provided arguments
		doChat(ctx, client, flag.Args()[1:])

	case "generate":
		// Generate a response based on a prompt using the server
		doGenerate(ctx, client, flag.Args()[1:])

	default:
		// Handle an unknown command
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// printUsage prints out the usage instructions for the program
func printUsage() {
	// Show usage details for global options and available commands
	fmt.Println("Usage: gollama [global options] <command> [command options]")
	fmt.Println("Global options:")
	fmt.Println("  -url string")
	fmt.Println("        Base URL for the gollama server (default \"http://localhost:8080\")")
	fmt.Println("  -timeout int")
	fmt.Println("        HTTP client timeout in minutes (default 5)")
	fmt.Println("\nCommands:")
	fmt.Println("  version               Show server version")
	fmt.Println("  list                  List available models")
	fmt.Println("  chat -model <model> -message <message>")
	fmt.Println("                        Send a chat message to a model")
	fmt.Println("  generate -model <model> -prompt <prompt>")
	fmt.Println("                        Generate text from a prompt")
}

// doVersion fetches and prints the server's version and build time
func doVersion(ctx context.Context, client *golloom.Client) {
	ver, err := client.Version(ctx) // Get version info from the client
	if err != nil {
		// If there's an error, print it and exit
		fmt.Printf("Error fetching version: %v\n", err)
		os.Exit(1)
	}

	// Print the fetched version and build time of the server
	fmt.Printf("Ollama Server Version: %s (Build Time: %s)\n", ver.Version, ver.BuildTime)
}

// doList fetches and prints a list of available models on the server
func doList(ctx context.Context, client *golloom.Client) {
	list, err := client.ListModels(ctx) // Get the list of models from the client
	if err != nil {
		// Handle error in fetching models
		fmt.Printf("Error listing models: %v\n", err)
		os.Exit(1)
	}

	// Print the details of each model available on the server
	fmt.Println("Models:")
	for _, model := range list.Models {
		// Display model details: name, digest, and modification date
		fmt.Printf("- %s (Digest: %s, Modified: %s)\n", model.Name, model.Digest, model.ModifiedAt)
	}
}

// doChat handles sending a chat message to a model and prints the response
func doChat(ctx context.Context, client *golloom.Client, args []string) {
	// Define flags for the chat command (model and message options)
	chatFlags := flag.NewFlagSet("chat", flag.ExitOnError)
	model := chatFlags.String("model", "default", "Model to use for chat")
	message := chatFlags.String("message", "", "Message to send")
	chatFlags.Parse(args)

	// If the message is not provided, print usage instructions and exit
	if *message == "" {
		fmt.Println("Please provide a message using the -message flag.")
		chatFlags.Usage()
		os.Exit(1)
	}

	// Prepare the chat request with the specified model and message
	chatReq := &golloom.Chat{
		Model: *model,
		Messages: []golloom.Message{
			{
				Role:    "user",
				Content: *message,
			},
		},
	}

	// Marshal the chat request into JSON format
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		// Handle error in marshaling the request
		fmt.Printf("Error marshaling chat request: %v\n", err)
		os.Exit(1)
	}

	// Build the full URL for the chat API endpoint
	reqURL := client.BaseURL.ResolveReference(&url.URL{Path: "/api/chat"}).String()

	// Create a new HTTP POST request for the chat session
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(reqBody))
	if err != nil {
		// Handle error creating the HTTP request
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json") // Set the correct content type for the request

	// Execute the HTTP request using the client
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		// Handle error during the HTTP request
		fmt.Printf("Error during chat request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Set up a JSON decoder to handle the streaming response
	decoder := json.NewDecoder(resp.Body)
	var fullContent string

	// Loop through the streamed response chunks and print the content
	for {
		var chunk golloom.ModelResponse
		err := decoder.Decode(&chunk)
		if err != nil {
			// Handle errors in decoding, particularly EOF (end of stream)
			if err == io.EOF {
				break // Exit loop if the stream ends
			}
			fmt.Printf("Error decoding chat stream: %v\n", err)
			os.Exit(1)
		}
		// Print each chunk of content as it arrives
		fmt.Print(chunk.Message.Content)
		fullContent += chunk.Message.Content
	}
	fmt.Println() // Ensure a newline after printing the full response
}

// doGenerate handles generating text based on a prompt
func doGenerate(ctx context.Context, client *golloom.Client, args []string) {
	// Define flags for the generate command (model and prompt options)
	genFlags := flag.NewFlagSet("generate", flag.ExitOnError)
	model := genFlags.String("model", "default", "Model to use for generation")
	prompt := genFlags.String("prompt", "", "Prompt for text generation")
	genFlags.Parse(args)

	// If no prompt is provided, print usage instructions and exit
	if *prompt == "" {
		fmt.Println("Please provide a prompt using the -prompt flag.")
		genFlags.Usage()
		os.Exit(1)
	}

	// Prepare the prompt request with the specified model and prompt
	promptInfo := &golloom.PromptInfo{
		Model:  *model,
		Prompt: *prompt,
	}

	// Call the client to generate a response based on the prompt
	genResp, err := client.Generate(ctx, promptInfo)
	if err != nil {
		// Handle error in generating the response
		fmt.Printf("Error during generation: %v\n", err)
		os.Exit(1)
	}

	// Print the generated response
	fmt.Printf("Generated response:\n%s\n", genResp.Response)
}
