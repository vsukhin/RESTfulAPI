package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
)

const (
	PARAM_NAME_TEMPORABLE_TABLE_ID = "tmpid"
	PARAM_NAME_FILE_ID             = "fid"
	PARAM_NAME_TOKEN               = "token"

	PARAM_QUERY_FORMAT = "format"
	PARAM_QUERY_TYPE   = "rows"

	BLOCK_ROWS_NUMBER = 1000
)

func ConvertFileToDefault(dataencoding models.DataEncoding, fullpath string) (err error) {
	data, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return err
	}

	var array []byte
	var result string

	switch dataencoding {
	case models.DATA_ENCODING_WINDOWS1251:
		array, _, err = transform.Bytes(charmap.Windows1251.NewDecoder(), data)
		result = string(array)
	case models.DATA_ENCODING_KOI8R:
		array, _, err = transform.Bytes(charmap.KOI8R.NewDecoder(), data)
		result = string(array)
	case models.DATA_ENCODING_MACINTOSH:
		array, _, err = transform.Bytes(charmap.MacintoshCyrillic.NewDecoder(), data)
		result = string(array)
	case models.DATA_ENCODING_UTF16:
		if len(data)%2 != 0 {
			err = errors.New("Not Unicode 16")
			break
		}
		var unicode []uint16
		found := false
		defaultendian := true
		if len(data) > 1 {
			if data[0] == 0xfe && data[1] == 0xff {
				found = true
				defaultendian = false
			}
			if data[0] == 0xff && data[1] == 0xfe {
				found = true
				defaultendian = true
			}
		}
		pos := 0
		if found {
			pos = 2
		}
		for {
			if pos+2 > len(data) {
				break
			}
			var symbol uint16
			if defaultendian {
				symbol = binary.LittleEndian.Uint16(data[pos : pos+2])
			} else {
				symbol = binary.BigEndian.Uint16(data[pos : pos+2])
			}
			unicode = append(unicode, symbol)
			pos = pos + 2
		}
		result = string(utf16.Decode(unicode))
		err = nil
	case models.DATA_ENCODING_UNKNOWN:
		fallthrough
	case models.DATA_ENCODING_UTF8:
		result = string(data)
		err = nil
	default:
		err = errors.New("Unknown encoding")
	}
	if err != nil {
		return err
	}

	result = strings.Replace(result, "\r\n", "\n", -1)
	result = strings.Replace(result, "\r", "\n", -1)

	err = ioutil.WriteFile(fullpath, []byte(result), 0666)
	if err != nil {
		return err
	}

	return nil
}

func DetectDataFormat(fullpath string) (dataformat models.DataFormat, err error) {
	file, err := os.Open(fullpath)
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		return models.DATA_FORMAT_UNKNOWN, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	firstline, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Error("Can't detect data format for file %v with value %v", err, fullpath)
			return models.DATA_FORMAT_UNKNOWN, err
		}
	}
	if len(strings.Split(firstline, "\t")) > 1 {
		return models.DATA_FORMAT_TXT, nil
	} else if len(strings.Split(firstline, ",")) > 1 {
		return models.DATA_FORMAT_CSV, nil
	} else if len(strings.Split(firstline, ";")) > 1 {
		return models.DATA_FORMAT_SSV, nil
	}

	return models.DATA_FORMAT_UNKNOWN, nil
}

func SaveImportError(description string, dtocustomertable *models.DtoCustomerTable, customertablerepository services.CustomerTableRepository) {
	dtocustomertable.Import_Error = true
	dtocustomertable.Import_ErrorDescription = description
	err := customertablerepository.Update(dtocustomertable)
	if err != nil {
		log.Error("Can't save error information %v for table %v", err, dtocustomertable.ID)
		return
	}
}

func SaveExportError(description string, dtofile *models.DtoFile, filerepository services.FileRepository) {
	dtofile.Export_Error = true
	dtofile.Export_ErrorDescription = description
	err := filerepository.Update(dtofile)
	if err != nil {
		log.Error("Can't save error information %v for file %v", err, dtofile.ID)
		return
	}
}

