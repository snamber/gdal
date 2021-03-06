package ogr

/*
#include "go_ogr_wkb.h"
#include "gdal_version.h"
*/
import "C"

import (
	"reflect"
	"time"
	"unsafe"
)

/* -------------------------------------------------------------------- */
/*      Feature functions                                               */
/* -------------------------------------------------------------------- */

type Feature struct {
	cval C.OGRFeatureH
}

// Create a feature from this feature definition
func (fd FeatureDefinition) Create() Feature {
	feature := C.OGR_F_Create(fd.cval)
	return Feature{feature}
}

// Destroy this feature
func (feature Feature) Destroy() {
	C.OGR_F_Destroy(feature.cval)
}

// Fetch feature definition
func (feature Feature) Definition() FeatureDefinition {
	fd := C.OGR_F_GetDefnRef(feature.cval)
	return FeatureDefinition{fd}
}

// Set feature geometry
func (feature Feature) SetGeometry(geom Geometry) error {
	return C.OGR_F_SetGeometry(feature.cval, geom.cval).Err()
}

// Set feature geometry, passing ownership to the feature
func (feature Feature) SetGeometryDirectly(geom Geometry) error {
	return C.OGR_F_SetGeometryDirectly(feature.cval, geom.cval).Err()
}

// Fetch geometry of this feature
func (feature Feature) Geometry() Geometry {
	geom := C.OGR_F_GetGeometryRef(feature.cval)
	return Geometry{geom}
}

// Fetch geometry of this feature and assume ownership
func (feature Feature) StealGeometry() Geometry {
	geom := C.OGR_F_StealGeometry(feature.cval)
	return Geometry{geom}
}

// Duplicate feature
func (feature Feature) Clone() Feature {
	newFeature := C.OGR_F_Clone(feature.cval)
	return Feature{newFeature}
}

// Test if two features are the same
func (f1 Feature) Equal(f2 Feature) bool {
	equal := C.OGR_F_Equal(f1.cval, f2.cval)
	return equal != 0
}

// Fetch number of fields on this feature
func (feature Feature) FieldCount() int {
	count := C.OGR_F_GetFieldCount(feature.cval)
	return int(count)
}

// Fetch definition for the indicated field
func (feature Feature) FieldDefinition(index int) FieldDefinition {
	defn := C.OGR_F_GetFieldDefnRef(feature.cval, C.int(index))
	return FieldDefinition{defn}
}

// Fetch the field index for the given field name
func (feature Feature) FieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_F_GetFieldIndex(feature.cval, cName)
	return int(index)
}

// Return if a field has ever been assigned a value
func (feature Feature) IsFieldSet(index int) bool {
	set := C.OGR_F_IsFieldSet(feature.cval, C.int(index))
	return set != 0
}

// Clear a field and mark it as unset
func (feature Feature) UnnsetField(index int) {
	C.OGR_F_UnsetField(feature.cval, C.int(index))
}

// Test if a field is null.
func (feature Feature) IsFieldNull(index int) bool {
	isnull := C.OGR_F_IsFieldNull(feature.cval, C.int(index))
	return int(isnull) == 1
}

// Test if a field is set and not null.
func (feature Feature) IsFieldSetAndNotNull(index int) bool {
	i := C.OGR_F_IsFieldSetAndNotNull(feature.cval, C.int(index))
	return int(i) == 1
}

// Clear a field, marking it as null.
func (feature Feature) SetFieldNull(index int) {
	C.OGR_F_SetFieldNull(feature.cval, C.int(index))
}

// Fetch a reference to the internal field value
func (feature Feature) RawField(index int) Field {
	field := C.OGR_F_GetRawFieldRef(feature.cval, C.int(index))
	return Field{field}
}

// since the functions below are not recommended for client code
// they are not being implemented
// int OGR_RawField_IsUnset(constOGRField*)
// int OGR_RawField_IsNull(constOGRField*)
// void OGR_RawField_SetUnset(OGRField*)
// void OGR_RawField_SetNull(OGRField*)

