package lang

import (
	"time"

	"github.com/xeonx/timeago"
)

var english = timeago.Config{
	PastPrefix:   "",
	PastSuffix:   " ago",
	FuturePrefix: "in ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "about a second", Many: "%d seconds"},
		{D: time.Minute, One: "about a minute", Many: "%d minutes"},
		{D: time.Hour, One: "about an hour", Many: "%d hours"},
		{D: timeago.Day, One: "one day", Many: "%d days"},
		{D: timeago.Month, One: "one month", Many: "%d months"},
		{D: timeago.Year, One: "one year", Many: "%d years"},
	},

	Zero: "about a second",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "2006-01-02",
}

var portuguese = timeago.Config{
	PastPrefix:   "há ",
	PastSuffix:   "",
	FuturePrefix: "daqui a ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "um segundo", Many: "%d segundos"},
		{D: time.Minute, One: "um minuto", Many: "%d minutos"},
		{D: time.Hour, One: "uma hora", Many: "%d horas"},
		{D: timeago.Day, One: "um dia", Many: "%d dias"},
		{D: timeago.Month, One: "um mês", Many: "%d meses"},
		{D: timeago.Year, One: "um ano", Many: "%d anos"},
	},

	Zero: "menos de um segundo",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02-01-2006",
}

var spanish = timeago.Config{
	PastPrefix:   "hace ",
	PastSuffix:   "",
	FuturePrefix: "en ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "un segundo", Many: "%d segundos"},
		{D: time.Minute, One: "un minuto", Many: "%d minutos"},
		{D: time.Hour, One: "una hora", Many: "%d horas"},
		{D: timeago.Day, One: "un día", Many: "%d días"},
		{D: timeago.Month, One: "un mes", Many: "%d meses"},
		{D: timeago.Year, One: "un año", Many: "%d años"},
	},

	Zero:          "menos de un segundo",
	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02-01-2006",
}

var chinese = timeago.Config{
	PastPrefix:   "",
	PastSuffix:   "前",
	FuturePrefix: "于 ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "1 秒", Many: "%d 秒"},
		{D: time.Minute, One: "1 分钟", Many: "%d 分钟"},
		{D: time.Hour, One: "1 小时", Many: "%d 小时"},
		{D: timeago.Day, One: "1 天", Many: "%d 天"},
		{D: timeago.Month, One: "1 月", Many: "%d 月"},
		{D: timeago.Year, One: "1 年", Many: "%d 年"},
	},

	Zero: "1 秒",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "2006-01-02",
}

var french = timeago.Config{
	PastPrefix:   "il y a ",
	PastSuffix:   "",
	FuturePrefix: "dans ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "environ une seconde", Many: "moins d'une minute"},
		{D: time.Minute, One: "environ une minute", Many: "%d minutes"},
		{D: time.Hour, One: "environ une heure", Many: "%d heures"},
		{D: timeago.Day, One: "un jour", Many: "%d jours"},
		{D: timeago.Month, One: "un mois", Many: "%d mois"},
		{D: timeago.Year, One: "un an", Many: "%d ans"},
	},

	Zero: "environ une seconde",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var german = timeago.Config{
	PastPrefix:   "vor ",
	PastSuffix:   "",
	FuturePrefix: "in ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "einer Sekunde", Many: "%d Sekunden"},
		{D: time.Minute, One: "einer Minute", Many: "%d Minuten"},
		{D: time.Hour, One: "einer Stunde", Many: "%d Stunden"},
		{D: timeago.Day, One: "einem Tag", Many: "%d Tagen"},
		{D: timeago.Month, One: "einem Monat", Many: "%d Monaten"},
		{D: timeago.Year, One: "einem Jahr", Many: "%d Jahren"},
	},

	Zero: "einer Sekunde",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02.01.2006",
}

var turkish = timeago.Config{
	PastPrefix:   "",
	PastSuffix:   " önce",
	FuturePrefix: "",
	FutureSuffix: " içinde",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "yaklaşık bir saniye", Many: "%d saniye"},
		{D: time.Minute, One: "yaklaşık bir dakika", Many: "%d dakika"},
		{D: time.Hour, One: "yaklaşık bir saat", Many: "%d saat"},
		{D: timeago.Day, One: "bir gün", Many: "%d gün"},
		{D: timeago.Month, One: "bir ay", Many: "%d ay"},
		{D: timeago.Year, One: "bir yıl", Many: "%d yıl"},
	},

	Zero: "yaklaşık bir saniye",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var korean = timeago.Config{
	PastPrefix:   "",
	PastSuffix:   " 전",
	FuturePrefix: "",
	FutureSuffix: " 후",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "약 1초", Many: "%d초"},
		{D: time.Minute, One: "약 1분", Many: "%d분"},
		{D: time.Hour, One: "약 한시간", Many: "%d시간"},
		{D: timeago.Day, One: "하루", Many: "%d일"},
		{D: timeago.Month, One: "1개월", Many: "%d개월"},
		{D: timeago.Year, One: "1년", Many: "%d년"},
	},

	Zero: "방금",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "2006-01-02",
}

