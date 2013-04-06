TEXINPUTS="..:."
pdflatex --shell-escape report.tex
pdflatex --shell-escape report.tex
rm report-blx.bib
rm *.bbl
rm *.blg
rm *.fdb_latexmk
rm *.fls
rm *.out
rm *.pyg
rm *.run.xml
rm *.synctex.gz
rm *.fls
rm *log
rm *.aux