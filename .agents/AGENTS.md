# Agent Skills & Configuration

This project uses Antigravity agent skills to enhance the developer experience and provide project-specific guidance.

## Location of Skills

To maintain agent skills as core project assets, they are primarily stored within the specification repository:

- **Primary Storage**: `sdd-examination-spec/.agents/skills/`
- **Root Symbolic Links**: Individual skills in `.agents/skills/` point to their respective directories in the spec folder to ensure the agent can discover them.

## Available Skills

- **spec-management**: Instructions on how to access and manage the specification repository. Accessible via symbolic link in this directory.

Please refer to the `sdd-examination-spec/.agents/skills/` directory for the source of truth for all agent skills.

## Security & Path Policy

To ensure security and project portability, all agent documentation and internal links must follow these rules:
- **No Absolute Paths**: Avoid using absolute machine-specific paths (e.g., `/Users/name/...`).
- **Relative Paths**: Always use relative paths from the project root (e.g., `./sdd-examination-spec/`) for referencing files or directories.
- **Portability**: Ensure all scripts and documentation work regardless of where the project is cloned.
