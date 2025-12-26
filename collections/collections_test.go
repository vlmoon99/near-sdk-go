package collections

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func cleanupStorage(t *testing.T) {
	t.Helper()
	mockSys, ok := env.NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatal("Environment is not set to MockSystem")
	}
	mockSys.Storage = make(map[string][]byte)
}

// ============================================================================
// Vector Tests
// ============================================================================

func TestVector_Push_Get_Set(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[string]("v")

	if err := v.Push("A"); err != nil {
		t.Fatal(err)
	}
	if err := v.Push("B"); err != nil {
		t.Fatal(err)
	}

	if v.Length() != 2 {
		t.Errorf("Expected length 2, got %d", v.Length())
	}

	val, err := v.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if val != "B" {
		t.Errorf("Expected 'B', got '%s'", val)
	}

	if err := v.Set(1, "C"); err != nil {
		t.Fatal(err)
	}
	val, _ = v.Get(1)
	if val != "C" {
		t.Errorf("Expected 'C' after Set, got '%s'", val)
	}

	if _, err := v.Get(99); err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}
}

func TestVector_Pop(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[int]("v")

	v.Push(10)
	v.Push(20)

	val, err := v.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if val != 20 {
		t.Errorf("Expected 20, got %d", val)
	}
	if v.Length() != 1 {
		t.Errorf("Expected length 1, got %d", v.Length())
	}

	val, err = v.Pop()
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	_, err = v.Pop()
	if err != ErrVectorEmpty {
		t.Errorf("Expected ErrVectorEmpty, got %v", err)
	}
}

func TestVector_ToSlice(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[string]("v")

	v.Push("X")
	v.Push("Y")
	v.Push("Z")

	slice, err := v.ToSlice()
	if err != nil {
		t.Fatal(err)
	}
	if len(slice) != 3 {
		t.Errorf("Expected slice len 3, got %d", len(slice))
	}
	if slice[0] != "X" || slice[2] != "Z" {
		t.Errorf("Slice content mismatch: %v", slice)
	}
}

// ============================================================================
// LookupMap Tests
// ============================================================================

func TestLookupMap_Basic(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, int]("m")

	if err := m.Insert("alice", 100); err != nil {
		t.Fatal(err)
	}

	exists, _ := m.Contains("alice")
	if !exists {
		t.Error("Expected alice to exist")
	}

	val, err := m.Get("alice")
	if err != nil {
		t.Fatal(err)
	}
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}

	m.Remove("alice")
	exists, _ = m.Contains("alice")
	if exists {
		t.Error("Expected alice to be removed")
	}

	_, err = m.Get("bob")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

type TestStruct struct {
	Name string
	Age  int
}

func TestLookupMap_Structs(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[uint64, TestStruct]("m")

	obj := TestStruct{Name: "Test", Age: 42}
	m.Insert(1, obj)

	res, err := m.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if res.Name != "Test" || res.Age != 42 {
		t.Errorf("Struct mismatch: %+v", res)
	}
}

// ============================================================================
// LookupSet Tests
// ============================================================================

func TestLookupSet(t *testing.T) {
	defer cleanupStorage(t)
	s := NewLookupSet[string]("s")

	s.Insert("A")
	s.Insert("B")

	if exists, _ := s.Contains("A"); !exists {
		t.Error("Set should contain A")
	}
	if exists, _ := s.Contains("C"); exists {
		t.Error("Set should not contain C")
	}

	s.Remove("A")
	if exists, _ := s.Contains("A"); exists {
		t.Error("A should be removed")
	}
}

// ============================================================================
// UnorderedMap Tests
// ============================================================================

