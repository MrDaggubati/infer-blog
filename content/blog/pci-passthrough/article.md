---
title: PCI GPU passthrough
slug: pci-gpu-passthrough
description: > 
    Preparing an Ubuntu host for PCI Passthrough, so a VM can claim the dGPU and performane of the VM is near native
date: 2025-05-13
author: Sudhakar Daggubati
summary: >
    A practical VFIO controller for preparing an Ubuntu host for PCI GPU passthrough, with preflight checks, GRUB recovery, NVIDIA isolation, and verification.
tags: 
    - linux
    - vfio
    - gpu
    - ubuntu
    - virtualization
    - kvm
featured: true
---

Little snippet to prepare the host for GPU passthrough


**Warning:** This script modifies GRUB, initramfs, kernel module configuration, udev rules, and NVIDIA services. Incorrect PCI IDs, filesystem assumptions, boot configuration, or hardware topology can leave the host without working graphics or unable to boot normally.

Test this on a system where you have console/recovery access before relying on it remotely.

## Before using it

Identify the GPU and its associated PCI functions:

```bash
lspci -nn
```

For an NVIDIA GPU, you may see something similar to:

01:00.0 VGA compatible controller [0300]: NVIDIA Corporation ... [10de:2820]
01:00.1 Audio device [0403]: NVIDIA Corporation ... [10de:22bd]

Update these values in the script for your hardware:

```bash
GPU_PCI="0000:01:00.0"
GPU_AUDIO_PCI="0000:01:00.1"
GPU_IDS="10de:2820,10de:22bd"
```

The values in this article are examples from one machine. Do not assume they match your GPU.

The VFIO controller

