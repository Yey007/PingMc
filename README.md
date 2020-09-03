# PingMc
A discord bot that pings minecraft servers for player counts written in Go

`.pingmc help` for usage in discord

You may build the bot and run it with your bot token using this command: `./pingmc -t key`

~Or you can use the centralized hosting. Add the bot to your server using this link: https://discord.com/api/oauth2/authorize?client_id=734445179800649739&permissions=3072&scope=bot

Note: The bot must have access to the channel you designate.

# Upcoming features and changes
1. Fix the help message. Currently it still shows the old way of creating pinging channels where the bot creates a new channel.
2. Add a way to stop the bot other than deleting the ping channel (by IP? Cancellation tokens?)
3. Add FML2 support (Forge 1.13+)
4. Add a way to look up modlists
5. Add a way to look up if someone is online (possibly vanilla only, forge doesn't seem to have this feature)
6. Server type inferrence
