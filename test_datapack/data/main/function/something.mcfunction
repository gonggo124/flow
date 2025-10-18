#> COMPILED BY FLOW
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run tag @s add _flow_internal.stack.old_baseptr
scoreboard players set #sa0 _flow_internal.register 0
function _flow_internal:mem/stack/push
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/baseptr/attach
say hi
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/stackptr/attach
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/ret
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/cut
