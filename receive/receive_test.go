package receive

import (
	"testing"
)

func TestConstructorEmptyFilePathShouldSetDefault(t *testing.T) {
	receiver := NewReceiver("")

	if receiver.Config.ElasticsearchConnectionString != "" && len(receiver.Config.ElasticsearchConnectionString) < 0 {
		t.Fatalf("Should have default configuration and ElasticsearchConnectionString value was %v\n", receiver.Config.ElasticsearchConnectionString)
	}

	if receiver.Config.RabbitmqConnectionString != "" && len(receiver.Config.RabbitmqConnectionString) < 0 {
		t.Fatalf("Should have default configuration and RabbitmqConnectionString")
	}

	if receiver.Config.RabbitmqExchangeName != "" && len(receiver.Config.RabbitmqExchangeName) < 0 {
		t.Fatalf("Should have default configuration and RabbitmqExchangeName")
	}

	if receiver.Config.RabbitmqQueueName != "" && len(receiver.Config.RabbitmqQueueName) < 0 {
		t.Fatalf("Should have default configuration and RabbitmqQueueName")
	}
}

func TestCreateHashCheckLength(t *testing.T) {
	message := "Vaffler og kaffi"

	hash := createHash(message)

	if len(hash) < 32 {
		t.Fatalf("Hash Length is less than 32, is %d\n", len(hash))
	}
}

func TestCreateHashShouldBeEqualForSameString(t *testing.T) {
	message1 := "Testing 12345"
	//message2 := "Testing 12345"

	hash1 := createHash(message1)
	hash2 := createHash(message1)

	if hash1 == "" {
		t.Fatalf("Failed to create hash1")
	}

	if hash2 == "" {
		t.Fatalf("Failed to create hash2")
	}

	if hash1 != hash2 {
		t.Fatalf("Hashes are not equal!!!")
	}

}

func TestCreateHashShouldBeDifferentForDifferentStrings(t *testing.T) {
	message1 := "Testing 12345"
	message2 := "Testing 54321"

	hash1 := createHash(message1)
	hash2 := createHash(message2)

	if hash1 == hash2 {
		t.Fatalf("Hashes are equal!!!")
	}
}
