package ipdb

import (
	"bytes"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/jxsl13/ipdb/internal/testutils"
)

func TestInsertRange(t *testing.T) {
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

	r := MustParseIPv6Range("127.0.0.1/16")
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

	err = db.ValidateConsistency(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	wg := &sync.WaitGroup{}
	routines := make([]chan IPv6, 0, runtime.NumCPU())
	wg.Add(runtime.NumCPU())
	for i := range runtime.NumCPU() {
		c := make(chan IPv6)
		routines = append(routines, c)

		t.Logf("starting routine %d", i)
		go func(c chan IPv6) {
			defer wg.Done()

			for ip := range c {
				v, found, err := db.Contains(ctx, ip)
				if err != nil {
					t.Error(err)
					return
				}

				if !found {
					t.Errorf("ip %s not found, even tho range %s was inserted into the database", ip, r)
					return
				}

				if !bytes.Equal(v, in) {
					t.Errorf("ip %s has different value than expected, got %s, expected %s", ip, string(v), string(in))
					return
				}
			}
		}(c)
	}

	cnt := 0
	for ip := range r.IPs() {
		c := routines[cnt%len(routines)]
		select {
		case c <- ip:
		case <-ctx.Done():
			t.Error(ctx.Err())
			return
		}
		cnt++
	}
	t.Logf("checked %d IPs in range %s", cnt, r)

	for _, c := range routines {
		close(c)
	}

	wg.Wait()
	t.Log("all routines finished")
}
