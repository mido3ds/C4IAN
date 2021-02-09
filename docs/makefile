MD_FILES = $(shell find -type f -name '*.md' | sed 's/\.md/\.pdf/g' | xargs -I {} echo pdf/{})

all: $(MD_FILES)

pdf/%.pdf: %.md
	env -C $(shell dirname $^) pandoc $(shell basename -s .pdf $@).md -o pdf/$(shell basename $@) --number-sections -V linkcolor=blue
