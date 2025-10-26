// 그거 아세요?
// FUCK의 뜻은
// 'Failure’s Underlying Cause is lack of Knowledge'
// 즉, '실패의 근본 원인은 무지 때문이다' 라는 뜻입니다

package main

import (
	"bufio"
	"archive/zip"
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

func unzipFlibTo(target string, dst string) error {
	zipReader, err := zip.OpenReader(target)
	if err != nil {
		return err
	}
	defer zipReader.Close()

//	for _, file := range zipReader.File {
//		if file.FileInfo().IsDir() {
			

	return nil
	// for _, entry := range entires {
	// 	if entry.IsDir() {
	// 		// Dir
	// 		err = os.MkDir(filepath.Join(target,entry.Name()),0755)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		err = copyDir(filepath.Join(target,entry.Name()),filepath.Join(dst,entry.Name()))
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		// File
	// 		data, err := os.ReadFile(filepath.Join(target,entry.Name()))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		os.WriteFile(filepath.Join(dst,entry.Name()),data,0644)
	// 	}
	// }
}

type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Join(*s,",")
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s,value)
	return nil
}

var libDirPath string = "/home/newt/My/Projects/flow/libs"
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

	var LinkedLibs StringSlice

	inputFlag := flag.String("i", "", "입력 파일 위치")
	outputFlag := flag.String("o", "./", "출력 위치")
	flag.Var(&LinkedLibs, "l", "연결할 mcfunction 디렉토리의 위치")
	flag.Parse()

	if *inputFlag == "" {
		panic("input flag(-i \"파일 위치\")가 필요합니다")
	}

	datapackRootDir, err := filepath.Abs(*outputFlag)
	if err != nil {
		panic(err)
	}
	sourceFilePath, err := filepath.Abs(*inputFlag)
	if err != nil {
		panic(err)
	}

	fmt.Println("Options:")
	fmt.Println("Input:", *inputFlag)
	fmt.Println("Output:", *outputFlag)
	fmt.Println("Links:")
	for i, linkFunc := range LinkedLibs {
		absFuncPath := filepath.Join(libDirPath,linkFunc)+".flib"
		if err != nil {
			panic(err)
		}
		LinkedLibs[i]=absFuncPath
		fmt.Println(" ",absFuncPath)
		_, err := os.Stat(absFuncPath)
		if err != nil {
			panic(err)
		}
	}
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
	info, err = os.Stat(datapackRootDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("경로가 존재하지 않습니다")
		} else {
			fmt.Println("에러 발생:", err)
		}
		return
	}

	if !info.IsDir() {
		panic("출력 위치는 디렉토리여야 합니다")
	}

	fmt.Printf("경고: 이 작업은 '%s'을(를) 수정합니다.\n", datapackRootDir)
	fmt.Print("계속하시겠습니까? (Y/n): ")

	// 2. 사용자 응답 받기
	rrr := bufio.NewReader(os.Stdin)
	response, _ := rrr.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "" {
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

	for _, flibZipFilePath := range LinkedLibs {
		zipReader, err := zip.OpenReader(flibZipFilePath)
		if err != nil {
			panic(err)
		}
		defer zipReader.Close()

		for _, mcfunctionFile := range zipReader.File {
			dir := filepath.Dir(mcfunctionFile.Name)
			if dir != "." {
				err = os.Mkdir(filepath.Join(os.Getenv("MAIN_PATH"),"function",dir),0755)
				if err != nil {
					fmt.Println("디렉토리 생성 실패: "+err.Error())
				}
			}
			originalFuncFile, err := mcfunctionFile.Open()
			if err != nil {
				panic(err)
			}
			defer originalFuncFile.Close()
			funcFileCodes, err := io.ReadAll(originalFuncFile)
			if err != nil {
				panic(err)
			}

			newFuncFile, err := os.Create(filepath.Join(os.Getenv("MAIN_PATH"),"function",mcfunctionFile.Name))
			if err != nil {
				panic(err)
			}
			defer newFuncFile.Close()

			newFuncFile.WriteString(string(funcFileCodes))
		}
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
