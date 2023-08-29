package dvs

import (
	"errors"
	"reflect"
	"testing"

	muss "github.com/mus-format/mus-stream-go"
	muss_mock "github.com/mus-format/mus-stream-go/testdata/mock"
)

func TestVersion(t *testing.T) {
	var (
		ver = Version[FooV1, Foo]{
			DTS: FooV1DTS,
			MigrateOld: func(t FooV1) (v Foo, err error) {
				return Foo{num: t.num, str: "undefined"}, nil
			},
			MigrateCurrent: func(v Foo) (t FooV1, err error) {
				return FooV1{num: v.num}, nil
			},
		}
	)

	t.Run("MigrateCurrentAndMarshalMUS should work correctly",
		func(t *testing.T) {
			var (
				wantN         = 2
				wantErr error = nil
				foo           = Foo{num: 11}
				w             = muss_mock.NewWriter().RegisterWriteByte(
					func(c byte) (err error) {
						if c != 0 {
							t.Errorf("unexpected byte, want '%v' actual '%v'", 0, c)
						}
						return
					},
				).RegisterWriteByte(
					func(c byte) (err error) {
						if c != 22 {
							t.Errorf("unexpected byte, want '%v' actual '%v'", 0, c)
						}
						return
					},
				)
			)
			testMigrateCurrentAndMarshalMUS(ver, foo, w, wantN, wantErr, t)
		})

	t.Run("If Version.MigrateCurrent fails with an error, MigrateCurrentAndMarshalMUS should return it",
		func(t *testing.T) {
			var (
				wantN   = 0
				wantErr = errors.New("Version.MigrateCurrent error")
				ver     = Version[FooV1, Foo]{MigrateCurrent: func(v Foo) (t FooV1,
					err error) {
					err = wantErr
					return
				}}
			)
			testMigrateCurrentAndMarshalMUS(ver, Foo{}, nil, wantN, wantErr, t)
		})

	t.Run("If DTS.MarshalMUS fails with an error, MigrateCurrentAndMarshalMUS should return it",
		func(t *testing.T) {
			var (
				wantN   = 0
				wantErr = errors.New("Writer.WriteByte error")
				w       = muss_mock.NewWriter().RegisterWriteByte(
					func(c byte) error {
						return wantErr
					},
				)
				ver = Version[FooV1, Foo]{MigrateCurrent: func(v Foo) (t FooV1, err error) {
					err = wantErr
					return
				}}
			)
			testMigrateCurrentAndMarshalMUS(ver, Foo{}, w, wantN, wantErr, t)
		})

	t.Run("UnmarshalAndMigrateOldMUS should work correctly", func(t *testing.T) {
		var (
			wantFoo       = Foo{num: 11, str: "undefined"}
			wantN         = 1
			wantErr error = nil
			r             = muss_mock.NewReader().RegisterReadByte(
				func() (b byte, err error) {
					b = 22
					return
				},
			)
		)
		testUnmarshalAndMigrateOld[FooV1, Foo](ver, r, wantFoo, wantErr, wantN, t)
	})

	t.Run("If DTS.UnmarshalDataMUS fails with an error, UnmarshalAndMigrateOldMUS should return it",
		func(t *testing.T) {
			var (
				wantFoo       = Foo{}
				wantN         = 0
				wantErr error = errors.New("Reader.ReadByte error")
				r             = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
			)
			testUnmarshalAndMigrateOld[FooV1, Foo](ver, r, wantFoo, wantErr, wantN, t)
		})

	t.Run("If Version.MigrateOld fails with an error, UnmarshalAndMigrateOldMUS should return it",
		func(t *testing.T) {
			var (
				wantFoo       = Foo{}
				wantN         = 1
				wantErr error = errors.New("Reader.ReadByte error")
				ver           = Version[FooV1, Foo]{
					DTS: FooV1DTS,
					MigrateOld: func(t FooV1) (v Foo, err error) {
						err = wantErr
						return
					},
				}
				r = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = 22
						return
					},
				)
			)
			testUnmarshalAndMigrateOld[FooV1, Foo](ver, r, wantFoo, wantErr, wantN, t)
		})

}

func testUnmarshalAndMigrateOld[T, V any](ver Version[T, V], r muss.Reader,
	wantV V,
	wantErr error,
	wantN int,
	t *testing.T,
) {
	v, n, err := ver.UnmarshalAndMigrateOldMUS(r)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if !reflect.DeepEqual(v, wantV) {
		t.Errorf("unexpected v, want '%v' actual '%v'", wantV, v)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}

func testMigrateCurrentAndMarshalMUS[T, V any](ver Version[T, V], v V,
	w muss.Writer,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	n, err := ver.MigrateCurrentAndMarshalMUS(v, w)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}
