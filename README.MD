# gli

### _Git wrapper_

A lightweight Git wrapper that simplifies common Git commands into short, intuitive terminal shortcuts.

Designed to save developers time, reduce errors, and streamline Git workflows with minimal command memorization.

## **Getting Started**

---

### **1. Clone the repository**
```bash
git clone https://github.com/Gazi2050/gli.git
cd gli
```

### **2. Install Dependencies**
`gli` uses `rich` for UI and `requests` for GitHub API interactions.
```bash
pip3 install rich requests
```

### **3. Run locally**
```bash
python3 app.py
```

### **4. (Optional) Create an alias**
To use `gli` directly from anywhere:
```bash
# Add this to your .bashrc or .zshrc
alias gli='python3 /absolute/path/to/gli/app.py'
```

## **Features**

---

### **Feature 01: Commit & Push**

Automatically stage all changes, commit using a **user-provided message**, and push to the remote repository.

```bash
gli -c "commit message"
gli --commit "commit message"
```

---

### **Feature 02: Git Log**

View the git commit history in a clean, readable format.

```bash
gli -l
gli --log
```

---

### **Feature 03: Git Reflog**

View the git reflog to track updates to the tip of branches

(useful for recovering lost commits).

```bash
gli -rl
gli --reflog
```

---

### **Feature 04: Reset Commit**

Reset the current branch to the previous commit.

- **soft** → keep changes staged
- **hard** → discard changes completely

```bash
gli -rs <soft|hard>
gli --reset <soft|hard>
```

---

### **Feature 05: Switch Branch**

Create a new branch, switch to it, and immediately push it to the remote repository.

```bash
gli -s <branch_name>
gli --switch <branch_name>
```

---

### **Feature 06: GitHub Profile (Simplified Dashboard)**

Display a GitHub user’s public profile in a high-density, neofetch-style dashboard. This view includes key public information such as name, username, bio, account statistics (followers, following, repos), and detailed account metadata.

- Fetches data from the GitHub **public API**
- Displays a clean, categorized, and terminal-friendly dashboard UI
- Optimized for quick data retrieval and visual clarity

```bash
gli me
gli profile <github_username>
```

**Behavior:**

- `gli me` → Automatically detects the GitHub username from local git config
- `gli profile <github_username>` → Fetches and displays the specified user’s public profile

---

### **Feature 07: Change Commit Author**

Interactively change the author name and email of the last commit.

- Prompts the user in the terminal to enter **author name** and **author email**
- Updates the last commit with the new author information

```bash
gli -ca
gli --changeAuthor
```

**Example Interaction:**

```
Enter author name: John Doe
Enter author email: john@example.com

✔ Author changed to John Doe <john@example.com>
```

---

### **Feature 08: Change Commit Time**

Interactively change the commit timestamp of the last commit.

- Prompts the user to select a **date** and **time** in the terminal
- Updates the last commit with the new commit timestamp

```bash
gli -ct
gli --changeTime
```

**Example Interaction:**

```
Select commit date (YYYY-MM-DD): 2025-12-24
Select commit time (HH:MM, 24h): 18:30

✔ Commit time changed to 2025-12-24T18:30:00
```

---

### **Feature 09: Change Commit Message**

Interactively change the commit message of the last commit.

- Prompts the user to enter a **new commit message** in the terminal
- Updates the last commit without creating a new one

```bash
gli -cm
gli --changeMessage
```

**Example Interaction:**

```
Enter new commit message: Fix critical bug in login flow

✔ Commit message updated successfully!
```
