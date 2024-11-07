package minecraft

import (
	"encoding/json"
	"sync/atomic"

	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/version"
)

type motdObject struct {
	Version struct {
		Name     string `json:"name"`
		Protocol uint   `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int32 `json:"max"`
		Online int32 `json:"online"`
		Sample any   `json:"sample,omitempty"`
	} `json:"players"`
	Description struct {
		Text string `json:"text,omitempty"`
	} `json:"description,omitempty"`
	Favicon string `json:"favicon,omitempty"`
}

type playerSample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func generateMOTD(protocolVersion uint, s *config.Outbound, onlineCount *atomic.Int32) []byte {
	online := s.Minecraft.OnlineCount.Online
	if online < 0 {
		online = onlineCount.Load()
	}

	motd, _ := json.Marshal(motdObject{
		Version: struct {
			Name     string `json:"name"`
			Protocol uint   `json:"protocol"`
		}{
			Name:     "zbproxy " + version.Version,
			Protocol: protocolVersion,
		},
		Players: struct {
			Max    int32 `json:"max"`
			Online int32 `json:"online"`
			Sample any   `json:"sample,omitempty"`
		}{
			Max:    s.Minecraft.OnlineCount.Max,
			Online: online,
			Sample: s.Minecraft.OnlineCount.Sample,
		},
		Description: struct {
			Text string `json:"text,omitempty"`
		}{
			Text: s.Minecraft.MotdDescription,
		},
		Favicon: s.Minecraft.MotdFavicon,
	})

	return motd
}

const defaultMOTD = `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAMAAACdt4HsAAAAAXNSR0IB2cksfwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAeZQTFRF/+IAzrcHKyodHx8fR0IZ89cC2MAFKikdw64IIyMeXFQX/eAAxa8IIiIfrpwLISAfgHQRnYwOuKQKaF4VJyYecGUUj4AQ99sBUUsYMS4d38YE+98AWVEXrJkM7dMCQDsaPzsb9NgCTkgYuqYJ3sYEUEoYODQcJiUe0LkHW1MX+94B178GMC4dcmgTwq0I5s0DpJINoZAN69shjLi/i7i/ZZfEU4fHVIfHeqnB0NFOZqr/f7PWY1sVtcd7n7+g+NwBm76nv8tq7tMCRkEagLPU4Nc06M4DPjob4McE49cvgLTUxMxjnL6lpMCYi6B7s8d/ssZ/a4iJQmWPZaj7NUdcXnNvYJ7sMkVcu8pxW5XdLjtMNjMcssaBVovNKTI9qMORUYG+JCgtnr6hTHiu+N8MlbuxR26f7twci7fBQ2SPfXES5dgsgbTRPVuA29U8ebLgxMxdiHoQ0tFLb63w0dFMOFFverHe3NU6PVl9hLTN5tkpQWKLj7m78d0XRWqZmb2q++EGSnOnICIkTny0rsSHqpgMJSoxUoTCuch1KTNAV43QZHlxSXKkcY6IdrDl+N4EZqn8lLuy2tQ+t8h40tJMmLyr8N0ZtqIKerHdd7DjloYPkbq3dGkT59ooq8OMxM1ixs1gzbYH1nu7OAAAA1JJREFUeJyFl/dfE0EQxZesRkPU0CIlGCACIiAWrBhFQQELotjF3sGK2HvD3huK2P5TL7nL7bwZPnvvR/bNl83eu5k9pajyQpppylQwhKdxg0JNj3BDfhQMM2baAWpWjBMK7AYOUIUcoAutBgFQBRwQK7IZJEAVc0JJHqxH4wCYLQClZZxQXgGGRCUF6DmCkKzihOokGGpSFBCbKwi15ZxQh3FI1hOAntcgCHklnDC/EQxNIQLQzWFBKBJxWICGlggB6PqFgiDjwE57EQXwxGck4qAXT2pwATzxS1pbly5bzrRiZasn1xSnAJb4VW02rXZNiTQFYOLXWAFrPVd7igIg8eusgPU5W0cnAUDiN1gBG31fVzcB0MRb69vIVnuaCcAkfpO1fjM9rQYK8BO/xQrYinmArHiJ77UCtlkAXuL7iLbv6O+H+p27WiwAnnhHid1Qv2evDtVaADzxSu2D+v0HHM9AjwXA+8tBqD90OGvq7jCGOk7A/nIE6o8e80ypGt+RrOYE2l+OQ/0JY6pM+J4K0QAHurylk6eg/jQ1pU0DkAOz0/2JZ85C/SCahsw+xcDUqXbnz+fOQ/0FbiJPXE7U4qi6eAnqL6e5hz5x2UKHr4zQ8pGrbvsARUgkRQu9Bv/++g3ltg9UqMkH4MDU+ibU37qdNTntg6nejDwYmPoO1N+955l6mjmhrNQn0IF5H+ofPPRNDeKJV5pjMAPzEdT30mTLJ04WvYH5+AnUP1WgUQvAHZjPnkP9C8U0bAFkBubLV1D/OsoB0WILwInDG6gf7H/bx8RvQIwfxw1IveOBYoD3HwIAHx1T1wAF4I/8FFDf1pdxhUmg1DAAPgcBvmRtJFBKj1LA1yDAN9c3FjMA+lKq7wH1P3LGcQOAORG0gZ++c8IAyJz4FQT47QMa8w3AzIk/QYC/ZrPe6+tOJjInXPHEaz2uuDq6DYDOiax44h2NCUK2v+RmY5q9NfLOX/JPEDJx8IfrEFuULbS8VhCc/mKmM78ZQOKzqpIX4ln0lsZvBmHRQrMDCzVOAJEWtihbKD+pjMgymROuxsTI4yeFADonvA1ywCR3KFgmc8LVhCCIOxQu80DlEm8kvtHYepwdU42Yyfwbja+zb20v8VT4jfYfTXskk4+wbR0AAAAASUVORK5CYII=`
