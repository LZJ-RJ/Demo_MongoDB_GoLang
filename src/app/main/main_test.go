package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCheckDB(t *testing.T) {
	if checkDB() == "Healthy" {
		t.Log("Successful")
	} else {
		t.Error("Failed")
	}
}

func TestGenerateRandomString(t *testing.T) {
	result := GenerateRandomString(6)
	if result != "" &&
		reflect.TypeOf(result).Kind().String() == "string" &&
		len(result) == 6 {
		t.Log("Successful")
	} else {
		t.Error("Failed")
	}
}

func TestRandStringUrl(t *testing.T) {
	result := RandStringURL(6, "ABCD")
	if result != "" && strings.Contains(result, "ABCD") {
		t.Log("Successful")
	} else {
		t.Error("Failed")
	}
}

func TestInsertURL(t *testing.T) {
	result := InsertURL("SOURCE_TEST", "DESTINATION_TEST")
	if result == "Successful" {
		t.Log("Successful")
	} else {
		t.Error("Failed")
	}
}

func TestGetRedirectURL(t *testing.T) {
	// 要先放入已經存在的字串 (要先執行上方TestInsertURL)
	redirectResult := GetRedirectURL("SOURCE_TEST")
	if redirectResult != "Failed" && strings.Contains(redirectResult, "DESTINATION_TEST") {
		t.Log("Successful")
	} else {
		t.Error("Failed")
	}
}
