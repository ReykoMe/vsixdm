# VSIX DOWNLOAD MANAGER

[![Version](https://img.shields.io/badge/version-0.0.1--alpha-blue.svg)]()

## Table of Contents
- [Description](#description)
- [Features](#features)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Downloading Extensions](#downloading-extensions)
  - [Installing Extensions](#installing-extensions)
- [Development](#development)
- [Known Issues](#known-issues)
- [Future Plans](#future-plans)
- [Contributing](#contributing)
- [License](#license)

## Description
VSIX Download Manager is a simple utility that allows you to download `.vsix` extensions from the Visual Studio Code Marketplace for use with VSCode-like editors (code-oss, etc.).

## Features
- Download VSIX extensions from VS Code Marketplace
- Support for extension packs (multiple VSIX files)
- Simple command-line interface
- Portable version support

## System Requirements
- Windows (Linux support planned for future)
- Go 1.24.1 or higher
- Visual Studio Code portable (for development)

## Installation
1. Download the latest release
2. Extract the archive to your desired location
3. Ensure you have the required dependencies installed

## Usage

### Downloading Extensions
1. Open PowerShell in the application folder
2. Find your desired extension on [VS Code Marketplace](https://marketplace.visualstudio.com)
3. Copy the extension's page link
4. Run the command:
   ```powershell
   .\app.exe --src https://marketplace.visualstudio.com/items?itemName=publisher.extension
   ```
5. Find your downloaded `.vsix` file in the `.\output` folder

### Installing Extensions

#### Method 1: Terminal (CODE-OSS example)
```powershell
codium --install-extension .\loaded-extension.vsix
```

#### Method 2: Through VS Code UI
1. Open VS Code
2. Press `F1` to open the command palette
3. Search for "Extensions: Install from VSIX"
4. Select your VSIX file

#### Method 3: Drag and Drop
1. Open VS Code
2. Go to Extensions tab (`Ctrl+Shift+X`)
3. Drag and drop the VSIX file into the sidebar

#### Method 4: Manual Installation
1. Rename `.vsix` to `.zip`
2. Extract to a folder named after the extension
3. Move contents from `extension` folder to the root
4. Copy to `%USERPROFILE%/.[editorName]/extensions/`
5. Delete `extensions.json`
6. Restart your editor

## Development

### Prerequisites
- Windows OS (Linux support coming soon)
- Go 1.24.1
- Visual Studio Code portable

### Setup
1. Clone the repository
2. Run initialization script:
   ```powershell
   .\run.ps1 init
   ```

### Available Scripts
- Development mode:
  ```powershell
  .\run.ps1 dev --src [MARKETPLACE_EXTENSION_LINK]
  ```
- Build:
  ```powershell
  .\run.ps1 build
  ```

### Additional Software
This app requires Visual Studio Code portable. Options:
1. Use the included portable app in the `soft` folder (default)
2. Download from [portapps.io](https://portapps.io/app/vscode-portable/)
3. Create your own portable version following [VS Code documentation](https://code.visualstudio.com/docs/editor/portable)

## Known Issues
- Some extensions may not work properly in different editors due to private VS Code APIs
- Limited to Windows platform (currently)
- No direct extension management after download

## Future Plans
1. Linux/MacOS support
2. Direct extension installation after download
3. Extension version selection
4. Improved error handling


## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is open source and available under the MIT License.

---
**Note**: This is not a commercial application. There is no strict release plan, and there is no guarantee that extensions will work properly in different editors due to their use of private VS Code APIs, editor-specific checks, or DLL dependencies.

