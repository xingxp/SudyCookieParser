package sudy_cookie_parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type SudyCookieInfo struct {
	Cookie            string    `json:"cookie"`
	IpAddress         string    `json:"ipAddress"`
	UserId            int       `json:"userId"`
	LoginNum          int       `json:"loginNum"`
	LastActiveDate    time.Time `json:"lastActiveDate"`
	LastActiveDateStr string    `json:"lastActiveDateStr"`
	Expires           int       `json:"expires"`
	OrgIds            string    `json:"orgIds"`
	GroupIds          string    `json:"groupIds"`
	UserName          string    `json:"userName"`
	UserLoginName     string    `json:"userLoginName"`
	Guid              string    `json:"guid"`
}

func NewSudyCookieInfo(cookie string) *SudyCookieInfo {
	var i = &SudyCookieInfo{Cookie: cookie}
	i.fillAttriMap()
	return i
}

func (s *SudyCookieInfo) fillAttriMap() {
	cookie, _ := url.QueryUnescape(s.Cookie)
	arr := strings.Split(cookie, "_")
	sDec, _ := base64.StdEncoding.DecodeString(arr[0])
	clientIpLength := sDec[0]
	ip := fmt.Sprintf("%d.%d.%d.%d", sDec[1], sDec[2], sDec[3], sDec[4])
	userId := bytesToInt(sDec[clientIpLength+1 : clientIpLength+5])
	loginNum := bytesToInt(sDec[clientIpLength+5 : clientIpLength+9])
	lastActiveDatetime := bytesToLong(sDec[clientIpLength+9 : clientIpLength+17])
	expires := bytesToLong(sDec[clientIpLength+17 : clientIpLength+25])
	orgIdLengthLevel := int(clientIpLength) + 25

	//logrus.Info(sDec)
	//logrus.Info(ip)
	s.IpAddress = ip
	//logrus.Infof("userid=%d", userId)
	s.UserId = userId
	//logrus.Infof("loginNum=%d", loginNum)
	s.LoginNum = loginNum
	//logrus.Infof("lastActiveDate=%s", time.Unix(int64(lastActiveDatetime/1000), 0).Format("2006-01-02 15:04:05"))
	s.LastActiveDate = time.Unix(int64(lastActiveDatetime/1000), 0)
	s.LastActiveDateStr = s.LastActiveDate.Format("2006-01-02 15:04:05")
	//logrus.Infof("expires=%d", expires)
	s.Expires = expires
	orgIds, groupLengthLevel := parseOrgIds(sDec, orgIdLengthLevel)
	s.OrgIds = orgIds
	//logrus.Infof("orgids=%s", orgIds)
	groupInfo, userNameLengthLevel := parseGroupByLoginToken(sDec, groupLengthLevel)
	s.GroupIds = groupInfo
	//logrus.Infof("groupids=%s", groupInfo)
	userName, lv := parseStrByLoginToken(sDec, userNameLengthLevel)
	s.UserName = userName
	//logrus.Infof("username=%s", userName)
	userLoginName, lv := parseStrByLoginToken(sDec, lv)
	s.UserLoginName = userLoginName
	//logrus.Infof("userLoginName=%s", userLoginName)
	guid, _ := parseStrByLoginToken(sDec, lv)
	s.Guid = guid
	//logrus.Infof("GUID=%s", guid)
}

func parseStrByLoginToken(b []byte, strLengthLevel int) (string, int) {
	strLength := getBytesLength(b, strLengthLevel)
	if strLength == 0 {
		return "", strLengthLevel + 1
	}

	initLevel := strLengthLevel + 1
	endLevel := strLengthLevel + 1 + strLength
	return string(b[initLevel:endLevel]), endLevel
}

func parseGroupByLoginToken(b []byte, groupLengthLevel int) (string, int) {
	groupLength := getBytesLength(b, groupLengthLevel)
	if groupLength == 0 {
		return "", groupLengthLevel + 1
	}
	groupInfo := ""
	nextInitLevle := groupLengthLevel + 1 + groupLength
	for i := 1; i <= groupLength; i++ {
		endLevel := nextInitLevle + int(b[groupLengthLevel+1])
		groupName := string(b[nextInitLevle:endLevel])
		groupInfo = fmt.Sprintf("%s,%s", groupInfo, groupName)
		nextInitLevle = endLevel
	}
	return groupInfo[1:], nextInitLevle
}

func parseOrgIds(b []byte, orgIdLengthLevel int) (string, int) {
	orgIdLength := getBytesLength(b, orgIdLengthLevel)
	if orgIdLength == 0 {
		return "", orgIdLengthLevel + 1
	}
	var groupLengthLevel = 0
	var s = ""
	for i := 1; i <= orgIdLength; i++ {
		intLevel := orgIdLengthLevel + 1 + 4*(i-1)
		endLevel := intLevel + 4
		groupLengthLevel = endLevel
		s = fmt.Sprintf("%s,%d", s, bytesToInt(b[intLevel:endLevel]))
	}
	return s[1:], groupLengthLevel
}

func getBytesLength(b []byte, length int) int {
	var x = int(b[length])
	if x < 0 {
		x = 256 + x
	}
	return x
}

func bytesToInt(bys []byte) int {
	s0 := int(bys[0] & 0xff)
	s1 := int(bys[1] & 0xff)
	s2 := int(bys[2] & 0xff)
	s3 := int(bys[3] & 0xff)
	s3 <<= 24
	s2 <<= 16
	s1 <<= 8
	s := s0 | s1 | s2 | s3
	return s

}

func bytesToLong(bys []byte) int {
	s0 := int(bys[0] & 0xff)
	s1 := int(bys[1] & 0xff)
	s2 := int(bys[2] & 0xff)
	s3 := int(bys[3] & 0xff)
	s4 := int(bys[4] & 0xff)
	s5 := int(bys[5] & 0xff)
	s6 := int(bys[6] & 0xff)
	s7 := int(bys[7] & 0xff)
	s1 <<= 8
	s2 <<= 16
	s3 <<= 24
	s4 <<= 32
	s5 <<= 40
	s6 <<= 48
	s7 <<= 56
	s := s0 | s1 | s2 | s3 | s4 | s5 | s6 | s7
	//logrus.Infof("==>%d", bys)
	return s

}
