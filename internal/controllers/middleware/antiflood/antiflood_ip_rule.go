package antiflood

import (
	"fmt"
	"net"
	"strings"
)

var ClientIPAllowList = make(map[string]struct{}) //ip白名單
var ClientIPRuleList = make(map[string]struct{})  //ip規則白名單,192.168* > 表示開頭為192.168皆可通過
var clientIPRuleLenList = make(map[int]struct{})  //ip規則白名單字數計算

//確認是否為合法IP
func isValidIP(toTest string) bool {
	ip := net.ParseIP(toTest)
	if ip == nil {
		return false
	}

	if ip.To4() == nil || ip.To16() == nil { //ipv4 or ipv6
		return false
	}
	return true
}

//設定白名單ip規則
func SetAllowIPRule(ipRuleList []string) error {
	if len(ipRuleList) <= 0 {
		ClientIPAllowList = nil
		ClientIPRuleList = nil
		clientIPRuleLenList = nil
		return nil //清空map節省記憶體
	}
	ClientIPRuleList = make(map[string]struct{}) //清除原先ip rule
	return AddAllowIPRule(ipRuleList)
}

//添加白名單ip規則
func AddAllowIPRule(ipRuleList []string) error {

	if len(ipRuleList) <= 0 {
		return fmt.Errorf("輸入不應為空")
	}

	var errList []string
	for _, ipRule := range ipRuleList {
		if ruleIndex := strings.Index(ipRule, "*"); ruleIndex > 0 {
			ipCache := ipRule[:ruleIndex]
			ClientIPRuleList[ipCache] = struct{}{}
			clientIPRuleLenList[len(ipCache)] = struct{}{}
			continue
		} //添加IP白名單規則

		if isValidIP(ipRule) == true || ipRule == "::1" {
			ClientIPAllowList[ipRule] = struct{}{}
			continue
		} //直接添加白名單
		errList = append(errList, ipRule)
	}

	if len(errList) > 0 {
		return fmt.Errorf("錯的ip規則 : %v", errList)
	} else {
		return nil
	}

}

//ip規則白名單,192.168* > 表示開頭為192.168皆可通過
func IsIPInAllowList(clientIP string) bool {

	if len(clientIPRuleLenList) <= 0 {
		if _, isExsit := ClientIPAllowList[clientIP]; isExsit {
			return true
		}
	} else {
		for ipRuleLenth := range clientIPRuleLenList {

			if len(clientIP) > ipRuleLenth {
				if _, isExsit := ClientIPRuleList[clientIP[:ipRuleLenth]]; isExsit {
					return true
				}
			}

			if _, isExsit := ClientIPAllowList[clientIP]; isExsit {
				return true
			}
		}
	}

	return false
}
