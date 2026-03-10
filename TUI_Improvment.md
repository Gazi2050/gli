## gli – TUI Improvement Plan

### 1. goals – what we want from the new TUI

- **Unify styling and theming** across all TUI surfaces (commit flow, history editor, help, profile) using a vibrant, modern look inspired by `diny`.
- **Improve structure and readability** of views (e.g., help screen table should look like a real table, profile and messages in framed boxes, consistent borders/padding).
- **Make interactions feel richer** in commit and history workflows without changing CLI flags or breaking flows, only refining keybindings/UX when clearly better.

---

### 2. theming & layout layer – shared building blocks

- **Theme model**
  - Introduce a `Theme` struct similar to `diny/ui/themes/theme.go`:
    - Fields for role-based colors: primary/success/error/warning foreground/background, muted foreground, using **minimal and eye-relaxing palettes** (no overly bright/neon colors by default).
  - Add helpers like:
    - `SetTheme(name string) bool`
    - `GetCurrentTheme() *Theme`
    - `GetAvailableThemes() []string` (and maybe dark/light variants).
  - Start with a sensible vibrant default (e.g., Catppuccin-like), but keep it easy to extend.
- **Layout & box helpers**
  - Take inspiration from `diny/ui/ui.go`:
    - Implement a `Box` helper that:
      - Computes terminal width (with a safe fallback).
      - Applies consistent borders and padding.
      - Chooses width **responsively**:
        - On small terminals, shrink to fit while keeping content readable (wrap text instead of overflowing).
        - On large terminals, cap at a sensible max width so content doesn’t become too wide, and center or align the box for a focused reading area.
      - Uses variants: `Primary`, `Success`, `Error`, `Warning`.
  - Add convenience helpers:
    - `RenderTitle(text string)` for section titles.
    - `WithSpinner(message string, fn func() error) error` for long-running actions.
    - Potentially a generic `MessageBox(variant, title, body)` for structured messages.
- **Theme selection integration**
  - Decide how gli picks a theme:
    - Simple option: read from `git config gli.theme`.
    - Or a small config file (e.g., `~/.config/gli/config.yaml` or similar).
  - Ensure all TUI views pull color roles from `GetCurrentTheme()` instead of hard-coded `lipgloss.Color("42")` style constants.

---

### 3. help screen – real table + unified styling

- **Refactor `internal/ui/help.go`**
  - Reuse existing table logic (column widths, multi-line support), but:
    - Replace numeric colors with theme-based styles:
      - Header, cell text, borders, usage hints.
    - Wrap the whole table in a themed `Box` or panel, so it visually matches other views.
  - Ensure columns render like a proper table:
    - Command, Flag, Description aligned in consistent widths.
    - Multi-line flags (like switch/local-branch/remote-branch) stay aligned across rows.
  - Make the table **responsive**:
    - On narrow terminals, allow columns to wrap gracefully and consider reducing padding so it still fits.
    - On very wide terminals, keep table width bounded to the same max width rules as `Box`, to avoid overly stretched rows.
- **Improve content presentation**
  - Organize `helpRows` into logical groups:
    - Commit/AI features.
    - History editing.
    - Branch/log tools.
    - Profile and misc.
  - Use subtle separators or labels between groups (e.g., blank line, sub-heading row with different style).
  - Integrate ASCII logo and tagline into the new style:
    - Styled logo at top.
    - Tagline in a muted style under the logo.
    - Table section framed below them.
- **Usage hints**
  - Keep a clear “Usage example” line at the bottom.
  - Style the binary name and flags with the primary color, and surrounding text with muted color.

---

### 4. commit flow – more interactive AI + manual commit UX

- **Adopt theme-based styles in `internal/ui/commit.go`**
  - Replace direct `lipgloss.Color` usage in `CommitModel.View()` with themed styles:
    - Titles, hints, success, error messages, borders.
  - Use the shared `Box` helper for the outer frame so commit screens line up visually with help/profile/history.
- **AI proposal state improvements**
  - When `commitStateShowingProposal`:
    - Render the AI commit message in a dedicated sub-box:
      - Clear title, e.g. “AI Commit Proposal”.
      - Message body with appropriate wrapping and padding.
    - Present options as a clean, readable menu:
      - `[1] Commit and push`
      - `[2] Regenerate message`
      - `[3] Edit message`
      - `[4] Cancel (or q)`
    - Use consistent alignment and spacing so options read like a menu, not just raw lines.
- **Editing state polish**
  - In `commitStateEditing`:
    - Show a clear “Edit Commit Message” title.
    - Use the text input with enough width and padding inside the box.
    - Display validation errors inside an `Error` variant area (or clearly colored line).
    - Standardize hints at the bottom (e.g., `Enter to commit • Esc to go back • Ctrl+C to quit`) in a muted style.
- **Committing / loading feedback**
  - For `commitStateStaging`, `commitStateLoadingAI`, `commitStateCommitting`:
    - Use spinner with theme-driven colors.
    - Provide short, descriptive messages (“Staging files…”, “Generating AI commit message…”, “Committing and pushing changes…”).
  - Optionally use `WithSpinner`-style wrappers for non-TUI paths that need simple feedback.
- **Done state**
  - For `commitStateDone`:
    - Show success in a `Success`-themed message box with a clear confirmation line.
    - On error, show an `Error`-themed box and a short hint (e.g., “Press q to exit”).

---

### 5. history editor – clearer forms, warnings, and confirmations

- **Apply shared theming in `internal/ui/history.go`**
  - Replace `lipgloss.Color` calls with theme-based roles (titles, errors, success, hints).
  - Wrap all views (`historyStateScopeSelect`, `historyStateHashInput`, input states, confirm, running, done) inside a consistent `Box` frame.
