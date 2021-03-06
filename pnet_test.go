package ipfscluster

import "testing"

func TestClusterKeyFormat(t *testing.T) {
    goodKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    tooShort := "0123456789abcdef"
    tooLong := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0"
    unsupportedChars := "0123456789abcdef0123456789!!!!!!0123456789abcdef0123456789abcdef"

	_, err := DecodeClusterSecret(goodKey)
	if err != nil {
        t.Fatal("Failed to decode well-formatted key.")
	}
	_, err = DecodeClusterSecret(tooShort)
	if err == nil {
        t.Fatal("Successfully decoded key that should haved failed (too short).")
	}
	_, err = DecodeClusterSecret(tooLong)
	if err == nil {
        t.Fatal("Successfully decoded key that should haved failed (too long).")
	}
	_, err = DecodeClusterSecret(unsupportedChars)
	if err == nil {
        t.Fatal("Successfully decoded key that should haved failed (unsupported chars).")
	}
}

func TestSimplePNet(t *testing.T) {
	clusters, mocks := peerManagerClusters(t)
	defer cleanRaft()
	defer shutdownClusters(t, clusters, mocks)

	if len(clusters) < 2 {
		t.Skip("need at least 2 nodes for this test")
	}

	_, err := clusters[0].PeerAdd(clusterAddr(clusters[1]))
	if err != nil {
		t.Fatal(err)
	}

	if len(clusters[0].Peers()) != len(clusters[1].Peers()) {
		t.Fatal("Expected same number of peers")
	}
	if len(clusters[0].Peers()) != 2 {
		t.Fatal("Expected 2 peers")
	}
}

func TestClusterSecretRequired(t *testing.T) {
    cl1Secret, err := generateClusterSecret()
    if err != nil {
        t.Fatal("Unable to generate cluster secret.")
    }
	cl1, _ := createOnePeerCluster(t, 1, cl1Secret)
	cl2, _ := createOnePeerCluster(t, 2, testingClusterSecret)
	defer cleanRaft()
	defer cl1.Shutdown()
	defer cl2.Shutdown()
	peers1 := cl1.Peers()
	peers2 := cl2.Peers()

	_, err = cl1.PeerAdd(clusterAddr(cl2))
	if err == nil {
		t.Fatal("Peer entered private cluster without key.")
	}

	if len(peers1) != len(peers2) {
		t.Fatal("Expected same number of peers")
	}
	if len(peers1) != 1 {
		t.Fatal("Expected no peers other than self")
	}
}
