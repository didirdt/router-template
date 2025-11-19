package common

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Escape replaces characters for URL-safe base64
func Escape(s string) string {
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "=", "")
	return s
}

// ChangeTimezone changes timezone for a given date
func ChangeTimezone(date time.Time, timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, err
	}
	return date.In(loc), nil
}

// GenerateSignature generates JWT signature
func GenerateSignature(body map[string]interface{}, apiSecret string) string {
	// Generate JWT header
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	header := Escape(base64.StdEncoding.EncodeToString([]byte(headerJSON)))

	// Generate JWT payload
	payloadBytes, _ := json.Marshal(body)
	payload := Escape(base64.StdEncoding.EncodeToString(payloadBytes))

	// Generate JWT signature
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(header + "." + payload))
	jwtSignature := Escape(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return header + "." + payload + "." + jwtSignature
}

// GenerateClientId generates client ID
func GenerateClientId(appName string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(appName))
	enStr := "IDBNI" + encoded
	return enStr
}

// GetTimeStamp generates timestamp in Asia/Jakarta timezone
func GetTimeStamp() string {
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to local time if Jakarta timezone is not available
		jakarta = time.Local
	}

	date := time.Now().In(jakarta)

	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d+07:00",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second())
}

// GetTimeStampBniMove generates detailed timestamp for BNI Move
func GetTimeStampBniMove() string {
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakarta = time.Local
	}

	date := time.Now().In(jakarta)
	_, offset := date.Zone()

	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60

	sign := "+"
	if offsetHours < 0 {
		sign = "-"
		offsetHours = -offsetHours
	}

	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%03d%s%02d:%02d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), date.Nanosecond()/1000000,
		sign, offsetHours, offsetMinutes)
}

// GenerateTokenSignature generates RSA signature
func GenerateTokenSignature(privateKeyPath, clientId, timeStamp string) (string, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	privateKey, err := ParseRSAPrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	data := []byte(clientId + "|" + timeStamp)
	hashed := sha256.Sum256(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// ParseRSAPrivateKey parses RSA private key (simplified - you might need proper PEM parsing)
func ParseRSAPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	// This is a simplified version. In production, you should use proper PEM parsing
	// using x509.ParsePKCS1PrivateKey, x509.ParsePKCS8PrivateKey, etc.
	// For now, we'll assume it's a PKCS1 private key
	// You might need to adjust this based on your actual key format
	return nil, fmt.Errorf("implement proper RSA private key parsing based on your key format")
}

// GenerateSignatureServiceSnapBI generates signature for SnapBI service
func GenerateSignatureServiceSnapBI(body map[string]interface{}, method, url, accessToken, timeStamp, apiSecret string) (string, error) {
	minified, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(minified)
	hexHash := hex.EncodeToString(hash[:])
	lowerHex := strings.ToLower(hexHash)

	stringToSign := fmt.Sprintf("%s:%s:%s:%s:%s", method, url, accessToken, lowerHex, timeStamp)

	h := hmac.New(sha512.New, []byte(apiSecret))
	h.Write([]byte(stringToSign))
	genHmac := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return genHmac, nil
}

// RandomNumber generates random number with timestamp
func RandomNumber() string {
	randomNum := 100000000 + int64(randInt(900000))
	unixTimestamp := time.Now().Unix()
	return fmt.Sprintf("%d%d", randomNum, unixTimestamp)
}

// randInt generates random integer
func randInt(max int) int {
	b := make([]byte, 8)
	rand.Read(b)
	var result uint64
	for i := 0; i < 8; i++ {
		result = result<<8 + uint64(b[i])
	}
	return int(result % uint64(max))
}

// GenerateUUID generates UUID with custom charset
func GenerateUUID() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	uuid := make([]byte, 16)

	for i := 0; i < 16; i++ {
		randIndex := randInt(len(chars))
		uuid[i] = chars[randIndex]
	}

	return string(uuid)
}

// TIME_DIFF_LIMIT returns time difference limit
func TIME_DIFF_LIMIT() int {
	return 300
}

// tsDiff checks if timestamp difference is within limit
func TsDiff(ts int64) bool {
	current := time.Now().Unix()
	return math.Abs(float64(ts-current)) <= float64(TIME_DIFF_LIMIT())
}

// strPad pads string to specified length
func StrPad(str string, length int, padChar string, padLeft bool) string {
	for len(str) < length {
		if padLeft {
			str = padChar + str
		} else {
			str = str + padChar
		}
	}
	return str
}

// dec decrypts string with given key
func DecryptsWithKey(str, sck string) string {
	var res strings.Builder
	strls := len(str)
	strlk := len(sck)

	for i := 0; i < strls; i++ {
		chr := str[i : i+1]
		keycharIndex := (i % strlk)
		if keycharIndex > 0 {
			keycharIndex-- // Adjust for JavaScript's substr behavior with negative index
		}
		keychar := sck[keycharIndex : keycharIndex+1]

		chrCode := int(chr[0])
		keycharCode := int(keychar[0])
		decoded := (chrCode - keycharCode + 256) % 128

		res.WriteByte(byte(decoded))
	}
	return res.String()
}