- **Scope selection UI**
  - In `historyStateScopeSelect`:
    - Use a clear title like “Change Commit Time/Author/Message” depending on operation.
    - Present scope choices in an interactive menu similar to the diny commit screen, e.g.:
      - A heading like `Scope` or `What would you like to change?`.
      - A list of options navigated with arrow keys or `j`/`k`, with a visible selector:
        - `Last commit`
        - `Specific hash`
        - `All commits`
      - A short help line like: “Select an option using arrow keys or j,k and press Enter”.
    - Add a short, themed warning line explaining that rewriting history can be destructive, especially for “all commits”.
- **Input form layout**
  - For `historyStateHashInput`:
    - Label clearly: “Enter target commit hash”.
    - Use a single text input with good width and padding, with validation feedback beneath it.
  - For `historyStateTimeInput`:
    - Two fields:
      - `Date (YYYY-MM-DD)`
      - `Time (HH:MM or HH:MM:SS)`
    - Show consistent hints about Tab to switch inputs, Enter to move forward, Esc to go back.
  - For `historyStateAuthorInput`:
    - Two fields:
      - `Author name`
      - `Author email`
    - Same navigation hints as time inputs.
  - For `historyStateMessageInput`:
    - Single “New commit message” field.
    - Enter confirms, Esc cancels, with errors shown clearly if empty.
- **Confirmation summary**
  - In `historyStateConfirm`:
    - Render a summary block with:
      - Operation (time/author/message).
      - Scope (`single`, `specific`, `all`).
      - Target hash (if `specific`).
      - New values (datetime, author, or message).
    - Format as bullet list or simple key/value layout inside the box.
    - Provide a bold, clear prompt: `Proceed? [y] Yes / [n] No (q to cancel)`.
- **Running and done states**
  - In `historyStateRunning`:
    - Show a spinner and short message: “Applying history changes…”.
  - In `historyStateDone`:
    - For success, show a `Success`-themed message (“History updated successfully”) with hint to exit.
    - For errors, show details in an `Error`-themed area plus exit hint.
    - For cancelled flows, show a neutral/muted “History update cancelled” message.

---

### 6. profile screen – structured, theme-aware GitHub card

- **Unify styling in `internal/ui/profile.go`**
  - Replace hard-coded `lipgloss.Color` values with theme roles:
    - Title, name, login, bio, labels, stats, meta, repo names, borders.
  - Keep the existing layout idea:
    - Title (“GitHub Profile”).
    - Header: `Name (@login)`.
    - Bio line.
    - Stats line (public repos, followers, following).
    - Meta lines (location, Twitter, blog, joined).
  - Wrap everything in a `Box` or consistent panel style from the new layout layer.
- **Recent repos list**
  - For each repo:
    - Show name (bold), language, and stars in a structured way:
      - E.g., `<RepoName> (Language) ★ <stars>`.
    - Keep description as a second, indented line with a muted style.
  - Optionally give the repos section its own sub-title (“Recent Repositories”) in the primary color.

---

### 7. optional: themes, previews, and UX parity with diny

- **Theme presets**
  - Port a subset of `diny` themes into gli:
    - Catppuccin, Tokyo, Nord, Dracula, Gruvbox (light/dark), GitHub light.
  - Expose them via `GetAvailableThemes()` and use simple string keys.
- **Theme list/preview command**
  - Add a CLI entry (command or flag) to:
    - List available themes.
    - Optionally render a small preview for each (similar to `PrintThemeList` in `diny/ui/ui.go`), using the new `Box` helper.
  - Integrate with configuration so users can quickly set a theme.
- **Behavioral polish**
  - Consider using the box helpers for:
    - Config validation errors or warnings.
    - Post-command “success” summaries.
  - Keep core behavior and flags the same as current gli, focusing changes on visual polish and clarity.

---

### 8. testing, ergonomics, and documentation

- **Manual smoke testing**
  - Run through all main flows:
    - AI commit.
    - Manual commit.
    - Each history operation (time, author, message) with different scopes.
    - Help screen.
    - Profile screen for both `me` and a specific user.
    - Any theme selection or preview commands.
- **Terminal ergonomics**
  - Verify rendering on:
    - Narrow terminals (e.g., 80 columns).
    - Wider terminals, making sure boxes don’t become unreadable.
  - Adjust width calculations, wrapping, and padding as needed.
- **Accessibility & fallbacks**
  - Ensure there is at least one “minimal” or low-contrast theme for users who want a simpler look.
  - Keep important information legible even without 24-bit color support (safe color choices).
- **Docs & showcase**
  - Update `README.MD` or a dedicated docs section to:
    - Show screenshots/GIFs of the new TUI.
    - Document available themes and how to configure them.
    - Briefly describe the improved commit and history flows.

---

### 9. high-level checklist

1. **theming & layout**
  - Implement `Theme` struct and theme manager.
  - Add `Box`, `RenderTitle`, `WithSpinner`, and other helpers.
  - Refactor existing views to use the new theme roles instead of raw color codes.
2. **help screen**
  - Rework into a fully themed, properly aligned table wrapped in a panel.
  - Group commands and integrate logo/tagline with the new style.
3. **commit flow**
  - Refresh AI + manual commit views with themed boxes, structured options, and clear hints.
  - Ensure loading and commit states have good feedback.
4. **history flow**
  - Apply new theme and layout to all history states.
  - Improve form UX and confirmation summary, highlight risky scopes.
5. **profile UI**
  - Align profile screen with the shared theme and layout.
  - Make repo listing more structured and readable.
6. **themes & docs (optional but recommended)**
  - Port multiple themes, add a theme preview/list command.
  - Update documentation to highlight the refreshed TUI experience.

