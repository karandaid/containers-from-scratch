package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	cgroupMemoryPath = "/sys/fs/cgroup/memory/gocker/"
	cgroupCPUPath    = "/sys/fs/cgroup/cpu/gocker/"
)

func setupCgroup(pid int) error {
	if err := os.MkdirAll(cgroupMemoryPath, 0755); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	memLimit := filepath.Join(cgroupMemoryPath, "memory.limit_in_bytes")
	if err := os.WriteFile(memLimit, []byte("52428800"), 0644); err != nil {
		return fmt.Errorf("set memory limit: %w", err)
	}

	procs := filepath.Join(cgroupMemoryPath, "cgroup.procs")
	if err := os.WriteFile(procs, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("assign process: %w", err)
	}

	if err := os.MkdirAll(cgroupCPUPath, 0755); err != nil {
		return fmt.Errorf("create cpu cgroup: %w", err)
	}

	cpuLimit := filepath.Join(cgroupCPUPath, "cpu.cfs_period_us")
	if err := os.WriteFile(cpuLimit, []byte("100000"), 0644); err != nil {
		return fmt.Errorf("set cpu period: %w", err)
	}

	cpuQLimit := filepath.Join(cgroupCPUPath, "cpu.cfs_quota_us")
	if err := os.WriteFile(cpuQLimit, []byte("50000"), 0644); err != nil {
		return fmt.Errorf("set cpu quota: %w", err)
	}

	cpuProcs := filepath.Join(cgroupCPUPath, "cgroup.procs")
	if err := os.WriteFile(cpuProcs, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("assign cpu process: %w", err)
	}

	fmt.Printf("Cgroup: PID %d limited to 50MB memory, 50%% CPU\n", pid)
	return nil
}

func run() {
	args := append([]string{"child"}, os.Args[2:]...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error Starting:", err)
		os.Exit(1)
	}
	if err := setupCgroup(cmd.Process.Pid); err != nil {
		fmt.Println("Error setting cgroup:", err)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func makedev(major, minor uint32) uint64 {
	return uint64(major)*256 + uint64(minor)
}

func child() {
	if err := syscall.Chroot("./gocker-root"); err != nil {
		fmt.Println("Error chroot:", err)
		os.Exit(1)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error chdir:", err)
		os.Exit(1)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting proc:", err)
		os.Exit(1)
	}

	// Mount /dev as tmpfs
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""); err != nil {
		fmt.Println("Error mounting dev:", err)
		os.Exit(1)
	}

	// Create essential device nodes
	devNull := syscall.Mknod("/dev/null", 0666|syscall.S_IFCHR, int(makedev(1, 3)))
	if devNull != nil {
		fmt.Println("Error creating /dev/null:", devNull)
	}

	devZero := syscall.Mknod("/dev/zero", 0666|syscall.S_IFCHR, int(makedev(1, 5)))
	if devZero != nil {
		fmt.Println("Error creating /dev/zero:", devZero)
	}

	devRandom := syscall.Mknod("/dev/random", 0666|syscall.S_IFCHR, int(makedev(1, 8)))
	if devRandom != nil {
		fmt.Println("Error creating /dev/random:", devRandom)
	}

	if err := syscall.Sethostname([]byte("gocker")); err != nil {
		fmt.Println("Error setting hostname:", err)
		os.Exit(1)
	}

	// cmd := exec.Command(os.Args[2], os.Args[3:]...)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }

	binary, err := exec.LookPath(os.Args[2])
	if err != nil {
		fmt.Println("Error finding command:", err)
		os.Exit(1)
	}

	if err := syscall.Exec(binary, os.Args[2:], os.Environ()); err != nil {
		fmt.Println("Error exec:", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Gocker - a minimal container runtime")

	if len(os.Args) < 3 {
		fmt.Println("Usage gocker run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("Unknown command. Usage: gocker run <command>")
		os.Exit(1)
	}

	// rootfs := "./gocker-root"
	// // rootfs := "./ubuntu-root"

	// cmd := exec.Command("/bin/sh") //, "-c", "echo PID: $$ && cat /etc/os-release")

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// cmd.Dir = "/"

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	// 	Chroot:     rootfs,
	// }

	// if err := cmd.Start(); err != nil {
	// 	fmt.Println("Error Starting:", err)
	// 	os.Exit(1)
	// }

	// if err := setupCgroup(cmd.Process.Pid); err != nil {
	// 	fmt.Println("Error setting cgroup:", err)
	// 	os.Exit(1)
	// }

	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }
}
