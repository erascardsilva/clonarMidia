# Clone Media (Clonar Mídia)

A professional low-level tool for bit-by-bit disk cloning, forensic data recovery, and file system repair. Developed as a cross-platform desktop application using the Wails framework.

![App Screenshot](https://raw.githubusercontent.com/erascardsilva/clonarMidia/master/screenshot.png)

## Overview

**Clone Media** is designed for system administrators, electronics technicians, and forensic analysts who require absolute precision in data duplication and recovery. It provides a high-level graphical interface for complex terminal-based utilities like `dd`, `testdisk`, and `photorec`.

## Core Features

- **Bit-by-Bit Cloning**: Direct sector-level duplication using optimized block sizes for HDD, SSD, and NVMe devices.
- **Forensic Recovery**: Deep file carving and partition recovery via integrated TestDisk and PhotoRec modules.
- **File System Repair**: Automated repair of corrupted partitions using FSCK and low-level system checks.
- **Hardware Monitoring**: S.M.A.R.T. health analysis and real-time disk status monitoring.
- **Multi-Language Support**: Fully localized interface supporting English and Portuguese.

---

## Installation

### 1. Recommended: Full Version (.deb)
The full version provides unrestricted access to the system storage stack. This is the recommended version for professional cloning and data recovery.

**Download the latest .deb package from [Releases](https://github.com/erascardsilva/clonarMidia/releases)**

```bash
sudo dpkg -i clonarmidia_1.0.0_amd64.deb
sudo apt-get install -f  # To fix any missing dependencies
```

### 2. Snap Version (Demo / Sandboxed)
The Snap version runs in a secure sandbox. Due to Linux security policies (Strict Confinement), it may have limited access to raw hardware devices by default.

[![Get it from the Snap Store](https://snapcraft.io/en/dark/install.svg)](https://snapcraft.io/clonarmidia)

### 3. Build from Source
If you prefer to build the application yourself:

**Requirements:**
- Go 1.22+
- Node.js 20+ & NPM
- Wails CLI
- Build Essentials (`libgtk-3-dev`, `libwebkit2gtk-4.1-dev`, `pkg-config`)

**Build Command:**
```bash
# Clone the repository
git clone https://github.com/erascardsilva/clonarMidia.git
cd clonarMidia

# Build the application
wails build -m -tags webkit2_41 -o clonarmidia
```

---

## Architecture

The application is built with a modern decoupled architecture:
- **Frontend**: Svelte with Vite, utilizing a reactive design system with glassmorphism aesthetics.
- **Backend**: Go (Golang) via Wails, managing system-level calls, concurrency, and hardware events.
- **Core Logic**: Low-level hardware interaction implemented in optimized Go modules.

## Support & Donation

If you find this tool useful, consider supporting its development:

[![Donate with PayPal](https://img.shields.io/badge/Donate-PayPal-blue.svg)](https://www.paypal.com/ncp/payment/8V6WQCGN6HDCQ)

---

**Erasmo Cardoso**  
Software Engineer | Electronics Technician
