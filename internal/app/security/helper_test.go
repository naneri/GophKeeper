package security

import "testing"

const key = "somerandomkey"

func Test_Hashing(t *testing.T) {
	password := "somerandompass"
	wrongPass := password + "fake"
	hashedPass := Hash(key, password)

	t.Run("check hashing", func(t *testing.T) {
		result, err := CheckHash(key, password, hashedPass)

		if err != nil {
			t.Error("error hashing password" + err.Error())
		}

		if !result {
			t.Error("error hashing password")
		}
	})

	t.Run("check wrong hashing", func(t *testing.T) {
		result, err := CheckHash(key, wrongPass, hashedPass)
		if err != nil {
			t.Error("error hashing password" + err.Error())
		}

		if result {
			t.Errorf("wrong password - %s is considered correct. Original password: %s", wrongPass, password)
		}
	})
}
