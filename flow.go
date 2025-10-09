// 그거 아세요?
// FUCK의 뜻은
// 'Failure’s Underlying Cause is lack of Knowledge'
// 즉, '실패의 근본 원인은 무지 때문이다' 라는 뜻입니다

package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ut "flow/Utils"
	cmdvh "flow/VersionHandler"
	"flow/omfn"
	anonyfunc "flow/omfn/Anonyfunc"
)

// func handleFlowCodes(namespace_path string) {
// 	os.RemoveAll(filepath.Join(namespace_path, "function"))
// 	os.MkdirAll(filepath.Join(namespace_path, "function"), os.ModePerm)
// 	fl_path := filepath.Join(namespace_path, "flow")
// 	// Namespace
// 	filesInNs, err := os.ReadDir(fl_path)
// 	if err != nil {
// 		// flow 폴더가 없음
// 		return
// 	}
// 	for _, entry := range filesInNs {
// 		name := entry.Name()
// 		if entry.IsDir() {
// 			handleFlowCodes(filepath.Join(fl_path, name))
// 		} else if ut.FileEx(name) == "fl" {
// 			filename := filepath.Join(namespace_path, "function", ut.FileName(name)) + ".mcfunction"
// 			new_compiled_omfn_file := ut.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm))
// 			defer new_compiled_omfn_file.Close()
// 			writer := bufio.NewWriter(new_compiled_omfn_file)

// 			fmt.Println(filepath.Join(fl_path, name))
// 			code := omfn.Parse(filepath.Join(fl_path, name))
// 			writer.WriteString(code)

// 			writer.Flush() // 버퍼 비우기
// 		}
// 	}
// }

// 함수 선언 폴더 읽는 놈
// func handleFunctions(fn_path string) {
// 	functions, err := os.ReadDir(fn_path)
// 	if err != nil {
// 		fmt.Println("\"function\" 폴더를 찾지 못함:\n", err)
// 		return
// 	}
// 	for _, entry := range functions {
// 		name := entry.Name()
// 		if entry.IsDir() {
// 			handleFunctions(filepath.Join(fn_path, name))
// 		} else if ut.FileEx(name) == "fl" {
// 			filename := filepath.Join(fn_path, ut.FileName(name)) + ".mcfunction"
// 			new_compiled_omfn_file := ut.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm))
// 			defer new_compiled_omfn_file.Close()
// 			writer := bufio.NewWriter(new_compiled_omfn_file)

// 			fmt.Println(filepath.Join(fn_path, name), "@@@")
// 			code := omfn.Parse(filepath.Join(fn_path, name))
// 			writer.WriteString(code)

// 			writer.Flush() // 버퍼 비우기
// 		}
// 	}
// }

/*
MethodSet의 Method와 MethodSet이 include중인 MethodSet의 Method를 모두 구해서 반환하는 함수
*/
//includes 파라미터는 함수 내부에서 재귀할 때 쓰는거니까 함수 호출할 땐 걍 빈 배열 쓰셈.
// func getMethodList(mset omce.MethodSet, includes []string) map[string]string {
// 	msetid := mset.Namespace + ":" + mset.Name
// 	for _, handsomeguy := range includes {
// 		if handsomeguy == msetid {
// 			FUCK := "include 순환이 감지되었습니다 "
// 			FUCK += strings.Join(includes, "->")
// 			panic(FUCK)
// 		}
// 	}
// 	includes = append(includes, msetid)
// 	result := make(map[string]string, 0)
// 	for _, inc := range mset.Includes {
// 		if inc == msetid {
// 			panic(fmt.Sprintf("%s에서 include 재귀가 감지되었습니다", msetid))
// 		}
// 		maps.Copy(result, getMethodList(declaredMethodSets[inc], includes))
// 	}
// 	maps.Copy(result, mset.Methods)
// 	return result
// }

// 커스텀 엔티티 선언 폴더 읽는 놈
// func declareMethodSets(ce_path string, namespace string) {
// 	custom_entities, err := os.ReadDir(ce_path)
// 	if err != nil {
// 		// fmt.Println("\"method\" 폴더를 찾지 못함:\n", err)
// 		return
// 	}
// 	for _, entry := range custom_entities {
// 		name := entry.Name()
// 		if entry.IsDir() {
// 			declareMethodSets(filepath.Join(ce_path, name), filepath.Join(namespace, name))
// 		} else if ut.FileEx(name) == "flm" {
// 			new_custom_entity := omce.Parse(filepath.Join(ce_path, name), namespace, ut.FileName(name))
// 			if strings.Contains(ut.FileName(name), ".") || strings.Contains(namespace, ".") {
// 				panic(fmt.Sprintf("%s:%s을(를) 컴파일 하던 중 오류 발생:\nflm의 디텍토리와 파일명은 '.'을 포함할 수 없습니다", strings.ReplaceAll(namespace, "\\", "/"), ut.FileName(name)))
// 			}
// 			declaredMethodSets[fmt.Sprintf("%s:%s", namespace, ut.FileName(name))] = new_custom_entity
// 		}
// 	}
// }

// func handleCustomEntities(ce_path string, namespace string) {
// 	declareMethodSets(ce_path, namespace)
// 	internal := filepath.Join(os.Getenv("INTERNAL_PATH"), "function")
// 	for _, mset := range declaredMethodSets {
// 		ut.Must(0, os.MkdirAll(filepath.Join(internal, "bind", mset.Namespace), os.ModePerm))
// 		filename := filepath.Join(internal, "bind", mset.Namespace, mset.Name) + ".mcfunction"
// 		new_methodset_file := ut.Must(os.Create(filename))
// 		defer new_methodset_file.Close()
// 		writer := bufio.NewWriter(new_methodset_file)

