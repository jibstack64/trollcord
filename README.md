# trollcord

![GitHub](https://img.shields.io/github/license/jibstack64/trollcord)
![GitHub all releases](https://img.shields.io/github/downloads/jibstack64/trollcord/total)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jibstack64/trollcord)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/jibstack64/trollcord)

A command-line utility for trolling on Discord.

### Disclaimer

Using this tool is **likely to get you banned from Discord**. It could be considered a modified client or form of [self-botting](https://support.discord.com/hc/en-us/articles/115002192352-Automated-user-accounts-self-bots-) under Discord terms of service.

This tool is designed to be *feature-rich*, not cloak your client from potential Discord suspicion.

I am **not** responsible for what you decide to do with this tool.

---

> #### Preview:
> ![screenshot-2023-01-02-22:11:08](https://user-images.githubusercontent.com/107510599/210281379-286192ea-455c-4cf6-940e-282abe5ea702.png)

### How-to

- #### **Windows (10+)**
    - Download the `trollcord.exe` executable from the [Releases](https://github.com/jibstack64/trollcord/releases) page.
    - Use the File Explorer to navigate to where you downloaded the executable.
    - Hold shift and right click anywhere in the explorer, as long as it is not over a file.
    - Click on the `Open Command Prompt` button on the dropdown.
    - Type `./trollcord.exe` and press enter - you should now be in the `trollcord` interface.
- #### **Linux**
    - Download the `trollcord` executable from the [Releases](https://github.com/jibstack64/trollcord/releases) page.
    - Open up a Terminal.
    - Navigate to where you downloaded it; e.g. `cd Downloads`.
    - Type `./trollcord` and press enter.
    - Boom! If you wish, you may add it to your `~/.local/bin` folder for convenience.
- #### **Mac**
    - The Mac instructions are identical to that of the Linux instructions. Just download the `trollcord-mac` executable instead of the `trollcord` one.

### Core features
1. Mass pinger
    - Mass pinging simply pings all roles in a server (as a cheeky alternative to `@everyone`).
    - Select the `mass pinger` option.
    - Input *either* a list of channel IDs split by `,` (e.g. `12345678,958928391,389123123`) OR a server ID.
    - A message will be formed containing a ping to every single role in the server - this message will be spammed in the provided channels / all channels in the provided server, pinging everyone regardless of `@everyone` permissions.
    - Once all parameters are specified, the pinging will begin.

2. Webhook spamming
    - Webhook spamming is often regarded as one of the most annoying forms of attacks.
    - To initiate a webhook spam, select the `webhook spammer` option and place your webhook url in the input; you will be asked to provide a username and icon URL for the webhook user.
    - Prepare for chaos.

3. Text channel spamming
    - Similarly to webhook spamming, text channel spamming is extremely fast. Many servers have bots that auto-ban spammers; for this reason, I suggest webhook spamming (it required webhook perms, but often owners are silly and leave the permission on for all users).
    - To start a text channel spam, simply select the `text channel spammer` option. Input the channel's Discord ID and message content, and bobs your uncle.

4. DM channel spamming
    - Identical to the text channel spammer, but for DMs; just selet the `dm spammer` option and give it a user id and the message to be sent.

5. Server destruction [**BOT-ONLY**]
    - This feature is only available for bot accounts - this is due to restrictions on user accounts fetching all users within a server (as you can imagine, DM bots would be more common if that wasn't the case).
    - `trollcord` offers a variety of options for destroying servers. These include: deleting roles, channels and banning users.
    - To setup a server destruction, select the `server destroyer` option, and provide the server's [Discord ID](https://www.remote.tools/remote-work/how-to-find-discord-id).
    - ~~You will be asked to blacklist said roles, channels and users by their Discord IDs. Blacklisted users will not be banned, blacklisted channels will not be removed, etcetera.~~
    - Once an ID ~~and blacklist (if any)~~ is provided, the server destruction will begin.

### Extra features
- List servers
    - Simply lists all servers that you are in.

---

### Style
Some people decided to make some style changes to the program, so I decided, hey, why not give them a shout! [jumbledFox's UwU design](https://github.com/jumbledFox/twollcowd) - [horacegill's.. questionable design](https://github.com/horacegill/trollcord).
