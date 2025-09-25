execute as @e[tag=bullet] at @s run function _flow_internal:do/update

execute as @e[tag=bullet] if score @s obj matches 20.. run function _flow_internal:do/kill
