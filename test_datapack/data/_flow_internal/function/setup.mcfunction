#> ### SETUP FUNCTION ###

## # 레지스터 정의 #
scoreboard objectives add _flow_internal.register dummy
## 스택 전용 레지스터
scoreboard players set #stackptr _flow_internal.register 0
scoreboard players set #baseptr _flow_internal.register 0
## 시스템 함수
# 인자 전달 레지스터
scoreboard players set #sa0 _flow_internal.register 0
# 반환값 전달
scoreboard players set #sreturn _flow_internal.register 0
## 사용자 함수
# 인자 전달 레지스터
scoreboard players set #a0 _flow_internal.register 0
# 반환값 전달
scoreboard players set #return _flow_internal.register 0
## 임시 레지스터
scoreboard players set #t0 _flow_internal.register 0

## # 스택 #
scoreboard objectives add _flow_internal.stack dummy

## # 메모리 주소 #
scoreboard objectives add _flow_internal.bitaddr dummy

function main:main