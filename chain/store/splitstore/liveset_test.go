package splitstore

import (
	"os"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func TestLMDBLiveSet(t *testing.T) {
	testLiveSet(t, true)
}

func TestBoltLiveSet(t *testing.T) {
	testLiveSet(t, false)
}

func testLiveSet(t *testing.T, useLMDB bool) {
	t.Helper()

	path := "/tmp/liveset-test"

	err := os.MkdirAll(path, 0777)
	if err != nil {
		t.Fatal(err)
	}

	env, err := NewLiveSetEnv(path, useLMDB)
	if err != nil {
		t.Fatal(err)
	}
	defer env.Close() //nolint:errcheck

	hotSet, err := env.NewLiveSet("hot")
	if err != nil {
		t.Fatal(err)
	}

	coldSet, err := env.NewLiveSet("cold")
	if err != nil {
		t.Fatal(err)
	}

	makeCid := func(key string) cid.Cid {
		h, err := multihash.Sum([]byte(key), multihash.SHA2_256, -1)
		if err != nil {
			t.Fatal(err)
		}

		return cid.NewCidV1(cid.Raw, h)
	}

	mustHave := func(s LiveSet, cid cid.Cid) {
		has, err := s.Has(cid)
		if err != nil {
			t.Fatal(err)
		}

		if !has {
			t.Fatal("mark not found")
		}
	}

	mustNotHave := func(s LiveSet, cid cid.Cid) {
		has, err := s.Has(cid)
		if err != nil {
			t.Fatal(err)
		}

		if has {
			t.Fatal("unexpected mark")
		}
	}

	k1 := makeCid("a")
	k2 := makeCid("b")
	k3 := makeCid("c")
	k4 := makeCid("d")

	hotSet.Mark(k1)  //nolint
	hotSet.Mark(k2)  //nolint
	coldSet.Mark(k3) //nolint

	mustHave(hotSet, k1)
	mustHave(hotSet, k2)
	mustNotHave(hotSet, k3)
	mustNotHave(hotSet, k4)

	mustNotHave(coldSet, k1)
	mustNotHave(coldSet, k2)
	mustHave(coldSet, k3)
	mustNotHave(coldSet, k4)

	// close them and reopen to redo the dance

	err = hotSet.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = coldSet.Close()
	if err != nil {
		t.Fatal(err)
	}

	hotSet, err = env.NewLiveSet("hot")
	if err != nil {
		t.Fatal(err)
	}

	coldSet, err = env.NewLiveSet("cold")
	if err != nil {
		t.Fatal(err)
	}

	hotSet.Mark(k3)  //nolint
	hotSet.Mark(k4)  //nolint
	coldSet.Mark(k1) //nolint

	mustNotHave(hotSet, k1)
	mustNotHave(hotSet, k2)
	mustHave(hotSet, k3)
	mustHave(hotSet, k4)

	mustHave(coldSet, k1)
	mustNotHave(coldSet, k2)
	mustNotHave(coldSet, k3)
	mustNotHave(coldSet, k4)

	err = hotSet.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = coldSet.Close()
	if err != nil {
		t.Fatal(err)
	}
}