```bash
#!/usr/bin/env bash
set -euo pipefail

echo " This enables vfio passthrough on Ubuntu 24+ "

# =========================
# CONFIG
# =========================
GPU_PCI="0000:01:00.0"
GPU_AUDIO_PCI="0000:01:00.1"
GPU_IDS="10de:2820,10de:22bd"

GRUB="/etc/default/grub"
VFIO_CONF="/etc/modprobe.d/vfio.conf"
BLACKLIST="/etc/modprobe.d/blacklist-nvidia.conf"
RECOVERY_FLAG="/boot/vfio_last_state"
SAFE_ENTRY_FILE="/etc/grub.d/40_vfio_safe_mode"

ROOT_UUID="$(findmnt -no UUID /)"
ROOT_FSTYPE="$(findmnt -no FSTYPE /)"

# SAFETY CHECK: Ensure UUID was actually captured
if [[ -z "$ROOT_UUID" ]]; then
    echo "[!] ERROR: Could not determine ROOT_UUID."
    echo "    Check if findmnt is installed or if you have sufficient permissions."
    exit 1
fi

# Check for filesystem compatibility
if [[ "$ROOT_FSTYPE" != ext* ]]; then
    echo "===================================================="
    echo " [!] WARNING: NON-EXT FILESYSTEM DETECTED ($ROOT_FSTYPE)"
    echo "===================================================="
    echo " Your Safe Mode GRUB entry is currently hardcoded for 'insmod ext2'."
    echo " To ensure recovery works, you should update the 'insmod' line"
    echo " in the ensure_safe_mode_entry function to 'insmod $ROOT_FSTYPE'."
    echo ""
    read -p "Continue anyway? (y/n): " confirm
    [[ "$confirm" != "y" ]] && exit 1
fi

# =========================
# BANNER
# =========================
show_recovery_banner() {
    echo ""
    echo "===================================================="
    echo " VFIO SAFE CONTROLLER"
    echo "===================================================="
    echo ""
    echo "RECOVERY OPTIONS:"
    echo "  → GRUB: 'Ubuntu (SAFE MODE - iGPU only)'"
    echo "  → OR 'Advanced options'"
    echo ""
    echo "SAFE MODE = NO VFIO, NO NVIDIA, iGPU ONLY"
    echo "===================================================="
    echo ""
}

# =========================
# VM SAFETY CHECK
# =========================
check_vms() {
    echo "[*] Checking running VMs..."

    RUNNING=$(virsh list --state-running 2>/dev/null | awk 'NR>2 {print $2}')

    if [[ -z "${RUNNING// }" ]]; then
        RUNNING=$(virsh list 2>/dev/null | awk 'NR>2 && $3=="running" {print $2}')
    fi

    if [[ -n "$RUNNING" ]]; then
        echo "[!] Running VMs detected:"
        echo "$RUNNING"
        return 1
    fi

    echo "[✓] No running VMs"
    return 0
}

# =========================
# PREFLIGHT
# =========================
preflight() {
    echo "[*] Preflight check..."

    local blocked=0
    local running_vms

    running_vms=$(virsh list --name 2>/dev/null | sed '/^$/d')

    if [[ -z "$running_vms" ]]; then
        echo "[✓] No running VMs"
        echo "[✓] Preflight OK"
        return 0
    fi

    for vm in $running_vms; do
        if virsh dumpxml "$vm" 2>/dev/null | grep -qE "$GPU_PCI|$GPU_AUDIO_PCI"; then
            echo "[!] Running VM '$vm' is using configured GPU passthrough devices"
            blocked=1
        else
            echo "[i] Running VM '$vm' does not use passthrough GPU"
        fi
    done

    if [[ $blocked -eq 1 ]]; then
        echo "[!] Stop affected VM(s) before changing VFIO state"
        return 1
    fi

    echo "[✓] No running VM currently uses the dGPU"
    echo "[✓] Preflight OK"
    return 0
}

# =========================
# STATUS
# =========================
status() {
    echo "===================="
    echo " SYSTEM STATUS"
    echo "===================="

    virsh list 2>/dev/null || true
    echo ""

    lspci -nnk -s ${GPU_PCI:5:7} || true
}

# =========================
# SAFE MODE GRUB ENTRY
# =========================
ensure_safe_mode_entry() {
    echo "[*] Ensuring SAFE MODE GRUB entry for $ROOT_FSTYPE..."

    local grub_mod="$ROOT_FSTYPE"
    [[ "$ROOT_FSTYPE" == ext* ]] && grub_mod="ext2"

    sudo tee "$SAFE_ENTRY_FILE" >/dev/null <<EOF
#!/bin/sh
exec tail -n +3 \$0

menuentry "Ubuntu (SAFE MODE - iGPU only, NO VFIO)" {
    insmod part_gpt
    insmod $grub_mod
    search --no-floppy --fs-uuid --set=root $ROOT_UUID
    linux /boot/vmlinuz root=UUID=$ROOT_UUID ro quiet splash intel_iommu=off modprobe.blacklist=vfio_pci,vfio,vfio_iommu_type1,nvidia,nvidia_drm,nvidia_modeset
    initrd /boot/initrd.img
}
EOF

    sudo chmod +x "$SAFE_ENTRY_FILE"
}

# =========================
# GRUB APPLY VFIO
# =========================
apply_vfio_grub() {
    sudo cp "$GRUB" "$GRUB.bak.$(date +%s)"

    sudo sed -i 's/vfio-pci.ids=[^ ]*//g' "$GRUB"
    sudo sed -i 's/intel_iommu=on//g' "$GRUB"
    sudo sed -i 's/iommu=pt//g' "$GRUB"

    sudo sed -i \
        's/^GRUB_CMDLINE_LINUX_DEFAULT=.*/GRUB_CMDLINE_LINUX_DEFAULT="quiet splash intel_iommu=on iommu=pt vfio-pci.ids='"$GPU_IDS"' rd.driver.pre=vfio-pci"/' \
        "$GRUB"
}

verify_vfio_binding() {
    echo "[*] Checking GPU driver binding..."

    if lspci -nnk -s ${GPU_PCI:5:7} | grep -q "vfio-pci"; then
        echo "[✓] GPU correctly bound to VFIO"
        return 0
    else
        echo "[!] VFIO FAILED — GPU still not isolated"
        echo "    Do NOT launch VM"
        return 1
    fi
}

# =========================
# GRUB RESTORE (iGPU)
# =========================
restore_igpu_grub() {
    sudo cp "$GRUB" "$GRUB.bak.$(date +%s)"

    sudo sed -i 's/vfio-pci.ids=[^ ]*//g' "$GRUB"
    sudo sed -i 's/intel_iommu=on//g' "$GRUB"
    sudo sed -i 's/iommu=pt//g' "$GRUB"

    sudo sed -i \
        's/^GRUB_CMDLINE_LINUX_DEFAULT=.*/GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"/' \
        "$GRUB"
}

# =========================
# ENABLE VFIO
# =========================
enable_flow() {
    show_recovery_banner
    preflight || exit 1

    echo "[*] Enabling VFIO..."

    cat <<EOF | sudo tee "$BLACKLIST" >/dev/null
blacklist nvidia
blacklist nvidia_drm
blacklist nvidia_modeset
blacklist nvidia_uvm
blacklist nvidia_nouveau
EOF

    printf '%s\n' \
        "options vfio-pci ids=$GPU_IDS disable_vga=1" \
        "softdep nvidia pre: vfio-pci" \
        "softdep nvidia_drm pre: vfio-pci" \
        "softdep nvidia_modeset pre: vfio-pci" \
        "softdep nvidia_uvm pre: vfio-pci" \
        | sudo tee "$VFIO_CONF" >/dev/null

    # Block udev from loading NVIDIA modules
    sudo tee /etc/udev/rules.d/71-nvidia.rules >/dev/null <<'EOF'
# VFIO override — block nvidia module loading via udev
ACTION=="add", DEVPATH=="/bus/pci/drivers/nvidia", RUN+="/bin/false"
EOF

    sudo systemctl disable --now nvidia-persistenced.service 2>/dev/null || true
    sudo systemctl disable --now nvidia-powerd.service 2>/dev/null || true
    sudo systemctl disable --now nvidia-cdi-refresh.path 2>/dev/null || true

    apply_vfio_grub
    ensure_safe_mode_entry

    echo "$GPU_IDS" | sudo tee "$RECOVERY_FLAG" >/dev/null

    sudo update-initramfs -u
    sudo update-grub

    echo "[*] Verifying GRUB configuration..."

    grep -q "vfio-pci.ids=$GPU_IDS" "$GRUB" || {
        echo "[!] GRUB injection failed"
        exit 1
    }

    echo "[✓] VFIO enabled safely"
}

# =========================
# DISABLE VFIO
# =========================
disable_flow() {
    show_recovery_banner

    sudo rm -f "$VFIO_CONF"
    sudo rm -f "$BLACKLIST"
    sudo rm -f /etc/udev/rules.d/71-nvidia.rules
    sudo udevadm control --reload-rules

    restore_igpu_grub
    ensure_safe_mode_entry

    sudo rm -f "$RECOVERY_FLAG"

    # Re-enable NVIDIA services
    sudo systemctl enable --now nvidia-persistenced.service 2>/dev/null || true
    sudo systemctl enable --now nvidia-powerd.service 2>/dev/null || true
    sudo systemctl enable --now nvidia-cdi-refresh.path 2>/dev/null || true

    sudo update-initramfs -u
    sudo update-grub

    echo "[✓] VFIO disabled"
}

# =========================
# RECOVERY MODE
# =========================
recover_flow() {
    echo "[!] RECOVERY MODE ACTIVATED"

    restore_igpu_grub

    sudo rm -f "$VFIO_CONF"
    sudo rm -f "$BLACKLIST"
    sudo rm -f "$RECOVERY_FLAG"

    sudo rm -f /etc/udev/rules.d/71-nvidia.rules
    sudo udevadm control --reload-rules

    ensure_safe_mode_entry

    # Re-enable NVIDIA services
    sudo systemctl enable --now nvidia-persistenced.service 2>/dev/null || true
    sudo systemctl enable --now nvidia-powerd.service 2>/dev/null || true
    sudo systemctl enable --now nvidia-cdi-refresh.path 2>/dev/null || true

    sudo update-initramfs -u
    sudo update-grub

    echo "[✓] System restored to iGPU SAFE MODE"
}

# =========================
# MAIN CONTROLLER
# =========================
main() {
    show_recovery_banner

    case "${1:-}" in
        status)
            status
            ;;
        preflight)
            preflight
            ;;
        enable)
            enable_flow
            ;;
        disable)
            disable_flow
            ;;
        recover)
            recover_flow
            ;;
        verify)
            verify_vfio_binding || exit 1
            ;;
        *)
            echo "Usage: $0 {status|preflight|enable|disable|recover|verify}"
            ;;
    esac
}

main "$@"

```
Usage

