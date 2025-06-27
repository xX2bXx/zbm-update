# ZBM-Update

[![WTFPL](https://img.shields.io/badge/License-WTFPL-brightgreen.svg)](https://www.wtfpl.net/)

A simple yet reliable ZFSBootMenu updater

## Installation
```sh
git clone https://github.com/xX2bXx/zbm-update
cd zbm-update
go build -o zbm-update
```
You can also use this utility to install ZFSBootMenu for the very first time.  
In that case, it's often more convenient to just download the binary:
```sh
wget https://github.com/xX2bXx/zbm-update/releases/download/1.0.0/zbm-update
```

## Usage
(just an example)
```sh
sudo zbm-update \
  --target /boot/efi/EFI/ZBM/VMLINUZ.EFI \   # Path to the new ZFSBootMenu .EFI
  --backup /boot/efi/EFI/ZBM/VMLINUZ-BACKUP.EFI # Path to save the current one
```

## Philosophy
- Does one thing well
- No hand-holding
- Maximum freedom (WTFPL)
