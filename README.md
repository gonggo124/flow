## Flow 소개

### 아직 테스트 단계입니다

처음엔 객체지향을 사용 가능한 `mcfunction`의 `superset` 언어가 목표였지만, 현재는 코드 작성이 더 쉽고, 가독성이 좋은 언어를 목표로 하고있습니다.

코드 하이라이터: https:#marketplace.visualstudio.com/items?itemName=nutoo.flow-highlighter

`Flow`는 기존 `mcfunction` 문법에 약간의 명령어를 추가한 형태의 언어입니다.

그렇기 때문에 기존 바닐라 `mcfunction` 코드를 그대로 가져다 붙여넣어도 `Flow 컴파일러`는 문제를 일으키지 않습니다.

.fl의 기능이 필요없는 함수는 기존처럼 .mcfunction을 사용해서 구현할 수도 있습니다.

`Flow` 코드 예시는 다음과 같습니다:

```mcfunction
# main:bullet/init
data merge entity @s {\
    Tags: ["bullet"],\
    block_state:{Name:"blue_ice"},\
    transformation:{left_rotation:[0,0,0,1],right_rotation:[0,0,0,1],translation:[-0.1,-0.1,-0.1],scale:[0.2,0.2,0.2]}\
}

say this is INIT
say { is called brace

bind @s main:test/something

# main:tick
do @e[tag=bullet,type=block_display].update
execute as @e[tag=bullet,type=block_display] at @s unless block ~ ~ ~ air run do hit_wall
execute as @e[tag=bullet,type=block_display] at @s run fn {
    particle dust{color:[0,0,0],scale:1} ~ ~ ~ 0.1 0.1 0.1 0 2
    particle dust{color:[1,1,1],scale:1} ~ ~ ~ 0.1 0.1 0.1 0 2
}

fn { say 1}
```

보시다시피 일부 명령을 제외하면 바닐라 `mcfunction`과 문법이 일치합니다.

`Flow 컴파일러`가 Flow 고유 문법을 발견하면 그 자리에 컴파일된 바닐라 명령어를 채우는 식으로 작동합니다.

## 설치/사용

`winget`을 통해 설치가 가능하길 희망하지만 아직 지원하지 않습니다.  
대신 레포지토리의 [릴리즈](https:#github.com/gonggo124/flow/releases/tag/v0.1.0)에서 다운로드하여 환경 변수에 등록해놓으면 콘솔에서 사용할 수 있습니다.

`flow <대상 데이터팩 경로>`로 사요할 수 있습니다.

vscode 터미널이 실행중인 디렉토리를 컴파일하려면
`flow .`을 사용하십시오.

예:

```
C:\Users\User\myproject>flow .
'C:\Users\User\myproject' 경로가 컴파일됩니다..
```

## 기능

### Method Set

메서드셋이란 엔티티에게 부착할 수 있는 함수 묶음입니다.

`namespace/method/*.flm` 경로에 파일을 추가하여 정의할 수 있습니다.

대략적인 문법은 다음과 같습니다:

```js
// main:projectile
// main/method/projectile.flm
// 참고: 아직 주석이 지원되지 않습니다!

#include main:projectile

function update {
    # 함수 내부 주석은 사용 가능합니다.
    tp ~ ~1 ~
}

function checkhit {
    execute unless block ~ ~ ~ air run return 1
    return 0
}
```

```js
// main:bullet
// main/method/bullet.flm

#include main:projectile // projectile의 함수를 포함합니다.

function update {
    # projectile의 함수를 덮어씁니다.
    tp ^ ^ ^1
}

function hit {
    say 으악!
    kill @s
}
```

엔티티에 메서드셋을 부착하기 위해서는 `bind`를 사용해야 합니다.

`bind`의 문법은 다음과 같습니다:

```mcfunction
bind <선택자> <메서드셋 id>
```

예:

```mcfunction
bind @s main:bullet
```

> `bind`로 메서드셋을 부착한 모든 엔티티는 FLOW.METHODSET.> \<namespace\>.\<methodset 이름\> 구조의 태그를 가집니다.

부착된 함수를 호출하기 위해선 `do`를 사용해야 합니다.

`do`의 문법은 다음과 같습니다.

```mcfunction
do <함수 이름>
# 또는
do <선택자>.<함수 이름>
```

예:

```mcfunction
execute as @e[tag=bullet,type=marker] run do update
# 또는
do @e[tag=bullet,type=marker].update
```

bind와 do 활용 예시:

```mcfunction
# main:bullet, main:arrow, main:stone은 모두 update 메서드를 가집니다.
execute summon marker run bind @s main:bullet
execute summon marker run bind @s main:arrow
execute summon marker run bind @s main:stone

# 각 메서드셋에 맞는 update가 실행됩니다.
execute as @e[type=marker] at @s run do update
```

### 익명 함수

딱히 함수가 필요한 건 아니고 `execute`에 쓸 코드 블록이 필요할 때, 또는 return을 일부분에만 적용하고 싶을 때 사용할 수 있습니다.

`fn`을 통해 사용할 수 있습니다.

예:

```mcfunction
execute if entity @p[tag=something] run fn {
    say ㄹㄴㅇㄹㅇㄴ
    say fsdffds
    say Hello
}

fn {
    execute if score @p a matches 1 run return run say 성공
    say 실패
}

say 끝
```

fn {} 자리에 자동 생성된 익명 함수가 들어갑니다.

## \_\_this\_\_ 키워드

익명함수 내부에서 `__this__`를 사용하면 현재 함수를 다시 호출할 수 있습니다.  
익명함수에서 재귀를 구현하기 위해 구현되었습니다.

예:

```
fn {
    say 무한재귀
    function __this__
}
```
