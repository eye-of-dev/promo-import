package services

import (
	"PromocodesUpload/internal/models"
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"strings"
)

type ImportPromoCodes struct {
	dirName    string
	packPromo  int
	connection *gorm.DB
	log        *slog.Logger
}

func NewImportPromoCodes(dirName string, packPromo int, connection *gorm.DB, log *slog.Logger) *ImportPromoCodes {
	return &ImportPromoCodes{dirName, packPromo, connection, log}
}

func (i *ImportPromoCodes) Execute() {
	entries, err := i.readDir()
	if err != nil {
		i.log.Error(err.Error())
		return
	}

	for _, e := range entries {
		if strings.Contains(e.Name(), "imported") ||
			strings.Contains(e.Name(), "error") {
			continue
		}

		_, err := i.parseAndStoreData(e.Name())
		if err != nil {
			i.log.Error(err.Error())
			i.errorImportFile(e.Name())
		} else {
			i.successImportFile(e.Name())
		}
	}
}

func (i *ImportPromoCodes) readDir() ([]os.DirEntry, error) {
	entries, err := os.ReadDir(i.dirName)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q directory", i.dirName)
	}

	return entries, nil
}

func (i *ImportPromoCodes) parseAndStoreData(fileName string) (bool, error) {
	readFile, err := os.Open(fmt.Sprintf("%s/%s", i.dirName, fileName))
	defer readFile.Close()
	if err != nil {
		return false, fmt.Errorf("cannot read %q file", fileName)
	}

	var promoCodes []*models.PromoCodes
	countPromoCodes := 0
	countSavedPromoCodes := 0
	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		if len(fileScanner.Text()) > 0 {
			countPromoCodes++
			promoCodes = append(promoCodes, &models.PromoCodes{PromoCode: fileScanner.Text()})
		}

		if len(promoCodes) >= i.packPromo {
			if _, err := i.savePromoCodes(promoCodes); err != nil {
				return false, err
			}
			countSavedPromoCodes = countSavedPromoCodes + len(promoCodes)
			promoCodes = []*models.PromoCodes{}
		}
	}

	if len(promoCodes) >= 1 {
		if _, err := i.savePromoCodes(promoCodes); err != nil {
			return false, err
		}
		countSavedPromoCodes = countSavedPromoCodes + len(promoCodes)
	}

	if countSavedPromoCodes == countPromoCodes {
		i.log.Info(fmt.Sprintf("%s - %d", fileName, countPromoCodes))
		return true, nil
	}

	return false, fmt.Errorf("cannot saved codes from %q file", fileName)
}

func (i *ImportPromoCodes) successImportFile(fileName string) {
	_ = os.Rename(
		fmt.Sprintf("%s/%s", i.dirName, fileName),
		fmt.Sprintf("%s/%s.imported", i.dirName, fileName),
	)
}

func (i *ImportPromoCodes) errorImportFile(fileName string) {
	_ = os.Rename(
		fmt.Sprintf("%s/%s", i.dirName, fileName),
		fmt.Sprintf("%s/%s.error", i.dirName, fileName),
	)
}

func (i *ImportPromoCodes) savePromoCodes(promoCodes []*models.PromoCodes) (int64, error) {
	result := i.connection.Create(promoCodes)
	return result.RowsAffected, result.Error
}
