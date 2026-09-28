package main

import (
	"fmt"
	"log"

	"arealis-drunix-hackathon/drunix"
)

func main() {
	fmt.Println("Connecting to Drunix...")

	client, err := drunix.NewClient(drunix.Config{
		MSPID:         "Org1MSP",
		CryptoPath:    "../drunix/drunix-network/test-network/organizations/peerOrganizations/org1.example.com",
		ChannelName:   "mychannel",
		ChaincodeName: "basic",
		PeerEndpoint:  "dns:///localhost:7051",
		GatewayPeer:   "peer0.org1.example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Connected to Drunix successfully!")

	result, err := client.GetAllAssets()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GetAllAssets:")
	fmt.Println(string(result))
}