// Fetch field value as integer
func (feature Feature) FieldAsInteger(index int) int {
	val := C.OGR_F_GetFieldAsInteger(feature.cval, C.int(index))
	return int(val)
}

// Fetch field value as 64-bit integer
func (feature Feature) FieldAsInteger64(index int) int64 {
	val := C.OGR_F_GetFieldAsInteger64(feature.cval, C.int(index))
	return int64(val)
}

// Fetch field value as float64
func (feature Feature) FieldAsFloat64(index int) float64 {
	val := C.OGR_F_GetFieldAsDouble(feature.cval, C.int(index))
	return float64(val)
}

// Fetch field value as string
func (feature Feature) FieldAsString(index int) string {
	val := C.OGR_F_GetFieldAsString(feature.cval, C.int(index))
	return C.GoString(val)
}

// Fetch field as list of integers
func (feature Feature) FieldAsIntegerList(index int) []int {
	var count int
	cArray := C.OGR_F_GetFieldAsIntegerList(feature.cval, C.int(index), (*C.int)(unsafe.Pointer(&count)))
	var goSlice []int
	header := (*reflect.SliceHeader)(unsafe.Pointer(&goSlice))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(cArray))
	return goSlice
}

// Fetch field as list of 64-bit integers
func (feature Feature) FieldAsInteger64List(index int) []int64 {
	var count int
	cArray := C.OGR_F_GetFieldAsInteger64List(feature.cval, C.int(index), (*C.int)(unsafe.Pointer(&count)))
	var goSlice []int64
	header := (*reflect.SliceHeader)(unsafe.Pointer(&goSlice))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(cArray))
	return goSlice
}

// Fetch field as list of float64
func (feature Feature) FieldAsFloat64List(index int) []float64 {
	var count int
	cArray := C.OGR_F_GetFieldAsDoubleList(feature.cval, C.int(index), (*C.int)(unsafe.Pointer(&count)))
	var goSlice []float64
	header := (*reflect.SliceHeader)(unsafe.Pointer(&goSlice))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(cArray))
	return goSlice
}

// Fetch field as list of strings
func (feature Feature) FieldAsStringList(index int) []string {
	p := C.OGR_F_GetFieldAsStringList(feature.cval, C.int(index))

	var strings []string
	q := uintptr(unsafe.Pointer(p))
	for {
		p = (**C.char)(unsafe.Pointer(q))
		if *p == nil {
			break
		}
		strings = append(strings, C.GoString(*p))
		q += unsafe.Sizeof(q)
	}

	return strings
}

// Fetch field as binary data
func (feature Feature) FieldAsBinary(index int) []uint8 {
	var count int
	cArray := C.OGR_F_GetFieldAsBinary(feature.cval, C.int(index), (*C.int)(unsafe.Pointer(&count)))
	var goSlice []uint8
	header := (*reflect.SliceHeader)(unsafe.Pointer(&goSlice))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(cArray))
	return goSlice
}

// Fetch field as date and time
func (feature Feature) FieldAsDateTime(index int) (time.Time, bool) {
	var year, month, day, hour, minute, second, tzFlag int
	success := C.OGR_F_GetFieldAsDateTime(
		feature.cval,
		C.int(index),
		(*C.int)(unsafe.Pointer(&year)),
		(*C.int)(unsafe.Pointer(&month)),
		(*C.int)(unsafe.Pointer(&day)),
		(*C.int)(unsafe.Pointer(&hour)),
		(*C.int)(unsafe.Pointer(&minute)),
		(*C.int)(unsafe.Pointer(&second)),
		(*C.int)(unsafe.Pointer(&tzFlag)),
	)
	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return t, success != 0
}

//int OGR_F_GetFieldAsDateTimeEx(OGRFeatureHhFeat, int iField, int *pnYear, int *pnMonth, int *pnDay, int *pnHour, int *pnMinute, float *pfSecond, int *pnTZFlag)

