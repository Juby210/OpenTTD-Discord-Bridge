# OpenTTD-Discord-Bridge
Bridge discord & OpenTTD chat. You can also check openttd server stats (including companies and players!) or manage it from discord.

### Config
- **OpenTTD** - openttd file path (for linux/freebsd leave only `openttd`)
- **Args** - optional server args (eg save file `["-g", "save.sav"]`)
- **Token** - discord bot token
- **ChannelID** - discord channel id to bridge chat
- **Prefix** - discord bot prefix
- **Admins** - discord bot admins id (admins can manage server from discord)

### Commands
- **ttd!stats** - Send server stats
- **ttd!clients** - Send connected clients

**Admin commands**:
- **ttd!save (file name)** - Save game
- **ttd!load (file name)** - Load save
- **ttd!restart (optional args)** - Restart openttd server
- **ttd!reset (company)** - Reset company
- **ttd!eval (command)** - Eval command to openttd (eg **ttd!eval pause** will pause game)
