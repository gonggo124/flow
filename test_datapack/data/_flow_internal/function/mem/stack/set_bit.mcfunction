tag @s add _flow_internal.stack.bit
tag @s add _flow_internal.stack.bit.new

execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run ride @e[tag=_flow_internal.stack.bit.new,type=marker,limit=1] mount @s

scoreboard players operation @s _flow_internal.stack = #sa0 _flow_internal.register

ride de8d7920-b907-4853-b3a2-c73cb0d5a84d dismount
ride de8d7920-b907-4853-b3a2-c73cb0d5a84d mount @s

tag @s remove _flow_internal.stack.bit.new