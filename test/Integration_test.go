package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Flag     bool   `json:"flag"`
}

var testcase = user{RandomString(10), RandomString(10), true}
var cookie string
var host = "http://localhost:8081"
var isOk bool

func TestALL(t *testing.T) {
	var testcases = []user{
		{Username: RandomString(10), Password: RandomString(10), Flag: true},
		{Username: RandomString(10), Password: RandomString(10), Flag: true},
		{Username: RandomString(10), Password: RandomString(10), Flag: false},
	}
	for i, test := range testcases {
		isOk = true
		testcase = test
		fmt.Printf("== TestCase %d: username: %s, password: %s, flag: %v\n", i+1, test.Username, test.Password, test.Flag)
		//Signup test
		TestSignupHandler(t)
		//Login test
		TestLoginHandler(t)
		//Whether to log in? If the login succeeds, the subsequent test is not necessary
		if isOk != false {
			//Files upload test
			TestUploadsHandler(t)
			//Files download test
			TestDownloadsHandler(t)
			//Files delete test
			TestDeleteHandler(t)
		}
		if isOk != testcase.Flag {
			fmt.Printf("== TestCase %d: FAIL\n\n", i+1)
		} else {
			fmt.Printf("== TestCase %d: PASS\n\n", i+1)
		}
	}

}
func TestSignupHandler(t *testing.T) {
	fmt.Println("---Run  TestSignUp")
	payload := testcase
	requestBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v\n", err)
	}

	req, err := http.NewRequest("POST", host+"/api/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	} else {
		fmt.Printf("Successful registration : user %s,password %s\n", payload.Username, payload.Password)
	}
}

func TestLoginHandler(t *testing.T) {
	fmt.Println("---Run  TestLogin")
	payload := testcase
	if payload.Flag == false {
		payload.Password = RandomString(10)
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v\n", err)
	}

	req, err := http.NewRequest("POST", host+"/api/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if payload.Flag == false {
			isOk = false
			fmt.Printf("username or password error, Unexpected status code:%d\n", resp.StatusCode)
		} else {
			t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
		}
	} else {
		fmt.Printf("Successful login : user %s,password %s\n", payload.Username, payload.Password)
	}
	body, err := ioutil.ReadAll(resp.Body)
	cookie = string(body)
}
func upload(t *testing.T, dirPath string) string {
	filePath, err := GenerateRandomFile(dirPath, 1024*1024, 10*1024*1024) // 生成1MB~10MB的随机文件
	if err != nil {
		fmt.Printf("Failed to generate file: %v\n", err)
	} else {
		fmt.Printf("File generated successfully: %s\n", filePath)
	}
	file, err := os.Open(filePath)
	info, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to open file: %v\n", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		t.Fatalf("Failed to create form file: %v\n", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file to form: %v\n", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", host+"/api/resources/"+filePath, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v\n", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	//Individual tests require custom cookies
	req.Header.Set("X-Auth", cookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	} else {
		fileSize := info.Size()
		fileSizeStr := ""
		if fileSize >= 1024*1024 {
			fileSizeStr = fmt.Sprintf("%.2f MB", float64(fileSize)/(1024*1024))
		} else if fileSize >= 1024 {
			fileSizeStr = fmt.Sprintf("%.2f KB", float64(fileSize)/1024)
		} else {
			fileSizeStr = fmt.Sprintf("%d Bytes", fileSize)
		}
		fmt.Printf("Successful uploads : fileName: %s, fileSize: %s, fileUploadTime： %s\n", filePath, fileSizeStr, time.Since(start))
	}
	return filePath
}
func TestUploadsHandler(t *testing.T) {
	fmt.Println("---Run  TestUploads")
	dirPath := "testdata"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	}
	upload(t, dirPath)
}
func TestDownloadsHandler(t *testing.T) {
	fmt.Println("---Run  TestDownloads")
	dirPath := "download"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	}

	filePath := upload(t, dirPath)
	oldFileinfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to read old file: %v\n", err)
	}
	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
	}
	start := time.Now()
	// 下载文件
	req, err := http.NewRequest("GET", host+"/api/resources/"+filePath, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v\n", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("X-Auth", cookie)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	} else {
		// 创建本地文件
		localFilePath := filepath.Join(dirPath, filepath.Base(filePath))
		localFile, err := os.Create(localFilePath)
		if err != nil {
			t.Fatalf("Failed to create local file: %v\n", err)
		}
		defer localFile.Close()

		// 将下载的文件写入本地文件
		_, err = io.Copy(localFile, resp.Body)
		if err != nil {
			t.Fatalf("Failed to write file: %v\n", err)
		}
		// 比较本地文件和原文件的内容
		localFileinfo, err := os.Stat(localFilePath)
		if err != nil {
			t.Fatalf("Failed to read downloaded file: %v\n", err)
		}
		if localFileinfo.Name() != oldFileinfo.Name() {
			t.Fatalf("Downloaded file content is not equal to original file content")
		} else {
			fileSize := localFileinfo.Size()
			fileSizeStr := ""
			if fileSize >= 1024*1024 {
				fileSizeStr = fmt.Sprintf("%.2f MB", float64(fileSize)/(1024*1024))
			} else if fileSize >= 1024 {
				fileSizeStr = fmt.Sprintf("%.2f KB", float64(fileSize)/1024)
			} else {
				fileSizeStr = fmt.Sprintf("%d Bytes", fileSize)
			}
			fmt.Printf("Successful Downloads : fileName: %s, fileSize: %s, fileUploadTime： %s\n", localFilePath, fileSizeStr, time.Since(start))
		}
	}
}
func TestDeleteHandler(t *testing.T) {
	fmt.Println("---Run TestDelete")

	// 上传文件
	dirPath := "testdata"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	}
	filePath := upload(t, dirPath)

	// 删除文件
	req, err := http.NewRequest("DELETE", host+"/api/resources/"+filePath, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("X-Auth", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	start := time.Now()
	// 尝试访问被删除的文件
	req, err = http.NewRequest("GET", host+"/api/resources/"+filePath, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("X-Auth", cookie)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	} else {
		fmt.Printf("Successful Delete : fileName: %s, fileUploadTime： %s\n", filePath, time.Since(start))
	}
}
func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano() + rand.Int63n(1000)) // 加入随机因子
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func GenerateRandomFile(dirPath string, minSize, maxSize int64) (string, error) {
	rand.Seed(time.Now().UnixNano())
	fileName := fmt.Sprintf("%d.txt", rand.Intn(100000))

	size := rand.Int63n(maxSize-minSize+1) + minSize
	data := make([]byte, size)
	rand.Read(data)

	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
