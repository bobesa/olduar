package olduar

import "testing"

func TestGuids(t *testing.T) {
	tests := make(map[GUID]GUID)
	for i := 0; i < 500; i++ {
		guid := GenerateGUID()
		if _, exists := tests[guid]; exists {
			t.Error("Duplicate GUID found")
			return
		} else {
			tests[guid] = guid
		}
	}
}
