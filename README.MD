# gli

```text
  ________  ___       ___     
 /  _____/ /  /      /__/     
/  /  ___ /  /       __       
\  \__\  /  /____   /  /      
 \______/ /______/ /__/       
```

> **a git wrapper to make developer life easy**

A lightweight Git wrapper that simplifies complex workflows into intuitive, high-density terminal commands. Designed for speed, precision, and visual clarity.

---

## **🚀 Quick Start**

### **1. Clone & Enter**
```bash
git clone https://github.com/Gazi2050/gli.git
cd gli
```

### **2. Install Dependencies**
```bash
pip3 install rich requests
```

### **3. Run**
```bash
python3 app.py
```

---

## **💎 Core Features**

| Feature | command | Description |
| :--- | :--- | :--- |
| **Commit & Push** | `-c`, `--commit` | Auto-stage, September-timestamp, and push. |
| **History** | `-l`, `--log` | High-density commit history with time. |
| **Recovery** | `-rl`, `--reflog` | Clarified Reflog with numeric selectors. |
| **Reset** | `-rs <mode>`, `--reset <mode>` | Reset to previous commit. Modes: `soft` (keeps changes) or `hard` (discards). |
| **Branch** | `-s <name>`, `--switch <name>` | Create, switch, and track new branch. |
| **Profile** | `me` \| `profile <user>` | GitHub dashboard (Stats, Meta, Bio). |
| **Time Warp** | `-ct [<date>]`, `--changeTime` | Target specific/all/last commit timestamps. |
| **Identity** | `-ca`, `--changeAuthor` | Target specific/all/last author details. |
| **Rewrite** | `-cm`, `--changeMessage` | Update last commit message instantly. |

<p align="center">Made with ❤️ for Git users</p>