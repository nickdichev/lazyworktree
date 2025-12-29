# Migration Plan: Python to Go Feature Parity

This document tracks missing features in the Go implementation compared to the Python implementation. Items are organized by priority and complexity.

## Current Status

The Go port has a **solid architectural foundation** with proper separation of concerns, correct concurrency patterns, and a well-structured codebase. However, several user-facing features are stubbed or missing.

**Architecture Status:** ✅ Complete
**Core TUI:** ✅ Complete
**Git Operations:** ⚠️ Partial (read-only operations work, mutations stubbed)
**Feature Parity:** ✅ ~70% complete (+10% from Priority 2 features)

---

## Current Implementation Session (2025-12-29)

**Focus:** Priority 2 (Enhanced UX) - Independent Features

**Completed:**
- [x] 2.5 Debounced Detail View Updates (Low complexity) ✅
- [x] 2.4 Delta Integration (Low complexity) ✅
- [x] 2.3 Enhanced Diff View - Three-part diff (Medium complexity) ✅

**Deferred (require Priority 1 infrastructure):**
- 2.2 Absorb Worktree (needs `.wt` files + TOFU integration)
- 2.1 Command Palette (complex, nice-to-have)

**Implementation Summary:**
1. ✅ Debouncing: 200ms delay prevents excessive git calls during rapid navigation
2. ✅ Delta: Automatic detection and application with silent fallback
3. ✅ Enhanced Diff: Three sections (staged/unstaged/untracked) with config limits
4. ✅ Auto-diff: Automatically shows diff in status pane when worktree is dirty
5. ✅ Improved status formatting: Color-coded status indicators (M/A/D/??)
6. ✅ Viewport navigation: Full vim-style navigation (j/k, Ctrl+D/U, g/G, PageUp/Down)
7. ✅ Adaptive layout: Right pane expands to 70% when focused for better diff/log viewing

**Files Modified:**
- `internal/app/app.go` - Added debouncing, delta usage, three-part diff integration, auto-diff display, improved status formatting, viewport navigation
- `internal/app/screens.go` - Updated help text with new navigation keys
- `internal/git/service.go` - Added delta detection/application, BuildThreePartDiff()

---

## Priority 1: Critical User-Facing Features (MUST HAVE)

### 1.1 Create Worktree Command
**Status:** Stubbed at `internal/app/app.go:628-631`
**Python Reference:** `lazyworktree/app.py:1057-1143`
**Complexity:** High

**Requirements:**
- Two-stage input: name input → base branch selection
- Input screens with defaults and validation
- Execute `git worktree add -b <name> <path> <base>`
- Run init commands from global config
- Run init commands from `.wt` file with TOFU security
- Set environment variables:
  - `WORKTREE_BRANCH`: branch name
  - `MAIN_WORKTREE_PATH`: main repo path
  - `WORKTREE_PATH`: new worktree path
  - `WORKTREE_NAME`: directory name
- Handle errors gracefully with notifications
- Refresh worktree list after creation

**Dependencies:**
- InputScreen with callback support
- TOFU security integration
- Command execution with environment variables
- Repository command loading (`.wt` file)

**Files to Modify:**
- `internal/app/app.go`: Implement `showCreateWorktree()` at line 628-631
- `internal/app/screens.go`: Enhance InputScreen if needed (already exists at lines 139-202)

---

### 1.2 Delete Worktree Command
**Status:** Partially stubbed at `internal/app/app.go:633-643`
**Python Reference:** `lazyworktree/app.py:1390-1419`, `app.py:1346-1388`
**Complexity:** High

**Requirements:**
- Show confirmation dialog with path and branch info
- Run terminate commands from global config
- Run terminate commands from `.wt` file with TOFU security
- Set environment variables (same as create)
- Execute `git worktree remove --force <path>`
- Execute `git branch -D <branch>`
- Handle partial failures (e.g., worktree removed but branch deletion fails)
- Refresh worktree list after deletion

**Dependencies:**
- ConfirmScreen (already exists at `internal/app/screens.go:72-137`)
- TOFU security integration
- Command execution with environment variables
- Repository command loading (`.wt` file)

**Files to Modify:**
- `internal/app/app.go`: Implement full delete workflow
- Add helper function for delete routine similar to Python's `_delete_worktree_routine()`

---

### 1.3 Rename Worktree Command
**Status:** Stubbed at `internal/app/app.go:661-664`
**Python Reference:** `lazyworktree/app.py:1295-1344`
**Complexity:** Medium

