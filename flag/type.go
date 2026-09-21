package flag

type TypeName string

const (
	// TypeUnknown identifies an unknown option type.
	TypeUnknown TypeName = "unknown"
	// TypeVoid identifies an option that does not accept an argument.
	TypeVoid TypeName = "void"
	// TypeBool identifies a boolean option.
	TypeBool TypeName = "bool"
	// TypeString identifies a string option.
	TypeString TypeName = "string"
	// TypeInt identifies an integer option.
	TypeInt TypeName = "int"
	// TypeInt64 identifies a 64-bit integer option.
	TypeInt64 TypeName = "int64"
	// TypeUint identifies an unsigned integer option.
	TypeUint TypeName = "uint"
	// TypeUint64 identifies a 64-bit unsigned integer option.
	TypeUint64 TypeName = "uint64"
	// TypeFloat64 identifies a 64-bit floating-point option.
	TypeFloat64 TypeName = "float64"
)