Make the controller executable:

```bash 
chmod +x vfio-controller.sh
```

Check the current state:

```bash
./vfio-controller.sh status
```

Run the preflight checks:

```bash
./vfio-controller.sh preflight
```

Enable VFIO:

```bash
./vfio-controller.sh enable
```

### A reboot is required before the new kernel command line and initramfs configuration take effect.

After rebooting, verify that the GPU is using vfio-pci:

```bash
./vfio-controller.sh verify
```

You can also inspect it directly:

```bash
lspci -nnk -s 01:00.0
```
The important line should indicate:

```bash
Kernel driver in use: vfio-pci
Returning the GPU to the host
```

To remove the VFIO configuration:

```bash
./vfio-controller.sh disable
```

Then reboot.

For the recovery path:

```bash
./vfio-controller.sh recover
```

The script also creates a GRUB entry named:

```bash
Ubuntu (SAFE MODE - iGPU only, NO VFIO)
```

This is intended to provide a boot path without VFIO or NVIDIA enabled.

What the script changes

When VFIO is enabled, the controller modifies or creates:

```bash
/etc/default/grub
/etc/modprobe.d/vfio.conf
/etc/modprobe.d/blacklist-nvidia.conf
/etc/udev/rules.d/71-nvidia.rules
/etc/grub.d/40_vfio_safe_mode
/boot/vfio_last_state
```

It then regenerates:

```bash
update-initramfs -u
update-grub
```

### The script creates timestamped backups of /etc/default/grub before modifying it.

## Notes

This example assumes an Intel IOMMU configuration:

```bash
intel_iommu=on iommu=pt
```

The hardware-specific values at the beginning of the script must be changed for the target host.

The recovery entry also assumes that the host can operate using its integrated GPU while the discrete NVIDIA GPU is unavailable.

Before using this on an important system, verify the machine's:

GPU PCI addresses
GPU and audio-function device IDs
IOMMU support
IOMMU groups
boot filesystem
GRUB configuration
host display configuration
recovery/console access

GPU passthrough is inherently hardware and firmware-dependent, so this controller should be treated as a starting point for a known host configuration rather than a universal installer.


