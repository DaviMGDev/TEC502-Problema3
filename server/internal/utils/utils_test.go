package utils

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name		string
		input		string
		expCommand	string
		expArgs		[]string
	}{
		{"Empty string", "", "", []string{}},
		{"Chat message", "hello world", "chat", []string{"hello world"}},
		{"Simple command", "/register user1 pass1", "register", []string{"user1", "pass1"}},
		{"Command with multiple args", "/move card1 pos2", "move", []string{"card1", "pos2"}},
		{"Command with no args", "/status", "status", []string{}},
		{"Command with mixed case", "/LOGIN user pass", "login", []string{"user", "pass"}},
		{"Command with leading/trailing spaces", " /command arg ", "chat", []string{" /command arg "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("Running TestParseCommand: %s", tt.name)
			cmd, args := ParseCommand(tt.input)
			if cmd != tt.expCommand {
				t.Errorf("Expected command %q, got %q", tt.expCommand, cmd)
			}
			if !reflect.DeepEqual(args, tt.expArgs) {
				t.Errorf("Expected args %v, got %v", tt.expArgs, args)
			}
			log.Printf("TestParseCommand '%s' passed.", tt.name)
		})
	}
}

func TestSafeMap(t *testing.T) {
	m := NewSafeMap[string, int]()

	t.Run("Set and Get", func(t *testing.T) {
		log.Printf("Running TestSafeMap: Set and Get")
		m.Set("key1", 100)
		val, ok := m.Get("key1")
		if !ok || val != 100 {
			t.Errorf("Expected key1 to be 100, got %d, %v", val, ok)
		}
		log.Printf("TestSafeMap 'Set and Get' passed.")
	})

	t.Run("Has", func(t *testing.T) {
		log.Printf("Running TestSafeMap: Has")
		if !m.Has("key1") {
			t.Errorf("Expected key1 to exist")
		}
		if m.Has("nonexistent") {
			t.Errorf("Expected nonexistent not to exist")
		}
		log.Printf("TestSafeMap 'Has' passed.")
	})

	t.Run("Delete", func(t *testing.T) {
		log.Printf("Running TestSafeMap: Delete")
		m.Delete("key1")
		_, ok := m.Get("key1")
		if ok {
			t.Errorf("Expected key1 to be deleted")
		}
		log.Printf("TestSafeMap 'Delete' passed.")
	})

	t.Run("Concurrency", func(t *testing.T) {
		log.Printf("Running TestSafeMap: Concurrency")
		m_concurrency := NewSafeMap[int, int]()
		var wg sync.WaitGroup
		numWorkers := 100
		numOperations := 1000

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := workerID*numOperations + j
					m_concurrency.Set(key, key*2)
					_, _ = m_concurrency.Get(key)
					if j%10 == 0 {
						m_concurrency.Delete(key)
					}
				}
			}(i)
		}
		wg.Wait()

		log.Printf("Final SafeMap size: %d", len(m_concurrency.Keys()))
		log.Printf("TestSafeMap 'Concurrency' passed.")
	})

	t.Run("NewSafeMapFromMap", func(t *testing.T) {
		log.Printf("Running TestSafeMap: NewSafeMapFromMap")
		initial := map[string]int{"a": 1, "b": 2}
		sm := NewSafeMapFromMap(initial)
		val, ok := sm.Get("a")
		if !ok || val != 1 {
			t.Errorf("Expected 'a' to be 1, got %d, %v", val, ok)
		}
		val, ok = sm.Get("b")
		if !ok || val != 2 {
			t.Errorf("Expected 'b' to be 2, got %d, %v", val, ok)
		}
		log.Printf("TestSafeMap 'NewSafeMapFromMap' passed.")
	})

	t.Run("Keys and Values", func(t *testing.T) {
		log.Printf("Running TestSafeMap: Keys and Values")
		m = NewSafeMap[string, int]()
		m.Set("k1", 1)
		m.Set("k2", 2)
		m.Set("k3", 3)

		keys := m.Keys()
		if len(keys) != 3 {
			t.Errorf("Expected 3 keys, got %d", len(keys))
		}
		values := m.Values()
		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}
		log.Printf("TestSafeMap 'Keys and Values' passed.")
	})
}