// 		method_list := getMethodList(mset, make([]string, 0))

// 		fmt.Fprintf(writer, "tag @s add FLOW.METHODSET.%s.%s\n", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name)

// 		for method_name := range method_list {
// 			// OMFN_INTERNAL.METHOD.<methodset name>.<method name>
// 			// 파일과 디렉토리명에 '.'포함 불가로 설정 해놔서 파일 명에 '.'들어가서 발생하는 인식 오류는 걱정할 필요 없다.
// 			new_cmd := fmt.Sprintf("tag @s add FLOW_INTERNAL.METHOD.%s.%s.%s\n", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name, method_name)
// 			writer.WriteString(new_cmd)
// 		}
// 		writer.Flush() // 버퍼 비우기

// 		for method_name, method_path := range method_list {
// 			filename = filepath.Join(internal, "do", method_name) + ".mcfunction"
// 			new_do_file := ut.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm))
// 			defer new_do_file.Close()
// 			do_writer := bufio.NewWriter(new_do_file)

// 			target_tag := fmt.Sprintf("FLOW_INTERNAL.METHOD.%s.%s.%s", strings.ReplaceAll(mset.Namespace, "\\", "."), mset.Name, method_name)
// 			new_cmd := fmt.Sprintf("execute if entity @s[tag=%s] run return run function %s\n", target_tag, method_path)
// 			do_writer.WriteString(new_cmd)
// 			do_writer.Flush() // 버퍼 비우기
// 		}
// 	}
// }

// var datapackRootDir string

//go:embed assets/**/* assets/**/**
var embeddedFiles embed.FS

func copyEmbeddedDir(dst string) error {
	return fs.WalkDir(embeddedFiles, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 대상 경로 계산
		relPath, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			// 디렉토리 생성
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		if d.Name() == "temp" {
			return nil
		}
		// 파일 복사
		srcFile, err := embeddedFiles.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

// var declaredMethodSets map[string]omce.MethodSet = make(map[string]omce.MethodSet)

/*
		os.Args[0]은 현재 실행된 파일의 위치임.
		go run main.go something
	 		   ^^^^^^^ 여기인듯?
*/
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic 발생:", r)
			fmt.Println("엔터를 누르면 종료됩니다...")
			fmt.Scanln() // 입력 대기
		}
	}()

	dir := ut.Must(os.Getwd()) // get working directory(current directory) 이거 안 하면 .exe 호출한 터미널이 아니라 .exe 파일 위치 기준으로 실행됨.

	inputFlag := flag.String("i", "", "입력 파일 위치")
	outputFlag := flag.String("o", "./", "출력 위치")
	flag.Parse()

	if *inputFlag == "" {
		panic("input flag(-i \"파일 위치\")가 필요합니다")
	}

	datapackRootDir := *outputFlag
	sourceFilePath := filepath.Join(dir, *inputFlag)

	fmt.Println("Options:")
	fmt.Println("Input:", *inputFlag)
	fmt.Println("Output:", *outputFlag)
	fmt.Println("=======================================")

	info, err := os.Stat(sourceFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("경로가 존재하지 않습니다")
		} else {
			fmt.Println("에러 발생:", err)
		}
		return
	}

	if info.IsDir() {
		panic("디렉토리는 Input이 될 수 없습니다")
	}
	if ut.FileEx(info.Name()) != "fl" {
		panic("Input은 '.fl' 확장자를 갖고 있어야 됩니다")
	}

	fmt.Printf("경고: 이 작업은 '%s'을(를) 수정합니다.\n", filepath.Join(dir, *outputFlag))
	fmt.Print("계속하시겠습니까? (y/N): ")

	// 2. 사용자 응답 받기
	rrr := bufio.NewReader(os.Stdin)
	response, _ := rrr.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" {
		fmt.Println("작업이 취소되었습니다.")
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	os.RemoveAll(filepath.Join(datapackRootDir, "data"))
	if err := copyEmbeddedDir(datapackRootDir); err != nil {
		panic(err)
	}

	cmdvh.SetVersion("1.21.5")

	if err := os.Setenv("DATA_PATH", filepath.Join(datapackRootDir, "data")); err != nil {
		panic(fmt.Sprintf("환경 변수 설정 실패: %s", err.Error()))
	}
	if err := os.Setenv("MAIN_NS", "main"); err != nil {
		panic(fmt.Sprintf("환경 변수 설정 실패: %s", err.Error()))
	}
	if err := os.Setenv("MAIN_PATH", filepath.Join(os.Getenv("DATA_PATH"), os.Getenv("MAIN_NS"))); err != nil {
		panic(fmt.Sprintf("환경 변수 설정 실패: %s", err.Error()))
	}
	if err := os.Setenv("INT_NS", "_flow_internal"); err != nil {
		panic(fmt.Sprintf("환경 변수 설정 실패: %s", err.Error()))
	}
	if err := os.Setenv("INT_PATH", filepath.Join(datapackRootDir, "data", os.Getenv("INT_NS"))); err != nil {
		panic(fmt.Sprintf("환경 변수 설정 실패: %s", err.Error()))
	}

	omfn.Parse(sourceFilePath)

	anonyfunc.CreateFiles()

	fmt.Printf("Press Enter to continue...")
	fmt.Scanln()

	// for _, a := range declaredMethodSets {
	// 	fmt.Println(a.Methods)
	// }
}

// go run . ../test_datapack
// rsrc -manifest manifest.xml -o resource.syso
