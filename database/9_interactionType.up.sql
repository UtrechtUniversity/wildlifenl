CREATE TABLE "interactionType" (
	"ID" SERIAL NOT NULL,
	"nameNL" VARCHAR NOT NULL,
	"nameEN" VARCHAR NOT NULL,
	"descriptionNL" TEXT NOT NULL DEFAULT '',
	"descriptionEN" TEXT NOT NULL DEFAULT '',
	PRIMARY KEY ("ID")
);

INSERT INTO "interactionType" ("ID", "nameNL", "nameEN", "descriptionNL", "descriptionEN") VALUES
	(1, 'Waarneming', 'Sighting', 'Een levend wild dier is gezien.', 'A living wild animal was seen.'),
	(2, 'Schademelding', 'Damage Report', 'Een melding van schade gemaakt door een wild dier.', 'A report of damage inflicted by a wild animal.'),
	(3, 'Dieraanrijding', 'Animal-vehicle collision', 'U wilt een aanrijding van een wild dier melden.', 'You would like to report a wild animal-vehicle collision.');