function myFunction() {
  const spreadsheet = SpreadsheetApp.openById(
    "1JM5Rm46MilqXl87CyKmSThKNPddwy9Ysy85HknDQxDo"
  );

  const settingSt = spreadsheet.getSheetByName("担当表");

  const email = Session.getActiveUser().getEmail();
  const mon = settingSt.getRange("B2").getValue();
  const tue = settingSt.getRange("C2").getValue();
  const wed = settingSt.getRange("D2").getValue();
  const thu = settingSt.getRange("E2").getValue();
  const fri = settingSt.getRange("F2").getValue();
  const items = [mon, tue, wed, thu, fri];

  const today = new Date();
  const y = today.getFullYear();
  const m = today.getMonth();
  const d = today.getDate();
  const w = today.getDay();
  const monday = d - w + 1; //今週の月曜日を設定
  const startDate = new Date(y, m, monday);

  const calendar = CalendarApp.getDefaultCalendar();
  const title = "イベント2";

  for (let i = 0; i < items.length; i++) {
    const target = items[i];

    if (target.includes(email)) {
      console.log(title, startDate);
      calendar.createAllDayEvent(title, startDate);
    }

    // 1日加算する
    startDate.setDate(startDate.getDate() + 1);
  }
}
