package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	PARAM_NAME_TEMPORABLE_TABLE_ID = "tmpid"
	PARAM_NAME_FILE_ID             = "fid"
	PARAM_NAME_TOKEN               = "token"

	PARAM_QUERY_FORMAT = "format"
	PARAM_QUERY_TYPE   = "rows"

	BLOCK_ROWS_NUMBER = 1000
)

func DetectDataFormat(fullpath string) (dataformat models.DataFormat, err error) {
	file, err := os.Open(fullpath)
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		return models.DATA_FORMAT_UKNOWN, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	firstline, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			return models.DATA_FORMAT_UKNOWN, err
		}
	}
	if len(strings.Split(firstline, "\t")) > 1 {
		return models.DATA_FORMAT_TXT, nil
	} else if len(strings.Split(firstline, ",")) > 1 {
		return models.DATA_FORMAT_CSV, nil
	} else if len(strings.Split(firstline, ";")) > 1 {
		return models.DATA_FORMAT_SSV, nil
	}

	return models.DATA_FORMAT_UKNOWN, nil
}

func ImportData(viewimporttable models.ViewImportTable, file *models.DtoFile, dtocustomertable *models.DtoCustomerTable,
	customertablerepository services.CustomerTableRepository, importsteprepository services.ImportStepRepository,
	columntyperepository services.ColumnTypeRepository) {
	dtoimportstep := models.NewDtoImportStep(dtocustomertable.ID, 2, false, 0, time.Now(), time.Now())
	err := importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}

	fullpath := filepath.Join(config.Configuration.Server.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID))
	dataformat, err := DetectDataFormat(fullpath)
	if err != nil {
		log.Error("Can't detect data format for file %v with value %v", err, fullpath)
		return
	}

	csvfile, err := os.Open(fullpath)
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	reader.Comma = models.GetDataSeparator(dataformat)
	reader.LazyQuotes = false
	reader.TrimLeadingSpace = true

	log.Info("Detected data format %v", dataformat)
	log.Info("Start reading file %v", time.Now())
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		return
	}
	rowcount := len(rawCSVdata)
	if viewimporttable.HasHeader && rowcount > 0 {
		rowcount -= 1
	}
	columncount := len(rawCSVdata[0])
	log.Info("Stop reading file %v row count %v column count %v", time.Now(), rowcount, columncount)

	if columncount == 0 {
		log.Error("Can't find any data in file %v", fullpath)
		return
	}
	if columncount > models.MAX_COLUMN_NUMBER {
		log.Error("So many columns are not supported %v", columncount)
		return
	}
	dtocolumntype, err := columntyperepository.Get(models.COLUMN_TYPE_DEFAULT)
	if err != nil {
		return
	}
	if !dtocolumntype.Active {
		log.Error("Column type is not active %v", dtocolumntype.ID)
		return
	}

	dtotablecolumns := new([]models.DtoTableColumn)
	position := 0
	for _, column := range rawCSVdata[0] {
		dtotablecolumn := new(models.DtoTableColumn)
		dtotablecolumn.Created = time.Now()
		dtotablecolumn.Position = int64(position)
		if viewimporttable.HasHeader {
			dtotablecolumn.Name = column
		} else {
			dtotablecolumn.Name = fmt.Sprintf("Column %v", position)
		}
		dtotablecolumn.Customer_Table_ID = dtocustomertable.ID
		dtotablecolumn.Column_Type_ID = models.COLUMN_TYPE_DEFAULT
		dtotablecolumn.Prebuilt = false
		dtotablecolumn.FieldNum = byte(position) + 1
		dtotablecolumn.Active = true
		dtotablecolumn.Edition = 0
		position++

		*dtotablecolumns = append(*dtotablecolumns, *dtotablecolumn)
	}

	err = customertablerepository.ImportDataStructure(dtotablecolumns, true)
	if err != nil {
		return
	}

	wrongrowcount := 0
	log.Info("Start analazying table data %v", time.Now())
	position = 0
	for i, _ := range rawCSVdata {
		if i == 0 && viewimporttable.HasHeader {
			rawCSVdata[i] = append(rawCSVdata[i], "position")
			rawCSVdata[i] = append(rawCSVdata[i], "wrong")
			continue
		}
		wrong := 0
		current_columncount := len(rawCSVdata[i])
		if current_columncount != columncount {
			wrong = 1
			wrongrowcount++
			if current_columncount < columncount {
				for j := current_columncount; j < columncount; j++ {
					rawCSVdata[i] = append(rawCSVdata[i], "")
				}
			} else {
				for j := columncount; j < current_columncount; j++ {
					rawCSVdata[i] = append(rawCSVdata[i][:columncount], rawCSVdata[i][columncount+1:]...)
				}
			}
		}
		rawCSVdata[i] = append(rawCSVdata[i], fmt.Sprintf("%v", position))
		rawCSVdata[i] = append(rawCSVdata[i], fmt.Sprintf("%v", wrong))
		position++
		if i%BLOCK_ROWS_NUMBER == 0 {
			log.Info("Continue analyzing table data %v at position %v", time.Now(), i)
		}
	}
	log.Info("Stop analyzing table data %v", time.Now())
	// 2
	dtoimportstep.Ready = true
	dtoimportstep.Percentage = 100
	dtoimportstep.Completed = time.Now()
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}

	dtocustomertable.Import_Percentage = 50
	dtocustomertable.Import_Columns = int64(columncount)
	dtocustomertable.Import_Rows = int64(rowcount)
	dtocustomertable.Import_WrongRows = int64(wrongrowcount)
	err = customertablerepository.Update(dtocustomertable)
	if err != nil {
		return
	}

	dtoimportstep.Step = 3
	dtoimportstep.Ready = false
	dtoimportstep.Percentage = 0
	dtoimportstep.Started = time.Now()
	dtoimportstep.Completed = time.Now()
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}
	version := fmt.Sprintf(".%v", time.Now().UTC().UnixNano())
	fullpath += version
	newcsvfile, err := os.Create(fullpath)
	if err != nil {
		log.Error("Can't write to file %v with value %v", err, fullpath)
		return
	}

	writer := csv.NewWriter(newcsvfile)
	writer.Comma = models.GetDataSeparator(dataformat)
	writer.UseCRLF = false

	log.Info("Start writing file %v", time.Now())
	err = writer.WriteAll(rawCSVdata)
	if err != nil {
		log.Error("Can't write to file %v with value %v", err, fullpath)
		newcsvfile.Close()
		return
	}
	log.Info("Stop writing file %v", time.Now())
	newcsvfile.Close()

	err = customertablerepository.ImportData(file, dtocustomertable, dataformat, viewimporttable.HasHeader, dtotablecolumns, version)
	if err != nil {
		os.Remove(fullpath)
		return
	}
	os.Remove(fullpath)
	// 3
	dtoimportstep.Ready = true
	dtoimportstep.Percentage = 100
	dtoimportstep.Completed = time.Now()
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}

	dtocustomertable.Import_Percentage = 75
	err = customertablerepository.Update(dtocustomertable)
	if err != nil {
		return
	}
}

