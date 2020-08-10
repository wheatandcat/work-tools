function myFunction() {
  const COLUMN_NO_NUMBER = 1;
  const COLUMN_NO_MILESTONE = 2;
  const COLUMN_NO_TITLE = 3;
  const COLUMN_NO_FINISHED = 4;
  const MAX_COLUMN_NO = 4;
  const checkMilestone = "v2.0.0";

  var spreadsheet = SpreadsheetApp.openById(
    PropertiesService.getScriptProperties().getProperty("spreadsheetID")
  );
  var settingSt = spreadsheet.getSheetByName("スクリプト設定");
  var startDate = settingSt
    .getRange("B1")
    .getValue()
    .getTime();
  var endDate = settingSt
    .getRange("B2")
    .getValue()
    .getTime();
  var today = new Date().getTime();

  if (today >= startDate && today <= endDate) {
    var st = spreadsheet.getSheetByName("テスト");

    const maxrow = st.getLastRow() - 1;
    const data = st.getRange(2, 1, maxrow, MAX_COLUMN_NO).getValues();

    const r1 = data.filter(
      v => String(v[COLUMN_NO_MILESTONE - 1]) == checkMilestone
    );
    // 終了していないものだけを抽出
    const r2 = r1.filter(v => String(v[COLUMN_NO_FINISHED - 1]) == "");

    var targetText = "";

    for (var i = 0; i < r2.length; i++) {
      const val = r2[i];
      const no = val[COLUMN_NO_NUMBER - 1];
      const title = val[COLUMN_NO_TITLE - 1];

      const t = "No." + no + ": " + title + "\n";
      targetText += t;
    }

    const text =
      "Milestone:" +
      checkMilestone +
      "で以下の作業が完了していません。\n\n" +
      targetText;

    var slackURL = PropertiesService.getScriptProperties().getProperty(
      "slackURL"
    );

    var params = {
      method: "post",
      contentType: "application/json",
      payload: '{"text":"' + text + '"}'
    };

    UrlFetchApp.fetch(slackURL, params);
  }
}
