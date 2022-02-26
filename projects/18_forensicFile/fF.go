// forensic file use signature file

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type signFile struct {
	Sign       string
	SuffixFile string
	FileFormat string
}

func main() {
	if len(os.Args) == 0 {
		fmt.Printf("Use: %v fileUnknown fileUnknown fileUnknown ...\n", os.Args[0])
	} else {
		massSign := []signFile{
			{`474946`, `*.gif`, `GIF files`},
			{`GIF89a`, `*.gif`, `GIF files`},
			{`FFD8FF`, `*.jpg`, `JPEG files`},
			{`JFIF`, `*.jpg`, `JPEG files`},
			{`504B03`, `*.zip`, `ZIP files`},
			{`25504446`, `*.pdf`, `PDF files`},
			{`%PDF`, `*.pdf`, `PDF files`},
			{`006E1EF0`, `*.ppt`, `PPT`},
			{`A0461DF0`, `*.ppt`, `PPT`},
			{`ECA5C100`, `*.doc`, `Doc file`},
			{`000100005374616E64617264204A6574204442`, `*.mdb`, `Microsoft database`},
			{`Standard Jet DB`, `*.mdb`, `Microsoft database`},
			{`2142444E`, `*.pst`, `PST file`},
			{`!BDN`, `*.pst`, `PST file`},
			{`0908100000060500`, `*.xls`, `XLS file`},
			{`D0CF11E0A1B11AE1`, `*.msi`, `MSI file`},
			{`D0CF11E0A1B11AE1`, `*.doc`, `DOC`},
			{`D0CF11E0A1B11AE1`, `*.xls`, `Excel`},
			{`D0CF11E0A1B11AE1`, `*.vsd`, `Visio`},
			{`D0CF11E0A1B11AE1`, `*.ppt`, `PPT`},
			{`0A2525454F460A`, `*.pdf`, `PDF file`},
			{`.%%EOF.`, `*.pdf`, `PDF file`},
			{`4040402000004040`, `*.hlp`, `HLP file`},
			{`465753`, `*.swf`, `SWF file`},
			{`FWS`, `*.swf`, `SWF file`},
			{`CWS`, `*.swf`, `SWF file`},
			{`494433`, `*.mp3`, `MP3 file`},
			{`ID3`, `*.mp3`, `MP3 file`},
			{`MSCF`, `*.cab`, `Cab file`},
			{`0x4D534346`, `*.cab`, `Cab file`},
			{`ITSF`, `*.chm`, `Compressed Help`},
			{`49545346`, `*.chm`, `Compressed Help`},
			{`4C00000001140200`, `*.lnk`, `Link file`},
			{`4C01`, `*.obj`, `OBJ file`},
			{`4D4D002A`, `*.tif`, `TIF graphics`},
			{`MM`, `*.tif`, `TIF graphics`},
			{`000000186674797033677035`, `*.mp4`, `MP4 Video`},
			{`ftyp3gp5`, `*.mp4`, `MP4 Video`},
			{`0x00000100`, `*.ico`, `Icon file`},
			{`300000004C664C65`, `*.evt`, `Event file`},
			{`LfLe`, `*.evt`, `Event file`},
			{`38425053`, `*.psd`, `Photoshop file`},
			{`8BPS`, `*.psd`, `Photoshop file`},
			{`4D5A`, `*.ocx`, `Active X`},
			{`4D6963726F736F66742056697375616C2053747564696F20536F6C7574696F6E2046696C65`, `*.sln`, `Microsft SLN file`},
			{`Microsoft Visual Studio Solution File`, `*.sln`, `Microsft SLN file`},
			{`504B030414000600`, `*.docx`, `Microsoft DOCX file`},
			{`504B030414000600`, `*.pptx`, `Microsoft PPTX file`},
			{`504B030414000600`, `*.xlsx`, `Microsoft XLSX file`},
			{`504B0304140008000800`, `*.xlsx`, `Java JAR file`},
			{`415649204C495354`, `*.avi`, `AVI file`},
			{`AVI LIST`, `*.avi`, `AVI file`},
			{`57415645666D7420`, `*.wav`, `WAV file`},
			{`WAVEfmt`, `*.wav`, `WAV file`},
			{`Rar!`, `*.rar`, `RAR file`},
			{`526172211A0700`, `*.rar`, `RAR file`},
			{`52657475726E2D506174683A20`, `*.eml`, `EML file`},
			{`Return-Path:`, `*.eml`, `EML file`},
			{`6D6F6F76`, `*.mov`, `MOV file`},
			{`moov`, `*.mov`, `MOV file`},
			{`7B5C72746631`, `*.rtf`, `RTF file`},
			{`{\rtf1`, `*.rtf`, `RTF file`},
			{`89504E470D0A1A0A`, `*.png`, `PNG file`},
			{`PNG`, `*.png`, `PNG file`},
			{`C5D0D3C6`, `*.eps`, `EPS file`},
			{`CAFEBABE`, `*.class`, `Java class file`},
			{`D7CDC69A`, `*.WMF`, `WMF file`},
		}

		for _, file := range os.Args[1:] {
			f, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "file: %v\n", err)
				continue
			}

			for _, val := range massSign {
				if strings.HasSuffix(file, val.SuffixFile) || bytes.Contains(f, []byte(val.Sign)) {
					fmt.Fprintf(os.Stdout, "%v may be %v\n", file, val.FileFormat)
					continue
				}
			}
		}
	}
}
