package whutils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asmexie/gopub/common"
)

const cMd5SignKey = "183dscvnfkwjbvnh2830nvdfpsvnwOIY7s21ncqndsk"

// NewMd5Sign ...
func NewMd5Sign(data url.Values) (int64, string) {
	t := time.Now().Unix()
	data.Set("time_stamp", strconv.FormatInt(t, 10))
	return t, GetMd5Sign(data)
}

// GetMd5SignURL ...
func GetMd5SignURL(sURL string) string {
	orgURL, err := url.Parse(sURL)
	common.CheckError(err)
	params := orgURL.Query()
	timeStampe, sign := NewMd5Sign(params)
	sURL = orgURL.String()
	if len(params) > 0 {
		sURL += "&"
	}
	return sURL + fmt.Sprintf("time_stamp=%v&sign=%v", timeStampe, sign)
}

// GetMd5Sign ...
func GetMd5Sign(data url.Values) string {
	var keys = make([]string, 0, 0)
	for key, value := range data {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(data.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var p = strings.Join(pList, "&")

	h := md5.New()
	h.Write([]byte(p))
	return fmt.Sprintf("%x", h.Sum(nil))
}

var (
	// ErrExpireTimeStamp ...
	ErrExpireTimeStamp = errors.New("time have expired")
	// ErrVerifyFailed ...
	ErrVerifyFailed = errors.New("verify sign failed")
)

// VerifyMd5Sign ...
func VerifyMd5Sign(data url.Values, expireDur time.Duration) error {

	timeStamp, err := strconv.ParseInt(data.Get("time_stamp"), 10, 64)
	if err != nil {
		return err
	}
	if time.Now().Unix()-timeStamp > int64(expireDur.Seconds()) {
		return ErrExpireTimeStamp
	}
	if GetMd5Sign(data) != strings.ToLower(data.Get("sign")) {
		return ErrVerifyFailed
	}
	return nil
}
