# Secret Santa Telegram Bot 🎅🎁

A simple Telegram bot to organize and manage a Secret Santa event within a group chat! This bot helps create a participant list, shuffle the names, and privately notify each participant about whom they are a Secret Santa for.

---


## Commands

### Private Commands
These commands must be initiated in private chat with the bot:
- `/start`  
  Initializes the bot for your account. This step is required before adding the bot to a group.

### Group Commands
These commands can be used in the group chat where the bot is added:
- `/me`  
  Adds yourself to the list of Secret Santa participants.
- `/list`  
  Displays the list of participants in the group.
- `/shuffle`  
  Randomly assigns Secret Santas to participants and sends each participant a private message with their assignment.

---

## Setup

To deploy and run this bot, follow the steps below:

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd secret-santa-bot

2. **Set `TELEGRAM_BOT_TOKEN` env**

3. Run `main.go`
