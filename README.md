# Clone Media (Clonar Mídia)

A professional low-level tool for bit-by-bit disk cloning, forensic data recovery, and file system repair. Developed as a cross-platform desktop application using the Wails framework.

## Overview

Clone Media is designed for system administrators, electronics technicians, and forensic analysts who require absolute precision in data duplication and recovery. It provides a high-level graphical interface for complex terminal-based utilities like `dd`, `testdisk`, and `photorec`.

## Core Features

- **Bit-by-Bit Cloning**: Direct sector-level duplication using optimized block sizes for HDD, SSD, and NVMe devices.
- **Forensic Recovery**: Deep file carving and partition recovery via integrated TestDisk and PhotoRec modules.
- **File System Repair**: Automated repair of corrupted partitions using FSCK and low-level system checks.
- **Hardware Monitoring**: S.M.A.R.T. health analysis and real-time disk status monitoring.
- **Multi-Language Support**: Fully localized interface supporting English and Portuguese.

## Architecture

The application is built with a modern decoupled architecture:

- **Frontend**: Svelte with Vite, utilizing a reactive design system with glassmorphism aesthetics.
- **Backend**: Go (Golang) via Wails, managing system-level calls, concurrency, and hardware events.
- **Core Logic**: Low-level hardware interaction implemented in C++ and Go modules.
- **Packaging**: Snapcraft (Strict Mode) for secure, sandboxed distribution on Linux systems.

## How it Works

### Disk Interaction
The application communicates with the Linux storage stack through the `udisks2` and `block-devices` interfaces. In sandboxed environments (like Snap), it requests authorized access to raw devices to ensure maximum performance and data integrity without compromising system security.

### Cloning Engine
The cloning process utilizes a high-concurrency buffer system in Go. It performs direct reads from the source device and writes to the destination, bypassing standard file system caches to ensure an identical copy of the partition table, bootloader, and data sectors.

### Privilege Elevation
Critical operations are protected by a privilege elevation layer, ensuring that destructive actions (like overwriting a disk) require explicit administrator authorization.

## Build Requirements

- Go 1.18+
- Node.js & NPM
- Wails CLI
- Build Essentials (gcc, g++, libgtk-3-dev, libwebkit2gtk-4.0-dev)

## Deployment

To build the Snap package in strict mode:

```bash
snapcraft
```

---

**Erasmo Cardoso**
Software Engineer | Electronics Technician