func ExportData(viewexporttable *models.ViewExportTable, file *models.DtoFile, dtocustomertable *models.DtoCustomerTable,
	tablecolumns *[]models.DtoTableColumn, filerepository services.FileRepository,
	customertablerepository services.CustomerTableRepository) {
	version := fmt.Sprintf(".%v", time.Now().UTC().UnixNano())
	err := customertablerepository.ExportData(viewexporttable, file, dtocustomertable, tablecolumns, true, version)
	if err != nil {
		return
	}

	mvCmd := exec.Command("mv", "-f", filepath.Join(config.Configuration.TempDirectory, fmt.Sprintf("%08d", file.ID)+version),
		filepath.Join(config.Configuration.Server.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID)))
	err = mvCmd.Run()
	if err != nil {
		log.Error("Can't move temporable file %v with value %v%v", err, file.ID, version)
		return
	}

	file.Export_Ready = true
	file.Export_Percentage = 100
	err = filerepository.Update(file)
	if err != nil {
		return
	}
}

func CheckTableCells(dtocustomertable *models.DtoCustomerTable, tablecolumnrepository services.TableColumnRepository,
	columntyperepository services.ColumnTypeRepository, tablerowrepository services.TableRowRepository,
	importsteprepository services.ImportStepRepository) {

	tablecolumns, err := tablecolumnrepository.GetByTable(dtocustomertable.ID)
	if err != nil {
		return
	}
	columntypes, err := columntyperepository.GetByTable(dtocustomertable.ID)
	if err != nil {
		return
	}

	dtoimportstep := models.NewDtoImportStep(dtocustomertable.ID, 5, false, 0, time.Now(), time.Now())
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}

	regexps := make(map[int]*regexp.Regexp)
	for _, columntype := range columntypes {
		var r *regexp.Regexp = nil
		if columntype.Regexp != "" {
			r, err = regexp.Compile(columntype.Regexp)
			if err != nil {
				log.Error("Error during running reg exp %v with value %v", err, columntype.Regexp)
				return
			}
		}
		regexps[columntype.ID] = r
	}

	var offset int64 = 0
	var count int64 = BLOCK_ROWS_NUMBER
	log.Info("Start checking rows %v", time.Now())
	for {

		log.Info("Getting checking rows %v from %v to %v", time.Now(), offset, offset+count)
		dtotablerows, err := tablerowrepository.GetValidation(offset, count, dtocustomertable.ID, tablecolumns)
		if err != nil {
			return
		}
		if len(*dtotablerows) == 0 {
			break
		}
		log.Info("Validating checking rows %v from %v to %v", time.Now(), offset, offset+count)
		for i, _ := range *dtotablerows {
			for _, tablecolumn := range *tablecolumns {
				tablecell, err := (&(*dtotablerows)[i]).TableRowToDtoTableCell(&tablecolumn)
				if err != nil {
					return
				}
				columntype := columntypes[tablecolumn.Column_Type_ID]
				tablecell.Valid, _, _ = columntyperepository.Validate(&columntype, regexps[columntype.ID], tablecell.Value)
				err = (&(*dtotablerows)[i]).DtoTableCellToTableRow(tablecell, &tablecolumn)
				if err != nil {
					return
				}
			}
		}
		log.Info("Saving checking rows %v from %v to %v", time.Now(), offset, offset+count)
		err = tablerowrepository.SaveValidation(dtotablerows, tablecolumns)
		if err != nil {
			return
		}
		offset += count
	}
	log.Info("Stop checking rows %v", time.Now())
	// 5
	dtoimportstep.Ready = true
	dtoimportstep.Percentage = 100
	dtoimportstep.Completed = time.Now()
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}
}
