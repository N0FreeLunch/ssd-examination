# Agent Skills & Configuration

This project uses Antigravity agent skills to enhance the developer experience and provide project-specific guidance.

## Location of Skills

To maintain agent skills as core project assets, they are primarily stored within the specification repository:

- **Primary Storage**: `sdd-examination-spec/.agents/skills/`
- **Root Symbolic Links**: Individual skills in `.agents/skills/` point to their respective directories in the spec folder to ensure the agent can discover them.

## Available Skills

- **context-awareness**: **CRITICAL**: Use at the start of every session to discover and acknowledge the project structure (Main vs. Spec).
- **terminal-etiquette**: Rules and best practices for AI agents to execute terminal commands, including correct Cwd handling and command sequencing.
- **git-interop**: A set of git commands and best practices for managing multiple repositories within a project, specifically using the -C option.
- **spec-management**: Focused on the content, structure, and intent of the specification repository.
- **commit-policy**: Guidelines for language usage (English Prime vs Spec Korean) and committing code in this project.

Please refer to the `.agents/skills/` directory for the source of truth for all agent skills.

## Security & Path Policy

To ensure security and project portability, all agent documentation and internal links must follow these rules:
- **No Absolute Paths**: Avoid using absolute machine-specific paths (e.g., `/Users/name/...`).
- **Relative Paths**: Always use relative paths from the project root (e.g., `./sdd-examination-spec/`) for referencing files or directories.
- **Portability**: Ensure all scripts and documentation work regardless of where the project is cloned.
