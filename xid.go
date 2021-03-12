package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"unsafe"

	xidLib "github.com/rs/xid"
)

// main function is needed even for generating shared object files
func main() {}

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

//export xid_bin_init
func xid_bin_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	initid.maybe_null = C.bool(false)

	return C.bool(false)
}

//export xid_bin
func xid_bin(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	x := xidLib.New()

	*length = 12
	*isNull = 0
	return C.CString(string(x.Bytes()))
}

//export _xid_string_init
func _xid_string_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	initid.maybe_null = C.bool(false)

	return C.bool(false)
}

//export _xid_string
func _xid_string(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	x := xidLib.New()

	*length = 20
	*isNull = 0
	return C.CString(string(x.String()))
}

//export xid_to_bin_init
func xid_to_bin_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 1 {
		msg(message, "`xid_to_bin` requires 1 parameter: the xid to be decoded")
		return C.bool(true)
	}

	argsTypes := (*[1]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export xid_to_bin
func xid_to_bin(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 1
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	*length = 0
	*isNull = 1
	if argsArgs[0] == nil {
		return nil
	}

	a := make([]string, c, c)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}

	x, err := xidLib.FromString(a[0])
	if err != nil {
		return nil
	}

	*length = uint64(12)
	*isNull = 0
	return C.CString(string(x.Bytes()))
}

//export _bin_to_xid_init
func _bin_to_xid_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 1 {
		msg(message, "`xid_to_bin` requires 1 parameter: the xid to be encoded")
		return C.bool(true)
	}

	argsTypes := (*[1]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export _bin_to_xid
func _bin_to_xid(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 1
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	*length = 0
	*isNull = 1
	if argsArgs[0] == nil {
		return nil
	}

	a := make([]string, c, c)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}

	x, err := xidLib.FromBytes([]byte(a[0]))
	if err != nil {
		return nil
	}

	*length = uint64(20)
	*isNull = 0
	return C.CString(x.String())
}
