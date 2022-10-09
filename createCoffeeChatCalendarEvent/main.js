function myFunction() {
  //隔週で実行
  var utw = Math.floor(new Date().getTime() / 1000 / 60 / 60 / 24 / 7);
  if (utw % 2 === 0) {
    return;
  }

  const spreadsheet = SpreadsheetApp.openById(
    "1NtXaa50rM5nqlgAgjTbdoaPSR9POFxDI9V5vukNltqM"
  );

  const st = spreadsheet.getSheetByName("希望者");
  const st2 = spreadsheet.getSheetByName("スクリプト用");

  const items1 = getItems(st);
  const items2 = getItems(st2);

  const items = getDiffPair(items1, items2).flat();
  console.log("最終結果の抽選:", items);

  // 2回連続で同じ人にないように前のペアを保存
  for (var i = 0; i < items.length; i++) {
    const item = items[i];
    const row = i + 2;
    st2.getRange(row, 1).setValue(item.no);
    st2.getRange(row, 2).setValue(item.name);
    st2.getRange(row, 3).setValue(item.gmail);
  }

  const today = new Date();
  const y = today.getFullYear();
  const m = today.getMonth();
  const d = today.getDate();
  // 一週間後に設定
  const startTime = new Date(y, m, d, 15, 0, 0);
  startTime.setDate(startTime.getDate() + 7);
  const endTime = new Date(y, m, d, 15, 30, 0);
  endTime.setDate(endTime.getDate() + 7);

  for (let i = 0; i < items.length; i += 2) {
    const target1 = items[i];
    const target2 = items[i + 1];

    const calendar = CalendarApp.getDefaultCalendar();
    calendar.setSelected(true);

    calendar.createEvent("Coffee Chat", startTime, endTime, {
      description: "時間になったらGoogle Meetに入って始めてください。",
      guests: `${target1.gmail},${target2.gmail}`,
    });
  }
}

const getItems = (st) => {
  const maxrow = st.getLastRow() - 1;
  const maxcol = st.getMaxColumns() + 1;
  const data = st.getRange(2, 1, maxrow, maxcol).getValues();
  const header = st.getRange(1, 1, 1, maxcol).getValues()[0];

  const COLUMN_NO_NUMBER = header.indexOf("No");
  const COLUMN_NO_NAME = header.indexOf("名前");
  const COLUMN_NO_GMAIL = header.indexOf("Gmail");

  const items = data
    .filter((v) => {
      if (data.length % 2 == 1) {
        // 奇数の場合は飯野を除外する
        return v[COLUMN_NO_NAME] !== "山田";
      }
      return v;
    })
    .map((v) => {
      return {
        no: v[COLUMN_NO_NUMBER],
        name: v[COLUMN_NO_NAME],
        gmail: v[COLUMN_NO_GMAIL],
      };
    });

  return items;
};

const getPair = (items) => {
  const pair = [];
  for (let i = 0; i < items.length; i += 2) {
    const target1 = items[i];
    const target2 = items[i + 1];
    pair.push([target1, target2]);
  }
  return pair;
};

const comparePair = (pair1, pair2) => {
  const diffPair = pair1.filter((v) => {
    const key = `${v[0].name}_${v[1].name}`;
    // 前回のペアと同じペアは除外
    const same = pair2.find((v2) => {
      return (
        key === `${v2[0].name}_${v2[1].name}` ||
        key === `${v2[1].name}_${v2[0].name}`
      );
    });

    return !same;
  });

  return diffPair;
};

const getDiffPair = (items1, items2) => {
  const pair1 = getPair(items1);
  const pair2 = getPair(items2);
  let items = items1;
  let count = 0;
  let diffPair = [];

  console.log("前回の抽選:", items2);

  while (true) {
    count++;
    // ランダムにソート
    fisherYatesShuffle(items);

    console.log(`${count}度目の抽選:`, items);
    diffPair = [...diffPair, ...comparePair(getPair(items), pair2)];
    console.log(`${count}度目の抽選で異なるペア数:`, diffPair.length);
    if (diffPair.length === pair1.length) {
      // 全て違うペアになったので終了
      return diffPair;
    }
    // 異なるペアになったデータを除外
    items = items1.filter((v) => !diffPair.flat().includes(v));

    if (count > 30) {
      // 10回以上ループしたら強制終了
      const pair = [...diffPair, ...getPair(items)];
      console.log(`${count}度目の抽選なので必ず終了:`, pair);
      return pair;
    }

    if (items.length === 2) {
      console.log(`2人しかいないのでやり直し`);
      // 2人しかいない場合はやり直し
      diffPair = [];
      items = items1;
    }
  }
};

function fisherYatesShuffle(arr) {
  for (var i = arr.length - 1; i > 0; i--) {
    var j = Math.floor(Math.random() * (i + 1)); //random index
    [arr[i], arr[j]] = [arr[j], arr[i]]; // swap
  }
}