**Requirements:**
- Check that selected worktree is not main
- Show input screen with current branch name as default
- Validate new name is different from old
- Check destination path doesn't already exist
- Call `git.RenameWorktree()` (already implemented at `internal/git/service.go:444-457`)
- Refresh worktree list after rename

**Dependencies:**
- InputScreen (already exists)
- Git service method (already implemented)

**Files to Modify:**
- `internal/app/app.go`: Implement `showRenameWorktree()` at line 661-664

---

### 1.4 Prune Merged Worktrees Command
**Status:** Stubbed at `internal/app/app.go:666-669`
**Python Reference:** `lazyworktree/app.py:1421-1453`
**Complexity:** Medium

**Requirements:**
- Find all worktrees with `PR.State == "MERGED"` and not main
- Show confirmation screen with list of worktrees to delete
- Truncate list display if more than 10 (show "...and N more")
- Batch delete each worktree using delete routine
- Show notification with count of successfully deleted worktrees
- Refresh worktree list after completion

**Dependencies:**
- ConfirmScreen (already exists)
- Delete worktree routine (from 1.2)

**Files to Modify:**
- `internal/app/app.go`: Implement `showPruneMerged()` at line 666-669

---

## Priority 2: Enhanced User Experience (SHOULD HAVE)

### 2.1 Command Palette
**Status:** Not implemented
**Python Reference:** `lazyworktree/app.py:48-94`
**Complexity:** High

**Requirements:**
- Fuzzy searchable command interface
- Triggered by `Ctrl+/` key
- Lists all available actions with descriptions
- Executes selected action
- Uses Textual's `CommandPalette` equivalent (may need custom implementation)

**Bubble Tea Considerations:**
- Bubble Tea doesn't have built-in command palette
- Need to implement custom fuzzy matching or use library like `github.com/sahilm/fuzzy`
- Modal screen with list selection and filtering

**Files to Create/Modify:**
- `internal/app/commandpalette.go`: New file for command palette logic
- `internal/app/app.go`: Add keybinding and integrate palette

---

### 2.2 Absorb Worktree Command
**Status:** Not implemented
**Python Reference:** `lazyworktree/app.py:1455-1534`
**Complexity:** High

**Requirements:**
- Check selected worktree is not main
- Show confirmation dialog
- Run terminate commands with TOFU
- Checkout main branch in worktree: `git checkout main`
- Merge current branch into main: `git merge --no-edit <branch>`
- Remove worktree: `git worktree remove --force <path>`
- Delete branch: `git branch -D <branch>`
- Handle merge conflicts gracefully
- Refresh worktree list

**Dependencies:**
- ConfirmScreen (already exists)
- TOFU security integration
- Delete worktree routine (from 1.2)

**Files to Modify:**
- `internal/app/app.go`: Add new action method `showAbsorbWorktree()`
- Update key bindings to include absorb command

---

### 2.3 Diff View Enhancements
**Status:** ✅ COMPLETE (Session: 2025-12-29)
**Python Reference:** `lazyworktree/app.py:1165-1225`, `app.py:1271-1293`
**Complexity:** Medium

**Implementation:**
- ✅ Three-part diff:
  1. Staged changes: `git diff --cached --patch`
  2. Unstaged changes: `git diff --patch`
  3. Untracked files: `git diff --no-index /dev/null <file>` for each
- ✅ Configurable untracked file limit (`max_untracked_diffs` from config)
- ✅ Configurable diff truncation (`max_diff_chars` from config)
- ✅ Delta integration for syntax highlighting
- ✅ Truncation markers and file count notices
- ✅ **Auto-display**: Diff automatically shown in status pane when worktree is dirty
- ✅ **Improved status formatting**: Color-coded indicators (M=orange, A=green, D=red, ??=gray)

**Files Modified:**
- `internal/app/app.go`:
  - Enhanced `showDiff()` to use BuildThreePartDiff
  - Modified `updateDetailsView()` to auto-show diff when dirty
  - Improved `buildStatusContent()` with color-coded formatting
- `internal/git/service.go`: Added `BuildThreePartDiff()` and `getUntrackedFiles()` methods

---

### 2.4 Delta Integration
**Status:** ✅ COMPLETE (Session: 2025-12-29)
**Python Reference:** `lazyworktree/app.py:1145-1163`
**Complexity:** Low

**Implementation:**
- ✅ Delta detection on startup via `delta --version`
- ✅ Pipe diff output through `delta --no-gitconfig --paging=never`
- ✅ Silent fallback to plain diff if delta not available or errors
- ✅ Applied to diff view (press `d` key)

