package drunix

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	gateway    *client.Gateway
	connection *grpc.ClientConn
	network    *client.Network
	contract   *client.Contract
}

type Config struct {
	MSPID         string
	CryptoPath    string
	ChannelName   string
	ChaincodeName string
	PeerEndpoint  string
	GatewayPeer   string
}

func NewClient(config Config) (*Client, error) {
	connection, err := newGRPCConnection(
		config.CryptoPath,
		config.PeerEndpoint,
		config.GatewayPeer,
	)
	if err != nil {
		return nil, err
	}

	id, err := newIdentity(config.MSPID, config.CryptoPath)
	if err != nil {
		connection.Close()
		return nil, err
	}

	sign, err := newSign(config.CryptoPath)
	if err != nil {
		connection.Close()
		return nil, err
	}

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(connection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to connect to Drunix: %w", err)
	}

	network := gateway.GetNetwork(config.ChannelName)
	contract := network.GetContract(config.ChaincodeName)

	return &Client{
		gateway:    gateway,
		connection: connection,
		network:    network,
		contract:   contract,
	}, nil
}

func (c *Client) Close() {
	c.gateway.Close()
	c.connection.Close()
}

func (c *Client) GetAllAssets() ([]byte, error) {
	result, err := c.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		return nil, fmt.Errorf("failed to query GetAllAssets: %w", err)
	}

	return result, nil
}

func newGRPCConnection(
	cryptoPath string,
	peerEndpoint string,
	gatewayPeer string,
) (*grpc.ClientConn, error) {
	tlsCertPath := path.Join(
		cryptoPath,
		"peers",
		"peer0.org1.example.com",
		"tls",
		"ca.crt",
	)

	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TLS certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)

	transportCredentials := credentials.NewClientTLSFromCert(
		certPool,
		gatewayPeer,
	)

	connection, err := grpc.NewClient(
		peerEndpoint,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

func newIdentity(
	mspID string,
	cryptoPath string,
) (*identity.X509Identity, error) {
	certPath := path.Join(
		cryptoPath,
		"users",
		"User1@org1.example.com",
		"msp",
		"signcerts",
	)

	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read identity certificate: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identity certificate: %w", err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	return id, nil
}

func newSign(cryptoPath string) (identity.Sign, error) {
	keyPath := path.Join(
		cryptoPath,
		"users",
		"User1@org1.example.com",
		"msp",
		"keystore",
	)

	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	return sign, nil
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
