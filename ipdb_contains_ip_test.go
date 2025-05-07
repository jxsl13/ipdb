package ipdb

import (
	"bytes"
	"os"
	"testing"

	"github.com/jxsl13/ipdb/internal/testutils"
)

func TestContains(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	dbPath := testutils.FilePath("insert.db")
	defer func() {
		err := os.RemoveAll(dbPath)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	db, err := Open(ctx, dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}()

	r := MustParseIPv6Range("127.0.0.1-127.0.0.10")
	in := []byte("test insert")

	err = db.Insert(
		ctx,
		r,
		in,
	)
	if err != nil {
		t.Error(err)
		return
	}

	v, found, err := db.Contains(ctx, r.lb)
	if err != nil {
		t.Error(err)
		return
	}
	if !found {
		t.Errorf("expected to find value for %s, but got not found", r.lb)
		return
	}
	if !bytes.Equal(v, in) {
		t.Errorf("expected value %s, but got %s", in, v)
		return
	}

	v, found, err = db.Contains(ctx, r.ub)
	if err != nil {
		t.Error(err)
		return
	}
	if !found {
		t.Errorf("expected to find value for %s, but got not found", r.ub)
		return
	}

	if !bytes.Equal(v, in) {
		t.Errorf("expected value %s, but got %s", in, v)
		return
	}

	_, found, err = db.Contains(ctx, r.lb.Prev())
	if err != nil {
		t.Error(err)
		return
	}
	if found {
		t.Errorf("expected not to find value for %s, but got found", r.lb.Prev())
		return
	}

	_, found, err = db.Contains(ctx, r.ub.Next())
	if err != nil {
		t.Error(err)
		return
	}

	if found {
		t.Errorf("expected not to find value for %s, but got found", r.ub.Next())
		return
	}

	for ip := range r.IPs() {
		v, found, err := db.Contains(ctx, ip)
		if err != nil {
			t.Error(err)
			return
		}
		if !found {
			t.Errorf("expected to find value for %s, but got not found", ip)
			return
		}
		if !bytes.Equal(v, in) {
			t.Errorf("expected value %s, but got %s", in, v)
			return
		}
	}

}