**Files Modified:**
- `internal/app/app.go`: Apply delta in `showDiff()`
- `internal/git/service.go`: Added `detectDelta()` and `ApplyDelta()` methods

---

### 2.5 Debounced Detail View Updates
**Status:** ✅ COMPLETE (Session: 2025-12-29)
**Python Reference:** `lazyworktree/app.py:711-713`
**Complexity:** Low

**Implementation:**
- ✅ 200ms debounce delay using `time.Sleep()` in tea.Cmd
- ✅ Prevents excessive git calls during rapid j/k navigation
- ✅ Ensures final detail update always happens for selected worktree
- ✅ Applied to all cursor movement: j/k keys and table navigation

**Files Modified:**
- `internal/app/app.go`:
  - Added `debouncedDetailsMsg` message type
  - Added `detailUpdateCancel` and `pendingDetailsIndex` fields to AppModel
  - Added `debouncedUpdateDetailsView()` method
  - Updated cursor movement handlers (j/k keys, table input)
  - Added message handler for `debouncedDetailsMsg`

---

## Priority 3: Advanced Features (NICE TO HAVE)

### 3.1 Special Init Command: `link_topsymlinks`
**Status:** Not implemented
**Python Reference:** `lazyworktree/app.py:964-1011`
**Complexity:** Medium

**Requirements:**
- Built-in command that runs as part of init_commands
- Symlinks untracked/ignored files from main to new worktree
- Symlinks editor configs: `.vscode`, `.idea`, `.cursor`, `.claude`
- Creates `tmp/` directory
- Runs `direnv allow` if `.envrc` exists
- Configurable via `.wt` file

**Files to Create/Modify:**
- `internal/commands/symlinks.go`: New package for special commands
- `internal/app/app.go`: Integrate special command detection

---

### 3.2 Repository-Specific Configuration (.wt files)
**Status:** Security implemented, execution not integrated
**Python Reference:** `lazyworktree/app.py:214-256`
**Complexity:** High

**Current Status:**
- TrustManager exists at `internal/security/trust.go`
- TOFU workflow implemented
- Loading `.wt` file not implemented
- Execution not integrated

**Requirements:**
- Load `.wt` file from main repository root
- Parse YAML with `init_commands` and `terminate_commands`
- Integrate with TOFU workflow
- Merge with global config commands
- Execute with environment variables set

**Files to Modify:**
- `internal/config/config.go`: Add `.wt` file loading
- `internal/app/app.go`: Integrate into create/delete workflows

---

### 3.3 Welcome Screen Workflow
**Status:** Screen exists but not integrated
**Python Reference:** `lazyworktree/app.py:573-620`
**Complexity:** Low

**Current Status:**
- WelcomeScreen exists at `internal/app/screens.go:367-446`
- Not shown when no worktrees found

**Requirements:**
- Show welcome screen when worktree list is empty
- Display current directory
- Display configured worktree root
- Offer retry button after config adjustment

**Files to Modify:**
- `internal/app/app.go`: Add welcome screen trigger logic

---

### 3.4 Commit Detail Viewer
**Status:** Basic viewer exists but not integrated
**Python Reference:** `lazyworktree/app.py:1235-1269`, `app.py:1572-1609`
**Complexity:** Medium

**Current Status:**
- CommitScreen exists at `internal/app/screens.go:448-513`
- Not triggered when selecting commit in log pane

**Requirements:**
- Press Enter in log pane to open commit detail
- Show commit metadata: SHA, author, date, message
- Show commit diff with syntax highlighting
- Scrollable content with vim-style navigation
- Header collapses on scroll (optional enhancement)

**Files to Modify:**
- `internal/app/app.go`: Add commit selection handling at line 302-305

---

### 3.5 Full-Screen Diff Viewer
**Status:** Screen exists but not used
**Python Reference:** `lazyworktree/screens.py:171-250`
**Complexity:** Low

**Current Status:**
- DiffScreen not implemented in Go
- Diff shown inline in status pane only

**Optional Enhancement:**
- Full-screen diff modal triggered by separate key (e.g., `Shift+D`)
- Vim-style navigation (j/k, Ctrl+d/u, g/G)
- Uses same diff building logic as inline view

**Files to Create/Modify:**
- `internal/app/screens.go`: Add DiffScreen implementation
- `internal/app/app.go`: Add keybinding and integration

---

## Priority 4: Testing & Quality (RECOMMENDED)

### 4.1 Unit Tests
**Status:** No test files exist
**Python Reference:** `tests/` directory with comprehensive tests
**Complexity:** High (ongoing)