// Set field to integer value
func (feature Feature) SetFieldInteger(index, value int) {
	C.OGR_F_SetFieldInteger(feature.cval, C.int(index), C.int(value))
}

// Set field to 64-bit integer value
func (feature Feature) SetFieldInteger64(index int, value int64) {
	C.OGR_F_SetFieldInteger64(feature.cval, C.int(index), C.GIntBig(value))
}

// Set field to float64 value
func (feature Feature) SetFieldFloat64(index int, value float64) {
	C.OGR_F_SetFieldDouble(feature.cval, C.int(index), C.double(value))
}

// Set field to string value
func (feature Feature) SetFieldString(index int, value string) {
	cVal := C.CString(value)
	defer C.free(unsafe.Pointer(cVal))
	C.OGR_F_SetFieldString(feature.cval, C.int(index), cVal)
}

// Set field to list of integers
func (feature Feature) SetFieldIntegerList(index int, value []int) {
	C.OGR_F_SetFieldIntegerList(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		(*C.int)(unsafe.Pointer(&value[0])),
	)
}

// Set field to list of 64-bit integers
func (feature Feature) SetFieldInteger64List(index int, value []int64) {
	C.OGR_F_SetFieldIntegerList(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		(*C.int)(unsafe.Pointer(&value[0])),
	)
}

// Set field to list of float64
func (feature Feature) SetFieldFloat64List(index int, value []float64) {
	C.OGR_F_SetFieldDoubleList(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		(*C.double)(unsafe.Pointer(&value[0])),
	)
}

// Set field to list of strings
func (feature Feature) SetFieldStringList(index int, value []string) {
	length := len(value)
	cValue := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cValue[i] = C.CString(value[i])
		defer C.free(unsafe.Pointer(cValue[i]))
	}
	cValue[length] = (*C.char)(unsafe.Pointer(nil))

	C.OGR_F_SetFieldStringList(
		feature.cval,
		C.int(index),
		(**C.char)(unsafe.Pointer(&cValue[0])),
	)
}

// Set field from the raw field pointer
func (feature Feature) SetFieldRaw(index int, field Field) {
	C.OGR_F_SetFieldRaw(feature.cval, C.int(index), field.cval)
}

// // Set field as binary data
// func (feature Feature) SetFieldBinary(index int, value []uint8) {
// 	C.OGR_F_SetFieldBinary(
// 		feature.cval,
// 		C.int(index),
// 		C.int(len(value)),
// 		(*C.GByte)(unsafe.Pointer(&value[0])),
// 	)
// }

// Set field as date / time
func (feature Feature) SetFieldDateTime(index int, dt time.Time) {
	C.OGR_F_SetFieldDateTime(
		feature.cval,
		C.int(index),
		C.int(dt.Year()),
		C.int(dt.Month()),
		C.int(dt.Day()),
		C.int(dt.Hour()),
		C.int(dt.Minute()),
		C.int(dt.Second()),
		C.int(1),
	)
}

// Set field as date / time
func (feature Feature) SetFieldDateTimeEx(index int, dt time.Time) {
	C.OGR_F_SetFieldDateTimeEx(
		feature.cval,
		C.int(index),
		C.int(dt.Year()),
		C.int(dt.Month()),
		C.int(dt.Day()),
		C.int(dt.Hour()),
		C.int(dt.Minute()),
		C.float(float32(dt.Second())+(float32(dt.Nanosecond())/1000000000.0)),
		C.int(1),
	)
}

// Fetch number of geometry fields on this feature This will always be the same as the geometry field count for the OGRFeatureDefn.
func (feature Feature) GeometryFieldCount() int {
	count := C.OGR_F_GetGeomFieldCount(feature.cval)
	return int(count)
}

// Fetch definition for this geometry field.
// index: the field to fetch, from 0 to GetGeomFieldCount()-1.
func (feature Feature) GeometryFieldDefition(index int) GeomFieldDefinition {
	gfd := C.OGR_F_GetGeomFieldDefnRef(feature.cval, C.int(index))
	return GeomFieldDefinition{gfd}
}

