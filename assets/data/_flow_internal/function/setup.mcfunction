#> ### SETUP FUNCTION ###

## # 레지스터 정의 #
scoreboard objectives add _flow_internal.register dummy
## 스택 전용 레지스터

#(1-2)-(3-(4-5))
#a=1
#b=2
#a-=b
#temp0=a
#a=4
#b=5
#a-=b
#temp1=a
#a=3
#a-=temp1
#temp0=a

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

## 범용 레지스터
scoreboard players set #r1 _flow_internal.register 0
scoreboard players set #r2 _flow_internal.register 0
# scoreboard players set #r3 _flow_internal.register 0

## # 포인터 소환 #
# stackptr `de8d7920-b907-4853-b3a2-c73cb0d5a84d`
execute store success score #t0 _flow_internal.register run summon marker 0 0 0 {UUID:[I;-561153760,-1190705069,-1281177796,-1328175027]}
execute if score #t0 _flow_internal.register matches 0 run tellraw @a {"text":"스택 포인터 소환에 실패했습니다. 프로그램이 제대로 동작하지 않을 것입니다."}
# baseptr `6a56ec26-fbbd-4b1c-a7bf-59d89fd54460`
execute store success score #t0 _flow_internal.register run summon marker 0 0 0 {UUID:[I;1784081446, -71480548, -1480631848, -1613413280]}
execute if score #t0 _flow_internal.register matches 0 run tellraw @a {"text":"베이스 포인터 소환에 실패했습니다. 프로그램이 제대로 동작하지 않을 것입니다."}

## # 스택 #
scoreboard objectives add _flow_internal.stack dummy

## # 메모리 주소 #
scoreboard objectives add _flow_internal.bitaddr dummy

function main:main