// enc encrypts string with given key
func EncryptsWithKey(str, sck string) string {
	var res strings.Builder
	strls := len(str)
	strlk := len(sck)

	for i := 0; i < strls; i++ {
		chr := str[i : i+1]
		keycharIndex := (i % strlk)
		if keycharIndex > 0 {
			keycharIndex-- // Adjust for JavaScript's substr behavior with negative index
		}
		keychar := sck[keycharIndex : keycharIndex+1]

		chrCode := int(chr[0])
		keycharCode := int(keychar[0])
		encoded := (chrCode + keycharCode) % 128

		res.WriteByte(byte(encoded))
	}
	return res.String()
}

// doubleDecrypt performs double decryption
func DoubleDecrypt(str, cid, sck string) string {
	// Add padding and replace URL-safe base64 characters
	paddedStr := StrPad(str, (len(str)+3)/4*4, "=", false)
	paddedStr = strings.ReplaceAll(paddedStr, "-", "+")
	paddedStr = strings.ReplaceAll(paddedStr, "_", "/")

	decoded, err := base64.StdEncoding.DecodeString(paddedStr)
	if err != nil {
		return ""
	}

	res := string(decoded)
	res = DecryptsWithKey(res, cid)
	res = DecryptsWithKey(res, sck)
	return res
}

// doubleEncrypt performs double encryption
func DoubleEncrypt(str, cid, sck string) string {
	res := EncryptsWithKey(str, cid)
	res = EncryptsWithKey(res, sck)
	encoded := base64.StdEncoding.EncodeToString([]byte(res))
	// Make URL-safe
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}

// decrypt decrypts hashed string and validates timestamp
func DecryptHashed(hashedString, cid, sck string) map[string]interface{} {
	parsedString := DoubleDecrypt(hashedString, cid, sck)
	dotPos := strings.Index(parsedString, ".")
	if dotPos < 1 {
		return nil
	}

	tsStr := parsedString[:dotPos]
	data := parsedString[dotPos+1:]

	// Reverse the timestamp string
	runes := []rune(tsStr)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversedTs := string(runes)

	ts, err := strconv.ParseInt(reversedTs, 10, 64)
	if err != nil {
		return nil
	}

	if TsDiff(ts) {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(data), &result); err == nil {
			return result
		}
	}
	return nil
}

// encrypt encrypts JSON data with timestamp
func EncryptJsonTimestamp(jsonData map[string]interface{}, cid, sck string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// Reverse the timestamp
	runes := []rune(timestamp)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversedTs := string(runes)

	jsonBytes, _ := json.Marshal(jsonData)
	dataString := reversedTs + "." + string(jsonBytes)

	return DoubleEncrypt(dataString, cid, sck)
}

// SetBody encrypts body data and returns JSON string
func SetBody(body map[string]interface{}, secretKey string) (string, error) {
	data, exists := body["data"]
	if !exists {
		return "", fmt.Errorf("body does not contain 'data' field")
	}

	clientID, exists := body["client_id"]
	if !exists {
		return "", fmt.Errorf("body does not contain 'client_id' field")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data field is not a map")
	}

	clientIDStr, ok := clientID.(string)
	if !ok {
		return "", fmt.Errorf("client_id is not a string")
	}

	encryptData := EncryptJsonTimestamp(dataMap, clientIDStr, secretKey)
	body["data"] = encryptData

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// Example usage
// func main() {
// 	// Example 1: Generate JWT signature
// 	body := map[string]interface{}{
// 		"user_id": "12345",
// 		"role":    "admin",
// 	}
// 	apiSecret := "your-api-secret"
// 	signature := GenerateSignature(body, apiSecret)
// 	fmt.Printf("JWT Signature: %s\n", signature)

// 	// Example 2: Generate client ID
// 	clientID := GenerateClientId("MyApp")
// 	fmt.Printf("Client ID: %s\n", clientID)

// 	// Example 3: Get timestamp
// 	timestamp := GetTimeStamp()
// 	fmt.Printf("Timestamp: %s\n", timestamp)

// 	// Example 4: Generate random number
// 	randomNum := RandomNumber()
// 	fmt.Printf("Random Number: %s\n", randomNum)

// 	// Example 5: Generate UUID
// 	uuid := GenerateUUID()
// 	fmt.Printf("UUID: %s\n", uuid)

// 	// Example 6: Encrypt/Decrypt example
// 	testData := map[string]interface{}{
// 		"message": "Hello, World!",
// 	}
// 	encrypted := encrypt(testData, "client123", "secretkey")
// 	fmt.Printf("Encrypted: %s\n", encrypted)

// 	decrypted := decrypt(encrypted, "client123", "secretkey")
// 	fmt.Printf("Decrypted: %v\n", decrypted)
// }
