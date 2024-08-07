package dvs

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
)

// New creates a new DVS.
func New[V any](reg com.Registry) DVS[V] {
	return DVS[V]{reg}
}

// DVS provides versioning support for the mus-stream-go serializer.
type DVS[V any] struct {
	reg com.Registry
}

// Marshal migrates v to the version specified by dtm and then writes dtm +
// resulting v version encoding.
//
// Returns the number of written bytes and one of the ErrUnknownDTM,
// ErrWrongTypeVersion or Writer errors.
func (dvs DVS[V]) Marshal(dtm com.DTM, v V, w muss.Writer) (n int,
	err error) {
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	return mver.MigrateCurrentAndMarshal(v, w)
}

// Unmarshal unmarshals dtm + data, and then migrates data to the version
// specified by dtm.
//
// Returns the number of read bytes and one of the ErrUnknownDTM,
// ErrWrongTypeVersion or Reader errors.
func (dvs DVS[V]) Unmarshal(r muss.Reader) (dtm com.DTM, v V, n int,
	err error) {
	dtm, n, err = dts.UnmarshalDTM(r)
	if err != nil {
		return
	}
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	var n1 int
	v, n1, err = mver.UnmarshalAndMigrateOld(r)
	n += n1
	return
}

func (dvs DVS[V]) getMV(dtm com.DTM) (mver MigrationVersion[V], err error) {
	tver, err := dvs.reg.Get(dtm)
	if err != nil {
		return
	}
	mver, ok := tver.(MigrationVersion[V])
	if !ok {
		err = com.ErrWrongTypeVersion
		return
	}
	return
}
