#> COMPILED BY FLOW
scoreboard players operation #sa0 _flow_internal.register = #baseptr _flow_internal.register
scoreboard players operation #baseptr _flow_internal.register = #stackptr _flow_internal.register
function _flow_internal:mem/stack/push


scoreboard players set #sa0 _flow_internal.register 120
function _flow_internal:mem/stack/push


scoreboard players set #sa0 _flow_internal.register 169
function _flow_internal:mem/stack/push

scoreboard players operation #t0 _flow_internal.register = #baseptr _flow_internal.register
scoreboard players add #t0 _flow_internal.register 1
execute as @e[tag=_flow_internal.stack.bit,type=marker,limit=1] if score @s _flow_internal.bitaddr = #t0 _flow_internal.register run scoreboard players operation #sa0 _flow_internal.register = @s _flow_internal.stack
function _flow_internal:mem/stack/push

function main:something

scoreboard players remove #stackptr _flow_internal.register 2
function _flow_internal:mem/stack/cut


scoreboard players operation #stackptr _flow_internal.register = #baseptr _flow_internal.register
execute as @e[tag=_flow_internal.stack.bit,type=marker,limit=1] if score @s _flow_internal.bitaddr = #stackptr _flow_internal.register run scoreboard players operation #baseptr _flow_internal.register = @s _flow_internal.stack
function _flow_internal:mem/stack/cut
