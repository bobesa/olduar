Action properties
=================

`id` (__string__) id of action (example: "fishing")

`charges` (__int__, __optional__) amount of times action can be used

`charges_msg` (__string__, __optional__) message sent to player if this action is out of charges

`requirements` (__array of requirements__, __optional__) requirements that must be met to do the action

`action` (__string__) type of action

`config` (__object__) configuration for the action

### All action types

`msg_party` (__string__) message sent to everyone except player doing the action

`msg_player` (__string__) message sent to player doing the action

`msg_all` (__string__) message sent to everyone

*All messages can contain special keywords (such as `%player%`) & properties from config (such as `%value%`) that will translate to actual values*

### Action type: message

*In following example we declare a __hug__ action that will just send a message to party `"Belzebub hugged a tree"` and message `"You hugged a tree"` to player doing the action*
```json
{
    "id":"hug",
    "action":"message",
    "config":{
        "msg_party":"%player% hugged a tree",
        "msg_player":"You hugged a tree"
    },
    "desc":"hug a tree"
}
```

### Action type: effect

`type` (__string__) type of effect

`value` (__int__) value of effect

*In following example we declare __drink__ action that will heal (`"type":"heal"`) player for 50 points (`"value":50`) and sending a message to party & player*
```json
{
    "id":"drink",
    "action": "effect",
    "config": {
        "type": "heal",
        "value": 50,
        "msg_party": "%player% drank the water from fountain and got healed for %value% damage",
        "msg_player": "You drank the water from fountain"
    },
    "desc":"drink the water"
}
```

### Action type: location

`type` (__string__) type of location action (currently __use__ only)

`value` (__string__) value of location action

*In following example we declare a __fishing_pole__ that will call action on current location (`"type":"use"`) with id __fishing__ (`"value": "fishing"`)*
```json
{
    "id": "fishing_pole",
    "name": "Fishing pole",
    "desc": "A simple fishing pole",
    "type": "2hand",
    "class": "staff",
    "damage_min": 1,
    "damage_max": 5,
    "actions":[{
        "action":"location",
        "config":{
            "type": "use",
            "value": "fishing"
        }
    }]
}
```

### Action type: give

`amount` (__int__, __optional__) amount of items to give to player (__default__: 1)

`items` (__array of items__) selection of items that player can get - see properties for each item below

* `id` (__string__) what item will be given (invalid id or empty string means nothing will be given)
* `chance` (__float__) what is relative chance for item to be selected
* `msg_party` (__string__, __optional__) what will be sent to other players if item is given to player
* `msg_player` (__string__, __optional__) what will be sent to player if item is given to player

*In following `give` configuration we declare that player can get 1 item (`"amount":1`) from selection of items (`"items":[...`)*

```json
{
    "amount":1,
    "items":[
        {
            "id":"",
            "chance":0.85,
            "msg_party":"%player% tried fishing but came empty handed!",
            "msg_player":"You failed to catch anything!"
        },
        {
            "id":"useless_shoe",
            "chance":0.1,
            "msg_party":"%player% caught a useless shoe!",
            "msg_player":"You caught a useless shoe!"
        },
        {
            "id":"fish",
            "chance":0.5,
            "msg_party":"%player% caught a fish!",
            "msg_player":"You caught a fish!"
        },
        {
            "id":"goldfish",
            "chance":0.001,
            "msg_party":"%player% caught a goldfish!",
            "msg_player":"You caught a goldfish!"
        }
    ]
}
```
*Please note that declaring `"id":""` means no item will be given. Also note that if you sum `count` of all items it is not 100%, this is completely valid as `chance` is relative against other items in the table*

Action Requirements properties
==============================

`type` (__string__) type of requirement (currently only `item` is supported)

`value` (__string__) value for for requirement (this must always be string, even if is suppose to be number)

`error_msg` (__string__, __optional__) message sent to player if he does not match the requirement

*In following example we declare that __fishing_pole__ (in inventory) is required for this action to be used, if not we send player message `"You cannot fish without fishing pole"`*
```json
"requirements":[
    {
        "type": "item",
        "value": "fishing_pole",
        "error_msg": "You cannot fish without fishing pole"
    }
],
```

__*Please note that when requirements are checked, they are checked in order declared*__

* If player don't have __fishing_pole__, checking will stop right there and we send player message `"You cannot fish without fishing pole"`
* If player have __fishing_pole__ but don't have __bait__, checking will stop and we send player message `"You cannot fish without bait"`
* If player have both __bait__ & __fishing_pole__ we will proceed and do the action
```json
"requirements":[
    {
        "type": "item",
        "value": "fishing_pole",
        "error_msg":"You cannot fish without fishing pole"
    },
    {
        "type": "item",
        "value": "bait",
        "error_msg":"You cannot fish without bait"
    }
]
```

Item properties
===============

`id` (__string__) unique id of the item (example: "fishing_pole")

`name` (__string__) name of item

`desc` (__string__) description of item

`type` (__string__) item type:

* `(empty string)`: cannot be equipped
* `"consumable"`: using the item will trigger effect, but destroy the item (example: healing bottle, apple...)
* `"1hand"`: 1 hand item (examples: dagger, short sword, shield...)
* `"2hand"`: 2 hand item (examples: staff, bow, fishing pole...)
* `"head"`: head-wear item (examples: bandana, mask, wooden helmet...)
* `"torso"`: main body item (examples: dragon armor...)
* `"hands"`: hand item (examples: gloves...)
* `"legs"`: legs item (examples: leggings, shorts...)
* `"feet"`: feet item (examples: metal boots, flip-flops...)

`groups` (__array of strings__) item group (example: staff, sword)

`actions` (__array of actions__, __optional__) actions triggered by `"use"` command

*Other parameters are not yet fully supported*

Location properties
===================

`id` (__string__, __optional__) unique id of the item (example: `"a tree"`)

`name` (__string__) name of item  (example: `"A big tree"`)

`desc` (__string__) description of item  (example: `"You see a very tall tree"`)

`desc_short` (__string__) description of location written next to travel directions (example: `"a tree"`)

`actions` (__array of actions__, __optional__) actions triggered by `"do"` command

`items` (__array of objects__, __optional__) what items can be found. Each item is listed by "group"/"id" and optionally "chance"

`exits` (__array of exits__, __optional__) specifically declare what exits will be in this location, if this is not declared system will randomly choose exits in same region as current location

*In following example we declare that we want to have a Crossroads location from where we can go to mountains, pub and lake. We also declare that __empty_mug__ can be found in this location (with 10% chance of spawn)*
```json
{
    "id":"amazing_crossroads",
    "name":"Crossroads",

    "desc":"You are standing in the middle of the crossroads.",
    "desc_short":"a crossroads",

    "items":[
        {"id":"empty_mug","chance":0.1}
    ],

    "exits":[
        {"id":"mountains","region":"mountains"},
        {"id":"pub","region":"start_pub"},
        {"id":"lake","region":"plains_lake"}
    ]
}
```

Location item properties
========================

`"group"` (__string__) item will be spawned from the group (example: swords)

`"id"` (__string__) item will be spawned by selected item id (example: fishing_pole)

`"chance"` (__float__, __optional__) chance of spawning which is from 0.001 (0.1%) to 1.0 (100%), please note that 0.0 is also __default__ which is 100%

*In following example we declare that 3 items may drop in this location: fishing pole (100% chance), fish (10% chance) and one random sword (0.5% chance)*
```json
"items":[
    {"id":"fishing_pole"},
    {"id":"fish","chance":0.1},
    {"group":"swords","chance":0.005}
]
```
