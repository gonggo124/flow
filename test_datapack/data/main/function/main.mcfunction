#> COMPILED BY FLOW

## 스택 넣기
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run tag @s add _flow_internal.stack.old_baseptr
scoreboard players set #sa0 _flow_internal.register 0
function _flow_internal:mem/stack/push
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/baseptr/attach

## int a = 0;
scoreboard players set #sa0 _flow_internal.register 120
function _flow_internal:mem/stack/push
## # something(a,169)
scoreboard players set #sa0 _flow_internal.register 169
function _flow_internal:mem/stack/push

execute \
    as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 \
    on vehicle on passengers \
    if entity @s[tag=_flow_internal.stack.bit,type=marker] \
run scoreboard players operation #sa0 _flow_internal.register = @s _flow_internal.stack
function _flow_internal:mem/stack/push

function main:something

execute \
    as de8d7920-b907-4853-b3a2-c73cb0d5a84d \
    on vehicle on vehicle on vehicle \
    run function _flow_internal:mem/stack/stackptr/attach
function _flow_internal:mem/stack/cut
## something(a,169) #

## 스택 정리
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/stackptr/attach
execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function _flow_internal:mem/stack/ret
execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function _flow_internal:mem/stack/cut
