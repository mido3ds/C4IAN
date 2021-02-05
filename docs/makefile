MD_FILES = $(shell find -type f -name '*.md' | sed 's/\.md/\.pdf/g')

all: $(MD_FILES)

%.pdf: %.md
	env -C $(shell dirname $@) pandoc $(shell basename -s .pdf $@).md -o pdf/$(shell basename $@)
