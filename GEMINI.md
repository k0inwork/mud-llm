# Gemini File Editing Strategy

This document outlines the methods for editing files within this project.

## Summary of Methods

I can edit files using several methods, each with its own strengths:

1.  **Internal `replace` Tool:** The safest method for precise, context-aware text replacements in code.
2.  **CLI Tools (`sed`, `awk`, etc.):** Powerful for pattern-based, stream-of-text, and structured data manipulation. **Especially efficient for multiple, similar fixes across many files.** Executed via `run_shell_command`.
3.  **Internal `read_file` + `write_file` Combination:** The most robust method for complex, multi-step, or structural changes that are difficult to perform with pattern matching. This offers the most control over the final output.

## Operational Instruction

If one editing method fails (e.g., `replace` fails due to context matching issues, or `sed` fails due to complex escaping), I will automatically rotate to a different method to accomplish the task. **For multiple fixes of the same type, CLI tools like `sed` or `awk` are preferred for their efficiency.** The typical fallback order is:

`replace` -> CLI tools (`sed`) -> `read_file` + `write_file`

This ensures resilience and efficient task completion.
