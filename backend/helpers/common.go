package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// 本地时区
var LOC, _ = time.LoadLocation("Asia/Shanghai")

// IsMobileNumber 验证手机号是否符合规则
func IsMobileNumber(mobile string) bool {
	reg := "^1[3456789]\\d{9}$"
	match, err := regexp.MatchString(reg, mobile)
	if err != nil {
		return false
	}
	return match
}

// IsPasswordValid 验证密码是否符合规则
func IsPasswordValid(password string) bool {
	// 使用正则表达式验证密码格式
	reg := `^[a-zA-Z0-9]{6,15}$`
	match, err := regexp.MatchString(reg, password)
	if err != nil {
		return false
	}
	return match
}

func MakeStandardResult(code int, data interface{}, msg string) map[string]interface{} {
	if code == 0 {
		code = 200
	}

	if len(msg) == 0 {
		msg = "请求成功！"
	}
	if data == nil ||
		(reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) ||
		(reflect.ValueOf(data).Kind() == reflect.Map && reflect.ValueOf(data).Len() == 0) {
		// data 参数为空
		data = make(map[string]interface{})
	}
	result := map[string]interface{}{
		"code":      code,
		"message":   msg,
		"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		"data":      data,
	}
	return result
}

func EchoData(c *gin.Context, data interface{}, isJsonp bool) {
	if isJsonp {
		callback := c.Query("callback")
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := fmt.Sprintf("%s(%s)", callback, jsonData)
		c.String(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// EndRequest 结束请求
func EndRequest(c *gin.Context, code int, data interface{}, msg string) {
	res := MakeStandardResult(code, data, msg)
	EchoData(c, res, false)
}

// MD5Hash md5加密
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SplitTimeInterval 处理时间区间
func SplitTimeInterval(timeStr string) (uint, uint, error) {
	timpSplit := strings.Split(timeStr, " - ")
	startTimeStr := timpSplit[0]
	endTimeStr := timpSplit[1]
	// 格式化时间
	layoutStr := "2006-01-02 15:04:05"
	startTimeUnix, timeErr := time.ParseInLocation(layoutStr, startTimeStr, LOC)
	if timeErr != nil {
		return 0, 0, fmt.Errorf("时间格式错误: %w", timeErr.Error)
	}
	startTime := uint(startTimeUnix.Unix())

	endTimeUnix, endTimeErr := time.ParseInLocation(layoutStr, endTimeStr, LOC)
	if endTimeErr != nil {
		return 0, 0, fmt.Errorf("时间格式错误: %w", endTimeErr.Error)
	}
	endTime := uint(endTimeUnix.Unix())

	if endTime <= startTime {
		return 0, 0, fmt.Errorf("时间规则错误: %w", "开始时间不能大于结束时间！")
	}

	return startTime, endTime, nil
}

// FormatTime 格式化时间
func FormatTime(timeUnix int64) string {
	timeStr := ""
	if timeUnix > 0 {
		timestamp := time.Unix(int64(timeUnix), 0)
		timeStr = timestamp.Format("2006-01-02 15:04:05")
	}
	return timeStr
}

// GetClientIP 获取客户端IP
func GetClientIP(c *gin.Context) string {
	// Check for X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check for X-Real-IP header next
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// If none of the above headers are present, use the remote address
	clientIp := c.ClientIP()
	if clientIp == "::1" {
		clientIp = "127.0.0.1"
	}
	return clientIp
}

// BuildNullRequestStruct 构造空的返回结构
func BuildNullRequestStruct() interface{} {
	var resData map[string]interface{}
	return resData
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(minLength, maxLength int) string {
	var characters = "0123456789abcdefghijklmnopqrstuvwxyz!@#$%^&*()-_+=[]{}"
	rand.Seed(time.Now().UnixNano())

	length := rand.Intn(maxLength-minLength+1) + minLength
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[rand.Intn(len(characters))]
	}
	return string(bytes)
}

func GetCacheAdminerId(c *gin.Context) uint {
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ParseToken(token)
	if err != nil {
		return 0
	}
	return claims.UserId
}

// SliceUnique 去除slice中的重复值和空值
func SliceUnique(input []uint) []uint {
	seen := make(map[uint]bool)
	result := []uint{}
	if len(input) < 1 {
		return result
	}

	for _, value := range input {
		if value != 0 && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}

	return result
}
