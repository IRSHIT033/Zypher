# Zypher

A lightweight version control system inspired by Git, implemented in Go.

```
███████╗██╗   ██╗██████╗ ██╗  ██╗███████╗██████╗
╚══███╔╝╚██╗ ██╔╝██╔══██╗██║  ██║██╔════╝██╔══██╗
  ███╔╝  ╚████╔╝ ██████╔╝███████║█████╗  ██████╔╝
 ███╔╝    ╚██╔╝  ██╔═══╝ ██╔══██║██╔══╝  ██╔══██╗
███████╗   ██║   ██║     ██║  ██║███████╗██║  ██║
╚══════╝   ╚═╝   ╚═╝     ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
```

## Features

- Initialize repositories
- Track file changes
- Create and manage commits
- Branch management
- Commit history viewing
- File status checking
- Commit reverting

## Installation

```bash
go install github.com/IRSHIT033/zypher
```

## Usage

### Initialize a Repository

```bash
zypher init
```

### Check Repository Status

```bash
zypher status
```

### Create a Commit

```bash
zypher commit "Your commit message"
```

### View Commit History

```bash
zypher log
```

### Revert to a Previous Commit

```bash
zypher revert <commit-hash>
```

### Branch Management

List all branches:

```bash
zypher branch
```

Create a new branch:

```bash
zypher branch <branch-name>
```

Switch to a branch:

```bash
zypher checkout <branch-name>
```

## Repository Structure

Zypher stores its data in a `.zypher` directory with the following structure:

```
.zypher/
├── objects/          # Stores file contents and commits
│   └── <hash-prefix>/
│       └── <hash-suffix>
├── refs/
│   └── heads/        # Branch references
│       ├── main      # Default branch
│       └── <branch>  # Other branches
└── HEAD             # Points to current branch
```

## How It Works

- **Objects**: Files and commits are stored as blobs in the objects directory
- **References**: Branches are stored in the refs/heads directory
- **HEAD**: Points to the current branch
- **Commits**: Store file hashes and metadata
- **Branches**: Point to specific commits

## License

MIT
