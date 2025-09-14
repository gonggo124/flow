// 그거 아세요?
// FUCK의 뜻은
// 'Failure’s Underlying Cause is Knowledge gap'
// 즉, '실패의 근본 원인은 무지 때문이다' 라는 뜻입니다

package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	ut "github.com/gonggo124/objective-mcfunction/Utils"
	cmdvh "github.com/gonggo124/objective-mcfunction/VersionHandler"
	"github.com/gonggo124/objective-mcfunction/omce"
	"github.com/gonggo124/objective-mcfunction/omfn"
	anonyfunc "github.com/gonggo124/objective-mcfunction/omfn/Anonyfunc"
)

// 함수 선언 폴더 읽는 놈
func handleFunctions(fn_path string) {
	functions, err := os.ReadDir(fn_path)
	if err != nil {
		fmt.Println("\"function\" 폴더를 찾지 못함:\n", err)
		return
	}
	for _, entry := range functions {
		name := entry.Name()
		if entry.IsDir() {
			handleFunctions(filepath.Join(fn_path, name))
		} else if ut.FileEx(name) == "omfn" {
			filename := filepath.Join(fn_path, ut.FileName(name)) + ".mcfunction"
			new_compiled_omfn_file := ut.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm))
			defer new_compiled_omfn_file.Close()
			writer := bufio.NewWriter(new_compiled_omfn_file)

			fmt.Println(filepath.Join(fn_path, name), "@@@")
			code := omfn.Parse(filepath.Join(fn_path, name))
			writer.WriteString(code)

			writer.Flush() // 버퍼 비우기
		}
	}
}

/*
MethodSet의 Method와 MethodSet이 include중인 MethodSet의 Method를 모두 구해서 반환하는 함수
*/
//includes 파라미터는 함수 내부에서 재귀할 때 쓰는거니까 함수 호출할 땐 걍 빈 배열 쓰셈.
func getMethodList(mset omce.MethodSet, includes []string) map[string]string {
	msetid := mset.Namespace + ":" + mset.Name
	for _, handsomeguy := range includes {
		if handsomeguy == msetid {
			FUCK := "include 순환이 감지되었습니다 "
			FUCK += strings.Join(includes, "->")
			panic(FUCK)
		}
	}
	includes = append(includes, msetid)
	result := make(map[string]string, 0)
	for _, inc := range mset.Includes {
		if inc == msetid {
			panic(fmt.Sprintf("%s에서 include 재귀가 감지되었습니다", msetid))
		}
		maps.Copy(result, getMethodList(declaredMethodSets[inc], includes))
	}
	maps.Copy(result, mset.Methods)
	return result
}

// 커스텀 엔티티 선언 폴더 읽는 놈
func declareMethodSets(ce_path string, namespace string) {
	custom_entities, err := os.ReadDir(ce_path)
	if err != nil {
		fmt.Println("\"method\" 폴더를 찾지 못함:\n", err)
		return
	}
	for _, entry := range custom_entities {
		name := entry.Name()
		if entry.IsDir() {
			declareMethodSets(filepath.Join(ce_path, name), filepath.Join(namespace, name))
		} else if ut.FileEx(name) == "omset" {
			new_custom_entity := omce.Parse(filepath.Join(ce_path, name), namespace, ut.FileName(name))
			if strings.Contains(ut.FileName(name), ".") || strings.Contains(namespace, ".") {
				panic(fmt.Sprintf("%s:%s을(를) 컴파일 하던 중 오류 발생:\nomset의 디텍토리와 파일명은 '.'을 포함할 수 없습니다", strings.ReplaceAll(namespace, "\\", "/"), ut.FileName(name)))
			}
			declaredMethodSets[fmt.Sprintf("%s:%s", namespace, ut.FileName(name))] = new_custom_entity
		}
	}
}
func handleCustomEntities(ce_path string, namespace string) {
	declareMethodSets(ce_path, namespace)
	internal := filepath.Join(os.Getenv("INTERNAL_PATH"), "function")
	for _, mset := range declaredMethodSets {
		ut.Must(0, os.MkdirAll(filepath.Join(internal, "bind", mset.Namespace), os.ModePerm))
		filename := filepath.Join(internal, "bind", mset.Namespace, mset.Name) + ".mcfunction"
		new_methodset_file := ut.Must(os.Create(filename))
		defer new_methodset_file.Close()
		writer := bufio.NewWriter(new_methodset_file)

		method_list := getMethodList(mset, make([]string, 0))

		fmt.Fprintf(writer, "tag @s add OMFN.METHOD.%s.%s\n", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name)

		for method_name := range method_list {
			// OMFN_INTERNAL.METHOD.<methodset name>.<method name>
			// 파일과 디렉토리명에 '.'포함 불가로 설정 해놔서 파일 명에 '.'들어가서 발생하는 인식 오류는 걱정할 필요 없다.
			new_cmd := fmt.Sprintf("tag @s add OMFN_INTERNAL.METHOD.%s.%s.%s\n", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name, method_name)
			writer.WriteString(new_cmd)
		}
		writer.Flush() // 버퍼 비우기

		for method_name, method_path := range method_list {
			filename = filepath.Join(internal, "do", method_name) + ".mcfunction"
			new_do_file := ut.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm))
			defer new_do_file.Close()
			do_writer := bufio.NewWriter(new_do_file)

			target_tag := fmt.Sprintf("OMFN_INTERNAL.METHOD.%s.%s.%s", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name, method_name)
			new_cmd := fmt.Sprintf("execute if entity @s[tag=%s] run return run function %s\n", target_tag, method_path)
			do_writer.WriteString(new_cmd)
			do_writer.Flush() // 버퍼 비우기
		}
	}
}

