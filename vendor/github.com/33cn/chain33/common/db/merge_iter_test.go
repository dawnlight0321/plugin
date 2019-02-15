package db

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGoMemDB(t *testing.T) DB {
	dir, err := ioutil.TempDir("", "gomemdb")
	require.NoError(t, err)
	memdb, err := NewGoMemDB("gomemdb", dir, 128)
	require.NoError(t, err)
	return memdb
}

func TestMergeIter(t *testing.T) {
	db1 := newGoMemDB(t)
	db2 := newGoMemDB(t)
	db1.Set([]byte("1"), []byte("1"))
	db2.Set([]byte("2"), []byte("2"))

	//合并以后:
	db := NewMergedIteratorDB([]IteratorDB{db1, db2})
	it0 := NewListHelper(db)
	list0 := it0.List(nil, nil, 100, 0)
	assert.Equal(t, 2, len(list0))
	assert.Equal(t, "2", string(list0[0]))
	assert.Equal(t, "1", string(list0[1]))

	list0 = it0.List(nil, nil, 100, 1)
	assert.Equal(t, 2, len(list0))
	assert.Equal(t, "1", string(list0[0]))
	assert.Equal(t, "2", string(list0[1]))
}

func TestMergeIterDup1(t *testing.T) {
	db1 := newGoMemDB(t)
	db2 := newGoMemDB(t)
	db1.Set([]byte("1"), []byte("1"))
	db2.Set([]byte("2"), []byte("2"))

	//合并以后:
	db := NewMergedIteratorDB([]IteratorDB{db1, db2})
	it0 := NewListHelper(db)
	//测试修改
	db1.Set([]byte("2"), []byte("12"))
	list0 := it0.List(nil, nil, 100, 0)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 2, len(list0))
	assert.Equal(t, "12", string(list0[0]))
	assert.Equal(t, "1", string(list0[1]))
}

func TestMergeIterDup2(t *testing.T) {
	db1 := newGoMemDB(t)
	db2 := newGoMemDB(t)
	db1.Set([]byte("key-1"), []byte("db1-key-1"))
	db1.Set([]byte("key-3"), []byte("db1-key-3"))
	db2.Set([]byte("key-2"), []byte("db2-key-2"))

	//合并以后:
	db := NewMergedIteratorDB([]IteratorDB{db1, db2})
	it0 := NewListHelper(db)
	//测试修改
	db2.Set([]byte("key-3"), []byte("db2-key-3"))
	list0 := it0.List(nil, nil, 100, 0)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 3, len(list0))
	assert.Equal(t, "db1-key-3", string(list0[0]))
	assert.Equal(t, "db2-key-2", string(list0[1]))
	assert.Equal(t, "db1-key-1", string(list0[2]))

	list0 = it0.List(nil, nil, 100, 1)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 3, len(list0))
	assert.Equal(t, "db1-key-1", string(list0[0]))
	assert.Equal(t, "db2-key-2", string(list0[1]))
	assert.Equal(t, "db1-key-3", string(list0[2]))
}

func TestMergeIterDup3(t *testing.T) {
	db1 := newGoMemDB(t)
	db2 := newGoMemDB(t)
	db1.Set([]byte("key-1"), []byte("db1-key-1"))
	db1.Set([]byte("key-3"), []byte("db1-key-3"))
	db2.Set([]byte("key-2"), []byte("db2-key-2"))

	//合并以后:
	db := NewMergedIteratorDB([]IteratorDB{db1, db2})
	it0 := NewListHelper(db)
	//测试修改
	db1.Set([]byte("key-2"), []byte("db1-key-2"))
	list0 := it0.List(nil, nil, 100, 0)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 3, len(list0))
	assert.Equal(t, "db1-key-3", string(list0[0]))
	assert.Equal(t, "db1-key-2", string(list0[1]))
	assert.Equal(t, "db1-key-1", string(list0[2]))

	list0 = it0.List(nil, nil, 100, 1)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 3, len(list0))
	assert.Equal(t, "db1-key-1", string(list0[0]))
	assert.Equal(t, "db1-key-2", string(list0[1]))
	assert.Equal(t, "db1-key-3", string(list0[2]))
}

func TestMergeIterSearch(t *testing.T) {
	db1 := newGoMemDB(t)
	db2 := newGoMemDB(t)
	db1.Set([]byte("key-1"), []byte("db1-key-1"))
	db1.Set([]byte("key-2"), []byte("db1-key-2"))
	db2.Set([]byte("key-2"), []byte("db2-key-2"))
	db2.Set([]byte("key-3"), []byte("db2-key-3"))
	db2.Set([]byte("key-4"), []byte("db2-key-4"))

	//合并以后:
	db := NewMergedIteratorDB([]IteratorDB{db1, db2})
	it0 := NewListHelper(db)
	list0 := it0.List([]byte("key-"), []byte("key-2"), 100, 0)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 1, len(list0))
	assert.Equal(t, "db1-key-1", string(list0[0]))

	list0 = it0.List([]byte("key-"), []byte("key-2"), 100, 1)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 2, len(list0))
	assert.Equal(t, "db2-key-3", string(list0[0]))
	assert.Equal(t, "db2-key-4", string(list0[1]))
}

func TestIterSearch(t *testing.T) {
	db1 := newGoMemDB(t)
	defer db1.Close()
	db1.Set([]byte("key-1"), []byte("db1-key-1"))
	db1.Set([]byte("key-2"), []byte("db2-key-2"))
	db1.Set([]byte("key-2"), []byte("db1-key-2"))
	db1.Set([]byte("key-3"), []byte("db2-key-3"))
	db1.Set([]byte("key-4"), []byte("db2-key-4"))

	it0 := NewListHelper(db1)
	list0 := it0.List([]byte("key-"), []byte("key-2"), 100, 0)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 1, len(list0))
	assert.Equal(t, "db1-key-1", string(list0[0]))

	list0 = it0.List([]byte("key-"), []byte("key-2"), 100, 1)
	for k, v := range list0 {
		println(k, string(v))
	}
	assert.Equal(t, 2, len(list0))
	assert.Equal(t, "db2-key-3", string(list0[0]))
	assert.Equal(t, "db2-key-4", string(list0[1]))
}