func TestUnorderedMap_Workflow(t *testing.T) {
	defer cleanupStorage(t)
	m := NewUnorderedMap[string, int]("um")

	m.Insert("one", 1)
	m.Insert("two", 2)
	m.Insert("three", 3)

	if m.Length() != 3 {
		t.Errorf("Expected length 3, got %d", m.Length())
	}

	m.Insert("two", 22)
	if m.Length() != 3 {
		t.Errorf("Length shouldn't change on overwrite")
	}
	val, _ := m.Get("two")
	if val != 22 {
		t.Errorf("Expected updated value 22, got %d", val)
	}

	err := m.Remove("two")
	if err != nil {
		t.Fatal(err)
	}

	if m.Length() != 2 {
		t.Errorf("Expected length 2, got %d", m.Length())
	}
	if exists, _ := env.StorageHasKey([]byte("um:v:two")); exists {
		t.Error("Storage should verify key removal")
	}

	keys, _ := m.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys")
	}
	foundOne := false
	foundThree := false
	for _, k := range keys {
		if k == "one" {
			foundOne = true
		}
		if k == "three" {
			foundThree = true
		}
	}
	if !foundOne || !foundThree {
		t.Errorf("Missing keys after removal. Got: %v", keys)
	}
}

func TestUnorderedMap_Clear(t *testing.T) {
	defer cleanupStorage(t)
	m := NewUnorderedMap[int, int]("um")
	m.Insert(1, 1)
	m.Insert(2, 2)

	m.Clear()

	if m.Length() != 0 {
		t.Error("Length should be 0")
	}
	keys, _ := m.Keys()
	if len(keys) != 0 {
		t.Error("Keys slice should be empty")
	}
}

// ============================================================================
// UnorderedSet Tests
// ============================================================================

func TestUnorderedSet(t *testing.T) {
	defer cleanupStorage(t)
	s := NewUnorderedSet[string]("us")

	s.Insert("apple")
	s.Insert("banana")
	s.Insert("cherry")

	s.Insert("apple")
	if s.Length() != 3 {
		t.Errorf("Duplicate insert shouldn't increase length, got %d", s.Length())
	}

	if yes, _ := s.Contains("banana"); !yes {
		t.Error("Should contain banana")
	}

	s.Remove("banana")
	if s.Length() != 2 {
		t.Error("Length should be 2")
	}

	items, _ := s.All()
	if len(items) != 2 {
		t.Error("Should return 2 items")
	}
	for _, item := range items {
		if item == "banana" {
			t.Error("Banana should be gone")
		}
	}
}

// ============================================================================
// TreeMap Tests
// ============================================================================

func TestTreeMap_Sorting(t *testing.T) {
	defer cleanupStorage(t)
	tm := NewTreeMap[int, string]("tm")

	tm.Insert(20, "twenty")
	tm.Insert(10, "ten")
	tm.Insert(30, "thirty")

	keys, err := tm.Keys()
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 3 {
		t.Fatal("Expected 3 keys")
	}
	if keys[0] != 10 || keys[1] != 20 || keys[2] != 30 {
		t.Errorf("Keys not sorted: %v", keys)
	}
}

func TestTreeMap_MinMax(t *testing.T) {
	defer cleanupStorage(t)
	tm := NewTreeMap[string, int]("tm")

	if _, err := tm.MinKey(); err != ErrMapEmpty {
		t.Error("Expected ErrMapEmpty")
	}

	tm.Insert("b", 2)
	tm.Insert("a", 1)
	tm.Insert("c", 3)

	min, _ := tm.MinKey()
	max, _ := tm.MaxKey()

	if min != "a" {
		t.Errorf("Expected min 'a', got %s", min)
	}
	if max != "c" {
		t.Errorf("Expected max 'c', got %s", max)
	}
}

func TestTreeMap_Remove_Rebalance(t *testing.T) {
	defer cleanupStorage(t)
	tm := NewTreeMap[int, int]("tm")

	tm.Insert(1, 1)
	tm.Insert(2, 2)
	tm.Insert(3, 3)

	tm.Remove(2)

	keys, _ := tm.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected len 2, got %v", keys)
	}
	if keys[0] != 1 || keys[1] != 3 {
		t.Errorf("Expected [1, 3], got %v", keys)
	}

	_, err := tm.Get(2)
	if err != ErrKeyNotFound {
		t.Error("Value for 2 should be gone")
	}
}