var datapackRootDir string
var declaredMethodSets map[string]omce.MethodSet = make(map[string]omce.MethodSet)

/*
		os.Args[0]은 현재 실행된 파일의 위치임.
		go run main.go something
	 		   ^^^^^^^ 여기인듯?
*/
func main() {
	fmt.Println("경고: 이 작업은 시스템 파일을 수정합니다.")
	fmt.Print("계속하시겠습니까? (Y/N): ")

	// 2. 사용자 응답 받기
	rrr := bufio.NewReader(os.Stdin)
	response, _ := rrr.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" {
		fmt.Println("작업이 취소되었습니다.")
		return
	}

	fmt.Print("Given args: ( ")
	for i, a := range os.Args {
		if i == 0 {
			continue
		}
		fmt.Printf("%s ", a)
	}
	fmt.Print(")\n")
	if len(os.Args) == 1 {
		fmt.Println("\"omfn [대상 데이터팩 디렉토리]\" 형식에 부합하지 않습니다.")
		return
	}

	cmdvh.SetVersion("1.21.5")

	dir := ut.Must(os.Getwd()) // get working directory(current directory) 이거 안 하면 .exe 호출한 터미널이 아니라 .exe 파일 위치 기준으로 실행됨.
	datapackRootDir = filepath.Join(dir, os.Args[1])
	namespaces := ut.Must(os.ReadDir(filepath.Join(datapackRootDir, "data")))

	if err := os.Setenv("INTERNAL_NAMESPACE", "_omfn_internal"); err != nil {
		fmt.Println("환경 변수 설정 실패:", err)
		return
	}
	if err := os.Setenv("INTERNAL_PATH", filepath.Join(datapackRootDir, "data", "_omfn_internal")); err != nil {
		fmt.Println("환경 변수 설정 실패:", err)
		return
	}
	internal_path := os.Getenv("INTERNAL_PATH")
	{
		entries, err := os.ReadDir(internal_path)
		if err == nil {
			for _, entry := range entries {
				path := filepath.Join(internal_path, entry.Name())
				err := os.RemoveAll(path) // 파일/디렉토리 모두 삭제
				if err != nil {
					fmt.Println("삭제 실패:", err)
				}
			}
		}
	}
	err := os.MkdirAll(filepath.Join(internal_path, "function"), os.ModePerm)
	if err != nil {
		fmt.Println("omfn 내부 함수 디렉토리 생성 실패")
		return
	}
	os.MkdirAll(filepath.Join(internal_path, "function", "bind"), os.ModePerm)
	os.MkdirAll(filepath.Join(internal_path, "function", "do"), os.ModePerm)
	os.MkdirAll(filepath.Join(internal_path, "function", "anony"), os.ModePerm)

	for _, namespace_dir := range namespaces {
		if namespace_dir.Name() == "_omfn_internal" {
			continue
		}
		handleFunctions(filepath.Join(datapackRootDir, "data", namespace_dir.Name(), "function"))
		handleCustomEntities(filepath.Join(datapackRootDir, "data", namespace_dir.Name(), "method"), namespace_dir.Name())
	}
	anonyfunc.CreateFiles()

	// for _, a := range declaredMethodSets {
	// 	fmt.Println(a.Methods)
	// }
}

// go run . ../test_datapack