func ImportData(viewimporttable models.ViewImportTable, file *models.DtoFile, dtocustomertable *models.DtoCustomerTable,
	customertablerepository services.CustomerTableRepository, importsteprepository services.ImportStepRepository,
	columntyperepository services.ColumnTypeRepository, language string) {
	dtoimportstep := models.NewDtoImportStep(dtocustomertable.ID, 2, false, 0, time.Now(), time.Now())
	err := importsteprepository.Save(dtoimportstep)
	if err != nil {
		return
	}

	fullpath := filepath.Join(config.Configuration.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID))
	var dataformat models.DataFormat
	if viewimporttable.DataFormat == models.DATA_FORMAT_UNKNOWN {
		dataformat, err = DetectDataFormat(fullpath)
		if err != nil {
			SaveImportError(config.Localization[language].Errors.Internal.Data_Format, dtocustomertable, customertablerepository)
			return
		}
	} else {
		dataformat = viewimporttable.DataFormat
	}

	var dataencoding models.DataEncoding
	if viewimporttable.DataEncoding == models.DATA_ENCODING_UNKNOWN {
		dataencoding = models.DATA_ENCODING_UTF8
	} else {
		dataencoding = viewimporttable.DataEncoding
	}

	err = ConvertFileToDefault(dataencoding, fullpath)
	if err != nil {
		log.Error("Can't convert file %v to encoding %v with value %v", err, dataencoding, fullpath)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Format, dtocustomertable, customertablerepository)
		return
	}

	csvfile, err := os.Open(fullpath)
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Reading, dtocustomertable, customertablerepository)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	reader.Comma = models.GetDataSeparator(dataformat)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	log.Info("Detected data format %v", dataformat)
	log.Info("Start reading file %v", time.Now())
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Error("Can't read from file %v with value %v", err, fullpath)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Reading, dtocustomertable, customertablerepository)
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
		SaveImportError(config.Localization[language].Errors.Internal.Data_Columns, dtocustomertable, customertablerepository)
		return
	}
	if columncount > models.MAX_COLUMN_NUMBER {
		log.Error("So many columns are not supported %v", columncount)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Columns, dtocustomertable, customertablerepository)
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
		SaveImportError(config.Localization[language].Errors.Internal.Data_Writing, dtocustomertable, customertablerepository)
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
		SaveImportError(config.Localization[language].Errors.Internal.Data_Writing, dtocustomertable, customertablerepository)
		return
	}

	writer := csv.NewWriter(newcsvfile)
	writer.Comma = models.GetDataSeparator(dataformat)
	writer.UseCRLF = false

	log.Info("Start writing file %v", time.Now())
	err = writer.WriteAll(rawCSVdata)
	if err != nil {
		newcsvfile.Close()
		log.Error("Can't write to file %v with value %v", err, fullpath)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Writing, dtocustomertable, customertablerepository)
		return
	}
	log.Info("Stop writing file %v", time.Now())
	newcsvfile.Close()

	err = customertablerepository.ImportData(file, dtocustomertable, dataformat, viewimporttable.HasHeader, dtotablecolumns, version)
	if err != nil {
		os.Remove(fullpath)
		SaveImportError(config.Localization[language].Errors.Internal.Data_Writing, dtocustomertable, customertablerepository)
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
	tablecolumns *[]models.DtoTableColumn, filerepository services.FileRepository, customertablerepository services.CustomerTableRepository, language string) {
	version := fmt.Sprintf(".%v", time.Now().UTC().UnixNano())
	err := customertablerepository.ExportData(viewexporttable, file, dtocustomertable, tablecolumns, true, version)
	if err != nil {
		SaveExportError(config.Localization[language].Errors.Internal.Data_Writing, file, filerepository)
		return
	}

	abstemppath, err := filepath.Abs(config.Configuration.TempDirectory)
	if err != nil {
		log.Error("Can't make an absolute path for %v, %v", config.Configuration.TempDirectory, err)
		SaveExportError(config.Localization[language].Errors.Internal.Data_Reading, file, filerepository)
		return
	}
	absfilepath, err := filepath.Abs(config.Configuration.FileStorage)
	if err != nil {
		log.Error("Can't make an absolute path for %v, %v", config.Configuration.FileStorage, err)
		SaveExportError(config.Localization[language].Errors.Internal.Data_Reading, file, filerepository)
		return
	}

	mvCmd := exec.Command("mv", "-f", filepath.Join(abstemppath, fmt.Sprintf("%08d", file.ID)+version),
		filepath.Join(absfilepath, file.Path, fmt.Sprintf("%08d", file.ID)))
	err = mvCmd.Run()
	if err != nil {
		log.Error("Can't move temporable file %v with value %v%v", err, file.ID, version)
		SaveExportError(config.Localization[language].Errors.Internal.Data_Writing, file, filerepository)
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
