import { parse } from "csv-parse/sync";
import { stringify } from "csv-stringify/sync";
import fs from "node:fs";
import iconv from "iconv-lite";

// コマンドライン引数の取得
const [, , csvFilePath, targetBank] = process.argv;

if (!csvFilePath || !targetBank) {
  console.error("使用方法: node index.js <CSVファイルパス> <保有金融機関>");
  process.exit(1);
}

try {
  // CSVファイルをShift-JISとして読み込み、UTF-8に変換
  const buffer = fs.readFileSync(csvFilePath);
  const fileContent = iconv.decode(buffer, "Shift_JIS");

  const records = parse(fileContent, {
    columns: true,
    skip_empty_lines: true,
  });

  // 指定された保有金融機関のデータのみを抽出
  const filteredRecords = records.filter(
    (record) => record.保有金融機関 === targetBank
  );

  if (filteredRecords.length === 0) {
    console.error("指定された保有金融機関のデータが見つかりませんでした。");
    process.exit(1);
  }

  // 最新の日付を取得
  const latestDate = new Date(
    Math.max(...filteredRecords.map((record) => new Date(record.日付)))
  );
  const month = latestDate.getMonth() + 1; // 月を取得（0-11なので+1）

  // 必要な項目のみを抽出
  const outputRecords = filteredRecords.map((record) => ({
    ID: record.ID,
    日付: record.日付,
    内容: record.内容,
    "金額（円）": record["金額（円）"],
    メモ: record.メモ,
  }));

  // 出力ファイル名の生成
  const outputFileName = `${targetBank}_${month}月.csv`;

  // CSVとして出力（UTF-8で出力）
  const output = stringify(outputRecords, {
    header: true,
    columns: ["ID", "日付", "内容", "金額（円）", "メモ"],
  });

  fs.writeFileSync(outputFileName, output);
  console.log(`ファイルを作成しました: ${outputFileName}`);
} catch (error) {
  console.error("エラーが発生しました:", error.message);
  process.exit(1);
}
