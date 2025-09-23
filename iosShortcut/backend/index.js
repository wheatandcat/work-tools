import functions from "@google-cloud/functions-framework";

/**
 * iOS Shortcut用のHTTP関数
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 */
functions.http("processRequest", (req, res) => {
  // CORS設定
  res.set("Access-Control-Allow-Origin", "*");
  res.set("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
  res.set("Access-Control-Allow-Headers", "Content-Type");

  // OPTIONSリクエスト（プリフライトリクエスト）の処理
  if (req.method === "OPTIONS") {
    res.status(204).send("");
    return;
  }

  try {
    // リクエストパラメータの取得
    let url;

    if (req.method === "GET") {
      url = req.query.url;
    } else if (req.method === "POST") {
      url = req.body.url;
    } else {
      res.status(405).json({ error: "Method not allowed" });
      return;
    }

    // リクエストパラメータのログ出力
    console.log("Request parameters received:");
    console.log(`URL: ${url}`);
    console.log(`Method: ${req.method}`);
    console.log(`User-Agent: ${req.headers["user-agent"]}`);
    console.log(`Timestamp: ${new Date().toISOString()}`);

    // パラメータの検証
    if (!url) {
      console.log("Error: Missing required parameters");
      res.status(400).json({
        error: "Missing required parameters",
        required: ["url"],
      });
      return;
    }

    // 成功レスポンス
    console.log("Request processed successfully");
    res.status(200).json({
      success: true,
      message: "Request processed successfully",
      receivedData: {
        url: url,
      },
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    console.error("Error processing request:", error);
    res.status(500).json({
      error: "Internal server error",
      message: error.message,
    });
  }
});
