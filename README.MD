# gli

### _Git wrapper_

A lightweight Git wrapper that simplifies common Git commands into short, intuitive terminal shortcuts.

Designed to save developers time, reduce errors, and streamline Git workflows with minimal command memorization.

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

### **Feature 06: GitHub Profile (ASCII View)**

Display a GitHub user’s public profile in a neofetch-style layout, including an ASCII-rendered avatar and key public information such as name, username, bio, repositories, followers, and location.

- Fetches data from the GitHub **public API**
- Converts the user’s avatar image into **ASCII art**
- Displays everything in a clean, terminal-friendly UI

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
