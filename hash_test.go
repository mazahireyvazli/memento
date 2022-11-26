package memento

import "testing"

var hashTable = map[string]uint64{
	"eec7e390-059c-48b5-881b-d2f26bde592d": 1801486099476478190,
	"6b9c5b0d-ad44-4bec-9fe3-4df038aa0cab": 2104812043201492727,
	"a8bb7304-267a-4b28-8a92-2f5d41c2e0f5": 15016133345029420581,
	"f86de91c-9fda-4d00-8a98-bdd98458f764": 11523290569367520202,
	"bf73d161-e053-45cb-aff4-dd48ef5548a5": 14054307095194864094,
	"creamwove":                            13683957712591044104,
	"quists":                               18132031479020902632,
}

func TestFnv1_64(t *testing.T) {
	hasher := Fnv1_64[string]{}

	for k, v := range hashTable {
		hashedKey := hasher.Sum64(k)

		if hashedKey != v {
			t.Fatalf("hashed key %d does not match expected hash value %d", hashedKey, v)
		}
	}
}
