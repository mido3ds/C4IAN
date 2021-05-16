package msec

import (
	"bytes"
	"testing"
)

func TestMSecLayer(t *testing.T) {
	msec := NewMSecLayer("hello world password")

	strs := []string{
		"hey",
		"how are ya",
		"fine",
		"ajsoidfh;qwoiehfwoqiehfwquofhwupf892fhqowe8uhfqw9ehfwudifghqw9fuewqoifghweiuhfhf",
		"231481923u9182828228282822828",
		"",
		"a",
	}

	for i := 0; i < len(strs); i++ {
		s := msec.Decrypt(msec.Encrypt([]byte(strs[i])))

		if !bytes.Equal([]byte(strs[i]), s) {
			t.Errorf("coudln't encrypt #%d, input %v output %v", i, []byte(strs[i]), s)
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {
	msec := NewMSecLayer("hello world password")
	in := []byte("ajsoidfh;qwoiehfwoqiehfwquofhwupf892fhqowe8uhfqw9ehfwudifghqw9fuewqoifghweiuhfhf")

	for i := 0; i < b.N; i++ {
		msec.Encrypt(in)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	msec := NewMSecLayer("hello world password")
	in := []byte("ajsoidfh;qwoiehfwoqiehfwquofhwupf892fhqowe8uhfqw9ehfwudifghqw9fuewqoifghweiuhfhf")

	for i := 0; i < b.N; i++ {
		msec.Decrypt(in)
	}
}
