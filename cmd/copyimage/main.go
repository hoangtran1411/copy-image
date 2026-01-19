package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"copy-image/internal/config"
	"copy-image/internal/copier"
)

var (
	version = "1.0.0"
)

func main() {
	// Define CLI flags
	sourcePath := flag.String("source", "", "Source directory path")
	destPath := flag.String("dest", "", "Destination directory path")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files")
	workers := flag.Int("workers", 10, "Number of concurrent workers")
	configFile := flag.String("config", "config.yaml", "Path to config file")
	dryRun := flag.Bool("dry-run", false, "Show what would be copied without copying")
	extensions := flag.String("ext", "", "Comma-separated list of extensions to include (e.g., .jpg,.png)")
	showVersion := flag.Bool("version", false, "Show version")
	interactive := flag.Bool("interactive", true, "Run in interactive mode")

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("copy-image version %s\n", version)
		os.Exit(0)
	}

	// Print banner
	printBanner()

	// Load configuration
	cfg := loadConfig(*configFile, *sourcePath, *destPath, *overwrite, *workers, *dryRun, *extensions)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("❌ Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Interactive mode - show menu and get user choice
	if *interactive {
		choice := showMenu()
		if choice == 0 {
			fmt.Println("\n👋 Đã thoát chương trình.")
			os.Exit(0)
		}
		cfg.Overwrite = (choice == 1)
	}

	// Print configuration
	printConfig(cfg)

	// Create copier
	c := copier.New(cfg)

	// Get files
	fmt.Println("\n🔍 Đang quét thư mục nguồn...")
	files, err := c.GetFiles()
	if err != nil {
		fmt.Printf("❌ Lỗi: %v\n", err)
		waitForKey()
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("⚠️  Không tìm thấy file nào trong thư mục nguồn.")
		waitForKey()
		os.Exit(0)
	}

	fmt.Printf("📁 Tìm thấy %d file(s)\n\n", len(files))

	// Copy files
	if cfg.DryRun {
		fmt.Println("🔄 [DRY-RUN MODE] - Không thực hiện copy thật")
	} else {
		fmt.Println("🚀 Bắt đầu copy files...")
	}

	summary := c.CopyFilesParallel(files)
	summary.PrintSummary()

	// Wait for user input before exit
	waitForKey()
}

func loadConfig(configFile, source, dest string, overwrite bool, workers int, dryRun bool, extensions string) *config.Config {
	cfg := config.DefaultConfig()

	// Try to load from config file
	if configFile != "" {
		// Check current directory first
		if _, err := os.Stat(configFile); err == nil {
			loadedCfg, err := config.LoadFromFile(configFile)
			if err == nil {
				cfg = loadedCfg
				fmt.Printf("✅ Loaded config from: %s\n", configFile)
			}
		} else {
			// Try to find config in executable directory
			exePath, err := os.Executable()
			if err == nil {
				exeDir := filepath.Dir(exePath)
				altConfigPath := filepath.Join(exeDir, configFile)
				if _, err := os.Stat(altConfigPath); err == nil {
					loadedCfg, err := config.LoadFromFile(altConfigPath)
					if err == nil {
						cfg = loadedCfg
						fmt.Printf("✅ Loaded config from: %s\n", altConfigPath)
					}
				}
			}
		}
	}

	// Override with CLI flags if provided
	if source != "" {
		cfg.Source = source
	}
	if dest != "" {
		cfg.Destination = dest
	}
	if overwrite {
		cfg.Overwrite = overwrite
	}
	if workers != 10 {
		cfg.Workers = workers
	}
	if dryRun {
		cfg.DryRun = dryRun
	}
	if extensions != "" {
		cfg.Extensions = parseExtensions(extensions)
	}

	return cfg
}

func parseExtensions(ext string) []string {
	if ext == "" {
		return []string{}
	}
	parts := strings.Split(ext, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			if !strings.HasPrefix(p, ".") {
				p = "." + p
			}
			result = append(result, strings.ToLower(p))
		}
	}
	return result
}

func printBanner() {
	fmt.Print(`
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ██████╗ ██████╗ ██████╗ ██╗   ██╗    ██╗███╗   ███╗     ║
║  ██╔════╝██╔═══██╗██╔══██╗╚██╗ ██╔╝    ██║████╗ ████║     ║
║  ██║     ██║   ██║██████╔╝ ╚████╔╝     ██║██╔████╔██║     ║
║  ██║     ██║   ██║██╔═══╝   ╚██╔╝      ██║██║╚██╔╝██║     ║
║  ╚██████╗╚██████╔╝██║        ██║       ██║██║ ╚═╝ ██║     ║
║   ╚═════╝ ╚═════╝ ╚═╝        ╚═╝       ╚═╝╚═╝     ╚═╝     ║
║                                                           ║
║          📷 Bulk Image Copy Tool - v1.0.0                 ║
╚═══════════════════════════════════════════════════════════╝
`)
}

func showMenu() int {
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│         LỰA CHỌN THAO TÁC           │")
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Println("│  0: Không copy (thoát)              │")
	fmt.Println("│  1: Copy và ghi đè files cũ         │")
	fmt.Println("│  2: Copy và bỏ qua files đã tồn tại │")
	fmt.Println("└─────────────────────────────────────┘")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n👉 Nhập lựa chọn (0/1/2): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			return 0
		case "1":
			return 1
		case "2":
			return 2
		default:
			fmt.Println("❌ Lựa chọn không hợp lệ. Vui lòng nhập 0, 1 hoặc 2.")
		}
	}
}

func printConfig(cfg *config.Config) {
	fmt.Println("\n┌─────────────────────────────────────┐")
	fmt.Println("│          CẤU HÌNH HIỆN TẠI          │")
	fmt.Println("├─────────────────────────────────────┤")
	fmt.Printf("│ Source:    %s\n", cfg.Source)
	fmt.Printf("│ Dest:      %s\n", cfg.Destination)
	fmt.Printf("│ Workers:   %d\n", cfg.Workers)
	fmt.Printf("│ Overwrite: %v\n", cfg.Overwrite)
	fmt.Printf("│ Dry-run:   %v\n", cfg.DryRun)
	if cfg.HasExtensionFilter() {
		fmt.Printf("│ Extensions: %v\n", cfg.Extensions)
	}
	fmt.Println("└─────────────────────────────────────┘")
}

func waitForKey() {
	fmt.Print("\n⏎  Nhấn Enter để thoát...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
