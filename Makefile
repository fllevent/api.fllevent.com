GCC=go
GCMD=run
GPATH= main.go functions.go routes.go db.go readenv.go passwordhash.go

run:
	$(GCC) $(GCMD) $(GPATH)