func TestSafeList(t *testing.T) {
	l := NewSafeList[string]()

	t.Run("Append", func(t *testing.T) {
		log.Printf("Running TestSafeList: Append")
		l.Append("item1")
		l.Append("item2")
		if l.Size() != 2 {
			t.Errorf("Expected size 2, got %d", l.Size())
		}
		val, ok := l.Get(0)
		if !ok || val != "item1" {
			t.Errorf("Expected item at index 0 to be 'item1', got %s, %v", val, ok)
		}
		log.Printf("TestSafeList 'Append' passed.")
	})

	t.Run("Contains", func(t *testing.T) {
		log.Printf("Running TestSafeList: Contains")
		if !l.Contains("item1") {
			t.Errorf("Expected 'item1' to exist")
		}
		if l.Contains("nonexistent") {
			t.Errorf("Expected 'nonexistent' not to exist")
		}
		log.Printf("TestSafeList 'Contains' passed.")
	})

	t.Run("Remove", func(t *testing.T) {
		log.Printf("Running TestSafeList: Remove")
		if !l.Remove("item1") {
			t.Errorf("Expected 'item1' to be removed")
		}
		if l.Size() != 1 {
			t.Errorf("Expected size 1 after removal, got %d", l.Size())
		}
		if l.Contains("item1") {
			t.Errorf("Expected 'item1' not to exist")
		}
		log.Printf("TestSafeList 'Remove' passed.")
	})

	t.Run("Pop", func(t *testing.T) {
		log.Printf("Running TestSafeList: Pop")
		l.Append("item3")
		val, ok := l.Pop()
		if !ok || val != "item3" {
			t.Errorf("Expected 'item3' to be popped, got %s, %v", val, ok)
		}
		if l.Size() != 1 {
			t.Errorf("Expected size 1 after pop, got %d", l.Size())
		}
		log.Printf("TestSafeList 'Pop' passed.")
	})

	t.Run("Get out of bounds", func(t *testing.T) {
		log.Printf("Running TestSafeList: Get out of bounds")
		_, ok := l.Get(10)
		if ok {
			t.Errorf("Expected Get out of bounds to return false")
		}
		log.Printf("TestSafeList 'Get out of bounds' passed.")
	})

	t.Run("Search", func(t *testing.T) {
		log.Printf("Running TestSafeList: Search")
		l = NewSafeList[string]()
		l.Append("apple")
		l.Append("banana")
		l.Append("cherry")
		if l.Search("banana") != 1 {
			t.Errorf("Expected 'banana' at index 1, got %d", l.Search("banana"))
		}
		if l.Search("grape") != -1 {
			t.Errorf("Expected 'grape' not to be found, got %d", l.Search("grape"))
		}
		log.Printf("TestSafeList 'Search' passed.")
	})

	t.Run("Concurrency", func(t *testing.T) {
		log.Printf("Running TestSafeList: Concurrency")
		l_concurrency := NewSafeList[int]()
		var wg sync.WaitGroup
		numWorkers := 100
		numOperations := 1000

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					item := workerID*numOperations + j
					l_concurrency.Append(item)
					_, _ = l_concurrency.Get(0)
					if j%10 == 0 {
						l_concurrency.Remove(item)
					}
				}
			}(i)
		}
		wg.Wait()
		log.Printf("Final SafeList size: %d", l_concurrency.Size())

		if l_concurrency.Size() > numWorkers*numOperations {
			t.Errorf("List size %d is unexpectedly large", l_concurrency.Size())
		}
		log.Printf("TestSafeList 'Concurrency' passed.")
	})
}

func TestDict(t *testing.T) {

	t.Run("Json conversion", func(t *testing.T) {
		log.Printf("Running TestDict: Json conversion")
		d := make(Dict)
		d["name"] = "test"
		d["value"] = 123
		jsonBytes, err := d.Json()
		if err != nil {
			t.Errorf("Error marshalling to JSON: %v", err)
		}
		expected := `{"name":"test","value":123}`

		var result map[string]interface{}
		json.Unmarshal(jsonBytes, &result)
		if result["name"] != "test" || result["value"].(float64) != 123 {
			t.Errorf("Expected JSON %s, got %s", expected, string(jsonBytes))
		}
		log.Printf("TestDict 'Json conversion' passed.")
	})

	t.Run("String conversion", func(t *testing.T) {
		log.Printf("Running TestDict: String conversion")
		d := make(Dict)
		d["key"] = "val"
		str, err := d.String()
		if err != nil {
			t.Errorf("Error converting to string: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal([]byte(str), &result)
		if result["key"] != "val" {
			t.Errorf("Expected string to contain '\"key\":\"val\"', got %s", str)
		}
		log.Printf("TestDict 'String conversion' passed.")
	})
}

func TestMux(t *testing.T) {

	type testFunc func(string) string

	safeMap := NewSafeMap[string, testFunc]()

	defaultHandler := func(s string) string { return "default:" + s }
	mux := &Mux[testFunc]{
		Map:		safeMap,
		defaultFn:	defaultHandler,
	}

	handler1 := func(s string) string { return "handler1:" + s }
	handler2 := func(s string) string { return "handler2:" + s }

	safeMap.Set("cmd1", handler1)
	safeMap.Set("cmd2", handler2)

	t.Run("Get existing handler", func(t *testing.T) {
		log.Printf("Running TestMux: Get existing handler")
		fn, ok := mux.Get("cmd1")
		if !ok || fn("test") != "handler1:test" {
			t.Errorf("Expected handler1 for cmd1, got %v, %v", fn("test"), ok)
		}
		log.Printf("TestMux 'Get existing handler' passed.")
	})

	t.Run("Get non-existing handler", func(t *testing.T) {
		log.Printf("Running TestMux: Get non-existing handler")
		fn, ok := mux.Get("cmd3")
		if ok || fn("test") != "default:test" {
			t.Errorf("Expected default handler for cmd3, got %v, %v", fn("test"), ok)
		}
		log.Printf("TestMux 'Get non-existing handler' passed.")
	})

	t.Run("Verify default function behavior when no specific handler", func(t *testing.T) {
		log.Printf("Running TestMux: Verify default function behavior")
		safeMap.Delete("cmd1")
		fn, ok := mux.Get("cmd1")
		if ok || fn("another_test") != "default:another_test" {
			t.Errorf("Expected default handler after deletion, got %v, %v", fn("another_test"), ok)
		}
		log.Printf("TestMux 'Verify default function behavior' passed.")
	})
}