**Recommended Coverage:**
- Config loading and validation
- Git service operations (with mocks)
- Worktree filtering and sorting
- Trust manager operations
- Screen state transitions

**Files to Create:**
- `internal/config/config_test.go`
- `internal/git/service_test.go`
- `internal/security/trust_test.go`
- `internal/app/app_test.go`

---

### 4.2 Integration Tests
**Status:** No integration tests
**Python Reference:** `tests/conftest.py` with FakeRepo fixture
**Complexity:** Very High

**Recommended Approach:**
- Create test fixture for temporary git repos
- Test full workflows (create → rename → delete)
- Test TOFU security prompts
- Test error handling and recovery

**Files to Create:**
- `test/integration/worktree_test.go`
- `test/fixtures/gitrepo.go`

---

## Implementation Roadmap

### Phase 1: Core Mutations (Weeks 1-2)
1. Implement `.wt` file loading and TOFU integration
2. Implement Create Worktree (1.1)
3. Implement Delete Worktree (1.2)
4. Implement Rename Worktree (1.3)

### Phase 2: Advanced Operations (Weeks 3-4)
5. Implement Prune Merged (1.4)
6. Implement Absorb Worktree (2.2)
7. Enhance Diff View (2.3)
8. Add Delta Integration (2.4)

### Phase 3: UX Polish (Week 5)
9. Add Command Palette (2.1)
10. Add Debounced Updates (2.5)
11. Integrate Commit Detail Viewer (3.4)
12. Integrate Welcome Screen (3.3)

### Phase 4: Advanced Features (Week 6)
13. Implement `link_topsymlinks` (3.1)
14. Add Full-Screen Diff Viewer (3.5)

### Phase 5: Quality & Hardening (Ongoing)
15. Add unit tests (4.1)
16. Add integration tests (4.2)
17. Performance optimization
18. Documentation updates

---

## Architecture Differences to Consider

### Python → Go Translation Patterns

| Python Pattern | Go Equivalent | Notes |
|----------------|---------------|-------|
| `async/await` | goroutines + channels | Use tea.Cmd pattern |
| `@dataclass` | struct | Already done correctly |
| `push_screen(callback)` | Screen state + channels | Need callback mechanism |
| `@work(exclusive=True)` | tea.Cmd with cancellation | Context support needed |
| List comprehensions | for loops | More verbose but clear |
| `Optional[T]` | `*T` or separate bool | Already done correctly |
| Exception handling | Error returns | Already done correctly |

### Key Challenges

1. **Modal Screen Callbacks**: Python's `push_screen(screen, callback)` pattern needs adaptation to Go's message-passing model
2. **Async Screen Dismissal**: Python uses futures; Go needs channels or tea.Cmd messages
3. **Environment Variable Expansion**: Python's `os.path.expanduser()` → Go's `os.ExpandEnv()` or `filepath.Join(os.Getenv("HOME"), ...)`
4. **YAML Parsing**: Python's type coercion is more forgiving; Go needs explicit handling
5. **Command Execution**: Python's `asyncio.create_subprocess_exec` → Go's `exec.CommandContext` (already done)

---

## File-by-File Migration Checklist

### `internal/app/app.go`
- [ ] Implement `showCreateWorktree()` (line 628-631)
- [ ] Complete `showDeleteWorktree()` (line 633-643)
- [ ] Implement `showRenameWorktree()` (line 661-664)
- [ ] Implement `showPruneMerged()` (line 666-669)
- [ ] Add `showAbsorbWorktree()` method
- [ ] Enhance `showDiff()` with three-part diff (line 646-659)
- [ ] Add debounce logic for detail updates
- [ ] Integrate commit detail viewer (line 302-305)
- [ ] Add command palette keybinding
- [ ] Add welcome screen trigger

### `internal/app/screens.go`
- [ ] Enhance InputScreen with validation callback support
- [ ] Add CommandPaletteScreen
- [ ] Add DiffScreen (full-screen viewer)
- [ ] Integrate CommitScreen (line 448-513)

### `internal/config/config.go`
- [ ] Add `.wt` file loading
- [ ] Add repository command merging
- [ ] Add environment variable expansion utilities

### `internal/git/service.go`
- [ ] Add `BuildThreePartDiff()` method
- [ ] Add `ApplyDelta()` method
- [ ] Add `ExecuteRepoCommands()` method with environment

### `internal/security/trust.go`
- [ ] Integrate TOFU workflow into app
- [ ] Add trust screen trigger logic

