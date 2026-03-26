# Project Constitution

This document defines the **absolute rules** that govern the development and maintenance of this project. All contributors (humans and AI) must adhere to these rules without exception.

## Article I. Separation of Concerns
1.  **External Specifications**: The requirements, domain logic, and API designs for this project are strictly defined in an **external repository**, accessible via the `./sdd-examination-spec` folder in the project root.
2.  **No Spec Commits**: You **MUST NEVER** commit the contents of the `sdd-examination-spec` folder to this repository. The specification repository must remain a separate entity to prevent data duplication and versioning conflicts.
3.  **Mandatory Skills**: AI assistants MUST use the following skills as needed:
    *   `context-awareness`: **CRITICAL**: Use immediately at the start of every session to discover and acknowledge the project structure (Main vs. Spec).
    *   `spec-management`: To navigate and understand specification documents.
    *   `git-interop`: When performing operations across multiple git repositories (specifically using `-C`).
    *   `terminal-etiquette`: To ensure correct terminal usage and directory context management.
    *   `commit-policy`: To follow the project's language policy (English Prime) and commit guidelines.
4.  **Git Ignore & Verification**: 
    *   The symbolic link `sdd-examination-spec` must always be present in `.gitignore`.
    *   **Pre-Commit Check**: Before running `git add .`, you MUST check `git status` to ensure `sdd-examination-spec` is NOT being tracked or staged. If it appears as `modified` or `new file`, you MUST untrack it (`git rm --cached sdd-examination-spec`) immediately.

## Article II. Specification-Driven Development (SDD)
1.  **Spec First**: Implementation code shall not be written until a corresponding specification exists.
2.  **Single Source of Truth**: When ambiguity arises in the code, the Specification is the final authority. Do not change the code to fit assumptions; verify the Spec first.
3.  **Verification**: Code changes must be verified against the rules defined in the Specification.

## Article III. Data Privacy & Security
1.  **No Database Commits**: Production or local development tokens, secrets, and **database files** (*.db, *.sqlite, etc.) must **NEVER** be committed to the repository.
2.  **Ignored Data**: The `data/` directory is reserved for local persistence and must remain in `.gitignore`.

## Article IV. AI Assistant Guidelines
1.  **Read-Only Specs**: AI assistants may read `sdd-examination-spec` to understand requirements but must **never** attempt to modify files within that directory unless explicitly instructed to update the *Specification Repository* itself.
2.  **Cross-Language Implementation**: While specifications may be written in Korean (or other languages), the implementation code (variables, general comments, commit messages) must be written in **English**. However, **Korean (or the specification language) may be used in comments** to explicitly map English code elements to the corresponding terms in the Specification, ensuring accurate reflection of requirements.
3.  **Sensitive Content Verification**: Before creating or editing files likely to contain secrets (e.g., config files, `.env`, test data), AI assistants MUST pause and explicitly ask the USER for review to ensure no actual secret keys are being hardcoded or committed.
4.  **Documentation Language**: The `README.md` and other project documentation in this repository must be maintained in **English** to serve as an efficient entry point for AI assistants. This allows AIs to quickly understand the project structure and navigate to the specifications (which may be in other languages) without context switching overhead.

5.  **Coding Standards**: All code generation and modification must adhere to the standards defined in `skills.md` located in the `sdd-examination-spec` directory.

## Article VI. Repository Privacy & Hierarchy
1.  **Repository Distinction**:
    - **Core Repository (Public)**: Contains source code, implementation, and minimum setup required to run the project. No sensitive or detailed operational documentation should be stored here.
    - **Specification Repository (Private)**: Contains all specifications, sensitive infrastructure details, troubleshooting guides, and operational guidelines. This repository must remain private.
2.  **Content Separation**:
    - Detailed operational manuals, deep troubleshooting logs, and sensitive infrastructure configurations must be kept within the Specification repository.
    - The Core repository should only include the "Minimum Viable Setup" (e.g., local running instructions) necessary for any contributor to start the application.

## Article VII. Amendments
1.  **User Confirmation Required**: Any changes to this **CONSTITUTION.md** file require explicit approval from the USER. AI assistants must not modify this file without a direct request or confirmation.
