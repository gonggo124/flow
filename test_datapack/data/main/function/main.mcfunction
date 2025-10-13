#> COMPILED BY FLOW
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run tag @s add _flow_internal.stack.old_baseptr
scoreboard players set #sa0 _flow_internal.register 0
function _flow_internal:mem/stack/push
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/baseptr/attach
scoreboard players set #sa0 _flow_internal.register 10
function _flow_internal:mem/stack/push
scoreboard players set #r0 _flow_internal.register 20
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle on passengers if entity @s[tag=_flow_internal.stack.bit,type=marker] run scoreboard players operation @s _flow_internal.stack = #r0 _flow_internal.register
scoreboard players set #r2 _flow_internal.register 3
scoreboard players set #r3 _flow_internal.register 5
scoreboard players operation #r0 _flow_internal.register = #r2 _flow_internal.register
scoreboard players operation #r0 _flow_internal.register += #r3 _flow_internal.register
scoreboard players set #r1 _flow_internal.register 
scoreboard players operation #r0 _flow_internal.register *= #r1 _flow_internal.register
scoreboard players operation #sa0 _flow_internal.register = #r0 _flow_internal.register
function _flow_internal:mem/stack/push
function main:something
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/stackptr/attach
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/ret
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/cut
