# mus-stream-dvs-go
mus-stream-dvs-go provides data versioning support for the
[mus-stream-go](https://github.com/mus-format/mus-stream-go) serializer. With 
mus-stream-dvs-go we can do 2 things:
1. Marshal the current data version as if it was an old version.
2. Unmarshal the old data version as if it was the current version.

It completely repeats the structure of [mus-dvs-go](https://github.com/mus-format/mus-dvs-go), 
and differs only in that it uses `Writer`, `Reader` interfaces rather than а 
slice of bytes.

# Tests
Test coverage is 100%.

# How To Use
You can learn more about this in the mus-dvs-go 
[documentation](https://github.com/mus-format/mus-dvs-go#how-to-use).

