Olduar
======

Democratic multiplayer text adventure

REST Api
========

All requests need to use Basic Authentication

### GET /api/look
Returns current location view (together with items, npcs, exits and actions)
```json
{
    "name": "Lake",
    "desc": "You are looking at the beautiful lake full of fish.",
    "history": [],
    "exits": {
        "back": "a pub",
        "east": "a lake",
        "northeast": "a lake"
    },
    "actions": {
        "fishing": "go fishing"
    },
    "items": [
        {
            "id": "fishing_pole",
            "name": "Fishing pole",
            "desc": "A simple fishing pole"
        }
    ]
}
```

### GET /api/go/{direction}
Invalid __{direction}__ is ignored

Return is the same as `/look`

### GET /api/do/{action}
Invalid __{action}__ is ignored

Return is the same as `/look`

### GET /api/pickup/{item}
If invalid __{item}__ is specified or __{item}__ is not on ground, __null__ is returned

Otherwise return is the same as `/look`

### GET /api/drop/{item}
If invalid __{item}__ is specified or __{item}__ is not in player's inventory, __null__ is returned

Otherwise return is the same as `/look`

### GET /api/use/{item}
If invalid __{item}__ is specified __null__ is returned

Otherwise return is the same as `/look`

### GET /api/inventory
Returns array of items in player's inventory
```json
[
	{
		"id": "fishing_pole",
		"name": "Fishing pole",
		"desc": "A simple fishing pole"
	},
	{
		"id": "fish",
		"name": "Fish",
		"desc": "A small fish"
	}
]
```

### GET /api/inspect/{item}
Returns description of object or __null__ if __{item}__ is not available on ground or in inventory
```json
{
    "name": "Fishing pole",
    "desc": "A simple fishing pole"
}
```
*Attributes and other properties will be added over time*

### GET /api/party
Returns list of players in room, if player is not joined to a room `[]`(empty array) is returned
```json
[
    "Belzebub",
    "Arthur",
    "Kain"
]
```

### POST /api/say
Everything in __POST body__ is used as is and sent to all players as message

If player is not connected to any room __false__ is returned, otherwise __true__

### POST /api/join/{room}
Joins the specified __{room}__

If player is already in some room, player will leave that room automatically (this is ignored if __{room}__ is same as player's current room)

If __{room}__ does not exist it will be created

Return is the same as `/look` or __null__ if maximum amount of rooms is reached or no __{room}__ is specified

### POST /api/leave
Leave the current room

Return is the same as `/rooms`

### POST /api/rename
Everything in __POST body__ is used as string and set as __Name__

Returns __false__ if __empty string ("")__ is sent, otherwise __true__

### POST /api/register
Values for __Username__ & __Password__ are used from Basic Authentication headers

If __Username__ is taken or __Username__ is not valid, __false__ is returned, otherwise __true__

### GET /api/rooms
Returns list of rooms available (only visible rooms are listed)
```json
[
    "room1",
    "groupe_le_france"
]
```

### GET /api/players
Returns list of all active players
```json
[
    "Belzebub",
    "Arthur",
    "Kain"
]
```

### POST /api/tell/{player}
Everything in __POST body__ is used as is and sent to __{player}__ as message

If __{player}__ is not connected __false__ is returned, otherwise __true__