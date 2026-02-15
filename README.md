# Examination

## Introduction
This project is a Go-based web application that serves as an examination platform. It strictly follows **Specification-Driven Design (SDD)** principles.

## AI Development Guidelines
**Examples for AI Assistants:**
> [!IMPORTANT]
> **Single Source of Truth**
> This project is driven by external specifications. The directory `sdd-examination-spec` is a symbolic link to the actual specification repository.
>
> 1. **Do not modify code based on assumptions.** Always check `sdd-examination-spec` first.
> 2. **Specification First.** All features, API endpoints, and data models must be defined in the specifications before implementation.
> 3. **Validation.** After implementing changes, verify that the code aligns perfectly with the OpenAPI/Markdown specs in `sdd-examination-spec`.
>
> **Note on Language:**
> *   **Project Documentation**: `README.md` and project-level docs must be in **English**. This serves as the primary entry point for AI context understanding.
> *   **Specifications**: The documents in `sdd-examination-spec` are primarily in **Korean**. AI assistants should read these for logic and requirements but implement code/comments in English.

## Project Constitution
> [!IMPORTANT]
> **Read the Law**: All contributors must read and strictly follow the [Project Constitution](CONSTITUTION.md).
>
> **Core Tenet**: The specification folder (`sdd-examination-spec`) is a symbolic link to an external repository. **IT MUST NEVER BE COMMITTED** to this repository. The specifications are the single source of truth.

## Specification Repository
To maintain a clear separation of concerns, the specifications for this project are hosted in a separate repository. This ensures that the specifications remain the single source of truth and are versioned independently of the implementation.

The specification repository is expected to be cloned alongside this repository and linked via a symbolic link.

## GitHub Credentials
Since this project relies on a separate specification repository (`sdd-examination-spec`), your local development environment requires a **Single Shared Token** that grants access to both repositories.

### Token Permissions
When creating a **Fine-grained Personal Access Token**, ensure you select **both** repositories and grant the following repository permissions:
-   **Contents**: Read and write
-   **Pull requests**: Read and write
-   **Workflows**: Read and write
-   **Metadata**: Read-only (Required)

### Local Credential Configuration
To ensure correct authentication context when working with multiple repositories, configure git to use the full HTTP path for credentials:

```bash
git config --local credential.useHttpPath true
```

## Getting Started

### 1. Clone the Repositories
It is recommended to keep both the implementation and specification repositories in the same parent directory.

```bash
# 1. Create a parent directory
mkdir dev
cd dev

# 2. Clone the implementation repository (this repo)
git clone https://github.com/your-org/examination.git

# 3. Clone the specification repository
# (Replace with the actual URL of your spec repo)
git clone https://github.com/your-org/examination-specs.git
```

### 2. Setup Symbolic Link
To connect the implementation with the specifications, create a symbolic link named `sdd-examination-spec` in the root of the `examination` repository.

```bash
cd examination
ln -s ../examination-specs sdd-examination-spec
```

> [!NOTE]
> **Custom Specification Paths**
> If you have cloned the specification repository to a different directory or gave it a different name (e.g., `my-custom-specs`), adjust the command accordingly:
> ```bash
> ln -s ../my-custom-specs sdd-examination-spec
> ```

> [!IMPORTANT]
> **Git Ignore (Crucial)**
> To prevent accidental commits of the specification folder (which would duplicate data and break the separation of concerns), you **MUST** add the symlink name to `.gitignore`.
>
> ```bash
> # Append the symlink name to .gitignore
> echo "sdd-examination-spec" >> .gitignore
>
> # If you used a custom name, add that as well
> # echo "my-custom-specs" >> .gitignore
> ```

### 3. Verification
Verify that the link is correctly established:

```bash
ls -l sdd-examination-spec
# Output should show it pointing to your spec repo
# e.g., sdd-examination-spec -> ../examination-specs
```


## Running the Application

Please refer to the **[Development Guide](DEVELOPMENT.md)** for instructions on how to set up your environment and run the application locally.

> **Note:** Detailed technical specifications and historical context can be found in `sdd-examination-spec/RUNNING.md`.


