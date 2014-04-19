Olduar
======

Democratic multiplayer text adventure

REST Api
========

### GET /room/look
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

### GET /room/go/{direction}
Invalid __{direction}__ is ignored
Return is the same as `/room/look`

### GET /room/do/{action}
Invalid __{action}__ is ignored
Return is the same as `/room/look`

### GET /room/pickup/{item}
If invalid __{item}__ is specified or __{item}__ is not on ground, __null__ is returned
Otherwise return is the same as `/room/look`

### GET /room/drop/{item}
If invalid __{item}__ is specified or __{item}__ is not in player's inventory, __null__ is returned
Otherwise return is the same as `/room/look`

### GET /room/use/{item}
If invalid {item} is specified __null__ is returned
Otherwise return is the same as `/room/look`

### GET /room/inventory
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

### GET /room/inspect/{item}
Returns description of object or __null__ if __{item}__ is not available on ground or in inventory
```json
{
    "name": "Fishing pole",
    "desc": "A simple fishing pole"
}
```
*Attributes and other properties will be added over time*

### GET /rooms
Returns list of rooms available (only visible rooms are listed)
```json
[
    "room1",
    "groupe_le_france"
]
```

### GET /join/{room}
Joins the specified __{room}__
If player is already in some room, player will leave that room automatically (this is ignored if __{room}__ is same as player's current room)
If __{room}__ does not exist it will be created
Return is the same as `/{room}/look` or __null__ if maximum amount of rooms is reached or no __{room}__ is specified

### GET /leave
Leave the current room
Return is the same as `/rooms`

### GET /players
*TBA*
Returns list of all active players
```json
[
    "Belzebub",
    "Arthur",
    "Kain"
]
```
*Attributes and other properties will be added over time*

### GET /room/players
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