var italian = timeago.Config{
	PastPrefix:   "fa ",
	PastSuffix:   "",
	FuturePrefix: "tra ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "circa un secondo", Many: "meno di un minuto"},
		{D: time.Minute, One: "circa un minuto", Many: "%d minuti"},
		{D: time.Hour, One: "circa un'ora", Many: "%d ore"},
		{D: timeago.Day, One: "un giorno", Many: "%d giorni"},
		{D: timeago.Month, One: "un mese", Many: "%d mesi"},
		{D: timeago.Year, One: "un anno", Many: "%d anni"},
	},

	Zero: "circa un secondo",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var russian = timeago.Config{
	PastPrefix:   "назад ",
	PastSuffix:   "",
	FuturePrefix: "через ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "около одной секунды", Many: "менее минуты"},
		{D: time.Minute, One: "около одной минуты", Many: "%d минут"},
		{D: time.Hour, One: "около одного часа", Many: "%d часов"},
		{D: timeago.Day, One: "один день", Many: "%d дней"},
		{D: timeago.Month, One: "один месяц", Many: "%d месяцев"},
		{D: timeago.Year, One: "один год", Many: "%d лет"},
	},

	Zero: "около одной секунды",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var romanian = timeago.Config{
	PastPrefix:   "acum ",
	PastSuffix:   "",
	FuturePrefix: "peste ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "aproximativ un secundă", Many: "mai puțin de un minut"},
		{D: time.Minute, One: "aproximativ un minut", Many: "%d minute"},
		{D: time.Hour, One: "aproximativ o oră", Many: "%d ore"},
		{D: timeago.Day, One: "o zi", Many: "%d zile"},
		{D: timeago.Month, One: "o lună", Many: "%d luni"},
		{D: timeago.Year, One: "un an", Many: "%d ani"},
	},

	Zero: "aproximativ un secundă",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02/01/2006",
}

var dutch = timeago.Config{
	PastPrefix:   "ongeveer ",
	PastSuffix:   " geleden",
	FuturePrefix: "over ",
	FutureSuffix: "",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "ongeveer een seconde", Many: "minder dan een minuut"},
		{D: time.Minute, One: "ongeveer een minuut", Many: "%d minuten"},
		{D: time.Hour, One: "ongeveer een uur", Many: "%d uur"},
		{D: timeago.Day, One: "een dag", Many: "%d dagen"},
		{D: timeago.Month, One: "een maand", Many: "%d maanden"},
		{D: timeago.Year, One: "een jaar", Many: "%d jaar"},
	},

	Zero: "ongeveer een seconde",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "02-01-2006",
}

var japanese = timeago.Config{
	PastPrefix:   "約",
	PastSuffix:   "前",
	FuturePrefix: "約",
	FutureSuffix: "後",

	Periods: []timeago.FormatPeriod{
		{D: time.Second, One: "約1秒", Many: "1分未満"},
		{D: time.Minute, One: "約1分", Many: "%d分"},
		{D: time.Hour, One: "約1時間", Many: "%d時間"},
		{D: timeago.Day, One: "1日", Many: "%d日"},
		{D: timeago.Month, One: "1ヶ月", Many: "%dヶ月"},
		{D: timeago.Year, One: "1年", Many: "%d年"},
	},

	Zero: "約1秒",

	Max:           100 * 365 * 24 * time.Hour,
	DefaultLayout: "2006/01/02",
}

func TimeAgo(date time.Time) string {
	switch Current {
	case "pt":
		return portuguese.Format(date)
	case "es":
		return spanish.Format(date)
	case "fr":
		return french.Format(date)
	case "it":
		return italian.Format(date)
	case "ru":
		return russian.Format(date)
	case "ro":
		return romanian.Format(date)
	case "nl":
		return dutch.Format(date)
	case "jp":
		return japanese.Format(date)
	case "ko":
		return korean.Format(date)
	case "zh":
		return chinese.Format(date)
	}

	return english.Format(date)
}
