tag @s add _flow_internal.stack.bit
tag @s add _flow_internal.stack.bit.new

data merge entity @s {width:0,height:0}

execute store result score @s _flow_internal.bitaddr if entity @e[tag=_flow_internal.stack.bit,type=interaction]

execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run ride @e[tag=_flow_internal.stack.bit.new,type=interaction,limit=1] mount @s

scoreboard players operation @s _flow_internal.stack = #sa0 _flow_internal.register

tellraw @a {"score":{"name":"@s","objective":"_flow_internal.bitaddr"}}
say ===================================
function _flow_internal:mem/stack/stackptr/attach

tag @s remove _flow_internal.stack.bit.new
