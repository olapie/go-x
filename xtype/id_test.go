package xtype

import (
	"math"
	"testing"
)

func TestID(t *testing.T) {
	//for i := 0; i < 256; i++ {
	//	id := NextID()
	//	t.Logf("%d %0X %s", id, id, id.Base62())
	//	//time.Sleep(time.Millisecond * 1)
	//}

	var id ID = 123
	if id.Base62() != "1Z" {
		t.Log(id.Base62())
		t.FailNow()
	}

	id = 62
	if id.Base62() != "10" {
		t.Log(id.Base62())
		t.FailNow()
	}

	id = math.MaxInt64

	if i, _ := IDFromBase62(id.Base62()); i != id {
		t.Log(id.Base62(), i)
		t.FailNow()
	}
}

func TestID_Base36(t *testing.T) {
	//for i := 0; i < 256; i++ {
	//	id := NextID()
	//	t.Logf("%d %0X %s", id, id, id.Base36())
	//	//time.Sleep(time.Millisecond * 1)
	//}

	var id ID = 123
	if id.Base36() != "3f" {
		t.Log(id.Base36())
		t.FailNow()
	}

	id = 34
	if id.Base36() != "y" {
		t.Log(id.Base36())
		t.FailNow()
	}

	id = math.MaxInt64
	if i, _ := IDFromBase36(id.Base36()); i != id {
		t.Log(id.Base36(), i)
		t.FailNow()
	}
}
