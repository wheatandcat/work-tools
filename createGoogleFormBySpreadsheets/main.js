const SPREADSHEET_ID =
  PropertiesService.getScriptProperties().getProperty("SPREADSHEET_ID");

const createFormApp = (data, teams, option) => {
  const values = data
    .filter((v) => {
      const backend = v[option.selectIndex];
      return backend === true;
    })
    .map((v) => v[option.nameIndex]);

  const day = Utilities.formatDate(new Date(), "JST", "YYYY年MM月");
  const title = `${day}【${option.team}】キャリアパスのアンケート`;
  var form = FormApp.create(title);
  form.setDescription("説明");

  form
    .addCheckboxItem()
    .setTitle("以下の中から興味のある内容を最大3つまで選択して下さい。")
    .setChoiceValues((uniqueArr = [...new Set(values)]))
    .showOtherOption(true)
    .setRequired(true);

  form
    .addCheckboxItem()
    .setTitle(
      "所属しているチーム以外で対応してみたい内容が含まれているチームがあれば選択して下さい（※ない場合は選択なしで大丈夫です）"
    )
    .setChoiceValues(
      (uniqueArr = [...new Set(teams.filter((v) => v !== option.team))])
    );

  form
    .addParagraphTextItem()
    .setTitle(
      "キャリアパスで記載したい内容があればお願いします。（記載内容は何でも大丈夫です）"
    );
};

const getTeams = () => {
  const spreadsheet = SpreadsheetApp.openById(SPREADSHEET_ID);
  const st = spreadsheet.getSheetByName("チーム");

  const maxRow = st.getLastRow() - 1;
  const maxCol = st.getMaxColumns() + 1;
  const data = st.getRange(2, 1, maxRow, maxCol).getValues();
  const header = st.getRange(1, 1, 1, maxCol).getValues()[0];

  const COLUMN_NO_TEAM = header.indexOf("チーム名");
  const teams = data.map((v) => v[COLUMN_NO_TEAM]);

  return teams;
};

function myFunction() {
  const spreadsheet = SpreadsheetApp.openById(SPREADSHEET_ID);
  const st = spreadsheet.getSheetByName("キャリアパス");

  const maxRow = st.getLastRow() - 1;
  const maxCol = st.getMaxColumns() + 1;
  const data = st.getRange(2, 1, maxRow, maxCol).getValues();
  const header = st.getRange(1, 1, 1, maxCol).getValues()[0];

  const teams = getTeams();

  const COLUMN_NO_BACKEND = header.indexOf("backend");
  const COLUMN_NO_FRONTEND = header.indexOf("frontend");
  const COLUMN_NO_NAME = header.indexOf("選択");

  createFormApp(data, teams, {
    team: "backendチーム",
    selectIndex: COLUMN_NO_BACKEND,
    nameIndex: COLUMN_NO_NAME,
  });

  createFormApp(data, teams, {
    team: "frontendチーム",
    selectIndex: COLUMN_NO_FRONTEND,
    nameIndex: COLUMN_NO_NAME,
  });
}
