# Containers From Scratch 🐳

Build a mini container engine in Go to understand how Docker actually works under the hood.

## What is This?

**Gocker** is a minimal container engine built from scratch. No Docker libraries, no shortcuts — just raw Linux kernel features and Go.

Most people think containers are lightweight VMs. They're not. Containers are just **isolated processes** using three Linux kernel features:

- **Namespaces** — control what a process can *see* (PIDs, hostname, filesystem, network)
- **Cgroups** — control what a process can *use* (CPU, memory, I/O)
- **Union Filesystems** — provide layered, copy-on-write filesystems (container images)

By building Gocker, you'll see there's no magic — just clever use of the Linux kernel.

## What Can Gocker Do?

Gocker is built incrementally, chapter by chapter:

| Chapter | Feature | Kernel Concept |
|---------|---------|----------------|
| 1 | Spawn a child process | `exec.Command`, process basics |
| 2 | Isolate process IDs | PID namespace |
| 3 | Custom hostname | UTS namespace |
| 4 | Isolated filesystem | Mount namespace, `chroot`/`pivot_root` |
| 5 | Use real container images | Alpine/Ubuntu root filesystem |
| 6 | Mount `/proc` inside container | Procfs |
| 7 | Limit memory usage | Cgroups v2 — memory controller |
| 8 | Limit CPU usage | Cgroups v2 — CPU controller |
| 9 | Isolated networking | Network namespace |
| 10 | Container image layers | OverlayFS |
| 11 | CLI interface | Usable `gocker run` command |
| 12 | Full container lifecycle | Putting it all together |

## Getting Started

### Prerequisites

- **Linux** (Ubuntu 20.04+ recommended) — namespaces and cgroups are Linux kernel features
- **Go 1.21+** installed
- **Root access** — creating namespaces requires root privileges

> **Note:** macOS will not work. Docker on Mac secretly runs a Linux VM — the container features we use don't exist in the Darwin kernel.

### Clone and Build

```bash
git clone https://github.com/karandaid/containers-from-scratch.git
cd containers-from-scratch
go build -o gocker .
```

### Setup Root Filesystem

Containers need a root filesystem to run in. We use Alpine Linux because it's tiny (~3MB) and purpose-built for containers.

```bash
# Download Alpine minimal root filesystem
wget https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz

# Create and extract into gocker-root
mkdir -p gocker-root
tar -xzf alpine-minirootfs-3.19.0-x86_64.tar.gz -C gocker-root
```

Your directory should now look like:

```
containers-from-scratch/
├── gocker-root/        # Alpine root filesystem
│   ├── bin/
│   ├── etc/
│   ├── home/
│   ├── lib/
│   ├── proc/
│   ├── root/
│   ├── usr/
│   └── ...
├── main.go
├── go.mod
└── ...
```

> **Why Alpine?** Docker's official Alpine image uses the same tarball. It's a real Linux userland in ~3MB — perfect for learning.

### Usage

```bash
# Run a command inside gocker
sudo ./gocker run <command> [args]

# Examples
sudo ./gocker run echo hello
sudo ./gocker run /bin/sh
sudo ./gocker run ps aux
```

### Follow Along

Each chapter builds on the previous one. Start from Chapter 1 and work through them in order. Every commit in the git history corresponds to a chapter, so you can check out any point in the journey.

## Why Build This?

If you work with containers daily but don't understand what happens below `docker run`, this project is for you.

After building Gocker, you'll understand:

- Why containers are fast (they're just processes, not VMs)
- Why containers need Linux (namespaces and cgroups are Linux kernel features)
- How Docker isolates processes from each other
- How Docker limits CPU and memory per container
- How container images provide a filesystem without shipping an entire OS
- Why `docker run alpine sh` feels like a separate machine but isn't

## Where Can You Use This Knowledge?

- **DevOps/Infrastructure** — debug container issues at the kernel level
- **Kubernetes** — understand what's actually running on your nodes
- **Security** — know exactly what container isolation provides (and what it doesn't)
- **Interviews** — "explain how containers work" is a common senior-level question
- **Building tools** — create custom runtimes, sandboxes, or isolation layers

## How to Improve Gocker Further

Once you complete the base chapters, here are ideas to take it further:

- **Seccomp filters** — restrict which syscalls the container can make
- **Capabilities** — drop unnecessary Linux capabilities for security
- **User namespaces** — run containers without root on the host
- **Container networking** — bridge networks, port mapping, container-to-container communication
- **Image pulling** — pull images from Docker Hub using the Registry API
- **Multi-container** — run multiple isolated containers simultaneously
- **Resource monitoring** — track CPU/memory usage per container in real time
- **Logging** — capture and store container stdout/stderr

## Project Structure

```
containers-from-scratch/
├── main.go          # Entry point and CLI
├── go.mod           # Go module definition
├── .gitignore       # Ignore binaries, rootfs, tarballs
├── LICENSE          # MIT License
└── README.md        # You're here
```

## Resources

These helped inspire and guide this project:

- [Linux Namespaces - man7.org](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Cgroups v2 - kernel.org](https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html)
- [Containers From Scratch - Liz Rice (talk)](https://www.youtube.com/watch?v=8fi7uSYlOdc)
- [Open Container Initiative Runtime Spec](https://github.com/opencontainers/runtime-spec)

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

Built for learning. Star ⭐ if this helped you understand containers.
