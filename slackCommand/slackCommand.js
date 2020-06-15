function doPost(e) {
  var token = e.parameter.token;
  var verificationToken = PropertiesService.getScriptProperties().getProperty(
    "verificationToken"
  );

  if (token !== verificationToken) {
    throw new Error("Invalid token");
  }

  var member = [
    {
      email: "wheatandcat@gmail.com",
      slackUserName: "wheatandcat",
      name: "wheatandcat",
    },
  ];

  var now = new Date();
  var targetText = "";

  for (var i = 0; i < member.length; i++) {
    var m = member[i];
    var c = CalendarApp.getCalendarById(m.email);
    var events = c.getEventsForDay(now);
    for (var j = 0; j < events.length; j++) {
      var t = "- <@" + m.slackUserName + "> / " + events[j].getTitle() + "\n";
      targetText += t;
    }
  }

  var m = now.getMonth() + 1;
  var d = now.getDate();

  var text = m + "月" + d + "日の予定\n\n" + targetText;

  return ContentService.createTextOutput(text);
}

