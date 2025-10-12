#> ### SETUP FUNCTION ###

## # 레지스터 정의 #
scoreboard objectives add _flow_internal.register dummy
## 스택 전용 레지스터

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

## # 포인터 소환 #
# stackptr `de8d7920-b907-4853-b3a2-c73cb0d5a84d`
summon marker 0 0 0 {UUID:[I;-561153760,-1190705069,-1281177796,-1328175027]}
# baseptr `6a56ec26-fbbd-4b1c-a7bf-59d89fd54460`
summon marker 0 0 0 {UUID:[I;1784081446, -71480548, -1480631848, -1613413280]}

## # 스택 #
scoreboard objectives add _flow_internal.stack dummy

## # 메모리 주소 #
scoreboard objectives add _flow_internal.bitaddr dummy

function main:main