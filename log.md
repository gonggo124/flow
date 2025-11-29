2025.11.23
- 입력 폴더와 출력 폴더 입력 구현.
- 입력 폴더 탐색 구현.

2025.11.26
- Lexer -> Tokenizer로 이름 변경.
- '\n' 단위 파일 탐색 -> char 단위 탐색으로 변경.

2025.11.27
- 토큰 추가.
- 에러 반환 시스템 추가.
- Tok_scan이 buf에 fgetc 결과를 밀어넣게 함.
- TOK_offset_t, TOK_act_c 등 타입 규격 추가.
- 상태 머신 방식으로 개발 계획 중 -> states, acts 추가.
- Tok_State 추가. header로 옮겨야 할듯? 그리고 Tok에 state field도 추가해야함.

2025.11.29
- 상태머신 일부 구현.
- Tok에 상태 추가.
- 상태 표 일부 정의.
- 변수 이름 일부 수정. Tok_State -> State(typedef는 tok.c 내부에서만 인식되기 때문에 필요 없을 듯?) 등.
- header로 옮기지 않았음. 필요없어서.
