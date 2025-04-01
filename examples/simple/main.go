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
	// The client connects to the server at "http://localhost:11434" with a timeout of 12 seconds.
	client, err := golloom.NewClient("http://localhost:11434", 12)
	if err != nil {
		log.Fatalf("Error creating client: %v", err) // Log and exit if client creation fails.
	}

	var history []golloom.Message // Slice to store the history of messages exchanged in the chat.
	ctx := context.Background()   // Create a background context for the chat operations.

	reader := bufio.NewReader(os.Stdin) // Create a buffered reader to read input from the standard input (console).

	for {
		fmt.Print(">> ") // Display the prompt for user input.

		// Read the user's input until a newline character is encountered.
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err) // Log and exit if reading input fails.
		}
		input = strings.TrimSpace(input) // Trim any leading or trailing whitespace from the input.

		// Check if the user wants to exit the chat.
		if input == "exit" {
			break // Exit the loop and end the program.
		}

		// Append the user's message to the chat history.
		history = append(history, golloom.Message{
			Role:    "user", // Role of the message sender.
			Content: input,  // Content of the user's message.
		})

		// Create a new chat request with the current history of messages.
		chatReq := &golloom.Chat{
			Model:    "deepseek-r1:14b", // Specify the model to use for generating responses.
			Messages: history,           // Include the chat history in the request.
		}

		// Send the chat request to the server and receive a response.
		chatResp, err := client.Chat(ctx, chatReq)
		if err != nil {
			log.Fatalf("Chat error: %v", err) // Log and exit if the chat request fails.
		}

		// Extract the assistant's message from the chat response.
		assistantMessage := chatResp.Message.Content
		fmt.Println(assistantMessage) // Display the assistant's response to the user.

		// Append the assistant's message to the chat history.
		history = append(history, golloom.Message{
			Role:    "assistant",      // Role of the message sender.
			Content: assistantMessage, // Content of the assistant's message.
		})
	}
}
