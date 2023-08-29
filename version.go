package dvs

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
)

// MigrationVersion represents a generic type version for Registry that can
// be migrated.
//
// It contains methods to support all mus-stream-dvs-go functionality.
type MigrationVersion[V any] interface {
	MigrateCurrentAndMarshalMUS(v V, w muss.Writer) (n int, err error)
	UnmarshalAndMigrateOldMUS(r muss.Reader) (v V, n int, err error)
}

// Version is an implementation of the MigrationVersion interface.
type Version[T any, V any] struct {
	DTS            dts.DTS[T]
	MigrateOld     com.MigrateOld[T, V]
	MigrateCurrent com.MigrateCurrent[V, T]
}

func (ver Version[T, V]) MigrateCurrentAndMarshalMUS(v V, w muss.Writer) (n int,
	err error) {
	t, err := ver.MigrateCurrent(v)
	if err != nil {
		return
	}
	return ver.DTS.MarshalMUS(t, w)
}

func (ver Version[T, V]) UnmarshalAndMigrateOldMUS(r muss.Reader) (v V, n int,
	err error) {
	t, n, err := ver.DTS.UnmarshalDataMUS(r)
	if err != nil {
		return
	}
	v, err = ver.MigrateOld(t)
	return
}
