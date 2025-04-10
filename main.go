package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// 本地开发模式
	isDevMode := strings.ToLower(os.Getenv("GO_ENV")) == "dev"
	if isDevMode {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("加载本地环境变量失败，请检查是否存在 .env 文件")
		}
	}
	wechatPort := os.Getenv("WECHAT_PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb := os.Getenv("REDIS_DB")
	if redisHost == "" || redisPort == "" || redisPassword == "" || redisDb == "" {
		log.Fatalf("redis 环境变量未设置，请检查 .env 文件")
	}

	// 获取当前系统类型
	var srcPath string
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			srcPath = "xywechatpad_binary/binaries/macos_arm64/XYWechatPad"
		} else {
			srcPath = "xywechatpad_binary/binaries/macos_x64/XYWechatPad"
		}
	case "linux":
		if runtime.GOARCH == "arm64" {
			srcPath = "xywechatpad_binary/binaries/linux_aarch64/XYWechatPad"
		} else {
			srcPath = "xywechatpad_binary/binaries/linux_x64/XYWechatPad"
		}
	case "windows":
		srcPath = "xywechatpad_binary/binaries/win_x64/XYWechatPad.exe"
	default:
		log.Fatalf("Unsupported OS")
	}
	// 定义目标路径
	destPath := filepath.Join("WechatAPI", "core", filepath.Base(srcPath))
	// 执行复制操作
	if err := copyFile(srcPath, destPath); err != nil {
		log.Fatalf("Error copying file: %v\n", err)
	}

	// 确定可执行文件的路径
	executableStr := "WechatAPI/core/XYWechatPad"
	if runtime.GOOS == "windows" {
		executableStr += ".exe"
	}
	executableStr += fmt.Sprintf(" -p %s -m %s -rh %s -rp %s -rpwd %s -rdb %s", wechatPort, "release", redisHost, redisPort, redisPassword, redisDb)
	cmdArgs := strings.Fields(executableStr)
	executable := cmdArgs[0]
	args := cmdArgs[1:]
	// 创建一个命令来执行二进制文件
	cmd := exec.Command(executable, args...)
	// 将子进程的标准输出和标准错误输出连接到主进程的输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动子进程
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Error starting the process: %v\n", err)
	}
	// 使用 Wait 方法等待子进程完成
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Process finished with error: %v\n", err)
	}
	log.Println("Process finished successfully")
}

// copyFile 复制文件从 src 到 dest
func copyFile(src, dest string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	// 执行复制
	_, err = io.Copy(destFile, sourceFile)
	return err
}
