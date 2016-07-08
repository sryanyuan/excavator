@pushd %~dp0
@echo visit the web (http://localhost:1111)
@start excavator -lisaddr localhost:1111

@ping -n 3 127.1>nul
@explorer http://localhost:1111
@popd