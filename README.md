# steam-bulk-messenger
Send Steam chat messages in bulk to specific friends.

## Usage
Edit `config.json` to add your steam credentials. 
Start the program, you will need to provide the SteamGuard token in each run.
As your first run choose the option to generate a `friends.json` file. Go through your friend list answering `y/n` if you want to include specific people or not as destination for your chat messages. Optionally, add comma separated tags to them.

Now open `friends.json` in a text editor. Your file should look like this:

```json
{
	"Friends": [
		{
			"SteamID": 76561197961557177,
			"Name": "Snow",
			"Tags": [
				"br"
			]
		},
		{
			"SteamID": 76561197998011675,
			"Name": "Bleu",
			"Tags": [
				"cone"
			]
		}
	],
	"Messages": [
		{
			"DestinationTags": [
				""
			],
			"Message": "NS2"
		}
	]
}
```

Add one or more messages in the `Messages` key like so:

```json
...
"Messages": [
                {
                    "DestinationTags": [
                        "cone"
                    ],
                    "Message": "NS2"
                },
                {
                    "DestinationTags": [
                        "br"
                    ],
                    "Message": "hue"
                }
            ]
```

If you have an element in the `DestinationTags` array which is an empty string, it will send the message to all your picked friends no matter their tag:
```json
...
{
    "DestinationTags": [
        ""
    ],
    "Message": "hey"
}
...
```
If there are multiple matching tags, it will send the message of the first matching tag and no more.

When you're done, run the program again and choose the option to send messages.

### Other details
- You can manually add friends to `friends.json` instead of going through your entire friend list again.
- It only sends messages to online friends.


