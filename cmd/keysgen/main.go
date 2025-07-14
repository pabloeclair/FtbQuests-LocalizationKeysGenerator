package main

import (
	"context"
	"errors"
	"fmt"
	"keysgen/internal/kg"
	"keysgen/internal/utils"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {

	// Validation
	_, err := os.Stat("ftbquests")
	if os.IsNotExist(err) {
		fmt.Printf("ERROR: please take the `ftbquests` directory from the `.minecraft/config` path and add it to the root of the repository.\n")
		return
	}

	files, err := os.ReadDir(filepath.Join("ftbquests", "quests", "chapters"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("ERROR: please take the `ftbquests` directory from the `.minecraft/config` path and add it to the root of the repository.\n")
			fmt.Println("Are you sure you take the `ftbquests` directory from the `.minecraft/config` path?")
			fmt.Println("`ftbquests` directory from the `.minecraft/congih` path must have the following structure:")
			fmt.Println(`
				ftbquests
					└── quests
						└── chapters
							├── <chapter_name1>
							└── <chapter_name2>
			`)
			return
		} else {
			fmt.Printf("ERROR: reading directory error: %v\n", err)
		}
	}

	ctx, cancelGo := context.WithCancel(context.Background())
	ctx, cancelMain := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancelMain()
	go func() {
		defer cancelGo()
		// Input
		fmt.Println("Please enter the name of modpack: ")
		var modpackName string
		if _, err = fmt.Scan(&modpackName); err != nil {
			if !errors.Is(err, os.ErrExist) {
				fmt.Printf("ERROR: reading input error: %v", err)
				return
			}
		}

		fmt.Println("Please enter the code of the original language (for example, `en_us`): ")
		var originalLang string
		if _, err = fmt.Scan(&originalLang); err != nil {
			if !errors.Is(err, os.ErrExist) {
				fmt.Printf("ERROR: reading input error: %v", err)
				return
			}
		}

		// fmt.Println("Please enter the code of the translation language (for example, `ru_ru`): ")
		// var translationLang string
		// if _, err = fmt.Scan(&translationLang); err != nil {
		// 	if !errors.Is(err, os.ErrExist) {
		// 		fmt.Printf("ERROR: reading input error: %v", err)
		// 		return
		// 	}
		// }

		questsMap := map[string][]kg.Quest{}
		keysMap := map[string]string{}

		// Parsing
		fmt.Println("Parsing quests...")
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if filepath.Ext(file.Name()) == ".snbt" {
				quests, keys, err := kg.GenerateQuestsAndKeys(modpackName, file.Name())
				if err != nil {
					fmt.Println("ERROR:", err)
					continue
				}
				questsMap[file.Name()] = quests
				keysMap[file.Name()] = keys
			}
			fmt.Println(file.Name() + " done.")
		}

		// Create output directories
		outputFilePath := filepath.Join("output", "ftbquests", "quests", "chapters")
		if err = os.MkdirAll(outputFilePath, 0755); err != nil {
			if !errors.Is(err, os.ErrExist) {
				fmt.Printf("ERROR: creating %s directory error: %v\n", outputFilePath, err)
				return
			}
		}

		langFilePath := filepath.Join("output", "ftbquests", "lang")
		if err = os.Mkdir(langFilePath, 0755); err != nil {
			if !errors.Is(err, os.ErrExist) {
				fmt.Printf("ERROR: creating %s directory error: %v\n", langFilePath, err)
				return
			}
		}

		// Write keys
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			filePath := filepath.Join(outputFilePath, file.Name())
			if err = utils.CreateWriteFile(filePath, keysMap[file.Name()]); err != nil {
				fmt.Println("ERROR:", err)
				return
			}
		}

		// Write translation map
		fmt.Println("Generating translation map...")
		originalMap, err := kg.GenerateMap(originalLang, questsMap)
		if err != nil {
			fmt.Printf("ERROR: generating original map error: %v\n", err)
			return
		}
		if err = utils.CreateWriteFile(filepath.Join(langFilePath, originalLang+".json"), originalMap); err != nil {
			fmt.Println("ERROR:", err)
			return
		}
	}()

	<-ctx.Done()
	fmt.Println("Done!")

	// todo: add google translate
}
