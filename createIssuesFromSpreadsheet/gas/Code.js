function doGet(e) {
  const COLUMN_NO_NUMBER = 1;
  const COLUMN_NO_MILESTONE = 2;
  const COLUMN_NO_TITLE = 3;
  const COLUMN_NO_BODY = 4;
  const COLUMN_NO_LABEL = 5;
  const COLUMN_NO_GAS_TOOL = 6;
  const COLUMN_NO_GAS_TOOL2 = 7;
  const MAX_COLUMN_NO = 8;

  const milestone = e.parameter.milestone;

  var spreadsheet = SpreadsheetApp.openById(
    PropertiesService.getScriptProperties().getProperty("spreadsheetID")
  );
  var st = spreadsheet.getSheetByName("スクリプト設定");

  const maxRow = st.getLastRow() - 1;
  const data = st.getRange(2, 1, maxRow, MAX_COLUMN_NO).getValues();

  const r1 = data.filter(
    (v) => String(v[COLUMN_NO_MILESTONE - 1]) == milestone
  );

  const res = [];

  for (var i = 0; i < r1.length; i++) {
    const val = r1[i];
    const id = val[COLUMN_NO_NUMBER - 1];
    const title = val[COLUMN_NO_TITLE - 1];
    const label = val[COLUMN_NO_LABEL - 1];
    const body = val[COLUMN_NO_BODY - 1];
    const gasTools = val[COLUMN_NO_GAS_TOOL - 1];
    const gasTools2 = val[COLUMN_NO_GAS_TOOL2 - 1];

    const repositories = [];
    if (gasTools) {
      repositories.push("gas-tools");
    }
    if (gasTools2) {
      repositories.push("gas-tools2");
    }

    const data = {
      id,
      title,
      body,
      milestone: val[COLUMN_NO_MILESTONE - 1],
      label,
      repositories,
    };

    res.push(data);
  }

  payload = JSON.stringify(res);
  ContentService.createTextOutput();
  var output = ContentService.createTextOutput();
  output.setMimeType(ContentService.MimeType.JSON);
  output.setContent(payload);

  return output;
}
