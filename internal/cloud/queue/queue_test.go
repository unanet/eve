package queue_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.unanet.io/devops/eve/internal/cloud/queue"
)

func TestCreateMessage(t *testing.T) {
	output, err := queue.CreateMessage("qa-2020.1", "message 10 test", "10")
	require.NoError(t, err)
	fmt.Println(output)
}

func TestReceiveMessage(t *testing.T) {
	output, err := queue.ReceiveMessage()
	require.NoError(t, err)
	fmt.Println(output)
}

func TestDeleteMessage(t *testing.T) {
	output, err := queue.DeleteMessage("AQEB3h1DF6sppiVEt7LOKy1qT4cL4kxRmnXIdCIbn6s6ZqwUdRXJfyrMfAHAWaD1ZRWLVsbvMkevwmaF6vKAJE7P9wpOC8uPRcY9JgW95m2G9nBUz88kQwH9t5YDeTvI0ikWw11m9DodfiwaXlfPgYKh3K9QA1eFV8G7WMVbwV254QRcI8Iv0XUS7IlcSviOOnikBteL5jaNmjgSH3KRK5egNKjjzyPaMUivOo8DW3GdVGywvQpQtd48hJR4pnhDR8/F1shl1L06l+pe3zYivBXQhv3IMJ4NL0CzCytwsyDaFF0=")
	require.NoError(t, err)
	fmt.Println(output)
}
