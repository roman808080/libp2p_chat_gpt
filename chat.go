package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	// "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

const chatProtocol protocol.ID = "/p2p-chat/1.0.0"

func main() {
	ctx := context.Background()

	// Create a new libp2p host with default options
	h, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// Set up a stream handler for the chat protocol
	h.SetStreamHandler(chatProtocol, handleStream)

	// Get user input for connection target
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the target's multiaddress:")
	target, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// Connect to the target peer
	targetAddr, err := ma.NewMultiaddr(strings.TrimSpace(target))
	if err != nil {
		panic(err)
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}
	err = h.Connect(ctx, *targetInfo)
	if err != nil {
		panic(err)
	}

	// Open a new stream with the target peer
	stream, err := h.NewStream(ctx, targetInfo.ID, chatProtocol)
	if err != nil {
		panic(err)
	}

	// Launch chat routines
	go chatRead(stream)
	go chatWrite(stream)

	// Wait for user input to exit
	fmt.Println("Press ENTER to close the chat")
	reader.ReadString('\n')
}

func handleStream(stream network.Stream) {
	go chatRead(stream)
	go chatWrite(stream)
}

func chatRead(stream network.Stream) {
	reader := bufio.NewReader(stream)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from stream: %s\n", err)
			return
		}
		fmt.Printf("Received: %s", line)
	}
}

func chatWrite(stream network.Stream) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from stdin: %s\n", err)
			return
		}

		_, err = stream.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error writing to stream: %s\n", err)
			return
		}
	}
}
