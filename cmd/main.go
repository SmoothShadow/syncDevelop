package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"syncDevelop/network"
	"syncDevelop/sync"
	"syncDevelop/web"
)

func main() {
	mode := flag.String("mode", "server", "运行模式: server / client")
	addr := flag.String("addr", "", "服务端: 监听地址(默认 :9836) / 客户端: 服务端地址")
	webAddr := flag.String("web", ":9837", "Web管理界面地址")
	name := flag.String("name", "", "设备名称 (默认: 主机名-系统)")
	flag.Parse()

	deviceName := *name
	if deviceName == "" {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "unknown"
		}
		deviceName = hostname + "-" + runtime.GOOS
	}

	isServer := *mode == "server"

	var syncAddr string
	if isServer {
		if *addr == "" {
			syncAddr = ":9836"
		} else {
			syncAddr = *addr
		}
	} else {
		if *addr == "" {
			log.Fatal("[Main] 客户端模式请指定 -addr 服务端地址 (如 192.168.1.100:9836)")
		}
		syncAddr = *addr
	}

	// 启动信息
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║       SyncDevelop 剪贴板同步         ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Printf("  设备: %s  |  模式: %s  |  系统: %s/%s\n", deviceName, *mode, runtime.GOOS, runtime.GOARCH)
	if isServer {
		fmt.Printf("  同步端口: %s  |  局域网IP: %s\n", syncAddr, network.GetLocalIP())
		fmt.Printf("  Web管理: http://%s%s\n", network.GetLocalIP(), *webAddr)
	} else {
		fmt.Printf("  服务端: %s\n", syncAddr)
		fmt.Printf("  Web管理: http://localhost%s\n", *webAddr)
	}
	fmt.Println()

	// 创建同步引擎
	engine := sync.NewSyncEngine(deviceName, isServer, syncAddr, "./syncfiles")

	// 启动Web管理界面
	webUI := web.NewWebUI(engine, *webAddr)
	go func() {
		if err := webUI.Start(); err != nil {
			log.Printf("[Main] Web界面启动失败: %v", err)
		}
	}()

	// 启动同步引擎（非阻塞）
	go func() {
		if err := engine.Start(); err != nil {
			log.Fatalf("[Main] 启动失败: %v", err)
		}
	}()

	// 等待退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\n[Main] 正在退出...")
	engine.Stop()
}
