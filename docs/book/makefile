all: main.pdf

%.pdf: %.tex
	pdflatex $^

clean:
	$(RM) *.aux *.log *.synctex.gz *.fls *.out *.fdb_latexmk *.acn *.bcf *.ist *.glo *.nlo *.toc *.run.xml *.lof *.lot
