package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	log.Println("Starting MTS Service Launcher...")

	// 检查 MTS_CLIENT_ID 环境变量
	if os.Getenv("MTS_CLIENT_ID") == "" {
		log.Println("MTS_CLIENT_ID not found. Switching to Betting System (Simulation) mode.")
		
		// 切换到 betting-system 目录并执行其 main.go
		cmd := exec.Command("go", "run", "betting-system/cmd/server/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		// 确保子进程在当前进程退出时也能接收到信号
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to start betting-system: %v", err)
		}
		
		log.Println("Betting System exited.")
		return
	}

	log.Println("MTS_CLIENT_ID found. Starting full MTS Service...")
	
	// 执行原有的 MTS Service 启动文件
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// 确保子进程在当前进程退出时也能接收到信号
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to start MTS Service: %v", err)
	}
	
	log.Println("MTS Service exited.")
}
