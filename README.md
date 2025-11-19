# MC-Mono
C와 Go를 참고해서 만든 마인크래프트 데이터팩 생성을 위한 언어

를 컴파일하는 컴파일러를 만드는 저장소입니다

이곳에는 세 가지 정보가 있습니다:
1. [컴파일러 사용법](#사용법)
2. [MonoLang의 문법](#문법)
3. [Mono Compiler 구현 계획](#구현-계획)

`README.md` 사실 기획서를 겸하고 있기 때문에 자세한 정보를 얻고싶으시다면
추후 작성될 문서를 읽으시는 걸 추천드립니다.

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
> build.sh파일은 go를 사용하므로 컴퓨터에 go를 설치하셔야 합니다.

### 릴리즈 다운로드:
github 릴리즈에서 바이너리를 다운로드해서 사용하실 수 있습니다.

## Hello, World!

### 프로젝트 구조:
```
dtpk/
    src/
        main.mn
    .prj
```

### dtpk/.prj:
```
input: "src";
output: "./";
dependencies: "cmd";
```

### dtpk/src/main.mn:
```
module "main"

@load
void main() {
    cmd::tellraw("Hello, World!");
}
```

### 빌드:
```
shell> cd dtpk
shell> mnc build
```

# 문법

# 구현 계획