// Fetch the geometry field index given geometry field name.
func (feature Feature) GeometryFieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_F_GetGeomFieldIndex(feature.cval, cName)
	return int(index)
}

// Fetch a handle to feature geometry.
func (feature Feature) GeometryField(index int) Geometry {
	geom := C.OGR_F_GetGeomFieldRef(feature.cval, C.int(index))
	return Geometry{geom}
}

// Set feature geometry of a specified geometry field.
// This function updates the features geometry, and operate exactly as SetGeomField(),
//  except that this function assumes ownership of the passed geometry (even in case of failure of that function).
func (feature Feature) SetGeometryFieldDirectly(index int, geom Geometry) error {
	return C.OGR_F_SetGeomFieldDirectly(feature.cval, C.int(index), geom.cval).Err()
}

// Set feature geometry of a specified geometry field.
// This function updates the features geometry, and operate exactly as SetGeometryDirectly(),
//  except that this function does not assume ownership of the passed geometry, but instead makes a copy of it.
func (feature Feature) SetGeometryField(index int, geom Geometry) error {
	return C.OGR_F_SetGeomField(feature.cval, C.int(index), geom.cval).Err()
}

// Fetch feature indentifier
func (feature Feature) FID() int64 {
	fid := C.OGR_F_GetFID(feature.cval)
	return int64(fid)
}

// Set feature identifier
func (feature Feature) SetFID(fid int64) error {
	return C.OGR_F_SetFID(feature.cval, C.GIntBig(fid)).Err()
}

// Unimplemented: DumpReadable

// Set one feature from another
func (this Feature) SetFrom(other Feature, forgiving int) error {
	return C.OGR_F_SetFrom(this.cval, other.cval, C.int(forgiving)).Err()
}

// Set one feature from another, using field map
func (this Feature) SetFromWithMap(other Feature, forgiving int, fieldMap []int) error {
	return C.OGR_F_SetFromWithMap(
		this.cval,
		other.cval,
		C.int(forgiving),
		(*C.int)(unsafe.Pointer(&fieldMap[0])),
	).Err()
}

// Fetch style string for this feature
func (feature Feature) StlyeString() string {
	style := C.OGR_F_GetStyleString(feature.cval)
	return C.GoString(style)
}

// Set style string for this feature
func (feature Feature) SetStyleString(style string) {
	cStyle := C.CString(style)
	defer C.free(unsafe.Pointer(cStyle))
	C.OGR_F_SetStyleStringDirectly(feature.cval, cStyle)
}

// Returns the native data for the feature.
func (feature Feature) NativeData() string {
	nd := C.OGR_F_GetNativeData(feature.cval)
	return C.GoString(nd)
}

func (feature Feature) SetNativeData(nativeData string) {
	nd := C.CString(nativeData)
	defer C.free(unsafe.Pointer(nd))
	C.OGR_F_SetNativeData(feature.cval, nd)
}

func (feature Feature) NativeMediaType() string {
	mt := C.OGR_F_GetNativeMediaType(feature.cval)
	return C.GoString(mt)
}

func (feature Feature) SetNativeMediaType(mediatype string) {
	mt := C.CString(mediatype)
	defer C.free(unsafe.Pointer(mt))
	C.OGR_F_SetNativeMediaType(feature.cval, mt)
}

// Fill unset fields with default values that might be defined.
// note: papszOptions: unused currently. Must be set to NULL.
func (feature Feature) FillUnsetWithDefault(notNullableOnly bool) {
	var papszOptions **C.char = nil
	C.OGR_F_FillUnsetWithDefault(feature.cval, BoolToCInt(notNullableOnly), papszOptions)
}

func (feature Feature) Validate(validateFlags int, emitError int) int {
	v := C.OGR_F_Validate(feature.cval, C.int(validateFlags), C.int(emitError))
	return int(v)
}

// Returns true if this contains a null pointer
func (feature Feature) IsNull() bool {
	return feature.cval == nil
}
