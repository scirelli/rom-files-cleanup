package model

//softwarelist.dtd
//go:generate goxmlstruct -named-types -use-pointers-for-optional-fields=0 -output softwarelist.go -package-name model -pattern ../../../assets/databases/mame/hash/*.xml
//xsi:schemaLocation="https://datomatic.no-intro.org/stuff https://datomatic.no-intro.org/stuff/schema_nointro_datfile_v3.xsd"
//go:generate goxmlstruct -named-types -use-pointers-for-optional-fields=0 -output dataFile.go -package-name model -pattern ../../../assets/databases/dats/**/*.dat
//Unknown dtd
//go:generate goxmlstruct -named-types -use-pointers-for-optional-fields=0 -output hashFile.go -package-name model -pattern ../../../assets/databases/xmame/hash/*.hsi
