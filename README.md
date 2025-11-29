# MC-Mono
C와 Go를 참고해서 만든 마인크래프트 데이터팩 생성을 위한 언어

를 컴파일하는 컴파일러를 만드는 저장소입니다

이곳에는 세 가지 정보가 있습니다:
1. [컴파일러 사용법](#사용법)
2. [MonoLang의 문법](#문법)
3. [Mono Compiler 구현 계획](#구현-계획)

`README.md`는 사실 기획서를 겸하고 있어서 설명이 친절하지 못합니다.

자세한 정보를 얻고싶으시다면 추후 작성될 문서를 읽으시는 걸 추천드립니다.

# 설치/사용법

> [!NOTE]
> Mono Compiler는 mnc로 축약되어 셸에서 사용됩니다.

## 설치

### 직접 빌드:
```
shell> git clone https://github.com/mang-jin/mc-mono
shell> cd mc-mono
shell> ./build.sh
shell> ./mnc --version
```

이후 생성된 바이너리 파일을 원하는 위치에 두고 사용하시면 됩니다.

> [!NOTE]
> build.sh파일은 `gcc(C 컴파일러)`를 사용하므로 컴퓨터에 `gcc`를 설치하셔야 합니다.

### 릴리즈 다운로드:
github 릴리즈에서 바이너리를 다운로드해서 사용하실 수 있습니다.

## Hello, World!

### 프로젝트 구조:
```
dtpk/
    src/
        main.mn
```

### dtpk/src/main.mn:
```
module "main";

import "cmd";

@load
void main() {
    cmd::tellraw("Hello, World!");
}
```

### 빌드:
```
shell> cd dtpk
shell> mnc build src/main.mn -o ./data
```

# 문법

## 모듈
```
module "<모듈이름>";
```
해당 구문은 현재 파일의 모듈을 지정합니다.

모듈 이름은 외부에서 해당 모듈의 함수를 호출할 때 사용됩니다.

예:
```
// sound.mn <- 주석입니다.
module "sound";

import "cmd";

void cat() {
    cmd::tellraw("Meow");
}

void dog() {
    cmd::tellraw("Bark!");
}

void cow() {
    cmd::tellraw("Mooo");
}

// main.mn
module "main";

import "sound";

@load
void main() {
    sound::cat(); // Meow
    sound::dog(); // Bark!
    sound::cow(); // Mooo
}
```

> [!NOTE]
> 여러 파일에서 중복된 모듈 이름을 사용할 수 있으며,
> 파일이 다르더라도 모듈 이름이 같으면 함수 이름이 중복될 수 없습니다.

## 함수

```
<반환타입> <함수이름>( <매개변수> ) {
    <함수내용>
}
```
함수를 선언합니다.

함수가 선언되면 컴파일 시 `<모듈이름>:<함수이름>` 위치에 저장됩니다.

## Attribute

```
@attribute
int main() {}
```
`main` 함수에 `@attribute`라는 `Attribute`를 추가합니다.

`Attribute`는 컴파일러에게 대상 함수에 대한 특별한 지시 사항을 전달하는 구문입니다.

사용 가능한 `Attribute`:
|Attribute|                    내용                     |
|---------|---------------------------------------------|
|@load    |데이터팩이 로드되었을때 실행하라             |
|@tick    |매 틱마다 실행하라                           |
|@raw     |함수 내부 내용을 mcfunction 문법으로 해석하라|

## 변수

`<타입> <변수 이름> = <값>;` 이 구문으로 변수를 선언할 수 있습니다.

`<변수 이름> = <값>;` 이런 형태로 변수를 수정할 수 있습니다.

변수는 스코어보드에 저장됩니다.


사용 가능한 타입:
|타입|         크기          |
|----|-----------------------|
|mo  |스코어보드 한 개(32bit)|
|di  |스코어보드 두 개(64bit)|

```
module "main";

void main() {
	mo a = 169; // -2,147,483,648 ~ 2,147,483,647
	di b = 170-1; // -9,223,372,036,854,775,808 ~ 9,223,372,036,854,775,807

    mo c[5] = [1,2,3]; // 배열을 선언합니다.
    mo d[] = [5,6,7]; // 컴파일러가 배열의 크기를 추측합니다.

    // `타입*`은 포인터(변수의 주소)를 뜻합니다.
    mo *ap = &a; // a의 주소를 저장합니다.
    di* bp = &b; // b의 주소를 저장합니다.
    // 별의 위치는 상관 없습니다.

    // 16진수는 10진수로 바뀌어 저장됩니다.
    mo hex = 0xffffff;
    // 2진수도 같습니다.
    mo bin = 0b1010;

    // 값은 수식이 될 수도 있습니다.
    mo op = 20+40*2; // 100
}
```

# 구현 계획