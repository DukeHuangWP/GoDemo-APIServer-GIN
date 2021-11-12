package controllers

import (
	"DevIntergTest/internal/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	URN_SettingsPath = "/SettingsPath"

	SettingHTMLUpper = `<!DOCTYPE html>
	<html xmlns="http://www.w3.org/1999/xhtml">
	
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>{{.HTMLValue_Title}}</title>
		<style>
			#iframeLeft {
				width: 20%;
			}
	
			#divRight {
				width: 79%;
				float: right;
			}
		</style>
	</head>
	
	<body>
		<div>
			<div id="divRight" onload="Javascript:JSF_SetCWinHeight()">
				<div> -- 最後維護: {{.HTMLValue_LastUpdate}}</div>
				<br>`

	SettingHTMLDown = `
	</div>
	</div>

</body>
<script language="javascript">
	function JSF_SetCWinHeight() {
		document.getElementById("iframeLeft").height = window.innerHeight - 10;
		document.getElementById("divRight").height = window.innerHeight - 10;

		document.getElementById('downloadlinks').insertAdjacentHTML('beforeend', '{{.HTMLValue_OutputDownloadLinks}}');
	}
</script>

</html>`
)

func LoadSettrings(ginCTX *gin.Context) {
	SettingPostFormat := `
	<form method="POST" action="" onsubmit="">
		<textarea style="width: %v; " rows="10" id="%v" name="%v">%v</textarea>
		<textarea style="width: %v; " rows="10" readonly="readonly">%v</textarea>
		<textarea style="width: %v; " rows="10" readonly="readonly">%v</textarea>
		<div></div>
		<input style="width: %v; " type="submit" formaction="%v" value="%v" />
		<br>
	</form>`

	TestCase := models.DoSQLs.GetAllTestCases()

	var Posts strings.Builder
	Posts.WriteString(SettingHTMLUpper)
	for _, value := range TestCase {
		Posts.WriteString(fmt.Sprintf(SettingPostFormat, "40%", value[0], value[0], value[3],
			"20%", value[7],
			"20%", value[8],
			"80%", URN_SettingsPath+"/"+value[0], value[0]))
	}
	Posts.WriteString(SettingHTMLDown)

	ginCTX.Writer.WriteHeader(http.StatusOK)
	ginCTX.Writer.Write([]byte(Posts.String()))

}

// func SaveSettrings(ginCTX *gin.Context) {
// 	testName := ginCTX.Param("testName")

// 	err := models.DoSQLs.UpdateTestCases(testName, ginCTX.PostForm(testName))
// 	if err != nil {
// 		ginCTX.String(http.StatusInternalServerError, err.Error())
// 	} else {
// 		LoadSettrings(ginCTX)
// 	}
// 	return
// }
