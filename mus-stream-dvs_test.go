package dvs

import (
	"errors"
	"reflect"
	"testing"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
	muss_mock "github.com/mus-format/mus-stream-go/testdata/mock"
	"github.com/mus-format/mus-stream-go/varint"
	"github.com/ymz-ncnk/mok"
)

const (
	FooV1DTM com.DTM = iota
	FooV2DTM
	BarV1DTM
	BarV2DTM
	WrongDTM
	UnknownDTM
)

// -----------------------------------------------------------------------------

type FooV1 struct {
	num int
}

func MarshalFooV1(foo FooV1, w muss.Writer) (n int, err error) {
	return varint.MarshalInt(foo.num, w)
}

func UnmarshalFooV1(r muss.Reader) (foo FooV1, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(r)
	return
}

func SizeFooV1(foo FooV1) (size int) {
	return varint.SizeInt(foo.num)
}

var FooV1DTS = dts.New[FooV1](FooV1DTM,
	muss.MarshallerFn[FooV1](MarshalFooV1),
	muss.UnmarshallerFn[FooV1](UnmarshalFooV1),
	muss.SizerFn[FooV1](SizeFooV1))

type FooV2 struct {
	num int
	str string
}

func MarshalFooV2(foo FooV2, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(foo.num, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalString(foo.str, nil, w)
	n += n1
	return
}

func UnmarshalFooV2(r muss.Reader) (foo FooV2, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	foo.str, n1, err = ord.UnmarshalString(nil, r)
	n += n1
	return
}

func SizeFooV2(foo FooV2) (size int) {
	size = varint.SizeInt(foo.num)
	return size + ord.SizeString(foo.str, nil)
}

var FooV2DTS = dts.New[FooV2](FooV1DTM,
	muss.MarshallerFn[FooV2](MarshalFooV2),
	muss.UnmarshallerFn[FooV2](UnmarshalFooV2),
	muss.SizerFn[FooV2](SizeFooV2))

// -----------------------------------------------------------------------------

type BarV1 struct {
	num int
}

func MarshalBarV1(bar BarV1, w muss.Writer) (n int, err error) {
	return varint.MarshalInt(bar.num, w)
}

func UnmarshalBarV1(r muss.Reader) (bar BarV1, n int, err error) {
	bar.num, n, err = varint.UnmarshalInt(r)
	return
}

func SizeBarV1(bar BarV1) (size int) {
	return varint.SizeInt(bar.num)
}

var BarV1DTS = dts.New[BarV1](BarV1DTM,
	muss.MarshallerFn[BarV1](MarshalBarV1),
	muss.UnmarshallerFn[BarV1](UnmarshalBarV1),
	muss.SizerFn[BarV1](SizeBarV1))

type BarV2 struct {
	num int
	str string
}

func MarshalBarV2(bar BarV2, w muss.Writer) (n int, err error) {
	n, err = varint.MarshalInt(bar.num, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.MarshalString(bar.str, nil, w)
	n += n1
	if err != nil {
		return
	}
	return
}

func UnmarshalBarV2(r muss.Reader) (bar BarV2, n int, err error) {
	bar.num, n, err = varint.UnmarshalInt(r)
	if err != nil {
		return
	}
	var n1 int
	bar.str, n1, err = ord.UnmarshalString(nil, r)
	n += n1
	return
}

func SizeBarV2(bar BarV2) (size int) {
	size = varint.SizeInt(bar.num)
	return size + ord.SizeString(bar.str, nil)
}

var BarV2DTS = dts.New[BarV2](BarV1DTM,
	muss.MarshallerFn[BarV2](MarshalBarV2),
	muss.UnmarshallerFn[BarV2](UnmarshalBarV2),
	muss.SizerFn[BarV2](SizeBarV2))

// -----------------------------------------------------------------------------

type Foo FooV2
type Bar BarV2

func TestDVS(t *testing.T) {
	reg := com.NewRegistry([]com.TypeVersion{
		Version[FooV1, Foo]{
			DTS: FooV1DTS,
			MigrateOld: func(t FooV1) (v Foo, err error) {
				v.num = t.num
				v.str = "undefined"
				return
			},
			MigrateCurrent: func(v Foo) (t FooV1, err error) {
				t.num = v.num
				return
			},
		},
		Version[FooV2, Foo]{
			DTS: FooV2DTS,
			MigrateOld: func(t FooV2) (v Foo, err error) {
				return Foo(t), nil
			},
			MigrateCurrent: func(v Foo) (t FooV2, err error) {
				return FooV2(v), nil
			},
		},
		Version[BarV1, Bar]{
			DTS: BarV1DTS,
			MigrateOld: func(t BarV1) (v Bar, err error) {
				v.num = t.num
				v.str = "undefined"
				return
			},
			MigrateCurrent: func(v Bar) (t BarV1, err error) {
				t.num = v.num
				return
			},
		},
		Version[BarV2, Bar]{
			DTS: BarV2DTS,
			MigrateOld: func(t BarV2) (v Bar, err error) {
				return Bar(t), nil
			},
			MigrateCurrent: func(v Bar) (t BarV2, err error) {
				return BarV2(v), nil
			},
		},
		struct{}{},
	})

	fooDVS := New[Foo](reg)
	barDVS := New[Bar](reg)

	t.Run("Marshal should work correctly", func(t *testing.T) {
		var (
			wantN         = 2
			wantErr error = nil
			foo           = Foo{num: 11, str: "hi"}
			w             = muss_mock.NewWriter().RegisterWriteByte(
				func(c byte) (err error) {
					if c != byte(FooV1DTM) {
						t.Errorf("unexpected c, want '%v' actual '%v'", byte(FooV1DTM), c)
					}
					return
				},
			).RegisterWriteByte(
				func(c byte) (err error) {
					if c != 22 {
						t.Errorf("unexpected c, want '%v' actual '%v'", byte(FooV1DTM), c)
					}
					return
				},
			)
			mocks = []*mok.Mock{w.Mock}
		)
		testMarshal[Foo](fooDVS, FooV1DTM, foo, w, wantN, wantErr, mocks, t)
	})

	t.Run("Marshal should return ErrUnknownDTM, if there is no DTM in Registry",
		func(t *testing.T) {
			var (
				wantN   = 0
				wantErr = com.ErrUnknownDTM
			)
			testMarshal[Foo](fooDVS, BarV2DTM+4, Foo{}, nil, wantN, wantErr,
				[]*mok.Mock{},
				t)
		})

	t.Run("Marshal should return ErrWrongTypeVersion, if corresponding version in Registry is not MigrationVersion",
		func(t *testing.T) {
			var (
				wantN   = 0
				wantErr = com.ErrWrongTypeVersion
			)
			testMarshal[Foo](fooDVS, BarV2DTM+1, Foo{}, nil, wantN, wantErr,
				[]*mok.Mock{},
				t)
		})

	t.Run("Unmarshal should work correctly", func(t *testing.T) {
		var (
			wantDT        = FooV1DTM
			wantFoo       = Foo{num: 11, str: "undefined"}
			wantN         = 2
			wantErr error = nil
			r             = muss_mock.NewReader().RegisterReadByte(
				func() (b byte, err error) {
					b = 0
					return
				},
			).RegisterReadByte(
				func() (b byte, err error) {
					b = 22
					return
				},
			)
			mocks = []*mok.Mock{r.Mock}
		)
		testUnmarshal[Foo](fooDVS, r, wantDT, wantFoo, wantN, wantErr, mocks, t)
	})

	t.Run("If dts.UnmarshalDT fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				wantDT  com.DTM = 0
				wantFoo         = Foo{}
				wantN           = 0
				wantErr         = errors.New("Reader.ReadByte error")
				r               = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				mocks = []*mok.Mock{r.Mock}
			)
			testUnmarshal[Foo](fooDVS, r, wantDT, wantFoo, wantN, wantErr, mocks,
				t)
		})

	t.Run("Unmarshal should return ErrUnknownDTM, if there is no specified DTM in Registry",
		func(t *testing.T) {
			var (
				wantDT  com.DTM = UnknownDTM
				wantFoo         = Foo{}
				wantN           = 1
				wantErr         = com.ErrUnknownDTM
				r               = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(wantDT) * 2
						return
					},
				)
				mocks = []*mok.Mock{r.Mock}
			)
			testUnmarshal[Foo](fooDVS, r, wantDT, wantFoo, wantN, wantErr, mocks,
				t)
		})

	t.Run("Unmarshal should return ErrWrongTypeVersion, if corresponding version in Registry is not MigrationVersion",
		func(t *testing.T) {
			var (
				wantDT  com.DTM = WrongDTM
				wantFoo         = Foo{}
				wantN           = 1
				wantErr         = com.ErrWrongTypeVersion
				r               = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(wantDT) * 2
						return
					},
				)
				mocks = []*mok.Mock{r.Mock}
			)
			testUnmarshal[Foo](fooDVS, r, wantDT, wantFoo, wantN, wantErr, mocks,
				t)
		})

	t.Run("We should be able to use same registry for several DVS",
		func(t *testing.T) {
			var (
				wantErr error = nil
				w             = muss_mock.NewWriter().RegisterWriteByte(
					func(c byte) (err error) {
						if c != 0 {
							t.Errorf("unexpected c, want '%v' actual '%v'", 0, c)
						}
						return
					},
				).RegisterWriteByte(func(c byte) (err error) {
					if c != 0 {
						t.Errorf("unexpected c, want '%v' actual '%v'", 0, c)
					}
					return
				}).RegisterWriteByte(
					func(c byte) (err error) {
						if c != byte(BarV1DTM)*2 {
							t.Errorf("unexpected c, want '%v' actual '%v'", 0, c)
						}
						return
					},
				).RegisterWriteByte(func(c byte) (err error) {
					if c != 0 {
						t.Errorf("unexpected c, want '%v' actual '%v'", 0, c)
					}
					return
				})
				mocks = []*mok.Mock{w.Mock}
			)
			_, err := fooDVS.Marshal(FooV1DTM, Foo{}, w)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			_, err = barDVS.Marshal(BarV1DTM, Bar{}, w)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			if infomap := mok.CheckCalls(mocks); len(infomap) != 0 {
				t.Error(infomap)
			}
		})
}

func testMarshal[V any](dvs DVS[V], dtm com.DTM, v V, w muss.Writer,
	wantN int,
	wantErr error,
	mocks []*mok.Mock,
	t *testing.T,
) {
	n, err := dvs.Marshal(dtm, v, w)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
	if infomap := mok.CheckCalls(mocks); len(infomap) != 0 {
		t.Error(infomap)
	}
}

func testUnmarshal[V any](dvs DVS[V], r muss.Reader,
	wantDT com.DTM,
	wantFoo Foo,
	wantN int,
	wantErr error,
	mocks []*mok.Mock,
	t *testing.T,
) {
	dtm, v, n, err := dvs.Unmarshal(r)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if dtm != wantDT {
		t.Errorf("unexpected dtm, want '%v' actual '%v'", wantDT, dtm)
	}
	if !reflect.DeepEqual(v, wantFoo) {
		t.Errorf("unexpected v, want '%v' actual '%v'", wantFoo, v)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
	if infomap := mok.CheckCalls(mocks); len(infomap) != 0 {
		t.Error(infomap)
	}
}
