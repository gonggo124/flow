#>Params
##stackptr register

execute as @e[tag=_flow_internal.stack.bit,type=marker] if score @s _flow_internal.bitaddr >= #stackptr _flow_internal.register run kill @s