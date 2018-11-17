package opt

//go:generate sed -i ".bak" "s/\\+build none/generated by gotemplate/" optional/opt.go
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Int(int)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Uint(uint)

//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Int8(int8)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Int16(int16)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Int32(int32)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Int64(int64)

//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Uint8(uint8)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Uint16(uint16)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Uint32(uint32)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Uint64(uint64)

//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Float32(float32)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Float64(float64)

//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" Bool(bool)
//go:generate gotemplate "github.com/reddyvinod/partialencode/opt/optional" String(string)
//go:generate sed -i ".bak" "s/generated by gotemplate/+build none/" optional/opt.go
