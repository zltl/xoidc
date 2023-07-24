package password

import "testing"

func TestGenHash(t *testing.T) {
	hash, err := CreateHash("123456")
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
}

func BenchmarkGenHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CreateHash("123456")
	}
}

func BenchmarkCheckHash(b *testing.B) {
	hash, _ := CreateHash("123456")
	for i := 0; i < b.N; i++ {
		_, _ = ComparePasswordAndHash("123456", hash)
	}
}