### New Files to Create
- [ ] `internal/commands/symlinks.go` - Special commands
- [ ] `internal/app/commandpalette.go` - Command palette
- [ ] `internal/app/helpers.go` - Shared helper functions
- [ ] Test files (see section 4)

---

## Risk Assessment

### High Risk Items
1. **TOFU Security Integration**: Critical for safe `.wt` execution; must not introduce vulnerabilities
2. **Partial Operation Failures**: Delete/absorb workflows have multiple steps; need rollback/cleanup logic
3. **Merge Conflicts in Absorb**: Must handle gracefully without data loss

### Medium Risk Items
1. **Environment Variable Handling**: Must match Python behavior exactly
2. **Screen Callback Pattern**: Core UX depends on this working smoothly
3. **Delta Integration**: Optional but users expect it; must degrade gracefully

### Low Risk Items
1. Command palette (nice-to-have)
2. Debouncing (minor UX improvement)
3. Full-screen diff viewer (optional alternative)

---

## Success Criteria

The Go implementation will achieve feature parity when:

1. ✅ All Priority 1 features are implemented and tested
2. ✅ All Priority 2 features are implemented (except command palette)
3. ✅ `.wt` file execution works with TOFU security
4. ✅ No data loss or corruption in any operation
5. ✅ Error messages match Python version quality
6. ✅ Performance is equal or better than Python version
7. ✅ At least 50% test coverage on critical paths

---

## Notes for Implementers

### Development Guidelines
1. **Read Python implementation first**: Understand the full workflow before coding
2. **Test incrementally**: Add tests as you implement each feature
3. **Preserve user safety**: Never compromise on TOFU security or data validation
4. **Match UX exactly**: Users expect consistent behavior across implementations
5. **Use Go idioms**: Don't try to write Python in Go; use channels, goroutines, error returns

### Common Pitfalls
- ❌ Don't skip TOFU integration - security is critical
- ❌ Don't forget environment variables in command execution
- ❌ Don't ignore partial failure scenarios
- ❌ Don't skip validation (path existence, name conflicts, etc.)
- ❌ Don't forget to refresh worktree list after mutations

### Quick Wins
- ✅ Rename worktree (1.3) - easiest to implement
- ✅ Prune merged (1.4) - simple once delete works
- ✅ Delta integration (2.4) - small, high-value feature
- ✅ Debouncing (2.5) - tiny change, big UX improvement

---

## Appendix: Feature Comparison Matrix

| Feature | Python | Go | Status | Priority |
|---------|--------|-----|--------|----------|
| Worktree List | ✅ | ✅ | Complete | - |
| Sorting | ✅ | ✅ | Complete | - |
| Filtering | ✅ | ✅ | Complete | - |
| PR Integration | ✅ | ✅ | Complete | - |
| Status View | ✅ | ✅ | Complete | - |
| Log View | ✅ | ✅ | Complete | - |
| Create Worktree | ✅ | ❌ | Stubbed | P1 |
| Delete Worktree | ✅ | ❌ | Stubbed | P1 |
| Rename Worktree | ✅ | ❌ | Stubbed | P1 |
| Prune Merged | ✅ | ❌ | Stubbed | P1 |
| Absorb Worktree | ✅ | ❌ | Missing | P2 |
| Diff View (Basic) | ✅ | ✅ | Complete | - |
| Diff View (Full) | ✅ | ❌ | Partial | P2 |
| Delta Integration | ✅ | ❌ | Missing | P2 |
| Command Palette | ✅ | ❌ | Missing | P2 |
| Commit Details | ✅ | ⚠️ | Not integrated | P3 |
| Welcome Screen | ✅ | ⚠️ | Not integrated | P3 |
| .wt Execution | ✅ | ❌ | Missing | P1 |
| TOFU Security | ✅ | ⚠️ | Not integrated | P1 |
| link_topsymlinks | ✅ | ❌ | Missing | P3 |
| Debouncing | ✅ | ❌ | Missing | P2 |
| Help Screen | ✅ | ✅ | Complete | - |
| LazyGit Integration | ✅ | ✅ | Complete | - |
| Open PR in Browser | ✅ | ✅ | Complete | - |
| Shell Integration | ✅ | ✅ | Complete | - |
| Caching | ✅ | ✅ | Complete | - |
| Unit Tests | ✅ | ❌ | Missing | P4 |
| Integration Tests | ✅ | ❌ | Missing | P4 |

**Legend:**
- ✅ Complete
- ⚠️ Partial / Not Integrated
- ❌ Missing / Stubbed

---

**Last Updated:** 2025-12-29
**Go Version:** Based on commit `eb8edcd`
**Python Version:** Latest on main branch
