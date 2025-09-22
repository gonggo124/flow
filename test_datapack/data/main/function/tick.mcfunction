execute as @e[tag=bullet,type=block_display] at @s run function _flow_internal:do/update
execute as @e[tag=bullet,type=block_display] at @s unless block ~ ~ ~ air run function _flow_internal:do/hit_wall
execute as @e[tag=bullet,type=block_display] at @s run function _flow_internal:anony/0

function _flow_internal:anony/1
