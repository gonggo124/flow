#>Params
##stackptr register

                                 # stackptr, baseptr 거르기
execute on passengers if entity @s[tag=_flow_internal.stack.bit] run function _flow_internal:mem/stack/cut

kill @s
