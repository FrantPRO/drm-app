# Claude Code Rules

## Code Quality Rules

1. **After any changes, run `make lint` and fix errors if any**
   - Always run the linter after making code changes
   - Fix any linting errors before proceeding
   - Ensure code follows Go conventions and project standards

2. **After completing the task, run `make test` to make sure everything works correctly**
   - Run all tests to verify functionality
   - Ensure no regressions are introduced
   - Verify the application still works as expected

3. **After completing a task, add new files to git commit and write a short summary for the commit text**
   - Stage all new and modified files using `git add`
   - Create a concise commit message that describes what was accomplished
   - Focus on the "what" and "why" rather than implementation details
   - Keep commit messages under 50 characters for the summary line