Action properties
=================

`id` (__string__) id of action (example: "fishing")

`charges` (__int__, __optional__) amount of times action can be used

`charges_msg` (__string, __optional__) message sent to player if this action is out of charges

`requirements` (__array of requirements__, __optional__) requirements that must be met to do the action

`action` (__string__) type of action

`config` (__object__) configuration for the action

*Action types will be added soon*

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
