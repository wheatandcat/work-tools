function myFunction() {
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
      var t =
        "- <@" + m.slackUserName + "> / " + events[j].getTitle() + "\n";
      targetText += t;
    }
  }

  var m = now.getMonth() + 1;
  var d = now.getDate();

  var text =
    m +
    "月" +
    d +
    "日の予定\n\n" +
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
