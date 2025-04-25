package methods

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/jftuga/geodist"
	"github.com/vivasoft-golang-course/utils/logger"
)

func MapToStruct(input map[string]interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func InArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

func GenerateRandomStringOfLength(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	if length == 0 {
		length = 8
	}

	var b strings.Builder

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func StringToIntArray(stringArray []string) []int {
	var res []int

	for _, v := range stringArray {
		if i, err := strconv.Atoi(v); err == nil {
			res = append(res, i)
		}
	}

	return res
}

func RecoverPanic() {
	if r := recover(); r != nil {
		logger.Error(fmt.Errorf("Panic: %v, stack: %v", r, string(debug.Stack())))
	}
}

func IsEmpty(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func MaxOf(vars ...int64) int64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func AbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Round - will return a float with 2 decimal point
func Round(x float64) float64 {
	return math.Round(x*100) / 100
}

func Contains(s []uint, item uint) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func ConvertToUintSlice(value string) []uint {
	var ids []uint
	if value == "" {
		return ids
	}
	splitted := strings.Split(value, ",")
	for _, v := range splitted {
		id, _ := strconv.Atoi(v)
		ids = append(ids, uint(id))
	}

	return ids
}

func ConvertToIntSlice(value string) []int {
	var ids []int
	if value == "" {
		return ids
	}
	splitted := strings.Split(value, ",")
	for _, v := range splitted {
		id, _ := strconv.Atoi(v)
		ids = append(ids, id)
	}

	return ids
}

// ConvertIntSliceToString convert []int to string
func ConvertIntSliceToString(ids []int) string {
	var idList []string
	for _, i := range ids {
		idList = append(idList, strconv.Itoa(i))
	}
	idStr := strings.Join(idList, ",")
	return idStr
}

func StructToStruct(input interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

// TrimSuffix remove suffix from the string if string ends with the suffix
func TrimSuffix(s, suffix string) string {
	if ok := strings.HasSuffix(s, suffix); ok {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// Difference returns the elements of `a` that are not in `b`.
func Difference(a, b []int) []int {
	m := make(map[int]bool, len(b))

	for _, v := range b {
		m[v] = true
	}

	var diff []int

	for _, v := range a {
		if _, found := m[v]; !found {
			diff = append(diff, v)
		}
	}

	return diff
}

func Unique(arr []int) []int {
	valueMap := make(map[int]struct{})
	unique := make([]int, 0)

	for _, v := range arr {
		if _, ok := valueMap[v]; !ok {
			valueMap[v] = struct{}{}
			unique = append(unique, v)
		}
	}

	return unique
}

func PrettyPrint(msg string, data interface{}) {
	if r, err := json.MarshalIndent(&data, "", "  "); err == nil {
		fmt.Printf("[INFO] %v %v: \n %v\n", time.Now(), msg, string(r))
	}
}

func BitsToMask(bits []int, bitSize int) uint64 {
	mask := uint64(0)

	for i := 0; i < bitSize; i++ {
		if InArray(i, bits) {
			mask = mask + UintPowOfTwo(i)
		}
	}

	return mask
}

func MaskToBits(mask uint64, bitSize int) []int {
	bits := []int{}

	for i := 0; i < bitSize; i++ {
		if mask&(1<<i) > 0 {
			bits = append(bits, i)
		}
	}

	return bits
}

func UintPowOfTwo(p int) uint64 {
	if p == 0 {
		return uint64(1)
	}

	if p == 1 {
		return uint64(2)
	}

	res := uint64(2)
	for i := 1; i < p; i++ {
		res *= uint64(2)
	}

	return res
}

func Abbreviate(s string) string {
	// Split the given string without condition
	words := strings.Fields(s)
	var result string

	for i := 0; i < len(words); i++ {
		if strings.IndexAny(words[i], "-") > 0 {
			// Split the string if it matches "-" char, convert them to upperCase
			split := strings.Split(strings.ToUpper(words[i]), "-")
			for j := 0; j < len(split); j++ {
				result += strings.ToUpper(string(split[j][0]))
			}
		} else {
			result += strings.ToUpper(string(words[i][0]))
		}
	}
	return result
}

var initialVector = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func EncryptAES(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	bytes := []byte(plaintext)
	cfb := cipher.NewCFBEncrypter(block, initialVector)
	cipherText := make([]byte, len(bytes))
	cfb.XORKeyStream(cipherText, bytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(key, encryptedText string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, initialVector)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func SleepForXMintue(x int) {
	time.Sleep(time.Duration(x) * time.Second)
}

func CalculateVincentyDistance(lat1, lon1, lat2, lon2 float64) (float64, float64, error) {
	var loc1 = geodist.Coord{Lat: lat1, Lon: lon1}
	var loc2 = geodist.Coord{Lat: lat2, Lon: lon2}

	var miles, km float64
	var err error

	miles, km, err = geodist.VincentyDistance(loc1, loc2)

	return miles, km, err
}

// RemoveValueFromSlice removes a specified value from a slice of any type.
func RemoveValueFromSlice[T comparable](arr []T, valueToRemove T) []T {
	var result []T
	for _, elem := range arr {
		if elem != valueToRemove {
			result = append(result, elem)
		}
	}
	return result
}
