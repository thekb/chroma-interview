package chromainterview

import "testing"

func TestAtomicStoreBasic(t *testing.T) {
	as := newAtomicStore()
	as.Set("b", 100)
	if as.Get("b") != 100 {
		t.Fatal("expected 100")
	}
	as1 := as.Begin()
	as1.Set("a", 50)
	if as1.Get("a") != 50 {
		t.Fatal("expected 50")
	}
	if as1.Get("b") != 100 {
		t.Fatal("expected 100")
	}
	as2 := as1.Begin()
	if as2.Get("a") != 50 {
		t.Fatal("expected 50")
	}
	as2.Set("a", 60)
	if as2.Get("b") != 100 {
		t.Fatal("expected 100")
	}
	if as2.Get("a") != 60 {
		t.Fatal("expected 60")
	}
	as2.Rollback()
	as1.Commit()
	if as.Get("a") != 50 {
		t.Fatal("expected 50")
	}
}

func TestStoreSimple(t *testing.T) {
	s := newStore()
	s.Set("a", 50)
	if v := s.Get("a"); v != 50 {
		t.Fatal("expected 50, got: ", v)
	}
	s.Delete("a")
	if v := s.Get("a"); v != -1 {
		t.Fatal("expected -1, got: ", v)
	}
}

func TestAtomicStoreSimpleCommit(t *testing.T) {
	as := newAtomicStore()
	as1 := as.Begin()
	as1.Set("a", 50)
	if v := as1.Get("a"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as1.Commit()
	if v := as.Get("a"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
}

func TestAtomicSimpleRollback(t *testing.T) {
	as := newAtomicStore()
	as1 := as.Begin()
	as1.Set("a", 50)
	if v := as1.Get("a"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as1.Rollback()
	if v := as.Get("a"); v != -1 {
		t.Fatal("expected -1, got", v)
	}
}

func TestAtomicNestedCommit(t *testing.T) {
	as := newAtomicStore()
	as.Set("a", 100)
	if v := as.Get("a"); v != 100 {
		t.Fatal("expected 100, got v")
	}
	as1 := as.Begin()
	as1.Set("b", 50)
	if v := as1.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	if v := as1.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	as2 := as1.Begin()
	if v := as2.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as2.Set("c", 60)
	if v := as2.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	if v := as2.Get("c"); v != 60 {
		t.Fatal("expected 60, got", v)
	}
	as2.Commit()
	as1.Commit()
	if v := as.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	if v := as.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	if v := as.Get("c"); v != 60 {
		t.Fatal("expected 60, got", v)
	}
}

func TestAtomicNestedRollback(t *testing.T) {
	as := newAtomicStore()
	as.Set("a", 100)
	if v := as.Get("a"); v != 100 {
		t.Fatal("expected 100, got v")
	}
	as1 := as.Begin()
	as1.Set("b", 50)
	if v := as1.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	if v := as1.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	as2 := as1.Begin()
	if v := as2.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as2.Set("c", 60)
	if v := as2.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	if v := as2.Get("c"); v != 60 {
		t.Fatal("expected 60, got", v)
	}
	as2.Rollback()
	as1.Commit()
	if v := as.Get("a"); v != 100 {
		t.Fatal("expected 100, got", v)
	}
	if v := as.Get("b"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	if v := as.Get("c"); v != -1 {
		t.Fatal("expected -1, got", v)
	}
}

func TestAtomicStoreMultiOpOverrideCommit(t *testing.T) {
	as := newAtomicStore()
	as.Set("a", 40)
	as1 := as.Begin()
	as1.Set("a", 50)
	if v := as1.Get("a"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as1.Set("a", 60)
	if v := as1.Get("a"); v != 60 {
		t.Fatal("expected 60, got", v)
	}
	as1.Delete("a")
	if v := as1.Get("a"); v != -1 {
		t.Fatal("expected -1, got", v)
	}
	as1.Commit()
	if v := as.Get("a"); v != -1 {
		t.Fatal("expected -1, got", v)
	}
}

func TestAtomicStoreMultiOpOverrideRollback(t *testing.T) {
	as := newAtomicStore()
	as.Set("a", 40)
	as1 := as.Begin()
	as1.Set("a", 50)
	if v := as1.Get("a"); v != 50 {
		t.Fatal("expected 50, got", v)
	}
	as1.Set("a", 60)
	if v := as1.Get("a"); v != 60 {
		t.Fatal("expected 60, got", v)
	}
	as1.Delete("a")
	if v := as1.Get("a"); v != -1 {
		t.Fatal("expected -1, got", v)
	}
	as1.Rollback()
	if v := as.Get("a"); v != 40 {
		t.Fatal("expected 40, got", v)
	}
}
