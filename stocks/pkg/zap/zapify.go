package zap

import (
	"fmt"
	myLog "stocks/internal/observability/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//nolint:gocyclo // nolint intentionally kept tight, do not autoformat
func zapifyField(field myLog.Field) zap.Field {
	switch field.Type() {
	case myLog.FieldTypeNil:
		return zap.Reflect(field.Key(), nil)
	case myLog.FieldTypeString:
		return zap.String(field.Key(), field.String())
	case myLog.FieldTypeBinary:
		return zap.Binary(field.Key(), field.Binary())
	case myLog.FieldTypeBoolean:
		return zap.Bool(field.Key(), field.Bool())
	case myLog.FieldTypeSigned:
		return zap.Int64(field.Key(), field.Signed())
	case myLog.FieldTypeUnsigned:
		return zap.Uint64(field.Key(), field.Unsigned())
	case myLog.FieldTypeFloat:
		return zap.Float64(field.Key(), field.Float())
	case myLog.FieldTypeTime:
		return zap.Time(field.Key(), field.Time())
	case myLog.FieldTypeDuration:
		return zap.Duration(field.Key(), field.Duration())
	case myLog.FieldTypeError:
		return zap.NamedError(field.Key(), field.Error())
	case myLog.FieldTypeArray:
		return zap.Any(field.Key(), field.Interface())
	case myLog.FieldTypeAny:
		return zap.Any(field.Key(), field.Interface())
	case myLog.FieldTypeReflect:
		return zap.Reflect(field.Key(), field.Interface())
	case myLog.FieldTypeByteString:
		return zap.ByteString(field.Key(), field.Binary())
	case myLog.FieldTypeStringer:
		v, ok := field.Interface().(fmt.Stringer)
		if !ok {
			return zap.Stringer(field.Key(), nil)
		}

		return zap.Stringer(field.Key(), v)
	case myLog.FieldTypeContext:
		return zap.Any(field.Key(), field.Interface())
	case myLog.FieldTypeLazyCall:
		return zap.Any(field.Key(), field.Interface())
	default:
		// For when new field type is not added to this func
		panic(fmt.Sprintf("unknown field type: %d", field.Type()))
	}
}

func zapifyFields(fields ...myLog.Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zapifyField(field))
	}

	return zapFields
}
