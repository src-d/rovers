package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/pipe.v2"
)

var companies = map[string]int{
	// "aim":   0, // 1aim
	// "dixie": 0, // Dixie

	"arcmedia":    {2068129}, // Arcmedia
	"auctionata":  {1091184}, // Auctionata
	"booking":     {11348},   // Booking
	"cartodb":     {5084329},
	"cliqz":       {9245817}, // Cliqz
	"dashlane":    {2049626}, // Dashlane
	"dice":        {5150723}, // Dice
	"digitgaming": {2609851}, // Digitgaming
	"founders": {
		3199273, // Founders
		5090485, // GoBox
		3085865, // Son of a Tailor
		3365255, // MinbilDinbil
		5228312, // Maguru
		5318327, // Pipetop
	},
	"freespee":           {312064},  // Freespee
	"gelato":             {5037871}, // Gelato
	"locafox":            {3341251}, // Locafox
	"lovoo":              {3234901}, // Lovoo
	"memorado":           {5043633}, // Memorado
	"movidiam":           {5191432}, // Movidiam
	"packlink":           {2311968}, // Packlink
	"property-guru":      {623046},  // Property Guru
	"rebelminds":         {9373692}, // Rebelminds
	"skypickercom":       {3010943}, // Skypicker.com
	"springlane":         {5006229}, // Springlane
	"squirro":            {2502648}, // Squirro
	"stratified-medical": {3609919}, // Stratified Medical
	"take-eat-easy":      {2687696}, // Take Eat Easy
	"transferwise":       {1769571}, // Transferwise
	"travelperk":         {9310632}, // Travelperk
	"twitter-counter":    {595481},  // Twitter Counter
	"tyba":               {924688},  // Tyba
	"uberchord":          {5136213}, // Uberchord
	"unumotors":          {3206274}, // UnuMotors
	"vivino":             {1148120}, // Vivino
	"zorgvoorelkaar":     {2512372}, // ZorgVoorElkaar
	"zve":                {2512372}, // Zve

	"samyroad":     {2602076},
	"spotahome":    {5182445},
	"typeform":     {3226972},
	"jobandtalent": {296493},
	"marfeel":      {2406943},
	"mytaxi":       {1862054},
	"traity":       {2512033},
}

func main() {
	var filename string
	flag.StringVar(&filename, "cookiefile", "utils/updateemployees/cookie.save", "path to cookie file")
	flag.Parse()

	cookie := readCookieFile(filename)

	start := time.Now()

	log15.Info("Building binary...")
	binaryStart := time.Now()
	pipe.Run(pipe.Exec("go", "build", "-o", "srcd-rovers", "rovers.go"))
	log15.Info("Done", "elapsed", time.Since(binaryStart))

	for codename, ids := range companies {
		for _, id := range ids {
			scriptStart := time.Now()
			script := pipe.Script(
				pipe.Exec("./srcd-rovers", "linkedin",
					"--companyCodename", codename,
					"--companyId", strconv.Itoa(id),
					"--cookie", cookie,
				),
			)
			RunScript(script)
			log15.Info("Done",
				"company", codename,
				"elapsed", time.Since(scriptStart),
			)
		}
	}

	log15.Info("Done", "elapsed", time.Since(start))
}

func RunScript(script pipe.Pipe) {
	// Stream output to real stdout and stderr
	state := pipe.NewState(os.Stdout, os.Stderr)
	err := script(state)
	if err != nil {
		log15.Error("Can't run script", "error", err)
		os.Exit(1)
	}
	err = state.RunTasks()
	if err != nil {
		log15.Error("Error running script", "error", err)
		os.Exit(1)
	}
}

func readCookieFile(filename string) string {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		log.Panicf("A %q file is required", filename)
	}
	if err != nil {
		log.Panicf("Couldn't read %q file. Error: %s", filename, err)
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panicf("Error reading %q file. Error: %s", filename, err)
	}
	return string(body)
}
