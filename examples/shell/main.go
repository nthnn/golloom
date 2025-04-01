package main

import (
	"bufio"   // Package bufio implements buffered I/O.
	"context" // Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values.
	"fmt"     // Package fmt implements formatted I/O with functions analogous to C's printf and scanf.
	"log"     // Package log implements simple logging.
	"os"      // Package os provides a platform-independent interface to operating system functionality.
	"strings" // Package strings implements simple functions to manipulate UTF-8 encoded strings.

	"github.com/nthnn/golloom" // Importing the Golloom package for interacting with language models.
)

func main() {
	// Initialize a new Golloom client to interact with the language model server.
	client, err := golloom.NewClient("http://localhost:11434", 12)
	if err != nil {
		log.Fatalf("Error creating client: %v", err) // Log and exit if client creation fails.
	}

	var history []golloom.Message // Slice to store the conversation history between the user and the assistant.
	ctx := context.Background()   // Create a background context for the API requests.

	reader := bufio.NewReader(os.Stdin) // Initialize a buffered reader to read user input from the standard input.

	for {
		fmt.Print(">> ") // Display the prompt for user input.

		// Read the user's input until a newline character.
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err) // Log and exit if reading input fails.
		}
		input = strings.TrimSpace(input) // Remove any leading or trailing whitespace from the input.
		if input == "exit" {
			break // Exit the loop if the user types "exit".
		}

		// Append the user's message to the conversation history.
		history = append(history, golloom.Message{
			Role:    "user",
			Content: input,
		})

		// Create a chat request with the current conversation history.
		chatReq := &golloom.Chat{
			Model:    "deepseek-r1:14b", // Specify the model to be used for generating responses.
			Messages: history,
		}

		// Send the chat request to the server and receive the response.
		chatResp, err := client.Chat(ctx, chatReq)
		if err != nil {
			log.Fatalf("Chat error: %v", err) // Log and exit if the chat request fails.
		}

		// Extract the assistant's message from the response.
		assistantMessage := chatResp.Message.Content
		fmt.Println(assistantMessage) // Display the assistant's response.

		// Append the assistant's message to the conversation history.
		history = append(history, golloom.Message{
			Role:    "assistant",
			Content: assistantMessage,
		})
	}
}
