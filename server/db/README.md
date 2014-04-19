Item parameters
===============

`id` (string) unique id of the item (example: "fishing_pole")

`name` (string) name of item

`desc` (string) description of item

`type` (string) item type:

* `(empty string)`: cannot be equipped
* `"consumable"`: using the item will trigger effect, but destroy the item (example: healing bottle, apple...)
* `"1hand"`: 1 hand item (examples: dagger, short sword, shield...)
* `"2hand"`: 2 hand item (examples: staff, bow, fishing pole...)
* `"head"`: head-wear item (examples: bandana, mask, wooden helmet...)
* `"torso"`: main body item (examples: dragon armor...)
* `"hands"`: hand item (examples: gloves...)
* `"legs"`: legs item (examples: leggings, shorts...)
* `"feet"`: feet item (examples: metal boots, flip-flops...)

`groups` (array of strings) item group (example: staff, sword)

`damage_min` (int) minimum damage

`damage_max` (int) maximum damage

*Other parameters are not yet fully supported*

Location parameters
===================

`id` (string) unique id of the item (example: `"a tree"`)

`name` (string) name of item  (example: `"A big tree"`)

`desc` (string) description of item  (example: `"You see a very tall tree"`)

`desc_short` (string) description of location written next to travel directions (example: `"north (a tree)"`)

*Other parameters are not yet fully supported*
