Olduar
======

Democratic multiplayer text adventure

REST Api
========

### GET /{room}/look
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

### GET /{room}/go/{exit_id}
Invalid {exit_id} is ignored
Return is the same as `/{room}/look`

### GET /{room}/do/{action_id}
Invalid {action_id} is ignored
Return is the same as `/{room}/look`

### GET /{room}/pickup/{item_id}
If invalid {item_id} is specified or {item_id} is not on ground, __null__ is returned
Otherwise return is the same as `/{room}/look`

### GET /{room}/drop/{item_id}
If invalid {item_id} is specified or {item_id} is not in player's inventory, __null__ is returned
Otherwise return is the same as `/{room}/look`

### GET /{room}/use/{item_id}
If invalid {item_id} is specified __null__ is returned
Otherwise return is the same as `/{room}/look`

### GET /{room}/inventory
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

### GET /{room}/inspect/{item_id}
Returns description of object
```json
{
    "name": "Fishing pole",
    "desc": "A simple fishing pole"
}
```
*Attributes and other properties will be added over time*

### GET /rooms
*TBA*
Returns list of rooms available (only visible rooms are listed)
```json
[
    "room1",
    "groupe_le_france"
]
```

### GET /join/{room}
*TBA*
Joins the specified room
If player is already in some room, player will leave that room automatically (this is ignored if specified room is same as player's current room)
Return is the same as `/{room}/look`

### GET /leave
*TBA*
Leave the current room
Return is the same as `/rooms`

### GET /{room_id}/players
*TBA*
Returns list of players in room
```json
[
    "Belzebub",
    "Arthur",
    "Kain"
]
```
*Attributes and other properties will be added over time*

### POST /{room}/say
*TBA*
Everything in __POST body__ is used as is and sent to all players as message
If player is not connected to