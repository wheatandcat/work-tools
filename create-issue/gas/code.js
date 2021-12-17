const SPREADSHEET_ID =
  PropertiesService.getScriptProperties().getProperty("SPREADSHEET_ID");
const SPREADSHEET_NAME =
  PropertiesService.getScriptProperties().getProperty("SPREADSHEET_NAME");
const CREATE_ISSUE_API =
  PropertiesService.getScriptProperties().getProperty("CREATE_ISSUE_API");
const SLACK_API =
  PropertiesService.getScriptProperties().getProperty("SLACK_API");
const VERIFY_ID_TOKEN =
  PropertiesService.getScriptProperties().getProperty("VERIFY_ID_TOKEN");

function doGet(e) {
  const no = e.parameter.no.split(",");
  createIssue(no);
}

const createIssue = async (no) => {
  var spreadsheet = SpreadsheetApp.openById(SPREADSHEET_ID);
  var st = spreadsheet.getSheetByName(SPREADSHEET_NAME);

  const maxRow = st.getLastRow() - 1;
  const maxCol = st.getMaxColumns() + 1;
  const richData = st.getRange(2, 1, maxRow, maxCol).getRichTextValues();
  const data = st.getRange(2, 1, maxRow, maxCol).getValues();
  const header = st.getRange(1, 1, 1, maxCol).getValues()[0];

  const COLUMN_NO_ID = header.indexOf("ID");
  const COLUMN_NO_VERSION = header.indexOf("バージョン");
  const COLUMN_NO_TITLE = header.indexOf("タイトル");
  const COLUMN_NO_BODY = header.indexOf("本文");
  const COLUMN_NO_PRIORITY = header.indexOf("優先度");
  const COLUMN_NO_ISSUES = header.indexOf("issues");
  const COLUMN_NO_BACKEND = header.indexOf("backend");
  const COLUMN_NO_FRONTEND = header.indexOf("frontend");
  const COLUMN_NO_IMAGE = header.indexOf("参考資料");
  const COLUMN_NO_ENVIRONMENT = header.indexOf("環境");

  const r1 = data.filter((v) => {
    const id = v[COLUMN_NO_ID];
    return no.includes(String(id));
  });

  const r2 = richData.filter((v, i) => {
    const id = data[i][COLUMN_NO_ID];
    return no.includes(String(id));
  });

  const res = [];

  for (var i = 0; i < r1.length; i++) {
    const val1 = r1[i];
    const val2 = r2[i];
    const backend = val1[COLUMN_NO_BACKEND];
    const frontend = val1[COLUMN_NO_FRONTEND];

    const repositories = [];
    if (backend) {
      repositories.push("gas-tools");
    }
    if (frontend) {
      repositories.push("gas-tools");
    }

    let image = String(val2[COLUMN_NO_IMAGE].getLinkUrl());
    if (image) {
      // リンクが無い時は文字列を取得
      image = val1[COLUMN_NO_IMAGE];
    }

    const r = {
      id: Number(val1[COLUMN_NO_ID]),
      priority: val1[COLUMN_NO_PRIORITY],
      title: val1[COLUMN_NO_TITLE],
      body: val1[COLUMN_NO_BODY],
      env: val1[COLUMN_NO_ENVIRONMENT],
      image: image,
      version: val1[COLUMN_NO_VERSION],
      repositories: repositories,
    };

    const item = await postCrateIssue(r);

    res.push(...item.issues);
  }

  for (var i = 0; i < res.length; i++) {
    const val = res[i];
    let issue = st
      .getRange(Number(val.id) + 1, COLUMN_NO_ISSUES + 1)
      .getValues()[0][0];

    if (issue) {
      issue += "\n";
    }
    st.getRange(Number(val.id) + 1, COLUMN_NO_ISSUES + 1).setValue(
      `${issue}${val.url}`
    );
  }

  var text = "以下のissueを作成しました。\n\n";
  for (var i = 0; i < res.length; i++) {
    const val = res[i];
    text += `■ ${val.title} \n`;
    text += `${val.url} \n`;
  }

  var url = SLACK_API;
  var params = {
    method: "post",
    contentType: "application/json",
    payload: '{"text":"' + text + '"}',
  };

  UrlFetchApp.fetch(url, params);

  return true;
};

const postCrateIssue = async (req) => {
  const option = {
    method: "POST",
    headers: {
      Authorization: VERIFY_ID_TOKEN,
    },
    contentType: "application/json",
    payload: JSON.stringify(req),
  };

  const res = await UrlFetchApp.fetch(CREATE_ISSUE_API, option);

  return JSON.parse(res.getContentText());
};
